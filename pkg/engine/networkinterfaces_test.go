// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"fmt"
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
	"github.com/google/go-cmp/cmp"
)

func TestCreateNIC(t *testing.T) {

	// Test Master NIC
	cs := &api.ContainerService{
		Properties: &api.Properties{
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ClientID: "barClientID",
				Secret:   "bazSecret",
			},
			MasterProfile: &api.MasterProfile{
				Count:               1,
				DNSPrefix:           "myprefix1",
				VMSize:              "Standard_DS2_v2",
				AvailabilityProfile: api.VirtualMachineScaleSets,
				IPAddressCount:      5,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorType:    api.Kubernetes,
				OrchestratorVersion: "1.10.2",
				KubernetesConfig: &api.KubernetesConfig{
					NetworkPlugin: "azure",
				},
			},
			LinuxProfile: &api.LinuxProfile{},
			AgentPoolProfiles: []*api.AgentPoolProfile{
				{
					Name:                "agentpool",
					VMSize:              "Standard_D2_v2",
					Count:               1,
					AvailabilityProfile: api.AvailabilitySet,
				},
			},
		},
	}

	nic := CreateMasterVMNetworkInterfaces(cs)

	expected := NetworkInterfaceARM{

		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			Copy: map[string]string{
				"count": "[sub(variables('masterCount'), variables('masterOffset'))]",
				"name":  "nicLoopNode",
			},
			DependsOn: []string{
				"[variables('vnetID')]",
				"[variables('masterLbName')]",
			},
		},
		Interface: network.Interface{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[concat(variables('masterVMNamePrefix'), 'nic-', copyIndex(variables('masterOffset')))]"),
			InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
				IPConfigurations: &[]network.InterfaceIPConfiguration{
					{
						Name: helpers.PointerToString("ipconfig1"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
								{
									ID: helpers.PointerToString("[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
								},
							},
							LoadBalancerInboundNatRules: &[]network.InboundNatRule{
								{
									ID: helpers.PointerToString("[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'),copyIndex(variables('masterOffset')))]"),
								},
							},
							PrivateIPAddress:          helpers.PointerToString("[variables('masterPrivateIpAddrs')[copyIndex(variables('masterOffset'))]]"),
							Primary:                   helpers.PointerToBool(true),
							PrivateIPAllocationMethod: network.Static,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
					{
						Name: helpers.PointerToString("ipconfig2"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							Primary:                   helpers.PointerToBool(false),
							PrivateIPAllocationMethod: network.Dynamic,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
					{
						Name: helpers.PointerToString("ipconfig3"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							Primary:                   helpers.PointerToBool(false),
							PrivateIPAllocationMethod: network.Dynamic,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
					{
						Name: helpers.PointerToString("ipconfig4"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							Primary:                   helpers.PointerToBool(false),
							PrivateIPAllocationMethod: network.Dynamic,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
					{
						Name: helpers.PointerToString("ipconfig5"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							Primary:                   helpers.PointerToBool(false),
							PrivateIPAllocationMethod: network.Dynamic,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
				},
			},
			Type: helpers.PointerToString("Microsoft.Network/networkInterfaces"),
		},
	}

	diff := cmp.Diff(nic, expected)

	if diff != "" {
		t.Errorf("unexpected diff while expecting equal structs: %s", diff)
	}

	// Test Master NIC with custom Vnet
	cs.Properties.MasterProfile.VnetSubnetID = "fooSubnet"

	expected.DependsOn = []string{
		"[variables('nsgID')]",
		"[variables('masterLbName')]",
	}

	expected.NetworkSecurityGroup = &network.SecurityGroup{
		ID: helpers.PointerToString("[variables('nsgID')]"),
	}

	nic = CreateMasterVMNetworkInterfaces(cs)

	diff = cmp.Diff(nic, expected)

	if diff != "" {
		t.Errorf("unexpected diff while expecting equal structs: %s", diff)
	}

	// Test Master NIC with MultiMaster

	cs.Properties.MasterProfile.Count = 3

	nic = CreateMasterVMNetworkInterfaces(cs)

	expected.DependsOn = []string{
		"[variables('nsgID')]",
		"[variables('masterInternalLbName')]",
		"[variables('masterLbName')]",
	}

	expected.IPConfigurations = &[]network.InterfaceIPConfiguration{
		{
			Name: helpers.PointerToString("ipconfig1"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
					{
						ID: helpers.PointerToString("[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
					},
					{
						ID: helpers.PointerToString("[concat(variables('masterInternalLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
					},
				},
				LoadBalancerInboundNatRules: &[]network.InboundNatRule{
					{
						ID: helpers.PointerToString("[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'),copyIndex(variables('masterOffset')))]"),
					},
				},
				PrivateIPAddress:          helpers.PointerToString("[variables('masterPrivateIpAddrs')[copyIndex(variables('masterOffset'))]]"),
				Primary:                   helpers.PointerToBool(true),
				PrivateIPAllocationMethod: network.Static,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("ipconfig2"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				Primary:                   helpers.PointerToBool(false),
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("ipconfig3"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				Primary:                   helpers.PointerToBool(false),
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("ipconfig4"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				Primary:                   helpers.PointerToBool(false),
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("ipconfig5"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				Primary:                   helpers.PointerToBool(false),
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		},
	}

	diff = cmp.Diff(nic, expected)

	if diff != "" {
		t.Errorf("unexpected diff while expecting equal structs: %s", diff)
	}

	// Test Master NIC with Cosmos etcd

	cs.Properties.MasterProfile.CosmosEtcd = helpers.PointerToBool(true)
	cs.Properties.MasterProfile.Count = 3

	nic = CreateMasterVMNetworkInterfaces(cs)
	expected.DependsOn = []string{
		"[variables('nsgID')]",
		"[variables('masterInternalLbName')]",
		"[resourceId('Microsoft.DocumentDB/databaseAccounts/', variables('cosmosAccountName'))]",
		"[variables('masterLbName')]",
	}
	diff = cmp.Diff(nic, expected)

	if diff != "" {
		t.Errorf("unexpected diff while expecting equal structs: %s", diff)
	}

	// Test Master NIC without AzureCNI and customNodes DNS

	cs.Properties.MasterProfile.IPAddressCount = 5
	cs.Properties.LinuxProfile = &api.LinuxProfile{
		CustomNodesDNS: &api.CustomNodesDNS{
			DNSServer: "barServer",
		},
	}
	nic = CreateMasterVMNetworkInterfaces(cs)
	expected.Interface.DNSSettings = &network.InterfaceDNSSettings{
		DNSServers: &[]string{
			"[parameters('dnsServer')]",
		},
	}
	diff = cmp.Diff(nic, expected)

	if diff != "" {
		t.Errorf("unexpected diff while expecting equal structs: %s", diff)
	}
}

func TestCreatePrivateClusterNetworkInterface(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ClientID: "barClientID",
				Secret:   "bazSecret",
			},
			MasterProfile: &api.MasterProfile{
				Count:               1,
				DNSPrefix:           "myprefix1",
				VMSize:              "Standard_DS2_v2",
				AvailabilityProfile: api.VirtualMachineScaleSets,
				IPAddressCount:      5,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorType:    api.Kubernetes,
				OrchestratorVersion: "1.10.2",
				KubernetesConfig: &api.KubernetesConfig{
					NetworkPlugin: "azure",
				},
			},
			LinuxProfile: &api.LinuxProfile{},
			AgentPoolProfiles: []*api.AgentPoolProfile{
				{
					Name:                "agentpool",
					VMSize:              "Standard_D2_v2",
					Count:               1,
					AvailabilityProfile: api.AvailabilitySet,
				},
			},
		},
	}

	actual := createPrivateClusterMasterVMNetworkInterface(cs)

	expected := NetworkInterfaceARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			Copy: map[string]string{
				"count": "[sub(variables('masterCount'), variables('masterOffset'))]",
				"name":  "nicLoopNode",
			},
			DependsOn: []string{
				"[variables('vnetID')]",
			},
		},
		Interface: network.Interface{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[concat(variables('masterVMNamePrefix'), 'nic-', copyIndex(variables('masterOffset')))]"),
			InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
				IPConfigurations: &[]network.InterfaceIPConfiguration{
					{
						Name: helpers.PointerToString("ipconfig1"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							PrivateIPAddress:          helpers.PointerToString("[variables('masterPrivateIpAddrs')[copyIndex(variables('masterOffset'))]]"),
							Primary:                   helpers.PointerToBool(true),
							PrivateIPAllocationMethod: network.Static,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
					{
						Name: helpers.PointerToString("ipconfig2"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							Primary:                   helpers.PointerToBool(false),
							PrivateIPAllocationMethod: network.Dynamic,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
					{
						Name: helpers.PointerToString("ipconfig3"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							Primary:                   helpers.PointerToBool(false),
							PrivateIPAllocationMethod: network.Dynamic,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
					{
						Name: helpers.PointerToString("ipconfig4"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							Primary:                   helpers.PointerToBool(false),
							PrivateIPAllocationMethod: network.Dynamic,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
					{
						Name: helpers.PointerToString("ipconfig5"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							Primary:                   helpers.PointerToBool(false),
							PrivateIPAllocationMethod: network.Dynamic,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
				},
			},
			Type: helpers.PointerToString("Microsoft.Network/networkInterfaces"),
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}

	// Test private cluster NIC customVnet and multiple masters

	cs.Properties.MasterProfile.VnetSubnetID = "fooSubnet"
	cs.Properties.MasterProfile.Count = 3

	actual = createPrivateClusterMasterVMNetworkInterface(cs)

	expected.DependsOn = []string{
		"[variables('nsgID')]",
		"[variables('masterInternalLbName')]",
	}

	expected.Interface.NetworkSecurityGroup = &network.SecurityGroup{
		ID: helpers.PointerToString("[variables('nsgID')]"),
	}

	expected.IPConfigurations = &[]network.InterfaceIPConfiguration{
		{
			Name: helpers.PointerToString("ipconfig1"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
					{
						ID: helpers.PointerToString("[concat(variables('masterInternalLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
					},
				},
				LoadBalancerInboundNatRules: &[]network.InboundNatRule{},
				PrivateIPAddress:            helpers.PointerToString("[variables('masterPrivateIpAddrs')[copyIndex(variables('masterOffset'))]]"),
				Primary:                     helpers.PointerToBool(true),
				PrivateIPAllocationMethod:   network.Static,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("ipconfig2"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				Primary:                   helpers.PointerToBool(false),
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("ipconfig3"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				Primary:                   helpers.PointerToBool(false),
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("ipconfig4"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				Primary:                   helpers.PointerToBool(false),
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("ipconfig5"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				Primary:                   helpers.PointerToBool(false),
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		},
	}

	diff = cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}

	// Test Master NIC without AzureCNI and customNodes DNS

	cs.Properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "notazure"
	cs.Properties.LinuxProfile = &api.LinuxProfile{
		CustomNodesDNS: &api.CustomNodesDNS{
			DNSServer: "barServer",
		},
	}
	actual = createPrivateClusterMasterVMNetworkInterface(cs)
	expected.EnableIPForwarding = helpers.PointerToBool(true)
	expected.IPConfigurations = &[]network.InterfaceIPConfiguration{
		{
			Name: helpers.PointerToString("ipconfig1"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
					{
						ID: helpers.PointerToString("[concat(variables('masterInternalLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
					},
				},
				LoadBalancerInboundNatRules: &[]network.InboundNatRule{},
				PrivateIPAddress:            helpers.PointerToString("[variables('masterPrivateIpAddrs')[copyIndex(variables('masterOffset'))]]"),
				Primary:                     helpers.PointerToBool(true),
				PrivateIPAllocationMethod:   network.Static,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
				},
			},
		},
	}

	expected.Interface.InterfacePropertiesFormat.DNSSettings = &network.InterfaceDNSSettings{
		DNSServers: &[]string{
			"[parameters('dnsServer')]",
		},
	}

	diff = cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}
}

func TestCreateJumpboxNIC(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ClientID: "barClientID",
				Secret:   "bazSecret",
			},
			MasterProfile: &api.MasterProfile{
				Count:               1,
				DNSPrefix:           "myprefix1",
				VMSize:              "Standard_DS2_v2",
				AvailabilityProfile: api.VirtualMachineScaleSets,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorType:    api.Kubernetes,
				OrchestratorVersion: "1.10.2",
				KubernetesConfig: &api.KubernetesConfig{
					NetworkPlugin: "azure",
				},
			},
			LinuxProfile: &api.LinuxProfile{},
			AgentPoolProfiles: []*api.AgentPoolProfile{
				{
					Name:                "agentpool",
					VMSize:              "Standard_D2_v2",
					Count:               1,
					AvailabilityProfile: api.AvailabilitySet,
				},
			},
		},
	}

	actual := createJumpboxNetworkInterface(cs)

	expected := NetworkInterfaceARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			DependsOn: []string{
				"[concat('Microsoft.Network/publicIpAddresses/', variables('jumpboxPublicIpAddressName'))]",
				"[concat('Microsoft.Network/networkSecurityGroups/', variables('jumpboxNetworkSecurityGroupName'))]",
				"[variables('vnetID')]",
			},
		},
		Interface: network.Interface{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[variables('jumpboxNetworkInterfaceName')]"),
			InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
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
			},
			Type: helpers.PointerToString("Microsoft.Network/networkInterfaces"),
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}
}
func TestCreateAgentVMASNICWithSLB(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ClientID: "barClientID",
				Secret:   "bazSecret",
			},
			MasterProfile: &api.MasterProfile{
				Count:               1,
				DNSPrefix:           "myprefix1",
				VMSize:              "Standard_DS2_v2",
				AvailabilityProfile: api.VirtualMachineScaleSets,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorType:    api.Kubernetes,
				OrchestratorVersion: "1.10.2",
				KubernetesConfig: &api.KubernetesConfig{
					NetworkPlugin:   "azure",
					LoadBalancerSku: StandardLoadBalancerSku,
				},
			},
			LinuxProfile: &api.LinuxProfile{},
			AgentPoolProfiles: []*api.AgentPoolProfile{
				{
					Name:                "agentpool",
					VMSize:              "Standard_D2_v2",
					Count:               1,
					AvailabilityProfile: api.AvailabilitySet,
				},
			},
		},
	}

	profile := &api.AgentPoolProfile{
		Name:           "fooAgent",
		OSType:         "Linux",
		IPAddressCount: 1,
	}

	// Test AgentVMAS NIC with Standard LB, should add dependsOn for agentLbID and adds agentLbBackendPoolName as backendaddress pool

	actual := createAgentVMASNetworkInterface(cs, profile)

	expected := NetworkInterfaceARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			Copy: map[string]string{
				"count": "[sub(variables('fooAgentCount'), variables('fooAgentOffset'))]",
				"name":  "loop",
			},
			DependsOn: []string{
				"[variables('vnetID')]",
				"[variables('agentLbID')]",
			},
		},
		Interface: network.Interface{
			Type:     helpers.PointerToString("Microsoft.Network/networkInterfaces"),
			Name:     helpers.PointerToString("[concat(variables('fooAgentVMNamePrefix'), 'nic-', copyIndex(variables('fooAgentOffset')))]"),
			Location: helpers.PointerToString("[variables('location')]"),
			InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
				IPConfigurations: &[]network.InterfaceIPConfiguration{
					{
						Name: helpers.PointerToString("ipconfig1"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
								{
									ID: helpers.PointerToString("[concat(variables('agentLbID'), '/backendAddressPools/', variables('agentLbBackendPoolName'))]"),
								},
							},
							Primary:                   helpers.PointerToBool(true),
							PrivateIPAllocationMethod: network.Dynamic,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString(fmt.Sprintf("[variables('%sVnetSubnetID')]", profile.Name)),
							},
						},
					},
				},
			},
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}
}

func TestCreateAgentVMASNIC(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ClientID: "barClientID",
				Secret:   "bazSecret",
			},
			MasterProfile: &api.MasterProfile{
				Count:               1,
				DNSPrefix:           "myprefix1",
				VMSize:              "Standard_DS2_v2",
				AvailabilityProfile: api.VirtualMachineScaleSets,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorType:    api.Kubernetes,
				OrchestratorVersion: "1.10.2",
				KubernetesConfig: &api.KubernetesConfig{
					NetworkPlugin: "azure",
				},
			},
			LinuxProfile: &api.LinuxProfile{},
			AgentPoolProfiles: []*api.AgentPoolProfile{
				{
					Name:                "agentpool",
					VMSize:              "Standard_D2_v2",
					Count:               1,
					AvailabilityProfile: api.AvailabilitySet,
				},
			},
		},
	}

	profile := &api.AgentPoolProfile{
		Name:                              "fooAgent",
		OSType:                            "Linux",
		Role:                              "Infra",
		LoadBalancerBackendAddressPoolIDs: []string{"/subscriptions/123/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/mySLB/backendAddressPools/mySLBBEPool"},
	}

	actual := createAgentVMASNetworkInterface(cs, profile)

	var ipConfigurations []network.InterfaceIPConfiguration

	expected := NetworkInterfaceARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			Copy: map[string]string{
				"count": "[sub(variables('fooAgentCount'), variables('fooAgentOffset'))]",
				"name":  "loop",
			},
			DependsOn: []string{
				"[variables('vnetID')]",
			},
		},
		Interface: network.Interface{
			Type:     helpers.PointerToString("Microsoft.Network/networkInterfaces"),
			Name:     helpers.PointerToString("[concat(variables('fooAgentVMNamePrefix'), 'nic-', copyIndex(variables('fooAgentOffset')))]"),
			Location: helpers.PointerToString("[variables('location')]"),
			InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
				IPConfigurations: &ipConfigurations,
			},
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}

	// Test AgentVMAS NIC with Custom Vnet
	profile.VnetSubnetID = "fooSubnet"

	actual = createAgentVMASNetworkInterface(cs, profile)

	expected.DependsOn = []string{
		"[variables('nsgID')]",
	}

	expected.NetworkSecurityGroup = &network.SecurityGroup{
		ID: helpers.PointerToString("[variables('nsgID')]"),
	}

	diff = cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}

	// Test AgentVMAS NIC with Custom Vnet and multipleIPAddresses
	profile.IPAddressCount = 5

	actual = createAgentVMASNetworkInterface(cs, profile)

	expected.IPConfigurations = &[]network.InterfaceIPConfiguration{
		{
			Name: helpers.PointerToString("ipconfig1"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
					{
						ID: helpers.PointerToString("[concat(resourceId('Microsoft.Network/loadBalancers', variables('routerLBName')), '/backendAddressPools/backend')]"),
					},
				},
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('fooAgentVnetSubnetID')]"),
				},
				Primary: helpers.PointerToBool(true),
			},
		},
		{
			Name: helpers.PointerToString("ipconfig2"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
					{
						ID: helpers.PointerToString("[concat(resourceId('Microsoft.Network/loadBalancers', variables('routerLBName')), '/backendAddressPools/backend')]"),
					},
				},
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('fooAgentVnetSubnetID')]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("ipconfig3"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
					{
						ID: helpers.PointerToString("[concat(resourceId('Microsoft.Network/loadBalancers', variables('routerLBName')), '/backendAddressPools/backend')]"),
					},
				},
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('fooAgentVnetSubnetID')]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("ipconfig4"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
					{
						ID: helpers.PointerToString("[concat(resourceId('Microsoft.Network/loadBalancers', variables('routerLBName')), '/backendAddressPools/backend')]"),
					},
				},
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('fooAgentVnetSubnetID')]"),
				},
			},
		},
		{
			Name: helpers.PointerToString("ipconfig5"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
					{
						ID: helpers.PointerToString("[concat(resourceId('Microsoft.Network/loadBalancers', variables('routerLBName')), '/backendAddressPools/backend')]"),
					},
				},
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('fooAgentVnetSubnetID')]"),
				},
			},
		},
	}

	diff = cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}
}

