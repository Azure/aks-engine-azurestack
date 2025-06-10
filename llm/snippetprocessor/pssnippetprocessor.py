#!/usr/bin/env python3
"""
PowerShell-specific snippet processor for extracting PowerShell code elements.
"""

from typing import Tuple, Optional
from .snippetprocessor import SnippetProcessor
from .location import Location


class PsSnippetProcessor(SnippetProcessor):
    """
    PowerShell-specific implementation of SnippetProcessor.
    
    This class provides functionality to extract PowerShell-specific code elements:
    - Function definitions
    - Cmdlet definitions
    - Script blocks
    - Variable declarations
    - Class definitions
    - Module definitions
    """
    
    def extract_snippet(self) -> str:
        """
        Extract and return a formatted PowerShell code snippet.
        
        Returns:
            The formatted code snippet as a string with preserved spacing and indentation
            
        Raises:
            SnippetNotFoundError: If the snippet is not found in the source code
        """
        from .exceptions import SnippetNotFoundError
        
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
        Get the start and end location of a specific PowerShell code snippet.
        
        Returns:
            A tuple containing (start_location, end_location) if found, None otherwise.
            The locations represent the beginning and end of the code snippet.
        """
        # Check if "begin_with" is in code_elements to extract function
        if "begin_with" in self.code_elements:
            return self._extract_function_location(self.code_elements["begin_with"])
        
        # If no specific extraction method found, extract entire file
        return self._extract_entire_file_location()

    def _extract_function_location(self, begin_with_pattern: str) -> Optional[Tuple[Location, Location]]:
        """
        Extract the location of a PowerShell function that begins with a specific pattern
        and is delimited by curly braces { }.
        
        Args:
            begin_with_pattern: Pattern to match the beginning of the function
            
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
            in_function = False
            
            for line_num, line in enumerate(lines, 1):
                # Check if this line contains the function beginning pattern
                if line.strip().startswith(begin_with_pattern):
                    start_location = Location(line_num, 1)
                    in_function = True
                    # Count opening braces in the declaration line
                    brace_count += line.count('{') - line.count('}')
                    
                    # If function completes on same line, return immediately
                    if brace_count == 0:
                        end_location = Location(line_num, len(line))
                        return (start_location, end_location)
                    continue
                
                if in_function:
                    # Count braces to find the end of the function
                    brace_count += line.count('{') - line.count('}')
                    
                    # Function ends when brace count returns to 0
                    if brace_count == 0 and start_location is not None:
                        end_location = Location(line_num, len(line))
                        return (start_location, end_location)
            
            # Function declaration found but no closing brace (malformed code)
            if start_location is not None and in_function:
                # Return up to the end of file
                end_line = len(lines)
                end_location = Location(end_line, len(lines[-1]) if lines else 1)
                return (start_location, end_location)
            
            return None
            
        except (FileNotFoundError, IOError) as e:
            raise e

    def _extract_entire_file_location(self) -> Optional[Tuple[Location, Location]]:
        """
        Extract the location representing the entire PowerShell source code file.
        
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

    def apply_snippet_change(self, modified_code: str) -> None:
        """
        Apply modified PowerShell code to the source file.
        
        Args:
            modified_code: The modified PowerShell code content to write to the file
            
        Returns:
            None
            
        Raises:
            IOError: If there's an error writing to the file
            FileNotFoundError: If the source file doesn't exist
            PermissionError: If there are permission issues writing to the file
        """
        from .exceptions import SnippetNotFoundError
        
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
                
        except (FileNotFoundError, PermissionError) as e:
            raise e
        except Exception as e:
            raise IOError(f"Error applying snippet change: {str(e)}") from e
   
    def get_system_prompt(self) -> str:
        """
        Get the PowerShell-specific system prompt for AI-assisted code processing.
        
        Returns:
            A string containing the PowerShell-specific system prompt with context about
            PowerShell language features, best practices, and processing instructions
        """
        return """You are an expert PowerShell developer assistant. When processing PowerShell code snippets, adhere to the following guidelines:

**LANGUAGE SPECIFICS**
- Use PowerShell-approved verbs for function and cmdlet names (Get-, Set-, New-, Remove-, etc.).
- Define functions with the `function` keyword or using the `Function:` drive syntax.
- Implement advanced functions with `[CmdletBinding()]` attribute for cmdlet-like behavior.
- Define parameters with proper type annotations: `[Parameter()][string]$Name`.
- Use PowerShell variables with the `$` prefix (e.g., `$MyVariable`).
- Follow PowerShell naming conventions: PascalCase for functions, camelCase for parameters.
- Ensure proper use of PowerShell operators (`-eq`, `-ne`, `-like`, `-match`, etc.).

**CODE STRUCTURE**
- Maintain correct PowerShell module structure and manifest files (.psd1).
- Preserve and properly organize `using` statements and module imports.
- Follow PowerShell naming conventions: Verb-Noun pattern for functions and cmdlets.
- Implement proper parameter validation using PowerShell validation attributes.
- Use pipeline-aware functions with proper `Begin`, `Process`, and `End` blocks when appropriate.
- Structure scripts with proper comment-based help documentation.

**BEST PRACTICES**
- Use PowerShell's built-in error handling with `try-catch-finally` blocks and `$ErrorActionPreference`.
- Leverage the PowerShell pipeline for data processing and transformation.
- Use Write-Verbose, Write-Debug, and Write-Information for proper logging.
- Implement proper parameter sets and parameter validation.
- Use splatting for passing parameters to improve readability.
- Follow PowerShell security best practices (avoid hardcoded credentials, use SecureString).
- Prefer PowerShell native cmdlets over external executables when possible.

**WHEN MODIFYING CODE, ENSURE:**
1. The code follows PowerShell syntax and language rules.
2. All variables are properly scoped and typed.
3. Functions include appropriate parameter validation and help documentation.
4. Error handling follows PowerShell best practices with proper exception types.
5. Code follows PowerShell style guidelines and formatting conventions.
6. Comments use PowerShell documentation standards (`#` for line comments, `<# #>` for block comments).
7. Pipeline usage is optimized and follows PowerShell conventions.

**INSTRUCTION PROCESSING RULE**
- When provided with an `INSTRUCTION` XML tag and a `CURRENTCODE` XML tag, strictly follow the instruction in the `INSTRUCTION` tag and modify the code in the `CURRENTCODE` tag accordingly.
- Do not abort the process prematurely; ensure the instruction is fully applied.
- **When returning code, follow these rules:**
    - return only the code as plain text, with no explanations, comments, or extra formatting.

**GENERAL INSTRUCTIONS**
- Maintain the existing code style and patterns when making modifications.
- Focus on clarity, correctness, and idiomatic PowerShell usage in all changes.
- Ensure compatibility with the target PowerShell version (Windows PowerShell 5.1 or PowerShell 7+)."""
