// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"
	"strings"

	compute "github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/pkg/errors"
)

// ListVirtualMachines returns (the first page of) the machines in the specified resource group.
func (az *AzureClient) ListVirtualMachines(ctx context.Context, resourceGroup string) ([]*compute.VirtualMachine, error) {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	pager := az.virtualMachinesClient.NewListPager(resourceGroup, nil)
	list := []*compute.VirtualMachine{}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "listing virtual machines for resource group %s", resourceGroup)
		}
		list = append(list, page.Value...)
	}
	return list, nil
}

// GetVirtualMachine returns the specified machine in the specified resource group.
func (az *AzureClient) GetVirtualMachine(ctx context.Context, resourceGroup, name string) (compute.VirtualMachine, error) {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	vm, err := az.virtualMachinesClient.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		return compute.VirtualMachine{}, errors.Wrapf(err, "getting virtual machine %s/%s", resourceGroup, name)
	}
	return vm.VirtualMachine, nil
}

// RestartVirtualMachine restarts the specified virtual machine.
func (az *AzureClient) RestartVirtualMachine(ctx context.Context, resourceGroup, name string) error {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	poller, err := az.virtualMachinesClient.BeginRestart(ctx, resourceGroup, name, nil)
	if err != nil {
		return errors.Wrapf(err, "restarting virtual machine %s/%s", resourceGroup, name)
	}
	if _, err = poller.PollUntilDone(ctx, nil); err != nil {
		return errors.Wrapf(err, "restarting virtual machine %s/%s", resourceGroup, name)
	}
	return err
}

// DeleteVirtualMachine handles deletion of a CRP/VMAS VM (aka, not a VMSS VM).
func (az *AzureClient) DeleteVirtualMachine(ctx context.Context, resourceGroup, name string) error {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	poller, err := az.virtualMachinesClient.BeginDelete(ctx, resourceGroup, name, nil)
	if err != nil {
		return errors.Wrapf(err, "deleting virtual machine %s/%s", resourceGroup, name)
	}
	if _, err = poller.PollUntilDone(ctx, nil); err != nil {
		return errors.Wrapf(err, "deleting virtual machine %s/%s", resourceGroup, name)
	}
	return err
}

// GetVirtualMachinePowerState returns the virtual machine's PowerState status code
func (az *AzureClient) GetVirtualMachinePowerState(ctx context.Context, resourceGroup, name string) (string, error) {
	ctx = policy.WithHTTPHeader(ctx, az.acceptLanguageHeader)
	vm, err := az.GetVirtualMachine(ctx, resourceGroup, name)
	if err != nil {
		return "", errors.Wrapf(err, "fetching virtual machine %s/%s", resourceGroup, name)
	}
	for _, status := range vm.Properties.InstanceView.Statuses {
		if strings.HasPrefix(*status.Code, "PowerState") {
			return *status.Code, nil
		}
	}
	return "", nil
}
