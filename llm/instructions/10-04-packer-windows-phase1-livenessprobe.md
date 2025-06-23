
# Input 
<KubernetesVersion>{{k8s_version_revision}}</KubernetesVersion>
<KubernetesPreviousVersion>{{k8s_previous_version_revision}}</KubernetesPreviousVersion>
<AzureCloudManagerImages>{{cloud_provider_image_versions}}</AzureCloudManagerImages>
<CSIImages>{{csi_image_versions}}</CSIImages>

# Input Validation
  - **Do not add code to implement the Input Validation logic.**
  - Retrieve the desired `livenessprobe` container image version from the `<CSIImages>` XML tag.
  - Review the current entries in the `$imagesToPull` list.
  - **Version Existence Check**: Search the `$imagesToPull` find the list of versions for `livenessprobe` . 
    - If the desired version DO NOT exist in the list, return "True".
    - If the desired version exists in the array, return "False". 

# Code snippet Filter:
   - source code path: `vhd/packer/configure-windows-vhd.ps1`
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
- [ ] For `livenessprobe`:
  - Retrieve the desired `livenessprobe` container image version from the `<CSIImages>` XML tag.
  - Review the current entries in the `$imagesToPull` list.
  - **Check if the `livenessprobe` container image version already exists in the `$imagesToPull` list.**
    - If it exists, **skip the update** for `livenessprobe`.
    - If it does not exist, proceed with the following steps:      
    - Remove the first occurrence of the `livenessprobe` entry from `$imagesToPull`.
      - Add the new `livenessprobe` entry with the updated image version directly below the previous version's position in `$imagesToPull` (maintain correct indentation).
      - **CRITICAL**: Ensure each array entry remains on its own separate line with proper indentation. Do not accidentally merge lines during the edit process.
      - Double-check the list for proper formatting and indentation.


# Examples
## **Example: Add `livenessprobe` v2.15.0 when it does not exist (following preceding rules)**

**Before:**

        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.8.0",
        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.10.0",
        
**After:**

        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.10.0",
        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.15.0",

## **Example: Attempting to add `livenessprobe` v2.15.0 when it already exists (no change needed)**

**Before:**

        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.10.0",
        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.15.0",
        
**After:**

        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.10.0",
        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.15.0",

## **Example: Incorrect (should NOT happen)**

**Before:**

        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.10.0",
        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.15.0",
        
**After (Incorrect):**

        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.15.0",
        "mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v2.15.0",
        
*This is incorrect, it should no change because the target version already in the list*

**IMPORTANT FORMATTING NOTE**: When performing the replacement, ensure that:
1. Each array entry stays on its own separate line
2. Proper indentation (8 spaces) is maintained for each entry
3. No lines are accidentally merged together during the edit process
4. The comma placement and string quotes remain consistent

**You must review and ensure that all items on the **Component version Check list** are checked. If any items are not checked, make the necessary changes to ensure all checkboxes are checked.**


**After making changes, you MUST review the **Component version Check list** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**

