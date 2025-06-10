"""
Snippet processor package for extracting and filtering code elements.
"""

from .snippetfilter import SnippetFilter
from .filetype import FileType
from .location import Location
from .snippetprocessor import SnippetProcessor
from .gosnippetprocessor import GoSnippetProcessor
from .pssnippetprocessor import PsSnippetProcessor
from .bashsnippetprocessor import BashSnippetProcessor
from .snippetprocessorfactory import SnippetProcessorFactory, UnsupportedFileTypeError
from .exceptions import SnippetNotFoundError

__all__ = [
    'SnippetFilter', 
    'FileType', 
    'Location', 
    'SnippetProcessor', 
    'GoSnippetProcessor',
    'PsSnippetProcessor',
    'BashSnippetProcessor',
    'SnippetProcessorFactory',
    'UnsupportedFileTypeError',
    'SnippetNotFoundError'
]
