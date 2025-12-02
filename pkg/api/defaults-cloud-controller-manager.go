// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"strconv"
	"strings"

	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
)

func (cs *ContainerService) setCloudControllerManagerConfig() {
	o := cs.Properties.OrchestratorProfile
	isAzureCNIDualStack := cs.Properties.IsAzureCNIDualStack()
	clusterCidr := o.KubernetesConfig.ClusterSubnet
	if isAzureCNIDualStack {
		clusterSubnets := strings.Split(clusterCidr, ",")
		if len(clusterSubnets) > 1 {
			clusterCidr = clusterSubnets[1]
		}
	}
	staticCloudControllerManagerConfig := map[string]string{
		"--allocate-node-cidrs":         strconv.FormatBool(!o.IsAzureCNI() || isAzureCNIDualStack),
		"--configure-cloud-routes":      strconv.FormatBool(cs.Properties.RequireRouteTable()),
		"--cloud-provider":              "azure",
		"--cloud-config":                "/etc/kubernetes/azure.json",
		"--cluster-cidr":                clusterCidr,
		"--kubeconfig":                  "/var/lib/kubelet/kubeconfig",
		"--leader-elect":                "true",
		"--route-reconciliation-period": "10s",
		"--v":                           "2",
	}

	// Disable cloud-node controller
	staticCloudControllerManagerConfig["--controllers"] = "*,-cloud-node"

	// Set --cluster-name based on appropriate DNS prefix
	if cs.Properties.MasterProfile != nil {
		staticCloudControllerManagerConfig["--cluster-name"] = cs.Properties.MasterProfile.DNSPrefix
	}

	// Default cloud-controller-manager config
	defaultCloudControllerManagerConfig := map[string]string{
		"--route-reconciliation-period": DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
	}

	// If no user-configurable cloud-controller-manager config values exists, use the defaults
	if o.KubernetesConfig.CloudControllerManagerConfig == nil {
		o.KubernetesConfig.CloudControllerManagerConfig = defaultCloudControllerManagerConfig
	} else {
		for key, val := range defaultCloudControllerManagerConfig {
			// If we don't have a user-configurable cloud-controller-manager config for each option
			if _, ok := o.KubernetesConfig.CloudControllerManagerConfig[key]; !ok {
				// then assign the default value
				o.KubernetesConfig.CloudControllerManagerConfig[key] = val
			}
		}
	}

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	for key, val := range staticCloudControllerManagerConfig {
		o.KubernetesConfig.CloudControllerManagerConfig[key] = val
	}

	invalidFeatureGates := []string{}
	// Remove --feature-gate VolumeSnapshotDataSource starting with 1.22
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.22.0-alpha.1") {
		invalidFeatureGates = append(invalidFeatureGates, "VolumeSnapshotDataSource")
	}
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.27.0") {
		// Remove --feature-gate ControllerManagerLeaderMigration starting with 1.27
		// Reference: https://github.com/kubernetes/kubernetes/pull/113534
		invalidFeatureGates = append(invalidFeatureGates, "ControllerManagerLeaderMigration")
		// Remove --feature-gate ExpandCSIVolumes, ExpandInUsePersistentVolumes, ExpandPersistentVolumes starting with 1.27
		// Reference: https://github.com/kubernetes/kubernetes/pull/113942
		invalidFeatureGates = append(invalidFeatureGates, "ExpandCSIVolumes", "ExpandInUsePersistentVolumes", "ExpandPersistentVolumes")
		// Remove --feature-gate CSIInlineVolume, CSIMigration, CSIMigrationAzureDisk, DaemonSetUpdateSurge, EphemeralContainers, IdentifyPodOS, LocalStorageCapacityIsolation, NetworkPolicyEndPort, StatefulSetMinReadySeconds starting with 1.27
		// Reference: https://github.com/kubernetes/kubernetes/pull/114410
		invalidFeatureGates = append(invalidFeatureGates, "CSIInlineVolume", "CSIMigration", "CSIMigrationAzureDisk", "DaemonSetUpdateSurge", "EphemeralContainers", "IdentifyPodOS", "LocalStorageCapacityIsolation", "NetworkPolicyEndPort", "StatefulSetMinReadySeconds")
	}

	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.28.0") {
		// Remove --feature-gate AdvancedAuditing starting with 1.28
		invalidFeatureGates = append(invalidFeatureGates, "AdvancedAuditing", "DisableAcceleratorUsageMetrics", "DryRun", "PodSecurity")

		invalidFeatureGates = append(invalidFeatureGates, "NetworkPolicyStatus", "PodHasNetworkCondition", "UserNamespacesStatelessPodsSupport")

		// Remove --feature-gate CSIMigrationGCE starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/117055
		invalidFeatureGates = append(invalidFeatureGates, "CSIMigrationGCE")

		// Remove --feature-gate CSIStorageCapacity starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/118018
		invalidFeatureGates = append(invalidFeatureGates, "CSIStorageCapacity")

		// Remove --feature-gate DelegateFSGroupToCSIDriver starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/117655
		invalidFeatureGates = append(invalidFeatureGates, "DelegateFSGroupToCSIDriver")

		// Remove --feature-gate DevicePlugins starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/117656
		invalidFeatureGates = append(invalidFeatureGates, "DevicePlugins")

		// Remove --feature-gate KubeletCredentialProviders starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/116901
		invalidFeatureGates = append(invalidFeatureGates, "KubeletCredentialProviders")

		// Remove --feature-gate MixedProtocolLBService, ServiceInternalTrafficPolicy, ServiceIPStaticSubrange, EndpointSliceTerminatingCondition  starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/117237
		invalidFeatureGates = append(invalidFeatureGates, "MixedProtocolLBService", "ServiceInternalTrafficPolicy", "ServiceIPStaticSubrange", "EndpointSliceTerminatingCondition")

		// Remove --feature-gate WindowsHostProcessContainers starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/117570
		invalidFeatureGates = append(invalidFeatureGates, "WindowsHostProcessContainers")
	}
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.29.0") {
		// Remove --feature-gate CSIMigrationvSphere starting with 1.29
		// Reference: https://github.com/kubernetes/kubernetes/pull/121291
		invalidFeatureGates = append(invalidFeatureGates, "CSIMigrationvSphere")

		// Remove --feature-gate ProbeTerminationGracePeriod starting with 1.29
		// Reference: https://github.com/kubernetes/kubernetes/pull/121257
		invalidFeatureGates = append(invalidFeatureGates, "ProbeTerminationGracePeriod")

		// Remove --feature-gate JobTrackingWithFinalizers starting with 1.29
		// Reference: https://github.com/kubernetes/kubernetes/pull/119100
		invalidFeatureGates = append(invalidFeatureGates, "JobTrackingWithFinalizers")

		// Remove --feature-gate TopologyManager starting with 1.29
		// Reference: https://github.com/kubernetes/kubernetes/pull/121252
		invalidFeatureGates = append(invalidFeatureGates, "TopologyManager")

		// Remove --feature-gate OpenAPIV3 starting with 1.29
		// Reference: https://github.com/kubernetes/kubernetes/pull/121255
		invalidFeatureGates = append(invalidFeatureGates, "OpenAPIV3")

		// Remove --feature-gate SeccompDefault starting with 1.29
		// Reference: https://github.com/kubernetes/kubernetes/pull/121246
		invalidFeatureGates = append(invalidFeatureGates, "SeccompDefault")

		// Remove --feature-gate CronJobTimeZone, JobMutableNodeSchedulingDirectives, LegacyServiceAccountTokenNoAutoGeneration starting with 1.29
		// Reference: https://github.com/kubernetes/kubernetes/pull/120192
		invalidFeatureGates = append(invalidFeatureGates, "CronJobTimeZone", "JobMutableNodeSchedulingDirectives", "LegacyServiceAccountTokenNoAutoGeneration")

		// Remove --feature-gate DownwardAPIHugePages starting with 1.29
		// Reference: https://github.com/kubernetes/kubernetes/pull/120249
		invalidFeatureGates = append(invalidFeatureGates, "DownwardAPIHugePages")

		// Remove --feature-gate GRPCContainerProbe starting with 1.29
		// Reference: https://github.com/kubernetes/kubernetes/pull/120248
		invalidFeatureGates = append(invalidFeatureGates, "GRPCContainerProbe")

		// Remove --feature-gate RetroactiveDefaultStorageClass starting with 1.29
		// Reference: https://github.com/kubernetes/kubernetes/pull/120861
		invalidFeatureGates = append(invalidFeatureGates, "RetroactiveDefaultStorageClass")
	}
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.30.0") {
		// Remove --feature-gate KubeletPodResources starting with 1.30
		// Reference: https://github.com/kubernetes/kubernetes/pull/122139
		invalidFeatureGates = append(invalidFeatureGates, "KubeletPodResources")

		// Remove --feature-gate KubeletPodResourcesGetAllocatable starting with 1.30
		// Reference: https://github.com/kubernetes/kubernetes/pull/122138
		invalidFeatureGates = append(invalidFeatureGates, "KubeletPodResourcesGetAllocatable")

		// Remove --feature-gate LegacyServiceAccountTokenTracking starting with 1.30
		// Reference: https://github.com/kubernetes/kubernetes/pull/122409
		invalidFeatureGates = append(invalidFeatureGates, "LegacyServiceAccountTokenTracking")

		// Remove --feature-gate MinimizeIPTablesRestore starting with 1.30
		// Reference: https://github.com/kubernetes/kubernetes/pull/122136
		invalidFeatureGates = append(invalidFeatureGates, "MinimizeIPTablesRestore")

		// Remove --feature-gate ProxyTerminatingEndpoints starting with 1.30
		// Reference: https://github.com/kubernetes/kubernetes/pull/122134
		invalidFeatureGates = append(invalidFeatureGates, "ProxyTerminatingEndpoints")

		// Remove --feature-gate RemoveSelfLink starting with 1.30
		// Reference: https://github.com/kubernetes/kubernetes/pull/122468
		invalidFeatureGates = append(invalidFeatureGates, "RemoveSelfLink")

		// Remove --feature-gate SecurityContextDeny starting with 1.30
		// Reference: https://github.com/kubernetes/kubernetes/pull/122612
		invalidFeatureGates = append(invalidFeatureGates, "SecurityContextDeny")

		// Remove --feature-gate APISelfSubjectReview starting with 1.30
		// Reference: https://github.com/kubernetes/kubernetes/pull/122032
		invalidFeatureGates = append(invalidFeatureGates, "APISelfSubjectReview")

		// Remove --feature-gate CSIMigrationAzureFile  starting with 1.30
		// Reference: https://github.com/kubernetes/kubernetes/pull/122576
		invalidFeatureGates = append(invalidFeatureGates, "CSIMigrationAzureFile")

		// Remove --feature-gate ExpandedDNSConfig starting with 1.30
		// Reference: https://github.com/kubernetes/kubernetes/pull/122086
		invalidFeatureGates = append(invalidFeatureGates, "ExpandedDNSConfig")

		// Remove --feature-gate ExperimentalHostUserNamespaceDefaulting starting with 1.30
		// Reference: https://github.com/kubernetes/kubernetes/pull/122088
		invalidFeatureGates = append(invalidFeatureGates, "ExperimentalHostUserNamespaceDefaulting")

		// Remove --feature-gate IPTablesOwnershipCleanup starting with 1.30
		// Reference: https://github.com/kubernetes/kubernetes/pull/122137
		invalidFeatureGates = append(invalidFeatureGates, "IPTablesOwnershipCleanup")
	}
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.31.0") {
		// Remove --feature-gate APIPriorityAndFairness starting with 1.31
		invalidFeatureGates = append(invalidFeatureGates, "APIPriorityAndFairness")

		// Remove --feature-gate ConsistentHTTPGetHandlers starting with 1.31
		invalidFeatureGates = append(invalidFeatureGates, "ConsistentHTTPGetHandlers")

		// Remove --feature-gate CSIMigrationRBD starting with 1.31
		invalidFeatureGates = append(invalidFeatureGates, "CSIMigrationRBD")

		// Remove --feature-gate CSINodeExpandSecret starting with 1.31
		invalidFeatureGates = append(invalidFeatureGates, "CSINodeExpandSecret")

		// Remove --feature-gate CustomResourceValidationExpressions starting with 1.31
		invalidFeatureGates = append(invalidFeatureGates, "CustomResourceValidationExpressions")

		// Remove --feature-gate DefaultHostNetworkHostPortsInPodTemplates starting with 1.31
		invalidFeatureGates = append(invalidFeatureGates, "DefaultHostNetworkHostPortsInPodTemplates")

		// Remove --feature-gate InTreePluginRBDUnregister starting with 1.31
		invalidFeatureGates = append(invalidFeatureGates, "InTreePluginRBDUnregister")

		// Remove --feature-gate JobReadyPods starting with 1.31
		invalidFeatureGates = append(invalidFeatureGates, "JobReadyPods")

		// Remove --feature-gate ReadWriteOncePod starting with 1.31
		invalidFeatureGates = append(invalidFeatureGates, "ReadWriteOncePod")

		// Remove --feature-gate ServiceNodePortStaticSubrange starting with 1.31
		invalidFeatureGates = append(invalidFeatureGates, "ServiceNodePortStaticSubrange")

		// Remove --feature-gate SkipReadOnlyValidationGCE starting with 1.31
		invalidFeatureGates = append(invalidFeatureGates, "SkipReadOnlyValidationGCE")
	}
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.32.0") {
		// Remove --feature-gate CloudDualStackNodeIPs starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "CloudDualStackNodeIPs")

		// Remove --feature-gate DRAControlPlaneController starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "DRAControlPlaneController")

		// Remove --feature-gate HPAContainerMetrics starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "HPAContainerMetrics")

		// Remove --feature-gate KMSv2 starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "KMSv2")

		// Remove --feature-gate KMSv2KDF starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "KMSv2KDF")

		// Remove --feature-gate LegacyServiceAccountTokenCleanUp starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "LegacyServiceAccountTokenCleanUp")

		// Remove --feature-gate MinDomainsInPodTopologySpread starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "MinDomainsInPodTopologySpread")

		// Remove --feature-gate NewVolumeManagerReconstruction starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "NewVolumeManagerReconstruction")

		// Remove --feature-gate NodeOutOfServiceVolumeDetach starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "NodeOutOfServiceVolumeDetach")

		// Remove --feature-gate PodHostIPs starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "PodHostIPs")

		// Remove --feature-gate ServerSideApply starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "ServerSideApply")

		// Remove --feature-gate ServerSideFieldValidation starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "ServerSideFieldValidation")

		// Remove --feature-gate StableLoadBalancerNodeSet starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "StableLoadBalancerNodeSet")

		// Remove --feature-gate ValidatingAdmissionPolicy starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "ValidatingAdmissionPolicy")

		// Remove --feature-gate ZeroLimitedNominalConcurrencyShares starting with 1.32
		invalidFeatureGates = append(invalidFeatureGates, "ZeroLimitedNominalConcurrencyShares")
	}
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.33.0") {
		// Remove --feature-gate AdmissionWebhookMatchConditions starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "AdmissionWebhookMatchConditions")

		// Remove --feature-gate AggregatedDiscoveryEndpoint starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "AggregatedDiscoveryEndpoint")

		// Remove --feature-gate APIListChunking starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "APIListChunking")

		// Remove --feature-gate AppArmor starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "AppArmor")

		// Remove --feature-gate AppArmorFields starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "AppArmorFields")

		// Remove --feature-gate CPUManager starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "CPUManager")

		// Remove --feature-gate DisableCloudProviders starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "DisableCloudProviders")

		// Remove --feature-gate DisableKubeletCloudCredentialProviders starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "DisableKubeletCloudCredentialProviders")

		// Remove --feature-gate EfficientWatchResumption starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "EfficientWatchResumption")

		// Remove --feature-gate JobPodFailurePolicy starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "JobPodFailurePolicy")

		// Remove --feature-gate KubeProxyDrainingTerminatingNodes starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "KubeProxyDrainingTerminatingNodes")

		// Remove --feature-gate PDBUnhealthyPodEvictionPolicy starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "PDBUnhealthyPodEvictionPolicy")

		// Remove --feature-gate PersistentVolumeLastPhaseTransitionTime starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "PersistentVolumeLastPhaseTransitionTime")

		// Remove --feature-gate RemainingItemCount starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "RemainingItemCount")

		// Remove --feature-gate VolumeCapacityPriority starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "VolumeCapacityPriority")

		// Remove --feature-gate WatchBookmark starting with 1.33
		invalidFeatureGates = append(invalidFeatureGates, "WatchBookmark")
	}
	removeInvalidFeatureGates(o.KubernetesConfig.CloudControllerManagerConfig, invalidFeatureGates)

	// TODO add RBAC support
	/*if *o.KubernetesConfig.EnableRbac {
		o.KubernetesConfig.CloudControllerManagerConfig["--use-service-account-credentials"] = "true"
	}*/
}
