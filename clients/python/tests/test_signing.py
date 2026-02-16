"""Tests for signing utilities."""

import subprocess
from pathlib import Path

import pytest

from model_registry.signing import (
    ImageSigner,
    InitializationError,
    ModelSigner,
    Signer,
    SigningConfig,
    SigningError,
    VerificationError,
)

TUF_URL = "https://tuf.example.com"
ROOT_URL = f"{TUF_URL}/root.json"
FULCIO_URL = "https://fulcio.example.com"
REKOR_URL = "https://rekor.example.com"
TSA_URL = "https://tsa.example.com"
OIDC_ISSUER = "https://oidc.example.com"


class TestImageSigner:
    """Test ImageSigner class."""

    def test_initialize_with_all_args(self, tmp_path, mocker):
        """Test initialize with all arguments provided."""
        signer = ImageSigner()
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        # Mock sigstore dir to safe temp location
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=tmp_path / ".sigstore")

        signer.initialize(
            tuf_url=TUF_URL,
            root=ROOT_URL,
            root_checksum="abc123"
        )

        mock_runner.run.assert_called_once()
        cmd = mock_runner.run.call_args[0][0]
        assert cmd == [
            "cosign", "initialize",
            "--mirror", "https://tuf.example.com",
            "--root", ROOT_URL,
            "--root-checksum", "abc123"
        ]

    def test_initialize_with_instance_defaults(self, tmp_path, mocker):
        """Test initialize uses instance defaults when method args not provided."""
        signer = ImageSigner(
            tuf_url="https://default-tuf.example.com",
            root="https://default-tuf.example.com/root.json",
            root_checksum="def456"
        )
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        # Mock sigstore dir to safe temp location
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=tmp_path / ".sigstore")

        signer.initialize()

        cmd = mock_runner.run.call_args[0][0]
        assert cmd.count("https://default-tuf.example.com") == 1
        assert cmd.count("https://default-tuf.example.com/root.json") == 1
        assert cmd.count("def456") == 1

    def test_initialize_method_args_override_defaults(self, tmp_path, mocker):
        """Test method arguments override instance defaults."""
        signer = ImageSigner(
            tuf_url="https://default-tuf.example.com",
            root_checksum="default-checksum"
        )
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        # Mock sigstore dir to safe temp location
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=tmp_path / ".sigstore")

        signer.initialize(
            tuf_url="https://override-tuf.example.com",
            root_checksum="override-checksum"
        )

        cmd = mock_runner.run.call_args[0][0]
        assert cmd.count("https://override-tuf.example.com") == 1
        assert cmd.count("override-checksum") == 1
        assert cmd.count("https://default-tuf.example.com") == 0
        assert cmd.count("default-checksum") == 0

    def test_initialize_raises_if_dir_exists_without_force(self, tmp_path, mocker):
        """Test initialize raises FileExistsError if directory exists and force=False."""
        # Create actual sigstore dir in temp location
        temp_sigstore = tmp_path / ".sigstore"
        temp_sigstore.mkdir()
        (temp_sigstore / "config.json").write_text("{}")

        signer = ImageSigner()
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        # Mock _get_sigstore_dir to return our temp directory
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=temp_sigstore)

        # Should raise because directory exists and force=False
        with pytest.raises(FileExistsError, match="Sigstore directory already exists"):
            signer.initialize()

        # Verify directory was NOT removed
        assert temp_sigstore.exists()

    def test_initialize_removes_sigstore_dir_with_force(self, tmp_path, mocker):
        """Test initialize removes existing .sigstore directory when force=True."""
        # Create actual sigstore dir in temp location
        temp_sigstore = tmp_path / ".sigstore"
        temp_sigstore.mkdir()
        (temp_sigstore / "config.json").write_text("{}")

        signer = ImageSigner()
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        # Mock _get_sigstore_dir to return our temp directory
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=temp_sigstore)

        signer.initialize(force=True)

        # Verify it was actually removed
        assert not temp_sigstore.exists()

    def test_initialize_succeeds_if_dir_not_exists(self, tmp_path, mocker):
        """Test initialize succeeds when directory doesn't exist."""
        temp_sigstore = tmp_path / ".sigstore"

        signer = ImageSigner()
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        # Mock _get_sigstore_dir to return our temp directory (which doesn't exist)
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=temp_sigstore)

        # Should succeed without force since directory doesn't exist
        signer.initialize()

        mock_runner.run.assert_called_once()

    def test_sign_with_all_args(self, tmp_path, mocker):
        """Test sign with all arguments provided."""
        # Create a temporary token file
        token_file = tmp_path / "token"
        token_file.write_text("test-token")

        signer = ImageSigner()
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        # Mock _get_sigstore_dir to return existing directory so initialization is skipped
        sigstore_dir = tmp_path / ".sigstore"
        sigstore_dir.mkdir()
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=sigstore_dir)

        signer.sign(
            image="quay.io/example/image@sha256:abc123",
            identity_token_path=str(token_file),
            fulcio_url=FULCIO_URL,
            rekor_url=REKOR_URL
        )

        mock_runner.run.assert_called_once()
        cmd = mock_runner.run.call_args[0][0]
        assert cmd == [
            "cosign", "sign", "-y",
            "--identity-token", "test-token",
            "--fulcio-url", FULCIO_URL,
            "--rekor-url", REKOR_URL,
            "quay.io/example/image@sha256:abc123"
        ]

    def test_sign_with_instance_defaults(self, tmp_path, mocker):
        """Test sign uses instance defaults when method args not provided."""
        # Create a temporary token file
        token_file = tmp_path / "default-token"
        token_file.write_text("test-token")

        signer = ImageSigner(
            identity_token_path=str(token_file),
            fulcio_url="https://default-fulcio.example.com",
            rekor_url="https://default-rekor.example.com"
        )
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        signer.sign(image="quay.io/example/image@sha256:abc123")

        cmd = mock_runner.run.call_args[0][0]
        assert cmd.count("test-token") == 1
        assert cmd.count("https://default-fulcio.example.com") == 1
        assert cmd.count("https://default-rekor.example.com") == 1

    def test_sign_accepts_pathlike(self, tmp_path, mocker):
        """Test sign accepts PathLike objects."""
        # Create a temporary token file
        token_file = tmp_path / "token"
        token_file.write_text("test-token")

        signer = ImageSigner()
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        # Mock _get_sigstore_dir to return existing directory so initialization is skipped
        sigstore_dir = tmp_path / ".sigstore"
        sigstore_dir.mkdir()
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=sigstore_dir)

        # Pass Path object instead of string
        signer.sign(
            image="quay.io/example/image@sha256:abc123",
            identity_token_path=Path(token_file),
            fulcio_url=FULCIO_URL,
        )

        mock_runner.run.assert_called_once()
        cmd = mock_runner.run.call_args[0][0]
        # File content should be read and passed
        assert cmd.count("test-token") == 1

    def test_sign_raises_if_token_file_not_found(self, mocker):
        """Test sign raises FileNotFoundError if token file doesn't exist."""
        signer = ImageSigner()
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        with pytest.raises(FileNotFoundError, match="Identity token file not found"):
            signer.sign(
                image="quay.io/example/image@sha256:abc123",
                identity_token_path="/nonexistent/token/file"  # noqa: S106
            )

    def test_init_raises_if_token_file_not_found(self):
        """Test __init__ raises FileNotFoundError if token file doesn't exist."""
        with pytest.raises(FileNotFoundError, match="Identity token file not found"):
            ImageSigner(identity_token_path="/nonexistent/token/file")  # noqa: S106

    def test_verify_with_all_args(self, tmp_path, mocker):
        """Test verify with all arguments provided."""
        signer = ImageSigner()
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        # Mock _get_sigstore_dir to return existing directory so initialization is skipped
        sigstore_dir = tmp_path / ".sigstore"
        sigstore_dir.mkdir()
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=sigstore_dir)

        signer.verify(
            image="quay.io/example/image@sha256:abc123",
            certificate_identity="https://kubernetes.io/namespaces/test/serviceaccounts/default",
            oidc_issuer=OIDC_ISSUER
        )

        mock_runner.run.assert_called_once()
        cmd = mock_runner.run.call_args[0][0]
        assert cmd == [
            "cosign", "verify",
            "--certificate-identity", "https://kubernetes.io/namespaces/test/serviceaccounts/default",
            "--certificate-oidc-issuer", OIDC_ISSUER,
            "quay.io/example/image@sha256:abc123"
        ]

    def test_verify_with_instance_defaults(self, mocker):
        """Test verify uses instance defaults when method args not provided."""
        signer = ImageSigner(
            certificate_identity="https://kubernetes.io/namespaces/default/serviceaccounts/default",
            oidc_issuer="https://default-oidc.example.com"
        )
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        signer.verify(image="quay.io/example/image@sha256:abc123")

        cmd = mock_runner.run.call_args[0][0]
        assert cmd.count("https://kubernetes.io/namespaces/default/serviceaccounts/default") == 1
        assert cmd.count("https://default-oidc.example.com") == 1

    def test_initialize_with_no_args_sends_minimal_command(self, tmp_path, mocker):
        """Test initialize with no arguments sends minimal command."""
        signer = ImageSigner()
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        # Mock sigstore dir to safe temp location
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=tmp_path / ".sigstore")

        signer.initialize()

        cmd = mock_runner.run.call_args[0][0]
        assert cmd == ["cosign", "initialize"]

    def test_sign_image_is_last_argument(self, tmp_path, mocker):
        """Test that image reference is always the last argument for sign."""
        token_file = tmp_path / "token"
        token_file.write_text("test-token")

        signer = ImageSigner(identity_token_path=str(token_file))
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        signer.sign(image="quay.io/example/image@sha256:abc123")

        cmd = mock_runner.run.call_args[0][0]
        assert cmd[-1] == "quay.io/example/image@sha256:abc123"

    def test_verify_image_is_last_argument(self, mocker):
        """Test that image reference is always the last argument for verify."""
        signer = ImageSigner(certificate_identity="test-id")
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        signer.verify(image="quay.io/example/image@sha256:abc123")

        cmd = mock_runner.run.call_args[0][0]
        assert cmd[-1] == "quay.io/example/image@sha256:abc123"

    def test_initialize_raises_initialization_error_on_command_failure(self, tmp_path, mocker):
        """Test initialize raises InitializationError when cosign command fails."""
        signer = ImageSigner()
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=tmp_path / ".sigstore")

        # Mock runner to raise CalledProcessError
        mock_runner = mocker.MagicMock()
        mock_runner.run.side_effect = subprocess.CalledProcessError(1, ["cosign", "initialize"])
        signer.runner = mock_runner

        with pytest.raises(InitializationError, match="Failed to initialize sigstore"):
            signer.initialize(tuf_url="https://tuf.example.com")

    def test_sign_raises_signing_error_on_command_failure(self, tmp_path, mocker):
        """Test sign raises SigningError when cosign command fails."""
        token_file = tmp_path / "token"
        token_file.write_text("test-token")

        signer = ImageSigner()

        # Mock _get_sigstore_dir to return existing directory so initialization is skipped
        sigstore_dir = tmp_path / ".sigstore"
        sigstore_dir.mkdir()
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=sigstore_dir)

        # Mock runner to raise CalledProcessError
        mock_runner = mocker.MagicMock()
        mock_runner.run.side_effect = subprocess.CalledProcessError(1, ["cosign", "sign"])
        signer.runner = mock_runner

        with pytest.raises(SigningError, match="Failed to sign image"):
            signer.sign(
                image="quay.io/example/image@sha256:abc123",
                identity_token_path=str(token_file)
            )

    def test_verify_raises_verification_error_on_command_failure(self, tmp_path, mocker):
        """Test verify raises VerificationError when cosign command fails."""
        signer = ImageSigner()

        # Mock _get_sigstore_dir to return existing directory so initialization is skipped
        sigstore_dir = tmp_path / ".sigstore"
        sigstore_dir.mkdir()
        mocker.patch.object(signer, "_get_sigstore_dir", return_value=sigstore_dir)

        # Mock runner to raise CalledProcessError
        mock_runner = mocker.MagicMock()
        mock_runner.run.side_effect = subprocess.CalledProcessError(1, ["cosign", "verify"])
        signer.runner = mock_runner

        with pytest.raises(VerificationError, match="Failed to verify image"):
            signer.verify(
                image="quay.io/example/image@sha256:abc123",
                certificate_identity="https://kubernetes.io/namespaces/test/serviceaccounts/default"
            )


