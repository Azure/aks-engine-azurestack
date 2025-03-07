// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

const (
	// Kubernetes is the string constant for the Kubernetes orchestrator type
	Kubernetes string = "Kubernetes"
)

const (
	// DefaultVNETCIDR is the default CIDR block for the VNET
	DefaultVNETCIDR = "10.0.0.0/8"
	// DefaultVNETCIDRIPv6 is the default IPv6 CIDR block for the VNET
	DefaultVNETCIDRIPv6 = "2001:1234:5678:9a00::/56"
	// DefaultInternalLbStaticIPOffset specifies the offset of the internal LoadBalancer's IP
	// address relative to the first consecutive Kubernetes static IP
	DefaultInternalLbStaticIPOffset = 10
	// NetworkPolicyNone is the string expression for the deprecated NetworkPolicy usage pattern "none"
	NetworkPolicyNone = "none"
	// NetworkPolicyCalico is the string expression for calico network policy config option
	NetworkPolicyCalico = "calico"
	// NetworkPolicyCilium is the string expression for cilium network policy config option
	NetworkPolicyCilium = "cilium"
	// NetworkPluginCilium is the string expression for cilium network plugin config option
	NetworkPluginCilium = NetworkPolicyCilium
	// NetworkPolicyAntrea is the string expression for antrea network policy config option
	NetworkPolicyAntrea = "antrea"
	// NetworkPluginAntrea is the string expression for antrea network plugin config option
	NetworkPluginAntrea = NetworkPolicyAntrea
	// NetworkPolicyAzure is the string expression for Azure CNI network policy manager
	NetworkPolicyAzure = "azure"
	// NetworkPluginAzure is the string expression for Azure CNI plugin
	NetworkPluginAzure = "azure"
	// NetworkPluginKubenet is the string expression for kubenet network plugin
	NetworkPluginKubenet = "kubenet"
	// NetworkPluginFlannel is the string expression for flannel network plugin
	NetworkPluginFlannel = "flannel"
	// DefaultGeneratorCode specifies the source generator of the cluster template.
	DefaultGeneratorCode = "aksengine"
	// DefaultKubernetesKubeletMaxPods is the max pods per kubelet
	DefaultKubernetesKubeletMaxPods = 110
	// DefaultMasterEtcdServerPort is the default etcd server port for Kubernetes master nodes
	DefaultMasterEtcdServerPort = 2380
	// DefaultMasterEtcdClientPort is the default etcd client port for Kubernetes master nodes
	DefaultMasterEtcdClientPort = 2379
	// etcdAccountNameFmt is the name format for a typical etcd account on Cosmos
	etcdAccountNameFmt = "%sk8s"
	// BasicLoadBalancerSku is the string const for Azure Basic Load Balancer
	BasicLoadBalancerSku = "Basic"
	// StandardLoadBalancerSku is the string const for Azure Standard Load Balancer
	StandardLoadBalancerSku = "Standard"
)

const (
	//DefaultExtensionsRootURL  Root URL for extensions
	DefaultExtensionsRootURL = "https://raw.githubusercontent.com/Azure/aks-engine/master/"
	// DefaultDockerEngineRepo for grabbing docker engine packages
	DefaultDockerEngineRepo = "https://download.docker.com/linux/ubuntu"
	// DefaultDockerComposeURL for grabbing docker images
	DefaultDockerComposeURL = "https://github.com/docker/compose/releases/download"
)

const (
	kubeConfigJSON = "k8s/kubeconfig.json"
	// Windows custom scripts. These should all be listed in template_generator.go:func GetKubernetesWindowsAgentFunctions
	kubernetesWindowsAgentCustomDataPS1           = "k8s/kuberneteswindowssetup.ps1"
	kubernetesWindowsAgentFunctionsPS1            = "k8s/kuberneteswindowsfunctions.ps1"
	kubernetesWindowsConfigFunctionsPS1           = "k8s/windowsconfigfunc.ps1"
	kubernetesWindowsContainerdFunctionsPS1       = "k8s/windowscontainerdfunc.ps1"
	kubernetesWindowsCsiProxyFunctionsPS1         = "k8s/windowscsiproxyfunc.ps1"
	kubernetesWindowsKubeletFunctionsPS1          = "k8s/windowskubeletfunc.ps1"
	kubernetesWindowsCniFunctionsPS1              = "k8s/windowscnifunc.ps1"
	kubernetesWindowsAzureCniFunctionsPS1         = "k8s/windowsazurecnifunc.ps1"
	kubernetesWindowsHostsConfigAgentFunctionsPS1 = "k8s/windowshostsconfigagentfunc.ps1"
	kubernetesWindowsOpenSSHFunctionPS1           = "k8s/windowsinstallopensshfunc.ps1"
	kubernetesWindowsHypervtemplatetoml           = "k8s/containerdtemplate.toml"
)

