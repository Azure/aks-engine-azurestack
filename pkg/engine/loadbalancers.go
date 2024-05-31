// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"fmt"
	"strconv"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
	"github.com/blang/semver"
)

// CreateClusterLoadBalancerForIPv6 creates the cluster loadbalancer with IPv4 and IPv6 FE config
// this loadbalancer is created for the ipv6 dual stack feature and configured with 1 ipv4 FE, 1 ipv6 FE
// and 2 backend address pools - v4 and v6, 2 rules - v4 and v6. Atleast existence of 1 rule is a
// requirement now to allow egress. This can be removed later.
// TODO (aramase)
func CreateClusterLoadBalancerForIPv6() LoadBalancerARM {
	loadbalancer := LoadBalancerARM{
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
						// cluster name used as backend addr pool name for ipv4 to ensure backward compat
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
	return loadbalancer
}

// CreateMasterLoadBalancer creates a master LB
// In a private cluster scenario, we don't attach the inbound foo, e.g., TCP 443 and SSH access
func CreateMasterLoadBalancer(prop *api.Properties) LoadBalancerARM {
	loadBalancer := LoadBalancerARM{
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
			},
			Sku: &network.LoadBalancerSku{
				Name: "[variables('loadBalancerSku')]",
			},
			Type: helpers.PointerToString("Microsoft.Network/loadBalancers"),
		},
	}

	if !prop.OrchestratorProfile.IsPrivateCluster() {
		loadBalancingRules := &[]network.LoadBalancingRule{
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
		}
		probes := &[]network.Probe{
			{
				Name: helpers.PointerToString("tcpHTTPSProbe"),
				ProbePropertiesFormat: &network.ProbePropertiesFormat{
					Protocol:          network.ProbeProtocolTCP,
					Port:              helpers.PointerToInt32(443),
					IntervalInSeconds: helpers.PointerToInt32(5),
					NumberOfProbes:    helpers.PointerToInt32(2),
				},
			},
		}
		loadBalancer.LoadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules = loadBalancingRules
		loadBalancer.LoadBalancer.LoadBalancerPropertiesFormat.Probes = probes
		if prop.OrchestratorProfile.KubernetesConfig.LoadBalancerSku == api.StandardLoadBalancerSku {
			udpRule := network.LoadBalancingRule{
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
			}
			*loadBalancer.LoadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules = append(*loadBalancer.LoadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules, udpRule)
		}
		var inboundNATRules []network.InboundNatRule
		sshNATPorts := []int32{
			22,
			2201,
			2202,
			2203,
			2204,
		}
		for i := 0; i < prop.MasterProfile.Count; i++ {
			inboundNATRule := network.InboundNatRule{
				Name: helpers.PointerToString(fmt.Sprintf("[concat('SSH-', variables('masterVMNamePrefix'), %d)]", i)),
				InboundNatRulePropertiesFormat: &network.InboundNatRulePropertiesFormat{
					BackendPort:      helpers.PointerToInt32(22),
					EnableFloatingIP: helpers.PointerToBool(false),
					FrontendIPConfiguration: &network.SubResource{
						ID: helpers.PointerToString("[variables('masterLbIPConfigID')]"),
					},
					FrontendPort: helpers.PointerToInt32(sshNATPorts[i]),
					Protocol:     network.TransportProtocolTCP,
				},
			}
			inboundNATRules = append(inboundNATRules, inboundNATRule)
		}
		loadBalancer.InboundNatRules = &inboundNATRules
	} else {
		outboundRules := createOutboundRules(prop)
		outboundRule := (*outboundRules)[0]
		outboundRule.OutboundRulePropertiesFormat.BackendAddressPool.ID = helpers.PointerToString("[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]")
		(*outboundRule.OutboundRulePropertiesFormat.FrontendIPConfigurations)[0].ID = helpers.PointerToString("[variables('masterLbIPConfigID')]")
		loadBalancer.LoadBalancer.LoadBalancerPropertiesFormat.OutboundRules = outboundRules
	}

	return loadBalancer
}

