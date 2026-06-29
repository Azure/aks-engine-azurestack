# Attention!

Notable changes in this release:

- Chocolatey is a package manager for Windows. Starting from this release, aks-engine-azurestack will no longer be released on chocolatey. 
  - Binary downloads for the latest version of AKS Engine are available on Github. Download the package for your operating system, and extract the aks-engine-azurestack file (optionally add it to your %PATH% environment variable for more convenient CLI usage).

<a name="v0.81.0"></a>
# [v0.81.0] - 2024-11-08
### Bug Fixes 🐞
- specify securityContext in cloud node manager ([#274](https://github.com/Azure/aks-engine-azurestack/issues/274))
- remove guest agent before sysprep to run windows cse on hub ([#273](https://github.com/Azure/aks-engine-azurestack/issues/273))
- specify instanceView to avoid nil pointer during rotate-certs ([#272](https://github.com/Azure/aks-engine-azurestack/issues/272))

### Build 🏭
- **deps:** bump github.com/Azure/azure-sdk-for-go/sdk/azidentity from 1.5.2 to 1.6.0 ([#258](https://github.com/Azure/aks-engine-azurestack/issues/258))

### Continuous Integration 💜
- fix Hub E2E tests REQUESTS_CA_BUNDLE ([#271](https://github.com/Azure/aks-engine-azurestack/issues/271))
- Replace go-dev base image ([#264](https://github.com/Azure/aks-engine-azurestack/issues/264))
- Output AIB release notes and cgmanifest to customization.log  ([#263](https://github.com/Azure/aks-engine-azurestack/issues/263))
- Fix disconnected tests ([#262](https://github.com/Azure/aks-engine-azurestack/issues/262))
- bring back deployment error unit tests ([#251](https://github.com/Azure/aks-engine-azurestack/issues/251))
- Add SkipAzLogin and CheckIngressIPOnly parameters ([#250](https://github.com/Azure/aks-engine-azurestack/issues/250))
- fix release github action ([#278](https://github.com/Azure/aks-engine-azurestack/issues/278))

### Documentation 📘
- Azure Stack Hub doc update for v0.80.2 ([#223](https://github.com/Azure/aks-engine-azurestack/issues/223))

### Maintenance 🔧
- update Linux and Windows VHDs for Oct 2024 ([#275](https://github.com/Azure/aks-engine-azurestack/issues/275))
- update Linux and Windows VHDs for Oct 2024 ([#270](https://github.com/Azure/aks-engine-azurestack/issues/270))
- Update Windows VHD AIB job to use Oct 2024 patches ([#269](https://github.com/Azure/aks-engine-azurestack/issues/269))
- support Kubernetes v1.28.15 & v1.29.10 ([#268](https://github.com/Azure/aks-engine-azurestack/issues/268))
- update containerd to v1.6.36 and runc to v1.1.14 ([#267](https://github.com/Azure/aks-engine-azurestack/issues/267))
- remove chocolatey package manager ([#266](https://github.com/Azure/aks-engine-azurestack/issues/266))
- update to golang v1.23 ([#265](https://github.com/Azure/aks-engine-azurestack/issues/265))
- support Kubernetes v1.29.8 ([#261](https://github.com/Azure/aks-engine-azurestack/issues/261))
- support Kubernetes v1.28.13 ([#260](https://github.com/Azure/aks-engine-azurestack/issues/260))
- remove usage of module autorest/azure, golint fix ([#257](https://github.com/Azure/aks-engine-azurestack/issues/257))
- remove usage of module autorest/azure ([#256](https://github.com/Azure/aks-engine-azurestack/issues/256))
- remove usage of module autorest/to, golint fix ([#255](https://github.com/Azure/aks-engine-azurestack/issues/255))
- remove usage of module autorest/to ([#254](https://github.com/Azure/aks-engine-azurestack/issues/254))
- remove container monitoring addon ([#253](https://github.com/Azure/aks-engine-azurestack/issues/253))
- use image builder for AKSe VHDs ([#248](https://github.com/Azure/aks-engine-azurestack/issues/248))
- upgrade azure sdk to use MISE and MSAL ([#247](https://github.com/Azure/aks-engine-azurestack/issues/247))
- single ARM client for Azure and Hub ([#246](https://github.com/Azure/aks-engine-azurestack/issues/246))
- Azure ARM client targets hybrid API versions ([#245](https://github.com/Azure/aks-engine-azurestack/issues/245))
- remove support for CLI & DEVICE auth methods ([#244](https://github.com/Azure/aks-engine-azurestack/issues/244))
- remove Virtual Machine ScaleSets support ([#243](https://github.com/Azure/aks-engine-azurestack/issues/243))
- update to ginkgo v2, golang v1.20, and k8s client v1.27 ([#242](https://github.com/Azure/aks-engine-azurestack/issues/242))
- remove skus, locations and update commands ([#240](https://github.com/Azure/aks-engine-azurestack/issues/240))
- grant github action bot proper write permission ([#233](https://github.com/Azure/aks-engine-azurestack/issues/233))

### Security Fix 🛡️
- Move E2E clusters to user assigned managed identity ([#241](https://github.com/Azure/aks-engine-azurestack/issues/241))
- Add advanced security policy (Github Policy Service) ([#239](https://github.com/Azure/aks-engine-azurestack/issues/239))
- Add branch protection (Github Policy Service) ([#236](https://github.com/Azure/aks-engine-azurestack/issues/236))

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.81.0...HEAD
[v0.81.0]: https://github.com/Azure/aks-engine-azurestack/compare/v0.80.2...v0.81.0
