// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers/to"
	"github.com/google/go-cmp/cmp"
)

func TestKubeletConfigDefaults(t *testing.T) {
	cs := CreateMockContainerService("testcluster", common.RationalizeReleaseAndVersion(Kubernetes, common.KubernetesDefaultRelease, "", false, false, false), 3, 2, false)
	winProfile := &AgentPoolProfile{}
	winProfile.Count = 1
	winProfile.Name = "agentpool2"
	winProfile.VMSize = "Standard_D2_v2"
	winProfile.OSType = Windows
	cs.Properties.AgentPoolProfiles = append(cs.Properties.AgentPoolProfiles, winProfile)
	cs.Properties.OrchestratorProfile.KubernetesConfig.Addons = []KubernetesAddon{
		{
			Name:    common.AADPodIdentityAddonName,
			Enabled: to.BoolPtr(true),
		},
	}
	cs.setKubeletConfig(false)
	kubeletConfig := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	expected := getDefaultLinuxKubeletConfig(cs)
	for key, val := range kubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}
	expected["--register-with-taints"] = common.MasterNodeTaint
	masterKubeletConfig := cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig
	for key, val := range masterKubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected masterProfile kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}
	cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--register-with-taints"] = "node-role.kubernetes.io/customtaint=true:NoSchedule"
	cs.setKubeletConfig(false)
	expected["--register-with-taints"] = fmt.Sprintf("node-role.kubernetes.io/customtaint=true:NoSchedule,%s", common.MasterNodeTaint)
	masterKubeletConfig = cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig
	for key, val := range masterKubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected masterProfile kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}
	cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--register-with-taints"] = fmt.Sprintf("node-role.kubernetes.io/customtaint=true:NoSchedule,node-role.kubernetes.io/customtaint2=true:NoSchedule,%s", common.MasterNodeTaint)
	cs.setKubeletConfig(false)
	expected["--register-with-taints"] = fmt.Sprintf("node-role.kubernetes.io/customtaint=true:NoSchedule,node-role.kubernetes.io/customtaint2=true:NoSchedule,%s", common.MasterNodeTaint)
	masterKubeletConfig = cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig
	for key, val := range masterKubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected masterProfile kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}
	cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--register-with-taints"] = fmt.Sprintf("%s,node-role.kubernetes.io/customtaint=true:NoSchedule,node-role.kubernetes.io/customtaint2=true:NoSchedule", common.MasterNodeTaint)
	cs.setKubeletConfig(false)
	expected["--register-with-taints"] = fmt.Sprintf("%s,node-role.kubernetes.io/customtaint=true:NoSchedule,node-role.kubernetes.io/customtaint2=true:NoSchedule", common.MasterNodeTaint)
	masterKubeletConfig = cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig
	for key, val := range masterKubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected masterProfile kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}
	expected["--register-with-taints"] = fmt.Sprintf("%s=true:NoSchedule", common.AADPodIdentityTaintKey)
	linuxProfileKubeletConfig := cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
	for key, val := range linuxProfileKubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected Linux agent profile kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}
	linuxProfileKubeletConfig["--register-with-taints"] = "node-role.kubernetes.io/customtaint=true:NoSchedule,node-role.kubernetes.io/customtaint2=true:NoSchedule"
	cs.setKubeletConfig(false)
	expected["--register-with-taints"] = fmt.Sprintf("node-role.kubernetes.io/customtaint=true:NoSchedule,node-role.kubernetes.io/customtaint2=true:NoSchedule,%s", fmt.Sprintf("%s=true:NoSchedule", common.AADPodIdentityTaintKey))
	for key, val := range linuxProfileKubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected Linux agent profile kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}
	cs.Properties.OrchestratorProfile.KubernetesConfig.ContainerRuntime = Containerd
	cs.setKubeletConfig(false)
	expected["--container-runtime"] = "remote" // todo: remove when default kubernetes version >= 1.27
	expected["--runtime-request-timeout"] = "15m"
	expected["--container-runtime-endpoint"] = "unix:///run/containerd/containerd.sock"
	for key, val := range linuxProfileKubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected Linux agent profile kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}
	delete(expected, "--register-with-taints")

	windowsProfileKubeletConfig := cs.Properties.AgentPoolProfiles[1].KubernetesConfig.KubeletConfig
	expected["--image-credential-provider-config"] = "c:\\k\\credential-provider\\credential-provider-config.yaml"
	expected["--image-credential-provider-bin-dir"] = "c:\\k\\credential-provider"
	expected["--pod-infra-container-image"] = "kubletwin/pause"
	expected["--kubeconfig"] = "c:\\k\\config"
	expected["--cloud-config"] = "c:\\k\\azure.json"
	expected["--cgroups-per-qos"] = "false"
	expected["--enforce-node-allocatable"] = "\"\"\"\""
	expected["--system-reserved"] = "memory=2Gi"
	expected["--client-ca-file"] = "c:\\k\\ca.crt"
	expected["--hairpin-mode"] = "promiscuous-bridge"
	expected["--image-pull-progress-deadline"] = "20m"
	expected["--resolv-conf"] = "\"\"\"\""
	expected["--eviction-hard"] = "\"\"\"\""
	delete(expected, "--pod-manifest-path")
	delete(expected, "--protect-kernel-defaults")
	delete(expected, "--tls-cert-file")
	delete(expected, "--tls-private-key-file")
	for key, val := range windowsProfileKubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected Windows agent profile kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}
	delete(expected, "--container-runtime") // todo: remove when default kubernetes version >= 1.27
	delete(expected, "--runtime-request-timeout")
	delete(expected, "--container-runtime-endpoint")

	// validate aad-pod-identity disabled scenario
	cs = CreateMockContainerService("testcluster", common.RationalizeReleaseAndVersion(Kubernetes, common.KubernetesDefaultRelease, "", false, false, false), 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.Addons = []KubernetesAddon{
		{
			Name:    common.AADPodIdentityAddonName,
			Enabled: to.BoolPtr(false),
		},
	}
	cs.setKubeletConfig(false)
	linuxProfileKubeletConfig = cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
	expected = getDefaultLinuxKubeletConfig(cs)
	for key, val := range linuxProfileKubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected Linux agent profile kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}

	cs = CreateMockContainerService("testcluster", common.RationalizeReleaseAndVersion(Kubernetes, common.KubernetesDefaultRelease, "", false, false, false), 3, 2, false)
	// check when ip-masq-agent is explicitly disabled in kubernetes config
	cs.Properties.OrchestratorProfile.KubernetesConfig.Addons = []KubernetesAddon{
		{
			Name:    common.IPMASQAgentAddonName,
			Enabled: to.BoolPtr(false),
		},
	}

	cs.setKubeletConfig(false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig

	// Removed kubelet --non-masquerade-cidr deprecated CLI flag (https://github.com/kubernetes/kubernetes/pull/97252)
	for key := range map[string]string{"--non-masquerade-cidr": DefaultKubernetesSubnet} {
		if _, ok := k[key]; ok {
			t.Fatal("got unexpected (removed) '--non-masquerade-cidr' kubelet config value")
		}
	}

	cs = CreateMockContainerService("testcluster", "", 3, 2, false)
	// TODO test all default overrides
	overrideVal := "/etc/override"
	cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig = map[string]string{
		"--image-credential-provider-config": overrideVal,
	}
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	for key, val := range map[string]string{"--image-credential-provider-config": overrideVal} {
		if k[key] != val {
			t.Fatalf("got unexpected kubelet config value for %s: %s, expected %s",
				key, k[key], val)
		}
	}

	cs = CreateMockContainerService("testcluster", "", 3, 2, false)
	cs.setKubeletConfig(false)
	kubeletConfig = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	expectedKeys := []string{
		"--authentication-token-webhook",
	}
	for _, key := range expectedKeys {
		if _, ok := kubeletConfig[key]; !ok {
			t.Fatalf("could not find expected kubelet config value for %s", key)
		}
	}
}

