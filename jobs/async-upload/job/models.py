"""
Configuration models for the async upload job using Pydantic for type safety and validation.
"""
from __future__ import annotations
from enum import StrEnum
import logging
from typing import Any
from pydantic import BaseModel, Field, model_validator, ConfigDict


class SourceType(StrEnum):
    """Supported source types for model artifacts."""
    S3 = "s3"
    OCI = "oci"
    URI = "uri"


class DestinationType(StrEnum):
    """Supported destination types for model artifacts."""
    S3 = "s3"
    OCI = "oci"


class S3Config(BaseModel):
    """S3 storage configuration."""
    bucket: str | None = None
    key: str | None = None
    region: str | None = None
    access_key_id: str | None = None
    secret_access_key: str | None = None
    endpoint_url: str | None = None

    @model_validator(mode='after')
    def validate_s3_required_fields(self) -> 'S3Config':
        """Validate that required S3 fields are present when S3 is used."""
        # This validation will be called by the parent models when needed
        return self


class OCIConfig(BaseModel):
    """OCI registry configuration."""
    uri: str | None = None
    registry: str | None = None
    username: str | None = None
    password: str | None = None
    email: str | None = None
    base_image: str = "busybox:latest"
    enable_tls_verify: bool = True

    @model_validator(mode='after')
    def validate_oci_credentials(self) -> 'OCIConfig':
        """Validate OCI credentials consistency."""
        if (self.username is None) != (self.password is None):
            raise ValueError("Both OCI username and password must be set together or both be None")
        return self


class SourceConfig(BaseModel):
    """Source configuration for model artifacts."""
    type: SourceType
    uri: str | None = None  # For URI type sources
    s3: S3Config = Field(default_factory=S3Config)
    oci: OCIConfig = Field(default_factory=OCIConfig)
    credentials_path: str | None = None

    @model_validator(mode='after')
    def validate_source_config(self) -> 'SourceConfig':
        """Validate source configuration based on type."""
        if self.type == SourceType.S3:
            if not all([self.s3.access_key_id, self.s3.secret_access_key, self.s3.bucket, self.s3.key]):
                raise ValueError("S3 credentials (access_key_id, secret_access_key), bucket, and key must be set for S3 sources")
        elif self.type == SourceType.OCI:
            if not all([self.oci.registry, self.oci.uri]):
                raise ValueError("OCI registry and URI must be set for OCI sources")
        elif self.type == SourceType.URI:
            if not self.uri:
                raise ValueError("URI must be set for URI type sources")
        return self


class DestinationConfig(BaseModel):
    """Destination configuration for model artifacts."""
    type: DestinationType
    s3: S3Config = Field(default_factory=S3Config)
    oci: OCIConfig = Field(default_factory=OCIConfig)
    credentials_path: str | None = None

    @model_validator(mode='after')
    def validate_destination_config(self) -> 'DestinationConfig':
        """Validate destination configuration based on type."""
        if self.type == DestinationType.S3:
            if not all([self.s3.access_key_id, self.s3.secret_access_key, self.s3.bucket, self.s3.key]):
                raise ValueError("S3 credentials (access_key_id, secret_access_key), bucket, and key must be set for S3 destinations")
        elif self.type == DestinationType.OCI:
            if not all([self.oci.registry, self.oci.uri]):
                raise ValueError("OCI registry and URI must be set for OCI destinations")
        return self


class ModelConfig(BaseModel):
    """Model registry model information."""
    id: str = Field(..., description="Model ID")
    version_id: str = Field(..., description="Model version ID")
    artifact_id: str = Field(..., description="Model artifact ID")

    @model_validator(mode='after')
    def validate_model_ids(self) -> 'ModelConfig':
        """Validate that all model IDs are provided."""
        if not all([self.id, self.version_id, self.artifact_id]):
            raise ValueError("Model ID, version ID and artifact ID must be set")
        return self


class StorageConfig(BaseModel):
    """Storage configuration for temporary files."""
    path: str = Field(default="/tmp/model-sync", description="Local storage path for temporary files")


class RegistryConfig(BaseModel):
    """Model registry client configuration."""
    server_address: str = Field(..., description="Model registry server address")
    port: int = Field(default=443, description="Model registry server port")
    is_secure: bool = Field(default=True, description="Use secure connection")
    author: str | None = Field(default=None, description="Author for model registration")
    user_token: str | None = Field(default=None, description="User authentication token")
    user_token_envvar: str | None = Field(default=None, description="Environment variable containing user token")
    custom_ca: str | None = Field(default=None, description="Custom CA certificate")
    custom_ca_envvar: str | None = Field(default=None, description="Environment variable containing custom CA")
    log_level: int = Field(default=logging.WARNING, description="Logging level for registry client")

    @model_validator(mode='after')
    def validate_registry_config(self) -> 'RegistryConfig':
        """Validate registry configuration."""
        if not self.server_address:
            raise ValueError("Registry server address must be set")
        return self


class AsyncUploadConfig(BaseModel):
    """Main configuration for the async upload job."""
    model_config = ConfigDict(
        # Allow extra fields for backward compatibility
        extra='forbid',
        # Validate assignments
        validate_assignment=True,
        # Populate by name (allows both snake_case and field names)
        populate_by_name=True
    )

    source: SourceConfig
    destination: DestinationConfig
    model: ModelConfig
    storage: StorageConfig = Field(default_factory=StorageConfig)
    registry: RegistryConfig

    def to_dict(self) -> dict[str, Any]:
        """Convert to dictionary for backward compatibility."""
        return self.model_dump()

    def __getitem__(self, key: str) -> Any:
        """Allow dictionary-style access for backward compatibility."""
        attr = getattr(self, key)
        # If the attribute is a Pydantic model, wrap it to support dict-style access
        if isinstance(attr, BaseModel):
            return DictCompatibilityWrapper(attr)
        return attr

    def get(self, key: str, default: Any = None) -> Any:
        """Allow dictionary-style get for backward compatibility."""
        try:
            return self[key]
        except (AttributeError, KeyError):
            return default


class DictCompatibilityWrapper:
    """Wrapper to make Pydantic models compatible with dictionary-style access."""
    
    def __init__(self, model: BaseModel):
        self._model = model
    
    def __getitem__(self, key: str) -> Any:
        attr = getattr(self._model, key)
        # If the attribute is a Pydantic model, wrap it to support dict-style access
        if isinstance(attr, BaseModel):
            return DictCompatibilityWrapper(attr)
        return attr
    
    def get(self, key: str, default: Any = None) -> Any:
        try:
            return self[key]
        except (AttributeError, KeyError):
            return default
    
    def __getattr__(self, name: str) -> Any:
        """Allow direct attribute access too."""
        return getattr(self._model, name)
    
    def __eq__(self, other) -> bool:
        """Support equality checks."""
        if hasattr(other, '_model'):
            return self._model == other._model
        return getattr(self._model, 'value', self._model) == other
    
    def __str__(self) -> str:
        """Support string representation."""
        return str(self._model)
    
    def __repr__(self) -> str:
        """Support repr."""
        return repr(self._model)


# Type alias for backward compatibility
Config = AsyncUploadConfig 