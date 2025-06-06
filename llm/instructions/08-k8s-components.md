

# Input 
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>
<CSIImages>{{csi_image_versions}}</CSIImages>

# Input Validation
- Get Kubernetes version in xml tag <KubernetesVersion>
- Ensure the Kubernetes version is in the format [MAJOR].[MINOR].[REVISION]. If the version starts with a leading 'v' (e.g., v1.31.8), remove the 'v'.

# Code Snippt Filter:
   - source code path: `pkg/api/k8s_versions.go`
   - object name: kubernetesImageBaseVersionedImages
   - object type: map
   - begin with: `var kubernetesImageBaseVersionedImages = map[string]map[string]map[string]string{`


## Component version Check list
	
	- [ ] For each Kubernetes version `[MAJOR].[MINOR].[REVISION]`, both `common.KubernetesImageBaseTypeGCR` and `common.KubernetesImageBaseTypeMCR` contain lists of CSI-related container images.
    - [ ] **Add New Version Entry:** Create a new entry with the key `[MAJOR].[MINOR]` (e.g., `"1.30"`) in both `common.KubernetesImageBaseTypeGCR` and `common.KubernetesImageBaseTypeMCR`.
        - Populate all CSI-related container images for this entry using the versions specified in the `<CSIImages>` XML tag.
    - [ ] **Reuse Non-CSI Images:** For the following non-CSI images, copy the values from the previous version (`[MAJOR].[MINOR-1]`):
        - `common.AddonResizerComponentName`
        - `common.MetricsServerAddonName`
        - `common.AddonManagerComponentName`
        - `common.ClusterAutoscalerAddonName`
    - [ ] **Map Placement:** Insert the new version key at the top of the map.
    - [ ] **Map Removal:**  
        - Retain only the five most recent versions in the map (i.e., `[MAJOR].[MINOR]` through `[MAJOR].[MINOR-4]`).  
        - For example, when adding `1.30`, remove `1.24` if present.

**You must review and ensure that all items on the **Version Update Check list** are checked. If any items are not checked, make the necessary changes to ensure all checkboxes are checked.**

## Example: Adding Kubernetes 1.30.10

**Example Entry for "1.30":**

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


**After making changes, you MUST review the **Version Update Check list** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**

