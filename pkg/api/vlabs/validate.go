// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package vlabs

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers/to"
	"github.com/Azure/aks-engine-azurestack/pkg/versions"
	compute "github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/compute/armcompute"
	"github.com/blang/semver"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validate                       *validator.Validate
	keyvaultIDRegex                *regexp.Regexp
	labelValueRegex                *regexp.Regexp
	labelKeyRegex                  *regexp.Regexp
	diskEncryptionSetIDRegex       *regexp.Regexp
	proximityPlacementGroupIDRegex *regexp.Regexp
	// Any version has to be available in a container image from mcr.microsoft.com/oss/etcd-io/etcd:v[Version]
	etcdValidVersions = [...]string{"2.2.5", "2.3.0", "2.3.1", "2.3.2", "2.3.3", "2.3.4", "2.3.5", "2.3.6", "2.3.7", "2.3.8",
		"3.0.0", "3.0.1", "3.0.2", "3.0.3", "3.0.4", "3.0.5", "3.0.6", "3.0.7", "3.0.8", "3.0.9", "3.0.10", "3.0.11", "3.0.12", "3.0.13", "3.0.14", "3.0.15", "3.0.16", "3.0.17",
		"3.1.0", "3.1.1", "3.1.2", "3.1.2", "3.1.3", "3.1.4", "3.1.5", "3.1.6", "3.1.7", "3.1.8", "3.1.9", "3.1.10",
		"3.2.0", "3.2.1", "3.2.2", "3.2.3", "3.2.4", "3.2.5", "3.2.6", "3.2.7", "3.2.8", "3.2.9", "3.2.11", "3.2.12",
		"3.2.13", "3.2.14", "3.2.15", "3.2.16", "3.2.23", "3.2.24", "3.2.25", "3.2.26", "3.3.0", "3.3.1", "3.3.8", "3.3.9", "3.3.10", "3.3.13", "3.3.15", "3.3.18", "3.3.19", "3.3.22", "3.3.25"}
	containerdValidVersions              = [...]string{"1.3.2", "1.3.3", "1.3.4", "1.3.5", "1.3.6", "1.3.7", "1.3.8", "1.3.9", "1.4.4", "1.4.6", "1.4.7", "1.4.8", "1.4.9", "1.4.11", "1.5.11", "1.5.13", "1.5.16", "1.6.21", "1.6.28", "1.6.36"}
	kubernetesImageBaseTypeValidVersions = [...]string{"", common.KubernetesImageBaseTypeGCR, common.KubernetesImageBaseTypeMCR}
	cachingTypesValidValues              = [...]string{"", string(compute.CachingTypesNone), string(compute.CachingTypesReadWrite), string(compute.CachingTypesReadOnly)}
	linuxEth0MTUAllowedValues            = [...]int{1500, 3900}
	networkPluginPlusPolicyAllowed       = []k8sNetworkConfig{
		{
			networkPlugin: "",
			networkPolicy: "",
		},
		{
			networkPlugin: "azure",
			networkPolicy: "",
		},
		{
			networkPlugin: "azure",
			networkPolicy: "azure",
		},
		{
			networkPlugin: "kubenet",
			networkPolicy: "",
		},
		{
			networkPlugin: "flannel",
			networkPolicy: "",
		},
		{
			networkPlugin: NetworkPluginCilium,
			networkPolicy: NetworkPolicyCilium,
		},
		{
			networkPlugin: "kubenet",
			networkPolicy: "calico",
		},
		{
			networkPlugin: "azure",
			networkPolicy: "calico",
		},
		{
			networkPlugin: "",
			networkPolicy: "calico",
		},
		{
			networkPlugin: "",
			networkPolicy: NetworkPolicyCilium,
		},
		{
			networkPlugin: NetworkPluginAntrea,
			networkPolicy: NetworkPolicyAntrea,
		},
		{
			networkPlugin: "azure",
			networkPolicy: NetworkPolicyAntrea,
		},
		{
			networkPlugin: "",
			networkPolicy: NetworkPolicyAntrea,
		},
		{
			networkPlugin: "",
			networkPolicy: "azure", // for backwards-compatibility w/ prior networkPolicy usage
		},
		{
			networkPlugin: "",
			networkPolicy: "none", // for backwards-compatibility w/ prior networkPolicy usage
		},
	}
)

const (
	labelKeyPrefixMaxLength = 253
	labelValueFormat        = "^([A-Za-z0-9][-A-Za-z0-9_.]{0,61})?[A-Za-z0-9]$"
	labelKeyFormat          = "^(([a-zA-Z0-9-]+[.])*[a-zA-Z0-9-]+[/])?([A-Za-z0-9][-A-Za-z0-9_.]{0,61})?[A-Za-z0-9]$"
)

type k8sNetworkConfig struct {
	networkPlugin string
	networkPolicy string
}

func init() {
	validate = validator.New()
	keyvaultIDRegex = regexp.MustCompile(`^/subscriptions/\S+/resourceGroups/\S+/providers/Microsoft.KeyVault/vaults/[^/\s]+$`)
	labelValueRegex = regexp.MustCompile(labelValueFormat)
	labelKeyRegex = regexp.MustCompile(labelKeyFormat)
	diskEncryptionSetIDRegex = regexp.MustCompile(`^/subscriptions/\S+/resourceGroups/\S+/providers/Microsoft.Compute/diskEncryptionSets/[^/\s]+$`)
	proximityPlacementGroupIDRegex = regexp.MustCompile(`^/subscriptions/\S+/resourceGroups/\S+/providers/Microsoft.Compute/proximityPlacementGroups/[^/\s]+$`)
}

// Validate implements APIObject
func (a *Properties) validate(isUpdate bool) error {
	if e := validate.Struct(a); e != nil {
		return handleValidationErrors(e.(validator.ValidationErrors))
	}
	if e := a.ValidateOrchestratorProfile(isUpdate); e != nil {
		return e
	}
	if e := a.validateMasterProfile(isUpdate); e != nil {
		return e
	}
	if e := a.validateAgentPoolProfiles(isUpdate); e != nil {
		return e
	}
	if e := a.validateZones(); e != nil {
		return e
	}
	if e := a.validateLinuxProfile(); e != nil {
		return e
	}
	if e := a.validateAddons(isUpdate); e != nil {
		return e
	}
	if e := a.validateExtensions(); e != nil {
		return e
	}
	if e := a.validateVNET(); e != nil {
		return e
	}
	if e := a.validateServicePrincipalProfile(); e != nil {
		return e
	}

	if e := a.validateAADProfile(); e != nil {
		return e
	}

	if e := a.validateCustomKubeComponent(); e != nil {
		return e
	}

	if e := a.validateAzureStackSupport(); e != nil {
		return e
	}

	if e := a.validateWindowsProfile(isUpdate); e != nil {
		return e
	}
	return nil
}

func handleValidationErrors(e validator.ValidationErrors) error {
	// Override any version specific validation error message
	// common.HandleValidationErrors if the validation error message is general
	return common.HandleValidationErrors(e)
}