func getDefaultLinuxKubeletConfig(cs *ContainerService) map[string]string {
	k8sComponentsByVersionMap := GetK8sComponentsByVersionMap(&KubernetesConfig{KubernetesImageBaseType: common.KubernetesImageBaseTypeGCR})
	return map[string]string{
		"--address":                           "0.0.0.0",
		"--allow-privileged":                  "true", // validate that we delete this key for >= 1.15 clusters
		"--anonymous-auth":                    "false",
		"--authorization-mode":                "Webhook",
		"--authentication-token-webhook":      "true",
		"--image-credential-provider-config":  "/var/lib/kubelet/credential-provider-config.yaml",
		"--image-credential-provider-bin-dir": "/var/lib/kubelet/credential-provider",
		"--cadvisor-port":                     "", // Validate that we delete this key for >= 1.12 clusters
		"--cgroups-per-qos":                   "true",
		"--client-ca-file":                    "/etc/kubernetes/certs/ca.crt",
		"--cloud-provider":                    "azure",
		"--cloud-config":                      "/etc/kubernetes/azure.json",
		"--cluster-dns":                       DefaultKubernetesDNSServiceIP,
		"--cluster-domain":                    "cluster.local",
		"--enforce-node-allocatable":          "pods",
		"--event-qps":                         DefaultKubeletEventQPS,
		"--eviction-hard":                     DefaultKubernetesHardEvictionThreshold,
		"--image-gc-high-threshold":           strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"--image-gc-low-threshold":            strconv.Itoa(DefaultKubernetesGCLowThreshold),
		"--image-pull-progress-deadline":      "30m",
		"--keep-terminated-pod-volumes":       "false",
		"--kubeconfig":                        "/var/lib/kubelet/kubeconfig",
		"--max-pods":                          strconv.Itoa(DefaultKubernetesMaxPods),
		"--network-plugin":                    NetworkPluginKubenet,
		"--node-status-update-frequency":      k8sComponentsByVersionMap[cs.Properties.OrchestratorProfile.OrchestratorVersion]["nodestatusfreq"],
		"--non-masquerade-cidr":               DefaultNonMasqueradeCIDR,
		"--pod-manifest-path":                 "/etc/kubernetes/manifests",
		"--pod-infra-container-image":         cs.Properties.OrchestratorProfile.KubernetesConfig.MCRKubernetesImageBase + k8sComponentsByVersionMap[cs.Properties.OrchestratorProfile.OrchestratorVersion][common.PauseComponentName],
		"--pod-max-pids":                      strconv.Itoa(DefaultKubeletPodMaxPIDs),
		"--protect-kernel-defaults":           "true",
		"--read-only-port":                    "0",
		"--rotate-certificates":               "true",
		"--seccomp-default":                   "true",
		"--streaming-connection-idle-timeout": "5m",
		"--feature-gates":                     "ExecProbeTimeout=true,RotateKubeletServerCertificate=true",
		"--tls-cipher-suites":                 TLSStrongCipherSuitesKubelet,
		"--tls-cert-file":                     "/etc/kubernetes/certs/kubeletserver.crt",
		"--tls-private-key-file":              "/etc/kubernetes/certs/kubeletserver.key",
		"--v":                                 "2",
		"--volume-plugin-dir":                 "/etc/kubernetes/volumeplugins",
		"--healthz-port":                      DefaultKubeletHealthzPort,
	}
}

func TestKubeletConfigAzureStackDefaults(t *testing.T) {
	cs := CreateMockContainerService("testcluster", common.RationalizeReleaseAndVersion(Kubernetes, "", "", false, false, false), 3, 2, false)
	cs.Properties.CustomCloudProfile = &CustomCloudProfile{}
	winProfile := &AgentPoolProfile{}
	winProfile.Count = 1
	winProfile.Name = "agentpool2"
	winProfile.VMSize = "Standard_D2_v2"
	winProfile.OSType = Windows
	cs.Properties.AgentPoolProfiles = append(cs.Properties.AgentPoolProfiles, winProfile)
	cs.setKubeletConfig(false)
	k8sComponentsByVersionMap := GetK8sComponentsByVersionMap(&KubernetesConfig{KubernetesImageBaseType: common.KubernetesImageBaseTypeGCR})
	kubeletConfig := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	expected := map[string]string{
		"--address":                           "0.0.0.0",
		"--allow-privileged":                  "true", // validate that we delete this key for >= 1.15 clusters
		"--anonymous-auth":                    "false",
		"--authentication-token-webhook":      "true",
		"--authorization-mode":                "Webhook",
		"--image-credential-provider-config":  "/var/lib/kubelet/credential-provider-config.yaml",
		"--image-credential-provider-bin-dir": "/var/lib/kubelet/credential-provider",
		"--cadvisor-port":                     "", // Validate that we delete this key for >= 1.12 clusters
		"--cgroups-per-qos":                   "true",
		"--client-ca-file":                    "/etc/kubernetes/certs/ca.crt",
		"--cloud-config":                      "/etc/kubernetes/azure.json",
		"--cloud-provider":                    "azure",
		"--cluster-dns":                       DefaultKubernetesDNSServiceIP,
		"--cluster-domain":                    "cluster.local",
		"--enforce-node-allocatable":          "pods",
		"--event-qps":                         DefaultKubeletEventQPS,
		"--eviction-hard":                     DefaultKubernetesHardEvictionThreshold,
		"--feature-gates":                     "ExecProbeTimeout=true,RotateKubeletServerCertificate=true",
		"--image-gc-high-threshold":           strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"--image-gc-low-threshold":            strconv.Itoa(DefaultKubernetesGCLowThreshold),
		"--image-pull-progress-deadline":      "30m",
		"--keep-terminated-pod-volumes":       "false",
		"--kubeconfig":                        "/var/lib/kubelet/kubeconfig",
		"--max-pods":                          strconv.Itoa(DefaultKubernetesMaxPods),
		"--network-plugin":                    NetworkPluginKubenet,
		"--node-status-update-frequency":      DefaultAzureStackKubernetesNodeStatusUpdateFrequency,
		"--non-masquerade-cidr":               DefaultNonMasqueradeCIDR,
		"--pod-infra-container-image":         cs.Properties.OrchestratorProfile.KubernetesConfig.MCRKubernetesImageBase + k8sComponentsByVersionMap[cs.Properties.OrchestratorProfile.OrchestratorVersion][common.PauseComponentName],
		"--pod-manifest-path":                 "/etc/kubernetes/manifests",
		"--pod-max-pids":                      strconv.Itoa(DefaultKubeletPodMaxPIDs),
		"--protect-kernel-defaults":           "true",
		"--read-only-port":                    "0",
		"--register-with-taints":              common.MasterNodeTaint,
		"--rotate-certificates":               "true",
		"--seccomp-default":                   "true",
		"--streaming-connection-idle-timeout": "5m",
		"--tls-cert-file":                     "/etc/kubernetes/certs/kubeletserver.crt",
		"--tls-cipher-suites":                 TLSStrongCipherSuitesKubelet,
		"--tls-private-key-file":              "/etc/kubernetes/certs/kubeletserver.key",
		"--v":                                 "2",
		"--volume-plugin-dir":                 "/etc/kubernetes/volumeplugins",
		"--healthz-port":                      DefaultKubeletHealthzPort,
	}
	for key, val := range kubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}
	masterKubeletConfig := cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig
	for key, val := range masterKubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected masterProfile kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}
	linuxProfileKubeletConfig := cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
	for key, val := range linuxProfileKubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected Linux agent profile kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}

	windowsProfileKubeletConfig := cs.Properties.AgentPoolProfiles[1].KubernetesConfig.KubeletConfig
	expected["--image-credential-provider-config"] = "c:\\k\\credential-provider\\credential-provider-config.yaml"
	expected["--image-credential-provider-bin-dir"] = "c:\\k\\credential-provider"
	expected["--pod-infra-container-image"] = "kubletwin/pause"
	expected["--kubeconfig"] = "c:\\k\\config"
	expected["--cloud-config"] = "c:\\k\\azure.json"
	expected["--cgroups-per-qos"] = "false"
	expected["--enforce-node-allocatable"] = "\"\"\"\""
	expected["--system-reserved"] = "memory=2Gi"
	expected["--client-ca-file"] = "c:\\k\\ca.crt"
	expected["--hairpin-mode"] = "promiscuous-bridge"
	expected["--image-pull-progress-deadline"] = "20m"
	expected["--resolv-conf"] = "\"\"\"\""
	expected["--eviction-hard"] = "\"\"\"\""
	delete(expected, "--pod-manifest-path")
	delete(expected, "--protect-kernel-defaults")
	delete(expected, "--tls-cert-file")
	delete(expected, "--tls-private-key-file")
	for key, val := range windowsProfileKubeletConfig {
		if expected[key] != val {
			t.Fatalf("got unexpected Windows agent profile kubelet config value for %s: %s, expected %s",
				key, val, expected[key])
		}
	}

	cs = CreateMockContainerService("testcluster", common.RationalizeReleaseAndVersion(Kubernetes, common.KubernetesDefaultRelease, "", false, false, false), 3, 2, false)
	// check when ip-masq-agent is explicitly disabled in kubernetes config
	cs.Properties.OrchestratorProfile.KubernetesConfig.Addons = []KubernetesAddon{
		{
			Name:    common.IPMASQAgentAddonName,
			Enabled: to.BoolPtr(false),
		},
	}

	cs.setKubeletConfig(false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig

	// Removed kubelet --non-masquerade-cidr deprecated CLI flag (https://github.com/kubernetes/kubernetes/pull/97252)
	for key := range map[string]string{"--non-masquerade-cidr": DefaultKubernetesSubnet} {
		if _, ok := k[key]; ok {
			t.Fatal("got unexpected (removed) '--non-masquerade-cidr' kubelet config value")
		}
	}

	cs = CreateMockContainerService("testcluster", "", 3, 2, false)
	// TODO test all default overrides
	overrideVal := "/etc/override"
	cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig = map[string]string{
		"--image-credential-provider-config": overrideVal,
	}
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	for key, val := range map[string]string{"--image-credential-provider-config": overrideVal} {
		if k[key] != val {
			t.Fatalf("got unexpected kubelet config value for %s: %s, expected %s",
				key, k[key], val)
		}
	}

	cs = CreateMockContainerService("testcluster", common.RationalizeReleaseAndVersion(Kubernetes, "", "", false, false, false), 3, 2, false)
	cs.setKubeletConfig(false)
	kubeletConfig = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	expectedKeys := []string{
		"--authentication-token-webhook",
	}
	for _, key := range expectedKeys {
		if _, ok := kubeletConfig[key]; !ok {
			t.Fatalf("could not find expected kubelet config value for %s", key)
		}
	}
}

