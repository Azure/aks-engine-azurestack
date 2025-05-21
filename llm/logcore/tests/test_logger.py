"""
Test cases for the logger module.
"""

import logging
import unittest
from unittest import mock
import socket
import threading

# Use relative imports for testing
from logcore.logger import Logger
from logcore.log_level import LogLevel
from logcore.destinations import LogDestination


class MockDestination(LogDestination):
    """Mock implementation of LogDestination for testing."""
    
    def __init__(self):
        self.handler = mock.MagicMock(spec=logging.Handler)
        self.formatter = None
    
    def setup(self, formatter):
        self.formatter = formatter
        return self.handler
    
    def get_handler(self):
        return self.handler


class TestLogger(unittest.TestCase):
    """Test cases for the Logger class."""
    
    def setUp(self):
        """Set up test fixtures."""
        # Reset Logger class state before each test
        Logger._loggers = {}
        Logger._root_logger = None
        Logger._destinations = []
        
        # Create a mock destination
        self.mock_destination = MockDestination()
    
    def tearDown(self):
        """Tear down test fixtures."""
        # Restore the root logger to avoid affecting other tests
        logging.getLogger().handlers = []
    
    def test_configure_root_logger(self):
        """Test that _configure_root_logger sets up logging properly."""
        # Configure the root logger
        Logger._configure_root_logger(LogLevel.INFO, [self.mock_destination])
        
        # Check that the root logger is configured
        self.assertIsNotNone(Logger._root_logger)
        self.assertEqual(Logger._root_logger.level, LogLevel.INFO)
        self.assertEqual(len(Logger._root_logger.handlers), 1)
        self.assertEqual(Logger._root_logger.handlers[0], self.mock_destination.handler)
        
        # Check that the formatter was created and set
        self.assertIsNotNone(self.mock_destination.formatter)
    
    def test_get_logger_returns_same_instance(self):
        """Test that get_logger returns the same logger instance for the same name."""
        # Configure the root logger first
        Logger._configure_root_logger(LogLevel.INFO, [self.mock_destination])
        
        # Get loggers with the same name
        logger1 = Logger.get_logger("test")
        logger2 = Logger.get_logger("test")
        
        # Check they are the same instance
        self.assertIs(logger1, logger2)
        
        # But different from another logger
        logger3 = Logger.get_logger("other")
        self.assertIsNot(logger1, logger3)
    
    def test_get_logger_initializes_correctly(self):
        """Test that get_logger initializes the logger correctly."""
        # Configure the root logger first
        Logger._configure_root_logger(LogLevel.DEBUG, [self.mock_destination])
        
        # Get a logger
        logger = Logger.get_logger("test")
        
        # Check the logger is initialized correctly        self.assertEqual(logger.name, "test")
        self.assertEqual(logger._logger.level, LogLevel.DEBUG)
    
    def test_get_caller_info(self):
        """Test that _get_caller_info returns the correct caller information."""
        # Configure the root logger first
        Logger._configure_root_logger(LogLevel.INFO, [self.mock_destination])
        
        # Get a logger
        logger = Logger.get_logger("test")
        
        # Mock inspect.currentframe to return a predictable frame
        with mock.patch('inspect.currentframe') as mock_currentframe:
            # Create a mock frame with our desired values
            mock_frame = mock.MagicMock()
            mock_outer_frame = mock.MagicMock()
            mock_outer_frame.filename = "/path/to/test_logger.py"
            mock_outer_frame.lineno = 123
            mock_outer_frame.function = "test_get_caller_info"
            
            # Set up the frame chain
            mock_currentframe.return_value = mock_frame
            mock_frame.__enter__.return_value = mock_frame
            mock_frame.__exit__.return_value = None
            
            # Make getouterframes return our mock frame
            with mock.patch('inspect.getouterframes', return_value=[None, None, mock_outer_frame]):
                # Call _get_caller_info directly
                filename, lineno, func_name = logger._get_caller_info()
                
                # The caller should match our mock values
                self.assertEqual(filename, "test_logger.py")
                self.assertEqual(lineno, 123)
                self.assertEqual(func_name, "test_get_caller_info")
    
    def test_prepare_log_data(self):
        """Test that _prepare_log_data adds the correct contextual information."""
        # Configure the root logger first
        Logger._configure_root_logger(LogLevel.INFO, [self.mock_destination])
        
        # Get a logger
        logger = Logger.get_logger("test")
        
        # Prepare log data with some custom fields
        log_data = logger._prepare_log_data(
            "Test message",
            user_id="user123",
            action="login"
        )
        
        # Check the log data contains all expected fields
        self.assertEqual(log_data["message"], "Test message")
        self.assertEqual(log_data["user_id"], "user123")
        self.assertEqual(log_data["action"], "login")
        
        # Check the automatic context fields
        self.assertIn("sourceFile", log_data)
        self.assertIn("lineNumber", log_data)
        self.assertIn("functionName", log_data)
        self.assertIn("processId", log_data)
        self.assertIn("threadId", log_data)
        self.assertEqual(log_data["hostName"], socket.gethostname())
        self.assertEqual(log_data["threadId"], threading.get_ident())
    
    def test_logging_methods(self):
        """Test that the logging methods call the underlying logger correctly."""
        # Configure the root logger first
        Logger._configure_root_logger(LogLevel.DEBUG, [self.mock_destination])
        
        # Get a logger
        logger = Logger.get_logger("test")
          # Mock the _prepare_log_data method to return a fixed dict
        with mock.patch.object(logger, '_prepare_log_data') as mock_prepare:
            mock_prepare.return_value = {"message": "Test message", "key": "value"}
            # Override the logger with a mock
            logger._logger = mock.MagicMock()
            
            # Test each logging method
            logger.debug("Test message", key="value")
            logger._logger.debug.assert_called_with({"message": "Test message", "key": "value"})
            
            logger.info("Test message", key="value")
            logger._logger.info.assert_called_with({"message": "Test message", "key": "value"})
            
            logger.warning("Test message", key="value")
            logger._logger.warning.assert_called_with({"message": "Test message", "key": "value"})
            
            logger.error("Test message", key="value")
            logger._logger.error.assert_called_with({"message": "Test message", "key": "value"})


if __name__ == '__main__':
    unittest.main()
