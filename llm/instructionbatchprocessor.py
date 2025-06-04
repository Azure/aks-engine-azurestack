"""
Instruction batch processor module for processing multiple instructions.

This module provides the InstructionBatchProcessor class that handles
batch processing of instruction files with proper error handling and reporting.
"""

import os
import re
import html
from typing import Dict, List, Optional
from urllib.request import urlopen
from urllib.error import URLError

from .instruction import InstructionLoader, Instruction
from .instructionprocessor import InstructionProcessor
from .snippetprocessor import SnippetProcessorFactory, UnsupportedFileTypeError, SnippetProcessor


class InstructionBatchProcessor:
    """
    Processes multiple instruction files in batch with comprehensive error handling.
    
    This class encapsulates the logic for loading, validating, and processing
    multiple instruction files from a directory, providing detailed progress
    reporting and error handling for each instruction.
    """
    
    def __init__(self, code_root_path: str, instruction_path: str, verbose: bool = False, k8s_version: Optional[str] = None) -> None:
        """
        Initialize the batch processor with required paths.
        
        Args:
            code_root_path: The root path of the source code directory
            instruction_path: The directory path containing instruction files
            verbose: Whether to enable verbose output
            k8s_version: The Kubernetes version to use for processing (optional)
            
        Raises:
            FileNotFoundError: If code_root_path or instruction_path doesn't exist
            NotADirectoryError: If the paths are not directories
        """
        self._code_root_path = self._validate_and_resolve_path(code_root_path, "Code root")
        self._instruction_path = self._validate_and_resolve_path(instruction_path, "Instruction")
        self._verbose = verbose
        self._k8s_version = k8s_version
        self._processed_count = 0
        self._failed_count = 0
        self._k8s_component_data: Optional[Dict[str, str]] = None  # Late binding field
    
    @property
    def code_root_path(self) -> str:
        """Get the absolute code root path."""
        return self._code_root_path
    
    @property
    def instruction_path(self) -> str:
        """Get the absolute instruction path."""
        return self._instruction_path
    
    @property
    def verbose(self) -> bool:
        """Get the verbose setting."""
        return self._verbose
    
    @property
    def k8s_version(self) -> Optional[str]:
        """Get the Kubernetes version."""
        return self._k8s_version
    
    @property
    def processed_count(self) -> int:
        """Get the number of successfully processed instructions."""
        return self._processed_count
    
    @property
    def failed_count(self) -> int:
        """Get the number of failed instruction processing attempts."""
        return self._failed_count
    
    def process_instructions(self) -> None:
        """
        Process all instruction files from the instruction directory.
        
        This method loads all instructions, validates them, and processes each one
        with comprehensive error handling and progress reporting.
        """
        # Load instructions from the directory
        instructions = self._load_instructions()
        
        if not instructions:
            print(f"No instruction files found in: {self._instruction_path}")
            return
        
        self._print_header(instructions)
        
        # Process each instruction
        for i, instruction in enumerate(instructions, 1):
            self._process_single_instruction(instruction, i, len(instructions))
        
        self._print_summary()
    
    def _validate_and_resolve_path(self, path: str, path_type: str) -> str:
        """
        Validate and resolve a path to its absolute form.
        
        Args:
            path: The path to validate
            path_type: Description of the path type for error messages
            
        Returns:
            The absolute path
            
        Raises:
            FileNotFoundError: If the path doesn't exist
            NotADirectoryError: If the path is not a directory
        """
        absolute_path = os.path.abspath(path)
        
        if not os.path.exists(absolute_path):
            raise FileNotFoundError(f"{path_type} path not found: {absolute_path}")
        
        if not os.path.isdir(absolute_path):
            raise NotADirectoryError(f"{path_type} path is not a directory: {absolute_path}")
        
        return absolute_path
    
    def _load_instructions(self) -> List[Instruction]:
        """
        Load all instruction files from the instruction directory.
        
        Returns:
            List of loaded instructions
        """
        instruction_loader = InstructionLoader(self._instruction_path)
        return instruction_loader.load_instructions()
    
    def _print_header(self, instructions: List[Instruction]) -> None:
        """
        Print the processing header with summary information.
        
        Args:
            instructions: List of instructions to be processed
        """
        print(f"Found {len(instructions)} instruction(s) to process...")
        print(f"Code root path: {self._code_root_path}")
        print(f"Instruction path: {self._instruction_path}")
        print("-" * 60)
    
    def _process_single_instruction(self, instruction: Instruction, current: int, total: int) -> None:
        """
        Process a single instruction with comprehensive error handling.
        
        Args:
            instruction: The instruction to process
            current: Current instruction number (1-based)
            total: Total number of instructions
        """
        try:
            print(f"\nProcessing instruction {current}/{total}: {instruction.file_path}")
            
            # Create InstructionProcessor with code root path
            processor = InstructionProcessor(instruction, self._code_root_path)
            
            # Display basic information about the instruction
            self._print_instruction_info(processor)
            
            # Process snippet filter and snippet processing
            self._process_snippet_operations(processor)
            
            self._processed_count += 1
            
        except Exception as e:
            print(f"  ✗ Error processing instruction: {e}")
            self._failed_count += 1
    
    def _print_instruction_info(self, processor: InstructionProcessor) -> None:
        """
        Print basic information about the instruction being processed.
        
        Args:
            processor: The instruction processor instance
        """
        print(f"  Filename: {processor.get_filename()}")
        print(f"  Content length: {processor.get_content_length()} characters")
        print(f"  File type: {processor.get_file_type()}")
        
        if self._verbose:
            print(f"  Code root path: {processor.code_root_path}")
    
    def _process_snippet_operations(self, processor: InstructionProcessor) -> None:
        """
        Process snippet filter extraction and snippet processing operations.
        
        Args:
            processor: The instruction processor instance
        """
        try:
            snippet_filter = processor.get_snippet_filter()
            self._print_snippet_filter_info(snippet_filter)
            
            # Create and use snippet processor
            self._process_with_snippet_processor(snippet_filter, processor.instruction)
            
            print("  ✓ Successfully extracted and resolved snippet filter information")
            
        except Exception as e:
            print(f"  ⚠ Warning: Failed to extract snippet filter: {e}")
    
    def _print_snippet_filter_info(self, snippet_filter) -> None:
        """
        Print information about the snippet filter.
        
        Args:
            snippet_filter: The snippet filter object
        """
        print(f"  Snippet filter: {snippet_filter}")
        print(f"  Resolved source code path: {snippet_filter.source_code_path}")
        print(f"  Code elements: {len(snippet_filter.code_elements)} found")
        
        # Display code elements if any and verbose mode is enabled
        if snippet_filter.code_elements and self._verbose:
            print("  Code elements:")
            for element_name, element_desc in snippet_filter.code_elements.items():
                print(f"    - {element_name}: {element_desc}")
    
    def _process_with_snippet_processor(self, snippet_filter, instruction: Instruction) -> None:
        """
        Create and use snippet processor for the given filter.
        
        Args:
            snippet_filter: The snippet filter to process
            instruction: The instruction object being processed
        """
        try:
            snippet_processor = SnippetProcessorFactory.create_processor(snippet_filter)
            modified_code = self._get_modified_code_from_llm(snippet_processor, instruction)
            snippet_processor.apply_snippet_change(modified_code=modified_code)
            
        except UnsupportedFileTypeError as e:
            print(f"  ⚠ Warning: Unsupported file type for snippet processing: {e}")
        except Exception as e:
            print(f"  ⚠ Warning: Failed to create snippet processor: {e}")
    
    def _print_summary(self) -> None:
        """Print the final processing summary."""
        print(f"\n{'-' * 60}")
        print("Instruction processing completed.")
        print(f"Successfully processed: {self._processed_count}")
        if self._failed_count > 0:
            print(f"Failed to process: {self._failed_count}")
    
    def _get_modified_code_from_llm(self, snippet_processor: 'SnippetProcessor', instruction: Instruction) -> str:
        """
        Call LLM to retrieve the modified code based on the snippet and instruction.
        
        Args:
            snippet_processor: The snippet processor instance containing the code to modify
            instruction: The instruction object containing the modification requirements
            
        Returns:
            The modified code as returned by the LLM
            
        Raises:
            ValueError: If LLM response is invalid or empty
            RuntimeError: If LLM API call fails
        """
        try:
            # Import AzureLLMClient here to avoid circular imports
            from .client.llm_client import AzureLLMClient
            
            # Extract the current code snippet
            current_snippet = snippet_processor.extract_snippet()
            
            # Get the system prompt from the snippet processor
            system_prompt = snippet_processor.get_system_prompt()
            
            # Normalize the instruction content by replacing placeholders
            normalized_instruction = self._normalize_instruction(instruction.get_content())
            
            # Create the user prompt with normalized instruction and current code
            user_prompt = f"""
Please modify the following code according to the instruction provided.

<INSTRUCTION>
{normalized_instruction}
</INSTRUCTION>

<CURRENTCODE>
```
{current_snippet}
```
<CURRENTCODE>

Please provide only the modified code in your response, without any additional explanation or formatting.
"""
            
            # Initialize LLM client and get completion
            llm_client = AzureLLMClient()
            modified_code = llm_client.get_completion(
                prompt=user_prompt,
                system_message=system_prompt,
                temperature=0.0  # Low temperature for more consistent code generation
            )
            
            if not modified_code or not modified_code.strip():
                raise ValueError("LLM returned empty or invalid modified code")
            
            if self._verbose:
                print(f"  LLM generated {len(modified_code)} characters of modified code")
            
            return modified_code.strip()
            
        except Exception as e:
            raise RuntimeError(f"Failed to get modified code from LLM: {e}") from e
    
    def _get_k8s_component_data(self) -> Dict[str, str]:
        """
        Get Kubernetes component data using late binding pattern.
        
        This method implements lazy loading of Kubernetes component data,
        initializing the data only when first accessed. The data includes
        component versions, configurations, and metadata based on the
        current Kubernetes version.
        
        Returns:
            Dictionary mapping component names to their corresponding data/versions
            
        Raises:
            ValueError: If k8s_version is not set or is invalid
            RuntimeError: If component data cannot be loaded
        """
        # Late binding: initialize data only when first accessed
        if self._k8s_component_data is None:
            self._k8s_component_data = self._initialize_k8s_component_data()
        
        return self._k8s_component_data
    
    def _initialize_k8s_component_data(self) -> Dict[str, str]:
        """
        Initialize Kubernetes component data based on the current version.
        
        Returns:
            Dictionary with component data
            
        Raises:
            ValueError: If k8s_version is not set
            RuntimeError: If component data cannot be determined
        """
        if not self._k8s_version:
            raise ValueError("Kubernetes version is required to initialize component data")
        
        try:
            # Component data mapping based on Kubernetes version
            component_data = {
                "k8s_version": self._k8s_version,
            }
            
            # Add removed feature gates if available
            try:
                removed_feature_gates = self._find_removed_feature_gates()
                component_data["removed_feature_gates"] = removed_feature_gates
                
                if self._verbose and removed_feature_gates:
                    print(f"  Added removed feature gates: {removed_feature_gates}")
                    
            except Exception as e:
                if self._verbose:
                    print(f"  Warning: Could not retrieve removed feature gates: {e}")
                component_data["removed_feature_gates"] = ""
            
            if self._verbose:
                print(f"  Initialized component data for Kubernetes {self._k8s_version}")
                print(f"  Component count: {len(component_data)}")
            
            return component_data
            
        except Exception as e:
            raise RuntimeError(f"Failed to initialize Kubernetes component data: {e}") from e
    
    def _normalize_instruction(self, instruction_content: str) -> str:
        """
        Normalize instruction content by replacing placeholders with actual values.
        
        This method processes instruction content and replaces various placeholders
        with corresponding values from the Kubernetes component data and other
        system information. It supports both explicit placeholders and XML-style tags.
        
        Args:
            instruction_content: The raw instruction content containing placeholders
            
        Returns:
            The normalized instruction content with placeholders replaced
            
        Raises:
            ValueError: If required component data is not available
            RuntimeError: If placeholder replacement fails
        """
        if not instruction_content or not instruction_content.strip():
            return instruction_content
        
        try:
            normalized_content = instruction_content
            
            # Get component data for placeholder replacement
            if self._k8s_version:
                component_data = self._get_k8s_component_data()
                
                # Replace standard placeholders with component data
                for key, value in component_data.items():
                    placeholder_patterns = [
                        f"__{key}__",  # Standard placeholder format: __key__
                        f"{{{{{key}}}}}",  # Braces format: {{key}}
                        f"${{{key}}}",  # Shell-style format: ${key}
                    ]
                    
                    for pattern in placeholder_patterns:
                        normalized_content = normalized_content.replace(pattern, value)
            
            if self._verbose and normalized_content != instruction_content:
                original_length = len(instruction_content)
                normalized_length = len(normalized_content)
                print(f"  Instruction normalized: {original_length} → {normalized_length} characters")
            
            return normalized_content
            
        except Exception as e:
            raise RuntimeError(f"Failed to normalize instruction content: {e}") from e
    
    def _replace_system_placeholders(self, content: str) -> str:
        """
        Replace system-level placeholders in the content.
        
        Args:
            content: The content to process
            
        Returns:
            Content with system placeholders replaced
        """
        system_replacements = {
            "__code_root_path__": self._code_root_path,
            "__instruction_path__": self._instruction_path,
            "{{code_root_path}}": self._code_root_path,
            "{{instruction_path}}": self._instruction_path,
            "${code_root_path}": self._code_root_path,
            "${instruction_path}": self._instruction_path,
        }
        
        normalized_content = content
        for placeholder, value in system_replacements.items():
            normalized_content = normalized_content.replace(placeholder, value)
        
        return normalized_content
    
    def _get_previous_minor_version(self, k8s_version: str) -> str:
        """
        Extract the previous minor version from a Kubernetes version string.
        
        Args:
            k8s_version: Kubernetes version string (e.g., "1.31.8")
            
        Returns:
            Previous minor version string (e.g., "1.30")
            
        Raises:
            ValueError: If the version format is invalid or minor version is 0
        """
        # Match version pattern like "1.31.8" or "v1.31.8"
        version_pattern = r'^v?(\d+)\.(\d+)(?:\.(\d+))?'
        match = re.match(version_pattern, k8s_version.strip())
        
        if not match:
            raise ValueError(f"Invalid Kubernetes version format: {k8s_version}")
        
        major = int(match.group(1))
        minor = int(match.group(2))
        
        if minor == 0:
            raise ValueError(f"Cannot get previous minor version for {k8s_version} (minor version is 0)")
        
        return f"{major}.{minor - 1}"
    
    def _fetch_html_content(self, url: str) -> str:
        """
        Fetch HTML content from the given URL using urllib.
        
        Args:
            url: URL to fetch content from
            
        Returns:
            HTML content as string
            
        Raises:
            URLError: If the request fails
        """
        try:
            with urlopen(url, timeout=30) as response:
                content = response.read()
                # Decode content, handling potential encoding issues
                if isinstance(content, bytes):
                    try:
                        return content.decode('utf-8')
                    except UnicodeDecodeError:
                        return content.decode('utf-8', errors='ignore')
                return content
        except URLError as e:
            raise URLError(f"Failed to fetch content from {url}: {e}")
    
    def _extract_table_rows(self, html_content: str) -> List[List[str]]:
        """
        Extract table rows from HTML content using regex patterns.
        
        Args:
            html_content: HTML content containing the feature gates table
            
        Returns:
            List of table rows, where each row is a list of cell contents
        """
        # Find the table with "Feature Gates Removed" caption
        # Look for caption with the text, case-insensitive
        caption_pattern = r'<caption[^>]*>.*?Feature\s+Gates\s+Removed.*?</caption>'
        caption_match = re.search(caption_pattern, html_content, re.IGNORECASE | re.DOTALL)
        
        if not caption_match:
            raise ValueError("Feature Gates Removed table not found")
        
        # Find the table that contains this caption
        # Look for the table tag that comes before or after the caption
        table_start = html_content.rfind('<table', 0, caption_match.start())
        if table_start == -1:
            # Caption might be inside the table, look for table after caption
            table_start = html_content.find('<table', caption_match.end())
            if table_start == -1:
                raise ValueError("Table containing Feature Gates Removed caption not found")
        
        # Find the end of this table
        table_end = html_content.find('</table>', table_start)
        if table_end == -1:
            raise ValueError("Table end tag not found")
        
        table_html = html_content[table_start:table_end + 8]  # +8 for </table>
        
        # Extract all table rows
        rows = []
        row_pattern = r'<tr[^>]*>(.*?)</tr>'
        row_matches = re.findall(row_pattern, table_html, re.DOTALL | re.IGNORECASE)
        
        for row_content in row_matches:
            # Extract cell content from each row
            cell_pattern = r'<(?:td|th)[^>]*>(.*?)</(?:td|th)>'
            cell_matches = re.findall(cell_pattern, row_content, re.DOTALL | re.IGNORECASE)
            
            # Clean up cell content - remove HTML tags and decode entities
            clean_cells = []
            for cell in cell_matches:
                # Remove HTML tags
                clean_cell = re.sub(r'<[^>]+>', '', cell)
                # Decode HTML entities
                clean_cell = html.unescape(clean_cell)
                # Clean up whitespace
                clean_cell = re.sub(r'\s+', ' ', clean_cell).strip()
                clean_cells.append(clean_cell)
            
            if clean_cells:  # Only add non-empty rows
                rows.append(clean_cells)
        
        return rows
    
    def _parse_removed_features_table(self, html_content: str, target_version: str) -> List[str]:
        """
        Parse the HTML content to extract removed feature gates for the target version.
        
        Args:
            html_content: HTML content containing the feature gates table
            target_version: Target Kubernetes version to look for
            
        Returns:
            List of removed feature gate names
        """
        rows = self._extract_table_rows(html_content)
        
        if not rows:
            raise ValueError("No table rows found")
        
        # The first row should be the header
        headers = rows[0]
        
        # Find column indices
        feature_col_idx = None
        to_col_idx = None
        stage_col_idx = None
        
        for idx, header in enumerate(headers):
            header_lower = header.lower()
            if 'feature' in header_lower or 'gate' in header_lower:
                feature_col_idx = idx
            elif 'to' in header_lower:
                to_col_idx = idx
            elif 'stage' in header_lower:
                stage_col_idx = idx
        
        if feature_col_idx is None or to_col_idx is None or stage_col_idx is None:
            raise ValueError(f"Required columns not found. Headers: {headers}")
        
        # Extract removed features
        removed_features = []
        
        for row in rows[1:]:  # Skip header row
            if len(row) <= max(feature_col_idx, to_col_idx, stage_col_idx):
                continue
            
            feature_name = row[feature_col_idx].strip()
            to_version = row[to_col_idx].strip()
            stage = row[stage_col_idx].strip()
            
            # Check if this row matches our criteria
            if (to_version == target_version and 
                stage.lower() in ['deprecated', 'ga']):
                removed_features.append(feature_name)
        
        return removed_features
    
    def _find_removed_feature_gates(self) -> str:
        """
        Find removed feature gates for the previous Kubernetes minor release.
        
        Uses the instance's k8s_version to determine removed feature gates.
        
        Returns:
            Comma-separated string of removed feature gates with =true suffix
            
        Raises:
            ValueError: If k8s_version is not set or format is invalid
            URLError: If fetching web content fails
        """
        if not self._k8s_version:
            raise ValueError("Kubernetes version is required but not set")
            
        feature_gates_url = "https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates-removed/"
        
        # Get previous minor version
        previous_version = self._get_previous_minor_version(self._k8s_version)
        
        if self._verbose:
            print(f"  Looking for removed feature gates for version: {previous_version}")
        
        # Fetch HTML content
        html_content = self._fetch_html_content(feature_gates_url)
        
        # Parse and extract removed features
        removed_features = self._parse_removed_features_table(html_content, previous_version)
        
        # Format as comma-separated string with =true suffix
        if not removed_features:
            if self._verbose:
                print(f"  No removed feature gates found for version: {previous_version}")
            return ""
        
        result = ",".join(f"{feature}=true" for feature in removed_features)
        
        if self._verbose:
            print(f"  Found {len(removed_features)} removed feature gates: {result}")
        
        return result