"""Exceptions for signing operations."""


class SigningError(Exception):
    """Raised when image signing fails."""


class VerificationError(Exception):
    """Raised when image verification fails."""


class InitializationError(Exception):
    """Raised when sigstore initialization fails."""
