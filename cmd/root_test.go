// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package cmd

import (
	"os"
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/armhelpers"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"
	"github.com/Azure/aks-engine-azurestack/pkg/i18n"
	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

// mockAuthProvider implements AuthProvider and allows in particular to stub out getClient()
type mockAuthProvider struct {
	getClientMock armhelpers.AKSEngineClient
	*authArgs
}

func (provider *mockAuthProvider) getClient(env *api.Environment) (armhelpers.AKSEngineClient, error) {
	if provider.getClientMock == nil {
		return &armhelpers.MockAKSEngineClient{}, nil
	}
	return provider.getClientMock, nil

}
func (provider *mockAuthProvider) getAuthArgs() *authArgs {
	return provider.authArgs
}

func TestNewRootCmd(t *testing.T) {
	t.Parallel()

	command := NewRootCmd()
	if command.Use != rootName || command.Short != rootShortDescription || command.Long != rootLongDescription {
		t.Fatalf("root command should have use %s equal %s, short %s equal %s and long %s equal to %s", command.Use, rootName, command.Short, rootShortDescription, command.Long, rootLongDescription)
	}
	// The commands need to be listed in alphabetical order
	expectedCommands := []*cobra.Command{newAddPoolCmd(), getCompletionCmd(command), newDeployCmd(), newGenerateCmd(), newGetLogsCmd(), newGetVersionsCmd(), newOrchestratorsCmd(), newRotateCertsCmd(), newScaleCmd(), newUpgradeCmd(), newVersionCmd()}
	rc := command.Commands()

	for i, c := range expectedCommands {
		if rc[i].Use != c.Use {
			t.Fatalf("root command should have command %s, but found %s", c.Use, rc[i].Use)
		}
	}

	command.SetArgs([]string{"--debug"})
	err := command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}

func TestShowDefaultModelArg(t *testing.T) {
	t.Parallel()

	command := NewRootCmd()
	command.SetArgs([]string{"--show-default-model"})
	err := command.Execute()
	if err != nil {
		t.Fatal(err)
	}
	// TODO: examine command output
}

func TestDebugArg(t *testing.T) {
	t.Parallel()

	command := NewRootCmd()
	command.SetArgs([]string{"--show-default-model"})
	err := command.Execute()
	if err != nil {
		t.Fatal(err)
	}
	// TODO: examine command output
}

func TestCompletionCommand(t *testing.T) {
	t.Parallel()

	command := getCompletionCmd(NewRootCmd())
	command.SetArgs([]string{})
	err := command.Execute()
	if err != nil {
		t.Fatal(err)
	}
	// TODO: examine command output
}

func TestGetSelectedCloudFromAzConfig(t *testing.T) {
	for _, test := range []struct {
		desc   string
		data   []byte
		expect string
	}{
		{"nil file", nil, "AzureCloud"},
		{"empty file", []byte{}, "AzureCloud"},
		{"no cloud section", []byte(`
		[key]
		foo = bar
		`), "AzureCloud"},
		{"cloud section empty", []byte(`
		[cloud]
		[foo]
		foo = bar
		`), "AzureCloud"},
		{"AzureCloud selected", []byte(`
		[cloud]
		name = AzureCloud
		`), "AzureCloud"},
		{"custom cloud", []byte(`
		[cloud]
		name = myCloud
		`), "myCloud"},
	} {
		t.Run(test.desc, func(t *testing.T) {
			test := test
			t.Parallel()

			f, err := ini.Load(test.data)
			if err != nil {
				t.Fatal(err)
			}

			cloud := getSelectedCloudFromAzConfig(f)
			if cloud != test.expect {
				t.Fatalf("exepcted %q, got %q", test.expect, cloud)
			}
		})
	}
}

func TestGetCloudSubFromAzConfig(t *testing.T) {
	goodUUID, err := uuid.Parse("ccabad21-ea42-4ea1-affc-17ae73f9df66")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		desc   string
		data   []byte
		expect uuid.UUID
		err    bool
	}{
		{"empty file", []byte{}, uuid.UUID{}, true},
		{"no entry for cloud", []byte(`
		[SomeCloud]
		subscription = 00000000-0000-0000-0000-000000000000
		`), uuid.UUID{}, true},
		{"invalid UUID", []byte(`
		[AzureCloud]
		subscription = not-a-good-value
		`), uuid.UUID{}, true},
		{"real UUID", []byte(`
		[AzureCloud]
		subscription = ` + goodUUID.String() + `
		`), goodUUID, false},
	} {
		t.Run(test.desc, func(t *testing.T) {
			f, err := ini.Load(test.data)
			if err != nil {
				t.Fatal(err)
			}

			id, err := getCloudSubFromAzConfig("AzureCloud", f)

			if test.err != (err != nil) {
				t.Fatalf("expected err=%v, got: %v", test.err, err)
			}
			if test.err {
				return
			}
			if id.String() != test.expect.String() {
				t.Fatalf("expected %s, got %s", test.expect, id)
			}
		})
	}
}

