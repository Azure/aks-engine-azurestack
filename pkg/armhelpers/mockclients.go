// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
	"github.com/Azure/aks-engine-azurestack/pkg/kubernetes"
	"github.com/Azure/go-autorest/autorest/to"

	authorization "github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/authorization/armauthorization"
	compute "github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/compute/armcompute"
	resources "github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/resources/armresources"
	azStorage "github.com/Azure/azure-sdk-for-go/storage"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	//DefaultFakeVMName is the default name assigned to VMs part of FakeListVirtualMachineResult
	DefaultFakeVMName = "k8s-agentpool1-12345678-0"
)

var defaultK8sVersionForFakeVMs string

func init() {
	initialVersion := common.RationalizeReleaseAndVersion(common.Kubernetes, "", "", false, false, false)
	defaultK8sVersionForFakeVMs = fmt.Sprintf("Kubernetes:%s", initialVersion)
}

// MockAKSEngineClient is an implementation of AKSEngineClient where all requests error out
type MockAKSEngineClient struct {
	FailDeployTemplate                     bool
	FailDeployTemplateQuota                bool
	FailDeployTemplateConflict             bool
	FailDeployTemplateWithProperties       bool
	FailEnsureResourceGroup                bool
	FailListVirtualMachines                bool
	FailListVirtualMachinesTags            bool
	FailGetVirtualMachine                  bool
	FailRestartVirtualMachine              bool
	FailDeleteVirtualMachine               bool
	FailDeleteNetworkInterface             bool
	FailGetKubernetesClient                bool
	FailListProviders                      bool
	ShouldSupportVMIdentity                bool
	FailDeleteRoleAssignment               bool
	FailEnsureDefaultLogAnalyticsWorkspace bool
	FailAddContainerInsightsSolution       bool
	FailGetLogAnalyticsWorkspaceInfo       bool
	MockKubernetesClient                   *MockKubernetesClient
	FakeListVirtualMachineResult           func() []*compute.VirtualMachine
}

// MockStorageClient mock implementation of StorageClient
type MockStorageClient struct {
	FailCreateContainer bool
	FailSaveBlockBlob   bool
}

// MockKubernetesClient mock implementation of KubernetesClient
type MockKubernetesClient struct {
	FailListPods              bool
	FailListNodes             bool
	FailListServiceAccounts   bool
	FailListPodSecurityPolicy bool
	FailGetNode               bool
	UpdateNodeFunc            func(*v1.Node) (*v1.Node, error)
	GetNodeFunc               func(name string) (*v1.Node, error)
	FailUpdateNode            bool
	FailDeleteNode            bool
	FailDeleteServiceAccount  bool
	FailSupportEviction       bool
	FailDeletePod             bool
	FailDeleteClusterRole     bool
	FailDeleteDaemonSet       bool
	FailDeleteDeployment      bool
	FailEvictPod              bool
	FailWaitForDelete         bool
	ShouldSupportEviction     bool
	PodsList                  *v1.PodList
	ServiceAccountList        *v1.ServiceAccountList
	PodSecurityPolicyList     *policyv1beta1.PodSecurityPolicyList
	FailGetDeploymentCount    int
	FailUpdateDeploymentCount int
}

// ListPods returns Pods running on the passed in node
func (mkc *MockKubernetesClient) ListPods(node *v1.Node) (*v1.PodList, error) {
	if mkc.FailListPods {
		return nil, errors.New("ListPods failed")
	}
	if mkc.PodsList != nil {
		return mkc.PodsList, nil
	}
	return &v1.PodList{}, nil
}

// ListAllPods returns all Pods running
func (mkc *MockKubernetesClient) ListAllPods() (*v1.PodList, error) {
	if mkc.FailListPods {
		return nil, errors.New("ListAllPods failed")
	}
	if mkc.PodsList != nil {
		return mkc.PodsList, nil
	}
	return &v1.PodList{}, nil
}

