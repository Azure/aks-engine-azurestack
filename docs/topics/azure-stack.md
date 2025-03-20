# AKS Engine on Azure Stack Hub

* [Introduction](#introduction)
* [Marketplace Prerequisites](#marketplace-prerequisites)
* [Service Principals and Identity Providers](#service-principals-and-identity-providers)
* [CLI flags](#cli-flags)
* [Cluster Definition (aka API Model)](#cluster-definition-aka-api-model)
  * [location](#location)
  * [kubernetesConfig](#kubernetesConfig)
  * [customCloudProfile](#customCloudProfile)
  * [masterProfile](#masterProfile)
  * [agentPoolProfiles](#agentPoolProfiles)
* [Azure Stack Hub Instances Registered with Azure's China cloud](#azure-stack-hub-instances-registered-with-azures-china-cloud)
* [Disconnected Azure Stack Hub Instances](#disconnected-azure-stack-hub-instances)
* [AKS Engine Versions](#aks-engine-versions)
* [Cloud Provider for Azure](#cloud-provider-for-azure)
* [Volume Provisioner: Container Storage Interface Drivers (preview)](#volume-provisioner-container-storage-interface-drivers-preview)
* [Known Issues and Limitations](#known-issues-and-limitations)
* [Frequently Asked Questions](#frequently-asked-questions)

## Introduction

Specific AKS Engine [versions](#aks-engine-versions) can be used to provision self-managed Kubernetes clusters on [Azure Stack Hub](https://azure.microsoft.com/overview/azure-stack/). AKS Engine's `generate`, [deploy](../tutorials/quickstart.md#deploy), [upgrade](upgrade.md), and [scale](scale.md) commands can be executed as if you were targeting Azure's public cloud. You are only required to slightly update your cluster definition to provide some extra information about your Azure Stack Hub instance.

The goal of this guide is to explain how to provision Kubernetes clusters to Azure Stack Hub using AKS Engine and to capture the differences between Azure and Azure Stack Hub. Bear in mind as well that not every AKS Engine feature or configuration option is currently supported on Azure Stack Hub. In most cases, these are not available because dependent Azure components are not part of Azure Stack Hub.

## Marketplace prerequisites

Because Azure Stack Hub instances do not have infinite storage available, Azure Stack Hub administrators are in charge of managing it by selecting which marketplace items are downloaded from Azure's marketplace. The Azure Stack Hub administrator can follow this [guide](https://docs.microsoft.com/azure-stack/operator/azure-stack-download-azure-marketplace-item) for a general explanation about how to download marketplace items from Azure.

Before you try to deploy the first Kubernetes cluster, make sure these marketplace items were made available to the target subscription by the Azure Stack Hub administrator.

* `Custom Script for Linux 2.0` virtual machine extension
* [Required](#aks-engine-versions) `AKS Base Image` virtual machine

## Service Principals and Identity Providers

Kubernetes uses a `service principal` identity to talk to Azure Stack Hub APIs to dynamically manage resources such as storage or load balancers. Therefore, you will need to create a service principal before you can provision a Kubernetes cluster using AKS Engine.

This [guide](https://docs.microsoft.com/azure-stack/operator/azure-stack-create-service-principals) explains how to create and manage service principals on Azure Stack Hub for both Azure Active Directory (AAD) and Active Directory Federation Services (ADFS) identity providers. This other [guide](../../docs/topics/service-principals.md) is a good resource to understand the permissions that the service principal requires to deploy under your subscription.

Once you have created the required service principal, make sure to assign it the `contributor` role at the target subscription scope.

## CLI flags

To indicate to AKS Engine that your target platform is Azure Stack Hub, all commands require CLI flag `azure-env` to be set to `"AzureStackCloud"`.

If your Azure Stack Hub instance uses ADFS to authenticate identities, then flag `identity-system` is also required.

``` bash
aks-engine-azurestack deploy \
  --azure-env AzureStackCloud \
  --api-model kubernetes.json \
  --location local \
  --resource-group kube-rg \
  --identity-system adfs \
  --client-id $SPN_CLIENT_ID \
  --client-secret $SPN_CLIENT_SECRET \
  --subscription-id $TENANT_SUBSCRIPTION_ID \
  --output-directory kube-rg
```

## Cluster Definition (aka API Model)

This section details how to tailor your cluster definitions in order to make them compatible with Azure Stack Hub. You can start off from this [template](../../examples/azure-stack/kubernetes-azurestack.json).

Unless otherwise specified down below, standard [cluster definition](../../docs/topics/clusterdefinitions.md) properties should also work with Azure Stack Hub. Please create an [issue](https://github.com/Azure/aks-engine-azurestack/issues/new) if you find that we missed a property that should be called out.

### location

| Name       | Required | Description                                                   |
| ---------- | -------- | ------------------------------------------------------------- |
| location   | yes      | The region name of the target Azure Stack Hub. |

### kubernetesConfig

`kubernetesConfig` describes Kubernetes specific configuration.

| Name                            | Required | Description                          |
| ------------------------------- | -------- | ------------------------------------ |
| addons                          | no       | A few addons are not supported on Azure Stack Hub. See the [complete list](#unsupported-addons) down below.|
| kubernetesImageBase             | no       | For AKS Engine versions lower than v0.48.0, this is a required field. It specifies the default image base URL to be used for all Kubernetes-related containers such as hyperkube, cloud-controller-manager, pause, addon-manager, etc. This property should be set to `"mcr.microsoft.com/k8s/azurestack/core/"`. |
| networkPlugin                   | yes      | Specifies the network plugin implementation for the cluster. Valid values are `"kubenet"` (default) for k8s software networking implementation and `"azure"`, which provides an Azure native networking experience. |
| networkPolicy                   | no      | Specifies the network policy enforcement tool for the cluster (currently Linux-only). Valid values are: `"azure"` (experimental) for Azure CNI-compliant network policy (note: Azure CNI-compliant network policy requires explicit `"networkPlugin": "azure"` configuration as well). |
| useInstanceMetadata             | no      | Use the Azure cloud provider instance metadata service for appropriate resource discovery operations. This property should be always set to `"false"`. |

### customCloudProfile

`customCloudProfile` contains information specific to the target Azure Stack Hub instance.

| Name                            | Required | Description|
| ------------------------------- | -------- | ---------- |
| environment                     | no       | The custom cloud type. This property should be always set to `"AzureStackCloud"`. |
| identitySystem                  | yes      | Specifies the identity provider used by the Azure Stack Hub instance. Valid values are `"azure_ad"` (default) and `"adfs"`. |
| portalUrl                       | yes      | The tenant portal URL. |
| dependenciesLocation            | no       | Specifies where to locate the dependencies required to during the provision/upgrade process. Valid values are `"public"` (default), `"china"`, `"german"` and `"usgovernment".` |

### masterProfile

`masterProfile` describes the settings for control plane configuration.

| Name                            | Required | Description|
| ------------------------------- | -------- | ---------- |
| vmsize                          | yes      | Specifies a valid [Azure Stack Hub VM size](https://docs.microsoft.com/azure-stack/user/azure-stack-vm-sizes). |
| distro                          | yes      | Specifies the control plane's Linux distribution. `"aks-ubuntu-18.04"` is supported. This is a custom image based on UbuntuServer that come with pre-installed software necessary for Kubernetes deployments. |

### agentPoolProfiles

`agentPoolProfiles` are used to create agents with different capabilities.

| Name                            | Required | Description|
| ------------------------------- | -------- | ---------- |
| vmsize                          | yes      | Describes a valid [Azure Stack Hub VM size](https://docs.microsoft.com/azure-stack/user/azure-stack-vm-sizes). |
| osType                          | no       | Specifies the agent pool's Operating System. Supported values are `"Windows"` and `"Linux"`. Defaults to `"Linux"`. |
| distro                          | yes      | Specifies the control plane's Linux distribution. `"aks-ubuntu-18.04"` is supported. This is a custom image based on UbuntuServer that come with pre-installed software necessary for Kubernetes deployments. |
| availabilityProfile             | yes      | Only `"AvailabilitySet"` is currently supported. |
| acceleratedNetworkingEnabled    | yes      | Use `Azure Accelerated Networking` feature for Linux agents. This property should be always set to `"false"`. |

`linuxProfile` provides the linux configuration for each linux node in the cluster
| Name                            | Required | Description|
| ------------------------------- | -------- | ---------- |
| enableUnattendedUpgrades | no    | Configure each Linux node VM (including control plane node VMs) to run `/usr/bin/unattended-upgrade` in the background according to a daily schedule. If enabled, the default `unattended-upgrades` package configuration will be used as provided by the Ubuntu distro version running on the VM. More information [here](https://help.ubuntu.com/community/AutomaticSecurityUpdates). By default, `enableUnattendedUpgrades` is set to `true`.|
| runUnattendedUpgradesOnBootstrap| no       | Invoke an unattended-upgrade when each Linux node VM comes online for the first time. In practice this is accomplished by performing an `apt-get update`, followed by a manual invocation of `/usr/bin/unattended-upgrade`, to fetch updated apt configuration, and install all package updates provided by the unattended-upgrade facility, respectively. Defaults to `"false"`. |

## Azure Stack Hub Instances Registered with Azure's China cloud

If your Azure Stack Hub instance is located in China, then the `dependenciesLocation` property of your cluster definition should be set to `"china"`. This switch ensures that the provisioning process fetches software dependencies from reachable hosts within China's mainland.

## Disconnected Azure Stack Hub Instances

By default, the AKS Engine provisioning process relies on an internet connection to download the software dependencies required to create or upgrade a cluster (Kubernetes images, etcd binaries, network plugins and so on).

If your Azure Stack Hub instance is air-gapped or if network connectivity in your geographical location is not reliable, then the default approach will not work, take a long time or timeout due to transient networking issues.

To overcome these issues, you should set the `distro` property of your cluster definition to `"aks-ubuntu-18.04"`. This will instruct AKS Engine to deploy VM nodes using a base OS image called `AKS Base Image`. This custom image, generally based on Ubuntu Server, already contains the required software dependencies in its file system. Hence, internet connectivity won't be required during the provisioning process.

The `AKS Base Image` marketplace item has to be available in your Azure Stack Hub's Marketplace before it could be used by AKS Engine. Your Azure Stack Hub administrator can follow this [guide](https://docs.microsoft.com/azure-stack/operator/azure-stack-download-azure-marketplace-item) for a general explanation about how to download marketplace items from Azure.

Each AKS Engine release is validated and tied to a specific version of the AKS Base Image. Therefore, you need to take note of the base image version required by the AKS Engine release that you plan to use, and then download exactly that base image version. New builds of the `AKS Base Image` are frequently released to ensure that your disconnected cluster can be upgraded to the latest supported version of each component.

Make sure `linuxProfile.runUnattendedUpgradesOnBootstrap` is set to `"false"` when you deploy, or upgrade, a cluster to air-gapped Azure Stack Hub clouds.

## AKS Engine Versions

| AKS Engine                 | AKS Base Image     | Kubernetes versions | Notes |
|----------------------------|--------------------|---------------------|-------|
| [v0.43.1](https://github.com/Azure/aks-engine/releases/tag/v0.43.1)   | [AKS Base Ubuntu 16.04-LTS Image Distro, October 2019 (2019.10.24)](https://github.com/Azure/aks-engine/blob/v0.43.0/releases/vhd-notes/aks-ubuntu-1604/aks-ubuntu-1604-201910_2019.10.24.txt) | 1.15.5, 1.15.4, 1.14.8, 1.14.7 |  |
| [v0.48.0](https://github.com/Azure/aks-engine/releases/tag/v0.48.0)   | [AKS Base Ubuntu 16.04-LTS Image Distro, March 2020 (2020.03.19)](https://github.com/Azure/aks-engine/blob/v0.48.0/vhd/release-notes/aks-engine-ubuntu-1604/aks-engine-ubuntu-1604-202003_2020.03.19.txt) | 1.15.10, 1.14.7 |  |
| [v0.51.0](https://github.com/Azure/aks-engine/releases/tag/v0.51.0)   | [AKS Base Ubuntu 16.04-LTS Image Distro, May 2020 (2020.05.13)](https://github.com/Azure/aks-engine/blob/v0.51.0/vhd/release-notes/aks-engine-ubuntu-1604/aks-engine-ubuntu-1604-202005_2020.05.13.txt), [AKS Base Windows Image (17763.1217.200513)](https://github.com/Azure/aks-engine/blob/v0.51.0/vhd/release-notes/aks-windows/2019-datacenter-core-smalldisk-17763.1217.200513.txt)  | 1.15.12, 1.16.8, 1.16.9 | API Model Samples ([Linux](https://github.com/Azure/aks-engine/blob/v0.51.0/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine/blob/v0.51.0/examples/azure-stack/kubernetes-windows.json)) |
| [v0.55.0](https://github.com/Azure/aks-engine/releases/tag/v0.55.0)   | [AKS Base Ubuntu 16.04-LTS Image Distro, August 2020 (2020.08.24)](https://github.com/Azure/aks-engine/blob/v0.55.0/vhd/release-notes/aks-engine-ubuntu-1604/aks-engine-ubuntu-1604-202007_2020.08.24.txt), [AKS Base Windows Image (17763.1397.200820)](https://github.com/Azure/aks-engine/blob/v0.55.0/vhd/release-notes/aks-windows/2019-datacenter-core-smalldisk-17763.1397.200820.txt)  | 1.15.12, 1.16.14, 1.17.11 | API Model Samples ([Linux](https://github.com/Azure/aks-engine/blob/v0.55.0/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine/blob/v0.55.0/examples/azure-stack/kubernetes-windows.json)) |
| [v0.55.4](https://github.com/Azure/aks-engine/releases/tag/v0.55.4)   | [AKS Base Ubuntu 16.04-LTS Image Distro, September 2020 (2020.09.14)](https://github.com/Azure/aks-engine/blob/v0.55.0/vhd/release-notes/aks-engine-ubuntu-1604/aks-engine-ubuntu-1604-202007_2020.08.24.txt), [AKS Base Windows Image (17763.1397.200820)](https://github.com/Azure/aks-engine/blob/v0.55.0/vhd/release-notes/aks-windows/2019-datacenter-core-smalldisk-17763.1397.200820.txt)  | 1.15.12, 1.16.14, 1.17.11 | API Model Samples ([Linux](https://github.com/Azure/aks-engine/blob/v0.55.0/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine/blob/v0.55.0/examples/azure-stack/kubernetes-windows.json)) |
| [v0.60.1](https://github.com/Azure/aks-engine/releases/tag/v0.60.1)   | [AKS Base Ubuntu 18.04-LTS Image Distro, 2021 Q1 (2021.01.28)](https://github.com/Azure/aks-engine/blob/v0.60.1/vhd/release-notes/aks-engine-ubuntu-1804/aks-engine-ubuntu-1804-202007_2021.01.28.txt), [AKS Base Ubuntu 16.04-LTS Image Distro, January 2021 (2021.01.28)](https://github.com/Azure/aks-engine/blob/v0.60.1/vhd/release-notes/aks-engine-ubuntu-1604/aks-engine-ubuntu-1604-202007_2021.01.28.txt), [AKS Base Windows Image (17763.1697.210129)](https://github.com/Azure/aks-engine/blob/v0.60.1/vhd/release-notes/aks-windows/2109-datacenter-core-smalldisk-17763.1697.210129.txt)  | 1.16.14, 1.16.15, 1.17.17, 1.18.15 | API Model Samples ([Linux](https://github.com/Azure/aks-engine/blob/v0.60.1/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine/blob/v0.60.1/examples/azure-stack/kubernetes-windows.json)) |
| [v0.63.0](https://github.com/Azure/aks-engine/releases/tag/v0.63.0)   | [AKS Base Ubuntu 18.04-LTS Image Distro, 2021 Q2 (2021.05.24)](https://github.com/Azure/aks-engine/blob/v0.63.0/vhd/release-notes/aks-engine-ubuntu-1804/aks-engine-ubuntu-1804-202007_2021.05.24.txt), [AKS Base Windows Image (17763.1935.210520)](https://github.com/Azure/aks-engine/blob/v0.63.0/vhd/release-notes/aks-windows/2109-datacenter-core-smalldisk-17763.1935.210520.txt)  | 1.18.18, 1.19.10, 1.20.6 | API Model Samples ([Linux](https://github.com/Azure/aks-engine/blob/v0.65.0/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine/blob/v0.65.0/examples/azure-stack/kubernetes-windows.json)) |
| [v0.67.0](https://github.com/Azure/aks-engine/releases/tag/v0.67.0) | [AKS Base Ubuntu 18.04-LTS Image Distro, 2021 Q3 (2021.09.27)](https://github.com/Azure/aks-engine/blob/v0.67.0/vhd/release-notes/aks-engine-ubuntu-1804/aks-engine-ubuntu-1804-202007_2021.09.27.txt), [AKS Base Windows Image (17763.2213.210927)](https://github.com/Azure/aks-engine/blob/v0.67.0/vhd/release-notes/aks-windows/2019-datacenter-core-smalldisk-17763.2213.210927.txt) | 1.19.15, 1.20.11 | API Model Samples ([Linux](https://github.com/Azure/aks-engine/blob/v0.70.0/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine/blob/v0.70.0/examples/azure-stack/kubernetes-windows.json)) |
| [v0.67.3](https://github.com/Azure/aks-engine/releases/tag/v0.67.3)   | [AKS Base Ubuntu 18.04-LTS Image Distro, 2021 Q3 (2021.09.27)](https://github.com/Azure/aks-engine/blob/v0.67.3/vhd/release-notes/aks-engine-ubuntu-1804/aks-engine-ubuntu-1804-202007_2021.09.27.txt), [AKS Base Windows Image (17763.2213.210927)](https://github.com/Azure/aks-engine/blob/v0.67.3/vhd/release-notes/aks-windows/2019-datacenter-core-smalldisk-17763.2213.210927.txt)  | 1.19.15, 1.20.11 | API Model Samples ([Linux](https://github.com/Azure/aks-engine/blob/v0.70.0/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine/blob/v0.70.0/examples/azure-stack/kubernetes-windows.json)) |
| [v0.70.0](https://github.com/Azure/aks-engine/releases/tag/v0.70.0)   | [AKS Base Ubuntu 18.04-LTS Image Distro, 2022 Q2 (2022.04.07)](https://github.com/Azure/aks-engine/blob/v0.70.0/vhd/release-notes/aks-engine-ubuntu-1804/aks-engine-ubuntu-1804-202112_2022.04.07.txt), [AKS Base Windows Image (17763.2565.220408)](https://github.com/Azure/aks-engine/blob/v0.70.0/vhd/release-notes/aks-windows/2019-datacenter-core-smalldisk-17763.2565.220408.txt)  | 1.21.10*, 1.22.7* | API Model Samples ([Linux](https://github.com/Azure/aks-engine/blob/v0.71.0/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine/blob/v0.71.0/examples/azure-stack/kubernetes-windows.json)) |
| [v0.71.0](https://github.com/Azure/aks-engine/releases/tag/v0.71.0)   | [AKS Base Ubuntu 18.04-LTS Image Distro, 2022 Q3 (2022.08.12)](https://github.com/Azure/aks-engine/blob/v0.71.0/vhd/release-notes/aks-engine-ubuntu-1804/aks-engine-ubuntu-1804-202112_2022.08.12.txt), [AKS Base Windows Image (17763.3232.220805)](https://github.com/Azure/aks-engine/blob/v0.71.0/vhd/release-notes/aks-windows/2019-datacenter-core-smalldisk-17763.3232.220805.txt)  | 1.22.7*, 1.23.6* | API Model Samples ([Linux](https://github.com/Azure/aks-engine/blob/v0.73.0/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine/blob/v0.73.0/examples/azure-stack/kubernetes-windows.json)) |
| [v0.73.0](https://github.com/Azure/aks-engine/releases/tag/v0.73.0)   | [AKS Base Ubuntu 18.04-LTS Image Distro, 2022 Q4 (2022.11.02)](https://github.com/Azure/aks-engine/blob/v0.73.0/vhd/release-notes/aks-engine-ubuntu-1804/aks-engine-ubuntu-1804-202112_2022.11.02.txt), [AKS Base Windows Image (17763.3532.221102)](https://github.com/Azure/aks-engine/blob/v0.73.0/vhd/release-notes/aks-windows/2019-datacenter-core-smalldisk-17763.3532.221102.txt)  | 1.22.15*, 1.23.13* | API Model Samples ([Linux](https://github.com/Azure/aks-engine-azurestack/blob/v0.75.3/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine-azurestack/blob/v0.75.3/examples/azure-stack/kubernetes-windows.json)) |
| [v0.75.3](https://github.com/Azure/aks-engine-azurestack/releases/tag/v0.75.3)   | [AKS Base Ubuntu 20.04-LTS Image Distro (2023.032.2)](https://github.com/Azure/aks-engine-azurestack/blob/v0.75.3/vhd/release-notes/aks-engine-ubuntu-2004/aks-engine-azurestack-ubuntu-2004_2023.032.2.txt), [AKS Base Windows Server 2019 Image Docker (17763.3887.20230332)](https://github.com/Azure/aks-engine-azurestack/blob/v0.75.3/vhd/release-notes/aks-windows/2019-datacenter-core-azurestack-smalldisk-17763.3887.20230332.txt), [AKS Base Windows Server 2019 Image Containerd (17763.3887.20230332)](https://github.com/Azure/aks-engine-azurestack/blob/v0.75.3/vhd/release-notes/aks-windows-2019-containerd/2019-datacenter-core-azurestack-ctrd-17763.3887.20230332.txt) | 1.23.15*, 1.24.9** | API Model Samples ([Linux](https://github.com/Azure/aks-engine-azurestack/blob/v0.76.0/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine-azurestack/blob/v0.76.0/examples/azure-stack/kubernetes-windows.json)) |
| [v0.76.0](https://github.com/Azure/aks-engine-azurestack/releases/tag/v0.76.0)   | [AKS Base Ubuntu 20.04-LTS Image Distro (2023.116.3)](https://github.com/Azure/aks-engine-azurestack/blob/v0.76.0/vhd/release-notes/aks-engine-ubuntu-2004/aks-engine-azurestack-ubuntu-2004_2023.116.3.txt), [AKS Base Windows Server 2019 Image Containerd (17763.4252.20231163)](https://github.com/Azure/aks-engine-azurestack/blob/v0.76.0/vhd/release-notes/aks-windows-2019-containerd/2019-datacenter-core-azurestack-ctrd-17763.4252.20231163.txt) | 1.24.11**, 1.25.7** | API Model Samples ([Linux](https://github.com/Azure/aks-engine-azurestack/blob/v0.77.0/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine-azurestack/blob/v0.77.0/examples/azure-stack/kubernetes-windows.json)) |
| [v0.77.0](https://github.com/Azure/aks-engine-azurestack/releases/tag/v0.77.0) | [AKS Base Ubuntu 20.04-LTS Image Distro (2023.206.1)](https://github.com/Azure/aks-engine-azurestack/blob/v0.77.0/vhd/release-notes/aks-engine-ubuntu-2004/aks-engine-azurestack-ubuntu-2004_2023.206.1.txt), [AKS Base Windows Server 2019 Image Containerd (17763.4645.20232061)](https://github.com/Azure/aks-engine-azurestack/blob/v0.77.0/vhd/release-notes/aks-windows-2019-containerd/2019-datacenter-core-azurestack-ctrd-17763.4645.20232061.txt) | 1.25.7**, 1.26.6** | API Model Samples ([Linux](https://github.com/Azure/aks-engine-azurestack/blob/e726bc30ece616a5282c1d027048e60d0696f7ad/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine-azurestack/blob/e726bc30ece616a5282c1d027048e60d0696f7ad/examples/azure-stack/kubernetes-windows.json)) |
| [v0.78.0](https://github.com/Azure/aks-engine-azurestack/releases/tag/v0.78.0)   | [AKS Base Ubuntu 20.04-LTS Image Distro (2023.242.3)](https://github.com/Azure/aks-engine-azurestack/blob/v0.78.0/vhd/release-notes/aks-engine-ubuntu-2004/aks-engine-azurestack-ubuntu-2004_2023.242.3.txt), [AKS Base Windows Server 2019 Image Containerd (17763.4737.20232423)](https://github.com/Azure/aks-engine-azurestack/blob/v0.78.0/vhd/release-notes/aks-windows-2019-containerd/2019-datacenter-core-azurestack-ctrd-17763.4737.20232423.txt) | 1.25.13**, 1.26.8** | API Model Samples ([Linux](https://github.com/Azure/aks-engine-azurestack/blob/3d5c7916e9995045a0868945d1f732aaa4c5b5bc/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine-azurestack/blob/3d5c7916e9995045a0868945d1f732aaa4c5b5bc/examples/azure-stack/kubernetes-windows.json)) |
| [v0.79.0](https://github.com/Azure/aks-engine-azurestack/releases/tag/v0.79.0)   | [AKS Base Ubuntu 20.04-LTS Image Distro (2023.296.1)](https://github.com/Azure/aks-engine-azurestack/blob/v0.79.0/vhd/release-notes/aks-engine-ubuntu-2004/aks-engine-azurestack-ubuntu-2004_2023.296.1.txt), [AKS Base Windows Server 2019 Image Containerd (17763.4974.20232961)](https://github.com/Azure/aks-engine-azurestack/blob/v0.79.0/vhd/release-notes/aks-windows-2019-containerd/2019-datacenter-core-azurestack-ctrd-17763.4974.20232961.txt) | 1.26.9**, 1.27.6** | API Model Samples ([Linux](https://github.com/Azure/aks-engine-azurestack/blob/af3e6b25e6304f7c68eb90965b8e9a36887cb457/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine-azurestack/blob/af3e6b25e6304f7c68eb90965b8e9a36887cb457/examples/azure-stack/kubernetes-windows.json)) |
| [v0.80.3](https://github.com/Azure/aks-engine-azurestack/releases/tag/v0.80.3)   | [AKS Base Ubuntu 20.04-LTS Image Distro (2024.032.1)](https://github.com/Azure/aks-engine-azurestack/blob/v0.80.2/vhd/release-notes/aks-engine-ubuntu-2004/aks-engine-azurestack-ubuntu-2004_2024.032.1.txt), [AKS Base Windows Server 2019 Image Containerd (17763.5329.20240321)](https://github.com/Azure/aks-engine-azurestack/blob/v0.80.2/vhd/release-notes/aks-windows-2019-containerd/2019-datacenter-core-azurestack-ctrd-17763.5329.20240321.txt) | 1.27.10**, 1.28.6** | API Model Samples ([Linux](https://github.com/Azure/aks-engine-azurestack/blob/96b45e116aef6ae2e5031561b0e10621ae86a068/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine-azurestack/blob/96b45e116aef6ae2e5031561b0e10621ae86a068/examples/azure-stack/kubernetes-windows.json)) |
| [v0.81.1](https://github.com/Azure/aks-engine-azurestack/releases/tag/v0.81.1)   | [AKS Base Ubuntu 20.04-LTS Image Distro (2024.311.1)](https://github.com/Azure/aks-engine-azurestack/blob/master/vhd/release-notes/aks-engine-ubuntu-2004/aks-engine-azurestack-ubuntu-2004_2024.311.1.txt), [AKS Base Windows Server 2019 Image Containerd (17763.6414.20243111)](https://github.com/Azure/aks-engine-azurestack/blob/master/vhd/release-notes/aks-windows-2019-containerd/2019-datacenter-core-azurestack-ctrd-17763.6414.20243111.txt) | 1.28.15**, 1.29.10** | API Model Samples ([Linux](https://github.com/Azure/aks-engine-azurestack/blob/master/examples/azure-stack/kubernetes-azurestack.json), [Windows](https://github.com/Azure/aks-engine-azurestack/blob/master/examples/azure-stack/kubernetes-windows.json)) |

>\* Starting from Kubernetes v1.21, **only** [Cloud Provider for Azure](#cloud-provider-for-azure) is supported on Azure Stack Hub.

>\** Starting from Kubernetes v1.24, **ONLY** the `containerd` container runtime is supported. Please refer to the section [*Upgrading Kubernetes clusters created with docker container runtime*](#upgrading-kubernetes-clusters-created-with-docker-container-runtime) for more details.

## Cloud Provider for Azure

The [Cloud Provider for Azure](https://github.com/kubernetes-sigs/cloud-provider-azure) project (aka `cloud-controller-manager`, out-of-tree cloud provider or external cloud provider) implements the [Kubernetes cloud provider interface](https://github.com/kubernetes/cloud-provider) for Azure clouds. The out-of-tree implementation is the replacement for the deprecated [in-tree implementation](https://github.com/kubernetes/kubernetes/tree/master/staging/src/k8s.io/legacy-cloud-providers/azure).

On Azure Stack Hub, starting from Kubernetes v1.21, AKS Engine-based clusters will exclusively use `cloud-controller-manager`. Hence, to deploy a Kubernetes v1.21+ cluster, it is required to set `orchestratorProfile.kubernetesConfig.useCloudControllerManager` to `true` in the API Model ([example](/examples/azure-stack/kubernetes-azurestack.json)). AKS Engine's upgrade process will automatically update the `useCloudControllerManager` flag.

> **Upgrade considerations:** the process of upgrading a Kubernetes cluster from v1.20 (or lower version) to v1.21 (or greater version) will cause downtime to workloads relying on the `kubernetes.io/azure-disk` in-tree volume provisioner. Before upgrading to Kubernetes v1.21+, it is **highly recommended** to perform a full backup of the application data and validate in a **pre-production environment** that the cluster storage resources (PV and PVC) can be migrated to the a new volume provisioner. Learn how to migrate to the Azure Disk CSI driver [here](#migrate-persistent-storage-to-the-azure-disk-csi-driver).

## Volume Provisioners: Container Storage Interface Drivers (Azure Disk CSI Driver)

The [in-tree volume provisioner](https://kubernetes.io/blog/2019/12/09/kubernetes-1-17-feature-csi-migration-beta/) is only compatible with the in-tree cloud provider. Therefore, a v1.21+ cluster has to include a Container Storage Interface (CSI) Driver if user workloads rely on persistent storage.

AKS Engine will **not** enable any CSI driver by default on Azure Stack Hub. For workloads that require a CSI driver, it is possible to either:

* For AKS Engine versions 0.75.3 and above: explicitly enable the `azuredisk-csi-driver` [addon](../topics/clusterdefinitions.md#addons) (Linux and/or Windows cluster)
  * If you are using Windows HostProcess containers in your workloads, do not use the addon. Instead use `Helm` to [install the `azuredisk-csi-driver` chart](#1-install-azure-disk-csi-driver-manually) and set the following value: `--set windows.useHostProcessContainers=true`. (Windows cluster only)
* For AKS Engine versions 0.70.0 and above: explicitly enable the `azuredisk-csi-driver` [addon](../topics/clusterdefinitions.md#addons) (Linux cluster only)
* For AKS Engine versions 0.70.0 and above: use `Helm` to [install the `azuredisk-csi-driver` chart](#1-install-azure-disk-csi-driver-manually) (Linux and/or Windows clusters).

### Azure Disk CSI Driver: Version Mapping
The following table shows the supported Azure Disk CSI Driver version for each Kubernetes version.
- If using addon, the Azure Disk CSI Driver version is automatically selected based on the cluster's Kubernetes version. 
- If using Helm, the Azure Disk CSI Driver version has to be specified manually based on the cluster's Kubernetes version.

| Kubernetes Version | Azure Disk CSI Driver Version |
| ------------------ | ----------------------------- |
| <= 1.25            | 1.10.0                        |
| 1.26               | 1.26.5                        |
| 1.27               | 1.28.3                        |
| >= 1.28            | 1.29.1                        |

* Exception: For AKS Engine v0.77.0, kubernetes version 1.25 uses Azure Disk CSI Driver version v1.26.5
* Note: Only the versions in the mapping have been tested. While other version combinations may work, they are not officially supported.

### Migrate Persistent Storage to the Azure Disk CSI driver

The process of upgrading an AKS Engine-based cluster from v1.20 (or lower version) to v1.21 (or greater version) will cause downtime to workloads relying on the `kubernetes.io/azure-disk` in-tree volume provisioner as this provisioner is not part of the [Cloud Provider for Azure](#cloud-provider-for-azure).

If the data persisted in the underlying Azure disks should be preserved, then the following extra steps are required once the cluster upgrade process is completed:

1. [Install the Azure Disk CSI driver](#1-install-azure-disk-csi-driver-manually).
1. [Remove the deprecated in-tree storage classes](#2-replace-storage-classes).
1. [Recreate the persistent volumes and claims](#3-recreate-persistent-volumes).

#### 1. Install Azure Disk CSI driver manually

The following script uses `Helm` to install the Azure Disk CSI Driver:

```bash
DRIVER_VERSION=v1.26.5 # if using k8s v1.26
DRIVER_VERSION=v1.28.3 # if using k8s v1.27
DRIVER_VERSION=v1.29.1 # if using k8s >= v1.28
DRIVER_VERSION=v1.32.0 # if using k8s >= v1.30
helm repo add azuredisk-csi-driver https://raw.githubusercontent.com/kubernetes-sigs/azuredisk-csi-driver/master/charts
helm install azuredisk-csi-driver azuredisk-csi-driver/azuredisk-csi-driver \
  --namespace kube-system \
  --set cloud=AzureStackCloud \
  --set controller.runOnMaster=true \
  --set node.supportZone=false \
  --set windows.getNodeInfoFromLabels=true \
  --set linux.getNodeInfoFromLabels=true \
  # --set windows.useHostProcessContainers=true \ # if using Windows HostProcess containers
  --version ${DRIVER_VERSION}
```

#### 2. Replace Storage Classes

The `kube-addon-manager` will automatically create the Azure Disk CSI driver storage classes (`disk.csi.azure.com`) once the in-tree storage classes (`kubernetes.io/azure-disk`) are manually deleted:

```bash
IN_TREE_SC="default managed-premium managed-standard"

# Delete deprecated "kubernetes.io/azure-disk" storage classes
kubectl delete storageclasses ${IN_TREE_SC}

# Wait for addon manager to create the "disk.csi.azure.com" storage class resources
kubectl get --watch storageclasses
```

#### 3. Recreate Persistent Volumes

Once the Azure Disk CSI Driver is installed and the storage classes replaced, the next step is to recreate the persistent volumes (PV) and persistent volumes claims (PVC) using the Azure Disk CSI driver (or alternative CSI solution).

This is a multi-step process that can be different depending on how these resources were initially deployed. The high level steps are:

* Delete the deployment or statefulset that references the PV + PVC pairs to migrate (backup resource definition if necessary).
* Ensure the PVs' `persistentVolumeReclaimPolicy` property is set to value `Retain` ([example](https://learn.microsoft.com/en-us/azure/aks/csi-storage-drivers#migrate-in-tree-persistent-volumes)).
* Delete the PV + PVC pairs to migrate (backup resource definitions if necessary).
* To migrate, update the PVs' resource definition by removing the `azureDisk` object and adding a `csi` object with reference to the original AzureDisk ([example](https://github.com/kubernetes-sigs/azuredisk-csi-driver/blob/master/docs/driver-parameters.md#static-provisioning-bring-your-own-azure-disk)).
* Recreate, in the following order, the PV resource/s, PVC resource/s (if necessary), and finally the deployment or statefulset.

The following migration [script](../../examples/azure-stack/migratepv.sh) is provided as a template.

> After running the migration script, if the pod is stuck with error "Unable to attach or mount volumes", make sure [Azure Disk CSI Driver was installed](#1-install-azure-disk-csi-driver-manually) and [storage classes were recreated](#2-replace-storage-classes).

### Azure Disk CSI Driver: Details

As a [replacement of the current in-tree volume provisioner](https://kubernetes.io/blog/2019/12/09/kubernetes-1-17-feature-csi-migration-beta/), Azure Disk Container Storage Interface (CSI) Driver is available on Azure Stack Hub. Please find details in the following table.

|                       | Azure Disk CSI Driver
|-----------------------|------------------------------------------------------------------------------------------------------------------------------
| Stage on Azure Stack Hub | GA
| Project Repository    | [azuredisk-csi-driver](https://github.com/kubernetes-sigs/azuredisk-csi-driver)
| CSI Driver Version    | v1.0.0+
| Access Mode           | ReadWriteOnce
| Windows Agent Node    | Support
| Dynamic Provisioning  | Support
| Considerations        | [Azure Disk CSI Driver Limitations](https://github.com/kubernetes-sigs/azuredisk-csi-driver/blob/master/docs/limitations.md)
| Slack Support Channel | [#provider-azure](https://kubernetes.slack.com/archives/C5HJXTT9Q)

> To deploy a CSI driver to an air-gapped cluster, make sure that your `helm` chart is referencing container images that are reachable from the cluster nodes.

### Azure Disk CSI Driver: Requirements

* Azure Stack build 2011 and later.
* AKS Engine version v0.60.1 and later.
* Kubernetes version 1.18 and later.
* Since the Controller server of CSI Drivers requires 2 replicas, a single node master pool is not recommended.
* [Helm 3](https://helm.sh/docs/intro/install/)

### Azure Disk CSI Driver: Examples

In this section, please follow the example commands to deploy a StatefulSet application consuming CSI Driver.

```bash
# Install CSI Driver
DRIVER_VERSION=v1.26.5 # if using k8s v1.26
DRIVER_VERSION=v1.28.3 # if using k8s v1.27
DRIVER_VERSION=v1.29.1 # if using k8s >= v1.28
DRIVER_VERSION=v1.32.0 # if using k8s >= v1.30
helm repo add azuredisk-csi-driver https://raw.githubusercontent.com/kubernetes-sigs/azuredisk-csi-driver/master/charts
helm install azuredisk-csi-driver azuredisk-csi-driver/azuredisk-csi-driver \
  --namespace kube-system \
  --set cloud=AzureStackCloud \
  --set controller.runOnMaster=true \
  --set node.supportZone=false \
  --set windows.getNodeInfoFromLabels=true \
  --set linux.getNodeInfoFromLabels=true \
  # --set windows.useHostProcessContainers=true \ # if using Windows HostProcess containers
  --version ${DRIVER_VERSION}

# Deploy Storage Class
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/azuredisk-csi-driver/master/deploy/example/storageclass-azuredisk-csi-azurestack.yaml

# Deploy example StatefulSet application
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/azuredisk-csi-driver/master/deploy/example/statefulset.yaml

# Validate volumes and applications
# You should see a sequence of timestamps are persisted in the volume.
kubectl exec statefulset-azuredisk-0 -- tail /mnt/azuredisk/outfile

# Delete example StatefulSet application
kubectl delete -f https://raw.githubusercontent.com/kubernetes-sigs/azuredisk-csi-driver/master/deploy/example/statefulset.yaml

# Delete Storage Class
# Before delete the Storage Class, please make sure Pods that consume the Storage Class have been terminated.
kubectl delete -f https://raw.githubusercontent.com/kubernetes-sigs/azuredisk-csi-driver/master/deploy/example/storageclass-azuredisk-csi-azurestack.yaml

# Uninstall CSI Driver
helm uninstall azuredisk-csi-driver --namespace kube-system
helm repo remove azuredisk-csi-driver
```

## Known Issues and Limitations

This section lists all known issues you may find when you use the GA version.

### Unsupported Addons

AKS Engine includes a number of optional [addons](../topics/clusterdefinitions.md#addons) that can be deployed as part of the cluster provisioning process.

The list below includes the addons currently unsupported on Azure Stack Hub:

* AAD Pod Identity
* Cluster Autoscaler
* KeyVault Flex Volume
* SMB Flex Volume

### Chrony daemon fails to restart

The `chrony` daemon on a Linux node may fail to restart with error message: `"Could not open /dev/ptp_hyperv: No such file or directory"`.
The workaround for this issue is to manually reboot the affected Linux nodes.
This operation will regenerate the `/dev/ptp_hyperv` symlink and allow the chrony daemon to restart successfully.

### OSProfile exceeds maximum characters length error

Addons enabled in the API Model are Base64 encoded and included in the VMs ARM template. There is a length limit of 87380 characters for the custom data, thus if too many addons are enabled in the API Model, the `aks-engine-azurestack` operations could fail with the below error:

```console
Custom data in OSProfile must be in Base64 encoding and with a maximum length of 87380 characters
```

In such cases, try reduce the number of enabled addons or remove all of them in the API Model.

### get-versions command

By default, `aks-engine-azurestack get-versions` shows which Kubernetes versions are supported by each AKS Engine release on Azure's public cloud. Include flag `--azure-env` to get the list of supported Kubernetes versions on a custom cloud such as an Azure Stack Hub cloud (`aks-engine-azurestack get-versions --azure-env AzureStackCloud`). Upgrade paths for Azure Stack Hub can also be found [here](https://docs.microsoft.com/azure-stack/user/kubernetes-aks-engine-release-notes).

### Upgrade from private-preview Kubernetes cluster with Windows nodes

There is no official support for private-preview Kubernetes cluster with Windows nodes created with AKS Engine v0.43.1 to upgrade with AKS Engine v0.55.0. Users are encouraged to deploy new Kubernetes cluster with Windows nodes with the latest AKS Engine version.

### Upgrading Kubernetes clusters created with the Ubuntu 18.04 distro

Starting with AKS Engine v0.75.3, the Ubuntu 18.04 distro is not longer a supported option as the OS reached its end-of-life. For AKS Engine v0.75.3 or later versions, aks-engine-azurestack upgrade will automatically overwrite the unsupported `aks-ubuntu-18.04` distro value with `aks-ubuntu-20.04`.

### Upgrading Kubernetes clusters created with the Ubuntu 16.04 distro

Starting with AKS Engine v0.63.0, the Ubuntu 16.04 distro is not longer a supported option as the OS reached its end-of-life. For AKS Engine v0.67.0 or later versions, aks-engine upgrade will automatically overwrite the unsupported `aks-ubuntu-16.04` distro value with with `aks-ubuntu-18.04`. For AKS Engine v0.75.3 or later versions, aks-engine-azurestack upgrade will automatically overwrite the unsupported `aks-ubuntu-16.04` distro value with `aks-ubuntu-20.04`.

### Upgrading Kubernetes clusters created with docker container runtime

In Kubernetes v1.24, the dockershim component was removed from kubelet. As a result, the docker container runtime is no longer a supported option. See [Kubernetes v1.24 release notes](https://kubernetes.io/blog/2022/05/03/kubernetes-1-24-release-announcement/#dockershim-removed-from-kubelet) for more information. For AKS Engine v0.75.3 or later versions, aks-engine-azurestack upgrade will automatically overwrite the unsupported `docker` containerRuntime value with `containerd`.

For AKS Engine release v0.75.3, clusters with windows nodes on Kubernetes v1.23 can use the Windows base image with Docker runtime. Clusters with windows nodes on Kubernetes v1.24 can use the Windows base image with Containerd runtime.

## Frequently Asked Questions

### Sample extensions are not working

Extensions in AKS Engine provide an easy way to include your own customization at provisioning time.

Because Azure and Azure Stack Hub currently rely on a different version of the Compute Resource Provider API, you may find that some of sample [extensions](https://github.com/Azure/aks-engine/tree/master/extensions) fail to deploy correctly.

This can be resolved by making a small modification to the extension `template.json` file. Replacing all usages of template parameter `apiVersionDeployments` by the hard-code value `2017-12-01` (or whatever API version Azure Stack Hub runs at the time you try to deploy) should be all you need.

Once you are done updating the extension template, host the extension directory in your own Github repository or storage account. Finally, at deployment time, make sure that your cluster definition points to the new [rootURL](https://github.com/Azure/aks-engine/blob/master/docs/topics/extensions.md#rooturl).

### The cluster nodes do not contain the latest Ubuntu OS security patches

If an `aks-ubuntu-18.04` image is created by the AKS Engine team prior to the release of an OS security patch, then the image won't include those security patches until an unattended upgrade is triggered.

When `linuxProfile.enableUnattendedUpgrades` is set to `true`, unattended upgrades will be checked and/or installed once a day. To ensure that the nodes are rebooted in a non-disruptive way, you can deploy [kured](https://github.com/weaveworks/kured) or similar solutions.

To deploy a cluster that includes the latests OS security patches right from the beginning, set `linuxProfile.runUnattendedUpgradesOnBootstrap` to `"true"` (see [example](../../examples/azure-stack/kubernetes-azurestack.json)).

To apply the latest OS security patches to an existing cluster, you can either do it manually or use the `aks-engine-azurestack upgrade` command. A manual upgrade can be done by executing `apt-get update && apt-get upgrade` and rebooting the node if necessary. If you use the `aks-engine-azurestack upgrade` command, set `linuxProfile.runUnattendedUpgradesOnBootstrap` to `"true"` in the generated `apimodel.json` and execute `aks-engine-azurestack upgrade` (a forced upgrade to the current Kubernetes version also works).

### Troubleshooting

This [how-to guide](/docs/howto/troubleshooting.md) has a good high-level explanation of how AKS Engine interacts with the Azure Resource Manager (ARM) and lists a few potential issues that can cause AKS Engine commands to fail.

Please refer to the [get-logs](../topics/get-logs.md) command documentation to simplify the logs collection task.

## Next Steps

* [What is the AKS engine on Azure Stack Hub?](https://docs.microsoft.com/azure-stack/user/azure-stack-kubernetes-aks-engine-overview)
