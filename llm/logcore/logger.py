"""
Logger implementation for code update workflow.

This module provides a Logger class that wraps Python's standard logging
module to provide structured JSON logs with automatic contextual information.
"""

import json
import logging
import os
import socket
import sys
import threading
import time
import traceback
from abc import ABC, abstractmethod
from datetime import datetime, timezone
from inspect import currentframe, getframeinfo
from pathlib import Path
from typing import Any, Dict, Literal, Optional, Union

# Define log levels as a LiteralType for type hints
LogLevel = Literal["DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"]

class LogFormatter(logging.Formatter):
    """Custom log formatter that outputs log entries as JSON objects."""
    
    def __init__(self) -> None:
        """Initialize the formatter."""
        super().__init__()
        self.hostname = socket.gethostname()
    
    def format(self, record: logging.LogRecord) -> str:
        """Format the log record as a JSON string.
        
        Args:
            record: The log record to format.
            
        Returns:
            A JSON string representation of the log record.
        """
        # Extract the caller information if available in extra data
        source_file_path = getattr(record, "source_file_path", record.pathname)
        line_number = getattr(record, "line_number", record.lineno)
        function_name = getattr(record, "function_name", record.funcName)
        
        # Create the base log entry
        log_entry = {
            "timestamp": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z",
            "sourceFilePath": source_file_path,
            "lineNumber": line_number,
            "functionName": function_name,
            "processId": os.getpid(),
            "threadId": threading.get_ident(),
            "hostName": self.hostname,
            "level": record.levelname,
            "message": record.getMessage(),
        }
        
        # Add any additional context from the record
        if hasattr(record, "context") and isinstance(record.context, dict):
            log_entry.update(record.context)
            
        return json.dumps(log_entry)


class LogHandler(ABC):
    """Abstract base class for log handlers."""
    
    @abstractmethod
    def setup(self, formatter: logging.Formatter) -> logging.Handler:
        """Set up and return a logging.Handler.
        
        Args:
            formatter: The formatter to use for the handler.
            
        Returns:
            A configured logging.Handler instance.
        """
        pass


class ConsoleHandler(LogHandler):
    """Handler that outputs logs to the console."""
    
    def setup(self, formatter: logging.Formatter) -> logging.Handler:
        """Set up and return a StreamHandler for console output.
        
        Args:
            formatter: The formatter to use for the handler.
            
        Returns:
            A configured StreamHandler instance.
        """
        handler = logging.StreamHandler(sys.stdout)
        handler.setFormatter(formatter)
        return handler


class FileHandler(LogHandler):
    """Handler that outputs logs to a file."""
    
    def __init__(self, file_path: str) -> None:
        """Initialize the file handler.
        
        Args:
            file_path: Path to the log file.
        """
        self.file_path = file_path
    
    def setup(self, formatter: logging.Formatter) -> logging.Handler:
        """Set up and return a FileHandler.
        
        Args:
            formatter: The formatter to use for the handler.
            
        Returns:
            A configured FileHandler instance.
        """
        # Ensure the directory exists
        Path(self.file_path).parent.mkdir(parents=True, exist_ok=True)
        
        handler = logging.FileHandler(self.file_path)
        handler.setFormatter(formatter)
        return handler


