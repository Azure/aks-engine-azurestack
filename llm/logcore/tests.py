"""
Tests for the logging component.

This module contains unit tests for the logging component to verify its functionality,
including log formatting, handlers, contextual information, and configuration options.
The tests ensure that the logging component meets all requirements specified in the
logger.md document, including JSON formatting, automatic contextual information capture,
and proper handling of additional context.
"""

import json
import logging
import os
import tempfile
import unittest
from io import StringIO
from pathlib import Path
from unittest.mock import patch

from . import get_logger, setup_logging
from .logger import ConsoleHandler, FileHandler, Logger, LogFormatter


class TestLogFormatter(unittest.TestCase):
    """Test the LogFormatter class.
    
    These tests verify that the LogFormatter correctly formats log records as JSON objects
    with all the required fields as specified in the requirements:
    - timestamp
    - sourceFilePath
    - lineNumber
    - functionName
    - processId
    - threadId
    - hostName
    - level
    - message
    - additional contextual information as key-value pairs
    """
    
    def test_format(self) -> None:
        """Test that the formatter correctly formats log records as JSON.
        
        This test creates a mock log record with sample data, adds context to it,
        and verifies that the formatter produces a valid JSON string with all
        required fields correctly populated. It also checks that additional
        context is properly included in the JSON output.
        """
        formatter = LogFormatter()
        record = logging.LogRecord(
            name="test",
            level=logging.INFO,
            pathname="/path/to/file.py",
            lineno=42,
            msg="Test message",
            args=(),
            exc_info=None
        )
        
        # Add context to the record
        record.context = {"key": "value"}
        
        # Format the record
        result = formatter.format(record)
        
        # Parse the result back to JSON
        log_entry = json.loads(result)
        
        # Check the structure
        self.assertEqual(log_entry["level"], "INFO")
        self.assertEqual(log_entry["message"], "Test message")
        self.assertEqual(log_entry["sourceFilePath"], "/path/to/file.py")
        self.assertEqual(log_entry["lineNumber"], 42)
        self.assertEqual(log_entry["key"], "value")
        
        # Check that required fields are present
        required_fields = [
            "timestamp", "sourceFilePath", "lineNumber", "functionName",
            "processId", "threadId", "hostName", "level", "message"
        ]
        for field in required_fields:
            self.assertIn(field, log_entry)


class TestLogHandlers(unittest.TestCase):
    """Test the log handlers.
    
    These tests verify that the log handlers (ConsoleHandler and FileHandler)
    correctly create and configure the appropriate logging.Handler instances.
    The handlers are responsible for directing log output to the desired
    destinations (console or file) with the correct formatting.
    """
    
    def test_console_handler(self) -> None:
        """Test that ConsoleHandler creates a StreamHandler.
        
        This test verifies that the ConsoleHandler.setup() method returns a properly
        configured logging.StreamHandler instance, which directs log output to sys.stdout.
        It also confirms that the handler uses the provided formatter.
        """
        formatter = LogFormatter()
        handler = ConsoleHandler().setup(formatter)
        
        self.assertIsInstance(handler, logging.StreamHandler)
        self.assertEqual(handler.formatter, formatter)
    
    def test_file_handler(self) -> None:
        """Test that FileHandler creates a FileHandler.
        
        This test creates a temporary directory and file path, then verifies that
        the FileHandler.setup() method returns a properly configured logging.FileHandler
        instance pointing to the correct file path. It ensures that log entries will be
        written to the specified file with the correct formatter.
        """
        with tempfile.TemporaryDirectory() as temp_dir:
            log_file = Path(temp_dir) / "test.log"
            
            formatter = LogFormatter()
            handler = FileHandler(str(log_file)).setup(formatter)
            
            self.assertIsInstance(handler, logging.FileHandler)
            self.assertEqual(handler.formatter, formatter)
            self.assertEqual(handler.baseFilename, str(log_file))


