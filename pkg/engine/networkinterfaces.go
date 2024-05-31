// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"fmt"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
)

// CreateMasterVMNetworkInterfaces returns an ARM resource for the master VM NIC
func CreateMasterVMNetworkInterfaces(cs *api.ContainerService) NetworkInterfaceARM {
	var dependencies []string
	if cs.Properties.MasterProfile != nil && cs.Properties.MasterProfile.IsCustomVNET() {
		dependencies = append(dependencies, "[variables('nsgID')]")
	} else {
		dependencies = append(dependencies, "[variables('vnetID')]")
	}

	if cs.Properties.MasterProfile != nil && cs.Properties.MasterProfile.HasMultipleNodes() {
		dependencies = append(dependencies, "[variables('masterInternalLbName')]")
	}

	if cs.Properties.MasterProfile != nil && cs.Properties.MasterProfile.HasCosmosEtcd() {
		dependencies = append(dependencies, "[resourceId('Microsoft.DocumentDB/databaseAccounts/', variables('cosmosAccountName'))]")
	}

	lbBackendAddressPools := []network.BackendAddressPool{}
	dependencies = append(dependencies, "[variables('masterLbName')]")
	publicLbPool := network.BackendAddressPool{
		ID: helpers.PointerToString("[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
	}
	lbBackendAddressPools = append(lbBackendAddressPools, publicLbPool)

	armResource := ARMResource{
		APIVersion: "[variables('apiVersionNetwork')]",
		Copy: map[string]string{
			"count": "[sub(variables('masterCount'), variables('masterOffset'))]",
			"name":  "nicLoopNode",
		},
		DependsOn: dependencies,
	}

	if cs.Properties.MasterProfile != nil && cs.Properties.MasterProfile.HasMultipleNodes() {
		internalLbPool := network.BackendAddressPool{
			ID: helpers.PointerToString("[concat(variables('masterInternalLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
		}
		lbBackendAddressPools = append(lbBackendAddressPools, internalLbPool)
	}

	loadBalancerIPConfig := network.InterfaceIPConfiguration{
		Name: helpers.PointerToString("ipconfig1"),
		InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
			LoadBalancerBackendAddressPools: &lbBackendAddressPools,
			PrivateIPAddress:                helpers.PointerToString("[variables('masterPrivateIpAddrs')[copyIndex(variables('masterOffset'))]]"),
			Primary:                         helpers.PointerToBool(true),
			PrivateIPAllocationMethod:       network.Static,
			Subnet: &network.Subnet{
				ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
			},
		},
	}

	publicNatRules := []network.InboundNatRule{
		{
			ID: helpers.PointerToString("[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'),copyIndex(variables('masterOffset')))]"),
		},
	}
	loadBalancerIPConfig.LoadBalancerInboundNatRules = &publicNatRules

	isAzureCNI := cs.Properties.OrchestratorProfile.IsAzureCNI()

	ipConfigurations := []network.InterfaceIPConfiguration{loadBalancerIPConfig}

	nicProperties := network.InterfacePropertiesFormat{
		IPConfigurations: &ipConfigurations,
	}

	if isAzureCNI {
		ipConfigurations = append(ipConfigurations, getSecondaryNICIPConfigs(cs.Properties.MasterProfile.IPAddressCount)...)
		if cs.Properties.FeatureFlags.IsFeatureEnabled("EnableIPv6DualStack") {
			nicProperties.EnableIPForwarding = helpers.PointerToBool(true)
		}
	} else {
		if !cs.Properties.IsAzureStackCloud() {
			nicProperties.EnableIPForwarding = helpers.PointerToBool(true)
		}
	}

	// add ipv6 nic config for dual stack
	if cs.Properties.FeatureFlags.IsFeatureEnabled("EnableIPv6DualStack") || cs.Properties.FeatureFlags.IsFeatureEnabled("EnableIPv6Only") {
		ipv6Config := network.InterfaceIPConfiguration{
			Name: helpers.PointerToString("ipconfigv6"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				PrivateIPAddressVersion: "IPv6",
				Primary:                 helpers.PointerToBool(false),
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		}

		ipConfigurations = append(ipConfigurations, ipv6Config)
	}

	linuxProfile := cs.Properties.LinuxProfile
	if linuxProfile != nil && linuxProfile.HasCustomNodesDNS() {
		nicProperties.DNSSettings = &network.InterfaceDNSSettings{
			DNSServers: &[]string{
				"[parameters('dnsServer')]",
			},
		}
	}

	if cs.Properties.MasterProfile != nil && cs.Properties.MasterProfile.IsCustomVNET() {
		nicProperties.NetworkSecurityGroup = &network.SecurityGroup{
			ID: helpers.PointerToString("[variables('nsgID')]"),
		}
	}

	networkInterface := network.Interface{
		Location:                  helpers.PointerToString("[variables('location')]"),
		Name:                      helpers.PointerToString("[concat(variables('masterVMNamePrefix'), 'nic-', copyIndex(variables('masterOffset')))]"),
		InterfacePropertiesFormat: &nicProperties,
		Type:                      helpers.PointerToString("Microsoft.Network/networkInterfaces"),
	}

	return NetworkInterfaceARM{
		ARMResource: armResource,
		Interface:   networkInterface,
	}
}

// createPrivateClusterMasterVMNetworkInterface returns an ARM resource for the master VM NIC in a private cluster configuration scenario
func createPrivateClusterMasterVMNetworkInterface(cs *api.ContainerService) NetworkInterfaceARM {
	var dependencies []string
	if cs.Properties.MasterProfile.IsCustomVNET() {
		dependencies = append(dependencies, "[variables('nsgID')]")
	} else {
		dependencies = append(dependencies, "[variables('vnetID')]")
	}

	loadBalancerIPConfig := network.InterfaceIPConfiguration{
		Name: helpers.PointerToString("ipconfig1"),
		InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
			PrivateIPAddress:          helpers.PointerToString("[variables('masterPrivateIpAddrs')[copyIndex(variables('masterOffset'))]]"),
			Primary:                   helpers.PointerToBool(true),
			PrivateIPAllocationMethod: network.Static,
			Subnet: &network.Subnet{
				ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
			},
		},
	}

	if cs.Properties.MasterProfile.HasMultipleNodes() {
		dependencies = append(dependencies, "[variables('masterInternalLbName')]")
		var lbBackendAddressPools []network.BackendAddressPool
		internalLbPool := network.BackendAddressPool{
			ID: helpers.PointerToString("[concat(variables('masterInternalLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
		}
		lbBackendAddressPools = append(lbBackendAddressPools, internalLbPool)
		if cs.Properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku == api.StandardLoadBalancerSku {
			dependencies = append(dependencies, "[variables('masterLbName')]")
			publicLbPool := network.BackendAddressPool{
				ID: helpers.PointerToString("[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
			}
			lbBackendAddressPools = append(lbBackendAddressPools, publicLbPool)
		}
		loadBalancerIPConfig.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools = &lbBackendAddressPools
		loadBalancerIPConfig.InterfaceIPConfigurationPropertiesFormat.LoadBalancerInboundNatRules = &[]network.InboundNatRule{}
	}
	ipConfigurations := []network.InterfaceIPConfiguration{loadBalancerIPConfig}

	isAzureCNI := cs.Properties.OrchestratorProfile.IsAzureCNI()

	if isAzureCNI {
		for i := 2; i <= cs.Properties.MasterProfile.IPAddressCount; i++ {
			ipConfig := network.InterfaceIPConfiguration{
				Name: helpers.PointerToString(fmt.Sprintf("ipconfig%d", i)),
				InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
					Primary:                   helpers.PointerToBool(false),
					PrivateIPAllocationMethod: network.Dynamic,
					Subnet: &network.Subnet{
						ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
					},
				},
			}
			ipConfigurations = append(ipConfigurations, ipConfig)
		}
	}

	nicProperties := network.InterfacePropertiesFormat{
		IPConfigurations: &ipConfigurations,
	}

	if !isAzureCNI && !cs.Properties.IsAzureStackCloud() {
		nicProperties.EnableIPForwarding = helpers.PointerToBool(true)
	}

	// Enable IPForwarding on NetworkInterface for azurecni dualstack
	if isAzureCNI {
		if cs.Properties.FeatureFlags.IsFeatureEnabled("EnableIPv6DualStack") {
			nicProperties.EnableIPForwarding = helpers.PointerToBool(true)
		}
	}

	linuxProfile := cs.Properties.LinuxProfile
	if linuxProfile != nil && linuxProfile.HasCustomNodesDNS() {
		nicProperties.DNSSettings = &network.InterfaceDNSSettings{
			DNSServers: &[]string{
				"[parameters('dnsServer')]",
			},
		}
	}

	if cs.Properties.MasterProfile.IsCustomVNET() {
		nicProperties.NetworkSecurityGroup = &network.SecurityGroup{
			ID: helpers.PointerToString("[variables('nsgID')]"),
		}
	}

	networkInterface := network.Interface{
		Location:                  helpers.PointerToString("[variables('location')]"),
		Name:                      helpers.PointerToString("[concat(variables('masterVMNamePrefix'), 'nic-', copyIndex(variables('masterOffset')))]"),
		InterfacePropertiesFormat: &nicProperties,
		Type:                      helpers.PointerToString("Microsoft.Network/networkInterfaces"),
	}

	armResource := ARMResource{
		APIVersion: "[variables('apiVersionNetwork')]",
		Copy: map[string]string{
			"count": "[sub(variables('masterCount'), variables('masterOffset'))]",
			"name":  "nicLoopNode",
		},
		DependsOn: dependencies,
	}

	return NetworkInterfaceARM{
		ARMResource: armResource,
		Interface:   networkInterface,
	}
}

func createJumpboxNetworkInterface(cs *api.ContainerService) NetworkInterfaceARM {
	dependencies := []string{
		"[concat('Microsoft.Network/publicIpAddresses/', variables('jumpboxPublicIpAddressName'))]",
		"[concat('Microsoft.Network/networkSecurityGroups/', variables('jumpboxNetworkSecurityGroupName'))]",
	}

	if !cs.Properties.MasterProfile.IsCustomVNET() {
		dependencies = append(dependencies, "[variables('vnetID')]")
	}

	armResource := ARMResource{
		APIVersion: "[variables('apiVersionNetwork')]",
		DependsOn:  dependencies,
	}

	nicProperties := network.InterfacePropertiesFormat{
		IPConfigurations: &[]network.InterfaceIPConfiguration{
			{
				Name: helpers.PointerToString("ipconfig1"),
				InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
					Subnet: &network.Subnet{
						ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
					},
					Primary:                   helpers.PointerToBool(true),
					PrivateIPAllocationMethod: network.Dynamic,
					PublicIPAddress: &network.PublicIPAddress{
						ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses', variables('jumpboxPublicIpAddressName'))]"),
					},
				},
			},
		},
		NetworkSecurityGroup: &network.SecurityGroup{
			ID: helpers.PointerToString("[resourceId('Microsoft.Network/networkSecurityGroups', variables('jumpboxNetworkSecurityGroupName'))]"),
		},
	}

	networkInterface := network.Interface{
		Location:                  helpers.PointerToString("[variables('location')]"),
		Name:                      helpers.PointerToString("[variables('jumpboxNetworkInterfaceName')]"),
		InterfacePropertiesFormat: &nicProperties,
		Type:                      helpers.PointerToString("Microsoft.Network/networkInterfaces"),
	}

	return NetworkInterfaceARM{
		ARMResource: armResource,
		Interface:   networkInterface,
	}
}

func createAgentVMASNetworkInterface(cs *api.ContainerService, profile *api.AgentPoolProfile) NetworkInterfaceARM {
	isWindows := profile.IsWindows()
	isCustomVNet := profile.IsCustomVNET()
	isAzureCNI := cs.Properties.OrchestratorProfile.IsAzureCNI()

	armResource := ARMResource{
		APIVersion: "[variables('apiVersionNetwork')]",
		Copy: map[string]string{
			"count": fmt.Sprintf("[sub(variables('%[1]sCount'), variables('%[1]sOffset'))]", profile.Name),
			"name":  "loop",
		},
	}

	var dependencies []string

	if isCustomVNet {
		dependencies = append(dependencies, "[variables('nsgID')]")
	} else {
		dependencies = append(dependencies, "[variables('vnetID')]")
	}
	if profile.LoadBalancerBackendAddressPoolIDs == nil &&
		cs.Properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku == api.StandardLoadBalancerSku {
		dependencies = append(dependencies, "[variables('agentLbID')]")
	}

	armResource.DependsOn = dependencies

	networkInterface := network.Interface{
		Type:     helpers.PointerToString("Microsoft.Network/networkInterfaces"),
		Name:     helpers.PointerToString("[concat(variables('" + profile.Name + "VMNamePrefix'), 'nic-', copyIndex(variables('" + profile.Name + "Offset')))]"),
		Location: helpers.PointerToString("[variables('location')]"),
	}

	networkInterface.InterfacePropertiesFormat = &network.InterfacePropertiesFormat{}

	if isCustomVNet {
		networkInterface.NetworkSecurityGroup = &network.SecurityGroup{
			ID: helpers.PointerToString("[variables('nsgID')]"),
		}
	}

	if isWindows {
		networkInterface.EnableAcceleratedNetworking = profile.AcceleratedNetworkingEnabledWindows
	} else {
		networkInterface.EnableAcceleratedNetworking = profile.AcceleratedNetworkingEnabled
	}

	var ipConfigurations []network.InterfaceIPConfiguration
	for i := 1; i <= profile.IPAddressCount; i++ {
		ipConfig := network.InterfaceIPConfiguration{
			Name:                                     helpers.PointerToString(fmt.Sprintf("ipconfig%d", i)),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{},
		}
		if i == 1 {
			ipConfig.Primary = helpers.PointerToBool(true)
			backendPools := make([]network.BackendAddressPool, 0)
			if profile.LoadBalancerBackendAddressPoolIDs != nil {
				for _, lbBackendPoolID := range profile.LoadBalancerBackendAddressPoolIDs {
					backendPools = append(backendPools,
						network.BackendAddressPool{
							ID: helpers.PointerToString(lbBackendPoolID),
						},
					)
				}
			} else {
				if cs.Properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku == api.StandardLoadBalancerSku {
					agentLbBackendAddressPools := network.BackendAddressPool{
						ID: helpers.PointerToString("[concat(variables('agentLbID'), '/backendAddressPools/', variables('agentLbBackendPoolName'))]"),
					}
					backendPools = append(backendPools, agentLbBackendAddressPools)
				}
			}
			ipConfig.LoadBalancerBackendAddressPools = &backendPools
		}
		ipConfig.PrivateIPAllocationMethod = network.Dynamic
		ipConfig.Subnet = &network.Subnet{
			ID: helpers.PointerToString(fmt.Sprintf("[variables('%sVnetSubnetID')]", profile.Name)),
		}
		if !isWindows {
			if profile.Role == "Infra" {
				ipConfig.LoadBalancerBackendAddressPools = &[]network.BackendAddressPool{
					{
						ID: helpers.PointerToString("[concat(resourceId('Microsoft.Network/loadBalancers', variables('routerLBName')), '/backendAddressPools/backend')]"),
					},
				}
			}
		}

		if cs.Properties.FeatureFlags.IsFeatureEnabled("EnableIPv6DualStack") || cs.Properties.FeatureFlags.IsFeatureEnabled("EnableIPv6Only") {
			if cs.Properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku != api.StandardLoadBalancerSku {
				var backendPools []network.BackendAddressPool
				if ipConfig.LoadBalancerBackendAddressPools != nil {
					backendPools = *ipConfig.LoadBalancerBackendAddressPools
				}
				backendPools = append(backendPools, network.BackendAddressPool{
					ID: helpers.PointerToString("[concat(resourceId('Microsoft.Network/loadBalancers',parameters('masterEndpointDNSNamePrefix')), '/backendAddressPools/', parameters('masterEndpointDNSNamePrefix'))]"),
				})
				ipConfig.LoadBalancerBackendAddressPools = &backendPools
			}
		}
		ipConfigurations = append(ipConfigurations, ipConfig)
	}

	// add ipv6 nic config for dual stack
	if cs.Properties.FeatureFlags.IsFeatureEnabled("EnableIPv6DualStack") || cs.Properties.FeatureFlags.IsFeatureEnabled("EnableIPv6Only") {
		ipv6Config := network.InterfaceIPConfiguration{
			Name: helpers.PointerToString("ipconfigv6"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				PrivateIPAddressVersion: "IPv6",
				Primary:                 helpers.PointerToBool(false),
				Subnet: &network.Subnet{
					ID: helpers.PointerToString(fmt.Sprintf("[variables('%sVnetSubnetID')]", profile.Name)),
				},
			},
		}
		ipConfigurations = append(ipConfigurations, ipv6Config)
	}

	networkInterface.IPConfigurations = &ipConfigurations

	if !isAzureCNI && !cs.Properties.IsAzureStackCloud() {
		networkInterface.EnableIPForwarding = helpers.PointerToBool(true)
	}

	// Enable IPForwarding on NetworkInterface for azurecni dualstack
	if isAzureCNI {
		if cs.Properties.FeatureFlags.IsFeatureEnabled("EnableIPv6DualStack") {
			networkInterface.EnableIPForwarding = helpers.PointerToBool(true)
		}
	}

	return NetworkInterfaceARM{
		ARMResource: armResource,
		Interface:   networkInterface,
	}
}

func getSecondaryNICIPConfigs(n int) []network.InterfaceIPConfiguration {
	var ipConfigurations []network.InterfaceIPConfiguration
	for i := 2; i <= n; i++ {
		ipConfig := network.InterfaceIPConfiguration{
			Name: helpers.PointerToString(fmt.Sprintf("ipconfig%d", i)),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				Primary:                   helpers.PointerToBool(false),
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		}
		ipConfigurations = append(ipConfigurations, ipConfig)
	}
	return ipConfigurations
}
