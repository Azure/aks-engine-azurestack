
<a name="v0.75.4"></a>

# [v0.75.4] - 2023-03-09

### Bug Fixes 🐞
- kubernetes-azurestack.json uses distro aks-ubuntu-20.04 ([#87](https://github.com/Azure/aks-engine-azurestack/issues/87))
- use cross-platform pause image as the containerd sandbox image on Windows ([#106](https://github.com/Azure/aks-engine-azurestack/issues/106))

### Continuous Integration 💜
- release workflow creates pre-release
- release workflows run no-egress scenarios ([#85](https://github.com/Azure/aks-engine-azurestack/issues/85))

### Documentation 📘
- Enable azurestack-csi-driver addon for mixed clusters ([#96](https://github.com/Azure/aks-engine-azurestack/issues/96))
- Azure Stack Hub doc update for v0.75.3 ([#81](https://github.com/Azure/aks-engine-azurestack/issues/81))

### Testing 💚
- e2e suite validates an existing PV works after a cluster upgrade ([#92](https://github.com/Azure/aks-engine-azurestack/issues/92))

#### Please report any issues here: https://github.com/Azure/aks-engine-azurestack/issues/new
[v0.75.4]: https://github.com/Azure/aks-engine-azurestack/compare/v0.75.3...v0.75.4
