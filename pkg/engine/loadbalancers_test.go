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

func TestCreateMasterLoadBalancer(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{
				Count: 1,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				KubernetesConfig: &api.KubernetesConfig{
					LoadBalancerSku: BasicLoadBalancerSku,
				},
			},
		},
	}
	actual := CreateMasterLoadBalancer(cs.Properties)

	expected := LoadBalancerARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			DependsOn: []string{
				"[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]",
			},
		},
		LoadBalancer: network.LoadBalancer{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[variables('masterLbName')]"),
			LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
				BackendAddressPools: &[]network.BackendAddressPool{
					{
						Name: helpers.PointerToString("[variables('masterLbBackendPoolName')]"),
					},
				},
				FrontendIPConfigurations: &[]network.FrontendIPConfiguration{
					{
						Name: helpers.PointerToString("[variables('masterLbIPConfigName')]"),
						FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
							PublicIPAddress: &network.PublicIPAddress{
								ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('masterPublicIPAddressName'))]"),
							},
						},
					},
				},
				LoadBalancingRules: &[]network.LoadBalancingRule{
					{
						Name: helpers.PointerToString("LBRuleHTTPS"),
						LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
							FrontendIPConfiguration: &network.SubResource{
								ID: helpers.PointerToString("[variables('masterLbIPConfigID')]"),
							},
							BackendAddressPool: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
							},
							Protocol:             network.TransportProtocolTCP,
							FrontendPort:         helpers.PointerToInt32(443),
							BackendPort:          helpers.PointerToInt32(443),
							EnableFloatingIP:     helpers.PointerToBool(false),
							IdleTimeoutInMinutes: helpers.PointerToInt32(5),
							LoadDistribution:     network.LoadDistributionDefault,
							Probe: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('masterLbID'),'/probes/tcpHTTPSProbe')]"),
							},
						},
					},
				},
				InboundNatRules: &[]network.InboundNatRule{
					{
						Name: helpers.PointerToString("[concat('SSH-', variables('masterVMNamePrefix'), 0)]"),
						InboundNatRulePropertiesFormat: &network.InboundNatRulePropertiesFormat{
							FrontendIPConfiguration: &network.SubResource{
								ID: helpers.PointerToString("[variables('masterLbIPConfigID')]"),
							},
							Protocol:         network.TransportProtocol("Tcp"),
							FrontendPort:     helpers.PointerToInt32(22),
							BackendPort:      helpers.PointerToInt32(22),
							EnableFloatingIP: helpers.PointerToBool(false),
						},
					},
				},
				Probes: &[]network.Probe{
					{
						Name: helpers.PointerToString("tcpHTTPSProbe"),
						ProbePropertiesFormat: &network.ProbePropertiesFormat{
							Protocol:          network.ProbeProtocolTCP,
							Port:              helpers.PointerToInt32(443),
							IntervalInSeconds: helpers.PointerToInt32(5),
							NumberOfProbes:    helpers.PointerToInt32(2),
						},
					},
				},
			},
			Sku: &network.LoadBalancerSku{
				Name: "[variables('loadBalancerSku')]",
			},
			Type: helpers.PointerToString("Microsoft.Network/loadBalancers"),
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing load balancers: %s", diff)
	}

}

func TestCreateMasterLoadBalancerPrivate(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{
				Count: 1,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorType:    Kubernetes,
				OrchestratorVersion: "1.16.4",
				KubernetesConfig: &api.KubernetesConfig{
					LoadBalancerSku: BasicLoadBalancerSku,
					PrivateCluster: &api.PrivateCluster{
						Enabled: helpers.PointerToBool(true),
					},
				},
			},
		},
	}
	actual := CreateMasterLoadBalancer(cs.Properties)

	expected := LoadBalancerARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			DependsOn: []string{
				"[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]",
			},
		},
		LoadBalancer: network.LoadBalancer{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[variables('masterLbName')]"),
			LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
				BackendAddressPools: &[]network.BackendAddressPool{
					{
						Name: helpers.PointerToString("[variables('masterLbBackendPoolName')]"),
					},
				},
				FrontendIPConfigurations: &[]network.FrontendIPConfiguration{
					{
						Name: helpers.PointerToString("[variables('masterLbIPConfigName')]"),
						FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
							PublicIPAddress: &network.PublicIPAddress{
								ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('masterPublicIPAddressName'))]"),
							},
						},
					},
				},
				OutboundRules: &[]network.OutboundRule{
					{
						Name: helpers.PointerToString("LBOutboundRule"),
						OutboundRulePropertiesFormat: &network.OutboundRulePropertiesFormat{
							FrontendIPConfigurations: &[]network.SubResource{
								{
									ID: helpers.PointerToString("[variables('masterLbIPConfigID')]"),
								},
							},
							BackendAddressPool: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
							},
							Protocol:               network.Protocol1All,
							IdleTimeoutInMinutes:   helpers.PointerToInt32(0),
							AllocatedOutboundPorts: helpers.PointerToInt32(0),
							EnableTCPReset:         helpers.PointerToBool(true),
						},
					},
				},
			},
			Sku: &network.LoadBalancerSku{
				Name: "[variables('loadBalancerSku')]",
			},
			Type: helpers.PointerToString("Microsoft.Network/loadBalancers"),
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing load balancers: %s", diff)
	}

}

