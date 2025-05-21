"""
Test cases for the log_level module.
"""

import unittest
import logging

# Use relative import for testing
from logcore.log_level import LogLevel


class TestLogLevel(unittest.TestCase):
    """Test cases for the LogLevel class."""
    
    def test_log_level_values(self):
        """Test that log levels have the correct values."""
        self.assertEqual(LogLevel.DEBUG, logging.DEBUG)
        self.assertEqual(LogLevel.INFO, logging.INFO)
        self.assertEqual(LogLevel.WARNING, logging.WARNING)
        self.assertEqual(LogLevel.ERROR, logging.ERROR)
    
    def test_log_level_comparison(self):
        """Test log level comparison operations."""
        self.assertLess(LogLevel.DEBUG, LogLevel.INFO)
        self.assertLess(LogLevel.INFO, LogLevel.WARNING)
        self.assertLess(LogLevel.WARNING, LogLevel.ERROR)
        
        self.assertGreater(LogLevel.ERROR, LogLevel.WARNING)
        self.assertGreater(LogLevel.WARNING, LogLevel.INFO)
        self.assertGreater(LogLevel.INFO, LogLevel.DEBUG)
    
    def test_log_level_is_enum(self):
        """Test that LogLevel is an IntEnum."""
        self.assertTrue(issubclass(LogLevel, int))
        self.assertEqual(LogLevel.DEBUG, 10)
        self.assertEqual(LogLevel.INFO, 20)
        self.assertEqual(LogLevel.WARNING, 30)
        self.assertEqual(LogLevel.ERROR, 40)


if __name__ == '__main__':
    unittest.main()
