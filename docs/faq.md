# AKS Engine FAQ

This page provides help with the most common questions about AKS Engine.

## Is AKS Engine Currently Maintained?

This project is stable, meaning the pace of new features is intentionally low. New releases are expected to include support for the latest security updates and Kubernetes versions.
AKS Engine is maintained by teams who depend on it including [AKS Engine on Azure Stack Hub](https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-kubernetes-aks-engine-overview) and is intended for exclusive use on Azure Stack Hub.

## What's the Difference Between AKS and AKS Engine?

Azure Kubernetes Service ([AKS][]) is a Microsoft Azure service that supports fully managed Kubernetes clusters. [AKS Engine][], now only available for use on Azure Stack Hub, is an Azure open source project that allows you to create your own self-managed Kubernetes clusters with lots of user-configurable options.

## What's the Difference Between `Azure/aks-engine` and `Azure/aks-engine-azurestack`?

The deprecated [Azure/aks-engine] project supports ARM-based clouds such as Azure and Azure Stack Hub.

While Azure users are encouraged to migrate to [Cluster API Provider Azure](https://github.com/kubernetes-sigs/cluster-api-provider-azure) to provision sself-managed Kubernetes on Azure, the [Azure/aks-engine-azurestack] project is the path forward for Azure Stack Hub users.

Other than the fact that the deliverable included on each `Azure/aks-engine-azurestack` release will be named `aks-engine-azurestack` (instead of `aks-engine`), the user experience for Azure Stack Hub users remains the same.

[AKS]: https://azure.microsoft.com/en-us/services/kubernetes-service/
[AKS Engine]: https://github.com/Azure/aks-engine-azurestack
[Azure command-line tool]: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest
[Azure/aks-engine]: https://github.com/Azure/aks-engine
[Azure/aks-engine-azurestack]: https://github.com/Azure/aks-engine-azurestack
