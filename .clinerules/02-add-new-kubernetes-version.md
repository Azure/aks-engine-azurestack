# Cline Rule: Sequential AI Code and Configuration Modification

## System Role and Cline Rules

### Trigger Keywords

- The process is initiated only when the keyword `AddNewVersion 1.x.y` (e.g., `AddNewVersion 1.31.8`) is present in the user input. If this keyword is absent, the AI must skip all steps.
- To start from a specific step, the keyword is `AddNewVersion Step X` (e.g., `AddNewVersion Step 2`). The AI must begin execution from the specified step.
- To list all steps in summary, the keyword is `AddNewVersion steps`. The AI must provide a numbered list of all steps, each with a one-sentence summary, and instruct the user on how to start from any step.
- **This cline rule apply only if a valid trigger keyword is present in the user input. If you do not see a trigger keyword do not apply any of the following rules.**

### General Principles

- The AI assistant must strictly follow the user's explicit instructions for all code and configuration changes.
- All modifications must be thoroughly documented, with clear, concise explanations that directly align with the user's directives.
- All code provided must be complete, executable, and free from omissions or abbreviations.
- The AI must not skip, abbreviate, or summarize any code or configuration changes; all relevant details must be included.

### Stepwise Execution

- Each major step is defined as a Header 2 under "Steps".
- The AI must perform only one major step at a time.
- Upon completing each step, the AI must confirm its completion and then pause, waiting for the user to provide the keyword **"continue"** before proceeding.
- The AI must not anticipate, summarize, or execute any subsequent steps until the "continue" command is received.
- If the user prompt specifies a starting step, the AI must begin at that step; otherwise, it must start from Step 1 and proceed sequentially.

### User Interaction and Control

- If additional information or clarification is needed, the AI must pause and request the necessary details from the user before proceeding.
- The AI must not conclude the analysis or solution prematurely; it should continue analyzing and refining the solution unless explicitly instructed otherwise by the user.

### Best Practices

- Ensure clarity, accuracy, and adherence to technical conventions in all responses.
- Maintain a consistent and logical structure throughout the process.
- Validate all user input to ensure it aligns with these rules before proceeding.

# Steps

## Step 1: Update Kubernetes Version Support in `versions.go` (Completed Example)

### Objective

**Note:** This step should be skipped if the new Kubernetes version already exists in the target maps (`AllKubernetesSupportedVersions`, `AllKubernetesSupportedVersionsAzureStack`, `AllKubernetesWindowsSupportedVersionsAzureStack`) and is configured correctly as per the instructions (i.e., the new version and the previous latest version are `true`, and other relevant older versions are `false`).

Modify the specified Go source file to add a new Kubernetes version to the supported versions lists. This involves updating several boolean maps to reflect that the new version and the immediately preceding supported version are available for new cluster creation, while older versions are not.

### File to Modify

- **Path:** `pkg/api/common/versions.go`

### Variables to Update

The following map-type variables within the specified file need to be updated:

- `AllKubernetesSupportedVersions`
- `AllKubernetesSupportedVersionsAzureStack`
- `AllKubernetesWindowsSupportedVersionsAzureStack`

### Context: Kubernetes Version Maps

These maps define the Kubernetes versions supported by AKS Engine.

- **Key:** A string representing the Kubernetes version (e.g., `"1.29.15"`).
- **Value:** A boolean:
  - `true`: Indicates that a new Kubernetes cluster can be created with this version.
  - `false`: Indicates that the version was previously supported, but new clusters can no longer be created with it. Existing clusters on this version might still be upgradeable _from_.

### Detailed Instructions

For each of the maps listed in "Variables to Update":

1.  **Identify the Current Latest Supported Version:**

    - Scan the map to find the highest Kubernetes version string that currently has a value of `true`. Let's call this `PREVIOUS_LATEST_VERSION`.

2.  **Add the New Kubernetes Version:**

    - Introduce the new Kubernetes version (e.g., `NEW_VERSION`) as a new key in the map.
    - Set the value for `NEW_VERSION` to `true`.

3.  **Update Existing Version States:**
    - Ensure `PREVIOUS_LATEST_VERSION` remains `true`.
    - For all other existing Kubernetes version keys in the map (i.e., versions that are neither `NEW_VERSION` nor `PREVIOUS_LATEST_VERSION`), set their values to `false`.

### Example (Completed for all maps)

Suppose we are adding `NEW_VERSION = "1.31.8"`.

**File:** `pkg/api/common/versions.go`

**1. Target Map: `AllKubernetesSupportedVersions`**

- **Before Modification (Illustrative):**
  ```golang
  var AllKubernetesSupportedVersions = map[string]bool{
      "1.28.5":  false,
      "1.29.10": false,
      "1.29.15": true,  // PREVIOUS_LATEST_VERSION before 1.30.10
      "1.30.10": true,  // PREVIOUS_LATEST_VERSION before 1.31.8
  }
  ```