// cloud-init (i.e. ARM customData) source file references
const (
	kubernetesMasterNodeCustomDataYaml = "k8s/cloud-init/masternodecustomdata.yml"
	kubernetesNodeCustomDataYaml       = "k8s/cloud-init/nodecustomdata.yml"
	kubernetesJumpboxCustomDataYaml    = "k8s/cloud-init/jumpboxcustomdata.yml"
	kubernetesCSEMainScript            = "k8s/cloud-init/artifacts/cse_main.sh"
	kubernetesCSEHelpersScript         = "k8s/cloud-init/artifacts/cse_helpers.sh"
	kubernetesCSEInstall               = "k8s/cloud-init/artifacts/cse_install.sh"
	kubernetesCSEConfig                = "k8s/cloud-init/artifacts/cse_config.sh"
	kubernetesCISScript                = "k8s/cloud-init/artifacts/cis.sh"
	kubernetesCSECustomCloud           = "k8s/cloud-init/artifacts/cse_customcloud.sh"
	kubernetesHealthMonitorScript      = "k8s/cloud-init/artifacts/health-monitor.sh"
	// kubernetesKubeletMonitorSystemdTimer     = "k8s/cloud-init/artifacts/kubelet-monitor.timer" // TODO enable
	kubernetesKubeletMonitorSystemdService   = "k8s/cloud-init/artifacts/kubelet-monitor.service"
	apiServerAdmissionConfiguration          = "k8s/cloud-init/artifacts/apiserver-admission-control.yaml"
	apiserverMonitorSystemdService           = "k8s/cloud-init/artifacts/apiserver-monitor.service"
	kubernetesDockerMonitorSystemdService    = "k8s/cloud-init/artifacts/docker-monitor.service"
	etcdMonitorSystemdService                = "k8s/cloud-init/artifacts/etcd-monitor.service"
	labelNodesScript                         = "k8s/cloud-init/artifacts/label-nodes.sh"
	labelNodesSystemdService                 = "k8s/cloud-init/artifacts/label-nodes.service"
	untaintNodesScript                       = "k8s/cloud-init/artifacts/untaint-nodes.sh"
	untaintNodesSystemdService               = "k8s/cloud-init/artifacts/untaint-nodes.service"
	kubernetesMasterGenerateProxyCertsScript = "k8s/cloud-init/artifacts/generateproxycerts.sh"
	kubernetesCustomSearchDomainsScript      = "k8s/cloud-init/artifacts/setup-custom-search-domains.sh"
	kubeletSystemdService                    = "k8s/cloud-init/artifacts/kubelet.service"
	aptPreferences                           = "k8s/cloud-init/artifacts/apt-preferences"
	dockerClearMountPropagationFlags         = "k8s/cloud-init/artifacts/docker_clear_mount_propagation_flags.conf"
	systemdBPFMount                          = "k8s/cloud-init/artifacts/sys-fs-bpf.mount"
	etcdSystemdService                       = "k8s/cloud-init/artifacts/etcd.service"
	auditdRules                              = "k8s/cloud-init/artifacts/auditd-rules"
	// scripts and service for enabling ipv6 dual stack
	dhcpv6SystemdService      = "k8s/cloud-init/artifacts/dhcpv6.service"
	dhcpv6ConfigurationScript = "k8s/cloud-init/artifacts/enable-dhcpv6.sh"
	// script for getting key version from keyvault for kms
	kmsKeyvaultKeySystemdService = "k8s/cloud-init/artifacts/kms-keyvault-key.service"
	kmsKeyvaultKeyScript         = "k8s/cloud-init/artifacts/kms-keyvault-key.sh"
)

