package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/spf13/pflag"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
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
	"events.k8s.io/v1":       {"events"},
	"v1":                     {"events", "bindings", "componentstatuses"},
}

type options struct {
	outputDir         string
	quiet             bool
	dumpSecrets       bool
	dumpManagedFields bool
}

func main() {
	opts := &options{}
	pflag.StringVarP(&opts.outputDir, "output-dir", "o", "out", "Output directory")
	pflag.BoolVarP(&opts.quiet, "quiet", "q", false, "Quiet, suppress output")
	pflag.BoolVarP(&opts.dumpSecrets, "dump-secrets", "s", false, "Dump secrets (disabled by default)")
	pflag.BoolVarP(&opts.dumpManagedFields, "dump-managed-fields", "m", false, "Dump managed fields (disabled by default)")
	pflag.Parse()
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
		fmt.Printf("Failed to create dynamic client: %v\n", err)
		return
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		fmt.Printf("Failed to create discovery client: %v\n", err)
		return
	}

	resourceList, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		fmt.Printf("Failed to discover resources: %v\n", err)
		return
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
		filePath := filepath.Join(dirPath, fmt.Sprintf("%s.yaml", name))
		err := writeYAML(filePath, item.Object, options)
		if err != nil {
			fmt.Printf("Failed to write YAML for %s: %v\n", filePath, err)
		}
		fileCount++
	}
	return int64(fileCount), nil
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