- **After Adding `NEW_VERSION = "1.31.8"`:**
  ```golang
  var AllKubernetesSupportedVersions = map[string]bool{
      "1.28.5":  false,
      "1.29.10": false,
      "1.29.15": false, // Becomes false
      "1.30.10": true,  // PREVIOUS_LATEST_VERSION, remains true
      "1.31.8":  true,  // NEW_VERSION, set to true
  }
  ```

## Step 2: Update Default Kubernetes Versions in `const.go`

### Objective

**Note:** This step should be skipped if the default Kubernetes version constants in `pkg/api/common/const.go` (i.e., `KubernetesDefaultRelease`, `KubernetesDefaultReleaseWindows`, `KubernetesDefaultReleaseAzureStack`, `KubernetesDefaultReleaseWindowsAzureStack`) are already set to the `major.minor` part of the `PREVIOUS_LATEST_VERSION` (as determined by the logic in Step 1) and align with the project's versioning strategy.

Modify the specified Go source file to update the default Kubernetes version string constants. These defaults are used when a user does not specify a Kubernetes version during cluster creation.

### File to Modify

- **Path:** `pkg/api/common/const.go`

### Variables to Update

The following string constant variables within the specified file need to be updated:

- `KubernetesDefaultRelease`
- `KubernetesDefaultReleaseWindows`
- `KubernetesDefaultReleaseAzureStack`
- `KubernetesDefaultReleaseWindowsAzureStack`

### Context: Default Version Logic

In Step 1, both the `NEW_VERSION` (the version being added) and the `PREVIOUS_LATEST_VERSION` (the version that was latest before `NEW_VERSION`, and remains supported) are set to `true` in the Kubernetes version support maps (e.g., `AllKubernetesSupportedVersions`).

For this step (Step 2), the default Kubernetes version constants (such as `KubernetesDefaultRelease`) must be updated to point to this `PREVIOUS_LATEST_VERSION`. This strategy ensures that the default version for new clusters is the second most recent, actively supported Kubernetes release. This provides users with a default that is both mature and relevant, balancing access to newer features with established stability.

### Detailed Instructions

1.  **Identify the Target `PREVIOUS_LATEST_VERSION` for Each Default Constant:**

    - For `KubernetesDefaultRelease` and `KubernetesDefaultReleaseWindows`: Use the `PREVIOUS_LATEST_VERSION` identified from the `AllKubernetesSupportedVersions` map in Step 1 (this is the version that was latest before `NEW_VERSION`, and remains `true`).
    - For `KubernetesDefaultReleaseAzureStack`: Use the `PREVIOUS_LATEST_VERSION` identified from the `AllKubernetesSupportedVersionsAzureStack` map in Step 1.
    - For `KubernetesDefaultReleaseWindowsAzureStack`: Use the `PREVIOUS_LATEST_VERSION` identified from the `AllKubernetesWindowsSupportedVersionsAzureStack` map in Step 1.

2.  **Update the String Constants:**
    - In `pkg/api/common/const.go`, change the string value assigned to each of the listed constants to their corresponding `PREVIOUS_LATEST_VERSION` identified above.
    - **Important:** Ensure that only the `major.minor` portion of the version (e.g., "1.30" from "1.30.10") is used for these constants.

### Example

Continuing with the scenario from Step 1 where `NEW_VERSION = "1.31.8"` was added:

**Assumed `PREVIOUS_LATEST_VERSION` values from Step 1's example:**

- From `AllKubernetesSupportedVersions`: `"1.30.10"`
- From `AllKubernetesSupportedVersionsAzureStack`: `"1.30.10"`
- From `AllKubernetesWindowsSupportedVersionsAzureStack`: `"1.30.10"`

**File:** `pkg/api/common/const.go`

**Before Modification (Illustrative):**

```golang
// KubernetesDefaultRelease is the default Kubernetes version for Linux agent pools
const KubernetesDefaultRelease = "1.29"
// KubernetesDefaultReleaseWindows is the default Kubernetes version for Windows agent pools
const KubernetesDefaultReleaseWindows = "1.29"
// KubernetesDefaultReleaseAzureStack is the default Kubernetes version for Azure Stack
const KubernetesDefaultReleaseAzureStack = "1.29"
// KubernetesDefaultReleaseWindowsAzureStack is the default Kubernetes version for Windows on Azure Stack
const KubernetesDefaultReleaseWindowsAzureStack = "1.29"

```

**After Modification (Illustrative):**

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

## Step 3: Remove Deprecated Kubernetes Feature Gates

### Objective

