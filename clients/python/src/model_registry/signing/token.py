"""Token utility functions for model signing."""

import base64
import json


def decode_jwt_payload(token: str) -> dict:
    """Decode JWT token payload to dictionary.

    Args:
        token: JWT token string in format header.payload.signature

    Returns:
        Dictionary of claims from the token payload

    Raises:
        ValueError: If token format is invalid or decoding fails

    Note:
        JWT format per OpenID spec: "A JWT is represented as the concatenation of
        the Encoded JWT Header, the JWT Second Part, and the JWT Third Part, in that
        order, with the parts being separated by period ('.') characters."
        https://openid.net/specs/draft-jones-json-web-token-07.html
    """
    try:
        # Validate JWT has exactly 3 parts (header.payload.signature)
        parts = token.split(".")
        if len(parts) != 3:
            msg = f"Invalid JWT format: expected 3 parts separated by '.', got {len(parts)}"
            raise ValueError(msg)

        payload = parts[1]
        # Add proper padding for base64 decoding
        padding = len(payload) % 4
        if padding:
            payload += "=" * (4 - padding)
        return json.loads(base64.b64decode(payload))
    except (IndexError, json.JSONDecodeError, Exception) as e:
        msg = f"Failed to decode JWT token payload: {e}"
        raise ValueError(msg) from e


def extract_client_id(claims: dict) -> str | None:
    """Extract client_id from decoded JWT claims' aud (audience) field.

    The aud claim can be either a string or an array of strings.
    If it's an array, returns the first element.

    Args:
        claims: Decoded JWT claims dictionary

    Returns:
        Client ID from aud claim, or None if not present or empty
    """
    aud = claims.get("aud")

    if isinstance(aud, list):
        return aud[0] if aud else None
    return aud


def k8s_sub_to_certificate_identity(sub: str) -> str:
    """Convert Kubernetes service account sub claim to certificate SAN format.

    Converts from: system:serviceaccount:namespace:serviceaccount
    To: https://kubernetes.io/namespaces/namespace/serviceaccounts/serviceaccount

    Args:
        sub: The sub claim from a Kubernetes service account token

    Returns:
        Certificate identity in SAN format

    Example:
        >>> k8s_sub_to_certificate_identity("system:serviceaccount:project2:wb2")
        'https://kubernetes.io/namespaces/project2/serviceaccounts/wb2'
    """
    # Format: system:serviceaccount:namespace:serviceaccount
    parts = sub.split(":")
    if len(parts) == 4 and parts[0] == "system" and parts[1] == "serviceaccount":
        namespace = parts[2]
        serviceaccount = parts[3]
        return f"https://kubernetes.io/namespaces/{namespace}/serviceaccounts/{serviceaccount}"
    # If not in expected format, return as-is
    return sub