class TestModelSigner:
    """Test ModelSigner class."""

    def test_get_trust_config_path_with_tuf_url(self, tmp_path):
        """Test get_trust_config_path returns correct path when tuf_url is set."""
        signer = ModelSigner(
            tuf_url=TUF_URL,
            cache_dir=tmp_path
        )

        result = signer.get_trust_config_path()

        # Should return path within cache_dir
        assert result.parent.parent == tmp_path
        assert result.name == "trust_config.json"
        assert str(result).count("https%3A%2F%2Ftuf.example.com") == 1

    def test_get_trust_config_path_without_tuf_url_raises_error(self):
        """Test get_trust_config_path raises SigningError when tuf_url not set."""
        signer = ModelSigner()

        with pytest.raises(SigningError, match="tuf_url is not configured"):
            signer.get_trust_config_path()

    def test_ensure_trust_initialized_creates_config_if_missing(self, tmp_path, mocker):
        """Test _ensure_trust_initialized calls initialize if config doesn't exist."""
        signer = ModelSigner(
            tuf_url=TUF_URL,
            fulcio_url=FULCIO_URL,
            rekor_url=REKOR_URL,
            tsa_url=TSA_URL,
            oidc_issuer=OIDC_ISSUER,
            cache_dir=tmp_path
        )

        # Mock initialize method
        mock_initialize = mocker.patch.object(signer, "initialize")

        # Call _ensure_trust_initialized
        signer._ensure_trust_initialized()

        # Should have called initialize
        mock_initialize.assert_called_once_with(
            fulcio_url=FULCIO_URL,
            rekor_url=REKOR_URL,
            tsa_url=TSA_URL,
            oidc_issuer=OIDC_ISSUER,
            force=True
        )

    def test_ensure_trust_initialized_skips_if_config_exists(self, tmp_path, mocker):
        """Test _ensure_trust_initialized doesn't call initialize if config exists."""
        signer = ModelSigner(
            tuf_url=TUF_URL,
            cache_dir=tmp_path
        )

        # Create the trust config file
        config_path = signer.get_trust_config_path()
        config_path.parent.mkdir(parents=True, exist_ok=True)
        config_path.write_text('{"test": "config"}')

        # Mock initialize method
        mock_initialize = mocker.patch.object(signer, "initialize")

        # Call _ensure_trust_initialized
        signer._ensure_trust_initialized()

        # Should NOT have called initialize
        mock_initialize.assert_not_called()


