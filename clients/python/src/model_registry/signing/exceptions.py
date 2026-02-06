"""Exceptions for signing operations."""


class BaseSigningError(Exception):
    """Base exception for signing operations."""


class SigningError(BaseSigningError):
    """Raised when signing fails."""


class VerificationError(BaseSigningError):
    """Raised when verification fails."""


class InitializationError(BaseSigningError):
    """Raised when initialization fails."""
