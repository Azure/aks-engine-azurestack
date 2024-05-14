// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"

	compute "github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/pkg/errors"
)

// ListVirtualMachineImages returns the list of images available in the current environment
func (az *AzureClient) ListVirtualMachineImages(ctx context.Context, location, publisherName, offer, skus string) ([]*compute.VirtualMachineImageResource, error) {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	list, err := az.virtualMachineImagesClient.List(ctx, location, publisherName, offer, skus, &compute.VirtualMachineImagesClientListOptions{
		Top: to.Ptr(int32(10)),
	})
	if err != nil {
		return nil, errors.Wrap(err, "listing virtual machine images")
	}
	return list.VirtualMachineImageResourceArray, nil
}

// GetVirtualMachineImage returns an image or an error where there is no image
func (az *AzureClient) GetVirtualMachineImage(ctx context.Context, location, publisherName, offer, skus, version string) (compute.VirtualMachineImage, error) {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	image, err := az.virtualMachineImagesClient.Get(ctx, location, publisherName, offer, skus, version, nil)
	if err != nil {
		return compute.VirtualMachineImage{}, errors.Wrap(err, "fetching virtual machine image")
	}
	return image.VirtualMachineImage, err
}
