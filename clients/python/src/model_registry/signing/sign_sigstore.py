"""Sigstore signer implementation for model signing."""
# mypy: ignore-errors

import pathlib

from model_signing._signing.sign_sigstore import _DEFAULT_CLIENT_ID, _DEFAULT_CLIENT_SECRET
from model_signing._signing.sign_sigstore import Signer as BaseSigner
from sigstore import models as sigstore_models
from sigstore import oidc as sigstore_oidc
from sigstore import sign as sigstore_signer


class Signer(BaseSigner):
    """Signing using Sigstore."""

    def __init__(
        self,
        *,
        oidc_issuer: str | None = None,
        use_ambient_credentials: bool = True,
        use_staging: bool = False,
        identity_token: str | None = None,
        force_oob: bool = False,
        client_id: str | None = None,
        client_secret: str | None = None,
        trust_config: pathlib.Path | None = None,
    ):
        """Initializes Sigstore signers.

        Needs to set-up a signing context to use the public goods instance and
        machinery for getting an identity token to use in signing.

        Args:
            oidc_issuer: An optional OpenID Connect issuer to use instead of the
              default production one. Only relevant if `use_staging = False`.
              Default is empty, relying on the Sigstore configuration.
            use_ambient_credentials: Use ambient credentials (also known as
              Workload Identity). Default is True. If ambient credentials cannot
              be used (not available, or option disabled), a flow to get signer
              identity via OIDC will start.
            use_staging: Use staging configurations, instead of production. This
              is supposed to be set to True only when testing. Default is False.
            force_oob: If True, forces an out-of-band (OOB) OAuth flow. If set,
              the OAuth authentication will not attempt to open the default web
              browser. Instead, it will display a URL and code for manual
              authentication. Default is False, which means the browser will be
              opened automatically if possible.
            identity_token: An explicit identity token to use when signing,
              taking precedence over any ambient credential or OAuth workflow.
            client_id: An optional client ID to use when performing OIDC-based
              authentication. This is typically used to identify the
              application making the request to the OIDC provider. If not
              provided, the default client ID configured by Sigstore will be
              used.
            client_secret: An optional client secret to use along with the
              client ID when authenticating with the OIDC provider. This is
              required for confidential clients that need to prove their
              identity to the OIDC provider. If not provided, it is assumed
              that the client is public or the provider does not require a
              secret.
            trust_config: A path to a custom trust configuration. When
              provided, the signature verification process will rely on the
              supplied PKI and trust configurations, instead of the default
              Sigstore setup. If not specified, the default Sigstore
              configuration is used.
        """
        # Initializes the signing and issuer contexts based on provided
        # configuration.
        if use_staging:
            trust_config = sigstore_models.ClientTrustConfig.staging()
        elif trust_config:
            trust_config = sigstore_models.ClientTrustConfig.from_json(trust_config.read_text())
        else:
            trust_config = sigstore_models.ClientTrustConfig.production()

        if not oidc_issuer:
            oidc_issuer = trust_config.signing_config.get_oidc_url()

        self._oidc_issuer = oidc_issuer
        self._signing_context = sigstore_signer.SigningContext.from_trust_config(trust_config)
        self._use_ambient_credentials = use_ambient_credentials
        self._identity_token = identity_token
        self._force_oob = force_oob
        self._client_id = client_id or _DEFAULT_CLIENT_ID
        self._client_secret = client_secret or _DEFAULT_CLIENT_SECRET

    def _get_identity_token(self) -> sigstore_oidc.IdentityToken:
        """Obtains an identity token to use in signing.

        The precedence matches that of sigstore-python:
        1) Explicitly supplied identity token
        2) Ambient credential detected in the environment, if enabled
        3) Interactive OAuth flow
        """
        if self._identity_token:
            return sigstore_oidc.IdentityToken(self._identity_token, self._client_id)
        if self._use_ambient_credentials:
            token = sigstore_oidc.detect_credential(self._client_id)
            if token:
                return sigstore_oidc.IdentityToken(token, self._client_id)

        issuer = sigstore_oidc.Issuer(self._oidc_issuer)
        return issuer.identity_token(
            force_oob=self._force_oob,
            client_id=self._client_id,
            client_secret=self._client_secret,
        )
