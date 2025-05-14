// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"testing"
)

func TestSchedulerDefaultConfig(t *testing.T) {
	cs := CreateMockContainerService("testcluster", "", 3, 2, false)
	cs.setSchedulerConfig()
	s := cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	for key, val := range staticSchedulerConfig {
		if val != s[key] {
			t.Fatalf("got unexpected kube-scheduler static config value for %s. Expected %s, got %s",
				key, val, s[key])
		}
	}
	for key, val := range defaultSchedulerConfig {
		if val != s[key] {
			t.Fatalf("got unexpected kube-scheduler default config value for %s. Expected %s, got %s",
				key, val, s[key])
		}
	}
}

func TestSchedulerUserConfig(t *testing.T) {
	cs := CreateMockContainerService("testcluster", "", 3, 2, true)
	assignmentMap := map[string]string{
		"--scheduler-name": "my-custom-name",
		"--feature-gates":  "APIListChunking=true,APIResponseCompression=true,Accelerators=true,AdvancedAuditing=true",
	}
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = assignmentMap
	cs.setSchedulerConfig()
	for key, val := range assignmentMap {
		if val != cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig[key] {
			t.Fatalf("got unexpected kube-scheduler config value for %s. Expected %s, got %s",
				key, val, cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig[key])
		}
	}
}

func TestSchedulerFlagReplacement(t *testing.T) {
	cs := CreateMockContainerService("testcluster", "", 3, 2, true)

	// Verify that the flag is replaced after 1.28
	flagChange := map[string]string{}
	// The deprecated flag --lock-object-namespace and --lock-object-name have been removed from kube-scheduler.
	// Please use --leader-elect-resource-namespace and --leader-elect-resource-name or ComponentConfig instead to configure those parameters. (#119130, @SataQiu) [SIG Scheduling]
	flagChange["--lock-object-namespace"] = "--leader-elect-resource-namespace"
	flagChange["--lock-object-name"] = "--leader-elect-resource-name"

	assignmentMap := map[string]string{
		"--lock-object-namespace": "system-lock",
		"--lock-object-name":      "leader",
		"--scheduler-name":        "my-custom-name",
	}

	assignment128MapSanitized := map[string]string{
		"--leader-elect-resource-namespace": "system-lock",
		"--leader-elect-resource-name":      "leader",
		"--scheduler-name":                  "my-custom-name",
	}
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.28.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = assignmentMap
	cs.setSchedulerConfig()
	for key, val := range assignment128MapSanitized {
		if val != cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig[key] {
			t.Fatalf("got unexpected kube-scheduler config value for %s. Expected %s, got %s",
				key, val, cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig[key])
		}
	}

	for key, val := range flagChange {
		if _, ok := cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig[key]; ok {
			t.Fatalf("got unexpected kube-scheduler config value. The flag should be replaced from %s to %s",
				key, val)
		}
	}

	// Verify that the flag is not replaced before 1.28
	cs = CreateMockContainerService("testcluster", "", 3, 2, true)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.27.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = assignmentMap
	cs.setSchedulerConfig()
	for key, val := range assignmentMap {
		if val != cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig[key] {
			t.Fatalf("got unexpected kube-scheduler config value for %s. Expected %s, got %s",
				key, val, cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig[key])
		}
	}
}