func TestWriteCustomCloudProfile(t *testing.T) {
	t.Parallel()

	cs, err := prepareCustomCloudProfile()
	if err != nil {
		t.Fatalf("failed to prepare custom cloud profile: %v", err)
	}
	if err = writeCustomCloudProfile(cs); err != nil {
		t.Fatalf("failed to write custom cloud profile: %v", err)
	}

	environmentFilePath := os.Getenv("AZURE_ENVIRONMENT_FILEPATH")
	if environmentFilePath == "" {
		t.Fatal("failed to write custom cloud profile: err - AZURE_ENVIRONMENT_FILEPATH is empty")
	}

	if _, err = os.Stat(environmentFilePath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		t.Fatalf("failed to write custom cloud profile: file %s does not exist", environmentFilePath)
	}

	azurestackenvironment, err := os.ReadFile(environmentFilePath)
	if err != nil {
		t.Fatalf("failed to write custom cloud profile: can not read file %s ", environmentFilePath)
	}
	azurestackenvironmentStr := string(azurestackenvironment)
	expectedResult := `{"name":"azurestackcloud","managementPortalURL":"https://management.local.azurestack.external/","publishSettingsURL":"https://management.local.azurestack.external/publishsettings/index","serviceManagementEndpoint":"https://management.azurestackci15.onmicrosoft.com/36f71706-54df-4305-9847-5b038a4cf189","resourceManagerEndpoint":"https://management.local.azurestack.external/","activeDirectoryEndpoint":"https://login.windows.net/","galleryEndpoint":"https://portal.local.azurestack.external=30015/","keyVaultEndpoint":"https://vault.azurestack.external/","graphEndpoint":"https://graph.windows.net/","serviceBusEndpoint":"https://servicebus.azurestack.external/","batchManagementEndpoint":"https://batch.azurestack.external/","storageEndpointSuffix":"core.azurestack.external","sqlDatabaseDNSSuffix":"database.azurestack.external","trafficManagerDNSSuffix":"trafficmanager.cn","keyVaultDNSSuffix":"vault.azurestack.external","serviceBusEndpointSuffix":"servicebus.azurestack.external","serviceManagementVMDNSSuffix":"chinacloudapp.cn","resourceManagerVMDNSSuffix":"cloudapp.azurestack.external","containerRegistryDNSSuffix":"azurecr.io","cosmosDBDNSSuffix":"","tokenAudience":"https://management.azurestack.external/","apiManagementHostNameSuffix":"","synapseEndpointSuffix":"","resourceIdentifiers":{"graph":"","keyVault":"","datalake":"","batch":"","operationalInsights":"","storage":"","synapse":"","serviceBus":""}}`
	if azurestackenvironmentStr != expectedResult {
		t.Fatalf("failed to write custom cloud profile: expected %s , got %s ", expectedResult, azurestackenvironmentStr)
	}
}

