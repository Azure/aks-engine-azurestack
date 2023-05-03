
<a name="v0.76.0"></a>
# [v0.76.0] - 2023-04-26
### Bug Fixes 🐞
- enableUnattendedUpgrades not honored ([#124](https://github.com/Azure/aks-engine-azurestack/issues/124))
- shorten custom data in cloud init files ([#121](https://github.com/Azure/aks-engine-azurestack/issues/121))
- add kube-addon-manager v9.1.6 to vhd ([#118](https://github.com/Azure/aks-engine-azurestack/issues/118))
- remove invalid k8s v1.24 flags ([#114](https://github.com/Azure/aks-engine-azurestack/issues/114))
- use cross-platform pause image as the containerd sandbox image on Windows ([#106](https://github.com/Azure/aks-engine-azurestack/issues/106))
- enforce each addon manager pod ([#99](https://github.com/Azure/aks-engine-azurestack/issues/99))
- kubernetes-azurestack.json uses distro aks-ubuntu-20.04 ([#87](https://github.com/Azure/aks-engine-azurestack/issues/87))
- get-akse.sh pulls using the correct file name ([#84](https://github.com/Azure/aks-engine-azurestack/issues/84))

### Documentation 📘
- clarify that the "azuredisk-csi-driver" addon works now on both Linux and Windows nodes ([#96](https://github.com/Azure/aks-engine-azurestack/issues/96))

### Continuous Integration 💜
- call test-vhd-no-egress from release github workflow ([#123](https://github.com/Azure/aks-engine-azurestack/issues/123))
- exclude version control information from test binary ([#122](https://github.com/Azure/aks-engine-azurestack/issues/122))
- Add -buildvcs=false for go build ([#120](https://github.com/Azure/aks-engine-azurestack/issues/120))
- call test-vhd-no-egress github workflow from create-release-branch ([#119](https://github.com/Azure/aks-engine-azurestack/issues/119))
- add no-egress GitHub action ([#108](https://github.com/Azure/aks-engine-azurestack/issues/108))
- release workflow tags the correct commit ([#113](https://github.com/Azure/aks-engine-azurestack/issues/113))
- gen-release-changelog wf creates branch and commit ([#112](https://github.com/Azure/aks-engine-azurestack/issues/112))
- update actions/checkout to v3 ([#111](https://github.com/Azure/aks-engine-azurestack/issues/111))
- chocolatey workflow ([#86](https://github.com/Azure/aks-engine-azurestack/issues/86))
- release workflows run no-egress scenarios ([#85](https://github.com/Azure/aks-engine-azurestack/issues/85))

### Features 🌈
- get-logs collects etcd metrics ([#126](https://github.com/Azure/aks-engine-azurestack/issues/126))
- add support for Kubernetes v1.24.11 ([#109](https://github.com/Azure/aks-engine-azurestack/issues/109))
- migrate from Pod Security Policy to Pod Security admission ([#94](https://github.com/Azure/aks-engine-azurestack/issues/94))
- DISA Ubuntu 20.04 STIG compliance ([#83](https://github.com/Azure/aks-engine-azurestack/issues/83))

### Maintenance 🔧
- update Linux and Windows VHDs for April 2023 ([#129](https://github.com/Azure/aks-engine-azurestack/issues/129))
- Update Windows VHD packer job to use Apr 2023 patches ([#127](https://github.com/Azure/aks-engine-azurestack/issues/127))
- update go-dev image to v1.36.2 ([#128](https://github.com/Azure/aks-engine-azurestack/issues/128))
- update Linux and Windows VHDs for March 2023 ([#115](https://github.com/Azure/aks-engine-azurestack/issues/115))
- support Kubernetes v1.25.7 ([#105](https://github.com/Azure/aks-engine-azurestack/issues/105))
- upgrade coredns to v1.9.4 ([#98](https://github.com/Azure/aks-engine-azurestack/issues/98))
- replace usage of deprecated "io/ioutil" golang package ([#97](https://github.com/Azure/aks-engine-azurestack/issues/97))
- upgrade containerd to 1.5.16 ([#95](https://github.com/Azure/aks-engine-azurestack/issues/95))
- upgrade pause to v3.8 ([#93](https://github.com/Azure/aks-engine-azurestack/issues/93))
- update golang toolchain to v1.19 ([#90](https://github.com/Azure/aks-engine-azurestack/issues/90))
- update registries for nvidia and k8s.io components ([#88](https://github.com/Azure/aks-engine-azurestack/issues/88))
- remove package apache2-utils from VHD ([#82](https://github.com/Azure/aks-engine-azurestack/issues/82))

### Security Fix 🛡️
- bump x/net and x/crypto ([#104](https://github.com/Azure/aks-engine-azurestack/issues/104))

### Testing 💚
- remove auditd from packages validate script ([#116](https://github.com/Azure/aks-engine-azurestack/issues/116))
- e2e suite validates an existing PV works after a cluster upgrade ([#92](https://github.com/Azure/aks-engine-azurestack/issues/92))

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.76.0...HEAD
[v0.76.0]: https://github.com/Azure/aks-engine-azurestack/compare/v0.75.4...v0.76.0
