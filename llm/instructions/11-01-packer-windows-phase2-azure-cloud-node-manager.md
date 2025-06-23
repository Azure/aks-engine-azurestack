

# Input 
<KubernetesVersion>{{k8s_version_revision}}</KubernetesVersion>
<KubernetesPreviousVersion>{{k8s_previous_version_revision}}</KubernetesPreviousVersion>
<AzureCloudManagerImages>{{cloud_provider_image_versions}}</AzureCloudManagerImages>
<CSIImages>{{csi_image_versions}}</CSIImages>

# Input Validation
  - **Do not add code to implement the Input Validation logic.**
  - Retrieve the desired `azure-cloud-node-manager` container image version from the `<AzureCloudManagerImages>` XML tag.
  - Review the current entries in the `$imagesToPull` list.
  - **Version Existence Check**: Search the `$imagesToPull` find the list of versions for `azure-cloud-node-manager` . 
    - If the desired version DO NOT exist in the list, return "True".
    - If the desired version exists in the array, return "False". 
      
# Code snippet Filter:
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

## Newline Preservation Guidelines for String Replacement

**CRITICAL**: To prevent accidental line merging when using `replace_string_in_file`, follow these strategies:

### Strategy 2: Target Single Line Only
Replace only the specific line that needs to change:

```powershell
# CORRECT - Replace only the target line:
oldString: "                \"mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v1.29.9\","
newString: "                \"mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v1.31.7\","
```

## Component version Check list
- [ ] For `azure-cloud-node-manager`:
  - Retrieve the `azure-cloud-node-manager` container image version from the `<AzureCloudManagerImages>` XML tag.
  - Review the current entries in the `$imagesToPull` list.
  - **Check if the `azure-cloud-node-manager` container image version already exists in the `$imagesToPull` list.**
    - If it exists, **skip the update** for `azure-cloud-node-manager`.
    - If it does not exist, proceed with the following steps:      
    - Replace the oldest occurrence of the `azure-cloud-node-manager` entry from `$imagesToPull` with the new version.
      - Update the `azure-cloud-node-manager` entry with the new image version while maintaining the same position in `$imagesToPull` (maintain correct indentation).
      - **CRITICAL**: Ensure each array entry remains on its own separate line with proper indentation. Do not accidentally merge lines during the edit process.
      - Double-check the list for proper formatting and indentation.

Example to add `azure-cloud-node-manager` for v1.30.8

**Before:**

        "mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v1.28.5",
        "mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v1.29.9",
**After:**
        "mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v1.30.8",
        "mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v1.29.9",

**IMPORTANT FORMATTING NOTE**: When performing the replacement, ensure that:
1. Each array entry stays on its own separate line
2. Proper indentation (8 spaces) is maintained for each entry
3. No lines are accidentally merged together during the edit process
4. The comma placement and string quotes remain consistent

**You must review and ensure that all items on the **Component version Check list** are checked. If any items are not checked, make the necessary changes to ensure all checkboxes are checked.**


**After making changes, you MUST review the **Component version Check list** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**