// cloud-init destination file references
const (
	apiServerAdmissionConfigurationFilepath    = "/etc/kubernetes/apiserver-admission-control.yaml"
	customCloudConfigCSEScriptFilepath         = "/opt/azure/containers/provision_configs_custom_cloud.sh"
	customCloudAzureCNIConfigCSEScriptFilepath = "/opt/azure/containers/provision_azurestack_cni.sh"
	cseHelpersScriptFilepath                   = "/opt/azure/containers/provision_source.sh"
	cseInstallScriptFilepath                   = "/opt/azure/containers/provision_installs.sh"
	cseConfigScriptFilepath                    = "/opt/azure/containers/provision_configs.sh"
	cseUbuntu2204StigScriptFilepath            = "/opt/azure/containers/provision_stig_ubuntu2204.sh"
	customSearchDomainsCSEScriptFilepath       = "/opt/azure/containers/setup-custom-search-domains.sh"
	dhcpV6ServiceCSEScriptFilepath             = "/etc/systemd/system/dhcpv6.service"
	dhcpV6ConfigCSEScriptFilepath              = "/opt/azure/containers/enable-dhcpv6.sh"
	kmsKeyvaultKeyServiceCSEScriptFilepath     = "/etc/systemd/system/kms-keyvault-key.service"
	kmsKeyvaultKeyCSEScriptFilepath            = "/opt/azure/containers/kms-keyvault-key.sh"
)

const (
	agentOutputs     = "agentoutputs.t"
	agentParams      = "agentparams.t"
	armParameters    = "k8s/armparameters.t"
	iaasOutputs      = "iaasoutputs.t"
	kubernetesParams = "k8s/kubernetesparams.t"
	masterOutputs    = "masteroutputs.t"
	masterParams     = "masterparams.t"
	windowsParams    = "windowsparams.t"
)

