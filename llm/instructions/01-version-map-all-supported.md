
# Input Validation
- Get Kubernetes version in xml tag <KubernetesVersion>
- Ensure the Kubernetes version is in the format [MAJOR].[MINOR].[REVISION]. If the version starts with a leading 'v' (e.g., v1.31.8), remove the 'v'.

# Code Snippt Filter:
   - source code path: `pkg/api/common/versions.go`
   - object name: AllKubernetesSupportedVersions
   - begin with: `var AllKubernetesSupportedVersions = map[string]bool`


## Version Update Check list

1. Version Map Updates (in `pkg/api/common/versions.go`):

   - [ ] Add new version [MAJRO].[Minor].[REVISION] to all version maps with `true`
   - [ ] Keep previous version N-1 ([MAJRO].[Minor-1].x) as `true`
   - [ ] Set all older versions ([MAJRO].[Minor-2].x and below) to `false`
   - [ ] DO NOT remove any older version key in
   - [ ] Update maps in AllKubernetesSupportedVersions and put new Kubernetes Version end of map AllKubernetesSupportedVersions

**You must review and ensure that all items on the **Version Update Check list** are checked. If any items are not checked, make the necessary changes to ensure all checkboxes are checked.**

## Example: Adding Kubernetes 1.30

Here's an example of required changes when adding Kubernetes 1.30.10:

1. **Version Map Updates**:

```go
// Before
var AllKubernetesSupportedVersions = map[string]bool{
    ...
    "1.27.x": false,  // N-2 (unsupported)
    "1.28.20": true,  // N-1
    "1.29.15": true,  // N (current latest)
}

// After
var AllKubernetesSupportedVersions = map[string]bool{
     ...
     "1.27.x": false,  // N-3 (unsupported)
     "1.28.20": false, // N-2 (now unsupported)
     "1.29.15": true,  // N-1
     "1.30.10": true,  // N (new latest)

}
```

**After making changes, you MUST review the **Version Update Check list** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**