func TestCreateNICWithIPv6DualStackFeature(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ClientID: "barClientID",
				Secret:   "bazSecret",
			},
			MasterProfile: &api.MasterProfile{
				Count:               1,
				DNSPrefix:           "myprefix1",
				VMSize:              "Standard_DS2_v2",
				AvailabilityProfile: api.VirtualMachineScaleSets,
				IPAddressCount:      5,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorType:    api.Kubernetes,
				OrchestratorVersion: "1.15.0-beta.2",
				KubernetesConfig: &api.KubernetesConfig{
					NetworkPlugin: "kubenet",
				},
			},
			LinuxProfile: &api.LinuxProfile{},
			AgentPoolProfiles: []*api.AgentPoolProfile{
				{
					Name:                "agentpool",
					VMSize:              "Standard_D2_v2",
					Count:               1,
					AvailabilityProfile: api.AvailabilitySet,
				},
			},
			FeatureFlags: &api.FeatureFlags{
				EnableIPv6DualStack: true,
			},
		},
	}

	nic := CreateMasterVMNetworkInterfaces(cs)
	expected := NetworkInterfaceARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			Copy: map[string]string{
				"count": "[sub(variables('masterCount'), variables('masterOffset'))]",
				"name":  "nicLoopNode",
			},
			DependsOn: []string{
				"[variables('vnetID')]",
				"[variables('masterLbName')]",
			},
		},
		Interface: network.Interface{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[concat(variables('masterVMNamePrefix'), 'nic-', copyIndex(variables('masterOffset')))]"),
			InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
				IPConfigurations: &[]network.InterfaceIPConfiguration{
					{
						Name: helpers.PointerToString("ipconfig1"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
								{
									ID: helpers.PointerToString("[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"),
								},
							},
							LoadBalancerInboundNatRules: &[]network.InboundNatRule{
								{
									ID: helpers.PointerToString("[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'),copyIndex(variables('masterOffset')))]"),
								},
							},
							PrivateIPAddress:          helpers.PointerToString("[variables('masterPrivateIpAddrs')[copyIndex(variables('masterOffset'))]]"),
							Primary:                   helpers.PointerToBool(true),
							PrivateIPAllocationMethod: network.Static,
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
					{
						Name: helpers.PointerToString("ipconfigv6"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							PrivateIPAddressVersion: "IPv6",
							Primary:                 helpers.PointerToBool(false),
							Subnet: &network.Subnet{
								ID: helpers.PointerToString("[variables('vnetSubnetID')]"),
							},
						},
					},
				},
			},
			Type: helpers.PointerToString("Microsoft.Network/networkInterfaces"),
		},
	}
	expected.EnableIPForwarding = helpers.PointerToBool(true)
	diff := cmp.Diff(nic, expected)

	if diff != "" {
		t.Errorf("unexpected diff while expecting equal structs: %s", diff)
	}
}

