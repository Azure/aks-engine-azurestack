# Input
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>
<CSIImages>{{csi_image_versions}}</CSIImages>

# Input Validation
- Extract Kubernetes version ([MAJOR].[MINOR].[REVISION]) from xml tag <KubernetesVersion>
- Extract the target `azuredisk-csi` version from the `<CSIImages>` XML tag, specifically from the JSON key `azuredisk-csi`
- Use extracted CSI driver version as [CSIDRIVERVERSION]
- Validate current CSI driver version assignment:
   - if line `DRIVER_VERSION=v[CSIDRIVERVERSION] # if using k8s >= v[MAJOR].[MINOR]` already exists, then skip update (return False)
   - else proceed with update (return True)

# Target File Location:
- source code path: `docs/topics/azure-stack.md`
- section: Azure Disk CSI Driver Examples

# Azure Stack CSI Driver Update Checklist

1. CSI Driver Version Update (in `docs/topics/azure-stack.md`):
   - [ ] Locate the version assignment for previous minor release `# if using k8s >= v[MAJOR].[MINOR-1]`
   - [ ] Add new line after it with proper indentation: `DRIVER_VERSION=v[CSIDRIVERVERSION] # if using k8s >= v[MAJOR].[MINOR]`
   - [ ] Preserve the original code's spacing and formatting
   - [ ] Maintain consistent commenting format

# Example: Updating CSI Driver version for Kubernetes 1.30.10 with azuredisk-csi v1.31.5

Input validation:
- k8s_version: `1.30.10` → [MAJOR].[MINOR] = `1.30`
- CSI driver version: `1.31.5` → [CSIDRIVERVERSION] = `1.31.5`

Required transformation when updating CSI driver for k8s v1.30:

**Before:**
```bash
DRIVER_VERSION=v1.26.5 # if using k8s v1.26
DRIVER_VERSION=v1.28.3 # if using k8s v1.27
DRIVER_VERSION=v1.29.1 # if using k8s >= v1.28
```

**After:**
```bash
DRIVER_VERSION=v1.26.5 # if using k8s v1.26
DRIVER_VERSION=v1.28.3 # if using k8s v1.27
DRIVER_VERSION=v1.29.1 # if using k8s >= v1.28
DRIVER_VERSION=v1.31.5 # if using k8s >= v1.30
```


# Verification
- **All items on the Azure Stack CSI Driver Update Checklist have been verified and completed.**
- **The new CSI driver version line has been added with proper indentation and format.**
- **The previous versions remain unchanged and properly formatted.**
