# Input 
<CSIImages>{{csi_image_versions}}</CSIImages>

# Input Validation Requirements

- Extract the target `azuredisk-csi` version from the `<CSIImages>` XML tag, specifically from the JSON key `azuredisk-csi`.
- Examine the current `AZUREDISK_CSI_VERSIONS` list.
- **Version Existence Check**:
  - If the desired version is **not present** in `AZUREDISK_CSI_VERSIONS`, return `"True"`.
  - If the desired version **is present** in the list, return `"False"`.
- **Do not add code to implement the Input Validation logic.**
  
# Code Snippt Filter:
   - source code path: `vhd/packer/install-dependencies.sh`
   - object name: installCSIAzureDisk
   - object type: func
   - begin with: `installCSIAzureDisk() {`


# Fundamental Rules

- **It is crucial to keep the indentation consistent with the existing format when making any changes.**

## Component Version Checklist

- [ ] Extract the target `azuredisk-csi` version from the `<CSIImages>` XML tag, specifically from the JSON key `azuredisk-csi`.
    - Review the current `AZUREDISK_CSI_VERSIONS` list.
    - Remove the lowest (oldest) Kubernetes version from `AZUREDISK_CSI_VERSIONS`.
    - Add the new Kubernetes version entry at the top of the `AZUREDISK_CSI_VERSIONS` list, ensuring correct indentation.
    - Verify that the list is properly formatted and indented.

# Examples
## **Example: Add Kubernetes v1.31.0 when it does not exist (following preceding rules)**

**Before:**
AZUREDISK_CSI_VERSIONS="
1.30.10
1.29.15
"
**After:**
AZUREDISK_CSI_VERSIONS="
1.31.0
1.30.10
"

## **Example: Attempting to add Kubernetes v1.30.10 when it already exists (no change needed)**

**Before:**
AZUREDISK_CSI_VERSIONS="
1.31.0
1.30.10
"
**After:**
AZUREDISK_CSI_VERSIONS="
1.31.0
1.30.10
"

## **Example: Incorrect (should NOT happen)**

**Before:**
AZUREDISK_CSI_VERSIONS="
1.30.10
1.29.15
"
**After (Incorrect):**
AZUREDISK_CSI_VERSIONS="
1.30.10
1.30.10
"
*This is incorrect, it should no change because the target version already in the list*



**After making changes, you MUST review the **Component version Check ist** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**