class TestSigningConfig:
    """Test SigningConfig class."""

    def test_loads_environment_variables(self, monkeypatch):
        """Test constructor loads configuration from environment variables for the 4 service URLs."""
        monkeypatch.setenv("SIGSTORE_TUF_URL", "https://env-tuf.example.com")
        monkeypatch.setenv("SIGSTORE_FULCIO_URL", "https://env-fulcio.example.com")
        monkeypatch.setenv("SIGSTORE_REKOR_URL", "https://env-rekor.example.com")
        monkeypatch.setenv("SIGSTORE_TSA_URL", "https://env-tsa.example.com")

        config = SigningConfig.create()

        assert config.tuf_url == "https://env-tuf.example.com"
        assert config.fulcio_url == "https://env-fulcio.example.com"
        assert config.rekor_url == "https://env-rekor.example.com"
        assert config.tsa_url == "https://env-tsa.example.com"

    def test_creates_config_from_kwargs(self):
        """Test creating config from keyword arguments."""
        config = SigningConfig(
            tuf_url="https://tuf.example.com",
            root_url="https://tuf.example.com/root.json",
            fulcio_url="https://fulcio.example.com",
            rekor_url="https://rekor.example.com",
            certificate_identity="user@example.com",
        )

        assert config.tuf_url == "https://tuf.example.com"
        assert config.root_url == "https://tuf.example.com/root.json"
        assert config.fulcio_url == "https://fulcio.example.com"
        assert config.rekor_url == "https://rekor.example.com"
        assert config.certificate_identity == "user@example.com"

    def test_model_dump_exports_all_fields(self):
        """Test model_dump exports all configuration fields."""
        config = SigningConfig(
            tuf_url="https://tuf.example.com",
            fulcio_url="https://fulcio.example.com",
        )

        result = config.model_dump()

        assert "tuf_url" in result
        assert "fulcio_url" in result
        assert "rekor_url" in result
        assert result["tuf_url"] == "https://tuf.example.com"
        assert result["fulcio_url"] == "https://fulcio.example.com"

    def test_model_dump_exclude_none(self):
        """Test model_dump with exclude_none excludes None values."""
        config = SigningConfig(
            tuf_url="https://tuf.example.com",
            fulcio_url=None,
        )

        result = config.model_dump(exclude_none=True)

        assert "tuf_url" in result
        assert "fulcio_url" not in result
        assert result["tuf_url"] == "https://tuf.example.com"

    def test_empty_string_env_vars_treated_as_none(self, monkeypatch):
        """Test empty string environment variables are treated as None for the 4 URL env vars."""
        # Set URL env vars to empty strings (only these 4 use env vars now)
        monkeypatch.setenv("SIGSTORE_TUF_URL", "")
        monkeypatch.setenv("SIGSTORE_FULCIO_URL", "")
        monkeypatch.setenv("SIGSTORE_TSA_URL", "")
        # Set one to a real value to verify it still works
        monkeypatch.setenv("SIGSTORE_REKOR_URL", "https://rekor.example.com")

        config = SigningConfig.create()

        # Empty string env vars should be None
        assert config.tuf_url is None
        assert config.fulcio_url is None
        assert config.tsa_url is None
        # Non-empty env var should work
        assert config.rekor_url == "https://rekor.example.com"


