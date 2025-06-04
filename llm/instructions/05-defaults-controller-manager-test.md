
# Input 
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>
<RemovedFeatureGate>{{removed_feature_gates}}</RemovedFeatureGate>

# Code Snippt Filter:
   - source code path: `pkg/api/defaults-controller-manager_test.go`
   - object name: TestControllerManagerConfigFeatureGates
   - object type: func
   - begin with: `func TestControllerManagerConfigFeatureGates`

# Input Validation
- Get Kubernetes version in xml tag <KubernetesVersion>
- Ensure the Kubernetes version is in the format [MAJOR].[MINOR].[REVISION]. If the version starts with a leading 'v' (e.g., v1.31.8), remove the 'v'.
- Get Removed feature gate in xml tag <RemovedFeatureGate>
- If a test for [MAJOR].[MINOR].[REVISION] does not exist, add the test logic for [MAJOR].[MINOR].[REVISION] at the end of the function TestControllerManagerConfigFeatureGates.

Function **TestControllerManagerConfigFeatureGates** verify the elimination of removed feature gates for [MAJOR].[MINOR].[REVISION]

Here's the checklist for adding a test logic at the end of the function **TestControllerManagerConfigFeatureGates** for [MAJOR].[MINOR].[REVISION]:

Code Structure:

- [ ] Location function **TestControllerManagerConfigFeatureGates**
- [ ] Add test logic at end of function **TestControllerManagerConfigFeatureGates**
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
	cs.setControllerManagerConfig()
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
	cs.setControllerManagerConfig()
	if a["--feature-gates"] != featuregate[MAJOR][MINOR-1]Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"[MAJOR][MINOR-1].0", featuregate[MAJOR][MINOR], a["--feature-gates"], featuregate[MAJOR][MINOR-1]Sanitized)
	}


**After making changes, you MUST review the checklist to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**

