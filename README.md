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

Attention: This dumps secrets, too!

## Via `go run`

The easiest way is to run the code like this:

```terminal
go run github.com/guettli/dumpall@latest

Written: out/cert-manager/v1.Service/cert-manager.yaml
Written: out/cert-manager/v1.Service/cert-manager-webhook.yaml
Written: out/default/v1.Service/kubernetes.yaml
Written: out/_cluster/v1.Namespace/cert-manager.yaml    <-- non-namespaces resources use the directory "_cluster"
...
```

## See Changes

After running dumpall you can modify your cluster, or just wait some time.

Then you can compare the changes with your favorite diff tool. I like [Meld](https://meldmerge.org/):

```terminal
mv out out-1

go run github.com/guettli/dumpall@latest

meld out-1 out
```

## Feedback is welcome

Please create an issue if you have a question or a feature request.
