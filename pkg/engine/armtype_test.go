// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute"
	. "github.com/onsi/gomega"
)

func TestMarshalJSON(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ClientID: "barClientID",
				Secret:   "bazSecret",
			},
			MasterProfile: &api.MasterProfile{
				Count:                     3,
				DNSPrefix:                 "myprefix1",
				VMSize:                    "Standard_DS2_v2",
				AvailabilityProfile:       api.VirtualMachineScaleSets,
				PlatformUpdateDomainCount: helpers.PointerToInt(3),
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorType:    api.Kubernetes,
				OrchestratorVersion: "1.10.2",
				KubernetesConfig: &api.KubernetesConfig{
					NetworkPlugin: "azure",
				},
			},
			FeatureFlags: &api.FeatureFlags{
				BlockOutboundInternet: false,
			},
		},
	}
	armObject := CreateCustomScriptExtension(cs)

	jsonObj, err := json.MarshalIndent(armObject, "", "   ")
	if err != nil {
		t.Error(err)
	}
	// TODO: why print this? Let's validate it.
	fmt.Println(string(jsonObj))
}

func TestMarshalJSONAvailabilitySetARM(t *testing.T) {
	g := NewGomegaWithT(t)

	type VMASTestDatum struct {
		avSet AvailabilitySetARM
		json  string
	}

	vmasTestData := []VMASTestDatum{
		{
			avSet: AvailabilitySetARM{
				ARMResource: ARMResource{
					APIVersion: "[variables('apiVersionCompute')]",
				},
				AvailabilitySet: compute.AvailabilitySet{
					Name:     helpers.PointerToString("[variables('masterAvailabilitySet')]"),
					Location: helpers.PointerToString("[variables('location')]"),
					Type:     helpers.PointerToString("Microsoft.Compute/availabilitySets"),
					Sku: &compute.Sku{
						Name: helpers.PointerToString("Aligned"),
					},
					AvailabilitySetProperties: &compute.AvailabilitySetProperties{
						PlatformFaultDomainCount:  helpers.PointerToInt32(3),
						PlatformUpdateDomainCount: helpers.PointerToInt32(3),
						ProximityPlacementGroup: &compute.SubResource{
							ID: helpers.PointerToString("ProximityPlacementGroupResourceID"),
						},
					},
				},
			},
			json: `{
			"apiVersion": "[variables('apiVersionCompute')]",
			"properties": {
				"platformFaultDomainCount": 3,
				"platformUpdateDomainCount": 3,
				"proximityPlacementGroup": {
					"id": "ProximityPlacementGroupResourceID"
				}
			},
			"sku": {
				"name": "Aligned"
			},
			"name": "[variables('masterAvailabilitySet')]",
			"type": "Microsoft.Compute/availabilitySets",
			"location": "[variables('location')]",
			"tags": null
		}`},
		{
			avSet: AvailabilitySetARM{
				ARMResource: ARMResource{
					APIVersion: "[variables('apiVersionCompute')]",
				},
				AvailabilitySet: compute.AvailabilitySet{
					Name:     helpers.PointerToString("[variables('masterAvailabilitySet')]"),
					Location: helpers.PointerToString("[variables('location')]"),
					Type:     helpers.PointerToString("Microsoft.Compute/availabilitySets"),
					Sku: &compute.Sku{
						Name: helpers.PointerToString("Aligned"),
					},
					AvailabilitySetProperties: &compute.AvailabilitySetProperties{
						PlatformUpdateDomainCount: helpers.PointerToInt32(3),
						ProximityPlacementGroup: &compute.SubResource{
							ID: helpers.PointerToString("ProximityPlacementGroupResourceID"),
						},
					},
				},
			},
			json: `{
			"apiVersion": "[variables('apiVersionCompute')]",
			"properties": {
				"platformFaultDomainCount": "[if(contains(split('canadacentral,centralus,eastus,eastus2,northcentralus,northeurope,southcentralus,westeurope,westus',','),variables('location')),3,if(equals('centraluseuap',variables('location')),1,2))]",
				"platformUpdateDomainCount": 3,
				"proximityPlacementGroup": {  
					"id": "ProximityPlacementGroupResourceID"
   				}
			},
			"sku": {
				"name": "Aligned"
			},
			"name": "[variables('masterAvailabilitySet')]",
			"type": "Microsoft.Compute/availabilitySets",
			"location": "[variables('location')]",
			"tags": null
		}`},
		{
			avSet: AvailabilitySetARM{
				ARMResource: ARMResource{
					APIVersion: "[variables('apiVersionCompute')]",
				},
				AvailabilitySet: compute.AvailabilitySet{
					Name:     helpers.PointerToString("[variables('masterAvailabilitySet')]"),
					Location: helpers.PointerToString("[variables('location')]"),
					Type:     helpers.PointerToString("Microsoft.Compute/availabilitySets"),
					Sku: &compute.Sku{
						Name: helpers.PointerToString("Aligned"),
					},
					AvailabilitySetProperties: &compute.AvailabilitySetProperties{},
				},
			},
			json: `{
			"apiVersion": "[variables('apiVersionCompute')]",
			"properties": {},
			"sku": {
				"name": "Aligned"
			},
			"name": "[variables('masterAvailabilitySet')]",
			"type": "Microsoft.Compute/availabilitySets",
			"location": "[variables('location')]",
			"tags": null
		}`},
	}

	for _, vmasTestDatum := range vmasTestData {
		output, err := json.Marshal(vmasTestDatum.avSet)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(string(output)).To(MatchJSON(vmasTestDatum.json))
	}
}
