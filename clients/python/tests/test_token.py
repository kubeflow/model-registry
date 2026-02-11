"""Tests for token utility functions."""

import base64
import json

import pytest

from model_registry.signing.token import decode_jwt_payload, extract_client_id


def create_jwt_token(payload_data: dict) -> str:
    """Helper to create a JWT token for testing."""
    header = base64.urlsafe_b64encode(json.dumps({"alg": "RS256"}).encode()).decode().rstrip("=")
    payload = base64.urlsafe_b64encode(json.dumps(payload_data).encode()).decode().rstrip("=")
    signature = "fake-signature"
    return f"{header}.{payload}.{signature}"


class TestDecodeJwtPayload:
    """Test decode_jwt_payload function."""

    @pytest.mark.parametrize("payload_data", [
        {"sub": "test-user", "aud": "test-audience", "iss": "test-issuer"},
        {"sub": "x"},  # Short payload to test padding
        {
            "sub": "system:serviceaccount:test:sa",
            "aud": ["aud1", "aud2"],
            "kubernetes.io": {
                "namespace": "test",
                "serviceaccount": {"name": "sa", "uid": "123"}
            }
        },
    ])
    def test_decode_valid_jwt(self, payload_data):
        """Test decoding valid JWT tokens with various payloads."""
        token = create_jwt_token(payload_data)
        result = decode_jwt_payload(token)
        assert result == payload_data

    @pytest.mark.parametrize(("invalid_token", "error_pattern"), [
        ("not-a-valid-token", "Failed to decode JWT token payload"),
        ("header.payload", "Failed to decode JWT token payload"),  # Missing signature
        ("header..signature", "Failed to decode JWT token payload"),  # Empty payload
    ])
    def test_decode_invalid_format_raises_error(self, invalid_token, error_pattern):
        """Test that invalid token formats raise ValueError."""
        with pytest.raises(ValueError, match=error_pattern):
            decode_jwt_payload(invalid_token)

    def test_decode_invalid_base64_raises_error(self):
        """Test that invalid base64 in payload raises ValueError."""
        header = base64.urlsafe_b64encode(json.dumps({"alg": "RS256"}).encode()).decode().rstrip("=")
        payload = "!!!invalid-base64!!!"
        signature = "fake-signature"
        token = f"{header}.{payload}.{signature}"

        with pytest.raises(ValueError, match="Failed to decode JWT token payload"):
            decode_jwt_payload(token)

    def test_decode_invalid_json_raises_error(self):
        """Test that invalid JSON in payload raises ValueError."""
        header = base64.urlsafe_b64encode(json.dumps({"alg": "RS256"}).encode()).decode().rstrip("=")
        payload = base64.urlsafe_b64encode(b"not-valid-json").decode().rstrip("=")
        signature = "fake-signature"
        token = f"{header}.{payload}.{signature}"

        with pytest.raises(ValueError, match="Failed to decode JWT token payload"):
            decode_jwt_payload(token)


class TestExtractClientId:
    """Test extract_client_id function."""

    @pytest.mark.parametrize(("claims", "expected"), [
        # String aud
        ({"aud": "https://example.com/issuer"}, "https://example.com/issuer"),
        # Array with single element
        ({"aud": ["https://example.com/issuer"]}, "https://example.com/issuer"),
        # Array with multiple elements (returns first)
        ({"aud": ["first", "second", "third"]}, "first"),
        # Empty array
        ({"aud": []}, None),
        # Missing aud
        ({"sub": "test-user", "iss": "test-issuer"}, None),
        # None aud
        ({"aud": None}, None),
        # Kubernetes token
        ({
            "aud": ["https://kubernetes.default.svc"],
            "iss": "https://kubernetes.default.svc",
            "kubernetes.io": {"namespace": "default"}
        }, "https://kubernetes.default.svc"),
    ])
    def test_extract_client_id(self, claims, expected):
        """Test extracting client_id from various claim structures."""
        result = extract_client_id(claims)
        assert result == expected

    def test_extract_preserves_claims(self):
        """Test that extract_client_id doesn't modify the claims dict."""
        original_claims = {"aud": "test", "sub": "subject"}
        claims = original_claims.copy()

        extract_client_id(claims)

        assert claims == original_claims
