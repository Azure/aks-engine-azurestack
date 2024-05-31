// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/services/preview/msi/mgmt/2015-08-31-preview/msi"
)

func TestCreateUserAssignedIdentities(t *testing.T) {
	expectedAssignedIdentity := UserAssignedIdentitiesARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionManagedIdentity')]",
		},
		Identity: msi.Identity{
			Type:     helpers.PointerToString("Microsoft.ManagedIdentity/userAssignedIdentities"),
			Name:     helpers.PointerToString("[variables('userAssignedID')]"),
			Location: helpers.PointerToString("[variables('location')]"),
		},
	}

	actual := createUserAssignedIdentities()

	diff := cmp.Diff(expectedAssignedIdentity, actual)

	if diff != "" {
		t.Errorf("unexpected diff while comparing structs: %s", diff)
	}
}

func TestCreateAppGwUserAssignedIdentities(t *testing.T) {
	expectedAssignedIdentity := UserAssignedIdentitiesARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionManagedIdentity')]",
		},
		Identity: msi.Identity{
			Type:     helpers.PointerToString("Microsoft.ManagedIdentity/userAssignedIdentities"),
			Name:     helpers.PointerToString("[variables('appGwICIdentityName')]"),
			Location: helpers.PointerToString("[variables('location')]"),
		},
	}

	actual := createAppGwUserAssignedIdentities()

	diff := cmp.Diff(expectedAssignedIdentity, actual)

	if diff != "" {
		t.Errorf("unexpected diff while comparing structs: %s", diff)
	}
}
