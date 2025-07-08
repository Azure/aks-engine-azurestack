# Azure Container Registry (ACR) Credential Provider

## Overview

Starting with Kubernetes v1.29, you **must** use [Out-of-Tree Credential Providers][KEP] to pull container images from private Azure Container Registry (ACR) instances. The legacy `--azure-container-registry-config` kubelet flag has been deprecated and removed. See [Cloud Provider Azure][CPA] documentation for more details.

## What This Means for Your Cluster

When you deploy an AKS Engine cluster, the migration to ACR credential providers happens automatically:

 **Automatic Migration**: AKS Engine removes the deprecated `--azure-container-registry-config` flag  
 **New Configuration**: Replaces it with modern credential provider flags:
- `--image-credential-provider-bin-dir`
- `--image-credential-provider-config`

> **No Action Required**: This transition is transparent - your cluster will continue to authenticate with ACR without any manual intervention.

## How It Works

### Configuration File Location

The credential provider uses a configuration file automatically deployed to:
```
/var/lib/kubelet/credential-provider-config.yaml
```

This file is **automatically provisioned** to all nodes (both master and worker) in your cluster.

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