class TestSigner:
    """Test unified Signer class."""

    def test_init_with_individual_args(self):
        """Test __init__ with individual arguments."""
        signer = Signer(
            tuf_url="https://tuf.example.com",
            fulcio_url="https://fulcio.example.com",
        )

        assert signer.config.tuf_url == "https://tuf.example.com"
        assert signer.config.fulcio_url == "https://fulcio.example.com"

    def test_init_loads_create(self, monkeypatch):
        """Test __init__ loads from environment variables."""
        monkeypatch.setenv("SIGSTORE_TUF_URL", "https://env-tuf.example.com")
        signer = Signer()

        assert signer.config.tuf_url == "https://env-tuf.example.com"

    def test_init_with_overrides(self, monkeypatch):
        """Test __init__ with override arguments over environment."""
        monkeypatch.setenv("SIGSTORE_TUF_URL", "https://env-tuf.example.com")
        monkeypatch.setenv("SIGSTORE_FULCIO_URL", "https://env-fulcio.example.com")

        signer = Signer(
            fulcio_url="https://override-fulcio.example.com",
        )

        # tuf_url should come from env
        assert signer.config.tuf_url == "https://env-tuf.example.com"
        # fulcio_url should be overridden
        assert signer.config.fulcio_url == "https://override-fulcio.example.com"

    def test_model_signer_attribute_creates_instance(self):
        """Test model_signer attribute creates ModelSigner instance."""
        signer = Signer(tuf_url="https://tuf.example.com")

        model_signer = signer.model_signer

        assert isinstance(model_signer, ModelSigner)
        assert model_signer.tuf_url == "https://tuf.example.com"

    def test_model_signer_attribute_is_same_instance(self):
        """Test model_signer attribute is same instance."""
        signer = Signer()

        model_signer1 = signer.model_signer
        model_signer2 = signer.model_signer

        assert model_signer1 is model_signer2

    def test_image_signer_attribute_creates_instance(self):
        """Test image_signer attribute creates ImageSigner instance."""
        signer = Signer(tuf_url="https://tuf.example.com")

        image_signer = signer.image_signer

        assert isinstance(image_signer, ImageSigner)
        assert image_signer.tuf_url == "https://tuf.example.com"

    def test_image_signer_attribute_is_same_instance(self):
        """Test image_signer attribute is same instance."""
        signer = Signer()

        image_signer1 = signer.image_signer
        image_signer2 = signer.image_signer

        assert image_signer1 is image_signer2

    def test_signature_filename_passed_to_model_signer(self):
        """Test signature_filename is passed to ModelSigner."""
        signer = Signer(signature_filename="custom.sig")

        assert signer.model_signer.signature_filename == "custom.sig"

    def test_ignore_paths_passed_to_model_signer(self):
        """Test ignore_paths is passed to ModelSigner."""
        ignore = [".git", "__pycache__"]
        signer = Signer(ignore_paths=ignore)

        # Config converts to Path objects
        assert signer.model_signer.ignore_paths == [Path(".git"), Path("__pycache__")]

    def test_sign_model_delegates_to_model_signer(self, tmp_path, mocker):
        """Test sign_model delegates to ModelSigner.sign."""
        signer = Signer(tuf_url="https://tuf.example.com")

        # Mock ModelSigner.sign
        mock_sign = mocker.patch.object(signer.model_signer, "sign")
        mock_sign.return_value = tmp_path / "model.sig"

        result = signer.sign_model(tmp_path)

        mock_sign.assert_called_once()
        assert result == tmp_path / "model.sig"

    def test_verify_model_delegates_to_model_signer(self, tmp_path, mocker):
        """Test verify_model delegates to ModelSigner.verify."""
        signer = Signer(tuf_url="https://tuf.example.com")

        # Mock ModelSigner.verify
        mock_verify = mocker.patch.object(signer.model_signer, "verify")

        signer.verify_model(tmp_path)

        mock_verify.assert_called_once()

    def test_sign_image_delegates_to_image_signer(self, mocker):
        """Test sign_image delegates to ImageSigner.sign."""
        signer = Signer(tuf_url="https://tuf.example.com")

        # Mock ImageSigner.sign
        mock_sign = mocker.patch.object(signer.image_signer, "sign")

        signer.sign_image("quay.io/user/image@sha256:abc123")

        mock_sign.assert_called_once()

    def test_verify_image_delegates_to_image_signer(self, mocker):
        """Test verify_image delegates to ImageSigner.verify."""
        signer = Signer(tuf_url="https://tuf.example.com")

        # Mock ImageSigner.verify
        mock_verify = mocker.patch.object(signer.image_signer, "verify")

        signer.verify_image("quay.io/user/image@sha256:abc123")

        mock_verify.assert_called_once()

    def test_initialize_delegates_to_both_signers(self, tmp_path, mocker):
        """Test initialize delegates to both ModelSigner and ImageSigner."""
        signer = Signer(tuf_url="https://tuf.example.com")

        # Mock both initialize methods
        mock_model_init = mocker.patch.object(signer.model_signer, "initialize")
        mock_image_init = mocker.patch.object(signer.image_signer, "initialize")

        signer.initialize(force=True)

        mock_model_init.assert_called_once()
        mock_image_init.assert_called_once()
        assert mock_model_init.call_args[1]["force"] is True
        assert mock_image_init.call_args[1]["force"] is True


