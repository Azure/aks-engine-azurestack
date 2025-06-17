# Input
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>

# Input Validation
- Extract Kubernetes version ([MAJOR].[MINOR].[REVISION]) from xml tag <KubernetesVersion>
- Calculate target release: [MAJOR].[MINOR] (without revision)
- Validate target file existence at `examples/kubernetes-releases/kubernetes[MAJOR].[MINOR].json`:
   - if file already exists, then skip creation (return False)
   - else proceed with creation (return True)

# Target File Location
- source code path: `examples/kubernetes-releases/`
- File naming pattern: `kubernetes[MAJOR].[MINOR].json`

# Kubernetes Release Example Creation Checklist

1. **Version Resolution** (from `<KubernetesVersion>` input):
   - [ ] Extract Kubernetes version ([MAJOR].[MINOR].[REVISION]) from XML tag
   - [ ] Calculate target release: [MAJOR].[MINOR] (exclude revision)
   - [ ] Determine previous minor version: [MAJOR].[MINOR-1]
   - [ ] Check if target file `kubernetes[MAJOR].[MINOR].json` already exists

2. **Source Template Selection**:
   - [ ] Locate previous version file: `kubernetes[MAJOR].[MINOR-1].json`
   - [ ] Copy entire contents from previous version file as template
   - [ ] Preserve original JSON formatting and indentation

3. **JSON File Creation** (in `examples/kubernetes-releases/kubernetes[MAJOR].[MINOR].json`):
   - [ ] Update orchestratorRelease to `[MAJOR].[MINOR]`
   - [ ] Preserve all other properties unchanged from template
   - [ ] Maintain original JSON formatting and indentation

# Version Resolution Process
1. **Input**: Extract version from `<KubernetesVersion>` tag
2. **Calculate Target Release**: [MAJOR].[MINOR] (exclude revision)
3. **Find Previous Template**: Look for `kubernetes[MAJOR].[MINOR-1].json` 
4. **Copy and Update**: Copy previous version file and update orchestratorRelease

# Example: Creating Kubernetes 1.31 release example

**Input Analysis:**
- Input version: `1.30.10`
- Target release: `1.30` (31 major, 31 minor)
- Previous release: `1.39` (31 major, 30 minor)
- Source template: `kubernetes1.39.json`
- Target file: `kubernetes1.30.json`

**Required transformation:**

**Source (kubernetes1.39.json):**
```json
{
    "apiVersion": "vlabs",
    "properties": {
      "orchestratorProfile": {
        "orchestratorRelease": "1.29"
      },
      "masterProfile": {
        "count": 1,
        "dnsPrefix": "",
        "vmSize": "Standard_D2_v3"
      },
      "agentPoolProfiles": [
        {
          "name": "agentpool1",
          "count": 3,
          "vmSize": "Standard_D2_v3"
        }
      ],
      "linuxProfile": {
        "adminUsername": "azureuser",
        "ssh": {
          "publicKeys": [
            {
              "keyData": ""
            }
          ]
        }
      }
    }
  }
```

**Target (kubernetes1.31.json):**
```json
{
    "apiVersion": "vlabs",
    "properties": {
      "orchestratorProfile": {
        "orchestratorRelease": "1.30"
      },
      "masterProfile": {
        "count": 1,
        "dnsPrefix": "",
        "vmSize": "Standard_D2_v3"
      },
      "agentPoolProfiles": [
        {
          "name": "agentpool1",
          "count": 3,
          "vmSize": "Standard_D2_v3"
        }
      ],
      "linuxProfile": {
        "adminUsername": "azureuser",
        "ssh": {
          "publicKeys": [
            {
              "keyData": ""
            }
          ]
        }
      }
    }
  }
```

# Verification
- **All items on the Kubernetes Release Example Creation Checklist have been verified and completed.**
- **Target file `kubernetes[MAJOR].[MINOR].json` has been created with updated orchestratorRelease value.**
- **All other properties remain unchanged from the source template.**
