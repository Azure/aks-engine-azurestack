// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package kubernetesupgrade

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
	"github.com/Azure/aks-engine-azurestack/pkg/armhelpers"
	"github.com/Azure/aks-engine-azurestack/pkg/armhelpers/utils"
	"github.com/Azure/aks-engine-azurestack/pkg/i18n"
	"github.com/Azure/aks-engine-azurestack/pkg/kubernetes"
	compute "github.com/Azure/azure-sdk-for-go/profile/p20200901/resourcemanager/compute/armcompute"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	util "k8s.io/client-go/util/retry"
)

// ClusterTopology contains resources of the cluster the upgrade operation
// is targeting
type ClusterTopology struct {
	DataModel      *api.ContainerService
	SubscriptionID string
	Location       string
	ResourceGroup  string
	NameSuffix     string

	AgentPoolsToUpgrade map[string]bool
	AgentPools          map[string]*AgentPoolTopology

	MasterVMs         *[]*compute.VirtualMachine
	UpgradedMasterVMs *[]*compute.VirtualMachine
}

// AgentPoolTopology contains agent VMs in a single pool
type AgentPoolTopology struct {
	Identifier       *string
	Name             *string
	AgentVMs         *[]*compute.VirtualMachine
	UpgradedAgentVMs *[]*compute.VirtualMachine
}

// UpgradeCluster upgrades a cluster with Orchestrator version X.X to version Y.Y.
// Right now upgrades are supported for Kubernetes cluster only.
type UpgradeCluster struct {
	Translator *i18n.Translator
	Logger     *logrus.Entry
	ClusterTopology
	Client             armhelpers.AKSEngineClient
	StepTimeout        *time.Duration
	CordonDrainTimeout *time.Duration
	UpgradeWorkFlow    UpgradeWorkFlow
	Force              bool
	ControlPlaneOnly   bool
	CurrentVersion     string
}

// MasterPoolName pool name
const MasterPoolName = "master"

// UpgradeCluster runs the workflow to upgrade a Kubernetes cluster.
func (uc *UpgradeCluster) UpgradeCluster(az armhelpers.AKSEngineClient, kubeConfig string, aksEngineVersion string) error {
	uc.MasterVMs = &[]*compute.VirtualMachine{}
	uc.UpgradedMasterVMs = &[]*compute.VirtualMachine{}
	uc.AgentPools = make(map[string]*AgentPoolTopology)

	var kubeClient kubernetes.Client
	if az != nil {
		timeout := time.Duration(60) * time.Minute
		k, err := az.GetKubernetesClient("", kubeConfig, interval, timeout)
		if err != nil {
			uc.Logger.Warnf("Failed to get a Kubernetes client: %v", err)
		}
		kubeClient = k
	}

	if err := uc.setNodesToUpgrade(kubeClient, uc.ResourceGroup); err != nil {
		return uc.Translator.Errorf("Error while querying ARM for resources: %+v", err)
	}

	if kubeClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
		defer cancel()
		notReadyStream := uc.upgradedNotReadyStream(kubeClient, wait.Backoff{Steps: 15, Duration: 10 * time.Second})
		if err := uc.checkControlPlaneNodesStatus(ctx, notReadyStream); err != nil {
			uc.Logger.Error("Aborting the upgrade process to avoid potential control plane downtime")
			return errors.Wrap(err, "checking status of upgraded control plane nodes")
		}
	}

	kc := uc.DataModel.Properties.OrchestratorProfile.KubernetesConfig
	if kc != nil && kc.IsClusterAutoscalerEnabled() && !uc.ControlPlaneOnly {
		// pause the cluster-autoscaler before running upgrade and resume it afterward
		uc.Logger.Info("Pausing cluster autoscaler, replica count: 0")
		count, err := uc.SetClusterAutoscalerReplicaCount(kubeClient, 0)
		if err != nil {
			uc.Logger.Errorf("Failed to pause cluster-autoscaler: %v", err)
			if !uc.Force {
				return err
			}
		} else {
			if err == nil {
				defer func() {
					uc.Logger.Infof("Resuming cluster autoscaler, replica count: %d", count)
					if _, err = uc.SetClusterAutoscalerReplicaCount(kubeClient, count); err != nil {
						uc.Logger.Errorf("Failed to resume cluster-autoscaler: %v", err)
					}
				}()
			}
		}
	}

	upgradeVersion := uc.DataModel.Properties.OrchestratorProfile.OrchestratorVersion
	what := "control plane and all nodes"
	if uc.ControlPlaneOnly {
		what = "control plane nodes"
	}
	uc.Logger.Infof("Upgrading %s to Kubernetes version %s", what, upgradeVersion)

	if err := uc.getUpgradeWorkflow(kubeConfig, aksEngineVersion).RunUpgrade(); err != nil {
		return err
	}

	what = "Cluster"
	if uc.ControlPlaneOnly {
		what = "Control plane"
	}
	uc.Logger.Infof("%s upgraded successfully to Kubernetes version %s", what, upgradeVersion)
	return nil
}

