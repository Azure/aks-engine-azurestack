# Input 
<KubernetesVersions>{{k8s_version}}</KubernetesVersions>
<AzureCloudManagerImages>{{cloud_provider_image_versions}}</AzureCloudManagerImages>
<CSIImages>{{csi_image_versions}}</CSIImages>

# Input Validation
  - **Do not add code to implement the Input Validation logic.**
  - Extract all supported Kubernetes versions from the `<KubernetesVersions>` XML tag.
  - The versions are provided as a comma-separated string.
  - Inspect the current entries in the $map for the key `c:\akse-cache\win-k8s\`.
  - Parse the version from the value, omitting any leading 'v' (e.g., from https://packages.aks.azure.com/kubernetes/v1.29.15/windowszip/v1.29.15-1int.zip, the version is 1.29.15).
  - **Version Presence Check**:
    - If any of the desired versions are NOT present in the list, return "True".
    - If all desired versions are present in the list, return "False".
  
# Code snippet Filter:
   - source code path: `vhd/packer/configure-windows-vhd-phase2.ps1`
   - object name: Get-FilesToCacheOnVHD
   - object type: func   - begin with: `function Get-FilesToCacheOnVHD {`

# Kubernetes Binaries  - Extract all supported Kubernetes versions from the `<KubernetesVersions>` XML tag.
  - The versions are provided as a comma-separated string.
  - Sort the versions in ascending order.
  - Replace the existing entries in the $map for the key `c:\akse-cache\win-k8s\` with the new supported versions.
  - For each version in the sorted list, directly calculate and construct the new entry string (do not write implementation code):
    - Format: https://packages.aks.azure.com/kubernetes/v[MAJOR].[MINOR].[PATCH]/windowszip/v[MAJOR].[MINOR].[PATCH]-1int.zip
    - Example: For version 1.29.15, the URL is https://packages.aks.azure.com/kubernetes/v1.29.15/windowszip/v1.29.15-1int.zip
  - Ensure the final list is correctly formatted and properly indented.

## Newline Preservation Guidelines for String Replacement

**CRITICAL**: To prevent accidental line merging when using `replace_string_in_file`, follow this strategy:

### Target Single Line Only
Replace only the specific line that needs to change:

```powershell
# CORRECT - Replace only the target line:
oldString: "        \"https://packages.aks.azure.com/kubernetes/v1.28.12/windowszip/v1.28.12-1int.zip\","
newString: "        \"https://packages.aks.azure.com/kubernetes/v1.30.10/windowszip/v1.30.10-1int.zip\","
```

# Examples:
For New Kuberntes version: "1.30.10"

Existing entries in $map for `c:\akse-cache\win-k8s\`:
    "c:\akse-cache\win-k8s\" = @(
        "https://packages.aks.azure.com/kubernetes/v1.28.12/windowszip/v1.28.12-1int.zip",
        "https://packages.aks.azure.com/kubernetes/v1.29.8/windowszip/v1.29.8-1int.zip"
    );

The resulting value (replacing existing versions with new supported versions):
    "c:\akse-cache\win-k8s\" = @(
        "https://packages.aks.azure.com/kubernetes/v1.30.10/windowszip/v1.30.10-1int.zip",
        "https://packages.aks.azure.com/kubernetes/v1.29.8/windowszip/v1.29.8-1int.zip"
    );

**IMPORTANT FORMATTING NOTE**: When performing the replacement, ensure that:
1. Each key stays on its own separate line
2. Proper indentation (8 spaces) is maintained for each entry
3. No lines are accidentally merged together during the edit process
4. The comma placement and string quotes remain consistent
5. Each URL entry is properly formatted and aligned

**You must review and ensure that the new key entries are properly formatted with each entry on its own separate line with correct indentation.**

**After making changes, you MUST verify that all URL entries in the array remain on separate lines and maintain proper indentation and formatting.**