class TestLogger(unittest.TestCase):
    """Test the Logger class.
    
    These tests verify the core functionality of the Logger class, including:
    - Basic logging at different levels
    - Adding contextual information to log entries
    - Creating loggers with persistent context
    - Logging exceptions with traceback information
    - Setting and respecting log levels
    
    The tests use a StringIO object to capture log output for verification without
    writing to actual files or the console.
    """
    
    def setUp(self) -> None:
        """Set up the test environment.
        
        This method is called before each test. It resets the logging system to
        ensure tests don't affect each other, and sets up a StringIO object to
        capture log output for verification. It also configures a handler with
        the LogFormatter to format logs as JSON.
        """
        # Reset the logging system before each test
        logging.root.handlers = []
        
        # Set up a StringIO to capture log output
        self.log_output = StringIO()
        handler = logging.StreamHandler(self.log_output)
        handler.setFormatter(LogFormatter())
        logging.root.addHandler(handler)
        logging.root.setLevel(logging.DEBUG)
    
    def test_basic_logging(self) -> None:
        """Test basic logging functionality.
        
        This test creates a Logger instance and logs a simple message at the INFO level.
        It then verifies that the log output contains a valid JSON object with the correct
        level and message. This test ensures that the most basic logging functionality
        works as expected with proper JSON formatting.
        """
        logger = Logger("test")
        
        logger.info("Test message")
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entry = json.loads(output)
        
        # Check the structure
        self.assertEqual(log_entry["level"], "INFO")
        self.assertEqual(log_entry["message"], "Test message")
    
    def test_debug_logging(self) -> None:
        """Test debug logging functionality.
        
        This test creates a Logger instance and logs a message at the DEBUG level.
        It verifies that the log output contains a valid JSON object with the
        correct DEBUG level and message. This test ensures that debug logging works
        as expected and produces properly formatted JSON output.
        """
        logger = Logger("test")
        logger.set_level("DEBUG")
        # Log a debug message
        logger.debug("Debug test message")
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entry = json.loads(output)
        
        # Check the structure
        self.assertEqual(log_entry["level"], "DEBUG")
        self.assertEqual(log_entry["message"], "Debug test message")
    
    def test_log_with_context(self) -> None:
        """Test logging with additional context.
        
        This test verifies that context provided as a dictionary to a log method
        (in this case info()) is properly included in the JSON output. This fulfills
        the requirement for including supplementary context as key-value pairs
        alongside the main message, which is essential for operations like tracking
        correlation IDs, HTTP methods, etc.
        """
        logger = Logger("test")
        
        logger.info("Test message", {"key": "value"})
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entry = json.loads(output)
        
        # Check the context
        self.assertEqual(log_entry["key"], "value")
    
    def test_debug_with_context(self) -> None:
        """Test debug logging with additional context.
        
        This test verifies that context provided as a dictionary to the debug() method
        is properly included in the JSON output. It confirms that the debug level
        correctly handles contextual information in the same way as other log levels.
        """
        logger = Logger("test")
        logger.set_level("DEBUG")
        # Log a debug message with context
        logger.debug("Debug test message", {"debug_key": "debug_value", "trace_id": "abc-123"})
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entry = json.loads(output)
        
        # Check the structure and context
        self.assertEqual(log_entry["level"], "DEBUG")
        self.assertEqual(log_entry["message"], "Debug test message")
        self.assertEqual(log_entry["debug_key"], "debug_value")
        self.assertEqual(log_entry["trace_id"], "abc-123")
    
    def test_warning_logging(self) -> None:
        """Test warning logging functionality.
        
        This test creates a Logger instance and logs a message at the WARNING level.
        It verifies that the log output contains a valid JSON object with the
        correct WARNING level and message. This test ensures that warning logging works
        as expected and produces properly formatted JSON output.
        """
        logger = Logger("test")
        
        # Log a warning message
        logger.warning("Warning test message")
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entry = json.loads(output)
        
        # Check the structure
        self.assertEqual(log_entry["level"], "WARNING")
        self.assertEqual(log_entry["message"], "Warning test message")
    
    def test_warning_with_context(self) -> None:
        """Test warning logging with additional context.
        
        This test verifies that context provided as a dictionary to the warning() method
        is properly included in the JSON output. It confirms that the warning level
        correctly handles contextual information in the same way as other log levels.
        """
        logger = Logger("test")
        
        # Log a warning message with context
        logger.warning("Warning test message", {"alert_type": "security", "severity": "medium"})
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entry = json.loads(output)
        
        # Check the structure and context
        self.assertEqual(log_entry["level"], "WARNING")
        self.assertEqual(log_entry["message"], "Warning test message")
        self.assertEqual(log_entry["alert_type"], "security")
        self.assertEqual(log_entry["severity"], "medium")
    
    def test_logger_with_context(self) -> None:
        """Test creating a logger with persistent context.
        
        This test verifies the with_context() method of the Logger class, which creates
        a new logger with persistent context that will be included in all log entries.
        It checks that both the persistent context and additional per-log context are
        properly included in the JSON output. This feature is important for adding
        consistent context (like correlation IDs) to all logs within a particular
        operation or request.
        """
        logger = Logger("test").with_context({"persistent": "value"})
        
        logger.info("Test message", {"key": "value"})
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entry = json.loads(output)
        
        # Check both contexts
        self.assertEqual(log_entry["persistent"], "value")
        self.assertEqual(log_entry["key"], "value")
    
    def test_log_error_with_exception(self) -> None:
        """Test logging an error with an exception.
        
        This test verifies that when an exception is provided to the error() method,
        the exception message is included in the log message and a traceback is added
        to the context. This is crucial for error tracking and debugging, as it provides
        detailed information about what went wrong and where in the code the error occurred.
        The test confirms that the traceback is properly formatted as a list in the JSON output.
        """
        logger = Logger("test")
        
        try:
            raise ValueError("Test error")
        except Exception as e:
            logger.error("Error occurred", exception=e)
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entry = json.loads(output)
        
        # Check the message and traceback
        self.assertEqual(log_entry["level"], "ERROR")
        self.assertIn("Test error", log_entry["message"])
        self.assertIn("traceback", log_entry)
        self.assertIsInstance(log_entry["traceback"], list)
    
    def test_critical_logging(self) -> None:
        """Test critical logging functionality.
        
        This test creates a Logger instance and logs a message at the CRITICAL level.
        It verifies that the log output contains a valid JSON object with the
        correct CRITICAL level and message. This test ensures that critical logging works
        as expected and produces properly formatted JSON output.
        """
        logger = Logger("test")
        
        # Log a critical message
        logger.critical("Critical test message")
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entry = json.loads(output)
        
        # Check the structure
        self.assertEqual(log_entry["level"], "CRITICAL")
        self.assertEqual(log_entry["message"], "Critical test message")
    
    def test_critical_with_context(self) -> None:
        """Test critical logging with additional context.
        
        This test verifies that context provided as a dictionary to the critical() method
        is properly included in the JSON output. It confirms that the critical level
        correctly handles contextual information in the same way as other log levels.
        """
        logger = Logger("test")
        
        # Log a critical message with context
        logger.critical("Critical test message", context={"alert_type": "system_failure", "impact": "high"})
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entry = json.loads(output)
        
        # Check the structure and context
        self.assertEqual(log_entry["level"], "CRITICAL")
        self.assertEqual(log_entry["message"], "Critical test message")
        self.assertEqual(log_entry["alert_type"], "system_failure")
        self.assertEqual(log_entry["impact"], "high")
    
    def test_critical_with_exception(self) -> None:
        """Test logging a critical message with an exception.
        
        This test verifies that when an exception is provided to the critical() method,
        the exception message is included in the log message and a traceback is added
        to the context, similar to the error() method. This ensures that critical errors
        are properly logged with all the necessary debugging information.
        """
        logger = Logger("test")
        
        try:
            raise RuntimeError("Critical system error")
        except Exception as e:
            logger.critical("System failure", exception=e, context={"service": "database"})
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entry = json.loads(output)
        
        # Check the message, traceback, and context
        self.assertEqual(log_entry["level"], "CRITICAL")
        self.assertIn("Critical system error", log_entry["message"])
        self.assertIn("traceback", log_entry)
        self.assertIsInstance(log_entry["traceback"], list)
        self.assertEqual(log_entry["service"], "database")
    
    def test_set_level(self) -> None:
        """Test setting the logging level.
        
        This test verifies that when a Logger is initialized with a specific level (ERROR),
        messages at lower levels (INFO) are not logged. This tests the filtering capability
        of the logging system, which is essential for controlling verbosity in different
        environments (e.g., more verbose in development, less verbose in production).
        The test confirms that only the ERROR message appears in the output.
        """
        logger = Logger("test", level="ERROR")
        
        logger.info("This should not be logged")
        logger.error("This should be logged")
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entries = [json.loads(line) for line in output.strip().split("\n") if line]
        
        # Should only have one log entry
        self.assertEqual(len(log_entries), 1)
        self.assertEqual(log_entries[0]["level"], "ERROR")
    
    def test_set_level_method(self) -> None:
        """Test the set_level method for changing log levels dynamically.
        
        This test verifies that the set_level() method correctly changes the logging level
        of an existing Logger instance. It creates a logger, sets it to a restrictive level
        (ERROR), attempts to log at lower levels, then changes the level to a more permissive
        one (INFO) and verifies that messages at that level are now logged.
        """
        logger = Logger("test", level="ERROR")
        
        # At ERROR level, INFO messages should not be logged
        logger.info("This should not be logged")
        
        # Change the level to INFO
        logger.set_level("INFO")
        
        # Now INFO messages should be logged
        logger.info("This should be logged")
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entries = [json.loads(line) for line in output.strip().split("\n") if line]
        
        # Should have one log entry (the INFO message after setting level to INFO)
        self.assertEqual(len(log_entries), 1)
        self.assertEqual(log_entries[0]["level"], "INFO")
        self.assertEqual(log_entries[0]["message"], "This should be logged")
    
    def test_set_level_hierarchy(self) -> None:
        """Test that the logging level hierarchy is respected.
        
        This test verifies that the logging system correctly respects the level hierarchy:
        DEBUG < INFO < WARNING < ERROR < CRITICAL
        
        It creates a logger with WARNING level and verifies that:
        1. Messages at lower levels (DEBUG, INFO) are not logged
        2. Messages at equal or higher levels (WARNING, ERROR, CRITICAL) are logged
        """
        logger = Logger("test", level="WARNING")
        
        # These should not be logged
        logger.debug("Debug message")
        logger.info("Info message")
        
        # These should be logged
        logger.warning("Warning message")
        logger.error("Error message")
        logger.critical("Critical message")
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Parse the output as JSON
        log_entries = [json.loads(line) for line in output.strip().split("\n") if line]
        
        # Should have three log entries (WARNING, ERROR, CRITICAL)
        self.assertEqual(len(log_entries), 3)
        self.assertEqual(log_entries[0]["level"], "WARNING")
        self.assertEqual(log_entries[1]["level"], "ERROR")
        self.assertEqual(log_entries[2]["level"], "CRITICAL")
    
    def test_caller_function_name(self) -> None:
        """Test that the logger correctly captures the caller's function name.
        
        This test verifies that the Logger correctly identifies and records the name
        of the function that called the logger, not just the name of the logging method.
        It uses two nested function calls to ensure that the function name reported
        is actually the caller's name, not an intermediate function.
        """
        logger = Logger("test")
        
        # Define a nested function structure to test caller detection
        def outer_function() -> None:
            def inner_function() -> None:
                # Log from the inner function
                logger.info("Log from inner function")
            
            # Call the inner function
            inner_function()
            
            # Log directly from outer function
            logger.info("Log from outer function")
        
        # Call the outer function to generate logs
        outer_function()
        
        # Get the log output
        output = self.log_output.getvalue()
        
        # Split the output into individual log entries
        log_entries = [json.loads(line) for line in output.strip().split('\n')]
        
        # First log entry should be from inner_function
        self.assertEqual(log_entries[0]["functionName"], "inner_function")
        
        # Second log entry should be from outer_function
        self.assertEqual(log_entries[1]["functionName"], "outer_function")


