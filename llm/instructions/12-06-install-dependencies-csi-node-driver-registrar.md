# Input 
<CSIImages>{{csi_image_versions}}</CSIImages>

# Input Validation Requirements

- Extract the target `csi-node-driver-registrar` version from the `<CSIImages>` XML tag, specifically from the JSON key `csi-node-driver-registrar`.
- Examine the current `CSI_NODE_DRIVER_REGISTRAR_VERSIONS` list.
- **Version Existence Check**:
  - If the desired version is **not present** in `CSI_NODE_DRIVER_REGISTRAR_VERSIONS`, return `"True"`.
  - If the desired version **is present** in the list, return `"False"`.
- **Do not add code to implement the Input Validation logic.**
  
# Code Snippt Filter:
   - source code path: `vhd/packer/install-dependencies.sh`
   - object name: installCSINodeDriverRegistrar
   - object type: func
   - begin with: `installCSINodeDriverRegistrar() {`


# Fundamental Rules

- **It is crucial to keep the indentation consistent with the existing format when making any changes.**

## Component Version Checklist

- [ ] Extract the target `csi-node-driver-registrar` version from the `<CSIImages>` XML tag, specifically from the JSON key `csi-node-driver-registrar`.
    - Review the current `CSI_NODE_DRIVER_REGISTRAR_VERSIONS` list.
    - Remove the lowest (oldest) Kubernetes version from `CSI_NODE_DRIVER_REGISTRAR_VERSIONS`.
    - Add the new Kubernetes version entry at the top of the `CSI_NODE_DRIVER_REGISTRAR_VERSIONS` list, ensuring correct indentation.
    - Verify that the list is properly formatted and indented.

# Examples
## **Example: Add Kubernetes v1.31.0 when it does not exist (following preceding rules)**

**Before:**
CSI_NODE_DRIVER_REGISTRAR_VERSIONS="
1.30.10
1.29.15
"
**After:**
CSI_NODE_DRIVER_REGISTRAR_VERSIONS="
1.31.0
1.30.10
"

## **Example: Attempting to add Kubernetes v1.30.10 when it already exists (no change needed)**

**Before:**
CSI_NODE_DRIVER_REGISTRAR_VERSIONS="
1.31.0
1.30.10
"
**After:**
CSI_NODE_DRIVER_REGISTRAR_VERSIONS="
1.31.0
1.30.10
"

## **Example: Incorrect (should NOT happen)**

**Before:**
CSI_NODE_DRIVER_REGISTRAR_VERSIONS="
1.30.10
1.29.15
"
**After (Incorrect):**
CSI_NODE_DRIVER_REGISTRAR_VERSIONS="
1.30.10
1.30.10
"
*This is incorrect, it should no change because the target version already in the list*



**After making changes, you MUST review the **Component version Check ist** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**
