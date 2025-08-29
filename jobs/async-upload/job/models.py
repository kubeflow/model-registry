"""
Configuration models for the async upload job using Pydantic for type safety and validation.
"""
from __future__ import annotations
from enum import StrEnum
import logging
from typing import Any, Union
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


class BaseStorageConfig(BaseModel):
    """Base configuration for storage types."""
    credentials_path: str | None = None


class S3Config(BaseModel):
    """Basic S3 storage configuration. To be used as an intermediary model for the S3StorageConfig, allowing for additional values to be overlaid until the final config is created and validated via the S3StorageConfig model."""
    bucket: str | None = None
    key: str | None = None # 'path' in bucket
    region: str | None = None
    access_key_id: str | None = None
    secret_access_key: str | None = None
    endpoint: str | None = None


class OCIConfig(BaseModel):
    """Basic OCI registry configuration. To be used as an intermediary model for the OCIStorageConfig, allowing for additional values to be overlaid until the final config is created and validated via the OCIStorageConfig model."""
    uri: str
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


class S3StorageConfig(BaseStorageConfig, S3Config):
    """S3 storage configuration with validation - can be used for both source and destination."""
    
    @model_validator(mode='after')
    def validate_s3_storage(self) -> 'S3StorageConfig':
        """Validate that required S3 fields are present."""
        if not all([self.access_key_id, self.secret_access_key, self.bucket, self.key]):
            raise ValueError("S3 credentials (access_key_id, secret_access_key), bucket, and key must be set")
        return self


class OCIStorageConfig(BaseStorageConfig, OCIConfig):
    """OCI storage configuration with validation - can be used for both source and destination."""
    
    @model_validator(mode='after')
    def validate_oci_storage(self) -> 'OCIStorageConfig':
        """Validate that required OCI fields are present."""
        if not all([self.registry, self.uri]):
            raise ValueError("OCI registry and URI must be set")
        return self


class URISourceConfig(BaseModel):
    """Basic URI source configuration. To be used as an intermediary model for the URISourceStorageConfig, allowing for additional values to be overlaid until the final config is created and validated via the URISourceStorageConfig model."""
    uri: str | None = None


class URISourceStorageConfig(BaseStorageConfig, URISourceConfig):
    """URI source storage configuration with validation - only used for sources, not destinations."""
    
    @model_validator(mode='after')
    def validate_uri_storage(self) -> 'URISourceStorageConfig':
        """Validate that required URI field is present."""
        if not self.uri:
            raise ValueError("URI must be set for URI type sources")
        return self


# Union types for source and destination configurations - this enables isinstance() checks
SourceConfig = Union[S3StorageConfig, OCIStorageConfig, URISourceStorageConfig]
DestinationConfig = Union[S3StorageConfig, OCIStorageConfig]


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
