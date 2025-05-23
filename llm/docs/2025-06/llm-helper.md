
## LLM Client Wrapper Code Generation Requirements

### 1. **Programming Language**
- **Python 3.8+**

### 2. **Purpose**
- Create a reusable, object-oriented wrapper class for interacting with the Azure LLM chat completions API, based on the provided code.
- The wrapper should simplify sending user prompts and receiving model responses, abstracting away authentication and message formatting.

### 3. **Inputs**
- **Initialization parameters:**
  - `endpoint` (str): The Azure inference endpoint URL.
  - `model` (str): The model name or identifier.
  - `token` (str, optional): The Azure API token. If not provided, the class should attempt to read it from the `GITHUB_TOKEN` environment variable.
- **Method parameters:**
  - `prompt` (str): The user’s input message.
  - `system_message` (str, optional): The system message to set context for the model (default: empty string).
  - `temperature` (float, optional): Sampling temperature (default: 1.0).
  - `top_p` (float, optional): Nucleus sampling parameter (default: 1.0).

### 4. **Outputs**
- The main method should return the model’s response as a string (the content of the first choice).

### 5. **Class and Method Design**
- **Class Name:** `AzureLLMClient`
- **Constructor:** Accepts `endpoint`, `model`, and optional `token`.
- **Method:** `get_completion(prompt, system_message="", temperature=1.0, top_p=1.0)`  
  - Returns the model’s response as a string.

### 6. **Constraints**
- Use the Azure SDK imports as in the original code.
- Handle missing or invalid tokens with a clear exception.
- Validate that `prompt` is a non-empty string.
- The wrapper should be thread-safe for concurrent use (stateless per call).
- Do not print output; return values only.
- Include docstrings for the class and all methods.
- Do not hardcode endpoint, model, or token values.
- Do not include any code outside the class definition except for necessary imports.

### 7. **Additional Context/Requirements**
- The wrapper should be easy to extend for future parameters or message types.
- Include type hints for all method signatures.
- Raise meaningful exceptions for API errors or invalid input.
- The code should be ready for use in production scripts or applications.

---

**Example Usage (for documentation only, not to be included in the code):**
```python
client = AzureLLMClient(endpoint="https://models.github.ai/inference", model="openai/gpt-4.1")
response = client.get_completion("What is the capital of France?")
print(response)  # Should print: "Paris"
```
