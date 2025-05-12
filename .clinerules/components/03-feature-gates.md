# Remove Deprecated Kubernetes Feature Gates

## Objective

Modify the specified Go source files to add logic for removing feature gates that have been deprecated or have become generally available (GA) in the `NEW_VERSION`. This involves appending the names of these feature gates to an `invalidFeatureGates` slice if the cluster's Kubernetes version is `NEW_VERSION`'s `major.minor.0` or newer.

## Skip Conditions

For each file listed below, skip modification if the logic to remove the specified deprecated/GA feature gates for the `NEW_VERSION` already exists and is correctly configured (e.g., an `if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.31.0")` block containing the correct feature gates).

## Files to Modify

The logic for removing deprecated feature gates needs to be applied to:

- `pkg/api/defaults-kubelet.go`
- `pkg/api/defaults-apiserver.go`
- `pkg/api/defaults-cloud-controller-manager.go`
- `pkg/api/defaults-controller-manager.go`
- `pkg/api/defaults-scheduler.go`

## Detailed Instructions

1. **Identify Deprecated/GA Feature Gates:**

   - Acquire the list of feature gates no longer valid for the `NEW_VERSION`
   - Note any reference URLs (e.g., Kubernetes PR links) for documentation
   - Identify which component(s) each feature gate applies to

2. **Construct Code Blocks:**
   For each file:
   - Locate feature gates section with `invalidFeatureGates` slice
   - Add new version check block:
     ```golang
     if common.IsKubernetesVersionGe(o.OrchestratorVersion, "NEW_VERSION.0") {
         // Add relevant feature gates for this component
         invalidFeatureGates = append(invalidFeatureGates, "FeatureGateName")
     }
     ```
   - Include reference URLs as comments

## Example Feature Gates

For Kubernetes 1.31.0:

```golang
// In pkg/api/defaults-apiserver.go
if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.31.0") {
    // Remove feature-gate ServerSideApply starting with 1.31
    // Reference: https://github.com/kubernetes/kubernetes/pull/xxxxx
    invalidFeatureGates = append(invalidFeatureGates, "ServerSideApply")

    // Remove feature-gate ValidatingAdmissionPolicy starting with 1.31
    // Reference: https://github.com/kubernetes/kubernetes/pull/yyyyy
    invalidFeatureGates = append(invalidFeatureGates, "ValidatingAdmissionPolicy")
}
```

## User Interaction Required

The user must provide:

1. List of feature gates to remove (as JSON array)
2. Any reference URLs for documentation
3. Component specificity (if known)

## Validation Checks

1. Version check uses correct format (`major.minor.0`)
2. Feature gates added to correct component files
3. Reference URLs included as comments
4. Consistent code style with existing blocks
