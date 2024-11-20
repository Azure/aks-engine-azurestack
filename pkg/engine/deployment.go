// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"github.com/Azure/aks-engine-azurestack/pkg/helpers/to"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/resources/mgmt/resources"
)

func createAzureStackTelemetry(azureTelemetryPID string) DeploymentARM {
	properties := resources.DeploymentPropertiesExtended{
		Mode: "Incremental",
		Template: map[string]interface{}{
			"resources":      []interface{}{},
			"contentVersion": "1.0.0.0",
			"$schema":        "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
		},
	}

	return DeploymentARM{
		ARMResource: ARMResource{
			APIVersion: "2018-05-01",
		},
		DeploymentExtended: resources.DeploymentExtended{
			Name:       to.StringPtr(azureTelemetryPID),
			Type:       to.StringPtr("Microsoft.Resources/deployments"),
			Properties: &properties,
		},
	}
}
