"""Cosign signing and verification for container images."""

from __future__ import annotations

import os
import shutil
import subprocess
import sys
from functools import partial
from pathlib import Path
from typing import TYPE_CHECKING

from typing_extensions import Self

from model_registry.signing.exceptions import InitializationError, SigningError, VerificationError

if TYPE_CHECKING:
    from .config import SigningConfig


class CommandRunner:
    """Command execution wrapper for running CLI tools."""

    def __init__(self, **kwargs):
        """Initialize CommandRunner with default subprocess.run arguments.

        Args:
            **kwargs: Additional default keyword arguments to pass to subprocess.run
        """
        self._run = partial(subprocess.run, check=True, capture_output=True, text=True, **kwargs)

    def run(self, cmd: list[str]) -> subprocess.CompletedProcess:
        """Run a CLI command and handle output.

        Args:
            cmd: Command and arguments as a list

        Returns:
            CompletedProcess instance

        Raises:
            subprocess.CalledProcessError: If the command fails
        """
        # Explicitly pass current environment to ensure DOCKER_CONFIG and other vars are inherited
        result = self._run(cmd, env=os.environ.copy())
        if result.stderr:
            print(result.stderr, file=sys.stderr)
        return result


class ImageSigner:
    """Tool for signing and verifying container images."""

    def __init__(
        self,
        tuf_url: str | None = None,
        root: str | None = None,
        root_checksum: str | None = None,
        identity_token_path: str | os.PathLike[str] | None = None,
        fulcio_url: str | None = None,
        rekor_url: str | None = None,
        certificate_identity: str | None = None,
        oidc_issuer: str | None = None,
        client_id: str | None = None,
    ):
        """Initialize ImageSigner tool.

        Args:
            tuf_url: Default TUF mirror URL
            root: Default root JSON URL
            root_checksum: Default root JSON checksum
            identity_token_path: Default path to OIDC identity token file
            fulcio_url: Default Fulcio server URL
            rekor_url: Default Rekor server URL
            certificate_identity: Default certificate identity
            oidc_issuer: Default OIDC issuer URL
            client_id: Default OIDC client ID

        Raises:
            FileNotFoundError: If identity_token_path is provided but doesn't exist
        """
        self.runner = CommandRunner()
        self.tuf_url = tuf_url
        self.root = root
        self.root_checksum = root_checksum

        # Validate identity token path if provided
        if identity_token_path is not None and not os.path.exists(identity_token_path):
            msg = f"Identity token file not found: {identity_token_path}"
            raise FileNotFoundError(msg)
        self.identity_token_path = identity_token_path

        self.fulcio_url = fulcio_url
        self.rekor_url = rekor_url
        self.certificate_identity = certificate_identity
        self.oidc_issuer = oidc_issuer
        self.client_id = client_id

    @classmethod
    def from_config(cls, config: SigningConfig) -> Self:
        """Create ImageSigner from a SigningConfig instance."""
        return cls(
            tuf_url=config.tuf_url,
            root=config.root_url,
            root_checksum=config.root_checksum,
            identity_token_path=config.identity_token_path,
            fulcio_url=config.fulcio_url,
            rekor_url=config.rekor_url,
            certificate_identity=config.certificate_identity,
            oidc_issuer=config.oidc_issuer,
            client_id=config.client_id,
        )

    def _get_sigstore_dir(self) -> Path:
        """Get the sigstore directory path.

        Returns:
            Path to the sigstore directory
        """
        return Path.home() / ".sigstore"

    def _ensure_initialized(self):
        """Ensure sigstore configuration is initialized.

        Auto-initializes sigstore config if not already present using instance defaults.
        """
        sigstore_dir = self._get_sigstore_dir()
        if not sigstore_dir.exists():
            self.initialize()

    def initialize(  # noqa: C901
        self,
        tuf_url: str | None = None,
        root: str | None = None,
        root_checksum: str | None = None,
        force: bool = False,
    ):
        """Initialize sigstore configuration.

        Args:
            tuf_url: TUF mirror URL (optional, overrides instance default)
            root: Root JSON URL (optional, overrides instance default)
            root_checksum: Root JSON checksum (optional, overrides instance default)
            force: If True, remove existing sigstore directory; if False, raise error if exists

        Raises:
            FileExistsError: If sigstore directory exists and force=False
        """
        # Check for existing sigstore directory
        sigstore_dir = self._get_sigstore_dir()
        if sigstore_dir.exists():
            if not force:
                msg = f"Sigstore directory already exists: {sigstore_dir}. Use force=True to overwrite."
                raise FileExistsError(msg)
            shutil.rmtree(sigstore_dir)

        # Use method args or fall back to instance attributes
        if tuf_url is None:
            tuf_url = self.tuf_url
        if root is None:
            root = self.root
        if root_checksum is None:
            root_checksum = self.root_checksum

        # Initialize cosign with TUF mirror
        cmd = ["cosign", "initialize"]

        if tuf_url is not None:
            cmd.extend(["--mirror", tuf_url])

        if root is not None:
            cmd.extend(["--root", root])

        if root_checksum is not None:
            cmd.extend(["--root-checksum", root_checksum])

        try:
            self.runner.run(cmd)
        except subprocess.CalledProcessError as e:
            msg = f"Failed to initialize sigstore: {e}"
            raise InitializationError(msg) from e

    def sign(  # noqa: C901
        self,
        image: str,
        identity_token_path: str | os.PathLike[str] | None = None,
        fulcio_url: str | None = None,
        rekor_url: str | None = None,
        client_id: str | None = None,
    ):
        """Sign a container image and upload the signature.

        Args:
            image: Full image reference including digest (e.g., quay.io/user/image@sha256:...)
            identity_token_path: Path to OIDC identity token file (optional, overrides instance default)
            fulcio_url: Fulcio server URL (optional, overrides instance default)
            rekor_url: Rekor server URL (optional, overrides instance default)
            client_id: OIDC client ID (optional, overrides instance default)

        Raises:
            FileNotFoundError: If identity_token_path is provided but doesn't exist
        """
        # Use method args or fall back to instance attributes
        if identity_token_path is None:
            identity_token_path = self.identity_token_path
        if fulcio_url is None:
            fulcio_url = self.fulcio_url
        if rekor_url is None:
            rekor_url = self.rekor_url
        if client_id is None:
            client_id = self.client_id

        # Validate identity token path if provided
        if identity_token_path is not None and not os.path.exists(identity_token_path):
            msg = f"Identity token file not found: {identity_token_path}"
            raise FileNotFoundError(msg)

        self._ensure_initialized()

        cmd = ["cosign", "sign", "-y"]

        if identity_token_path is not None:
            # Read token content from file (cosign doesn't actually accept file paths despite docs)
            with open(identity_token_path) as f:
                token_content = f.read().strip()
            cmd.extend(["--identity-token", token_content])

        if fulcio_url is not None:
            cmd.extend(["--fulcio-url", fulcio_url])

        if rekor_url is not None:
            cmd.extend(["--rekor-url", rekor_url])

        # Use oidc_issuer for --oidc-client-id parameter (matches cosign behavior)
        if self.oidc_issuer is not None:
            cmd.extend(["--oidc-client-id", self.oidc_issuer])

        cmd.append(image)

        try:
            self.runner.run(cmd)
        except subprocess.CalledProcessError as e:
            msg = f"Failed to sign image {image}: {e}"
            raise SigningError(msg) from e

    def verify(self, image: str, certificate_identity: str | None = None, oidc_issuer: str | None = None):
        """Verify a signed container image.

        Args:
            image: Full image reference including digest (e.g., quay.io/user/image@sha256:...)
            certificate_identity: Expected certificate identity (optional, overrides instance default)
            oidc_issuer: OIDC issuer URL (optional, overrides instance default)
        """
        # Use method args or fall back to instance attributes
        if certificate_identity is None:
            certificate_identity = self.certificate_identity
        if oidc_issuer is None:
            oidc_issuer = self.oidc_issuer

        self._ensure_initialized()

        cmd = ["cosign", "verify"]

        if certificate_identity is not None:
            cmd.extend(["--certificate-identity", certificate_identity])

        if oidc_issuer is not None:
            cmd.extend(["--certificate-oidc-issuer", oidc_issuer])

        cmd.append(image)

        try:
            self.runner.run(cmd)
        except subprocess.CalledProcessError as e:
            msg = f"Failed to verify image {image}: {e}"
            raise VerificationError(msg) from e
