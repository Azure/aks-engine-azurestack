// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/authorization/mgmt/2015-07-01/authorization"
)

// DeleteRoleAssignmentByID deletes a roleAssignment via its unique identifier
func (az *AzureClient) DeleteRoleAssignmentByID(ctx context.Context, roleAssignmentID string) (authorization.RoleAssignment, error) {
	return az.authorizationClient.DeleteByID(ctx, roleAssignmentID)
}

// ListRoleAssignmentsForPrincipal (e.g. a VM) via the scope and the unique identifier of the principal
func (az *AzureClient) ListRoleAssignmentsForPrincipal(ctx context.Context, scope string, principalID string) (RoleAssignmentListResultPage, error) {
	page, err := az.authorizationClient.ListForScope(ctx, scope, fmt.Sprintf("principalId eq '%s'", principalID))
	return &page, err
}