// ValidateOrchestratorProfile validates the orchestrator profile and the addons dependent on the version of the orchestrator
func (a *Properties) ValidateOrchestratorProfile(isUpdate bool) error {
	o := a.OrchestratorProfile
	// On updates we only need to make sure there is a supported patch version for the minor version
	if !isUpdate {
		version := common.RationalizeReleaseAndVersion(
			o.OrchestratorType,
			o.OrchestratorRelease,
			o.OrchestratorVersion,
			isUpdate,
			a.HasWindows(),
			a.IsAzureStackCloud())
		if a.IsAzureStackCloud() {
			if version == "" && a.HasWindows() {
				return errors.Errorf("the following OrchestratorProfile configuration is not supported on Azure Stack with OsType \"Windows\": OrchestratorType: \"%s\", OrchestratorRelease: \"%s\", OrchestratorVersion: \"%s\". Please use one of the following versions: %v", o.OrchestratorType, o.OrchestratorRelease, o.OrchestratorVersion, common.GetAllSupportedKubernetesVersions(false, true, true))
			} else if version == "" {
				return errors.Errorf("the following OrchestratorProfile configuration is not supported on Azure Stack: OrchestratorType: \"%s\", OrchestratorRelease: \"%s\", OrchestratorVersion: \"%s\". Please use one of the following versions: %v", o.OrchestratorType, o.OrchestratorRelease, o.OrchestratorVersion, common.GetAllSupportedKubernetesVersions(false, false, true))
			}
		} else {
			if version == "" && a.HasWindows() {
				return errors.Errorf("the following OrchestratorProfile configuration is not supported with OsType \"Windows\": OrchestratorType: \"%s\", OrchestratorRelease: \"%s\", OrchestratorVersion: \"%s\". Please use one of the following versions: %v", o.OrchestratorType, o.OrchestratorRelease, o.OrchestratorVersion, common.GetAllSupportedKubernetesVersions(false, true, false))
			} else if version == "" {
				return errors.Errorf("the following OrchestratorProfile configuration is not supported: OrchestratorType: \"%s\", OrchestratorRelease: \"%s\", OrchestratorVersion: \"%s\". Please use one of the following versions: %v", o.OrchestratorType, o.OrchestratorRelease, o.OrchestratorVersion, common.GetAllSupportedKubernetesVersions(false, false, false))
			}
		}

		sv, err := semver.Make(version)
		if err != nil {
			return errors.Errorf("could not validate version %s", version)
		}

		if a.HasAvailabilityZones() {
			minVersion, err := semver.Make("1.12.0")
			if err != nil {
				return errors.New("could not validate version")
			}

			if sv.LT(minVersion) {
				return errors.New("availabilityZone is only available in Kubernetes version 1.12 or greater")
			}
		}

		if o.KubernetesConfig != nil {
			err := o.KubernetesConfig.Validate(version, a.HasWindows(), a.FeatureFlags.IsIPv6DualStackEnabled(), a.FeatureFlags.IsIPv6OnlyEnabled(), isUpdate)
			if err != nil {
				return err
			}

			if o.KubernetesConfig.EnableAggregatedAPIs {
				if !o.KubernetesConfig.IsRBACEnabled() {
					return errors.New("enableAggregatedAPIs requires the enableRbac feature as a prerequisite")
				}
			}

			if to.Bool(o.KubernetesConfig.EnableDataEncryptionAtRest) {
				if o.KubernetesConfig.EtcdEncryptionKey != "" {
					_, err = base64.StdEncoding.DecodeString(o.KubernetesConfig.EtcdEncryptionKey)
					if err != nil {
						return errors.New("etcdEncryptionKey must be base64 encoded. Please provide a valid base64 encoded value or leave the etcdEncryptionKey empty to auto-generate the value")
					}
				}
			}

			if to.Bool(o.KubernetesConfig.EnableEncryptionWithExternalKms) {
				if to.Bool(a.OrchestratorProfile.KubernetesConfig.UseManagedIdentity) && a.OrchestratorProfile.KubernetesConfig.UserAssignedID == "" {
					log.Warnf("Clusters with enableEncryptionWithExternalKms=true and system-assigned identity are not upgradable! You will not be able to upgrade your cluster using `aks-engine-azurestack upgrade`")
				}
			}

			if to.Bool(o.KubernetesConfig.EnablePodSecurityPolicy) {
				log.Warnf("EnablePodSecurityPolicy is deprecated in favor of the addon pod-security-policy.")
				if !o.KubernetesConfig.IsRBACEnabled() {
					return errors.Errorf("enablePodSecurityPolicy requires the enableRbac feature as a prerequisite")
				}
				if len(o.KubernetesConfig.PodSecurityPolicyConfig) > 0 {
					log.Warnf("Raw manifest for PodSecurityPolicy using PodSecurityPolicyConfig is deprecated in favor of the addon pod-security-policy. This will be ignored.")
				}
			}

			if o.KubernetesConfig.LoadBalancerSku != "" {
				if !strings.EqualFold(o.KubernetesConfig.LoadBalancerSku, StandardLoadBalancerSku) && !strings.EqualFold(o.KubernetesConfig.LoadBalancerSku, BasicLoadBalancerSku) {
					return errors.Errorf("Invalid value for loadBalancerSku, only %s and %s are supported", StandardLoadBalancerSku, BasicLoadBalancerSku)
				}
			}

			if o.KubernetesConfig.LoadBalancerSku == StandardLoadBalancerSku {
				if !to.Bool(a.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB) {
					return errors.Errorf("standard loadBalancerSku should exclude master nodes. Please set KubernetesConfig \"ExcludeMasterFromStandardLB\" to \"true\"")
				}
			}

			if o.KubernetesConfig.LoadBalancerSku == BasicLoadBalancerSku {
				if o.KubernetesConfig.LoadBalancerOutboundIPs != nil {
					return errors.Errorf("kubernetesConfig.loadBalancerOutboundIPs configuration only supported for Standard loadBalancerSku=Standard")
				}
			}

			if o.KubernetesConfig.DockerEngineVersion != "" {
				log.Warnf("docker-engine is deprecated in favor of moby, but you passed in a dockerEngineVersion configuration. This will be ignored.")
			}

			if o.KubernetesConfig.MaximumLoadBalancerRuleCount < 0 {
				return errors.New("maximumLoadBalancerRuleCount shouldn't be less than 0")
			}

			if o.KubernetesConfig.LoadBalancerOutboundIPs != nil {
				if to.Int(o.KubernetesConfig.LoadBalancerOutboundIPs) > common.MaxLoadBalancerOutboundIPs {
					return errors.Errorf("kubernetesConfig.loadBalancerOutboundIPs was set to %d, the maximum allowed is %d", to.Int(o.KubernetesConfig.LoadBalancerOutboundIPs), common.MaxLoadBalancerOutboundIPs)
				}
			}

			// https://docs.microsoft.com/en-us/azure/load-balancer/load-balancer-outbound-rules-overview
			if o.KubernetesConfig.LoadBalancerSku == StandardLoadBalancerSku && o.KubernetesConfig.OutboundRuleIdleTimeoutInMinutes != 0 && (o.KubernetesConfig.OutboundRuleIdleTimeoutInMinutes < 4 || o.KubernetesConfig.OutboundRuleIdleTimeoutInMinutes > 120) {
				return errors.New("outboundRuleIdleTimeoutInMinutes shouldn't be less than 4 or greater than 120")
			}

			if a.IsAzureStackCloud() {
				if common.IsKubernetesVersionGe(a.OrchestratorProfile.OrchestratorVersion, "1.21.0") && !to.Bool(o.KubernetesConfig.UseCloudControllerManager) {
					return errors.New("useCloudControllerManager should be set to true for Kubernetes v1.21+ clusters on Azure Stack Hub")
				}

				if common.IsKubernetesVersionGe(a.OrchestratorProfile.OrchestratorVersion, "1.24.0") && o.KubernetesConfig.ContainerRuntime == Docker {
					return errors.Errorf("Docker runtime is no longer supported for v1.24+ clusters, use %s containerRuntime value instead", Containerd)
				}

				if to.Bool(o.KubernetesConfig.UseInstanceMetadata) {
					return errors.New("useInstanceMetadata shouldn't be set to true as feature not yet supported on Azure Stack")
				}

				if o.KubernetesConfig.EtcdDiskSizeGB != "" {
					etcdDiskSizeGB, err := strconv.Atoi(o.KubernetesConfig.EtcdDiskSizeGB)
					if err != nil {
						return errors.Errorf("could not convert EtcdDiskSizeGB to int")
					}
					if etcdDiskSizeGB > MaxAzureStackManagedDiskSize {
						return errors.Errorf("EtcdDiskSizeGB max size supported on Azure Stack is %d", MaxAzureStackManagedDiskSize)
					}
				}
			}

			if o.KubernetesConfig.EtcdStorageLimitGB != 0 {
				if o.KubernetesConfig.EtcdStorageLimitGB > 8 {
					log.Warnf("EtcdStorageLimitGB of %d is larger than the recommended maximum of 8", o.KubernetesConfig.EtcdStorageLimitGB)
				}
				if o.KubernetesConfig.EtcdStorageLimitGB < 2 {
					return errors.Errorf("EtcdStorageLimitGB value of %d is too small, the minimum allowed is 2", o.KubernetesConfig.EtcdStorageLimitGB)
				}
			}
		}
	} else {
		version := common.RationalizeReleaseAndVersion(
			o.OrchestratorType,
			o.OrchestratorRelease,
			o.OrchestratorVersion,
			false,
			a.HasWindows(),
			a.IsAzureStackCloud())
		if version == "" {
			patchVersion := common.GetValidPatchVersion(o.OrchestratorType, o.OrchestratorVersion, isUpdate, a.HasWindows(), a.IsAzureStackCloud())
			// if there isn't a supported patch version for this version fail
			if patchVersion == "" {
				if a.HasWindows() {
					return errors.Errorf("the following OrchestratorProfile configuration is not supported with Windows agentpools: OrchestratorType: \"%s\", OrchestratorRelease: \"%s\", OrchestratorVersion: \"%s\". Please check supported Release or Version for this build of aks-engine", o.OrchestratorType, o.OrchestratorRelease, o.OrchestratorVersion)
				}
				return errors.Errorf("the following OrchestratorProfile configuration is not supported: OrchestratorType: \"%s\", OrchestratorRelease: \"%s\", OrchestratorVersion: \"%s\". Please check supported Release or Version for this build of aks-engine", o.OrchestratorType, o.OrchestratorRelease, o.OrchestratorVersion)
			}
		}
	}

	if a.HasFlatcar() && o.KubernetesConfig.NetworkPlugin == "azure" && o.KubernetesConfig.NetworkMode == NetworkModeBridge {
		return errors.Errorf("Flatcar node pools require 'transparent' networkMode with Azure CNI")
	}

	return a.validateContainerRuntime(isUpdate)
}

func (a *Properties) validateMasterProfile(isUpdate bool) error {
	m := a.MasterProfile

	if m.Count == 1 && !isUpdate {
		log.Warnf("Running only 1 control plane VM not recommended for production clusters, use 3 or 5 for control plane redundancy")
	}
	if m.IsVirtualMachineScaleSets() && m.VnetSubnetID != "" && m.FirstConsecutiveStaticIP != "" {
		return errors.New("when masterProfile's availabilityProfile is VirtualMachineScaleSets and a vnetSubnetID is specified, the firstConsecutiveStaticIP should be empty and will be determined by an offset from the first IP in the vnetCidr")
	}

	if m.ImageRef != nil {
		if err := m.ImageRef.validateImageNameAndGroup(); err != nil {
			return err
		}
	}

	if m.IsVirtualMachineScaleSets() {
		if !isUpdate {
			log.Warnf("Clusters with a VMSS control plane are not upgradable! You will not be able to upgrade your cluster using `aks-engine-azurestack upgrade`")
		}
		e := validateVMSS(a.OrchestratorProfile, false, m.StorageProfile, a.HasWindows(), a.IsAzureStackCloud())
		if e != nil {
			return e
		}
		if !a.IsClusterAllVirtualMachineScaleSets() {
			return errors.New("VirtualMachineScaleSets for master profile must be used together with virtualMachineScaleSets for agent profiles. Set \"availabilityProfile\" to \"VirtualMachineScaleSets\" for agent profiles")
		}

		if a.OrchestratorProfile.KubernetesConfig != nil && to.Bool(a.OrchestratorProfile.KubernetesConfig.UseManagedIdentity) && a.OrchestratorProfile.KubernetesConfig.UserAssignedID == "" {
			return errors.New("virtualMachineScaleSets for master profile can be used only with user assigned MSI ! Please specify \"userAssignedID\" in \"kubernetesConfig\"")
		}
	}
	if m.SinglePlacementGroup != nil && m.AvailabilityProfile == AvailabilitySet {
		return errors.New("singlePlacementGroup is only supported with VirtualMachineScaleSets")
	}

	if e := validateProximityPlacementGroupID(m.ProximityPlacementGroupID); e != nil {
		return e
	}

	distroValues := DistroValues
	if isUpdate {
		distroValues = append(distroValues, AKSDockerEngine, AKS1604Deprecated, AKS1804Deprecated)
	}
	if !validateDistro(m.Distro, distroValues) {
		switch m.Distro {
		case AKSDockerEngine, AKS1604Deprecated:
			return errors.Errorf("The %s distro is deprecated, please use %s instead", m.Distro, AKSUbuntu1604)
		case AKS1804Deprecated:
			return errors.Errorf("The %s distro is deprecated, please use %s instead", m.Distro, AKSUbuntu1804)
		default:
			return errors.Errorf("The %s distro is not supported", m.Distro)
		}
	}

	if to.Bool(m.AuditDEnabled) {
		if m.Distro != "" && !m.IsUbuntu() {
			return errors.Errorf("auditd was enabled for master vms, but an Ubuntu-based distro was not selected")
		}
	} else {
		if a.FeatureFlags.IsEnforceUbuntuDisaStigEnabled() && m.Distro != "" && m.IsUbuntu() {
			return errors.New("AuditD should be enabled in all Ubuntu-based pools if feature flag 'EnforceUbuntu2004DisaStig' or 'EnforceUbuntu2204DisaStig' is set")
		}
	}

	var validOSDiskCachingType bool
	for _, valid := range cachingTypesValidValues {
		if valid == m.OSDiskCachingType {
			validOSDiskCachingType = true
		}
	}
	if !validOSDiskCachingType {
		return errors.Errorf("Invalid masterProfile osDiskCachingType value \"%s\", please use one of the following versions: %s", m.OSDiskCachingType, cachingTypesValidValues)
	}

	return common.ValidateDNSPrefix(m.DNSPrefix)
}