class Logger:
    """Logger class that provides structured JSON logging with contextual information."""
    
    def __init__(
        self, 
        name: str, 
        level: LogLevel = "INFO",
        context: Optional[Dict[str, Any]] = None
    ) -> None:
        """Initialize a logger instance.
        
        Args:
            name: The logger name.
            level: The logging level. Defaults to "INFO".
            context: Additional context to include in all log entries. Defaults to None.
        """
        self._logger = logging.getLogger(name)
        self._logger.setLevel(getattr(logging, level))
        self._base_context = context or {}
    
    def _get_caller_info(self) -> tuple:
        """Get information about the caller of the log method.
        
        Returns:
            A tuple containing (file_path, line_number, function_name).
        """
        # Get the frame that called the log method (3 frames up from this function)
        # This goes beyond the wrapper methods (debug, info, etc.) to get the actual caller
        frame = currentframe().f_back.f_back.f_back
        frame_info = getframeinfo(frame)
        
        return frame_info.filename, frame_info.lineno, frame_info.function
    
    def _log(
        self, 
        level: LogLevel, 
        message: str, 
        context: Optional[Dict[str, Any]] = None
    ) -> None:
        """Log a message with the specified level and context.
        
        Args:
            level: The logging level for the message.
            message: The message to log.
            context: Additional context to include with this log entry. Defaults to None.
        """
        # Get caller information
        source_file_path, line_number, function_name = self._get_caller_info()
        
        # Combine base context with the provided context
        log_context = {**self._base_context}
        if context:
            log_context.update(context)
        
        # Create extra data for the log record
        extra = {
            "source_file_path": source_file_path,
            "line_number": line_number,
            "function_name": function_name,
            "context": log_context
        }
        
        # Log the message with the extra data
        getattr(self._logger, level.lower())(message, extra=extra)
    
    def debug(self, message: str, context: Optional[Dict[str, Any]] = None) -> None:
        """Log a debug message.
        
        Args:
            message: The message to log.
            context: Additional context to include with this log entry. Defaults to None.
        """
        self._log("DEBUG", message, context)
    
    def info(self, message: str, context: Optional[Dict[str, Any]] = None) -> None:
        """Log an info message.
        
        Args:
            message: The message to log.
            context: Additional context to include with this log entry. Defaults to None.
        """
        self._log("INFO", message, context)
    
    def warning(self, message: str, context: Optional[Dict[str, Any]] = None) -> None:
        """Log a warning message.
        
        Args:
            message: The message to log.
            context: Additional context to include with this log entry. Defaults to None.
        """
        self._log("WARNING", message, context)
    
    def error(
        self, 
        message: str, 
        exception: Optional[Exception] = None, 
        context: Optional[Dict[str, Any]] = None
    ) -> None:
        """Log an error message.
        
        Args:
            message: The message to log.
            exception: Optional exception to include in the log. Defaults to None.
            context: Additional context to include with this log entry. Defaults to None.
        """
        if exception:
            message = f"{message}: {str(exception)}"
            if context is None:
                context = {}
            context["traceback"] = traceback.format_exception(
                type(exception), exception, exception.__traceback__
            )
        
        self._log("ERROR", message, context)
    
    def critical(
        self, 
        message: str, 
        exception: Optional[Exception] = None, 
        context: Optional[Dict[str, Any]] = None
    ) -> None:
        """Log a critical message.
        
        Args:
            message: The message to log.
            exception: Optional exception to include in the log. Defaults to None.
            context: Additional context to include with this log entry. Defaults to None.
        """
        if exception:
            message = f"{message}: {str(exception)}"
            if context is None:
                context = {}
            context["traceback"] = traceback.format_exception(
                type(exception), exception, exception.__traceback__
            )
        
        self._log("CRITICAL", message, context)
    
    def with_context(self, context: Dict[str, Any]) -> 'Logger':
        """Create a new logger instance with additional context.
        
        Args:
            context: Additional context to include in all log entries.
            
        Returns:
            A new Logger instance with the combined context.
        """
        new_context = {**self._base_context, **context}
        return Logger(self._logger.name, context=new_context)
    
    def set_level(self, level: LogLevel) -> None:
        """Set the logging level for this logger.
        
        Args:
            level: The logging level to set.
        """
        self._logger.setLevel(getattr(logging, level))


# Global dictionary to store configured loggers
_loggers: Dict[str, Logger] = {}


def setup_logging(
    app_name: str = "code-updator",
    level: LogLevel = "INFO",
    handlers: Optional[list[LogHandler]] = None,
    default_context: Optional[Dict[str, Any]] = None
) -> None:
    """Set up the logging system.
    
    Args:
        app_name: The name of the application. Defaults to "code-updator".
        level: The default logging level. Defaults to "INFO".
        handlers: List of log handlers to use. Defaults to None which uses ConsoleHandler.
        default_context: Default context to include in all log entries. Defaults to None.
    """
    # Use a console handler if no handlers are provided
    if handlers is None:
        handlers = [ConsoleHandler()]
    
    # Configure root logger
    root_logger = logging.getLogger()
    root_logger.setLevel(getattr(logging, level))
    
    # Clear existing handlers
    for handler in root_logger.handlers[:]:
        root_logger.removeHandler(handler)
    
    # Create formatter
    formatter = LogFormatter()
    
    # Add handlers
    for handler in handlers:
        root_logger.addHandler(handler.setup(formatter))
    
    # Create application logger
    _loggers[app_name] = Logger(app_name, level, default_context)


def get_logger(
    name: Optional[str] = None, 
    context: Optional[Dict[str, Any]] = None
) -> Logger:
    """Get a logger instance.
    
    Args:
        name: The logger name. Defaults to None, which uses the default app name.
        context: Additional context to include in all log entries. Defaults to None.
        
    Returns:
        A Logger instance.
        
    Raises:
        RuntimeError: If setup_logging has not been called yet.
    """
    if not _loggers:
        # If no loggers are configured, set up the default logger
        setup_logging()
    
    app_name = next(iter(_loggers.keys()))
    base_logger = _loggers.get(name) if name in _loggers else _loggers[app_name]
    
    if context:
        return base_logger.with_context(context)
    
    return base_logger
