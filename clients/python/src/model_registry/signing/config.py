"""Configuration management for signing operations."""

from __future__ import annotations

import os
from collections.abc import Callable, Iterable, Sequence
from pathlib import Path
from typing import TypeAlias, TypeVar, overload

from pydantic import BaseModel
from typing_extensions import Self

PathLike: TypeAlias = str | os.PathLike[str]

T = TypeVar("T")
U = TypeVar("U")


class SigningConfig(BaseModel):
    """Configuration for model and image signing operations.

    Loads configuration from environment variables with SIGNING_ prefix.
    All fields are optional to allow flexible configuration.

    Example:
        >>> # Load from environment variables
        >>> config = SigningConfig.from_env()
        >>> # Or provide values directly
        >>> config = SigningConfig(
        ...     tuf_url="https://tuf.sigstore.dev",
        ...     identity_token_path="/path/to/token.jwt"
        ... )
    """

    tuf_url: str | None = None
    root_url: str | None = None
    root_checksum: str | None = None
    identity_token_path: Path | None = None
    fulcio_url: str | None = None
    rekor_url: str | None = None
    tsa_url: str | None = None
    certificate_identity: str | None = None
    oidc_issuer: str | None = None
    client_id: str | None = None
    cache_dir: Path | None = None
    cosign_bin_url: str | None = None
    signature_filename: str | None = None
    ignore_paths: Sequence[Path] | None = None

    @classmethod
    def from_env(
        cls,
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
        cosign_bin_url: str | None = None,
        signature_filename: str | None = None,
        ignore_paths: Iterable[PathLike] | None = None,
    ) -> Self:
        """Create SigningConfig from environment variables with optional overrides.

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
            cosign_bin_url: Cosign binary URL
            signature_filename: Signature filename (no env var support)
            ignore_paths: Paths to ignore during signing (no env var support)

        Returns:
            SigningConfig instance with values from environment variables and overrides
        """
        return cls(
            tuf_url=resolve(tuf_url, "SIGNING_TUF_URL"),
            root_url=resolve(root_url, "SIGNING_ROOT_URL"),
            root_checksum=resolve(root_checksum, "SIGNING_ROOT_CHECKSUM"),
            identity_token_path=resolve(identity_token_path, "SIGNING_IDENTITY_TOKEN_PATH", Path),
            fulcio_url=resolve(fulcio_url, "SIGNING_FULCIO_URL"),
            rekor_url=resolve(rekor_url, "SIGNING_REKOR_URL"),
            tsa_url=resolve(tsa_url, "SIGNING_TSA_URL"),
            certificate_identity=resolve(certificate_identity, "SIGNING_CERTIFICATE_IDENTITY"),
            oidc_issuer=resolve(oidc_issuer, "SIGNING_OIDC_ISSUER"),
            client_id=resolve(client_id, "SIGNING_CLIENT_ID"),
            cache_dir=resolve(cache_dir, "SIGNING_CACHE_DIR", Path),
            cosign_bin_url=resolve(cosign_bin_url, "SIGNING_COSIGN_BIN_URL"),
            signature_filename=signature_filename,
            ignore_paths=[Path(p) for p in ignore_paths] if ignore_paths else None,
        )


@overload
def resolve(value: None, env_var: str) -> str | None: ...


@overload
def resolve(value: T, env_var: str) -> T: ...


@overload
def resolve(value: None, env_var: str, converter: Callable[[str | T], U]) -> U | None: ...


@overload
def resolve(value: T, env_var: str, converter: Callable[[str | T], U]) -> U: ...


def resolve(
    value: T | None,
    env_var: str,
    converter: Callable[[str | T], U] | None = None,
) -> T | U | str | None:
    """Resolve configuration value with fallback to environment variable.

    Resolution logic:
    1. Use value if not None; otherwise read env_var from environment variables
    2. If result is None or empty string (no value or empty env var), return None
    3. Otherwise, if converter is provided, apply converter to result

    Args:
        value: Direct value that takes precedence over environment variable
        env_var: Environment variable name to read if value is None
        converter: Optional callable to convert the result to desired type

    Returns:
        The resolved and optionally converted value, or None

    Example:
        >>> resolve(None, "PORT")  # env var as string
        "8080"
        >>> resolve(None, "PORT", int)  # env var converted to int
        8080
        >>> resolve("3000", "PORT", int)  # value converted
        3000
        >>> resolve(None, "UNSET_VAR", int)  # missing value and env var
        None
        >>> resolve(None, "EMPTY_VAR", int)  # empty string env var
        None
    """
    result: T | str | None
    if value is not None:
        result = value
    else:
        result = os.getenv(env_var) or None
        if result is None:
            return None

    if converter is not None:
        return converter(result)
    return result
