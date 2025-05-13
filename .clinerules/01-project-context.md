# AKS Engine - Units of Kubernetes on Azure Stack Hub

## Overview

AKS Engine is an ARM template-driven way to provision a self-managed Kubernetes cluster on Azure Stack Hub. By leveraging [ARM (Azure Resource Manager)][ARM], AKS Engine helps you create, destroy and maintain clusters provisioned with basic IaaS resources in Azure Stack Hub.

AKS Engine is designed to be used as a CLI tool (`aks-engine-azurestack`). This document outlines the functionality that `aks-engine-azurestack` provides to create and maintain a Kubernetes cluster on Azure.

## `aks-engine-azurestack` commands

To get a quick overview of the commands available via the `aks-engine-azurestack` CLI tool, just run `aks-engine-azurestack` with no arguments (or include the `--help` argument):

```sh
$ aks-engine
Usage:
  aks-engine-azurestack [flags]
  aks-engine-azurestack [command]

Available Commands:
  addpool       Add a node pool to an existing AKS Engine-created Kubernetes cluster
  completion    Generates bash completion scripts
  deploy        Deploy an Azure Resource Manager template
  generate      Generate an Azure Resource Manager template
  get-logs      Collect logs and current cluster nodes configuration.
  get-versions  Display info about supported Kubernetes versions
  help          Help about any command
  rotate-certs  (experimental) Rotate certificates on an existing AKS Engine-created Kubernetes cluster
  scale         Scale an existing AKS Engine-created Kubernetes cluster
  upgrade       Upgrade an existing AKS Engine-created Kubernetes cluster
  version       Print the version of aks-engine

Flags:
      --debug                enable verbose debug logs
  -h, --help                 help for aks-engine
      --show-default-model   Dump the default API model to stdout

Use "aks-engine-azurestack [command] --help" for more information about a command.
```

# Technical Stack Documentation

## Programming Language

The AKS engine is developed using Go (Golang).

## Testing Framework

Ginkgo is used as the primary testing framework for the project.

## Overall Process

1. **Command-Line Interface Development**  
   The CLI is developed using `github.com/spf13/cobra`.

2. **API Model Loading and Deserialization**

   - Load the API model from a JSON file.
   - Deserialize the JSON into objects defined in `pkg/api/types.go`.

3. **Default Value Assignment**

   - Set default values for the deserialized API objects as needed.

4. **ARM Template Generation**

   - Use `text/template` with custom functions to generate Azure Resource Manager (ARM) templates for deploying, scaling, or upgrading a Kubernetes cluster.

5. **Deployment**
   - Deploy the generated ARM template to provision or update the Kubernetes cluster.