func TestSchedulerStaticConfig(t *testing.T) {
	cs := CreateMockContainerService("testcluster", "", 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = map[string]string{
		"--kubeconfig":      "user-override",
		"--leader-elect":    "user-override",
		"--profiling":       "user-override",
		"--bind-address":    "user-override",
		"--tls-min-version": "user-override",
	}
	cs.setSchedulerConfig()
	for key, val := range staticSchedulerConfig {
		if val != cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig[key] {
			t.Fatalf("kube-scheduler static config did not override user values for %s. Expected %s, got %s",
				key, val, cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig)
		}
	}
}

func TestSchedulerConfigEnableProfiling(t *testing.T) {
	// Test
	// "schedulerConfig": {
	// 	"--profiling": "true"
	// },
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = map[string]string{
		"--profiling": "true",
	}
	cs.setSchedulerConfig()
	s := cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	if s["--profiling"] != "true" {
		t.Fatalf("got unexpected '--profiling' Scheduler config value for \"--profiling\": \"true\": %s",
			s["--profiling"])
	}

	// Test default
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setSchedulerConfig()
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	if s["--profiling"] != DefaultKubernetesSchedulerEnableProfiling {
		t.Fatalf("got unexpected default value for '--profiling' Scheduler config: %s",
			s["--profiling"])
	}
}

func TestSchedulerFeatureGates(t *testing.T) {
	// test defaultTestClusterVer
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setSchedulerConfig()
	s := cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	if s["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' Scheduler config value for k8s v%s: %s",
			defaultTestClusterVer, s["--feature-gates"])
	}

	// test 1.19.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.19.0"
	cs.setSchedulerConfig()
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	if s["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' Scheduler config value for k8s v%s: %s",
			"1.19.0", s["--feature-gates"])
	}

	// test 1.22.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.22.0"
	cs.setSchedulerConfig()
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	if s["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' Scheduler config value for k8s v%s: %s",
			"1.22.0", s["--feature-gates"])
	}

	// test 1.25.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.25.0"
	cs.setSchedulerConfig()
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	if s["--feature-gates"] != "PodSecurity=true" {
		t.Fatalf("got unexpected '--feature-gates' Scheduler config value for k8s v%s: %s",
			"1.25.0", s["--feature-gates"])
	}

	// test user-overrides, removal of VolumeSnapshotDataSource for k8s versions >= 1.22
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.22.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = make(map[string]string)
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	s["--feature-gates"] = "VolumeSnapshotDataSource=true"
	cs.setSchedulerConfig()
	if s["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' Scheduler config value for \"--feature-gates\": \"VolumeSnapshotDataSource=true\": %s for k8s v%s",
			s["--feature-gates"], "1.22.0")
	}

	// test user-overrides, no removal of VolumeSnapshotDataSource for k8s versions < 1.22
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.19.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = make(map[string]string)
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	s["--feature-gates"] = "VolumeSnapshotDataSource=true"
	cs.setSchedulerConfig()
	if s["--feature-gates"] != "VolumeSnapshotDataSource=true" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for \"--feature-gates\": \"VolumeSnapshotDataSource=true\": %s for k8s v%s",
			s["--feature-gates"], "1.19.0")
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.27.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = make(map[string]string)
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	s["--feature-gates"] = "ControllerManagerLeaderMigration=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,IdentifyPodOS=true,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,StatefulSetMinReadySeconds=true"
	cs.setSchedulerConfig()
	if s["--feature-gates"] != "PodSecurity=true" {
		t.Fatalf("got unexpected '--feature-gates' Scheduler config value for \"--feature-gates\": %s for k8s v%s",
			s["--feature-gates"], "1.27.0")
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.26.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = make(map[string]string)
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	s["--feature-gates"] = "ControllerManagerLeaderMigration=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,IdentifyPodOS=true,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,StatefulSetMinReadySeconds=true"
	cs.setSchedulerConfig()
	if s["--feature-gates"] != "CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,ControllerManagerLeaderMigration=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,IdentifyPodOS=true,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,PodSecurity=true,StatefulSetMinReadySeconds=true" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for \"--feature-gates\": %s for k8s v%s",
			s["--feature-gates"], "1.26.0")
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.28
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.28.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = make(map[string]string)
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	featuregate128 := "AdvancedAuditing=true,CSIMigrationGCE=true,CSIStorageCapacity=true,DelegateFSGroupToCSIDriver=true,DevicePlugins=true,DisableAcceleratorUsageMetrics=true,DryRun=true,EndpointSliceTerminatingCondition=true,KubeletCredentialProviders=true,MixedProtocolLBService=true,NetworkPolicyStatus=true,PodHasNetworkCondition=true,PodSecurity=true,ServiceIPStaticSubrange=true,ServiceInternalTrafficPolicy=true,UserNamespacesStatelessPodsSupport=true,WindowsHostProcessContainers=true"
	s["--feature-gates"] = featuregate128
	featuregate128Sanitized := ""
	cs.setSchedulerConfig()
	if s["--feature-gates"] != featuregate128Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n kubelet config original value  %s \n, expected sanitized value: %s \n, actual sanitized value: %s \n ",
			"1.28.0", featuregate128, s["--feature-gates"], featuregate128Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.27.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = make(map[string]string)
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	s["--feature-gates"] = featuregate128
	featuregate127Sanitized := featuregate128
	cs.setSchedulerConfig()
	if s["--feature-gates"] != featuregate127Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n kubelet config original value  %s \n, expected sanitized value: %s \n, actual sanitized value: %s \n ",
			"1.27.0", featuregate128, s["--feature-gates"], featuregate127Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.29
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.29.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = make(map[string]string)
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	featuregate129 := "CSIMigrationvSphere=true,CronJobTimeZone=true,DownwardAPIHugePages=true,GRPCContainerProbe=true,JobMutableNodeSchedulingDirectives=true,JobTrackingWithFinalizers=true,LegacyServiceAccountTokenNoAutoGeneration=true,OpenAPIV3=true,ProbeTerminationGracePeriod=true,RetroactiveDefaultStorageClass=true,SeccompDefault=true,TopologyManager=true"
	s["--feature-gates"] = featuregate129
	featuregate129Sanitized := ""
	cs.setSchedulerConfig()
	if s["--feature-gates"] != featuregate129Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n Scheduler config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.29.0", featuregate129, s["--feature-gates"], featuregate129Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.29
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.28.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = make(map[string]string)
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	s["--feature-gates"] = featuregate129
	featuregate128Sanitized = featuregate129
	cs.setSchedulerConfig()
	if s["--feature-gates"] != featuregate128Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n Scheduler config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.28.0", featuregate129, s["--feature-gates"], featuregate128Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.30
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.30.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	featuregate130 := "APISelfSubjectReview=true,CSIMigrationAzureFile=true,ExpandedDNSConfig=true,ExperimentalHostUserNamespaceDefaulting=true,IPTablesOwnershipCleanup=true,KubeletPodResources=true,KubeletPodResourcesGetAllocatable=true,LegacyServiceAccountTokenTracking=true,MinimizeIPTablesRestore=true,ProxyTerminatingEndpoints=true,RemoveSelfLink=true,SecurityContextDeny=true"
	s["--feature-gates"] = featuregate130
	featuregate130Sanitized := ""
	cs.setAPIServerConfig()
	if s["--feature-gates"] != featuregate130Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.30.0", featuregate130, s["--feature-gates"], featuregate130Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.30
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.29.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	s["--feature-gates"] = featuregate130
	featuregate129Sanitized = featuregate130
	cs.setAPIServerConfig()
	if s["--feature-gates"] != featuregate129Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.29.0", featuregate130, s["--feature-gates"], featuregate129Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.31
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.31.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = make(map[string]string)
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	featuregate131 := "CloudDualStackNodeIPs=true,DRAControlPlaneController=true,HPAContainerMetrics=true,KMSv2=true,KMSv2KDF=true,LegacyServiceAccountTokenCleanUp=true,MinDomainsInPodTopologySpread=true,NewVolumeManagerReconstruction=true,NodeOutOfServiceVolumeDetach=true,PodHostIPs=true,ServerSideApply=true,ServerSideFieldValidation=true,StableLoadBalancerNodeSet=true,ValidatingAdmissionPolicy=true,ZeroLimitedNominalConcurrencyShares=true"
	s["--feature-gates"] = featuregate131
	featuregate131Sanitized := ""
	cs.setSchedulerConfig()
	if s["--feature-gates"] != featuregate131Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n Scheduler config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.31.0", featuregate131, s["--feature-gates"], featuregate131Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.31
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.30.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = make(map[string]string)
	s = cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	s["--feature-gates"] = featuregate131
	featuregate130Sanitized = featuregate131
	cs.setSchedulerConfig()
	if s["--feature-gates"] != featuregate130Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n Scheduler config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.30.0", featuregate131, s["--feature-gates"], featuregate130Sanitized)
	}
}
