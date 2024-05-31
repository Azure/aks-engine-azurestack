// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/services/preview/msi/mgmt/2015-08-31-preview/msi"
)

func createUserAssignedIdentities() UserAssignedIdentitiesARM {
	return UserAssignedIdentitiesARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionManagedIdentity')]",
		},
		Identity: msi.Identity{
			Type:     helpers.PointerToString("Microsoft.ManagedIdentity/userAssignedIdentities"),
			Name:     helpers.PointerToString("[variables('userAssignedID')]"),
			Location: helpers.PointerToString("[variables('location')]"),
		},
	}
}

func createAppGwUserAssignedIdentities() UserAssignedIdentitiesARM {
	return UserAssignedIdentitiesARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionManagedIdentity')]",
		},
		Identity: msi.Identity{
			Type:     helpers.PointerToString("Microsoft.ManagedIdentity/userAssignedIdentities"),
			Name:     helpers.PointerToString("[variables('appGwICIdentityName')]"),
			Location: helpers.PointerToString("[variables('location')]"),
		},
	}
}
