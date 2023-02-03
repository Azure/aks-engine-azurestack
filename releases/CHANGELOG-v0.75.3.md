
<a name="v0.75.3"></a>
# [v0.75.3] - 2023-02-03
### Bug Fixes 🐞
- CoreDNS image not updated after cluster upgrade ([#75](https://github.com/Azure/aks-engine-azurestack/issues/75))
- change reference of cni config to scripts dir ([#71](https://github.com/Azure/aks-engine-azurestack/issues/71))
- add Azure CNI config script to Ubuntu VHD ([#70](https://github.com/Azure/aks-engine-azurestack/issues/70))
- unit test checking e2e configs ([#49](https://github.com/Azure/aks-engine-azurestack/issues/49))
- syntax error in Windows VHD script ([#41](https://github.com/Azure/aks-engine-azurestack/issues/41))
- ensure eth0 addr is set to NIC's primary addr ([#39](https://github.com/Azure/aks-engine-azurestack/issues/39))

### Continuous Integration 💜
- remove no-egress job from create branch action ([#79](https://github.com/Azure/aks-engine-azurestack/issues/79))
- PR gate runs E2E suite ([#69](https://github.com/Azure/aks-engine-azurestack/issues/69))
- PR checks consume SIG images ([#64](https://github.com/Azure/aks-engine-azurestack/issues/64))
- E2E PR check uses user assigned identity ([#54](https://github.com/Azure/aks-engine-azurestack/issues/54))
- fix variable name in e2e PR check ([#52](https://github.com/Azure/aks-engine-azurestack/issues/52))
- e2e PR check sets tenant ([#51](https://github.com/Azure/aks-engine-azurestack/issues/51))
- e2e PR check does not use AvailabilitySets ([#47](https://github.com/Azure/aks-engine-azurestack/issues/47))
- e2e PR check does not use custom VNET ([#46](https://github.com/Azure/aks-engine-azurestack/issues/46))
- e2e PR check does not use MSI ([#45](https://github.com/Azure/aks-engine-azurestack/issues/45))

### Documentation 📘
- remove Azure as a target cloud ([#43](https://github.com/Azure/aks-engine-azurestack/issues/43))
- rename binary name in all markdown files ([#42](https://github.com/Azure/aks-engine-azurestack/issues/42))

### Maintenance 🔧
- update default windows image to jan 2023 ([#77](https://github.com/Azure/aks-engine-azurestack/issues/77))
- Update Windows VHD packer job to use Jan 2023 patches ([#76](https://github.com/Azure/aks-engine-azurestack/issues/76))
- change base image sku and version to azurestack ([#74](https://github.com/Azure/aks-engine-azurestack/issues/74))
- set fsType to ext4 in supported storage classes ([#73](https://github.com/Azure/aks-engine-azurestack/issues/73))
- enable v1.23.15 & v1.24.9, use ubuntu 20.04 as default, force containerd runtime ([#68](https://github.com/Azure/aks-engine-azurestack/issues/68))
- include relevant updates from v0.75.0 ([#56](https://github.com/Azure/aks-engine-azurestack/issues/56))
- include relevant updates from v0.74.0 ([#55](https://github.com/Azure/aks-engine-azurestack/issues/55))
- include relevant updates from v0.73.0 ([#53](https://github.com/Azure/aks-engine-azurestack/issues/53))
- remove kv-fluxvolume addon ([#48](https://github.com/Azure/aks-engine-azurestack/issues/48))
- prefer ADO for PR E2E check ([#38](https://github.com/Azure/aks-engine-azurestack/issues/38))
- added e2e to PR workflow ([#36](https://github.com/Azure/aks-engine-azurestack/issues/36))

### Testing 💚
- e2e sets ImageRef in all linux nodepools ([#65](https://github.com/Azure/aks-engine-azurestack/issues/65))

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.75.3...HEAD
[v0.75.3]: https://github.com/Azure/aks-engine-azurestack/compare/v0.71.1...v0.75.3
