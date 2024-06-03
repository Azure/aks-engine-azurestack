// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"github.com/Azure/aks-engine-azurestack/pkg/helpers/to"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
)

func createRouteTable() RouteTableARM {
	return RouteTableARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
		},
		RouteTable: network.RouteTable{
			Location: to.StringPtr("[variables('location')]"),
			Name:     to.StringPtr("[variables('routeTableName')]"),
			Type:     to.StringPtr("Microsoft.Network/routeTables"),
		},
	}
}
