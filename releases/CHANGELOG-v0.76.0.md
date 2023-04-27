
<a name="v0.76.0"></a>
# [v0.76.0] - 2023-04-26
### Bug Fixes 🐞
- enableUnattendedUpgrades not honored ([#124](https://github.com/Azure/aks-engine-azurestack/issues/124))
- shorten custom data in cloud init files ([#121](https://github.com/Azure/aks-engine-azurestack/issues/121))
- add kube-addon-manager v9.1.6 to vhd ([#118](https://github.com/Azure/aks-engine-azurestack/issues/118))
- remove invalid k8s v1.24 flags ([#114](https://github.com/Azure/aks-engine-azurestack/issues/114))

### Continuous Integration 💜
- call test-vhd-no-egress from release github workflow ([#123](https://github.com/Azure/aks-engine-azurestack/issues/123))
- exclude version control information from test binary ([#122](https://github.com/Azure/aks-engine-azurestack/issues/122))
- Add -buildvcs=false for go build ([#120](https://github.com/Azure/aks-engine-azurestack/issues/120))
- call test-vhd-no-egress github workflow from create-release-branch ([#119](https://github.com/Azure/aks-engine-azurestack/issues/119))
- add no-egress GitHub action ([#108](https://github.com/Azure/aks-engine-azurestack/issues/108))
- release workflow tags the correct commit ([#113](https://github.com/Azure/aks-engine-azurestack/issues/113))
- gen-release-changelog wf creates branch and commit ([#112](https://github.com/Azure/aks-engine-azurestack/issues/112))
- update actions/checkout to v3 ([#111](https://github.com/Azure/aks-engine-azurestack/issues/111))

### Features 🌈
- get-logs collects etcd metrics ([#126](https://github.com/Azure/aks-engine-azurestack/issues/126))
- add support for Kubernetes v1.24.11 ([#109](https://github.com/Azure/aks-engine-azurestack/issues/109))

### Maintenance 🔧
- update Linux and Windows VHDs for April 2023 ([#129](https://github.com/Azure/aks-engine-azurestack/issues/129))
- Update Windows VHD packer job to use Apr 2023 patches ([#127](https://github.com/Azure/aks-engine-azurestack/issues/127))
- update go-dev image to v1.36.2 ([#128](https://github.com/Azure/aks-engine-azurestack/issues/128))
- update Linux and Windows VHDs for March 2023 ([#115](https://github.com/Azure/aks-engine-azurestack/issues/115))

### Testing 💚
- remove auditd from packages validate script ([#116](https://github.com/Azure/aks-engine-azurestack/issues/116))

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.76.0...HEAD
[v0.76.0]: https://github.com/Azure/aks-engine-azurestack/compare/v0.75.4...v0.76.0
