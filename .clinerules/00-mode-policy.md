# General Mode Policy: Plan Before Acting

## Objective

This rule defines the standard operating procedure for PLAN MODE and ACT MODE to ensure that all system-altering actions (like file modifications or command executions that change state) are explicitly planned and approved before execution.

## Core Principles

1.  **PLAN MODE is for Planning:**

    - In PLAN MODE, the AI's primary role is to understand the user's request, gather necessary information (e.g., using `read_file`, `search_files`, `list_files`), ask clarifying questions (`ask_followup_question`), and develop a detailed plan of action.
    - If the plan involves modifying files or executing commands that alter the system state, the AI must clearly describe these intended actions in PLAN MODE. This includes specifying which files/commands are involved and the exact changes or operations to be performed (e.g., showing diffs for file changes).
    - The AI **must not** use tools that modify the file system (e.g., `write_to_file`, `replace_in_file`) or execute state-altering commands while in PLAN MODE.
    - Before any such actions are taken, the AI must explicitly state the plan and instruct the user: "To apply these planned changes/actions, please switch to ACT MODE."

2.  **ACT MODE is for Execution:**
    - The AI will only use tools to modify files or execute state-altering commands when in ACT MODE.
    - Actions taken in ACT MODE must directly correspond to the plan that was presented and implicitly or explicitly approved by the user in the preceding PLAN MODE discussion.
    - After each significant action or step is completed in ACT MODE, the AI should confirm its successful completion.

## Interaction with Specific Task Rules (e.g., AddNewVersion)

- This general mode policy takes precedence. Even if a specific task rule (like in `02-add-new-kubernetes-version.md`) outlines steps that involve modifications, those modifications must still adhere to this Plan-then-Act cycle.
- For example, for a multi-step task defined in another rule file:
  1.  In PLAN MODE, the AI will describe the plan for Step 1 (including any file changes).
  2.  The AI will ask the user to switch to ACT MODE to apply Step 1.
  3.  In ACT MODE, the AI will apply Step 1.
  4.  After successful application, the AI will ask the user to switch back to PLAN MODE and issue a "continue" (or similar) command to plan Step 2.
  5.  This cycle repeats for all subsequent steps.

## User Control

- The user retains control over when to switch between PLAN MODE and ACT MODE. The AI will prompt but not force mode changes.
