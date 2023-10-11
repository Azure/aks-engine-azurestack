// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"strconv"
	"strings"

	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
)

func (cs *ContainerService) setCloudControllerManagerConfig() {
	o := cs.Properties.OrchestratorProfile
	isAzureCNIDualStack := cs.Properties.IsAzureCNIDualStack()
	clusterCidr := o.KubernetesConfig.ClusterSubnet
	if isAzureCNIDualStack {
		clusterSubnets := strings.Split(clusterCidr, ",")
		if len(clusterSubnets) > 1 {
			clusterCidr = clusterSubnets[1]
		}
	}
	staticCloudControllerManagerConfig := map[string]string{
		"--allocate-node-cidrs":         strconv.FormatBool(!o.IsAzureCNI() || isAzureCNIDualStack),
		"--configure-cloud-routes":      strconv.FormatBool(cs.Properties.RequireRouteTable()),
		"--cloud-provider":              "azure",
		"--cloud-config":                "/etc/kubernetes/azure.json",
		"--cluster-cidr":                clusterCidr,
		"--kubeconfig":                  "/var/lib/kubelet/kubeconfig",
		"--leader-elect":                "true",
		"--route-reconciliation-period": "10s",
		"--v":                           "2",
	}

	// Disable cloud-node controller
	staticCloudControllerManagerConfig["--controllers"] = "*,-cloud-node"

	// Set --cluster-name based on appropriate DNS prefix
	if cs.Properties.MasterProfile != nil {
		staticCloudControllerManagerConfig["--cluster-name"] = cs.Properties.MasterProfile.DNSPrefix
	}

	// Default cloud-controller-manager config
	defaultCloudControllerManagerConfig := map[string]string{
		"--route-reconciliation-period": DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
	}

	// If no user-configurable cloud-controller-manager config values exists, use the defaults
	if o.KubernetesConfig.CloudControllerManagerConfig == nil {
		o.KubernetesConfig.CloudControllerManagerConfig = defaultCloudControllerManagerConfig
	} else {
		for key, val := range defaultCloudControllerManagerConfig {
			// If we don't have a user-configurable cloud-controller-manager config for each option
			if _, ok := o.KubernetesConfig.CloudControllerManagerConfig[key]; !ok {
				// then assign the default value
				o.KubernetesConfig.CloudControllerManagerConfig[key] = val
			}
		}
	}

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	for key, val := range staticCloudControllerManagerConfig {
		o.KubernetesConfig.CloudControllerManagerConfig[key] = val
	}

	invalidFeatureGates := []string{}
	// Remove --feature-gate VolumeSnapshotDataSource starting with 1.22
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.22.0-alpha.1") {
		invalidFeatureGates = append(invalidFeatureGates, "VolumeSnapshotDataSource")
	}
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.27.0") {
		// Remove --feature-gate ControllerManagerLeaderMigration starting with 1.27
		// Reference: https://github.com/kubernetes/kubernetes/pull/113534
		invalidFeatureGates = append(invalidFeatureGates, "ControllerManagerLeaderMigration")
		// Remove --feature-gate ExpandCSIVolumes, ExpandInUsePersistentVolumes, ExpandPersistentVolumes starting with 1.27
		// Reference: https://github.com/kubernetes/kubernetes/pull/113942
		invalidFeatureGates = append(invalidFeatureGates, "ExpandCSIVolumes", "ExpandInUsePersistentVolumes", "ExpandPersistentVolumes")
		// Remove --feature-gate CSIInlineVolume, CSIMigration, CSIMigrationAzureDisk, DaemonSetUpdateSurge, EphemeralContainers, IdentifyPodOS, LocalStorageCapacityIsolation, NetworkPolicyEndPort, StatefulSetMinReadySeconds starting with 1.27
		// Reference: https://github.com/kubernetes/kubernetes/pull/114410
		invalidFeatureGates = append(invalidFeatureGates, "CSIInlineVolume", "CSIMigration", "CSIMigrationAzureDisk", "DaemonSetUpdateSurge", "EphemeralContainers", "IdentifyPodOS", "LocalStorageCapacityIsolation", "NetworkPolicyEndPort", "StatefulSetMinReadySeconds")
	}
	removeInvalidFeatureGates(o.KubernetesConfig.CloudControllerManagerConfig, invalidFeatureGates)

	// TODO add RBAC support
	/*if *o.KubernetesConfig.EnableRbac {
		o.KubernetesConfig.CloudControllerManagerConfig["--use-service-account-credentials"] = "true"
	}*/
}
