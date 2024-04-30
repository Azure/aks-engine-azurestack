// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"
	"github.com/pkg/errors"
)

// ListVirtualMachines returns (the first page of) the machines in the specified resource group.
func (az *AzureClient) ListVirtualMachines(ctx context.Context, resourceGroup string) (VirtualMachineListResultPage, error) {
	page, err := az.virtualMachinesClient.List(ctx, resourceGroup)
	return &page, err
}

// GetVirtualMachine returns the specified machine in the specified resource group.
func (az *AzureClient) GetVirtualMachine(ctx context.Context, resourceGroup, name string) (compute.VirtualMachine, error) {
	return az.virtualMachinesClient.Get(ctx, resourceGroup, name, "")
}

// RestartVirtualMachine restarts the specified virtual machine.
func (az *AzureClient) RestartVirtualMachine(ctx context.Context, resourceGroup, name string) error {
	future, err := az.virtualMachinesClient.Restart(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, az.virtualMachinesClient.Client); err != nil {
		return err
	}

	_, err = future.Result(az.virtualMachinesClient)
	return err
}

// DeleteVirtualMachine handles deletion of a CRP/VMAS VM (aka, not a VMSS VM).
func (az *AzureClient) DeleteVirtualMachine(ctx context.Context, resourceGroup, name string) error {
	future, err := az.virtualMachinesClient.Delete(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, az.virtualMachinesClient.Client); err != nil {
		return err
	}

	_, err = future.Result(az.virtualMachinesClient)
	return err
}

// GetAvailabilitySet retrieves the specified VM availability set.
func (az *AzureClient) GetAvailabilitySet(ctx context.Context, resourceGroup, availabilitySetName string) (compute.AvailabilitySet, error) {
	return az.availabilitySetsClient.Get(ctx, resourceGroup, availabilitySetName)
}

// GetAvailabilitySetFaultDomainCount returns the first existing fault domain count it finds from the IDs provided.
func (az *AzureClient) GetAvailabilitySetFaultDomainCount(ctx context.Context, resourceGroup string, vmasIDs []string) (int, error) {
	var count int
	if len(vmasIDs) > 0 {
		id := vmasIDs[0]
		// extract the last element of the id for VMAS name
		ss := strings.Split(id, "/")
		name := ss[len(ss)-1]
		vmas, err := az.GetAvailabilitySet(ctx, resourceGroup, name)
		if err != nil {
			return 0, err
		}
		// Assume that all VMASes in the cluster share a value for platformFaultDomainCount
		count = int(*vmas.AvailabilitySetProperties.PlatformFaultDomainCount)
	}
	return count, nil
}

// GetVirtualMachinePowerState returns the virtual machine's PowerState status code
func (az *AzureClient) GetVirtualMachinePowerState(ctx context.Context, resourceGroup, name string) (string, error) {
	vm, err := az.virtualMachinesClient.Get(ctx, resourceGroup, name, compute.InstanceView)
	if err != nil {
		return "", errors.Wrapf(err, "fetching virtual machine resource")
	}
	for _, status := range *vm.VirtualMachineProperties.InstanceView.Statuses {
		if strings.HasPrefix(*status.Code, "PowerState") {
			return *status.Code, nil
		}
	}
	return "", nil
}
