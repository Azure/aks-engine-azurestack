"""
Configuration functions for the logger.

This module provides functions to configure the logger.
"""

from pathlib import Path
from typing import Optional, Union, List

from .log_level import LogLevel
from .logger import Logger
from .destinations import StdoutDestination, FileDestination, LogDestination


def configure_logger(
    level: LogLevel = LogLevel.INFO,
    log_to_stdout: bool = True,
    log_to_file: bool = False,
    log_dir: Optional[Union[str, Path]] = None,
    max_size_bytes: int = 5 * 1024 * 1024,
    max_total_size_bytes: int = 500 * 1024 * 1024,
    app_name: str = "codeupdater"
) -> None:
    """
    Configure the logger with the specified settings.
    
    Args:
        level: Minimum log level to record
        log_to_stdout: Whether to log to standard output
        log_to_file: Whether to log to a file
        log_dir: Directory for log files (required if log_to_file is True)
        max_size_bytes: Maximum size of each log file before rotation
        max_total_size_bytes: Maximum total size of all log files
        app_name: Application name prefix for log file names
    """
    destinations: List[LogDestination] = []
    
    if log_to_stdout:
        destinations.append(StdoutDestination())
    
    if log_to_file:
        if log_dir is None:
            raise ValueError("log_dir must be provided when log_to_file is True")
        destinations.append(
            FileDestination(
                log_dir=log_dir, 
                max_size_bytes=max_size_bytes,
                max_total_size_bytes=max_total_size_bytes,
                app_name=app_name
            )
        )
    
    Logger._configure_root_logger(level, destinations)


def get_logger(name: str) -> Logger:
    """
    Get or create a logger with the given name.
    
    Args:
        name: Name of the logger
        
    Returns:
        Logger instance
    """
    return Logger.get_logger(name)
