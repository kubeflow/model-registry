"""Exceptions for the model registry."""


class StoreException(Exception):
    """Storage related error."""


class MissingMetadata(Exception):
    """Not enough metadata to complete operation."""


class UnsupportedTypeException(StoreException):
    """Raised when an unsupported type is encountered."""


class TypeNotFoundException(StoreException):
    """Raised when a type cannot be found."""


class ServerException(StoreException):
    """Raised when the server returns a bad response."""


class DuplicateException(StoreException):
    """Raised when the user tries to put an object with a conflicting property."""
