#!/usr/bin/env python3
"""
File type module for determining file types from extensions.
"""

from typing import Optional
from enum import Enum


class FileType(Enum):
    """
    Represents different file types and provides functionality to determine file type from extensions.
    """
    GOLANG = "golang"
    JSON = "json"
    PS1 = "ps1"
    BASH = "bash"
    MAKEFILE = "makefile"
    MARKDOWN = "markdown"
    
    @classmethod
    def get_filetype(cls, file_path: str) -> Optional['FileType']:
        """
        Determine the file type from a file path or extension.
        
        Args:
            file_path: The file path or filename to analyze
            
        Returns:
            The corresponding FileType enum, or None if not recognized
        """
        if not file_path:
            return None
            
        # Extract filename from path
        filename = file_path.split('/')[-1]
        
        # Handle special case for Makefile (no extension, exact name match)
        if filename == 'Makefile':
            return cls.MAKEFILE
            
        # Get the file extension for other file types
        if '.' in filename:
            extension = '.' + filename.split('.')[-1].lower()
            
            # File extension mappings
            extension_map = {
                '.go': cls.GOLANG,
                '.json': cls.JSON,
                '.ps1': cls.PS1,
                '.sh': cls.BASH,
                '.bash': cls.BASH,
                '.md': cls.MARKDOWN,
                '.markdown': cls.MARKDOWN,
                '.mk': cls.MAKEFILE,
            }
            
            return extension_map.get(extension)
            
        return None
    
    @classmethod
    def get_all_types(cls) -> list['FileType']:
        """
        Get a list of all available file types.
        
        Returns:
            List of all FileType enums
        """
        return [cls.GOLANG, cls.JSON, cls.PS1, cls.BASH, cls.MAKEFILE, cls.MARKDOWN]
