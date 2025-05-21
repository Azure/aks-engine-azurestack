# General rule

- You are a helpful and useful AI code assistant with experience as a software developer.
- Strictly follow the user's instructions and only generate code or responses based on the exact context or information explicitly provided in the prompt.
- **Do not add, infer, or assume** any functionality, libraries, or code that is not specified in the given context.
- If the context does not contain enough information to fulfill the user's request, **ask the user to provide the missing details needed**. Clearly specify what information you need.
- Do not use prior knowledge, best practices, or external resources. Only use what is given in the prompt.
- Do not deviate from the user's instructions under any circumstances.

# Input Validation

- Ensure the Kubernetes version is in the format [MAJOR].[MINOR].[REVISION]. If the version starts with a leading 'v' (e.g., v1.31.8), remove the 'v'.
- If the removed feature list is not available, prompt the user to provide the removed feature list as a comma-separated string. Use this input to generate the feature gate list in the code.
  - For example: - User input: "REMOVED-FEATURE-GATE01,REMOVED-FEATURE-GATE02" - Convert to: "[REMOVED-FEATURE-GATE01=true,REMOVED-FEATURE-GATE02=true]"
    In pkg\api\defaults-apiserver_test.go, add test cases at the end of the function TestAPIServerFeatureGates.
    If a test for [MAJOR].[MINOR].[REVISION] does not exist, add the test logic for [MAJOR].[MINOR].[REVISION] at the end of the function TestAPIServerFeatureGates.

pkg\api\defaults-apiserver_test.go , function **TestAPIServerFeatureGates** verify the elimination of removed feature gates for [MAJOR].[MINOR].[REVISION]

Here's the checklist for adding a test logic at the end of the function **TestAPIServerFeatureGates** for [MAJOR].[MINOR].[REVISION]:

Code Structure:

- [ ] Location function **TestAPIServerFeatureGates**
- [ ] Add test logic at end of function **TestAPIServerFeatureGates**
- [ ] Follow naming convention `featuregate[MAJOR][MINOR]` and `featuregate[MAJOR][MINOR-1]Sanitized`

## Example

Here's an example test case for version [MAJOR].[MINOR].[REVISION] (e.g., 1.30.10):

```go
    // test user-overrides, removal of feature gates for k8s versions >= [MAJOR].[MINOR]
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "[MAJOR].[MINOR].0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	featuregate[MAJOR][MINOR] := "[REMOVED-FEATURE-GATE01=true,REMOVED-FEATURE-GATE02=true]"
	a["--feature-gates"] = featuregate[MAJOR][MINOR]
	featuregate[MAJOR][MINOR]Sanitized := ""
	cs.setAPIServerConfig()
	if a["--feature-gates"] != featuregate[MAJOR][MINOR]Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.30.0", featuregate[MAJOR][MINOR], a["--feature-gates"], featuregate[MAJOR][MINOR]Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < [MAJOR].[MINOR]
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "[MAJOR].[MINOR-1].0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	a["--feature-gates"] = featuregate[MAJOR][MINOR]
	featuregate[MAJOR][MINOR-1]Sanitized = featuregate[MAJOR][MINOR]
	cs.setAPIServerConfig()
	if a["--feature-gates"] != featuregate129Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.29.0", featuregate[MAJOR][MINOR], a["--feature-gates"], featuregate[MAJOR][MINOR-1]Sanitized)
	}
```

**After making changes, you MUST review the checklist to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**

the following is the content of pkg\api\defaults-apiserver_test.go

