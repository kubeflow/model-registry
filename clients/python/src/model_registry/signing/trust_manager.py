"""TrustManager for TUF trust configuration initialization and caching."""

import hashlib
import json
import logging
import os
import tempfile
import urllib.error
import urllib.request
from pathlib import Path
from urllib import parse

import platformdirs
from sigstore._internal.tuf import TrustUpdater

from .exceptions import InitializationError

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)


class TrustManager:
    """Manages TUF trust configuration initialization and caching.

    Handles downloading, validating, and caching root.json from TUF servers,
    with support for explicit bootstrap (like cosign initialize).
    """

    def __init__(self, cache_dir: str | os.PathLike[str] | None = None):
        """Initialize TrustManager.

        Args:
            cache_dir: Base cache directory (default: platformdirs user_data_dir)
                Each TUF URL gets URL-encoded subdirectory for isolation
        """
        if cache_dir:
            self.base_cache_dir = Path(cache_dir)
        else:
            self.base_cache_dir = Path(platformdirs.user_data_dir("model-registry")) / "signing"

    def get_tuf_cache_dir(self, tuf_url: str) -> Path:
        """Get cache directory for a specific TUF URL with URL-encoded isolation."""
        repo_base = parse.quote(tuf_url, safe="")
        return self.base_cache_dir / repo_base

    def get_root_json_path(self, tuf_url: str) -> Path:
        """Get cached root.json path for bootstrap."""
        return self.get_tuf_cache_dir(tuf_url) / "root.json"

    def get_trust_root_path(self, tuf_url: str) -> Path:
        """Get trusted_root.json cache path."""
        return self.get_tuf_cache_dir(tuf_url) / "root" / "targets" / "trusted_root.json"

    def get_trust_config_path(self, tuf_url: str) -> Path:
        """Get trust_config.json cache path."""
        return self.get_tuf_cache_dir(tuf_url) / "trust_config.json"

    def initialize(
        self,
        tuf_url: str,
        root_url: str | None = None,
        root_checksum: str | None = None,
        fulcio_url: str | None = None,
        rekor_url: str | None = None,
        tsa_url: str | None = None,
        oidc_issuer: str | None = None,
        force: bool = False,
    ) -> dict:
        """Download and cache trust configuration from TUF.

        Args:
            tuf_url: TUF repository URL
            root_url: Explicit root.json URL (like cosign --root)
            root_checksum: SHA256 checksum to validate root.json
            fulcio_url: Fulcio CA URL (optional - may be extracted from TUF)
            rekor_url: Rekor transparency log URL (optional)
            tsa_url: Timestamp authority URL (optional)
            oidc_issuer: OIDC issuer URL
            force: If True, overwrite existing cached config

        Returns:
            Trust configuration dict ready for signing/verification

        Raises:
            InitializationError: If initialization fails
            FileExistsError: If config exists and force=False
        """
        config_path = self.get_trust_config_path(tuf_url)
        if config_path.exists() and not force:
            msg = f"Trust configuration already exists at {config_path}. Use force=True to overwrite."
            raise FileExistsError(msg)

        try:
            # Download trusted root
            trusted_root = self._download_trusted_root(tuf_url, root_url=root_url, root_checksum=root_checksum)

            # Transform and create config
            trusted_root = self._transform_and_fix_trust_root(trusted_root)
            trust_config = self._create_trust_config(trusted_root, fulcio_url, rekor_url, tsa_url, oidc_issuer)

            # Cache results
            self._cache_trust_config(tuf_url, trust_config)
            return trust_config

        except Exception as e:
            if isinstance(e, (InitializationError, FileExistsError)):
                raise
            msg = f"Trust initialization failed: {e}"
            raise InitializationError(msg) from e

    def _try_direct_trust_updater(self, tuf_url: str) -> dict | None:
        """Try TrustUpdater without bootstrap (auto-detect mode only)."""
        try:
            updater = TrustUpdater(url=tuf_url, offline=False)
            trusted_root_path = updater.get_trusted_root_path()
            logger.info(f"Downloaded (TrustUpdater): {Path(trusted_root_path).name}")
            with open(trusted_root_path) as f:
                return json.load(f)
        except Exception:
            logger.info("Direct TrustUpdater failed, trying with bootstrap...")
            return None

    def _load_cached_root(self, cache_path: Path, checksum: str | None) -> bytes | None:
        """Load cached root.json and validate checksum if provided."""
        if not cache_path.exists():
            return None

        try:
            cached_data = cache_path.read_bytes()
            if checksum:
                cached_checksum = hashlib.sha256(cached_data).hexdigest()
                if cached_checksum.lower() == checksum.lower():
                    logger.info("Using cached root.json (checksum verified)")
                    return cached_data
                logger.info("Cached root.json checksum mismatch, re-downloading...")
                return None
            logger.info("Using cached root.json")
            return cached_data
        except Exception as e:
            logger.info(f"Could not use cached root.json: {e}, re-downloading...")
            return None

    def _download_root_json(self, urls: list[str]) -> bytes:
        """Download root.json from list of URLs."""
        for url in urls:
            try:
                logger.info(f"Downloading from: {url}")
                with urllib.request.urlopen(url) as response:  # noqa: S310
                    return response.read()
            except urllib.error.HTTPError:
                continue
        msg = f"Could not download root.json from: {urls}"
        raise InitializationError(msg)

    def _validate_checksum(self, data: bytes, expected_checksum: str | None):
        """Validate SHA256 checksum if provided."""
        if expected_checksum:
            computed_checksum = hashlib.sha256(data).hexdigest()
            if computed_checksum.lower() != expected_checksum.lower():
                msg = f"Root.json checksum mismatch!\n  Expected: {expected_checksum}\n  Got:      {computed_checksum}"
                raise InitializationError(msg)
            logger.info(f"Checksum verified: {computed_checksum}")
        else:
            logger.info("No checksum provided (skipping validation)")

    def _bootstrap_trust_updater(self, tuf_url: str, root_data: bytes) -> dict:
        """Bootstrap TrustUpdater with root.json and return trusted_root."""
        with tempfile.TemporaryDirectory() as tmpdir:
            bootstrap_root = Path(tmpdir) / "root.json"
            bootstrap_root.write_bytes(root_data)

            updater = TrustUpdater(url=tuf_url, offline=False, bootstrap_root=bootstrap_root)
            trusted_root_path = updater.get_trusted_root_path()
            logger.info(f"Initialized TUF client: {Path(trusted_root_path).name}")

            with open(trusted_root_path) as f:
                return json.load(f)

    def _download_trusted_root(
        self,
        tuf_url: str,
        root_url: str | None = None,
        root_checksum: str | None = None,
    ) -> dict:
        """Download trusted_root.json from TUF server.

        Supports explicit root.json URL (cosign-style) or auto-detection.

        Args:
            tuf_url: TUF repository URL
            root_url: Explicit URL to root.json (optional, auto-detect if not provided)
            root_checksum: SHA256 checksum to validate root.json (optional)

        Returns:
            Parsed trusted_root.json as dict

        Raises:
            InitializationError: If download or initialization fails
        """
        try:
            # Try direct TrustUpdater first (auto-detect mode only)
            if not root_url:
                logger.info(f"Downloading trusted root from TUF: {tuf_url}")
                result = self._try_direct_trust_updater(tuf_url)
                if result:
                    return result

            # Prepare URLs to try
            if root_url:
                logger.info(f"Downloading root.json from explicit URL: {root_url}")
                root_urls = [root_url]
            else:
                tuf_url_clean = tuf_url.rstrip("/")
                root_urls = [
                    f"{tuf_url_clean}/root.json",
                    f"{tuf_url_clean}/1.root.json",
                ]

            # Try cache first
            root_cache_path = self.get_root_json_path(tuf_url)
            root_data = self._load_cached_root(root_cache_path, root_checksum)

            # Download if not cached
            if not root_data:
                root_data = self._download_root_json(root_urls)
                self._validate_checksum(root_data, root_checksum)
                logger.info(f"Downloaded root.json ({len(root_data)} bytes)")

                # Cache for future use
                root_cache_path.parent.mkdir(parents=True, exist_ok=True)
                root_cache_path.write_bytes(root_data)
                logger.info(f"Cached to: {root_cache_path}")

            # Bootstrap TrustUpdater with root.json
            return self._bootstrap_trust_updater(tuf_url, root_data)

        except InitializationError:
            raise
        except Exception as e:
            msg = f"Failed to download trusted root: {e}"
            raise InitializationError(msg) from e

    def _transform_and_fix_trust_root(self, trusted_root: dict) -> dict:
        """Apply transformations to trusted root."""
        logger.info("Transforming checkpointKeyId to logId...")
        return self._transform_checkpoint_to_logid(trusted_root)

    @staticmethod
    def _transform_checkpoint_to_logid(obj):
        """Transform checkpointKeyId to logId recursively (TUF v0.1 → v0.2)."""
        if isinstance(obj, dict):
            if "checkpointKeyId" in obj:
                obj["logId"] = obj.pop("checkpointKeyId")
            for key, value in obj.items():
                obj[key] = TrustManager._transform_checkpoint_to_logid(value)
        elif isinstance(obj, list):
            for i, item in enumerate(obj):
                obj[i] = TrustManager._transform_checkpoint_to_logid(item)
        return obj

    def _create_trust_config(
        self,
        trusted_root: dict,
        fulcio_url: str | None = None,
        rekor_url: str | None = None,
        tsa_url: str | None = None,
        oidc_issuer: str | None = None,
    ) -> dict:
        """Create complete trust configuration from trusted root."""
        tsa_url_full = None
        if tsa_url:
            tsa_url_full = tsa_url.rstrip("/") + "/api/v1/timestamp"

        return {
            "mediaType": "application/vnd.dev.sigstore.clienttrustconfig.v0.1+json",
            "trustedRoot": trusted_root,
            "signingConfig": {
                "mediaType": "application/vnd.dev.sigstore.signingconfig.v0.2+json",
                "caUrls": (
                    [
                        {
                            "url": fulcio_url,
                            "majorApiVersion": 1,
                            "validFor": {"start": "2020-01-01T00:00:00Z"},
                            "operator": "",
                        }
                    ]
                    if fulcio_url
                    else []
                ),
                "oidcUrls": (
                    [
                        {
                            "url": oidc_issuer,
                            "majorApiVersion": 1,
                            "validFor": {"start": "2020-01-01T00:00:00Z"},
                            "operator": "",
                        }
                    ]
                    if oidc_issuer
                    else []
                ),
                "rekorTlogUrls": (
                    [
                        {
                            "url": rekor_url,
                            "majorApiVersion": 1,
                            "validFor": {"start": "2020-01-01T00:00:00Z"},
                            "operator": "",
                        }
                    ]
                    if rekor_url
                    else []
                ),
                "tsaUrls": (
                    [
                        {
                            "url": tsa_url_full,
                            "majorApiVersion": 1,
                            "validFor": {"start": "2020-01-01T00:00:00Z"},
                            "operator": "",
                        }
                    ]
                    if tsa_url_full
                    else []
                ),
                "rekorTlogConfig": {"selector": "ANY"},
                "tsaConfig": {"selector": "ANY"},
            },
        }

    def _cache_trust_config(self, tuf_url: str, trust_config: dict):
        """Cache trust configuration to disk."""
        trust_root_path = self.get_trust_root_path(tuf_url)
        trust_config_path = self.get_trust_config_path(tuf_url)

        # Create cache directories
        trust_root_path.parent.mkdir(parents=True, exist_ok=True)
        trust_config_path.parent.mkdir(parents=True, exist_ok=True)

        # Cache trusted root
        trusted_root = trust_config.get("trustedRoot", {})
        trust_root_path.write_text(json.dumps(trusted_root, indent=2))
        logger.info(f"Cached trusted root to: {trust_root_path}")

        # Cache complete trust config
        trust_config_path.write_text(json.dumps(trust_config, indent=2))
        logger.info(f"Cached trust config to: {trust_config_path}")
