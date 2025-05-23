#!/usr/bin/env python3
"""
Test script for the GoFunctionParser functionality.
"""

import sys
import os

# Add the current directory to Python path
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from add_k8s_version import GoFunctionParser

def test_function_extraction():
    """Test the function extraction capability."""
    
    # Sample Go code with multiple functions
    sample_go_code = '''package main

import (
    "fmt"
    "log"
)

// Simple function
func simpleFunction() {
    fmt.Println("Hello, World!")
}

// Function with parameters and return value
func addNumbers(a int, b int) int {
    return a + b
}

// Function with receiver (method)
func (p *Person) getName() string {
    return p.name
}

// Generic function
func genericFunction[T any](value T) T {
    return value
}

// Complex function with nested blocks
func complexFunction(x int) error {
    if x > 0 {
        for i := 0; i < x; i++ {
            if i%2 == 0 {
                fmt.Printf("Even: %d\\n", i)
            } else {
                fmt.Printf("Odd: %d\\n", i)
            }
        }
    }
    return nil
}

func main() {
    simpleFunction()
    result := addNumbers(5, 3)
    fmt.Println(result)
}
'''

    parser = GoFunctionParser()
    
    # Test extracting different functions
    test_functions = [
        "simpleFunction",
        "addNumbers", 
        "getName",
        "genericFunction",
        "complexFunction",
        "main"
    ]
    
    for func_name in test_functions:
        try:
            function_code, start_line, end_line = parser.extract_function(sample_go_code, func_name)
            print(f"✓ Successfully extracted function '{func_name}':")
            print(f"  Lines: {start_line} to {end_line}")
            print(f"  Code length: {len(function_code)} characters")
            print(f"  First line: {function_code.split(chr(10))[0]}")
            print()
        except ValueError as e:
            print(f"✗ Failed to extract function '{func_name}': {e}")
            print()
    
    # Test function replacement
    try:
        new_simple_function = '''func simpleFunction() {
    fmt.Println("Hello, Modified World!")
    fmt.Println("This function has been updated!")
}'''
        
        modified_code = parser.replace_function(sample_go_code, "simpleFunction", new_simple_function)
        print("✓ Successfully replaced function 'simpleFunction'")
        print(f"Modified code length: {len(modified_code)} characters")
        
        # Verify the replacement worked
        if "Hello, Modified World!" in modified_code:
            print("✓ Function replacement verified")
        else:
            print("✗ Function replacement verification failed")
            
    except ValueError as e:
        print(f"✗ Failed to replace function: {e}")

if __name__ == "__main__":
    test_function_extraction()