func TestKubeletConfigDefaultsRemovals(t *testing.T) {
	cs := CreateMockContainerService("testcluster", common.RationalizeReleaseAndVersion(Kubernetes, "", "", false, false, false), 3, 2, false)
	poolProfile := &AgentPoolProfile{}
	poolProfile.Count = 1
	poolProfile.Name = "agentpool2"
	poolProfile.VMSize = "Standard_D2_v2"
	poolProfile.OSType = Linux
	cs.Properties.AgentPoolProfiles = append(cs.Properties.AgentPoolProfiles, poolProfile)
	cs.setKubeletConfig(false)
	kubeletConfig := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	unexpected := []string{
		"--cadvisor-port",
	}
	for _, key := range unexpected {
		if _, ok := kubeletConfig[key]; ok {
			t.Fatalf("got unexpected kubelet config value for %s, expected it not to be present",
				key)
		}
	}
	cs = CreateMockContainerService("testcluster", "", 3, 2, false)
	cs.Properties.AgentPoolProfiles = append(cs.Properties.AgentPoolProfiles, poolProfile)
	cs.setKubeletConfig(false)
	kubeletConfig = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	unexpected = []string{
		"--allow-privileged",
		"--cadvisor-port",
	}
	for _, key := range unexpected {
		if _, ok := kubeletConfig[key]; ok {
			t.Fatalf("got unexpected kubelet config value for %s, expected it not to be present",
				key)
		}
	}
}

func TestKubeletConfigUseCloudControllerManager(t *testing.T) {
	// Test UseCloudControllerManager = true
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager = to.BoolPtr(true)
	cs.setKubeletConfig(false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--cloud-provider"] != "external" {
		t.Fatalf("got unexpected '--cloud-provider' kubelet config value for UseCloudControllerManager=true: %s",
			k["--cloud-provider"])
	}

	// Test UseCloudControllerManager = false
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager = to.BoolPtr(false)
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--cloud-provider"] != "azure" {
		t.Fatalf("got unexpected '--cloud-provider' kubelet config value for UseCloudControllerManager=false: %s",
			k["--cloud-provider"])
	}

}

func TestKubeletConfigCloudConfig(t *testing.T) {
	// Test default value and custom value for --cloud-config
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setKubeletConfig(false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--cloud-config"] != "/etc/kubernetes/azure.json" {
		t.Fatalf("got unexpected '--cloud-config' kubelet config default value: %s",
			k["--cloud-config"])
	}

	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig["--cloud-config"] = "custom.json"
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--cloud-config"] != "custom.json" {
		t.Fatalf("got unexpected '--cloud-config' kubelet config default value: %s",
			k["--cloud-config"])
	}
}

func TestKubeletConfigAzureContainerRegistryConfig(t *testing.T) {
	// Test default value and custom value for --image-credential-provider-config
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setKubeletConfig(false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--image-credential-provider-config"] != "/var/lib/kubelet/credential-provider-config.yaml" {
		t.Fatalf("got unexpected '--image-credential-provider-config' kubelet config default value: %s",
			k["--image-credential-provider-config"])
	}
	if k["--image-credential-provider-bin-dir"] != "/var/lib/kubelet/credential-provider" {
		t.Fatalf("got unexpected '--image-credential-provider-bin-dir' kubelet config default value: %s",
			k["--image-credential-provider-bin-dir"])
	}

	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig["--image-credential-provider-config"] = "custom.json"
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--image-credential-provider-config"] != "custom.json" {
		t.Fatalf("got unexpected '--image-credential-provider-config' kubelet config default value: %s",
			k["--image-credential-provider-config"])
	}
	if k["--image-credential-provider-bin-dir"] != "/var/lib/kubelet/credential-provider" {
		t.Fatalf("got unexpected '--image-credential-provider-bin-dir' kubelet config default value: %s",
			k["--image-credential-provider-bin-dir"])
	}
}