func (a *Properties) validateAgentPoolProfiles(isUpdate bool) error {

	profileNames := make(map[string]bool)
	for i, agentPoolProfile := range a.AgentPoolProfiles {
		if e := validatePoolName(agentPoolProfile.Name); e != nil {
			return e
		}

		// validate os type is linux if dual stack feature is enabled
		if a.FeatureFlags.IsIPv6DualStackEnabled() || a.FeatureFlags.IsIPv6OnlyEnabled() {
			if agentPoolProfile.OSType == Windows {
				if a.FeatureFlags.IsIPv6DualStackEnabled() && !common.IsKubernetesVersionGe(a.OrchestratorProfile.OrchestratorVersion, "1.19.0") {
					return errors.Errorf("Dual stack IPv6 feature is supported on Windows only from Kubernetes version 1.19, but OrchestratorProfile.OrchestratorVersion is '%s'", a.OrchestratorProfile.OrchestratorVersion)
				}
				if a.FeatureFlags.IsIPv6OnlyEnabled() {
					return errors.Errorf("Single stack IPv6 feature is supported only with Linux, but agent pool '%s' is of os type %s", agentPoolProfile.Name, agentPoolProfile.OSType)
				}
			}
			if agentPoolProfile.Distro == Flatcar {
				return errors.Errorf("Dual stack and single stack IPv6 feature is currently supported only with Ubuntu, but agent pool '%s' is of distro type %s", agentPoolProfile.Name, agentPoolProfile.Distro)
			}
		}

		// validate that each AgentPoolProfile Name is unique
		if _, ok := profileNames[agentPoolProfile.Name]; ok {
			return errors.Errorf("profile name '%s' already exists, profile names must be unique across pools", agentPoolProfile.Name)
		}
		profileNames[agentPoolProfile.Name] = true

		if e := validatePoolOSType(agentPoolProfile.OSType); e != nil {
			return e
		}

		if to.Bool(agentPoolProfile.AcceleratedNetworkingEnabled) || to.Bool(agentPoolProfile.AcceleratedNetworkingEnabledWindows) {
			if a.IsAzureStackCloud() {
				return errors.Errorf("AcceleratedNetworkingEnabled or AcceleratedNetworkingEnabledWindows shouldn't be set to true as feature is not yet supported on Azure Stack")
			} else if e := validatePoolAcceleratedNetworking(agentPoolProfile.VMSize); e != nil {
				return e
			}
		}

		if to.Bool(agentPoolProfile.VMSSOverProvisioningEnabled) {
			if agentPoolProfile.AvailabilityProfile == AvailabilitySet {
				return errors.Errorf("You have specified VMSS Overprovisioning in agent pool %s, but you did not specify VMSS", agentPoolProfile.Name)
			}
		}

		if to.Bool(agentPoolProfile.AuditDEnabled) {
			if agentPoolProfile.Distro != "" && !agentPoolProfile.IsUbuntu() {
				return errors.Errorf("You have enabled auditd in agent pool %s, but you did not specify an Ubuntu-based distro", agentPoolProfile.Name)
			}
		} else {
			if a.FeatureFlags.IsEnforceUbuntuDisaStigEnabled() && agentPoolProfile.IsUbuntu() {
				return errors.New("AuditD should be enabled in all Ubuntu-based pools if feature flag 'EnforceUbuntu2004DisaStig' or 'EnforceUbuntu2204DisaStig' is set")
			}
		}

		if to.Bool(agentPoolProfile.EnableVMSSNodePublicIP) {
			if agentPoolProfile.AvailabilityProfile == AvailabilitySet {
				return errors.Errorf("You have enabled VMSS node public IP in agent pool %s, but you did not specify VMSS", agentPoolProfile.Name)
			}
			if !strings.EqualFold(a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku, BasicLoadBalancerSku) {
				return errors.Errorf("You have enabled VMSS node public IP in agent pool %s, but you did not specify Basic Load Balancer SKU", agentPoolProfile.Name)
			}
		}

		if e := agentPoolProfile.validateOrchestratorSpecificProperties(); e != nil {
			return e
		}

		if agentPoolProfile.ImageRef != nil {
			if e := agentPoolProfile.ImageRef.validateImageNameAndGroup(); e != nil {
				return e
			}
		}

		if e := agentPoolProfile.validateAvailabilityProfile(); e != nil {
			return e
		}

		if e := agentPoolProfile.validateRoles(); e != nil {
			return e
		}

		if e := agentPoolProfile.validateCustomNodeLabels(); e != nil {
			return e
		}

		if agentPoolProfile.AvailabilityProfile != AvailabilitySet {
			e := validateVMSS(a.OrchestratorProfile, isUpdate, agentPoolProfile.StorageProfile, a.HasWindows(), a.IsAzureStackCloud())
			if e != nil {
				return e
			}
		}

		if a.AgentPoolProfiles[i].AvailabilityProfile != a.AgentPoolProfiles[0].AvailabilityProfile {
			return errors.New("mixed mode availability profiles are not allowed. Please set either VirtualMachineScaleSets or AvailabilitySet in availabilityProfile for all agent pools")
		}

		if a.AgentPoolProfiles[i].SinglePlacementGroup != nil && a.AgentPoolProfiles[i].AvailabilityProfile == AvailabilitySet {
			return errors.New("singlePlacementGroup is only supported with VirtualMachineScaleSets")
		}

		distroValues := DistroValues
		if isUpdate {
			distroValues = append(distroValues, AKSDockerEngine, AKS1604Deprecated, AKS1804Deprecated)
		}
		if !validateDistro(agentPoolProfile.Distro, distroValues) {
			switch agentPoolProfile.Distro {
			case AKSDockerEngine, AKS1604Deprecated:
				return errors.Errorf("The %s distro is deprecated, please use %s instead", agentPoolProfile.Distro, AKSUbuntu1604)
			case AKS1804Deprecated:
				return errors.Errorf("The %s distro is deprecated, please use %s instead", agentPoolProfile.Distro, AKSUbuntu1804)
			default:
				return errors.Errorf("The %s distro is not supported", agentPoolProfile.Distro)
			}
		}

		if e := agentPoolProfile.validateLoadBalancerBackendAddressPoolIDs(); e != nil {
			return e
		}

		if agentPoolProfile.IsEphemeral() {
			log.Warnf("Ephemeral disks are enabled for Agent Pool %s. This feature in AKS-Engine is experimental, and data could be lost in some cases.", agentPoolProfile.Name)
		}

		if e := validateProximityPlacementGroupID(agentPoolProfile.ProximityPlacementGroupID); e != nil {
			return e
		}
		var validOSDiskCachingType, validDataDiskCachingType bool
		for _, valid := range cachingTypesValidValues {
			if valid == agentPoolProfile.OSDiskCachingType {
				validOSDiskCachingType = true
			}
			if valid == agentPoolProfile.DataDiskCachingType {
				validDataDiskCachingType = true
			}
		}
		if !validOSDiskCachingType {
			return errors.Errorf("Invalid osDiskCachingType value \"%s\" for agentPoolProfile \"%s\", please use one of the following versions: %s", agentPoolProfile.OSDiskCachingType, agentPoolProfile.Name, cachingTypesValidValues)
		}
		if !validDataDiskCachingType {
			return errors.Errorf("Invalid dataDiskCachingType value \"%s\" for agentPoolProfile \"%s\", please use one of the following versions: %s", agentPoolProfile.DataDiskCachingType, agentPoolProfile.Name, cachingTypesValidValues)
		}
		if agentPoolProfile.IsEphemeral() {
			if agentPoolProfile.OSDiskCachingType != "" && agentPoolProfile.OSDiskCachingType != string(compute.CachingTypesReadOnly) {
				return errors.Errorf("Invalid osDiskCachingType value \"%s\" for agentPoolProfile \"%s\" using Ephemeral Disk, you must use: %s", agentPoolProfile.OSDiskCachingType, agentPoolProfile.Name, string(compute.CachingTypesReadOnly))
			}
		}
	}

	return nil
}

func (a *Properties) validateZones() error {
	if a.HasAvailabilityZones() {
		var poolsWithZones, poolsWithoutZones []string
		for _, pool := range a.AgentPoolProfiles {
			if pool.HasAvailabilityZones() {
				poolsWithZones = append(poolsWithZones, pool.Name)
			} else {
				poolsWithoutZones = append(poolsWithoutZones, pool.Name)
			}
		}
		if !a.MastersAndAgentsUseAvailabilityZones() {
			poolsWithZonesPrefix := "pool"
			poolsWithoutZonesPrefix := "pool"
			if len(poolsWithZones) > 1 {
				poolsWithZonesPrefix = "pools"
			}
			if len(poolsWithoutZones) > 1 {
				poolsWithoutZonesPrefix = "pools"
			}
			poolsWithZonesString := helpers.GetEnglishOrderedQuotedListWithOxfordCommas(poolsWithZones)
			poolsWithoutZonesString := helpers.GetEnglishOrderedQuotedListWithOxfordCommas(poolsWithoutZones)
			if !a.MasterProfile.HasAvailabilityZones() {
				if len(poolsWithZones) == len(a.AgentPoolProfiles) {
					log.Warnf("This cluster is using Availability Zones for %s %s, but not for master VMs", poolsWithZonesPrefix, poolsWithZonesString)
				} else {
					log.Warnf("This cluster is using Availability Zones for %s %s, but not for %s %s, nor for master VMs", poolsWithZonesPrefix, poolsWithZonesString, poolsWithoutZonesPrefix, poolsWithoutZonesString)
				}
			} else {
				if len(poolsWithoutZones) > 0 {
					log.Warnf("This cluster is using Availability Zones for master VMs, but not for %s %s", poolsWithoutZonesPrefix, poolsWithoutZonesString)
				}
			}
		} else {
			// agent pool profiles
			for _, agentPoolProfile := range a.AgentPoolProfiles {
				if agentPoolProfile.AvailabilityProfile == AvailabilitySet {
					return errors.New("Availability Zones are not supported with an AvailabilitySet. Please either remove availabilityProfile or set availabilityProfile to VirtualMachineScaleSets")
				}
			}
			if a.OrchestratorProfile.KubernetesConfig != nil && a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku != "" && !strings.EqualFold(a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku, StandardLoadBalancerSku) {
				return errors.New("Availability Zones requires Standard LoadBalancer. Please set KubernetesConfig \"LoadBalancerSku\" to \"Standard\"")
			}
		}
	}
	return nil
}

func (a *Properties) validateLinuxProfile() error {
	var validEth0MTU bool
	if a.LinuxProfile.Eth0MTU != 0 {
		if a.OrchestratorProfile != nil &&
			a.OrchestratorProfile.KubernetesConfig != nil &&
			a.OrchestratorProfile.KubernetesConfig.NetworkPlugin == NetworkPluginKubenet {
			return errors.Errorf("Custom linuxProfile eth0MTU value not allowed when using Kubenet")
		}
		for _, valid := range linuxEth0MTUAllowedValues {
			if valid == a.LinuxProfile.Eth0MTU {
				validEth0MTU = true
				break
			}
		}
		if !validEth0MTU {
			allowedMTUs := ""
			for _, mtu := range linuxEth0MTUAllowedValues {
				allowedMTUs += strconv.Itoa(mtu) + ", "
			}
			allowedMTUs = strings.TrimRight(allowedMTUs, ", ")
			return errors.Errorf("Invalid linuxProfile eth0MTU value \"%d\", please use one of the following values: %s", a.LinuxProfile.Eth0MTU, allowedMTUs)
		}
	}
	for _, publicKey := range a.LinuxProfile.SSH.PublicKeys {
		if e := validate.Var(publicKey.KeyData, "required"); e != nil {
			return errors.New("KeyData in LinuxProfile.SSH.PublicKeys cannot be empty string")
		}
	}
	if a.LinuxProfile.EnableUnattendedUpgrades == nil {
		log.Warnf("linuxProfile.enableUnattendedUpgrades configuration was not declared, your cluster nodes will be configured to run unattended-upgrade by default")
	}
	return validateKeyVaultSecrets(a.LinuxProfile.Secrets, false)
}

