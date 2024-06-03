// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"fmt"
	"strconv"

	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers/to"
)

func (cs *ContainerService) setAPIServerConfig() {
	o := cs.Properties.OrchestratorProfile
	staticAPIServerConfig := map[string]string{
		"--bind-address":                "0.0.0.0",
		"--advertise-address":           "<advertiseAddr>",
		"--allow-privileged":            "true",
		"--audit-log-path":              "/var/log/kubeaudit/audit.log",
		"--secure-port":                 "443",
		"--service-account-lookup":      "true",
		"--etcd-certfile":               "/etc/kubernetes/certs/etcdclient.crt",
		"--etcd-keyfile":                "/etc/kubernetes/certs/etcdclient.key",
		"--tls-cert-file":               "/etc/kubernetes/certs/apiserver.crt",
		"--tls-private-key-file":        "/etc/kubernetes/certs/apiserver.key",
		"--client-ca-file":              "/etc/kubernetes/certs/ca.crt",
		"--service-account-key-file":    "/etc/kubernetes/certs/apiserver.key",
		"--kubelet-client-certificate":  "/etc/kubernetes/certs/client.crt",
		"--kubelet-client-key":          "/etc/kubernetes/certs/client.key",
		"--service-cluster-ip-range":    o.KubernetesConfig.ServiceCIDR,
		"--storage-backend":             o.GetAPIServerEtcdAPIVersion(),
		"--enable-bootstrap-token-auth": "true",
	}

	if cs.Properties.MasterProfile != nil {
		if cs.Properties.MasterProfile.HasCosmosEtcd() {
			// Configuration for cosmos etcd
			staticAPIServerConfig["--etcd-servers"] = fmt.Sprintf("https://%s:%s", cs.Properties.MasterProfile.GetCosmosEndPointURI(), strconv.Itoa(DefaultMasterEtcdClientPort))
		} else {
			// Configuration for local etcd
			staticAPIServerConfig["--etcd-cafile"] = "/etc/kubernetes/certs/ca.crt"
			staticAPIServerConfig["--etcd-servers"] = fmt.Sprintf("https://127.0.0.1:%s", strconv.Itoa(DefaultMasterEtcdClientPort))
		}
	}

	// Data Encryption at REST configuration conditions
	if to.Bool(o.KubernetesConfig.EnableDataEncryptionAtRest) || to.Bool(o.KubernetesConfig.EnableEncryptionWithExternalKms) {
		staticAPIServerConfig["--encryption-provider-config"] = "/etc/kubernetes/encryption-config.yaml"
	}

	// Enable cloudprovider if we're not using cloud controller manager
	if !to.Bool(o.KubernetesConfig.UseCloudControllerManager) {
		staticAPIServerConfig["--cloud-provider"] = "azure"
		staticAPIServerConfig["--cloud-config"] = "/etc/kubernetes/azure.json"
	}

	// Default apiserver config
	defaultAPIServerConfig := map[string]string{
		"--admission-control-config-file": "/etc/kubernetes/apiserver-admission-control.yaml",
		"--anonymous-auth":                "false",
		"--audit-log-maxage":              "30",
		"--audit-log-maxbackup":           "10",
		"--audit-log-maxsize":             "100",
		"--profiling":                     DefaultKubernetesAPIServerEnableProfiling,
		"--request-timeout":               "1m", // STIG Rule ID: SV-242438r879806_rule
		"--tls-cipher-suites":             TLSStrongCipherSuitesAPIServer,
		"--tls-min-version":               "VersionTLS12", // STIG Rule ID: SV-242468r879889_rule
		"--v":                             DefaultKubernetesAPIServerVerbosity,
	}

	// Aggregated API configuration
	if o.KubernetesConfig.EnableAggregatedAPIs {
		defaultAPIServerConfig["--requestheader-client-ca-file"] = "/etc/kubernetes/certs/proxy-ca.crt"
		defaultAPIServerConfig["--proxy-client-cert-file"] = "/etc/kubernetes/certs/proxy.crt"
		defaultAPIServerConfig["--proxy-client-key-file"] = "/etc/kubernetes/certs/proxy.key"
		defaultAPIServerConfig["--requestheader-allowed-names"] = ""
		defaultAPIServerConfig["--requestheader-extra-headers-prefix"] = "X-Remote-Extra-"
		defaultAPIServerConfig["--requestheader-group-headers"] = "X-Remote-Group"
		defaultAPIServerConfig["--requestheader-username-headers"] = "X-Remote-User"
	}

	// AAD configuration
	if cs.Properties.HasAadProfile() {
		defaultAPIServerConfig["--oidc-username-claim"] = "oid"
		defaultAPIServerConfig["--oidc-groups-claim"] = "groups"
		defaultAPIServerConfig["--oidc-client-id"] = "spn:" + cs.Properties.AADProfile.ServerAppID
		issuerHost := "sts.windows.net"
		if helpers.GetTargetEnv(cs.Location, cs.Properties.GetCustomCloudName()) == "AzureChinaCloud" {
			issuerHost = "sts.chinacloudapi.cn"
		}
		defaultAPIServerConfig["--oidc-issuer-url"] = "https://" + issuerHost + "/" + cs.Properties.AADProfile.TenantID + "/"
	}

	// Audit Policy configuration
	defaultAPIServerConfig["--audit-policy-file"] = "/etc/kubernetes/addons/audit-policy.yaml"

	// RBAC configuration
	if to.Bool(o.KubernetesConfig.EnableRbac) {
		defaultAPIServerConfig["--authorization-mode"] = "Node,RBAC"
	}

	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.20.0-alpha.1") {
		defaultAPIServerConfig["--service-account-issuer"] = "https://kubernetes.default.svc.cluster.local"
		defaultAPIServerConfig["--service-account-signing-key-file"] = "/etc/kubernetes/certs/apiserver.key"
	}

	if !common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.20.0-alpha.0") {
		defaultAPIServerConfig["--insecure-port"] = "0"
	}

	// Set default admission controllers
	admissionControlKey, admissionControlValues := getDefaultAdmissionControls(cs)
	defaultAPIServerConfig[admissionControlKey] = admissionControlValues

	// If no user-configurable apiserver config values exists, use the defaults
	if o.KubernetesConfig.APIServerConfig == nil {
		o.KubernetesConfig.APIServerConfig = defaultAPIServerConfig
	} else {
		for key, val := range defaultAPIServerConfig {
			// If we don't have a user-configurable apiserver config for each option
			if _, ok := o.KubernetesConfig.APIServerConfig[key]; !ok {
				// then assign the default value
				o.KubernetesConfig.APIServerConfig[key] = val
			} else {
				// Manual override of "--audit-policy-file" for back-compat
				if key == "--audit-policy-file" {
					if o.KubernetesConfig.APIServerConfig[key] == "/etc/kubernetes/manifests/audit-policy.yaml" {
						o.KubernetesConfig.APIServerConfig[key] = val
					}
				}
			}
		}
	}

	// STIG Rule ID: SV-254801r879719_rule
	addDefaultFeatureGates(o.KubernetesConfig.APIServerConfig, o.OrchestratorVersion, "1.25.0", "PodSecurity=true")

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	for key, val := range staticAPIServerConfig {
		o.KubernetesConfig.APIServerConfig[key] = val
	}

	// Remove flags for secure communication to kubelet, if configured
	if !to.Bool(o.KubernetesConfig.EnableSecureKubelet) {
		for _, key := range []string{"--kubelet-client-certificate", "--kubelet-client-key"} {
			delete(o.KubernetesConfig.APIServerConfig, key)
		}
	}

	// Enforce flags removal that don't work with specific versions, to accommodate upgrade
	// Remove flags that are not compatible with any supported versions
	for _, key := range []string{"--admission-control", "--repair-malformed-updates"} {
		delete(o.KubernetesConfig.APIServerConfig, key)
	}

	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.24.0") {
		// https://github.com/kubernetes/kubernetes/pull/106859
		removedFlags124 := []string{"--address", "--insecure-bind-address", "--port", "--insecure-port"}
		for _, key := range removedFlags124 {
			delete(o.KubernetesConfig.APIServerConfig, key)
		}
	}

	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.25.0") {
		// https://github.com/kubernetes/kubernetes/pull/108624
		removedFlags125 := []string{"--service-account-api-audiences"}
		for _, key := range removedFlags125 {
			delete(o.KubernetesConfig.APIServerConfig, key)
		}
	}

	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.27.0") {
		// https://github.com/kubernetes/kubernetes/pull/114446
		removedFlags127 := []string{"--master-service-namespace"}
		for _, key := range removedFlags127 {
			delete(o.KubernetesConfig.APIServerConfig, key)
		}
	}

	// Set bind address to prefer IPv6 address for single stack IPv6 cluster
	// Remove --advertise-address so that --bind-address will be used
	if cs.Properties.FeatureFlags.IsFeatureEnabled("EnableIPv6Only") {
		o.KubernetesConfig.APIServerConfig["--bind-address"] = "::"
		for _, key := range []string{"--advertise-address"} {
			delete(o.KubernetesConfig.APIServerConfig, key)
		}
	}

	// Manual override of "--service-account-issuer" starting with 1.20
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.20.0-alpha.1") && o.KubernetesConfig.APIServerConfig["--service-account-issuer"] == "kubernetes.default.svc" {
		o.KubernetesConfig.APIServerConfig["--service-account-issuer"] = "https://kubernetes.default.svc.cluster.local"
	}

	cs.overrideAPIServerConfig()
}

