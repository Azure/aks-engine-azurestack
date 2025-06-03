from pathlib import Path
import re
from typing import Optional, Tuple, NamedTuple


class Position(NamedTuple):
    """Represents a position in a text file with row and column numbers."""
    row: int
    column: int


class VariableDefinition(NamedTuple):
    """Represents a variable definition with its content and positions."""
    content: str
    start_position: Position
    end_position: Position


def load_file_content(file_path: str) -> str:
    """
    Load the content of a file into a string.
    
    Args:
        file_path: Path to the file to load
        
    Returns:
        The file content as a string
        
    Raises:
        FileNotFoundError: If the file doesn't exist
        IOError: If there's an error reading the file
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            return file.read()
    except FileNotFoundError:
        raise FileNotFoundError(f"File not found: {file_path}")
    except IOError as e:
        raise IOError(f"Error reading file {file_path}: {e}")


def get_position_from_offset(content: str, offset: int) -> Position:
    """
    Convert a character offset to row and column position.
    
    Args:
        content: The text content
        offset: Character offset from the beginning of the content
        
    Returns:
        Position with row and column numbers (1-based indexing)
    """
    lines_before = content[:offset].split('\n')
    row = len(lines_before)
    column = len(lines_before[-1]) + 1 if lines_before else 1
    return Position(row=row, column=column)


def extract_variable_definition(content: str, variable_name: str) -> Optional[VariableDefinition]:
    """
    Extract a variable definition from Go source code with position information.
    
    Args:
        content: The source code content
        variable_name: The name of the variable to extract
        
    Returns:
        VariableDefinition containing the content and start/end positions,
        or None if the variable is not found
    """
    # Pattern to match variable declarations like:
    # var VariableName = map[string]bool{
    # var VariableName = []string{
    # var VariableName = SomeType{
    pattern = rf'var\s+{re.escape(variable_name)}\s*=\s*[^{{]*\{{'
    
    # Find the start of the variable definition
    match = re.search(pattern, content)
    if not match:
        return None
    
    start_pos = match.start()
    
    # Find the matching closing brace
    brace_count = 0
    in_variable = False
    end_pos = start_pos
    
    for i, char in enumerate(content[start_pos:], start_pos):
        if char == '{':
            brace_count += 1
            in_variable = True
        elif char == '}':
            brace_count -= 1
            if in_variable and brace_count == 0:
                end_pos = i + 1
                break
    
    if brace_count != 0:
        return None  # Unmatched braces
    
    variable_content = content[start_pos:end_pos]
    start_position = get_position_from_offset(content, start_pos)
    end_position = get_position_from_offset(content, end_pos)
    
    return VariableDefinition(
        content=variable_content,
        start_position=start_position,
        end_position=end_position
    )


def main() -> None:
    """Main function to test loading the Go file content and extracting variable definitions."""
    go_file_path = "/home/honcao/code/akse05/pkg/api/common/versions.go"
    variable_name = "AllKubernetesSupportedVersions"
    
    try:
        # Load the Go file content into a string
        go_file_content = load_file_content(go_file_path)
        
        # Print some basic info about the loaded content
        print(f"Successfully loaded file: {go_file_path}")
        print(f"File size: {len(go_file_content)} characters")
        print(f"Number of lines: {go_file_content.count(chr(10)) + 1}")
        
        # Extract the specific variable definition
        variable_definition = extract_variable_definition(go_file_content, variable_name)
        
        if variable_definition:
            print(f"\n=== Variable Definition: {variable_name} ===")
            print(f"Start position: Row {variable_definition.start_position.row}, Column {variable_definition.start_position.column}")
            print(f"End position: Row {variable_definition.end_position.row}, Column {variable_definition.end_position.column}")
            print(f"Content:\n{variable_definition.content}")
            
            # Print some statistics about the variable
            lines = variable_definition.content.split('\n')
            print(f"\nVariable definition statistics:")
            print(f"- Lines: {len(lines)}")
            print(f"- Characters: {len(variable_definition.content)}")
            
            # Count the number of entries in the map
            entry_pattern = r'"[^"]+"\s*:\s*(true|false)'
            entries = re.findall(entry_pattern, variable_definition.content)
            print(f"- Map entries: {len(entries)}")
            
            # Count true vs false values
            true_count = entries.count('true')
            false_count = entries.count('false')
            print(f"- Enabled versions (true): {true_count}")
            print(f"- Disabled versions (false): {false_count}")
            
        else:
            print(f"Variable '{variable_name}' not found in the file")
            
    except (FileNotFoundError, IOError) as e:
        print(f"Error: {e}")


if __name__ == "__main__":
    main()