class TestModelSignerFromConfig:
    """Test ModelSigner.from_config classmethod."""

    def test_from_config_creates_model_signer(self):
        """Test from_config creates ModelSigner with config values."""
        config = SigningConfig(
            tuf_url="https://tuf.example.com",
            fulcio_url="https://fulcio.example.com",
            rekor_url="https://rekor.example.com",
        )

        signer = ModelSigner.from_config(config)

        assert signer.tuf_url == "https://tuf.example.com"
        assert signer.fulcio_url == "https://fulcio.example.com"
        assert signer.rekor_url == "https://rekor.example.com"

    def test_from_config_uses_all_config_values(self):
        """Test from_config uses all config values."""
        config = SigningConfig(
            tuf_url="https://tuf.example.com",
            fulcio_url="https://fulcio.example.com",
            rekor_url="https://rekor.example.com",
        )

        signer = ModelSigner.from_config(config)

        assert signer.tuf_url == "https://tuf.example.com"
        assert signer.fulcio_url == "https://fulcio.example.com"
        assert signer.rekor_url == "https://rekor.example.com"


class TestImageSignerFromConfig:
    """Test ImageSigner.from_config classmethod."""

    def test_from_config_creates_image_signer(self):
        """Test from_config creates ImageSigner with config values."""
        config = SigningConfig(
            tuf_url="https://tuf.example.com",
            root_url="https://tuf.example.com/root.json",
            fulcio_url="https://fulcio.example.com",
        )

        signer = ImageSigner.from_config(config)

        assert signer.tuf_url == "https://tuf.example.com"
        assert signer.root == "https://tuf.example.com/root.json"
        assert signer.fulcio_url == "https://fulcio.example.com"

    def test_from_config_uses_all_config_values(self):
        """Test from_config uses all config values."""
        config = SigningConfig(
            tuf_url="https://tuf.example.com",
            fulcio_url="https://fulcio.example.com",
            rekor_url="https://rekor.example.com",
        )

        signer = ImageSigner.from_config(config)

        assert signer.tuf_url == "https://tuf.example.com"
        assert signer.fulcio_url == "https://fulcio.example.com"
        assert signer.rekor_url == "https://rekor.example.com"

    def test_from_config_maps_root_url_to_root(self):
        """Test from_config correctly maps root_url to root parameter."""
        config = SigningConfig(root_url="https://tuf.example.com/root.json")

        signer = ImageSigner.from_config(config)

        # ImageSigner uses 'root' not 'root_url'
        assert signer.root == "https://tuf.example.com/root.json"

    def test_default_root_url_from_tuf_url(self, monkeypatch):
        """Test root_url defaults to {tuf_url}/root.json."""
        monkeypatch.setenv("SIGSTORE_TUF_URL", "https://tuf.example.com")

        config = SigningConfig.create()

        assert config.tuf_url == "https://tuf.example.com"
        assert config.root_url == "https://tuf.example.com/root.json"

    def test_explicit_root_url_overrides_default(self, monkeypatch):
        """Test explicit root_url overrides the default."""
        monkeypatch.setenv("SIGSTORE_TUF_URL", "https://tuf.example.com")

        config = SigningConfig.create(root_url="https://custom.example.com/custom-root.json")

        assert config.tuf_url == "https://tuf.example.com"
        assert config.root_url == "https://custom.example.com/custom-root.json"

    def test_default_k8s_token_path(self, tmp_path, monkeypatch):
        """Test identity_token_path defaults to k8s service account token if it exists."""
        # Create a fake k8s token file
        k8s_token_dir = tmp_path / "var" / "run" / "secrets" / "kubernetes.io" / "serviceaccount"
        k8s_token_dir.mkdir(parents=True)
        k8s_token_path = k8s_token_dir / "token"
        k8s_token_path.write_text("fake-k8s-token")

        # Mock Path to use our temp directory
        import model_registry.signing.config as config_module
        original_path = config_module.Path

        class MockPath(type(tmp_path)):  # type: ignore[misc]
            def __new__(cls, *args):
                if args and args[0] == "/var/run/secrets/kubernetes.io/serviceaccount/token":
                    return original_path(k8s_token_path)
                return original_path(*args)

        monkeypatch.setattr(config_module, "Path", MockPath)

        config = SigningConfig.create()

        assert config.identity_token_path == k8s_token_path

    def test_extract_values_from_token(self, tmp_path):
        """Test extracting oidc_issuer, client_id, and certificate_identity from token."""
        # Create a fake JWT token with the claims we want to extract
        import base64
        import json

        header = base64.urlsafe_b64encode(b'{"alg":"RS256"}').decode().rstrip("=")
        payload_dict = {
            "iss": "https://issuer.example.com",
            "aud": ["client-123", "client-456"],
            "email": "user@example.com",
            "sub": "user-subject-id",
        }
        payload = base64.urlsafe_b64encode(json.dumps(payload_dict).encode()).decode().rstrip("=")
        signature = base64.urlsafe_b64encode(b"fake-signature").decode().rstrip("=")
        token = f"{header}.{payload}.{signature}"

        # Write token to a file
        token_file = tmp_path / "token"
        token_file.write_text(token)

        config = SigningConfig.create(identity_token_path=token_file)

        assert config.oidc_issuer == "https://issuer.example.com"
        assert config.client_id == "client-123"  # First element of aud array
        assert config.certificate_identity == "user@example.com"  # Prefers email over sub

    def test_extract_sub_when_no_email_in_token(self, tmp_path):
        """Test certificate_identity falls back to sub when email is not present."""
        import base64
        import json

        header = base64.urlsafe_b64encode(b'{"alg":"RS256"}').decode().rstrip("=")
        payload_dict = {
            "iss": "https://issuer.example.com",
            "aud": "client-123",
            "sub": "user-subject-id",
        }
        payload = base64.urlsafe_b64encode(json.dumps(payload_dict).encode()).decode().rstrip("=")
        signature = base64.urlsafe_b64encode(b"fake-signature").decode().rstrip("=")
        token = f"{header}.{payload}.{signature}"

        token_file = tmp_path / "token"
        token_file.write_text(token)

        config = SigningConfig.create(identity_token_path=token_file)

        assert config.certificate_identity == "user-subject-id"  # Falls back to sub

    def test_explicit_values_override_token_extraction(self, tmp_path):
        """Test explicit values override token extraction."""
        import base64
        import json

        header = base64.urlsafe_b64encode(b'{"alg":"RS256"}').decode().rstrip("=")
        payload_dict = {
            "iss": "https://issuer.example.com",
            "aud": "client-123",
            "email": "user@example.com",
        }
        payload = base64.urlsafe_b64encode(json.dumps(payload_dict).encode()).decode().rstrip("=")
        signature = base64.urlsafe_b64encode(b"fake-signature").decode().rstrip("=")
        token = f"{header}.{payload}.{signature}"

        token_file = tmp_path / "token"
        token_file.write_text(token)

        config = SigningConfig.create(
            identity_token_path=token_file,
            oidc_issuer="https://custom-issuer.example.com",
            client_id="custom-client",
            certificate_identity="custom@example.com",
        )

        # Explicit values should win
        assert config.oidc_issuer == "https://custom-issuer.example.com"
        assert config.client_id == "custom-client"
        assert config.certificate_identity == "custom@example.com"


    def test_k8s_token_uses_sub_for_certificate_identity(self, tmp_path):
        """Test Kubernetes service account token converts sub to certificate SAN format."""
        import base64
        import json

        header = base64.urlsafe_b64encode(b'{"alg":"RS256"}').decode().rstrip("=")
        payload_dict = {
            "iss": "https://kubernetes.default.svc.cluster.local",
            "aud": ["https://kubernetes.default.svc.cluster.local"],
            "sub": "system:serviceaccount:my-namespace:my-serviceaccount",
            "kubernetes.io": {
                "namespace": "my-namespace",
                "serviceaccount": {"name": "my-serviceaccount"},
            },
        }
        payload = base64.urlsafe_b64encode(json.dumps(payload_dict).encode()).decode().rstrip("=")
        signature = base64.urlsafe_b64encode(b"fake-signature").decode().rstrip("=")
        token = f"{header}.{payload}.{signature}"

        token_file = tmp_path / "token"
        token_file.write_text(token)

        config = SigningConfig.create(identity_token_path=token_file)

        # For K8s tokens, should convert sub to certificate SAN format
        assert config.certificate_identity == "https://kubernetes.io/namespaces/my-namespace/serviceaccounts/my-serviceaccount"
        assert config.oidc_issuer == "https://kubernetes.default.svc.cluster.local"

    def test_k8s_token_prefers_sub_over_email(self, tmp_path):
        """Test Kubernetes token prefers sub over email even if email is present."""
        import base64
        import json

        header = base64.urlsafe_b64encode(b'{"alg":"RS256"}').decode().rstrip("=")
        payload_dict = {
            "iss": "https://kubernetes.default.svc.cluster.local",
            "aud": "https://kubernetes.default.svc.cluster.local",
            "sub": "system:serviceaccount:my-namespace:my-serviceaccount",
            "email": "user@example.com",  # This should be ignored for K8s tokens
            "kubernetes.io": {
                "namespace": "my-namespace",
                "serviceaccount": {"name": "my-serviceaccount"},
            },
        }
        payload = base64.urlsafe_b64encode(json.dumps(payload_dict).encode()).decode().rstrip("=")
        signature = base64.urlsafe_b64encode(b"fake-signature").decode().rstrip("=")
        token = f"{header}.{payload}.{signature}"

        token_file = tmp_path / "token"
        token_file.write_text(token)

        config = SigningConfig.create(identity_token_path=token_file)

        # Should use converted sub, not email, for K8s tokens
        assert config.certificate_identity == "https://kubernetes.io/namespaces/my-namespace/serviceaccounts/my-serviceaccount"