func getDefaultAdmissionControls(cs *ContainerService) (string, string) {
	o := cs.Properties.OrchestratorProfile
	admissionControlKey := "--enable-admission-plugins"
	// Only include admission controllers that are not enabled by default
	admissionControlValues := "ExtendedResourceToleration"

	// Pod Security Policy configuration
	if o.KubernetesConfig.IsAddonEnabled(common.PodSecurityPolicyAddonName) {
		admissionControlValues += ",PodSecurityPolicy"
	}

	return admissionControlKey, admissionControlValues
}

// overrideAPIServerConfig fixes the kube-apiserver configuration,
// mostly by cleaning up removed features (flags, gates or admission controllers)
func (cs *ContainerService) overrideAPIServerConfig() {
	o := cs.Properties.OrchestratorProfile

	invalidFeatureGates := []string{}
	// Remove --feature-gate VolumeSnapshotDataSource starting with 1.22
	// Reference: https://github.com/kubernetes/kubernetes/pull/101531
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
		// Remove --feature-gate AdvancedAuditing,DisableAcceleratorUsageMetrics,DryRun,PodSecurity starting with 1.28
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
	removeInvalidFeatureGates(o.KubernetesConfig.APIServerConfig, invalidFeatureGates)

	if common.ShouldDisablePodSecurityPolicyAddon(o.OrchestratorVersion) {
		curPlugins := o.KubernetesConfig.APIServerConfig["--enable-admission-plugins"]
		newPlugins := common.RemoveFromCommaSeparatedList(curPlugins, "PodSecurityPolicy")
		o.KubernetesConfig.APIServerConfig["--enable-admission-plugins"] = newPlugins
	}
}