// ListNodes returns a list of Nodes registered in the api server
func (mkc *MockKubernetesClient) ListNodes() (*v1.NodeList, error) {
	if mkc.FailListNodes {
		return nil, errors.New("ListNodes failed")
	}
	node := &v1.Node{}
	node.Name = fmt.Sprintf("%s-1234", common.LegacyControlPlaneVMPrefix)
	node.Status.Conditions = append(node.Status.Conditions, v1.NodeCondition{Type: v1.NodeReady, Status: v1.ConditionTrue})
	node.Status.NodeInfo.KubeletVersion = "1.9.10"
	node2 := &v1.Node{}
	node2.Name = "k8s-agentpool3-1234"
	node2.Status.Conditions = append(node2.Status.Conditions, v1.NodeCondition{Type: v1.NodeMemoryPressure, Status: v1.ConditionTrue})
	node2.Status.NodeInfo.KubeletVersion = "1.9.9"
	nodeList := &v1.NodeList{}
	nodeList.Items = append(nodeList.Items, *node)
	nodeList.Items = append(nodeList.Items, *node2)
	return nodeList, nil
}

// ListNodesByOptions returns a list of Nodes registered in the api server
func (mkc *MockKubernetesClient) ListNodesByOptions(opts metav1.ListOptions) (*v1.NodeList, error) {
	return &v1.NodeList{}, nil
}

// ListServiceAccounts returns a list of Service Accounts in the provided namespace
func (mkc *MockKubernetesClient) ListServiceAccounts(namespace string) (*v1.ServiceAccountList, error) {
	if mkc.FailListServiceAccounts {
		return nil, errors.New("ListServiceAccounts failed")
	}
	if mkc.ServiceAccountList != nil {
		return mkc.ServiceAccountList, nil
	}
	sa := &v1.ServiceAccount{}
	sa.Namespace = namespace
	sa.Name = "service-account-1"
	sa2 := &v1.ServiceAccount{}
	sa2.Namespace = namespace
	sa.Name = "service-account-2"
	saList := &v1.ServiceAccountList{}
	saList.Items = append(saList.Items, *sa)
	saList.Items = append(saList.Items, *sa2)
	return saList, nil
}

// ListPodSecurityPolices returns the list of Pod Security Policies
func (mkc *MockKubernetesClient) ListPodSecurityPolices(opts metav1.ListOptions) (*policyv1beta1.PodSecurityPolicyList, error) {
	if mkc.FailListPodSecurityPolicy {
		return nil, errors.New("ListPodSecurityPolices failed")
	}
	if mkc.PodSecurityPolicyList != nil {
		return mkc.PodSecurityPolicyList, nil
	}
	psp1 := &policyv1beta1.PodSecurityPolicy{}
	psp1.Name = "privileged"
	psp2 := &policyv1beta1.PodSecurityPolicy{}
	psp2.Name = "restricted"
	pspList := &policyv1beta1.PodSecurityPolicyList{}
	pspList.Items = append(pspList.Items, *psp1)
	pspList.Items = append(pspList.Items, *psp2)
	return pspList, nil
}

// GetNode returns details about node with passed in name
func (mkc *MockKubernetesClient) GetNode(name string) (*v1.Node, error) {
	if mkc.GetNodeFunc != nil {
		return mkc.GetNodeFunc(name)
	}
	if mkc.FailGetNode {
		return nil, errors.New("GetNode failed")
	}
	node := &v1.Node{}
	node.Status.Conditions = append(node.Status.Conditions, v1.NodeCondition{Type: v1.NodeReady, Status: v1.ConditionTrue})
	node.Status.NodeInfo.KubeletVersion = common.RationalizeReleaseAndVersion(common.Kubernetes, "", "", false, false, false)
	return node, nil
}

// UpdateNode updates the node in the api server with the passed in info
func (mkc *MockKubernetesClient) UpdateNode(node *v1.Node) (*v1.Node, error) {
	if mkc.UpdateNodeFunc != nil {
		return mkc.UpdateNodeFunc(node)
	}
	if mkc.FailUpdateNode {
		return nil, errors.New("UpdateNode failed")
	}
	return node, nil
}

// DeleteNode deregisters node in the api server
func (mkc *MockKubernetesClient) DeleteNode(name string) error {
	if mkc.FailDeleteNode {
		return errors.New("DeleteNode failed")
	}
	return nil
}

// DeleteServiceAccount deletes the provided service account
func (mkc *MockKubernetesClient) DeleteServiceAccount(sa *v1.ServiceAccount) error {
	if mkc.FailDeleteServiceAccount {
		return errors.New("DeleteServiceAccount failed")
	}
	return nil
}

// SupportEviction queries the api server to discover if it supports eviction, and returns supported type if it is supported
func (mkc *MockKubernetesClient) SupportEviction() (string, error) {
	if mkc.FailSupportEviction {
		return "", errors.New("SupportEviction failed")
	}
	if mkc.ShouldSupportEviction {
		return "version", nil
	}
	return "", nil
}

