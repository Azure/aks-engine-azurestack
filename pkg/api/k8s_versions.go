// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"strconv"
	"strings"

	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
)

const (
	pauseImageReference                               string = "oss/kubernetes/pause:3.8"
	smbFlexVolumeImageReference                       string = "mcr.microsoft.com/k8s/flexvolume/smb-flexvolume:1.0.2"
	calicoTyphaImageReference                         string = "typha:v3.8.9"
	calicoCNIImageReference                           string = "cni:v3.8.9"
	calicoNodeImageReference                          string = "node:v3.8.9"
	calicoPod2DaemonImageReference                    string = "pod2daemon-flexvol:v3.8.0"
	calicoClusterProportionalAutoscalerImageReference string = "mcr.microsoft.com/oss/kubernetes/autoscaler/cluster-proportional-autoscaler:1.1.2-r2"
	ciliumAgentImageReference                         string = "docker.io/cilium/cilium:v1.4"
	ciliumCleanStateImageReference                    string = "docker.io/cilium/cilium-init:2018-10-16"
	ciliumOperatorImageReference                      string = "docker.io/cilium/operator:v1.4"
	ciliumEtcdOperatorImageReference                  string = "docker.io/cilium/cilium-etcd-operator:v2.0.5"
	antreaControllerImageReference                    string = "projects.registry.vmware.com/antrea/antrea-ubuntu:v1.3.0"
	antreaAgentImageReference                                = antreaControllerImageReference
	antreaOVSImageReference                                  = antreaControllerImageReference
	antreaInstallCNIImageReference                           = antreaControllerImageReference
	azureNPMContainerImageReference                   string = "mcr.microsoft.com/containernetworking/azure-npm:v1.4.59"
	aadPodIdentityNMIImageReference                   string = "mcr.microsoft.com/k8s/aad-pod-identity/nmi:1.6.1"
	aadPodIdentityMICImageReference                   string = "mcr.microsoft.com/k8s/aad-pod-identity/mic:1.6.1"
	azurePolicyImageReference                         string = "mcr.microsoft.com/azure-policy/policy-kubernetes-addon-prod:prod_20201023.1"
	gatekeeperImageReference                          string = "mcr.microsoft.com/oss/open-policy-agent/gatekeeper:v3.2.3"
	nodeProblemDetectorImageReference                 string = "registry.k8s.io/node-problem-detector/node-problem-detector:v0.8.4"
	csiAzureFileImageReference                        string = "oss/kubernetes-csi/azurefile-csi:v1.9.0"
	azureCloudControllerManagerImageReference         string = "oss/kubernetes/azure-cloud-controller-manager:v1.1.1"
	azureCloudNodeManagerImageReference               string = "oss/kubernetes/azure-cloud-node-manager:v1.1.1"
	dashboardImageReference                           string = "mcr.microsoft.com/oss/kubernetes/dashboard:v2.0.4" // deprecated
	dashboardMetricsScraperImageReference             string = "mcr.microsoft.com/oss/kubernetes/metrics-scraper:v1.0.4"
	kubeFlannelImageReference                         string = "quay.io/coreos/flannel:v0.8.0-amd64"
	flannelInstallCNIImageReference                   string = "quay.io/coreos/flannel:v0.10.0-amd64"
	KubeRBACProxyImageReference                       string = "gcr.io/kubebuilder/kube-rbac-proxy:v0.4.0"
	ScheduledMaintenanceManagerImageReference         string = "quay.io/awesomenix/drainsafe-manager:latest"
	nvidiaDevicePluginImageReference                  string = "oss/nvidia/k8s-device-plugin:1.0.0-beta6"
	virtualKubeletImageReference                      string = "virtual-kubelet:1.2.1.2" // Deprecated
	tillerImageReference                              string = "oss/kubernetes/tiller:v2.13.1"
	csiSecretsStoreProviderAzureImageReference        string = "oss/azure/secrets-store/provider-azure:0.0.12"
	csiSecretsStoreDriverImageReference               string = "oss/kubernetes-csi/secrets-store/driver:v0.0.19"
	clusterProportionalAutoscalerImageReference       string = "mcr.microsoft.com/oss/kubernetes/autoscaler/cluster-proportional-autoscaler:1.8.5"
	azureArcOnboardingImageReference                  string = "arck8sonboarding.azurecr.io/arck8sonboarding:v0.1.0"
	azureKMSProviderImageReference                    string = "k8s/kms/keyvault:v0.0.10"
)

var kubernetesImageBaseDefaultImages = map[string]map[string]string{
	common.KubernetesImageBaseTypeGCR: {
		common.DashboardAddonName:                   "kubernetes-dashboard-amd64:v1.10.1", // deprecated
		common.DashboardMetricsScraperContainerName: "",
		common.ExecHealthZComponentName:             "exechealthz-amd64:1.2",
		common.CoreDNSAddonName:                     "coredns:1.6.7",
		common.KubeDNSAddonName:                     "k8s-dns-kube-dns-amd64:1.15.4",
		common.DNSMasqComponentName:                 "k8s-dns-dnsmasq-nanny-amd64:1.15.4",
		common.DNSSidecarComponentName:              "k8s-dns-sidecar-amd64:1.14.10",
		common.ReschedulerAddonName:                 "rescheduler:v0.4.0", // Deprecated
		common.IPMASQAgentAddonName:                 "networking/ip-masq-agent:v2.8.0",
		common.KubeProxyAddonName:                   "kube-proxy",
		common.ControllerManagerComponentName:       "kube-controller-manager",
		common.APIServerComponentName:               "kube-apiserver",
		common.SchedulerComponentName:               "kube-scheduler",
		common.Hyperkube:                            "hyperkube-amd64",
	},
	common.KubernetesImageBaseTypeMCR: {
		common.DashboardAddonName:                   "oss/kubernetes/dashboard:v2.0.4", // deprecated
		common.DashboardMetricsScraperContainerName: "oss/kubernetes/metrics-scraper:v1.0.4",
		common.ExecHealthZComponentName:             "oss/kubernetes/exechealthz:1.2",
		common.CoreDNSAddonName:                     "oss/kubernetes/coredns:v1.9.4",
		common.KubeDNSAddonName:                     "oss/kubernetes/k8s-dns-kube-dns:1.15.4",
		common.DNSMasqComponentName:                 "oss/kubernetes/k8s-dns-dnsmasq-nanny:1.15.4",
		common.DNSSidecarComponentName:              "oss/kubernetes/k8s-dns-sidecar:1.14.10",
		common.ReschedulerAddonName:                 "oss/kubernetes/rescheduler:v0.4.0", // Deprecated
		common.IPMASQAgentAddonName:                 "oss/kubernetes/ip-masq-agent:v2.8.0",
		common.KubeProxyAddonName:                   "oss/kubernetes/kube-proxy",
		common.ControllerManagerComponentName:       "oss/kubernetes/kube-controller-manager",
		common.APIServerComponentName:               "oss/kubernetes/kube-apiserver",
		common.SchedulerComponentName:               "oss/kubernetes/kube-scheduler",
		common.Hyperkube:                            "oss/kubernetes/hyperkube",
	},
}

var csiSidecarComponentsOverrides = map[string]map[string]string{
	common.AzureFileCSIDriverAddonName: {
		common.CSIProvisionerContainerName: "oss/kubernetes-csi/csi-provisioner:v2.2.2",
		common.CSISnapshotterContainerName: "oss/kubernetes-csi/csi-snapshotter:v4.2.1",
	},
}

func getDefaultImage(image, kubernetesImageBaseType string) string {
	return kubernetesImageBaseDefaultImages[kubernetesImageBaseType][image]
}

