# Update Default Kubernetes Version Constants

## Objective

Modify the specified Go source file to update the default Kubernetes version string constants. These defaults are used when a user does not specify a Kubernetes version during cluster creation.

## Skip Conditions

Skip this step if the default Kubernetes version constants in `pkg/api/common/const.go` are already set to the `major.minor` part of the `PREVIOUS_LATEST_VERSION` and align with the project's versioning strategy:

- `KubernetesDefaultRelease`
- `KubernetesDefaultReleaseWindows`
- `KubernetesDefaultReleaseAzureStack`
- `KubernetesDefaultReleaseWindowsAzureStack`

## File to Modify

- **Path:** `pkg/api/common/const.go`

## Variables to Update

The following string constant variables need to be updated:

- `KubernetesDefaultRelease`
- `KubernetesDefaultReleaseWindows`
- `KubernetesDefaultReleaseAzureStack`
- `KubernetesDefaultReleaseWindowsAzureStack`

## Default Version Logic

- From Step 1, both the `NEW_VERSION` and the `PREVIOUS_LATEST_VERSION` are set to `true` in the version support maps
- Default version constants must point to the `PREVIOUS_LATEST_VERSION`
- This ensures the default version is the second most recent, actively supported Kubernetes release
- This strategy balances access to newer features with established stability

## Detailed Instructions

1. **Identify Target Versions:**

   - For `KubernetesDefaultRelease` and `KubernetesDefaultReleaseWindows`:
     Use `PREVIOUS_LATEST_VERSION` from `AllKubernetesSupportedVersions`
   - For `KubernetesDefaultReleaseAzureStack`:
     Use `PREVIOUS_LATEST_VERSION` from `AllKubernetesSupportedVersionsAzureStack`
   - For `KubernetesDefaultReleaseWindowsAzureStack`:
     Use `PREVIOUS_LATEST_VERSION` from `AllKubernetesWindowsSupportedVersionsAzureStack`

2. **Update Constants:**
   - Change each constant to use only the `major.minor` portion of its corresponding `PREVIOUS_LATEST_VERSION`
   - Example: If `PREVIOUS_LATEST_VERSION` is "1.30.10", use "1.30"

## Example

Given `NEW_VERSION = "1.31.8"` and `PREVIOUS_LATEST_VERSION = "1.30.10"`:

```golang
// KubernetesDefaultRelease is the default Kubernetes version for Linux agent pools
const KubernetesDefaultRelease = "1.30"

// KubernetesDefaultReleaseWindows is the default Kubernetes version for Windows agent pools
const KubernetesDefaultReleaseWindows = "1.30"

// KubernetesDefaultReleaseAzureStack is the default Kubernetes version for Azure Stack
const KubernetesDefaultReleaseAzureStack = "1.30"

// KubernetesDefaultReleaseWindowsAzureStack is the default Kubernetes version for Windows on Azure Stack
const KubernetesDefaultReleaseWindowsAzureStack = "1.30"
```

## Validation Checks

1. All constants should use only `major.minor` format (e.g., "1.30")
2. All constants should point to the `PREVIOUS_LATEST_VERSION` from Step 1
3. Constants should maintain their original comments
4. Changes should be consistent across all four constants
