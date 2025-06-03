#!/usr/bin/env python3
"""
Base snippet processor module for processing code snippets.
"""

from abc import ABC, abstractmethod
from typing import Dict, Tuple, Optional
from .snippetfilter import SnippetFilter
from .location import Location
from .exceptions import SnippetNotFoundError


class SnippetProcessor(ABC):
    """
    Abstract base class for processing code snippets from source files.
    
    This class defines the interface for snippet processors that extract
    and process code elements based on specific file types and requirements.
    """
    
    def __init__(self, snippet_filter: SnippetFilter) -> None:
        """
        Initialize a SnippetProcessor instance.
        
        Args:
            snippet_filter: A SnippetFilter instance to work with
        """
        self._snippet_filter: SnippetFilter = snippet_filter
    
    @property
    def snippet_filter(self) -> SnippetFilter:
        """
        Get the snippet filter (readonly property).
        
        Returns:
            The SnippetFilter instance
        """
        return self._snippet_filter
    
    @property
    def source_code_path(self) -> str:
        """
        Get the source code path from the snippet filter.
        
        Returns:
            The source code file path
        """
        return self._snippet_filter.source_code_path
    
    @property
    def code_elements(self) -> Dict[str, str]:
        """
        Get the extracted code elements from the snippet filter.
        
        Returns:
            Dictionary mapping element names to their content
        """
        return self._snippet_filter.code_elements
    
    @abstractmethod
    def extract_snippet(self) -> str:
        """
        Extract and return a formatted code snippet.
        
        Returns:
            The formatted code snippet as a string
            
        Raises:
            SnippetNotFoundError: If the snippet is not found in the source code
        """
        pass

    @abstractmethod
    def get_snippet_location(self) -> Optional[Tuple[Location, Location]]:
        """
        Get the start and end location of a specific code snippet.
        
        Returns:
            A tuple containing (start_location, end_location) if found, None otherwise.
            The locations represent the beginning and end of the code snippet.
        """
        pass

    @abstractmethod
    def apply_snippet_change(self, modified_code: str) -> None:
        """
        Apply modified code to the source file.
        
        Args:
            modified_code: The modified code content to write to the file
            
        Returns:
            None
            
        Raises:
            IOError: If there's an error writing to the file
            FileNotFoundError: If the source file doesn't exist
            PermissionError: If there are permission issues writing to the file
        """
        pass

    @abstractmethod
    def get_system_prompt(self) -> str:
        """
        Get the system prompt for AI-assisted code processing.
        
        This method should return a language-specific system prompt that provides
        context and instructions for AI models when processing code snippets.
        Each concrete implementation should define its own prompt based on the
        programming language and specific requirements.
        
        Returns:
            A string containing the system prompt for the specific programming language
        """
        pass