func (a *Properties) validateAddons(isUpdate bool) error {
	if a.OrchestratorProfile.KubernetesConfig != nil && a.OrchestratorProfile.KubernetesConfig.Addons != nil {
		var isAvailabilitySets bool
		var kubeDNSEnabled bool
		var corednsEnabled bool

		for _, agentPool := range a.AgentPoolProfiles {
			if agentPool.IsAvailabilitySets() {
				isAvailabilitySets = true
			}
		}
		for _, addon := range a.OrchestratorProfile.KubernetesConfig.Addons {
			if addon.Data != "" {
				if len(addon.Config) > 0 || len(addon.Containers) > 0 {
					return errors.New("Config and containers should be empty when addon.Data is specified")
				}
				if _, err := base64.StdEncoding.DecodeString(addon.Data); err != nil {
					return errors.Errorf("Addon %s's data should be base64 encoded", addon.Name)
				}
			}

			if addon.Mode != "" {
				if addon.Mode != AddonModeEnsureExists && addon.Mode != AddonModeReconcile {
					return errors.Errorf("addon %s has a mode configuration '%s', must be either %s or %s", addon.Name, addon.Mode, AddonModeEnsureExists, AddonModeReconcile)
				}
			}

			// Validation for addons if they are enabled
			if to.Bool(addon.Enabled) {
				switch addon.Name {
				case "cluster-autoscaler":
					if isAvailabilitySets {
						return errors.Errorf("cluster-autoscaler addon can only be used with VirtualMachineScaleSets. Please specify \"availabilityProfile\": \"%s\"", VirtualMachineScaleSets)
					}
					for _, pool := range addon.Pools {
						if pool.Name == "" {
							return errors.Errorf("cluster-autoscaler addon pools configuration must have a 'name' property that correlates with a pool name in the agentPoolProfiles array")
						}
						if a.GetAgentPoolByName(pool.Name) == nil {
							return errors.Errorf("cluster-autoscaler addon pool 'name' %s does not match any agentPoolProfiles nodepool name", pool.Name)
						}
						if pool.Config != nil {
							var min, max int
							var err error
							if pool.Config["min-nodes"] != "" {
								min, err = strconv.Atoi(pool.Config["min-nodes"])
								if err != nil {
									return errors.Errorf("cluster-autoscaler addon pool 'name' %s has invalid 'min-nodes' config, must be a string int, got %s", pool.Name, pool.Config["min-nodes"])
								}
							}
							if pool.Config["max-nodes"] != "" {
								max, err = strconv.Atoi(pool.Config["max-nodes"])
								if err != nil {
									return errors.Errorf("cluster-autoscaler addon pool 'name' %s has invalid 'max-nodes' config, must be a string int, got %s", pool.Name, pool.Config["max-nodes"])
								}
							}
							if min > max {
								return errors.Errorf("cluster-autoscaler addon pool 'name' %s has invalid config, 'max-nodes' %d must be greater than or equal to 'min-nodes' %d", pool.Name, max, min)
							}
						}
					}
				case "aad":
					if !a.HasAADAdminGroupID() {
						return errors.New("aad addon can't be enabled without a valid aadProfile w/ adminGroupID")
					}
				case "appgw-ingress":
					if (a.ServicePrincipalProfile == nil || len(a.ServicePrincipalProfile.ObjectID) == 0) &&
						!to.Bool(a.OrchestratorProfile.KubernetesConfig.UseManagedIdentity) {
						return errors.New("appgw-ingress add-ons requires 'objectID' to be specified or UseManagedIdentity to be true")
					}

					if a.OrchestratorProfile.KubernetesConfig.NetworkPlugin != "azure" {
						return errors.New("appgw-ingress add-ons can only be used with Network Plugin as 'azure'")
					}

					if len(addon.Config["appgw-subnet"]) == 0 {
						return errors.New("appgw-ingress add-ons requires 'appgw-subnet' in the Config. It is used to provision the subnet for Application Gateway in the vnet")
					}
				case "cloud-node-manager":
					if !to.Bool(a.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager) {
						return errors.Errorf("%s add-on requires useCloudControllerManager to be true", addon.Name)
					}
					if !a.ShouldEnableAzureCloudAddon(addon.Name) {
						minVersion := "1.16.0"
						if a.HasWindows() {
							minVersion = "1.18.0"
						}
						return errors.Errorf("%s add-on can only be used Kubernetes %s or above", addon.Name, minVersion)
					}
				case common.CiliumAddonName:
					if !common.IsKubernetesVersionGe(a.OrchestratorProfile.OrchestratorVersion, "1.16.0") {
						if a.OrchestratorProfile.KubernetesConfig.NetworkPolicy != NetworkPolicyCilium {
							return errors.Errorf("%s addon may only be enabled if the networkPolicy=%s", common.CiliumAddonName, NetworkPolicyCilium)
						}
					} else {
						return errors.Errorf("%s addon is not supported on Kubernetes v1.16.0 or greater", common.CiliumAddonName)
					}
				case common.AntreaAddonName:
					if a.OrchestratorProfile.KubernetesConfig.NetworkPolicy != NetworkPolicyAntrea {
						return errors.Errorf("%s addon may only be enabled if the networkPolicy=%s", common.AntreaAddonName, NetworkPolicyAntrea)
					}
				case common.FlannelAddonName:
					if isUpdate {
						if a.OrchestratorProfile.KubernetesConfig.NetworkPolicy != "" {
							return errors.Errorf("%s addon does not support NetworkPolicy, replace %s with \"\"", common.FlannelAddonName, a.OrchestratorProfile.KubernetesConfig.NetworkPolicy)
						}
						networkPlugin := a.OrchestratorProfile.KubernetesConfig.NetworkPlugin
						if networkPlugin != "" {
							if networkPlugin != NetworkPluginFlannel {
								return errors.Errorf("%s addon is not supported with networkPlugin=%s, please use networkPlugin=%s", common.FlannelAddonName, networkPlugin, NetworkPluginFlannel)
							}
						}
						if a.OrchestratorProfile.KubernetesConfig.ContainerRuntime != Containerd {
							return errors.Errorf("%s addon is only supported with containerRuntime=%s", common.FlannelAddonName, Containerd)
						}
					} else {
						return errors.Errorf("%s addon is deprecated for new clusters", common.FlannelAddonName)
					}
				case common.KubeDNSAddonName:
					kubeDNSEnabled = true
				case common.CoreDNSAddonName:
					corednsEnabled = true
				case common.SecretsStoreCSIDriverAddonName:
					if !common.IsKubernetesVersionGe(a.OrchestratorProfile.OrchestratorVersion, "1.16.0") {
						return errors.Errorf("%s add-on can only be used in 1.16+", addon.Name)
					}
				case common.PodSecurityPolicyAddonName:
					if common.ShouldDisablePodSecurityPolicyAddon(a.OrchestratorProfile.OrchestratorVersion) {
						log.Warn("The PodSecurityPolicy admission was removed in Kubernetes v1.25+. " +
							"The pod security standards will be enforced by the built-in PodSecurity admission controller instead. " +
							"See https://github.com/Azure/aks-engine-azurestack/blob/master/docs/topics/pod-security.md")
					}
				case common.AzureArcOnboardingAddonName:
					if err := addon.validateArcAddonConfig(); err != nil {
						return err
					}
				case common.ReschedulerAddonName:
					if isUpdate {
						log.Warnf("The rescheduler addon has been deprecated and disabled, it will be removed during this update")
					}
					return errors.Errorf("The rescheduler addon has been deprecated and disabled, please remove it from your cluster configuration before creating a new cluster")
				case common.ContainerMonitoringAddonName:
					if isUpdate {
						log.Warnf("The container monitoring addon has been deprecated and disabled, it will be removed during this update")
					}
					return errors.Errorf("The container monitoring addon has been deprecated and disabled, please remove it from your cluster configuration before creating a new cluster")
				case common.DashboardAddonName:
					log.Warnf("The kube-dashboard addon is deprecated, we recommend you install the dashboard yourself, see https://github.com/kubernetes/dashboard")
				case common.AzureCNINetworkMonitorAddonName:
					if isUpdate {
						log.Warnf("The Azure CNI networkmonitor addon has been deprecated, it will be marked as disabled")
					}
				}
			} else {
				// Validation for addons if they are disabled
				switch addon.Name {
				case "cloud-node-manager":
					if a.ShouldEnableAzureCloudAddon(addon.Name) && !a.IsAzureStackCloud() {
						minVersion := "1.16.0"
						if a.HasWindows() {
							minVersion = "1.18.0"
						}
						return errors.Errorf("%s add-on is required when useCloudControllerManager is true in Kubernetes %s or above", addon.Name, minVersion)
					}
				case common.AzureCloudProviderAddonName:
					return errors.Errorf("%s add-on is required, it cannot be disabled", addon.Name)
				}
			}
		}
		if kubeDNSEnabled && corednsEnabled {
			return errors.New("Both kube-dns and coredns addons are enabled, only one of these may be enabled on a cluster")
		}
	}
	return nil
}

func (a *Properties) validateExtensions() error {
	for _, agentPool := range a.AgentPoolProfiles {
		if len(agentPool.Extensions) != 0 && (len(agentPool.AvailabilityProfile) == 0 || agentPool.IsVirtualMachineScaleSets()) {
			return errors.Errorf("Extensions are currently not supported with VirtualMachineScaleSets. Please specify \"availabilityProfile\": \"%s\"", AvailabilitySet)
		}

		if agentPool.OSType == Windows && len(agentPool.Extensions) != 0 {
			for _, e := range agentPool.Extensions {
				if e.Name == "prometheus-grafana-k8s" {
					return errors.Errorf("prometheus-grafana-k8s extension is currently not supported for Windows agents")
				}
			}
		}
	}

	for _, extension := range a.ExtensionProfiles {
		if extension.ExtensionParametersKeyVaultRef != nil {
			if e := validate.Var(extension.ExtensionParametersKeyVaultRef.VaultID, "required"); e != nil {
				return errors.Errorf("the Keyvault ID must be specified for Extension %s", extension.Name)
			}
			if e := validate.Var(extension.ExtensionParametersKeyVaultRef.SecretName, "required"); e != nil {
				return errors.Errorf("the Keyvault Secret must be specified for Extension %s", extension.Name)
			}
			if !keyvaultIDRegex.MatchString(extension.ExtensionParametersKeyVaultRef.VaultID) {
				return errors.Errorf("Extension %s's keyvault secret reference is of incorrect format", extension.Name)
			}
		}
	}
	return nil
}

func (a *Properties) validateVNET() error {
	isCustomVNET := a.MasterProfile.IsCustomVNET()
	for _, agentPool := range a.AgentPoolProfiles {
		if agentPool.IsCustomVNET() != isCustomVNET {
			return errors.New("Multiple VNET Subnet configurations specified.  The master profile and each agent pool profile must all specify a custom VNET Subnet, or none at all")
		}
	}
	if isCustomVNET {
		if a.MasterProfile.IsVirtualMachineScaleSets() && a.MasterProfile.AgentVnetSubnetID == "" {
			return errors.New("when master profile is using VirtualMachineScaleSets and is custom vnet, set \"vnetsubnetid\" and \"agentVnetSubnetID\" for master profile")
		}

		subscription, resourcegroup, vnetname, _, e := common.GetVNETSubnetIDComponents(a.MasterProfile.VnetSubnetID)
		if e != nil {
			return e
		}

		for _, agentPool := range a.AgentPoolProfiles {
			agentSubID, agentRG, agentVNET, _, err := common.GetVNETSubnetIDComponents(agentPool.VnetSubnetID)
			if err != nil {
				return err
			}
			if agentSubID != subscription ||
				agentRG != resourcegroup ||
				agentVNET != vnetname {
				return errors.New("Multiple VNETS specified.  The master profile and each agent pool must reference the same VNET (but it is ok to reference different subnets on that VNET)")
			}
		}

		masterFirstIP := net.ParseIP(a.MasterProfile.FirstConsecutiveStaticIP)
		if masterFirstIP == nil && !a.MasterProfile.IsVirtualMachineScaleSets() {
			return errors.Errorf("MasterProfile.FirstConsecutiveStaticIP (with VNET Subnet specification) '%s' is an invalid IP address", a.MasterProfile.FirstConsecutiveStaticIP)
		}

		if a.MasterProfile.VnetCidr != "" {
			_, _, err := net.ParseCIDR(a.MasterProfile.VnetCidr)
			if err != nil {
				return errors.Errorf("MasterProfile.VnetCidr '%s' contains invalid cidr notation", a.MasterProfile.VnetCidr)
			}
		}
	}
	return nil
}

