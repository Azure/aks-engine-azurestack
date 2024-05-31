// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"fmt"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute"
)

func CreateAvailabilitySet(cs *api.ContainerService, isManagedDisks bool) AvailabilitySetARM {

	armResource := ARMResource{
		APIVersion: "[variables('apiVersionCompute')]",
	}

	avSet := compute.AvailabilitySet{
		Name:     helpers.PointerToString("[variables('masterAvailabilitySet')]"),
		Location: helpers.PointerToString("[variables('location')]"),
		Type:     helpers.PointerToString("Microsoft.Compute/availabilitySets"),
	}

	if !cs.Properties.MasterProfile.HasAvailabilityZones() {
		if isManagedDisks {
			avSet.AvailabilitySetProperties = &compute.AvailabilitySetProperties{}
			if cs.Properties.MasterProfile.PlatformFaultDomainCount != nil {
				p := int32(*cs.Properties.MasterProfile.PlatformFaultDomainCount)
				avSet.PlatformFaultDomainCount = helpers.PointerToInt32(p)
			}
			if cs.Properties.MasterProfile.PlatformUpdateDomainCount != nil {
				p := int32(*cs.Properties.MasterProfile.PlatformUpdateDomainCount)
				avSet.PlatformUpdateDomainCount = helpers.PointerToInt32(p)
			}
			if cs.Properties.MasterProfile.ProximityPlacementGroupID != "" {
				avSet.ProximityPlacementGroup = &compute.SubResource{
					ID: helpers.PointerToString(cs.Properties.MasterProfile.ProximityPlacementGroupID),
				}
			}
			avSet.Sku = &compute.Sku{
				Name: helpers.PointerToString("Aligned"),
			}
		} else if cs.Properties.MasterProfile.IsStorageAccount() {
			avSet.AvailabilitySetProperties = &compute.AvailabilitySetProperties{}
		}
	}

	return AvailabilitySetARM{
		ARMResource:     armResource,
		AvailabilitySet: avSet,
	}
}

func createAgentAvailabilitySets(profile *api.AgentPoolProfile) AvailabilitySetARM {

	armResource := ARMResource{
		APIVersion: "[variables('apiVersionCompute')]",
	}

	avSet := compute.AvailabilitySet{
		Name:                      helpers.PointerToString(fmt.Sprintf("[variables('%sAvailabilitySet')]", profile.Name)),
		Location:                  helpers.PointerToString("[variables('location')]"),
		Type:                      helpers.PointerToString("Microsoft.Compute/availabilitySets"),
		AvailabilitySetProperties: &compute.AvailabilitySetProperties{},
	}

	if profile.IsManagedDisks() {
		if profile.PlatformFaultDomainCount != nil {
			p := int32(*profile.PlatformFaultDomainCount)
			avSet.PlatformFaultDomainCount = helpers.PointerToInt32(p)
		}
		if profile.PlatformUpdateDomainCount != nil {
			p := int32(*profile.PlatformUpdateDomainCount)
			avSet.PlatformUpdateDomainCount = helpers.PointerToInt32(p)
		}
		if profile.ProximityPlacementGroupID != "" {
			avSet.ProximityPlacementGroup = &compute.SubResource{
				ID: helpers.PointerToString(profile.ProximityPlacementGroupID),
			}
		}

		avSet.Sku = &compute.Sku{
			Name: helpers.PointerToString("Aligned"),
		}
	}

	return AvailabilitySetARM{
		ARMResource:     armResource,
		AvailabilitySet: avSet,
	}
}
