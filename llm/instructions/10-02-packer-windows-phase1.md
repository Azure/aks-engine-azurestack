# Input 
<KubernetesVersion>{{k8s_version_revision}}</KubernetesVersion>
<KubernetesPreviousVersion>{{k8s_previous_version_revision}}</KubernetesPreviousVersion>
<AzureCloudManagerImages>{{cloud_provider_image_versions}}</AzureCloudManagerImages>
<CSIImages>{{csi_image_versions}}</CSIImages>

# Input Validation
  - Retrieve the desired `azuredisk-csi` container image version from the `<CSIImages>` XML tag.
  - Review the current entries in the `$imagesToPull` list.
  - **Version Existence Check**: Search the `$imagesToPull` find the list of versions for `azuredisk-csi` . 
    - If the desired version DO NOT exist in the list, return "True".
    - If the desired version exists in the array, return "False". 

# Code Snippt Filter:
   - source code path: `vhd/packer/configure-windows-vhd-phase2.ps1`
   - object name: Get-ContainerImages
   - object type: func
   - begin with: `function Get-ContainerImages {`

# Fundamental Rules

- [ ] **Container Images**  
      When adding a new Kubernetes version, you must also add the corresponding Kubernetes component images for that version.  
      You will receive precise instructions regarding the container image key and how to determine the version.  
      Examine the current pattern, add the new container image version to the list, and remove the oldest version from the list.
      **It is crucial to keep the indentation consistent with the existing format when making any changes.**

## Component version Check list
  - [ ] For `azuredisk-csi`:
    - Retrieve the  `azuredisk-csi` container image version from the `<CSIImages>` XML tag.
    - Review the current entries in the `$imagesToPull` list.
    - Identify the first existing `azuredisk-csi` version and its `-windows-hp` variant in the list and remove them
    - Insert the new `azuredisk-csi` and `azuredisk-csi-windows-hp` entries, each with the updated image version, immediately after the existing version in the `$imagesToPull` list. Ensure indentation matches the surrounding entries.
    - Double-check the list for proper formatting and indentation.


Example to add `azuredisk-csi` for v1.30.8

**Before:**
	"mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.28.3",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.28.3-windows-hp",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.29.1",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.29.1-windows-hp",
**After:**
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.29.1",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.29.1-windows-hp",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.31.5",
        "mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v1.31.5-windows-hp",
**You must review and ensure that all items on the **Component version Check list** are checked. If any items are not checked, make the necessary changes to ensure all checkboxes are checked.**


**After making changes, you MUST review the **Component version Check list** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**

