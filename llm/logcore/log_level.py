"""
Log level definitions for the logger.
"""

import enum
import logging


class LogLevel(enum.IntEnum):
    """Log level enumeration matching standard logging levels."""
    
    DEBUG = logging.DEBUG
    INFO = logging.INFO
    WARNING = logging.WARNING
    ERROR = logging.ERROR