```golang
func TestAPIServerFeatureGates(t *testing.T) {
	// test defaultTestClusterVer
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setAPIServerConfig()
	a := cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	if a["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for k8s v%s: %s",
			defaultTestClusterVer, a["--feature-gates"])
	}

	// test 1.19.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.19.0"
	cs.setAPIServerConfig()
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	if a["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for k8s v%s: %s",
			"1.19.0", a["--feature-gates"])
	}

	// test 1.22.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.22.0"
	cs.setAPIServerConfig()
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	if a["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for k8s v%s: %s",
			"1.22.0", a["--feature-gates"])
	}

	// test 1.25.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.25.0"
	cs.setAPIServerConfig()
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	if a["--feature-gates"] != "PodSecurity=true" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for k8s v%s: %s",
			"1.25.0", a["--feature-gates"])
	}

	// test 1.26.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.26.0"
	cs.setAPIServerConfig()
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	if a["--feature-gates"] != "PodSecurity=true" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for k8s v%s: %s",
			"1.26.0", a["--feature-gates"])
	}

	// test user-overrides, removal of VolumeSnapshotDataSource for k8s versions >= 1.22
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.22.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	a["--feature-gates"] = "VolumeSnapshotDataSource=true"
	cs.setAPIServerConfig()
	if a["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for \"--feature-gates\": \"VolumeSnapshotDataSource=true\": %s for k8s v%s",
			a["--feature-gates"], "1.22.0")
	}

	// test user-overrides, no removal of VolumeSnapshotDataSource for k8s versions < 1.22
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.19.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	a["--feature-gates"] = "VolumeSnapshotDataSource=true"
	cs.setAPIServerConfig()
	if a["--feature-gates"] != "VolumeSnapshotDataSource=true" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for \"--feature-gates\": \"VolumeSnapshotDataSource=true\": %s for k8s v%s",
			a["--feature-gates"], "1.19.0")
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.27.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	a["--feature-gates"] = "ControllerManagerLeaderMigration=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,IdentifyPodOS=true,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,StatefulSetMinReadySeconds=true"
	cs.setAPIServerConfig()
	if a["--feature-gates"] != "PodSecurity=true" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for \"--feature-gates\": %s for k8s v%s",
			a["--feature-gates"], "1.27.0")
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.26.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	a["--feature-gates"] = "ControllerManagerLeaderMigration=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,IdentifyPodOS=true,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,StatefulSetMinReadySeconds=true"
	cs.setAPIServerConfig()
	if a["--feature-gates"] != "CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,ControllerManagerLeaderMigration=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,IdentifyPodOS=true,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,PodSecurity=true,StatefulSetMinReadySeconds=true" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for \"--feature-gates\": %s for k8s v%s",
			a["--feature-gates"], "1.26.0")
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.28
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.28.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	featuregate128 := "AdvancedAuditing=true,CSIMigrationGCE=true,CSIStorageCapacity=true,DelegateFSGroupToCSIDriver=true,DevicePlugins=true,DisableAcceleratorUsageMetrics=true,DryRun=true,EndpointSliceTerminatingCondition=true,KubeletCredentialProviders=true,MixedProtocolLBService=true,NetworkPolicyStatus=true,PodHasNetworkCondition=true,PodSecurity=true,ServiceIPStaticSubrange=true,ServiceInternalTrafficPolicy=true,UserNamespacesStatelessPodsSupport=true,WindowsHostProcessContainers=true"
	a["--feature-gates"] = featuregate128
	featuregate128Sanitized := ""
	cs.setAPIServerConfig()
	if a["--feature-gates"] != featuregate128Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, expected sanitized value: %s \n, actual sanitized value: %s \n ",
			"1.28.0", featuregate128, a["--feature-gates"], featuregate128Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.27.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	a["--feature-gates"] = featuregate128
	featuregate127Sanitized := featuregate128
	cs.setAPIServerConfig()
	if a["--feature-gates"] != featuregate127Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, expected sanitized value: %s \n, actual sanitized value: %s \n ",
			"1.27.0", featuregate128, a["--feature-gates"], featuregate127Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.29
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.29.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	featuregate129 := "CSIMigrationvSphere=true,CronJobTimeZone=true,DownwardAPIHugePages=true,GRPCContainerProbe=true,JobMutableNodeSchedulingDirectives=true,JobTrackingWithFinalizers=true,LegacyServiceAccountTokenNoAutoGeneration=true,OpenAPIV3=true,ProbeTerminationGracePeriod=true,RetroactiveDefaultStorageClass=true,SeccompDefault=true,TopologyManager=true"
	a["--feature-gates"] = featuregate129
	featuregate129Sanitized := ""
	cs.setAPIServerConfig()
	if a["--feature-gates"] != featuregate129Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.29.0", featuregate129, a["--feature-gates"], featuregate129Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.29
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.28.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	a["--feature-gates"] = featuregate129
	featuregate128Sanitized = featuregate129
	cs.setAPIServerConfig()
	if a["--feature-gates"] != featuregate128Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.28.0", featuregate129, a["--feature-gates"], featuregate128Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.30
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.30.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	featuregate130 := "APISelfSubjectReview=true,CSIMigrationAzureFile=true,ExpandedDNSConfig=true,ExperimentalHostUserNamespaceDefaulting=true,IPTablesOwnershipCleanup=true,KubeletPodResources=true,KubeletPodResourcesGetAllocatable=true,LegacyServiceAccountTokenTracking=true,MinimizeIPTablesRestore=true,ProxyTerminatingEndpoints=true,RemoveSelfLink=true,SecurityContextDeny=true"
	a["--feature-gates"] = featuregate130
	featuregate130Sanitized := ""
	cs.setAPIServerConfig()
	if a["--feature-gates"] != featuregate130Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.30.0", featuregate130, a["--feature-gates"], featuregate130Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.30
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.29.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	a["--feature-gates"] = featuregate130
	featuregate129Sanitized = featuregate130
	cs.setAPIServerConfig()
	if a["--feature-gates"] != featuregate129Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.29.0", featuregate130, a["--feature-gates"], featuregate129Sanitized)
	}
}
```

make the necessary changes and return the full TestAPIServerFeatureGates code

Add Kuberetnes Version 1.31.8, removed feature gates: "CloudDualStackNodeIPs=true,DRAControlPlaneController=true,HPAContainerMetrics=true,KMSv2=true,KMSv2KDF=true,LegacyServiceAccountTokenCleanUp=true,MinDomainsInPodTopologySpread=true,NewVolumeManagerReconstruction=true,NodeOutOfServiceVolumeDetach=true,PodHostIPs=true,ServerSideApply=true,ServerSideFieldValidation=true,StableLoadBalancerNodeSet=true,ValidatingAdmissionPolicy=true,ZeroLimitedNominalConcurrencyShares=true"