func TestCreateAgentVMASNICWithIPv6DualStackFeature(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ClientID: "barClientID",
				Secret:   "bazSecret",
			},
			MasterProfile: &api.MasterProfile{
				Count:               1,
				DNSPrefix:           "myprefix1",
				VMSize:              "Standard_DS2_v2",
				AvailabilityProfile: api.VirtualMachineScaleSets,
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorType:    api.Kubernetes,
				OrchestratorVersion: "1.15.0-beta.2",
				KubernetesConfig: &api.KubernetesConfig{
					NetworkPlugin: "kubenet",
				},
			},
			LinuxProfile: &api.LinuxProfile{},
			AgentPoolProfiles: []*api.AgentPoolProfile{
				{
					Name:                "agentpool",
					VMSize:              "Standard_D2_v2",
					Count:               1,
					AvailabilityProfile: api.AvailabilitySet,
				},
			},
			FeatureFlags: &api.FeatureFlags{
				EnableIPv6DualStack: true,
			},
		},
	}

	profile := &api.AgentPoolProfile{
		Name:                              "fooAgent",
		OSType:                            "Linux",
		Role:                              "Infra",
		LoadBalancerBackendAddressPoolIDs: []string{"/subscriptions/123/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/mySLB/backendAddressPools/mySLBBEPool"},
		IPAddressCount:                    1,
	}

	actual := createAgentVMASNetworkInterface(cs, profile)

	var ipConfigurations []network.InterfaceIPConfiguration

	expected := NetworkInterfaceARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			Copy: map[string]string{
				"count": "[sub(variables('fooAgentCount'), variables('fooAgentOffset'))]",
				"name":  "loop",
			},
			DependsOn: []string{
				"[variables('vnetID')]",
			},
		},
		Interface: network.Interface{
			Type:     helpers.PointerToString("Microsoft.Network/networkInterfaces"),
			Name:     helpers.PointerToString("[concat(variables('fooAgentVMNamePrefix'), 'nic-', copyIndex(variables('fooAgentOffset')))]"),
			Location: helpers.PointerToString("[variables('location')]"),
			InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
				IPConfigurations: &ipConfigurations,
			},
		},
	}
	expected.IPConfigurations = &[]network.InterfaceIPConfiguration{
		{
			Name: helpers.PointerToString("ipconfig1"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				LoadBalancerBackendAddressPools: &[]network.BackendAddressPool{
					{
						ID: helpers.PointerToString("[concat(resourceId('Microsoft.Network/loadBalancers', variables('routerLBName')), '/backendAddressPools/backend')]"),
					},
					{
						ID: helpers.PointerToString("[concat(resourceId('Microsoft.Network/loadBalancers',parameters('masterEndpointDNSNamePrefix')), '/backendAddressPools/', parameters('masterEndpointDNSNamePrefix'))]"),
					},
				},
				PrivateIPAllocationMethod: network.Dynamic,
				Subnet: &network.Subnet{
					ID: helpers.PointerToString("[variables('fooAgentVnetSubnetID')]"),
				},
				Primary: helpers.PointerToBool(true),
			},
		},
		{
			Name: helpers.PointerToString("ipconfigv6"),
			InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
				PrivateIPAddressVersion: "IPv6",
				Primary:                 helpers.PointerToBool(false),
				Subnet: &network.Subnet{
					ID: helpers.PointerToString(fmt.Sprintf("[variables('%sVnetSubnetID')]", profile.Name)),
				},
			},
		},
	}
	expected.EnableIPForwarding = helpers.PointerToBool(true)

	diff := cmp.Diff(actual, expected)
	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}
}