func TestKubeletConfigNetworkPlugin(t *testing.T) {
	// Test NetworkPlugin = "kubenet"
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = NetworkPluginKubenet
	cs.setKubeletConfig(false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if _, ok := k["--network-plugin"]; ok {
		t.Fatalf("got unexpected (removed) '--network-plugin' kubelet config value")
	}

	// Test NetworkPlugin = "azure"
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = NetworkPluginAzure
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if _, ok := k["--network-plugin"]; ok {
		t.Fatalf("got unexpected (removed) '--network-plugin' kubelet config value")
	}
}

func TestKubeletConfigEnableSecureKubelet(t *testing.T) {
	// Test EnableSecureKubelet = true
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableSecureKubelet = to.BoolPtr(true)
	cs.setKubeletConfig(false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	ka := cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
	for _, kubernetesConfig := range []map[string]string{k, ka} {
		if kubernetesConfig["--anonymous-auth"] != "false" {
			t.Fatalf("got unexpected '--anonymous-auth' kubelet config value for EnableSecureKubelet=true: %s",
				kubernetesConfig["--anonymous-auth"])
		}
		if kubernetesConfig["--authorization-mode"] != "Webhook" {
			t.Fatalf("got unexpected '--authorization-mode' kubelet config value for EnableSecureKubelet=true: %s",
				kubernetesConfig["--authorization-mode"])
		}
		if kubernetesConfig["--client-ca-file"] != "/etc/kubernetes/certs/ca.crt" {
			t.Fatalf("got unexpected '--client-ca-file' kubelet config value for EnableSecureKubelet=true: %s",
				kubernetesConfig["--client-ca-file"])
		}
	}

	// Test EnableSecureKubelet = false
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableSecureKubelet = to.BoolPtr(false)
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	ka = cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
	for _, kubernetesConfig := range []map[string]string{k, ka} {
		for _, key := range []string{"--anonymous-auth", "--client-ca-file"} {
			if _, ok := kubernetesConfig[key]; ok {
				t.Fatalf("got unexpected '%s' kubelet config value for EnableSecureKubelet=false: %s",
					key, kubernetesConfig[key])
			}
		}
	}

	// Test default (EnableSecureKubelet = false) for Windows
	cs = CreateMockContainerService("testcluster", "", 3, 1, false)
	p := GetK8sDefaultProperties(true)
	cs.Properties = p
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableSecureKubelet = to.BoolPtr(false)
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	ka = cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
	for _, kubernetesConfig := range []map[string]string{k, ka} {
		for _, key := range []string{"--anonymous-auth", "--client-ca-file"} {
			if _, ok := kubernetesConfig[key]; ok {
				t.Fatalf("got unexpected '%s' kubelet config value for EnableSecureKubelet=false: %s",
					key, kubernetesConfig[key])
			}
		}
	}

	// Test explicit EnableSecureKubelet = false for Windows
	cs = CreateMockContainerService("testcluster", "", 3, 1, false)
	p = GetK8sDefaultProperties(true)
	cs.Properties = p
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableSecureKubelet = to.BoolPtr(false)
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	ka = cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
	for _, kubernetesConfig := range []map[string]string{k, ka} {
		for _, key := range []string{"--anonymous-auth", "--client-ca-file"} {
			if _, ok := kubernetesConfig[key]; ok {
				t.Fatalf("got unexpected '%s' kubelet config value for EnableSecureKubelet=false: %s",
					key, kubernetesConfig[key])
			}
		}
	}

	// Test EnableSecureKubelet = true for Windows
	cs = CreateMockContainerService("testcluster", "", 3, 1, false)
	p = GetK8sDefaultProperties(true)
	cs.Properties = p
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableSecureKubelet = to.BoolPtr(true)
	cs.setKubeletConfig(false)
	kubernetesConfig := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	kubernetesConfigWindowsAgentPool := cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
	if kubernetesConfig["--anonymous-auth"] != "false" {
		t.Fatalf("got unexpected '--anonymous-auth' kubelet config value for EnableSecureKubelet=true: %s",
			kubernetesConfig["--anonymous-auth"])
	}
	if kubernetesConfig["--client-ca-file"] != "/etc/kubernetes/certs/ca.crt" {
		t.Fatalf("got unexpected '--client-ca-file' kubelet config value for EnableSecureKubelet=true: %s",
			kubernetesConfig["--client-ca-file"])
	}
	if kubernetesConfigWindowsAgentPool["--anonymous-auth"] != "false" {
		t.Fatalf("got unexpected '--anonymous-auth' kubelet config value for EnableSecureKubelet=true: %s",
			kubernetesConfig["--anonymous-auth"])
	}
	if kubernetesConfigWindowsAgentPool["--client-ca-file"] != "c:\\k\\ca.crt" {
		t.Fatalf("got unexpected '--client-ca-file' kubelet config value for EnableSecureKubelet=true: %s",
			kubernetesConfig["--client-ca-file"])
	}
}

func TestKubeletMaxPods(t *testing.T) {
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = NetworkPluginAzure
	cs.setKubeletConfig(false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--max-pods"] != strconv.Itoa(DefaultKubernetesMaxPodsVNETIntegrated) {
		t.Fatalf("got unexpected '--max-pods' kubelet config value for NetworkPolicy=%s: %s",
			NetworkPluginAzure, k["--max-pods"])
	}

	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = NetworkPluginKubenet
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--max-pods"] != strconv.Itoa(DefaultKubernetesMaxPods) {
		t.Fatalf("got unexpected '--max-pods' kubelet config value for NetworkPolicy=%s: %s",
			NetworkPluginKubenet, k["--max-pods"])
	}

	// Test that user-overrides for --max-pods work as intended
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = NetworkPluginKubenet
	cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig["--max-pods"] = "99"
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--max-pods"] != "99" {
		t.Fatalf("got unexpected '--max-pods' kubelet config value for NetworkPolicy=%s: %s",
			NetworkPluginKubenet, k["--max-pods"])
	}

	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = NetworkPluginAzure
	cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig["--max-pods"] = "99"
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--max-pods"] != "99" {
		t.Fatalf("got unexpected '--max-pods' kubelet config value for NetworkPolicy=%s: %s",
			NetworkPluginKubenet, k["--max-pods"])
	}
}

func TestKubeletCalico(t *testing.T) {
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy = NetworkPolicyCalico
	cs.setKubeletConfig(false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if _, ok := k["--network-plugin"]; ok {
		t.Fatal("got unexpected (removed) '--network-plugin' kubelet config value")
	}
}

func TestKubeletIPMasqAgentEnabledOrDisabled(t *testing.T) {
	subnet := "172.16.0.0/16"
	// MasterIPMasqAgent disabled, --non-masquerade-cidr should be subnet
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	b := false
	cs.Properties.OrchestratorProfile.KubernetesConfig = &KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    common.IPMASQAgentAddonName,
				Enabled: &b,
			},
		},
	}
	cs.Properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet = subnet
	cs.setKubeletConfig(false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	// Removed kubelet --non-masquerade-cidr deprecated CLI flag (https://github.com/kubernetes/kubernetes/pull/97252)
	if _, ok := k["--non-masquerade-cidr"]; ok {
		t.Fatal("got unexpected (removed) '--non-masquerade-cidr' kubelet config value")
	}

	// MasterIPMasqAgent enabled, --non-masquerade-cidr should be 0.0.0.0/0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	b = true
	cs.Properties.OrchestratorProfile.KubernetesConfig = &KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    common.IPMASQAgentAddonName,
				Enabled: &b,
			},
		},
	}
	cs.Properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet = subnet
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if _, ok := k["--non-masquerade-cidr"]; ok {
		t.Fatalf("got unexpected (removed) '--non-masquerade-cidr' kubelet config value")
	}

	// No ip-masq-agent addon configuration specified, --non-masquerade-cidr should be 0.0.0.0/0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if _, ok := k["--non-masquerade-cidr"]; ok {
		t.Fatalf("got unexpected (removed) '--non-masquerade-cidr' kubelet config value")
	}
}

func TestEnforceNodeAllocatable(t *testing.T) {
	// Validate default
	cs := CreateMockContainerService("testcluster", "", 3, 2, false)
	cs.setKubeletConfig(false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--enforce-node-allocatable"] != "pods" {
		t.Fatalf("got unexpected '--enforce-node-allocatable' kubelet config value %s, the expected value is %s",
			k["--enforce-node-allocatable"], "pods")
	}

	// Validate that --enforce-node-allocatable is overridable
	cs = CreateMockContainerService("testcluster", "", 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig = &KubernetesConfig{
		KubeletConfig: map[string]string{
			"--enforce-node-allocatable": "kube-reserved/system-reserved",
		},
	}
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--enforce-node-allocatable"] != "kube-reserved/system-reserved" {
		t.Fatalf("got unexpected '--enforce-node-allocatable' kubelet config value %s, the expected value is %s",
			k["--enforce-node-allocatable"], "kube-reserved/system-reserved")
	}
}

