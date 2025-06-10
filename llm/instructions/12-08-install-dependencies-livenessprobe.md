# Input 
<CSIImages>{{csi_image_versions}}</CSIImages>

# Input Validation Requirements

- Extract the target `livenessprobe` version from the `<CSIImages>` XML tag, specifically from the JSON key `livenessprobe`.
- Examine the current `LIVENESSPROBE_VERSION` list.
- **Version Existence Check**:
  - If the desired version is **not present** in `LIVENESSPROBE_VERSION`, return `"True"`.
  - If the desired version **is present** in the list, return `"False"`.
- **Do not add code to implement the Input Validation logic.**
  
# Code Snippt Filter:
   - source code path: `vhd/packer/install-dependencies.sh`
   - object name: installCSILivenessProbe
   - object type: func
   - begin with: `installCSILivenessProbe() {`


# Fundamental Rules

- **It is crucial to keep the indentation consistent with the existing format when making any changes.**

## Component Version Checklist

- [ ] Extract the target `livenessprobe` version from the `<CSIImages>` XML tag, specifically from the JSON key `livenessprobe`.
    - Review the current `LIVENESSPROBE_VERSION` list.
    - Remove the lowest (oldest) Kubernetes version from `LIVENESSPROBE_VERSION`.
    - Add the new Kubernetes version entry at the top of the `LIVENESSPROBE_VERSION` list, ensuring correct indentation.
    - Verify that the list is properly formatted and indented.

# Examples
## **Example: Add Kubernetes v1.31.0 when it does not exist (following preceding rules)**

**Before:**
LIVENESSPROBE_VERSION="
1.30.10
1.29.15
"
**After:**
LIVENESSPROBE_VERSION="
1.31.0
1.30.10
"

## **Example: Attempting to add Kubernetes v1.30.10 when it already exists (no change needed)**

**Before:**
LIVENESSPROBE_VERSION="
1.31.0
1.30.10
"
**After:**
LIVENESSPROBE_VERSION="
1.31.0
1.30.10
"

## **Example: Incorrect (should NOT happen)**

**Before:**
LIVENESSPROBE_VERSION="
1.30.10
1.29.15
"
**After (Incorrect):**
LIVENESSPROBE_VERSION="
1.30.10
1.30.10
"
*This is incorrect, it should no change because the target version already in the list*



**After making changes, you MUST review the **Component version Check ist** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**