class TestSetupLogging(unittest.TestCase):
    """Test the setup_logging function.
    
    These tests verify that the setup_logging function correctly configures the logging
    system according to the provided parameters. The function is responsible for
    initializing the root logger, setting up handlers, and creating the default logger
    instances that will be returned by get_logger().
    """
    
    def test_setup_logging_defaults(self) -> None:
        """Test setup_logging with default parameters.
        
        This test verifies that when setup_logging is called without any parameters,
        it correctly initializes the logging system with default values:
        - The default app name should be "code-updator"
        - The default logging level should be INFO
        - The default log handler should output to the console
        
        It confirms that a Logger instance can be obtained using get_logger() and that
        the default app name is registered in the _loggers dictionary.
        """
        # Reset loggers
        from logcore.logger import _loggers
        _loggers.clear()
        
        # Set up logging
        setup_logging()
        
        # Check that the default logger exists
        logger = get_logger()
        self.assertIsInstance(logger, Logger)
        
        # Check that the default app name is used
        from logcore.logger import _loggers
        self.assertIn("code-updator", _loggers)
    
    def test_setup_logging_custom(self) -> None:
        """Test setup_logging with custom parameters.
        
        This test verifies that setup_logging correctly applies custom configuration:
        - Custom app name ("test-app")
        - Custom logging level (DEBUG)
        - Custom default context ({"app": "test"})
        
        It confirms that a Logger instance can be obtained using get_logger() and that
        the custom app name is registered in the _loggers dictionary. This test ensures
        that the logging system is configurable to meet the specific needs of different
        applications or environments.
        """
        # Reset loggers
        from logcore.logger import _loggers
        _loggers.clear()
        
        # Set up logging with custom parameters
        setup_logging(
            app_name="test-app",
            level="DEBUG",
            default_context={"app": "test"}
        )
        
        # Check that the custom logger exists
        logger = get_logger()
        self.assertIsInstance(logger, Logger)
        
        # Check that the custom app name is used
        from logcore.logger import _loggers
        self.assertIn("test-app", _loggers)
    
    def test_file_output(self) -> None:
        """Test logging to a file.
        
        This test verifies that setup_logging can configure the logging system to write
        logs to a file instead of the console. It:
        1. Creates a temporary directory and file path
        2. Sets up logging with a FileHandler pointing to that path
        3. Logs a test message
        4. Verifies that the file exists and contains a valid JSON log entry
        5. Confirms that the log entry has the correct level and message
        
        This test is crucial for ensuring that logs can be persisted to disk for later
        analysis, which is often required in production environments.
        """
        with tempfile.TemporaryDirectory() as temp_dir:
            log_file = Path(temp_dir) / "test.log"
            
            # Reset loggers
            from logcore.logger import _loggers
            _loggers.clear()
            
            # Set up logging with a file handler
            setup_logging(
                app_name="test-app",
                level="INFO",
                handlers=[FileHandler(str(log_file))]
            )
            
            # Get a logger and log a message
            logger = get_logger()
            logger.info("Test message")
            
            # Check that the file exists and contains the log message
            self.assertTrue(log_file.exists())
            
            with open(log_file, "r") as f:
                content = f.read()
            
            # Parse the content as JSON
            log_entry = json.loads(content)
            
            # Check the structure
            self.assertEqual(log_entry["level"], "INFO")
            self.assertEqual(log_entry["message"], "Test message")


if __name__ == "__main__":
    """
    Entry point for running the tests directly.
    
    When this script is executed directly (rather than imported), it runs all the tests
    defined in this module. This is useful for development and debugging of the logging
    component, allowing quick verification that changes don't break existing functionality.
    """
    unittest.main()
