# Input 
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>

# Input Validation
  - Retrieve the desired Kubernetes version from the `<KubernetesVersion>` XML tag.
  - Review the current entries in the `K8S_VERSIONS` list.
  - **Version Existence Check**: Search the `K8S_VERSIONS` find the list of versions for Kubernetes. 
    - If the desired version DO NOT exist in the list, return "True".
    - If the desired version exists in the array, return "False". 
  - **Do not add code to implement the Input Validation logic.**
  
# Code snippet Filter:
   - source code path: `vhd/packer/install-dependencies.sh`
   - object name: installKubeBinaries
   - object type: func
   - begin with: `installKubeBinaries() {`


# Fundamental Rules

- **It is crucial to keep the indentation consistent with the existing format when making any changes.**

## Component Version Checklist

- [ ] For Kubernetes version specified in the `<KubernetesVersion>` XML tag:
    - Review the current `K8S_VERSIONS` list.
    - Remove the lowest (oldest) Kubernetes version from `K8S_VERSIONS`.
    - Add the new Kubernetes version entry at the top of the `K8S_VERSIONS` list, ensuring correct indentation.
    - Verify that the list is properly formatted and indented.

# Examples
## **Example: Add Kubernetes v1.31.0 when it does not exist (following preceding rules)**

**Before:**
K8S_VERSIONS="
1.30.10
1.29.15
"
**After:**
K8S_VERSIONS="
1.31.0
1.30.10
"

## **Example: Attempting to add Kubernetes v1.30.10 when it already exists (no change needed)**

**Before:**
K8S_VERSIONS="
1.31.0
1.30.10
"
**After:**
K8S_VERSIONS="
1.31.0
1.30.10
"

## **Example: Incorrect (should NOT happen)**

**Before:**
K8S_VERSIONS="
1.30.10
1.29.15
"
**After (Incorrect):**
K8S_VERSIONS="
1.30.10
1.30.10
"
*This is incorrect, it should no change because the target version already in the list*



**After making changes, you MUST review the **Component version Check ist** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**