// addons source and destination file references
const (
	metricsServerAddonSourceFilename              string = "metrics-server.yaml"
	metricsServerAddonDestinationFilename         string = "metrics-server.yaml"
	tillerAddonSourceFilename                     string = "tiller.yaml"
	tillerAddonDestinationFilename                string = "tiller.yaml"
	aadPodIdentityAddonSourceFilename             string = "aad-pod-identity.yaml"
	aadPodIdentityAddonDestinationFilename        string = "aad-pod-identity.yaml"
	azureDiskCSIAddonSourceFilename               string = "azuredisk-csi-driver-deployment.yaml"
	azureDiskCSIAddonDestinationFilename          string = "azuredisk-csi-driver-deployment.yaml"
	azureFileCSIAddonSourceFilename               string = "azurefile-csi-driver-deployment.yaml"
	azureFileCSIAddonDestinationFilename          string = "azurefile-csi-driver-deployment.yaml"
	clusterAutoscalerAddonSourceFilename          string = "cluster-autoscaler.yaml"
	clusterAutoscalerAddonDestinationFilename     string = "cluster-autoscaler.yaml"
	smbFlexVolumeAddonSourceFilename              string = "smb-flexvolume.yaml"
	smbFlexVolumeAddonDestinationFilename         string = "smb-flexvolume.yaml"
	dashboardAddonSourceFilename                  string = "kubernetes-dashboard.yaml" // Deprecated
	dashboardAddonDestinationFilename             string = "kubernetes-dashboard.yaml" // Deprecated
	nvidiaAddonSourceFilename                     string = "nvidia-device-plugin.yaml"
	nvidiaAddonDestinationFilename                string = "nvidia-device-plugin.yaml"
	ipMasqAgentAddonSourceFilename                string = "ip-masq-agent.yaml"
	ipMasqAgentAddonDestinationFilename           string = "ip-masq-agent.yaml"
	calicoAddonSourceFilename                     string = "calico.yaml"
	calicoAddonDestinationFilename                string = "calico.yaml"
	azureNetworkPolicyAddonSourceFilename         string = "azure-network-policy.yaml"
	azureNetworkPolicyAddonDestinationFilename    string = "azure-network-policy.yaml"
	azurePolicyAddonSourceFilename                string = "azure-policy-deployment.yaml"
	azurePolicyAddonDestinationFilename           string = "azure-policy-deployment.yaml"
	cloudNodeManagerAddonSourceFilename           string = "cloud-node-manager.yaml"
	cloudNodeManagerAddonDestinationFilename      string = "cloud-node-manager.yaml"
	nodeProblemDetectorAddonSourceFilename        string = "node-problem-detector.yaml"
	nodeProblemDetectorAddonDestinationFilename   string = "node-problem-detector.yaml"
	kubeDNSAddonSourceFilename                    string = "kube-dns.yaml"
	kubeDNSAddonDestinationFilename               string = "kube-dns.yaml"
	corednsAddonSourceFilename                    string = "coredns.yaml"
	corednsAddonDestinationFilename               string = "coredns.yaml"
	kubeProxyAddonSourceFilename                  string = "kube-proxy.yaml"
	kubeProxyAddonDestinationFilename             string = "kube-proxy.yaml"
	podSecurityPolicyAddonSourceFilename          string = "pod-security-policy.yaml"
	podSecurityPolicyAddonDestinationFilename     string = "pod-security-policy.yaml"
	aadDefaultAdminGroupAddonSourceFilename       string = "aad-default-admin-group-rbac.yaml"
	aadDefaultAdminGroupDestinationFilename       string = "aad-default-admin-group-rbac.yaml"
	ciliumAddonSourceFilename                     string = "cilium.yaml"
	ciliumAddonDestinationFilename                string = "cilium.yaml"
	antreaAddonSourceFilename                     string = "antrea.yaml"
	antreaAddonDestinationFilename                string = "antrea.yaml"
	auditPolicyAddonSourceFilename                string = "audit-policy.yaml"
	auditPolicyAddonDestinationFilename           string = "audit-policy.yaml"
	cloudProviderAddonSourceFilename              string = "azure-cloud-provider.yaml"
	cloudProviderAddonDestinationFilename         string = "azure-cloud-provider.yaml"
	flannelAddonSourceFilename                    string = "flannel.yaml"
	flannelAddonDestinationFilename               string = "flannel.yaml"
	scheduledMaintenanceAddonSourceFilename       string = "scheduled-maintenance-deployment.yaml"
	scheduledMaintenanceAddonDestinationFilename  string = "scheduled-maintenance-deployment.yaml"
	secretsStoreCSIDriverAddonSourceFileName      string = "secrets-store-csi-driver.yaml"
	secretsStoreCSIDriverAddonDestinationFileName string = "secrets-store-csi-driver.yaml"
	connectedClusterAddonSourceFilename           string = "arc-onboarding.yaml"
	connectedClusterAddonDestinationFilename      string = "arc-onboarding.yaml"
)

// components source and destination file references
const (
	schedulerComponentSourceFilename                   string = "kubernetesmaster-kube-scheduler.yaml"
	schedulerComponentDestinationFilename              string = "kube-scheduler.yaml"
	controllerManagerComponentSourceFilename           string = "kubernetesmaster-kube-controller-manager.yaml"
	controllerManagerComponentDestinationFilename      string = "kube-controller-manager.yaml"
	cloudControllerManagerComponentSourceFilename      string = "kubernetesmaster-cloud-controller-manager.yaml"
	cloudControllerManagerComponentDestinationFilename string = "cloud-controller-manager.yaml"
	apiServerComponentSourceFilename                   string = "kubernetesmaster-kube-apiserver.yaml"
	apiServerComponentDestinationFilename              string = "kube-apiserver.yaml"
	addonManagerComponentSourceFilename                string = "kubernetesmaster-kube-addon-manager.yaml"
	addonManagerComponentDestinationFilename           string = "kube-addon-manager.yaml"
	clusterInitComponentDestinationFilename            string = "cluster-init.yaml"
	azureKMSComponentSourceFilename                    string = "kubernetesmaster-azure-kubernetes-kms.yaml"
	azureKMSComponentDestinationFilename               string = "kube-azure-kms.yaml"
)

const linuxCSELogPath string = "/var/log/azure/cluster-provision.log"
