#!/usr/bin/env python3
"""
Go Code Modification Assistant using Azure LLM

This module provides functionality to modify Go source code using Azure LLM
based on provided instructions while maintaining Go best practices and standards.
"""

import sys
import os
import re
from typing import Optional, Tuple, Dict
import argparse
import logging

# Add the current directory to Python path to handle relative imports
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from client.llm_client import AzureLLMClient


class GoCodeModificationError(Exception):
    """Custom exception for Go code modification errors."""
    pass


class GoFunctionParser:
    """Utility class for parsing and extracting Go functions from source code."""
    
    @staticmethod
    def extract_function(source_code: str, function_name: str) -> Tuple[str, int, int]:
        """
        Extract a specific function from Go source code.
        
        Args:
            source_code: Complete Go source code
            function_name: Name of the function to extract
            
        Returns:
            Tuple of (function_code, start_line, end_line)
            
        Raises:
            ValueError: If function is not found or parsing fails
        """
        lines = source_code.split('\n')
        
        # Pattern to match function declarations
        # Handles: func name(...) ..., func (receiver) name(...) ..., func name[T](...) ...
        func_pattern = re.compile(
            rf'^\s*func\s+(?:\([^)]*\)\s+)?{re.escape(function_name)}\s*(?:\[[^\]]*\])?\s*\(',
            re.MULTILINE
        )
        
        match = func_pattern.search(source_code)
        if not match:
            raise ValueError(f"Function '{function_name}' not found in source code")
        
        # Find the line number where the function starts
        start_pos = match.start()
        start_line = source_code[:start_pos].count('\n')
        
        # Find the opening brace of the function
        brace_pos = source_code.find('{', start_pos)
        if brace_pos == -1:
            raise ValueError(f"Could not find opening brace for function '{function_name}'")
        
        # Count braces to find the end of the function
        brace_count = 0
        pos = brace_pos
        in_string = False
        in_char = False
        in_comment = False
        in_block_comment = False
        escape_next = False
        
        while pos < len(source_code):
            char = source_code[pos]
            
            # Handle escape sequences
            if escape_next:
                escape_next = False
                pos += 1
                continue
            
            if char == '\\' and (in_string or in_char):
                escape_next = True
                pos += 1
                continue
            
            # Handle comments
            if not in_string and not in_char:
                if pos < len(source_code) - 1:
                    if source_code[pos:pos+2] == '//':
                        in_comment = True
                        pos += 2
                        continue
                    elif source_code[pos:pos+2] == '/*':
                        in_block_comment = True
                        pos += 2
                        continue
                
                if in_block_comment and pos < len(source_code) - 1 and source_code[pos:pos+2] == '*/':
                    in_block_comment = False
                    pos += 2
                    continue
            
            if in_comment and char == '\n':
                in_comment = False
                pos += 1
                continue
            
            if in_comment or in_block_comment:
                pos += 1
                continue
            
            # Handle strings and characters
            if char == '"' and not in_char:
                in_string = not in_string
            elif char == '\'' and not in_string:
                in_char = not in_char
            
            # Count braces only when not in strings or comments
            if not in_string and not in_char:
                if char == '{':
                    brace_count += 1
                elif char == '}':
                    brace_count -= 1
                    if brace_count == 0:
                        # Found the end of the function
                        end_pos = pos + 1
                        end_line = source_code[:end_pos].count('\n')
                        function_code = source_code[start_pos:end_pos]
                        return function_code, start_line, end_line
            
            pos += 1
        
        raise ValueError(f"Could not find complete function body for '{function_name}'")
    
    @staticmethod
    def replace_function(source_code: str, function_name: str, new_function_code: str) -> str:
        """
        Replace a specific function in Go source code with new implementation.
        
        Args:
            source_code: Complete Go source code
            function_name: Name of the function to replace
            new_function_code: New function implementation
            
        Returns:
            Modified source code with replaced function
            
        Raises:
            ValueError: If function is not found
        """
        try:
            old_function_code, start_line, end_line = GoFunctionParser.extract_function(
                source_code, function_name
            )
            
            lines = source_code.split('\n')
            
            # Replace the function lines (end_line is inclusive, so we need end_line + 1)
            new_function_lines = new_function_code.rstrip('\n').split('\n')
            new_lines = lines[:start_line] + new_function_lines + lines[end_line + 1:]
            
            return '\n'.join(new_lines)
            
        except ValueError as e:
            raise ValueError(f"Failed to replace function '{function_name}': {str(e)}") from e