class TestTokenUtils:
    """Test token utility functions."""

    def test_k8s_sub_to_certificate_identity(self):
        """Test converting Kubernetes sub claim to certificate SAN format."""
        from model_registry.signing.token import k8s_sub_to_certificate_identity

        result = k8s_sub_to_certificate_identity("system:serviceaccount:project2:wb2")
        assert result == "https://kubernetes.io/namespaces/project2/serviceaccounts/wb2"

    def test_k8s_sub_to_certificate_identity_with_different_namespace(self):
        """Test conversion with different namespace and service account."""
        from model_registry.signing.token import k8s_sub_to_certificate_identity

        result = k8s_sub_to_certificate_identity("system:serviceaccount:my-namespace:my-serviceaccount")
        assert result == "https://kubernetes.io/namespaces/my-namespace/serviceaccounts/my-serviceaccount"

    def test_k8s_sub_to_certificate_identity_invalid_format(self):
        """Test conversion with invalid format returns as-is."""
        from model_registry.signing.token import k8s_sub_to_certificate_identity

        # Not a K8s format
        result = k8s_sub_to_certificate_identity("user@example.com")
        assert result == "user@example.com"

        # Missing parts
        result = k8s_sub_to_certificate_identity("system:serviceaccount:namespace")
        assert result == "system:serviceaccount:namespace"