func (a *Properties) validateServicePrincipalProfile() error {
	useManagedIdentityDisabled := a.OrchestratorProfile.KubernetesConfig != nil &&
		a.OrchestratorProfile.KubernetesConfig.UseManagedIdentity != nil && !to.Bool(a.OrchestratorProfile.KubernetesConfig.UseManagedIdentity)

	if useManagedIdentityDisabled {
		if a.ServicePrincipalProfile == nil {
			return errors.Errorf("ServicePrincipalProfile must be specified")
		}
		if e := validate.Var(a.ServicePrincipalProfile.ClientID, "required"); e != nil {
			return errors.Errorf("the service principal client ID must be specified")
		}
		if (len(a.ServicePrincipalProfile.Secret) == 0 && a.ServicePrincipalProfile.KeyvaultSecretRef == nil) ||
			(len(a.ServicePrincipalProfile.Secret) != 0 && a.ServicePrincipalProfile.KeyvaultSecretRef != nil) {
			return errors.Errorf("either the service principal client secret or keyvault secret reference must be specified")
		}

		if a.OrchestratorProfile.KubernetesConfig != nil && to.Bool(a.OrchestratorProfile.KubernetesConfig.EnableEncryptionWithExternalKms) && len(a.ServicePrincipalProfile.ObjectID) == 0 {
			return errors.Errorf("the service principal object ID must be specified when enableEncryptionWithExternalKms is true")
		}

		if a.ServicePrincipalProfile.KeyvaultSecretRef != nil {
			if e := validate.Var(a.ServicePrincipalProfile.KeyvaultSecretRef.VaultID, "required"); e != nil {
				return errors.Errorf("the Keyvault ID must be specified for the Service Principle")
			}
			if e := validate.Var(a.ServicePrincipalProfile.KeyvaultSecretRef.SecretName, "required"); e != nil {
				return errors.Errorf("the Keyvault Secret must be specified for the Service Principle")
			}
			if !keyvaultIDRegex.MatchString(a.ServicePrincipalProfile.KeyvaultSecretRef.VaultID) {
				return errors.Errorf("service principal client keyvault secret reference is of incorrect format")
			}
		}
	}
	return nil
}

func (a *Properties) validateAADProfile() error {
	if profile := a.AADProfile; profile != nil {
		if _, err := uuid.Parse(profile.ClientAppID); err != nil {
			return errors.Errorf("clientAppID '%v' is invalid", profile.ClientAppID)
		}
		if _, err := uuid.Parse(profile.ServerAppID); err != nil {
			return errors.Errorf("serverAppID '%v' is invalid", profile.ServerAppID)
		}
		if len(profile.TenantID) > 0 {
			if _, err := uuid.Parse(profile.TenantID); err != nil {
				return errors.Errorf("tenantID '%v' is invalid", profile.TenantID)
			}
		}
		if len(profile.AdminGroupID) > 0 {
			if _, err := uuid.Parse(profile.AdminGroupID); err != nil {
				return errors.Errorf("adminGroupID '%v' is invalid", profile.AdminGroupID)
			}
		}
	}
	return nil
}

func (a *AgentPoolProfile) validateAvailabilityProfile() error {
	switch a.AvailabilityProfile {
	case AvailabilitySet:
	case VirtualMachineScaleSets:
	case "":
	default:
		{
			return errors.Errorf("unknown availability profile type '%s' for agent pool '%s'.  Specify either %s, or %s", a.AvailabilityProfile, a.Name, AvailabilitySet, VirtualMachineScaleSets)
		}
	}

	return nil
}

func (a *AgentPoolProfile) validateRoles() error {
	validRoles := []AgentPoolProfileRole{AgentPoolProfileRoleEmpty}
	var found bool
	for _, validRole := range validRoles {
		if a.Role == validRole {
			found = true
			break
		}
	}
	if !found {
		return errors.Errorf("Role %q is not supported", a.Role)
	}
	return nil
}

func (a *AgentPoolProfile) validateCustomNodeLabels() error {
	if len(a.CustomNodeLabels) > 0 {
		for k, v := range a.CustomNodeLabels {
			if e := validateKubernetesLabelKey(k); e != nil {
				return e
			}
			if e := validateKubernetesLabelValue(v); e != nil {
				return e
			}
		}
	}
	return nil
}

func validateVMSS(o *OrchestratorProfile, isUpdate bool, storageProfile string, hasWindows bool, isAzureStackCloud bool) error {
	version := common.RationalizeReleaseAndVersion(
		o.OrchestratorType,
		o.OrchestratorRelease,
		o.OrchestratorVersion,
		isUpdate,
		hasWindows,
		isAzureStackCloud)
	if version == "" {
		return errors.Errorf("the following OrchestratorProfile configuration is not supported: OrchestratorType: %s, OrchestratorRelease: %s, OrchestratorVersion: %s. Please check supported Release or Version for this build of aks-engine", o.OrchestratorType, o.OrchestratorRelease, o.OrchestratorVersion)
	}

	sv, err := semver.Make(version)
	if err != nil {
		return errors.Errorf("could not validate version %s", version)
	}
	minVersion, err := semver.Make("1.10.0")
	if err != nil {
		return errors.New("could not validate version")
	}
	if sv.LT(minVersion) {
		return errors.Errorf("VirtualMachineScaleSets are only available in Kubernetes version %s or greater. Please set \"orchestratorVersion\" to %s or above", minVersion.String(), minVersion.String())
	}
	// validation for instanceMetadata using VMSS with Kubernetes
	minVersion, err = semver.Make("1.10.2")
	if err != nil {
		return errors.New("could not validate version")
	}
	if o.KubernetesConfig != nil && o.KubernetesConfig.UseInstanceMetadata != nil {
		if *o.KubernetesConfig.UseInstanceMetadata && sv.LT(minVersion) {
			return errors.Errorf("VirtualMachineScaleSets with instance metadata is supported for Kubernetes version %s or greater. Please set \"useInstanceMetadata\": false in \"kubernetesConfig\" or set \"orchestratorVersion\" to %s or above", minVersion.String(), minVersion.String())
		}
	}
	if storageProfile == StorageAccount {
		return errors.Errorf("VirtualMachineScaleSets does not support %s disks.  Please specify \"storageProfile\": \"%s\" (recommended) or \"availabilityProfile\": \"%s\"", StorageAccount, ManagedDisks, AvailabilitySet)
	}
	return nil
}

func (a *Properties) validateWindowsProfile(isUpdate bool) error {
	hasWindowsAgentPools := false
	for _, profile := range a.AgentPoolProfiles {
		if profile.OSType == Windows {
			hasWindowsAgentPools = true
			break
		}
	}

	if !hasWindowsAgentPools {
		return nil
	}

	o := a.OrchestratorProfile
	version := common.RationalizeReleaseAndVersion(
		o.OrchestratorType,
		o.OrchestratorRelease,
		o.OrchestratorVersion,
		isUpdate,
		hasWindowsAgentPools,
		a.IsAzureStackCloud())

	if version == "" {
		return errors.Errorf("Orchestrator %s version %s does not support Windows", o.OrchestratorType, o.OrchestratorVersion)
	}

	w := a.WindowsProfile
	if w == nil {
		return errors.New("WindowsProfile is required when the cluster definition contains Windows agent pools")
	}
	if e := validate.Var(w.AdminUsername, "required"); e != nil {
		return errors.New("WindowsProfile.AdminUsername is required, when agent pool specifies Windows")
	}
	if e := validate.Var(w.AdminPassword, "required"); e != nil {
		return errors.New("WindowsProfile.AdminPassword is required, when agent pool specifies Windows")
	}
	if !validatePasswordComplexity(w.AdminUsername, w.AdminPassword) {
		return errors.New("WindowsProfile.AdminPassword complexity not met. Windows password should contain 3 of the following categories - uppercase letters(A-Z), lowercase(a-z) letters, digits(0-9), special characters (~!@#$%^&*_-+=`|\\(){}[]:;<>,.?/')")
	}
	if e := validateKeyVaultSecrets(w.Secrets, true); e != nil {
		return e
	}
	if e := validateCsiProxyWindowsProperties(w, version); e != nil {
		return e
	}
	if e := validateWindowsRuntimes(w.WindowsRuntimes); e != nil {
		return e
	}

	return nil
}

func validateCsiProxyWindowsProperties(w *WindowsProfile, k8sVersion string) error {
	if w.IsCSIProxyEnabled() && !common.IsKubernetesVersionGe(k8sVersion, "1.18.0") {
		return errors.New("CSI proxy for Windows is only available in Kubernetes versions 1.18.0 or greater")
	}
	return nil
}

func validateWindowsRuntimes(r *WindowsRuntimes) error {
	if r == nil {
		// can be blank defaults will be applied
		return nil
	}

	if r.Default != "process" && r.Default != "hyperv" {
		return errors.New("Default runtime types are process or hyperv")
	}

	if r.HypervRuntimes != nil {
		handlersMap := make(map[string]bool)
		for _, h := range r.HypervRuntimes {
			if h.BuildNumber != "17763" && h.BuildNumber != "18362" && h.BuildNumber != "18363" && h.BuildNumber != "19041" {
				return errors.New("Current hyper-v build id values supported are 17763, 18362, 18363, 19041")
			}

			if _, ok := handlersMap[h.BuildNumber]; ok {
				return errors.Errorf("Hyper-v RuntimeHandlers have duplicate runtime with build number '%s', Windows Runtimes must be unique", h.BuildNumber)
			}
			handlersMap[h.BuildNumber] = true
		}
	}

	return nil
}

