"""Tests for signing utilities."""

from pathlib import Path

import pytest

from model_registry.signing import ImageSigner, InitializationError, SigningError, VerificationError


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
            tuf_url="https://tuf.example.com",
            root="https://tuf.example.com/root.json",
            root_checksum="abc123"
        )

        mock_runner.run.assert_called_once()
        cmd = mock_runner.run.call_args[0][0]
        assert cmd == [
            "cosign", "initialize",
            "--mirror", "https://tuf.example.com",
            "--root", "https://tuf.example.com/root.json",
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
        assert "https://default-tuf.example.com" in cmd
        assert "https://default-tuf.example.com/root.json" in cmd
        assert "def456" in cmd

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
        assert "https://override-tuf.example.com" in cmd
        assert "override-checksum" in cmd
        assert "https://default-tuf.example.com" not in cmd
        assert "default-checksum" not in cmd

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

        signer.sign(
            image="quay.io/example/image@sha256:abc123",
            identity_token_path=str(token_file),
            fulcio_url="https://fulcio.example.com",
            rekor_url="https://rekor.example.com"
        )

        mock_runner.run.assert_called_once()
        cmd = mock_runner.run.call_args[0][0]
        assert cmd == [
            "cosign", "sign", "-y",
            "--identity-token", str(token_file),
            "--fulcio-url", "https://fulcio.example.com",
            "--rekor-url", "https://rekor.example.com",
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
        assert str(token_file) in cmd
        assert "https://default-fulcio.example.com" in cmd
        assert "https://default-rekor.example.com" in cmd

    def test_sign_accepts_pathlike(self, tmp_path, mocker):
        """Test sign accepts PathLike objects."""
        # Create a temporary token file
        token_file = tmp_path / "token"
        token_file.write_text("test-token")

        signer = ImageSigner()
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        # Pass Path object instead of string
        signer.sign(
            image="quay.io/example/image@sha256:abc123",
            identity_token_path=Path(token_file),
            fulcio_url="https://fulcio.example.com",
        )

        mock_runner.run.assert_called_once()
        cmd = mock_runner.run.call_args[0][0]
        # os.fspath should have converted Path to string
        assert str(token_file) in cmd

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

    def test_verify_with_all_args(self, mocker):
        """Test verify with all arguments provided."""
        signer = ImageSigner()
        mock_runner = mocker.MagicMock()
        signer.runner = mock_runner

        signer.verify(
            image="quay.io/example/image@sha256:abc123",
            certificate_identity="https://kubernetes.io/namespaces/test/serviceaccounts/default",
            oidc_issuer="https://oidc.example.com"
        )

        mock_runner.run.assert_called_once()
        cmd = mock_runner.run.call_args[0][0]
        assert cmd == [
            "cosign", "verify",
            "--certificate-identity", "https://kubernetes.io/namespaces/test/serviceaccounts/default",
            "--certificate-oidc-issuer", "https://oidc.example.com",
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
        assert "https://kubernetes.io/namespaces/default/serviceaccounts/default" in cmd
        assert "https://default-oidc.example.com" in cmd

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
        import subprocess

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
        import subprocess

        token_file = tmp_path / "token"
        token_file.write_text("test-token")

        signer = ImageSigner()

        # Mock runner to raise CalledProcessError
        mock_runner = mocker.MagicMock()
        mock_runner.run.side_effect = subprocess.CalledProcessError(1, ["cosign", "sign"])
        signer.runner = mock_runner

        with pytest.raises(SigningError, match="Failed to sign image"):
            signer.sign(
                image="quay.io/example/image@sha256:abc123",
                identity_token_path=str(token_file)
            )

    def test_verify_raises_verification_error_on_command_failure(self, mocker):
        """Test verify raises VerificationError when cosign command fails."""
        import subprocess

        signer = ImageSigner()

        # Mock runner to raise CalledProcessError
        mock_runner = mocker.MagicMock()
        mock_runner.run.side_effect = subprocess.CalledProcessError(1, ["cosign", "verify"])
        signer.runner = mock_runner

        with pytest.raises(VerificationError, match="Failed to verify image"):
            signer.verify(
                image="quay.io/example/image@sha256:abc123",
                certificate_identity="https://kubernetes.io/namespaces/test/serviceaccounts/default"
            )