**Note:** For each file listed below, this step (or the part pertaining to that specific file) should be skipped if the logic to remove the specified deprecated/GA feature gates for the `NEW_VERSION` (e.g., an `if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.31.0")` block or similar, containing the correct feature gates) already exists and is correctly configured in that file.

Modify the specified Go source files to add logic for removing feature gates that have been deprecated or have become generally available (GA) in the `NEW_VERSION` (e.g., 1.31.x). This involves appending the names of these feature gates to an `invalidFeatureGates` slice (or similarly named variable if different per file) if the cluster's Kubernetes version is `NEW_VERSION`'s `major.minor.0` (e.g., "1.31.0") or newer.

### Files to Modify

The logic for removing deprecated feature gates needs to be applied to the following files:

- **Path:** `pkg/api/defaults-kubelet.go`
- **Path:** `pkg/api/defaults-apiserver.go`
- **Path:** `pkg/api/defaults-cloud-controller-manager.go`
- **Path:** `pkg/api/defaults-controller-manager.go`
- **Path:** `pkg/api/defaults-scheduler.go`

### Detailed Instructions

1.  **Identify Deprecated/GA Feature Gates:**

    - The user will provide a list of feature gates (as a JSON array of strings) that are no longer valid for the `NEW_VERSION` (e.g., for 1.31.x, this would typically apply to version "1.31.0" and onwards).
    - For each feature gate, if a reference URL (e.g., a Kubernetes PR link detailing its removal/deprecation) is provided, it should be included as a comment in the generated code.

