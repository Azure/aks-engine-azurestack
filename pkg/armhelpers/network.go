// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// DeleteNetworkInterface deletes the specified network interface.
func (az *AzureClient) DeleteNetworkInterface(ctx context.Context, resourceGroup, nicName string) error {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	poller, err := az.interfacesClient.BeginDelete(ctx, resourceGroup, nicName, nil)
	if err != nil {
		return err
	}
	if _, err = poller.PollUntilDone(ctx, nil); err != nil {
		return err
	}
	return err
}
