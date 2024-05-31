// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute"
	"github.com/google/go-cmp/cmp"
)

func TestCreateAvailabilitySet(t *testing.T) {

	//Test AvSet without ManagedDisk
	cs := &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{
				AvailabilityZones: []string{},
			},
		},
	}

	avSet := CreateAvailabilitySet(cs, false)

	expectedAvSet := AvailabilitySetARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionCompute')]",
		},
		AvailabilitySet: compute.AvailabilitySet{
			Name:     helpers.PointerToString("[variables('masterAvailabilitySet')]"),
			Location: helpers.PointerToString("[variables('location')]"),
			Type:     helpers.PointerToString("Microsoft.Compute/availabilitySets"),
		},
	}

	diff := cmp.Diff(avSet, expectedAvSet)

	if diff != "" {
		t.Errorf("unexpected error while comparing availability sets: %s", diff)
	}

	//Test AvSet with ManagedDisk

	cs = &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{
				PlatformUpdateDomainCount: helpers.PointerToInt(3),
			},
		},
	}

	avSet = CreateAvailabilitySet(cs, true)

	expectedAvSet = AvailabilitySetARM{
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
			},
		},
	}

	diff = cmp.Diff(avSet, expectedAvSet)

	if diff != "" {
		t.Errorf("unexpected error while comparing availability sets: %s", diff)
	}

	//Test AvSet with StorageAccount
	cs = &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{
				StorageProfile: api.StorageAccount,
			},
		},
	}

	avSet = CreateAvailabilitySet(cs, false)

	expectedAvSet = AvailabilitySetARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionCompute')]",
		},
		AvailabilitySet: compute.AvailabilitySet{
			Name:                      helpers.PointerToString("[variables('masterAvailabilitySet')]"),
			Location:                  helpers.PointerToString("[variables('location')]"),
			Type:                      helpers.PointerToString("Microsoft.Compute/availabilitySets"),
			AvailabilitySetProperties: &compute.AvailabilitySetProperties{},
		},
	}

	diff = cmp.Diff(avSet, expectedAvSet)

	if diff != "" {
		t.Errorf("unexpected error while comparing availability sets: %s", diff)
	}

	// Test availability set with platform fault domain+update count  set
	count := 3
	cs = &api.ContainerService{
		Properties: &api.Properties{
			MasterProfile: &api.MasterProfile{
				PlatformFaultDomainCount:  &count,
				PlatformUpdateDomainCount: &count,
			},
		},
	}

	avSet = CreateAvailabilitySet(cs, true)

	expectedAvSet = AvailabilitySetARM{
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
				PlatformFaultDomainCount:  helpers.PointerToInt32(int32(count)),
				PlatformUpdateDomainCount: helpers.PointerToInt32(3),
			},
		},
	}

	diff = cmp.Diff(avSet, expectedAvSet)

	if diff != "" {
		t.Errorf("unexpected error while comparing availability sets: %s", diff)
	}
}

func TestCreateAgentAvailabilitySets(t *testing.T) {
	//Test AvSet without ManagedDisk
	profile := &api.AgentPoolProfile{
		Name:           "foobar",
		StorageProfile: api.StorageAccount,
	}

	avSet := createAgentAvailabilitySets(profile)

	expectedAvSet := AvailabilitySetARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionCompute')]",
		},
		AvailabilitySet: compute.AvailabilitySet{
			Name:                      helpers.PointerToString("[variables('foobarAvailabilitySet')]"),
			Location:                  helpers.PointerToString("[variables('location')]"),
			Type:                      helpers.PointerToString("Microsoft.Compute/availabilitySets"),
			AvailabilitySetProperties: &compute.AvailabilitySetProperties{},
		},
	}

	diff := cmp.Diff(avSet, expectedAvSet)

	if diff != "" {
		t.Errorf("unexpected error while comparing availability sets: %s", diff)
	}

	//Test AvSet wit ManagedDisk
	profile = &api.AgentPoolProfile{
		Name:                      "foobar",
		StorageProfile:            api.ManagedDisks,
		PlatformUpdateDomainCount: helpers.PointerToInt(3),
	}

	avSet = createAgentAvailabilitySets(profile)

	expectedAvSet = AvailabilitySetARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionCompute')]",
		},
		AvailabilitySet: compute.AvailabilitySet{
			Name:     helpers.PointerToString("[variables('foobarAvailabilitySet')]"),
			Location: helpers.PointerToString("[variables('location')]"),
			Type:     helpers.PointerToString("Microsoft.Compute/availabilitySets"),
			AvailabilitySetProperties: &compute.AvailabilitySetProperties{
				PlatformUpdateDomainCount: helpers.PointerToInt32(3),
			},
			Sku: &compute.Sku{
				Name: helpers.PointerToString("Aligned"),
			},
		},
	}

	diff = cmp.Diff(avSet, expectedAvSet)

	if diff != "" {
		t.Errorf("unexpected error while comparing availability sets: %s", diff)
	}

	// Test availability set with platform fault+update domain count set
	count := 3
	profile = &api.AgentPoolProfile{
		Name:                      "foobar",
		StorageProfile:            api.ManagedDisks,
		PlatformFaultDomainCount:  &count,
		PlatformUpdateDomainCount: &count,
	}

	avSet = createAgentAvailabilitySets(profile)

	expectedAvSet = AvailabilitySetARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionCompute')]",
		},
		AvailabilitySet: compute.AvailabilitySet{
			Name:     helpers.PointerToString("[variables('foobarAvailabilitySet')]"),
			Location: helpers.PointerToString("[variables('location')]"),
			Type:     helpers.PointerToString("Microsoft.Compute/availabilitySets"),
			AvailabilitySetProperties: &compute.AvailabilitySetProperties{
				PlatformFaultDomainCount:  helpers.PointerToInt32(int32(count)),
				PlatformUpdateDomainCount: helpers.PointerToInt32(3),
			},
			Sku: &compute.Sku{
				Name: helpers.PointerToString("Aligned"),
			},
		},
	}

	diff = cmp.Diff(avSet, expectedAvSet)

	if diff != "" {
		t.Errorf("unexpected error while comparing availability sets: %s", diff)
	}
}
