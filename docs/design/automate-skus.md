# Azure VM SKU and Locations Automation

This document proposes a design to automate the updating of some data structures used to represent Azure constants in AKS Engine.

## Problem

AKS Engine represents some information about Azure in code constants or similar read-only data structures. Specifically, Azure data center location names and Virtual Machine SKU details are needed for validation and for enabling certain VM features.

To stay in sync with Azure, these data structures must be refreshed periodically. To date, that has been a task for maintainers to perform on-demand. AKS Engine has often lagged behind Azure's introduction of new locations or VM SKUs.

A python script will regenerate some, but not all of the data structures.

## Goals

- On a regular schedule, notify the AKS Engine team if new locations or SKUs exist
- Automate the creation of the pull request with the new locations or SKUs

## Non-goals

## Proposal: scriptable aks-engine-azurestack commands

New developer-only commands in the `aks-engine-azurestack` binary generate the necessary Go source code. The commands replace the existing python script and extend it to cover the accelerated networking capability list.

### Example locations output

```shell
$ aks-engine-azurestack get-locations --client-id=${AZURE_CLIENT_ID} --client-secret=${AZURE_CLIENT_SECRET}
Location            Name                      Latitude    Longitude
australiacentral    Australia Central         -35.3075    149.1244
australiacentral2   Australia Central 2       -35.3075    149.1244
centraluseuap       Central US EUAP (Canary)  N/A         N/A
chinaeast           China East                N/A         N/A
...
$ aks-engine-azurestack get-locations -o code
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package helpers

// AUTOGENERATED FILE

// GetAzureLocations provides all azure regions in prod.
// Related powershell to refresh this list:
//   Get-AzureRmLocation | Select-Object -Property Location
func GetAzureLocations() []string {
	return []string{
		"australiacentral",
		"australiacentral2",
		"centraluseuap",
		"chinaeast",
...

```

### Example SKUs output

```shell
$ aks-engine-azurestack get-skus -o human
Name                 Storage Account Type  Accelerated Networking Support
Standard_A0          Standard_LRS          true
Standard_A1          Standard_LRS          true
Standard_D8s_v3      Premium_LRS           true
Standard_DC1s_v2     Premium_LRS           false
...
$ aks-engine-azurestack get-skus -o code
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package helpers

// AUTOGENERATED FILE

// GetKubernetesAllowedVMSKUs returns the allowed sizes for Kubernetes agent
func GetKubernetesAllowedVMSKUs() string {
	return `      "allowedValues": [
		"Standard_A0",
		"Standard_A1",
		"Standard_D8s_v3",
		"Standard_DC1s_v2",
...
```

First we refactor the code so the data structures are in separate files. Then each file can be overwritten by redirecting the output of an `aks-engine-azurestack` command.

### Pros

- No change to existing data structures, just a new way to refresh them
- Keeps all knowledge of these Azure details in AKS Engine
- Modest implementation effort

### Cons

- Hidden developer-facing commands aren't the best UX
- Still requires a human to oversee the process

## Alternative implementations considered

### Dynamic API lookups

Instead of consulting a lookup table in code, AKS Engine could make Azure API calls just-in-time. The location and SKU-related details of a cluster could be validated in real time by `aks-engine-azurestack generate` or `aks-engine-azurestack deploy` against current Azure information.

#### Pros

- Location and SKU support data is always current
- Live Azure as the one source of truth

#### Cons

- Requires significant refactoring and implementation effort
- Making API calls during `aks-engine-azurestack generate` may be unwanted behavior
- Somewhat slower and less reliable than using the data in code
- Requires implementing most of the main proposal above (as a starting point)

### Improve the existing python script

The existing python script is crusty but functional. It could be extended to cover generating the remaining code.

#### Pros

- Modest implementation effort

#### Cons

- Having logic that's not in Go code makes it harder to locate and to grok
