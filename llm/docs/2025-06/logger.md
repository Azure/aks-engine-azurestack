# Logging Component Requirements

## Overview

The logging component is responsible for recording all actions, changes, and errors throughout the code update workflow. It must provide configurable verbosity levels, automatically capture contextual information, and structure log messages as JSON for traceability and machine parsing. The component should be extensible, support multiple output destinations, and handle sensitive data appropriately.

---

## Functional Requirements

### 1. Verbosity Levels

Use logging library

- Support the following log levels:
  - `DEBUG`: Detailed information for diagnosing problems.
  - `INFO`: Confirmation that things are working as expected.
  - `WARNING`: Indication of unexpected events or potential issues.
  - `ERROR`: Serious problems that prevent some functionality.
- The minimum log level must be configurable at runtime.

### 2. Automatic Contextual Information

Each log entry must automatically include:

- **Timestamp**: Date and time in ISO 8601 format (UTC).
- **Source File Name**: Name of the Python file where the logger was called.
- **Line Number**: Line number in the source file.
- **Function Name**: Name of the function or method where the logger was called (if available, otherwise null).
- **Process ID**: ID of the process (always included).
- **Thread ID**: ID of the thread (always included).
- **Host Name**: Name of the host machine (always included).

If any contextual information is unavailable, the field must be present with a null value.

### Log Message Format

All log messages must be structured as key-value pairs (`string`, `string`). Each log entry should be serialized as a JSON object before being written to the log, ensuring consistency and facilitating downstream parsing and analysis.

Each log entry must be output in the following format:

`[TIMESTAMP] [LEVEL] [FILENAME:LINENO] [FUNCTION] {"Key1":"Value1", "Key2":"Value2"}`

#### Example

`[2024-06-10T14:23:45.123Z] [INFO] [data_processor.py:42] [UpdateFunctionSnippet] {"correlationId":"abc123-xyz789","operationId":"op456","operationName":"UpdateFunctionSnippet","url":"https://abc.com","body":"request body","message":"Function 'process_data' updated successfully."}`

### 3. Output Destinations

The logger must support configurable output destinations, including:

1. **Standard Output (stdout):**

   - Log messages can be directed to standard output for real-time monitoring or containerized environments.

2. **File Output:**
   - Log messages can be written to a file on disk.
   - The logger must implement log file rotation:
     - Create a new log file when the current file reaches 5 MB in size.
     - Maintain a maximum total log storage size of 500 MB by deleting or archiving the oldest log files as needed.
   - File naming should include a timestamp or sequence number to distinguish rotated files.
   - File permissions must be set appropriately to prevent unauthorized access.

The output destination(s) must be configurable at runtime.

#### Log File Naming Pattern

- Log files should be named using a pattern that includes the date and time to ensure uniqueness and facilitate easy log management.
- The recommended pattern is:  
  `appname-YYYYMMDD-HHMMSS.log`

  - `appname`: (Optional) A short identifier for the application or service.
  - `YYYYMMDD`: 4-digit year, 2-digit month, 2-digit day.
  - `HHMMSS`: 2-digit hour (24-hour), 2-digit minute, 2-digit second (UTC or local time, specify in documentation).

- Use UTC time for consistency across distributed systems.
- Avoid spaces and special characters in file names.
- **Example:**  
   `codeupdater-20240611-153045.log`
