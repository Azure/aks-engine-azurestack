
# Input 
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>
<RemovedFeatureGate>{{removed_feature_gates}}</RemovedFeatureGate>

# Code Snippt Filter:
   - source code path: `pkg/api/defaults-kubelet_test.go`
   - object name: TestKubeletConfigFeatureGates
   - object type: func
   - begin with: `func TestKubeletConfigFeatureGates`

# Input Validation
- Extract the Kubernetes version from the <KubernetesVersion> XML tag.
- Retrieve the removed feature gate from the <RemovedFeatureGate> XML tag.
- In the **TestKubeletConfigFeatureGates** function, search for the string `featuregate[MAJOR][MINOR] :=` corresponding to the Kubernetes version (e.g., for version 1.30.10, look for `featuregate130 :=`).
    - If this string is NOT found, return "True"; if it is found, return "False".

Function **TestKubeletConfigFeatureGates** verify the elimination of removed feature gates for [MAJOR].[MINOR].[REVISION]

Here's the checklist for adding a test logic at the end of the function **TestKubeletConfigFeatureGates** for [MAJOR].[MINOR].[REVISION]:

Code Structure:

- [ ] Location function **TestKubeletConfigFeatureGates**
- [ ] Add test logic at end of function **TestKubeletConfigFeatureGates**
- [ ] Follow naming convention `featuregate[MAJOR][MINOR]` and `featuregate[MAJOR][MINOR-1]Sanitized`

## Example

Here's an example test case for version [MAJOR].[MINOR].[REVISION] (e.g., 1.30.10):


    // test user-overrides, removal of feature gates for k8s versions >= [MAJOR].[MINOR]
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "[MAJOR].[MINOR].0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	featuregate[MAJOR][MINOR] := "[REMOVED-FEATURE-GATE01=true,REMOVED-FEATURE-GATE02=true]"
	a["--feature-gates"] = featuregate[MAJOR][MINOR]
	featuregate[MAJOR][MINOR]Sanitized := ""
	cs.setKubeletConfig()
	if a["--feature-gates"] != featuregate[MAJOR][MINOR]Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"[MAJOR][MINOR]", featuregate[MAJOR][MINOR], a["--feature-gates"], featuregate[MAJOR][MINOR]Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < [MAJOR].[MINOR]
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "[MAJOR].[MINOR-1].0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	a["--feature-gates"] = featuregate[MAJOR][MINOR]
	featuregate[MAJOR][MINOR-1]Sanitized = featuregate[MAJOR][MINOR]
	cs.setKubeletConfig()
	if a["--feature-gates"] != featuregate[MAJOR][MINOR-1]Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"[MAJOR][MINOR-1].0", featuregate[MAJOR][MINOR], a["--feature-gates"], featuregate[MAJOR][MINOR-1]Sanitized)
	}


**After making changes, you MUST review the checklist to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**

