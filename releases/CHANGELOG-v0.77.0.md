
<a name="v0.77.0"></a>
# [v0.77.0] - 2023-07-26
### Bug Fixes 🐞
- update WindowsContainerdURL during upgrade and scale ([#167](https://github.com/Azure/aks-engine-azurestack/issues/167))
- disable powershell progress bar during ssh to collect windows logs ([#166](https://github.com/Azure/aks-engine-azurestack/issues/166))
- stop collecting docker engine logs for Windows nodes ([#153](https://github.com/Azure/aks-engine-azurestack/issues/153))
- upgrade azure network policy manager to v1.4.32 ([#151](https://github.com/Azure/aks-engine-azurestack/issues/151))
- persist runUnattendedUpgradesOnBootstrap ([#142](https://github.com/Azure/aks-engine-azurestack/issues/142))
- addpool generates prefixes based on pool count ([#138](https://github.com/Azure/aks-engine-azurestack/issues/138))

### Build 🏭
- produce cgmanifest.json per UbuntuVHD build ([#156](https://github.com/Azure/aks-engine-azurestack/issues/156))

### Continuous Integration 💜
- use cloud controller manager in disconnected pipeline ([#165](https://github.com/Azure/aks-engine-azurestack/issues/165))
- fix choco package release action ([#137](https://github.com/Azure/aks-engine-azurestack/issues/137))

### Documentation 📘
- how to rotate spn credentials ([#163](https://github.com/Azure/aks-engine-azurestack/issues/163))
- how to renew front-proxy certs ([#157](https://github.com/Azure/aks-engine-azurestack/issues/157))
- Azure Stack Hub doc update for v0.76.0 ([#133](https://github.com/Azure/aks-engine-azurestack/issues/133))

### Features 🌈
- DISA Kubernetes STIG compliance ([#143](https://github.com/Azure/aks-engine-azurestack/issues/143))

### Maintenance 🔧
- update Linux and Windows VHDs for July 2023 ([#170](https://github.com/Azure/aks-engine-azurestack/issues/170))
- upgrade CSI snapshot to v5.0.1, get node info from labels, disable zones support ([#164](https://github.com/Azure/aks-engine-azurestack/issues/164))
- Update Windows VHD packer job to use July 2023 patches ([#162](https://github.com/Azure/aks-engine-azurestack/issues/162))
- upgrade azuredisk-csi-driver v1.26.5 components ([#161](https://github.com/Azure/aks-engine-azurestack/issues/161))
- support Kubernetes v1.26.6 ([#160](https://github.com/Azure/aks-engine-azurestack/issues/160))
- upgrade containerd to v1.6.21 and runc to v1.1.7 ([#159](https://github.com/Azure/aks-engine-azurestack/issues/159))
- upgrade azuredisk-csi-driver to v1.26.5 ([#158](https://github.com/Azure/aks-engine-azurestack/issues/158))
- remove private preview warning for AzureCNI ([#155](https://github.com/Azure/aks-engine-azurestack/issues/155))
- update ip-masq-agent security context ([#144](https://github.com/Azure/aks-engine-azurestack/issues/144))
- fix Windows VHD version for April 2023 ([#131](https://github.com/Azure/aks-engine-azurestack/issues/131))

### Security Fix 🛡️
- use of insecure HostKeyCallback implementation ([#149](https://github.com/Azure/aks-engine-azurestack/issues/149))
- incorrect conversion between integer types ([#147](https://github.com/Azure/aks-engine-azurestack/issues/147))
- create CodeQL action ([#145](https://github.com/Azure/aks-engine-azurestack/issues/145))

### Testing 💚
- E2E suite pulls busybox from MCR ([#154](https://github.com/Azure/aks-engine-azurestack/issues/154))

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.77.0...HEAD
[v0.77.0]: https://github.com/Azure/aks-engine-azurestack/compare/v0.76.0...v0.77.0
