
<a name="v0.79.0"></a>
# [v0.79.0] - 2023-10-23
### Bug Fixes 🐞
- shorten custom data in cloud init files and azure disk csi driver addon ([#205](https://github.com/Azure/aks-engine-azurestack/issues/205))
- remove LegacyServiceAccountTokenNoAutoGeneration for k8s v1.27 ([#200](https://github.com/Azure/aks-engine-azurestack/issues/200))
- Remove unused bridge network ([#183](https://github.com/Azure/aks-engine-azurestack/issues/183))

### Build 🏭
- **deps:** bump golang.org/x/net from 0.10.0 to 0.17.0 ([#199](https://github.com/Azure/aks-engine-azurestack/issues/199))
- **deps:** bump golang.org/x/net from 0.10.0 to 0.17.0 in /test/e2e ([#198](https://github.com/Azure/aks-engine-azurestack/issues/198))

### Documentation 📘
- clarify cluster-autoscaler support on Azure Stack Hub ([#202](https://github.com/Azure/aks-engine-azurestack/issues/202))
- Add seccomp profile, csi driver, and dualstack documentation ([#196](https://github.com/Azure/aks-engine-azurestack/issues/196))
- Azure Stack Hub doc update for v0.78.0 ([#189](https://github.com/Azure/aks-engine-azurestack/issues/189))

### Features 🌈
- Enable seccomp profile defaulting ([#193](https://github.com/Azure/aks-engine-azurestack/issues/193))

### Maintenance 🔧
- update Linux and Windows VHDs for October 2023 ([#206](https://github.com/Azure/aks-engine-azurestack/issues/206))
- rotate-certs creates its own known_hosts copy ([#204](https://github.com/Azure/aks-engine-azurestack/issues/204))
- include Windows Server October 2023 patches ([#203](https://github.com/Azure/aks-engine-azurestack/issues/203))
- support Kubernetes v1.27.6 ([#195](https://github.com/Azure/aks-engine-azurestack/issues/195))
- remove invalid k8s v1.27 flags and feature gates ([#197](https://github.com/Azure/aks-engine-azurestack/issues/197))
- support Kubernetes v1.26.9 ([#194](https://github.com/Azure/aks-engine-azurestack/issues/194))
- set MTU to 1500 for the kubenet CNI ([#192](https://github.com/Azure/aks-engine-azurestack/issues/192))
- upgrade azuredisk csi driver to v1.28.3 ([#190](https://github.com/Azure/aks-engine-azurestack/issues/190))
- new Windows signed scripts and package version ([#176](https://github.com/Azure/aks-engine-azurestack/issues/176))

### Testing 💚
- e2e validates Windows HostProcess pods ([#182](https://github.com/Azure/aks-engine-azurestack/issues/182))

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.79.0...HEAD
[v0.79.0]: https://github.com/Azure/aks-engine-azurestack/compare/v0.78.0...v0.79.0