func TestProtectKernelDefaults(t *testing.T) {
	// Validate default
	cs := CreateMockContainerService("testcluster", "", 3, 2, false)
	_, err := cs.SetPropertiesDefaults(PropertiesDefaultsParams{
		IsScale:    false,
		IsUpgrade:  false,
		PkiKeySize: helpers.DefaultPkiKeySize,
	})
	if err != nil {
		t.Error(err)
	}
	km := cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig
	if km["--protect-kernel-defaults"] != "true" {
		t.Fatalf("got unexpected '--protect-kernel-defaults' kubelet config value %s, the expected value is %s",
			km["--protect-kernel-defaults"], "true")
	}
	ka := cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
	if ka["--protect-kernel-defaults"] != "true" {
		t.Fatalf("got unexpected '--protect-kernel-defaults' kubelet config value %s, the expected value is %s",
			ka["--protect-kernel-defaults"], "true")
	}

	// Validate that --protect-kernel-defaults is "true" by default for relevant distros
	for _, distro := range DistroValues {
		switch distro {
		case AKSUbuntu1604, AKSUbuntu1804, AKSUbuntu2004:
			cs = CreateMockContainerService("testcluster", "", 3, 2, false)
			cs.Properties.MasterProfile.Distro = distro
			cs.Properties.AgentPoolProfiles[0].Distro = distro
			_, err = cs.SetPropertiesDefaults(PropertiesDefaultsParams{
				IsScale:    false,
				IsUpgrade:  false,
				PkiKeySize: helpers.DefaultPkiKeySize,
			})
			if err != nil {
				t.Error(err)
			}
			km = cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig
			if km["--protect-kernel-defaults"] != "true" {
				t.Fatalf("got unexpected '--protect-kernel-defaults' kubelet config value %s, the expected value is %s",
					km["--protect-kernel-defaults"], "true")
			}
			ka = cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
			if ka["--protect-kernel-defaults"] != "true" {
				t.Fatalf("got unexpected '--protect-kernel-defaults' kubelet config value %s, the expected value is %s",
					ka["--protect-kernel-defaults"], "true")
			}

		// Validate that --protect-kernel-defaults is not enabled for relevant distros
		case Ubuntu, Ubuntu1804, Ubuntu1804Gen2, Ubuntu2004, ACC1604, Flatcar:
			cs = CreateMockContainerService("testcluster", "", 3, 2, false)
			if distro == Flatcar {
				cs.Properties.MasterProfile.Distro = Ubuntu1804
			} else {
				cs.Properties.MasterProfile.Distro = distro
			}
			cs.Properties.AgentPoolProfiles[0].Distro = distro
			_, err = cs.SetPropertiesDefaults(PropertiesDefaultsParams{
				IsScale:    false,
				IsUpgrade:  false,
				PkiKeySize: helpers.DefaultPkiKeySize,
			})
			if err != nil {
				t.Error(err)
			}
			km = cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig
			if _, ok := km["--protect-kernel-defaults"]; ok {
				t.Fatalf("got unexpected '--protect-kernel-defaults' kubelet config value %s",
					km["--protect-kernel-defaults"])
			}
			ka = cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
			if _, ok := ka["--protect-kernel-defaults"]; ok {
				t.Fatalf("got unexpected '--protect-kernel-defaults' kubelet config value %s",
					ka["--protect-kernel-defaults"])
			}
		}
	}

	// Validate that --protect-kernel-defaults is not enabled for Windows
	cs = CreateMockContainerService("testcluster", "", 3, 2, false)
	cs.Properties.MasterProfile.Distro = AKSUbuntu1604
	cs.Properties.AgentPoolProfiles[0].OSType = Windows
	_, err = cs.SetPropertiesDefaults(PropertiesDefaultsParams{
		IsScale:    false,
		IsUpgrade:  false,
		PkiKeySize: helpers.DefaultPkiKeySize,
	})
	if err != nil {
		t.Error(err)
	}
	km = cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig
	if km["--protect-kernel-defaults"] != "true" {
		t.Fatalf("got unexpected '--protect-kernel-defaults' kubelet config value %s, the expected value is %s",
			km["--protect-kernel-defaults"], "true")
	}
	ka = cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
	if _, ok := ka["--protect-kernel-defaults"]; ok {
		t.Fatalf("got unexpected '--protect-kernel-defaults' kubelet config value %s",
			ka["--protect-kernel-defaults"])
	}

	// Validate that --protect-kernel-defaults is overridable
	for _, distro := range DistroValues {
		switch distro {
		case Ubuntu, Ubuntu1804, Ubuntu1804Gen2, Ubuntu2004, AKSUbuntu1604, AKSUbuntu1804, AKSUbuntu2004:
			cs = CreateMockContainerService("testcluster", "", 3, 2, false)
			cs.Properties.MasterProfile.Distro = "ubuntu"
			cs.Properties.AgentPoolProfiles[0].Distro = "ubuntu"
			cs.Properties.OrchestratorProfile.KubernetesConfig = &KubernetesConfig{
				KubeletConfig: map[string]string{
					"--protect-kernel-defaults": "false",
				},
			}
			_, err = cs.SetPropertiesDefaults(PropertiesDefaultsParams{
				IsScale:    false,
				IsUpgrade:  false,
				PkiKeySize: helpers.DefaultPkiKeySize,
			})
			if err != nil {
				t.Error(err)
			}
			km = cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig
			if km["--protect-kernel-defaults"] != "false" {
				t.Fatalf("got unexpected '--protect-kernel-defaults' kubelet config value %s, the expected value is %s",
					km["--protect-kernel-defaults"], "false")
			}
			ka = cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig
			if ka["--protect-kernel-defaults"] != "false" {
				t.Fatalf("got unexpected '--protect-kernel-defaults' kubelet config value %s, the expected value is %s",
					ka["--protect-kernel-defaults"], "false")
			}
		}
	}
}

func TestStaticWindowsConfig(t *testing.T) {
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 1, false)
	p := GetK8sDefaultProperties(true)
	cs.Properties = p
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableSecureKubelet = to.BoolPtr(true)

	// Start with copy of Linux config
	staticLinuxKubeletConfig := map[string]string{
		"--address":                     "0.0.0.0",
		"--allow-privileged":            "true",
		"--anonymous-auth":              "false",
		"--authorization-mode":          "Webhook",
		"--client-ca-file":              "/etc/kubernetes/certs/ca.crt",
		"--pod-manifest-path":           "/etc/kubernetes/manifests",
		"--cluster-dns":                 cs.Properties.OrchestratorProfile.KubernetesConfig.DNSServiceIP,
		"--cgroups-per-qos":             "true",
		"--kubeconfig":                  "/var/lib/kubelet/kubeconfig",
		"--keep-terminated-pod-volumes": "false",
	}
	expected := make(map[string]string)
	for key, val := range staticLinuxKubeletConfig {
		if key != "--pod-manifest-path" {
			expected[key] = val
		}
	}

	// Add Windows-specific overrides
	// Eventually paths should not be hardcoded here. They should be relative to $global:KubeDir in the PowerShell script
	expected["--image-credential-provider-config"] = "c:\\k\\credential-provider\\credential-provider-config.yaml"
	expected["--image-credential-provider-bin-dir"] = "c:\\k\\credential-provider"
	expected["--pod-infra-container-image"] = "kubletwin/pause"
	expected["--kubeconfig"] = "c:\\k\\config"
	expected["--cloud-config"] = "c:\\k\\azure.json"
	expected["--cgroups-per-qos"] = "false"
	expected["--enforce-node-allocatable"] = "\"\"\"\""
	expected["--system-reserved"] = "memory=2Gi"
	expected["--client-ca-file"] = "c:\\k\\ca.crt"
	expected["--hairpin-mode"] = "promiscuous-bridge"
	expected["--image-pull-progress-deadline"] = "20m"
	expected["--resolv-conf"] = "\"\"\"\""
	expected["--eviction-hard"] = "\"\"\"\""

	cs.setKubeletConfig(false)
	for _, profile := range cs.Properties.AgentPoolProfiles {
		if profile.OSType == Windows {
			for key, val := range expected {
				if val != profile.KubernetesConfig.KubeletConfig[key] {
					t.Fatalf("got unexpected '%s' kubelet config value, expected %s, got %s",
						key, val, profile.KubernetesConfig.KubeletConfig[key])
				}
			}
		}
	}
	if _, ok := cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--register-with-taints"]; ok {
		t.Fatalf("got unexpected --register-with-taints kubelet config with no Linux node pools")
	}
}

func TestKubeletRotateCertificates(t *testing.T) {
	cs := CreateMockContainerService("testcluster", common.RationalizeReleaseAndVersion(Kubernetes, "1.18", "", false, false, false), 3, 2, false)
	cs.setKubeletConfig(false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--rotate-certificates"] != "true" {
		t.Fatalf("got unexpected '--rotate-certificates' kubelet config value for k8s version %s: %s",
			cs.Properties.OrchestratorProfile.OrchestratorVersion, k["--rotate-certificates"])
	}

	// Test user-override
	cs = CreateMockContainerService("testcluster", common.RationalizeReleaseAndVersion(Kubernetes, "1.18", "", false, false, false), 3, 2, false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	k["--rotate-certificates"] = "false"
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--rotate-certificates"] != "false" {
		t.Fatalf("got unexpected '--rotate-certificates' kubelet config value despite override value %s: %s",
			"false", k["--rotate-certificates"])
	}
}

