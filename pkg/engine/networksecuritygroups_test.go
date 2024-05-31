// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
	"github.com/google/go-cmp/cmp"
)

func TestCreateNetworkSecurityGroup(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			OrchestratorProfile: &api.OrchestratorProfile{
				KubernetesConfig: &api.KubernetesConfig{
					PrivateCluster: &api.PrivateCluster{
						Enabled: helpers.PointerToBool(false),
					},
				},
			},
			AgentPoolProfiles: []*api.AgentPoolProfile{
				{
					Name:   "fooAgent",
					OSType: "Linux",
				},
			},
			FeatureFlags: &api.FeatureFlags{},
		},
	}

	// Test create normal nsg

	actual := CreateNetworkSecurityGroup(cs)

	expected := NetworkSecurityGroupARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
		},
		SecurityGroup: network.SecurityGroup{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[variables('nsgName')]"),
			Type:     helpers.PointerToString("Microsoft.Network/networkSecurityGroups"),
			SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
				SecurityRules: &[]network.SecurityRule{
					{
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
					},
					{
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
					},
				},
			},
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing nsgs : %s", diff)
	}

	// Test Create NSG with windows and Block Outbound internet

	cs.Properties.AgentPoolProfiles = []*api.AgentPoolProfile{
		{
			Name:   "fooAgent",
			OSType: "Windows",
		},
	}

	cs.Properties.FeatureFlags.BlockOutboundInternet = true

	actual = CreateNetworkSecurityGroup(cs)

	rules := *expected.SecurityRules

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

	rules = append(rules, rdpRule, vnetRule, blockOutBoundRule, allowARMRule)

	expected.SecurityRules = &rules

	diff = cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing nsgs : %s", diff)
	}

	// Test private cluster

	cs.Properties.OrchestratorProfile.KubernetesConfig.PrivateCluster.Enabled = helpers.PointerToBool(true)

	actual = CreateNetworkSecurityGroup(cs)

	for _, rule := range rules {
		if helpers.String(rule.Name) == "allow_kube_tls" {
			source := "VirtualNetwork"
			rule.SourceAddressPrefix = &source
		}
	}

	diff = cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing nsgs : %s", diff)
	}
}

func TestCreateJumpboxNSG(t *testing.T) {
	expected := NetworkSecurityGroupARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
		},
		SecurityGroup: network.SecurityGroup{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[variables('jumpboxNetworkSecurityGroupName')]"),
			Type:     helpers.PointerToString("Microsoft.Network/networkSecurityGroups"),
			SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
				SecurityRules: &[]network.SecurityRule{
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
				},
			},
		},
	}

	actual := createJumpboxNSG()

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing nsgs : %s", diff)
	}
}
