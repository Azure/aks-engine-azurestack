// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/pkg/errors"
)

// ListDeploymentOperations gets all deployments operations for a deployment.
func (az *AzureClient) ListDeploymentOperations(ctx context.Context, resourceGroupName string, deploymentName string) ([]*armresources.DeploymentOperation, error) {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	pager := az.deploymentOperationsClient.NewListPager(resourceGroupName, deploymentName, nil)
	list := []*armresources.DeploymentOperation{}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "listing deployment operations for %s/%s", resourceGroupName, deploymentName)
		}
		list = append(list, page.Value...)
	}
	return list, nil
}
