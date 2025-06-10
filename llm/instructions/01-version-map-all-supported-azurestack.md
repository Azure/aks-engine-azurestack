

# Input 
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>

# Input Validation
- Get Kubernetes version in xml tag <KubernetesVersion>
- Ensure the Kubernetes version is in the format [MAJOR].[MINOR].[REVISION]. If the version starts with a leading 'v' (e.g., v1.31.8), remove the 'v'.
- **Version Existence Check**: Search the `AllKubernetesSupportedVersionsAzureStack` find the list of keys for `AllKubernetesSupportedVersionsAzureStack` . 
    - If the desired version DO NOT exist in the list, return "True".
    - If the desired version exists in the array, return "False". 
  - **Do not add code to implement the Input Validation logic.**
  
# Code Snippt Filter:
   - source code path: `pkg/api/common/versions.go`
   - object name: AllKubernetesSupportedVersionsAzureStack
   - object type: map
   - begin with: `var AllKubernetesSupportedVersionsAzureStack = map[string]bool`


## Version Update Check list

1. Version Map Updates (in `pkg/api/common/versions.go`):

   - [ ]**If the new Kubernetes version already exists in the map, make no changes and return the original code exactly as received.**
   - [ ] Add new version [MAJRO].[Minor].[REVISION] to AllKubernetesSupportedVersionsAzureStack with `true`
   - [ ] Keep previous version N-1 ([MAJRO].[Minor-1].x) as `true`
   - [ ] Set all older versions ([MAJRO].[Minor-2].x and below) to `false`
   - [ ] DO NOT remove any older version key in
   - [ ] Preserve the original code's spacing and formatting.
   - [ ] Update maps in AllKubernetesSupportedVersionsAzureStack and put new Kubernetes Version end of map AllKubernetesSupportedVersionsAzureStack

**You must review and ensure that all items on the **Version Update Check list** are checked. If any items are not checked, make the necessary changes to ensure all checkboxes are checked.**

## Example: Adding Kubernetes 1.30

Here's an example of required changes when adding Kubernetes 1.30.10:

1. **Version Map Updates**:

// Before
var AllKubernetesSupportedVersionsAzureStack = map[string]bool{
    ...
    "1.27.x": false,
    "1.28.20": true,
    "1.29.15": true,
}

// After
var AllKubernetesSupportedVersionsAzureStack = map[string]bool{
     ...
     "1.27.x": false,
     "1.28.20": false,
     "1.29.15": true,
     "1.30.10": true,
}

**After making changes, you MUST review the **Version Update Check list** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**
