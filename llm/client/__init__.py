"""
Azure LLM Client Package.

This package provides a high-level, object-oriented wrapper for interacting
with Azure's Large Language Model (LLM) chat completions API. It abstracts
away authentication, message formatting, and error handling to provide a
clean and simple interface for generating AI completions.

Classes:
    AzureLLMClient: Main client class for Azure LLM interactions.

Example:
    Basic usage:
    
    >>> from llm.client import AzureLLMClient
    >>> client = AzureLLMClient()
    >>> response = client.get_completion("Hello, how are you?")
    >>> print(response)
    
    With custom configuration:
    
    >>> client = AzureLLMClient(
    ...     endpoint="https://your-endpoint.com",
    ...     model="gpt-4",
    ...     token="your-api-token"
    ... )
    >>> response = client.get_completion(
    ...     prompt="Explain quantum computing",
    ...     system_message="You are a helpful physics tutor",
    ...     temperature=0.7
    ... )
"""

from .llm_client import AzureLLMClient

# Public API - what gets imported when using "from llm.client import *"
__all__ = [
    "AzureLLMClient",
]

# Package metadata
__version__ = "1.0.0"
__author__ = "AKSe"
__description__ = "Azure LLM Client wrapper for chat completions API"