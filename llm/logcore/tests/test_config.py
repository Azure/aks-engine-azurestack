"""
Test cases for the config module.
"""

import tempfile
import unittest
from pathlib import Path
from unittest import mock

# Use relative imports for testing
from logcore.config import configure_logger, get_logger
from logcore.log_level import LogLevel
from logcore.logger import Logger
from logcore.destinations import StdoutDestination, FileDestination


class TestConfig(unittest.TestCase):
    """Test cases for the config module functions."""
    
    def setUp(self):
        """Set up test fixtures."""
        # Reset Logger class state before each test
        Logger._loggers = {}
        Logger._root_logger = None
        Logger._destinations = []
        
        # Create a temporary directory
        self.temp_dir = tempfile.TemporaryDirectory()
        self.log_dir = Path(self.temp_dir.name)
    
    def tearDown(self):
        """Tear down test fixtures."""
        # Clean up the temporary directory
        self.temp_dir.cleanup()
    
    @mock.patch('logcore.logger.Logger._configure_root_logger')
    def test_configure_logger_with_defaults(self, mock_configure):
        """Test configure_logger with default arguments."""
        configure_logger()
        
        # Check that the root logger was configured with expected values
        mock_configure.assert_called_once()
        
        # Get the arguments passed to _configure_root_logger
        args, kwargs = mock_configure.call_args
        level, destinations = args
        
        # Check level
        self.assertEqual(level, LogLevel.INFO)
        
        # Check destinations
        self.assertEqual(len(destinations), 1)
        self.assertIsInstance(destinations[0], StdoutDestination)
    
    @mock.patch('logcore.logger.Logger._configure_root_logger')
    def test_configure_logger_without_stdout(self, mock_configure):
        """Test configure_logger without stdout."""
        configure_logger(log_to_stdout=False)
        
        # Check that the root logger was configured with expected values
        mock_configure.assert_called_once()
        
        # Get the arguments passed to _configure_root_logger
        args, kwargs = mock_configure.call_args
        level, destinations = args
        
        # Check destinations - should be empty
        self.assertEqual(len(destinations), 0)
    
    @mock.patch('logcore.logger.Logger._configure_root_logger')
    def test_configure_logger_with_file(self, mock_configure):
        """Test configure_logger with file output."""
        configure_logger(
            log_to_stdout=True,
            log_to_file=True,
            log_dir=self.log_dir,
            max_size_bytes=1000,
            max_total_size_bytes=5000,
            app_name="test-app"
        )
        
        # Check that the root logger was configured with expected values
        mock_configure.assert_called_once()
        
        # Get the arguments passed to _configure_root_logger
        args, kwargs = mock_configure.call_args
        level, destinations = args
        
        # Check destinations
        self.assertEqual(len(destinations), 2)
        self.assertIsInstance(destinations[0], StdoutDestination)
        self.assertIsInstance(destinations[1], FileDestination)
        
        # Check file destination properties
        file_dest = destinations[1]
        self.assertEqual(file_dest._log_dir, self.log_dir)
        self.assertEqual(file_dest._max_size_bytes, 1000)
        self.assertEqual(file_dest._max_total_size_bytes, 5000)
        self.assertEqual(file_dest._app_name, "test-app")
    
    def test_configure_logger_with_file_missing_log_dir(self):
        """Test configure_logger raises ValueError when log_dir is missing."""
        with self.assertRaises(ValueError):
            configure_logger(log_to_file=True)
    
    def test_get_logger_calls_logger_get_logger(self):
        """Test that get_logger calls Logger.get_logger with the correct name."""
        with mock.patch('logcore.logger.Logger.get_logger') as mock_get_logger:
            result = get_logger("test_name")
            
            # Check that Logger.get_logger was called with the correct name
            mock_get_logger.assert_called_once_with("test_name")
            
            # Check that the result is the return value of Logger.get_logger
            self.assertEqual(result, mock_get_logger.return_value)


if __name__ == '__main__':
    unittest.main()
