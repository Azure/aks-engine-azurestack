
# Input 
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>
<RemovedFeatureGate>{{removed_feature_gates}}</RemovedFeatureGate>

# Code snippet Filter:
   - source code path: `pkg/api/defaults-kubelet.go`
   - object name: setKubeletConfig
   - object type: func
   - begin with: `func (cs *ContainerService) setKubeletConfig`

# Input Validation
- Extract the Kubernetes version from the <KubernetesVersion> XML tag.
- Inspect the setKubeletConfig function to look for the statement `if common.IsKubernetesVersionGe(o.OrchestratorVersion, "[MAJOR][MINOR].0")`. For instance, for Kubernetes version 1.30.10, check for `if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.30.0")`.
    - If this statement is NOT present, return "True"; if it is present, return "False".

Function **setKubeletConfig** put removed feature gate into list invalidFeatureGates for [MAJOR].[MINOR].[REVISION]

# Checklist
Here's the checklist for adding to the function **setKubeletConfig** for [MAJOR].[MINOR].[REVISION]:

- [ ] Location function **setKubeletConfig**
- [ ] Extract removed feature gate names from the `<RemovedFeatureGate>` XML tag as a comma-separated string (e.g., from "FeatureGate01=true,FeatureGate02=true" extract "FeatureGate01" and "FeatureGate02").
- [ ] Add logic put removed feature gate into list invalidFeatureGates
  Use the following template to add logic for the new version and feature gates:

```
if common.IsKubernetesVersionGe(o.OrchestratorVersion, "[MAJOR][MINOR].0") {
    // Remove --feature-gate <FeatureGateName> starting with [MAJOR][MINOR]
    invalidFeatureGates = append(invalidFeatureGates, "<FeatureGateName>")
    // Repeat for each removed feature gate
}
```
## Example

Here's an example for version [MAJOR].[MINOR].[REVISION] (e.g., 1.30.10):


	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "[MAJOR][MINOR].0") {
		// Remove --feature-gate KubeletPodResources starting with [MAJOR][MINOR]
		invalidFeatureGates = append(invalidFeatureGates, "KubeletPodResources")

		// Remove --feature-gate KubeletPodResourcesGetAllocatable starting with [MAJOR][MINOR]
		invalidFeatureGates = append(invalidFeatureGates, "KubeletPodResourcesGetAllocatable")

		// Remove --feature-gate LegacyServiceAccountTokenTracking starting with [MAJOR][MINOR]
		invalidFeatureGates = append(invalidFeatureGates, "LegacyServiceAccountTokenTracking")

		// Remove --feature-gate MinimizeIPTablesRestore starting with [MAJOR][MINOR]
		invalidFeatureGates = append(invalidFeatureGates, "MinimizeIPTablesRestore")

		// Remove --feature-gate ProxyTerminatingEndpoints starting with [MAJOR][MINOR]
		invalidFeatureGates = append(invalidFeatureGates, "ProxyTerminatingEndpoints")

		// Remove --feature-gate RemoveSelfLink starting with [MAJOR][MINOR]
		invalidFeatureGates = append(invalidFeatureGates, "RemoveSelfLink")

		// Remove --feature-gate SecurityContextDeny starting with [MAJOR][MINOR]
		invalidFeatureGates = append(invalidFeatureGates, "SecurityContextDeny")

		// Remove --feature-gate APISelfSubjectReview starting with [MAJOR][MINOR]
		invalidFeatureGates = append(invalidFeatureGates, "APISelfSubjectReview")

		// Remove --feature-gate CSIMigrationAzureFile  starting with [MAJOR][MINOR]
		invalidFeatureGates = append(invalidFeatureGates, "CSIMigrationAzureFile")

		// Remove --feature-gate ExpandedDNSConfig starting with [MAJOR][MINOR]
		invalidFeatureGates = append(invalidFeatureGates, "ExpandedDNSConfig")

		// Remove --feature-gate ExperimentalHostUserNamespaceDefaulting starting with [MAJOR][MINOR]
		invalidFeatureGates = append(invalidFeatureGates, "ExperimentalHostUserNamespaceDefaulting")

		// Remove --feature-gate IPTablesOwnershipCleanup starting with [MAJOR][MINOR]
		invalidFeatureGates = append(invalidFeatureGates, "IPTablesOwnershipCleanup")
	}


**After making changes, you MUST review the checklist to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**