func (a *AgentPoolProfile) validateOrchestratorSpecificProperties() error {

	if e := validate.Var(a.DNSPrefix, "len=0"); e != nil {
		return errors.New("AgentPoolProfile.DNSPrefix must be empty for Kubernetes")
	}
	if e := validate.Var(a.Ports, "len=0"); e != nil {
		return errors.New("AgentPoolProfile.Ports must be empty for Kubernetes")
	}
	if validate.Var(a.ScaleSetPriority, "eq=Regular") == nil && validate.Var(a.ScaleSetEvictionPolicy, "len=0") != nil {
		return errors.New("property 'AgentPoolProfile.ScaleSetEvictionPolicy' must be empty for AgentPoolProfile.Priority of Regular")
	}

	if a.DNSPrefix != "" {
		if e := common.ValidateDNSPrefix(a.DNSPrefix); e != nil {
			return e
		}
		if len(a.Ports) > 0 {
			if e := validateUniquePorts(a.Ports, a.Name); e != nil {
				return e
			}
		} else {
			a.Ports = []int{80, 443, 8080}
		}
	} else if e := validate.Var(a.Ports, "len=0"); e != nil {
		return errors.Errorf("AgentPoolProfile.Ports must be empty when AgentPoolProfile.DNSPrefix is empty")
	}

	if len(a.DiskSizesGB) > 0 {
		if e := validate.Var(a.StorageProfile, "eq=StorageAccount|eq=ManagedDisks"); e != nil {
			return errors.Errorf("property 'StorageProfile' must be set to either '%s' or '%s' when attaching disks", StorageAccount, ManagedDisks)
		}
		if e := validate.Var(a.AvailabilityProfile, "eq=VirtualMachineScaleSets|eq=AvailabilitySet"); e != nil {
			return errors.Errorf("property 'AvailabilityProfile' must be set to either '%s' or '%s' when attaching disks", VirtualMachineScaleSets, AvailabilitySet)
		}
		if a.StorageProfile == StorageAccount && (a.AvailabilityProfile != AvailabilitySet) {
			return errors.Errorf("VirtualMachineScaleSets does not support storage account attached disks.  Instead specify 'StorageAccount': '%s' or specify AvailabilityProfile '%s'", ManagedDisks, AvailabilitySet)
		}
	}

	if a.DiskEncryptionSetID != "" {
		if !diskEncryptionSetIDRegex.MatchString(a.DiskEncryptionSetID) {
			return errors.Errorf("DiskEncryptionSetID(%s) is of incorrect format, correct format: %s", a.DiskEncryptionSetID, diskEncryptionSetIDRegex.String())
		}
	}
	return nil
}

func (a *AgentPoolProfile) validateLoadBalancerBackendAddressPoolIDs() error {

	if a.LoadBalancerBackendAddressPoolIDs != nil {
		for _, backendPoolID := range a.LoadBalancerBackendAddressPoolIDs {
			if len(backendPoolID) == 0 {
				return errors.Errorf("AgentPoolProfile.LoadBalancerBackendAddressPoolIDs can not contain empty string. Agent pool name: %s", a.Name)
			}
		}
	}

	return nil
}

func validateProximityPlacementGroupID(ppgID string) error {
	if ppgID != "" {
		if !proximityPlacementGroupIDRegex.MatchString(ppgID) {
			return errors.Errorf("ProximityPlacementGroupID(%s) is of incorrect format, correct format: %s", ppgID, proximityPlacementGroupIDRegex.String())
		}
	}
	return nil
}

func validateKeyVaultSecrets(secrets []KeyVaultSecrets, requireCertificateStore bool) error {
	for _, s := range secrets {
		if len(s.VaultCertificates) == 0 {
			return errors.New("Valid KeyVaultSecrets must have no empty VaultCertificates")
		}
		if s.SourceVault == nil {
			return errors.New("missing SourceVault in KeyVaultSecrets")
		}
		if s.SourceVault.ID == "" {
			return errors.New("KeyVaultSecrets must have a SourceVault.ID")
		}
		for _, c := range s.VaultCertificates {
			if _, e := url.Parse(c.CertificateURL); e != nil {
				return errors.Errorf("Certificate url was invalid. received error %s", e)
			}
			if e := validateName(c.CertificateStore, "KeyVaultCertificate.CertificateStore"); requireCertificateStore && e != nil {
				return errors.Errorf("%s for certificates in a WindowsProfile", e)
			}
		}
	}
	return nil
}

func validatePasswordComplexity(name string, password string) (out bool) {

	if strings.EqualFold(name, password) {
		return false
	}

	if len(password) == 0 {
		return false
	}

	hits := 0
	if regexp.MustCompile(`[0-9]+`).MatchString(password) {
		hits++
	}
	if regexp.MustCompile(`[A-Z]+`).MatchString(password) {
		hits++
	}
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		hits++
	}
	if regexp.MustCompile(`[~!@#\$%\^&\*_\-\+=\x60\|\(\){}\[\]:;"'<>,\.\?/]+`).MatchString(password) {
		hits++
	}
	return hits > 2
}

// Validate validates the KubernetesConfig
func (k *KubernetesConfig) Validate(k8sVersion string, hasWindows, ipv6DualStackEnabled, isIPv6, isUpdate bool) error {
	// number of minimum retries allowed for kubelet to post node status
	const minKubeletRetries = 4

	// enableIPv6DualStack and enableIPv6Only are mutually exclusive feature flags
	if ipv6DualStackEnabled && isIPv6 {
		return errors.Errorf("featureFlags.EnableIPv6DualStack and featureFlags.EnableIPv6Only can't be enabled at the same time")
	}

	sv, err := semver.Make(k8sVersion)
	if err != nil {
		return errors.Errorf("could not validate version %s", k8sVersion)
	}

	if ipv6DualStackEnabled {
		minVersion, err := semver.Make("1.16.0")
		if err != nil {
			return errors.New("could not validate version")
		}
		if sv.LT(minVersion) {
			return errors.Errorf("IPv6 dual stack not available in kubernetes version %s", k8sVersion)
		}
		// ipv6 dual stack feature is currently only supported with kubenet
		if k.NetworkPlugin != "kubenet" && k.NetworkPlugin != "azure" {
			return errors.Errorf("the OrchestratorProfile.KubernetesConfig.NetworkPlugin '%s' is invalid, IPv6 dual stack supported only with 'kubenet' and 'azure'", k.NetworkPlugin)
		}

		if k.NetworkPlugin == "azure" && k.NetworkPolicy != "" {
			return errors.Errorf("Network policy %s is not supported for azure cni dualstack", k.NetworkPolicy)
		}
	}

	if isIPv6 {
		minVersion, err := semver.Make("1.18.0")
		if err != nil {
			return errors.New("could not validate version")
		}
		if sv.LT(minVersion) {
			return errors.Errorf("IPv6 single stack not available in kubernetes version %s", k8sVersion)
		}
		// single stack IPv6 feature is currently only supported with kubenet
		if k.NetworkPlugin != "kubenet" {
			return errors.Errorf("the OrchestratorProfile.KubernetesConfig.NetworkPlugin '%s' is invalid, IPv6 single stack supported only with kubenet", k.NetworkPlugin)
		}
	}

	if k.ClusterSubnet != "" {
		clusterSubnets := strings.Split(k.ClusterSubnet, ",")
		if !ipv6DualStackEnabled && len(clusterSubnets) > 1 {
			return errors.Errorf("OrchestratorProfile.KubernetesConfig.ClusterSubnet '%s' is an invalid subnet", k.ClusterSubnet)
		}
		if ipv6DualStackEnabled && len(clusterSubnets) > 2 {
			return errors.Errorf("the OrchestratorProfile.KubernetesConfig.ClusterSubnet '%s' is an invalid subnet, not more than 2 subnets for ipv6 dual stack", k.ClusterSubnet)
		}

		for _, clusterSubnet := range clusterSubnets {
			_, subnet, err := net.ParseCIDR(clusterSubnet)
			if err != nil {
				return errors.Errorf("OrchestratorProfile.KubernetesConfig.ClusterSubnet '%s' is an invalid subnet", clusterSubnet)
			}

			if k.NetworkPlugin == "azure" {
				ones, bits := subnet.Mask.Size()
				if bits-ones <= 8 {
					return errors.Errorf("OrchestratorProfile.KubernetesConfig.ClusterSubnet '%s' must reserve at least 9 bits for nodes", clusterSubnet)
				}
			}
		}
	}

	if k.DockerBridgeSubnet != "" {
		_, _, err := net.ParseCIDR(k.DockerBridgeSubnet)
		if err != nil {
			return errors.Errorf("OrchestratorProfile.KubernetesConfig.DockerBridgeSubnet '%s' is an invalid subnet", k.DockerBridgeSubnet)
		}
	}

	if k.MaxPods != 0 {
		if k.MaxPods < KubernetesMinMaxPods {
			return errors.Errorf("OrchestratorProfile.KubernetesConfig.MaxPods '%v' must be at least %v", k.MaxPods, KubernetesMinMaxPods)
		}
	}

	if k.KubeletConfig != nil {
		if _, ok := k.KubeletConfig["--node-status-update-frequency"]; ok {
			val := k.KubeletConfig["--node-status-update-frequency"]
			_, err := time.ParseDuration(val)
			if err != nil {
				return errors.Errorf("--node-status-update-frequency '%s' is not a valid duration", val)
			}
		}
	}

	if _, ok := k.ControllerManagerConfig["--node-monitor-grace-period"]; ok {
		_, err := time.ParseDuration(k.ControllerManagerConfig["--node-monitor-grace-period"])
		if err != nil {
			return errors.Errorf("--node-monitor-grace-period '%s' is not a valid duration", k.ControllerManagerConfig["--node-monitor-grace-period"])
		}
	}

	if k.KubeletConfig != nil {
		if _, ok := k.KubeletConfig["--node-status-update-frequency"]; ok {
			if _, ok := k.ControllerManagerConfig["--node-monitor-grace-period"]; ok {
				nodeStatusUpdateFrequency, _ := time.ParseDuration(k.KubeletConfig["--node-status-update-frequency"])
				ctrlMgrNodeMonitorGracePeriod, _ := time.ParseDuration(k.ControllerManagerConfig["--node-monitor-grace-period"])
				kubeletRetries := ctrlMgrNodeMonitorGracePeriod.Seconds() / nodeStatusUpdateFrequency.Seconds()
				if kubeletRetries < minKubeletRetries {
					return errors.Errorf("aks-engine-azurestack requires that --node-monitor-grace-period(%f)s be larger than nodeStatusUpdateFrequency(%f)s by at least a factor of %d; ", ctrlMgrNodeMonitorGracePeriod.Seconds(), nodeStatusUpdateFrequency.Seconds(), minKubeletRetries)
				}
			}
		}
		// Re-enable this unit test if --non-masquerade-cidr is re-introduced
		/*if _, ok := k.KubeletConfig["--non-masquerade-cidr"]; ok {
			if _, _, err := net.ParseCIDR(k.KubeletConfig["--non-masquerade-cidr"]); err != nil {
				return errors.Errorf("--non-masquerade-cidr kubelet config '%s' is an invalid CIDR string", k.KubeletConfig["--non-masquerade-cidr"])
			}
		}*/
	}

	if _, ok := k.ControllerManagerConfig["--pod-eviction-timeout"]; ok {
		_, err := time.ParseDuration(k.ControllerManagerConfig["--pod-eviction-timeout"])
		if err != nil {
			return errors.Errorf("--pod-eviction-timeout '%s' is not a valid duration", k.ControllerManagerConfig["--pod-eviction-timeout"])
		}
	}

	if _, ok := k.ControllerManagerConfig["--route-reconciliation-period"]; ok {
		_, err := time.ParseDuration(k.ControllerManagerConfig["--route-reconciliation-period"])
		if err != nil {
			return errors.Errorf("--route-reconciliation-period '%s' is not a valid duration", k.ControllerManagerConfig["--route-reconciliation-period"])
		}
	}

	if k.DNSServiceIP != "" || k.ServiceCidr != "" {
		if k.DNSServiceIP == "" {
			return errors.New("OrchestratorProfile.KubernetesConfig.DNSServiceIP must be specified when ServiceCidr is")
		}
		if k.ServiceCidr == "" {
			return errors.New("OrchestratorProfile.KubernetesConfig.ServiceCidr must be specified when DNSServiceIP is")
		}

		dnsIP := net.ParseIP(k.DNSServiceIP)
		if dnsIP == nil {
			return errors.Errorf("OrchestratorProfile.KubernetesConfig.DNSServiceIP '%s' is an invalid IP address", k.DNSServiceIP)
		}

		primaryServiceCIDR := k.ServiceCidr
		if ipv6DualStackEnabled {
			// split the service cidr to see if there are multiple cidrs
			serviceCidrs := strings.Split(k.ServiceCidr, ",")
			if len(serviceCidrs) > 2 {
				return errors.Errorf("OrchestratorProfile.KubernetesConfig.ServiceCidr '%s' is an invalid CIDR subnet. More than 2 CIDRs not allowed for dualstack", k.ServiceCidr)
			}
			if len(serviceCidrs) == 2 {
				firstServiceCIDR, secondServiceCIDR := serviceCidrs[0], serviceCidrs[1]
				_, _, err := net.ParseCIDR(secondServiceCIDR)
				if err != nil {
					return errors.Errorf("OrchestratorProfile.KubernetesConfig.ServiceCidr '%s' is an invalid CIDR subnet", secondServiceCIDR)
				}
				// use the primary service cidr for further validation
				primaryServiceCIDR = firstServiceCIDR
			}
			// if # of service cidrs is 1, then continues with the default validation
		}

		_, serviceCidr, err := net.ParseCIDR(primaryServiceCIDR)
		if err != nil {
			return errors.Errorf("OrchestratorProfile.KubernetesConfig.ServiceCidr '%s' is an invalid CIDR subnet", primaryServiceCIDR)
		}

		// Finally validate that the DNS ip is within the subnet
		if !serviceCidr.Contains(dnsIP) {
			return errors.Errorf("OrchestratorProfile.KubernetesConfig.DNSServiceIP '%s' is not within the ServiceCidr '%s'", k.DNSServiceIP, primaryServiceCIDR)
		}

		// and that the DNS IP is _not_ the subnet broadcast address
		broadcast := common.IP4BroadcastAddress(serviceCidr)
		if dnsIP.Equal(broadcast) {
			return errors.Errorf("OrchestratorProfile.KubernetesConfig.DNSServiceIP '%s' cannot be the broadcast address of ServiceCidr '%s'", k.DNSServiceIP, primaryServiceCIDR)
		}

		// and that the DNS IP is _not_ the first IP in the service subnet
		firstServiceIP := common.CidrFirstIP(serviceCidr.IP)
		if firstServiceIP.Equal(dnsIP) {
			return errors.Errorf("OrchestratorProfile.KubernetesConfig.DNSServiceIP '%s' cannot be the first IP of ServiceCidr '%s'", k.DNSServiceIP, primaryServiceCIDR)
		}
	}

	if k.ProxyMode != "" && k.ProxyMode != KubeProxyModeIPTables && k.ProxyMode != KubeProxyModeIPVS {
		return errors.Errorf("Invalid KubeProxyMode %v. Allowed modes are %v and %v", k.ProxyMode, KubeProxyModeIPTables, KubeProxyModeIPVS)
	}

	// dualstack IPVS mode supported from 1.16+
	// dualstack IPtables mode supported from 1.18+
	if ipv6DualStackEnabled && k.ProxyMode == KubeProxyModeIPTables {
		minVersion, err := semver.Make("1.18.0")
		if err != nil {
			return errors.New("could not validate version")
		}
		if sv.LT(minVersion) {
			return errors.Errorf("KubeProxyMode %v in dualstack not supported with %s version", k.ProxyMode, k8sVersion)
		}
	}

	// Validate that we have a valid etcd version
	if e := validateEtcdVersion(k.EtcdVersion); e != nil {
		return e
	}

	// Validate containerd scenarios
	if k.ContainerRuntime == Docker || k.ContainerRuntime == "" {
		if k.MobyVersion != "" && k.ContainerdVersion != "" && versions.LessThan(k.MobyVersion, "19.03") {
			return errors.Errorf("containerdVersion is only valid in a non-docker context, use %s containerRuntime value instead if you wish to provide a containerdVersion", Containerd)
		}
	}
	if e := validateContainerdVersion(k.ContainerdVersion); e != nil {
		return e
	}

	if to.Bool(k.UseCloudControllerManager) || k.CustomCcmImage != "" {
		sv, err := semver.Make(k8sVersion)
		if err != nil {
			return errors.Errorf("could not validate version %s", k8sVersion)
		}
		minVersion, err := semver.Make("1.8.0")
		if err != nil {
			return errors.New("could not validate version")
		}
		if sv.LT(minVersion) {
			return errors.Errorf("OrchestratorProfile.KubernetesConfig.UseCloudControllerManager and OrchestratorProfile.KubernetesConfig.CustomCcmImage not available in kubernetes version %s", k8sVersion)
		}
	}

	if e := k.validateNetworkPlugin(hasWindows, isUpdate); e != nil {
		return e
	}
	if e := k.validateNetworkPolicy(k8sVersion, hasWindows); e != nil {
		return e
	}
	if e := k.validateNetworkPluginPlusPolicy(); e != nil {
		return e
	}
	if e := k.validateNetworkMode(); e != nil {
		return e
	}
	if e := k.validateKubernetesImageBaseType(); e != nil {
		return e
	}

	if to.Bool(k.EnableMultipleStandardLoadBalancers) && !common.IsKubernetesVersionGe(k8sVersion, "1.20.0-beta.1") {
		return errors.Errorf("OrchestratorProfile.KubernetesConfig.EnableMultipleStandardLoadBalancers is available since kubernetes version v1.20.0-beta.1, current version is %s", k8sVersion)
	}
	if k.Tags != "" && !common.IsKubernetesVersionGe(k8sVersion, "1.20.0-beta.1") {
		return errors.Errorf("OrchestratorProfile.KubernetesConfig.Tags is available since kubernetes version v1.20.0-beta.1, current version is %s", k8sVersion)
	}
	return k.validateContainerRuntimeConfig()
}

