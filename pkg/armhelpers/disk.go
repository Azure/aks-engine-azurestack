// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/pkg/errors"
)

// DeleteManagedDisk deletes a managed disk.
func (az *AzureClient) DeleteManagedDisk(ctx context.Context, resourceGroupName string, diskName string) error {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	poller, err := az.disksClient.BeginDelete(ctx, resourceGroupName, diskName, nil)
	if err != nil {
		return errors.Wrapf(err, "deleting managed disk %s/%s", resourceGroupName, diskName)
	}
	if _, err = poller.PollUntilDone(ctx, nil); err != nil {
		return errors.Wrapf(err, "deleting managed disk %s/%s", resourceGroupName, diskName)
	}
	return nil
}

// ListManagedDisksByResourceGroup lists managed disks in a resource group.
func (az *AzureClient) ListManagedDisksByResourceGroup(ctx context.Context, resourceGroupName string) ([]*armcompute.Disk, error) {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	pager := az.disksClient.NewListByResourceGroupPager(resourceGroupName, nil)
	list := []*armcompute.Disk{}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "listing managed disks for resource group %s", resourceGroupName)
		}
		list = append(list, page.Value...)
	}
	return list, nil
}
