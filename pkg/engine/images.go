// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"fmt"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute"
)

func createWindowsImageReference(agentPoolProfileName string, windowsProfile *api.WindowsProfile) *compute.ImageReference {
	var computeImageRef compute.ImageReference

	if windowsProfile.HasCustomImage() {
		computeImageRef = compute.ImageReference{
			ID: helpers.PointerToString(fmt.Sprintf("[resourceId('Microsoft.Compute/images', '%sCustomWindowsImage')]", agentPoolProfileName)),
		}
	} else if windowsProfile.HasImageRef() {
		imageRef := windowsProfile.ImageRef
		if windowsProfile.HasImageGallery() {
			computeImageRef = compute.ImageReference{
				ID: helpers.PointerToString(fmt.Sprintf("[concat('/subscriptions/', '%s', '/resourceGroups/', parameters('agentWindowsImageResourceGroup'), '/providers/Microsoft.Compute/galleries/', '%s', '/images/', parameters('agentWindowsImageName'), '/versions/', '%s')]", imageRef.SubscriptionID, imageRef.Gallery, imageRef.Version)),
			}
		} else {
			computeImageRef = compute.ImageReference{
				ID: helpers.PointerToString("[resourceId(parameters('agentWindowsImageResourceGroup'), 'Microsoft.Compute/images', parameters('agentWindowsImageName'))]"),
			}
		}
	} else {
		computeImageRef = compute.ImageReference{
			Offer:     helpers.PointerToString("[parameters('agentWindowsOffer')]"),
			Publisher: helpers.PointerToString("[parameters('agentWindowsPublisher')]"),
			Sku:       helpers.PointerToString("[parameters('agentWindowsSku')]"),
			Version:   helpers.PointerToString("[parameters('agentWindowsVersion')]"),
		}
	}

	return &computeImageRef
}

func createWindowsImage(profile *api.AgentPoolProfile) ImageARM {
	return ImageARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionCompute')]",
		},
		Image: compute.Image{
			Type:     helpers.PointerToString("Microsoft.Compute/images"),
			Name:     helpers.PointerToString(fmt.Sprintf("%sCustomWindowsImage", profile.Name)),
			Location: helpers.PointerToString("[variables('location')]"),
			ImageProperties: &compute.ImageProperties{
				StorageProfile: &compute.ImageStorageProfile{
					OsDisk: &compute.ImageOSDisk{
						OsType:             "Windows",
						OsState:            compute.Generalized,
						BlobURI:            helpers.PointerToString("[parameters('agentWindowsSourceUrl')]"),
						StorageAccountType: compute.StorageAccountTypesStandardLRS,
					},
				},
				// TODO: Expose Hyper-V generation for VHD URL refs in apimodel
				HyperVGeneration: compute.HyperVGenerationTypesV1,
			},
		},
	}
}
