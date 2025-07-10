# Azure Container Registry (ACR) Credential Provider

## Overview

Starting with Kubernetes v1.29, you **must** use [Out-of-Tree Credential Providers][KEP] to pull container images from private Azure Container Registry (ACR) instances. The legacy `--azure-container-registry-config` kubelet flag has been deprecated and removed. See [Cloud Provider Azure][CPA] documentation for more details.

## What This Means for Your Cluster

When you deploy an AKS Engine cluster, the ACR credential provider would be installed and configuraged by default.

When you update your AKS Engine cluster, deprecated kubelet `--azure-container-registry-config` flag is replaced
with flags `--image-credential-provider-bin-dir` and `--image-credential-provider-config`.
This transition should be transparent, the kubelet is expected to authenticate with ACR without any manual intervention.

## How It Works

### Configuration File Location

The credential provider configuration file is provisioned to all nodes in the following locations:

- Linux nodes: `/var/lib/kubelet/credential-provider-config.yaml`
- Windows nodes: `C:\k\credential-provider\credential-provider-config.yaml`

### Default Configuration

The default configuration handles authentication for all Azure Container Registry endpoints:

```yaml
kind: CredentialProviderConfig
apiVersion: kubelet.config.k8s.io/v1
providers:
  - name: azure-acr-credential-provider
    matchImages:
      - "*.azurecr.io"    # Azure Public Cloud
      - "*.azurecr.cn"    # Azure China Cloud  
      - "*.azurecr.de"    # Azure Germany Cloud
      - "*.azurecr.us"    # Azure US Government Cloud
```

> **Reference**: View the complete default configuration file at [`credential-provider-config.yaml`](../../parts/k8s/cloud-init/artifacts/credential-provider-config.yaml)

[KEP]: https://github.com/kubernetes/enhancements/tree/master/keps/sig-cloud-provider/2133-out-of-tree-credential-provider
[CPA]: https://cloud-provider-azure.sigs.k8s.io/topics/credential-provider/
