"""
Log destination implementations for the logger.

This module provides various output destinations for logs, including
standard output and file output with rotation.
"""

import datetime
import logging
import logging.handlers
import os
import sys
from abc import ABC, abstractmethod
from pathlib import Path
from typing import Optional, Union

# NOTE: Windows-specific modules previously used for file sharing were removed
# in favor of Python's built-in cross-platform file handling


class LogDestination(ABC):
    """Abstract base class for log destinations."""
    
    @abstractmethod
    def setup(self, formatter: logging.Formatter) -> None:
        """Set up the log destination with the given formatter."""
        pass
    
    @abstractmethod
    def get_handler(self) -> logging.Handler:
        """Get the logging handler for this destination."""
        pass
        
    @abstractmethod
    def close(self) -> None:
        """Close the log destination and release any resources."""
        pass


class StdoutDestination(LogDestination):
    """Standard output log destination."""
    
    def __init__(self) -> None:
        self._handler: Optional[logging.StreamHandler] = None
    
    def setup(self, formatter: logging.Formatter) -> None:
        """Set up stdout logging handler."""
        self._handler = logging.StreamHandler(sys.stdout)
        self._handler.setFormatter(formatter)
    
    def get_handler(self) -> logging.Handler:
        """Get the stdout logging handler."""
        if self._handler is None:
            raise RuntimeError("StdoutDestination not set up")
        return self._handler
        
    def close(self) -> None:
        """
        Close the stdout handler.
        
        For StdoutDestination, this just flushes the stream
        since we don't want to close stdout.
        """
        if self._handler is not None:
            self._handler.flush()


class FileDestination(LogDestination):
    """File output log destination with rotation support."""
    
    def __init__(
        self, 
        log_dir: Union[str, Path], 
        max_size_bytes: int = 5 * 1024 * 1024,  # 5 MB
        max_total_size_bytes: int = 500 * 1024 * 1024,  # 500 MB
        app_name: str = "codeupdater"
    ) -> None:
        """
        Initialize a file destination for logs.
        
        Args:
            log_dir: Directory where log files will be stored
            max_size_bytes: Maximum size of each log file before rotation
            max_total_size_bytes: Maximum total size of all log files
            app_name: Application name prefix for log file names
        """
        self._log_dir = Path(log_dir)
        self._max_size_bytes = max_size_bytes
        self._max_total_size_bytes = max_total_size_bytes
        self._app_name = app_name
        self._handler: Optional[logging.handlers.RotatingFileHandler] = None
        self._log_path: Optional[Path] = None
          # Create log directory if it doesn't exist
        os.makedirs(self._log_dir, exist_ok=True)
        
    def setup(self, formatter: logging.Formatter) -> None:
        """Set up file logging handler with rotation."""
        timestamp = datetime.datetime.now(datetime.timezone.utc).strftime('%Y%m%d-%H%M%S')
        log_filename = f"{self._app_name}-{timestamp}.log"
        self._log_path = self._log_dir / log_filename
        
        # Create the rotating file handler using Python's built-in functionality
        # Setting encoding explicitly and using 'w' mode to force creation of new file
        self._handler = logging.handlers.RotatingFileHandler(
            self._log_path,
            maxBytes=self._max_size_bytes,
            backupCount=int(self._max_total_size_bytes / self._max_size_bytes),
            delay=False,  # Open immediately
            mode='a',     # Append mode
            encoding='utf-8',
            errors='replace'
        )
        
        self._handler.setFormatter(formatter)
        
        # Enable immediate flushing to ensure logs are written immediately
        self._handler.setLevel(logging.NOTSET)
        
        # Configure the handler to flush after each log entry
        # This ensures other processes can see updates more quickly
        original_emit = self._handler.emit
        
        def emit_and_flush(record):
            original_emit(record)
            self._handler.flush()
            
        self._handler.emit = emit_and_flush
    
    def get_handler(self) -> logging.Handler:
        """Get the file logging handler."""
        if self._handler is None:
            raise RuntimeError("FileDestination not set up")
        return self._handler
        
    def close(self) -> None:
        """
        Close the file handler to release the file lock.
        
        This method should be called when the logger is no longer needed
        to ensure proper file cleanup and allow other processes to access the log.
        """
        if self._handler is not None:
            self._handler.flush()
            self._handler.close()
            self._handler = None