// kubernetesImageBaseVersionedImages is a convenience map for "kubernetesImageBase" image version references that are distinct across versions of Kubernetes
// For example, cluster-autoscaler generally ships a per-Kubernetes-version build
// The map supports GCR or MCR image string flavors
var kubernetesImageBaseVersionedImages = map[string]map[string]map[string]string{
	common.KubernetesImageBaseTypeGCR: {
		"1.31": {
			common.CSIProvisionerContainerName:                "oss/kubernetes-csi/csi-provisioner:v5.2.0",
			common.CSIAttacherContainerName:                   "oss/kubernetes-csi/csi-attacher:v4.8.1",
			common.CSILivenessProbeContainerName:              "oss/kubernetes-csi/livenessprobe:v2.15.0",
			common.CSILivenessProbeWindowsContainerName:       "oss/kubernetes-csi/livenessprobe:v2.15.0",
			common.CSINodeDriverRegistrarContainerName:        "oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
			common.CSINodeDriverRegistrarWindowsContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
			common.CSISnapshotterContainerName:                "oss/kubernetes-csi/csi-snapshotter:v8.2.0",
			common.CSISnapshotControllerContainerName:         "oss/kubernetes-csi/snapshot-controller:v8.2.0",
			common.CSIResizerContainerName:                    "oss/kubernetes-csi/csi-resizer:v1.13.2",
			common.CSIAzureDiskContainerName:                  "oss/kubernetes-csi/azuredisk-csi:v1.31.10",
			common.AddonResizerComponentName:                  "addon-resizer:1.8.7",
			common.MetricsServerAddonName:                     "metrics-server/metrics-server:v0.5.2",
			common.AddonManagerComponentName:                  "kube-addon-manager-amd64:v9.1.6",
			common.ClusterAutoscalerAddonName:                 "cluster-autoscaler:v1.22.1",
		},
		"1.30": {
			common.CSIProvisionerContainerName:                "oss/kubernetes-csi/csi-provisioner:v5.2.0",
			common.CSIAttacherContainerName:                   "oss/kubernetes-csi/csi-attacher:v4.8.0",
			common.CSILivenessProbeContainerName:              "oss/kubernetes-csi/livenessprobe:v2.15.0",
			common.CSILivenessProbeWindowsContainerName:       "oss/kubernetes-csi/livenessprobe:v2.15.0",
			common.CSINodeDriverRegistrarContainerName:        "oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
			common.CSINodeDriverRegistrarWindowsContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
			common.CSISnapshotterContainerName:                "oss/kubernetes-csi/csi-snapshotter:v8.2.0",
			common.CSISnapshotControllerContainerName:         "oss/kubernetes-csi/snapshot-controller:v8.2.0",
			common.CSIResizerContainerName:                    "oss/kubernetes-csi/csi-resizer:v1.13.1",
			common.CSIAzureDiskContainerName:                  "oss/kubernetes-csi/azuredisk-csi:v1.31.5",
			common.AddonResizerComponentName:                  "addon-resizer:1.8.7",
			common.MetricsServerAddonName:                     "metrics-server/metrics-server:v0.5.2",
			common.AddonManagerComponentName:                  "kube-addon-manager-amd64:v9.1.6",
			common.ClusterAutoscalerAddonName:                 "cluster-autoscaler:v1.22.1",
		},
		"1.29": {
			common.CSIProvisionerContainerName:                "oss/kubernetes-csi/csi-provisioner:v3.5.0",
			common.CSIAttacherContainerName:                   "oss/kubernetes-csi/csi-attacher:v4.3.0",
			common.CSILivenessProbeContainerName:              "oss/kubernetes-csi/livenessprobe:v2.10.0",
			common.CSILivenessProbeWindowsContainerName:       "oss/kubernetes-csi/livenessprobe:v2.10.0",
			common.CSINodeDriverRegistrarContainerName:        "oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
			common.CSINodeDriverRegistrarWindowsContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
			common.CSISnapshotterContainerName:                "oss/kubernetes-csi/csi-snapshotter:v6.2.2",
			common.CSISnapshotControllerContainerName:         "oss/kubernetes-csi/snapshot-controller:v6.2.2",
			common.CSIResizerContainerName:                    "oss/kubernetes-csi/csi-resizer:v1.8.0",
			common.CSIAzureDiskContainerName:                  "oss/kubernetes-csi/azuredisk-csi:v1.29.1",
			common.AddonResizerComponentName:                  "addon-resizer:1.8.7",
			common.MetricsServerAddonName:                     "metrics-server/metrics-server:v0.5.0",
			common.AddonManagerComponentName:                  "kube-addon-manager-amd64:v9.1.6",
			common.ClusterAutoscalerAddonName:                 "cluster-autoscaler:v1.18.0",
		},
		"1.28": {
			common.CSIProvisionerContainerName:                "oss/kubernetes-csi/csi-provisioner:v3.5.0",
			common.CSIAttacherContainerName:                   "oss/kubernetes-csi/csi-attacher:v4.3.0",
			common.CSILivenessProbeContainerName:              "oss/kubernetes-csi/livenessprobe:v2.10.0",
			common.CSILivenessProbeWindowsContainerName:       "oss/kubernetes-csi/livenessprobe:v2.10.0",
			common.CSINodeDriverRegistrarContainerName:        "oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
			common.CSINodeDriverRegistrarWindowsContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
			common.CSISnapshotterContainerName:                "oss/kubernetes-csi/csi-snapshotter:v6.2.2",
			common.CSISnapshotControllerContainerName:         "oss/kubernetes-csi/snapshot-controller:v6.2.2",
			common.CSIResizerContainerName:                    "oss/kubernetes-csi/csi-resizer:v1.8.0",
			common.CSIAzureDiskContainerName:                  "oss/kubernetes-csi/azuredisk-csi:v1.29.1",
			common.AddonResizerComponentName:                  "addon-resizer:1.8.7",
			common.MetricsServerAddonName:                     "metrics-server/metrics-server:v0.5.0",
			common.AddonManagerComponentName:                  "kube-addon-manager-amd64:v9.1.6",
			common.ClusterAutoscalerAddonName:                 "cluster-autoscaler:v1.18.0",
		},
		"1.27": {
			common.CSIProvisionerContainerName:                "oss/kubernetes-csi/csi-provisioner:v3.5.0",
			common.CSIAttacherContainerName:                   "oss/kubernetes-csi/csi-attacher:v4.3.0",
			common.CSILivenessProbeContainerName:              "oss/kubernetes-csi/livenessprobe:v2.10.0",
			common.CSILivenessProbeWindowsContainerName:       "oss/kubernetes-csi/livenessprobe:v2.10.0",
			common.CSINodeDriverRegistrarContainerName:        "oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
			common.CSINodeDriverRegistrarWindowsContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
			common.CSISnapshotterContainerName:                "oss/kubernetes-csi/csi-snapshotter:v6.2.2",
			common.CSISnapshotControllerContainerName:         "oss/kubernetes-csi/snapshot-controller:v6.2.2",
			common.CSIResizerContainerName:                    "oss/kubernetes-csi/csi-resizer:v1.8.0",
			common.CSIAzureDiskContainerName:                  "oss/kubernetes-csi/azuredisk-csi:v1.28.3",
			common.AddonResizerComponentName:                  "addon-resizer:1.8.7",
			common.MetricsServerAddonName:                     "metrics-server/metrics-server:v0.5.0",
			common.AddonManagerComponentName:                  "kube-addon-manager-amd64:v9.1.6",
			common.ClusterAutoscalerAddonName:                 "cluster-autoscaler:v1.18.0",
		},
	},
	common.KubernetesImageBaseTypeMCR: {
		"1.31": {
			common.CSIProvisionerContainerName:                "oss/kubernetes-csi/csi-provisioner:v5.2.0",
			common.CSIAttacherContainerName:                   "oss/kubernetes-csi/csi-attacher:v4.8.1",
			common.CSILivenessProbeContainerName:              "oss/kubernetes-csi/livenessprobe:v2.15.0",
			common.CSILivenessProbeWindowsContainerName:       "oss/kubernetes-csi/livenessprobe:v2.15.0",
			common.CSINodeDriverRegistrarContainerName:        "oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
			common.CSINodeDriverRegistrarWindowsContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
			common.CSISnapshotterContainerName:                "oss/kubernetes-csi/csi-snapshotter:v8.2.0",
			common.CSISnapshotControllerContainerName:         "oss/kubernetes-csi/snapshot-controller:v8.2.0",
			common.CSIResizerContainerName:                    "oss/kubernetes-csi/csi-resizer:v1.13.2",
			common.CSIAzureDiskContainerName:                  "oss/kubernetes-csi/azuredisk-csi:v1.31.10",
			common.AddonResizerComponentName:                  "oss/kubernetes/autoscaler/addon-resizer:1.8.7",
			common.MetricsServerAddonName:                     "oss/kubernetes/metrics-server:v0.5.2",
			common.AddonManagerComponentName:                  "oss/kubernetes/kube-addon-manager:v9.1.6",
			common.ClusterAutoscalerAddonName:                 "oss/kubernetes/autoscaler/cluster-autoscaler:v1.22.1",
		},
		"1.30": {
			common.CSIProvisionerContainerName:                "oss/kubernetes-csi/csi-provisioner:v5.2.0",
			common.CSIAttacherContainerName:                   "oss/kubernetes-csi/csi-attacher:v4.8.0",
			common.CSILivenessProbeContainerName:              "oss/kubernetes-csi/livenessprobe:v2.15.0",
			common.CSILivenessProbeWindowsContainerName:       "oss/kubernetes-csi/livenessprobe:v2.15.0",
			common.CSINodeDriverRegistrarContainerName:        "oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
			common.CSINodeDriverRegistrarWindowsContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.13.0",
			common.CSISnapshotterContainerName:                "oss/kubernetes-csi/csi-snapshotter:v8.2.0",
			common.CSISnapshotControllerContainerName:         "oss/kubernetes-csi/snapshot-controller:v8.2.0",
			common.CSIResizerContainerName:                    "oss/kubernetes-csi/csi-resizer:v1.13.1",
			common.CSIAzureDiskContainerName:                  "oss/kubernetes-csi/azuredisk-csi:v1.31.5",
			common.AddonResizerComponentName:                  "oss/kubernetes/autoscaler/addon-resizer:1.8.7",
			common.MetricsServerAddonName:                     "oss/kubernetes/metrics-server:v0.5.2",
			common.AddonManagerComponentName:                  "oss/kubernetes/kube-addon-manager:v9.1.6",
			common.ClusterAutoscalerAddonName:                 "oss/kubernetes/autoscaler/cluster-autoscaler:v1.22.1",
		},
		"1.29": {
			common.CSIProvisionerContainerName:                "oss/kubernetes-csi/csi-provisioner:v3.5.0",
			common.CSIAttacherContainerName:                   "oss/kubernetes-csi/csi-attacher:v4.3.0",
			common.CSILivenessProbeContainerName:              "oss/kubernetes-csi/livenessprobe:v2.10.0",
			common.CSILivenessProbeWindowsContainerName:       "oss/kubernetes-csi/livenessprobe:v2.10.0",
			common.CSINodeDriverRegistrarContainerName:        "oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
			common.CSINodeDriverRegistrarWindowsContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
			common.CSISnapshotterContainerName:                "oss/kubernetes-csi/csi-snapshotter:v6.2.2",
			common.CSISnapshotControllerContainerName:         "oss/kubernetes-csi/snapshot-controller:v6.2.2",
			common.CSIResizerContainerName:                    "oss/kubernetes-csi/csi-resizer:v1.8.0",
			common.CSIAzureDiskContainerName:                  "oss/kubernetes-csi/azuredisk-csi:v1.29.1",
			common.AddonResizerComponentName:                  "oss/kubernetes/autoscaler/addon-resizer:1.8.7",
			common.MetricsServerAddonName:                     "oss/kubernetes/metrics-server:v0.5.2",
			common.AddonManagerComponentName:                  "oss/kubernetes/kube-addon-manager:v9.1.6",
			common.ClusterAutoscalerAddonName:                 "oss/kubernetes/autoscaler/cluster-autoscaler:v1.22.1",
		},
		"1.28": {
			common.CSIProvisionerContainerName:                "oss/kubernetes-csi/csi-provisioner:v3.5.0",
			common.CSIAttacherContainerName:                   "oss/kubernetes-csi/csi-attacher:v4.3.0",
			common.CSILivenessProbeContainerName:              "oss/kubernetes-csi/livenessprobe:v2.10.0",
			common.CSILivenessProbeWindowsContainerName:       "oss/kubernetes-csi/livenessprobe:v2.10.0",
			common.CSINodeDriverRegistrarContainerName:        "oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
			common.CSINodeDriverRegistrarWindowsContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
			common.CSISnapshotterContainerName:                "oss/kubernetes-csi/csi-snapshotter:v6.2.2",
			common.CSISnapshotControllerContainerName:         "oss/kubernetes-csi/snapshot-controller:v6.2.2",
			common.CSIResizerContainerName:                    "oss/kubernetes-csi/csi-resizer:v1.8.0",
			common.CSIAzureDiskContainerName:                  "oss/kubernetes-csi/azuredisk-csi:v1.29.1",
			common.AddonResizerComponentName:                  "oss/kubernetes/autoscaler/addon-resizer:1.8.7",
			common.MetricsServerAddonName:                     "oss/kubernetes/metrics-server:v0.5.2",
			common.AddonManagerComponentName:                  "oss/kubernetes/kube-addon-manager:v9.1.6",
			common.ClusterAutoscalerAddonName:                 "oss/kubernetes/autoscaler/cluster-autoscaler:v1.22.1",
		},
		"1.27": {
			common.CSIProvisionerContainerName:                "oss/kubernetes-csi/csi-provisioner:v3.5.0",
			common.CSIAttacherContainerName:                   "oss/kubernetes-csi/csi-attacher:v4.3.0",
			common.CSILivenessProbeContainerName:              "oss/kubernetes-csi/livenessprobe:v2.10.0",
			common.CSILivenessProbeWindowsContainerName:       "oss/kubernetes-csi/livenessprobe:v2.10.0",
			common.CSINodeDriverRegistrarContainerName:        "oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
			common.CSINodeDriverRegistrarWindowsContainerName: "oss/kubernetes-csi/csi-node-driver-registrar:v2.8.0",
			common.CSISnapshotterContainerName:                "oss/kubernetes-csi/csi-snapshotter:v6.2.2",
			common.CSISnapshotControllerContainerName:         "oss/kubernetes-csi/snapshot-controller:v6.2.2",
			common.CSIResizerContainerName:                    "oss/kubernetes-csi/csi-resizer:v1.8.0",
			common.CSIAzureDiskContainerName:                  "oss/kubernetes-csi/azuredisk-csi:v1.28.3",
			common.AddonResizerComponentName:                  "oss/kubernetes/autoscaler/addon-resizer:1.8.7",
			common.MetricsServerAddonName:                     "oss/kubernetes/metrics-server:v0.5.2",
			common.AddonManagerComponentName:                  "oss/kubernetes/kube-addon-manager:v9.1.6",
			common.ClusterAutoscalerAddonName:                 "oss/kubernetes/autoscaler/cluster-autoscaler:v1.22.1",
		},
	},
}

