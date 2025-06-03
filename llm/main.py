#!/usr/bin/env python3
"""
Main module for processing instructions with InstructionBatchProcessor.
"""

import argparse
import os
import re
import sys

# Handle both direct execution and module import
try:
    from .instructionbatchprocessor import InstructionBatchProcessor
except ImportError:
    # Direct execution - add parent directory to path
    sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
    from llm.instructionbatchprocessor import InstructionBatchProcessor


def validate_k8s_version(version: str) -> str:
    """
    Validate and normalize Kubernetes version format.
    
    Args:
        version: The Kubernetes version string to validate
        
    Returns:
        The normalized version string in format [MAJOR].[MINOR].[REVISION]
        
    Raises:
        argparse.ArgumentTypeError: If the version format is invalid
    """
    if not version:
        raise argparse.ArgumentTypeError("Kubernetes version cannot be empty")
    
    # Remove leading 'v' if present
    normalized_version = version.lstrip('v')
    
    # Validate format: MAJOR.MINOR.REVISION (all must be numbers)
    version_pattern = r'^(\d+)\.(\d+)\.(\d+)$'
    if not re.match(version_pattern, normalized_version):
        raise argparse.ArgumentTypeError(
            f"Invalid Kubernetes version format: '{version}'. "
            f"Expected format: [MAJOR].[MINOR].[REVISION] (e.g., 1.31.8 or v1.31.8)"
        )
    
    return normalized_version


def main() -> None:
    """
    Main function to parse command line arguments and process instructions.
    """
    parser = argparse.ArgumentParser(
        description="Process instructions using InstructionBatchProcessor",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python -m llm.main --code-root-path /path/to/code --instruction-path /path/to/instructions
  python -m llm.main --code-root-path ./pkg --instruction-path ./llm/instructions --k8s-version 1.31.8
  python main.py --code-root-path ../pkg --instruction-path ./instructions --k8s-version v1.30.10
        """
    )
    
    parser.add_argument(
        "--code-root-path",
        required=True,
        help="Root path of the source code directory"
    )
    
    parser.add_argument(
        "--instruction-path",
        required=True, 
        help="Directory path containing instruction files (.md files)"
    )
    
    parser.add_argument(
        "-v", "--verbose",
        action="store_true",
        help="Enable verbose output"
    )
    
    parser.add_argument(
        "--k8s-version",
        type=validate_k8s_version,
        help="Kubernetes version in format [MAJOR].[MINOR].[REVISION] (e.g., 1.31.8 or v1.31.8)"
    )
    
    args = parser.parse_args()
    
    try:
        # Convert to absolute paths for clarity
        code_root_path = os.path.abspath(args.code_root_path)
        instruction_path = os.path.abspath(args.instruction_path)
        
        if args.verbose:
            print(f"Code root path (absolute): {code_root_path}")
            print(f"Instruction path (absolute): {instruction_path}")
            if args.k8s_version:
                print(f"Kubernetes version: {args.k8s_version}")
        
        # Create and use the batch processor
        batch_processor = InstructionBatchProcessor(
            code_root_path=code_root_path,
            instruction_path=instruction_path,
            verbose=args.verbose,
            k8s_version=args.k8s_version
        )
        batch_processor.process_instructions()
        
    except (FileNotFoundError, NotADirectoryError) as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)
    except KeyboardInterrupt:
        print("\nOperation cancelled by user.", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"Unexpected error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()