
# Input
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>

# Input Validation
- Extract Kubernetes version ([MAJOR].[MINOR].[REVISION]) from xml tag <KubernetesVersion>
- Validate current KUBECTL_VERSION value:
   - if already set to `KUBECTL_VERSION ?= v[MAJOR].[MINOR].[REVISION]`, then skip update (return False)
   - else proceed with update (return True)

# Target File Location:
- source code path: `hack/tools/Makefile`
- variable name: KUBECTL_VERSION

# Makefile Update Checklist

1. KUBECTL_VERSION Variable Update (in `hack/tools/Makefile`):
   - [ ] Update KUBECTL_VERSION to `v[MAJOR].[MINOR].[REVISION]`
   - [ ] Preserve the original code's spacing and formatting
   - [ ] Maintain the variable assignment format. Example: `KUBECTL_VERSION ?= v1.30.10`

# Example: Updating KUBECTL_VERSION to 1.30.10

Required transformation when updating kubectl version to 1.30.10:

**Before:**
```makefile
SHELLCHECK_VERSION ?= v0.8.0
AZCLI_VERSION   ?= 2.56.0
PYWINRM_VERSION ?= 0.4.3
KUBECTL_VERSION ?= v1.29.12
```

**After:**
```makefile
SHELLCHECK_VERSION ?= v0.8.0
AZCLI_VERSION   ?= 2.56.0
PYWINRM_VERSION ?= 0.4.3
KUBECTL_VERSION ?= v1.30.10
```

# Verification
- **All items on the Makefile Update Checklist have been verified and completed.**
