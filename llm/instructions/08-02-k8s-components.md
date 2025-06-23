

# Input 
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>
<CSIImages>{{csi_image_versions}}</CSIImages>

# Input Validation
- Get Kubernetes version in xml tag <KubernetesVersion>
- Examinate the map `common.KubernetesImageBaseTypeMCR` to search key `"[MAJOR][MINOR]": {")`
	- If DOES NOT exist, return "True"; otherwise, return "False".

# Code snippet Filter:
- **File Path:** `pkg/api/k8s_versions.go`
- **Target Object:** `kubernetesImageBaseVersionedImages` (map variable)
- **Object Type:** `map[string]map[string]map[string]string`
- **Declaration:** Begins with `var kubernetesImageBaseVersionedImages = map[string]map[string]map[string]string{`

# Instructions

## Overview
Add a new Kubernetes version entry to the `kubernetesImageBaseVersionedImages` map in `pkg/api/k8s_versions.go`. This map contains container image versions for different Kubernetes releases.

## Step-by-Step Process

### 1. Extract Version Information
- Parse the Kubernetes version from `<KubernetesVersion>` tag to get `[MAJOR].[MINOR]` format
- Extract CSI image versions from `<CSIImages>` tag for the new components

### 2. Create New Version Entry
- Add a new map entry with key `"[MAJOR].[MINOR]"` (e.g., `"1.30"`)
- Place this new entry at the **top** of the `kubernetesImageBaseVersionedImages` map

### 3. Populate CSI-Related Images
Use the versions from `<CSIImages>` tag for these components:
- `common.CSIProvisionerContainerName`
- `common.CSIAttacherContainerName`
- `common.CSILivenessProbeContainerName`
- `common.CSILivenessProbeWindowsContainerName`
- `common.CSINodeDriverRegistrarContainerName`
- `common.CSINodeDriverRegistrarWindowsContainerName`
- `common.CSISnapshotterContainerName`
- `common.CSISnapshotControllerContainerName`
- `common.CSIResizerContainerName`
- `common.CSIAzureDiskContainerName`

### 4. Copy Non-CSI Images
For these components, copy the values from the previous version `[MAJOR].[MINOR-1]`:
- `common.AddonResizerComponentName`
- `common.MetricsServerAddonName`
- `common.AddonManagerComponentName`
- `common.ClusterAutoscalerAddonName`

### 5. Maintain Map Size
- Keep only the **5 most recent versions** in the map
- Remove the oldest version when adding a new one (e.g., when adding "1.30", remove "1.25" if it exists)

# Example: Adding Kubernetes Version 1.30

When adding Kubernetes version 1.30.10, create the following entry at the top of the map:

```go
"1.30": {
    common.CSIProvisionerContainerName:                "oss/kubernetes-csi/csi-provisioner:v5.2.0",
    common.CSIAttacherContainerName:                   "oss/kubernetes-csi/csi-attacher:v4.8.0",
    common.CSILivenessProbeContainerName:              "oss/kubernetes-csi/livenessprobe:v2.15.0",
    common.CSILivenessProbeWindowsContainerName:       "oss/kubernetes-csi/livenessprobe:v2.15.0",
    common.CSINodeDriverRegistrarContainerName:        "oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
    common.CSINodeDriverRegistrarWindowsContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
    common.CSISnapshotterContainerName:                "oss/kubernetes-csi/csi-snapshotter:v8.2.0",
    common.CSISnapshotControllerContainerName:         "oss/kubernetes-csi/snapshot-controller:v8.2.0",
    common.CSIResizerContainerName:                    "oss/kubernetes-csi/csi-resizer:v1.13.1",
    common.CSIAzureDiskContainerName:                  "oss/kubernetes-csi/azuredisk-csi:v1.30.10",
    common.AddonResizerComponentName:                  "oss/kubernetes/autoscaler/addon-resizer:1.8.7",
    common.MetricsServerAddonName:                     "oss/kubernetes/metrics-server:v0.5.2",
    common.AddonManagerComponentName:                  "oss/kubernetes/kube-addon-manager:v9.1.6",
    common.ClusterAutoscalerAddonName:                 "oss/kubernetes/autoscaler/cluster-autoscaler:v1.22.1",
},
```

## Verification Checklist

After completing the changes, verify:
- [ ] New version entry `"[MAJOR].[MINOR]"` added at the top of the map
- [ ] All CSI-related images use versions from `<CSIImages>` tag
- [ ] Non-CSI images copied from previous version
- [ ] Map contains only 5 most recent versions
- [ ] Oldest version removed if necessary
- [ ] Syntax is valid Go code with proper formatting

