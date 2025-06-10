"""
Azure LLM Client Wrapper for chat completions API.

This module provides a reusable, object-oriented wrapper class for interacting
with the Azure LLM chat completions API, abstracting away authentication and
message formatting.
"""

import os
from typing import Optional

from azure.ai.inference import ChatCompletionsClient
from azure.ai.inference.models import SystemMessage, UserMessage
from azure.core.credentials import AzureKeyCredential
from dotenv import load_dotenv

# Try relative import first, fallback to absolute import
try:
    from ..logcore.logger import get_logger
except (ImportError, ValueError):
    from logcore.logger import get_logger

# Load environment variables from .env file
load_dotenv()


class AzureLLMClient:
    """A wrapper class for Azure LLM chat completions API.
    
    This class simplifies sending user prompts and receiving model responses,
    abstracting away authentication and message formatting.
    """
    
    def __init__(
        self, 
        endpoint: Optional[str] = None, 
        model: Optional[str] = None, 
        token: Optional[str] = None
    ) -> None:
        """Initialize the Azure LLM client.
        
        Args:
            endpoint: The Azure inference endpoint URL. If not provided, attempts to read
                     from AZURE_LLM_ENDPOINT environment variable, then defaults to
                     "https://models.github.ai/inference".
            model: The model name or identifier. If not provided, attempts to read
                  from AZURE_LLM_MODEL environment variable, then defaults to
                  "openai/gpt-4.1".
            token: The Azure API token. If not provided, attempts to read
                  from GITHUB_TOKEN environment variable.
                  
        Raises:
            ValueError: If endpoint or model is empty, or if token cannot be found.
        """
        # Get endpoint from parameter, environment variable, or default
        if endpoint is None:
            endpoint = os.getenv("AZURE_LLM_ENDPOINT", "https://models.github.ai/inference")
        
        if not endpoint or not endpoint.strip():
            raise ValueError("endpoint cannot be empty")
        
        # Get model from parameter, environment variable, or default
        if model is None:
            model = os.getenv("AZURE_LLM_MODEL", "openai/gpt-4.1")
        
        if not model or not model.strip():
            raise ValueError("model cannot be empty")
        
        # Get token from parameter or environment variable
        if token is None:
            token = os.getenv("GITHUB_TOKEN")
        
        if not token:
            raise ValueError(
                "token must be provided either as parameter or through "
                "GITHUB_TOKEN environment variable"
            )
        
        self._endpoint = endpoint.strip()
        self._model = model.strip()
        self._token = token
        self._logger = get_logger(context={"component": "AzureLLMClient"})
        self._logger.set_level("WARNING")
        
        # Initialize the Azure chat completions client
        self._client = ChatCompletionsClient(
            endpoint=self._endpoint,
            credential=AzureKeyCredential(self._token)
        )
        
        self._logger.info(
            "Initialized Azure LLM client",
            context={
                "endpoint": self._endpoint,
                "model": self._model
            }
        )
    
    def get_completion(
        self,
        prompt: str,
        system_message: str = "",
        temperature: float = 0.0,
        top_p: float = 1.0
    ) -> str:
        """Get a completion response from the Azure LLM model.
        
        Args:
            prompt: The user's input message.
            system_message: The system message to set context for the model.
                          Defaults to empty string.
            temperature: Sampling temperature (0.0 to 2.0). Defaults to 0.0.
            top_p: Nucleus sampling parameter (0.0 to 1.0). Defaults to 1.0.
            
        Returns:
            The model's response as a string.
            
        Raises:
            ValueError: If prompt is empty.
            RuntimeError: If the API call fails or returns no response.
        """
        if not prompt or not prompt.strip():
            raise ValueError("prompt cannot be empty")
        
        # Validate temperature range
        if not (0.0 <= temperature <= 2.0):
            raise ValueError("temperature must be between 0.0 and 2.0")
        
        # Validate top_p range
        if not (0.0 <= top_p <= 1.0):
            raise ValueError("top_p must be between 0.0 and 1.0")
        
        # Build the messages list
        messages = []
        if system_message.strip():
            messages.append(SystemMessage(content=system_message.strip()))
        
        messages.append(UserMessage(content=prompt.strip()))
        
        # Log request details before sending
        self._logger.info(
            "Preparing to send completion request",
            context={
                "endpoint": self._endpoint,
                "model": self._model,
                "prompt_length": len(prompt),
                "has_system_message": bool(system_message.strip()),
                "temperature": temperature,
                "top_p": top_p,
                "message_count": len(messages)
            }
        )
        
        # Log detailed request body for debugging
        request_body = {
            "messages": [
                {"role": "system" if isinstance(msg, SystemMessage) else "user", 
                 "content": msg.content} for msg in messages
            ],
            "model": self._model,
            "temperature": temperature,
            "top_p": top_p
        }
        
        self._logger.debug(
            "Request body details",
            context={
                "request_body": request_body,
                "full_prompt": prompt.strip(),
                "full_system_message": system_message.strip() if system_message.strip() else None
            }
        )
        
        try:
            # Log just before making the API call
            self._logger.info("Sending request to Azure LLM API")
            
            # Make the API call
            response = self._client.complete(
                messages=messages,
                model=self._model,
                temperature=temperature,
                top_p=top_p
            )
            
            # Log immediately after receiving response
            self._logger.info(
                "Received response from Azure LLM API",
                context={
                    "response_received": True,
                    "choices_available": len(response.choices) if response.choices else 0
                }
            )
            
            # Extract the response content
            if not response.choices:
                raise RuntimeError("API returned no choices in response")
            
            first_choice = response.choices[0]
            if not hasattr(first_choice, 'message') or not first_choice.message:
                raise RuntimeError("API response choice has no message")
            
            content = first_choice.message.content
            if content is None:
                raise RuntimeError("API response message has no content")
            
            # Log detailed response body for debugging
            response_details = {
                "choices_count": len(response.choices),
                "first_choice_role": getattr(first_choice.message, 'role', 'unknown'),
                "response_content_length": len(content),
                "response_content": content,
                "finish_reason": getattr(first_choice, 'finish_reason', 'unknown')
            }
            
            # Add additional response metadata if available
            if hasattr(response, 'usage'):
                response_details["usage"] = {
                    "prompt_tokens": getattr(response.usage, 'prompt_tokens', 'unknown'),
                    "completion_tokens": getattr(response.usage, 'completion_tokens', 'unknown'),
                    "total_tokens": getattr(response.usage, 'total_tokens', 'unknown')
                }
            
            self._logger.debug(
                "Response body details",
                context=response_details
            )
            
            self._logger.info(
                "Successfully processed completion response",
                context={
                    "response_length": len(content),
                    "choices_count": len(response.choices),
                    "finish_reason": getattr(first_choice, 'finish_reason', 'unknown')
                }
            )
            
            return content
            
        except Exception as e:
            self._logger.error(
                "Failed to get completion from Azure LLM API",
                exception=e,
                context={
                    "endpoint": self._endpoint,
                    "model": self._model,
                    "prompt_length": len(prompt)
                }
            )
            
            # Re-raise as RuntimeError for API-related errors
            if "API" in str(e) or "HTTP" in str(e) or "connection" in str(e).lower():
                raise RuntimeError(f"Azure LLM API error: {str(e)}") from e
            
            # Re-raise other exceptions as-is
            raise