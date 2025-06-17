
# Input
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>

# Input Validation
- Extract Kubernetes version ([MAJOR].[MINOR].[REVISION]) from xml tag <KubernetesVersion>
- Calculate target orchestratorRelease: [MAJOR].[MINOR-1] (one minor version below input)
- Validate current orchestratorRelease value in `examples/azure-stack/kubernetes-azurestack.json`:
   - if already set to `"orchestratorRelease": "[MAJOR].[MINOR-1]",`, then skip update (return False)
   - else proceed with update (return True)

# Target File Location
- source code path: `examples/azure-stack/kubernetes-azurestack.json`
- JSON properties: `orchestratorRelease` and `orchestratorVersion`

# Azure Stack Example Update Checklist

1. **Version Resolution** (from `<KubernetesVersion>` input):
   - [ ] Extract Kubernetes version ([MAJOR].[MINOR].[REVISION]) from XML tag
   - [ ] Calculate target orchestratorRelease: [MAJOR].[MINOR-1] (subtract 1 from minor version)
   - [ ] Search `AllKubernetesSupportedVersionsAzureStack` map in `pkg/api/common/versions.go`
   - [ ] Find highest patch version of the previous minor release ([MAJOR].[MINOR-1].x)

2. **JSON File Update** (in `examples/azure-stack/kubernetes-azurestack.json`):
   - [ ] Update orchestratorRelease to `[MAJOR].[MINOR-1]` format
   - [ ] Set orchestratorVersion to the latest supported patch version for the previous minor version release
   - [ ] Preserve original JSON formatting and indentation

# Version Resolution Process
1. **Input**: Extract version from `<KubernetesVersion>` tag
2. **Calculate Target Release**: [MAJOR].[MINOR-1] (subtract 1 from minor version)
3. **Find Latest Patch**: Search `AllKubernetesSupportedVersionsAzureStack` for highest patch version of the previous minor release ([MAJOR].[MINOR-1].x)
4. **Update JSON**: Set both orchestratorRelease and orchestratorVersion accordingly

# Example: Updating Azure Stack example for Kubernetes 1.30.10

**Input Analysis:**
- Input version: `1.30.10`
- Target orchestratorRelease: `1.29` (30-1=29)
- Latest supported 1.29.x version in AllKubernetesSupportedVersionsAzureStack: `1.29.15`

**Required transformation:**

**Before:**
```json
    "properties": {
        "orchestratorProfile": {
            "orchestratorRelease": "1.28",
            "orchestratorVersion": "1.28.15",
```

**After:**
```json
    "properties": {
        "orchestratorProfile": {
            "orchestratorRelease": "1.29",
            "orchestratorVersion": "1.29.15",
```

# Verification
- **All items on the Azure Stack Example Update Checklist have been verified and completed.**
