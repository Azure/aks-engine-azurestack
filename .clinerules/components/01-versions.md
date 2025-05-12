# Update Kubernetes Version Support Maps

## Objective

Update version support maps to add a new Kubernetes version. This involves modifying maps to reflect that the new version and the immediately preceding supported version are available for new cluster creation, while older versions are not.

## Skip Conditions

Skip this step if the new Kubernetes version already exists in ALL of these maps and is configured correctly:

- `AllKubernetesSupportedVersions`
- `AllKubernetesSupportedVersionsAzureStack`
- `AllKubernetesWindowsSupportedVersionsAzureStack`

## File to Modify

- **Path:** `pkg/api/common/versions.go`

## Variables to Update

Update these map-type variables:

- `AllKubernetesSupportedVersions`
- `AllKubernetesSupportedVersionsAzureStack`
- `AllKubernetesWindowsSupportedVersionsAzureStack`

## Detailed Instructions

For each map listed in "Variables to Update":

1. **Identify Current Latest Version:**

   - Find highest version with `true` value
   - This becomes `PREVIOUS_LATEST_VERSION`

2. **Add New Version:**

   - Add `NEW_VERSION` as key
   - Set value to `true`

3. **Update Existing Versions:**
   - Keep `PREVIOUS_LATEST_VERSION` as `true`
   - Set all other versions to `false`

## Example

For adding `NEW_VERSION = "1.31.8"`:

```golang
var AllKubernetesSupportedVersions = map[string]bool{
    "1.28.5":  false,
    "1.29.10": false,
    "1.29.15": false, // Becomes false
    "1.30.10": true,  // PREVIOUS_LATEST_VERSION, remains true
    "1.31.8":  true,  // NEW_VERSION, set to true
}
```

## Validation Checks

1. Only `NEW_VERSION` and `PREVIOUS_LATEST_VERSION` should be `true`
2. All other versions should be `false`
3. Version strings should match exactly (including patch version)
4. Changes must be applied consistently across all three maps
