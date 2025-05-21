"""
Integration tests for the logcore module.
"""

import json
import os
import re
import tempfile
import unittest
from io import StringIO
from pathlib import Path
import sys
from contextlib import redirect_stdout

# Use relative imports for testing
from logcore import configure_logger, get_logger
from logcore.log_level import LogLevel


class TestLoggingIntegration(unittest.TestCase):
    """Integration tests for the entire logging system."""
    
    def setUp(self):
        """Set up test fixtures."""
        # Create a temporary directory for log files
        self.temp_dir = tempfile.TemporaryDirectory()
        self.log_dir = Path(self.temp_dir.name)
        
        # Capture stdout
        self.stdout_capture = StringIO()
    
    def tearDown(self):
        """Tear down test fixtures."""
        # Clean up the temporary directory
        self.temp_dir.cleanup()
    
    def test_stdout_logging(self):
        """Test logging to stdout."""
        # Configure logger to only log to stdout
        configure_logger(
            level=LogLevel.DEBUG,
            log_to_stdout=True,
            log_to_file=False
        )
        
        # Get a logger
        logger = get_logger("test_integration")
        
        # Log messages at different levels
        with redirect_stdout(self.stdout_capture):
            logger.debug("Debug message", correlation_id="123")
            logger.info("Info message", correlation_id="456")
            logger.warning("Warning message", correlation_id="789")
            logger.error("Error message", correlation_id="abc")
        
        # Get captured output
        output = self.stdout_capture.getvalue()
        
        # Check that all messages were logged
        self.assertIn("Debug message", output)
        self.assertIn("Info message", output)
        self.assertIn("Warning message", output)
        self.assertIn("Error message", output)
        
        # Check log format
        log_lines = output.strip().split("\n")
        self.assertEqual(len(log_lines), 4)
        
        # Parse each log line to verify structure
        for line in log_lines:
            # Format should be: [TIMESTAMP] [LEVEL] [FILENAME:LINENO] [FUNCTION] {"key1": "value1", ...}
            match = re.match(r'\[(.*?)\] \[(.*?)\] \[(.*?)\] \[(.*?)\] (.*)', line)
            self.assertIsNotNone(match)
            
            timestamp, level, location, function, json_str = match.groups()
            
            # Verify timestamp format
            self.assertRegex(timestamp, r'\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z')
            
            # Verify level is one of DEBUG, INFO, WARNING, ERROR
            self.assertIn(level, ["DEBUG", "INFO", "WARNING", "ERROR"])
            
            # Verify location is filename:lineno
            self.assertRegex(location, r'test_integration\.py:\d+')
            
            # Parse JSON data
            json_data = json.loads(json_str)
            
            # Check required fields
            self.assertIn("message", json_data)
            self.assertIn("correlation_id", json_data)
            self.assertIn("sourceFile", json_data)
            self.assertIn("lineNumber", json_data)
            self.assertIn("functionName", json_data)
            self.assertIn("processId", json_data)
            self.assertIn("threadId", json_data)
            self.assertIn("hostName", json_data)
    
    def test_file_logging(self):
        """Test logging to file."""
        # Configure logger to log to file only
        configure_logger(
            level=LogLevel.INFO,
            log_to_stdout=False,
            log_to_file=True,
            log_dir=self.log_dir,
            app_name="test-app"
        )
        
        # Get a logger
        logger = get_logger("test_file")
        
        # Log messages at different levels
        logger.debug("Debug message")  # Should not be logged (below INFO)
        logger.info("Info message")
        logger.warning("Warning message")
        logger.error("Error message")
        
        # Find the log file (should be only one)
        log_files = list(self.log_dir.glob("test-app-*.log"))
        self.assertEqual(len(log_files), 1)
        
        log_file = log_files[0]
        
        # Read the log file
        with open(log_file, 'r') as f:
            log_content = f.read()
        
        # Check that only INFO and above messages were logged
        self.assertNotIn("Debug message", log_content)
        self.assertIn("Info message", log_content)
        self.assertIn("Warning message", log_content)
        self.assertIn("Error message", log_content)
    
    def test_file_rotation(self):
        """Test log file rotation."""
        # Configure logger with a small file size to trigger rotation
        configure_logger(
            level=LogLevel.INFO,
            log_to_stdout=False,
            log_to_file=True,
            log_dir=self.log_dir,
            max_size_bytes=100,  # Very small to force rotation
            max_total_size_bytes=400,  # Limit to 4 files
            app_name="rotate-test"
        )
        
        # Get a logger
        logger = get_logger("test_rotation")
        
        # Log enough messages to trigger rotation (at least 5 files)
        for i in range(100):
            logger.info(f"Message {i} with some extra content to make the log entry larger")
        
        # Find the log files
        log_files = list(self.log_dir.glob("rotate-test-*.log*"))
        
        # Should be limited by max_total_size_bytes/max_size_bytes (at most 5 files: 1 active + 4 backups)
        self.assertLessEqual(len(log_files), 5)
        
        # Check for backup files (.log.1, .log.2, etc.)
        backup_files = [f for f in log_files if ".log." in str(f)]
        self.assertGreater(len(backup_files), 0)
    
    def test_log_permissions(self):
        """Test log file permissions."""
        # Configure logger to log to file
        configure_logger(
            level=LogLevel.INFO,
            log_to_stdout=False,
            log_to_file=True,
            log_dir=self.log_dir,
            app_name="perm-test"
        )
        
        # Get a logger
        logger = get_logger("test_permissions")
        
        # Log a message
        logger.info("Test message")
        
        # Find the log file
        log_files = list(self.log_dir.glob("perm-test-*.log"))
        self.assertEqual(len(log_files), 1)
        
        log_file = log_files[0]
        
        # Check file permissions (should be 0o644 = -rw-r--r--)
        file_mode = os.stat(log_file).st_mode & 0o777
        self.assertEqual(file_mode, 0o644)


if __name__ == '__main__':
    unittest.main()
