#!/usr/bin/env python3
"""
Instruction module for handling instruction objects.
"""

import os
from typing import Optional


class Instruction:
    """
    Represents an instruction with file path and content.
    
    This class encapsulates instruction data with a file path and lazily loads
    content from the file when first accessed.
    """
    
    def __init__(self, file_path: str) -> None:
        """
        Initialize an Instruction instance.
        
        Args:
            file_path: The path to the instruction file
            
        Raises:
            FileNotFoundError: If the specified file does not exist
        """
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"Instruction file not found: {file_path}")
        
        self._file_path: str = file_path
        self._content: Optional[str] = None
    
    @property
    def file_path(self) -> str:
        """
        Get the file path (readonly property).
        
        Returns:
            The file path of the instruction
        """
        return self._file_path
    
    def get_content(self) -> str:
        """
        Get the content of the instruction.
        
        Lazily loads content from file on first access and caches it
        for subsequent calls.
        
        Returns:
            The content of the instruction
            
        Raises:
            IOError: If there's an error reading the file
        """
        if self._content is None:
            try:
                with open(self._file_path, 'r', encoding='utf-8') as file:
                    self._content = file.read()
            except IOError as e:
                raise IOError(f"Error reading instruction file '{self._file_path}': {e}")
        
        return self._content

class InstructionLoader:
    """
    Loads instructions from a specified directory.
    """

    def __init__(self, directory: str) -> None:
        """
        Initialize an InstructionLoader instance.

        Args:
            directory: The directory containing instruction files
            
        Raises:
            FileNotFoundError: If the directory does not exist
            NotADirectoryError: If the path exists but is not a directory
        """
        if not os.path.exists(directory):
            raise FileNotFoundError(f"Directory not found: {directory}")
        
        if not os.path.isdir(directory):
            raise NotADirectoryError(f"Path is not a directory: {directory}")
        
        self._directory: str = directory

    def load_instructions(self) -> list[Instruction]:
        """
        Load all instruction files from the specified directory.
        
        Only files with .md extension are loaded as instructions.

        Returns:
            A list of Instruction objects

        Raises:
            FileNotFoundError: If the directory does not exist
        """
        if not os.path.exists(self._directory):
            raise FileNotFoundError(f"Instruction directory not found: {self._directory}")

        instructions = []
        for filename in os.listdir(self._directory):
            if filename.endswith(".md"):
                file_path = os.path.join(self._directory, filename)
                instructions.append(Instruction(file_path))

        return instructions