// DeleteDeployment deletes the passed in daemonset
func (mkc *MockKubernetesClient) DeleteClusterRole(role *rbacv1.ClusterRole) error {
	if mkc.FailDeleteClusterRole {
		return errors.New("ClusterRole failed")
	}
	return nil
}

// DeleteDaemonSet deletes the passed in daemonset
func (mkc *MockKubernetesClient) DeleteDaemonSet(pod *appsv1.DaemonSet) error {
	if mkc.FailDeleteDaemonSet {
		return errors.New("DaemonSet failed")
	}
	return nil
}

// DeleteDeployment deletes the passed in daemonset
func (mkc *MockKubernetesClient) DeleteDeployment(pod *appsv1.Deployment) error {
	if mkc.FailDeleteDeployment {
		return errors.New("deployment failed")
	}
	return nil
}

// DeletePod deletes the passed in pod
func (mkc *MockKubernetesClient) DeletePod(pod *v1.Pod) error {
	if mkc.FailDeletePod {
		return errors.New("DeletePod failed")
	}
	return nil
}

// EvictPod evicts the passed in pod using the passed in api version
func (mkc *MockKubernetesClient) EvictPod(pod *v1.Pod, policyGroupVersion string) error {
	if mkc.FailEvictPod {
		return errors.New("EvictPod failed")
	}
	return nil
}

// WaitForDelete waits until all pods are deleted. Returns all pods not deleted and an error on failure
func (mkc *MockKubernetesClient) WaitForDelete(logger *log.Entry, pods []v1.Pod, usingEviction bool) ([]v1.Pod, error) {
	if mkc.FailWaitForDelete {
		return nil, errors.New("WaitForDelete failed")
	}
	return []v1.Pod{}, nil
}

// DaemonSet returns a given daemonset in a namespace.
func (mkc *MockKubernetesClient) GetDaemonSet(namespace, name string) (*appsv1.DaemonSet, error) {
	return &appsv1.DaemonSet{
		Spec: appsv1.DaemonSetSpec{},
	}, nil
}

// GetDeployment returns a given deployment in a namespace.
func (mkc *MockKubernetesClient) GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	if mkc.FailGetDeploymentCount > 0 {
		mkc.FailGetDeploymentCount--
		return nil, errors.New("GetDeployment failed")
	}
	var replicas int32 = 1
	return &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
	}, nil
}

// UpdateDeployment updates a deployment to match the given specification.
func (mkc *MockKubernetesClient) UpdateDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	if mkc.FailUpdateDeploymentCount > 0 {
		mkc.FailUpdateDeploymentCount--
		return nil, errors.New("UpdateDeployment failed")
	}
	return &appsv1.Deployment{}, nil
}

// DeleteBlob mock
func (msc *MockStorageClient) DeleteBlob(container, blob string, options *azStorage.DeleteBlobOptions) error {
	return nil
}

// AddAcceptLanguages mock
func (mc *MockAKSEngineClient) AddAcceptLanguages(languages []string) {}

