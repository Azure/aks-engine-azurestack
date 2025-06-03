#!/usr/bin/env python3
"""
Snippet filter module for extracting code elements from source files.
"""

import json
from dataclasses import dataclass, field
from typing import Dict, Union, Self


@dataclass(frozen=True)
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
    
    source_code_path: str
    code_elements: Dict[str, str] = field(default_factory=dict)
    
    def __post_init__(self) -> None:
        """Validate the instance after initialization."""
        if not self.source_code_path or not self.source_code_path.strip():
            raise ValueError("source_code_path cannot be empty")
    
    @classmethod
    def from_json(cls, json_data: Union[str, Dict]) -> Self:
        """
        Create a SnippetFilter instance from JSON data.
        
        Args:
            json_data: Either a JSON string or a dictionary containing the filter data
                      Expected format: {
                          "source_code_path": "<path>",
                          "code_elements": {
                              "<element_name>": "<element_value>",
                              ...
                          }
                      }
        
        Returns:
            A SnippetFilter instance initialized with the JSON data
            
        Raises:
            ValueError: If the JSON data is invalid or missing required fields
            TypeError: If json_data is neither string nor dict
        """
        # Parse JSON string if needed
        match json_data:
            case str():
                try:
                    data = json.loads(json_data)
                except json.JSONDecodeError as e:
                    raise ValueError(f"Failed to parse JSON string: {e}") from e
            case dict():
                data = json_data
            case _:
                raise TypeError(f"json_data must be str or dict, got {type(json_data)}")
        
        # Validate that data is a dictionary
        if not isinstance(data, dict):
            raise ValueError("Parsed JSON data must be a dictionary")
        
        # Extract and validate source_code_path
        source_code_path = data.get("source_code_path", "")
        if not source_code_path or not isinstance(source_code_path, str):
            raise ValueError("source_code_path is required and must be a non-empty string")
        
        # Extract and validate code_elements
        code_elements = data.get("code_elements", {})
        if code_elements and not isinstance(code_elements, dict):
            raise ValueError("code_elements must be a dictionary")
        
        # Validate that all keys and values are strings
        if code_elements:
            for key, value in code_elements.items():
                if not isinstance(key, str) or not isinstance(value, str):
                    raise ValueError("All code_elements keys and values must be strings")
        
        return cls(source_code_path=source_code_path, code_elements=code_elements)
    
    def to_dict(self) -> Dict[str, Union[str, Dict[str, str]]]:
        """
        Convert the SnippetFilter to a dictionary.
        
        Returns:
            Dictionary representation of the SnippetFilter
        """
        return {
            "source_code_path": self.source_code_path,
            "code_elements": self.code_elements.copy()
        }
    
    def to_json(self) -> str:
        """
        Convert the SnippetFilter to a JSON string.
        
        Returns:
            JSON string representation of the SnippetFilter
        """
        return json.dumps(self.to_dict(), indent=2)
    
    def __str__(self) -> str:
        """
        Return a string representation of the SnippetFilter.
        
        Returns:
            A formatted string showing the source code path and code elements content
        """
        if self.code_elements:
            elements_list = []
            for name, content in self.code_elements.items():
                # Truncate long content for readability
                display_content = content[:100] + "..." if len(content) > 100 else content
                # Replace newlines with \n for single-line display
                display_content = display_content.replace('\n', '\\n')
                elements_list.append(f"'{name}': '{display_content}'")
            elements_str = "{" + ", ".join(elements_list) + "}"
        else:
            elements_str = "{}"
        
        return f"SnippetFilter(source_code_path='{self.source_code_path}', code_elements={elements_str})"
