// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func (az *AzureClient) DeleteVirtualHardDisk(ctx context.Context, resourceGroup string, vhd *armcompute.VirtualHardDisk) error {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	parts, err := azblob.ParseURL(*vhd.URI)
	if err != nil {
		return err
	}
	accountName := strings.Split(parts.Host, ".")[0]
	keys, err := az.getStorageKeys(ctx, resourceGroup, accountName)
	if err != nil {
		return err
	}
	serviceURL := fmt.Sprintf("%s%s", parts.Scheme, parts.Host)
	client, err := az.storageBlobClientFactory(serviceURL, *keys[0].Value)
	if err != nil {
		return err
	}
	_, err = client.DeleteBlob(ctx, parts.ContainerName, parts.BlobName, nil)
	if err != nil {
		return err
	}
	return nil
}

func (az *AzureClient) getStorageKeys(ctx context.Context, resourceGroup, accountName string) ([]*armstorage.AccountKey, error) {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	storageKeysResult, err := az.storageAccountsClient.ListKeys(ctx, resourceGroup, accountName, nil)
	if err != nil {
		return nil, err
	}
	return storageKeysResult.Keys, nil
}
