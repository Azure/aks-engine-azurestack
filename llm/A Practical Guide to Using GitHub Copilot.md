# Leveraging AI Coding Assistants: A Practical Guide to Using GitHub Copilot

## Introduction


GitHub Copilot is an AI-powered coding assistant designed to help developers write code faster, reduce repetitive tasks, and improve productivity.

This guide focuses on using the GitHub Copilot extension in Visual Studio Code (VS Code).
It offers a step-by-step overview of how to get started, including account setup, integration, and best practices for maximizing its benefits in your daily workflow.


For example, a developer can use GitHub Copilot Chat in agent mode to interact directly with an AI assistant.
Here’s how you might use this feature:

- **Open Copilot Chat in your IDE.**
- **Switch to agent mode** to enable advanced AI interactions.
- **Select your preferred model** (e.g., GPT-4.1 or Claude Sonnet 4) for the session.
- **Ask the AI to review, explain, or modify your code**—for example, request a function refactor, bug fix, or code documentation.
- **Review and accept the AI’s suggestions** directly in your development environment, enabling seamless collaboration and rapid iteration.

## GitHub Copilot Account Setup

**Important:** Do not sign up for Copilot on the GitHub website. All Microsoft FTEs and interns already have access via their GitHub Enterprise account.

To get started:
1. Go to [aka.ms/copilot](https://aka.ms/copilot) with your Microsoft credential to set up and check your account status.
2. For automatic sign-in and more info, visit [aka.ms/copilot/explore](https://aka.ms/copilot/explore).

Always follow internal compliance and security guidelines as outlined in the [AI Guidance for Microsoft Developers](https://eng.ms/docs/initiatives/ai-guidance-for-microsoft-developers/github-copilot-guidelines).


## GitHub Copilot Context Management

Managing context in Copilot Chat helps you get more accurate and relevant AI responses. VS Code provides several ways to control what information Copilot uses:

### Implicit Context
- Copilot automatically includes your current file, selected text, and file name in the chat context.
- The active file is shown as a context item in the chat input. You can deselect it if you want to exclude it.

### Explicit Context with #-Mentions
- Use `#` in your prompt to reference specific files, folders, or symbols. For example:
  - `Explain the logic in #main.go`
  - `Refactor #cmd/deploy.go to improve readability`
- Type `#` in the chat input to see a list of available context items.


### Adding Files and Codebase Search
- Add files or folders as context by typing `#` followed by their name, or by dragging and dropping them into the Chat view.
- Use `#codebase` to let Copilot search your entire workspace for relevant information. Example:
  - `How is authentication handled? #codebase`

### Viewing File Changes with #changes
- Use `#changes` to reference and review all modified files in your workspace. Example:
  - `Show me the changes in #changes`

### Reference Web Content
- Use `#fetch` to include content from a web page. Example:
  - `Summarize the latest VS Code features #fetch https://code.visualstudio.com/updates`

### Reference Tools and Extensions
- Reference built-in or custom tools by typing `#` followed by the tool name. Example:
  - `Show open issues #github-mcp`

### Use @-Mentions for Specialized Assistants
- Type `@` to invoke domain-specific chat participants, like `@vscode` or `@terminal`. Example:
  - `@terminal list all Python files in the current directory`


### Managing Chat Sessions and History

- Start a new chat session to reset context when switching topics.
- Delete sessions that are not working for you to keep your workspace organized.

By managing context with these features, you can guide Copilot to generate code and answers that fit your project’s needs.

## GitHub Copilot Prompt Engineering 


### Be Explicit with Your Instructions

To get the best results from Copilot, always state exactly what you want, rather than what you do not want. Clear, informative prompts help Copilot generate code that matches your needs and expectations.

- **Say what you want:** Instead of telling Copilot what to avoid, describe the desired outcome directly. For example, use "Write a Python function that sorts a list of integers using the built-in `sorted()` function" rather than "Don’t use a custom sort implementation."
- **Be informative:** Include all relevant details, such as parameter lists, expected return types, libraries or frameworks to use, and any preferences for style or structure. The more specific your prompt, the more accurate Copilot’s output will be.
- **Iterate and refine:** If Copilot’s initial suggestion isn’t perfect, rephrase your prompt, add more details, or clarify your requirements. Iterative refinement helps you converge on the best solution.
- **Include necessary files in the context:** Reference or attach relevant files, code snippets, or documentation in your prompt. This gives Copilot the background it needs to generate code that fits your project.

**Example:**

> "Write a Python function `filter_active_users(users: list[dict]) -> list[dict]` that returns only users with `'active': True'`. Use the standard library only."

*How it works:* This prompt specifies the function name, parameter and return types, filtering logic, and library preference, making it easy for Copilot to generate exactly what you want.

### Break Complicated Tasks into Smaller Steps

Breaking down complicated tasks into smaller, manageable steps is a key strategy for success when working with Copilot or any AI coding assistant.

#### Benefits of Breaking Down Tasks

- **Clarity and accuracy:** Smaller, focused tasks are easier to describe, understand, and review—helping Copilot generate more correct and relevant code with less risk of errors or misunderstandings.
- **Iterative improvement:** Tackling one step at a time allows you to review, test, and refine each part before moving on.
- **Easier debugging and maintenance:** Isolating issues is simpler when each step is independent and well-defined, and modular code is easier to update or extend.
- **Best practices and collaboration:** Encourages modular design, single responsibility, and makes it easier for teams to collaborate and peer review, as each step can be discussed and validated separately.

#### Example

Suppose you want Copilot to help you build a web API that reads data from a database, processes it, and returns a JSON response. Instead of asking for the entire API at once, break it down into focused prompts:

- **Step 1:** Write a function to connect to a PostgreSQL database and fetch all rows from the `users` table.
- **Step 2:** Write a function to process a list of user records and filter out inactive users.
- **Step 3:** Write a function to serialize a list of user objects to JSON.
- **Step 4:** Combine the above functions into a Flask API endpoint.

*How it works:* By guiding Copilot through each step, you get more reliable, testable code and can easily adjust or debug individual parts. This approach also helps you learn and control the process, rather than relying on a single, complex AI-generated solution.

### Formatting Long or Structured Instructions

For complex tasks or detailed requirements:

- **Use Markdown (`.md`) files for instructions:**
  - Markdown is the standard and recommended format for long or comprehensive instructions.
  - It is easy to read, edit, and version control.

- **Use XML or JSON formats for parameters:**
  - XML and JSON are self-describing and easy to parse.
  - They clearly separate values and reduce ambiguity in parameter passing.

**Example:**

Suppose you want to pass a Kubernetes version as a parameter for a validation or mapping task. You can use XML like this:

```xml
<KubernetesVersion>1.30.10</KubernetesVersion>
```


### Few-Shot Prompting

Few-shot prompting is a technique where you provide Copilot with a few clear examples of both the input and the desired output.

Large language models excel at pattern recognition. When you show Copilot before-and-after examples, it can infer the transformation you want and apply it to similar cases.

This approach is especially effective for repetitive or structured changes—such as code refactoring, formatting, or applying project-specific conventions—because it helps Copilot generalize your intent across your codebase.

**How to use examples to show the desired changes:**

Suppose you want Copilot to help you update default version constants in your codebase. You can provide a before-and-after example:

**Before:**
```go
const KubernetesDefaultRelease = "1.28"
```

**After:**
```go
const KubernetesDefaultRelease = "1.29"
```

By including both the original and the desired result, you teach Copilot the exact change you want. This approach can be used for any coding task—show a few before/after pairs, and Copilot will generalize the pattern for similar changes throughout your codebase.

### Chain of Thought

Chain of Thought prompting encourages Copilot to reason step by step, especially for complex or multi-stage problems. You can guide Copilot by asking it to explain its reasoning or break down the solution into logical steps.

**Example:**

```python
# Explain step by step how this function works and suggest improvements
def process_data(data):
    # ...function code...
```
*How it works:* This prompt encourages Copilot to reason through the function’s logic step by step, providing a detailed explanation and actionable feedback. This is especially useful for complex or unfamiliar code, as it helps you understand and improve it.


### Using Formatted Data and Patterns

You can improve Copilot’s understanding and output by using structured formats and clear naming in your prompts and code.

- **Use formatted data:** When possible, provide data in structured formats such as XML or JSON. This helps Copilot recognize data schemas and relationships.

  **Example:**
  ```json
  {
    "user": {
      "id": 123,
      "name": "Alice"
    }
  }
  ```
  *How it works:* Providing structured data in your prompt helps Copilot recognize the schema and relationships, leading to more accurate code generation for tasks like parsing, validation, or transformation.

- **Create variables for future reference:** Assign meaningful variable names to important data or objects you will use later. This makes your code and prompts clearer for both you and Copilot.

  **Example:**
  ```python
  user_profile = get_user_profile(user_id)
  # Use user_profile in subsequent code
  ```
  *How it works:* By naming variables meaningfully and referencing them in comments, you help Copilot maintain context and continuity in multi-step code generation tasks.

- **Write patterns with meaningful words:** Use descriptive and self-explanatory names for functions, variables, and patterns. This helps Copilot generate code that is easier to read and maintain.

  **Example:**
  ```python
  def calculate_total_price(items):
      # ...function code...
  ```
  *How it works:* Using descriptive function and variable names makes your intent clear, enabling Copilot to generate code that is easier to read, maintain, and extend.

### Custom Instructions: Tailoring Copilot to Your Workflow

Custom instructions in Visual Studio Code let you define coding guidelines, best practices, and project requirements that Copilot will automatically consider when generating code, reviewing, or suggesting changes. This ensures that Copilot’s output aligns with your team’s standards and the specific needs of your project.

#### Types of Custom Instruction Files

- **.github/copilot-instructions.md**: Place this file at the root of your repository to define workspace-wide coding guidelines and requirements. All instructions in this file are automatically included in every Copilot chat request for the workspace. Use this file for general coding practices, preferred technologies, and project-wide standards.

- **.instructions.md files**: Create one or more `.instructions.md` files for more granular, task-specific, or language-specific guidance. These files can be placed in the `.github/instructions` folder or elsewhere in your workspace. Use the `applyTo` property in the file’s front matter to target specific files, folders, or patterns (e.g., only apply to test files or a particular language). These instructions are automatically included for matching files, or can be manually attached to a chat session as needed.

Both types of files use Markdown formatting and support the `applyTo` property for automatic application. You can combine multiple instruction files, but avoid conflicting rules.


#### How Custom Instructions Work

- **Settings Integration:** You can also specify instructions in your VS Code `settings.json` for tasks like code review or commit message generation.
- **Automatic and Manual Use:** Instructions are applied automatically when the instruction file includes an `applyTo` property in its front matter. This property specifies a glob pattern (such as `**/*.py` for all Python files), and Visual Studio Code will automatically apply the instructions to all files matching that pattern. For example, if you set `applyTo: '**'`, the instructions will be included in every Copilot chat request in your workspace. You can also manually attach instructions to a chat session via the Chat view’s Configure Chat menu.

#### Example: Workspace-Wide Coding Guidelines

Create a `.github/copilot-instructions.md` file at the root of your repository:

```markdown
---
applyTo: '**'
---
- Use descriptive variable and function names.
- Always write docstrings for public functions and classes.
- Prefer composition over inheritance unless a clear "is-a" relationship exists.
- Follow PEP 8 for Python code formatting.
```

*How it works:* These instructions will be included in every Copilot chat request in your workspace, guiding the AI to follow your team’s conventions and style automatically.

#### Example: Task-Specific Instructions

Create a `.github/instructions/test-writing.instructions.md` file:

```markdown
---
applyTo: '**/*_test.py'
description: 'Guidelines for writing Python tests.'
---
- Use pytest for all new tests.
- Name test functions with the `test_` prefix.
- Mock external dependencies using the `unittest.mock` library.
```

*How it works:* These instructions are only applied to files matching the pattern `*_test.py`, ensuring that Copilot follows your preferred testing practices when generating or reviewing test code.

#### Example: Code Review Instructions in Settings

Add to your `.vscode/settings.json`:

```json
{
  "github.copilot.chat.reviewSelection.instructions": [
    { "text": "Check for security vulnerabilities and suggest improvements." },
    { "file": ".github/instructions/secure-coding.instructions.md" }
  ]
}
```

*How it works:* When you ask Copilot to review code, it will use these instructions to focus on security and reference additional guidelines from the specified file.

#### Best Practices for Custom Instructions

- Keep instructions clear, concise, and self-contained.
- Organize instructions by topic or task for easier management.
- Use version control to share and track changes with your team.
- Avoid conflicting or ambiguous instructions across files.

For more details and advanced usage, see the [official documentation](https://code.visualstudio.com/docs/copilot/copilot-customization#_custom-instructions).

## Using AI to Write Prompts for You

If you’re unsure how to craft an effective prompt, you can use Copilot or other AI assistants to help generate prompts tailored to your coding goals. This can save time and help you learn best practices for prompt engineering.

### How to Use AI to Write Prompts

- **Describe your coding goal or problem in plain language.**
- **Ask Copilot Chat or another AI assistant to suggest a prompt for your scenario.**
- **Review and refine the suggested prompt to fit your specific needs.**

**Example:**

> "I want to write a Python function that parses a JSON file and returns a list of user names. Can you suggest a good prompt for Copilot?"

**AI-Suggested Prompt:**

> Write a Python function that reads a JSON file containing user data and returns a list of user names.

You can then use or further refine this prompt in Copilot to get high-quality code suggestions.

## Using AI to Review Its Output

You can leverage Copilot or other AI assistants to review and critique the code they generate. This helps catch mistakes, improve code quality, and learn from the AI’s reasoning process.

### How to Use AI for Code Review

- **Ask the AI to review its own output for correctness, security, and style.**
- **Request suggestions for improvements or alternative approaches.**
- **Use follow-up prompts to clarify or expand on the review.**

**Example:**

> "Review the following function for bugs, security issues, and code style. Suggest improvements."

**AI Review Output:**

> The function works as intended, but you should add error handling for invalid input. Consider using type hints and adding a docstring for clarity. Also, avoid using global variables for better maintainability.

By using AI to review its own output, you can quickly iterate on code quality and adopt best practices with minimal effort.


