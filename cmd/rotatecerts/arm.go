// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package rotatecerts

import (
	"context"
	"strings"
	"time"

	"github.com/Azure/aks-engine-azurestack/pkg/armhelpers"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

// ARMClientWrapper is an ARM client with simple retry logic
type ARMClientWrapper struct {
	client    armhelpers.AKSEngineClient
	timeout   time.Duration
	backoff   wait.Backoff
	retryFunc func(err error) bool
}

// NewARMClientWrapper returns an ARM client with simple retry logic
func NewARMClientWrapper(client armhelpers.AKSEngineClient, interval, timeout time.Duration) *ARMClientWrapper {
	return &ARMClientWrapper{
		client:  client,
		timeout: timeout,
		backoff: wait.Backoff{
			Steps:    int(int64(timeout/time.Millisecond) / int64(interval/time.Millisecond)),
			Duration: interval,
		},
		retryFunc: func(err error) bool { return err != nil },
	}
}

// GetVirtualMachinePowerState restarts the specified virtual machine.
func (arm *ARMClientWrapper) GetVirtualMachinePowerState(resourceGroup, vmName string) (string, error) {
	var err error
	status := ""
	err = retry.OnError(arm.backoff, arm.retryFunc, func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		status, err = arm.client.GetVirtualMachinePowerState(ctx, resourceGroup, vmName)
		if err != nil {
			return errors.Wrap(err, "fetching virtual machine resource")
		}
		return nil
	})
	return status, err
}

// RestartVirtualMachine returns the virtual machine's Power state
func (arm *ARMClientWrapper) RestartVirtualMachine(resourceGroup, vmName string) error {
	var err error
	err = retry.OnError(arm.backoff, arm.retryFunc, func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err = arm.client.RestartVirtualMachine(ctx, resourceGroup, vmName); err != nil {
			return errors.Wrap(err, "restarting virtual machine")
		}
		return nil
	})
	return err
}

func isVirtualMachineRunning(status string) bool {
	return strings.EqualFold(status, "PowerState/running")
}
