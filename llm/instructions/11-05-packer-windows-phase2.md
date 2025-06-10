

# Input 
<KubernetesVersions>{{all_supported_versions}}</KubernetesVersions>
<AzureCloudManagerImages>{{cloud_provider_image_versions}}</AzureCloudManagerImages>
<CSIImages>{{csi_image_versions}}</CSIImages>

# Input Validation
  - Extract all supported Kubernetes versions from the `<KubernetesVersions>` XML tag.
  - The versions are provided as a comma-separated string.
  - Inspect the current entries in the $map for the key `c:\akse-cache\win-k8s\`.
  - Parse the version from the value, omitting any leading 'v' (e.g., from https://packages.aks.azure.com/kubernetes/v1.29.15/windowszip/v1.29.15-1int.zip, the version is 1.29.15).
  - **Version Presence Check**:
    - If any of the desired versions are NOT present in the list, return "True".
    - If all desired versions are present in the list, return "False".
  - **Do not add code to implement the Input Validation logic.**
  
# Code Snippt Filter:
   - source code path: `vhd/packer/configure-windows-vhd-phase2.ps1`
   - object name: Get-FilesToCacheOnVHD
   - object type: func
   - begin with: `function Get-FilesToCacheOnVHD {`

# Kubernetes Binaries
  - Extract all supported Kubernetes versions from the `<KubernetesVersions>` XML tag.
  - The versions are provided as a comma-separated string.
  - Sort the versions in ascending order.
  - Remove all existing entries in the $map for the key `c:\akse-cache\win-k8s\`.
  - For each version in the sorted list, directly calculate and construct the new entry string (do not write implementation code):
    - Format: https://packages.aks.azure.com/kubernetes/v[MAJOR].[MINOR].[PATCH]/windowszip/v[MAJOR].[MINOR].[PATCH]-1int.zip
    - Example: For version 1.29.15, the URL is https://packages.aks.azure.com/kubernetes/v1.29.15/windowszip/v1.29.15-1int.zip
  - Ensure the final list is correctly formatted and properly indented.

# Examples:
Supported Versions: "1.30.10,1.29.15"

The resulting value:
    "c:\akse-cache\win-k8s\" = @(
        "https://packages.aks.azure.com/kubernetes/v1.29.15/windowszip/v1.29.15-1int.zip",
        "https://packages.aks.azure.com/kubernetes/v1.30.10/windowszip/v1.30.10-1int.zip"
    );