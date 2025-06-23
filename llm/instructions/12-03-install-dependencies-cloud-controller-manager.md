# Input 
<AzureCloudManagerImages>{{cloud_provider_image_versions}}</AzureCloudManagerImages>

# Input Validation Requirements

- Extract the target `azure-cloud-controller-manager` version from the `<AzureCloudManagerImages>` XML tag, specifically from the JSON key `azure-cloud-controller-manager`.
- Examine the current `CLOUD_CONTROLLER_MANAGER_VERSIONS` list.
- **Version Existence Check**:
  - If the desired version is **not present** in `CLOUD_CONTROLLER_MANAGER_VERSIONS`, return `"True"`.
  - If the desired version **is present** in the list, return `"False"`.
- **Do not add code to implement the Input Validation logic.**
  
# Code snippet Filter:
   - source code path: `vhd/packer/install-dependencies.sh`
   - object name: installCloudControllerManager
   - object type: func
   - begin with: `installCloudControllerManager() {`


# Fundamental Rules

- **It is crucial to keep the indentation consistent with the existing format when making any changes.**

## Component Version Checklist

- [ ] Extract the target `azure-cloud-controller-manager` version from the `<AzureCloudManagerImages>` XML tag, specifically from the JSON key `azure-cloud-controller-manager`.
    - Review the current `CLOUD_CONTROLLER_MANAGER_VERSIONS` list.
    - Remove the lowest (oldest) Kubernetes version from `CLOUD_CONTROLLER_MANAGER_VERSIONS`.
    - Add the new Kubernetes version entry at the top of the `CLOUD_CONTROLLER_MANAGER_VERSIONS` list, ensuring correct indentation.
    - Verify that the list is properly formatted and indented.

# Examples
## **Example: Add Kubernetes v1.31.0 when it does not exist (following preceding rules)**

**Before:**
CLOUD_CONTROLLER_MANAGER_VERSIONS="
1.30.10
1.29.15
"
**After:**
CLOUD_CONTROLLER_MANAGER_VERSIONS="
1.31.0
1.30.10
"

## **Example: Attempting to add Kubernetes v1.30.10 when it already exists (no change needed)**

**Before:**
CLOUD_CONTROLLER_MANAGER_VERSIONS="
1.31.0
1.30.10
"
**After:**
CLOUD_CONTROLLER_MANAGER_VERSIONS="
1.31.0
1.30.10
"

## **Example: Incorrect (should NOT happen)**

**Before:**
CLOUD_CONTROLLER_MANAGER_VERSIONS="
1.30.10
1.29.15
"
**After (Incorrect):**
CLOUD_CONTROLLER_MANAGER_VERSIONS="
1.30.10
1.30.10
"
*This is incorrect, it should no change because the target version already in the list*



**After making changes, you MUST review the **Component version Check ist** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**
