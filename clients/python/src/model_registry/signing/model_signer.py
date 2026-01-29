"""ModelSigner for signing and verifying ML models with Sigstore."""

import json
import logging
import os
import tempfile
from pathlib import Path

from google.protobuf import json_format
from model_signing import hashing, signing
from model_signing._hashing import io as hashing_io
from model_signing._hashing import memory
from model_signing._signing.sign_sigstore import (
    Signature,
)
from model_signing._signing.sign_sigstore import (
    Signature as ModelSignature,
)
from model_signing._signing.sign_sigstore import (
    Verifier as ModelVerifier,
)
from model_signing._signing.signing import Payload
from sigstore import dsse as sigstore_dsse
from sigstore import oidc as sigstore_oidc
from sigstore.models import ClientTrustConfig
from sigstore.sign import SigningContext

from .exceptions import InitializationError, SigningError, VerificationError
from .trust_manager import TrustManager

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)


class Signer:
    """Custom Signer that integrates model-signing with sigstore.

    Implements the signing.Signer interface to work with model-signing's
    payload and signature types.
    """

    def __init__(
        self,
        identity_token_path: str | os.PathLike[str],
        oidc_issuer: str,
        trust_config: dict,
    ):
        """Initialize Signer with identity token, OIDC issuer, and trust config.

        Args:
            identity_token_path: Path to identity token file
            oidc_issuer: OIDC issuer URL
            trust_config: Trust configuration dict from TUF
        """
        self.identity_token_path = identity_token_path
        self.oidc_issuer = oidc_issuer

        # Create signing context from trust config
        trust_config_obj = ClientTrustConfig.from_json(json.dumps(trust_config))
        self._signing_context = SigningContext.from_trust_config(trust_config_obj)

    def _get_identity_token(self):
        """Get identity token from file."""
        token_str = Path(self.identity_token_path).read_text().strip()
        return sigstore_oidc.IdentityToken(token_str, client_id=self.oidc_issuer)

    def sign(self, payload):
        """Sign payload using sigstore.

        Args:
            payload: model-signing Payload object

        Returns:
            Signature object with signed bundle
        """
        # Convert payload to DSSE statement
        statement = sigstore_dsse.Statement(json_format.MessageToJson(payload.statement.pb).encode("utf-8"))

        # Get identity token and sign
        token = self._get_identity_token()
        with self._signing_context.signer(token) as signer:
            bundle = signer.sign_dsse(statement)

        return Signature(bundle)


