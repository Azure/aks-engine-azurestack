
# Input 
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>

# Input Validation
- Get Kubernetes version in xml tag <KubernetesVersion>
- Ensure the Kubernetes version is in the format [MAJOR].[MINOR].[REVISION]. If the version starts with a leading 'v' (e.g., v1.31.8), remove the 'v'.

# Code Snippt Filter:
   - source code path: `pkg/api/common/const.go`
   - object name: AllKubernetesSupportedVersions
   - object type: const
   - begin with: `KubernetesDefaultRelease string`

# Version Constants


1. Default Version Updates:

   - [ ] Set all default version constants to [MAJRO].[Minor-1]:
     - KubernetesDefaultRelease
     - KubernetesDefaultReleaseWindows
     - KubernetesDefaultReleaseAzureStack
     - KubernetesDefaultReleaseWindowsAzureStack

**You must review and ensure that all items on the **Version Update Check list** are checked. If any items are not checked, make the necessary changes to ensure all checkboxes are checked.**

## Example: Adding Kubernetes 1.30

Here's an example of required changes when adding Kubernetes 1.30.10:


2. **Default Version Updates**:

// Before
const KubernetesDefaultRelease = "1.28"

// After
const KubernetesDefaultRelease = "1.29"

**After making changes, you MUST review the **Version Update Check list** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**
