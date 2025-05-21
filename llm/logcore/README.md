# LogCore Module

A structured logging module for Python applications that provides configurable verbosity levels, automatic contextual information, and multiple output destinations.

## Features

- Multiple log levels (`DEBUG`, `INFO`, `WARNING`, `ERROR`)
- Automatic inclusion of contextual information including:
  - Timestamp in ISO 8601 format (UTC)
  - Source file name
  - Line number
  - Function name
  - Process ID
  - Thread ID
  - Host name
- Structured log messages in JSON format
- Multiple output destinations (stdout and file)
- Log file rotation with configurable size limits
- Thread-safe implementation
- Support for custom context fields
- Configurable log retention policies

## Installation

```bash
# From PyPI (if published)
pip install llm-logcore

# From source
git clone https://github.com/yourusername/llm-logcore.git
cd llm-logcore
pip install .
```

## Usage

### Basic Usage

```python
from llm.logcore import configure_logger, get_logger, LogLevel

# Configure the logger (once in your application)
configure_logger(level=LogLevel.INFO, log_to_stdout=True)

# Get a logger instance
logger = get_logger("my_module")

# Log messages
logger.debug("Detailed debugging information")
logger.info("Normal operation information")
logger.warning("Warning about potential issues")
logger.error("Error condition")
```

### Adding Contextual Data

You can add any key-value pairs to your log messages:

```python
logger.info(
    "Processing data",
    correlationId="abc123",
    operationName="ProcessData",
    dataSize=1024
)
```

### Configuring File Output with Rotation

```python
from llm.logcore import configure_logger, LogLevel

configure_logger(
    level=LogLevel.INFO,
    log_to_stdout=True,
    log_to_file=True,
    log_dir="/var/log/myapp",
    max_size_bytes=5 * 1024 * 1024,  # 5 MB
    max_total_size_bytes=500 * 1024 * 1024,  # 500 MB
    app_name="myapp"
)
```

### Log Format

Logs are formatted as:

```
[TIMESTAMP] [LEVEL] [FILENAME:LINENO] [FUNCTION] {"Key1":"Value1", "Key2":"Value2"}
```

Example:

```
[2024-06-10T14:23:45.123Z] [INFO] [processor.py:42] [update_data] {"correlationId":"abc123","operationName":"UpdateData","message":"Processing completed successfully"}
```

## Thread Safety

The logger is designed to be thread-safe, allowing concurrent logging from multiple threads without data corruption or race conditions.

## Advanced Configuration

### Setting Default Context

You can set default context values that will be included in all log entries:

```python
from llm.logcore import configure_logger, set_default_context, LogLevel

configure_logger(level=LogLevel.INFO)
set_default_context(service="auth-service", environment="production")

# All subsequent log entries will include the service and environment fields
logger = get_logger("auth")
logger.info("User authenticated")  # Will include service and environment fields
```

### Using Context Managers

For operation-scoped context:

```python
from llm.logcore import get_logger, log_context

logger = get_logger("transactions")

with log_context(operation_id="tx-123456", user_id="user-789"):
    # All log messages within this context will include these fields
    logger.info("Starting transaction")
    # Process...
    logger.info("Transaction complete")
```

## Performance Considerations

- Log messages below the configured level are discarded with minimal processing overhead
- Context capture is optimized to minimize impact on application performance
- File rotation operations happen asynchronously to avoid blocking the application

## Contributing

Contributions are welcome! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