class ModelSigner:
    """Tool for signing and verifying ML models with Sigstore."""

    def __init__(
        self,
        tuf_url: str | None = None,
        root: str | None = None,
        root_checksum: str | None = None,
        identity_token_path: str | os.PathLike[str] | None = None,
        fulcio_url: str | None = None,
        rekor_url: str | None = None,
        tsa_url: str | None = None,
        certificate_identity: str | None = None,
        oidc_issuer: str | None = None,
        cache_dir: str | os.PathLike[str] | None = None,
        signature_filename: str = "model.sig",
    ):
        """Initialize ModelSigner with optional default parameters.

        Args:
            tuf_url: Default TUF server URL
            root: URL to TUF root.json for bootstrap (optional)
            root_checksum: SHA256 checksum of root.json for validation (optional)
            identity_token_path: Default path to identity token file
            fulcio_url: Default Fulcio (CA) server URL
            rekor_url: Default Rekor (transparency log) server URL
            tsa_url: Default TSA (timestamp authority) server URL
            certificate_identity: Default certificate identity for verification
            oidc_issuer: Default OIDC issuer URL
            cache_dir: Base cache directory (default: platformdirs user_data_dir)
                Each TUF URL gets URL-encoded subdirectory for isolation
            signature_filename: Filename for signatures (default: model.sig)
        """
        self.tuf_url = tuf_url
        self.root = root
        self.root_checksum = root_checksum
        self.identity_token_path = identity_token_path
        self.fulcio_url = fulcio_url
        self.rekor_url = rekor_url
        self.tsa_url = tsa_url
        self.certificate_identity = certificate_identity
        self.oidc_issuer = oidc_issuer
        self.signature_filename = signature_filename

        # Create trust initializer
        self.trust_manager = TrustManager(cache_dir)

    def _load_trust_config(self, tuf_url: str) -> dict:
        """Load cached trust configuration.

        Args:
            tuf_url: TUF URL to locate the trust config cache

        Returns:
            Trust config dict

        Raises:
            InitializationError: If trust config doesn't exist or can't be parsed
        """
        config_path = self.trust_manager._get_trust_config_path(tuf_url)
        if not config_path.exists():
            msg = "Trust configuration not initialized. Run initialize() first."
            raise InitializationError(msg)

        try:
            return json.loads(config_path.read_text())
        except (OSError, json.JSONDecodeError) as e:
            msg = f"Failed to load trust config: {e}"
            raise InitializationError(msg) from e

    def initialize(  # noqa: C901
        self,
        tuf_url: str | None = None,
        root: str | None = None,
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
            root: URL to TUF root.json for bootstrap (uses instance default if not provided)
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
        if root is None:
            root = self.root
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
            root_url=root,
            root_checksum=root_checksum,
            fulcio_url=fulcio_url,
            rekor_url=rekor_url,
            tsa_url=tsa_url,
            oidc_issuer=oidc_issuer,
            force=force,
        )

    def create_manifest(self, model_path: str | os.PathLike[str]):
        """Create a manifest from model directory.

        Args:
            model_path: Path to model directory

        Returns:
            Model manifest object

        Raises:
            FileNotFoundError: If model_path doesn't exist
        """
        model_path = Path(model_path)
        if not model_path.exists():
            msg = f"Model path not found: {model_path}"
            raise FileNotFoundError(msg)

        logger.info(f"Creating manifest for {model_path.name}...")

        # Use the high-level hashing API
        hashing_config = hashing.Config()
        model_manifest = hashing_config.hash(model_path)

        resources = list(model_manifest.resource_descriptors())
        logger.info(f"Created manifest with {len(resources)} files")
        return model_manifest

    def sign(  # noqa: C901
        self,
        model_path: str | os.PathLike[str],
        identity_token_path: str | os.PathLike[str] | None = None,
        fulcio_url: str | None = None,
        rekor_url: str | None = None,
        tsa_url: str | None = None,
    ) -> Path:
        """Sign a model directory.

        Args:
            model_path: Path to model directory
            identity_token_path: Path to identity token file (uses instance default if not provided)
            fulcio_url: Fulcio URL (uses instance default if not provided)
            rekor_url: Rekor URL (uses instance default if not provided)
            tsa_url: TSA URL (uses instance default if not provided)

        Returns:
            Path to the signature file

        Raises:
            SigningError: If signing fails
            FileNotFoundError: If model_path doesn't exist or trust config not initialized
        """
        model_path = Path(model_path)
        if not model_path.exists():
            msg = f"Model path not found: {model_path}"
            raise FileNotFoundError(msg)
        if not model_path.is_dir():
            msg = f"Model path must be a directory: {model_path}"
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

        # Validate identity token path
        if identity_token_path is not None:
            token_path = Path(identity_token_path)
            if not token_path.exists():
                msg = f"Identity token file not found: {token_path}"
                raise SigningError(msg)

        try:
            logger.info(f"Signing model: {model_path}")

            # Auto-initialize if not already done
            config_path = self.trust_manager._get_trust_config_path(self.tuf_url)
            if not config_path.exists():
                if not self.tuf_url:
                    msg = (
                        "Trust configuration not initialized and tuf_url not provided. "
                        "Call initialize() first or provide tuf_url."
                    )
                    raise SigningError(msg)

                logger.info("Trust configuration not initialized. Auto-initializing...")
                self.initialize(
                    fulcio_url=self.fulcio_url,
                    rekor_url=self.rekor_url,
                    tsa_url=self.tsa_url,
                    oidc_issuer=self.oidc_issuer,
                    force=True,
                )

            # Load trust config and create signing context
            logger.info("Initializing signing context...")
            trust_config = self._load_trust_config(self.tuf_url)

            # Create manifest and payload
            model_manifest = self.create_manifest(model_path)
            Payload(model_manifest)

            # Create custom signer (initializes signing context internally)
            signer = Signer(
                identity_token_path=identity_token_path,
                oidc_issuer=self.oidc_issuer,
                trust_config=trust_config,
            )

            # Sign using model-signing API
            config = signing.Config()
            config._signer = signer
            signature_path = model_path / self.signature_filename
            config.sign(model_path, signature_path)

            logger.info(f"Signed successfully: {signature_path}")
            return signature_path

        except (FileNotFoundError, SigningError):
            raise
        except Exception as e:
            msg = f"Signing failed: {e}"
            raise SigningError(msg) from e

    def verify(  # noqa: C901
        self,
        model_path: str | os.PathLike[str],
        certificate_identity: str | None = None,
        oidc_issuer: str | None = None,
    ) -> bool:
        """Verify a signed model.

        Args:
            model_path: Path to model directory
            certificate_identity: Expected certificate identity (uses instance default if not provided)
            oidc_issuer: OIDC issuer (uses instance default if not provided)

        Returns:
            True if verification passed, False otherwise

        Raises:
            VerificationError: If verification fails
            FileNotFoundError: If model_path or signature doesn't exist
        """
        model_path = Path(model_path)
        if not model_path.exists():
            msg = f"Model path not found: {model_path}"
            raise FileNotFoundError(msg)

        signature_path = model_path / self.signature_filename
        if not signature_path.exists():
            msg = f"Signature not found at {signature_path}"
            raise FileNotFoundError(msg)

        # Use method args or fall back to instance defaults
        if certificate_identity is None:
            certificate_identity = self.certificate_identity
        if oidc_issuer is None:
            oidc_issuer = self.oidc_issuer

        temp_file = None
        try:
            logger.info(f"Verifying model: {model_path}")

            # Auto-initialize if not already done
            config_path = self.trust_manager._get_trust_config_path(self.tuf_url)
            if not config_path.exists():
                if not self.tuf_url:
                    msg = (
                        "Trust configuration not initialized and tuf_url not provided. "
                        "Call initialize() first or provide tuf_url."
                    )
                    raise VerificationError(msg)

                logger.info("Trust configuration not initialized. Auto-initializing...")
                self.initialize(
                    fulcio_url=self.fulcio_url,
                    rekor_url=self.rekor_url,
                    tsa_url=self.tsa_url,
                    oidc_issuer=self.oidc_issuer,
                    force=True,
                )

            # Load trust config
            trust_config = self._load_trust_config(self.tuf_url)

            # Save trust config to temp file (required by verifier)
            temp_file = Path(tempfile.mkdtemp()) / "trust_config.json"
            temp_file.write_text(json.dumps(trust_config))

            # Load signature
            logger.info("Loading signature bundle...")
            signature = ModelSignature.read(signature_path)

            # Create verifier with identity policy
            logger.info("Creating verifier...")
            if certificate_identity:
                logger.info(f"Expected identity: {certificate_identity}")

            verifier = ModelVerifier(
                identity=certificate_identity,
                oidc_issuer=oidc_issuer,
                trust_config=temp_file,
            )

            # Verify signature
            logger.info("Verifying signature...")
            model_manifest = verifier.verify(signature)

            logger.info("Signature verified")
            if certificate_identity:
                logger.info(f"Identity: {certificate_identity}")
            if oidc_issuer:
                logger.info(f"Issuer: {oidc_issuer}")

            # Verify file hashes
            logger.info("Verifying file hashes...")
            resources = list(model_manifest.resource_descriptors())
            verified_count = 0
            failed_count = 0

            for descriptor in resources:
                file_path = model_path / descriptor.identifier

                if not file_path.exists():
                    logger.error(f"Missing: {descriptor.identifier}")
                    failed_count += 1
                    continue

                file_hasher = hashing_io.SimpleFileHasher(file_path, content_hasher=memory.SHA256())
                actual_digest = file_hasher.compute()

                if actual_digest.digest_hex == descriptor.digest.digest_hex:
                    verified_count += 1
                else:
                    logger.error(f"Hash mismatch: {descriptor.identifier}")
                    failed_count += 1

            logger.info("=" * 70)
            logger.info("VERIFICATION RESULTS")
            logger.info("=" * 70)
            logger.info(f"Files verified: {verified_count}/{len(resources)}")

            if failed_count > 0:
                logger.error(f"Files failed: {failed_count}")
                logger.error("VERIFICATION FAILED")
                return False
            logger.info("VERIFICATION PASSED")
            return True

        except (FileNotFoundError, VerificationError):
            raise
        except Exception as e:
            msg = f"Verification failed: {e}"
            raise VerificationError(msg) from e
        finally:
            # Cleanup temp file
            if temp_file and temp_file.exists():
                temp_file.unlink()
                if temp_file.parent.exists():
                    temp_file.parent.rmdir()
