// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute"
	"github.com/google/go-cmp/cmp"
)

func TestCreateWindowsImageReference(t *testing.T) {
	cases := []struct {
		name        string
		profileName string
		w           api.WindowsProfile
		expected    compute.ImageReference
	}{
		{
			name:        "CustomImageUrl",
			profileName: "foobar",
			w: api.WindowsProfile{
				WindowsImageSourceURL: "https://some/image.vhd",
			},
			expected: compute.ImageReference{
				ID: helpers.PointerToString("[resourceId('Microsoft.Compute/images', 'foobarCustomWindowsImage')]"),
			},
		},
		{
			name:        "Image gallery reference",
			profileName: "foo",
			w: api.WindowsProfile{
				ImageRef: &api.ImageReference{
					Gallery:        "gallery",
					Name:           "test",
					ResourceGroup:  "testRg",
					SubscriptionID: "00000000-0000-0000-0000-000000000000",
					Version:        "0.1.0",
				},
			},
			expected: compute.ImageReference{
				ID: helpers.PointerToString("[concat('/subscriptions/', '00000000-0000-0000-0000-000000000000', '/resourceGroups/', parameters('agentWindowsImageResourceGroup'), '/providers/Microsoft.Compute/galleries/', 'gallery', '/images/', parameters('agentWindowsImageName'), '/versions/', '0.1.0')]"),
			},
		},
		{
			name:        "Image reference",
			profileName: "bar",
			w: api.WindowsProfile{
				ImageRef: &api.ImageReference{
					Name:          "tead",
					ResourceGroup: "testRg",
				},
			},
			expected: compute.ImageReference{
				ID: helpers.PointerToString("[resourceId(parameters('agentWindowsImageResourceGroup'), 'Microsoft.Compute/images', parameters('agentWindowsImageName'))]"),
			},
		},
		{
			name:        "Marketplace image",
			profileName: "baz",
			w: api.WindowsProfile{
				WindowsOffer:     "offer",
				WindowsPublisher: "pub",
				WindowsSku:       "sku",
				ImageVersion:     "ver",
			},
			expected: compute.ImageReference{
				Offer:     helpers.PointerToString("[parameters('agentWindowsOffer')]"),
				Publisher: helpers.PointerToString("[parameters('agentWindowsPublisher')]"),
				Sku:       helpers.PointerToString("[parameters('agentWindowsSku')]"),
				Version:   helpers.PointerToString("[parameters('agentWindowsVersion')]"),
			},
		},
		{
			name:        "Default",
			profileName: "qux",
			w:           api.WindowsProfile{},
			expected: compute.ImageReference{
				Offer:     helpers.PointerToString("[parameters('agentWindowsOffer')]"),
				Publisher: helpers.PointerToString("[parameters('agentWindowsPublisher')]"),
				Sku:       helpers.PointerToString("[parameters('agentWindowsSku')]"),
				Version:   helpers.PointerToString("[parameters('agentWindowsVersion')]"),
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := createWindowsImageReference(c.profileName, &c.w)
			expected := &c.expected

			diff := cmp.Diff(actual, expected)

			if diff != "" {
				t.Errorf("unexpected diff while comparing compute.ImageRefernce: %s", diff)
			}
		})
	}
}

func TestCreateWindowsImage(t *testing.T) {
	profile := &api.AgentPoolProfile{
		Name: "foobar",
	}

	actual := createWindowsImage(profile)

	expected := ImageARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionCompute')]",
		},
		Image: compute.Image{
			Type:     helpers.PointerToString("Microsoft.Compute/images"),
			Name:     helpers.PointerToString("foobarCustomWindowsImage"),
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
				HyperVGeneration: compute.HyperVGenerationTypesV1,
			},
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected diff while comparing windows images: %s", diff)
	}
}
