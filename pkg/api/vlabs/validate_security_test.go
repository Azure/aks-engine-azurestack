// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package vlabs

import (
	"testing"
)

// TestValidateURLForSSRF_DangerousHostnames tests SSRF protection
func TestValidateURLForSSRF_DangerousHostnames(t *testing.T) {
	dangerousURLs := []string{
		"https://localHoSt/",
		"https://127.0.1.2/",
		"https://0177.0.23.19/", // Octal representation
		"https://2130706433/",   // Decimal representation
		"https://0x7f.00331.0246.174/",
		"https://[::1]/",
		"https://[fc00::]/",
		"https://169.254.169.254/", // AWS/Azure metadata service
		"http://localhost/test",
		"http://127.0.0.1/api",
		"http://0.0.0.0/",
		"https://[0:0:0:0:0:0:0:1]/",
	}

	for _, url := range dangerousURLs {
		t.Run(url, func(t *testing.T) {
			err := validateURLForSSRF(url)
			if err == nil {
				t.Errorf("Expected validation to fail for dangerous URL: %s", url)
			}
		})
	}
}

// TestValidateURLForSSRF_SafeHostnames tests legitimate URLs are allowed
func TestValidateURLForSSRF_SafeHostnames(t *testing.T) {
	safeURLs := []string{
		"https://vault.azure.net/",
		"https://example.com/path",
		"https://subdomain.example.com/",
		"https://management.azure.com/",
	}

	for _, url := range safeURLs {
		t.Run(url, func(t *testing.T) {
			err := validateURLForSSRF(url)
			if err != nil {
				t.Errorf("Expected validation to pass for safe URL: %s, got error: %v", url, err)
			}
		})
	}
}

// TestContainsPathTraversal_DangerousPaths tests path traversal detection
func TestContainsPathTraversal_DangerousPaths(t *testing.T) {
	dangerousPaths := []string{
		"/../../OtherPath/",
		"/..//OtherPath/",
		"/%2E%2E%2f/OtherPath/",
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32",
		"%2e%2e%2ftest",
		"%2E%2E%2Ftest",
		"%252e%252e%252ftest",
		"..%2ftest",
		"..%5ctest",
	}

	for _, path := range dangerousPaths {
		t.Run(path, func(t *testing.T) {
			result := containsPathTraversal(path)
			if !result {
				t.Errorf("Expected path traversal detection to return true for: %s", path)
			}
		})
	}
}

// TestContainsPathTraversal_SafePaths tests legitimate paths are allowed
func TestContainsPathTraversal_SafePaths(t *testing.T) {
	safePaths := []string{
		"/subscriptions/abc123/resourceGroups/myRG",
		"/test/path/to/resource",
		"mySecretName",
		"valid-secret-name-123",
	}

	for _, path := range safePaths {
		t.Run(path, func(t *testing.T) {
			result := containsPathTraversal(path)
			if result {
				t.Errorf("Expected path traversal detection to return false for safe path: %s", path)
			}
		})
	}
}

// TestDecodeURLRecursively tests multi-level URL decoding
func TestDecodeURLRecursively(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "Single encoding",
			input:    "%2E%2E%2F",
			expected: "../",
			hasError: false,
		},
		{
			name:     "Double encoding",
			input:    "%252E%252E%252F",
			expected: "../",
			hasError: false,
		},
		{
			name:     "No encoding",
			input:    "test",
			expected: "test",
			hasError: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := decodeURLRecursively(tt.input)
			if tt.hasError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestIsDangerousHostname(t *testing.T) {
	tests := []struct {
		hostname   string
		isDangerous bool
	}{
		{"localhost", true},
		{"LocalHost", true},
		{"LOCALHOST", true},
		{"127.0.0.1", true},
		{"127.0.1.1", true},
		{"169.254.169.254", true},
		{"::1", true},
		{"[::1]", true},
		{"0.0.0.0", true},
		{"example.com", false},
		{"vault.azure.net", false},
		{"192.168.1.1", true},
		{"10.0.0.1", true},
		{"172.16.0.1", true},
	}

	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			result := isDangerousHostname(tt.hostname)
			if result != tt.isDangerous {
				t.Errorf("For hostname %s: expected isDangerous=%v, got %v", tt.hostname, tt.isDangerous, result)
			}
		})
	}
}

func TestValidateStringLength(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		maxLength int
		hasError  bool
	}{
		{
			name:      "Within limit",
			value:     "short",
			fieldName: "TestField",
			maxLength: 10,
			hasError:  false,
		},
		{
			name:      "At limit",
			value:     "exactlyten",
			fieldName: "TestField",
			maxLength: 10,
			hasError:  false,
		},
		{
			name:      "Exceeds limit",
			value:     "this is way too long",
			fieldName: "TestField",
			maxLength: 10,
			hasError:  true,
		},
		{
			name:      "Empty string",
			value:     "",
			fieldName: "TestField",
			maxLength: 10,
			hasError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStringLength(tt.value, tt.fieldName, tt.maxLength)
			if tt.hasError && err == nil {
				t.Errorf("Expected error for value length %d with max %d", len(tt.value), tt.maxLength)
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidateKeyVaultSecretsWithSSRF(t *testing.T) {
	tests := []struct {
		name                     string
		secrets                  []KeyVaultSecrets
		requireCertificateStore bool
		expectError              bool
	}{
		{
			name: "Valid KeyVault URL",
			secrets: []KeyVaultSecrets{
				{
					SourceVault: &KeyVaultID{
						ID: "/subscriptions/sub123/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/myvault",
					},
					VaultCertificates: []KeyVaultCertificate{
						{
							CertificateURL:   "https://myvault.vault.azure.net/secrets/mycert",
							CertificateStore: "My",
						},
					},
				},
			},
			requireCertificateStore: true,
			expectError:              false,
		},
		{
			name: "Dangerous hostname in vault URL",
			secrets: []KeyVaultSecrets{
				{
					SourceVault: &KeyVaultID{
						ID: "/subscriptions/sub123/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/localhost",
					},
					VaultCertificates: []KeyVaultCertificate{
						{
							CertificateURL:   "https://localhost/secrets/mycert",
							CertificateStore: "My",
						},
					},
				},
			},
			requireCertificateStore: true,
			expectError:              true,
		},
		{
			name: "Path traversal in certificate URL",
			secrets: []KeyVaultSecrets{
				{
					SourceVault: &KeyVaultID{
						ID: "/subscriptions/sub123/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/myvault",
					},
					VaultCertificates: []KeyVaultCertificate{
						{
							CertificateURL:   "https://vault.azure.net/../../secrets/mycert",
							CertificateStore: "My",
						},
					},
				},
			},
			requireCertificateStore: true,
			expectError:              true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateKeyVaultSecrets(tt.secrets, tt.requireCertificateStore)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
