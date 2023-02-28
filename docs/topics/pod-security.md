# Pod Security Admission

Kubernetes v1.25 [removes the `PodSecurityPolicy` (PSP) admission controller][PSPDeprecation]
in favor of the newer [PodSecurity][PSA] admission controller (PSA).
The PodSecurity admission is enabled by default since Kubernetes v1.23.

Up until Kubernetes v1.25, AKS Engine creates a set of [PodSecurityPolicy resources][PSPParts]
via the `"pod-security-policy"` addon (enabled by default).
The `"pod-security-policy"` addon will be forcefully disabled on Kubernetes v1.25+ clusters.
Hence, cluster administrators interested in enforcing the [Pod Security Standards][PSS]
through the `PodSecurity` admission controller should use a combination of:

- [Namespace Labels][PSALabels]
- [Admission Controller configuration file][PSAConfig]
  - See how to upload a configuration file [below](#admission-controller-configuration-file).

## Migrating from PSP

Migrating from the PodSecurityPolicy admission is a **pre-requisite** to upgrade to Kubernetes v1.25+.

The [step-by-step guide][Migration] published by the Kubernetes project is a good resource
on how to migrate from the PodSecurityPolicy admission
if the target cluster includes customizations beyond AKS Engine's defaults.

If the target AKS Engine-based cluster only includes the configuration applied
by the `"pod-security-policy"` addon, then there are two migration alternatives:

### Option 1: Migrate using `aks-engine-azurestack upgrade`

To migrate using a cluster upgrade operation, the API Model should be updated in the following way:

- Disable the `"pod-security-policy"` addon
- Remove `PodSecurityPolicy` from the plugins list in `--enable-admission-plugins`

```json
{
  "kubernetesConfig": {
    "addons": [
      {
        "name": "pod-security-policy",
        "enabled": false
      }
    ],
    "apiServerConfig": {
      "--enable-admission-plugins": "...,ExtendedResourceToleration",
    }
  }
}
```

To speed up the upgrade process, use the `--control-plane-only`.

To validate pod security before upgrading to Kubernetes v1.25,
perform a "forced" upgrade to the cluster's current version (needs flag `--force`).

```bash
aks-engine-azurestack upgrade \
  --upgrade-version {CurrentKubernetesVersion} \
  --force \
  --control-plane-only \
  ...
```

Once upgrade is over, delete the resources created by the `"pod-security-policy"` addon:

```bash
kubectl delete clusterrolebinding default:restricted default:privileged
kubectl delete clusterrole psp:restricted psp:privileged 
kubectl delete psp restricted privileged
```

### Option 2: Manual Migration

An alternative, faster, migration process is to follow this **ordered** sequence of steps:

1. For each master node:
   1. Delete the `"pod-security-policy"` addon yaml from the `kube-addon-manager` directory
   1. Remove the `PodSecurityPolicy` admission controller from the `kube-apiserver` manifest
1. Update API Model as indicated in the [previous section](#migrate-using-aks-engine-azurestack-upgrade)
1. Delete resources created by the `"pod-security-policy"` addon

```bash
# Install node-shell if preferred over SSH => https://github.com/kvaps/kubectl-node-shell
# kubectl krew index add kvaps https://github.com/kvaps/krew-index
# kubectl krew install kvaps/node-shell

# 1.1. Delete the PSP addon yaml
kubectl node-shell {MasterNodeName} \
  -- sh -c 'rm -f /etc/kubernetes/addons/pod-security-policy.yaml'
# 1.2. Remove PodSecurityPolicy admission controller
kubectl node-shell {MasterNodeName} \
  -- sh -c 'sed -i s/,PodSecurityPolicy//1 /etc/kubernetes/manifests/kube-apiserver.yaml'

# 3. Delete resources created by the PSP addon
kubectl delete clusterrolebinding default:restricted default:privileged
kubectl delete clusterrole psp:restricted psp:privileged 
kubectl delete psp restricted privileged
```

## Admission Controller configuration file

The behavior of the PSA controller can be customized through the `AdmissionConfiguration` file.

The default AKS Engine configuration enforces the `privileged` stardard
(see [apiserver-admission-control.yaml][PSADefaultConfig]]).

### Custom Admission Controller configuration

To use a custom `AdmissionConfiguration` file, combine [CustomFiles](/examples/customfiles/README.md) 
and the Kubernetes API server flag `--admission-control-config-file`
(default location `/etc/kubernetes/apiserver-admission-control.yaml`).

```json
{
  "orchestratorProfile": {
    "kubernetesConfig": {
      "apiServerConfig": {
        "--admission-control-config-file": "/etc/kubernetes/apiserver-admission-control.yaml",
      }
    },
  },
  "masterProfile": {
    "customFiles": [
      {
        "source" : "/local/path/to/my/apiserver-admission-control.yaml",
        "dest" : "/etc/kubernetes/apiserver-admission-control.yaml"
      }
    ]
  }
}
```

[PSPDeprecation]: https://kubernetes.io/blog/2021/04/06/podsecuritypolicy-deprecation-past-present-and-future/
[PSA]: https://kubernetes.io/docs/concepts/security/pod-security-admission/
[PSPParts]: /parts/k8s/addons/pod-security-policy.yaml
[PSS]: https://kubernetes.io/docs/concepts/security/pod-security-standards/
[PSALabels]: https://kubernetes.io/docs/tasks/configure-pod-container/enforce-standards-namespace-labels/
[PSAConfig]: https://kubernetes.io/docs/tasks/configure-pod-container/enforce-standards-admission-controller/
[Migration]: https://kubernetes.io/docs/tasks/configure-pod-container/migrate-from-psp/
[PSADefaultConfig]: /parts/k8s/cloud-init/artifacts/apiserver-admission-control.yaml
