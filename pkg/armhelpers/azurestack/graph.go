// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package azurestack

import (
	"context"

	"github.com/Azure/aks-engine-azurestack/pkg/armhelpers"
	"github.com/Azure/azure-sdk-for-go/services/authorization/mgmt/2015-07-01/authorization"
	"github.com/pkg/errors"
)

// DeleteRoleAssignmentByID deletes a roleAssignment via its unique identifier
func (az *AzureClient) DeleteRoleAssignmentByID(ctx context.Context, roleAssignmentID string) (authorization.RoleAssignment, error) {
	errorMessage := "error azure stack does not support deleting role assignement"
	return authorization.RoleAssignment{}, errors.New(errorMessage)
}

// ListRoleAssignmentsForPrincipal (e.g. a VM) via the scope and the unique identifier of the principal
func (az *AzureClient) ListRoleAssignmentsForPrincipal(ctx context.Context, scope string, principalID string) (armhelpers.RoleAssignmentListResultPage, error) {
	errorMessage := "error azure stack does not support listing role assignement"
	return nil, errors.New(errorMessage)
}
