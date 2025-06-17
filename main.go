package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

const (
	clusterNamespace = "_cluster"
)

var toSkip = map[string][]string{
	"apps/v1": {"replicasets"},
	"authentication.k8s.io/v1": {
		"selfsubjectreviews", "tokenreviews",
	},
	"authorization.k8s.io/v1": {
		"selfsubjectaccessreviews", "subjectaccessreviews",
		"selfsubjectrulesreviews", "localsubjectaccessreviews",
	},
	"coordination.k8s.io/v1": {"leases"},
	"discovery.k8s.io/v1":    {"endpointslices"},
	"v1":                     {"bindings", "componentstatuses"},
}

type options struct {
	outputDir         string
	quiet             bool
	dumpSecrets       bool
	dumpManagedFields bool
	removeOutdir      bool
}

func main() {
	err := mainWithError()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mainWithError() error {
	opts := &options{}
	pflag.StringVarP(&opts.outputDir, "out-dir", "o", "out", "Output directory (must not exist)")
	pflag.BoolVarP(&opts.quiet, "quiet", "q", false, "Quiet, suppress output")
	pflag.BoolVarP(&opts.dumpSecrets, "dump-secrets", "s", true, "Dump secrets")
	pflag.BoolVarP(&opts.dumpManagedFields, "dump-managed-fields", "m", false, "Dump managed fields (disabled by default)")
	pflag.BoolVarP(&opts.removeOutdir, "remove-out-dir", "r", false, "Remove out-dir before dumping (disabled by default)")
	pflag.Parse()

	if opts.removeOutdir {
		if err := os.RemoveAll(opts.outputDir); err != nil {
			return fmt.Errorf("failed to remove out-dir %s: %w", opts.outputDir, err)
		}
	}

	if _, err := os.Stat(opts.outputDir); err == nil {
		return fmt.Errorf("output directory %q already exists", opts.outputDir)
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// 80 concurrent requests were served in roughly 200ms
	// This means 400 requests in one second (to local kind cluster)
	// But why reduce this? I don't want people with better hardware
	// to wait for getting results from an api-server running at localhost
	config.QPS = 1000
	config.Burst = 1000

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w\n", err)
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create discovery client: %w\n", err)
	}

	resourceList, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return fmt.Errorf("failed to discover resources: %w\n", err)
	}

	var globalFileCount int64 = 0
	sort.Slice(resourceList, func(i int, j int) bool {
		return resourceList[i].GroupVersion < resourceList[j].GroupVersion
	})
	for _, apiGroup := range resourceList {
		sort.Slice(apiGroup.APIResources, func(i int, j int) bool {
			return apiGroup.APIResources[i].Name < apiGroup.APIResources[j].Name
		})
		for _, resource := range apiGroup.APIResources {
			skipSlice := toSkip[apiGroup.GroupVersion]
			if slices.Contains(skipSlice, resource.Name) {
				continue
			}

			gvr := schema.GroupVersionResource{
				Group:    getGroup(apiGroup.GroupVersion),
				Version:  getVersion(apiGroup.GroupVersion),
				Resource: resource.Name,
			}

			fileCount, err := processResource(dynClient, gvr, resource.Namespaced, opts)
			if err != nil {
				fmt.Printf("Failed to process resource %q %s: %v\n", apiGroup.GroupVersion, resource.Name, err)
			}
			globalFileCount += fileCount
		}
	}
	fmt.Printf("Total files written: %d\n", globalFileCount)
	return nil
}

func processResource(client dynamic.Interface, gvr schema.GroupVersionResource, isNamespaced bool, options *options) (count int64, err error) {
	var fileCount int64 = 0
	list, err := client.Resource(gvr).List(context.TODO(), meta.ListOptions{})
	if err != nil {
		return fileCount, fmt.Errorf("failed to list resources for %s: %v", gvr.Resource, err)
	}
	for _, item := range list.Items {
		ns := item.GetNamespace()
		if !isNamespaced {
			ns = clusterNamespace
		}
		gvk := item.GroupVersionKind()
		name := item.GetName()

		var dirPath string
		if gvk.Group == "" {
			dirPath = filepath.Join(options.outputDir, ns, gvk.Kind)
		} else {
			dirPath = filepath.Join(options.outputDir, ns, fmt.Sprintf("%s_%s", gvk.Group, gvk.Kind))
		}
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return fileCount, fmt.Errorf("failed to create directory %s: %v", dirPath, err)
		}
		yamlFilePath := filepath.Join(dirPath, fmt.Sprintf("%s.yaml", strings.ReplaceAll(name, ":", "_")))
		err := writeYAML(yamlFilePath, item.Object, options)
		if err != nil {
			fmt.Printf("Failed to write YAML for %s: %v\n", yamlFilePath, err)
		}

		// Dump logs if the resource is a Pod
		if gvk.Kind == "Pod" {
			err := dumpPodLogs(ns, name, dirPath, options)
			if err != nil {
				fmt.Printf("Failed to dump logs for Pod %s/%s: %v\n", ns, name, err)
			}
		}
		fileCount++
	}
	return int64(fileCount), nil
}

func dumpPodLogs(namespace, podName, dirPath string, options *options) error {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return fmt.Errorf("failed to build Kubernetes config: no in-cluster or kubeconfig found. Error: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes clientset: %v", err)
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, meta.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get Pod %s/%s: %v", namespace, podName, err)
	}

	for _, container := range pod.Spec.Containers {
		logOptions := &corev1.PodLogOptions{
			Container: container.Name,
		}
		logRequest := clientset.CoreV1().Pods(namespace).GetLogs(podName, logOptions)

		logs, err := logRequest.Stream(context.TODO())
		if err != nil {
			fmt.Printf("Failed to get logs for container %s in Pod %s/%s: %v\n", container.Name, namespace, podName, err)
			continue
		}
		defer logs.Close()

		logFilePath := filepath.Join(dirPath, fmt.Sprintf("%s_%s_logs.txt", podName, container.Name))
		logFile, err := os.Create(logFilePath)
		if err != nil {
			return fmt.Errorf("failed to create log file %s: %v", logFilePath, err)
		}
		defer logFile.Close()

		_, err = io.Copy(logFile, logs)
		if err != nil {
			return fmt.Errorf("failed to write logs to file %s: %v", logFilePath, err)
		}

		if !options.quiet {
			fmt.Printf("Written logs: %s\n", logFilePath)
		}
	}
	return nil
}

func writeYAML(filePath string, obj map[string]interface{}, options *options) error {
	metadata, ok := obj["metadata"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("metadata not found in object")
	}
	if !options.dumpManagedFields {
		delete(metadata, "managedFields")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()
	encoder.SetIndent(2)
	if err := encoder.Encode(obj); err != nil {
		return fmt.Errorf("failed to write YAML to file %s: %v", filePath, err)
	}
	if !options.quiet {
		fmt.Printf("Written: %s\n", filePath)
	}
	return nil
}

func getGroup(groupVersion string) string {
	parts := strings.Split(groupVersion, "/")
	if len(parts) > 1 {
		return parts[0]
	}
	return ""
}

func getVersion(groupVersion string) string {
	parts := strings.Split(groupVersion, "/")
	return parts[len(parts)-1]
}