func TestCreateLoadBalancerStandard(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{
				Count: 1,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				KubernetesConfig: &api.KubernetesConfig{
					LoadBalancerSku: api.StandardLoadBalancerSku,
				},
			},
		},
	}
	actual := CreateMasterLoadBalancer(cs.Properties)

	expected := LoadBalancerARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			DependsOn: []string{
				"[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]",
			},
		},
		LoadBalancer: network.LoadBalancer{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[variables('masterLbName')]"),
			LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
				BackendAddressPools: &[]network.BackendAddressPool{
					{
						Name: helpers.PointerToString("[variables('masterLbBackendPoolName')]"),
					},
				},
				FrontendIPConfigurations: &[]network.FrontendIPConfiguration{
					{
						Name: helpers.PointerToString("[variables('masterLbIPConfigName')]"),
						FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
							PublicIPAddress: &network.PublicIPAddress{
								ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('masterPublicIPAddressName'))]"),
							},
						},
					},
				},
				LoadBalancingRules: &[]network.LoadBalancingRule{
					{
						Name: helpers.PointerToString("LBRuleHTTPS"),
						LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
							FrontendIPConfiguration: &network.SubResource{
								ID: helpers.PointerToString("[variables('masterLbIPConfigID')]"),
							},
							BackendAddressPool: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
							},
							Protocol:             network.TransportProtocolTCP,
							FrontendPort:         helpers.PointerToInt32(443),
							BackendPort:          helpers.PointerToInt32(443),
							EnableFloatingIP:     helpers.PointerToBool(false),
							IdleTimeoutInMinutes: helpers.PointerToInt32(5),
							LoadDistribution:     network.LoadDistributionDefault,
							Probe: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('masterLbID'),'/probes/tcpHTTPSProbe')]"),
							},
						},
					},
					{
						Name: helpers.PointerToString("LBRuleUDP"),
						LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
							FrontendIPConfiguration: &network.SubResource{
								ID: helpers.PointerToString("[variables('masterLbIPConfigID')]"),
							},
							BackendAddressPool: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
							},
							Protocol:             network.TransportProtocolUDP,
							FrontendPort:         helpers.PointerToInt32(1123),
							BackendPort:          helpers.PointerToInt32(1123),
							EnableFloatingIP:     helpers.PointerToBool(false),
							IdleTimeoutInMinutes: helpers.PointerToInt32(5),
							LoadDistribution:     network.LoadDistributionDefault,
							Probe: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('masterLbID'),'/probes/tcpHTTPSProbe')]"),
							},
						},
					},
				},
				InboundNatRules: &[]network.InboundNatRule{
					{
						Name: helpers.PointerToString("[concat('SSH-', variables('masterVMNamePrefix'), 0)]"),
						InboundNatRulePropertiesFormat: &network.InboundNatRulePropertiesFormat{
							FrontendIPConfiguration: &network.SubResource{
								ID: helpers.PointerToString("[variables('masterLbIPConfigID')]"),
							},
							Protocol:         network.TransportProtocol("Tcp"),
							FrontendPort:     helpers.PointerToInt32(22),
							BackendPort:      helpers.PointerToInt32(22),
							EnableFloatingIP: helpers.PointerToBool(false),
						},
					},
				},
				Probes: &[]network.Probe{
					{
						Name: helpers.PointerToString("tcpHTTPSProbe"),
						ProbePropertiesFormat: &network.ProbePropertiesFormat{
							Protocol:          network.ProbeProtocolTCP,
							Port:              helpers.PointerToInt32(443),
							IntervalInSeconds: helpers.PointerToInt32(5),
							NumberOfProbes:    helpers.PointerToInt32(2),
						},
					},
				},
			},
			Sku: &network.LoadBalancerSku{
				Name: "[variables('loadBalancerSku')]",
			},
			Type: helpers.PointerToString("Microsoft.Network/loadBalancers"),
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing load balancers: %s", diff)
	}

}

