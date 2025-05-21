"""
Common test fixtures and configurations for logcore tests.
"""

import logging
import os
import sys
import tempfile
from pathlib import Path

import pytest

# Add the parent directory to the Python path for proper imports
parent_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
if parent_dir not in sys.path:
    sys.path.insert(0, parent_dir)

# Reset root logger before tests to avoid interference from other tests
def reset_logging():
    """Reset the logging module to its initial state."""
    root = logging.getLogger()
    for handler in root.handlers[:]:
        root.removeHandler(handler)
    root.setLevel(logging.WARNING)


@pytest.fixture(autouse=True)
def reset_logging_before_test():
    """Reset logging before each test."""
    reset_logging()
    yield
    reset_logging()


@pytest.fixture
def temp_log_dir():
    """Create a temporary directory for log files."""
    with tempfile.TemporaryDirectory() as temp_dir:
        yield Path(temp_dir)
