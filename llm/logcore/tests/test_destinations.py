"""
Test cases for the destinations module.
"""

import logging
import os
import tempfile
import unittest
from datetime import timezone
from pathlib import Path
from unittest import mock

# Use relative import for testing
from logcore.destinations import LogDestination, StdoutDestination, FileDestination


class TestLogDestination(unittest.TestCase):
    """Test cases for the abstract LogDestination class."""
    
    def test_cannot_instantiate_abstract_class(self):
        """Test that LogDestination cannot be instantiated directly."""
        with self.assertRaises(TypeError):
            LogDestination()


class TestStdoutDestination(unittest.TestCase):
    """Test cases for the StdoutDestination class."""
    
    def test_setup_creates_handler(self):
        """Test that setup creates a StreamHandler."""
        destination = StdoutDestination()
        mock_formatter = mock.MagicMock(spec=logging.Formatter)
        
        destination.setup(mock_formatter)
        
        # Assert handler was created and formatter was set
        self.assertIsInstance(destination._handler, logging.StreamHandler)
        mock_formatter.assert_not_called()  # Formatter is only set on the handler
    
    def test_get_handler_without_setup_raises_error(self):
        """Test that get_handler raises an error if setup was not called."""
        destination = StdoutDestination()
        
        with self.assertRaises(RuntimeError):
            destination.get_handler()
    
    def test_get_handler_returns_correct_handler(self):
        """Test that get_handler returns the correct handler after setup."""
        destination = StdoutDestination()
        mock_formatter = mock.MagicMock(spec=logging.Formatter)
        
        destination.setup(mock_formatter)
        handler = destination.get_handler()
        
        self.assertIsInstance(handler, logging.StreamHandler)
        self.assertEqual(handler, destination._handler)


class TestFileDestination(unittest.TestCase):
    """Test cases for the FileDestination class."""
    
    def setUp(self):
        """Set up test fixtures."""
        # Create a temporary directory
        self.temp_dir = tempfile.TemporaryDirectory()
        self.log_dir = Path(self.temp_dir.name)
    
    def tearDown(self):
        """Tear down test fixtures."""
        # Clean up the temporary directory
        self.temp_dir.cleanup()
    
    def test_initialization_creates_log_dir(self):
        """Test that initialization creates the log directory."""
        # Create a nested directory path
        nested_dir = self.log_dir / "nested" / "logs"
        
        # Initialize destination with non-existent path
        destination = FileDestination(
            log_dir=nested_dir,
            max_size_bytes=1000,
            max_total_size_bytes=5000,
            app_name="test-app"
        )
        
        # Directory should be created
        self.assertTrue(nested_dir.exists())
        self.assertTrue(nested_dir.is_dir())
    
    def test_setup_creates_rotating_file_handler(self):
        """Test that setup creates a RotatingFileHandler."""
        destination = FileDestination(
            log_dir=self.log_dir,
            max_size_bytes=1000,
            max_total_size_bytes=5000,
            app_name="test-app"
        )
        
        mock_formatter = mock.MagicMock(spec=logging.Formatter)
        
        with mock.patch('datetime.datetime') as mock_datetime:
            # Mock the datetime to control the log filename
            mock_now = mock.MagicMock()
            mock_now.strftime.return_value = "20250521-123456"
            mock_datetime.now.return_value = mock_now
            mock_datetime.timezone = timezone
            
            destination.setup(mock_formatter)
        
        # Assert handler was created
        self.assertIsInstance(destination._handler, logging.handlers.RotatingFileHandler)
          # Check the log file exists
        expected_log_file = self.log_dir / "test-app-20250521-123456.log"
        self.assertTrue(expected_log_file.exists())
        
        # Check file permissions - platform dependent checks
        if os.name == 'posix':  # Unix/Linux/Mac
            # On Unix systems, check for exact permissions (0o644 = -rw-r--r--)
            file_mode = os.stat(expected_log_file).st_mode & 0o777
            self.assertEqual(file_mode, 0o644)
        else:  # Windows
            # On Windows, just verify file exists and is readable
            # The exact permission bits don't translate directly
            self.assertTrue(os.access(expected_log_file, os.R_OK))
        
        # Check max size and backup count
        self.assertEqual(destination._handler.maxBytes, 1000)
        self.assertEqual(destination._handler.backupCount, 5)  # 5000/1000 = 5
    
    def test_get_handler_without_setup_raises_error(self):
        """Test that get_handler raises an error if setup was not called."""
        destination = FileDestination(log_dir=self.log_dir)
        
        with self.assertRaises(RuntimeError):
            destination.get_handler()
    
    def test_get_handler_returns_correct_handler(self):
        """Test that get_handler returns the correct handler after setup."""
        destination = FileDestination(log_dir=self.log_dir)
        mock_formatter = mock.MagicMock(spec=logging.Formatter)
        
        destination.setup(mock_formatter)
        handler = destination.get_handler()
        
        self.assertIsInstance(handler, logging.handlers.RotatingFileHandler)
        self.assertEqual(handler, destination._handler)


if __name__ == '__main__':
    unittest.main()