class GoCodeModificationAssistant:
    """Main assistant class that uses Azure LLM to modify Go code based on instructions."""
    
    def __init__(self):
        self._logger = logging.getLogger(__name__)
        self._llm_client = AzureLLMClient()
    
    def _create_system_prompt(self) -> str:
        """Create the system prompt for the Go code modification task."""
        return """You are an expert Go (Golang) code modification assistant. Your task is to apply modifications to existing Go source code based on provided instructions.

## CORE REQUIREMENTS

**Input**: You will receive:
1. Original Go source code
2. Modification instructions in markdown format

**Output**: Provide ONLY the complete, modified Go source code that can be directly written to a .go file.

## GO CODE STANDARDS

- Follow official Go formatting (gofmt-compliant)
- Use proper Go naming conventions (camelCase/PascalCase)
- Implement proper error handling with Go idioms
- Organize imports correctly (standard, third-party, local)
- Remove unused imports and variables
- Ensure thread safety for concurrent code
- Follow Go best practices for performance and memory management

## MODIFICATION RULES

1. **Preserve Functionality**: Keep existing behavior unless explicitly instructed to change
2. **Maintain Compatibility**: Preserve public APIs unless specifically modified
3. **Code Quality**: Apply Go best practices during modifications
4. **Syntax Correctness**: Ensure the output compiles without errors
5. **Complete Output**: Provide the entire file content, not just changes

## CRITICAL CONSTRAINTS - FOLLOW STRICTLY

**MANDATORY RESPONSE FORMAT**: 
- Your response MUST contain ONLY the complete, modified Go source code
- DO NOT include any explanations, descriptions, or commentary
- DO NOT add markdown code blocks or formatting
- DO NOT provide partial code or summaries of changes
- DO NOT include any text before or after the Go code
- The response must be the exact file content that can be directly saved as a .go file

**STRICT COMPLIANCE REQUIREMENTS**:
- Output ONLY valid, compilable Go code
- Preserve existing comments unless they conflict with modifications
- Maintain proper package declarations and file structure
- Ensure all braces, parentheses, and syntax are correct
- Follow the instructions exactly as specified

If instructions are unclear or impossible to implement safely, preserve the original code structure while making minimal necessary changes.

**REMINDER**: Your entire response must be nothing but the complete Go source code. Any deviation from this requirement is unacceptable."""
    
    def _create_user_message(self, original_code: str, instructions_markdown: str) -> str:
        """Create the user message containing the original code and instructions."""
        template = """Original Go source code:
```go
{original_code}
```

Modification instructions:
{instructions_markdown}

Please provide the complete modified Go source code following the system requirements."""
        
        return template.format(
            original_code=original_code,
            instructions_markdown=instructions_markdown if instructions_markdown else "No specific instructions provided. Please format and optimize the code according to Go best practices."
        )
    
    def _create_function_system_prompt(self) -> str:
        """Create the system prompt for Go function modification task."""
        return """You are an expert Go (Golang) function modification assistant. Your task is to apply modifications to a specific Go function based on provided instructions.

## CORE REQUIREMENTS

**Input**: You will receive:
1. A specific Go function source code
2. Modification instructions in markdown format

**Output**: Provide ONLY the complete, modified Go function code.

## GO CODE STANDARDS

- Follow official Go formatting (gofmt-compliant)
- Use proper Go naming conventions (camelCase/PascalCase)
- Implement proper error handling with Go idioms
- Ensure thread safety for concurrent code
- Follow Go best practices for performance and memory management

## MODIFICATION RULES

1. **Preserve Function Signature**: Keep the function name, parameters, and return types unless explicitly instructed to change
2. **Maintain Functionality**: Keep existing behavior unless explicitly instructed to change
3. **Code Quality**: Apply Go best practices during modifications
4. **Syntax Correctness**: Ensure the output compiles without errors
5. **Complete Function**: Provide the entire function implementation, not just changes

## CRITICAL CONSTRAINTS - FOLLOW STRICTLY

**MANDATORY RESPONSE FORMAT**: 
- Your response MUST contain ONLY the complete, modified Go function code
- DO NOT include any explanations, descriptions, or commentary
- DO NOT add markdown code blocks or formatting
- DO NOT provide partial code or summaries of changes
- DO NOT include any text before or after the Go function code
- The response must be the exact function content that can be directly used in a .go file

**STRICT COMPLIANCE REQUIREMENTS**:
- Output ONLY valid, compilable Go function code
- Preserve existing comments unless they conflict with modifications
- Maintain proper function structure (func keyword, signature, body)
- Ensure all braces, parentheses, and syntax are correct
- Follow the instructions exactly as specified

If instructions are unclear or impossible to implement safely, preserve the original function structure while making minimal necessary changes.

**REMINDER**: Your entire response must be nothing but the complete Go function code. Any deviation from this requirement is unacceptable."""
    
    def _create_function_user_message(self, function_code: str, instructions_markdown: str) -> str:
        """Create the user message containing the function code and instructions."""
        template = """Original Go function code:
```go
{function_code}
```

Modification instructions:
{instructions_markdown}

Please provide the complete modified Go function code following the system requirements."""
        
        return template.format(
            function_code=function_code,
            instructions_markdown=instructions_markdown if instructions_markdown else "No specific instructions provided. Please format and optimize the function according to Go best practices."
        )
    
    def modify_code(self, original_code: str, instructions_markdown: str) -> str:
        """
        Modify Go code using Azure LLM based on markdown instructions.
        
        Args:
            original_code: Original Go source code
            instructions_markdown: Modification instructions in markdown format
            
        Returns:
            Modified Go source code from the LLM
            
        Raises:
            GoCodeModificationError: If modification fails
        """
        try:
            self._logger.info("Preparing to modify Go code using Azure LLM")
            
            # Create system prompt and user message
            system_prompt = self._create_system_prompt()
            user_message = self._create_user_message(original_code, instructions_markdown)
            
            self._logger.info(f"Sending request to LLM with {len(original_code)} characters of Go code")
            if instructions_markdown:
                self._logger.info(f"Instructions length: {len(instructions_markdown)} characters")
            
            # Call Azure LLM
            modified_code = self._llm_client.get_completion(
                prompt=user_message,
                system_message=system_prompt,
                temperature=0.0,  # Use deterministic output for code generation
                top_p=1.0
            )
            
            self._logger.info(f"Received modified code with {len(modified_code)} characters")
            
            return modified_code
            
        except Exception as e:
            self._logger.error(f"Failed to modify Go code using LLM: {str(e)}")
            raise GoCodeModificationError(f"LLM modification failed: {str(e)}") from e
        
    def modify_function(self, source_code: str, function_name: str, instructions_markdown: str) -> str:
        """
        Modify a specific function in Go code using Azure LLM.
        
        Args:
            source_code: Complete Go source code
            function_name: Name of the function to modify
            instructions_markdown: Modification instructions in markdown format
            
        Returns:
            Complete source code with the modified function
            
        Raises:
            GoCodeModificationError: If modification fails
        """
        try:
            self._logger.info(f"Preparing to modify function '{function_name}' using Azure LLM")
            
            # Extract the specific function
            function_parser = GoFunctionParser()
            function_code, start_line, end_line = function_parser.extract_function(
                source_code, function_name
            )
            
            self._logger.info(f"Extracted function '{function_name}' from lines {start_line} to {end_line}")
            self._logger.info(f"Function code length: {len(function_code)} characters")
            
            # Create system prompt and user message for function modification
            system_prompt = self._create_function_system_prompt()
            user_message = self._create_function_user_message(function_code, instructions_markdown)
            
            if instructions_markdown:
                self._logger.info(f"Instructions length: {len(instructions_markdown)} characters")
            
            # Call Azure LLM to modify the function
            modified_function_code = self._llm_client.get_completion(
                prompt=user_message,
                system_message=system_prompt,
                temperature=0.0,  # Use deterministic output for code generation
                top_p=1.0
            )
            
            self._logger.info(f"Received modified function with {len(modified_function_code)} characters")
            
            # Replace the function in the original source code
            modified_source_code = GoFunctionParser.replace_function(
                source_code, function_name, modified_function_code.strip()
            )
            
            self._logger.info("Successfully replaced function in source code")
            
            return modified_source_code
            
        except ValueError as e:
            self._logger.error(f"Function parsing error: {str(e)}")
            raise GoCodeModificationError(f"Function parsing failed: {str(e)}") from e
        except Exception as e:
            self._logger.error(f"Failed to modify function '{function_name}' using LLM: {str(e)}")
            raise GoCodeModificationError(f"Function modification failed: {str(e)}") from e
    
    def extract_and_modify_function_only(self, source_code: str, function_name: str, instructions_markdown: str) -> str:
        """
        Extract and modify a specific function, returning only the modified function code.
        
        Args:
            source_code: Complete Go source code
            function_name: Name of the function to modify
            instructions_markdown: Modification instructions in markdown format
            
        Returns:
            Only the modified function code (not the complete source file)
            
        Raises:
            GoCodeModificationError: If modification fails
        """
        try:
            self._logger.info(f"Extracting and modifying function '{function_name}' using Azure LLM")
            
            # Extract the specific function
            function_parser = GoFunctionParser()
            function_code, start_line, end_line = function_parser.extract_function(
                source_code, function_name
            )
            
            self._logger.info(f"Extracted function '{function_name}' from lines {start_line} to {end_line}")
            self._logger.info(f"Function code length: {len(function_code)} characters")
            
            # Create system prompt and user message for function modification
            system_prompt = self._create_function_system_prompt()
            user_message = self._create_function_user_message(function_code, instructions_markdown)
            
            if instructions_markdown:
                self._logger.info(f"Instructions length: {len(instructions_markdown)} characters")
            
            # Call Azure LLM to modify the function
            modified_function_code = self._llm_client.get_completion(
                prompt=user_message,
                system_message=system_prompt,
                temperature=0.0,  # Use deterministic output for code generation
                top_p=1.0
            )
            
            self._logger.info(f"Received modified function with {len(modified_function_code)} characters")
            
            return modified_function_code.strip()
            
        except ValueError as e:
            self._logger.error(f"Function parsing error: {str(e)}")
            raise GoCodeModificationError(f"Function parsing failed: {str(e)}") from e
        except Exception as e:
            self._logger.error(f"Failed to modify function '{function_name}' using LLM: {str(e)}")
            raise GoCodeModificationError(f"Function modification failed: {str(e)}") from e