type getK8sVersionComponentsOverrides func(string) map[string]string

func GetK8sComponentsByVersionMap(k *KubernetesConfig) map[string]map[string]string {
	var overrides getK8sVersionComponentsOverrides
	switch k.KubernetesImageBaseType {
	case common.KubernetesImageBaseTypeGCR:
		overrides = getVersionOverridesGCR
	case common.KubernetesImageBaseTypeMCR:
		overrides = getVersionOverridesMCR
	default:
		overrides = getVersionOverridesGCR
	}
	ret := make(map[string]map[string]string)
	for _, version := range common.GetAllSupportedKubernetesVersions(true, false, false) {
		ret[version] = getK8sVersionComponents(version, k.KubernetesImageBaseType, overrides(version))
	}
	return ret
}

func getVersionOverridesMCR(v string) map[string]string {
	switch v {
	case "1.18.6":
		return map[string]string{common.WindowsArtifactComponentName: "v1.18.6-hotfix.20200723/windowszip/v1.18.6-hotfix.20200723-1int.zip"}
	case "1.18.4":
		return map[string]string{common.WindowsArtifactComponentName: "v1.18.4-hotfix.20200626/windowszip/v1.18.4-hotfix.20200626-1int.zip"}
	case "1.18.2":
		return map[string]string{common.WindowsArtifactComponentName: "v1.18.2-hotfix.20200624/windowszip/v1.18.2-hotfix.20200624-1int.zip"}
	case "1.17.9":
		return map[string]string{common.WindowsArtifactComponentName: "v1.17.9-hotfix.20200817/windowszip/v1.17.9-hotfix.20200817-1int.zip"}
	case "1.17.7":
		return map[string]string{common.WindowsArtifactComponentName: "v1.17.7-hotfix.20200817/windowszip/v1.17.7-hotfix.20200817-1int.zip"}
	case "1.16.13":
		return map[string]string{common.WindowsArtifactComponentName: "v1.16.13-hotfix.20200817/windowszip/v1.16.13-hotfix.20200817-1int.zip"}
	case "1.16.11":
		return map[string]string{common.WindowsArtifactComponentName: "v1.16.11-hotfix.20200617/windowszip/v1.16.11-hotfix.20200617-1int.zip"}
	case "1.16.10":
		return map[string]string{common.WindowsArtifactComponentName: "v1.16.10-hotfix.20200817/windowszip/v1.16.10-hotfix.20200817-1int.zip"}
	case "1.15.12":
		return map[string]string{common.WindowsArtifactComponentName: "v1.15.12-hotfix.20200817/windowszip/v1.15.12-hotfix.20200817-1int.zip"}
	case "1.15.11":
		return map[string]string{common.WindowsArtifactComponentName: "v1.15.11-hotfix.20200817/windowszip/v1.15.11-hotfix.20200817-1int.zip"}
	default:
		return nil
	}
}

