"""
Test cases for the formatter module.
"""

import json
import logging
import unittest
from datetime import datetime, timezone

# Use relative import for testing
from logcore.formatter import StructuredLogFormatter


class TestStructuredLogFormatter(unittest.TestCase):
    """Test cases for the StructuredLogFormatter class."""
    
    def setUp(self):
        """Set up the test fixture."""
        self.formatter = StructuredLogFormatter()
        
    def test_format_basic_message(self):
        """Test formatting a basic log record with just a message."""
        # Create a log record
        record = logging.LogRecord(
            name="test_logger",
            level=logging.INFO,
            pathname="/path/to/test_file.py",
            lineno=42,
            msg="Test message",
            args=(),
            exc_info=None,
            func="test_function"
        )
        
        # Format the record
        formatted = self.formatter.format(record)
          # Verify the format
        self.assertRegex(formatted, r'\[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z\]')  # Timestamp
        self.assertIn('[INFO]', formatted)
        # We need to check only the filename since the full path might be different
        self.assertIn('test_file.py:42', formatted)
        self.assertIn('[test_function]', formatted)
        
        # Extract and verify JSON data
        json_part = formatted.split(' ', 4)[4]
        json_data = json.loads(json_part)
        self.assertEqual(json_data["message"], "Test message")
    
    def test_format_with_dict_message(self):
        """Test formatting a log record with a dict message."""
        # Create a log record with a dict message
        record = logging.LogRecord(
            name="test_logger",
            level=logging.WARNING,
            pathname="/path/to/other_file.py",
            lineno=99,
            msg={"message": "Complex message", "userId": "user123", "action": "login"},
            args=(),
            exc_info=None,
            func="other_function"
        )
        
        # Format the record
        formatted = self.formatter.format(record)
          # Verify the format
        self.assertIn('[WARNING]', formatted)
        # Check only the filename since the full path might be different
        self.assertIn('other_file.py:99', formatted)
        self.assertIn('[other_function]', formatted)
        
        # Extract and verify JSON data
        json_part = formatted.split(' ', 4)[4]
        json_data = json.loads(json_part)
        self.assertEqual(json_data["message"], "Complex message")
        self.assertEqual(json_data["userId"], "user123")
        self.assertEqual(json_data["action"], "login")


if __name__ == '__main__':
    unittest.main()
