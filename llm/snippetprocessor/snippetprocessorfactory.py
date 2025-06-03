#!/usr/bin/env python3
"""
Factory module for creating concrete snippet processor instances based on source code path.
"""

from typing import Optional
from .snippetfilter import SnippetFilter
from .filetype import FileType
from .snippetprocessor import SnippetProcessor
from .gosnippetprocessor import GoSnippetProcessor
from .exceptions import SnippetNotFoundError


class UnsupportedFileTypeError(Exception):
    """Raised when trying to create a processor for an unsupported file type."""
    
    def __init__(self, file_path: str, file_type: Optional[FileType] = None) -> None:
        """
        Initialize UnsupportedFileTypeError.
        
        Args:
            file_path: The file path that caused the error
            file_type: The detected file type (if any)
        """
        self.file_path: str = file_path
        self.file_type: Optional[FileType] = file_type
        
        if file_type:
            message = f"Unsupported file type '{file_type.value}' for file: {file_path}"
        else:
            message = f"Unable to determine file type for: {file_path}"
            
        super().__init__(message)


class SnippetProcessorFactory:
    """
    Factory class for creating concrete snippet processor instances.
    
    This factory examines the source code path from a SnippetFilter and creates
    the appropriate concrete snippet processor based on the file type.
    
    Following the Factory Pattern and Pythonic OOD principles:
    - Single responsibility: Only creates snippet processors
    - Open/closed principle: Easy to extend with new processor types
    - Dependency inversion: Depends on abstractions, not concrete classes
    """
    
    @classmethod
    def create_processor(cls, snippet_filter: SnippetFilter) -> SnippetProcessor:
        """
        Create a concrete snippet processor based on the source code path.
        
        Args:
            snippet_filter: A SnippetFilter instance containing source code path and elements
            
        Returns:
            A concrete SnippetProcessor instance appropriate for the file type
            
        Raises:
            ValueError: If snippet_filter is invalid
            UnsupportedFileTypeError: If the file type is not supported
            SnippetNotFoundError: If the source file doesn't exist or is inaccessible
        """
        if not snippet_filter:
            raise ValueError("snippet_filter cannot be None")
            
        if not snippet_filter.source_code_path:
            raise ValueError("snippet_filter must have a valid source_code_path")
        
        # Determine the file type from the source code path
        file_type = FileType.get_filetype(snippet_filter.source_code_path)
        
        # Create the appropriate processor based on file type
        return cls._create_processor_by_type(snippet_filter, file_type)
    
    @classmethod
    def _create_processor_by_type(
        cls, 
        snippet_filter: SnippetFilter, 
        file_type: Optional[FileType]
    ) -> SnippetProcessor:
        """
        Create a processor instance based on the detected file type.
        
        Args:
            snippet_filter: The SnippetFilter instance
            file_type: The detected FileType, or None if unrecognized
            
        Returns:
            A concrete SnippetProcessor instance
            
        Raises:
            UnsupportedFileTypeError: If the file type is not supported
        """
        # Use match/case syntax for modern Python pattern matching
        match file_type:
            case FileType.GOLANG:
                return GoSnippetProcessor(snippet_filter)
            
            case FileType.JSON:
                # TODO: Implement JsonSnippetProcessor when needed
                raise UnsupportedFileTypeError(snippet_filter.source_code_path, file_type)
            
            case FileType.PS1:
                # TODO: Implement PowerShellSnippetProcessor when needed
                raise UnsupportedFileTypeError(snippet_filter.source_code_path, file_type)
            
            case FileType.BASH:
                # TODO: Implement BashSnippetProcessor when needed
                raise UnsupportedFileTypeError(snippet_filter.source_code_path, file_type)
            
            case FileType.MAKEFILE:
                # TODO: Implement MakefileSnippetProcessor when needed
                raise UnsupportedFileTypeError(snippet_filter.source_code_path, file_type)
            
            case FileType.MARKDOWN:
                # TODO: Implement MarkdownSnippetProcessor when needed
                raise UnsupportedFileTypeError(snippet_filter.source_code_path, file_type)
            
            case None:
                # File type could not be determined
                raise UnsupportedFileTypeError(snippet_filter.source_code_path, None)
            
            case _:
                # Unknown file type (should not happen with current enum)
                raise UnsupportedFileTypeError(snippet_filter.source_code_path, file_type)
    
    @classmethod
    def get_supported_file_types(cls) -> list[FileType]:
        """
        Get a list of currently supported file types.
        
        Returns:
            List of FileType enums that have concrete processor implementations
        """
        return [FileType.GOLANG]
    
    @classmethod
    def is_file_type_supported(cls, file_path: str) -> bool:
        """
        Check if a file type is supported by examining its path.
        
        Args:
            file_path: The file path to check
            
        Returns:
            True if the file type is supported, False otherwise
        """
        if not file_path:
            return False
            
        file_type = FileType.get_filetype(file_path)
        return file_type in cls.get_supported_file_types()
    
    @classmethod
    def get_processor_type_name(cls, snippet_filter: SnippetFilter) -> str:
        """
        Get the name of the processor type that would be created for this filter.
        
        Args:
            snippet_filter: The SnippetFilter to analyze
            
        Returns:
            The name of the processor class that would be created
            
        Raises:
            ValueError: If snippet_filter is invalid
            UnsupportedFileTypeError: If the file type is not supported
        """
        if not snippet_filter or not snippet_filter.source_code_path:
            raise ValueError("Invalid snippet_filter")
            
        file_type = FileType.get_filetype(snippet_filter.source_code_path)
        
        match file_type:
            case FileType.GOLANG:
                return "GoSnippetProcessor"
            case FileType.JSON:
                return "JsonSnippetProcessor"
            case FileType.PS1:
                return "PowerShellSnippetProcessor"
            case FileType.BASH:
                return "BashSnippetProcessor"
            case FileType.MAKEFILE:
                return "MakefileSnippetProcessor"
            case FileType.MARKDOWN:
                return "MarkdownSnippetProcessor"
            case None:
                raise UnsupportedFileTypeError(snippet_filter.source_code_path, None)
            case _:
                raise UnsupportedFileTypeError(snippet_filter.source_code_path, file_type)
