
<a name="v0.71.0"></a>
# [v0.71.0] - 2022-08-31
### Bug Fixes 🐞
- remove invalid feature gates ([#4903](https://github.com/Azure/aks-engine-azurestack/issues/4903))
- Set enableLoopbackDSR in Windows Azure CNI configuration for WinDSR ([#4793](https://github.com/Azure/aks-engine-azurestack/issues/4793))
- opt out of 'Anitimalware Signatures SDP' extension for Windows VHDs ([#4888](https://github.com/Azure/aks-engine-azurestack/issues/4888))

### Continuous Integration 💜
- remove 1.20 no-egress requirement ([#4944](https://github.com/Azure/aks-engine-azurestack/issues/4944))
- really use eastus for release automation ([#4923](https://github.com/Azure/aks-engine-azurestack/issues/4923))
- use eastus for release automation ([#4921](https://github.com/Azure/aks-engine-azurestack/issues/4921))
- use westus3 for release automation ([#4920](https://github.com/Azure/aks-engine-azurestack/issues/4920))

### Documentation 📘
- print deprecation message to stderr ([#4905](https://github.com/Azure/aks-engine-azurestack/issues/4905))
- update Azure Stack Hub Linux VHD version for v0.70.0 ([#4889](https://github.com/Azure/aks-engine-azurestack/issues/4889))
- Azure Stack Hub doc update for v0.70.0 ([#4856](https://github.com/Azure/aks-engine-azurestack/issues/4856))

### Features 🌈
- add support for Kubernetes v1.24.3 ([#4915](https://github.com/Azure/aks-engine-azurestack/issues/4915))
- add support for Kubernetes v1.23.9 ([#4914](https://github.com/Azure/aks-engine-azurestack/issues/4914))
- add support for Kubernetes v1.22.12 ([#4913](https://github.com/Azure/aks-engine-azurestack/issues/4913))
- add support for Kubernetes v1.21.14 ([#4911](https://github.com/Azure/aks-engine-azurestack/issues/4911))
- add support for Kubernetes v1.23.8 ([#4909](https://github.com/Azure/aks-engine-azurestack/issues/4909))
- add support for Kubernetes v1.22.11 ([#4910](https://github.com/Azure/aks-engine-azurestack/issues/4910))
- add support for Kubernetes v1.24.2 ([#4908](https://github.com/Azure/aks-engine-azurestack/issues/4908))
- add support for Kubernetes v1.21.13 ([#4900](https://github.com/Azure/aks-engine-azurestack/issues/4900))
- add support for Kubernetes v1.22.10 ([#4899](https://github.com/Azure/aks-engine-azurestack/issues/4899))
- add support for Kubernetes v1.23.7 ([#4901](https://github.com/Azure/aks-engine-azurestack/issues/4901))
- add support for Kubernetes v1.24.1 ([#4902](https://github.com/Azure/aks-engine-azurestack/issues/4902))
- Deprecate Kubernetes versions 1.19 and 1.20 ([#4894](https://github.com/Azure/aks-engine-azurestack/issues/4894))
- add chinaeast3 and chinanorth2 to AzureChinaCloud ([#4895](https://github.com/Azure/aks-engine-azurestack/issues/4895))
- add support for Kubernetes v1.22.9 ([#4883](https://github.com/Azure/aks-engine-azurestack/issues/4883))
- add support for Kubernetes v1.23.6 ([#4884](https://github.com/Azure/aks-engine-azurestack/issues/4884))
- add support for Kubernetes v1.24.0 ([#4890](https://github.com/Azure/aks-engine-azurestack/issues/4890))
- add support for Kubernetes v1.21.12 ([#4882](https://github.com/Azure/aks-engine-azurestack/issues/4882))
- add support for Kubernetes v1.23.5 ([#4861](https://github.com/Azure/aks-engine-azurestack/issues/4861))
- add support for Kubernetes v1.22.8 ([#4860](https://github.com/Azure/aks-engine-azurestack/issues/4860))
- add support for Kubernetes v1.21.11 ([#4859](https://github.com/Azure/aks-engine-azurestack/issues/4859))

### Maintenance 🔧
- rename project to aks-engine-azurestack
- updating Windows VHDs for next aks-engine release ([#4936](https://github.com/Azure/aks-engine-azurestack/issues/4936))
- Update Azure constants ([#4935](https://github.com/Azure/aks-engine-azurestack/issues/4935))
- pin github actions go versions to go.mod ([#4934](https://github.com/Azure/aks-engine-azurestack/issues/4934))
- Update Azure constants ([#4931](https://github.com/Azure/aks-engine-azurestack/issues/4931))
- upgrade container runtimes ([#4932](https://github.com/Azure/aks-engine-azurestack/issues/4932))
- Update Azure constants ([#4925](https://github.com/Azure/aks-engine-azurestack/issues/4925))
- add qatarcentral region ([#4919](https://github.com/Azure/aks-engine-azurestack/issues/4919))
- install 2022 7C updates to help vaildate windows networking fixes ([#4916](https://github.com/Azure/aks-engine-azurestack/issues/4916))
- add support for Kubernetes v1.23.6 on Azure Stack ([#4912](https://github.com/Azure/aks-engine-azurestack/issues/4912))
- get-logs copies kube-system logs from file system ([#4904](https://github.com/Azure/aks-engine-azurestack/issues/4904))
- Update Azure constants ([#4898](https://github.com/Azure/aks-engine-azurestack/issues/4898))
- update containerd to 1.5.11 ([#4891](https://github.com/Azure/aks-engine-azurestack/issues/4891))
- Update Azure constants ([#4896](https://github.com/Azure/aks-engine-azurestack/issues/4896))
- update out-of-tree cloud-provider-azure images ([#4892](https://github.com/Azure/aks-engine-azurestack/issues/4892))

### Testing 💚
- remove make test-kubernetes from workflows

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.71.0...HEAD
[v0.71.0]: https://github.com/Azure/aks-engine-azurestack/compare/v0.70.0...v0.71.0