// DeployTemplate mock
func (mc *MockAKSEngineClient) DeployTemplate(ctx context.Context, resourceGroup, name string, template, parameters map[string]interface{}) (resources.DeploymentExtended, error) {
	switch {
	case mc.FailDeployTemplate:
		return resources.DeploymentExtended{}, errors.New("DeployTemplate failed")

	case mc.FailDeployTemplateQuota:
		errmsg := `resources.DeploymentsClient#CreateOrUpdate: Failure responding to request: StatusCode=400 -- Original Error: autorest/azure: Service returned an error`
		resp := `{
"error":{
	"code":"InvalidTemplateDeployment",
	"message":"The template deployment is not valid according to the validation procedure. The tracking id is 'b5bd7d6b-fddf-4ec3-a3b0-ce285a48bd31'. See inner errors for details. Please see https://aka.ms/arm-deploy for usage details.",
	"details":[{
		"code":"QuotaExceeded",
		"message":"Operation results in exceeding quota limits of Core. Maximum allowed: 10, Current in use: 10, Additional requested: 2. Please read more about quota increase at http://aka.ms/corequotaincrease."
}]}}`

		return resources.DeploymentExtended{
				Properties: &resources.DeploymentPropertiesExtended{
					Error: &resources.ErrorResponse{
						Code:    to.StringPtr("400"),
						Message: to.StringPtr(resp),
					},
				},
			},
			errors.New(errmsg)

	case mc.FailDeployTemplateConflict:
		errmsg := `resources.DeploymentsClient#CreateOrUpdate: Failure sending request: StatusCode=200 -- Original Error: Long running operation terminated with status 'Failed': Code="DeploymentFailed" Message="At least one resource deployment operation failed. Please list deployment operations for details. Please see https://aka.ms/arm-debug for usage details`
		resp := `{
"status":"Failed",
"error":{
	"code":"DeploymentFailed",
	"message":"At least one resource deployment operation failed. Please list deployment operations for details. Please see https://aka.ms/arm-debug for usage details.",
	"details":[{
		"code":"Conflict",
		"message":"{\r\n  \"error\": {\r\n    \"code\": \"PropertyChangeNotAllowed\",\r\n    \"target\": \"dataDisk.createOption\",\r\n    \"message\": \"Changing property 'dataDisk.createOption' is not allowed.\"\r\n  }\r\n}"
}]}}`
		return resources.DeploymentExtended{
				Properties: &resources.DeploymentPropertiesExtended{
					Error: &resources.ErrorResponse{
						Code:    to.StringPtr("200"),
						Message: to.StringPtr(resp),
					},
				},
			},
			errors.New(errmsg)

	case mc.FailDeployTemplateWithProperties:
		errmsg := `resources.DeploymentsClient#CreateOrUpdate: Failure sending request: StatusCode=200 -- Original Error: Long running operation terminated with status 'Failed': Code="DeploymentFailed" Message="At least one resource deployment operation failed. Please list deployment operations for details. Please see https://aka.ms/arm-debug for usage details`
		resp := `{
"status":"Failed",
"error":{
	"code":"DeploymentFailed",
	"message":"At least one resource deployment operation failed. Please list deployment operations for details. Please see https://aka.ms/arm-debug for usage details.",
	"details":[{
		"code":"Conflict",
		"message":"{\r\n  \"error\": {\r\n    \"code\": \"PropertyChangeNotAllowed\",\r\n    \"target\": \"dataDisk.createOption\",\r\n    \"message\": \"Changing property 'dataDisk.createOption' is not allowed.\"\r\n  }\r\n}"
}]}}`
		provisioningState := "Failed"
		return resources.DeploymentExtended{
				Properties: &resources.DeploymentPropertiesExtended{
					Error: &resources.ErrorResponse{
						Code:    to.StringPtr("200"),
						Message: to.StringPtr(resp),
					},
					ProvisioningState: &provisioningState,
				},
			},
			errors.New(errmsg)
	default:
		return resources.DeploymentExtended{}, nil
	}
}

// EnsureResourceGroup mock
func (mc *MockAKSEngineClient) EnsureResourceGroup(ctx context.Context, resourceGroup, location string, managedBy *string) (resources.ResourceGroup, error) {
	if mc.FailEnsureResourceGroup {
		return resources.ResourceGroup{}, errors.New("EnsureResourceGroup failed")
	}
	return resources.ResourceGroup{}, nil
}

// ListVirtualMachines mock
func (mc *MockAKSEngineClient) ListVirtualMachines(ctx context.Context, resourceGroup string) ([]*compute.VirtualMachine, error) {
	if mc.FailListVirtualMachines {
		return nil, errors.New("ListVirtualMachines failed")
	}
	if mc.FakeListVirtualMachineResult == nil {
		mc.FakeListVirtualMachineResult = func() []*compute.VirtualMachine {
			machine := mc.MakeFakeVirtualMachine(DefaultFakeVMName, defaultK8sVersionForFakeVMs)
			machine.Properties.AvailabilitySet = &compute.SubResource{
				ID: to.StringPtr("MockAvailabilitySet"),
			}
			return []*compute.VirtualMachine{&machine}
		}
	}
	return mc.FakeListVirtualMachineResult(), nil

}

// GetVirtualMachine mock
func (mc *MockAKSEngineClient) GetVirtualMachine(ctx context.Context, resourceGroup, name string) (compute.VirtualMachine, error) {
	if mc.FailGetVirtualMachine {
		return compute.VirtualMachine{}, errors.New("GetVirtualMachine failed")
	}
	return mc.MakeFakeVirtualMachine(DefaultFakeVMName, defaultK8sVersionForFakeVMs), nil
}

// RestartVirtualMachine mock
func (mc *MockAKSEngineClient) RestartVirtualMachine(ctx context.Context, resourceGroup, name string) error {
	if mc.FailRestartVirtualMachine {
		return errors.New("RestartVirtualMachine failed")
	}
	return nil
}