func (k *KubernetesConfig) validateContainerRuntimeConfig() error {
	if val, ok := k.ContainerRuntimeConfig[common.ContainerDataDirKey]; ok {
		if val == "" {
			return errors.Errorf("OrchestratorProfile.KubernetesConfig.ContainerRuntimeConfig.DataDir '%s' is invalid: must not be empty", val)
		}
		if !strings.HasPrefix(val, "/") {
			return errors.Errorf("OrchestratorProfile.KubernetesConfig.ContainerRuntimeConfig.DataDir '%s' is invalid: must be absolute path", val)
		}
	}

	// Validate base config here, and only allow predefined mutations to ensure invariant.
	if k.ContainerRuntime == Containerd {
		_, err := common.GetContainerdConfig(k.ContainerRuntimeConfig, nil)
		if err != nil {
			return err
		}
	} else {
		_, err := common.GetDockerConfig(k.ContainerRuntimeConfig, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k *KubernetesConfig) validateNetworkPlugin(hasWindows, isUpdate bool) error {

	networkPlugin := k.NetworkPlugin

	// Check NetworkPlugin has a valid value.
	valid := false
	for _, plugin := range NetworkPluginValues {
		if networkPlugin == plugin {
			if plugin == NetworkPluginFlannel {
				if isUpdate {
					valid = true
				}
			} else {
				valid = true
			}
			break
		}
	}
	if !valid {
		if networkPlugin == NetworkPluginFlannel {
			return errors.Errorf("networkPlugin '%s' has been deprecated and is no longer supported for new cluster creation", networkPlugin)
		}
		return errors.Errorf("unknown networkPlugin '%s' specified", networkPlugin)
	}

	// Temporary safety check, to be removed when Windows support is added.
	if (networkPlugin == NetworkPluginAntrea) && hasWindows {
		return errors.Errorf("networkPlugin '%s' is not supporting windows agents", networkPlugin)
	}

	if networkPlugin == NetworkPluginKubenet && hasWindows {
		log.Warnf("Windows + Kubenet is for development and testing only, not recommended for production")
	}

	return nil
}

func (k *KubernetesConfig) validateNetworkPolicy(k8sVersion string, hasWindows bool) error {

	networkPolicy := k.NetworkPolicy
	networkPlugin := k.NetworkPlugin

	// Check NetworkPolicy has a valid value.
	valid := false
	for _, plugin := range NetworkPolicyValues {
		if networkPolicy == plugin {
			valid = true
			break
		}
	}
	if !valid {
		return errors.Errorf("unknown networkPolicy '%s' specified", networkPolicy)
	}

	if networkPolicy == "azure" && networkPlugin == "azure" && !common.IsKubernetesVersionGe(k8sVersion, "1.8.0") {
		return errors.New("networkPolicy azure requires kubernetes version of 1.8 or higher")
	}

	// Temporary safety check, to be removed when Windows support is added.
	if (networkPolicy == "calico" || networkPolicy == NetworkPolicyCilium ||
		networkPolicy == NetworkPolicyAntrea) && hasWindows {
		return errors.Errorf("networkPolicy '%s' is not supporting windows agents", networkPolicy)
	}

	return nil
}

func (k *KubernetesConfig) validateNetworkPluginPlusPolicy() error {
	var config k8sNetworkConfig

	config.networkPlugin = k.NetworkPlugin
	config.networkPolicy = k.NetworkPolicy

	for _, c := range networkPluginPlusPolicyAllowed {
		if c.networkPlugin == config.networkPlugin && c.networkPolicy == config.networkPolicy {
			return nil
		}
	}
	return errors.Errorf("networkPolicy '%s' is not supported with networkPlugin '%s'", config.networkPolicy, config.networkPlugin)
}

func (k *KubernetesConfig) validateNetworkMode() error {
	networkPlugin := k.NetworkPlugin
	networkPolicy := k.NetworkPolicy
	networkMode := k.NetworkMode

	// Check NetworkMode has a valid value.
	valid := false
	for _, mode := range NetworkModeValues {
		if networkMode == mode {
			valid = true
			break
		}
	}
	if !valid {
		return errors.Errorf("unknown networkMode '%s' specified", networkMode)
	}

	if networkMode != "" {
		if networkPlugin != "azure" {
			return errors.New("networkMode requires network plugin to be 'azure'")
		}

		if networkPolicy == "calico" && networkMode != NetworkModeTransparent {
			return errors.Errorf("networkMode '%s' is not supported by calico", networkMode)
		}
	}

	return nil
}

func (k *KubernetesConfig) validateKubernetesImageBaseType() error {
	for _, valid := range kubernetesImageBaseTypeValidVersions {
		if valid == k.KubernetesImageBaseType {
			return nil
		}
	}
	return errors.Errorf("Invalid kubernetesImageBaseType value \"%s\", please use one of the following versions: %s", k.KubernetesImageBaseType, kubernetesImageBaseTypeValidVersions)
}

func (k *KubernetesConfig) isUsingCustomKubeComponent() bool {
	return k.CustomKubeAPIServerImage != "" || k.CustomKubeControllerManagerImage != "" || k.CustomKubeSchedulerImage != "" || k.CustomKubeBinaryURL != ""
}

func (a *Properties) validateContainerRuntime(isUpdate bool) error {
	var containerRuntime string

	if a.OrchestratorProfile.KubernetesConfig != nil {
		containerRuntime = a.OrchestratorProfile.KubernetesConfig.ContainerRuntime
	}

	// Check for deprecated, non-back-compat
	if isUpdate && containerRuntime == KataContainers {
		return errors.Errorf("%s containerRuntime has been deprecated, you will not be able to update this cluster with this version of aks-engine", KataContainers)
	}

	// Check ContainerRuntime has a valid value.
	valid := false
	for _, runtime := range ContainerRuntimeValues {
		if containerRuntime == runtime {
			valid = true
			break
		}
	}
	if !valid {
		return errors.Errorf("unknown containerRuntime %q specified", containerRuntime)
	}

	// TODO: These validations should be relaxed once ContainerD and CNI plugins are more readily available
	if containerRuntime == Containerd && a.HasWindows() {
		if a.OrchestratorProfile.KubernetesConfig.NetworkPlugin == "kubenet" {
			if a.OrchestratorProfile.KubernetesConfig.WindowsSdnPluginURL == "" {
				return errors.Errorf("WindowsSdnPluginURL must be provided when using Windows with ContainerRuntime=containerd and networkPlugin=kubenet")
			}
		}
	}

	return nil
}

func (a *Properties) validateCustomKubeComponent() error {
	k := a.OrchestratorProfile.KubernetesConfig
	if k == nil {
		return nil
	}

	if common.IsKubernetesVersionGe(a.OrchestratorProfile.OrchestratorVersion, "1.17.0") {
		if k.CustomHyperkubeImage != "" {
			return errors.New("customHyperkubeImage has no effect in Kubernetes version 1.17.0 or above")
		}
	} else {
		if k.isUsingCustomKubeComponent() {
			return errors.New("customKubeAPIServerImage, customKubeControllerManagerImage, customKubeSchedulerImage or customKubeBinaryURL have no effect in Kubernetes version 1.16 or earlier")
		}
	}
	if !common.IsKubernetesVersionGe(a.OrchestratorProfile.OrchestratorVersion, "1.16.0") {
		if k.CustomKubeProxyImage != "" {
			return errors.New("customKubeProxyImage has no effect in Kubernetes version 1.15 or earlier")
		}
	}

	return nil
}

func validateName(name string, label string) error {
	if name == "" {
		return errors.Errorf("%s must be a non-empty value", label)
	}
	return nil
}

func validatePoolName(poolName string) error {
	// we will cap at length of 12 and all lowercase letters since this makes up the VMName
	poolNameRegex := `^([a-z][a-z0-9]{0,11})$`
	re, err := regexp.Compile(poolNameRegex)
	if err != nil {
		return err
	}
	submatches := re.FindStringSubmatch(poolName)
	if len(submatches) != 2 {
		return errors.Errorf("pool name '%s' is invalid. A pool name must start with a lowercase letter, have max length of 12, and only have characters a-z0-9", poolName)
	}
	return nil
}

func validatePoolOSType(os OSType) error {
	if os != Linux && os != Windows && os != "" {
		return errors.New("AgentPoolProfile.osType must be either Linux or Windows")
	}
	return nil
}

func validatePoolAcceleratedNetworking(vmSize string) error {
	if !helpers.AcceleratedNetworkingSupported(vmSize) {
		return errors.Errorf("AgentPoolProfile.vmsize %s does not support AgentPoolProfile.acceleratedNetworking", vmSize)
	}
	return nil
}

func validateUniquePorts(ports []int, name string) error {
	portMap := make(map[int]bool)
	for _, port := range ports {
		if _, ok := portMap[port]; ok {
			return errors.Errorf("agent profile '%s' has duplicate port '%d', ports must be unique", name, port)
		}
		portMap[port] = true
	}
	return nil
}

func validateKubernetesLabelValue(v string) error {
	if !(len(v) == 0) && !labelValueRegex.MatchString(v) {
		return errors.Errorf("Label value '%s' is invalid. Valid label values must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between", v)
	}
	return nil
}

func validateKubernetesLabelKey(k string) error {
	if !labelKeyRegex.MatchString(k) {
		return errors.Errorf("Label key '%s' is invalid. Valid label keys have two segments: an optional prefix and name, separated by a slash (/). The name segment is required and must be 63 characters or less, beginning and ending with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between. The prefix is optional. If specified, the prefix must be a DNS subdomain: a series of DNS labels separated by dots (.), not longer than 253 characters in total, followed by a slash (/)", k)
	}
	prefix := strings.Split(k, "/")
	if len(prefix) != 1 && len(prefix[0]) > labelKeyPrefixMaxLength {
		return errors.Errorf("Label key prefix '%s' is invalid. If specified, the prefix must be no longer than 253 characters in total", k)
	}
	return nil
}

func validateEtcdVersion(etcdVersion string) error {
	// "" is a valid etcdVersion that maps to DefaultEtcdVersion
	if etcdVersion == "" {
		return nil
	}
	for _, ver := range etcdValidVersions {
		if ver == etcdVersion {
			return nil
		}
	}
	return errors.Errorf("Invalid etcd version \"%s\", please use one of the following versions: %s", etcdVersion, etcdValidVersions)
}

func validateContainerdVersion(containerdVersion string) error {
	// "" is a valid containerd that maps to DefaultContainerdVersion
	if containerdVersion == "" {
		return nil
	}
	for _, ver := range containerdValidVersions {
		if ver == containerdVersion {
			return nil
		}
	}
	return errors.Errorf("Invalid containerd version \"%s\", please use one of the following versions: %s", containerdVersion, containerdValidVersions)
}

// Check that distro has a valid value
func validateDistro(distro Distro, distroValues []Distro) bool {
	var ret bool
	for _, d := range distroValues {
		if distro == d {
			ret = true
		}
	}
	switch distro {
	case AKSUbuntu1604, Ubuntu:
		log.Warnf("The '%s' distro uses Ubuntu 16.04-LTS, which is End of Life (EOL) and will no longer receive security updates", distro)
	}
	return ret
}

func (i *ImageReference) validateImageNameAndGroup() error {
	if i.Name == "" && i.ResourceGroup != "" {
		return errors.New("imageName needs to be specified when imageResourceGroup is provided")
	}
	if i.Name != "" && i.ResourceGroup == "" {
		return errors.New("imageResourceGroup needs to be specified when imageName is provided")
	}
	return nil
}

func (cs *ContainerService) validateCustomCloudProfile() error {
	a := cs.Properties

	if a.IsCustomCloudProfile() {
		if a.IsAzureStackCloud() {
			if a.CustomCloudProfile.PortalURL == "" {
				return errors.New("portalURL needs to be specified when AzureStackCloud CustomCloudProfile is provided")
			}

			if !strings.HasPrefix(a.CustomCloudProfile.PortalURL, fmt.Sprintf("https://portal.%s.", cs.Location)) {
				return errors.Errorf("portalURL needs to start with https://portal.%s. ", cs.Location)
			}

			if a.CustomCloudProfile.AuthenticationMethod != "" && !(a.CustomCloudProfile.AuthenticationMethod == ClientSecretAuthMethod || a.CustomCloudProfile.AuthenticationMethod == ClientCertificateAuthMethod) {
				return errors.Errorf("authenticationMethod allowed values are '%s' and '%s'", ClientCertificateAuthMethod, ClientSecretAuthMethod)
			}

			if a.CustomCloudProfile.IdentitySystem != "" && !(a.CustomCloudProfile.IdentitySystem == AzureADIdentitySystem || a.CustomCloudProfile.IdentitySystem == ADFSIdentitySystem) {
				return errors.Errorf("identitySystem allowed values are '%s' and '%s'", AzureADIdentitySystem, ADFSIdentitySystem)
			}
		}

		dependenciesLocationValues := DependenciesLocationValues
		if !validateDependenciesLocation(a.CustomCloudProfile.DependenciesLocation, dependenciesLocationValues) {
			return errors.Errorf("The %s dependenciesLocation is not supported. The supported vaules are %s", a.CustomCloudProfile.DependenciesLocation, dependenciesLocationValues)
		}
	}
	return nil
}

// Validate implements validation for ContainerService
func (cs *ContainerService) Validate(isUpdate bool) error {
	if e := cs.validateProperties(); e != nil {
		return e
	}
	if e := cs.validateLocation(); e != nil {
		return e
	}
	if e := cs.validateCustomCloudProfile(); e != nil {
		return e
	}
	if e := cs.Properties.validate(isUpdate); e != nil {
		return e
	}
	return nil
}

func (cs *ContainerService) validateLocation() error {
	if cs.Properties != nil && cs.Properties.IsCustomCloudProfile() && cs.Location == "" {
		return errors.New("missing ContainerService Location")
	}
	if cs.Location == "" {
		log.Warnf("No \"location\" value was specified, AKS Engine will generate an ARM template configuration valid for regions in public cloud only")
	}
	return nil
}

func (cs *ContainerService) validateProperties() error {
	if cs.Properties == nil {
		return errors.New("missing ContainerService Properties")
	}
	return nil
}

// Check that dependenciesLocation has a valid value
func validateDependenciesLocation(dependenciesLocation DependenciesLocation, dependenciesLocationValues []DependenciesLocation) bool {
	for _, d := range dependenciesLocationValues {
		if dependenciesLocation == d {
			return true
		}
	}
	return false
}

// validateAzureStackSupport logs a warning if apimodel contains preview features and returns an error if a property is not supported on Azure Stack clouds
func (a *Properties) validateAzureStackSupport() error {
	if a.IsAzureStackCloud() {
		networkPlugin := a.OrchestratorProfile.KubernetesConfig.NetworkPlugin
		if networkPlugin != "azure" && networkPlugin != "kubenet" && networkPlugin != "" {
			return errors.Errorf("kubernetesConfig.networkPlugin '%s' is not supported on Azure Stack clouds", networkPlugin)
		}
		if a.MasterProfile.AvailabilityProfile == VirtualMachineScaleSets {
			return errors.Errorf("masterProfile.availabilityProfile should be set to '%s' on Azure Stack clouds", AvailabilitySet)
		}
		for _, agentPool := range a.AgentPoolProfiles {
			pool := agentPool
			if pool.AvailabilityProfile != AvailabilitySet {
				return errors.Errorf("agentPoolProfiles[%s].availabilityProfile should be set to '%s' on Azure Stack clouds", pool.Name, AvailabilitySet)
			}
		}
	}
	return nil
}

func (a *KubernetesAddon) validateArcAddonConfig() error {
	if a.Config == nil {
		a.Config = make(map[string]string)
	}
	err := []string{}
	if a.Config["location"] == "" {
		err = append(err, "azure-arc-onboarding addon configuration must have a 'location' property")
	}
	if a.Config["tenantID"] == "" {
		err = append(err, "azure-arc-onboarding addon configuration must have a 'tenantID' property")
	}
	if a.Config["subscriptionID"] == "" {
		err = append(err, "azure-arc-onboarding addon configuration must have a 'subscriptionID' property")
	}
	if a.Config["resourceGroup"] == "" {
		err = append(err, "azure-arc-onboarding addon configuration must have a 'resourceGroup' property")
	}
	if a.Config["clusterName"] == "" {
		err = append(err, "azure-arc-onboarding addon configuration must have a 'clusterName' property")
	}
	if a.Config["clientID"] == "" {
		err = append(err, "azure-arc-onboarding addon configuration must have a 'clientID' property")
	}
	if a.Config["clientSecret"] == "" {
		err = append(err, "azure-arc-onboarding addon configuration must have a 'clientSecret' property")
	}
	if len(err) > 0 {
		return errors.New(strings.Join(err, "; "))
	}
	return nil
}