func getVersionOverridesGCR(v string) map[string]string {
	switch v {
	case "1.18.6":
		return map[string]string{common.WindowsArtifactComponentName: "v1.18.6-hotfix.20200723/windowszip/v1.18.6-hotfix.20200723-1int.zip"}
	case "1.18.4":
		return map[string]string{common.WindowsArtifactComponentName: "v1.18.4-hotfix.20200626/windowszip/v1.18.4-hotfix.20200626-1int.zip"}
	case "1.18.2":
		return map[string]string{common.WindowsArtifactComponentName: "v1.18.2-hotfix.20200624/windowszip/v1.18.2-hotfix.20200624-1int.zip"}
	case "1.17.9":
		return map[string]string{common.WindowsArtifactComponentName: "v1.17.9-hotfix.20200817/windowszip/v1.17.9-hotfix.20200817-1int.zip"}
	case "1.17.7":
		return map[string]string{common.WindowsArtifactComponentName: "v1.17.7-hotfix.20200817/windowszip/v1.17.7-hotfix.20200817-1int.zip"}
	case "1.16.13":
		return map[string]string{common.WindowsArtifactComponentName: "v1.16.13-hotfix.20200817/windowszip/v1.16.13-hotfix.20200817-1int.zip"}
	case "1.16.11":
		return map[string]string{common.WindowsArtifactComponentName: "v1.16.11-hotfix.20200617/windowszip/v1.16.11-hotfix.20200617-1int.zip"}
	case "1.16.10":
		return map[string]string{common.WindowsArtifactComponentName: "v1.16.10-hotfix.20200817/windowszip/v1.16.10-hotfix.20200817-1int.zip"}
	case "1.15.12":
		return map[string]string{common.WindowsArtifactComponentName: "v1.15.12-hotfix.20200817/windowszip/v1.15.12-hotfix.20200817-1int.zip"}
	case "1.15.11":
		return map[string]string{common.WindowsArtifactComponentName: "v1.15.11-hotfix.20200817/windowszip/v1.15.11-hotfix.20200817-1int.zip"}
	case "1.8.11":
		return map[string]string{common.KubeDNSAddonName: "k8s-dns-kube-dns-amd64:1.14.9"}
	case "1.8.9":
		return map[string]string{common.WindowsArtifactComponentName: "v1.8.9-2int.zip"}
	case "1.8.6":
		return map[string]string{common.WindowsArtifactComponentName: "v1.8.6-2int.zip"}
	case "1.8.2":
		return map[string]string{common.WindowsArtifactComponentName: "v1.8.2-2int.zip"}
	case "1.8.1":
		return map[string]string{common.WindowsArtifactComponentName: "v1.8.1-2int.zip"}
	case "1.8.0":
		return map[string]string{common.WindowsArtifactComponentName: "v1.8.0-2int.zip"}
	case "1.7.16":
		return map[string]string{common.WindowsArtifactComponentName: "v1.7.16-1int.zip"}
	case "1.7.15":
		return map[string]string{common.WindowsArtifactComponentName: "v1.7.15-1int.zip"}
	case "1.7.14":
		return map[string]string{common.WindowsArtifactComponentName: "v1.7.14-1int.zip"}
	case "1.7.13":
		return map[string]string{common.WindowsArtifactComponentName: "v1.7.13-1int.zip"}
	case "1.7.12":
		return map[string]string{common.WindowsArtifactComponentName: "v1.7.12-2int.zip"}
	case "1.7.10":
		return map[string]string{common.WindowsArtifactComponentName: "v1.7.10-1int.zip"}
	case "1.7.9":
		return map[string]string{common.WindowsArtifactComponentName: "v1.7.9-2int.zip"}
	case "1.7.7":
		return map[string]string{common.WindowsArtifactComponentName: "v1.7.7-2int.zip"}
	case "1.7.5":
		return map[string]string{common.WindowsArtifactComponentName: "v1.7.5-4int.zip"}
	case "1.7.4":
		return map[string]string{common.WindowsArtifactComponentName: "v1.7.4-2int.zip"}
	case "1.7.2":
		return map[string]string{common.WindowsArtifactComponentName: "v1.7.2-1int.zip"}
	default:
		return nil
	}
}

