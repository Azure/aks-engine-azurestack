
<a name="v0.85.1"></a>
# [v0.85.1] - 2026-11-12
### Bug Fixes 🐞
- remove obsolete crypto experiment from E2E builds ([#417](https://github.com/Azure/aks-engine-azurestack/issues/417))
- downgrade Windows Server image version since AIB image list is not updated ([#415](https://github.com/Azure/aks-engine-azurestack/issues/415))
- remediate eMudhra CA compromise in Go binaries and VM images ([#413](https://github.com/Azure/aks-engine-azurestack/issues/413))
- update package name and scripts for new test debian image ([#409](https://github.com/Azure/aks-engine-azurestack/issues/409))
- install Azure Stack CA bundle as one cert per file ([#408](https://github.com/Azure/aks-engine-azurestack/issues/408))
- kube-apiserver cloud provider flags removal for k8s 1.33+ ([#397](https://github.com/Azure/aks-engine-azurestack/issues/397))
- update golang and image versions ([#391](https://github.com/Azure/aks-engine-azurestack/issues/391))
- additional unsupported flags to remove ([#389](https://github.com/Azure/aks-engine-azurestack/issues/389))
- update k8s components image reference paths to v2 ([#388](https://github.com/Azure/aks-engine-azurestack/issues/388))
- Update Token Audience Configuration for Azure Stack ([#379](https://github.com/Azure/aks-engine-azurestack/issues/379))

### Build 🏭
- **deps:** bump the github-actions group across 1 directory with 6 updates ([#410](https://github.com/Azure/aks-engine-azurestack/issues/410))
- **deps:** bump golang.org/x/crypto from 0.51.0 to 0.52.0 in /test/e2e ([#403](https://github.com/Azure/aks-engine-azurestack/issues/403))
- **deps:** bump golang.org/x/crypto from 0.51.0 to 0.52.0 ([#404](https://github.com/Azure/aks-engine-azurestack/issues/404))
- **deps:** bump golang.org/x/net from 0.47.0 to 0.55.0 ([#402](https://github.com/Azure/aks-engine-azurestack/issues/402))
- **deps:** bump golang.org/x/net from 0.47.0 to 0.55.0 in /test/e2e ([#401](https://github.com/Azure/aks-engine-azurestack/issues/401))
- **deps:** bump github.com/sirupsen/logrus in /test/e2e ([#378](https://github.com/Azure/aks-engine-azurestack/issues/378))
- **deps:** bump github.com/sirupsen/logrus from 1.8.0 to 1.8.3 ([#374](https://github.com/Azure/aks-engine-azurestack/issues/374))
- **deps:** bump golang.org/x/crypto from 0.36.0 to 0.45.0 ([#366](https://github.com/Azure/aks-engine-azurestack/issues/366))

### Maintenance 🔧
- update Linux and Windows VHDs for October 2026 ([#416](https://github.com/Azure/aks-engine-azurestack/issues/416))
- Update windows base image for AKSe Windows VHD ([#414](https://github.com/Azure/aks-engine-azurestack/issues/414))
- change k8s 1.34 patch version to 1.34.4 ([#387](https://github.com/Azure/aks-engine-azurestack/issues/387))
- support for k8s 1.35.0 ([#386](https://github.com/Azure/aks-engine-azurestack/issues/386))
- support for kubernetes 1.34.6 ([#385](https://github.com/Azure/aks-engine-azurestack/issues/385))
- support for kubernetes 1.33.10 ([#384](https://github.com/Azure/aks-engine-azurestack/issues/384))
- Update windows base image for AKSe Windows VHD ([#383](https://github.com/Azure/aks-engine-azurestack/issues/383))

### Testing 💚
- update test for k8s 1.35 ([#390](https://github.com/Azure/aks-engine-azurestack/issues/390))

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.85.1...HEAD
[v0.85.1]: https://github.com/Azure/aks-engine-azurestack/compare/release-v0.84.2...v0.85.1