// MakeFakeVirtualMachine returns a fake compute.VirtualMachine
func (mc *MockAKSEngineClient) MakeFakeVirtualMachine(vmName string, orchestratorVersion string) compute.VirtualMachine {
	vm1Name := vmName

	creationSourceString := "creationSource"
	orchestratorString := "orchestrator"
	resourceNameSuffixString := "resourceNameSuffix"
	poolnameString := "poolName"

	creationSource := "aksengine-k8s-agentpool1-12345678-0"
	orchestrator := orchestratorVersion
	resourceNameSuffix := "12345678"
	poolname := "agentpool1"

	principalID := "00000000-1111-2222-3333-444444444444"

	tags := map[string]*string{
		creationSourceString:     &creationSource,
		orchestratorString:       &orchestrator,
		resourceNameSuffixString: &resourceNameSuffix,
		poolnameString:           &poolname,
	}

	var vmIdentity *compute.VirtualMachineIdentity
	if mc.ShouldSupportVMIdentity {
		vmIdentity = &compute.VirtualMachineIdentity{PrincipalID: &principalID}
	}

	if mc.FailListVirtualMachinesTags {
		tags = nil
	}

	osType := compute.OperatingSystemTypesLinux
	return compute.VirtualMachine{
		Name:     &vm1Name,
		Tags:     tags,
		Identity: vmIdentity,
		Properties: &compute.VirtualMachineProperties{
			StorageProfile: &compute.StorageProfile{
				OSDisk: &compute.OSDisk{
					OSType: &osType,
					Vhd: &compute.VirtualHardDisk{
						URI: &validOSDiskResourceName},
				},
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: []*compute.NetworkInterfaceReference{
					{
						ID: &validNicResourceName,
					},
				},
			},
		},
	}
}

// DeleteVirtualMachine mock
func (mc *MockAKSEngineClient) DeleteVirtualMachine(ctx context.Context, resourceGroup, name string) error {
	if mc.FailDeleteVirtualMachine {
		return errors.New("DeleteVirtualMachine failed")
	}

	return nil
}

// GetStorageClient mock

func (mc *MockAKSEngineClient) DeleteVirtualHardDisk(ctx context.Context, resourceGroup string, vhd *compute.VirtualHardDisk) error {
	return nil
}

// DeleteNetworkInterface mock
func (mc *MockAKSEngineClient) DeleteNetworkInterface(ctx context.Context, resourceGroup, nicName string) error {
	if mc.FailDeleteNetworkInterface {
		return errors.New("DeleteNetworkInterface failed")
	}

	return nil
}

var validOSDiskResourceName = "https://00k71r4u927seqiagnt0.blob.core.windows.net/osdisk/k8s-agentpool1-12345678-0-osdisk.vhd"
var validNicResourceName = "/subscriptions/DEC923E3-1EF1-4745-9516-37906D56DEC4/resourceGroups/acsK8sTest/providers/Microsoft.Network/networkInterfaces/k8s-agent-12345678-nic-0"

// RBAC Mocks

// DeleteManagedDisk is a wrapper around disksClient.Delete
func (mc *MockAKSEngineClient) DeleteManagedDisk(ctx context.Context, resourceGroupName string, diskName string) error {
	return nil
}

// ListManagedDisksByResourceGroup is a wrapper around disksClient.ListManagedDisksByResourceGroup
func (mc *MockAKSEngineClient) ListManagedDisksByResourceGroup(ctx context.Context, resourceGroupName string) ([]*compute.Disk, error) {
	return []*compute.Disk{}, nil
}

// GetKubernetesClient mock
func (mc *MockAKSEngineClient) GetKubernetesClient(apiserverURL, kubeConfig string, interval, timeout time.Duration) (kubernetes.Client, error) {
	if mc.FailGetKubernetesClient {
		return nil, errors.New("GetKubernetesClient failed")
	}

	if mc.MockKubernetesClient == nil {
		mc.MockKubernetesClient = &MockKubernetesClient{}
	}
	return mc.MockKubernetesClient, nil
}

// ListProviders mock
func (mc *MockAKSEngineClient) ListProviders(ctx context.Context) ([]*resources.Provider, error) {
	if mc.FailListProviders {
		return []*resources.Provider{}, errors.New("ListProviders failed")
	}

	return []*resources.Provider{}, nil
}

