// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
)

func createRouteTable() RouteTableARM {
	return RouteTableARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
		},
		RouteTable: network.RouteTable{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[variables('routeTableName')]"),
			Type:     helpers.PointerToString("Microsoft.Network/routeTables"),
		},
	}
}
