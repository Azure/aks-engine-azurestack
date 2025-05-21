"""
Run tests for the logcore module using pytest.

This script provides more advanced testing capabilities using pytest,
including code coverage reports.
"""

import os
import sys
import subprocess


def main():
    """Run pytest on the logcore module."""
    # Get the directory of this script
    current_dir = os.path.dirname(os.path.abspath(__file__))
    
    # Construct the path to the tests directory
    tests_dir = os.path.join(current_dir, 'tests')    # Define pytest arguments
    pytest_args = [
        '--cov=logcore',
        '--cov-report=term',
        '--cov-report=html:coverage_html',
        '--no-cov-on-fail',
        tests_dir
    ]
    
    # Add coverage config to omit test files
    os.environ['PYTHONPATH'] = current_dir
    os.environ['COVERAGE_FILE'] = os.path.join(current_dir, '.coverage')
    os.environ['COVERAGE_RCFILE'] = os.path.join(current_dir, '.coveragerc')
    
    # Create a .coveragerc file if it doesn't exist
    coverage_rc = os.path.join(current_dir, '.coveragerc')
    if not os.path.exists(coverage_rc):
        with open(coverage_rc, 'w') as f:
            f.write("""
[run]
omit =
    **/test_*.py
    **/tests/**
    **/conftest.py
""")
    
    # Add any command line arguments
    pytest_args.extend(sys.argv[1:])
    
    # Run pytest
    try:
        cmd = [sys.executable, '-m', 'pytest'] + pytest_args
        print(f"Running: {' '.join(cmd)}")
        result = subprocess.run(cmd, check=False)
        return result.returncode
    except Exception as e:
        print(f"Error running pytest: {e}")
        return 1


if __name__ == '__main__':
    sys.exit(main())
