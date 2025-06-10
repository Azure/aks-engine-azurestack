

# Input 
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>
<AzureCloudManagerImages>{{cloud_provider_image_versions}}</AzureCloudManagerImages>

# Input Validation
- Get Kubernetes version in xml tag <KubernetesVersion>
- Examinate the switch case to search `case "[MAJOR][MINOR]":`
	- If DOES NOT exist, return "True"; otherwise, return "False".
# Code Snippt Filter:
   - source code path: `pkg/api/k8s_versions.go`
   - object name: getK8sVersionComponents
   - object type: func
   - begin with: `func getK8sVersionComponents(version,`


# Fundamental Rules

- [ ] **Default Images per Kubernetes Version:**  
  For each Kubernetes version in the format `[MAJOR].[MINOR].[REVISION]` (e.g., `1.30.0`), ensure there are predefined default images for all Kubernetes components.

- [ ] **Version String Handling:**  
  The version string is already validated and will always be in the form `[MAJOR].[MINOR].[REVISION]` (e.g., `1.30.0`).  
  - **Do not add any logic to check for or modify the version string format (such as removing a `v` prefix).**

## Component version Check list
	

- [ ] **Add New Version Entry:**  
  - Create new switch case for `[MAJOR].[MINOR]` (e.g., `"1.30"`).
  **Important:**  
	- Do **not** update, remove, or modify any other entries.
	- Do **not** perform any additional logic or changes.
	- The only permitted action is adding the new key and Prune the old key.
  - **Reuse Previous Images:**  
    - Copy image values from the previous version (`[MAJOR].[MINOR-1]`).
	**Important:**  
		- When copying entries for a new version, ensure that for each key, the value is copied from the previous version using the **same key**.
		- Do not change the key or the lookup.
  - **Update Azure Cloud Manager Images:**  
    - Populate the following images for the new entry using the versions specified in the `<AzureCloudManagerImages>` XML tag:
      - `common.CloudControllerManagerComponentName`
      - `common.CloudNodeManagerAddonName`

- [ ] **Map Placement:**  
  - Insert the new switch case at the top of the map.

- [ ] - **Validation:**  
  - After generating the code, verify that for every switch case, the value is assigned using the same key for lookup as in the previous version, except for the explicitly updated Azure Cloud Manager images.

- [ ] **Prune Old Versions (Mandatory): in `switch majorMinor {`**
  - **Objective:** Remove case statements for versions less than or equal to `[MAJOR].[MINOR-5]`.
  - **Execution Steps:**
    1. **Identify Case Statements to Prune:**  
       - For each update, determine which case statements correspond to versions ≤ `[MAJOR].[MINOR-5]`.
    2. **Locate Case Statements:**  
       - Check if these case statements exist in the switch block.
    3. **Remove Outdated Cases:**  
       - Delete any case statements matching the identified versions.
    4. **Verify Removal:**  
       - Ensure all targeted case statements have been removed.
       - If any remain, repeat the process until none are left.
  - **Important:**  
    - Pruning old versions is mandatory for every update.
    - Never skip this step.
    - Always verify the final set of case statements to ensure only supported versions remain.

**You must review and ensure that all items on the **Version Update Check list** are checked. If any items are not checked, make the necessary changes to ensure all checkboxes are checked.**

## Example: Adding Kubernetes 1.30.10

**Example Entry for "1.30":**