func TestKubeletConfigFeatureGates(t *testing.T) {
	// test user-overrides
	cs := CreateMockContainerService("testcluster", "", 3, 2, false)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	k["--feature-gates"] = "DynamicKubeletConfig=true,ExecProbeTimeout=false"
	cs.setKubeletConfig(false)
	if k["--feature-gates"] != "DynamicKubeletConfig=true,ExecProbeTimeout=false,RotateKubeletServerCertificate=true" {
		t.Fatalf("got unexpected '--feature-gates' kubelet config value for \"--feature-gates\": \"\": %s",
			k["--feature-gates"])
	}

	// test 1.25.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.25.0"
	cs.setKubeletConfig(false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--feature-gates"] != "ExecProbeTimeout=true,PodSecurity=true,RotateKubeletServerCertificate=true" {
		t.Fatalf("got unexpected '--feature-gates' kubelet config value for k8s v%s: %s",
			"1.25.0", k["--feature-gates"])
	}

	// test user-overrides, removal of VolumeSnapshotDataSource for k8s versions >= 1.22
	cs = CreateMockContainerService("testcluster", common.RationalizeReleaseAndVersion(Kubernetes, "1.22", "", false, false, false), 3, 2, false)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	k["--feature-gates"] = "VolumeSnapshotDataSource=true"
	cs.setKubeletConfig(false)
	if k["--feature-gates"] != "ExecProbeTimeout=true,RotateKubeletServerCertificate=true" {
		t.Fatalf("got unexpected '--feature-gates' kubelet config value for \"--feature-gates\": \"\": %s",
			k["--feature-gates"])
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.27.0"
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	k["--feature-gates"] = "ControllerManagerLeaderMigration=true,ExpandCSIVolumes=true,ExpandInUsePersistentVolumes=true,ExpandPersistentVolumes=true,CSIInlineVolume=true,CSIMigration=true,CSIMigrationAzureDisk=true,DaemonSetUpdateSurge=true,EphemeralContainers=true,IdentifyPodOS=true,LocalStorageCapacityIsolation=true,NetworkPolicyEndPort=true,StatefulSetMinReadySeconds=true"
	cs.setKubeletConfig(false)
	if k["--feature-gates"] != "ExecProbeTimeout=true,PodSecurity=true,RotateKubeletServerCertificate=true" {
		t.Fatalf("got unexpected '--feature-gates' kubelet config value for \"--feature-gates\": \"\": %s",
			k["--feature-gates"])
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.28
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.28.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	featuregate128 := "AdvancedAuditing=true,CSIMigrationGCE=true,CSIStorageCapacity=true,DelegateFSGroupToCSIDriver=true,DevicePlugins=true,DisableAcceleratorUsageMetrics=true,DryRun=true,EndpointSliceTerminatingCondition=true,KubeletCredentialProviders=true,MixedProtocolLBService=true,NetworkPolicyStatus=true,PodHasNetworkCondition=true,PodSecurity=true,ServiceIPStaticSubrange=true,ServiceInternalTrafficPolicy=true,UserNamespacesStatelessPodsSupport=true,WindowsHostProcessContainers=true"
	k["--feature-gates"] = featuregate128
	featuregate128Sanitized := "ExecProbeTimeout=true,RotateKubeletServerCertificate=true"
	cs.setKubeletConfig(false)
	if k["--feature-gates"] != featuregate128Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n kubelet config original value  %s \n, expected sanitized value: %s \n, actual sanitized value: %s \n ",
			"1.28.0", featuregate128, k["--feature-gates"], featuregate128Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.27
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.27.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	k["--feature-gates"] = featuregate128
	featuregate127Sanitized := "AdvancedAuditing=true,CSIMigrationGCE=true,CSIStorageCapacity=true,DelegateFSGroupToCSIDriver=true,DevicePlugins=true,DisableAcceleratorUsageMetrics=true,DryRun=true,EndpointSliceTerminatingCondition=true,ExecProbeTimeout=true,KubeletCredentialProviders=true,MixedProtocolLBService=true,NetworkPolicyStatus=true,PodHasNetworkCondition=true,PodSecurity=true,RotateKubeletServerCertificate=true,ServiceIPStaticSubrange=true,ServiceInternalTrafficPolicy=true,UserNamespacesStatelessPodsSupport=true,WindowsHostProcessContainers=true"
	cs.setKubeletConfig(false)
	if k["--feature-gates"] != featuregate127Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n kubelet config original value  %s \n, expected sanitized value: %s \n, actual sanitized value: %s \n ",
			"1.27.0", featuregate128, k["--feature-gates"], featuregate127Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.29
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.29.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig = make(map[string]string)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	featuregate129 := "CSIMigrationvSphere=true,CronJobTimeZone=true,DownwardAPIHugePages=true,GRPCContainerProbe=true,JobMutableNodeSchedulingDirectives=true,JobTrackingWithFinalizers=true,LegacyServiceAccountTokenNoAutoGeneration=true,OpenAPIV3=true,ProbeTerminationGracePeriod=true,RetroactiveDefaultStorageClass=true,SeccompDefault=true,TopologyManager=true"
	k["--feature-gates"] = featuregate129
	featuregate129Sanitized := "ExecProbeTimeout=true,RotateKubeletServerCertificate=true"
	cs.setKubeletConfig(false)
	if k["--feature-gates"] != featuregate129Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n kubelet config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.29.0", featuregate129, k["--feature-gates"], featuregate129Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.29
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.28.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig = make(map[string]string)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	k["--feature-gates"] = featuregate129
	featuregate128Sanitized = "CSIMigrationvSphere=true,CronJobTimeZone=true,DownwardAPIHugePages=true,ExecProbeTimeout=true,GRPCContainerProbe=true,JobMutableNodeSchedulingDirectives=true,JobTrackingWithFinalizers=true,LegacyServiceAccountTokenNoAutoGeneration=true,OpenAPIV3=true,ProbeTerminationGracePeriod=true,RetroactiveDefaultStorageClass=true,RotateKubeletServerCertificate=true,SeccompDefault=true,TopologyManager=true"
	cs.setKubeletConfig(false)
	if k["--feature-gates"] != featuregate128Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n kubelet config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.28.0", featuregate129, k["--feature-gates"], featuregate128Sanitized)
	}

	// test user-overrides, removal of feature gates for k8s versions >= 1.30
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.30.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	featuregate130 := "APISelfSubjectReview=true,CSIMigrationAzureFile=true,ExpandedDNSConfig=true,ExperimentalHostUserNamespaceDefaulting=true,IPTablesOwnershipCleanup=true,KubeletPodResources=true,KubeletPodResourcesGetAllocatable=true,LegacyServiceAccountTokenTracking=true,MinimizeIPTablesRestore=true,ProxyTerminatingEndpoints=true,RemoveSelfLink=true,SecurityContextDeny=true"
	k["--feature-gates"] = featuregate130
	featuregate130Sanitized := ""
	cs.setAPIServerConfig()
	if k["--feature-gates"] != featuregate130Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.30.0", featuregate130, k["--feature-gates"], featuregate130Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.30
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.29.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig = make(map[string]string)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.APIServerConfig
	k["--feature-gates"] = featuregate130
	featuregate129Sanitized = featuregate130
	cs.setAPIServerConfig()
	if k["--feature-gates"] != featuregate129Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n API server config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.29.0", featuregate130, k["--feature-gates"], featuregate129Sanitized)
	}
}

