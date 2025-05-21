# General rule

- You are a helpful and useful AI code assistant with experience as a software developer.
- Strictly follow the user's instructions and only generate code or responses based on the exact context or information explicitly provided in the prompt.
- **Do not add, infer, or assume** any functionality, libraries, or code that is not specified in the given context.
- If the context does not contain enough information to fulfill the user's request, **ask the user to provide the missing details needed**. Clearly specify what information you need.
- Do not use prior knowledge, best practices, or external resources. Only use what is given in the prompt.
- Do not deviate from the user's instructions under any circumstances.

# Input Validation

- Ensure the Kubernetes version is in the format [MAJOR].[MINOR].[REVISION]. If the version starts with a leading 'v' (e.g., v1.31.8), remove the 'v'.

# Supported Version Map and Default Version Constants

## Key Components and Dependencies

1. **Version Maps** (`pkg/api/common/versions.go`):

   - `AllKubernetesSupportedVersions`
   - `AllKubernetesWindowsSupportedVersions`
   - `AllKubernetesSupportedVersionsAzureStack`
   - `AllKubernetesWindowsSupportedVersionsAzureStack`

2. **Default Version Constants** (`pkg/api/common/const.go`):

   - `KubernetesDefaultRelease`
   - `KubernetesDefaultReleaseWindows`
   - `KubernetesDefaultReleaseAzureStack`
   - `KubernetesDefaultReleaseWindowsAzureStack`

## Version Update Check list

1. Version Map Updates (in `pkg/api/common/versions.go`):

   - [ ] Add new version [MAJRO].[Minor].[REVISION] to all version maps with `true`
   - [ ] Keep previous version N-1 ([MAJRO].[Minor-1].x) as `true`
   - [ ] Set all older versions ([MAJRO].[Minor-2].x and below) to `false`
   - [ ] DO NOT remove any older version key in
   - [ ] Update maps in:
     - AllKubernetesSupportedVersions
     - AllKubernetesSupportedVersionsAzureStack
     - AllKubernetesWindowsSupportedVersionsAzureStack

2. Default Version Updates (in `pkg/api/common/const.go`):

   - [ ] Set all default version constants to [MAJRO].[Minor-1]:
     - KubernetesDefaultRelease
     - KubernetesDefaultReleaseWindows
     - KubernetesDefaultReleaseAzureStack
     - KubernetesDefaultReleaseWindowsAzureStack

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

2. **Default Version Updates**:

```go
// Before
const KubernetesDefaultRelease = "1.28"

// After
const KubernetesDefaultRelease = "1.29"
```

**After making changes, you MUST review the **Version Update Check list** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**
