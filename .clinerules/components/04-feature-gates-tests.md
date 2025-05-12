# Add Unit Tests to Verify Feature Gate Removal

## Objective

Add unit tests across all components to verify feature gate removal behavior for NEW_VERSION. These tests will verify:

1. Feature gates are removed when k8s version >= NEW_VERSION's major.minor.0
2. Feature gates are preserved when k8s version < NEW_VERSION's major.minor.0

## Skip Conditions

For each test file listed below, skip modification if tests for NEW_VERSION already exist in that file.

## Files to Modify

1. `pkg/api/defaults-apiserver_test.go`
2. `pkg/api/defaults-cloud-controller-manager_test.go`
3. `pkg/api/defaults-controller-manager_test.go`
4. `pkg/api/defaults-kubelet_test.go`
5. `pkg/api/defaults-scheduler_test.go`

## Test Pattern

Each test file should include:

```go
// test user-overrides, removal of feature gates for k8s versions >= NEW_VERSION
cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
cs.Properties.OrchestratorProfile.OrchestratorVersion = "NEW_VERSION.0"  // e.g., "1.31.0"
cs.Properties.OrchestratorProfile.KubernetesConfig.<ComponentConfig> = make(map[string]string)
component := cs.Properties.OrchestratorProfile.KubernetesConfig.<ComponentConfig>
component["--feature-gates"] = "Gate1=true,Gate2=true"
featuregateNewSanitized := ""  // Empty string because gates should be removed for NEW_VERSION+
cs.set<Component>Config()
if component["--feature-gates"] != featuregateNewSanitized {
    t.Fatalf("got unexpected '--feature-gates' for %s \n <component> config original value %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
        "NEW_VERSION.0", featuregateNew, component["--feature-gates"], featuregateNewSanitized)
}

// test user-overrides, no removal of feature gates for k8s versions < NEW_VERSION
cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
cs.Properties.OrchestratorProfile.OrchestratorVersion = "PREVIOUS_VERSION.0"  // e.g., "1.30.0"
cs.Properties.OrchestratorProfile.KubernetesConfig.<ComponentConfig> = make(map[string]string)
component = cs.Properties.OrchestratorProfile.KubernetesConfig.<ComponentConfig>
component["--feature-gates"] = "Gate1=true,Gate2=true"
featuregateOldSanitized := "Gate1=true,Gate2=true"  // Same as input because gates should be preserved
cs.set<Component>Config()
if component["--feature-gates"] != featuregateOldSanitized {
    t.Fatalf("got unexpected '--feature-gates' for %s \n <component> config original value %s \n, actual sanitized value: %s \n, expected sanitized value: %s \n ",
        "PREVIOUS_VERSION.0", featuregateNew, component["--feature-gates"], featuregateOldSanitized)
}
```

## Component-Specific Adjustments

Replace placeholders in the test pattern for each component:

1. **APIServer:**

   - `<ComponentConfig>` = `APIServerConfig`
   - `<Component>` = `APIServer`
   - `set<Component>Config()` = `setAPIServerConfig()`

2. **CloudControllerManager:**

   - `<ComponentConfig>` = `CloudControllerManagerConfig`
   - `<Component>` = `CloudControllerManager`
   - `set<Component>Config()` = `setCloudControllerManagerConfig()`

3. **ControllerManager:**

   - `<ComponentConfig>` = `ControllerManagerConfig`
   - `<Component>` = `ControllerManager`
   - `set<Component>Config()` = `setControllerManagerConfig()`

4. **Kubelet:**

   - `<ComponentConfig>` = `KubeletConfig`
   - `<Component>` = `Kubelet`
   - `set<Component>Config()` = `setKubeletConfig()`

5. **Scheduler:**
   - `<ComponentConfig>` = `SchedulerConfig`
   - `<Component>` = `Scheduler`
   - `set<Component>Config()` = `setSchedulerConfig()`

## Validation Checks

1. Test cases cover both NEW_VERSION and PREVIOUS_VERSION scenarios
2. Feature gates string matches expected format
3. Component-specific config and setup is correct
4. Error messages are descriptive and include relevant details
5. Tests use correct version strings (major.minor.0 format)
