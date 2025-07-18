// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers/to"
)

func (cs *ContainerService) setKubeletConfig(isUpgrade bool) {
	o := cs.Properties.OrchestratorProfile
	staticLinuxKubeletConfig := map[string]string{
		"--address":                           "0.0.0.0",
		"--allow-privileged":                  "true",
		"--anonymous-auth":                    "false",
		"--authorization-mode":                "Webhook",
		"--client-ca-file":                    "/etc/kubernetes/certs/ca.crt",
		"--pod-manifest-path":                 "/etc/kubernetes/manifests",
		"--cluster-dns":                       o.KubernetesConfig.DNSServiceIP,
		"--cgroups-per-qos":                   "true",
		"--kubeconfig":                        "/var/lib/kubelet/kubeconfig",
		"--keep-terminated-pod-volumes":       "false",
		"--tls-cert-file":                     "/etc/kubernetes/certs/kubeletserver.crt",
		"--tls-private-key-file":              "/etc/kubernetes/certs/kubeletserver.key",
		"--v":                                 "2",
		"--volume-plugin-dir":                 "/etc/kubernetes/volumeplugins",
		"--image-credential-provider-config":  "/var/lib/kubelet/credential-provider-config.yaml",
		"--image-credential-provider-bin-dir": "/var/lib/kubelet/credential-provider",
	}

	for key := range staticLinuxKubeletConfig {
		switch key {
		case "--anonymous-auth", "--client-ca-file":
			if !to.Bool(o.KubernetesConfig.EnableSecureKubelet) { // Don't add if EnableSecureKubelet is disabled
				delete(staticLinuxKubeletConfig, key)
			}
		}
	}

	// Start with copy of Linux config
	staticWindowsKubeletConfig := make(map[string]string)
	for key, val := range staticLinuxKubeletConfig {
		switch key {
		case "--pod-manifest-path", "--tls-cert-file", "--tls-private-key-file": // Don't add Linux-specific config
			staticWindowsKubeletConfig[key] = ""
		case "--anonymous-auth":
			if !to.Bool(o.KubernetesConfig.EnableSecureKubelet) { // Don't add if EnableSecureKubelet is disabled
				staticWindowsKubeletConfig[key] = ""
			} else {
				staticWindowsKubeletConfig[key] = val
			}
		case "--client-ca-file":
			if !to.Bool(o.KubernetesConfig.EnableSecureKubelet) { // Don't add if EnableSecureKubelet is disabled
				staticWindowsKubeletConfig[key] = ""
			} else {
				staticWindowsKubeletConfig[key] = "c:\\k\\ca.crt"
			}
		default:
			staticWindowsKubeletConfig[key] = val
		}
	}

	// Add Windows-specific overrides
	// Eventually paths should not be hardcoded here. They should be relative to $global:KubeDir in the PowerShell script
	staticWindowsKubeletConfig["--image-credential-provider-config"] = "c:\\k\\credential-provider\\credential-provider-config.yaml"
	staticWindowsKubeletConfig["--image-credential-provider-bin-dir"] = "c:\\k\\credential-provider"
	staticWindowsKubeletConfig["--pod-infra-container-image"] = "kubletwin/pause"
	staticWindowsKubeletConfig["--kubeconfig"] = "c:\\k\\config"
	staticWindowsKubeletConfig["--cloud-config"] = "c:\\k\\azure.json"
	staticWindowsKubeletConfig["--cgroups-per-qos"] = "false"
	staticWindowsKubeletConfig["--enforce-node-allocatable"] = "\"\"\"\""
	staticWindowsKubeletConfig["--system-reserved"] = "memory=2Gi"
	staticWindowsKubeletConfig["--hairpin-mode"] = "promiscuous-bridge"
	staticWindowsKubeletConfig["--image-pull-progress-deadline"] = "20m"
	staticWindowsKubeletConfig["--resolv-conf"] = "\"\"\"\""
	staticWindowsKubeletConfig["--eviction-hard"] = "\"\"\"\""

	nodeStatusUpdateFrequency := GetK8sComponentsByVersionMap(o.KubernetesConfig)[o.OrchestratorVersion]["nodestatusfreq"]
	if cs.Properties.IsAzureStackCloud() {
		nodeStatusUpdateFrequency = DefaultAzureStackKubernetesNodeStatusUpdateFrequency
	}

	// Default Kubelet config
	defaultKubeletConfig := map[string]string{
		"--cluster-domain":                    "cluster.local",
		"--network-plugin":                    "cni",
		"--pod-infra-container-image":         o.KubernetesConfig.MCRKubernetesImageBase + GetK8sComponentsByVersionMap(o.KubernetesConfig)[o.OrchestratorVersion][common.PauseComponentName],
		"--max-pods":                          strconv.Itoa(DefaultKubernetesMaxPods),
		"--eviction-hard":                     DefaultKubernetesHardEvictionThreshold,
		"--node-status-update-frequency":      nodeStatusUpdateFrequency,
		"--image-gc-high-threshold":           strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"--image-gc-low-threshold":            strconv.Itoa(DefaultKubernetesGCLowThreshold),
		"--non-masquerade-cidr":               DefaultNonMasqueradeCIDR,
		"--cloud-provider":                    "azure",
		"--cloud-config":                      "/etc/kubernetes/azure.json",
		"--event-qps":                         DefaultKubeletEventQPS,
		"--cadvisor-port":                     DefaultKubeletCadvisorPort,
		"--pod-max-pids":                      strconv.Itoa(DefaultKubeletPodMaxPIDs),
		"--image-pull-progress-deadline":      "30m",
		"--enforce-node-allocatable":          "pods",
		"--streaming-connection-idle-timeout": "5m", // STIG Rule ID: SV-245541r879622_rule
		"--tls-cipher-suites":                 TLSStrongCipherSuitesKubelet,
		"--healthz-port":                      DefaultKubeletHealthzPort,
		"--seccomp-default":                   "true",
	}

	// Set --non-masquerade-cidr if ip-masq-agent is disabled on AKS or
	// explicitly disabled in kubernetes config
	if cs.Properties.IsIPMasqAgentDisabled() {
		defaultKubeletConfig["--non-masquerade-cidr"] = cs.Properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet
	}

	// Apply Azure CNI-specific --max-pods value
	if o.KubernetesConfig.NetworkPlugin == NetworkPluginAzure {
		defaultKubeletConfig["--max-pods"] = strconv.Itoa(DefaultKubernetesMaxPodsVNETIntegrated)
	}

	minVersionRotateCerts := "1.11.9"
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, minVersionRotateCerts) {
		defaultKubeletConfig["--rotate-certificates"] = "true"
	}

	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.16.0") {
		// for enabling metrics-server v0.3.0+
		defaultKubeletConfig["--authentication-token-webhook"] = "true"
		defaultKubeletConfig["--read-only-port"] = "0" // we only have metrics-server v0.3 support in 1.16.0 and above
	}

	if o.KubernetesConfig.NeedsContainerd() {
		// Kubelet flag --container-runtime has been removed from k8s 1.27
		// Reference: https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.27.md#other-cleanup-or-flake
		if !common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.27.0") {
			defaultKubeletConfig["--container-runtime"] = "remote"
		}
		defaultKubeletConfig["--runtime-request-timeout"] = "15m"
		defaultKubeletConfig["--container-runtime-endpoint"] = "unix:///run/containerd/containerd.sock"
	}

	// If no user-configurable kubelet config values exists, use the defaults
	setMissingKubeletValues(o.KubernetesConfig, defaultKubeletConfig)
	if isUpgrade {
		// if upgrade, force default "--pod-infra-container-image" value
		o.KubernetesConfig.KubeletConfig["--pod-infra-container-image"] = defaultKubeletConfig["--pod-infra-container-image"]
	}

	addDefaultFeatureGates(o.KubernetesConfig.KubeletConfig, o.OrchestratorVersion, minVersionRotateCerts, "RotateKubeletServerCertificate=true")
	addDefaultFeatureGates(o.KubernetesConfig.KubeletConfig, o.OrchestratorVersion, "1.20.0-rc.0", "ExecProbeTimeout=true")
	// STIG Rule ID: SV-254801r879719_rule
	addDefaultFeatureGates(o.KubernetesConfig.KubeletConfig, o.OrchestratorVersion, "1.25.0", "PodSecurity=true")

	// Override default cloud-provider?
	if to.Bool(o.KubernetesConfig.UseCloudControllerManager) {
		staticLinuxKubeletConfig["--cloud-provider"] = "external"
	}

	// Override default --network-plugin?
	if o.KubernetesConfig.NetworkPlugin == NetworkPluginKubenet {
		if o.KubernetesConfig.NetworkPolicy != NetworkPolicyCalico {
			o.KubernetesConfig.KubeletConfig["--network-plugin"] = NetworkPluginKubenet
		}
	}

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	for key, val := range staticLinuxKubeletConfig {
		o.KubernetesConfig.KubeletConfig[key] = val
	}

	if isUpgrade && common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.14.0") {
		hasSupportPodPidsLimitFeatureGate := strings.Contains(o.KubernetesConfig.KubeletConfig["--feature-gates"], "SupportPodPidsLimit=true")
		podMaxPids, err := strconv.Atoi(o.KubernetesConfig.KubeletConfig["--pod-max-pids"])
		if err != nil {
			o.KubernetesConfig.KubeletConfig["--pod-max-pids"] = strconv.Itoa(-1)
		} else {
			// If we don't have an explicit SupportPodPidsLimit=true, disable --pod-max-pids by setting to -1
			// To prevent older clusters from inheriting SupportPodPidsLimit=true implicitly starting w/ 1.14.0
			if !hasSupportPodPidsLimitFeatureGate || podMaxPids <= 0 {
				o.KubernetesConfig.KubeletConfig["--pod-max-pids"] = strconv.Itoa(-1)
			}
		}
	}

	removeKubeletFlags(o.KubernetesConfig.KubeletConfig, o.OrchestratorVersion)

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
	removeInvalidFeatureGates(o.KubernetesConfig.KubeletConfig, invalidFeatureGates)

	// Master-specific kubelet config changes go here
	if cs.Properties.MasterProfile != nil {
		if cs.Properties.MasterProfile.KubernetesConfig == nil {
			cs.Properties.MasterProfile.KubernetesConfig = &KubernetesConfig{}
		}
		if cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig == nil {
			cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig = make(map[string]string)
		}
		if isUpgrade {
			// if upgrade, force default "--pod-infra-container-image" value
			cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--pod-infra-container-image"] = o.KubernetesConfig.KubeletConfig["--pod-infra-container-image"]
		}
		//Ensure cloud-provider setting
		if to.Bool(o.KubernetesConfig.UseCloudControllerManager) {
			cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--cloud-provider"] = "external"
		}
		setMissingKubeletValues(cs.Properties.MasterProfile.KubernetesConfig, o.KubernetesConfig.KubeletConfig)
		addDefaultFeatureGates(cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig, o.OrchestratorVersion, "", "")

		if isUpgrade && common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.14.0") {
			hasSupportPodPidsLimitFeatureGate := strings.Contains(cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--feature-gates"], "SupportPodPidsLimit=true")
			podMaxPids, err := strconv.Atoi(cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--pod-max-pids"])
			if err != nil {
				cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--pod-max-pids"] = strconv.Itoa(-1)
			} else {
				// If we don't have an explicit SupportPodPidsLimit=true, disable --pod-max-pids by setting to -1
				// To prevent older clusters from inheriting SupportPodPidsLimit=true implicitly starting w/ 1.14.0
				if !hasSupportPodPidsLimitFeatureGate || podMaxPids <= 0 {
					cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--pod-max-pids"] = strconv.Itoa(-1)
				}
			}
		}
		// "--protect-kernel-defaults" is only true for VHD based VMs since the base Ubuntu distros don't have a /etc/sysctl.d/60-CIS.conf file.
		if cs.Properties.MasterProfile.IsVHDDistro() {
			if _, ok := cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--protect-kernel-defaults"]; !ok {
				cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--protect-kernel-defaults"] = "true"
			}
		}
		// Override the --resolv-conf kubelet config value for Ubuntu 18.04 after the distro value is set.
		if cs.Properties.MasterProfile.IsUbuntu1804() || cs.Properties.MasterProfile.IsUbuntu2004() || cs.Properties.MasterProfile.IsUbuntu2204() {
			cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--resolv-conf"] = "/run/systemd/resolve/resolv.conf"
		}

		removeKubeletFlags(cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig, o.OrchestratorVersion)
		removeInvalidFeatureGates(cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig, invalidFeatureGates)

		if cs.Properties.AnyAgentIsLinux() {
			if val, ok := cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--register-with-taints"]; !ok {
				cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--register-with-taints"] = common.MasterNodeTaint
			} else {
				if !strings.Contains(val, common.MasterNodeTaint) {
					cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--register-with-taints"] += fmt.Sprintf(",%s", common.MasterNodeTaint)
				}
			}
		}
	}

	// Agent-specific kubelet config changes go here
	for _, profile := range cs.Properties.AgentPoolProfiles {
		if profile.KubernetesConfig == nil {
			profile.KubernetesConfig = &KubernetesConfig{}
		}
		if profile.KubernetesConfig.KubeletConfig == nil {
			profile.KubernetesConfig.KubeletConfig = make(map[string]string)
		}

		if isUpgrade {
			// if upgrade, force default "--pod-infra-container-image" value
			profile.KubernetesConfig.KubeletConfig["--pod-infra-container-image"] = o.KubernetesConfig.KubeletConfig["--pod-infra-container-image"]
		}

		if profile.IsWindows() {
			for key, val := range staticWindowsKubeletConfig {
				profile.KubernetesConfig.KubeletConfig[key] = val
			}
		} else {
			for key, val := range staticLinuxKubeletConfig {
				profile.KubernetesConfig.KubeletConfig[key] = val
			}
		}

		setMissingKubeletValues(profile.KubernetesConfig, o.KubernetesConfig.KubeletConfig)

		if isUpgrade && common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.14.0") {
			hasSupportPodPidsLimitFeatureGate := strings.Contains(profile.KubernetesConfig.KubeletConfig["--feature-gates"], "SupportPodPidsLimit=true")
			podMaxPids, err := strconv.Atoi(profile.KubernetesConfig.KubeletConfig["--pod-max-pids"])
			if err != nil {
				profile.KubernetesConfig.KubeletConfig["--pod-max-pids"] = strconv.Itoa(-1)
			} else {
				// If we don't have an explicit SupportPodPidsLimit=true, disable --pod-max-pids by setting to -1
				// To prevent older clusters from inheriting SupportPodPidsLimit=true implicitly starting w/ 1.14.0
				if !hasSupportPodPidsLimitFeatureGate || podMaxPids <= 0 {
					profile.KubernetesConfig.KubeletConfig["--pod-max-pids"] = strconv.Itoa(-1)
				}
			}
		}

		// "--protect-kernel-defaults" is only true for VHD based VMs since the base Ubuntu distros don't have a /etc/sysctl.d/60-CIS.conf file.
		if profile.IsVHDDistro() {
			if _, ok := profile.KubernetesConfig.KubeletConfig["--protect-kernel-defaults"]; !ok {
				profile.KubernetesConfig.KubeletConfig["--protect-kernel-defaults"] = "true"
			}
		}
		// Override the --resolv-conf kubelet config value for Ubuntu 18.04 after the distro value is set.
		if profile.IsUbuntu1804() || profile.IsUbuntu2004() || profile.IsUbuntu2204() {
			profile.KubernetesConfig.KubeletConfig["--resolv-conf"] = "/run/systemd/resolve/resolv.conf"
		}

		removeKubeletFlags(profile.KubernetesConfig.KubeletConfig, o.OrchestratorVersion)
		removeInvalidFeatureGates(profile.KubernetesConfig.KubeletConfig, invalidFeatureGates)
		if cs.Properties.OrchestratorProfile.KubernetesConfig.IsAddonEnabled(common.AADPodIdentityAddonName) && !profile.IsWindows() {
			if val, ok := profile.KubernetesConfig.KubeletConfig["--register-with-taints"]; !ok {
				profile.KubernetesConfig.KubeletConfig["--register-with-taints"] = fmt.Sprintf("%s=true:NoSchedule", common.AADPodIdentityTaintKey)
			} else {
				if !strings.Contains(val, common.AADPodIdentityTaintKey) {
					profile.KubernetesConfig.KubeletConfig["--register-with-taints"] += fmt.Sprintf(",%s=true:NoSchedule", common.AADPodIdentityTaintKey)
				}
			}
		}
	}
}

