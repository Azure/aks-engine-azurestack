// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
)

func TestDeleteVirtualMachine(t *testing.T) {
	mc, err := NewHTTPMockClient()
	if err != nil {
		t.Fatalf("failed to create HttpMockClient - %s", err)
	}

	mc.RegisterLogin()
	mc.RegisterVirtualMachineEndpoint()
	mc.RegisterDeleteOperation()

	err = mc.Activate()
	if err != nil {
		t.Fatalf("failed to activate HttpMockClient - %s", err)
	}
	defer mc.DeactivateAndReset()

	options := &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			InsecureAllowCredentialWithHTTP: true,
			Cloud:                           mc.GetEnvironment(),
		},
	}
	azureClient, err := NewAzureClient(subscriptionID, &fake.TokenCredential{}, options)
	if err != nil {
		t.Fatalf("can not get client %s", err)
	}

	err = azureClient.DeleteVirtualMachine(context.Background(), resourceGroup, virtualMachineName)
	if err != nil {
		t.Error(err)
	}
}
