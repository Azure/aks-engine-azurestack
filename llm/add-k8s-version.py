#!/usr/bin/env python3


import os
import sys
from typing import Optional, Dict, Any
from pathlib import Path

# Add the parent directory to the path to import llm modules
sys.path.append(str(Path(__file__).parent.parent))

from llm.client.llm_client import AzureLLMClient
from llm.logcore.logger import get_logger

def main() -> None:
    """Main entry point for the Kubernetes version addition tool."""
    logger = get_logger(context={"component": "main"})
    
    try:
        logger.info("Starting Kubernetes version addition tool")
        
        # Initialize the LLM client
        llm_client = AzureLLMClient()
        r = llm_client.get_completion("hello world")  
        logger.info(r)
       
        
    except Exception as e:
        logger.error(
            "Kubernetes version addition tool failed",
            exception=e
        )
        print(f"\nError: {str(e)}")
        sys.exit(1)


if __name__ == "__main__":
    main()