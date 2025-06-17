# Input 
<KubernetesVersion>{{k8s_version}}</KubernetesVersion>

# Input Validation
- Get Kubernetes version in xml tag <KubernetesVersion>
- **Do not add code to implement the Input Validation logic.**
- **Version Existence Check**: Search the `TestExampleAPIModels` test cases for entries with `apiModelPath: "../examples/kubernetes-releases/kubernetes[MAJOR].[Minor].json"`. 
    - If the desired version DO NOT exist in the list, return "True".
    - If the desired version exists in the test cases, return "False". 
  
  
# Code Snippet Filter:
   - source code path: `cmd/generate_test.go`
   - function name: TestExampleAPIModels
   - object type: test function with slice of test cases
   - look for: test cases with `apiModelPath: "../examples/kubernetes-releases/kubernetes*.json"`


## Test Case Update Checklist

1. Test Case Updates (in `cmd/generate_test.go`):
   - [ ] **If the new Kubernetes version already exists in the test cases, make no changes and return the original code exactly as received.**
   - [ ] Add new test case for version [MAJOR].[Minor].[REVISION] with name "[MAJOR].[Minor] example"
   - [ ] Set apiModelPath to "../examples/kubernetes-releases/kubernetes[MAJOR].[Minor].json"
   - [ ] Set setArgs to defaultSet
   - [ ] Keep previous version N-1 ([MAJOR].[Minor-1].x) test case
   - [ ] Remove older version N-2 ([MAJOR].[Minor-2].x) test case
   - [ ] Preserve the original code's spacing and formatting
   - [ ] Place new test case after the previous version test case

**You must review and ensure that all items on the **Test Case Update Checklist** are checked. If any items are not checked, make the necessary changes to ensure all checkboxes are checked.**

## Example: Adding Kubernetes 1.30

Here's an example of required changes when adding Kubernetes 1.30.10:

1. **Test Case Updates**:

```go
// Before
		{
			name:         "1.28 example",
			apiModelPath: "../examples/kubernetes-releases/kubernetes1.28.json",
			setArgs:      defaultSet,
		},
		{
			name:         "1.29 example",
			apiModelPath: "../examples/kubernetes-releases/kubernetes1.29.json",
			setArgs:      defaultSet,
		},

// After
		{
			name:         "1.29 example",
			apiModelPath: "../examples/kubernetes-releases/kubernetes1.29.json",
			setArgs:      defaultSet,
		},
		{
			name:         "1.30 example",
			apiModelPath: "../examples/kubernetes-releases/kubernetes1.30.json",
			setArgs:      defaultSet,
		},
```

**After making changes, you MUST review the **Test Case Update Checklist** to ensure all items are checked. If any items remain unchecked, make the necessary changes until all checkboxes are checked.**
