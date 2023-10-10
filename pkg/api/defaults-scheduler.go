// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import "github.com/Azure/aks-engine-azurestack/pkg/api/common"

// staticSchedulerConfig is not user-overridable
var staticSchedulerConfig = map[string]string{
	"--kubeconfig":   "/var/lib/kubelet/kubeconfig",
	"--leader-elect": "true",
}

// defaultSchedulerConfig provides targeted defaults, but is user-overridable
var defaultSchedulerConfig = map[string]string{
	"--v":               "2",
	"--profiling":       DefaultKubernetesSchedulerEnableProfiling,
	"--bind-address":    "127.0.0.1",    // STIG Rule ID: SV-242384r879530_rule
	"--tls-min-version": "VersionTLS12", // STIG Rule ID: SV-242377r879519_rule
}

func (cs *ContainerService) setSchedulerConfig() {
	o := cs.Properties.OrchestratorProfile

	// If no user-configurable scheduler config values exists, make an empty map, and fill in with defaults
	if o.KubernetesConfig.SchedulerConfig == nil {
		o.KubernetesConfig.SchedulerConfig = make(map[string]string)
	}

	for key, val := range defaultSchedulerConfig {
		// If we don't have a user-configurable scheduler config for each option
		if _, ok := o.KubernetesConfig.SchedulerConfig[key]; !ok {
			// then assign the default value
			o.KubernetesConfig.SchedulerConfig[key] = val
		}
	}

	// STIG Rule ID: SV-254801r879719_rule
	addDefaultFeatureGates(o.KubernetesConfig.SchedulerConfig, o.OrchestratorVersion, "1.25.0", "PodSecurity=true")

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	for key, val := range staticSchedulerConfig {
		o.KubernetesConfig.SchedulerConfig[key] = val
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
	}
	removeInvalidFeatureGates(o.KubernetesConfig.SchedulerConfig, invalidFeatureGates)
}
