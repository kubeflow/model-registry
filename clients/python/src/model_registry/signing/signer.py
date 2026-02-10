"""Unified signer interface for models and container images."""

import os
from collections.abc import Iterable
from pathlib import Path

from model_registry.signing.config import SigningConfig
from model_registry.signing.image_signer import ImageSigner
from model_registry.signing.model_signer import ModelSigner

PathLike = str | os.PathLike[str]


class Signer:
    """An entrypoint for all signing operations.

    This provides methods for signing and verifying models and container images.
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
        client_id: str | None = None,
        cache_dir: PathLike | None = None,
        signature_filename: str = "model.sig",
        ignore_paths: Iterable[PathLike] | None = None,
    ):
        """Initialize Signer with configuration.

        Loads configuration from environment variables by default.
        Individual parameters override environment values.

        Args:
            tuf_url: TUF server URL
            root_url: URL to TUF root.json
            root_checksum: SHA256 checksum of root.json
            identity_token_path: Path to identity token file
            fulcio_url: Fulcio server URL
            rekor_url: Rekor server URL
            tsa_url: TSA server URL
            certificate_identity: Expected certificate identity
            oidc_issuer: OIDC issuer URL
            client_id: OIDC client ID
            cache_dir: Cache directory
            signature_filename: Default signature filename (default: model.sig)
            ignore_paths: Default paths to ignore during signing
        """
        self.config = SigningConfig.create(
            tuf_url=tuf_url,
            root_url=root_url,
            root_checksum=root_checksum,
            identity_token_path=identity_token_path,
            fulcio_url=fulcio_url,
            rekor_url=rekor_url,
            tsa_url=tsa_url,
            certificate_identity=certificate_identity,
            oidc_issuer=oidc_issuer,
            client_id=client_id,
            cache_dir=cache_dir,
            signature_filename=signature_filename,
            ignore_paths=ignore_paths,
        )
        self.model_signer: ModelSigner = ModelSigner.from_config(self.config)
        self.image_signer: ImageSigner = ImageSigner.from_config(self.config)

    def initialize(
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
        """Initialize trust configuration for both model and image signing.

        This downloads TUF metadata and sets up the necessary trust roots.

        Args:
            tuf_url: TUF URL (uses config default if not provided)
            root_url: URL to TUF root.json for bootstrap (uses config default if not provided)
            root_checksum: SHA256 checksum of root.json (uses config default if not provided)
            fulcio_url: Fulcio URL (uses config default if not provided)
            rekor_url: Rekor URL (uses config default if not provided)
            tsa_url: TSA URL (uses config default if not provided)
            oidc_issuer: OIDC issuer (uses config default if not provided)
            force: If True, overwrite existing configuration
        """
        # Initialize model signer trust
        self.model_signer.initialize(
            tuf_url=tuf_url,
            root_url=root_url,
            root_checksum=root_checksum,
            fulcio_url=fulcio_url,
            rekor_url=rekor_url,
            tsa_url=tsa_url,
            oidc_issuer=oidc_issuer,
            force=force,
        )

        # Initialize image signer trust
        self.image_signer.initialize(
            tuf_url=tuf_url,
            root=root_url,
            root_checksum=root_checksum,
            force=force,
        )

    def sign_model(
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
            signature_path: Path where signature file should be written
            identity_token_path: Path to identity token file (uses config default if not provided)
            fulcio_url: Fulcio URL (uses config default if not provided)
            rekor_url: Rekor URL (uses config default if not provided)
            tsa_url: TSA URL (uses config default if not provided)
            client_id: OIDC client ID (extracted from token if not provided)
            ignore_paths: Paths to ignore during signing (uses config default if not provided)

        Returns:
            Path to the signature file

        Raises:
            SigningError: If signing fails or required parameters missing
            FileNotFoundError: If model_path or identity_token_path doesn't exist
        """
        return self.model_signer.sign(
            model_path=model_path,
            signature_path=signature_path,
            identity_token_path=identity_token_path,
            fulcio_url=fulcio_url,
            rekor_url=rekor_url,
            tsa_url=tsa_url,
            client_id=client_id,
            ignore_paths=ignore_paths,
        )

    def verify_model(
        self,
        model_path: PathLike,
        signature_path: PathLike | None = None,
        certificate_identity: str | None = None,
        oidc_issuer: str | None = None,
    ) -> None:
        """Verify a signed model.

        Args:
            model_path: Path to model directory
            signature_path: Path to signature file (default: model_path/model.sig)
            certificate_identity: Expected certificate identity (uses config default if not provided)
            oidc_issuer: OIDC issuer (uses config default if not provided)

        Raises:
            VerificationError: If verification fails or required parameters missing
            FileNotFoundError: If model_path or signature doesn't exist
        """
        self.model_signer.verify(
            model_path=model_path,
            signature_path=signature_path,
            certificate_identity=certificate_identity,
            oidc_issuer=oidc_issuer,
        )

    def sign_image(
        self,
        image: str,
        identity_token_path: PathLike | None = None,
        fulcio_url: str | None = None,
        rekor_url: str | None = None,
        client_id: str | None = None,
    ):
        """Sign a container image and upload the signature.

        Args:
            image: Full image reference including digest (e.g., quay.io/user/image@sha256:...)
            identity_token_path: Path to OIDC identity token file (uses config default if not provided)
            fulcio_url: Fulcio server URL (uses config default if not provided)
            rekor_url: Rekor server URL (uses config default if not provided)
            client_id: OIDC client ID (uses config default if not provided)

        Raises:
            SigningError: If signing fails
            FileNotFoundError: If identity_token_path doesn't exist
        """
        self.image_signer.sign(
            image=image,
            identity_token_path=identity_token_path,
            fulcio_url=fulcio_url,
            rekor_url=rekor_url,
            client_id=client_id,
        )

    def verify_image(
        self,
        image: str,
        certificate_identity: str | None = None,
        oidc_issuer: str | None = None,
    ):
        """Verify a signed container image.

        Args:
            image: Full image reference including digest (e.g., quay.io/user/image@sha256:...)
            certificate_identity: Expected certificate identity (uses config default if not provided)
            oidc_issuer: OIDC issuer URL (uses config default if not provided)

        Raises:
            VerificationError: If verification fails
        """
        self.image_signer.verify(
            image=image,
            certificate_identity=certificate_identity,
            oidc_issuer=oidc_issuer,
        )
