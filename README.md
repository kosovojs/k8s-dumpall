# Dump all Kubernetes resources into a directory structure

Dumps all Kubernetes resources into a directory structure:

```text
out/NAMESPACE/GVK/NAME.yaml
```

For example:

```text
out/kube-system/v1.ConfigMap/kubelet-config.yaml
```

Attention: This dumps secrets, too!
