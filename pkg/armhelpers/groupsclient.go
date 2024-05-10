// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"

	resources "github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// EnsureResourceGroup ensures the named resource group exists in the given location.
func (az *AzureClient) EnsureResourceGroup(ctx context.Context, name, location string, managedBy *string) (resources.ResourceGroup, error) {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	var tags map[string]*string
	group, err := az.groupsClient.Get(ctx, name, nil)
	if err == nil {
		tags = group.Tags
	}
	response, err := az.groupsClient.CreateOrUpdate(ctx, name, resources.ResourceGroup{
		Name:      &name,
		Location:  &location,
		ManagedBy: managedBy,
		Tags:      tags,
	}, nil)
	if err != nil {
		return resources.ResourceGroup{}, err
	}
	return response.ResourceGroup, nil
}

// DeleteResourceGroup delete the named resource group
func (az *AzureClient) DeleteResourceGroup(ctx context.Context, name string) error {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	poller, err := az.groupsClient.BeginDelete(ctx, name, nil)
	if err != nil {
		return err
	}
	if _, err = poller.PollUntilDone(ctx, nil); err != nil {
		return err
	}
	return err
}
