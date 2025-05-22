# Logging Component Requirements

## Overview

The logging component is responsible for recording all actions, changes, and errors throughout the code update workflow. It must provide configurable verbosity levels, automatically capture contextual information, and structure log messages as JSON for traceability and machine parsing. The component should be extensible, support multiple output destinations, and handle sensitive data appropriately.

---

## Functional Requirements

### Functional Requirements

**Logging System Integration**:  
  Implement an application-level logging wrapper that leverages Python’s standard logging library. The wrapper should provide a unified interface for recording actions, changes, and errors throughout the code update workflow. It must support configurable log levels, structured log messages.

### 1. Automatic Contextual Information

Each log entry must automatically include:

- **Timestamp**: Date and time in ISO 8601 format (UTC).
- **Source File Path**: File path of the Python file where the logger was called.
- **Line Number**: Line number in the source file.
- **Function Name**: Name of the function or method where the logger was called (if available, otherwise null).
- **Process ID**: ID of the process (always included).
- **Thread ID**: ID of the thread (always included).
- **Host Name**: Name of the host machine (always included).

### 2. Additional Context

- The application must support the inclusion of supplementary context as key-value pairs alongside the main message for each operation.
- Example context keys include (but are not limited to):
    - `correlationId`
    - `operationId`
    - `httpMethod`
    - `uri`
    - `body`
    - `elapsedTime`
- The system should allow for flexible extension to accommodate additional context keys as needed.

### 3. Log Entry Format

- Each log entry **must be output as a single JSON object** with the following structure:

```json
{
  "timestamp": "2024-06-15T14:23:45.123Z",
  "sourceFilePath": "/app/code_updater/workflows/update_manager.py",
  "lineNumber": 87,
  "functionName": "apply_update_template",
  "processId": 4521,
  "threadId": 140735218045696,
  "hostName": "dev-machine-01",
  "level": "INFO",
  "message": "Applied update template to function 'process_data' in 'data_utils.py'.",
  "correlationId": "b7e2c1d0-8f3a-4e2a-9c1a-2e4f5b6c7d8e",
  "operationId": "update-20240615-001",
  "httpMethod": "POST",
  "uri": "/api/v1/code/update",
  "body": "{\"function\": \"process_data\", \"template\": \"add_logging\"}",
  "elapsedTime": 0.342,
  "customKey": "customValue"
}
```
### 4. Default values
- Set the default logging level to INFO.
- Set the default logging output to the console.
- Set the default app name to "code-updator"