func TestCreateMasterInternalLoadBalancer(t *testing.T) {
	// Test with Basic LB
	cs := &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{},
			OrchestratorProfile: &api.OrchestratorProfile{
				KubernetesConfig: &api.KubernetesConfig{
					LoadBalancerSku: BasicLoadBalancerSku,
				},
			},
		},
	}

	actual := CreateMasterInternalLoadBalancer(cs)

	expected := LoadBalancerARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			DependsOn: []string{
				"[variables('vnetID')]",
			},
		},
		LoadBalancer: network.LoadBalancer{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[variables('masterInternalLbName')]"),
			LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
				BackendAddressPools: &[]network.BackendAddressPool{
					{
						Name: helpers.PointerToString("[variables('masterLbBackendPoolName')]"),
					},
				},
				FrontendIPConfigurations: &[]network.FrontendIPConfiguration{
					{
						Name: helpers.PointerToString("[variables('masterInternalLbIPConfigName')]"),
						FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
							PrivateIPAddress:          helpers.PointerToString("[variables('kubernetesAPIServerIP')]"),
							PrivateIPAllocationMethod: network.Static,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
				},
				LoadBalancingRules: &[]network.LoadBalancingRule{
					{
						Name: helpers.PointerToString("InternalLBRuleHTTPS"),
						LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
							BackendAddressPool: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('masterInternalLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
							},
							BackendPort:      helpers.PointerToInt32(4443),
							EnableFloatingIP: helpers.PointerToBool(false),
							FrontendIPConfiguration: &network.SubResource{
								ID: helpers.PointerToString("[variables('masterInternalLbIPConfigID')]"),
							},
							FrontendPort:         helpers.PointerToInt32(443),
							IdleTimeoutInMinutes: helpers.PointerToInt32(5),
							Protocol:             network.TransportProtocolTCP,
							Probe: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('masterInternalLbID'),'/probes/tcpHTTPSProbe')]"),
							},
						},
					},
				},
				Probes: &[]network.Probe{
					{
						Name: helpers.PointerToString("tcpHTTPSProbe"),
						ProbePropertiesFormat: &network.ProbePropertiesFormat{
							IntervalInSeconds: helpers.PointerToInt32(5),
							NumberOfProbes:    helpers.PointerToInt32(2),
							Port:              helpers.PointerToInt32(4443),
							Protocol:          network.ProbeProtocolTCP,
						},
					},
				},
			},
			Sku: &network.LoadBalancerSku{
				Name: network.LoadBalancerSkuName("[variables('loadBalancerSku')]"),
			},
			Type: helpers.PointerToString("Microsoft.Network/loadBalancers"),
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing load balancers: %s", diff)
	}

	// Test with Standard LB
	cs = &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{},
			OrchestratorProfile: &api.OrchestratorProfile{
				KubernetesConfig: &api.KubernetesConfig{
					LoadBalancerSku: api.StandardLoadBalancerSku,
				},
			},
		},
	}

	actual = CreateMasterInternalLoadBalancer(cs)

	expected.LoadBalancingRules = &[]network.LoadBalancingRule{
		{
			Name: helpers.PointerToString("InternalLBRuleHTTPS"),
			LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
				BackendAddressPool: &network.SubResource{
					ID: helpers.PointerToString("[concat(variables('masterInternalLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
				},
				BackendPort:      helpers.PointerToInt32(4443),
				EnableFloatingIP: helpers.PointerToBool(false),
				FrontendIPConfiguration: &network.SubResource{
					ID: helpers.PointerToString("[variables('masterInternalLbIPConfigID')]"),
				},
				FrontendPort:         helpers.PointerToInt32(443),
				IdleTimeoutInMinutes: helpers.PointerToInt32(5),
				Protocol:             network.TransportProtocolTCP,
				Probe: &network.SubResource{
					ID: helpers.PointerToString("[concat(variables('masterInternalLbID'),'/probes/tcpHTTPSProbe')]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("LBRuleUDP"),
			LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
				BackendAddressPool: &network.SubResource{
					ID: helpers.PointerToString("[concat(variables('masterInternalLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
				},
				BackendPort:      helpers.PointerToInt32(1123),
				EnableFloatingIP: helpers.PointerToBool(false),
				FrontendIPConfiguration: &network.SubResource{
					ID: helpers.PointerToString("[variables('masterInternalLbIPConfigID')]"),
				},
				FrontendPort:         helpers.PointerToInt32(1123),
				IdleTimeoutInMinutes: helpers.PointerToInt32(5),
				Protocol:             network.TransportProtocolUDP,
				Probe: &network.SubResource{
					ID: helpers.PointerToString("[concat(variables('masterInternalLbID'),'/probes/tcpHTTPSProbe')]"),
				},
			},
		},
	}

	diff = cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing load balancers: %s", diff)
	}

	// Test with custom Vnet
	cs = &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{
				VnetSubnetID: "fooSubnet",
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				KubernetesConfig: &api.KubernetesConfig{
					LoadBalancerSku: api.StandardLoadBalancerSku,
				},
			},
		},
	}

	actual = CreateMasterInternalLoadBalancer(cs)

	expected.DependsOn = []string{
		"[variables('nsgID')]",
	}

	diff = cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing load balancers: %s", diff)
	}

	// Test with VMSS
	cs = &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{
				VnetSubnetID:        "fooSubnet",
				AvailabilityProfile: api.VirtualMachineScaleSets,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				KubernetesConfig: &api.KubernetesConfig{
					LoadBalancerSku: api.StandardLoadBalancerSku,
				},
			},
		},
	}

	actual = CreateMasterInternalLoadBalancer(cs)

	expected.FrontendIPConfigurations = &[]network.FrontendIPConfiguration{
		{
			Name: helpers.PointerToString("[variables('masterInternalLbIPConfigName')]"),
			FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
				PrivateIPAddress:          helpers.PointerToString("[variables('kubernetesAPIServerIP')]"),
				PrivateIPAllocationMethod: network.Static,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetIDMaster')]"),
				},
			},
		},
	}

	diff = cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing load balancers: %s", diff)
	}
}

