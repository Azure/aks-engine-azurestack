#!/usr/bin/env python3
"""
Exception classes for snippet processing.
"""


class SnippetNotFoundError(Exception):
    """
    Exception raised when a requested code snippet is not found.
    
    This exception is raised by snippet processors when they cannot
    locate the requested code element in the source file.
    """
    
    def __init__(self, snippet_filter: str, source_path: str) -> None:
        """
        Initialize a SnippetNotFoundError.
        
        Args:
            snippet_filter: The string representation of the SnippetFilter that was not found
            source_path: The path to the source file where the snippet was searched
        """
        self.snippet_filter = snippet_filter
        self.source_path = source_path
        message = f"Snippet not found in source file '{source_path}' for filter: {snippet_filter}"
        super().__init__(message)
