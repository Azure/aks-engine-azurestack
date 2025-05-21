# Pythonic OOD code generation rules for python code

- Use clear, descriptive class names in PascalCase (e.g., `BankAccount`), following PEP 8.
- Give each class a single, well-defined responsibility.
- Encapsulate data with private attributes (prefix with `_`) and expose only necessary interfaces.
- Prefer composition over inheritance unless a natural "is-a" relationship exists.
- Use abstract base classes (`abc.ABC`) for defining interfaces and contracts.
- Name methods with concise, descriptive snake_case verbs.
- Initialize all essential attributes in `__init__` and avoid mutable default arguments.
- Use type hints for all parameters and return values to improve readability and tooling.
- Write clear, concise docstrings for all classes and methods using triple quotes.
- Avoid global variables; encapsulate state and behavior within classes.
- Raise built-in or custom exceptions for error handling, not error codes or magic values.
- Eliminate code duplication by refactoring common logic and following DRY principles.
- Design code to be easily testable, using dependency injection and avoiding hard-coded values.
- Always adhere to PEP 8 for formatting, naming, and code layout to ensure readability and consistency.
- Favor readability, simplicity, and explicitness—“Readability counts” and “Explicit is better than implicit.”
- Ensure no single file exceeds 1500 lines of code.
- Ensure no single method exceeds 200 lines of code.
- Target Python version 3.13; use features and syntax available in Python 3.13+.
