// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
)

func CreateNetworkSecurityGroup(cs *api.ContainerService) NetworkSecurityGroupARM {
	armResource := ARMResource{
		APIVersion: "[variables('apiVersionNetwork')]",
	}

	sshRule := network.SecurityRule{
		Name: helpers.PointerToString("allow_ssh"),
		SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
			Access:                   network.SecurityRuleAccessAllow,
			Description:              helpers.PointerToString("Allow SSH traffic to master"),
			DestinationAddressPrefix: helpers.PointerToString("*"),
			DestinationPortRange:     helpers.PointerToString("22-22"),
			Direction:                network.SecurityRuleDirectionInbound,
			Priority:                 helpers.PointerToInt32(101),
			Protocol:                 network.SecurityRuleProtocolTCP,
			SourceAddressPrefix:      helpers.PointerToString("*"),
			SourcePortRange:          helpers.PointerToString("*"),
		},
	}

	kubeTLSRule := network.SecurityRule{
		Name: helpers.PointerToString("allow_kube_tls"),
		SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
			Access:                   network.SecurityRuleAccessAllow,
			Description:              helpers.PointerToString("Allow kube-apiserver (tls) traffic to master"),
			DestinationAddressPrefix: helpers.PointerToString("*"),
			DestinationPortRange:     helpers.PointerToString("443-443"),
			Direction:                network.SecurityRuleDirectionInbound,
			Priority:                 helpers.PointerToInt32(100),
			Protocol:                 network.SecurityRuleProtocolTCP,
			SourceAddressPrefix:      helpers.PointerToString("*"),
			SourcePortRange:          helpers.PointerToString("*"),
		},
	}

	if cs.Properties.OrchestratorProfile.IsPrivateCluster() {
		source := "VirtualNetwork"
		kubeTLSRule.SourceAddressPrefix = &source
	}

	securityRules := []network.SecurityRule{
		sshRule,
		kubeTLSRule,
	}

	if cs.Properties.HasWindows() {
		rdpRule := network.SecurityRule{
			Name: helpers.PointerToString("allow_rdp"),
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Access:                   network.SecurityRuleAccessAllow,
				Description:              helpers.PointerToString("Allow RDP traffic to master"),
				DestinationAddressPrefix: helpers.PointerToString("*"),
				DestinationPortRange:     helpers.PointerToString("3389-3389"),
				Direction:                network.SecurityRuleDirectionInbound,
				Priority:                 helpers.PointerToInt32(102),
				Protocol:                 network.SecurityRuleProtocolTCP,
				SourceAddressPrefix:      helpers.PointerToString("*"),
				SourcePortRange:          helpers.PointerToString("*"),
			},
		}

		securityRules = append(securityRules, rdpRule)
	}

	if cs.Properties.FeatureFlags.IsFeatureEnabled("BlockOutboundInternet") {
		vnetRule := network.SecurityRule{
			Name: helpers.PointerToString("allow_vnet"),
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Access:                   network.SecurityRuleAccessAllow,
				Description:              helpers.PointerToString("Allow outbound internet to vnet"),
				DestinationAddressPrefix: helpers.PointerToString("[parameters('masterSubnet')]"),
				DestinationPortRange:     helpers.PointerToString("*"),
				Direction:                network.SecurityRuleDirectionOutbound,
				Priority:                 helpers.PointerToInt32(110),
				Protocol:                 network.SecurityRuleProtocolAsterisk,
				SourceAddressPrefix:      helpers.PointerToString("VirtualNetwork"),
				SourcePortRange:          helpers.PointerToString("*"),
			},
		}

		blockOutBoundRule := network.SecurityRule{
			Name: helpers.PointerToString("block_outbound"),
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Access:                   network.SecurityRuleAccessDeny,
				Description:              helpers.PointerToString("Block outbound internet from master"),
				DestinationAddressPrefix: helpers.PointerToString("*"),
				DestinationPortRange:     helpers.PointerToString("*"),
				Direction:                network.SecurityRuleDirectionOutbound,
				Priority:                 helpers.PointerToInt32(120),
				Protocol:                 network.SecurityRuleProtocolAsterisk,
				SourceAddressPrefix:      helpers.PointerToString("*"),
				SourcePortRange:          helpers.PointerToString("*"),
			},
		}

		allowARMRule := network.SecurityRule{
			Name: helpers.PointerToString("allow_ARM"),
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Access:                   network.SecurityRuleAccessAllow,
				Description:              helpers.PointerToString("Allow outbound internet to ARM"),
				DestinationAddressPrefix: helpers.PointerToString("AzureResourceManager"),
				DestinationPortRange:     helpers.PointerToString("443"),
				Direction:                network.SecurityRuleDirectionOutbound,
				Priority:                 helpers.PointerToInt32(100),
				Protocol:                 network.SecurityRuleProtocolTCP,
				SourceAddressPrefix:      helpers.PointerToString("*"),
				SourcePortRange:          helpers.PointerToString("*"),
			},
		}

		securityRules = append(securityRules, vnetRule)
		securityRules = append(securityRules, blockOutBoundRule)
		securityRules = append(securityRules, allowARMRule)
	}

	nsg := network.SecurityGroup{
		Location: helpers.PointerToString("[variables('location')]"),
		Name:     helpers.PointerToString("[variables('nsgName')]"),
		Type:     helpers.PointerToString("Microsoft.Network/networkSecurityGroups"),
		SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
			SecurityRules: &securityRules,
		},
	}

	return NetworkSecurityGroupARM{
		ARMResource:   armResource,
		SecurityGroup: nsg,
	}
}

func createJumpboxNSG() NetworkSecurityGroupARM {
	armResource := ARMResource{
		APIVersion: "[variables('apiVersionNetwork')]",
	}

	securityRules := []network.SecurityRule{
		{
			Name: helpers.PointerToString("default-allow-ssh"),
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Priority:                 helpers.PointerToInt32(1000),
				Protocol:                 network.SecurityRuleProtocolTCP,
				Access:                   network.SecurityRuleAccessAllow,
				Direction:                network.SecurityRuleDirectionInbound,
				SourceAddressPrefix:      helpers.PointerToString("*"),
				SourcePortRange:          helpers.PointerToString("*"),
				DestinationAddressPrefix: helpers.PointerToString("*"),
				DestinationPortRange:     helpers.PointerToString("22"),
			},
		},
	}
	nsg := network.SecurityGroup{
		Location: helpers.PointerToString("[variables('location')]"),
		Name:     helpers.PointerToString("[variables('jumpboxNetworkSecurityGroupName')]"),
		Type:     helpers.PointerToString("Microsoft.Network/networkSecurityGroups"),
		SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
			SecurityRules: &securityRules,
		},
	}
	return NetworkSecurityGroupARM{
		ARMResource:   armResource,
		SecurityGroup: nsg,
	}
}
