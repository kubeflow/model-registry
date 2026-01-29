"""Exceptions for signing operations."""


class ModelSigningError(Exception):
    """Base exception for model signing operations."""


class SigningError(ModelSigningError):
    """Raised when signing fails."""


class VerificationError(ModelSigningError):
    """Raised when verification fails."""


class InitializationError(ModelSigningError):
    """Raised when initialization fails."""
