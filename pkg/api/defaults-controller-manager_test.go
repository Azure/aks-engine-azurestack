// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"strings"
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/helpers/to"
)

func TestControllerManagerConfigEnableRbac(t *testing.T) {
	// Test EnableRbac = true
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableRbac = to.BoolPtr(true)
	cs.setControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--use-service-account-credentials"] != "true" {
		t.Fatalf("got unexpected '--use-service-account-credentials' Controller Manager config value for EnableRbac=true: %s",
			cm["--use-service-account-credentials"])
	}

	// Test default
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableRbac = to.BoolPtr(false)
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--use-service-account-credentials"] != DefaultKubernetesCtrlMgrUseSvcAccountCreds {
		t.Fatalf("got unexpected '--use-service-account-credentials' Controller Manager config value for EnableRbac=false: %s",
			cm["--use-service-account-credentials"])
	}
}

func TestControllerManagerConfigCloudProvider(t *testing.T) {
	// Test UseCloudControllerManager = true
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager = to.BoolPtr(true)
	cs.setControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--cloud-provider"] != "external" {
		t.Fatalf("got unexpected '--cloud-provider' Controller Manager config value for UseCloudControllerManager=true: %s",
			cm["--cloud-provider"])
	}

	// Test UseCloudControllerManager = false
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager = to.BoolPtr(false)
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--cloud-provider"] != "azure" {
		t.Fatalf("got unexpected '--cloud-provider' Controller Manager config value for UseCloudControllerManager=false: %s",
			cm["--cloud-provider"])
	}
}

func TestControllerManagerConfigEnableProfiling(t *testing.T) {
	// Test
	// "controllerManagerConfig": {
	// 	"--profiling": "true"
	// },
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = map[string]string{
		"--profiling": "true",
	}
	cs.setControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--profiling"] != "true" {
		t.Fatalf("got unexpected '--profiling' Controller Manager config value for \"--profiling\": \"true\": %s",
			cm["--profiling"])
	}

	// Test default
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--profiling"] != DefaultKubernetesCtrMgrEnableProfiling {
		t.Fatalf("got unexpected default value for '--profiling' Controller Manager config: %s",
			cm["--profiling"])
	}
}

