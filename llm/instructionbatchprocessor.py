"""
Instruction batch processor module for processing multiple instructions.

This module provides the InstructionBatchProcessor class that handles
batch processing of instruction files with proper error handling and reporting.
"""

import os
import re
import html
import json
from typing import Dict, List, Optional, Tuple
from urllib.request import urlopen
from urllib.error import URLError

import requests
import yaml

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
            
            # Check if changes should be applied based on input validation
            if self._should_apply_changes(snippet_processor, instruction):
                modified_code = self._get_modified_code_from_llm(snippet_processor, instruction)
                snippet_processor.apply_snippet_change(modified_code=modified_code)
            else:
                print("  ✓ No changes detected by Input Validation, skipping snippet modification")
            
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
            
            # Add CSI driver version information if available
            try:
                csi_version = self._find_latest_csi_driver_version()
                if csi_version:
                    component_data["csi_driver_version"] = csi_version
                    
                    # Get CSI image versions
                    csi_image_versions = self._find_csi_image_versions(csi_version)
                    if csi_image_versions:
                        component_data["csi_image_versions"] = csi_image_versions
                    
                    if self._verbose:
                        print(f"  Added CSI driver version: {csi_version}")
                        if csi_image_versions:
                            print(f"  Added CSI image versions")
                else:      
                    component_data["csi_driver_version"] = ""
                    component_data["csi_image_versions"] = ""
                    
            except Exception as e:
                if self._verbose:
                    print(f"  Warning: Could not retrieve CSI driver information: {e}")
                component_data["csi_driver_version"] = ""
                component_data["csi_image_versions"] = ""
            
            # Add Cloud Provider image versions if available
            try:
                cloud_provider_versions = self._find_cloud_provider_image_versions()
                if cloud_provider_versions:
                    component_data["cloud_provider_image_versions"] = cloud_provider_versions
                    
                    if self._verbose:
                        print(f"  Added Azure Cloud Provider image versions")
                else:
                    component_data["cloud_provider_image_versions"] = ""
                    
            except Exception as e:
                if self._verbose:
                    print(f"  Warning: Could not retrieve Cloud Provider image versions: {e}")
                component_data["cloud_provider_image_versions"] = ""
            
            # Add all supported versions if available
            try:
                supported_versions_list = self._get_current_supported_versions()
                if supported_versions_list:
                    component_data["all_supported_versions"] = ','.join(supported_versions_list)
                    
                    if self._verbose:
                        print(f"  Added all supported versions: {component_data['all_supported_versions']}")
                else:
                    component_data["all_supported_versions"] = ""
                    
            except Exception as e:
                if self._verbose:
                    print(f"  Warning: Could not retrieve supported versions: {e}")
                component_data["all_supported_versions"] = ""
            
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
    
    def _find_latest_csi_driver_version(self, timeout: float = 30.0) -> Optional[str]:
        """
        Find the latest Azure Disk CSI driver version for the current Kubernetes version.
        
        Args:
            timeout: Request timeout in seconds for HTTP requests
            
        Returns:
            The latest CSI driver version string (e.g., 'v1.31.10') or None if not found
            
        Raises:
            requests.RequestException: If there's an error fetching data from GitHub API
        """
        if not self._k8s_version:
            return None
            
        api_url = "https://api.github.com/repos/kubernetes-sigs/azuredisk-csi-driver/contents/charts"
        
        # Extract major and minor version from k8s_version
        normalized_version = self._k8s_version.lstrip('v')
        parts = normalized_version.split('.')
        major, minor = int(parts[0]), int(parts[1])
        
        # Set up session with appropriate headers
        session = requests.Session()
        session.headers.update({
            'Accept': 'application/vnd.github.v3+json',
            'User-Agent': 'InstructionBatchProcessor/1.0'
        })
        
        try:
            # Fetch available versions from GitHub API
            response = session.get(api_url, timeout=timeout)
            response.raise_for_status()
            
            data = response.json()
            if not isinstance(data, list):
                return None
            
            # Parse directory names and filter for valid version format
            versions = []
            version_pattern = r'^v\d+\.\d+\.\d+$'
            
            for item in data:
                if isinstance(item, dict) and item.get('type') == 'dir':
                    name = item.get('name', '')
                    if re.match(version_pattern, name):
                        versions.append(name)
            
            # Filter versions that match the major and minor version
            target_prefix = f"v{major}.{minor}."
            matching_versions = [v for v in versions if v.startswith(target_prefix)]
            
            if not matching_versions:
                return None
            
            # Find the highest version
            def version_key(version: str) -> Tuple[int, int, int]:
                """Extract version tuple for sorting."""
                clean_version = version[1:]  # Remove 'v'
                parts = clean_version.split('.')
                return int(parts[0]), int(parts[1]), int(parts[2])
            
            return max(matching_versions, key=version_key)
            
        except requests.RequestException as e:
            raise requests.RequestException(
                f"Failed to fetch CSI driver versions from GitHub API: {e}"
            ) from e
        finally:
            session.close()
    
    def _find_csi_image_versions(self, csi_version: str, timeout: float = 30.0) -> Optional[str]:
        """
        Find the image versions for CSI components from the values.yaml file.
        
        Args:
            csi_version: CSI driver version (e.g., 'v1.31.10')
            timeout: Request timeout in seconds for HTTP requests
            
        Returns:
            JSON string with image versions or None if not found
            
        Raises:
            requests.RequestException: If there's an error fetching data from GitHub API
        """
        # Construct the URL for the values.yaml file
        values_url = f"https://raw.githubusercontent.com/kubernetes-sigs/azuredisk-csi-driver/master/charts/{csi_version}/azuredisk-csi-driver/values.yaml"
        
        # List of image components we're interested in
        target_components = [
            "oss/kubernetes-csi/csi-provisioner",
            "oss/kubernetes-csi/csi-attacher", 
            "oss/kubernetes-csi/livenessprobe",
            "oss/kubernetes-csi/csi-node-driver-registrar",
            "oss/kubernetes-csi/csi-snapshotter",
            "oss/kubernetes-csi/snapshot-controller",
            "oss/kubernetes-csi/csi-resizer",
            "oss/kubernetes-csi/azuredisk-csi"
        ]
        
        # Set up session
        session = requests.Session()
        session.headers.update({
            'User-Agent': 'InstructionBatchProcessor/1.0'
        })
        
        try:
            # Fetch the values.yaml file
            response = session.get(values_url, timeout=timeout)
            response.raise_for_status()
            
            yaml_content = response.text
            image_versions = {}
            
            # Parse the YAML content
            try:
                yaml_data = yaml.safe_load(yaml_content)
                image_versions = self._extract_all_image_versions(yaml_data, target_components)
            except yaml.YAMLError as e:
                raise ValueError(f"Failed to parse YAML content: {e}")
            
            return json.dumps(image_versions) if image_versions else None
            
        except requests.RequestException as e:
            raise requests.RequestException(
                f"Failed to fetch values.yaml from GitHub: {e}"
            ) from e
        finally:
            session.close()
    
    def _extract_all_image_versions(self, yaml_data: Dict, target_components: List[str]) -> Dict[str, str]:
        """
        Extract image versions for all target components from parsed YAML data.
        
        Args:
            yaml_data: The parsed YAML data as a dictionary
            target_components: List of component names to extract versions for
            
        Returns:
            Dictionary mapping component names to their versions
        """
        image_versions = {}
        
        # Component name mapping from target_components to YAML keys
        component_mapping = {
            "oss/kubernetes-csi/azuredisk-csi": ["image", "azuredisk"],
            "oss/kubernetes-csi/csi-provisioner": ["image", "csiProvisioner"],
            "oss/kubernetes-csi/csi-attacher": ["image", "csiAttacher"],
            "oss/kubernetes-csi/csi-resizer": ["image", "csiResizer"],
            "oss/kubernetes-csi/livenessprobe": ["image", "livenessProbe"],
            "oss/kubernetes-csi/csi-node-driver-registrar": ["image", "nodeDriverRegistrar"],
            "oss/kubernetes-csi/csi-snapshotter": ["snapshot", "image", "csiSnapshotter"],
            "oss/kubernetes-csi/snapshot-controller": ["snapshot", "image", "csiSnapshotController"]
        }
        
        for component in target_components:
            if component in component_mapping:
                path = component_mapping[component]
                version = self._get_nested_value(yaml_data, path + ["tag"])
                if version:
                    # Ensure version starts with 'v'
                    version = version if version.startswith('v') else f'v{version}'
                    image_versions[component] = version
        
        return image_versions
    
    def _get_nested_value(self, data: Dict, path: List[str]) -> Optional[str]:
        """
        Get a nested value from a dictionary using a path list.
        
        Args:
            data: The dictionary to search in
            path: List of keys representing the path to the value
            
        Returns:
            The value if found, None otherwise
        """
        current = data
        for key in path:
            if isinstance(current, dict) and key in current:
                current = current[key]
            else:
                return None
        
        return str(current) if current is not None else None
    
    def _find_cloud_provider_image_versions(self, timeout: float = 30.0) -> Optional[str]:
        """
        Find the latest image versions for Azure Cloud Controller Manager and 
        Azure Cloud Node Manager from GitHub tags.
        
        Args:
            timeout: Request timeout in seconds for HTTP requests
            
        Returns:
            JSON string with image versions or None if not found
            
        Raises:
            requests.RequestException: If there's an error fetching data from GitHub API
        """
        if not self._k8s_version:
            return None
            
        api_url = "https://api.github.com/repos/kubernetes-sigs/cloud-provider-azure/tags"
        
        # Extract major and minor version from k8s_version
        normalized_version = self._k8s_version.lstrip('v')
        parts = normalized_version.split('.')
        major, minor = int(parts[0]), int(parts[1])
        
        # Target components to find versions for
        target_components = [
            "oss/kubernetes/azure-cloud-controller-manager",
            "oss/kubernetes/azure-cloud-node-manager"
        ]
        
        # Set up session with appropriate headers
        session = requests.Session()
        session.headers.update({
            'Accept': 'application/vnd.github.v3+json',
            'User-Agent': 'InstructionBatchProcessor/1.0'
        })
        
        try:
            # Fetch tags from GitHub API
            response = session.get(api_url, timeout=timeout)
            response.raise_for_status()
            
            data = response.json()
            if not isinstance(data, list):
                return None
            
            # Parse tag names and filter for valid version format
            versions = []
            version_pattern = r'^v\d+\.\d+\.\d+$'
            
            for item in data:
                if isinstance(item, dict):
                    name = item.get('name', '')
                    if re.match(version_pattern, name):
                        versions.append(name)
            
            # Filter versions that match the major and minor version
            target_prefix = f"v{major}.{minor}."
            matching_versions = [v for v in versions if v.startswith(target_prefix)]
            
            if not matching_versions:
                return None
            
            # Find the highest version
            def version_key(version: str) -> Tuple[int, int, int]:
                """Extract version tuple for sorting."""
                clean_version = version[1:]  # Remove 'v'
                parts = clean_version.split('.')
                return int(parts[0]), int(parts[1]), int(parts[2])
            
            latest_version = max(matching_versions, key=version_key)
            
            # Create the result dictionary with both components having the same version
            image_versions = {component: latest_version for component in target_components}
            
            return json.dumps(image_versions) if image_versions else None
            
        except requests.RequestException as e:
            raise requests.RequestException(
                f"Failed to fetch cloud provider versions from GitHub API: {e}"
            ) from e
        finally:
            session.close()
    
    def _should_apply_changes(self, snippet_processor: 'SnippetProcessor', instruction: Instruction) -> bool:
        """
        Check if changes should be applied based on input validation section in instruction.
        
        This method evaluates the "# Input Validation" section in the instruction content.
        If the section doesn't exist, changes should be applied. If it exists, the LLM
        evaluates the validation criteria and returns whether changes are needed.
        
        Args:
            snippet_processor: The snippet processor instance containing the code to validate
            instruction: The instruction object containing validation requirements
            
        Returns:
            True if changes should be applied, False otherwise
            
        Raises:
            RuntimeError: If LLM evaluation fails
        """
        try:
            # Get the instruction content
            instruction_content = instruction.get_content()

            # Normalize the instruction content by replacing placeholders
            normalized_instruction = self._normalize_instruction(instruction.get_content())

            # Extract the current code snippet
            current_snippet = snippet_processor.extract_snippet()

            # Check if "# Input Validation" section exists
            if "# Input Validation" not in instruction_content:
                if self._verbose:
                    print("  No Input Validation section found, proceeding with changes")
                return True
            
            # Import AzureLLMClient here to avoid circular imports
            from .client.llm_client import AzureLLMClient
            
            # Create the validation prompt using the full instruction content
            validation_prompt = f"""

<INSTRUCTION>
{normalized_instruction}
</INSTRUCTION>

<CURRENTCODE>
{current_snippet}
</CURRENTCODE>

Process the instruction in xml tag `<INSTRUCTION>`. Find the "# Input Validation" section and execute its logic.
If "# Input Validation" section exists: Execute the validation logic and return its result ("True" or "False").
If "# Input Validation" section does not exist: Return "True".

Return only "True" or "False".
"""
            
            # System message for validation
            system_message = """You are a validation processor. Your task is to find and execute "# Input Validation" logic.

Instructions:
1. Look for "# Input Validation" section in the instruction content
2. If found: Execute the validation logic within that section and return its direct result
3. If not found: Return "True"

The validation logic will directly evaluate to "True" or "False" - return exactly that result.
Respond with only "True" or "False"."""
            
            # Initialize LLM client and get validation result
            llm_client = AzureLLMClient()
            validation_result = llm_client.get_completion(
                prompt=validation_prompt,
                system_message=system_message,
                temperature=0.0  # Low temperature for consistent evaluation
            )
            
            # Parse the validation result
            validation_result = validation_result.strip().lower()
            
            if validation_result == "true":
                print("  Input validation indicates changes are needed")
                return True
            else:
                # For all other responses (including "false" and unexpected responses), return False
                print(f"  Input validation result: '{validation_result}', no changes needed")
                return False
                
        except Exception as e:
            # If validation fails, default to applying changes to be safe
            if self._verbose:
                print(f"  Validation check failed: {e}, proceeding with changes")
            return True
    
    def _get_current_supported_versions(self) -> List[str]:
        """
        Get the currently supported Kubernetes versions from AllKubernetesSupportedVersionsAzureStack.
        
        This method reads the versions.go file and extracts versions where the value is True
        from the AllKubernetesSupportedVersionsAzureStack map.
        
        Returns:
            List of supported Kubernetes version strings
            
        Raises:
            FileNotFoundError: If the versions.go file doesn't exist
            ValueError: If the map cannot be parsed
        """
        versions_file_path = os.path.join(self._code_root_path, "pkg", "api", "common", "versions.go")
        
        if not os.path.exists(versions_file_path):
            raise FileNotFoundError(f"versions.go file not found at: {versions_file_path}")
        
        try:
            with open(versions_file_path, 'r', encoding='utf-8') as file:
                content = file.read()
            
            # Find the AllKubernetesSupportedVersionsAzureStack map
            map_pattern = r'var\s+AllKubernetesSupportedVersionsAzureStack\s*=\s*map\[string\]bool\s*\{([^}]+)\}'
            match = re.search(map_pattern, content, re.DOTALL)
            
            if not match:
                raise ValueError("AllKubernetesSupportedVersionsAzureStack map not found in versions.go")
            
            map_content = match.group(1)
            
            # Parse the map entries
            supported_versions = []
            entry_pattern = r'"([^"]+)":\s*(true|false)'
            entries = re.findall(entry_pattern, map_content)
            
            for version, is_supported in entries:
                if is_supported.lower() == 'true':
                    supported_versions.append(version)
            
            # Sort versions for consistent output
            supported_versions.sort(key=lambda v: [int(x) for x in v.split('.')])
            
            if self._verbose:
                print(f"  Found {len(supported_versions)} supported versions: {supported_versions}")
            
            return supported_versions
            
        except Exception as e:
            raise ValueError(f"Failed to parse AllKubernetesSupportedVersionsAzureStack map: {e}") from e