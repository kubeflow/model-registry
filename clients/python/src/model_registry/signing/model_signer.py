"""ModelSigner for signing and verifying ML models with Sigstore."""

from __future__ import annotations

import logging
import os
from collections.abc import Iterable
from pathlib import Path
from typing import TYPE_CHECKING, TypeAlias

from model_signing import signing, verifying
from typing_extensions import Self

from model_registry.signing import sign_sigstore

from .exceptions import InitializationError, SigningError, VerificationError
from .token import decode_jwt_payload, extract_client_id
from .trust_manager import TrustManager

if TYPE_CHECKING:
    from .config import SigningConfig

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

PathLike: TypeAlias = str | os.PathLike[str]


class ModelSigner:
    """Tool for signing and verifying models.

    This uses OpenSSF Model Signing (OMS) to produce a single model.sig file for a
    directory containing a collection of files.
    """

    def __init__(
        self,
        tuf_url: str | None = None,
        root_url: str | None = None,
        root_checksum: str | None = None,
        identity_token_path: PathLike | None = None,
        fulcio_url: str | None = None,
        rekor_url: str | None = None,
        tsa_url: str | None = None,
        certificate_identity: str | None = None,
        oidc_issuer: str | None = None,
        cache_dir: PathLike | None = None,
        signature_filename: str | None = None,
        ignore_paths: Iterable[PathLike] | None = None,
    ):
        """Initialize ModelSigner with optional default parameters.

        Args:
            tuf_url: Default TUF server URL
            root_url: URL to TUF root.json for bootstrap (optional)
            root_checksum: SHA256 checksum of root.json for validation (optional)
            identity_token_path: Default path to identity token file
            fulcio_url: Default Fulcio (CA) server URL
            rekor_url: Default Rekor (transparency log) server URL
            tsa_url: Default TSA (timestamp authority) server URL
            certificate_identity: Default certificate identity for verification
            oidc_issuer: Default OIDC issuer URL
            cache_dir: Base cache directory (default: platformdirs user_data_dir)
                Each TUF URL gets URL-encoded subdirectory for isolation
            signature_filename: Default filename for signatures (default: model.sig)
            ignore_paths: Default paths to ignore during signing (optional)
        """
        self.tuf_url = tuf_url
        self.root_url = root_url
        self.root_checksum = root_checksum
        self.identity_token_path = identity_token_path
        self.fulcio_url = fulcio_url
        self.rekor_url = rekor_url
        self.tsa_url = tsa_url
        self.certificate_identity = certificate_identity
        self.oidc_issuer = oidc_issuer
        self.signature_filename = signature_filename if signature_filename is not None else "model.sig"
        self.ignore_paths = ignore_paths
        # Create trust initializer
        self.trust_manager = TrustManager(cache_dir)

    @classmethod
    def from_config(cls, config: SigningConfig) -> Self:
        """Create ModelSigner from a SigningConfig instance."""
        return cls(
            tuf_url=config.tuf_url,
            root_url=config.root_url,
            root_checksum=config.root_checksum,
            identity_token_path=config.identity_token_path,
            fulcio_url=config.fulcio_url,
            rekor_url=config.rekor_url,
            tsa_url=config.tsa_url,
            certificate_identity=config.certificate_identity,
            oidc_issuer=config.oidc_issuer,
            cache_dir=config.cache_dir,
            signature_filename=config.signature_filename,
            ignore_paths=config.ignore_paths,
        )

    def get_trust_config_path(self) -> Path:
        """Get the path to the cached trust configuration file.

        Returns:
            Path to trust_config.json

        Raises:
            SigningError: If tuf_url is not configured
        """
        if not self.tuf_url:
            msg = "tuf_url is not configured. Provide it during instantiation."
            raise SigningError(msg)
        return self.trust_manager.get_trust_config_path(self.tuf_url)

    def _ensure_trust_initialized(self):
        """Ensure trust configuration is initialized.

        Auto-initializes trust config if not already present using instance defaults.

        Raises:
            SigningError: If tuf_url is not configured
        """
        trust_config_path = self.get_trust_config_path()
        if not trust_config_path.exists():
            logger.info("Trust configuration not initialized. Auto-initializing...")
            self.initialize(
                fulcio_url=self.fulcio_url,
                rekor_url=self.rekor_url,
                tsa_url=self.tsa_url,
                oidc_issuer=self.oidc_issuer,
                force=True,
            )

    def initialize(  # noqa: C901
        self,
        tuf_url: str | None = None,
        root_url: str | None = None,
        root_checksum: str | None = None,
        fulcio_url: str | None = None,
        rekor_url: str | None = None,
        tsa_url: str | None = None,
        oidc_issuer: str | None = None,
        force: bool = False,
    ):
        """Download and cache trust configuration from TUF.

        Args:
            tuf_url: TUF URL (uses instance default if not provided)
            root_url: URL to TUF root.json for bootstrap (uses instance default if not provided)
            root_checksum: SHA256 checksum of root.json (uses instance default if not provided)
            fulcio_url: Fulcio URL (optional - extracted from TUF if available, otherwise required)
            rekor_url: Rekor URL (optional - extracted from TUF if available, otherwise required)
            tsa_url: TSA URL (optional - extracted from TUF if available, otherwise required)
            oidc_issuer: OIDC issuer (uses instance default if not provided)
            force: If True, overwrite existing cached config

        Raises:
            InitializationError: If initialization fails
            FileExistsError: If trust config exists and force=False

        Note:
            Service URLs (fulcio_url, rekor_url, tsa_url) will be extracted from
            TUF trusted_root metadata if available. You can override them with
            explicit parameters for custom deployments or environments.
        """
        # Use method args or fall back to instance defaults
        if tuf_url is None:
            tuf_url = self.tuf_url
        if root_url is None:
            root_url = self.root_url
        if root_checksum is None:
            root_checksum = self.root_checksum
        if fulcio_url is None:
            fulcio_url = self.fulcio_url
        if rekor_url is None:
            rekor_url = self.rekor_url
        if tsa_url is None:
            tsa_url = self.tsa_url
        if oidc_issuer is None:
            oidc_issuer = self.oidc_issuer

        if not tuf_url:
            msg = "tuf_url is required"
            raise InitializationError(msg)

        # Delegate to TrustManager
        self.trust_manager.initialize(
            tuf_url=tuf_url,
            root_url=root_url,
            root_checksum=root_checksum,
            fulcio_url=fulcio_url,
            rekor_url=rekor_url,
            tsa_url=tsa_url,
            oidc_issuer=oidc_issuer,
            force=force,
        )

    def sign(  # noqa: C901
        self,
        model_path: PathLike,
        signature_path: PathLike | None = None,
        identity_token_path: PathLike | None = None,
        fulcio_url: str | None = None,
        rekor_url: str | None = None,
        tsa_url: str | None = None,
        client_id: str | None = None,
        ignore_paths: Iterable[PathLike] | None = None,
    ) -> Path:
        """Sign a model directory.

        Args:
            model_path: Path to model directory
            signature_path: Path where signature file should be written (default: model_path/signature_filename)
            identity_token_path: Path to identity token file (uses instance default if not provided)
            fulcio_url: Fulcio URL (uses instance default if not provided)
            rekor_url: Rekor URL (uses instance default if not provided)
            tsa_url: TSA URL (uses instance default if not provided)
            client_id: OIDC client ID (extracted from token if not provided)
            ignore_paths: Paths to ignore during signing (uses instance default if not provided)

        Returns:
            Path to the signature file

        Raises:
            SigningError: If signing fails or required parameters missing
            FileNotFoundError: If model_path or identity_token_path doesn't exist
            ValueError: If token format is invalid
        """
        model_path = Path(model_path)
        if not model_path.exists():
            msg = f"Model path not found: {model_path}. Expected directory containing model files."
            raise FileNotFoundError(msg)
        if not model_path.is_dir():
            msg = f"Model path must be a directory, got file: {model_path}"
            raise FileNotFoundError(msg)

        # Use method args or fall back to instance defaults
        if identity_token_path is None:
            identity_token_path = self.identity_token_path
        if fulcio_url is None:
            fulcio_url = self.fulcio_url
        if rekor_url is None:
            rekor_url = self.rekor_url
        if tsa_url is None:
            tsa_url = self.tsa_url
        if ignore_paths is None:
            ignore_paths = self.ignore_paths

        # Validate identity_token_path is provided
        if identity_token_path is None:
            msg = "Identity token path is required for signing. Provide via parameter or during instantiation."
            raise SigningError(msg)

        # Validate identity token file exists
        token_path = Path(identity_token_path)
        if not token_path.exists():
            msg = f"Identity token file not found: {token_path}. Expected JWT token file."
            raise SigningError(msg)

        try:
            logger.info(f"Signing model: {model_path}")

            # Ensure trust configuration is initialized
            self._ensure_trust_initialized()

            # Load trust config path
            logger.info("Initializing signing context...")
            trust_config_path = self.get_trust_config_path()

            # Read and parse identity token
            token_str = token_path.read_text().strip()

            # Extract client_id from token if not provided
            if client_id is None:
                claims = decode_jwt_payload(token_str)
                client_id = extract_client_id(claims)

            # Create signer with trust configuration
            signer = sign_sigstore.Signer(
                identity_token=token_str,
                oidc_issuer=self.oidc_issuer,
                client_id=client_id,
                trust_config=trust_config_path,
            )

            # Use provided signature_path or default to model_path/signature_filename
            signature_path = (
                model_path / self.signature_filename if signature_path is None else Path(signature_path)
            ).absolute()

            # Sign using model-signing API
            config = signing.Config()

            # paths added to _ignored_paths will be used only if they are a
            # subpath of the model_path dir
            config._hashing_config._ignored_paths |= {signature_path}
            if ignore_paths is not None:
                config._hashing_config._ignored_paths |= {Path(p) for p in ignore_paths}
            config._signer = signer

            config.sign(model_path, signature_path)

            logger.info(f"Signed successfully: {signature_path}")
            return signature_path

        except (FileNotFoundError, SigningError, ValueError):
            raise
        except Exception as e:
            msg = f"Signing failed: {e}"
            raise SigningError(msg) from e

    def verify(  # noqa: C901
        self,
        model_path: PathLike,
        signature_path: PathLike | None = None,
        certificate_identity: str | None = None,
        oidc_issuer: str | None = None,
    ) -> None:
        """Verify a signed model.

        Verifies the model signature and raises an exception if verification fails.

        Args:
            model_path: Path to model directory
            signature_path: Path to signature file (default: model_path/signature_filename)
            certificate_identity: Expected certificate identity (uses instance default if not provided)
            oidc_issuer: OIDC issuer (uses instance default if not provided)

        Raises:
            VerificationError: If verification fails or required parameters missing
            FileNotFoundError: If model_path or signature doesn't exist
        """
        model_path = Path(model_path)
        if not model_path.exists():
            msg = f"Model path not found: {model_path}. Expected directory containing signed model."
            raise FileNotFoundError(msg)

        # Use provided signature_path or default to model_path/signature_filename
        signature_path = model_path / self.signature_filename if signature_path is None else Path(signature_path)

        if not signature_path.exists():
            msg = f"Signature not found at {signature_path}. Model may not be signed."
            raise FileNotFoundError(msg)

        # Use method args or fall back to instance defaults
        if certificate_identity is None:
            certificate_identity = self.certificate_identity
        if oidc_issuer is None:
            oidc_issuer = self.oidc_issuer

        # Validate required parameters
        if certificate_identity is None:
            msg = "certificate_identity is required for verification. Provide via parameter or during instantiation."
            raise VerificationError(msg)

        if oidc_issuer is None:
            msg = "oidc_issuer is required for verification. Provide via parameter or during instantiation."
            raise VerificationError(msg)

        try:
            logger.info(f"Verifying model: {model_path}")

            # Ensure trust configuration is initialized
            self._ensure_trust_initialized()

            # Load trust config path
            trust_config_path = self.get_trust_config_path()

            logger.info(f"Expected identity: {certificate_identity}")
            logger.info(f"Expected issuer: {oidc_issuer}")

            # Verify using model_signing.verifying API
            verifying.Config().use_sigstore_verifier(
                identity=certificate_identity,
                oidc_issuer=oidc_issuer,
                trust_config=trust_config_path,
            ).verify(model_path, signature_path)

            logger.info("Successfully verified model signature")

        except FileNotFoundError:
            raise
        except ValueError as e:
            # verifying.Config().verify() raises ValueError on verification failure
            logger.error(f"Verification failed: {e}")
            msg = f"Verification failed: {e}"
            raise VerificationError(msg) from e
        except VerificationError:
            raise
        except Exception as e:
            msg = f"Verification failed: {e}"
            raise VerificationError(msg) from e
