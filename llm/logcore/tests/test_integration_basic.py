"""
Basic integration tests for the logcore module.
"""

import unittest
from io import StringIO
import sys
from contextlib import redirect_stdout

# Use relative imports for testing
from logcore import configure_logger, get_logger
from logcore.log_level import LogLevel


class TestBasicLoggingIntegration(unittest.TestCase):
    """Basic integration tests for the logging system."""
    
    def setUp(self):
        """Set up test fixtures."""
        # Reset all loggers to avoid interference between tests
        import logging
        root = logging.getLogger()
        for handler in root.handlers[:]:
            root.removeHandler(handler)
            
        # Capture stdout
        self.stdout_capture = StringIO()
    
    def test_logger_initialization(self):
        """Test that the logger can be initialized."""
        # Configure logger
        configure_logger(level=LogLevel.DEBUG)
        
        # Get a logger
        logger = get_logger("test_init")
        
        # Verify logger exists
        self.assertEqual(logger.name, "test_init")
    
    def test_basic_logging(self):
        """Test basic logging functionality."""
        # Configure logger
        configure_logger(level=LogLevel.DEBUG)
        
        # Get a logger
        logger = get_logger("test_basic")
        
        # Log a simple message with stdout redirection
        with redirect_stdout(self.stdout_capture):
            logger.info("This is a test message")
        
        # Get the output
        output = self.stdout_capture.getvalue()
        
        # Verify the message was logged
        self.assertIn("This is a test message", output)
        # Verify basic structure
        self.assertIn("[INFO]", output)


if __name__ == '__main__':
    unittest.main()
