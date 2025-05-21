# Summarize current file

Analyze the file `cmd/generate.go`. Write a flowchart in Mermaid format that illustrates the flow of the generate command, including its inputs, outputs, and major processes. Base your analysis solely on the code within this file; do not reference or infer from any external files or documentation.

---

# Drill down to add dependencies

```code
In the file `pkg/api/common/versions.go`, analyze the function:

`func GetSupportedVersions(orchType string, isUpdate, hasWindows bool, isAzureStackCloud bool) (versions []string, defaultVersion string)`

Drill down into its internal dependencies, including any functions, constants, or variables it directly calls or uses within the same file or package. For each key variable, function, or constant used, include its relative source code path. Draw a flowchart in Mermaid format that represents the function's logic, showing major decision points, dependencies (with their relative paths), inputs, and outputs.
```

---

# Based on Commit

This is the commit `ec6201313361b2a89713c0c04308cc0deb54410f` for adding version 1.30.10. Review the changes in this commit and derive the step-by-step logic required to add a new version, specifically version 1.31.8. Clearly outline the necessary modifications, following the same pattern and logic as in the referenced commit.

---

# Logic Pseudo code

Analyze the file pkg\api\defaults-apiserver.go. For each function in this file, create a section with a level 2 heading using the function's name, and summarize its essential logic in pseudocode.

---

# Steps Format and validateion

Present the required changes in checklist format. At the end, state that you must ensure all checkboxes are checked; if any are not, make the necessary changes until all are checked.

# System prompt - Refine my instructions

Transform the given input into a concise and precise instruction suitable for guiding a language model. Optimize the instruction for clarity and contextual relevance. Present your refined instruction in raw markdown, enclosed within a code block.

# Create a Python application

- Create a Python application that uses a large language model (LLM) to automate code update workflows. Your design should:
- Define the types of code updates to be automated (e.g., refactoring, dependency upgrades, bug fixes).
- Integrate with version control systems such as Git for tracking and applying changes.
- Implement robust error handling to manage LLM failures or invalid code suggestions.
- Specify the desired level of user interaction (e.g., fully automated, user approval required).
- Ensure clear input formatting for LLM prompts and structured output handling.
- Include safeguards to maintain code integrity, such as automated testing or validation before committing changes.

# Prmpt guidance:

- Demonstrate code changes using either a single example (one-shot) or multiple examples (few-shot),
- Validation steps to ensure all necessary changes are correctly implemented.
- Create a multi-step plan, provide the necessary context for each step, and execute each step one by one.

# Code Generation Rules

Save the resulting information to `projectcontext/defaultversion.md`.

- The file should not exceed 1000 lines.
- Log each condition, function entry and exit, and before and after each data processing step.
- Implement logging with the following verbose levels and categories:

  - Debug:

    - Log when entering a function. (Example: `log.Debug("Entering function X")`)
    - Log when entering and exiting minor conditions. (Example: `log.Debug("Minor condition met: Y")`)

  - Information:

    - Log before and after data processing steps. (Example: `log.Info("Starting data processing for user: %s", userID)`)
    - Log when entering and exiting major conditions. (Example: `log.Info("Entering main processing loop")`)

  - Warning:

    - Log when something is wrong, but the process can still proceed. (Example: `log.Warning("Configuration value is deprecated")`)

  - Error:
    - Log when something is wrong and the process cannot proceed. (Example: `log.Error("Failed to open file: %v", err)`)

## Files and Descriptions

| File Location                                  | Description                              |
| ---------------------------------------------- | ---------------------------------------- |
| projectcontext/01-versions.md                  | Kubernetes supported and default version |
| projectcontext/02-0-defaults-apiserver.md      | Defaults for apiserver                   |
| projectcontext/02-1-defaults-apiserver-test.md | Unit test for Defaults of apiserver      |
