"""
Core logging package for the application.

This package provides a structured logging implementation with configurable
verbosity levels, automatic contextual information, and multiple output destinations.
"""

from .log_level import LogLevel
from .destinations import LogDestination, StdoutDestination, FileDestination
from .formatter import StructuredLogFormatter
from .logger import Logger
from .config import configure_logger, get_logger

__all__ = [
    'LogLevel',
    'LogDestination',
    'StdoutDestination',
    'FileDestination',
    'StructuredLogFormatter',
    'Logger',
    'configure_logger',
    'get_logger',
]
