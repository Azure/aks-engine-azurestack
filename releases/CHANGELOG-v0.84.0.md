
<a name="v0.84.0"></a>
# [v0.84.0] - 2025-12-08
### Bug Fixes 🐞
- no liveness probe validation needed for csi windows hp deployment ([#375](https://github.com/Azure/aks-engine-azurestack/issues/375))
- update csi snapshot images ([#372](https://github.com/Azure/aks-engine-azurestack/issues/372))
- force update cloud-provider flag to 'external' during upgrade ([#357](https://github.com/Azure/aks-engine-azurestack/issues/357))
- Remove Vnet Config for Image Builder Service for S360 ([#354](https://github.com/Azure/aks-engine-azurestack/issues/354))
- additional property update in azurestack.json ([#353](https://github.com/Azure/aks-engine-azurestack/issues/353))
- Azure CLI installation check ([#339](https://github.com/Azure/aks-engine-azurestack/issues/339))
- Update acr credential provider setup ([#329](https://github.com/Azure/aks-engine-azurestack/issues/329))
- ensure --cloud-provider is always set to external in APIServerConfig ([#333](https://github.com/Azure/aks-engine-azurestack/issues/333))
- Fix ssh for windows nodes ([#303](https://github.com/Azure/aks-engine-azurestack/issues/303))
- Update Chrony Service Config ([#288](https://github.com/Azure/aks-engine-azurestack/issues/288))
- Add Ginkgo "timeout" param & Remove invalid azcopy link ([#287](https://github.com/Azure/aks-engine-azurestack/issues/287))

### Build 🏭
- **deps:** bump golang.org/x/crypto from 0.36.0 to 0.45.0 in /test/e2e ([#365](https://github.com/Azure/aks-engine-azurestack/issues/365))
- **deps:** bump golang.org/x/oauth2 from 0.7.0 to 0.27.0 ([#332](https://github.com/Azure/aks-engine-azurestack/issues/332))
- **deps:** bump golang.org/x/net from 0.36.0 to 0.38.0 ([#324](https://github.com/Azure/aks-engine-azurestack/issues/324))
- **deps:** bump golang.org/x/net from 0.36.0 to 0.38.0 in /test/e2e ([#316](https://github.com/Azure/aks-engine-azurestack/issues/316))
- **deps:** bump golang.org/x/net from 0.33.0 to 0.36.0 ([#313](https://github.com/Azure/aks-engine-azurestack/issues/313))
- **deps:** bump github.com/golang-jwt/jwt/v5 from 5.2.1 to 5.2.2 ([#312](https://github.com/Azure/aks-engine-azurestack/issues/312))
- **deps:** bump golang.org/x/net from 0.33.0 to 0.36.0 in /test/e2e ([#300](https://github.com/Azure/aks-engine-azurestack/issues/300))
- **deps:** bump golang.org/x/crypto from 0.24.0 to 0.31.0 ([#285](https://github.com/Azure/aks-engine-azurestack/issues/285))

### Continuous Integration 💜
- use of insecure HostKeyCallback implementation ([#284](https://github.com/Azure/aks-engine-azurestack/issues/284))
- update CodeQL action v3 ([#283](https://github.com/Azure/aks-engine-azurestack/issues/283))

### Documentation 📘
- Azure Stack Hub doc update for v0.82.1 ([#315](https://github.com/Azure/aks-engine-azurestack/issues/315))
- Azure Stack Hub doc update for v0.81.1 ([#280](https://github.com/Azure/aks-engine-azurestack/issues/280))

### Maintenance 🔧
- update Linux and Windows VHD version for November 2025 ([#373](https://github.com/Azure/aks-engine-azurestack/issues/373))
- support for kubernetes 1.33.5 ([#371](https://github.com/Azure/aks-engine-azurestack/issues/371))
- support for kubernetes 1.32.9 ([#370](https://github.com/Azure/aks-engine-azurestack/issues/370))
- support for kubernetes 1.31.13 ([#369](https://github.com/Azure/aks-engine-azurestack/issues/369))
- azuredisk csi to use Windows host process image ([#368](https://github.com/Azure/aks-engine-azurestack/issues/368))
- update Windows base image version ([#367](https://github.com/Azure/aks-engine-azurestack/issues/367))
- update Linux and Windows VHDs for August 2025 ([#342](https://github.com/Azure/aks-engine-azurestack/issues/342))
- Support Kubernetes 1.30.14 and 1.31.11 ([#328](https://github.com/Azure/aks-engine-azurestack/issues/328))
- move Linux files from cse to vhd ([#335](https://github.com/Azure/aks-engine-azurestack/issues/335))
- Support for Azure Container Registry (ACR) Credential Provider ([#326](https://github.com/Azure/aks-engine-azurestack/issues/326))
- update Azure Disk CSI driver to 1.31.5 ([#309](https://github.com/Azure/aks-engine-azurestack/issues/309))
- support Kubernetes v1.30.10 ([#297](https://github.com/Azure/aks-engine-azurestack/issues/297))
- Remove deprecated "azure-container-registry-config" flag for kubelet ([#306](https://github.com/Azure/aks-engine-azurestack/issues/306))
- update Azure CNI version to v1.4.59 ([#304](https://github.com/Azure/aks-engine-azurestack/issues/304))
- Switch to Ubuntu 22.04 as the default Linux version ([#294](https://github.com/Azure/aks-engine-azurestack/issues/294))

### Security Fix 🛡️
- security fixes for Ubuntu and Windows Base Images ([#325](https://github.com/Azure/aks-engine-azurestack/issues/325))
- Update ubuntu 22.04 stig config ([#298](https://github.com/Azure/aks-engine-azurestack/issues/298))
- upgrade golang.org/x/net to v0.33.0 and google.golang.org/grpc to v1.56.3 ([#295](https://github.com/Azure/aks-engine-azurestack/issues/295))

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.84.0...HEAD
[v0.84.0]: https://github.com/Azure/aks-engine-azurestack/compare/release-v0.81.2...v0.84.0
