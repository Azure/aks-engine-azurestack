#!/usr/bin/env python3
"""
Go-specific snippet processor for extracting Go code elements.
"""

import re
import subprocess
from typing import Tuple, Optional
from .snippetprocessor import SnippetProcessor
from .location import Location
from .exceptions import SnippetNotFoundError


class GoSnippetProcessor(SnippetProcessor):
    """
    Go-specific implementation of SnippetProcessor.
    
    This class provides functionality to extract Go-specific code elements:
    - Function definitions
    - Type definitions (structs, interfaces)
    - Method definitions
    - Variable declarations
    """
    
    def extract_snippet(self) -> str:
        """
        Extract and return a formatted Go code snippet.
        
        Returns:
            The formatted code snippet as a string with preserved spacing and indentation
            
        Raises:
            SnippetNotFoundError: If the snippet is not found in the source code
        """
        # Get the location of the snippet first
        location = self.get_snippet_location()
        if location is None:
            raise SnippetNotFoundError(str(self.snippet_filter), self.source_code_path)
        
        start_location, end_location = location
        
        try:
            with open(self.source_code_path, 'r', encoding='utf-8') as file:
                lines = file.readlines()
            
            # Extract lines based on location
            start_line_idx = start_location.line_number - 1  # Convert to 0-based index
            end_line_idx = end_location.line_number - 1      # Convert to 0-based index
            
            if start_line_idx < 0 or end_line_idx >= len(lines):
                raise SnippetNotFoundError("snippet", self.source_code_path)
            
            # Extract the snippet lines
            snippet_lines = []
            
            for line_idx in range(start_line_idx, end_line_idx + 1):
                line = lines[line_idx]
                
                if line_idx == start_line_idx and line_idx == end_line_idx:
                    # Single line snippet - use both start and end column positions
                    start_col = start_location.column_number - 1  # Convert to 0-based
                    end_col = end_location.column_number - 1      # Convert to 0-based
                    snippet_lines.append(line[start_col:end_col])
                elif line_idx == start_line_idx:
                    # First line - use start column to end of line
                    start_col = start_location.column_number - 1  # Convert to 0-based
                    snippet_lines.append(line[start_col:])
                elif line_idx == end_line_idx:
                    # Last line - use beginning of line to end column
                    end_col = end_location.column_number - 1      # Convert to 0-based
                    snippet_lines.append(line[:end_col])
                else:
                    # Middle lines - use entire line
                    snippet_lines.append(line)
            
            # Join the lines and return the snippet
            snippet = ''.join(snippet_lines)
            return snippet
            
        except (FileNotFoundError, IOError) as e:
            raise SnippetNotFoundError(str(self.snippet_filter), self.source_code_path) from e

    def get_snippet_location(self) -> Optional[Tuple[Location, Location]]:
        """
        Get the start and end location of a specific Go code snippet.
        
        Returns:
            A tuple containing (start_location, end_location) if found, None otherwise.
        """
        # Check if "object_type" is in code_elements to determine which extraction method to use
        if "object_type" in self.code_elements and "begin_with" in self.code_elements and self.code_elements["object_type"] in ("function","map"):
            # Get the first object name from code_elements for object extraction
            object_name = list(self.code_elements.keys())[0] if self.code_elements else ""
            return self._extract_block_location(self.code_elements["begin_with"])
        else:
            return self._extract_entire_file_location()

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
        # Get the location of the snippet to replace
        location = self.get_snippet_location()
        if location is None:
            raise SnippetNotFoundError(str(self.snippet_filter), self.source_code_path)
        
        start_location, end_location = location
        
        try:
            # Read the current file content
            with open(self.source_code_path, 'r', encoding='utf-8') as file:
                lines = file.readlines()
            
            # Convert to 0-based indices
            start_line_idx = start_location.line_number - 1
            end_line_idx = end_location.line_number - 1
            
            if start_line_idx < 0 or end_line_idx >= len(lines):
                raise IOError(f"Invalid line range in file: {self.source_code_path}")
            
            # Build the new file content
            new_lines = []
            
            # Add lines before the snippet
            new_lines.extend(lines[:start_line_idx])
            
            # Handle the replacement based on location boundaries
            if start_line_idx == end_line_idx:
                # Single line replacement - preserve parts of the line before/after the snippet
                original_line = lines[start_line_idx]
                start_col = start_location.column_number - 1  # Convert to 0-based
                end_col = end_location.column_number - 1      # Convert to 0-based
                
                # Construct the new line with replacement
                new_line = original_line[:start_col] + modified_code + original_line[end_col:]
                new_lines.append(new_line)
            else:
                # Multi-line replacement
                # Handle first line - keep content before start column
                first_line = lines[start_line_idx]
                start_col = start_location.column_number - 1  # Convert to 0-based
                line_prefix = first_line[:start_col]
                
                # Handle last line - keep content after end column
                last_line = lines[end_line_idx]
                end_col = end_location.column_number - 1      # Convert to 0-based
                line_suffix = last_line[end_col:]
                
                # Add the replacement content
                if modified_code.endswith('\n'):
                    # Modified code already has newline
                    replacement_content = line_prefix + modified_code + line_suffix
                else:
                    # Add newline if modified code doesn't have one
                    replacement_content = line_prefix + modified_code + line_suffix
                
                new_lines.append(replacement_content)
            
            # Add lines after the snippet
            new_lines.extend(lines[end_line_idx + 1:])
            
            # Write the modified content back to the file
            with open(self.source_code_path, 'w', encoding='utf-8') as file:
                file.writelines(new_lines)
            
            # Format the Go code using gofmt
            self._format_go_code()
                
        except (FileNotFoundError, PermissionError) as e:
            raise e
        except Exception as e:
            raise IOError(f"Error applying snippet change: {str(e)}") from e

    def _extract_block_location(self, begin_with_pattern: str) -> Optional[Tuple[Location, Location]]:
        """
        Extract the location of a Go code block that begins with a specific pattern
        and is delimited by curly braces { }.
        
        Args:
            begin_with_pattern: Regular expression pattern to match the beginning of the block
            
        Returns:
            A tuple containing (start_location, end_location) if found, None otherwise.
            
        Raises:
            FileNotFoundError: If the source file doesn't exist
            IOError: If there's an error reading the file
        """
        try:
            with open(self.source_code_path, 'r', encoding='utf-8') as file:
                lines = file.readlines()
            
            start_location = None
            brace_count = 0
            in_block = False
            
            for line_num, line in enumerate(lines, 1):
                # Check if this line contains the block beginning pattern
                if line.strip().startswith(begin_with_pattern):
                    start_location = Location(line_num, 1)
                    in_block = True
                    # Count opening braces in the declaration line
                    brace_count += line.count('{') - line.count('}')
                    
                    # If block completes on same line, return immediately
                    if brace_count == 0:
                        end_location = Location(line_num, len(line))
                        return (start_location, end_location)
                    continue
                
                if in_block:
                    # Count braces to find the end of the block
                    brace_count += line.count('{') - line.count('}')
                    
                    # Block ends when brace count returns to 0
                    if brace_count == 0 and start_location is not None:
                        end_location = Location(line_num, len(line))
                        return (start_location, end_location)
            
            # Block declaration found but no closing brace (malformed code)
            if start_location is not None and in_block:
                # Return up to the end of file
                end_line = len(lines)
                end_location = Location(end_line, len(lines[-1]) if lines else 1)
                return (start_location, end_location)
            
            return None
            
        except (FileNotFoundError, IOError) as e:
            raise e

    def _extract_entire_file_location(self) -> Optional[Tuple[Location, Location]]:
        """
        Extract the location representing the entire source code file.
        
        Returns:
            A tuple containing (start_location, end_location) representing the entire file,
            where start_location is (1, 1) and end_location is the last line and column.
            Returns None if the file is empty or doesn't exist.
            
        Raises:
            FileNotFoundError: If the source file doesn't exist
            IOError: If there's an error reading the file
        """
        try:
            with open(self.source_code_path, 'r', encoding='utf-8') as file:
                lines = file.readlines()
            
            if not lines:
                # Empty file
                return None
            
            # Start location is always the beginning of the file
            start_location = Location(1, 1)
            
            # End location is the last line and the end of that line
            last_line_number = len(lines)
            last_line_content = lines[-1]
            last_column = len(last_line_content)
            
            # If the last line doesn't end with a newline, add 1 to column position
            # to represent the position after the last character
            if not last_line_content.endswith('\n'):
                last_column += 1
            
            end_location = Location(last_line_number, last_column)
            
            return (start_location, end_location)
            
        except (FileNotFoundError, IOError) as e:
            raise e

    def _format_go_code(self) -> None:
        """
        Format the Go source file using gofmt.
        
        This method runs gofmt on the source file to ensure proper Go formatting.
        If gofmt is not available or encounters an error, it logs a warning but
        does not raise an exception to avoid breaking the workflow.
        
        Raises:
            IOError: Only if there's a critical file access issue that prevents formatting
        """
        try:
            # Run gofmt -w to format the file in place
            result = subprocess.run(
                ['gofmt', '-w', self.source_code_path],
                capture_output=True,
                text=True,
                timeout=30  # 30 second timeout
            )
            
            # gofmt returns 0 on success, non-zero on error
            if result.returncode != 0:
                # Log the error but don't raise an exception
                print(f"Warning: gofmt formatting failed for {self.source_code_path}: {result.stderr}")
                
        except subprocess.TimeoutExpired:
            print(f"Warning: gofmt formatting timed out for {self.source_code_path}")
        except FileNotFoundError:
            print("Warning: gofmt command not found. Please ensure Go is installed and gofmt is in PATH.")
        except Exception as e:
            print(f"Warning: Unexpected error during gofmt formatting: {str(e)}")

    def get_system_prompt(self) -> str:
        """
        Get the Go-specific system prompt for AI-assisted code processing.
        
        Returns:
            A string containing the Go-specific system prompt with context about
            Go language features, best practices, and processing instructions
        """
        return """You are an expert Go developer assistant. When processing Go code snippets, adhere to the following guidelines:

**LANGUAGE SPECIFICS**
- Use proper Go package declarations and import statements.
- Define functions with the `func` keyword.
- Implement methods using receiver functions on types.
- Define structs and interfaces with the `type` keyword (`type Name struct {}` and `type Name interface {}`).
- Declare variables using `var` or the short declaration `:=`.
- Ensure all code follows Go's strict formatting conventions (as enforced by `gofmt`).
- **Always run `gofmt` to format the code before returning it.**

**CODE STRUCTURE**
- Maintain correct Go package structure and naming conventions.
- Preserve and properly organize import statements.
- Follow Go naming conventions: capitalize exported identifiers, use lowercase for unexported ones.
- Implement explicit error handling with error returns, following Go idioms.
- Keep function signatures idiomatic and consistent with Go best practices.

**BEST PRACTICES**
- Use Go's built-in error handling patterns (e.g., `if err != nil { ... }`).
- Prefer communication via channels over shared memory, following Go's concurrency principles.
- Prioritize simplicity and readability over unnecessary complexity.
- Leverage Go's type system effectively.
- Adhere to standard Go project layout and conventions.

**WHEN MODIFYING CODE, ENSURE:**
1. The code is valid Go syntax.
2. Imports are accurate and properly managed.
3. Function signatures conform to Go conventions.
4. Error handling follows Go best practices.
5. Formatting matches `gofmt` standards.
6. Comments use Go documentation conventions (`//` for line comments, `/* */` for block comments).

**INSTRUCTION PROCESSING RULE**
- When provided with an `INSTRUCTION` XML tag and a `CURRENTCODE` XML tag, strictly follow the instruction in the `INSTRUCTION` tag and modify the code in the `CURRENTCODE` tag accordingly.
- Do not abort the process prematurely; ensure the instruction is fully applied.
- **Return only the modified Go code as plain text, with no explanations, comments, or extra text, and do not wrap the code in triple backticks or any other formatting.**

**GENERAL INSTRUCTIONS**
- Maintain the existing code style and patterns when making modifications.
- Focus on clarity, correctness, and idiomatic Go usage in all changes."""


