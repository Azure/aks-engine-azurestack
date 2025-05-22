# Logging Component

A structured JSON logging component for the code update workflow.

## Overview

This logging component is designed to provide a unified interface for recording actions, changes, and errors throughout the code update workflow. It leverages Python's standard logging library and extends it to provide structured JSON logs with automatic contextual information.

## Features

- **Structured JSON Logs**: All log entries are formatted as JSON objects for easy parsing and analysis.
- **Automatic Contextual Information**: Each log entry automatically includes timestamp, source file path, line number, function name, process ID, thread ID, and host name.
- **Configurable Verbosity Levels**: Supports standard logging levels (DEBUG, INFO, WARNING, ERROR, CRITICAL).
- **Multiple Output Destinations**: Supports logging to console, files, or custom handlers.
- **Flexible Context Extension**: Allows adding additional context as key-value pairs to log entries.

## Usage

### Basic Usage

```python
from logcore import setup_logging, get_logger

# Set up logging (once at application startup)
setup_logging(app_name="code-updator", level="INFO")

# Get a logger instance
logger = get_logger()

# Log a message
logger.info("Application started")
```

### Adding Context

```python
# Add context to a single log entry
logger.info("Processing request", {
    "correlationId": "b7e2c1d0-8f3a-4e2a-9c1a-2e4f5b6c7d8e",
    "operationId": "update-20240615-001",
    "httpMethod": "POST",
    "uri": "/api/v1/code/update"
})

# Create a logger with persistent context
request_logger = logger.with_context({
    "correlationId": "b7e2c1d0-8f3a-4e2a-9c1a-2e4f5b6c7d8e",
    "operationId": "update-20240615-001"
})

# Use the contextual logger
request_logger.info("Request received")
request_logger.info("Request processed")
```

### Logging Exceptions

```python
try:
    # Some operation that might fail
    result = process_data(data)
except Exception as e:
    # Log the exception with additional context
    logger.error(
        "Failed to process data",
        exception=e,
        context={
            "dataId": data.get("id"),
            "elapsedTime": 0.342
        }
    )
```

### Custom Output Handlers

```python
from logcore import setup_logging
from logcore.logger import ConsoleHandler, FileHandler

# Set up logging with multiple handlers
setup_logging(
    app_name="code-updator",
    level="INFO",
    handlers=[
        ConsoleHandler(),
        FileHandler("/path/to/logs/app.log")
    ]
)
```

## Output Format

Each log entry is a JSON object with the following structure:

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
