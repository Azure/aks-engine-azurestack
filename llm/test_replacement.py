#!/usr/bin/env python3
"""
Test script to verify function replacement functionality.
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from add_k8s_version import GoFunctionParser

def test_function_replacement():
    """Test function extraction and replacement."""
    
    # Sample Go code with multiple functions
    sample_go_code = '''package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    greet("Alice")
}

func greet(name string) {
    fmt.Printf("Hello, %s!\\n", name)
}

func add(a, b int) int {
    return a + b
}

func multiply(x, y int) int {
    result := x * y
    return result
}
'''

    print("Original Go code:")
    print("=" * 50)
    print(sample_go_code)
    print("=" * 50)
    
    try:
        # Test extracting the greet function
        function_code, start_line, end_line = GoFunctionParser.extract_function(sample_go_code, "greet")
        print(f"\\nExtracted 'greet' function (lines {start_line}-{end_line}):")
        print("-" * 30)
        print(function_code)
        print("-" * 30)
        
        # Test replacing the greet function
        new_greet_function = '''func greet(name string) {
    fmt.Printf("Greetings, %s! Welcome!\\n", name)
    fmt.Println("Have a great day!")
}'''
        
        modified_code = GoFunctionParser.replace_function(sample_go_code, "greet", new_greet_function)
        
        print("\\nModified Go code with replaced 'greet' function:")
        print("=" * 50)
        print(modified_code)
        print("=" * 50)
        
        # Verify the main function is still there
        main_function, _, _ = GoFunctionParser.extract_function(modified_code, "main")
        print("\\nVerification - 'main' function still exists:")
        print("-" * 30)
        print(main_function)
        print("-" * 30)
        
        # Verify the add function is still there
        add_function, _, _ = GoFunctionParser.extract_function(modified_code, "add")
        print("\\nVerification - 'add' function still exists:")
        print("-" * 30)
        print(add_function)
        print("-" * 30)
        
        print("\\n✅ Function replacement test completed successfully!")
        
    except Exception as e:
        print(f"❌ Error during function replacement test: {e}")
        return False
    
    return True

if __name__ == "__main__":
    success = test_function_replacement()
    sys.exit(0 if success else 1)
