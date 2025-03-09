# Clusterctl Manual

Get started with Cluster API using clusterctl to create a management cluster,
install providers, and create templates for your workload cluster.

```text
clusterctl [command] [global flags] [command flags]
```

### Global Flags

```text
      --config $XDG_CONFIG_HOME/cluster-api/clusterctl.yaml   Path to clusterctl configuration (default is $XDG_CONFIG_HOME/cluster-api/clusterctl.yaml) or to a remote location (i.e. https://example.com/clusterctl.yaml)
  -v, --v int                                                 Set the log level verbosity. This overrides the CLUSTERCTL_LOG_LEVEL environment variable.
```

### Cluster Management Commands:

* [clusterctl delete](#clusterctl-delete)
* [clusterctl generate](#clusterctl-generate)
* [clusterctl get](#clusterctl-get)
* [clusterctl init](#clusterctl-init)
* [clusterctl move](#clusterctl-move)
* [clusterctl upgrade](#clusterctl-upgrade)

### Troubleshooting and Debugging Commands:

* [clusterctl describe](#clusterctl-describe)

### Other Commands:

* [clusterctl alpha](#clusterctl-alpha)
* [clusterctl completion](#clusterctl-completion)
* [clusterctl config](#clusterctl-config)
* [clusterctl help](#clusterctl-help)
* [clusterctl version](#clusterctl-version)

### Additional Commands

* [clusterctl alpha help](#clusterctl-alpha-help)
* [clusterctl alpha rollout](#clusterctl-alpha-rollout)
* [clusterctl alpha rollout help](#clusterctl-alpha-rollout-help)
* [clusterctl alpha rollout pause](#clusterctl-alpha-rollout-pause)
* [clusterctl alpha rollout restart](#clusterctl-alpha-rollout-restart)
* [clusterctl alpha rollout resume](#clusterctl-alpha-rollout-resume)
* [clusterctl alpha topology](#clusterctl-alpha-topology)
* [clusterctl alpha topology help](#clusterctl-alpha-topology-help)
* [clusterctl config help](#clusterctl-config-help)
* [clusterctl config repositories](#clusterctl-config-repositories)
* [clusterctl describe cluster](#clusterctl-describe-cluster)
* [clusterctl describe help](#clusterctl-describe-help)
* [clusterctl generate cluster](#clusterctl-generate-cluster)
* [clusterctl generate help](#clusterctl-generate-help)
* [clusterctl generate provider](#clusterctl-generate-provider)
* [clusterctl generate yaml](#clusterctl-generate-yaml)
* [clusterctl get help](#clusterctl-get-help)
* [clusterctl get kubeconfig](#clusterctl-get-kubeconfig)
* [clusterctl init help](#clusterctl-init-help)
* [clusterctl init list-images](#clusterctl-init-list-images)
* [clusterctl upgrade apply](#clusterctl-upgrade-apply)
* [clusterctl upgrade help](#clusterctl-upgrade-help)
* [clusterctl upgrade plan](#clusterctl-upgrade-plan)

# Cluster Management Commands:

## `clusterctl delete`

Delete one or more providers from the management cluster.

```text
clusterctl delete [providers] [flags]
```

### Command Flags

```text
      --addon strings               Add-on providers and versions (e.g. helm:v0.1.0) to delete from the management cluster
      --all                         Force deletion of all the providers
  -b, --bootstrap strings           Bootstrap providers and versions (e.g. kubeadm:v1.1.5) to delete from the management cluster
  -c, --control-plane strings       ControlPlane providers and versions (e.g. kubeadm:v1.1.5) to delete from the management cluster
      --core string                 Core provider version (e.g. cluster-api:v1.1.5) to delete from the management cluster
  -h, --help                        help for delete
      --include-crd                 Forces the deletion of the provider's CRDs (and of all the related objects)
      --include-namespace           Forces the deletion of the namespace where the providers are hosted (and of all the contained objects)
  -i, --infrastructure strings      Infrastructure providers and versions (e.g. aws:v0.5.0) to delete from the management cluster
      --ipam strings                IPAM providers and versions (e.g. infoblox:v0.0.1) to delete from the management cluster
      --kubeconfig string           Path to the kubeconfig file to use for accessing the management cluster. If unspecified, default discovery rules apply.
      --kubeconfig-context string   Context to be used within the kubeconfig file. If empty, current context will be used.
      --runtime-extension strings   Runtime extension providers and versions (e.g. test:v0.0.1) to delete from the management cluster
```

## `clusterctl generate`

Generate yaml using clusterctl yaml processor.

```text
clusterctl generate [flags]
```

### Command Flags

```text
  -h, --help   help for generate
```

## `clusterctl get`

Get info from a management or workload cluster

```text
clusterctl get [flags]
```

### Command Flags

```text
  -h, --help   help for get
```

## `clusterctl init`

Initialize a management cluster.

Installs Cluster API core components, the kubeadm bootstrap provider,
and the selected bootstrap and infrastructure providers.

The management cluster must be an existing Kubernetes cluster, make sure
to have enough privileges to install the desired components.

Use 'clusterctl config repositories' to get a list of available providers and their configuration; if
necessary, edit $XDG_CONFIG_HOME/cluster-api/clusterctl.yaml file to add new provider or to customize existing ones.

Some providers require environment variables to be set before running clusterctl init.
Refer to the provider documentation, or use 'clusterctl generate provider --infrastructure [name] --describe'
to get a list of required variables.

See https://cluster-api.sigs.k8s.io for more details.

```text
clusterctl init [flags]
```

### Command Flags

```text
      --addon strings               Add-on providers and versions (e.g. helm:v0.1.0) to add to the management cluster.
  -b, --bootstrap strings           Bootstrap providers and versions (e.g. kubeadm:v1.1.5) to add to the management cluster. If unspecified, Kubeadm bootstrap provider's latest release is used.
  -c, --control-plane strings       Control plane providers and versions (e.g. kubeadm:v1.1.5) to add to the management cluster. If unspecified, the Kubeadm control plane provider's latest release is used.
      --core string                 Core provider version (e.g. cluster-api:v1.1.5) to add to the management cluster. If unspecified, Cluster API's latest release is used.
  -h, --help                        help for init
  -i, --infrastructure strings      Infrastructure providers and versions (e.g. aws:v0.5.0) to add to the management cluster.
      --ipam strings                IPAM providers and versions (e.g. in-cluster:v0.1.0) to add to the management cluster.
      --kubeconfig string           Path to the kubeconfig for the management cluster. If unspecified, default discovery rules apply.
      --kubeconfig-context string   Context to be used within the kubeconfig file. If empty, current context will be used.
      --runtime-extension strings   Runtime extension providers and versions to add to the management cluster; please note that clusterctl doesn't include any default runtime extensions and thus it is required to use custom configuration files to register runtime extensions.
  -n, --target-namespace string     The target namespace where the providers should be deployed. If unspecified, the provider components' default namespace is used.
      --validate                    If true, clusterctl will validate that the deployments will succeed on the management cluster. (default true)
      --wait-provider-timeout int   Wait timeout per provider installation in seconds. This value is ignored if --wait-providers is false (default 300)
      --wait-providers              Wait for providers to be installed.
```

## `clusterctl move`

Move Cluster API objects and all dependencies between management clusters.

Note: The destination cluster MUST have the required provider components installed.

```text
clusterctl move [flags]
```

### Command Flags

```text
      --dry-run                        Enable dry run, don't really perform the move actions
      --from-directory string          Read Cluster API objects and all dependencies from a directory into a management cluster.
  -h, --help                           help for move
      --hide-api-warnings string       Set of API server warnings to hide. Valid sets are "default" (includes metadata.finalizer warnings), "all" , and "none". (default "default")
      --kubeconfig string              Path to the kubeconfig file for the source management cluster. If unspecified, default discovery rules apply.
      --kubeconfig-context string      Context to be used within the kubeconfig file for the source management cluster. If empty, current context will be used.
  -n, --namespace string               The namespace where the workload cluster is hosted. If unspecified, the current context's namespace is used.
      --to-directory string            Write Cluster API objects and all dependencies from a management cluster to directory.
      --to-kubeconfig string           Path to the kubeconfig file to use for the destination management cluster.
      --to-kubeconfig-context string   Context to be used within the kubeconfig file for the destination management cluster. If empty, current context will be used.
```

## `clusterctl upgrade`

Upgrade core and provider components in a management cluster

```text
clusterctl upgrade [flags]
```

### Command Flags

```text
  -h, --help   help for upgrade
```

# Troubleshooting and Debugging Commands:

## `clusterctl describe`

Describe the status of workload clusters.

```text
clusterctl describe [flags]
```

### Command Flags

```text
  -h, --help   help for describe
```

# Other Commands:

## `clusterctl alpha`

These commands correspond to alpha features in clusterctl.

```text
clusterctl alpha [flags]
```

### Command Flags

```text
  -h, --help   help for alpha
```

## `clusterctl completion`

Output shell completion code for the specified shell (bash, zsh or fish).
The shell code must be evaluated to provide interactive completion of
clusterctl commands. This can be done by sourcing it from the
.bash_profile.

```text
clusterctl completion [bash|zsh|fish] [flags]
```

### Command Flags

```text
  -h, --help   help for completion
```

## `clusterctl config`

Display clusterctl configuration.

```text
clusterctl config [flags]
```

### Command Flags

```text
  -h, --help   help for config
```

## `clusterctl help`

Help about any command

```text
clusterctl help [command] [flags]
```

### Command Flags

```text
  -h, --help   help for help
```

## `clusterctl version`

Print clusterctl version

```text
clusterctl version [flags]
```

### Command Flags

```text
  -h, --help            help for version
  -o, --output string   Output format; available options are 'yaml', 'json' and 'short'
```

# Additional Commands

## `clusterctl alpha help`

Help about any command

```text
clusterctl alpha help [command] [flags]
```

### Command Flags

```text
  -h, --help   help for help
```

## `clusterctl alpha rollout`

Manage the rollout of a cluster-api resource.
Valid resource types include:

   * machinedeployment
   * kubeadmcontrolplane

```text
clusterctl alpha rollout SUBCOMMAND [flags]
```

### Command Flags

```text
  -h, --help   help for rollout
```

## `clusterctl alpha rollout help`

Help about any command

```text
clusterctl alpha rollout help [command] [flags]
```

### Command Flags

```text
  -h, --help   help for help
```

## `clusterctl alpha rollout pause`

Mark the provided cluster-api resource as paused.

       Paused resources will not be reconciled by a controller. Use "clusterctl alpha rollout resume" to resume a paused resource. Currently only MachineDeployments and KubeadmControlPlanes support being paused.

```text
clusterctl alpha rollout pause RESOURCE
```

### Command Flags

```text
  -h, --help                        help for pause
      --kubeconfig string           Path to the kubeconfig file to use for accessing the management cluster. If unspecified, default discovery rules apply.
      --kubeconfig-context string   Context to be used within the kubeconfig file. If empty, current context will be used.
  -n, --namespace string            Namespace where the resource(s) reside. If unspecified, the defult namespace will be used.
```

## `clusterctl alpha rollout restart`

Restart of cluster-api resources.

       Resources will be rollout restarted.

```text
clusterctl alpha rollout restart RESOURCE
```

### Command Flags

```text
  -h, --help                        help for restart
      --kubeconfig string           Path to the kubeconfig file to use for accessing the management cluster. If unspecified, default discovery rules apply.
      --kubeconfig-context string   Context to be used within the kubeconfig file. If empty, current context will be used.
  -n, --namespace string            Namespace where the resource(s) reside. If unspecified, the defult namespace will be used.
```

## `clusterctl alpha rollout resume`

Resume a paused cluster-api resource

       Paused resources will not be reconciled by a controller. By resuming a resource, we allow it to be reconciled again. Currently only MachineDeployments and KubeadmControlPlanes support being resumed.

```text
clusterctl alpha rollout resume RESOURCE
```

### Command Flags

```text
  -h, --help                        help for resume
      --kubeconfig string           Path to the kubeconfig file to use for accessing the management cluster. If unspecified, default discovery rules apply.
      --kubeconfig-context string   Context to be used within the kubeconfig file. If empty, current context will be used.
  -n, --namespace string            Namespace where the resource(s) reside. If unspecified, the defult namespace will be used.
```

## `clusterctl alpha topology`

Commands for ClusterClass based clusters

```text
clusterctl alpha topology [flags]
```

### Command Flags

```text
  -h, --help   help for topology
```

## `clusterctl alpha topology help`

Help about any command

```text
clusterctl alpha topology help [command] [flags]
```

### Command Flags

```text
  -h, --help   help for help
```

## `clusterctl config help`

Help about any command

```text
clusterctl config help [command] [flags]
```

### Command Flags

```text
  -h, --help   help for help
```

## `clusterctl config repositories`

Display the list of providers and their repository configurations.

clusterctl ships with a list of known providers; if necessary, edit
$XDG_CONFIG_HOME/cluster-api/clusterctl.yaml file to add a new provider or to customize existing ones.

```text
clusterctl config repositories [flags]
```

### Command Flags

```text
  -h, --help            help for repositories
  -o, --output string   Output format. Valid values: [yaml text]. (default "text")
```

## `clusterctl describe cluster`

Provide an "at glance" view of a Cluster API cluster designed to help the user in quickly
understanding if there are problems and where.
.

```text
clusterctl describe cluster NAME [flags]
```

### Command Flags

```text
  -c, --color                       Enable or disable color output; if not set color is enabled by default only if using tty. The flag is overridden by the NO_COLOR env variable if set.
      --echo                        Show MachineInfrastructure and BootstrapConfig when ready condition is true or it has the Status, Severity and Reason of the machine's object.
      --grouping                    Groups machines when ready condition has the same Status, Severity and Reason. (default true)
  -h, --help                        help for cluster
      --kubeconfig string           Path to a kubeconfig file to use for the management cluster. If empty, default discovery rules apply.
      --kubeconfig-context string   Context to be used within the kubeconfig file. If empty, current context will be used.
  -n, --namespace string            The namespace where the workload cluster is located. If unspecified, the current namespace will be used.
      --show-conditions string      list of comma separated kind or kind/name for which the command should show all the object's conditions (use 'all' to show conditions for everything).
      --show-machinesets            Show MachineSet objects.
      --show-resourcesets           Show cluster resource sets.
      --show-templates              Show infrastructure and bootstrap config templates associated with the cluster.
      --v1beta2                     Use V1Beta2 conditions..
```

## `clusterctl describe help`

Help about any command

```text
clusterctl describe help [command] [flags]
```

### Command Flags

```text
  -h, --help   help for help
```

## `clusterctl generate cluster`

Generate templates for creating workload clusters.

clusterctl ships with a list of known providers; if necessary, edit
$XDG_CONFIG_HOME/cluster-api/clusterctl.yaml to add new provider or to customize existing ones.

Each provider configuration links to a repository; clusterctl uses this information
to fetch templates when creating a new cluster.

```text
clusterctl generate cluster NAME [flags]
```

### Command Flags

```text
      --control-plane-machine-count int    The number of control plane machines for the workload cluster. (default 1)
  -f, --flavor string                      The workload cluster template variant to be used when reading from the infrastructure provider repository. If unspecified, the default cluster template will be used.
      --from string                        The URL to read the workload cluster template from. If unspecified, the infrastructure provider repository URL will be used. If set to '-', the workload cluster template is read from stdin.
      --from-config-map string             The ConfigMap to read the workload cluster template from. This can be used as alternative to read from the provider repository or from an URL
      --from-config-map-key string         The ConfigMap.Data key where the workload cluster template is hosted. If unspecified, "template" will be used
      --from-config-map-namespace string   The namespace where the ConfigMap exists. If unspecified, the current namespace will be used
  -h, --help                               help for cluster
  -i, --infrastructure string              The infrastructure provider to read the workload cluster template from. If unspecified, the default infrastructure provider will be used.
      --kubeconfig string                  Path to a kubeconfig file to use for the management cluster. If empty, default discovery rules apply.
      --kubeconfig-context string          Context to be used within the kubeconfig file. If empty, current context will be used.
      --kubernetes-version string          The Kubernetes version to use for the workload cluster. If unspecified, the value from OS environment variables or the $XDG_CONFIG_HOME/cluster-api/clusterctl.yaml config file will be used.
      --list-variables                     Returns the list of variables expected by the template instead of the template yaml
  -n, --target-namespace string            The namespace to use for the workload cluster. If unspecified, the current namespace will be used.
      --worker-machine-count int           The number of worker machines for the workload cluster. (default 0)
      --write-to string                    Specify the output file to write the template to, defaults to STDOUT if the flag is not set
```

## `clusterctl generate help`

Help about any command

```text
clusterctl generate help [command] [flags]
```

### Command Flags

```text
  -h, --help   help for help
```

## `clusterctl generate provider`

Generate templates for provider components.

clusterctl fetches the provider components from the provider repository and performs variable substitution.

Variable values are either sourced from the clusterctl config file or
from environment variables

```text
clusterctl generate provider [flags]
```

### Command Flags

```text
      --addon string               Add-on provider and version (e.g. helm:v0.1.0)
  -b, --bootstrap string           Bootstrap provider and version (e.g. kubeadm:v1.1.5)
  -c, --control-plane string       ControlPlane provider and version (e.g. kubeadm:v1.1.5)
      --core string                Core provider and version (e.g. cluster-api:v1.1.5)
      --describe                   Generate configuration without variable substitution.
  -h, --help                       help for provider
  -i, --infrastructure string      Infrastructure provider and version (e.g. aws:v0.5.0)
      --ipam string                IPAM provider and version (e.g. infoblox:v0.0.1)
      --raw                        Generate configuration without variable substitution in a yaml format.
      --runtime-extension string   Runtime extension provider and version (e.g. test:v0.0.1)
  -n, --target-namespace string    The target namespace where the provider should be deployed. If unspecified, the components default namespace is used.
      --write-to string            Specify the output file to write the template to, defaults to STDOUT if the flag is not set
```

## `clusterctl generate yaml`

Process yaml using clusterctl's yaml processor.

clusterctl ships with a simple yaml processor that performs variable
substitution that takes into account of default values.

Variable values are either sourced from the clusterctl config file or
from environment variables

```text
clusterctl generate yaml [flags]
```

### Command Flags

```text
      --from string      The URL to read the template from. It defaults to '-' which reads from stdin. (default "-")
  -h, --help             help for yaml
      --list-variables   Returns the list of variables expected by the template instead of the template yaml
```

## `clusterctl get help`

Help about any command

```text
clusterctl get help [command] [flags]
```

### Command Flags

```text
  -h, --help   help for help
```

## `clusterctl get kubeconfig`

Gets the kubeconfig file for accessing a workload cluster

```text
clusterctl get kubeconfig NAME [flags]
```

### Command Flags

```text
  -h, --help                        help for kubeconfig
      --kubeconfig string           Path to the kubeconfig file to use for accessing the management cluster. If unspecified, default discovery rules apply.
      --kubeconfig-context string   Context to be used within the kubeconfig file. If empty, current context will be used.
  -n, --namespace string            Namespace where the workload cluster exist.
```

## `clusterctl init help`

Help about any command

```text
clusterctl init help [command] [flags]
```

### Command Flags

```text
  -h, --help   help for help
```

## `clusterctl init list-images`

Lists the container images required for initializing the management cluster.

See https://cluster-api.sigs.k8s.io for more details.

```text
clusterctl init list-images [flags]
```

### Command Flags

```text
  -h, --help   help for list-images
```

## `clusterctl upgrade apply`

The upgrade apply command applies new versions of Cluster API providers as defined by clusterctl upgrade plan.

New version should be applied ensuring all the providers uses the same cluster API version
in order to guarantee the proper functioning of the management cluster.

	Specifying the provider using namespace/name:version is deprecated and will be dropped in a future release.

```text
clusterctl upgrade apply [flags]
```

### Command Flags

```text
      --addon strings               Add-on providers and versions (e.g. helm:v0.1.0) to upgrade to. This flag can be used as alternative to --contract.
  -b, --bootstrap strings           Bootstrap providers instance and versions (e.g. kubeadm:v1.1.5) to upgrade to. This flag can be used as alternative to --contract.
      --contract string             The API Version of Cluster API (contract, e.g. v1alpha4) the management cluster should upgrade to
  -c, --control-plane strings       ControlPlane providers instance and versions (e.g. kubeadm:v1.1.5) to upgrade to. This flag can be used as alternative to --contract.
      --core string                 Core provider instance version (e.g. cluster-api:v1.1.5) to upgrade to. This flag can be used as alternative to --contract.
  -h, --help                        help for apply
  -i, --infrastructure strings      Infrastructure providers instance and versions (e.g. aws:v2.0.1) to upgrade to. This flag can be used as alternative to --contract.
      --ipam strings                IPAM providers and versions (e.g. infoblox:v0.0.1) to upgrade to. This flag can be used as alternative to --contract.
      --kubeconfig string           Path to the kubeconfig file to use for accessing the management cluster. If unspecified, default discovery rules apply.
      --kubeconfig-context string   Context to be used within the kubeconfig file. If empty, current context will be used.
      --runtime-extension strings   Runtime extension providers and versions (e.g. test:v0.0.1) to upgrade to. This flag can be used as alternative to --contract.
      --wait-provider-timeout int   Wait timeout per provider upgrade in seconds. This value is ignored if --wait-providers is false (default 300)
      --wait-providers              Wait for providers to be upgraded.
```

## `clusterctl upgrade help`

Help about any command

```text
clusterctl upgrade help [command] [flags]
```

### Command Flags

```text
  -h, --help   help for help
```

## `clusterctl upgrade plan`

The upgrade plan command provides a list of recommended target versions for upgrading the
      Cluster API providers in a management cluster.

All the providers should be supporting the same API Version of Cluster API (contract) in order
      to guarantee the proper functioning of the management cluster.

Then, for each provider, the following upgrade options are provided:
- The latest patch release for the current API Version of Cluster API (contract).
- The latest patch release for the next API Version of Cluster API (contract), if available.

```text
clusterctl upgrade plan [flags]
```

### Command Flags

```text
  -h, --help                        help for plan
      --kubeconfig string           Path to the kubeconfig file to use for accessing the management cluster. If empty, default discovery rules apply.
      --kubeconfig-context string   Context to be used within the kubeconfig file. If empty, current context will be used.
```
