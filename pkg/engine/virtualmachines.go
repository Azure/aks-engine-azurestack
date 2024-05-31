// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute"
)

func CreateMasterVM(cs *api.ContainerService) VirtualMachineARM {
	hasAvailabilityZones := cs.Properties.MasterProfile.HasAvailabilityZones()
	isStorageAccount := cs.Properties.MasterProfile.IsStorageAccount()
	kubernetesConfig := cs.Properties.OrchestratorProfile.KubernetesConfig

	var useManagedIdentity, userAssignedIDEnabled bool
	if kubernetesConfig != nil {
		useManagedIdentity = helpers.Bool(kubernetesConfig.UseManagedIdentity)
		userAssignedIDEnabled = kubernetesConfig.UserAssignedIDEnabled()
	}

	var dependencies []string
	dependentNIC := "[concat('Microsoft.Network/networkInterfaces/', variables('masterVMNamePrefix'), 'nic-', copyIndex(variables('masterOffset')))]"
	dependencies = append(dependencies, dependentNIC)
	if !hasAvailabilityZones {
		dependencies = append(dependencies, "[concat('Microsoft.Compute/availabilitySets/',variables('masterAvailabilitySet'))]")
	}
	if isStorageAccount {
		dependencies = append(dependencies, "[variables('masterStorageAccountName')]")
	}

	armResource := ARMResource{
		APIVersion: "[variables('apiVersionCompute')]",
		Copy: map[string]string{
			"count": "[sub(variables('masterCount'), variables('masterOffset'))]",
			"name":  "vmLoopNode",
		},
		DependsOn: dependencies,
	}

	vmTags := map[string]*string{
		"creationSource":     helpers.PointerToString("[concat(parameters('generatorCode'), '-', variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')))]"),
		"resourceNameSuffix": helpers.PointerToString("[parameters('nameSuffix')]"),
		"orchestrator":       helpers.PointerToString("[variables('orchestratorNameVersionTag')]"),
		"aksEngineVersion":   helpers.PointerToString("[parameters('aksEngineVersion')]"),
		"poolName":           helpers.PointerToString("master"),
	}

	if kubernetesConfig != nil && kubernetesConfig.IsContainerMonitoringAddonEnabled() {
		addon := kubernetesConfig.GetAddonByName(common.ContainerMonitoringAddonName)
		clusterDNSPrefix := "aks-engine-cluster"
		if cs.Properties.MasterProfile != nil && cs.Properties.MasterProfile.DNSPrefix != "" {
			clusterDNSPrefix = cs.Properties.MasterProfile.DNSPrefix
		}
		vmTags["logAnalyticsWorkspaceResourceId"] = helpers.PointerToString(addon.Config["logAnalyticsWorkspaceResourceId"])
		vmTags["clusterName"] = helpers.PointerToString(clusterDNSPrefix)
	}

	virtualMachine := compute.VirtualMachine{
		Location: helpers.PointerToString("[variables('location')]"),
		Name:     helpers.PointerToString("[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')))]"),
		Tags:     vmTags,
		Type:     helpers.PointerToString("Microsoft.Compute/virtualMachines"),
	}

	addCustomTagsToVM(cs.Properties.MasterProfile.CustomVMTags, &virtualMachine)

	if hasAvailabilityZones {
		virtualMachine.Zones = &[]string{
			"[string(parameters('availabilityZones')[mod(copyIndex(variables('masterOffset')), length(parameters('availabilityZones')))])]",
		}
	}

	if useManagedIdentity {
		identity := &compute.VirtualMachineIdentity{}
		if userAssignedIDEnabled {
			identity.Type = compute.ResourceIdentityTypeUserAssigned
			identity.UserAssignedIdentities = map[string]*compute.VirtualMachineIdentityUserAssignedIdentitiesValue{
				"[variables('userAssignedIDReference')]": {},
			}
		} else {
			identity.Type = compute.ResourceIdentityTypeSystemAssigned
		}
		virtualMachine.Identity = identity
	}

	vmProperties := &compute.VirtualMachineProperties{}

	if !hasAvailabilityZones {
		vmProperties.AvailabilitySet = &compute.SubResource{
			ID: helpers.PointerToString("[resourceId('Microsoft.Compute/availabilitySets',variables('masterAvailabilitySet'))]"),
		}
	}

	vmProperties.HardwareProfile = &compute.HardwareProfile{
		VMSize: compute.VirtualMachineSizeTypes(cs.Properties.MasterProfile.VMSize),
	}

	vmProperties.NetworkProfile = &compute.NetworkProfile{
		NetworkInterfaces: &[]compute.NetworkInterfaceReference{
			{
				ID: helpers.PointerToString("[resourceId('Microsoft.Network/networkInterfaces',concat(variables('masterVMNamePrefix'),'nic-', copyIndex(variables('masterOffset'))))]"),
			},
		},
	}

	osProfile := &compute.OSProfile{
		AdminUsername: helpers.PointerToString("[parameters('linuxAdminUsername')]"),
		ComputerName:  helpers.PointerToString("[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')))]"),
		LinuxConfiguration: &compute.LinuxConfiguration{
			DisablePasswordAuthentication: helpers.PointerToBool(true),
		},
	}

	linuxProfile := cs.Properties.LinuxProfile
	if linuxProfile != nil && len(linuxProfile.SSH.PublicKeys) > 1 {
		publicKeyPath := "[variables('sshKeyPath')]"
		var publicKeys []compute.SSHPublicKey
		for _, publicKey := range linuxProfile.SSH.PublicKeys {
			publicKeyTrimmed := strings.TrimSpace(publicKey.KeyData)
			publicKeys = append(publicKeys, compute.SSHPublicKey{
				Path:    &publicKeyPath,
				KeyData: &publicKeyTrimmed,
			})
		}
		osProfile.LinuxConfiguration.SSH = &compute.SSHConfiguration{
			PublicKeys: &publicKeys,
		}

	} else {
		osProfile.LinuxConfiguration.SSH = &compute.SSHConfiguration{
			PublicKeys: &[]compute.SSHPublicKey{
				{
					KeyData: helpers.PointerToString("[parameters('sshRSAPublicKey')]"),
					Path:    helpers.PointerToString("[variables('sshKeyPath')]"),
				},
			},
		}
	}

	t, err := InitializeTemplateGenerator(Context{})

	customDataStr := getCustomDataFromJSON(t.GetMasterCustomDataJSONObject(cs))
	osProfile.CustomData = helpers.PointerToString(customDataStr)

	if err != nil {
		panic(err)
	}

	if linuxProfile != nil && linuxProfile.HasSecrets() {
		vsg := getVaultSecretGroup(linuxProfile)
		osProfile.Secrets = &vsg
	}
	vmProperties.OsProfile = osProfile

	storageProfile := &compute.StorageProfile{}
	imageRef := cs.Properties.MasterProfile.ImageRef
	etcdSizeGB, _ := strconv.ParseInt(kubernetesConfig.EtcdDiskSizeGB, 10, 32)
	if !cs.Properties.MasterProfile.HasCosmosEtcd() {
		dataDisk := compute.DataDisk{
			CreateOption: compute.DiskCreateOptionTypesEmpty,
			DiskSizeGB:   helpers.PointerToInt32(int32(etcdSizeGB)),
			Lun:          helpers.PointerToInt32(0),
			Name:         helpers.PointerToString("[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')),'-etcddisk')]"),
		}
		if cs.Properties.MasterProfile.IsStorageAccount() {
			dataDisk.Vhd = &compute.VirtualHardDisk{
				URI: helpers.PointerToString("[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('masterStorageAccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'vhds/', variables('masterVMNamePrefix'),copyIndex(variables('masterOffset')),'-etcddisk.vhd')]"),
			}
		}
		storageProfile.DataDisks = &[]compute.DataDisk{
			dataDisk,
		}
	}
	imgReference := &compute.ImageReference{}
	if cs.Properties.MasterProfile.HasImageRef() {
		if cs.Properties.MasterProfile.HasImageGallery() {
			imgReference.ID = helpers.PointerToString(fmt.Sprintf("[concat('/subscriptions/', '%s', '/resourceGroups/', parameters('osImageResourceGroup'), '/providers/Microsoft.Compute/galleries/', '%s', '/images/', parameters('osImageName'), '/versions/', '%s')]", imageRef.SubscriptionID, imageRef.Gallery, imageRef.Version))
		} else {
			imgReference.ID = helpers.PointerToString("[resourceId(parameters('osImageResourceGroup'), 'Microsoft.Compute/images', parameters('osImageName'))]")
		}
	} else {
		imgReference.Offer = helpers.PointerToString("[parameters('osImageOffer')]")
		imgReference.Publisher = helpers.PointerToString("[parameters('osImagePublisher')]")
		imgReference.Sku = helpers.PointerToString("[parameters('osImageSku')]")
		imgReference.Version = helpers.PointerToString("[parameters('osImageVersion')]")
	}

	osDisk := &compute.OSDisk{
		Caching:      compute.CachingTypes(cs.Properties.MasterProfile.OSDiskCachingType),
		CreateOption: compute.DiskCreateOptionTypesFromImage,
	}

	if isStorageAccount {
		osDisk.Name = helpers.PointerToString("[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')),'-osdisk')]")
		osDisk.Vhd = &compute.VirtualHardDisk{
			URI: helpers.PointerToString("[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('masterStorageAccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'vhds/',variables('masterVMNamePrefix'),copyIndex(variables('masterOffset')),'-osdisk.vhd')]"),
		}
	}

	if cs.Properties.MasterProfile.OSDiskSizeGB > 0 {
		osDisk.DiskSizeGB = helpers.PointerToInt32(int32(cs.Properties.MasterProfile.OSDiskSizeGB))
	}

	if helpers.Bool(cs.Properties.MasterProfile.UltraSSDEnabled) {
		vmProperties.AdditionalCapabilities = &compute.AdditionalCapabilities{
			UltraSSDEnabled: helpers.PointerToBool(true),
		}
	}

	storageProfile.OsDisk = osDisk
	storageProfile.ImageReference = imgReference
	vmProperties.StorageProfile = storageProfile

	virtualMachine.VirtualMachineProperties = vmProperties

	return VirtualMachineARM{
		ARMResource:    armResource,
		VirtualMachine: virtualMachine,
	}
}