// ListDeploymentOperations gets all deployments operations for a deployment.
func (mc *MockAKSEngineClient) ListDeploymentOperations(ctx context.Context, resourceGroupName string, deploymentName string) ([]*resources.DeploymentOperation, error) {
	provisioningState := "Failed"
	id := "00000000"
	operationID := "d5062e45-6e9f-4fd3-a0a0-6b2c56b15757"
	return []*resources.DeploymentOperation{
		{
			OperationID: &operationID,
			ID:          &id,
			Properties: &resources.DeploymentOperationProperties{
				ProvisioningState: &provisioningState,
				// Error: &resources.ErrorResponse{
				// 	Code: "DeploymentFailed", "message": "At least one resource deployment operation failed. Please list deployment operations for details. Please see http://aka.ms/arm-debug for usage details.",
				// 	Details: []*resources.ErrorResponse{
				// 		&resources.ErrorResponse{
				// 			Code:    "Conflict",
				// 			Message: "{\r\n  \"error\": {\r\n    \"message\": \"Conflict\",\r\n    \"code\": \"Conflict\"\r\n  }\r\n}",
				// 		},
				// 	},
				// },
			},
		},
	}, nil
}

// ListDeploymentOperationsNextResults retrieves the next set of results, if any.
func (mc *MockAKSEngineClient) ListDeploymentOperationsNextResults(lastResults resources.DeploymentOperationsListResult) (result resources.DeploymentOperationsListResult, err error) {
	return resources.DeploymentOperationsListResult{}, nil
}

// DeleteRoleAssignmentByID deletes a roleAssignment via its unique identifier
func (mc *MockAKSEngineClient) DeleteRoleAssignmentByID(ctx context.Context, roleAssignmentID string) (authorization.RoleAssignment, error) {
	if mc.FailDeleteRoleAssignment {
		return authorization.RoleAssignment{}, errors.New("DeleteRoleAssignmentByID failed")
	}
	return authorization.RoleAssignment{}, nil
}

// ListRoleAssignmentsForPrincipal (e.g. a VM) via the scope and the unique identifier of the principal
func (mc *MockAKSEngineClient) ListRoleAssignmentsForPrincipal(ctx context.Context, scope string, principalID string) ([]*authorization.RoleAssignment, error) {
	roleAssignments := []*authorization.RoleAssignment{}
	if mc.ShouldSupportVMIdentity {
		var assignmentID = "role-assignment-id"
		roleAssignments = append(roleAssignments, &authorization.RoleAssignment{
			ID: &assignmentID,
		})
	}
	return roleAssignments, nil
}

// EnsureDefaultLogAnalyticsWorkspace mock
func (mc *MockAKSEngineClient) EnsureDefaultLogAnalyticsWorkspace(ctx context.Context, resourceGroup, location string) (workspaceResourceID string, err error) {
	if mc.FailEnsureDefaultLogAnalyticsWorkspace {
		return "", errors.New("EnsureDefaultLogAnalyticsWorkspace failed")
	}

	return "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-workspace-rg/providers/Microsoft.OperationalInsights/workspaces/test-workspace", nil
}

// AddContainerInsightsSolution mock
func (mc *MockAKSEngineClient) AddContainerInsightsSolution(ctx context.Context, workspaceSubscriptionID, workspaceResourceGroup, workspaceName, workspaceLocation string) (result bool, err error) {
	if mc.FailAddContainerInsightsSolution {
		return false, errors.New("AddContainerInsightsSolution failed")
	}

	return true, nil
}

// GetLogAnalyticsWorkspaceInfo mock
func (mc *MockAKSEngineClient) GetLogAnalyticsWorkspaceInfo(ctx context.Context, workspaceSubscriptionID, workspaceResourceGroup, workspaceName string) (workspaceID string, workspaceKey, workspaceLocation string, err error) {
	if mc.FailGetLogAnalyticsWorkspaceInfo {
		return "", "", "", errors.New("GetLogAnalyticsWorkspaceInfo failed")
	}

	return "00000000-0000-0000-0000-000000000000", "4D+vyd5/jScBmsAwZOF/0GOBQ5kuFQc9JVaW+HlnJ58cyePJcwTpks+rVmvgcXGmmyujLDNEVPiT8pB274a9Yg==", "westus", nil
}

// GetVirtualMachinePowerState returns the virtual machine's PowerState status code
func (mc *MockAKSEngineClient) GetVirtualMachinePowerState(ctx context.Context, resourceGroup, name string) (string, error) {
	return "", nil
}
