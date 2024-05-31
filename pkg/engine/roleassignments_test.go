// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
	"github.com/google/go-cmp/cmp"
)

func TestCreateMSIRoleAssignment(t *testing.T) {
	// Test create Contributor role assignment
	actual := createMSIRoleAssignment(IdentityContributorRole)
	expected := RoleAssignmentARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionAuthorizationUser')]",
		},
		RoleAssignment: authorization.RoleAssignment{
			Type: helpers.PointerToString("Microsoft.Authorization/roleAssignments"),
			Name: helpers.PointerToString("[guid(concat(variables('userAssignedID'), 'roleAssignment', resourceGroup().id))]"),
			RoleAssignmentPropertiesWithScope: &authorization.RoleAssignmentPropertiesWithScope{
				RoleDefinitionID: helpers.PointerToString("[variables('contributorRoleDefinitionId')]"),
				PrincipalID:      helpers.PointerToString("[reference(variables('userAssignedIDReference'), variables('apiVersionManagedIdentity')).principalId]"),
				PrincipalType:    authorization.ServicePrincipal,
				Scope:            helpers.PointerToString("[resourceGroup().id]"),
			},
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}

	// Test create Reader role assignment
	actual = createMSIRoleAssignment(IdentityReaderRole)
	expected = RoleAssignmentARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionAuthorizationUser')]",
		},
		RoleAssignment: authorization.RoleAssignment{
			Type: helpers.PointerToString("Microsoft.Authorization/roleAssignments"),
			Name: helpers.PointerToString("[guid(concat(variables('userAssignedID'), 'roleAssignment', resourceGroup().id))]"),
			RoleAssignmentPropertiesWithScope: &authorization.RoleAssignmentPropertiesWithScope{
				RoleDefinitionID: helpers.PointerToString("[variables('readerRoleDefinitionId')]"),
				PrincipalID:      helpers.PointerToString("[reference(variables('userAssignedIDReference'), variables('apiVersionManagedIdentity')).principalId]"),
				PrincipalType:    authorization.ServicePrincipal,
				Scope:            helpers.PointerToString("[resourceGroup().id]"),
			},
		},
	}

	diff = cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}
}

func TestCreateKubernetesSpAppGIdentityOperatorAccessRoleAssignment(t *testing.T) {
	// using service principal
	cs := &api.ContainerService{
		Properties: &api.Properties{
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ObjectID: "xxxx",
			},
		},
	}

	actual := createKubernetesSpAppGIdentityOperatorAccessRoleAssignment(cs.Properties)
	expected := RoleAssignmentARM{
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
				PrincipalID:      helpers.PointerToString("xxxx"),
				PrincipalType:    authorization.ServicePrincipal,
				Scope:            helpers.PointerToString("[variables('appGwICIdentityId')]"),
			},
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}

	// using managed identity
	cs = &api.ContainerService{
		Properties: &api.Properties{
			OrchestratorProfile: &api.OrchestratorProfile{
				KubernetesConfig: &api.KubernetesConfig{
					UseManagedIdentity: helpers.PointerToBool(true),
				},
			},
		},
	}

	actual = createKubernetesSpAppGIdentityOperatorAccessRoleAssignment(cs.Properties)
	expected = RoleAssignmentARM{
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
				PrincipalID:      helpers.PointerToString("[reference(concat('Microsoft.ManagedIdentity/userAssignedIdentities/', variables('userAssignedID'))).principalId]"),
				PrincipalType:    authorization.ServicePrincipal,
				Scope:            helpers.PointerToString("[variables('appGwICIdentityId')]"),
			},
		},
	}

	diff = cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}
}

func TestCreateAppGwIdentityResourceGroupReadSysRoleAssignment(t *testing.T) {
	actual := createAppGwIdentityResourceGroupReadSysRoleAssignment()
	expected := RoleAssignmentARM{
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

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}
}

func TestCreateAppGwIdentityApplicationGatewayWriteSysRoleAssignment(t *testing.T) {
	actual := createAppGwIdentityApplicationGatewayWriteSysRoleAssignment()
	expected := RoleAssignmentARM{
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

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing: %s", diff)
	}
}
