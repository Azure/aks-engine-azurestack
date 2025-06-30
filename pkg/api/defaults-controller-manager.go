// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"strconv"
	"strings"

	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers/to"
)

func (cs *ContainerService) setControllerManagerConfig() {
	o := cs.Properties.OrchestratorProfile
	isAzureCNIDualStack := cs.Properties.IsAzureCNIDualStack()
	clusterCidr := o.KubernetesConfig.ClusterSubnet
	if isAzureCNIDualStack {
		clusterSubnets := strings.Split(clusterCidr, ",")
		if len(clusterSubnets) > 1 {
			clusterCidr = clusterSubnets[1]
		}
	}
	staticControllerManagerConfig := map[string]string{
		"--kubeconfig":                       "/var/lib/kubelet/kubeconfig",
		"--allocate-node-cidrs":              strconv.FormatBool(!o.IsAzureCNI() || isAzureCNIDualStack),
		"--configure-cloud-routes":           strconv.FormatBool(cs.Properties.RequireRouteTable()),
		"--cluster-cidr":                     clusterCidr,
		"--root-ca-file":                     "/etc/kubernetes/certs/ca.crt",
		"--cluster-signing-cert-file":        "/etc/kubernetes/certs/ca.crt",
		"--cluster-signing-key-file":         "/etc/kubernetes/certs/ca.key",
		"--service-account-private-key-file": "/etc/kubernetes/certs/apiserver.key",
		"--leader-elect":                     "true",
		"--v":                                "2",
		"--controllers":                      "*,bootstrapsigner,tokencleaner",
	}

	// Set --cluster-name based on appropriate DNS prefix
	if cs.Properties.MasterProfile != nil {
		staticControllerManagerConfig["--cluster-name"] = cs.Properties.MasterProfile.DNSPrefix
	}

	// Enable cloudprovider if we're not using cloud controller manager
	if !to.Bool(o.KubernetesConfig.UseCloudControllerManager) {
		staticControllerManagerConfig["--cloud-provider"] = "azure"
		staticControllerManagerConfig["--cloud-config"] = "/etc/kubernetes/azure.json"
	} else {
		staticControllerManagerConfig["--cloud-provider"] = "external"
	}

	ctrlMgrNodeMonitorGracePeriod := DefaultKubernetesCtrlMgrNodeMonitorGracePeriod
	ctrlMgrPodEvictionTimeout := DefaultKubernetesCtrlMgrPodEvictionTimeout
	ctrlMgrRouteReconciliationPeriod := DefaultKubernetesCtrlMgrRouteReconciliationPeriod

	if cs.Properties.IsAzureStackCloud() {
		ctrlMgrNodeMonitorGracePeriod = DefaultAzureStackKubernetesCtrlMgrNodeMonitorGracePeriod
		ctrlMgrPodEvictionTimeout = DefaultAzureStackKubernetesCtrlMgrPodEvictionTimeout
		ctrlMgrRouteReconciliationPeriod = DefaultAzureStackKubernetesCtrlMgrRouteReconciliationPeriod
	}

	// Default controller-manager config
	defaultControllerManagerConfig := map[string]string{
		"--bind-address":                    "127.0.0.1", // STIG Rule ID: SV-242385r879530_rule
		"--node-monitor-grace-period":       ctrlMgrNodeMonitorGracePeriod,
		"--pod-eviction-timeout":            ctrlMgrPodEvictionTimeout,
		"--route-reconciliation-period":     ctrlMgrRouteReconciliationPeriod,
		"--terminated-pod-gc-threshold":     DefaultKubernetesCtrlMgrTerminatedPodGcThreshold,
		"--tls-min-version":                 "VersionTLS12", // STIG Rule ID: SV-242376r879519_rule
		"--use-service-account-credentials": DefaultKubernetesCtrlMgrUseSvcAccountCreds,
		"--profiling":                       DefaultKubernetesCtrMgrEnableProfiling,
	}

	// If no user-configurable controller-manager config values exists, use the defaults
	if o.KubernetesConfig.ControllerManagerConfig == nil {
		o.KubernetesConfig.ControllerManagerConfig = defaultControllerManagerConfig
	} else {
		for key, val := range defaultControllerManagerConfig {
			// If we don't have a user-configurable controller-manager config for each option
			if _, ok := o.KubernetesConfig.ControllerManagerConfig[key]; !ok {
				// then assign the default value
				o.KubernetesConfig.ControllerManagerConfig[key] = val
			}
		}
	}

	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.24.0") {
		// https://github.com/kubernetes/kubernetes/pull/106860
		removedFlags124 := []string{"--address", "--port"}
		for _, key := range removedFlags124 {
			delete(o.KubernetesConfig.ControllerManagerConfig, key)
		}
	}

	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.25.0") {
		// https://github.com/kubernetes/kubernetes/pull/109612
		removedFlags125 := []string{"--deleting-pods-qps", "--deleting-pods-burst", "--register-retry-count"}
		for _, key := range removedFlags125 {
			delete(o.KubernetesConfig.ControllerManagerConfig, key)
		}
	}

	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.27.0") {
		// https://github.com/kubernetes/kubernetes/pull/115840
		removedFlags127 := []string{"--enable-taint-manager", "--pod-eviction-timeout"}
		for _, key := range removedFlags127 {
			delete(o.KubernetesConfig.ControllerManagerConfig, key)
		}
	}

	// Enables Node Exclusion from Services (toggled on agent nodes by the alpha.service-controller.kubernetes.io/exclude-balancer label).
	// ServiceNodeExclusion feature gate is GA in 1.19, removed in 1.22 (xref: https://github.com/kubernetes/kubernetes/pull/100776)
	if !common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.19.0") {
		addDefaultFeatureGates(o.KubernetesConfig.ControllerManagerConfig, o.OrchestratorVersion, "1.9.0", "ServiceNodeExclusion=true")
	}

	// Enable the consumption of local ephemeral storage and also the sizeLimit property of an emptyDir volume.
	addDefaultFeatureGates(o.KubernetesConfig.ControllerManagerConfig, o.OrchestratorVersion, "1.10.0", "LocalStorageCapacityIsolation=true")

	// LegacyServiceAccountTokenNoAutoGeneration feature gate is forced by Kubernetes to true in v1.27, and will be removed in v1.29 (https://github.com/kubernetes/kubernetes/pull/114522)
	if !common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.27.0") {
		// Enable legacy service account token autogeneration for >= v1.24.0 and < v1.27.0
		addDefaultFeatureGates(o.KubernetesConfig.ControllerManagerConfig, o.OrchestratorVersion, "1.24.0", "LegacyServiceAccountTokenNoAutoGeneration=false")
	}
	// STIG Rule ID: SV-254801r879719_rule
	addDefaultFeatureGates(o.KubernetesConfig.ControllerManagerConfig, o.OrchestratorVersion, "1.25.0", "PodSecurity=true")

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	for key, val := range staticControllerManagerConfig {
		o.KubernetesConfig.ControllerManagerConfig[key] = val
	}

	if o.KubernetesConfig.IsRBACEnabled() {
		o.KubernetesConfig.ControllerManagerConfig["--use-service-account-credentials"] = "true"
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
		// LegacyServiceAccountTokenNoAutoGeneration feature gate is forced by Kubernetes to true in v1.27, and will be removed in v1.29 (https://github.com/kubernetes/kubernetes/pull/114522).
		// Preemptively forcing removal of the feature gate now since the feature gate can only be true, and by default the token will not be autogenerated.
		invalidFeatureGates = append(invalidFeatureGates, "LegacyServiceAccountTokenNoAutoGeneration")

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
	removeInvalidFeatureGates(o.KubernetesConfig.ControllerManagerConfig, invalidFeatureGates)
}