2.  **Construct the Code Block (for each relevant file):**
    - For each file listed in "Files to Modify":
      - Locate the section where feature gates are conditionally added to an `invalidFeatureGates` slice (or a similarly purposed variable) based on Kubernetes versions.
      - Add a new `if` block that checks if the orchestrator version (e.g., `o.OrchestratorVersion`, or the relevant version variable in that file's context) is greater than or equal to the `major.minor.0` of the `NEW_VERSION` (e.g., `"1.31.0"`).
      - Inside this `if` block, for each feature gate identified in the previous step that is relevant to the component (Kubelet, API Server, etc.), add a line to append it to the `invalidFeatureGates` slice (or equivalent). Include a comment with the reference URL if available.
      - **Note:** The specific variable names (e.g., `o.OrchestratorVersion`, `invalidFeatureGates`) might differ slightly between files. The AI should adapt to the conventions within each file.

### Example

Suppose `NEW_VERSION = "1.31.8"`, and the user provides the following feature gates to be removed:
`["OldFeatureA", "AnotherOldFeatureB"]` (Note: The relevance of a feature gate might be component-specific).
With `AnotherOldFeatureB` having a reference: `https://github.com/kubernetes/kubernetes/pull/xxxxx`

**File (Illustrative for `pkg/api/defaults-kubelet.go`):** `pkg/api/defaults-kubelet.go`
(Similar logic would apply to other files like `pkg/api/defaults-apiserver.go`, etc., adapting to their specific structure and relevant feature gates)

**Proposed Code Addition (Illustrative for one file):**

```golang
// ... (existing code in the specific file, e.g., defaults-kubelet.go) ...

// Example for Kubelet feature gates
if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.30.0") {
    // ... (existing logic for 1.30.0) ...
}

// Add this new block for NEW_VERSION (e.g., 1.31.0)
if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.31.0") {
    // Remove --feature-gate OldFeatureA starting with 1.31 (if applicable to Kubelet)
    invalidFeatureGates = append(invalidFeatureGates, "OldFeatureA")

    // Remove --feature-gate AnotherOldFeatureB starting with 1.31 (if applicable to Kubelet)
    // Reference: https://github.com/kubernetes/kubernetes/pull/xxxxx
    invalidFeatureGates = append(invalidFeatureGates, "AnotherOldFeatureB")
}

// ... (rest of the function in that file) ...
```

3.  **User Interaction:**
    - The AI will first ask the user to provide the JSON array of feature gate names and any corresponding reference URLs. The user should ideally specify if a feature gate is specific to a component (Kubelet, API Server, etc.) if known.
    - After receiving this information, for each relevant file listed in "Files to Modify":
      - The AI will read the file.
      - The AI will determine if the feature gates provided are applicable to the component configured by that file.
      - The AI will propose the exact code changes for that file.
    - After proposing changes for all relevant files, the AI will wait for user confirmation before proceeding to apply them in ACT MODE.

## Step 4: Add Unit Tests to Verify Feature Gate Removal

### Objective

Add unit tests across all components to verify feature gate removal behavior for NEW_VERSION. These tests will verify:

1. Feature gates are removed when k8s version >= NEW_VERSION's major.minor.0
2. Feature gates are preserved when k8s version < NEW_VERSION's major.minor.0

**Note:** For each test file listed below, this step (or the part pertaining to that specific test file) should be skipped if tests for NEW_VERSION already exist in that file.

### Files to Modify

1. `pkg/api/defaults-apiserver_test.go`
2. `pkg/api/defaults-cloud-controller-manager_test.go`
3. `pkg/api/defaults-controller-manager_test.go`
4. `pkg/api/defaults-kubelet_test.go`
5. `pkg/api/defaults-scheduler_test.go`

### Example Feature Gates for 1.31.0

The following feature gates are to be removed in Kubernetes 1.31.0:

```
CloudDualStackNodeIPs
DRAControlPlaneController
HPAContainerMetrics
KMSv2
KMSv2KDF
LegacyServiceAccountTokenCleanUp
MinDomainsInPodTopologySpread
NewVolumeManagerReconstruction
NodeOutOfServiceVolumeDetach
PodHostIPs
ServerSideApply
ServerSideFieldValidation
StableLoadBalancerNodeSet
ValidatingAdmissionPolicy
ZeroLimitedNominalConcurrencyShares
```

### Test Pattern

Each test file should include:

```go
// test user-overrides, removal of feature gates for k8s versions >= NEW_VERSION
cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
cs.Properties.OrchestratorProfile.OrchestratorVersion = "NEW_VERSION.0"  // e.g., "1.31.0"
cs.Properties.OrchestratorProfile.KubernetesConfig.<ComponentConfig> = make(map[string]string)
component := cs.Properties.OrchestratorProfile.KubernetesConfig.<ComponentConfig>
featuregateNew := "CloudDualStackNodeIPs=true,DRAControlPlaneController=true,HPAContainerMetrics=true,KMSv2=true,KMSv2KDF=true,LegacyServiceAccountTokenCleanUp=true,MinDomainsInPodTopologySpread=true,NewVolumeManagerReconstruction=true,NodeOutOfServiceVolumeDetach=true,PodHostIPs=true,ServerSideApply=true,ServerSideFieldValidation=true,StableLoadBalancerNodeSet=true,ValidatingAdmissionPolicy=true,ZeroLimitedNominalConcurrencyShares=true"
component["--feature-gates"] = featuregateNew
featuregateNewSanitized := ""  // Empty string because gates should be removed for NEW_VERSION+
cs.set<Component>Config()
if component["--feature-gates"] != featuregateNewSanitized {
    t.Fatalf("got unexpected '--feature-gates' for %s \n <component> config original value %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
        "NEW_VERSION.0", featuregateNew, component["--feature-gates"], featuregateNewSanitized)
}

// test user-overrides, no removal of feature gates for k8s versions < NEW_VERSION
cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
cs.Properties.OrchestratorProfile.OrchestratorVersion = "PREVIOUS_VERSION.0"  // e.g., "1.30.0"
cs.Properties.OrchestratorProfile.KubernetesConfig.<ComponentConfig> = make(map[string]string)
component = cs.Properties.OrchestratorProfile.KubernetesConfig.<ComponentConfig>
component["--feature-gates"] = featuregateNew
featuregateOldSanitized := featuregateNew  // Same as input because gates should be preserved
cs.set<Component>Config()
if component["--feature-gates"] != featuregateOldSanitized {
    t.Fatalf("got unexpected '--feature-gates' for %s \n <component> config original value %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
        "PREVIOUS_VERSION.0", featuregateNew, component["--feature-gates"], featuregateOldSanitized)
}
```

### Detailed Instructions

For each test file:

1. **Check for Existing Tests:**

   - First verify that the file doesn't already have tests for the NEW_VERSION.
   - If tests for NEW_VERSION exist, skip modifying this file.

2. **Add Test Cases:**

   - If no tests exist for NEW_VERSION, add test cases following the pattern above.
   - Replace placeholder values:
     - `NEW_VERSION.0`: Use the major.minor.0 version (e.g., "1.31.0")
     - `PREVIOUS_VERSION.0`: Use the previous major.minor.0 version (e.g., "1.30.0")
     - `<ComponentConfig>`: Use the appropriate config map for the component
     - `<Component>`: Use the appropriate component name
     - `featuregateNew`: Use the complete list of feature gates that will be removed in NEW_VERSION

3. **Component-Specific Adjustments:**
   - APIServer: Use `APIServerConfig` and `setAPIServerConfig`
   - CloudControllerManager: Use `CloudControllerManagerConfig` and `setCloudControllerManagerConfig`
   - ControllerManager: Use `ControllerManagerConfig` and `setControllerManagerConfig`
   - Kubelet: Use `KubeletConfig` and `setKubeletConfig`
   - Scheduler: Use `SchedulerConfig` and `setSchedulerConfig`
