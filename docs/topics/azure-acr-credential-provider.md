# Image Credential Providers

Kubernetes v1.29 requires users to migrate to [Out-of-Tree Credential Providers][KEP]
in order to pull container images from private Azure Container Registry (ACR) instances.

For AKS Engine-based clusters migration to the ACR credential provider will be transparent.
AKS Engine will remove the deprecated kubelet flag `--azure-container-registry-config` from the cluster API Model
and replace it with flags `--image-credential-provider-bin-dir` and `--image-credential-provider-config`.

See [Cloud Provider Azure][CPA] documentation for more details.

## Image Credential Provider Configuration

The `Image Credential Provider Configuration` file instructs kubelet how to authenticate with Azure Container Registry
(see [credential-provider-config.yaml][CPC] for the default configuration file).

### Custom Image Credential Provider Configuration

To use a custom `Image Credential Provider Configuration` file, combine [CustomFiles](/examples/customfiles/README.md) 
and kublet flag `--image-credential-provider-config`
(default location `/var/lib/kubelet/credential-provider-config.yaml`).

> **Note:** The `Image Credential Provider Configuration` file will be provisioned to all nodes in the cluster.

```json
{
  "orchestratorProfile": {
    "kubernetesConfig": {
      "kubeletConfig": {
        "--image-credential-provider-config": "/var/lib/kubelet/credential-provider-config.yaml",
      }
    },
  },
  "masterProfile": {
    "customFiles": [
      {
        "source" : "/local/path/to/my/credential-provider-config.yaml",
        "dest" : "/var/lib/kubelet/credential-provider-config.yaml"
      }
    ]
  }
}
```

[KEP]: https://github.com/kubernetes/enhancements/tree/master/keps/sig-cloud-provider/2133-out-of-tree-credential-provider
[CPA]: https://cloud-provider-azure.sigs.k8s.io/topics/credential-provider/
[CPC]: /parts/k8s/cloud-init/artifacts/credential-provider-config.yaml
