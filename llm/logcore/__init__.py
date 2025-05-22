"""
Logging component for code update workflow.

This module provides a unified interface for recording actions, changes,
and errors throughout the code update workflow.
"""

from .logger import Logger, get_logger, setup_logging

__all__ = ["Logger", "get_logger", "setup_logging"]