func createJumpboxVirtualMachine(cs *api.ContainerService) VirtualMachineARM {
	armResource := ARMResource{
		APIVersion: "[variables('apiVersionCompute')]",
		DependsOn: []string{
			"[concat('Microsoft.Network/networkInterfaces/', variables('jumpboxNetworkInterfaceName'))]",
		},
	}

	kubernetesConfig := cs.Properties.OrchestratorProfile.KubernetesConfig

	vm := compute.VirtualMachine{
		Location: helpers.PointerToString("[variables('location')]"),
		Name:     helpers.PointerToString("[parameters('jumpboxVMName')]"),
		Type:     helpers.PointerToString("Microsoft.Compute/virtualMachines"),
	}

	storageProfile := compute.StorageProfile{
		ImageReference: &compute.ImageReference{
			Publisher: helpers.PointerToString("Canonical"),
			Offer:     helpers.PointerToString("UbuntuServer"),
			Sku:       helpers.PointerToString("16.04-LTS"),
			Version:   helpers.PointerToString("latest"),
		},
		DataDisks: &[]compute.DataDisk{},
	}

	var jumpBoxIsManagedDisks bool
	if kubernetesConfig != nil && kubernetesConfig.PrivateCluster != nil {
		jumpBoxIsManagedDisks = kubernetesConfig.PrivateJumpboxProvision() && kubernetesConfig.PrivateCluster.JumpboxProfile.StorageProfile == api.ManagedDisks
	}

	if jumpBoxIsManagedDisks {
		storageProfile.OsDisk = &compute.OSDisk{
			CreateOption: compute.DiskCreateOptionTypesFromImage,
			DiskSizeGB:   helpers.PointerToInt32(int32(kubernetesConfig.PrivateCluster.JumpboxProfile.OSDiskSizeGB)),
			ManagedDisk: &compute.ManagedDiskParameters{
				StorageAccountType: "[variables('vmSizesMap')[parameters('jumpboxVMSize')].storageAccountType]",
			},
		}
	} else {
		storageProfile.OsDisk = &compute.OSDisk{
			CreateOption: compute.DiskCreateOptionTypesFromImage,
			Vhd: &compute.VirtualHardDisk{
				URI: helpers.PointerToString("[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('jumpboxStorageAccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'vhds/',parameters('jumpboxVMName'),'jumpboxdisk.vhd')]"),
			},
			Name: helpers.PointerToString("[variables('jumpboxOSDiskName')]"),
		}
	}

	t, err := InitializeTemplateGenerator(Context{})

	if err != nil {
		panic(err)
	}

	customDataStr := getCustomDataFromJSON(t.GetJumpboxCustomDataJSON(cs))

	vmProperties := compute.VirtualMachineProperties{
		HardwareProfile: &compute.HardwareProfile{
			VMSize: "[parameters('jumpboxVMSize')]",
		},
		OsProfile: &compute.OSProfile{
			ComputerName:  helpers.PointerToString("[parameters('jumpboxVMName')]"),
			AdminUsername: helpers.PointerToString("[parameters('jumpboxUsername')]"),
			LinuxConfiguration: &compute.LinuxConfiguration{
				DisablePasswordAuthentication: helpers.PointerToBool(true),
				SSH: &compute.SSHConfiguration{
					PublicKeys: &[]compute.SSHPublicKey{
						{
							Path:    helpers.PointerToString("[concat('/home/', parameters('jumpboxUsername'), '/.ssh/authorized_keys')]"),
							KeyData: helpers.PointerToString("[parameters('jumpboxPublicKey')]"),
						},
					},
				},
			},
			CustomData: helpers.PointerToString(customDataStr),
		},
		NetworkProfile: &compute.NetworkProfile{
			NetworkInterfaces: &[]compute.NetworkInterfaceReference{
				{
					ID: helpers.PointerToString("[resourceId('Microsoft.Network/networkInterfaces', variables('jumpboxNetworkInterfaceName'))]"),
				},
			},
		},
		StorageProfile: &storageProfile,
	}

	vm.VirtualMachineProperties = &vmProperties

	return VirtualMachineARM{
		ARMResource:    armResource,
		VirtualMachine: vm,
	}
}

