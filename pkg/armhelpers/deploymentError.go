// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	resources "github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/resources/armresources"
	"github.com/sirupsen/logrus"
)

// DeploymentError contains the root deployment error along with deployment operation errors
type DeploymentError struct {
	DeploymentName    string
	ResourceGroup     string
	TopError          error
	StatusCode        string
	Response          []byte
	ProvisioningState string
	OperationsLists   []*resources.DeploymentOperation
}

func (e *DeploymentError) Error() string {
	var str string
	if e.TopError != nil {
		str = e.TopError.Error()
	}
	var ops []string
	for _, operation := range e.OperationsLists {
		if operation.Properties != nil && *operation.Properties.ProvisioningState == string(api.Failed) && operation.Properties.StatusMessage != nil {
			if b, err := json.MarshalIndent(operation.Properties.StatusMessage, "", "  "); err == nil {
				ops = append(ops, string(b))
			}
		}
	}
	return fmt.Sprintf("DeploymentName[%s] ResourceGroup[%s] TopError[%s] ProvisioningState[%s] Operations[%s]",
		e.DeploymentName, e.ResourceGroup, str, e.ProvisioningState, strings.Join(ops, " | "))
}

// DeploymentValidationError contains validation error
type DeploymentValidationError struct {
	Err error
}

// Error implements error interface
func (e *DeploymentValidationError) Error() string {
	return e.Err.Error()
}

// DeployTemplateSync deploys the template and returns ArmError
func DeployTemplateSync(az AKSEngineClient, logger *logrus.Entry, resourceGroupName, deploymentName string, template map[string]interface{}, parameters map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultARMOperationTimeout)
	defer cancel()
	deploymentExtended, err := az.DeployTemplate(ctx, resourceGroupName, deploymentName, template, parameters)
	if err == nil {
		return nil
	}

	logger.Infof("Getting detailed deployment errors for %s", deploymentName)
	deploymentErr := &DeploymentError{
		DeploymentName: deploymentName,
		ResourceGroup:  resourceGroupName,
		TopError:       err,
	}

	// try to extract error from ARM Response
	if deploymentExtended.Properties != nil && deploymentExtended.Properties.Error != nil && deploymentExtended.Properties.Error.Code != nil {
		logger.Infof("StatusCode: %s, Error: %s", *deploymentExtended.Properties.Error.Code, *deploymentExtended.Properties.Error.Message)
		deploymentErr.StatusCode = *deploymentExtended.Properties.Error.Code
		deploymentErr.Response, _ = json.Marshal(deploymentExtended.Properties.Error.Message)
	} else {
		logger.Errorf("Got error from Azure SDK without response from ARM")
		// This is the failed sdk validation before calling ARM path
		deploymentErr.StatusCode = "0"
		return deploymentErr
	}

	if deploymentExtended.Properties == nil || deploymentExtended.Properties.ProvisioningState == nil {
		logger.Warn("No resources.DeploymentExtended.Properties")
		return deploymentErr
	}
	operations, err := az.ListDeploymentOperations(ctx, resourceGroupName, deploymentName)
	if err != nil {
		logger.Errorf("unable to list deployment operations %s. error: %v", deploymentName, err)
	}

	deploymentErr.OperationsLists = append(deploymentErr.OperationsLists, operations...)
	deploymentErr.ProvisioningState = *deploymentExtended.Properties.ProvisioningState
	return deploymentErr
}