def setup_logging(level: str = "INFO") -> None:
    """Set up logging configuration."""
    logging.basicConfig(
        level=getattr(logging, level.upper()),
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )


def main() -> int:
    """Main entry point for the Go code modification assistant."""
    parser = argparse.ArgumentParser(
        description="Go Code Modification Assistant using Azure LLM",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Modify entire file
  python add-k8s-version.py --code-file input.go --instructions-file instructions.md --output output.go
  python add-k8s-version.py --code-file input.go --output output.go
  python add-k8s-version.py --code "package main..." --instructions "Add error handling"
  
  # Modify specific function (replace in full file)
  python add-k8s-version.py --code-file input.go --function "myFunction" --instructions "Add error handling" --output output.go
  python add-k8s-version.py --code-file input.go --function "myFunction" --instructions-file instructions.md
  
  # Modify specific function (output only the function)
  python add-k8s-version.py --code-file input.go --function "myFunction" --function-only --instructions "Add error handling"
  python add-k8s-version.py --code-file input.go --function "myFunction" --function-only --instructions-file instructions.md --output modified_function.go
        """
    )
    
    parser.add_argument(
        "--code",
        type=str,
        help="Go source code as string"
    )
    
    parser.add_argument(
        "--code-file",
        type=str,
        help="Path to Go source code file"
    )
    
    parser.add_argument(
        "--instructions",
        type=str,
        help="Modification instructions as markdown string"
    )
    
    parser.add_argument(
        "--instructions-file",
        type=str,
        help="Path to markdown file containing modification instructions"
    )
    
    parser.add_argument(
        "--function",
        type=str,
        help="Name of specific function to modify (if not provided, modifies entire file)"
    )
    
    parser.add_argument(
        "--function-only",
        action="store_true",
        help="When modifying a function, output only the modified function code (not the entire file)"
    )
    
    parser.add_argument(
        "--output",
        type=str,
        help="Output file path (default: stdout)"
    )
    
    parser.add_argument(
        "--log-level",
        choices=["DEBUG", "INFO", "WARNING", "ERROR"],
        default="INFO",
        help="Set logging level"
    )
    
    args = parser.parse_args()
    
    # Setup logging
    setup_logging(args.log_level)
    
    try:
        # Get source code
        if args.code:
            original_code = args.code
        elif args.code_file:
            with open(args.code_file, 'r', encoding='utf-8') as f:
                original_code = f.read()
        else:
            # Read from stdin
            original_code = sys.stdin.read()
        
        # Get instructions (markdown content)
        if args.instructions:
            instructions_markdown = args.instructions
        elif args.instructions_file:
            with open(args.instructions_file, 'r', encoding='utf-8') as f:
                instructions_markdown = f.read()
        else:
            instructions_markdown = ""
        
        # Log the markdown content for reference
        if instructions_markdown:
            logging.info(f"Loaded markdown instructions: {len(instructions_markdown)} characters")
        else:
            logging.info("No specific instructions provided - will use default Go best practices")
        
        # Create assistant and modify code using LLM
        assistant = GoCodeModificationAssistant()
        
        if args.function:
            # Modify specific function
            logging.info(f"Modifying specific function: {args.function}")
            if args.function_only:
                # Return only the modified function
                result_code = assistant.extract_and_modify_function_only(original_code, args.function, instructions_markdown)
            else:
                # Return the complete file with modified function
                result_code = assistant.modify_function(original_code, args.function, instructions_markdown)
        else:
            # Modify entire file
            if args.function_only:
                logging.warning("--function-only flag ignored when not specifying a function")
            result_code = assistant.modify_code(original_code, instructions_markdown)
        
        # Output result
        if args.output:
            with open(args.output, 'w', encoding='utf-8') as f:
                f.write(result_code)
        else:
            print(result_code)
        
        return 0
        
    except FileNotFoundError as e:
        logging.error(f"File not found: {e}")
        return 1
    except GoCodeModificationError as e:
        logging.error(f"Modification error: {e}")
        return 1
    except Exception as e:
        logging.error(f"Unexpected error: {e}")
        return 1


if __name__ == "__main__":
    sys.exit(main())