func createAgentAvailabilitySetVM(cs *api.ContainerService, profile *api.AgentPoolProfile) VirtualMachineARM {
	var dependencies []string

	isStorageAccount := profile.IsStorageAccount()
	hasDisks := profile.HasDisks()
	kubernetesConfig := cs.Properties.OrchestratorProfile.KubernetesConfig

	var useManagedIdentity, userAssignedIDEnabled bool

	if kubernetesConfig != nil {
		useManagedIdentity = helpers.Bool(kubernetesConfig.UseManagedIdentity)
		userAssignedIDEnabled = kubernetesConfig.UserAssignedIDEnabled()
	}

	if isStorageAccount {
		storageDep := fmt.Sprintf("[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(variables('%[1]sOffset')),variables('maxVMsPerStorageAccount')),variables('%[1]sStorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(variables('%[1]sOffset')),variables('maxVMsPerStorageAccount')),variables('%[1]sStorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('%[1]sAccountName'))]", profile.Name)
		dependencies = append(dependencies, storageDep)
		if hasDisks {
			dataDiskDep := fmt.Sprintf("[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(add(div(copyIndex(variables('%[1]sOffset')),variables('maxVMsPerStorageAccount')),variables('%[1]sStorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(add(div(copyIndex(variables('%[1]sOffset')),variables('maxVMsPerStorageAccount')),variables('%[1]sStorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('%[1]sDataAccountName'))]", profile.Name)
			dependencies = append(dependencies, dataDiskDep)
		}
	}

	dependencies = append(dependencies, fmt.Sprintf("[concat('Microsoft.Network/networkInterfaces/', variables('%[1]sVMNamePrefix'), 'nic-', copyIndex(variables('%[1]sOffset')))]", profile.Name))

	dependencies = append(dependencies, fmt.Sprintf("[concat('Microsoft.Compute/availabilitySets/', variables('%[1]sAvailabilitySet'))]", profile.Name))

	if profile.IsWindows() {
		windowsProfile := cs.Properties.WindowsProfile
		// Add dependency for Image resource created by createWindowsImage()
		if windowsProfile.HasCustomImage() {
			dependencies = append(dependencies, fmt.Sprintf("%sCustomWindowsImage", profile.Name))
		}
	}

	tags := map[string]*string{
		"creationSource":   helpers.PointerToString(fmt.Sprintf("[concat(parameters('generatorCode'), '-', variables('%[1]sVMNamePrefix'), copyIndex(variables('%[1]sOffset')))]", profile.Name)),
		"orchestrator":     helpers.PointerToString("[variables('orchestratorNameVersionTag')]"),
		"aksEngineVersion": helpers.PointerToString("[parameters('aksEngineVersion')]"),
		"poolName":         helpers.PointerToString(profile.Name),
	}

	if profile.IsWindows() {
		tags["resourceNameSuffix"] = helpers.PointerToString("[variables('winResourceNamePrefix')]")
	} else {
		tags["resourceNameSuffix"] = helpers.PointerToString("[parameters('nameSuffix')]")
	}

	armResource := ARMResource{
		APIVersion: "[variables('apiVersionCompute')]",
		DependsOn:  dependencies,
		Copy: map[string]string{
			"count": fmt.Sprintf("[sub(variables('%[1]sCount'), variables('%[1]sOffset'))]", profile.Name),
			"name":  "vmLoopNode",
		},
	}

	virtualMachine := compute.VirtualMachine{
		Location: helpers.PointerToString("[variables('location')]"),
		Name:     helpers.PointerToString(fmt.Sprintf("[concat(variables('%[1]sVMNamePrefix'), copyIndex(variables('%[1]sOffset')))]", profile.Name)),
		Type:     helpers.PointerToString("Microsoft.Compute/virtualMachines"),
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: helpers.PointerToString(fmt.Sprintf("[resourceId('Microsoft.Network/networkInterfaces',concat(variables('%[1]sVMNamePrefix'), 'nic-', copyIndex(variables('%[1]sOffset'))))]", profile.Name)),
					},
				},
			},
		},
		Tags: tags,
	}

	if profile.IsFlatcar() {
		virtualMachine.Plan = &compute.Plan{
			Publisher: helpers.PointerToString(fmt.Sprintf("[parameters('%sosImagePublisher')]", profile.Name)),
			Name:      helpers.PointerToString(fmt.Sprintf("[parameters('%sosImageSKU')]", profile.Name)),
			Product:   helpers.PointerToString(fmt.Sprintf("[parameters('%sosImageOffer')]", profile.Name)),
		}
	}

	addCustomTagsToVM(profile.CustomVMTags, &virtualMachine)

	if useManagedIdentity {
		if userAssignedIDEnabled && !profile.IsWindows() {
			virtualMachine.Identity = &compute.VirtualMachineIdentity{
				Type: compute.ResourceIdentityTypeUserAssigned,
				UserAssignedIdentities: map[string]*compute.VirtualMachineIdentityUserAssignedIdentitiesValue{
					"[variables('userAssignedIDReference')]": {},
				},
			}
		} else {
			virtualMachine.Identity = &compute.VirtualMachineIdentity{
				Type: compute.ResourceIdentityTypeSystemAssigned,
			}
		}
	}

	virtualMachine.AvailabilitySet = &compute.SubResource{
		ID: helpers.PointerToString(fmt.Sprintf("[resourceId('Microsoft.Compute/availabilitySets',variables('%sAvailabilitySet'))]", profile.Name)),
	}

	vmSize := fmt.Sprintf("[variables('%sVMSize')]", profile.Name)

	virtualMachine.HardwareProfile = &compute.HardwareProfile{
		VMSize: compute.VirtualMachineSizeTypes(vmSize),
	}

	osProfile := compute.OSProfile{
		ComputerName: helpers.PointerToString(fmt.Sprintf("[concat(variables('%[1]sVMNamePrefix'), copyIndex(variables('%[1]sOffset')))]", profile.Name)),
	}

	t, err := InitializeTemplateGenerator(Context{})

	if !profile.IsWindows() {
		osProfile.AdminUsername = helpers.PointerToString("[parameters('linuxAdminUsername')]")
		osProfile.LinuxConfiguration = &compute.LinuxConfiguration{
			DisablePasswordAuthentication: helpers.PointerToBool(true),
		}

		linuxProfile := cs.Properties.LinuxProfile
		if linuxProfile != nil && len(linuxProfile.SSH.PublicKeys) > 1 {
			publicKeyPath := "[variables('sshKeyPath')]"
			publicKeys := []compute.SSHPublicKey{}
			for _, publicKey := range linuxProfile.SSH.PublicKeys {
				publicKeyTrimmed := strings.TrimSpace(publicKey.KeyData)
				publicKeys = append(publicKeys, compute.SSHPublicKey{
					Path:    &publicKeyPath,
					KeyData: &publicKeyTrimmed,
				})
			}
			osProfile.LinuxConfiguration.SSH = &compute.SSHConfiguration{
				PublicKeys: &publicKeys,
			}

		} else {
			osProfile.LinuxConfiguration.SSH = &compute.SSHConfiguration{
				PublicKeys: &[]compute.SSHPublicKey{
					{
						KeyData: helpers.PointerToString("[parameters('sshRSAPublicKey')]"),
						Path:    helpers.PointerToString("[variables('sshKeyPath')]"),
					},
				},
			}
		}

		if err != nil {
			panic(err)
		}

		agentCustomData := getCustomDataFromJSON(t.GetKubernetesLinuxNodeCustomDataJSONObject(cs, profile))
		osProfile.CustomData = helpers.PointerToString(agentCustomData)

		if linuxProfile != nil && linuxProfile.HasSecrets() {
			vsg := getVaultSecretGroup(linuxProfile)
			osProfile.Secrets = &vsg
		}
	} else {
		osProfile.AdminUsername = helpers.PointerToString("[parameters('windowsAdminUsername')]")
		osProfile.AdminPassword = helpers.PointerToString("[parameters('windowsAdminPassword')]")
		osProfile.WindowsConfiguration = &compute.WindowsConfiguration{
			EnableAutomaticUpdates: helpers.PointerToBool(cs.Properties.WindowsProfile.GetEnableWindowsUpdate()),
		}
		agentCustomData := getCustomDataFromJSON(t.GetKubernetesWindowsNodeCustomDataJSONObject(cs, profile))
		osProfile.CustomData = helpers.PointerToString(agentCustomData)

		if cs.Properties.WindowsProfile.HasEnableAHUB() {
			licenseType := api.WindowsLicenseTypeNone
			if cs.Properties.WindowsProfile.GetEnableAHUB() {
				licenseType = api.WindowsLicenseTypeServer
			}
			virtualMachine.LicenseType = &licenseType
		}
	}

	virtualMachine.OsProfile = &osProfile

	storageProfile := compute.StorageProfile{}

	if profile.IsWindows() {
		storageProfile.ImageReference = createWindowsImageReference(profile.Name, cs.Properties.WindowsProfile)

		if profile.HasDisks() {
			storageProfile.DataDisks = getArmDataDisks(profile)
		}
	} else {
		imageRef := profile.ImageRef
		if profile.HasImageRef() {
			if profile.HasImageGallery() {
				storageProfile.ImageReference = &compute.ImageReference{
					ID: helpers.PointerToString(fmt.Sprintf("[concat('/subscriptions/', '%s', '/resourceGroups/', parameters('%sosImageResourceGroup'), '/providers/Microsoft.Compute/galleries/', '%s', '/images/', parameters('%sosImageName'), '/versions/', '%s')]", imageRef.SubscriptionID, profile.Name, imageRef.Gallery, profile.Name, imageRef.Version)),
				}
			} else {
				storageProfile.ImageReference = &compute.ImageReference{
					ID: helpers.PointerToString(fmt.Sprintf("[resourceId(variables('%[1]sosImageResourceGroup'), 'Microsoft.Compute/images', variables('%[1]sosImageName'))]", profile.Name)),
				}
			}
		} else {
			storageProfile.ImageReference = &compute.ImageReference{
				Offer:     helpers.PointerToString(fmt.Sprintf("[variables('%sosImageOffer')]", profile.Name)),
				Publisher: helpers.PointerToString(fmt.Sprintf("[variables('%sosImagePublisher')]", profile.Name)),
				Sku:       helpers.PointerToString(fmt.Sprintf("[variables('%sosImageSKU')]", profile.Name)),
				Version:   helpers.PointerToString(fmt.Sprintf("[variables('%sosImageVersion')]", profile.Name)),
			}
			storageProfile.DataDisks = getArmDataDisks(profile)
		}
	}

	osDisk := compute.OSDisk{
		CreateOption: compute.DiskCreateOptionTypesFromImage,
		Caching:      compute.CachingTypes(profile.OSDiskCachingType),
	}

	if profile.IsStorageAccount() {
		osDisk.Name = helpers.PointerToString(fmt.Sprintf("[concat(variables('%[1]sVMNamePrefix'), copyIndex(variables('%[1]sOffset')),'-osdisk')]", profile.Name))
		osDisk.Vhd = &compute.VirtualHardDisk{
			URI: helpers.PointerToString(fmt.Sprintf("[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(variables('%[1]sOffset')),variables('maxVMsPerStorageAccount')),variables('%[1]sStorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(variables('%[1]sOffset')),variables('maxVMsPerStorageAccount')),variables('%[1]sStorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('%[1]sAccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'osdisk/', variables('%[1]sVMNamePrefix'), copyIndex(variables('%[1]sOffset')), '-osdisk.vhd')]", profile.Name)),
		}
	}

	if profile.IsEphemeral() {
		osDisk.DiffDiskSettings = &compute.DiffDiskSettings{
			Option: compute.Local,
		}
	}

	if profile.OSDiskSizeGB > 0 {
		osDisk.DiskSizeGB = helpers.PointerToInt32(int32(profile.OSDiskSizeGB))
	}

	if profile.DiskEncryptionSetID != "" {
		osDisk.ManagedDisk = &compute.ManagedDiskParameters{
			DiskEncryptionSet: &compute.DiskEncryptionSetParameters{ID: helpers.PointerToString(profile.DiskEncryptionSetID)},
		}
	}

	if helpers.Bool(profile.UltraSSDEnabled) {
		virtualMachine.AdditionalCapabilities = &compute.AdditionalCapabilities{
			UltraSSDEnabled: helpers.PointerToBool(true),
		}
	}

	storageProfile.OsDisk = &osDisk

	virtualMachine.StorageProfile = &storageProfile

	return VirtualMachineARM{
		ARMResource:    armResource,
		VirtualMachine: virtualMachine,
	}
}