// TestCreateClusterLoadBalancerForIPv6 is a simple test..This setup and test will eventually
// be removed once the platform is enhanced and there'll be no requirement for having an ipv6
// fe to allow egress.
func TestCreateClusterLoadBalancerForIPv6(t *testing.T) {
	actual := CreateClusterLoadBalancerForIPv6()

	expected := LoadBalancerARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			DependsOn: []string{
				"[concat('Microsoft.Network/publicIPAddresses/', 'fee-ipv4')]",
			},
		},
		LoadBalancer: network.LoadBalancer{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[parameters('masterEndpointDNSNamePrefix')]"),
			LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
				BackendAddressPools: &[]network.BackendAddressPool{
					{
						Name: helpers.PointerToString("[parameters('masterEndpointDNSNamePrefix')]"),
					},
				},
				FrontendIPConfigurations: &[]network.FrontendIPConfiguration{
					{
						Name: helpers.PointerToString("LBFE-v4"),
						FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
							PublicIPAddress: &network.PublicIPAddress{
								ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses', 'fee-ipv4')]"),
							},
						},
					},
				},
				LoadBalancingRules: &[]network.LoadBalancingRule{
					{
						Name: helpers.PointerToString("LBRuleIPv4"),
						LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
							FrontendIPConfiguration: &network.SubResource{
								ID: helpers.PointerToString("[resourceId('Microsoft.Network/loadBalancers/frontendIpConfigurations', parameters('masterEndpointDNSNamePrefix'), 'LBFE-v4')]"),
							},
							BackendAddressPool: &network.SubResource{
								ID: helpers.PointerToString("[resourceId('Microsoft.Network/loadBalancers/backendAddressPools', parameters('masterEndpointDNSNamePrefix'), parameters('masterEndpointDNSNamePrefix'))]"),
							},
							Protocol:     network.TransportProtocolTCP,
							FrontendPort: helpers.PointerToInt32(9090),
							BackendPort:  helpers.PointerToInt32(9090),
						},
					},
				},
			},
			Sku: &network.LoadBalancerSku{
				Name: "[variables('loadBalancerSku')]",
			},
			Type: helpers.PointerToString("Microsoft.Network/loadBalancers"),
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing load balancers: %s", diff)
	}
}

