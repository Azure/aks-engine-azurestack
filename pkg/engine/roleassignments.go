// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
)

type IdentityRoleDefinition string

const (
	// IdentityContributorRole means created user assigned identity will have "Contributor" role in created resource group
	IdentityContributorRole IdentityRoleDefinition = "[variables('contributorRoleDefinitionId')]"
	// IdentityReaderRole means created user assigned identity will have "Reader" role in created resource group
	IdentityReaderRole IdentityRoleDefinition = "[variables('readerRoleDefinitionId')]"
	// IdentityManagedIdentityOperatorRole means created user assigned identity or service principal will have operator access on a different managed identity
	IdentityManagedIdentityOperatorRole IdentityRoleDefinition = "[variables('managedIdentityOperatorRoleDefinitionId')]"
)

func createMSIRoleAssignment(identityRoleDefinition IdentityRoleDefinition) RoleAssignmentARM {
	return RoleAssignmentARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionAuthorizationUser')]",
		},
		RoleAssignment: authorization.RoleAssignment{
			Type: helpers.PointerToString("Microsoft.Authorization/roleAssignments"),
			Name: helpers.PointerToString("[guid(concat(variables('userAssignedID'), 'roleAssignment', resourceGroup().id))]"),
			RoleAssignmentPropertiesWithScope: &authorization.RoleAssignmentPropertiesWithScope{
				RoleDefinitionID: helpers.PointerToString(string(identityRoleDefinition)),
				PrincipalID:      helpers.PointerToString("[reference(variables('userAssignedIDReference'), variables('apiVersionManagedIdentity')).principalId]"),
				PrincipalType:    authorization.ServicePrincipal,
				Scope:            helpers.PointerToString("[resourceGroup().id]"),
			},
		},
	}
}

// createKubernetesSpAppGIdentityOperatorAccessRoleAssignment gives identity operator access on AGIC Identity to the cluster identity
func createKubernetesSpAppGIdentityOperatorAccessRoleAssignment(prop *api.Properties) RoleAssignmentARM {
	kubernetesSpObjectID := ""
	// determine objectId of the cluster identity used by the kubernetes cluster
	if prop.OrchestratorProfile != nil &&
		prop.OrchestratorProfile.KubernetesConfig != nil &&
		helpers.Bool(prop.OrchestratorProfile.KubernetesConfig.UseManagedIdentity) {
		kubernetesSpObjectID = "[reference(concat('Microsoft.ManagedIdentity/userAssignedIdentities/', variables('userAssignedID'))).principalId]"
	} else if prop.ServicePrincipalProfile.ObjectID != "" {
		kubernetesSpObjectID = prop.ServicePrincipalProfile.ObjectID
	}

	return RoleAssignmentARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionAuthorizationSystem')]",
			DependsOn: []string{
				"[concat('Microsoft.Network/applicationgateways/', variables('appGwName'))]",
				"[concat('Microsoft.ManagedIdentity/userAssignedIdentities/', variables('appGwICIdentityName'))]",
			},
		},
		RoleAssignment: authorization.RoleAssignment{
			Type: helpers.PointerToString("Microsoft.ManagedIdentity/userAssignedIdentities/providers/roleAssignments"),
			Name: helpers.PointerToString("[concat(variables('appGwICIdentityName'), '/Microsoft.Authorization/', guid(resourceGroup().id, 'aksidentityaccess'))]"),
			RoleAssignmentPropertiesWithScope: &authorization.RoleAssignmentPropertiesWithScope{
				RoleDefinitionID: helpers.PointerToString(string(IdentityManagedIdentityOperatorRole)),
				PrincipalID:      helpers.PointerToString(kubernetesSpObjectID),
				PrincipalType:    authorization.ServicePrincipal,
				Scope:            helpers.PointerToString("[variables('appGwICIdentityId')]"),
			},
		},
	}
}

// createAppGwIdentityResourceGroupReadSysRoleAssignment gives read access to Resource Group for Identity used by AGIC
func createAppGwIdentityResourceGroupReadSysRoleAssignment() RoleAssignmentARM {
	return RoleAssignmentARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionAuthorizationSystem')]",
			DependsOn: []string{
				"[concat('Microsoft.Network/applicationgateways/', variables('appGwName'))]",
				"[concat('Microsoft.ManagedIdentity/userAssignedIdentities/', variables('appGwICIdentityName'))]",
			},
		},
		RoleAssignment: authorization.RoleAssignment{
			Type: helpers.PointerToString("Microsoft.Authorization/roleAssignments"),
			Name: helpers.PointerToString("[guid(resourceGroup().id, 'identityrgaccess')]"),
			RoleAssignmentPropertiesWithScope: &authorization.RoleAssignmentPropertiesWithScope{
				RoleDefinitionID: helpers.PointerToString(string(IdentityReaderRole)),
				PrincipalID:      helpers.PointerToString("[reference(variables('appGwICIdentityId'), variables('apiVersionManagedIdentity')).principalId]"),
				Scope:            helpers.PointerToString("[resourceGroup().id]"),
			},
		},
	}
}

// createAppGwIdentityApplicationGatewayWriteSysRoleAssignment gives write access to Application Gateway for Identity used by AGIC
func createAppGwIdentityApplicationGatewayWriteSysRoleAssignment() RoleAssignmentARM {
	return RoleAssignmentARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionAuthorizationSystem')]",
			DependsOn: []string{
				"[concat('Microsoft.Network/applicationgateways/', variables('appGwName'))]",
				"[concat('Microsoft.ManagedIdentity/userAssignedIdentities/', variables('appGwICIdentityName'))]",
			},
		},
		RoleAssignment: authorization.RoleAssignment{
			Type: helpers.PointerToString("Microsoft.Network/applicationgateways/providers/roleAssignments"),
			Name: helpers.PointerToString("[concat(variables('appGwName'), '/Microsoft.Authorization/', guid(resourceGroup().id, 'identityappgwaccess'))]"),
			RoleAssignmentPropertiesWithScope: &authorization.RoleAssignmentPropertiesWithScope{
				RoleDefinitionID: helpers.PointerToString(string(IdentityContributorRole)),
				PrincipalID:      helpers.PointerToString("[reference(variables('appGwICIdentityId'), variables('apiVersionManagedIdentity')).principalId]"),
				Scope:            helpers.PointerToString("[variables('appGwId')]"),
			},
		},
	}
}