func getArmDataDisks(profile *api.AgentPoolProfile) *[]compute.DataDisk {
	var dataDisks []compute.DataDisk
	for i, diskSize := range profile.DiskSizesGB {
		dataDisk := compute.DataDisk{
			DiskSizeGB:   helpers.PointerToInt32(int32(diskSize)),
			Lun:          helpers.PointerToInt32(int32(i)),
			CreateOption: compute.DiskCreateOptionTypesEmpty,
			Caching:      compute.CachingTypes(profile.DataDiskCachingType),
		}
		if profile.StorageProfile == api.StorageAccount {
			dataDisk.Name = helpers.PointerToString(fmt.Sprintf("[concat(variables('%sVMNamePrefix'), copyIndex(),'-datadisk%d')]", profile.Name, i))
			dataDisk.Vhd = &compute.VirtualHardDisk{
				URI: helpers.PointerToString(fmt.Sprintf("[concat('http://',variables('storageAccountPrefixes')[mod(add(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('%sStorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('%sStorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('%sDataAccountName'),'.blob.core.windows.net/vhds/',variables('%sVMNamePrefix'),copyIndex(), '--datadisk%d.vhd')]",
					profile.Name, profile.Name, profile.Name, profile.Name, i)),
			}
		}
		dataDisks = append(dataDisks, dataDisk)
	}
	return &dataDisks
}

