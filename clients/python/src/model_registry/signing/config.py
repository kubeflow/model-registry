"""Configuration management for signing operations."""

from __future__ import annotations

import os
from collections.abc import Callable, Iterable, Sequence
from pathlib import Path
from typing import Final, TypeAlias, TypeVar, overload

from pydantic import BaseModel
from typing_extensions import Self

from .token import decode_jwt_payload, extract_client_id, k8s_sub_to_certificate_identity

PathLike: TypeAlias = str | os.PathLike[str]

T = TypeVar("T")
U = TypeVar("U")

DEFAULT_IGNORE_PATHS: Final[tuple[PathLike]] = (".cache",)


class SigningConfig(BaseModel):
    """Configuration for model and image signing operations.

    Loads configuration from environment variables with SIGSTORE_ prefix.
    All fields are optional to allow flexible configuration.

    Example:
        >>> # Create with intelligent defaults and env var fallback
        >>> config = SigningConfig.create()
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
    def create(  # noqa: C901
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
        """Create SigningConfig with intelligent defaults and environment variable fallbacks.

        Only the 4 Sigstore service URLs use environment variables:
        - SIGSTORE_TUF_URL
        - SIGSTORE_FULCIO_URL
        - SIGSTORE_REKOR_URL
        - SIGSTORE_TSA_URL

        Other values are calculated or use sensible defaults:
        - identity_token_path defaults to /var/run/secrets/kubernetes.io/serviceaccount/token
          (the service account token from the workspace/workbench where code is running)
        - root_url defaults to {tuf_url}/root.json
        - oidc_issuer, client_id, certificate_identity are extracted from identity token if available
          For Kubernetes service account tokens, certificate_identity is converted from sub claim
          (system:serviceaccount:namespace:name) to certificate SAN format
          (https://kubernetes.io/namespaces/namespace/serviceaccounts/name)

        Args:
            tuf_url: TUF server URL (env: SIGSTORE_TUF_URL)
            root_url: URL to TUF root.json (default: {tuf_url}/root.json)
            root_checksum: SHA256 checksum of root.json
            identity_token_path: Path to identity token file (default: /var/run/secrets/kubernetes.io/serviceaccount/token)
            fulcio_url: Fulcio server URL (env: SIGSTORE_FULCIO_URL)
            rekor_url: Rekor server URL (env: SIGSTORE_REKOR_URL)
            tsa_url: TSA server URL (env: SIGSTORE_TSA_URL)
            certificate_identity: Expected certificate identity (extracted from token if not provided)
            oidc_issuer: OIDC issuer URL (extracted from token if not provided)
            client_id: OIDC client ID (extracted from token if not provided)
            cache_dir: Cache directory
            cosign_bin_url: Cosign binary URL
            signature_filename: Signature filename
            ignore_paths: Paths to ignore during signing

        Returns:
            SigningConfig instance with resolved values
        """
        # Resolve 4 URLs from env vars
        resolved_tuf_url = resolve(tuf_url, "SIGSTORE_TUF_URL")
        resolved_fulcio_url = resolve(fulcio_url, "SIGSTORE_FULCIO_URL")
        resolved_rekor_url = resolve(rekor_url, "SIGSTORE_REKOR_URL")
        resolved_tsa_url = resolve(tsa_url, "SIGSTORE_TSA_URL")

        # Default identity_token_path to Kubernetes service account token
        resolved_identity_token_path: Path | None = None
        if identity_token_path is not None:
            resolved_identity_token_path = Path(identity_token_path) if identity_token_path else None
        else:
            k8s_token_path = Path("/var/run/secrets/kubernetes.io/serviceaccount/token")
            if k8s_token_path.exists():
                resolved_identity_token_path = k8s_token_path

        # Default root_url to {tuf_url}/root.json
        resolved_root_url = root_url
        if resolved_root_url is None and resolved_tuf_url is not None:
            resolved_root_url = f"{resolved_tuf_url}/root.json"

        # Try to extract values from identity token
        resolved_oidc_issuer = oidc_issuer
        resolved_client_id = client_id
        resolved_certificate_identity = certificate_identity

        if resolved_identity_token_path and resolved_identity_token_path.exists():
            try:
                token = resolved_identity_token_path.read_text().strip()
                claims = decode_jwt_payload(token)

                if resolved_oidc_issuer is None:
                    resolved_oidc_issuer = claims.get("iss")

                if resolved_client_id is None:
                    resolved_client_id = extract_client_id(claims)

                if resolved_certificate_identity is None:
                    # For Kubernetes service account tokens, convert sub to certificate SAN format
                    # For other OIDC tokens, prefer email over sub
                    sub = claims.get("sub", "")
                    is_k8s_token = sub.startswith("system:serviceaccount:")
                    if is_k8s_token:
                        resolved_certificate_identity = k8s_sub_to_certificate_identity(sub)
                    else:
                        resolved_certificate_identity = claims.get("email") or sub
            except (OSError, ValueError):
                # If token reading/parsing fails, just use None values
                pass

        # Convert cache_dir and ignore_paths
        resolved_cache_dir = Path(cache_dir).resolve() if cache_dir else None

        if ignore_paths is None:
            ignore_paths = DEFAULT_IGNORE_PATHS
        resolved_ignore_paths = [Path(p) for p in ignore_paths] if ignore_paths else None

        return cls(
            tuf_url=resolved_tuf_url,
            root_url=resolved_root_url,
            root_checksum=root_checksum,
            identity_token_path=resolved_identity_token_path,
            fulcio_url=resolved_fulcio_url,
            rekor_url=resolved_rekor_url,
            tsa_url=resolved_tsa_url,
            certificate_identity=resolved_certificate_identity,
            oidc_issuer=resolved_oidc_issuer,
            client_id=resolved_client_id,
            cache_dir=resolved_cache_dir,
            cosign_bin_url=cosign_bin_url,
            signature_filename=signature_filename,
            ignore_paths=resolved_ignore_paths,
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