case "1.30":
		ret = map[string]string{
			common.APIServerComponentName:                 getDefaultImage(common.APIServerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.ControllerManagerComponentName:         getDefaultImage(common.ControllerManagerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.KubeProxyAddonName:                     getDefaultImage(common.KubeProxyAddonName, kubernetesImageBaseType) + ":v" + version,
			common.SchedulerComponentName:                 getDefaultImage(common.SchedulerComponentName, kubernetesImageBaseType) + ":v" + version,
			common.CloudControllerManagerComponentName:    "oss/kubernetes/azure-cloud-controller-manager:v1.30.7",
			common.CloudNodeManagerAddonName:              "oss/kubernetes/azure-cloud-node-manager:v1.30.8",
			common.WindowsArtifactComponentName:           "v" + version + "/windowszip/v" + version + "-1int.zip",
			common.WindowsArtifactAzureStackComponentName: "v" + version + "/windowszip/v" + version + "-1int.zip",
			common.DashboardAddonName:                     dashboardImageReference,
			common.DashboardMetricsScraperContainerName:   dashboardMetricsScraperImageReference,
			common.ExecHealthZComponentName:               getDefaultImage(common.ExecHealthZComponentName, kubernetesImageBaseType),
			common.AddonResizerComponentName:              k8sComponent[common.AddonResizerComponentName],
			common.MetricsServerAddonName:                 k8sComponent[common.MetricsServerAddonName],
			common.CoreDNSAddonName:                       getDefaultImage(common.CoreDNSAddonName, kubernetesImageBaseType),
			common.CoreDNSAutoscalerName:                  clusterProportionalAutoscalerImageReference,
			common.KubeDNSAddonName:                       getDefaultImage(common.KubeDNSAddonName, kubernetesImageBaseType),
			common.AddonManagerComponentName:              k8sComponent[common.AddonManagerComponentName],
			common.DNSMasqComponentName:                   getDefaultImage(common.DNSMasqComponentName, kubernetesImageBaseType),
			common.PauseComponentName:                     pauseImageReference,
			common.TillerAddonName:                        tillerImageReference,
			common.ReschedulerAddonName:                   getDefaultImage(common.ReschedulerAddonName, kubernetesImageBaseType),
			common.ACIConnectorAddonName:                  virtualKubeletImageReference,
			common.ClusterAutoscalerAddonName:             k8sComponent[common.ClusterAutoscalerAddonName],
			common.DNSSidecarComponentName:                getDefaultImage(common.DNSSidecarComponentName, kubernetesImageBaseType),

			common.SMBFlexVolumeAddonName:                     smbFlexVolumeImageReference,
			common.IPMASQAgentAddonName:                       getDefaultImage(common.IPMASQAgentAddonName, kubernetesImageBaseType),
			common.AzureNetworkPolicyAddonName:                azureNPMContainerImageReference,
			common.CalicoTyphaComponentName:                   calicoTyphaImageReference,
			common.CalicoCNIComponentName:                     calicoCNIImageReference,
			common.CalicoNodeComponentName:                    calicoNodeImageReference,
			common.CalicoPod2DaemonComponentName:              calicoPod2DaemonImageReference,
			common.CalicoClusterAutoscalerComponentName:       calicoClusterProportionalAutoscalerImageReference,
			common.CiliumAgentContainerName:                   ciliumAgentImageReference,
			common.CiliumCleanStateContainerName:              ciliumCleanStateImageReference,
			common.CiliumOperatorContainerName:                ciliumOperatorImageReference,
			common.CiliumEtcdOperatorContainerName:            ciliumEtcdOperatorImageReference,
			common.AntreaControllerContainerName:              antreaControllerImageReference,
			common.AntreaAgentContainerName:                   antreaAgentImageReference,
			common.AntreaOVSContainerName:                     antreaOVSImageReference,
			"antrea" + common.AntreaInstallCNIContainerName:   antreaInstallCNIImageReference,
			common.NMIContainerName:                           aadPodIdentityNMIImageReference,
			common.MICContainerName:                           aadPodIdentityMICImageReference,
			common.AzurePolicyAddonName:                       azurePolicyImageReference,
			common.GatekeeperContainerName:                    gatekeeperImageReference,
			common.NodeProblemDetectorAddonName:               nodeProblemDetectorImageReference,
			common.CSIProvisionerContainerName:                k8sComponent[common.CSIProvisionerContainerName],
			common.CSIAttacherContainerName:                   k8sComponent[common.CSIAttacherContainerName],
			common.CSILivenessProbeContainerName:              k8sComponent[common.CSILivenessProbeContainerName],
			common.CSILivenessProbeWindowsContainerName:       k8sComponent[common.CSILivenessProbeWindowsContainerName],
			common.CSINodeDriverRegistrarContainerName:        k8sComponent[common.CSINodeDriverRegistrarContainerName],
			common.CSINodeDriverRegistrarWindowsContainerName: k8sComponent[common.CSINodeDriverRegistrarWindowsContainerName],
			common.CSISnapshotterContainerName:                k8sComponent[common.CSISnapshotterContainerName],
			common.CSISnapshotControllerContainerName:         k8sComponent[common.CSISnapshotControllerContainerName],
			common.CSIResizerContainerName:                    k8sComponent[common.CSIResizerContainerName],
			common.CSIAzureDiskContainerName:                  k8sComponent[common.CSIAzureDiskContainerName],
			common.CSIAzureFileContainerName:                  csiAzureFileImageReference,
			common.KubeFlannelContainerName:                   kubeFlannelImageReference,
			"flannel" + common.FlannelInstallCNIContainerName: flannelInstallCNIImageReference,
			common.KubeRBACProxyContainerName:                 KubeRBACProxyImageReference,
			common.ScheduledMaintenanceManagerContainerName:   ScheduledMaintenanceManagerImageReference,
			"nodestatusfreq":                                  DefaultKubernetesNodeStatusUpdateFrequency,
			"nodegraceperiod":                                 DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
			"podeviction":                                     DefaultKubernetesCtrlMgrPodEvictionTimeout,
			"routeperiod":                                     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
			"backoffretries":                                  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
			"backoffjitter":                                   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
			"backoffduration":                                 strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
			"backoffexponent":                                 strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
			"ratelimitqps":                                    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
			"ratelimitqpswrite":                               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPSWrite, 'f', -1, 64),
			"ratelimitbucket":                                 strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
			"ratelimitbucketwrite":                            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucketWrite),
			"gchighthreshold":                                 strconv.Itoa(DefaultKubernetesGCHighThreshold),
			"gclowthreshold":                                  strconv.Itoa(DefaultKubernetesGCLowThreshold),
			common.NVIDIADevicePluginAddonName:                nvidiaDevicePluginImageReference,
			common.CSISecretsStoreProviderAzureContainerName:  csiSecretsStoreProviderAzureImageReference,
			common.CSISecretsStoreDriverContainerName:         csiSecretsStoreDriverImageReference,
			common.AzureArcOnboardingAddonName:                azureArcOnboardingImageReference,
			common.AzureKMSProviderComponentName:              azureKMSProviderImageReference,
		}

**DO NOT make incorrect copy list the follwoing**
  common.CSIAttacherContainerName: k8sComponent[common.CSIAttacherComponentName], // Do not do this


**After making changes, you MUST review the **Version Update Check list** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**

