"""
Pytest configuration file for logcore tests.

This file is automatically loaded by pytest and helps set up the Python path
to allow importing the logcore module from its parent directory.
"""

import os
import sys
from pathlib import Path

# Add the parent directory to the Python path so that imports work correctly
sys.path.insert(0, str(Path(__file__).parent.parent))
