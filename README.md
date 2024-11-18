# Dump all Kubernetes resources into a directory structure

Dumps all Kubernetes resources into a directory structure:

Attention: This dumps secrets, too!

```text
out/NAMESPACE/GVK/NAME.yaml
```

For example:

```text
out/kube-system/v1.ConfigMap/kubelet-config.yaml
```
For non-namespaces resources, the subdirectory is called `_cluster`.
