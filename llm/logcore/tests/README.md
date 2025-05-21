# Tests for LogCore Module

This directory contains tests for the LogCore module.

## Running Tests

There are several ways to run the tests:

### Using unittest directly

```bash
# From the tests directory
python run_tests.py

# From the logcore directory
python run_tests.py
```

### Using pytest (recommended)

```bash
# Install test dependencies first
pip install -r requirements-test.txt

# From the logcore directory
python run_pytest.py

# With specific options
python run_pytest.py -v --no-cov  # Verbose, no coverage

# Run specific test file
python run_pytest.py tests/test_logger.py

# Run specific test
python run_pytest.py tests/test_logger.py::TestLogger::test_logging_methods
```

## Test Structure

- `test_log_level.py` - Tests for the LogLevel enum
- `test_formatter.py` - Tests for the StructuredLogFormatter
- `test_destinations.py` - Tests for log destinations (stdout, file)
- `test_logger.py` - Tests for the Logger class
- `test_config.py` - Tests for the configuration functions
- `test_integration.py` - Integration tests for the entire logging system

## Code Coverage

When running tests with pytest, code coverage reports are generated. An HTML report is
created in the `coverage_html` directory that shows which lines of code are covered
by the tests.

## Adding New Tests

To add a new test:

1. Create a new file named `test_<module>.py` in this directory
2. Define test classes that inherit from `unittest.TestCase`
3. Define test methods that start with `test_`
4. Run the tests to verify they work correctly