func TestControllerManagerConfigFeatureGates(t *testing.T) {
	// test defaultTestClusterVer
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"\": %s",
			cm["--feature-gates"])
	}

	// test 1.19.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.19.0"
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--feature-gates"] != "LocalStorageCapacityIsolation=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LocalStorageCapacityIsolation=true\": %s",
			cm["--feature-gates"])
	}

	// test 1.22.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.22.0"
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--feature-gates"] != "LocalStorageCapacityIsolation=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LocalStorageCapacityIsolation=true\": %s",
			cm["--feature-gates"])
	}

	// test 1.24.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.24.0"
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--feature-gates"] != "LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true\": %s",
			cm["--feature-gates"])
	}

	// test 1.25.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.25.0"
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--feature-gates"] != "LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true,PodSecurity=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true,PodSecurity=true\": %s",
			cm["--feature-gates"])
	}

	// test 1.26.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.26.0"
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--feature-gates"] != "LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true,PodSecurity=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true,PodSecurity=true\": %s",
			cm["--feature-gates"])
	}

	// test user-overrides
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	cm["--feature-gates"] = "TaintBasedEvictions=true"
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != "TaintBasedEvictions=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"TaintBasedEvictions=true\": %s",
			cm["--feature-gates"])
	}

	// test user-overrides, removal of VolumeSnapshotDataSource for k8s versions >= 1.22
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.22.0"
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	cm["--feature-gates"] = "VolumeSnapshotDataSource=true"
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != "LocalStorageCapacityIsolation=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LocalStorageCapacityIsolation=true\": %s",
			cm["--feature-gates"])
	}

	// test user-overrides, no removal of VolumeSnapshotDataSource for k8s versions < 1.22
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.19.0"
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	cm["--feature-gates"] = "VolumeSnapshotDataSource=true"
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != "LocalStorageCapacityIsolation=true,VolumeSnapshotDataSource=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LocalStorageCapacityIsolation=true,VolumeSnapshotDataSource=true\": %s",
			cm["--feature-gates"])
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.27.0"
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	cm["--feature-gates"] = "ControllerManagerLeaderMigration=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,IdentifyPodOS=true,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,StatefulSetMinReadySeconds=true,LegacyServiceAccountTokenNoAutoGeneration=false"
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != "PodSecurity=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": %s",
			cm["--feature-gates"])
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.26.0"
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	cm["--feature-gates"] = "ControllerManagerLeaderMigration=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,IdentifyPodOS=true,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,StatefulSetMinReadySeconds=true"
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != "CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,ControllerManagerLeaderMigration=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,IdentifyPodOS=true,LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,PodSecurity=true,StatefulSetMinReadySeconds=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": %s",
			cm["--feature-gates"])
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.28
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.28.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = make(map[string]string)
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	featuregate128 := "AdvancedAuditing=true,CSIMigrationGCE=true,CSIStorageCapacity=true,DelegateFSGroupToCSIDriver=true,DevicePlugins=true,DisableAcceleratorUsageMetrics=true,DryRun=true,EndpointSliceTerminatingCondition=true,KubeletCredentialProviders=true,MixedProtocolLBService=true,NetworkPolicyStatus=true,PodHasNetworkCondition=true,PodSecurity=true,ServiceIPStaticSubrange=true,ServiceInternalTrafficPolicy=true,UserNamespacesStatelessPodsSupport=true,WindowsHostProcessContainers=true"
	cm["--feature-gates"] = featuregate128
	featuregate128Sanitized := ""
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != featuregate128Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n controller manager config original value  %s \n, expected sanitized value: %s \n, actual sanitized value: %s \n ",
			"1.28.0", featuregate128, cm["--feature-gates"], featuregate128Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.27.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = make(map[string]string)
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	cm["--feature-gates"] = featuregate128
	featuregate127Sanitized := featuregate128
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != featuregate127Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n controller manager config original value  %s \n, expected sanitized value: %s \n, actual sanitized value: %s \n ",
			"1.27.0", featuregate128, cm["--feature-gates"], featuregate127Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.29
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.29.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = make(map[string]string)
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	featuregate129 := "CSIMigrationvSphere=true,CronJobTimeZone=true,DownwardAPIHugePages=true,GRPCContainerProbe=true,JobMutableNodeSchedulingDirectives=true,JobTrackingWithFinalizers=true,LegacyServiceAccountTokenNoAutoGeneration=true,OpenAPIV3=true,ProbeTerminationGracePeriod=true,RetroactiveDefaultStorageClass=true,SeccompDefault=true,TopologyManager=true"
	cm["--feature-gates"] = featuregate129
	featuregate129Sanitized := ""
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != featuregate129Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n controller manager config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.29.0", featuregate129, cm["--feature-gates"], featuregate129Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.29
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.28.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = make(map[string]string)
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	cm["--feature-gates"] = featuregate129
	featuregate128Sanitized = "CSIMigrationvSphere=true,CronJobTimeZone=true,DownwardAPIHugePages=true,GRPCContainerProbe=true,JobMutableNodeSchedulingDirectives=true,JobTrackingWithFinalizers=true,OpenAPIV3=true,ProbeTerminationGracePeriod=true,RetroactiveDefaultStorageClass=true,SeccompDefault=true,TopologyManager=true"
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != featuregate128Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n controller manager config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.28.0", featuregate129, cm["--feature-gates"], featuregate128Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.30
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.30.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	featuregate130 := "APISelfSubjectReview=true,CSIMigrationAzureFile=true,ExpandedDNSConfig=true,ExperimentalHostUserNamespaceDefaulting=true,IPTablesOwnershipCleanup=true,KubeletPodResources=true,KubeletPodResourcesGetAllocatable=true,LegacyServiceAccountTokenTracking=true,MinimizeIPTablesRestore=true,ProxyTerminatingEndpoints=true,RemoveSelfLink=true,SecurityContextDeny=true"
	cm["--feature-gates"] = featuregate130
	featuregate130Sanitized := ""
	cs.setAPIServerConfig()
	if cm["--feature-gates"] != featuregate130Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.30.0", featuregate130, cm["--feature-gates"], featuregate130Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.30
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.29.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	cm["--feature-gates"] = featuregate130
	featuregate129Sanitized = featuregate130
	cs.setAPIServerConfig()
	if cm["--feature-gates"] != featuregate129Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.29.0", featuregate130, cm["--feature-gates"], featuregate129Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.31
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.31.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = make(map[string]string)
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	featuregate131 := "APIPriorityAndFairness=true,CSIMigrationRBD=true,CSINodeExpandSecret=true,ConsistentHTTPGetHandlers=true,CustomResourceValidationExpressions=true,DefaultHostNetworkHostPortsInPodTemplates=true,InTreePluginRBDUnregister=true,JobReadyPods=true,ReadWriteOncePod=true,ServiceNodePortStaticSubrange=true,SkipReadOnlyValidationGCE=true"
	cm["--feature-gates"] = featuregate131
	featuregate131Sanitized := ""
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != featuregate131Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n controller manager config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.31", featuregate131, cm["--feature-gates"], featuregate131Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.31
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.30.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = make(map[string]string)
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	cm["--feature-gates"] = featuregate131
	featuregate130Sanitized = featuregate131
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != featuregate130Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n controller manager config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.30.0", featuregate131, cm["--feature-gates"], featuregate130Sanitized)
	}
}

