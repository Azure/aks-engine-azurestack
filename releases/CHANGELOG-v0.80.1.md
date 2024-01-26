

# Attention!

Notable changes in this release:

- The control-plane nodes' taint has been changed from node-role.kubernetes.io/master to node-role.kubernetes.io/control-plane, requiring users to update tolerations in their applications to schedule pods on these nodes. Example:

```
tolerations:
- key: node-role.kubernetes.io/control-plane
  operator: "Exists"
  effect: NoSchedule

```
- Use the new AzureDisk CSI Driver v1.29.1 for k8s v1.28+. Use AzureDisk CSI Driver v1.26.5 for k8s v1.26.
  - See [Azure Disk CSI Driver: Version Mapping](../docs/topics/azure-stack.md#azure-disk-csi-driver-version-mapping) for more details.


<a name="v0.80.1"></a>
# [v0.80.1] - 2024-01-24

### Maintenance 🔧
- Replace taint node-role.kubernetes.io/master to node-role.kubernetes.io/control-plane ([#225](https://github.com/Azure/aks-engine-azurestack/issues/225))

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.80.1...HEAD
[v0.80.1]: https://github.com/Azure/aks-engine-azurestack/compare/v0.80.0...v0.80.1



<a name="v0.80.0"></a>
# [v0.80.0] - 2024-01-19
### Build 🏭
- **deps:** bump golang.org/x/crypto from 0.14.0 to 0.17.0 ([#213](https://github.com/Azure/aks-engine-azurestack/issues/213))
- **deps:** bump golang.org/x/crypto from 0.14.0 to 0.17.0 in /test/e2e ([#212](https://github.com/Azure/aks-engine-azurestack/issues/212))

### Documentation 📘
- Azure Stack Hub doc update for v0.79.0 ([#208](https://github.com/Azure/aks-engine-azurestack/issues/208))

### Maintenance 🔧
- update Linux and Windows VHDs for January 2024 ([#222](https://github.com/Azure/aks-engine-azurestack/issues/222))
- remove unsupported feature gates for kubelet config in master and agent profile. ([#220](https://github.com/Azure/aks-engine-azurestack/issues/220))
- include Windows Server December 2023 patches ([#219](https://github.com/Azure/aks-engine-azurestack/issues/219))
- support CSI driver 1.29.1 for k8s 1.28 ([#218](https://github.com/Azure/aks-engine-azurestack/issues/218))
- support Kubernetes v1.28.5 and 1.27.9 ([#217](https://github.com/Azure/aks-engine-azurestack/issues/217))
- remove invalid k8s v1.28 flags and feature gates ([#216](https://github.com/Azure/aks-engine-azurestack/issues/216))

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.80.0...HEAD
[v0.80.0]: https://github.com/Azure/aks-engine-azurestack/compare/v0.79.0...v0.80.0