func getCustomDataFromJSON(jsonStr string) string {
	var customDataObj map[string]string
	err := json.Unmarshal([]byte(jsonStr), &customDataObj)
	if err != nil {
		panic(err)
	}
	return customDataObj["customData"]
}

func getVaultSecretGroup(linuxProfile *api.LinuxProfile) []compute.VaultSecretGroup {
	var vaultSecretGroups []compute.VaultSecretGroup
	if linuxProfile.HasSecrets() {
		for idx, lVault := range linuxProfile.Secrets {
			computeVault := compute.VaultSecretGroup{
				SourceVault: &compute.SubResource{
					ID: helpers.PointerToString(fmt.Sprintf("[parameters('linuxKeyVaultID%d')]", idx)),
				},
			}
			var vaultCerts []compute.VaultCertificate
			for certIdx := range lVault.VaultCertificates {
				vaultCert := compute.VaultCertificate{
					CertificateURL: helpers.PointerToString(fmt.Sprintf("[parameters('linuxKeyVaultID%dCertificateURL%d')]", idx, certIdx)),
				}
				vaultCerts = append(vaultCerts, vaultCert)
			}
			computeVault.VaultCertificates = &vaultCerts
			vaultSecretGroups = append(vaultSecretGroups, computeVault)
		}
	}
	return vaultSecretGroups
}

func addCustomTagsToVM(tags map[string]string, vm *compute.VirtualMachine) {
	for key, value := range tags {
		_, found := vm.Tags[key]
		if !found {
			vm.Tags[key] = helpers.PointerToString(value)
		}
	}
}
