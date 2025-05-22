"""
Example usage of the logging component.

This module demonstrates how to use the logging component in a code update workflow.
"""

import time
import uuid
from typing import Dict, Any

from . import setup_logging, get_logger
from .logger import ConsoleHandler, FileHandler


def process_data(data: Dict[str, Any]) -> Dict[str, Any]:
    """Process data and return the result.
    
    Args:
        data: The data to process.
        
    Returns:
        The processed data.
    """
    # Get a logger with operation context
    operation_id = f"update-{time.strftime('%Y%m%d')}-{str(uuid.uuid4())[:8]}"
    correlation_id = str(uuid.uuid4())
    
    logger = get_logger(context={
        "operationId": operation_id,
        "correlationId": correlation_id
    })
    
    # Start timing the operation
    start_time = time.time()
    
    # Log the start of the operation
    logger.info(f"Processing data for function '{data.get('function')}'", {
        "httpMethod": "POST",
        "uri": "/api/v1/code/update",
        "body": f"{data}"
    })
    
    try:
        # Simulate processing
        time.sleep(0.3)
        
        result = {
            "success": True,
            "function": data.get("function"),
            "template": data.get("template"),
            "processed": True
        }
        
        # Calculate elapsed time
        elapsed_time = time.time() - start_time
        
        # Log the success
        logger.info(
            f"Applied update template to function '{data.get('function')}' in '{data.get('file')}'",
            {
                "elapsedTime": round(elapsed_time, 3),
                "result": result
            }
        )
        
        return result
    except Exception as e:
        # Calculate elapsed time
        elapsed_time = time.time() - start_time
        
        # Log the error
        logger.error(
            f"Failed to apply update template to function '{data.get('function')}'",
            exception=e,
            context={
                "elapsedTime": round(elapsed_time, 3)
            }
        )
        
        raise


def main() -> None:
    """Run the example."""
    # Set up logging with both console and file output
    setup_logging(
        app_name="code-updator",
        level="INFO",
        handlers=[
            ConsoleHandler(),
            FileHandler("/tmp/code-updator.log")
        ],
        default_context={"appName": "code-updator"}
    )
    
    # Get the logger
    logger = get_logger()
    
    # Log application start
    logger.info("Application started")
    
    # Process some data
    try:
        data = {
            "function": "process_data",
            "template": "add_logging",
            "file": "data_utils.py"
        }
        result = process_data(data)
        logger.info("Data processing completed", {"result": result})
    except Exception as e:
        logger.critical("Application encountered a critical error", exception=e)
    
    # Log application end
    logger.info("Application ended")


if __name__ == "__main__":
    main()