func TestControllerManagerDefaultConfig(t *testing.T) {
	// Azure defaults
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--node-monitor-grace-period"] != string(DefaultKubernetesCtrlMgrNodeMonitorGracePeriod) {
		t.Fatalf("expected controller-manager to have node-monitor-grace-period set to its default value")
	}
	// --pod-eviction-timeout is removed after 1.27

	if cm["--route-reconciliation-period"] != string(DefaultKubernetesCtrlMgrRouteReconciliationPeriod) {
		t.Fatalf("expected controller-manager to have route-reconciliation-period set to its default value")
	}
	if cm["--bind-address"] != "127.0.0.1" {
		t.Fatalf("expected controller-manager to have route-reconciliation-period set to its default value")
	}
	if cm["--tls-min-version"] != "VersionTLS12" {
		t.Fatalf("expected controller-manager to have route-reconciliation-period set to its default value")
	}

	// Azure Stack defaults
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.CustomCloudProfile = &CustomCloudProfile{}
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--node-monitor-grace-period"] != string(DefaultAzureStackKubernetesCtrlMgrNodeMonitorGracePeriod) {
		t.Fatalf("expected controller-manager to have node-monitor-grace-period set to its default value")
	}
	// --pod-eviction-timeout is removed after 1.27

	if cm["--route-reconciliation-period"] != string(DefaultAzureStackKubernetesCtrlMgrRouteReconciliationPeriod) {
		t.Fatalf("expected controller-manager to have route-reconciliation-period set to its default value")
	}
}

func TestControllerManagerInsecureFlag(t *testing.T) {
	type controllerManagerTest struct {
		version string
		found   bool
	}

	controllerManagerTestsForceDelete := []controllerManagerTest{
		{
			version: "1.23.0",
			found:   true,
		},
		{
			version: "1.24.0",
			found:   false,
		},
	}

	for _, tt := range controllerManagerTestsForceDelete {
		cs := CreateMockContainerService("testcluster", tt.version, 3, 2, false)
		cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = map[string]string{
			"--address": "0.0.0.0",
			"--port":    "443",
		}
		cs.setControllerManagerConfig()
		a := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig

		_, found := a["--address"]
		if found != tt.found {
			t.Fatalf("got --address found %t want %t", found, tt.found)
		}
		_, found = a["--port"]
		if found != tt.found {
			t.Fatalf("got --port found %t want %t", found, tt.found)
		}
	}

}

func TestControllerManagerEnableTaintManagerFlag(t *testing.T) {
	type controllerManagerTest struct {
		version string
		found   bool
	}

	controllerManagerTestsForceDelete := []controllerManagerTest{
		{
			version: "1.26.0",
			found:   true,
		},
		{
			version: "1.27.0",
			found:   false,
		},
	}

	for _, tt := range controllerManagerTestsForceDelete {
		cs := CreateMockContainerService("testcluster", tt.version, 3, 2, false)
		cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = map[string]string{
			"--enable-taint-manager": "true",
			"--pod-eviction-timeout": "5m0s",
		}
		cs.setControllerManagerConfig()
		a := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig

		_, found := a["--enable-taint-manager"]
		if found != tt.found {
			t.Fatalf("got --enable-taint-manager found %t want %t", found, tt.found)
		}
		_, found = a["--pod-eviction-timeout"]
		if found != tt.found {
			t.Fatalf("got --pod-eviction-timeout found %t want %t", found, tt.found)
		}
	}

}

func TestControllerManagerFeatureGates132(t *testing.T) {
	// test user-overrides, removal of feature gates for k8s versions >= 1.32
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.32.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = make(map[string]string)
	a := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	featuregate132 := "GATE01=true,GATE02=true"
	a["--feature-gates"] = featuregate132
	cs.setControllerManagerConfig()
	// split both strings by "," and ensure no original item exists in the sanitized list
	originalList := strings.Split(featuregate132, ",")
	sanitizedList := strings.Split(a["--feature-gates"], ",")
	for _, of := range originalList {
		for _, sf := range sanitizedList {
			if of == sf {
				t.Fatalf("feature-gate %q should not exist in sanitized list for %s\nfeaturegate132 (original): %q\nfeaturegate132Sanitized (actual): %q", sf, "1.32.0", featuregate132, a["--feature-gates"])
			}
		}
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.32
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.31.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = make(map[string]string)
	a = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	a["--feature-gates"] = featuregate132
	cs.setControllerManagerConfig()
	actualList := strings.Split(a["--feature-gates"], ",")
	expectedList := strings.Split(featuregate132, ",")
	for _, exp := range expectedList {
		found := false
		for _, act := range actualList {
			if act == exp {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing feature-gate %q in actual '--feature-gates' for %s\nfeaturegate131 (expected subset): %q\nactual: %q",
				exp, "1.31.0", featuregate132, a["--feature-gates"])
		}
	}
}
