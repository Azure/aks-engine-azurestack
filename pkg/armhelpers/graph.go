// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"
	"fmt"

	authorization "github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/authorization/armauthorization"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/pkg/errors"
)

// DeleteRoleAssignmentByID deletes a roleAssignment via its unique identifier
func (az *AzureClient) DeleteRoleAssignmentByID(ctx context.Context, roleAssignmentID string) (authorization.RoleAssignment, error) {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	response, err := az.authorizationClient.DeleteByID(ctx, roleAssignmentID, nil)
	if err != nil {
		return authorization.RoleAssignment{}, errors.Wrapf(err, "deleting role assignment %s", roleAssignmentID)
	}
	return response.RoleAssignment, err
}

// ListRoleAssignmentsForPrincipal (e.g. a VM) via the scope and the unique identifier of the principal
func (az *AzureClient) ListRoleAssignmentsForPrincipal(ctx context.Context, scope string, principalID string) ([]*authorization.RoleAssignment, error) {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	pager := az.authorizationClient.NewListForScopePager(scope, &authorization.RoleAssignmentsClientListForScopeOptions{
		Filter: to.Ptr(fmt.Sprintf("principalId eq '%s'", principalID)),
	})
	list := []*authorization.RoleAssignment{}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "listing roles assignments for principal %s", principalID)
		}
		list = append(list, page.Value...)
	}
	return list, nil
}
