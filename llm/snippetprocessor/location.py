#!/usr/bin/env python3
"""
Location module for representing code positions.
"""


class Location:
    """
    Represents a specific location in source code with line and column numbers.
    
    This class provides functionality to track positions in source code files
    for precise snippet extraction and navigation.
    """
    
    def __init__(self, line_number: int, column_number: int) -> None:
        """
        Initialize a Location instance.
        
        Args:
            line_number: The line number (1-based)
            column_number: The column number (1-based)
        """
        if line_number <= 0:
            raise ValueError("Line number must be positive")
        if column_number <= 0:
            raise ValueError("Column number must be positive")
            
        self._line_number: int = line_number
        self._column_number: int = column_number
    
    @property
    def line_number(self) -> int:
        """
        Get the line number (readonly property).
        
        Returns:
            The line number (1-based)
        """
        return self._line_number
    
    @property
    def column_number(self) -> int:
        """
        Get the column number (readonly property).
        
        Returns:
            The column number (1-based)
        """
        return self._column_number
    
    def __str__(self) -> str:
        """
        String representation of the location.
        
        Returns:
            String in format "line:column"
        """
        return f"{self._line_number}:{self._column_number}"
    
    def __repr__(self) -> str:
        """
        Official string representation of the location.
        
        Returns:
            String representation for debugging
        """
        return f"Location(line_number={self._line_number}, column_number={self._column_number})"
    
    def __eq__(self, other: object) -> bool:
        """
        Check equality with another Location instance.
        
        Args:
            other: Another object to compare with
            
        Returns:
            True if both locations have the same line and column numbers
        """
        if not isinstance(other, Location):
            return False
        return (self._line_number == other._line_number and 
                self._column_number == other._column_number)
    
    def __lt__(self, other: 'Location') -> bool:
        """
        Check if this location comes before another location.
        
        Args:
            other: Another Location instance to compare with
            
        Returns:
            True if this location comes before the other location
        """
        if self._line_number != other._line_number:
            return self._line_number < other._line_number
        return self._column_number < other._column_number