// SetClusterAutoscalerReplicaCount changes the replica count of a cluster-autoscaler deployment.
func (uc *UpgradeCluster) SetClusterAutoscalerReplicaCount(kubeClient kubernetes.Client, replicaCount int32) (int32, error) {
	if kubeClient == nil {
		return 0, errors.New("no kubernetes client")
	}
	var count int32
	var err error
	const namespace, name, retries = "kube-system", "cluster-autoscaler", 10
	for attempt := 0; attempt < retries; attempt++ {
		deployment, getErr := kubeClient.GetDeployment(namespace, name)
		err = getErr
		if getErr == nil {
			count = *deployment.Spec.Replicas
			deployment.Spec.Replicas = &replicaCount
			if _, err = kubeClient.UpdateDeployment(namespace, deployment); err == nil {
				break
			}
		}
		sleepTime := time.Duration(rand.Intn(5))
		uc.Logger.Warnf("Failed to update cluster-autoscaler deployment: %v", err)
		uc.Logger.Infof("Retry updating cluster-autoscaler after %d seconds", sleepTime)
		time.Sleep(sleepTime * time.Second)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (uc *UpgradeCluster) getUpgradeWorkflow(kubeConfig string, aksEngineVersion string) UpgradeWorkFlow {
	if uc.UpgradeWorkFlow != nil {
		return uc.UpgradeWorkFlow
	}
	u := &Upgrader{}
	u.Init(uc.Translator, uc.Logger, uc.ClusterTopology, uc.Client, kubeConfig, uc.StepTimeout, uc.CordonDrainTimeout, aksEngineVersion, uc.ControlPlaneOnly)
	u.CurrentVersion = uc.CurrentVersion
	u.Force = uc.Force
	return u
}

func (uc *UpgradeCluster) setNodesToUpgrade(kubeClient kubernetes.Client, resourceGroup string) error {
	goalVersion := uc.DataModel.Properties.OrchestratorProfile.OrchestratorVersion

	ctx, cancel := context.WithTimeout(context.Background(), armhelpers.DefaultARMOperationTimeout)
	defer cancel()

	vmList, err := uc.Client.ListVirtualMachines(ctx, resourceGroup)
	if err != nil {
		return err
	}
	for _, vm := range vmList {
		// Windows VMs contain a substring of the name suffix
		if !strings.Contains(*(vm.Name), uc.NameSuffix) && !strings.Contains(*(vm.Name), uc.NameSuffix[:4]+"k8s") {
			uc.Logger.Infof("Skipping VM: %s for upgrade as it does not belong to cluster with expected name suffix: %s",
				*vm.Name, uc.NameSuffix)
			continue
		}
		currentVersion := uc.getNodeVersion(kubeClient, strings.ToLower(*vm.Name), vm.Tags, true)

		if uc.Force {
			if currentVersion == "" {
				currentVersion = "Unknown"
			}
			uc.addVMToUpgradeSets(vm, currentVersion)
		} else {
			if currentVersion == "" {
				uc.Logger.Infof("Skipping VM: %s for upgrade as the orchestrator version could not be determined.", *vm.Name)
				continue
			}
			// If the current version is different than the desired version then we add the VM to the list of VMs to upgrade.
			if currentVersion != goalVersion {
				if err := uc.upgradable(currentVersion); err != nil {
					return err
				}
				uc.addVMToUpgradeSets(vm, currentVersion)
			} else if currentVersion == goalVersion {
				uc.addVMToFinishedSets(vm, currentVersion)
			}
		}
	}

	return nil
}

func (uc *UpgradeCluster) upgradable(currentVersion string) error {
	nodeVersion := &api.OrchestratorProfile{
		OrchestratorType:    api.Kubernetes,
		OrchestratorVersion: currentVersion,
	}
	targetVersion := uc.DataModel.Properties.OrchestratorProfile.OrchestratorVersion

	orch, err := api.GetOrchestratorVersionProfile(nodeVersion, uc.DataModel.Properties.HasWindows(), uc.DataModel.Properties.IsAzureStackCloud())
	if err != nil {
		return err
	}

	for _, up := range orch.Upgrades {
		if up.OrchestratorVersion == targetVersion {
			return nil
		}
	}
	return errors.Errorf("%s cannot be upgraded to %s", currentVersion, targetVersion)
}

// getNodeVersion returns a node's current Kubernetes version via Kubernetes API or VM tag.
// For VMSS nodes, make sure OsProfile.ComputerName instead of VM name is used as the name here
// because the former is used as the K8s node name.
// Also, if the latest VMSS model is applied, then we can get the version info from the tags.
// Otherwise, we have to get version via K8s API. This is because VMSS does not support tags
// for individual instances and old/new instances have the same tags.
func (uc *UpgradeCluster) getNodeVersion(client kubernetes.Client, name string, tags map[string]*string, getVersionFromTags bool) string {
	if getVersionFromTags {
		if tags != nil && tags["orchestrator"] != nil {
			parts := strings.Split(*tags["orchestrator"], ":")
			if len(parts) == 2 {
				return parts[1]
			}
		}

		uc.Logger.Warnf("Expected tag \"orchestrator\" not found for VM: %s. Using Kubernetes API to retrieve Kubernetes version.", name)
	}

	if client != nil {
		node, err := client.GetNode(name)
		if err == nil {
			return strings.TrimPrefix(node.Status.NodeInfo.KubeletVersion, "v")
		}
		uc.Logger.Warnf("Failed to get node %s: %v", name, err)
	}
	return ""
}

func (uc *UpgradeCluster) addVMToAgentPool(vm *compute.VirtualMachine, isUpgradableVM bool) error {
	var poolIdentifier string
	var poolPrefix string
	var err error
	var vmPoolName string

	if vm.Tags != nil && vm.Tags["poolName"] != nil {
		vmPoolName = *vm.Tags["poolName"]
	} else {
		uc.Logger.Infof("poolName tag not found for VM: %s.", *vm.Name)
		// If there's only one agent pool, assume this VM is a member.
		agentPools := []string{}
		for k := range uc.AgentPoolsToUpgrade {
			if !strings.HasPrefix(k, "master") {
				agentPools = append(agentPools, k)
			}
		}
		if len(agentPools) == 1 {
			vmPoolName = agentPools[0]
		}
	}
	if vmPoolName == "" {
		uc.Logger.Warnf("Couldn't determine agent pool membership for VM: %s.", *vm.Name)
		return nil
	}

	uc.Logger.Infof("Evaluating VM: %s in pool: %s...", *vm.Name, vmPoolName)
	if vmPoolName == "" {
		uc.Logger.Infof("VM: %s does not contain `poolName` tag, skipping.", *vm.Name)
		return nil
	} else if !uc.AgentPoolsToUpgrade[vmPoolName] {
		uc.Logger.Infof("Skipping upgrade of VM: %s in pool: %s.", *vm.Name, vmPoolName)
		return nil
	}

	if *vm.Properties.StorageProfile.OSDisk.OSType == compute.OperatingSystemTypesWindows {
		poolPrefix, _, _, _, err = utils.WindowsVMNameParts(*vm.Name)
		if err != nil {
			uc.Logger.Errorf(err.Error())
			return err
		}
		if !strings.Contains(uc.NameSuffix, poolPrefix) {
			uc.Logger.Infof("Skipping VM: %s for upgrade as it does not belong to cluster with expected name suffix: %s",
				*vm.Name, uc.NameSuffix)
			return nil
		}

		// The k8s Windows VM Naming Format was previously "^([a-fA-F0-9]{5})([0-9a-zA-Z]{3})([a-zA-Z0-9]{4,6})$" (i.e.: 50621k8s9000)
		// The k8s Windows VM Naming Format is now "^([a-fA-F0-9]{4})([0-9a-zA-Z]{3})([0-9]{3,8})$" (i.e.: 1708k8s020)
		// The pool identifier is made of the first 11 or 9 characters
		if string((*vm.Name)[8]) == "9" {
			poolIdentifier = (*vm.Name)[:11]
		} else {
			poolIdentifier = (*vm.Name)[:9]
		}
	} else { // vm.StorageProfile.OsDisk.OsType == compute.Linux
		poolIdentifier, poolPrefix, _, err = utils.K8sLinuxVMNameParts(*vm.Name)
		if err != nil {
			uc.Logger.Errorf(err.Error())
			return err
		}

		if !strings.EqualFold(uc.NameSuffix, poolPrefix) {
			uc.Logger.Infof("Skipping VM: %s for upgrade as it does not belong to cluster with expected name suffix: %s",
				*vm.Name, uc.NameSuffix)
			return nil
		}
	}

	if uc.AgentPools[poolIdentifier] == nil {
		uc.AgentPools[poolIdentifier] =
			&AgentPoolTopology{&poolIdentifier, &vmPoolName, &[]*compute.VirtualMachine{}, &[]*compute.VirtualMachine{}}
	}

	orchestrator := "unknown"
	if vm.Tags != nil && vm.Tags["orchestrator"] != nil {
		orchestrator = *vm.Tags["orchestrator"]
	}
	//TODO(sterbrec): extract this from add to agentPool
	// separate the upgrade/skip decision from the agentpool composition
	if isUpgradableVM {
		uc.Logger.Infof("Adding Agent VM: %s, orchestrator: %s to pool: %s (AgentVMs)",
			*vm.Name, orchestrator, poolIdentifier)
		*uc.AgentPools[poolIdentifier].AgentVMs = append(*uc.AgentPools[poolIdentifier].AgentVMs, vm)
	} else {
		uc.Logger.Infof("Adding Agent VM: %s, orchestrator: %s to pool: %s (UpgradedAgentVMs)",
			*vm.Name, orchestrator, poolIdentifier)
		*uc.AgentPools[poolIdentifier].UpgradedAgentVMs = append(*uc.AgentPools[poolIdentifier].UpgradedAgentVMs, vm)
	}

	return nil
}

func (uc *UpgradeCluster) addVMToUpgradeSets(vm *compute.VirtualMachine, currentVersion string) {
	if strings.Contains(*(vm.Name), fmt.Sprintf("%s-", common.LegacyControlPlaneVMPrefix)) {
		uc.Logger.Infof("Master VM name: %s, orchestrator: %s (MasterVMs)", *vm.Name, currentVersion)
		*uc.MasterVMs = append(*uc.MasterVMs, vm)
	} else {
		if err := uc.addVMToAgentPool(vm, true); err != nil {
			uc.Logger.Errorf("Failed to add VM %s to agent pool: %s", *vm.Name, err)
		}
	}
}

func (uc *UpgradeCluster) addVMToFinishedSets(vm *compute.VirtualMachine, currentVersion string) {
	if strings.Contains(*(vm.Name), fmt.Sprintf("%s-", common.LegacyControlPlaneVMPrefix)) {
		uc.Logger.Infof("Master VM name: %s, orchestrator: %s (UpgradedMasterVMs)", *vm.Name, currentVersion)
		*uc.UpgradedMasterVMs = append(*uc.UpgradedMasterVMs, vm)
	} else {
		if err := uc.addVMToAgentPool(vm, false); err != nil {
			uc.Logger.Errorf("Failed to add VM %s to agent pool: %s", *vm.Name, err)
		}
	}
}

// checkControlPlaneNodesStatus checks whether it is safe to proceed with the upgrade process
// by looking at the status of previously upgraded control plane nodes.
//
// It returns an error if more than 1 of the already-upgraded control plane nodes are in the NotReady state.
// To recreate the node, users have to manually update the "orchestrator" tag on the VM.
func (uc *UpgradeCluster) checkControlPlaneNodesStatus(ctx context.Context, upgradedNotReadyStream <-chan []string) error {
	if len(*uc.UpgradedMasterVMs) == 0 {
		return nil
	}
	uc.Logger.Infoln("Checking status of upgraded control plane nodes")
	upgradedNotReadyCount := 0
loop:
	for {
		select {
		case upgradedNotReady, ok := <-upgradedNotReadyStream:
			if !ok {
				break loop
			}
			upgradedNotReadyCount = len(upgradedNotReady)
		case <-ctx.Done():
			break loop
		}
	}
	// return error if more than 1 upgraded node is not ready
	if upgradedNotReadyCount > 1 {
		uc.Logger.Error("At least 2 of the previously upgraded control plane nodes did not reach the NodeReady status")
		return errors.New("too many upgraded nodes are not ready")
	}
	return nil
}

func (uc *UpgradeCluster) upgradedNotReadyStream(client kubernetes.Client, backoff wait.Backoff) <-chan []string {
	alwaysRetry := func(_ error) bool {
		return true
	}
	upgraded := []string{}
	for _, vm := range *uc.UpgradedMasterVMs {
		upgraded = append(upgraded, *vm.Name)
	}
	stream := make(chan []string)
	go func() {
		defer close(stream)
		util.OnError(backoff, alwaysRetry, func() error { //nolint:errcheck
			upgradedNotReady, err := uc.getUpgradedNotReady(client, upgraded)
			if err != nil {
				return err
			}
			stream <- upgradedNotReady
			if len(upgradedNotReady) > 0 {
				return errors.New("retry to give NotReady nodes some extra time")
			}
			return nil
		})
	}()
	return stream
}

func (uc *UpgradeCluster) getUpgradedNotReady(client kubernetes.Client, upgraded []string) ([]string, error) {
	//TODO, the controlplane node will have both node-role.kubernetes.io/master and node-role.kubernetes.io/control-plane label
	// if node-role.kubernetes.io/master is removed in future change, also update the following label selector
	cpNodes, err := client.ListNodesByOptions(metav1.ListOptions{LabelSelector: "node-role.kubernetes.io/master"})
	if err != nil {
		return nil, err
	}
	nodeStatusMap := make(map[string]bool)
	for _, n := range cpNodes.Items {
		nodeStatusMap[n.Name] = kubernetes.IsNodeReady(&n)
	}
	upgradedNotReady := []string{}
	for _, vm := range upgraded {
		if ready, found := nodeStatusMap[vm]; found && !ready {
			upgradedNotReady = append(upgradedNotReady, vm)
		}
	}
	return upgradedNotReady, nil
}