func createOutboundRules(prop *api.Properties) *[]network.OutboundRule {
	currentVersion, _ := semver.Make(prop.OrchestratorProfile.OrchestratorVersion)
	min13Version, _ := semver.Make("1.13.7")
	min14Version, _ := semver.Make("1.14.3")
	min15Version, _ := semver.Make("1.15.0")

	if currentVersion.LT(min13Version) ||
		(currentVersion.Major == min14Version.Major && currentVersion.Minor == min14Version.Minor && currentVersion.LT(min14Version)) ||
		(currentVersion.Major == min15Version.Major && currentVersion.Minor == min15Version.Minor && currentVersion.LT(min15Version)) {
		return &[]network.OutboundRule{
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
					IdleTimeoutInMinutes:   helpers.PointerToInt32(prop.OrchestratorProfile.KubernetesConfig.OutboundRuleIdleTimeoutInMinutes),
					AllocatedOutboundPorts: helpers.PointerToInt32(0),
				},
			},
		}
	}
	outboundRule := network.OutboundRule{
		Name: helpers.PointerToString("LBOutboundRule"),
		OutboundRulePropertiesFormat: &network.OutboundRulePropertiesFormat{
			BackendAddressPool: &network.SubResource{
				ID: helpers.PointerToString("[concat(variables('agentLbID'), '/backendAddressPools/', variables('agentLbBackendPoolName'))]"),
			},
			Protocol:               network.Protocol1All,
			IdleTimeoutInMinutes:   helpers.PointerToInt32(prop.OrchestratorProfile.KubernetesConfig.OutboundRuleIdleTimeoutInMinutes),
			EnableTCPReset:         helpers.PointerToBool(true),
			AllocatedOutboundPorts: helpers.PointerToInt32(0),
		},
	}
	numIps := 1
	if prop.OrchestratorProfile.KubernetesConfig.LoadBalancerOutboundIPs != nil {
		numIps = *prop.OrchestratorProfile.KubernetesConfig.LoadBalancerOutboundIPs
	}
	agentLbIPConfigIDPrefix := "agentLbIPConfigID"
	frontendIPConfigurations := &[]network.SubResource{}
	for i := 1; i <= numIps; i++ {
		name := agentLbIPConfigIDPrefix
		if i > 1 {
			name += strconv.Itoa(i)
		}
		*frontendIPConfigurations = append(*frontendIPConfigurations, network.SubResource{
			ID: helpers.PointerToString(fmt.Sprintf("[variables('%s')]", name)),
		})
	}
	outboundRule.OutboundRulePropertiesFormat.FrontendIPConfigurations = frontendIPConfigurations
	return &[]network.OutboundRule{outboundRule}
}

// CreateStandardLoadBalancerForNodePools returns an ARM resource for the Standard LB that has all nodes in its backend pool
func CreateStandardLoadBalancerForNodePools(prop *api.Properties, isVMSS bool) LoadBalancerARM {
	loadBalancer := LoadBalancerARM{
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
				OutboundRules: createOutboundRules(prop),
			},
			Sku: &network.LoadBalancerSku{
				Name: "[variables('loadBalancerSku')]",
			},
			Type: helpers.PointerToString("Microsoft.Network/loadBalancers"),
		},
	}

	numIps := 1
	if prop.OrchestratorProfile.KubernetesConfig.LoadBalancerOutboundIPs != nil {
		numIps = *prop.OrchestratorProfile.KubernetesConfig.LoadBalancerOutboundIPs
	}
	agentPublicIPAddressNamePrefix := "agentPublicIPAddressName"
	agentLbIPConfigNamePrefix := "agentLbIPConfigName"
	frontendIPConfigurations := &[]network.FrontendIPConfiguration{}
	for i := 1; i <= numIps; i++ {
		agentPublicIPAddressName := agentPublicIPAddressNamePrefix
		agentLbIPConfigName := agentLbIPConfigNamePrefix
		if i > 1 {
			agentPublicIPAddressName += strconv.Itoa(i)
			agentLbIPConfigName += strconv.Itoa(i)
		}
		*frontendIPConfigurations = append(*frontendIPConfigurations, network.FrontendIPConfiguration{
			Name: helpers.PointerToString(fmt.Sprintf("[variables('%s')]", agentLbIPConfigName)),
			FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
				PublicIPAddress: &network.PublicIPAddress{
					ID: helpers.PointerToString(fmt.Sprintf("[resourceId('Microsoft.Network/publicIpAddresses',variables('%s'))]", agentPublicIPAddressName)),
				},
			},
		})
	}
	loadBalancer.LoadBalancer.LoadBalancerPropertiesFormat.FrontendIPConfigurations = frontendIPConfigurations
	return loadBalancer
}

func CreateMasterInternalLoadBalancer(cs *api.ContainerService) LoadBalancerARM {
	var dependencies []string
	if cs.Properties.MasterProfile.IsCustomVNET() {
		dependencies = append(dependencies, "[variables('nsgID')]")
	} else {
		dependencies = append(dependencies, "[variables('vnetID')]")
	}

	armResource := ARMResource{
		APIVersion: "[variables('apiVersionNetwork')]",
		DependsOn:  dependencies,
	}
	subnet := "[variables('vnetSubnetID')]"
	if cs.Properties.MasterProfile.IsVirtualMachineScaleSets() {
		subnet = "[variables('vnetSubnetIDMaster')]"
	}

	loadBalancer := network.LoadBalancer{
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
							ID: helpers.PointerToString(subnet),
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
	}

	if cs.Properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku == api.StandardLoadBalancerSku {
		udpRule := network.LoadBalancingRule{
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
		}
		*loadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules = append(*loadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules, udpRule)
	}

	loadBalancerARM := LoadBalancerARM{
		ARMResource:  armResource,
		LoadBalancer: loadBalancer,
	}

	return loadBalancerARM
}
