
# Attention!
Notable changes in this release:
- The `chrony` daemon on a Linux node may fail to restart with error message: `"Could not open /dev/ptp_hyperv: No such file or directory"`.
The workaround for this issue is to manually reboot the affected Linux nodes.
This operation will regenerate the `/dev/ptp_hyperv` symlink and allow the chrony daemon to restart successfully.

- The control-plane nodes' taint has been changed from node-role.kubernetes.io/master to node-role.kubernetes.io/control-plane, requiring users to update tolerations in their applications to schedule pods on these nodes. Example:


```
tolerations:
- key: node-role.kubernetes.io/control-plane
  operator: "Exists"
  effect: NoSchedule

```
- Use the new AzureDisk CSI Driver v1.29.1 for k8s v1.28+. Use AzureDisk CSI Driver v1.28.3 for k8s v1.27. Use AzureDisk CSI Driver v1.26.5 for k8s v1.26.
  - See [Azure Disk CSI Driver: Version Mapping](../docs/topics/azure-stack.md#azure-disk-csi-driver-version-mapping) for more details.

<a name="v0.80.3"></a>
# [v0.80.3] - 2025-02-06
### Bug Fixes 🐞
- Update Chrony Service Config ([#288](https://github.com/Azure/aks-engine-azurestack/issues/288))

### Maintenance 🔧
- update change log for v0.80.2
- grant github action bot proper write permission ([#233](https://github.com/Azure/aks-engine-azurestack/issues/233))
- update the change log for v0.80.2

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[Unreleased]: https://github.com/Azure/aks-engine-azurestack/compare/v0.80.3...HEAD
[v0.80.3]: https://github.com/Azure/aks-engine-azurestack/compare/v0.81.1...v0.80.3
