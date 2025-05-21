"""
Core logger implementation for the application.

This module provides a structured logger that supports multiple verbosity levels
and automatic contextual information.
"""

import inspect
import logging
import os
import socket
import threading
from typing import Any, Dict, List, Optional, Tuple

from .log_level import LogLevel
from .destinations import LogDestination


class Logger:
    """
    Structured logger implementation that fulfills the requirements.
    
    This class provides methods for logging at different verbosity levels
    with automatic contextual information and supports multiple output destinations.
    """
    
    _loggers: Dict[str, 'Logger'] = {}  # Registry of logger instances
    _root_logger: Optional[logging.Logger] = None
    _destinations: List[LogDestination] = []
    
    def __init__(self, name: str) -> None:
        """
        Initialize a new logger instance.
        
        Args:
            name: Name of the logger
        """
        self.name = name
        self._logger = logging.getLogger(name)
        
        # Ensure we don't add handlers multiple times
        if not self._logger.handlers and Logger._root_logger:
            self._logger.setLevel(Logger._root_logger.level)
            self._logger.propagate = False
    
    @classmethod
    def _configure_root_logger(cls, level: LogLevel, destinations: List[LogDestination]) -> None:
        """
        Configure the root logger with the specified level and destinations.
        
        Args:
            level: Minimum log level to record
            destinations: List of output destinations
        """
        # Store destinations for future loggers
        cls._destinations = destinations
        
        # Configure the root logger
        cls._root_logger = logging.getLogger()
        cls._root_logger.setLevel(level.value)
        cls._root_logger.handlers.clear()  # Remove any existing handlers
        
        # Create and set up the formatter
        from .formatter import StructuredLogFormatter
        formatter = StructuredLogFormatter()
        
        # Set up and add handlers for each destination
        for destination in destinations:
            destination.setup(formatter)
            cls._root_logger.addHandler(destination.get_handler())
    
    @classmethod
    def get_logger(cls, name: str) -> 'Logger':
        """
        Get or create a logger with the given name.
        
        Args:
            name: Name of the logger
            
        Returns:
            Logger instance
        """
        if name not in cls._loggers:
            cls._loggers[name] = Logger(name)
        return cls._loggers[name]
    
    def _get_caller_info(self) -> Tuple[str, int, str]:
        """
        Get information about the caller of the logging method.
        
        Returns:
            Tuple of (filename, line number, function name)
        """
        frame = inspect.currentframe()
        if frame:
            try:
                # Get the caller's frame (2 frames up from current)
                caller_frame = inspect.getouterframes(frame, 2)[2]
                return (
                    os.path.basename(caller_frame.filename),
                    caller_frame.lineno,
                    caller_frame.function
                )
            finally:
                # Properly dispose of the frame
                del frame
        return ("unknown", 0, "unknown")
    
    def _prepare_log_data(self, message: str, **kwargs: Any) -> Dict[str, Any]:
        """
        Prepare log data with automatic contextual information.
        
        Args:
            message: Log message
            **kwargs: Additional key-value pairs for the log entry
            
        Returns:
            Dictionary with log data
        """
        filename, lineno, func_name = self._get_caller_info()
        
        # Start with the provided log data
        log_data = {
            "message": message,
            **kwargs,
        }
        
        # Add contextual information
        log_data.update({
            "sourceFile": filename,
            "lineNumber": lineno,
            "functionName": func_name,
            "processId": os.getpid(),
            "threadId": threading.get_ident(),
            "hostName": socket.gethostname(),
        })
        
        return log_data
    
    def debug(self, message: str, **kwargs: Any) -> None:
        """
        Log a debug message.
        
        Args:
            message: Log message
            **kwargs: Additional key-value pairs for the log entry
        """
        log_data = self._prepare_log_data(message, **kwargs)
        self._logger.debug(log_data)
    
    def info(self, message: str, **kwargs: Any) -> None:
        """
        Log an info message.
        
        Args:
            message: Log message
            **kwargs: Additional key-value pairs for the log entry
        """
        log_data = self._prepare_log_data(message, **kwargs)
        self._logger.info(log_data)
    
    def warning(self, message: str, **kwargs: Any) -> None:
        """
        Log a warning message.
        
        Args:
            message: Log message
            **kwargs: Additional key-value pairs for the log entry
        """
        log_data = self._prepare_log_data(message, **kwargs)
        self._logger.warning(log_data)
    
    def error(self, message: str, **kwargs: Any) -> None:
        """
        Log an error message.
        
        Args:
            message: Log message
            **kwargs: Additional key-value pairs for the log entry
        """
        log_data = self._prepare_log_data(message, **kwargs)
        self._logger.error(log_data)



