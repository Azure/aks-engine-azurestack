# Copilot Agent Mode Instructions

When operating in **agent mode**, you MUST strictly follow these instructions for every request. Ensure clarity, accuracy, and compliance with all technical conventions.

## Rules

1. **Always include the following files in the context:**

   - `.clinerules/01-project-context.md`
   - `.clinerules/02-add-new-kubernetes-version.md`

2. **Based on the command:**

   - Follow the instructions in `.clinerules/02-add-new-kubernetes-version.md`.
   - Add the relevant detailed instruction from the appropriate file under `.clinerules/components`.

3. **If the instructions require file modifications:**
   - Include the files to be modified in the context.

## Additional Requirements

- You MUST follow these instructions **verbatim**.
- **Do not add any other content** beyond what is specified above.

## Plan and Act

1. **Plan:**  
   Before executing any command, carefully analyze and plan the required steps. Ensure you fully understand the requirements.
2. **Show Plan:**  
   Present your plan in detail, including all intended code changes for each file.
3. **Seek Confirmation:**  
   Ask the user for confirmation before proceeding with any changes or execution.

---
