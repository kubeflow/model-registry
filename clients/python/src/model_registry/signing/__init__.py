"""Signing utilities for model registry."""

import logging

from model_registry.signing.config import SigningConfig
from model_registry.signing.exceptions import (
    BaseSigningError,
    InitializationError,
    SigningError,
    VerificationError,
)
from model_registry.signing.image_signer import CommandRunner, ImageSigner
from model_registry.signing.model_signer import ModelSigner
from model_registry.signing.signer import Signer
from model_registry.signing.token import decode_jwt_payload, extract_client_id
from model_registry.signing.trust_manager import TrustManager

__all__ = [
    "CommandRunner",
    "ImageSigner",
    "ModelSigner",
    "Signer",
    "SigningConfig",
    "TrustManager",
    "BaseSigningError",
    "InitializationError",
    "SigningError",
    "VerificationError",
    "decode_jwt_payload",
    "extract_client_id",
]

# Shared parent logger for all signing components, independent of root logger.
# propagate=False prevents messages from leaking to the root logger (no double-printing).
logger = logging.getLogger(__name__)
logger.propagate = False
if not logger.handlers:  # guard against duplicate handlers on re-import
    _handler = logging.StreamHandler()
    _handler.setFormatter(
        logging.Formatter(
            fmt="%(asctime)s.%(msecs)03d - %(name)s:%(levelname)s: %(message)s",
            datefmt="%H:%M:%S",
        )
    )
    logger.addHandler(_handler)
    # Set to DEBUG so the parent logger never gates messages.
    # Actual filtering is done per-instance by InstanceLevelAdapter.
    logger.setLevel(logging.DEBUG)
