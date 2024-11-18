package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	outputDir        = "out"
	clusterNamespace = "_cluster"
)

func main() {
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

	skippedResources := sets.NewString("leases", "events", "endpointslices", "selfsubjectreviews",
		"tokenreviews", "localsubjectaccessreviews", "selfsubjectrulesreviews",
		"subjectaccessreviews", "selfsubjectaccessreviews", "bindings", "componentstatuses",
	)

	for _, apiGroup := range resourceList {
		for _, resource := range apiGroup.APIResources {
			if skippedResources.Has(strings.ToLower(resource.Name)) {
				continue
			}

			gvr := schema.GroupVersionResource{
				Group:    getGroup(apiGroup.GroupVersion),
				Version:  getVersion(apiGroup.GroupVersion),
				Resource: resource.Name,
			}

			err = processResource(dynClient, gvr, resource.Namespaced)
			if err != nil {
				fmt.Printf("Failed to process resource %s: %v\n", resource.Name, err)
			}
		}
	}
}

func processResource(client dynamic.Interface, gvr schema.GroupVersionResource, isNamespaced bool) error {
	list, err := client.Resource(gvr).List(context.TODO(), meta.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list resources for %s: %v", gvr.Resource, err)
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
			dirPath = filepath.Join(outputDir, ns, fmt.Sprintf("%s.%s", gvk.Version, gvk.Kind))
		} else {
			dirPath = filepath.Join(outputDir, ns, fmt.Sprintf("%s.%s.%s", gvk.Group, gvk.Version, gvk.Kind))
		}
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
		}
		filePath := filepath.Join(dirPath, fmt.Sprintf("%s.yaml", name))
		err := writeYAML(filePath, item.Object)
		if err != nil {
			fmt.Printf("Failed to write YAML for %s: %v\n", filePath, err)
		}
	}

	return nil
}

func writeYAML(filePath string, obj map[string]interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(obj); err != nil {
		return fmt.Errorf("failed to write YAML to file %s: %v", filePath, err)
	}

	fmt.Printf("Written: %s\n", filePath)
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
