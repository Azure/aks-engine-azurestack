// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/pkg/errors"
)

// NewDefaultCredential returns an AzureClient
func NewDefaultCredential(cloud cloud.Configuration, subscriptionID string) (*azidentity.DefaultAzureCredential, error) {
	options := &azidentity.DefaultAzureCredentialOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: cloud,
		},
	}
	cred, err := azidentity.NewDefaultAzureCredential(options)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

// NewClientSecretCredential returns an AzureClient via client_id and client_secret
func NewClientSecretCredential(cloud cloud.Configuration, subscriptionID, clientID, clientSecret string) (*azidentity.ClientSecretCredential, error) {
	tenantID, err := getOAuthConfig(subscriptionID, cloud)
	if err != nil {
		return nil, err
	}
	options := &azidentity.ClientSecretCredentialOptions{
		DisableInstanceDiscovery: true,
		ClientOptions: policy.ClientOptions{
			Cloud: cloud,
		},
	}
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, options)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

// NewClientSecretCredentialExternalTenant returns an AzureClient via client_id and client_secret from a tenant
func NewClientSecretCredentialExternalTenant(cloud cloud.Configuration, subscriptionID, clientID, clientSecret string) (*azidentity.ClientSecretCredential, error) {
	options := &azidentity.ClientSecretCredentialOptions{
		DisableInstanceDiscovery: true,
		ClientOptions: policy.ClientOptions{
			Cloud: cloud,
		},
	}
	cred, err := azidentity.NewClientSecretCredential("adfs", clientID, clientSecret, options)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

// NewClientCertificateCredential returns an AzureClient via client_id and jwt certificate assertion
func NewClientCertificateCredential(cloud cloud.Configuration, subscriptionID, clientID, certificatePath, privateKeyPath string) (*azidentity.ClientCertificateCredential, error) {
	certificateData, err := os.ReadFile(certificatePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read certificate")
	}
	block, _ := pem.Decode(certificateData)
	if block == nil {
		return nil, errors.New("Failed to decode pem block from certificate")
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse certificate")
	}
	privateKey, err := parseRsaPrivateKey(privateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse rsa private key")
	}
	tenantID, err := getOAuthConfig(subscriptionID, cloud)
	if err != nil {
		return nil, err
	}
	options := &azidentity.ClientCertificateCredentialOptions{
		SendCertificateChain:     true,
		DisableInstanceDiscovery: true,
		ClientOptions: policy.ClientOptions{
			Cloud: cloud,
		},
	}
	return azidentity.NewClientCertificateCredential(tenantID, clientID, []*x509.Certificate{certificate}, privateKey, options)
}

// NewClientCertificateCredentialExternalTenant returns an AzureClient via client_id and jwt certificate assertion a 3rd party tenant
func NewClientCertificateCredentialExternalTenant(cloud cloud.Configuration, subscriptionID, clientID, certificatePath, privateKeyPath string) (*azidentity.ClientCertificateCredential, error) {
	certificateData, err := os.ReadFile(certificatePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read certificate")
	}
	block, _ := pem.Decode(certificateData)
	if block == nil {
		return nil, errors.New("Failed to decode pem block from certificate")
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse certificate")
	}
	privateKey, err := parseRsaPrivateKey(privateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse rsa private key")
	}
	options := &azidentity.ClientCertificateCredentialOptions{
		SendCertificateChain:     true,
		DisableInstanceDiscovery: true,
		ClientOptions: policy.ClientOptions{
			Cloud: cloud,
		},
	}
	return azidentity.NewClientCertificateCredential("adfs", clientID, []*x509.Certificate{certificate}, privateKey, options)
}

func getOAuthConfig(subscriptionID string, cloud cloud.Configuration) (string, error) {
	tenantID, err := GetTenantID(subscriptionID, cloud)
	if err != nil {
		return "", err
	}
	return tenantID, nil
}

func parseRsaPrivateKey(path string) (*rsa.PrivateKey, error) {
	privateKeyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return nil, errors.New("Failed to decode a pem block from private key")
	}

	privatePkcs1Key, errPkcs1 := x509.ParsePKCS1PrivateKey(block.Bytes)
	if errPkcs1 == nil {
		return privatePkcs1Key, nil
	}

	privatePkcs8Key, errPkcs8 := x509.ParsePKCS8PrivateKey(block.Bytes)
	if errPkcs8 == nil {
		privatePkcs8RsaKey, ok := privatePkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("pkcs8 contained non-RSA key. Expected RSA key")
		}
		return privatePkcs8RsaKey, nil
	}

	return nil, errors.Errorf("failed to parse private key as Pkcs#1 or Pkcs#8. (%s). (%s)", errPkcs1, errPkcs8)
}