func removeKubeletFlags(k map[string]string, v string) {
	// Get rid of values not supported until v1.10
	if !common.IsKubernetesVersionGe(v, "1.10.0") {
		for _, key := range []string{"--pod-max-pids"} {
			delete(k, key)
		}
	}

	// Get rid of values not supported in v1.12 and up
	if common.IsKubernetesVersionGe(v, "1.12.0") {
		for _, key := range []string{"--cadvisor-port"} {
			delete(k, key)
		}
	}

	// Get rid of values not supported in v1.15 and up
	if common.IsKubernetesVersionGe(v, "1.15.0-beta.1") {
		for _, key := range []string{"--allow-privileged"} {
			delete(k, key)
		}
	}

	// Remove dockershim related flags in v1.24 and up
	if common.IsKubernetesVersionGe(v, "1.24.0-alpha") {
		for _, key := range []string{
			"--cni-conf-dir",
			"--cni-bin-dir",
			"--cni-cache-dir",
			"--docker-endpoint",
			"--experimental-dockershim-root-directory",
			"--image-pull-progress-deadline",
			"--network-plugin",
			"--network-plugin-mtu",
			"--non-masquerade-cidr",
		} {
			delete(k, key)
		}
	}

	// Get rid of values not supported in v1.27 and up
	if common.IsKubernetesVersionGe(v, "1.27.0") {
		for _, key := range []string{"--master-service-namespace", "--container-runtime"} {
			delete(k, key)
		}
	}

	// Get rid of values not supported in v1.30 and up
	if common.IsKubernetesVersionGe(v, "1.30.0") {
		for _, key := range []string{"--azure-container-registry-config"} {
			delete(k, key)
		}
	}
}

func setMissingKubeletValues(p *KubernetesConfig, d map[string]string) {
	if p.KubeletConfig == nil {
		p.KubeletConfig = d
	} else {
		for key, val := range d {
			// If we don't have a user-configurable value for each option
			if _, ok := p.KubeletConfig[key]; !ok {
				// then assign the default value
				p.KubeletConfig[key] = val
			}
		}
	}
}
