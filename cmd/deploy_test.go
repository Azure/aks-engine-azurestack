// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package cmd

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/armhelpers"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const ExampleAPIModel = `{
  "apiVersion": "vlabs",
  "properties": {
		"orchestratorProfile": { "orchestratorType": "Kubernetes", "kubernetesConfig": { "useManagedIdentity": %s, "etcdVersion" : "2.3.8" } },
    "masterProfile": { "count": 1, "dnsPrefix": "", "vmSize": "Standard_D2_v2" },
    "agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
    "windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
    "linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
    },
    "servicePrincipalProfile": { "clientId": "%s", "secret": "%s" }
  }
}
`

const ExampleAPIModelWithDNSPrefix = `{
	"apiVersion": "vlabs",
	"properties": {
		  "orchestratorProfile": { "orchestratorType": "Kubernetes", "kubernetesConfig": { "useManagedIdentity": %s, "etcdVersion" : "2.3.8" } },
	  "masterProfile": { "count": 1, "dnsPrefix": "mytestcluster", "vmSize": "Standard_D2_v2" },
	  "agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
	  "windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
	  "linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
	  },
	  "servicePrincipalProfile": { "clientId": "%s", "secret": "%s" }
	}
  }
  `

const ExampleAPIModelWithoutDNSPrefix = `{
	"apiVersion": "vlabs",
	"properties": {
		  "orchestratorProfile": { "orchestratorType": "Kubernetes", "kubernetesConfig": { "useManagedIdentity": %s, "etcdVersion" : "2.3.8" } },
	  "masterProfile": { "count": 1, "vmSize": "Standard_D2_v2" },
	  "agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
	  "windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
	  "linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
	  },
	  "servicePrincipalProfile": { "clientId": "%s", "secret": "%s" }
	}
  }
  `

const ExampleAPIModelWithoutServicePrincipalProfile = `{
	"apiVersion": "vlabs",
	"properties": {
		  "orchestratorProfile": { "orchestratorType": "Kubernetes", "kubernetesConfig": { "useManagedIdentity": %s, "etcdVersion" : "2.3.8" } },
	  "masterProfile": { "count": 1, "dnsPrefix": "mytestcluster", "vmSize": "Standard_D2_v2" },
	  "agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
	  "windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
	  "linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
	  }
	}
  }
  `
const ExampleAPIModelWithContainerMonitoringAddonWithNoConfig = `{
	"apiVersion": "vlabs",
	"properties": {
		  "orchestratorProfile": { "orchestratorType": "Kubernetes", "kubernetesConfig": { "addons":[{"name": "container-monitoring","enabled": true }]}},
	  "masterProfile": { "count": 1, "dnsPrefix": "mytestcluster", "vmSize": "Standard_D2_v2" },
	  "agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
	  "windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
	  "linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
	  }
	}
  }
  `

const ExampleAPIModelWithContainerMonitoringAddonWithExistingWorkspaceConfig = `{
	"apiVersion": "vlabs",
	"properties": {
		  "orchestratorProfile": { "orchestratorType": "Kubernetes", "kubernetesConfig": { "addons":[{"name": "container-monitoring","enabled": true, "config":{"logAnalyticsWorkspaceResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-workspace-rg/providers/Microsoft.OperationalInsights/workspaces/test-workspace"} }]}},
	  "masterProfile": { "count": 1, "dnsPrefix": "mytestcluster", "vmSize": "Standard_D2_v2" },
	  "agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
	  "windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
	  "linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
	  }
	}
  }
  `
const ExampleAPIModelWithContainerMonitoringAddonWithWorkspaceGUIDAndKeyConfig = `{
	"apiVersion": "vlabs",
	"properties": {
		  "orchestratorProfile": { "orchestratorType": "Kubernetes", "kubernetesConfig": { "addons":[{"name": "container-monitoring","enabled": true, "config":{"workspaceGuid": "MDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAw", "workspaceKey": "NEQrdnlkNS9qU2NCbXNBd1pPRi8wR09CUTVrdUZRYzlKVmFXK0hsbko1OGN5ZVBKY3dUcGtzK3JWbXZnY1hHbW15dWpMRE5FVlBpVDhwQjI3NGE5WWc9PQ=="} }]}},
	  "masterProfile": { "count": 1, "dnsPrefix": "mytestcluster", "vmSize": "Standard_D2_v2" },
	  "agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
	  "windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
	  "linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
	  }
	}
  }
  `

func getExampleAPIModel(useManagedIdentity bool, clientID, clientSecret string) string {
	return getAPIModel(ExampleAPIModel, useManagedIdentity, clientID, clientSecret)
}

func getAPIModel(baseAPIModel string, useManagedIdentity bool, clientID, clientSecret string) string {
	return fmt.Sprintf(
		baseAPIModel,
		strconv.FormatBool(useManagedIdentity),
		clientID,
		clientSecret)
}

func getAPIModelWithoutServicePrincipalProfile(useManagedIdentity bool) string {
	return fmt.Sprintf(
		ExampleAPIModelWithoutServicePrincipalProfile,
		strconv.FormatBool(useManagedIdentity))
}

func TestNewDeployCmd(t *testing.T) {
	t.Parallel()

	command := newDeployCmd()
	if command.Use != deployName || command.Short != deployShortDescription || command.Long != deployLongDescription {
		t.Fatalf("deploy command should have use %s equal %s, short %s equal %s and long %s equal to %s", command.Use, deployName, command.Short, deployShortDescription, command.Long, versionLongDescription)
	}

	expectedFlags := []string{"api-model", "dns-prefix", "auto-suffix", "output-directory", "ca-private-key-path", "resource-group", "location", "force-overwrite"}
	for _, f := range expectedFlags {
		if command.Flags().Lookup(f) == nil {
			t.Fatalf("deploy command should have flag %s", f)
		}
	}

	command.SetArgs([]string{})
	if err := command.Execute(); err == nil {
		t.Fatalf("expected an error when calling deploy with no arguments")
	}
}

func TestValidate(t *testing.T) {
	r := &cobra.Command{}
	cases := []struct {
		dcFactory   func(string) (deployCmd, []string)
		expectedErr error
		args        []string
		name        string
	}{
		{
			dcFactory: func(_ string) (deployCmd, []string) {
				return deployCmd{
					apimodelPath:      "",
					dnsPrefix:         "test",
					outputDirectory:   "output/test",
					forceOverwrite:    false,
					caCertificatePath: "test",
					caPrivateKeyPath:  "test",
					location:          "west europe",
				}, []string{}
			},
			expectedErr: nil,
			name:        "ValidWithNoAPIModelPath",
		},
		{
			dcFactory: func(_ string) (deployCmd, []string) {
				return deployCmd{
					apimodelPath:      "",
					dnsPrefix:         "test",
					outputDirectory:   "output/test",
					caCertificatePath: "test",
					caPrivateKeyPath:  "test",
				}, []string{"wrong/path"}
			},
			expectedErr: errors.New("specified api model does not exist (wrong/path)"),
			name:        "InvalidWithWrongPath",
		},
		{
			dcFactory: func(_ string) (deployCmd, []string) {
				return deployCmd{
					apimodelPath:      "",
					dnsPrefix:         "test",
					outputDirectory:   "output/test",
					caCertificatePath: "test",
					caPrivateKeyPath:  "test",
				}, []string{"test/apimodel.json", "some_random_stuff"}
			},
			expectedErr: errors.New("too many arguments were provided to 'deploy'"),
			name:        "InvalidWithTooManyArguments",
		},
		{
			dcFactory: func(apimodelPath string) (deployCmd, []string) {
				return deployCmd{
					apimodelPath:      "",
					dnsPrefix:         "test",
					outputDirectory:   "output/test",
					caCertificatePath: "test",
					caPrivateKeyPath:  "test",
				}, []string{apimodelPath}
			},
			expectedErr: errors.New("--location must be specified"),
			name:        "InvalidWithNoLocationSpecified",
		},
		{
			dcFactory: func(apimodelPath string) (deployCmd, []string) {
				return deployCmd{
					apimodelPath:      "",
					dnsPrefix:         "test",
					outputDirectory:   "output/test",
					caCertificatePath: "test",
					caPrivateKeyPath:  "test",
					location:          "west europe",
				}, []string{apimodelPath}
			},
			expectedErr: nil,
			name:        "ValidWithAPIModelAsArg",
		},
		{
			dcFactory: func(apimodelPath string) (deployCmd, []string) {
				return deployCmd{
					apimodelPath:      apimodelPath,
					dnsPrefix:         "test",
					outputDirectory:   "output/test",
					caCertificatePath: "test",
					caPrivateKeyPath:  "test",
					location:          "canadaeast",
				}, []string{}
			},
			expectedErr: nil,
			name:        "ValidWithAPIModelInParams",
		},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fileName, del := makeTmpFile(t, "apimodel-unit-test.json")
			defer del()
			dc, args := tc.dcFactory(fileName)
			err := dc.validateArgs(r, args)
			if err != nil && tc.expectedErr != nil {
				if err.Error() != tc.expectedErr.Error() {
					t.Fatalf("expected validate deploy command to return error %s, but instead got %s", tc.expectedErr.Error(), err.Error())
				}
			} else {
				if tc.expectedErr != nil {
					t.Fatalf("expected validate deploy command to return error %s, but instead got no error", tc.expectedErr.Error())
				} else if err != nil {
					t.Fatalf("expected validate deploy command to return no error, but instead got %s", err.Error())
				}
			}
		})
	}
}

func TestAutofillApimodelWithManagedIdentitySkipsCreds(t *testing.T) {
	t.Parallel()

	testAutodeployCredentialHandling(t, true, "", "")
}

func TestAutofillApimodelAllowsPrespecifiedCreds(t *testing.T) {
	t.Parallel()

	testAutodeployCredentialHandling(t, false, "clientID", "clientSecret")
}

func TestAutoSufixWithDnsPrefixInApiModel(t *testing.T) {
	t.Parallel()

	apiloader := &api.Apiloader{
		Translator: nil,
	}

	apimodel := getAPIModel(ExampleAPIModelWithDNSPrefix, false, "clientID", "clientSecret")
	cs, ver, err := apiloader.DeserializeContainerService([]byte(apimodel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}

	outDir, del := makeTmpDir(t)
	defer del()

	deployCmd := &deployCmd{
		apimodelPath:     "./this/is/unused.json",
		outputDirectory:  outDir,
		forceOverwrite:   true,
		location:         "westus",
		autoSuffix:       true,
		containerService: cs,
		apiVersion:       ver,

		client: &armhelpers.MockAKSEngineClient{},
		authProvider: &mockAuthProvider{
			authArgs: &authArgs{},
		},
	}

	err = autofillApimodel(deployCmd)
	if err != nil {
		t.Fatalf("unexpected error autofilling the example apimodel: %s", err)
	}

	if deployCmd.containerService.Properties.MasterProfile.DNSPrefix == "mytestcluster" {
		t.Fatalf("expected %s-{timestampsuffix} but got %s", "mytestcluster", deployCmd.containerService.Properties.MasterProfile.DNSPrefix)
	}

}

func TestAPIModelWithoutServicePrincipalProfileAndClientIdAndSecretInCmd(t *testing.T) {
	t.Parallel()

	apiloader := &api.Apiloader{
		Translator: nil,
	}

	apimodel := getAPIModelWithoutServicePrincipalProfile(false)
	TestClientIDInCmd, err := uuid.Parse("DEC923E3-1EF1-4745-9516-37906D56DEC4")
	if err != nil {
		t.Fatalf("Invalid ClientID in Test: %s", err)
	}

	TestClientSecretInCmd := "DEC923E3-1EF1-4745-9516-37906D56DEC4"

	cs, ver, err := apiloader.DeserializeContainerService([]byte(apimodel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}

	outDir, del := makeTmpDir(t)
	defer del()

	deployCmd := &deployCmd{
		apimodelPath:     "./this/is/unused.json",
		outputDirectory:  outDir,
		forceOverwrite:   true,
		location:         "westus",
		containerService: cs,
		apiVersion:       ver,

		client: &armhelpers.MockAKSEngineClient{},
		authProvider: &mockAuthProvider{
			authArgs: &authArgs{},
		},
	}
	deployCmd.getAuthArgs().ClientID = TestClientIDInCmd
	deployCmd.getAuthArgs().ClientSecret = TestClientSecretInCmd

	err = autofillApimodel(deployCmd)
	if err != nil {
		t.Fatalf("unexpected error autofilling the example apimodel: %s", err)
	}

	if deployCmd.containerService.Properties.ServicePrincipalProfile == nil || deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID == "" || deployCmd.containerService.Properties.ServicePrincipalProfile.Secret == "" {
		t.Fatalf("expected service principal profile to be populated from deployment command arguments")
	}

	if deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID != TestClientIDInCmd.String() {
		t.Fatalf("expected service principal profile client id to be %s but got %s", TestClientIDInCmd.String(), deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID)
	}

	if deployCmd.containerService.Properties.ServicePrincipalProfile.Secret != TestClientSecretInCmd {
		t.Fatalf("expected service principal profile client secret to be %s but got %s", TestClientSecretInCmd, deployCmd.containerService.Properties.ServicePrincipalProfile.Secret)
	}
}

func TestAPIModelWithEmptyServicePrincipalProfileAndClientIdAndSecretInCmd(t *testing.T) {
	t.Parallel()

	apiloader := &api.Apiloader{
		Translator: nil,
	}

	apimodel := getAPIModel(ExampleAPIModelWithDNSPrefix, false, "", "")
	TestClientIDInCmd, err := uuid.Parse("DEC923E3-1EF1-4745-9516-37906D56DEC4")
	if err != nil {
		t.Fatalf("Invalid ClientID in Test: %s", err)
	}

	TestClientSecretInCmd := "DEC923E3-1EF1-4745-9516-37906D56DEC4"

	cs, ver, err := apiloader.DeserializeContainerService([]byte(apimodel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}

	outDir, del := makeTmpDir(t)
	defer del()

	deployCmd := &deployCmd{
		apimodelPath:     "./this/is/unused.json",
		outputDirectory:  outDir,
		forceOverwrite:   true,
		location:         "westus",
		containerService: cs,
		apiVersion:       ver,

		client: &armhelpers.MockAKSEngineClient{},
		authProvider: &mockAuthProvider{
			authArgs: &authArgs{},
		},
	}
	deployCmd.getAuthArgs().ClientID = TestClientIDInCmd
	deployCmd.getAuthArgs().ClientSecret = TestClientSecretInCmd
	err = autofillApimodel(deployCmd)
	if err != nil {
		t.Fatalf("unexpected error autofilling the example apimodel: %s", err)
	}

	if deployCmd.containerService.Properties.ServicePrincipalProfile == nil || deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID == "" || deployCmd.containerService.Properties.ServicePrincipalProfile.Secret == "" {
		t.Fatalf("expected service principal profile to be populated from deployment command arguments")
	}

	if deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID != TestClientIDInCmd.String() {
		t.Fatalf("expected service principal profile client id to be %s but got %s", TestClientIDInCmd.String(), deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID)
	}

	if deployCmd.containerService.Properties.ServicePrincipalProfile.Secret != TestClientSecretInCmd {
		t.Fatalf("expected service principal profile client secret to be %s but got %s", TestClientSecretInCmd, deployCmd.containerService.Properties.ServicePrincipalProfile.Secret)
	}
}

func TestAPIModelWithoutServicePrincipalProfileAndWithoutClientIdAndSecretInCmd(t *testing.T) {
	t.Parallel()

	apiloader := &api.Apiloader{
		Translator: nil,
	}

	apimodel := getAPIModelWithoutServicePrincipalProfile(false)

	cs, ver, err := apiloader.DeserializeContainerService([]byte(apimodel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}

	outDir, del := makeTmpDir(t)
	defer del()

	deployCmd := &deployCmd{
		apimodelPath:     "./this/is/unused.json",
		outputDirectory:  outDir,
		forceOverwrite:   true,
		location:         "westus",
		containerService: cs,
		apiVersion:       ver,

		client: &armhelpers.MockAKSEngineClient{},
		authProvider: &mockAuthProvider{
			authArgs: &authArgs{},
		},
	}
	err = autofillApimodel(deployCmd)
	if err != nil {
		t.Fatalf("unexpected error autofilling the example apimodel: %s", err)
	}

	if deployCmd.containerService.Properties.ServicePrincipalProfile != nil {
		t.Fatalf("expected service principal profile to be nil for unmanaged identity, where client id and secret are not supplied in api model and deployment command")
	}

}

func testAutodeployCredentialHandling(t *testing.T, useManagedIdentity bool, clientID, clientSecret string) {
	apiloader := &api.Apiloader{
		Translator: nil,
	}

	apimodel := getExampleAPIModel(useManagedIdentity, clientID, clientSecret)
	cs, ver, err := apiloader.DeserializeContainerService([]byte(apimodel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}

	outDir, del := makeTmpDir(t)
	defer del()

	// deserialization happens in validate(), but we are testing just the default
	// setting that occurs in autofillApimodel (which is called from validate)
	// Thus, it assumes that containerService/apiVersion are already populated
	deployCmd := &deployCmd{
		apimodelPath:    "./this/is/unused.json",
		dnsPrefix:       "dnsPrefix1",
		outputDirectory: outDir,
		forceOverwrite:  true,
		location:        "westus",

		containerService: cs,
		apiVersion:       ver,

		client: &armhelpers.MockAKSEngineClient{},
		authProvider: &mockAuthProvider{
			authArgs: &authArgs{},
		},
	}

	err = autofillApimodel(deployCmd)
	if err != nil {
		t.Fatalf("unexpected error autofilling the example apimodel: %s", err)
	}

	err = deployCmd.validateAPIModelAsVLabs()
	if err != nil {
		t.Fatalf("unexpected error validating apimodel after populating defaults: %s", err)
	}

	if useManagedIdentity {
		if cs.Properties.ServicePrincipalProfile != nil &&
			(cs.Properties.ServicePrincipalProfile.ClientID != "" || cs.Properties.ServicePrincipalProfile.Secret != "") {
			t.Fatalf("Unexpected credentials were populated even though MSI was active.")
		}
	} else {
		if cs.Properties.ServicePrincipalProfile == nil ||
			cs.Properties.ServicePrincipalProfile.ClientID == "" || cs.Properties.ServicePrincipalProfile.Secret == "" {
			t.Fatalf("Credentials were missing even though MSI was not active.")
		}
	}
}

func TestDeployCmdMergeAPIModel(t *testing.T) {
	t.Parallel()

	d := &deployCmd{}
	d.apimodelPath = "../pkg/engine/testdata/simple/kubernetes.json"
	err := d.mergeAPIModel()
	if err != nil {
		t.Fatalf("unexpected error calling mergeAPIModel with no --set flag defined: %s", err.Error())
	}

	d = &deployCmd{}
	d.apimodelPath = "../pkg/engine/testdata/simple/kubernetes.json"
	d.set = []string{"masterProfile.count=3,linuxProfile.adminUsername=testuser"}
	err = d.mergeAPIModel()
	if err != nil {
		t.Fatalf("unexpected error calling mergeAPIModel with one --set flag: %s", err.Error())
	}

	d = &deployCmd{}
	d.apimodelPath = "../pkg/engine/testdata/simple/kubernetes.json"
	d.set = []string{"masterProfile.count=3", "linuxProfile.adminUsername=testuser"}
	err = d.mergeAPIModel()
	if err != nil {
		t.Fatalf("unexpected error calling mergeAPIModel with multiple --set flags: %s", err.Error())
	}

	d = &deployCmd{}
	d.apimodelPath = "../pkg/engine/testdata/simple/kubernetes.json"
	d.set = []string{"agentPoolProfiles[0].count=1"}
	err = d.mergeAPIModel()
	if err != nil {
		t.Fatalf("unexpected error calling mergeAPIModel with one --set flag to override an array property: %s", err.Error())
	}
}

func TestDeployCmdRun(t *testing.T) {
	t.Parallel()

	outdir, del := makeTmpDir(t)
	defer del()

	d := &deployCmd{
		client: &armhelpers.MockAKSEngineClient{},
		authProvider: &mockAuthProvider{
			authArgs:      &authArgs{},
			getClientMock: &armhelpers.MockAKSEngineClient{},
		},
		apimodelPath:    "../pkg/engine/testdata/simple/kubernetes.json",
		outputDirectory: outdir,
		forceOverwrite:  true,
		location:        "westus",
	}

	r := &cobra.Command{}
	f := r.Flags()

	addAuthFlags(d.getAuthArgs(), f)

	fakeRawSubscriptionID := "6dc93fae-9a76-421f-bbe5-cc6460ea81cb"
	fakeSubscriptionID, err := uuid.Parse(fakeRawSubscriptionID)
	fakeClientID := "b829b379-ca1f-4f1d-91a2-0d26b244680d"
	fakeClientSecret := "0se43bie-3zs5-303e-aav5-dcf231vb82ds"
	if err != nil {
		t.Fatalf("Invalid SubscriptionId in Test: %s", err)
	}

	d.getAuthArgs().SubscriptionID = fakeSubscriptionID
	d.getAuthArgs().rawSubscriptionID = fakeRawSubscriptionID
	d.getAuthArgs().rawClientID = fakeClientID
	d.getAuthArgs().ClientSecret = fakeClientSecret

	err = d.loadAPIModel()
	if err != nil {
		t.Fatalf("Failed to call LoadAPIModel: %s", err)
	}

	err = d.run()
	if err != nil {
		t.Fatalf("Failed to call LoadAPIModel: %s", err)
	}
}

func TestDeployCmdWithoutMasterProfile(t *testing.T) {
	t.Parallel()

	outdir, del := makeTmpDir(t)
	defer del()

	d := &deployCmd{
		client: &armhelpers.MockAKSEngineClient{},
		authProvider: &mockAuthProvider{
			authArgs:      &authArgs{},
			getClientMock: &armhelpers.MockAKSEngineClient{},
		},
		apimodelPath:    "../pkg/engine/testdata/simple/kubernetes.json",
		outputDirectory: outdir,
		forceOverwrite:  true,
		location:        "westus",
	}
	d.set = []string{"masterProfile=nil"}

	err := d.loadAPIModel()
	if err == nil {
		t.Fatalf("expected error loading api model without MasterProfile: %s", err.Error())
	}
}

func TestLoadApiModelOnCustomCloud(t *testing.T) {
	t.Parallel()

	outdir, del := makeTmpDir(t)
	defer del()

	d := &deployCmd{
		client: &armhelpers.MockAKSEngineClient{},
		authProvider: &mockAuthProvider{
			authArgs:      &authArgs{},
			getClientMock: &armhelpers.MockAKSEngineClient{},
		},
		apimodelPath:    "../pkg/engine/testdata/customcloud/kubernetes.json",
		outputDirectory: outdir,
		forceOverwrite:  true,
		location:        "westus",
	}

	r := &cobra.Command{}
	f := r.Flags()

	addAuthFlags(d.getAuthArgs(), f)

	fakeRawSubscriptionID := "6dc93fae-9a76-421f-bbe5-cc6460ea81cb"
	fakeSubscriptionID, err := uuid.Parse(fakeRawSubscriptionID)
	fakeClientID := "b829b379-ca1f-4f1d-91a2-0d26b244680d"
	fakeClientSecret := "0se43bie-3zs5-303e-aav5-dcf231vb82ds"
	if err != nil {
		t.Fatalf("Invalid SubscriptionId in Test: %s", err)
	}

	d.getAuthArgs().SubscriptionID = fakeSubscriptionID
	d.getAuthArgs().rawSubscriptionID = fakeRawSubscriptionID
	d.getAuthArgs().rawClientID = fakeClientID
	d.getAuthArgs().ClientSecret = fakeClientSecret
	err = d.loadAPIModel()
	if err != nil {
		t.Fatalf("Failed to call LoadAPIModel: %s", err)
	}
}

func TestLoadApiModelOnAzureStack(t *testing.T) {
	t.Parallel()

	outdir, del := makeTmpDir(t)
	defer del()

	d := &deployCmd{
		client: &armhelpers.MockAKSEngineClient{},
		authProvider: &mockAuthProvider{
			authArgs:      &authArgs{},
			getClientMock: &armhelpers.MockAKSEngineClient{},
		},
		apimodelPath:    "../pkg/engine/testdata/azurestack/kubernetes.json",
		outputDirectory: outdir,
		forceOverwrite:  true,
		location:        "westus",
	}

	r := &cobra.Command{}
	f := r.Flags()

	addAuthFlags(d.getAuthArgs(), f)

	d.location = "local"
	fakeRawSubscriptionID := "6dc93fae-9a76-421f-bbe5-cc6460ea81cb"
	fakeSubscriptionID, err := uuid.Parse(fakeRawSubscriptionID)
	fakeClientID := "b829b379-ca1f-4f1d-91a2-0d26b244680d"
	fakeClientSecret := "0se43bie-3zs5-303e-aav5-dcf231vb82ds"
	if err != nil {
		t.Fatalf("Invalid SubscriptionId in Test: %s", err)
	}

	d.getAuthArgs().IdentitySystem = "adfs"
	d.getAuthArgs().SubscriptionID = fakeSubscriptionID
	d.getAuthArgs().rawSubscriptionID = fakeRawSubscriptionID
	d.getAuthArgs().rawClientID = fakeClientID
	d.getAuthArgs().ClientSecret = fakeClientSecret
	err = d.loadAPIModel()
	if err != nil {
		t.Fatalf("Failed to call LoadAPIModel: %s", err)
	}

	if d.getAuthArgs().IdentitySystem != d.containerService.Properties.CustomCloudProfile.IdentitySystem {
		t.Fatal("Failed to set cli Identity system as default ")
	}
}

func TestLoadApiModelOnAzureStackWithAKSUbuntu1804Distro(t *testing.T) {
	t.Parallel()

	outdir, del := makeTmpDir(t)
	defer del()

	d := &deployCmd{
		client: &armhelpers.MockAKSEngineClient{},
		authProvider: &mockAuthProvider{
			authArgs:      &authArgs{},
			getClientMock: &armhelpers.MockAKSEngineClient{},
		},
		apimodelPath:    "../pkg/engine/testdata/azurestack/kubernetes-ubuntu18.json",
		outputDirectory: outdir,
		forceOverwrite:  true,
		location:        "westus",
	}

	r := &cobra.Command{}
	f := r.Flags()

	addAuthFlags(d.getAuthArgs(), f)

	d.location = "local"
	fakeRawSubscriptionID := "6dc93fae-9a76-421f-bbe5-cc6460ea81cb"
	fakeSubscriptionID, err := uuid.Parse(fakeRawSubscriptionID)
	fakeClientID := "b829b379-ca1f-4f1d-91a2-0d26b244680d"
	fakeClientSecret := "0se43bie-3zs5-303e-aav5-dcf231vb82ds"
	if err != nil {
		t.Fatalf("Invalid SubscriptionId in Test: %s", err)
	}

	d.getAuthArgs().IdentitySystem = "adfs"
	d.getAuthArgs().SubscriptionID = fakeSubscriptionID
	d.getAuthArgs().rawSubscriptionID = fakeRawSubscriptionID
	d.getAuthArgs().rawClientID = fakeClientID
	d.getAuthArgs().ClientSecret = fakeClientSecret
	err = d.loadAPIModel()
	if err != nil && err.Error() != errors.New("invalid master profile distro 'aks-ubuntu-18.04'").Error() {
		t.Fatalf("Failed to identify invalid distro 'aks-ubuntu-18.04': %s", err)
	}
}

func TestOutputDirectoryWithDNSPrefix(t *testing.T) {
	t.Parallel()

	apiloader := &api.Apiloader{
		Translator: nil,
	}

	apimodel := getAPIModel(ExampleAPIModelWithoutDNSPrefix, false, "clientID", "clientSecret")
	cs, ver, err := apiloader.DeserializeContainerService([]byte(apimodel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}

	d := &deployCmd{
		apimodelPath:     "./this/is/unused.json",
		outputDirectory:  "",
		dnsPrefix:        "dnsPrefix1",
		forceOverwrite:   true,
		location:         "westus",
		containerService: cs,
		apiVersion:       ver,
		client:           &armhelpers.MockAKSEngineClient{},
		authProvider: &mockAuthProvider{
			authArgs: &authArgs{},
		},
	}

	err = autofillApimodel(d)
	if err != nil {
		t.Fatalf("unexpected error autofilling the example apimodel: %s", err)
	}

	defer os.RemoveAll("_output")

	if d.outputDirectory != path.Join("_output", d.dnsPrefix) {
		t.Fatalf("Calculated output directory should be %s, actual value %s", path.Join("_output", d.dnsPrefix), d.outputDirectory)
	}
}

func makeTmpFile(t *testing.T, name string) (string, func()) {
	tmpF, err := os.CreateTemp(os.TempDir(), name)
	if err != nil {
		t.Fatalf("unable to create file: %s", err.Error())
	}

	return tmpF.Name(), func() {
		defer os.Remove(tmpF.Name())
	}
}
