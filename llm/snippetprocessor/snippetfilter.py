#!/usr/bin/env python3
"""
Snippet filter module for extracting code elements from source files.
"""

from typing import Dict


class SnippetFilter:
    """
    Extracts and manages code elements from source code files.
    
    This class provides functionality to extract specific code elements based on file type:
    - Go files: function names, type definitions, struct definitions
    - PowerShell files: function names
    - Markdown files: headers
    - Other files: relevant code elements
    
    The extracted elements are stored as a mapping of element names to their content.
    """
    
    def __init__(self, source_code_path: str) -> None:
        """
        Initialize a SnippetFilter instance.
        
        Args:
            source_code_path: The path to the source code file
        """
        self._source_code_path: str = source_code_path
        self._code_elements: Dict[str, str] = {}
    
    @property
    def source_code_path(self) -> str:
        """
        Get the source code path (readonly property).
        
        Returns:
            The source code file path
        """
        return self._source_code_path
    
    @property
    def code_elements(self) -> Dict[str, str]:
        """
        Get the extracted code elements mapping (readonly property).
        
        Returns:
            Dictionary mapping element names to their content:
            - For Go: function names, type names, struct names
            - For PowerShell: function names
            - For Markdown: header text
            - For other files: relevant code elements
        """
        return self._code_elements.copy()  # Return a copy to maintain readonly behavior

    def __str__(self) -> str:
        """
        Return a string representation of the SnippetFilter.
        
        Returns:
            A formatted string showing the source code path and code elements content
        """
        elements_str = ""
        if self._code_elements:
            elements_list = []
            for name, content in self._code_elements.items():
                # Truncate long content for readability
                display_content = content[:100] + "..." if len(content) > 100 else content
                # Replace newlines with \n for single-line display
                display_content = display_content.replace('\n', '\\n')
                elements_list.append(f"'{name}': '{display_content}'")
            elements_str = "{" + ", ".join(elements_list) + "}"
        else:
            elements_str = "{}"
        
        return f"SnippetFilter(source_code_path='{self._source_code_path}', code_elements={elements_str})"