func TestKubeletStrongCipherSuites(t *testing.T) {
	// Test allowed versions
	for _, version := range []string{"1.17.0", "1.18.0"} {
		cs := CreateMockContainerService("testcluster", version, 3, 2, false)
		cs.setKubeletConfig(false)
		k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
		if k["--tls-cipher-suites"] != TLSStrongCipherSuitesKubelet {
			t.Fatalf("got unexpected default value for '--tls-cipher-suites' kubelet config for Kubernetes version %s: %s",
				version, k["--tls-cipher-suites"])
		}
	}

	allSuites := "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_RC4_128_SHA,TLS_RSA_WITH_3DES_EDE_CBC_SHA,TLS_RSA_WITH_AES_128_CBC_SHA,TLS_RSA_WITH_AES_128_CBC_SHA256,TLS_RSA_WITH_AES_128_GCM_SHA256,TLS_RSA_WITH_AES_256_CBC_SHA,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_RC4_128_SHA"
	// Test user-override
	for _, version := range []string{"1.17.0", "1.18.0"} {
		cs := CreateMockContainerService("testcluster", version, 3, 2, false)
		cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig = map[string]string{
			"--tls-cipher-suites": allSuites,
		}
		cs.setKubeletConfig(false)
		k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
		if k["--tls-cipher-suites"] != allSuites {
			t.Fatalf("got unexpected default value for '--tls-cipher-suites' API server config for Kubernetes version %s: %s",
				version, k["--tls-cipher-suites"])
		}
	}
}

func TestSupportPodPidsLimitFeatureGate(t *testing.T) {
	cases := []struct {
		name                                   string
		cs                                     *ContainerService
		isUpgrade                              bool
		expectedPodMaxPids                     string
		expectedSupportPodPidsLimitFeatureGate bool
	}{
		{
			name: "no --pod-max-pids defined",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig:    &KubernetesConfig{},
					},
				},
			},
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined, no SupportPodPidsLimit defined",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids": "100",
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "100",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "no --pod-max-pids defined, SupportPodPidsLimit=false",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--feature-gates": "SupportPodPidsLimit=false",
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "no --pod-max-pids defined, SupportPodPidsLimit=true",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--feature-gates": "SupportPodPidsLimit=true",
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined, SupportPodPidsLimit=false",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids":  "100",
								"--feature-gates": "SupportPodPidsLimit=false",
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "100",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined, SupportPodPidsLimit=true",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids":  "100",
								"--feature-gates": "SupportPodPidsLimit=true",
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "100",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "no --pod-max-pids defined, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig:    &KubernetesConfig{},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids": "100",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "no --pod-max-pids defined, SupportPodPidsLimit=false, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--feature-gates": "SupportPodPidsLimit=false",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "no --pod-max-pids defined, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--feature-gates": "SupportPodPidsLimit=true",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined, no SupportPodPidsLimit defined, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids": "-100",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined, SupportPodPidsLimit=false, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids":  "100",
								"--feature-gates": "SupportPodPidsLimit=false",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids":  "100",
								"--feature-gates": "SupportPodPidsLimit=true",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "100",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined as 100a, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids":  "100a",
								"--feature-gates": "SupportPodPidsLimit=true",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			c.cs.setKubeletConfig(c.isUpgrade)
			podMaxPids := c.cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig["--pod-max-pids"]
			if podMaxPids != c.expectedPodMaxPids {
				t.Fatalf("expected --pod-max-pids be equal to %s, got %s", c.expectedPodMaxPids, podMaxPids)
			}
			featureGates := c.cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig["--feature-gates"]
			hasSupportPodPidsLimitFeatureGate := strings.Contains(featureGates, "SupportPodPidsLimit=true")
			if hasSupportPodPidsLimitFeatureGate != c.expectedSupportPodPidsLimitFeatureGate {
				t.Fatalf("expected SupportPodPidsLimit=true presence in --feature gates to be %t, got %t", c.expectedSupportPodPidsLimitFeatureGate, hasSupportPodPidsLimitFeatureGate)
			}
		})
	}

}

func TestSupportPodPidsLimitFeatureGateInMasterProfile(t *testing.T) {
	cases := []struct {
		name                                   string
		cs                                     *ContainerService
		isUpgrade                              bool
		expectedPodMaxPids                     string
		expectedSupportPodPidsLimitFeatureGate bool
	}{
		{
			name: "no --pod-max-pids defined",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{},
					},
				},
			},
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined, no SupportPodPidsLimit defined",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids": "100",
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "100",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "no --pod-max-pids defined, SupportPodPidsLimit=false",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--feature-gates": "SupportPodPidsLimit=false",
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "no --pod-max-pids defined, SupportPodPidsLimit=true",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--feature-gates": "SupportPodPidsLimit=true",
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined, SupportPodPidsLimit=false",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids":  "100",
								"--feature-gates": "SupportPodPidsLimit=false",
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "100",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined, SupportPodPidsLimit=true",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids":  "100",
								"--feature-gates": "SupportPodPidsLimit=true",
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "100",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "no --pod-max-pids defined, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "no --pod-max-pids defined, SupportPodPidsLimit=false, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--feature-gates": "SupportPodPidsLimit=false",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "no --pod-max-pids defined, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--feature-gates": "SupportPodPidsLimit=true",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined as 100, no SupportPodPidsLimit defined, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids": "100",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined as -100, no SupportPodPidsLimit defined, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids": "-100",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined, SupportPodPidsLimit=false, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids":  "100",
								"--feature-gates": "SupportPodPidsLimit=false",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined as 100, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids":  "100",
								"--feature-gates": "SupportPodPidsLimit=true",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "100",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined as 0, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids":  "0",
								"--feature-gates": "SupportPodPidsLimit=true",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined as -1, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids":  "-1",
								"--feature-gates": "SupportPodPidsLimit=true",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined as 100a, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-max-pids":  "100a",
								"--feature-gates": "SupportPodPidsLimit=true",
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			c.cs.setKubeletConfig(c.isUpgrade)
			podMaxPids := c.cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--pod-max-pids"]
			if podMaxPids != c.expectedPodMaxPids {
				t.Fatalf("expected --pod-max-pids be equal to %s, got %s", c.expectedPodMaxPids, podMaxPids)
			}
			featureGates := c.cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig["--feature-gates"]
			hasSupportPodPidsLimitFeatureGate := strings.Contains(featureGates, "SupportPodPidsLimit=true")
			if hasSupportPodPidsLimitFeatureGate != c.expectedSupportPodPidsLimitFeatureGate {
				t.Fatalf("expected SupportPodPidsLimit=true presence in --feature gates to be %t, got %t", c.expectedSupportPodPidsLimitFeatureGate, hasSupportPodPidsLimitFeatureGate)
			}
		})
	}

}

func TestSupportPodPidsLimitFeatureGateInAgentPool(t *testing.T) {
	cases := []struct {
		name                                   string
		cs                                     *ContainerService
		isUpgrade                              bool
		expectedPodMaxPids                     string
		expectedSupportPodPidsLimitFeatureGate bool
	}{
		{
			name: "no --pod-max-pids defined",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{},
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined, no SupportPodPidsLimit defined",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--pod-max-pids": "100",
								},
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "100",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "no --pod-max-pids defined, SupportPodPidsLimit=false",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--feature-gates": "SupportPodPidsLimit=false",
								},
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "no --pod-max-pids defined, SupportPodPidsLimit=true",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--feature-gates": "SupportPodPidsLimit=true",
								},
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined, SupportPodPidsLimit=false",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--pod-max-pids":  "100",
									"--feature-gates": "SupportPodPidsLimit=false",
								},
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "100",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined, SupportPodPidsLimit=true",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--pod-max-pids":  "100",
									"--feature-gates": "SupportPodPidsLimit=true",
								},
							},
						},
					},
				},
			},
			expectedPodMaxPids:                     "100",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "no --pod-max-pids defined, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{},
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "no --pod-max-pids defined, SupportPodPidsLimit=false, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--feature-gates": "SupportPodPidsLimit=false",
								},
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "no --pod-max-pids defined, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--feature-gates": "SupportPodPidsLimit=true",
								},
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined as 100, no SupportPodPidsLimit defined, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--pod-max-pids": "100",
								},
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined as -100, no SupportPodPidsLimit defined, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--pod-max-pids": "-100",
								},
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined, SupportPodPidsLimit=false, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--pod-max-pids":  "100",
									"--feature-gates": "SupportPodPidsLimit=false",
								},
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: false,
		},
		{
			name: "--pod-max-pids defined as 100, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--pod-max-pids":  "100",
									"--feature-gates": "SupportPodPidsLimit=true",
								},
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "100",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined as 0, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--feature-gates": "SupportPodPidsLimit=true",
								},
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined as -1, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--pod-max-pids":  "-1",
									"--feature-gates": "SupportPodPidsLimit=true",
								},
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
		{
			name: "--pod-max-pids defined as 100a, SupportPodPidsLimit=true, upgrade scenario",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--pod-max-pids":  "100a",
									"--feature-gates": "SupportPodPidsLimit=true",
								},
							},
						},
					},
				},
			},
			isUpgrade:                              true,
			expectedPodMaxPids:                     "-1",
			expectedSupportPodPidsLimitFeatureGate: true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			c.cs.setKubeletConfig(c.isUpgrade)
			podMaxPids := c.cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig["--pod-max-pids"]
			if podMaxPids != c.expectedPodMaxPids {
				t.Fatalf("expected --pod-max-pids be equal to %s, got %s", c.expectedPodMaxPids, podMaxPids)
			}
			featureGates := c.cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig["--feature-gates"]
			hasSupportPodPidsLimitFeatureGate := strings.Contains(featureGates, "SupportPodPidsLimit=true")
			if hasSupportPodPidsLimitFeatureGate != c.expectedSupportPodPidsLimitFeatureGate {
				t.Fatalf("expected SupportPodPidsLimit=true presence in --feature gates to be %t, got %t", c.expectedSupportPodPidsLimitFeatureGate, hasSupportPodPidsLimitFeatureGate)
			}
		})
	}
}

