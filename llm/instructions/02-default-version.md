
# Input 
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>

# Input Validation
- Retrieve the Kubernetes version from the <KubernetesVersion> XML tag.
- For each of the following default version constants, check if they are set to [MAJOR].[MINOR-1]:
   - KubernetesDefaultRelease
   - KubernetesDefaultReleaseWindows
   - KubernetesDefaultReleaseAzureStack
   - KubernetesDefaultReleaseWindowsAzureStack
- If any of these constants are not set to [MAJOR].[MINOR-1], return "True"; otherwise, return "False".

# Code snippet Filter:
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
