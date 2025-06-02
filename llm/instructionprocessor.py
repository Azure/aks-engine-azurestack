#!/usr/bin/env python3
"""
Instruction processor module for processing instruction objects.
"""

import json
from typing import Optional, Dict
from .instruction import Instruction
from .snippetprocessor import SnippetFilter, FileType
from .client.llm_client import AzureLLMClient


class InstructionProcessor:
    """
    Processes a single instruction and provides functionality to work with instruction data.
    
    This class handles the processing of an Instruction object, providing methods
    to analyze and transform instruction data.
    """
    
    def __init__(self, instruction: Instruction) -> None:
        """
        Initialize an InstructionProcessor instance.
        
        Args:
            instruction: An Instruction object to process
        """
        self._instruction: Instruction = instruction
        self._azure_llm_client: AzureLLMClient = AzureLLMClient()
    
    @property
    def instruction(self) -> Instruction:
        """
        Get the instruction (readonly property).
        
        Returns:
            The Instruction object
        """
        return self._instruction
    
    def get_content(self) -> str:
        """
        Get the content from the instruction.
        
        Returns:
            The content string from the instruction
        """
        return self._instruction.get_content()
    
    def get_file_path(self) -> str:
        """
        Get the file path of the instruction.
        
        Returns:
            The file path of the instruction
        """
        return self._instruction.file_path
    
    def get_filename(self) -> str:
        """
        Get the filename from the instruction's file path.
        
        Returns:
            The filename extracted from the file path
        """
        return self._instruction.file_path.split('/')[-1]
   
    def get_content_length(self) -> int:
        """
        Get the length of the instruction content.
        
        Returns:
            The number of characters in the instruction content
        """
        return len(self.get_content())
    
    def get_file_type(self) -> Optional['FileType']:
        """
        Get the file type of the instruction based on its file path.
        
        Returns:
            The FileType enum, or None if not recognized
        """
        return FileType.get_filetype(self.get_file_path())
    
    def get_snippet_filter(self) -> SnippetFilter:
        """
        Get a SnippetFilter for this instruction's source code.
        
        Returns:
            A SnippetFilter instance initialized with the instruction's file path
        """
        # Get the system prompt for instruction extraction
        system_prompt = self._get_instruction_extractor_system_prompt()
        
        # Get the instruction content and create user prompt
        instruction_content = self.get_content()
        user_prompt = self._get_instruction_extractor_user_prompt(instruction_content)
        
        # Send request to LLM using Azure LLM client
        llm_response = self._azure_llm_client.get_completion(
            prompt=user_prompt,
            system_message=system_prompt
        )
        
        # Parse the LLM response as JSON
        try:
            response_data = json.loads(llm_response)
            
            # Extract source_code_path from the response
            source_code_path = response_data.get("source_code_path", "")
            
            # Validate that source_code_path is not empty
            if not source_code_path or source_code_path.strip() == "":
                raise ValueError("source_code_path is empty or missing in LLM response")
            
            # Create and return SnippetFilter with the extracted path
            return SnippetFilter(source_code_path)
            
        except json.JSONDecodeError as e:
            raise ValueError(f"Failed to parse LLM response as JSON: {e}")
        except Exception as e:
            raise ValueError(f"Error processing LLM response: {e}")
    
    def _get_instruction_extractor_system_prompt(self) -> str:
        """
        Get the system prompt for instruction extraction.
        
        Returns:
            A system prompt string for parsing markdown and extracting structured data
        """
        system_prompt = """
You are an expert in parsing markdown files and extracting structured data for code analysis.

**Task:**  
Parse only the markdown content found within the XML tag `<Instruction>`. Ignore any content outside this tag. Then extract information from the section with the header `# Code Snippt Filter` (case-insensitive, ignore trailing colons or whitespace).

**Extraction and Formatting Instructions:**

1. **Locate Section:**

   - Find the section starting with the header `# Code Snippt Filter`.
   - Only extract bullet points (`- ...`) directly under this header, stopping at the next header or end of file.

2. **Key Generation and Mapping:**

   - Only process bullet points under the specified header.
   - For each bullet point, split at the first colon (`:`) to separate the label and value.
   - Normalize the label:
     - Convert to lowercase.
     - Replace spaces with underscores.
     - Remove any trailing colons or extra whitespace.
   - If the label is "source code path", map its value (with backticks removed) to the top-level key `"source_code_path"`.
   - All other bullet points are mapped as key-value pairs in the `"code_elements"` object, using the normalized label as the key and the value (with backticks removed) as the value.
   - Ensure all extracted values have surrounding backticks removed.

3. **Output Format:**
   - Return a valid JSON object with the following structure:
     ```json
     {
       "source_code_path": "<value>",
       "code_elements": {
         "<normalized_label>": "<value>",
         ...
       }
     }
     ```
   - If the required section or bullet points are missing, return an empty JSON object: `{}`.

**Example Input:**

```
# Code Snippt Filter:

- source code path: `pkg/api/common/versions.go`
- object name: AllKubernetesSupportedVersions
- begin with: `var AllKubernetesSupportedVersions = map[string]bool`

```

**Example Output:**

```json
{
  "source_code_path": "pkg/api/common/versions.go",
  "code_elements": {
    "object_name": "AllKubernetesSupportedVersions",
    "begin_with": "var AllKubernetesSupportedVersions = map[string]bool"
  }
}
```
"""
        return system_prompt
    
    def _get_instruction_extractor_user_prompt(self, user_prompt: str) -> str:
        """
        Get the user prompt wrapped with XML instruction tags.
        
        Args:
            user_prompt: The user prompt string to wrap
            
        Returns:
            The user prompt wrapped in <Instruction> XML tags
        """
        return f"<Instruction>\n{user_prompt}\n</Instruction>"