func getK8sVersionComponents(version, kubernetesImageBaseType string, overrides map[string]string) map[string]string {
	s := strings.Split(version, ".")
	majorMinor := strings.Join(s[:2], ".")
	var ret map[string]string
	k8sComponent := kubernetesImageBaseVersionedImages[kubernetesImageBaseType][majorMinor]
	switch majorMinor {
	case "1.31":
		ret = map[string]string{
			common.APIServerComponentName:                 getDefaultImage(common.APIServerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.ControllerManagerComponentName:         getDefaultImage(common.ControllerManagerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.KubeProxyAddonName:                     getDefaultImage(common.KubeProxyAddonName, kubernetesImageBaseType) + ":v" + version,
			common.SchedulerComponentName:                 getDefaultImage(common.SchedulerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.CloudControllerManagerComponentName:    "oss/kubernetes/azure-cloud-controller-manager:v1.31.7",
			common.CloudNodeManagerAddonName:              "oss/kubernetes/azure-cloud-node-manager:v1.31.7",
			common.WindowsArtifactComponentName:           "v" + version + "/windowszip/v" + version + "-1int.zip",
			common.WindowsArtifactAzureStackComponentName: "v" + version + "/windowszip/v" + version + "-1int.zip",
			common.DashboardAddonName:                     dashboardImageReference,
			common.DashboardMetricsScraperContainerName:   dashboardMetricsScraperImageReference,
			common.ExecHealthZComponentName:               getDefaultImage(common.ExecHealthZComponentName, kubernetesImageBaseType),
			common.AddonResizerComponentName:              k8sComponent[common.AddonResizerComponentName],
			common.MetricsServerAddonName:                 k8sComponent[common.MetricsServerAddonName],
			common.CoreDNSAddonName:                       getDefaultImage(common.CoreDNSAddonName, kubernetesImageBaseType),
			common.CoreDNSAutoscalerName:                  clusterProportionalAutoscalerImageReference,
			common.KubeDNSAddonName:                       getDefaultImage(common.KubeDNSAddonName, kubernetesImageBaseType),
			common.AddonManagerComponentName:              k8sComponent[common.AddonManagerComponentName],
			common.DNSMasqComponentName:                   getDefaultImage(common.DNSMasqComponentName, kubernetesImageBaseType),
			common.PauseComponentName:                     pauseImageReference,
			common.TillerAddonName:                        tillerImageReference,
			common.ReschedulerAddonName:                   getDefaultImage(common.ReschedulerAddonName, kubernetesImageBaseType),
			common.ACIConnectorAddonName:                  virtualKubeletImageReference,
			common.ClusterAutoscalerAddonName:             k8sComponent[common.ClusterAutoscalerAddonName],
			common.DNSSidecarComponentName:                getDefaultImage(common.DNSSidecarComponentName, kubernetesImageBaseType),

			common.SMBFlexVolumeAddonName:                     smbFlexVolumeImageReference,
			common.IPMASQAgentAddonName:                       getDefaultImage(common.IPMASQAgentAddonName, kubernetesImageBaseType),
			common.AzureNetworkPolicyAddonName:                azureNPMContainerImageReference,
			common.CalicoTyphaComponentName:                   calicoTyphaImageReference,
			common.CalicoCNIComponentName:                     calicoCNIImageReference,
			common.CalicoNodeComponentName:                    calicoNodeImageReference,
			common.CalicoPod2DaemonComponentName:              calicoPod2DaemonImageReference,
			common.CalicoClusterAutoscalerComponentName:       calicoClusterProportionalAutoscalerImageReference,
			common.CiliumAgentContainerName:                   ciliumAgentImageReference,
			common.CiliumCleanStateContainerName:              ciliumCleanStateImageReference,
			common.CiliumOperatorContainerName:                ciliumOperatorImageReference,
			common.CiliumEtcdOperatorContainerName:            ciliumEtcdOperatorImageReference,
			common.AntreaControllerContainerName:              antreaControllerImageReference,
			common.AntreaAgentContainerName:                   antreaAgentImageReference,
			common.AntreaOVSContainerName:                     antreaOVSImageReference,
			"antrea" + common.AntreaInstallCNIContainerName:   antreaInstallCNIImageReference,
			common.NMIContainerName:                           aadPodIdentityNMIImageReference,
			common.MICContainerName:                           aadPodIdentityMICImageReference,
			common.AzurePolicyAddonName:                       azurePolicyImageReference,
			common.GatekeeperContainerName:                    gatekeeperImageReference,
			common.NodeProblemDetectorAddonName:               nodeProblemDetectorImageReference,
			common.CSIProvisionerContainerName:                k8sComponent[common.CSIProvisionerContainerName],
			common.CSIAttacherContainerName:                   k8sComponent[common.CSIAttacherContainerName],
			common.CSILivenessProbeContainerName:              k8sComponent[common.CSILivenessProbeContainerName],
			common.CSILivenessProbeWindowsContainerName:       k8sComponent[common.CSILivenessProbeWindowsContainerName],
			common.CSINodeDriverRegistrarContainerName:        k8sComponent[common.CSINodeDriverRegistrarContainerName],
			common.CSINodeDriverRegistrarWindowsContainerName: k8sComponent[common.CSINodeDriverRegistrarWindowsContainerName],
			common.CSISnapshotterContainerName:                k8sComponent[common.CSISnapshotterContainerName],
			common.CSISnapshotControllerContainerName:         k8sComponent[common.CSISnapshotControllerContainerName],
			common.CSIResizerContainerName:                    k8sComponent[common.CSIResizerContainerName],
			common.CSIAzureDiskContainerName:                  k8sComponent[common.CSIAzureDiskContainerName],
			common.CSIAzureFileContainerName:                  csiAzureFileImageReference,
			common.KubeFlannelContainerName:                   kubeFlannelImageReference,
			"flannel" + common.FlannelInstallCNIContainerName: flannelInstallCNIImageReference,
			common.KubeRBACProxyContainerName:                 KubeRBACProxyImageReference,
			common.ScheduledMaintenanceManagerContainerName:   ScheduledMaintenanceManagerImageReference,
			"nodestatusfreq":                                  DefaultKubernetesNodeStatusUpdateFrequency,
			"nodegraceperiod":                                 DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
			"podeviction":                                     DefaultKubernetesCtrlMgrPodEvictionTimeout,
			"routeperiod":                                     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
			"backoffretries":                                  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
			"backoffjitter":                                   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
			"backoffduration":                                 strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
			"backoffexponent":                                 strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
			"ratelimitqps":                                    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
			"ratelimitqpswrite":                               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPSWrite, 'f', -1, 64),
			"ratelimitbucket":                                 strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
			"ratelimitbucketwrite":                            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucketWrite),
			"gchighthreshold":                                 strconv.Itoa(DefaultKubernetesGCHighThreshold),
			"gclowthreshold":                                  strconv.Itoa(DefaultKubernetesGCLowThreshold),
			common.NVIDIADevicePluginAddonName:                nvidiaDevicePluginImageReference,
			common.CSISecretsStoreProviderAzureContainerName:  csiSecretsStoreProviderAzureImageReference,
			common.CSISecretsStoreDriverContainerName:         csiSecretsStoreDriverImageReference,
			common.AzureArcOnboardingAddonName:                azureArcOnboardingImageReference,
			common.AzureKMSProviderComponentName:              azureKMSProviderImageReference,
		}
	case "1.30":
		ret = map[string]string{
			common.APIServerComponentName:                 getDefaultImage(common.APIServerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.ControllerManagerComponentName:         getDefaultImage(common.ControllerManagerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.KubeProxyAddonName:                     getDefaultImage(common.KubeProxyAddonName, kubernetesImageBaseType) + ":v" + version,
			common.SchedulerComponentName:                 getDefaultImage(common.SchedulerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.CloudControllerManagerComponentName:    "oss/kubernetes/azure-cloud-controller-manager:v1.30.7",
			common.CloudNodeManagerAddonName:              "oss/kubernetes/azure-cloud-node-manager:v1.30.8",
			common.WindowsArtifactComponentName:           "v" + version + "/windowszip/v" + version + "-1int.zip",
			common.WindowsArtifactAzureStackComponentName: "v" + version + "/windowszip/v" + version + "-1int.zip",
			common.DashboardAddonName:                     dashboardImageReference,
			common.DashboardMetricsScraperContainerName:   dashboardMetricsScraperImageReference,
			common.ExecHealthZComponentName:               getDefaultImage(common.ExecHealthZComponentName, kubernetesImageBaseType),
			common.AddonResizerComponentName:              k8sComponent[common.AddonResizerComponentName],
			common.MetricsServerAddonName:                 k8sComponent[common.MetricsServerAddonName],
			common.CoreDNSAddonName:                       getDefaultImage(common.CoreDNSAddonName, kubernetesImageBaseType),
			common.CoreDNSAutoscalerName:                  clusterProportionalAutoscalerImageReference,
			common.KubeDNSAddonName:                       getDefaultImage(common.KubeDNSAddonName, kubernetesImageBaseType),
			common.AddonManagerComponentName:              k8sComponent[common.AddonManagerComponentName],
			common.DNSMasqComponentName:                   getDefaultImage(common.DNSMasqComponentName, kubernetesImageBaseType),
			common.PauseComponentName:                     pauseImageReference,
			common.TillerAddonName:                        tillerImageReference,
			common.ReschedulerAddonName:                   getDefaultImage(common.ReschedulerAddonName, kubernetesImageBaseType),
			common.ACIConnectorAddonName:                  virtualKubeletImageReference,
			common.ClusterAutoscalerAddonName:             k8sComponent[common.ClusterAutoscalerAddonName],
			common.DNSSidecarComponentName:                getDefaultImage(common.DNSSidecarComponentName, kubernetesImageBaseType),

			common.SMBFlexVolumeAddonName:                     smbFlexVolumeImageReference,
			common.IPMASQAgentAddonName:                       getDefaultImage(common.IPMASQAgentAddonName, kubernetesImageBaseType),
			common.AzureNetworkPolicyAddonName:                azureNPMContainerImageReference,
			common.CalicoTyphaComponentName:                   calicoTyphaImageReference,
			common.CalicoCNIComponentName:                     calicoCNIImageReference,
			common.CalicoNodeComponentName:                    calicoNodeImageReference,
			common.CalicoPod2DaemonComponentName:              calicoPod2DaemonImageReference,
			common.CalicoClusterAutoscalerComponentName:       calicoClusterProportionalAutoscalerImageReference,
			common.CiliumAgentContainerName:                   ciliumAgentImageReference,
			common.CiliumCleanStateContainerName:              ciliumCleanStateImageReference,
			common.CiliumOperatorContainerName:                ciliumOperatorImageReference,
			common.CiliumEtcdOperatorContainerName:            ciliumEtcdOperatorImageReference,
			common.AntreaControllerContainerName:              antreaControllerImageReference,
			common.AntreaAgentContainerName:                   antreaAgentImageReference,
			common.AntreaOVSContainerName:                     antreaOVSImageReference,
			"antrea" + common.AntreaInstallCNIContainerName:   antreaInstallCNIImageReference,
			common.NMIContainerName:                           aadPodIdentityNMIImageReference,
			common.MICContainerName:                           aadPodIdentityMICImageReference,
			common.AzurePolicyAddonName:                       azurePolicyImageReference,
			common.GatekeeperContainerName:                    gatekeeperImageReference,
			common.NodeProblemDetectorAddonName:               nodeProblemDetectorImageReference,
			common.CSIProvisionerContainerName:                k8sComponent[common.CSIProvisionerContainerName],
			common.CSIAttacherContainerName:                   k8sComponent[common.CSIAttacherContainerName],
			common.CSILivenessProbeContainerName:              k8sComponent[common.CSILivenessProbeContainerName],
			common.CSILivenessProbeWindowsContainerName:       k8sComponent[common.CSILivenessProbeWindowsContainerName],
			common.CSINodeDriverRegistrarContainerName:        k8sComponent[common.CSINodeDriverRegistrarContainerName],
			common.CSINodeDriverRegistrarWindowsContainerName: k8sComponent[common.CSINodeDriverRegistrarWindowsContainerName],
			common.CSISnapshotterContainerName:                k8sComponent[common.CSISnapshotterContainerName],
			common.CSISnapshotControllerContainerName:         k8sComponent[common.CSISnapshotControllerContainerName],
			common.CSIResizerContainerName:                    k8sComponent[common.CSIResizerContainerName],
			common.CSIAzureDiskContainerName:                  k8sComponent[common.CSIAzureDiskContainerName],
			common.CSIAzureFileContainerName:                  csiAzureFileImageReference,
			common.KubeFlannelContainerName:                   kubeFlannelImageReference,
			"flannel" + common.FlannelInstallCNIContainerName: flannelInstallCNIImageReference,
			common.KubeRBACProxyContainerName:                 KubeRBACProxyImageReference,
			common.ScheduledMaintenanceManagerContainerName:   ScheduledMaintenanceManagerImageReference,
			"nodestatusfreq":                                  DefaultKubernetesNodeStatusUpdateFrequency,
			"nodegraceperiod":                                 DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
			"podeviction":                                     DefaultKubernetesCtrlMgrPodEvictionTimeout,
			"routeperiod":                                     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
			"backoffretries":                                  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
			"backoffjitter":                                   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
			"backoffduration":                                 strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
			"backoffexponent":                                 strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
			"ratelimitqps":                                    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
			"ratelimitqpswrite":                               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPSWrite, 'f', -1, 64),
			"ratelimitbucket":                                 strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
			"ratelimitbucketwrite":                            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucketWrite),
			"gchighthreshold":                                 strconv.Itoa(DefaultKubernetesGCHighThreshold),
			"gclowthreshold":                                  strconv.Itoa(DefaultKubernetesGCLowThreshold),
			common.NVIDIADevicePluginAddonName:                nvidiaDevicePluginImageReference,
			common.CSISecretsStoreProviderAzureContainerName:  csiSecretsStoreProviderAzureImageReference,
			common.CSISecretsStoreDriverContainerName:         csiSecretsStoreDriverImageReference,
			common.AzureArcOnboardingAddonName:                azureArcOnboardingImageReference,
			common.AzureKMSProviderComponentName:              azureKMSProviderImageReference,
		}
	case "1.29":
		ret = map[string]string{
			common.APIServerComponentName:                 getDefaultImage(common.APIServerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.ControllerManagerComponentName:         getDefaultImage(common.ControllerManagerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.KubeProxyAddonName:                     getDefaultImage(common.KubeProxyAddonName, kubernetesImageBaseType) + ":v" + version,
			common.SchedulerComponentName:                 getDefaultImage(common.SchedulerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.CloudControllerManagerComponentName:    "oss/kubernetes/azure-cloud-controller-manager:v1.29.8",
			common.CloudNodeManagerAddonName:              "oss/kubernetes/azure-cloud-node-manager:v1.29.9",
			common.WindowsArtifactComponentName:           "v" + version + "/windowszip/v" + version + "-1int.zip",
			common.WindowsArtifactAzureStackComponentName: "v" + version + "/windowszip/v" + version + "-1int.zip",
			common.DashboardAddonName:                     dashboardImageReference,
			common.DashboardMetricsScraperContainerName:   dashboardMetricsScraperImageReference,
			common.ExecHealthZComponentName:               getDefaultImage(common.ExecHealthZComponentName, kubernetesImageBaseType),
			common.AddonResizerComponentName:              k8sComponent[common.AddonResizerComponentName],
			common.MetricsServerAddonName:                 k8sComponent[common.MetricsServerAddonName],
			common.CoreDNSAddonName:                       getDefaultImage(common.CoreDNSAddonName, kubernetesImageBaseType),
			common.CoreDNSAutoscalerName:                  clusterProportionalAutoscalerImageReference,
			common.KubeDNSAddonName:                       getDefaultImage(common.KubeDNSAddonName, kubernetesImageBaseType),
			common.AddonManagerComponentName:              k8sComponent[common.AddonManagerComponentName],
			common.DNSMasqComponentName:                   getDefaultImage(common.DNSMasqComponentName, kubernetesImageBaseType),
			common.PauseComponentName:                     pauseImageReference,
			common.TillerAddonName:                        tillerImageReference,
			common.ReschedulerAddonName:                   getDefaultImage(common.ReschedulerAddonName, kubernetesImageBaseType),
			common.ACIConnectorAddonName:                  virtualKubeletImageReference,
			common.ClusterAutoscalerAddonName:             k8sComponent[common.ClusterAutoscalerAddonName],
			common.DNSSidecarComponentName:                getDefaultImage(common.DNSSidecarComponentName, kubernetesImageBaseType),

			common.SMBFlexVolumeAddonName:                     smbFlexVolumeImageReference,
			common.IPMASQAgentAddonName:                       getDefaultImage(common.IPMASQAgentAddonName, kubernetesImageBaseType),
			common.AzureNetworkPolicyAddonName:                azureNPMContainerImageReference,
			common.CalicoTyphaComponentName:                   calicoTyphaImageReference,
			common.CalicoCNIComponentName:                     calicoCNIImageReference,
			common.CalicoNodeComponentName:                    calicoNodeImageReference,
			common.CalicoPod2DaemonComponentName:              calicoPod2DaemonImageReference,
			common.CalicoClusterAutoscalerComponentName:       calicoClusterProportionalAutoscalerImageReference,
			common.CiliumAgentContainerName:                   ciliumAgentImageReference,
			common.CiliumCleanStateContainerName:              ciliumCleanStateImageReference,
			common.CiliumOperatorContainerName:                ciliumOperatorImageReference,
			common.CiliumEtcdOperatorContainerName:            ciliumEtcdOperatorImageReference,
			common.AntreaControllerContainerName:              antreaControllerImageReference,
			common.AntreaAgentContainerName:                   antreaAgentImageReference,
			common.AntreaOVSContainerName:                     antreaOVSImageReference,
			"antrea" + common.AntreaInstallCNIContainerName:   antreaInstallCNIImageReference,
			common.NMIContainerName:                           aadPodIdentityNMIImageReference,
			common.MICContainerName:                           aadPodIdentityMICImageReference,
			common.AzurePolicyAddonName:                       azurePolicyImageReference,
			common.GatekeeperContainerName:                    gatekeeperImageReference,
			common.NodeProblemDetectorAddonName:               nodeProblemDetectorImageReference,
			common.CSIProvisionerContainerName:                k8sComponent[common.CSIProvisionerContainerName],
			common.CSIAttacherContainerName:                   k8sComponent[common.CSIAttacherContainerName],
			common.CSILivenessProbeContainerName:              k8sComponent[common.CSILivenessProbeContainerName],
			common.CSILivenessProbeWindowsContainerName:       k8sComponent[common.CSILivenessProbeWindowsContainerName],
			common.CSINodeDriverRegistrarContainerName:        k8sComponent[common.CSINodeDriverRegistrarContainerName],
			common.CSINodeDriverRegistrarWindowsContainerName: k8sComponent[common.CSINodeDriverRegistrarWindowsContainerName],
			common.CSISnapshotterContainerName:                k8sComponent[common.CSISnapshotterContainerName],
			common.CSISnapshotControllerContainerName:         k8sComponent[common.CSISnapshotControllerContainerName],
			common.CSIResizerContainerName:                    k8sComponent[common.CSIResizerContainerName],
			common.CSIAzureDiskContainerName:                  k8sComponent[common.CSIAzureDiskContainerName],
			common.CSIAzureFileContainerName:                  csiAzureFileImageReference,
			common.KubeFlannelContainerName:                   kubeFlannelImageReference,
			"flannel" + common.FlannelInstallCNIContainerName: flannelInstallCNIImageReference,
			common.KubeRBACProxyContainerName:                 KubeRBACProxyImageReference,
			common.ScheduledMaintenanceManagerContainerName:   ScheduledMaintenanceManagerImageReference,
			"nodestatusfreq":                                  DefaultKubernetesNodeStatusUpdateFrequency,
			"nodegraceperiod":                                 DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
			"podeviction":                                     DefaultKubernetesCtrlMgrPodEvictionTimeout,
			"routeperiod":                                     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
			"backoffretries":                                  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
			"backoffjitter":                                   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
			"backoffduration":                                 strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
			"backoffexponent":                                 strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
			"ratelimitqps":                                    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
			"ratelimitqpswrite":                               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPSWrite, 'f', -1, 64),
			"ratelimitbucket":                                 strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
			"ratelimitbucketwrite":                            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucketWrite),
			"gchighthreshold":                                 strconv.Itoa(DefaultKubernetesGCHighThreshold),
			"gclowthreshold":                                  strconv.Itoa(DefaultKubernetesGCLowThreshold),
			common.NVIDIADevicePluginAddonName:                nvidiaDevicePluginImageReference,
			common.CSISecretsStoreProviderAzureContainerName:  csiSecretsStoreProviderAzureImageReference,
			common.CSISecretsStoreDriverContainerName:         csiSecretsStoreDriverImageReference,
			common.AzureArcOnboardingAddonName:                azureArcOnboardingImageReference,
			common.AzureKMSProviderComponentName:              azureKMSProviderImageReference,
		}
	case "1.28":
		ret = map[string]string{
			common.APIServerComponentName:                 getDefaultImage(common.APIServerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.ControllerManagerComponentName:         getDefaultImage(common.ControllerManagerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.KubeProxyAddonName:                     getDefaultImage(common.KubeProxyAddonName, kubernetesImageBaseType) + ":v" + version,
			common.SchedulerComponentName:                 getDefaultImage(common.SchedulerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.CloudControllerManagerComponentName:    "oss/kubernetes/azure-cloud-controller-manager:v1.28.5",
			common.CloudNodeManagerAddonName:              "oss/kubernetes/azure-cloud-node-manager:v1.28.5",
			common.WindowsArtifactComponentName:           "v" + version + "/windowszip/v" + version + "-1int.zip",
			common.WindowsArtifactAzureStackComponentName: "v" + version + "/windowszip/v" + version + "-1int.zip",
			common.DashboardAddonName:                     dashboardImageReference,
			common.DashboardMetricsScraperContainerName:   dashboardMetricsScraperImageReference,
			common.ExecHealthZComponentName:               getDefaultImage(common.ExecHealthZComponentName, kubernetesImageBaseType),
			common.AddonResizerComponentName:              k8sComponent[common.AddonResizerComponentName],
			common.MetricsServerAddonName:                 k8sComponent[common.MetricsServerAddonName],
			common.CoreDNSAddonName:                       getDefaultImage(common.CoreDNSAddonName, kubernetesImageBaseType),
			common.CoreDNSAutoscalerName:                  clusterProportionalAutoscalerImageReference,
			common.KubeDNSAddonName:                       getDefaultImage(common.KubeDNSAddonName, kubernetesImageBaseType),
			common.AddonManagerComponentName:              k8sComponent[common.AddonManagerComponentName],
			common.DNSMasqComponentName:                   getDefaultImage(common.DNSMasqComponentName, kubernetesImageBaseType),
			common.PauseComponentName:                     pauseImageReference,
			common.TillerAddonName:                        tillerImageReference,
			common.ReschedulerAddonName:                   getDefaultImage(common.ReschedulerAddonName, kubernetesImageBaseType),
			common.ACIConnectorAddonName:                  virtualKubeletImageReference,
			common.ClusterAutoscalerAddonName:             k8sComponent[common.ClusterAutoscalerAddonName],
			common.DNSSidecarComponentName:                getDefaultImage(common.DNSSidecarComponentName, kubernetesImageBaseType),

			common.SMBFlexVolumeAddonName:                     smbFlexVolumeImageReference,
			common.IPMASQAgentAddonName:                       getDefaultImage(common.IPMASQAgentAddonName, kubernetesImageBaseType),
			common.AzureNetworkPolicyAddonName:                azureNPMContainerImageReference,
			common.CalicoTyphaComponentName:                   calicoTyphaImageReference,
			common.CalicoCNIComponentName:                     calicoCNIImageReference,
			common.CalicoNodeComponentName:                    calicoNodeImageReference,
			common.CalicoPod2DaemonComponentName:              calicoPod2DaemonImageReference,
			common.CalicoClusterAutoscalerComponentName:       calicoClusterProportionalAutoscalerImageReference,
			common.CiliumAgentContainerName:                   ciliumAgentImageReference,
			common.CiliumCleanStateContainerName:              ciliumCleanStateImageReference,
			common.CiliumOperatorContainerName:                ciliumOperatorImageReference,
			common.CiliumEtcdOperatorContainerName:            ciliumEtcdOperatorImageReference,
			common.AntreaControllerContainerName:              antreaControllerImageReference,
			common.AntreaAgentContainerName:                   antreaAgentImageReference,
			common.AntreaOVSContainerName:                     antreaOVSImageReference,
			"antrea" + common.AntreaInstallCNIContainerName:   antreaInstallCNIImageReference,
			common.NMIContainerName:                           aadPodIdentityNMIImageReference,
			common.MICContainerName:                           aadPodIdentityMICImageReference,
			common.AzurePolicyAddonName:                       azurePolicyImageReference,
			common.GatekeeperContainerName:                    gatekeeperImageReference,
			common.NodeProblemDetectorAddonName:               nodeProblemDetectorImageReference,
			common.CSIProvisionerContainerName:                k8sComponent[common.CSIProvisionerContainerName],
			common.CSIAttacherContainerName:                   k8sComponent[common.CSIAttacherContainerName],
			common.CSILivenessProbeContainerName:              k8sComponent[common.CSILivenessProbeContainerName],
			common.CSILivenessProbeWindowsContainerName:       k8sComponent[common.CSILivenessProbeWindowsContainerName],
			common.CSINodeDriverRegistrarContainerName:        k8sComponent[common.CSINodeDriverRegistrarContainerName],
			common.CSINodeDriverRegistrarWindowsContainerName: k8sComponent[common.CSINodeDriverRegistrarWindowsContainerName],
			common.CSISnapshotterContainerName:                k8sComponent[common.CSISnapshotterContainerName],
			common.CSISnapshotControllerContainerName:         k8sComponent[common.CSISnapshotControllerContainerName],
			common.CSIResizerContainerName:                    k8sComponent[common.CSIResizerContainerName],
			common.CSIAzureDiskContainerName:                  k8sComponent[common.CSIAzureDiskContainerName],
			common.CSIAzureFileContainerName:                  csiAzureFileImageReference,
			common.KubeFlannelContainerName:                   kubeFlannelImageReference,
			"flannel" + common.FlannelInstallCNIContainerName: flannelInstallCNIImageReference,
			common.KubeRBACProxyContainerName:                 KubeRBACProxyImageReference,
			common.ScheduledMaintenanceManagerContainerName:   ScheduledMaintenanceManagerImageReference,
			"nodestatusfreq":                                  DefaultKubernetesNodeStatusUpdateFrequency,
			"nodegraceperiod":                                 DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
			"podeviction":                                     DefaultKubernetesCtrlMgrPodEvictionTimeout,
			"routeperiod":                                     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
			"backoffretries":                                  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
			"backoffjitter":                                   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
			"backoffduration":                                 strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
			"backoffexponent":                                 strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
			"ratelimitqps":                                    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
			"ratelimitqpswrite":                               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPSWrite, 'f', -1, 64),
			"ratelimitbucket":                                 strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
			"ratelimitbucketwrite":                            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucketWrite),
			"gchighthreshold":                                 strconv.Itoa(DefaultKubernetesGCHighThreshold),
			"gclowthreshold":                                  strconv.Itoa(DefaultKubernetesGCLowThreshold),
			common.NVIDIADevicePluginAddonName:                nvidiaDevicePluginImageReference,
			common.CSISecretsStoreProviderAzureContainerName:  csiSecretsStoreProviderAzureImageReference,
			common.CSISecretsStoreDriverContainerName:         csiSecretsStoreDriverImageReference,
			common.AzureArcOnboardingAddonName:                azureArcOnboardingImageReference,
			common.AzureKMSProviderComponentName:              azureKMSProviderImageReference,
		}
	case "1.27":
		ret = map[string]string{
			common.APIServerComponentName:                 getDefaultImage(common.APIServerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.ControllerManagerComponentName:         getDefaultImage(common.ControllerManagerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.KubeProxyAddonName:                     getDefaultImage(common.KubeProxyAddonName, kubernetesImageBaseType) + ":v" + version,
			common.SchedulerComponentName:                 getDefaultImage(common.SchedulerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.CloudControllerManagerComponentName:    "oss/kubernetes/azure-cloud-controller-manager:v1.27.13",
			common.CloudNodeManagerAddonName:              "oss/kubernetes/azure-cloud-node-manager:v1.27.13",
			common.WindowsArtifactComponentName:           "v" + version + "/windowszip/v" + version + "-1int.zip",
			common.WindowsArtifactAzureStackComponentName: "v" + version + "/windowszip/v" + version + "-1int.zip",
			common.DashboardAddonName:                     dashboardImageReference,
			common.DashboardMetricsScraperContainerName:   dashboardMetricsScraperImageReference,
			common.ExecHealthZComponentName:               getDefaultImage(common.ExecHealthZComponentName, kubernetesImageBaseType),
			common.AddonResizerComponentName:              k8sComponent[common.AddonResizerComponentName],
			common.MetricsServerAddonName:                 k8sComponent[common.MetricsServerAddonName],
			common.CoreDNSAddonName:                       getDefaultImage(common.CoreDNSAddonName, kubernetesImageBaseType),
			common.CoreDNSAutoscalerName:                  clusterProportionalAutoscalerImageReference,
			common.KubeDNSAddonName:                       getDefaultImage(common.KubeDNSAddonName, kubernetesImageBaseType),
			common.AddonManagerComponentName:              k8sComponent[common.AddonManagerComponentName],
			common.DNSMasqComponentName:                   getDefaultImage(common.DNSMasqComponentName, kubernetesImageBaseType),
			common.PauseComponentName:                     pauseImageReference,
			common.TillerAddonName:                        tillerImageReference,
			common.ReschedulerAddonName:                   getDefaultImage(common.ReschedulerAddonName, kubernetesImageBaseType),
			common.ACIConnectorAddonName:                  virtualKubeletImageReference,
			common.ClusterAutoscalerAddonName:             k8sComponent[common.ClusterAutoscalerAddonName],
			common.DNSSidecarComponentName:                getDefaultImage(common.DNSSidecarComponentName, kubernetesImageBaseType),

			common.SMBFlexVolumeAddonName:                     smbFlexVolumeImageReference,
			common.IPMASQAgentAddonName:                       getDefaultImage(common.IPMASQAgentAddonName, kubernetesImageBaseType),
			common.AzureNetworkPolicyAddonName:                azureNPMContainerImageReference,
			common.CalicoTyphaComponentName:                   calicoTyphaImageReference,
			common.CalicoCNIComponentName:                     calicoCNIImageReference,
			common.CalicoNodeComponentName:                    calicoNodeImageReference,
			common.CalicoPod2DaemonComponentName:              calicoPod2DaemonImageReference,
			common.CalicoClusterAutoscalerComponentName:       calicoClusterProportionalAutoscalerImageReference,
			common.CiliumAgentContainerName:                   ciliumAgentImageReference,
			common.CiliumCleanStateContainerName:              ciliumCleanStateImageReference,
			common.CiliumOperatorContainerName:                ciliumOperatorImageReference,
			common.CiliumEtcdOperatorContainerName:            ciliumEtcdOperatorImageReference,
			common.AntreaControllerContainerName:              antreaControllerImageReference,
			common.AntreaAgentContainerName:                   antreaAgentImageReference,
			common.AntreaOVSContainerName:                     antreaOVSImageReference,
			"antrea" + common.AntreaInstallCNIContainerName:   antreaInstallCNIImageReference,
			common.NMIContainerName:                           aadPodIdentityNMIImageReference,
			common.MICContainerName:                           aadPodIdentityMICImageReference,
			common.AzurePolicyAddonName:                       azurePolicyImageReference,
			common.GatekeeperContainerName:                    gatekeeperImageReference,
			common.NodeProblemDetectorAddonName:               nodeProblemDetectorImageReference,
			common.CSIProvisionerContainerName:                k8sComponent[common.CSIProvisionerContainerName],
			common.CSIAttacherContainerName:                   k8sComponent[common.CSIAttacherContainerName],
			common.CSILivenessProbeContainerName:              k8sComponent[common.CSILivenessProbeContainerName],
			common.CSILivenessProbeWindowsContainerName:       k8sComponent[common.CSILivenessProbeWindowsContainerName],
			common.CSINodeDriverRegistrarContainerName:        k8sComponent[common.CSINodeDriverRegistrarContainerName],
			common.CSINodeDriverRegistrarWindowsContainerName: k8sComponent[common.CSINodeDriverRegistrarWindowsContainerName],
			common.CSISnapshotterContainerName:                k8sComponent[common.CSISnapshotterContainerName],
			common.CSISnapshotControllerContainerName:         k8sComponent[common.CSISnapshotControllerContainerName],
			common.CSIResizerContainerName:                    k8sComponent[common.CSIResizerContainerName],
			common.CSIAzureDiskContainerName:                  k8sComponent[common.CSIAzureDiskContainerName],
			common.CSIAzureFileContainerName:                  csiAzureFileImageReference,
			common.KubeFlannelContainerName:                   kubeFlannelImageReference,
			"flannel" + common.FlannelInstallCNIContainerName: flannelInstallCNIImageReference,
			common.KubeRBACProxyContainerName:                 KubeRBACProxyImageReference,
			common.ScheduledMaintenanceManagerContainerName:   ScheduledMaintenanceManagerImageReference,
			"nodestatusfreq":                                  DefaultKubernetesNodeStatusUpdateFrequency,
			"nodegraceperiod":                                 DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
			"podeviction":                                     DefaultKubernetesCtrlMgrPodEvictionTimeout,
			"routeperiod":                                     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
			"backoffretries":                                  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
			"backoffjitter":                                   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
			"backoffduration":                                 strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
			"backoffexponent":                                 strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
			"ratelimitqps":                                    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
			"ratelimitqpswrite":                               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPSWrite, 'f', -1, 64),
			"ratelimitbucket":                                 strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
			"ratelimitbucketwrite":                            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucketWrite),
			"gchighthreshold":                                 strconv.Itoa(DefaultKubernetesGCHighThreshold),
			"gclowthreshold":                                  strconv.Itoa(DefaultKubernetesGCLowThreshold),
			common.NVIDIADevicePluginAddonName:                nvidiaDevicePluginImageReference,
			common.CSISecretsStoreProviderAzureContainerName:  csiSecretsStoreProviderAzureImageReference,
			common.CSISecretsStoreDriverContainerName:         csiSecretsStoreDriverImageReference,
			common.AzureArcOnboardingAddonName:                azureArcOnboardingImageReference,
			common.AzureKMSProviderComponentName:              azureKMSProviderImageReference,
		}
	default:
		ret = nil
	}
	for k, v := range overrides {
		ret[k] = v
	}
	return ret
}
