
<a name="v0.82.1"></a>
# [v0.82.1] - 2025-03-24
### Bug Fixes 🐞
- Fix ssh for windows nodes ([#303](https://github.com/Azure/aks-engine-azurestack/issues/303))
- Update Chrony Service Config ([#288](https://github.com/Azure/aks-engine-azurestack/issues/288))
- Add Ginkgo "timeout" param & Remove invalid azcopy link ([#287](https://github.com/Azure/aks-engine-azurestack/issues/287))
- Fix enableTelemetry apiVersion and ARM template ([#281](https://github.com/Azure/aks-engine-azurestack/issues/281))
- specify securityContext in cloud node manager ([#274](https://github.com/Azure/aks-engine-azurestack/issues/274))
- remove guest agent before sysprep to run windows cse on hub ([#273](https://github.com/Azure/aks-engine-azurestack/issues/273))
- specify instanceView to avoid nil pointer during rotate-certs ([#272](https://github.com/Azure/aks-engine-azurestack/issues/272))

### Build 🏭
- **deps:** bump golang.org/x/net from 0.33.0 to 0.36.0 ([#313](https://github.com/Azure/aks-engine-azurestack/issues/313))
- **deps:** bump github.com/golang-jwt/jwt/v5 from 5.2.1 to 5.2.2 ([#312](https://github.com/Azure/aks-engine-azurestack/issues/312))
- **deps:** bump golang.org/x/net from 0.33.0 to 0.36.0 in /test/e2e ([#300](https://github.com/Azure/aks-engine-azurestack/issues/300))
- **deps:** bump golang.org/x/crypto from 0.24.0 to 0.31.0 ([#285](https://github.com/Azure/aks-engine-azurestack/issues/285))
- **deps:** bump github.com/Azure/azure-sdk-for-go/sdk/azidentity from 1.5.2 to 1.6.0 ([#258](https://github.com/Azure/aks-engine-azurestack/issues/258))

### Continuous Integration 💜
- use of insecure HostKeyCallback implementation ([#284](https://github.com/Azure/aks-engine-azurestack/issues/284))
- update CodeQL action v3 ([#283](https://github.com/Azure/aks-engine-azurestack/issues/283))
- fix release github action ([#277](https://github.com/Azure/aks-engine-azurestack/issues/277))
- fix Hub E2E tests REQUESTS_CA_BUNDLE ([#271](https://github.com/Azure/aks-engine-azurestack/issues/271))
- Replace go-dev base image ([#264](https://github.com/Azure/aks-engine-azurestack/issues/264))
- Output AIB release notes and cgmanifest to customization.log  ([#263](https://github.com/Azure/aks-engine-azurestack/issues/263))
- Fix disconnected tests ([#262](https://github.com/Azure/aks-engine-azurestack/issues/262))
- bring back deployment error unit tests ([#251](https://github.com/Azure/aks-engine-azurestack/issues/251))
- Add SkipAzLogin and CheckIngressIPOnly parameters ([#250](https://github.com/Azure/aks-engine-azurestack/issues/250))

### Documentation 📘
- Azure Stack Hub doc update for v0.81.1 ([#280](https://github.com/Azure/aks-engine-azurestack/issues/280))
- Azure Stack Hub doc update for v0.80.2 ([#223](https://github.com/Azure/aks-engine-azurestack/issues/223))

### Maintenance 🔧
- update Azure Disk CSI driver to 1.31.5 ([#309](https://github.com/Azure/aks-engine-azurestack/issues/309))
- support Kubernetes v1.30.10 ([#297](https://github.com/Azure/aks-engine-azurestack/issues/297))
- Remove deprecated "azure-container-registry-config" flag for kubelet ([#306](https://github.com/Azure/aks-engine-azurestack/issues/306))
- update Azure CNI version to v1.4.59 ([#304](https://github.com/Azure/aks-engine-azurestack/issues/304))
- Switch to Ubuntu 22.04 as the default Linux version ([#294](https://github.com/Azure/aks-engine-azurestack/issues/294))
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
- Update ubuntu 22.04 stig config ([#298](https://github.com/Azure/aks-engine-azurestack/issues/298))
- upgrade golang.org/x/net to v0.33.0 and google.golang.org/grpc to v1.56.3 ([#295](https://github.com/Azure/aks-engine-azurestack/issues/295))
- Move E2E clusters to user assigned managed identity ([#241](https://github.com/Azure/aks-engine-azurestack/issues/241))
- Add advanced security policy (Github Policy Service) ([#239](https://github.com/Azure/aks-engine-azurestack/issues/239))
- Add branch protection (Github Policy Service) ([#236](https://github.com/Azure/aks-engine-azurestack/issues/236))

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.82.1...HEAD
[v0.82.1]: https://github.com/Azure/aks-engine-azurestack/compare/v0.80.3...v0.82.1
