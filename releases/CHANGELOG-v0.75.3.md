
<a name="v0.75.3"></a>
# [v0.75.3] - 2023-02-03

# Attention!
AKS Engine release v0.75.3, and all future AKS Engine releases on Azure Stack Hub, will be from the new [Azure/aks-engine-azurestack repository](https://github.com/Azure/aks-engine-azurestack). As such, all `aks-engine` commands should be replaced with `aks-engine-azurestack`. Please create an [issue in the new repository](https://github.com/Azure/aks-engine-azurestack/issues/new) if you find any. 

AKS Engine release v0.75.3 on Azure Stack Hub includes a new [Ubuntu 20.04-LTS VHD distro](https://github.com/Azure/aks-engine-azurestack/blob/v0.75.3/vhd/release-notes/aks-engine-ubuntu-2004/aks-engine-azurestack-ubuntu-2004_2023.032.2.txt) to use in either your control plane and/or worker node pools. Starting from this release, Ubuntu 18.04 will no longer be supported. Please refer to the section [*Upgrading Kubernetes clusters created with the Ubuntu 18.04 Distro*](https://github.com/Azure/aks-engine-azurestack/blob/0d6163211891aba81b8f84e1fd4c021ed6a3d592/docs/topics/azure-stack.md#upgrading-kubernetes-clusters-created-with-the-ubuntu-1804-distro) for more details. 


Starting from Kubernetes v1.24, only the `containerd` runtime is supported. Please refer to the section [*Upgrading Kubernetes clusters created with docker runtime*](https://github.com/Azure/aks-engine-azurestack/blob/0d6163211891aba81b8f84e1fd4c021ed6a3d592/docs/topics/azure-stack.md#upgrading-kubernetes-clusters-created-with-docker-container-runtime) for more details. For AKS Engine release v0.75.3, clusters with Windows nodes on Kubernetes v1.23 can use [the Windows base image with Docker runtime](https://github.com/Azure/aks-engine-azurestack/blob/v0.75.3/vhd/release-notes/aks-windows/2019-datacenter-core-azurestack-smalldisk-17763.3887.20230332.txt). Clusters with Windows nodes on Kubernetes v1.24 can use [the Windows base image with Containerd runtime](https://github.com/Azure/aks-engine-azurestack/blob/v0.75.3/vhd/release-notes/aks-windows-2019-containerd/2019-datacenter-core-azurestack-ctrd-17763.3887.20230332.txt).

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
