# Add Unit Tests to Verify Feature Gate Removal

## Input Requirements

1. **Kubernetes Version**

   - The specific Kubernetes version from which the feature gates should be removed
   - Format: "major.minor.patch" (e.g., "1.31.0")
   - NEW_VERSION is the target version for feature gate removal
   - PREVIOUS_VERSION is automatically calculated as (NEW_VERSION_MAJOR).(NEW_VERSION_MINOR-1).0
     Example:
     - If NEW_VERSION = "1.31.0"
     - Then PREVIOUS_VERSION = "1.30.0"

2. **Feature Gates List**
   - REQUIRED: Must be provided as input, no hardcoding allowed
   - Must be provided in one of these formats:
     - JSON array of feature gate names: `["FeatureGateA", "FeatureGateB"]`
     - Comma-separated string: `"FeatureGateA=true,FeatureGateB=true"`
   - No assumptions about feature gates should be made in the code
   - Feature gates list must be validated before use

## Input Validation

1. **Feature Gates Input Validation**

   - Must convert JSON array input to the required format with "=true" appended
   - Must validate comma-separated string format
   - Must store input in a variable for consistent reuse
   - Example conversion:
     - Input: ["Featuregate01", "Featuregate02"]
     - Stored: "Featuregate01=true,Featuregate02=true"

2. **Version Format Validation**
   - Must use ".0" suffix for all version numbers
   - Must calculate previous version correctly
   - Must validate version string format

## Test Organization

1. **Test Placement**

   - New tests MUST be added at the end of each component's test function
   - Group tests by version and functionality
   - Follow existing test naming and comment patterns

2. **Test Order**

   - Tests should be organized from oldest to newest versions
   - Feature gate removal tests for NEW_VERSION should be the last tests in the function
   - Previous version compatibility tests should immediately precede NEW_VERSION tests

3. **Comments and Documentation**

   - Each test section must start with a clear comment indicating version and purpose
   - Use standard comment format: `// test user-overrides, <purpose> for k8s versions <condition>`
   - Maintain blank lines between test cases for readability
   - Include inline comments explaining key steps and expectations

4. **Test Case Validation**
   - Before adding new test cases:
     - Verify no existing test for the same version exists
     - Ensure previous version test exists or is being added
     - Check that feature gates input is properly validated
   - After adding test cases:
     - Verify placement at end of test function
     - Confirm both positive (NEW_VERSION) and negative (PREVIOUS_VERSION) cases are present
     - Validate error messages include all required information

## Files to Modify

1. `pkg/api/defaults-apiserver_test.go`
2. `pkg/api/defaults-cloud-controller-manager_test.go`
3. `pkg/api/defaults-controller-manager_test.go`
4. `pkg/api/defaults-kubelet_test.go`
5. `pkg/api/defaults-scheduler_test.go`

## Error Message Format

Error messages must follow this template:

```
"got unexpected '--feature-gates' for %s \n component config: \n  original value: %s \n  actual sanitized value: %s \n  expected sanitized value: %s"
```

This provides:

1. Version context (%s = version string)
2. Original input (feature gates string)
3. Actual result (what we got)
4. Expected result (what we wanted)

## Test Pattern

### **Identify the Test Method**

- Locate the existing test method whose name matches the pattern: `Test*FeatureGates`
  - The method name must **start with** `Test` and **end with** `FeatureGates`.
  - **Examples:**
    - `TestAPIServerFeatureGates`
    - `TestCloudControllerManagerFeatureGates`

---

### Append Code

- **Insert the specified code snippet at the end of the located test method, just before the closing brace.**

### Variable Naming Conventions

- **Feature Gates Variable for New Version:**

  - Name: `featuregate<NEW_VERSION_MAJOR><NEW_VERSION_MINOR>`
  - Example: For Kubernetes version `1.31`, use `featuregate131`.

- **Feature Gates Variable for New Version (Sanitized):**

  - Name: `featuregate<NEW_VERSION_MAJOR><NEW_VERSION_MINOR>Sanitized`
  - Example: For Kubernetes version `1.31`, use `featuregate131Sanitized`.

- **Feature Gates Variable for Previous Version (Sanitized):**

  - Name: `featuregate<PREVIOUS_VERSION_MAJOR><PREVIOUS_VERSION_MINOR>Sanitized`
  - Example: For Kubernetes version `1.31` (previous version is `1.30`), use `featuregate130Sanitized`.

- **Reuse** the variable for the previous version's sanitized feature gates if it already exists.

```go
// test user-overrides, removal of feature gates for k8s versions >= 1.31
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.31.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	featuregate131 := "<FEATURE_GATES_STRING>"
	ccm["--feature-gates"] = featuregate131
	featuregate131Sanitized := ""
	cs.set<Component>Config()
	if ccm["--feature-gates"] != featuregate131Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n Cloud Controller Manager config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.31.0", featuregate131, ccm["--feature-gates"], featuregate131Sanitized)
	}

	// test user-overrides, no removal of feature gates for k8s versions < 1.31
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.30.0"
	cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig = make(map[string]string)
	ccm = cs.Properties.OrchestratorProfile.KubernetesConfig.CloudControllerManagerConfig
	ccm["--feature-gates"] = featuregate131
	featuregate130Sanitized = featuregate131
	cs.set<Component>Config()
	if ccm["--feature-gates"] != featuregate130Sanitized {
		t.Fatalf("got unexpected '--feature-gates' for %s \n Cloud Controller Manager config original value  %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
			"1.30.0", featuregate131, ccm["--feature-gates"], featuregate130Sanitized)
	}
```

## Component-Specific Adjustments

Replace placeholders in the test pattern:

1. **APIServer:**

   - `<ComponentConfig>` = `APIServerConfig`
   - `<Component>` = `APIServer`

2. **CloudControllerManager:**

   - `<ComponentConfig>` = `CloudControllerManagerConfig`
   - `<Component>` = `CloudControllerManager`

3. **ControllerManager:**

   - `<ComponentConfig>` = `ControllerManagerConfig`
   - `<Component>` = `ControllerManager`

4. **Kubelet:**

   - `<ComponentConfig>` = `KubeletConfig`
   - `<Component>` = `Kubelet`

5. **Scheduler:**
   - `<ComponentConfig>` = `SchedulerConfig`
   - `<Component>` = `Scheduler`

## Validation Checks

1. **Test Location**

   - Verify tests are placed at the end of test functions
   - Ensure logical grouping with related tests

2. **Version Handling**

   - Tests use correct version strings
   - Both NEW_VERSION and PREVIOUS_VERSION scenarios are covered

3. **Feature Gates**

   - Feature gates string comes from input
   - No hardcoded feature gate assumptions
   - Proper handling of empty feature gates string

4. **Error Messages**

   - Include clear version information
   - Show original and expected values
   - Provide context about the test purpose

5. **Component Configuration**
   - Correct component-specific setup
   - Proper initialization of config maps
