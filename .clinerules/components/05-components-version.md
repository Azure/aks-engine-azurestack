# Update Kubernetes Component Versions

## Objective

Update component versions for a new Kubernetes version in `pkg/api/k8s_versions.go`. The component versions need to be updated in both GCR (Google Container Registry) and MCR (Microsoft Container Registry) sections.

## Input Requirements

1. **Kubernetes Version**

   - The specific Kubernetes version to add (e.g., "1.31")
   - Format: "major.minor"

2. **Component Versions**

   - REQUIRED: Must be provided as comma-separated string of component=version pairs
   - Format: "component1=version1,component2=version2"
   - Example: "csi-provisioner=v5.2.0,csi-attacher=v4.8.0"

3. **Default Behavior**
   - For components not explicitly provided in the input, use values from previous version
   - Previous version is determined as (NEW_VERSION_MAJOR).(NEW_VERSION_MINOR-1)
   - Example:
     - If NEW_VERSION = "1.31"
     - Then use values from version "1.30" for unspecified components

## Files to Modify

- **Path:** `pkg/api/k8s_versions.go`
- **Maps to Update:**
  1. `kubernetesImageBaseVersionedImages[common.KubernetesImageBaseTypeGCR]`
  2. `kubernetesImageBaseVersionedImages[common.KubernetesImageBaseTypeMCR]`

## Component Name Mapping

Map between component names in input string and constants:

```go
componentMapping := map[string]string{
    "csi-provisioner": "common.CSIProvisionerContainerName",
    "csi-attacher": "common.CSIAttacherContainerName",
    "csi-snapshotter": "common.CSISnapshotterContainerName",
    "snapshot-controller": "common.CSISnapshotControllerContainerName",
    "csi-resizer": "common.CSIResizerContainerName",
    "livenessprobe": "common.CSILivenessProbeContainerName",
    "csi-node-driver-registrar": "common.CSINodeDriverRegistrarContainerName",
    "azuredisk-csi": "common.CSIAzureDiskContainerName",
    "azurefile-csi": "common.CSIAzureFileContainerName",
    "addon-resizer": "common.AddonResizerComponentName",
    "metrics-server": "common.MetricsServerAddonName",
    "addon-manager": "common.AddonManagerComponentName",
    "cluster-autoscaler": "common.ClusterAutoscalerAddonName"
}
```

## Image Reference Pattern

Component images should follow these patterns:

1. For GCR:

```
"oss/kubernetes-csi/[component-name]:[version]"
```

2. For MCR:

```
"oss/kubernetes-csi/[component-name]:[version]"
```

## Special Components

1. **Windows Components**

   - For components that have Windows variants, add both regular and Windows versions
   - Example:
     ```go
     common.CSILivenessProbeContainerName: "oss/kubernetes-csi/livenessprobe:v2.15.0",
     common.CSILivenessProbeWindowsContainerName: "oss/kubernetes-csi/livenessprobe:v2.15.0",
     ```

2. **Cloud Provider Components**
   - For cloud-controller-manager and cloud-node-manager:
     ```go
     common.CloudControllerManagerComponentName: "oss/kubernetes/azure-cloud-controller-manager:v[version]",
     common.CloudNodeManagerAddonName: "oss/kubernetes/azure-cloud-node-manager:v[version]",
     ```

## Example Version Entry

```go
"1.31": {
    common.CSIProvisionerContainerName: "oss/kubernetes-csi/csi-provisioner:v5.2.0",
    common.CSIAttacherContainerName: "oss/kubernetes-csi/csi-attacher:v4.8.0",
    common.CSILivenessProbeContainerName: "oss/kubernetes-csi/livenessprobe:v2.15.0",
    common.CSILivenessProbeWindowsContainerName: "oss/kubernetes-csi/livenessprobe:v2.15.0",
    common.CSINodeDriverRegistrarContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
    common.CSINodeDriverRegistrarWindowsContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
    common.CSISnapshotterContainerName: "oss/kubernetes-csi/csi-snapshotter:v8.2.0",
    common.CSISnapshotControllerContainerName: "oss/kubernetes-csi/snapshot-controller:v8.2.0",
    common.CSIResizerContainerName: "oss/kubernetes-csi/csi-resizer:v1.13.2",
    common.CSIAzureDiskContainerName: "oss/kubernetes-csi/azuredisk-csi:v1.31.5",
    common.CSIAzureFileContainerName: "oss/kubernetes-csi/azurefile-csi:v1.31.5",
    common.AddonResizerComponentName: "oss/kubernetes/autoscaler/addon-resizer:1.8.7",
    common.MetricsServerAddonName: "oss/kubernetes/metrics-server:v0.5.2",
    common.AddonManagerComponentName: "oss/kubernetes/kube-addon-manager:v9.1.6",
    common.ClusterAutoscalerAddonName: "oss/kubernetes/autoscaler/cluster-autoscaler:v1.22.1",
}
```

## Validation Checks

1. **Version Format**

   - All versions should start with 'v'
   - Follow semantic versioning (vX.Y.Z)
   - Exception: addon-resizer uses format like "1.8.7"

2. **Image Path Format**

   - GCR format: "oss/kubernetes-csi/[component]:[version]"
   - MCR format: "oss/kubernetes-csi/[component]:[version]"
   - Some components use different paths (e.g., "oss/kubernetes/autoscaler/")

3. **Component Completeness**

   - All components from previous version must be present
   - New components must have both GCR and MCR entries
   - Windows variants must be included where applicable

4. **Version Consistency**
   - Components should use consistent versioning across related components
   - Windows variants should match their non-Windows counterparts
   - Cloud provider components should use consistent versioning pattern