func TestValidateAuthArgs(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	validID := "cc6b141e-6afc-4786-9bf6-e3b9a5601460"
	invalidID := "invalidID"

	for _, tc := range []struct {
		name     string
		authArgs authArgs
		expected error
	}{
		{
			name: "AuthMethodIsRequired",
			authArgs: authArgs{
				AuthMethod: "",
			},
			expected: errors.New("--auth-method is a required parameter"),
		},
		{
			name: "AlwaysExpectValidClientID",
			authArgs: authArgs{
				rawSubscriptionID:   validID,
				rawClientID:         invalidID,
				ClientSecret:        "secret",
				AuthMethod:          "client_secret",
				RawAzureEnvironment: "AZUREPUBLICCLOUD",
			},
			expected: errors.New(`parsing --client-id: invalid UUID length: 9`),
		},
		{
			name: "AlwaysExpectValidClientID",
			authArgs: authArgs{
				rawSubscriptionID:   validID,
				rawClientID:         invalidID,
				ClientSecret:        "secret",
				AuthMethod:          "client_certificate",
				RawAzureEnvironment: "AZUREPUBLICCLOUD",
			},
			expected: errors.New(`parsing --client-id: invalid UUID length: 9`),
		},
		{
			name: "ClientSecretAuthExpectsClientSecret",
			authArgs: authArgs{
				rawSubscriptionID:   validID,
				rawClientID:         validID,
				ClientSecret:        "",
				AuthMethod:          "client_secret",
				RawAzureEnvironment: "AZUREPUBLICCLOUD",
			},
			expected: errors.New(`--client-secret must be specified when --auth-method="client_secret"`),
		},
		{
			name: "ValidClientSecretAuth",
			authArgs: authArgs{
				rawSubscriptionID:   validID,
				rawClientID:         validID,
				ClientSecret:        "secret",
				AuthMethod:          "client_secret",
				RawAzureEnvironment: "AZUREPUBLICCLOUD",
			},
			expected: nil,
		},
		{
			name: "ClientCertificateAuthExpectsCertificatePath",
			authArgs: authArgs{
				rawSubscriptionID:   validID,
				rawClientID:         validID,
				CertificatePath:     "",
				PrivateKeyPath:      "/a/path",
				AuthMethod:          "client_certificate",
				RawAzureEnvironment: "AZUREPUBLICCLOUD",
			},
			expected: errors.New(`--certificate-path and --private-key-path must be specified when --auth-method="client_certificate"`),
		},
		{
			name: "ClientCertificateAuthExpectsPrivateKeyPath",
			authArgs: authArgs{
				rawSubscriptionID:   validID,
				rawClientID:         validID,
				CertificatePath:     "/a/path",
				PrivateKeyPath:      "",
				AuthMethod:          "client_certificate",
				RawAzureEnvironment: "AZUREPUBLICCLOUD",
			},
			expected: errors.New(`--certificate-path and --private-key-path must be specified when --auth-method="client_certificate"`),
		},
		{
			name: "ValidClientCertificateAuth",
			authArgs: authArgs{
				rawSubscriptionID:   validID,
				rawClientID:         validID,
				CertificatePath:     "/a/path",
				PrivateKeyPath:      "/a/path",
				AuthMethod:          "client_certificate",
				RawAzureEnvironment: "AZUREPUBLICCLOUD",
			},
			expected: nil,
		},
	} {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			err := test.authArgs.validateAuthArgs()
			if test.expected != nil {
				g.Expect(err.Error()).To(Equal(test.expected.Error()))
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func prepareCustomCloudProfile() (*api.ContainerService, error) {
	const (
		name                         = "azurestackcloud"
		managementPortalURL          = "https://management.local.azurestack.external/"
		publishSettingsURL           = "https://management.local.azurestack.external/publishsettings/index"
		serviceManagementEndpoint    = "https://management.azurestackci15.onmicrosoft.com/36f71706-54df-4305-9847-5b038a4cf189"
		resourceManagerEndpoint      = "https://management.local.azurestack.external/"
		activeDirectoryEndpoint      = "https://login.windows.net/"
		galleryEndpoint              = "https://portal.local.azurestack.external=30015/"
		keyVaultEndpoint             = "https://vault.azurestack.external/"
		graphEndpoint                = "https://graph.windows.net/"
		serviceBusEndpoint           = "https://servicebus.azurestack.external/"
		batchManagementEndpoint      = "https://batch.azurestack.external/"
		storageEndpointSuffix        = "core.azurestack.external"
		sqlDatabaseDNSSuffix         = "database.azurestack.external"
		trafficManagerDNSSuffix      = "trafficmanager.cn"
		keyVaultDNSSuffix            = "vault.azurestack.external"
		serviceBusEndpointSuffix     = "servicebus.azurestack.external"
		serviceManagementVMDNSSuffix = "chinacloudapp.cn"
		resourceManagerVMDNSSuffix   = "cloudapp.azurestack.external"
		containerRegistryDNSSuffix   = "azurecr.io"
		tokenAudience                = "https://management.azurestack.external/"
	)
	cs := &api.ContainerService{
		Properties: &api.Properties{
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ClientID: "barClientID",
				Secret:   "bazSecret",
			},
			MasterProfile: &api.MasterProfile{
				Count:     1,
				DNSPrefix: "blueorange",
				VMSize:    "Standard_D2_v2",
			},
			OrchestratorProfile: &api.OrchestratorProfile{
				OrchestratorType: api.Kubernetes,
			},
			LinuxProfile: &api.LinuxProfile{},
			CustomCloudProfile: &api.CustomCloudProfile{
				IdentitySystem:       api.AzureADIdentitySystem,
				AuthenticationMethod: api.ClientSecretAuthMethod,
				Environment: &api.Environment{
					Name:                         name,
					ManagementPortalURL:          managementPortalURL,
					PublishSettingsURL:           publishSettingsURL,
					ServiceManagementEndpoint:    serviceManagementEndpoint,
					ResourceManagerEndpoint:      resourceManagerEndpoint,
					ActiveDirectoryEndpoint:      activeDirectoryEndpoint,
					GalleryEndpoint:              galleryEndpoint,
					KeyVaultEndpoint:             keyVaultEndpoint,
					GraphEndpoint:                graphEndpoint,
					ServiceBusEndpoint:           serviceBusEndpoint,
					BatchManagementEndpoint:      batchManagementEndpoint,
					StorageEndpointSuffix:        storageEndpointSuffix,
					SQLDatabaseDNSSuffix:         sqlDatabaseDNSSuffix,
					TrafficManagerDNSSuffix:      trafficManagerDNSSuffix,
					KeyVaultDNSSuffix:            keyVaultDNSSuffix,
					ServiceBusEndpointSuffix:     serviceBusEndpointSuffix,
					ServiceManagementVMDNSSuffix: serviceManagementVMDNSSuffix,
					ResourceManagerVMDNSSuffix:   resourceManagerVMDNSSuffix,
					ContainerRegistryDNSSuffix:   containerRegistryDNSSuffix,
					TokenAudience:                tokenAudience,
				},
			},
			AgentPoolProfiles: []*api.AgentPoolProfile{
				{
					Name:   "agentpool1",
					VMSize: "Standard_D2_v2",
					Count:  2,
				},
			},
		},
	}

	_, err := cs.SetPropertiesDefaults(api.PropertiesDefaultsParams{
		IsScale:    false,
		IsUpgrade:  false,
		PkiKeySize: helpers.DefaultPkiKeySize,
	})
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func TestWriteArtifacts(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	cs := api.CreateMockContainerService("testcluster", "", 3, 2, false)
	_, err := cs.SetPropertiesDefaults(api.PropertiesDefaultsParams{
		IsScale:    false,
		IsUpgrade:  false,
		PkiKeySize: helpers.DefaultPkiKeySize,
	})
	g.Expect(err).NotTo(HaveOccurred())

	outdir, del := makeTmpDir(t)
	defer del()

	err = writeArtifacts(outdir, cs, "vlabs", &i18n.Translator{})
	g.Expect(err).NotTo(HaveOccurred())
}

func makeTmpDir(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "_tmp_dir")
	if err != nil {
		t.Fatalf("unable to create dir: %s", err.Error())
	}
	return tmpDir, func() { defer os.RemoveAll(tmpDir) }
}
