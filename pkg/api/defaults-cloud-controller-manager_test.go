// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
)

func TestCloudControllerManagerConfig(t *testing.T) {
	k8sVersion := common.RationalizeReleaseAndVersion(Kubernetes, "", "", false, false, true)
	cs := CreateMockContainerService("testcluster", k8sVersion, 3, 2, false)
	cs.setCloudControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	if cm["--controllers"] != "*,-cloud-node" {
		t.Fatalf("got unexpected '--controllers' Cloud Controller Manager config value for Kubernetes %s: %s",
			k8sVersion, cm["--controllers"])
	}
}

func TestCloudControllerManagerFeatureGates(t *testing.T) {
	// test defaultTestClusterVer
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setCloudControllerManagerConfig()
	ccm := cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	if ccm["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' Cloud Controller Manager config value for k8s v%s: %s",
			defaultTestClusterVer, ccm["--feature-gates"])
	}

	// test 1.19.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.19.0"
	cs.setCloudControllerManagerConfig()
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	if ccm["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' Cloud Controller Manager config value for k8s v%s: %s",
			"1.19.0", ccm["--feature-gates"])
	}

	// test 1.22.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.22.0"
	cs.setCloudControllerManagerConfig()
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	if ccm["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' Cloud Controller Manager config value for k8s v%s: %s",
			"1.22.0", ccm["--feature-gates"])
	}

	// test user-overrides, removal of VolumeSnapshotDataSource for k8s versions >= 1.22
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.22.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	ccm["--feature-gates"] = "VolumeSnapshotDataSource=true"
	cs.setCloudControllerManagerConfig()
	if ccm["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for \"--feature-gates\": \"VolumeSnapshotDataSource=true\": %s for k8s v%s",
			ccm["--feature-gates"], "1.22.0")
	}

	// test user-overrides, no removal of VolumeSnapshotDataSource for k8s versions < 1.22
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.19.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	ccm["--feature-gates"] = "VolumeSnapshotDataSource=true"
	cs.setCloudControllerManagerConfig()
	if ccm["--feature-gates"] != "VolumeSnapshotDataSource=true" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for \"--feature-gates\": \"VolumeSnapshotDataSource=true\": %s for k8s v%s",
			ccm["--feature-gates"], "1.19.0")
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.27.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	ccm["--feature-gates"] = "ControllerManagerLeaderMigration=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,IdentifyPodOS=true,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,StatefulSetMinReadySeconds=true"
	cs.setCloudControllerManagerConfig()
	if ccm["--feature-gates"] != "" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for \"--feature-gates\": %s for k8s v%s",
			ccm["--feature-gates"], "1.27.0")
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.26.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	ccm["--feature-gates"] = "ControllerManagerLeaderMigration=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,IdentifyPodOS=true,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,StatefulSetMinReadySeconds=true"
	cs.setCloudControllerManagerConfig()
	if ccm["--feature-gates"] != "CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,ControllerManagerLeaderMigration=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,IdentifyPodOS=true,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,StatefulSetMinReadySeconds=true" {
		t.Fatalf("got unexpected '--feature-gates' API server config value for \"--feature-gates\": %s for k8s v%s",
			ccm["--feature-gates"], "1.26.0")
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.28
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.28.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	featuregate128 := "AdvancedAuditing=true,CSIMigrationGCE=true,CSIStorageCapacity=true,DelegateFSGroupToCSIDriver=true,DevicePlugins=true,DisableAcceleratorUsageMetrics=true,DryRun=true,EndpointSliceTerminatingCondition=true,KubeletCredentialProviders=true,MixedProtocolLBService=true,NetworkPolicyStatus=true,PodHasNetworkCondition=true,PodSecurity=true,ServiceIPStaticSubrange=true,ServiceInternalTrafficPolicy=true,UserNamespacesStatelessPodsSupport=true,WindowsHostProcessContainers=true"
	ccm["--feature-gates"] = featuregate128
	featuregate128Sanitized := ""
	cs.setCloudControllerManagerConfig()
	if ccm["--feature-gates"] != featuregate128Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n controller manager config original value  %s \n, expected sanitized value: %s \n, actual sanitized value: %s \n ",
			"1.28.0", featuregate128, ccm["--feature-gates"], featuregate128Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.27.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	ccm["--feature-gates"] = featuregate128
	featuregate127Sanitized := featuregate128
	cs.setCloudControllerManagerConfig()
	if ccm["--feature-gates"] != featuregate127Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n controller manager config original value  %s \n, expected sanitized value: %s \n, actual sanitized value: %s \n ",
			"1.27.0", featuregate128, ccm["--feature-gates"], featuregate127Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.29
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.29.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	featuregate129 := "CSIMigrationvSphere=true,CronJobTimeZone=true,DownwardAPIHugePages=true,GRPCContainerProbe=true,JobMutableNodeSchedulingDirectives=true,JobTrackingWithFinalizers=true,LegacyServiceAccountTokenNoAutoGeneration=true,OpenAPIV3=true,ProbeTerminationGracePeriod=true,RetroactiveDefaultStorageClass=true,SeccompDefault=true,TopologyManager=true"
	ccm["--feature-gates"] = featuregate129
	featuregate129Sanitized := ""
	cs.setCloudControllerManagerConfig()
	if ccm["--feature-gates"] != featuregate129Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n Cloud Controller Manager config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.29.0", featuregate129, ccm["--feature-gates"], featuregate129Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.29
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.28.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	ccm["--feature-gates"] = featuregate129
	featuregate128Sanitized = featuregate129
	cs.setCloudControllerManagerConfig()
	if ccm["--feature-gates"] != featuregate128Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n Cloud Controller Manager config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.28.0", featuregate129, ccm["--feature-gates"], featuregate128Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.30
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.30.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	featuregate130 := "APISelfSubjectReview=true,CSIMigrationAzureFile=true,ExpandedDNSConfig=true,ExperimentalHostUserNamespaceDefaulting=true,IPTablesOwnershipCleanup=true,KubeletPodResources=true,KubeletPodResourcesGetAllocatable=true,LegacyServiceAccountTokenTracking=true,MinimizeIPTablesRestore=true,ProxyTerminatingEndpoints=true,RemoveSelfLink=true,SecurityContextDeny=true"
	ccm["--feature-gates"] = featuregate130
	featuregate130Sanitized := ""
	cs.setAPIServerConfig()
	if ccm["--feature-gates"] != featuregate130Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.30.0", featuregate130, ccm["--feature-gates"], featuregate130Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.30
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.29.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	ccm["--feature-gates"] = featuregate130
	featuregate129Sanitized = featuregate130
	cs.setAPIServerConfig()
	if ccm["--feature-gates"] != featuregate129Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.29.0", featuregate130, ccm["--feature-gates"], featuregate129Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.31
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.31.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	featuregate131 := "CloudDualStackNodeIPs=true,DRAControlPlaneController=true,HPAContainerMetrics=true,KMSv2=true,KMSv2KDF=true,LegacyServiceAccountTokenCleanUp=true,MinDomainsInPodTopologySpread=true,NewVolumeManagerReconstruction=true,NodeOutOfServiceVolumeDetach=true,PodHostIPs=true,ServerSideApply=true,ServerSideFieldValidation=true,StableLoadBalancerNodeSet=true,ValidatingAdmissionPolicy=true,ZeroLimitedNominalConcurrencyShares=true"
	ccm["--feature-gates"] = featuregate131
	featuregate131Sanitized := ""
	cs.setCloudControllerManagerConfig()
	if ccm["--feature-gates"] != featuregate131Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n Cloud Controller Manager config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.31.0", featuregate131, ccm["--feature-gates"], featuregate131Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.31
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.30.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	ccm["--feature-gates"] = featuregate131
	featuregate130Sanitized = featuregate131
	cs.setCloudControllerManagerConfig()
	if ccm["--feature-gates"] != featuregate130Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n Cloud Controller Manager config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.30.0", featuregate131, ccm["--feature-gates"], featuregate130Sanitized)
	}
}
