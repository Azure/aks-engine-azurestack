# Adding New Kubernetes Version - Overview

## System Role and Rules

### General Principles

- All modifications must be thoroughly documented with clear, concise explanations
- All code provided must be complete, executable, and free from omissions
- Changes must be validated against the project's standards and best practices
- Each step must be completed and confirmed before proceeding to the next

### User Interaction and Control

- If additional information is needed, pause and request necessary details
- Validate all user input before proceeding
- Wait for explicit confirmation before moving to next steps

### Best Practices

- Ensure clarity and accuracy in all responses
- Maintain consistent and logical structure
- Follow project coding standards
- Document all changes thoroughly

## Components

### 1. Version Maps (`components/01-versions.md`)

Updates version support maps to add new Kubernetes version.

- Command: `AddNewVersion maps <version>`
- Example: `AddNewVersion maps 1.31.8`

### 2. Default Constants (`components/02-defaults.md`)

Updates default version strings across the codebase.

- Command: `AddNewVersion defaults <version>`
- Example: `AddNewVersion defaults 1.31.8`

### 3. Feature Gates (`components/03-feature-gates.md`)

Removes deprecated feature gates for new version.

- Command: `AddNewVersion feature-gates <version> <json array of feature gate list>`
- Example: `AddNewVersion feature-gates 1.31.8 ["Featuregate01", "Featuregate02"]`

### 4. Feature Gates Unit Tests (`components/04-feature-gates-tests.md`)

Adds tests to verify deprecated feature gates for new version..

- Command: `AddNewVersion feature-gates-tests <version> <json array of feature gate list>`
- Example: `AddNewVersion feature-gates-tests 1.31.8 ["Featuregate01", "Featuregate02"]`