func TestCreateAgentLoadBalancer(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{
				Count: 1,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorVersion: "1.18.2",
				KubernetesConfig: &api.KubernetesConfig{
					LoadBalancerSku: StandardLoadBalancerSku,
				},
			},
		},
	}
	actual := CreateStandardLoadBalancerForNodePools(cs.Properties, false)

	expected := LoadBalancerARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			DependsOn: []string{
				"[concat('Microsoft.Network/publicIPAddresses/', variables('agentPublicIPAddressName'))]",
			},
		},
		LoadBalancer: network.LoadBalancer{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[variables('agentLbName')]"),
			LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
				BackendAddressPools: &[]network.BackendAddressPool{
					{
						Name: helpers.PointerToString("[variables('agentLbBackendPoolName')]"),
					},
				},
				FrontendIPConfigurations: &[]network.FrontendIPConfiguration{
					{
						Name: helpers.PointerToString("[variables('agentLbIPConfigName')]"),
						FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
							PublicIPAddress: &network.PublicIPAddress{
								ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('agentPublicIPAddressName'))]"),
							},
						},
					},
				},
				OutboundRules: &[]network.OutboundRule{
					{
						Name: helpers.PointerToString("LBOutboundRule"),
						OutboundRulePropertiesFormat: &network.OutboundRulePropertiesFormat{
							FrontendIPConfigurations: &[]network.SubResource{
								{
									ID: helpers.PointerToString("[variables('agentLbIPConfigID')]"),
								},
							},
							BackendAddressPool: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('agentLbID'), '/backendAddressPools/', variables('agentLbBackendPoolName'))]"),
							},
							Protocol:               network.Protocol1All,
							IdleTimeoutInMinutes:   helpers.PointerToInt32(cs.Properties.OrchestratorProfile.KubernetesConfig.OutboundRuleIdleTimeoutInMinutes),
							EnableTCPReset:         helpers.PointerToBool(true),
							AllocatedOutboundPorts: helpers.PointerToInt32(0),
						},
					},
				},
			},
			Sku: &network.LoadBalancerSku{
				Name: "[variables('loadBalancerSku')]",
			},
			Type: helpers.PointerToString("Microsoft.Network/loadBalancers"),
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing load balancers: %s", diff)
	}

	// Test with > 1 LB outbound IP address
	cs = &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{
				Count: 1,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorVersion: "1.18.2",
				KubernetesConfig: &api.KubernetesConfig{
					LoadBalancerSku:         StandardLoadBalancerSku,
					LoadBalancerOutboundIPs: helpers.PointerToInt(6),
				},
			},
		},
	}
	actual = CreateStandardLoadBalancerForNodePools(cs.Properties, false)

	expected.LoadBalancer.LoadBalancerPropertiesFormat.FrontendIPConfigurations = &[]network.FrontendIPConfiguration{
		{
			Name: helpers.PointerToString("[variables('agentLbIPConfigName')]"),
			FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
				PublicIPAddress: &network.PublicIPAddress{
					ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('agentPublicIPAddressName'))]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("[variables('agentLbIPConfigName2')]"),
			FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
				PublicIPAddress: &network.PublicIPAddress{
					ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('agentPublicIPAddressName2'))]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("[variables('agentLbIPConfigName3')]"),
			FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
				PublicIPAddress: &network.PublicIPAddress{
					ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('agentPublicIPAddressName3'))]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("[variables('agentLbIPConfigName4')]"),
			FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
				PublicIPAddress: &network.PublicIPAddress{
					ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('agentPublicIPAddressName4'))]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("[variables('agentLbIPConfigName5')]"),
			FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
				PublicIPAddress: &network.PublicIPAddress{
					ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('agentPublicIPAddressName5'))]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("[variables('agentLbIPConfigName6')]"),
			FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
				PublicIPAddress: &network.PublicIPAddress{
					ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('agentPublicIPAddressName6'))]"),
				},
			},
		},
	}
	expected.LoadBalancer.LoadBalancerPropertiesFormat.OutboundRules = &[]network.OutboundRule{
		{
			Name: helpers.PointerToString("LBOutboundRule"),
			OutboundRulePropertiesFormat: &network.OutboundRulePropertiesFormat{
				FrontendIPConfigurations: &[]network.SubResource{
					{
						ID: helpers.PointerToString("[variables('agentLbIPConfigID')]"),
					},
					{
						ID: helpers.PointerToString("[variables('agentLbIPConfigID2')]"),
					},
					{
						ID: helpers.PointerToString("[variables('agentLbIPConfigID3')]"),
					},
					{
						ID: helpers.PointerToString("[variables('agentLbIPConfigID4')]"),
					},
					{
						ID: helpers.PointerToString("[variables('agentLbIPConfigID5')]"),
					},
					{
						ID: helpers.PointerToString("[variables('agentLbIPConfigID6')]"),
					},
				},
				BackendAddressPool: &network.SubResource{
					ID: helpers.PointerToString("[concat(variables('agentLbID'), '/backendAddressPools/', variables('agentLbBackendPoolName'))]"),
				},
				Protocol:               network.Protocol1All,
				IdleTimeoutInMinutes:   helpers.PointerToInt32(cs.Properties.OrchestratorProfile.KubernetesConfig.OutboundRuleIdleTimeoutInMinutes),
				EnableTCPReset:         helpers.PointerToBool(true),
				AllocatedOutboundPorts: helpers.PointerToInt32(0),
			},
		},
	}

	diff = cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing load balancers: %s", diff)
	}
}