func TestInputOverridesDuringUpgrade(t *testing.T) {
	defaultConfig := map[string]string{
		"--pod-infra-container-image": pauseImageReference,
	}

	defaultWindowsConfig := map[string]string{
		"--pod-infra-container-image": "kubletwin/pause",
	}

	inputConfig := map[string]string{
		"--pod-infra-container-image": "input-pause",
	}

	cases := []struct {
		name                     string
		cs                       *ContainerService
		isUpgrade                bool
		expectedKubeletConfig    map[string]string
		expectedWinKubeletConfig map[string]string
	}{
		{
			name: "use default value if no input",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							OSType: "Linux",
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{},
							},
						},
						{
							OSType: "Windows",
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{},
							},
						},
					},
				},
			},
			isUpgrade:                false,
			expectedKubeletConfig:    defaultConfig,
			expectedWinKubeletConfig: defaultWindowsConfig,
		},
		{
			name: "use input if not upgrade",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-infra-container-image": "input-pause",
							},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-infra-container-image": "input-pause",
							},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							OSType: "Linux",
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--pod-infra-container-image": "input-pause",
								},
							},
						},
						{
							OSType: "Windows",
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{},
							},
						},
					},
				},
			},
			isUpgrade:             false,
			expectedKubeletConfig: inputConfig,
			// --pod-infra-container-image is static for windows nodes
			expectedWinKubeletConfig: defaultWindowsConfig,
		},
		{
			name: "override input if upgrade",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-infra-container-image": "input-pause",
							},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{
								"--pod-infra-container-image": "input-pause",
							},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							OSType: "Linux",
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{
									"--pod-infra-container-image": "input-pause",
								},
							},
						},
						{
							OSType: "Windows",
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{},
							},
						},
					},
				},
			},
			isUpgrade:                true,
			expectedKubeletConfig:    defaultConfig,
			expectedWinKubeletConfig: defaultWindowsConfig,
		},
		{
			name: "override input if no input and upgrade",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.18.2",
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					MasterProfile: &MasterProfile{
						KubernetesConfig: &KubernetesConfig{
							KubeletConfig: map[string]string{},
						},
					},
					AgentPoolProfiles: []*AgentPoolProfile{
						{
							OSType: "Linux",
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{},
							},
						},
						{
							OSType: "Windows",
							KubernetesConfig: &KubernetesConfig{
								KubeletConfig: map[string]string{},
							},
						},
					},
				},
			},
			isUpgrade:                true,
			expectedKubeletConfig:    defaultConfig,
			expectedWinKubeletConfig: defaultWindowsConfig,
		},
	}

	assert := func(profile string, expected, actual map[string]string, t *testing.T) {
		for ek, ev := range expected {
			av, ok := actual[ek]
			if !ok {
				t.Fatalf("%s missing expected config %s", profile, ek)
			}
			if av != ev {
				t.Fatalf("%s expected config %s to equal %s, got %s", profile, ek, ev, av)
			}
		}
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			c.cs.setKubeletConfig(c.isUpgrade)
			assert("OrchestratorProfile", c.expectedKubeletConfig, c.cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig, t)
			assert("MasterProfile", c.expectedKubeletConfig, c.cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig, t)
			assert("Linux AgentPoolProfile", c.expectedKubeletConfig, c.cs.Properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig, t)
			assert("Windows AgentPoolProfile", c.expectedWinKubeletConfig, c.cs.Properties.AgentPoolProfiles[1].KubernetesConfig.KubeletConfig, t)
		})
	}
}

func TestReadOnlyPort(t *testing.T) {
	cases := []struct {
		name                 string
		cs                   *ContainerService
		expectedReadOnlyPort string
	}{
		{
			name: "default pre-1.16",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.15.0",
						KubernetesConfig:    &KubernetesConfig{},
					},
				},
			},
			expectedReadOnlyPort: "",
		},
		{
			name: "default 1.16",
			cs: &ContainerService{
				Properties: &Properties{
					OrchestratorProfile: &OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: "1.16.0",
						KubernetesConfig:    &KubernetesConfig{},
					},
				},
			},
			expectedReadOnlyPort: "0",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			c.cs.setKubeletConfig(false)
			readOnlyPort := c.cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig["--read-only-port"]
			if readOnlyPort != c.expectedReadOnlyPort {
				t.Fatalf("expected --read-only-port be equal to %s, got %s", c.expectedReadOnlyPort, readOnlyPort)
			}
		})
	}

}

func TestRemoveKubeletFlags(t *testing.T) {
	cases := []struct {
		name          string
		kubeletConfig map[string]string
		version       string
		expected      map[string]string
	}{
		{
			name: "v1.17.0",
			kubeletConfig: map[string]string{
				"--pod-max-pids":     "100",
				"--cadvisor-port":    "1234",
				"--allow-privileged": "true",
			},
			expected: map[string]string{
				"--pod-max-pids": "100",
			},
			version: "1.17.0",
		},
		{
			name: "v1.9.0",
			kubeletConfig: map[string]string{
				"--pod-max-pids":     "100",
				"--cadvisor-port":    "1234",
				"--allow-privileged": "true",
			},
			expected: map[string]string{
				"--cadvisor-port":    "1234",
				"--allow-privileged": "true",
			},
			version: "1.9.0",
		},
		{
			name: "v1.24.0",
			kubeletConfig: map[string]string{
				"--cni-conf-dir":    "/opt/cni/conf",
				"--cni-bin-dir":     "/opt/cni/bin",
				"--cni-cache-dir":   "/opt/cni/cache",
				"--docker-endpoint": "/docker.sock",
				"--experimental-dockershim-root-directory": "/some/dir",
				"--image-pull-progress-deadline":           "30m",
				"--network-plugin":                         "cni",
				"--network-plugin-mtu":                     "2000",
				"--pod-max-pids":                           "100",
				"--non-masquerade-cidr":                    "0.0.0.0/0",
			},
			expected: map[string]string{
				"--pod-max-pids": "100",
			},
			version: "1.24.0-alpha",
		},
		{
			name: "v1.27.0",
			kubeletConfig: map[string]string{
				"--pod-max-pids":             "100",
				"--master-service-namespace": "default",
				"--container-runtime":        "remote",
			},
			expected: map[string]string{
				"--pod-max-pids": "100",
			},
			version: "1.27.0",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			removeKubeletFlags(c.kubeletConfig, c.version)
			diff := cmp.Diff(c.kubeletConfig, c.expected)
			if diff != "" {
				t.Errorf("unexpected diff while expecting equal structs: %s", diff)
			}
		})
	}
}
