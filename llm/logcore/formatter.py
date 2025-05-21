"""
Log formatter for structured logs.

This module provides a formatter that formats log records in a consistent
JSON-based structure for better traceability and machine parsing.
"""

import datetime
import json
import logging
import os
from typing import Dict, Any


class StructuredLogFormatter(logging.Formatter):
    """
    Formatter for structured logs in the specified format.
    
    Format: [TIMESTAMP] [LEVEL] [FILENAME:LINENO] [FUNCTION] {"Key1":"Value1", "Key2":"Value2"}
    """
    
    def format(self, record: logging.LogRecord) -> str:
        """Format the log record according to the required format."""
        # Get the timestamp in ISO 8601 format
        timestamp = datetime.datetime.fromtimestamp(
            record.created, datetime.timezone.utc
        ).isoformat(timespec='milliseconds').replace('+00:00', 'Z')
        
        # Get the source file name, line number, and function name
        filename = record.pathname.split(os.sep)[-1] if record.pathname else None
        lineno = record.lineno if hasattr(record, 'lineno') else None
        func_name = record.funcName if hasattr(record, 'funcName') else None
        
        # Create structured log message with additional context
        log_data: Dict[str, Any] = {}
        
        # If there's log data in the message (assuming it's already in JSON format)
        if isinstance(record.msg, dict):
            log_data.update(record.msg)
        else:
            log_data["message"] = record.msg
            
        # Format the header part of the log
        header_parts = [
            f"[{timestamp}]",
            f"[{record.levelname}]",
            f"[{filename}:{lineno}]",
            f"[{func_name}]"
        ]
        
        # Format the data part as JSON
        json_data = json.dumps(log_data, ensure_ascii=False)
        
        return " ".join(header_parts) + " " + json_data
