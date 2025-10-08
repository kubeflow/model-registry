"""MLflow AbstractStore implementation for Model Registry."""

import asyncio
import logging
import os
from urllib.parse import parse_qs, urlparse

from mlflow.entities.model_registry import (
    ModelVersion,
    ModelVersionTag,
    RegisteredModel,
    RegisteredModelTag,
)
from mlflow.exceptions import MlflowException
from mlflow.protos.databricks_pb2 import (
    INVALID_PARAMETER_VALUE,
    RESOURCE_ALREADY_EXISTS,
    RESOURCE_DOES_NOT_EXIST,
)
from mlflow.entities.model_registry.model_version_status import ModelVersionStatus
from mlflow.store.entities.paged_list import PagedList
from mlflow.store.model_registry.abstract_store import AbstractStore

from model_registry import ModelRegistry
from model_registry.exceptions import StoreError
from model_registry.types import RegisteredModelState, ModelVersionState
from mr_openapi.models.artifact_state import ArtifactState

logger = logging.getLogger(__name__)


class ModelRegistryStore(AbstractStore):
    """
    MLflow AbstractStore implementation that uses Model Registry as the backend.

    This allows MLflow model registry operations to be backed by a Model Registry server.
    """

    def __init__(self, store_uri=None, tracking_uri=None):
        super().__init__(store_uri, tracking_uri)

        if not store_uri:
            raise MlflowException("Store URI is required", INVALID_PARAMETER_VALUE)

        parsed_uri = urlparse(store_uri)

        if parsed_uri.scheme == "modelregistry":
            self.is_secure = os.getenv("MODEL_REGISTRY_SECURE", "true").lower() == "true"
        elif parsed_uri.scheme == "modelregistry+https":
            self.is_secure = True
        elif parsed_uri.scheme == "modelregistry+http":
            self.is_secure = False
        else:
            raise MlflowException(
                f"Invalid store URI scheme: {parsed_uri.scheme}. Expected: modelregistry, modelregistry+http, or modelregistry+https",
                INVALID_PARAMETER_VALUE,
            )

        hostname = parsed_uri.hostname or os.getenv("MODEL_REGISTRY_HOST", "localhost")
        self.port = parsed_uri.port or int(os.getenv("MODEL_REGISTRY_PORT", "443" if self.is_secure else "80"))

        query_params = parse_qs(parsed_uri.query)
        self.author = query_params.get("author", [os.getenv("MODEL_REGISTRY_AUTHOR", "unknown")])[0]

        if "is-secure" in query_params:
            self.is_secure = query_params["is-secure"][0].lower() == "true"

        self.user_token = query_params.get("user-token", [os.getenv("MODEL_REGISTRY_TOKEN")])[0]
        self.custom_ca = query_params.get("custom-ca", [os.getenv("MODEL_REGISTRY_CA")])[0]

        protocol = "https" if self.is_secure else "http"
        self.server_address = f"{protocol}://{hostname}"

        try:
            self._client = ModelRegistry(
                server_address=self.server_address,
                port=self.port,
                author=self.author,
                is_secure=self.is_secure,
                user_token=self.user_token,
                custom_ca=self.custom_ca,
            )
        except Exception as e:
            raise MlflowException(f"Failed to connect to Model Registry: {e}", INVALID_PARAMETER_VALUE) from e

    def _convert_mr_to_mlflow_registered_model(self, mr_model):
        if mr_model is None:
            return None

        tags = []
        if hasattr(mr_model, 'custom_properties') and mr_model.custom_properties:
            for key, value in mr_model.custom_properties.items():
                tags.append(RegisteredModelTag(key=key, value=str(value)))

        return RegisteredModel(
            name=mr_model.name,
            creation_timestamp=mr_model.create_time_since_epoch,
            last_updated_timestamp=mr_model.last_update_time_since_epoch,
            description=mr_model.description,
            tags=tags,
        )

    def _convert_mr_to_mlflow_model_version(self, mr_version, mr_artifact=None, registered_model_name=None):
        if mr_version is None:
            return None

        tags = []
        current_stage = None
        if hasattr(mr_version, 'custom_properties') and mr_version.custom_properties:
            for key, value in mr_version.custom_properties.items():
                if key == "mlflow.stage":
                    current_stage = str(value)
                tags.append(ModelVersionTag(key=key, value=str(value)))

        source = None
        status = ModelVersionStatus.to_string(ModelVersionStatus.READY)
        if mr_artifact:
            source = mr_artifact.uri
            if hasattr(mr_artifact, 'state') and mr_artifact.state:
                state_mapping = {
                    ArtifactState.PENDING: ModelVersionStatus.PENDING_REGISTRATION,
                    ArtifactState.LIVE: ModelVersionStatus.READY,
                    ArtifactState.ABANDONED: ModelVersionStatus.FAILED_REGISTRATION,
                    ArtifactState.DELETED: ModelVersionStatus.FAILED_REGISTRATION,
                    ArtifactState.MARKED_FOR_DELETION: ModelVersionStatus.FAILED_REGISTRATION,
                }
                status = ModelVersionStatus.to_string(state_mapping.get(mr_artifact.state, ModelVersionStatus.READY))

        return ModelVersion(
            name=registered_model_name,
            version=mr_version.name,
            creation_timestamp=mr_version.create_time_since_epoch,
            last_updated_timestamp=mr_version.last_update_time_since_epoch,
            description=mr_version.description,
            source=source,
            tags=tags,
            status=status,
            current_stage=current_stage,
        )

    def create_registered_model(self, name, tags=None, description=None, deployment_job_id=None):
        try:
            metadata = {}
            if tags:
                for tag in tags:
                    metadata[tag.key] = tag.value
            existing_model = self._client.get_registered_model(name)
            if existing_model:
                raise MlflowException(
                    f"Registered model with name '{name}' already exists",
                    RESOURCE_ALREADY_EXISTS,
                )
            mr_model = self._client.async_runner(
                self._client._register_model(name, description=description, **metadata)
            )

            return self._convert_mr_to_mlflow_registered_model(mr_model)

        except StoreError as e:
            if "already exists" in str(e):
                raise MlflowException(
                    f"Registered model with name '{name}' already exists",
                    RESOURCE_ALREADY_EXISTS,
                ) from e
            raise MlflowException(f"Failed to create registered model: {e}") from e

    def update_registered_model(self, name, description, deployment_job_id=None):
        try:
            mr_model = self._client.get_registered_model(name)
            if not mr_model:
                raise MlflowException(
                    f"Registered model with name '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            mr_model.description = description
            updated_model = self._client.update(mr_model)

            return self._convert_mr_to_mlflow_registered_model(updated_model)

        except StoreError as e:
            raise MlflowException(f"Failed to update registered model: {e}") from e

    def rename_registered_model(self, name, new_name):
        # Model Registry doesn't support renaming models directly
        # We'd need to create a new model and copy all versions
        raise MlflowException(
            "Renaming registered models is not supported by Model Registry",
            INVALID_PARAMETER_VALUE,
        )

    def delete_registered_model(self, name):
        try:
            mr_model = self._client.get_registered_model(name)
            if not mr_model:
                raise MlflowException(
                    f"Registered model with name '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            # Archive all model versions first
            versions_pager = self._client.get_model_versions(name)
            for version in versions_pager:
                if hasattr(version, 'state') and version.state != ModelVersionState.ARCHIVED:
                    version.state = ModelVersionState.ARCHIVED
                    self._client.update(version)

            mr_model.state = RegisteredModelState.ARCHIVED
            self._client.update(mr_model)

        except StoreError as e:
            raise MlflowException(f"Failed to delete registered model: {e}") from e

    def search_registered_models(
        self, filter_string=None, max_results=None, order_by=None, page_token=None
    ):
        try:
            models_pager = self._client.get_registered_models()

            if max_results:
                models_pager = models_pager.page_size(max_results)

            models = []
            for model in models_pager:
                if hasattr(model, 'state') and model.state == RegisteredModelState.ARCHIVED:
                    continue
                models.append(model)
                if max_results and len(models) >= max_results:
                    break
            mlflow_models = []
            for mr_model in models:
                mlflow_model = self._convert_mr_to_mlflow_registered_model(mr_model)
                if mlflow_model:
                    mlflow_models.append(mlflow_model)

            return PagedList(mlflow_models, None)

        except StoreError as e:
            raise MlflowException(f"Failed to search registered models: {e}") from e

    def get_registered_model(self, name):
        try:
            mr_model = self._client.get_registered_model(name)
            if not mr_model:
                raise MlflowException(
                    f"Registered model with name '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            if hasattr(mr_model, 'state') and mr_model.state == RegisteredModelState.ARCHIVED:
                raise MlflowException(
                    f"Registered model with name '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            return self._convert_mr_to_mlflow_registered_model(mr_model)
        except StoreError as e:
            raise MlflowException(f"Failed to get registered model: {e}") from e

    def get_latest_versions(self, name, stages=None):
        try:
            # Get all versions for the model
            versions_pager = self._client.get_model_versions(name)
            all_versions = []

            for version in versions_pager:
                if hasattr(version, 'state') and version.state == ModelVersionState.ARCHIVED:
                    continue
                all_versions.append(version)

            # Since Model Registry doesn't have stages concept, return the latest version
            # TODO: Implement stage concept using tags or metadata
            if all_versions:
                latest_version = max(all_versions, key=lambda v: v.create_time_since_epoch)

                # Get corresponding artifact
                mr_artifact = self._client.get_model_artifact(name, latest_version.name)
                mlflow_version = self._convert_mr_to_mlflow_model_version(latest_version, mr_artifact, name)
                if mlflow_version:
                    return [mlflow_version]

            return []

        except StoreError as e:
            raise MlflowException(f"Failed to get latest versions: {e}") from e

    def set_registered_model_tag(self, name, tag):
        try:
            mr_model = self._client.get_registered_model(name)
            if not mr_model:
                raise MlflowException(
                    f"Registered model with name '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            # Add tag to custom properties
            if not hasattr(mr_model, 'custom_properties') or not mr_model.custom_properties:
                mr_model.custom_properties = {}
            mr_model.custom_properties[tag.key] = tag.value

            self._client.update(mr_model)

        except StoreError as e:
            raise MlflowException(f"Failed to set registered model tag: {e}") from e

    def delete_registered_model_tag(self, name, key):
        try:
            mr_model = self._client.get_registered_model(name)
            if not mr_model:
                raise MlflowException(
                    f"Registered model with name '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            # Remove tag from custom properties
            if hasattr(mr_model, 'custom_properties') and mr_model.custom_properties:
                mr_model.custom_properties.pop(key, None)
                self._client.update(mr_model)

        except StoreError as e:
            raise MlflowException(f"Failed to delete registered model tag: {e}") from e

    # ModelVersion CRUD operations

    def create_model_version(
        self,
        name,
        source,
        run_id=None,
        tags=None,
        run_link=None,
        description=None,
        local_model_path=None,
        model_id: str | None = None,
    ):
        """Create a new model version."""
        try:
            # Ensure the registered model exists
            mr_model = self._client.get_registered_model(name)
            if not mr_model:
                # Create the registered model first
                mr_model = self._client.async_runner(
                    self._client._register_model(name)
                )

            # Generate a version name (Model Registry uses version names, not numbers)
            # For MLflow compatibility, we'll use sequential integers as strings
            # Get existing versions to determine next version number
            try:
                versions_pager = self._client.get_model_versions(name)
                existing_versions = []

                # Use proper iterator protocol which includes built-in loop detection
                for version in versions_pager:
                    existing_versions.append(version)

                # Extract version numbers and find the next one
                version_numbers = []
                for v in existing_versions:
                    try:
                        # Try to parse version name as integer
                        version_numbers.append(int(v.name))
                    except ValueError:
                        # If it's not an integer, skip it
                        pass

                next_version = max(version_numbers, default=0) + 1
                version_name = str(next_version)

            except Exception:
                # Fallback to simple sequential numbering
                version_name = "1"

            # Convert MLflow tags to metadata
            metadata = {}
            if tags:
                for tag in tags:
                    metadata[tag.key] = tag.value

            # Add MLflow-specific metadata
            if run_id:
                metadata["mlflow.run_id"] = run_id
            if run_link:
                metadata["mlflow.run_link"] = run_link
            if model_id:
                metadata["mlflow.model_id"] = model_id

            # Create the model version
            mr_version = self._client.async_runner(
                self._client._register_new_version(
                    mr_model,
                    version_name,
                    self.author,
                    description=description,
                    **metadata
                )
            )

            # Create the model artifact
            # Extract format info from source or use defaults
            mr_artifact = self._client.async_runner(
                self._client._register_model_artifact(
                    mr_version,
                    name,  # artifact name same as model name
                    source,
                    model_format_name="mlflow-model",  # Default format
                    model_format_version="1.0",  # Default version
                )
            )

            return self._convert_mr_to_mlflow_model_version(mr_version, mr_artifact, name)

        except StoreError as e:
            raise MlflowException(f"Failed to create model version: {e}") from e

    def update_model_version(self, name, version, description):
        try:
            mr_version = self._client.get_model_version(name, version)
            if not mr_version:
                raise MlflowException(
                    f"Model version '{version}' for model '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            # Update description
            mr_version.description = description
            updated_version = self._client.update(mr_version)

            mr_artifact = self._client.get_model_artifact(name, version)

            return self._convert_mr_to_mlflow_model_version(updated_version, mr_artifact)

        except StoreError as e:
            raise MlflowException(f"Failed to update model version: {e}") from e

    def transition_model_version_stage(self, name, version, stage, archive_existing_versions):
        # Model Registry doesn't have built-in stage concept
        # We'll implement this using tags
        try:
            mr_version = self._client.get_model_version(name, version)
            if not mr_version:
                raise MlflowException(
                    f"Model version '{version}' for model '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            # Add stage as a custom property
            if not hasattr(mr_version, 'custom_properties') or not mr_version.custom_properties:
                mr_version.custom_properties = {}
            mr_version.custom_properties["mlflow.stage"] = stage

            # Handle archive_existing_versions
            if archive_existing_versions and stage in ["staging", "production"]:
                # Get all versions and archive those in the same stage
                versions_pager = self._client.get_model_versions(name)
                all_versions = []

                # Use proper iterator protocol which includes built-in loop detection
                for other_version in versions_pager:
                    all_versions.append(other_version)

                # Archive existing versions in the same stage
                for other_version in all_versions:
                    if (other_version.name != version and
                        hasattr(other_version, 'custom_properties') and
                        other_version.custom_properties and
                        other_version.custom_properties.get("mlflow.stage") == stage):
                        other_version.custom_properties["mlflow.stage"] = "archived"
                        self._client.update(other_version)

            updated_version = self._client.update(mr_version)

            mr_artifact = self._client.get_model_artifact(name, version)

            return self._convert_mr_to_mlflow_model_version(updated_version, mr_artifact)

        except StoreError as e:
            raise MlflowException(f"Failed to transition model version stage: {e}") from e

    def delete_model_version(self, name, version):
        try:
            mr_version = self._client.get_model_version(name, version)
            if not mr_version:
                raise MlflowException(
                    f"Model version '{version}' for model '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            mr_version.state = ModelVersionState.ARCHIVED
            self._client.update(mr_version)

        except StoreError as e:
            raise MlflowException(f"Failed to delete model version: {e}") from e

    def get_model_version(self, name, version):
        try:
            mr_version = self._client.get_model_version(name, version)
            if not mr_version:
                return None

            if hasattr(mr_version, 'state') and mr_version.state == ModelVersionState.ARCHIVED:
                return None

            mr_artifact = self._client.get_model_artifact(name, version)

            return self._convert_mr_to_mlflow_model_version(mr_version, mr_artifact, name)

        except StoreError as e:
            raise MlflowException(f"Failed to get model version: {e}") from e

    def get_model_version_download_uri(self, name, version):
        try:
            # Check if version is archived before returning download URI
            mr_version = self._client.get_model_version(name, version)
            if not mr_version:
                raise MlflowException(
                    f"Model version '{version}' of model '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            if hasattr(mr_version, 'state') and mr_version.state == ModelVersionState.ARCHIVED:
                raise MlflowException(
                    f"Model version '{version}' of model '{name}' has been deleted",
                    RESOURCE_DOES_NOT_EXIST,
                )

            mr_artifact = self._client.get_model_artifact(name, version)
            if not mr_artifact:
                raise MlflowException(
                    f"Model artifact for version '{version}' of model '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            return mr_artifact.uri

        except StoreError as e:
            raise MlflowException(f"Failed to get model version download URI: {e}") from e

    def search_model_versions(
        self, filter_string=None, max_results=None, order_by=None, page_token=None
    ):
        try:
            # Model Registry doesn't support cross-model version search
            # This would require iterating through all models
            # For now, return empty list
            # TODO: Implement proper search across all models
            return PagedList([], None)

        except StoreError as e:
            raise MlflowException(f"Failed to search model versions: {e}") from e

    def set_model_version_tag(self, name, version, tag):
        try:
            mr_version = self._client.get_model_version(name, version)
            if not mr_version:
                raise MlflowException(
                    f"Model version '{version}' for model '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            # Add tag to custom properties
            if not hasattr(mr_version, 'custom_properties') or not mr_version.custom_properties:
                mr_version.custom_properties = {}
            mr_version.custom_properties[tag.key] = tag.value

            self._client.update(mr_version)

        except StoreError as e:
            raise MlflowException(f"Failed to set model version tag: {e}") from e

    def delete_model_version_tag(self, name, version, key):
        try:
            mr_version = self._client.get_model_version(name, version)
            if not mr_version:
                raise MlflowException(
                    f"Model version '{version}' for model '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            # Remove tag from custom properties
            if hasattr(mr_version, 'custom_properties') and mr_version.custom_properties:
                mr_version.custom_properties.pop(key, None)
                self._client.update(mr_version)

        except StoreError as e:
            raise MlflowException(f"Failed to delete model version tag: {e}") from e

    # Alias operations

    def set_registered_model_alias(self, name, alias, version):
        # Model Registry doesn't have built-in alias support
        # We'll implement this using tags on the registered model
        try:
            mr_model = self._client.get_registered_model(name)
            if not mr_model:
                raise MlflowException(
                    f"Registered model with name '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            # Convert version to string if it's an integer
            version_str = str(version)

            # Verify the version exists
            mr_version = self._client.get_model_version(name, version_str)
            if not mr_version:
                raise MlflowException(
                    f"Model version '{version_str}' for model '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            # Add alias as a custom property
            if not hasattr(mr_model, 'custom_properties') or not mr_model.custom_properties:
                mr_model.custom_properties = {}
            mr_model.custom_properties[f"mlflow.alias.{alias}"] = version_str

            self._client.update(mr_model)

        except StoreError as e:
            raise MlflowException(f"Failed to set registered model alias: {e}") from e

    def delete_registered_model_alias(self, name, alias):
        try:
            mr_model = self._client.get_registered_model(name)
            if not mr_model:
                raise MlflowException(
                    f"Registered model with name '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            # Remove alias from custom properties
            if hasattr(mr_model, 'custom_properties') and mr_model.custom_properties:
                mr_model.custom_properties.pop(f"mlflow.alias.{alias}", None)
                self._client.update(mr_model)

        except StoreError as e:
            raise MlflowException(f"Failed to delete registered model alias: {e}") from e

    def get_model_version_by_alias(self, name, alias):
        try:
            mr_model = self._client.get_registered_model(name)
            if not mr_model:
                raise MlflowException(
                    f"Registered model with name '{name}' not found",
                    RESOURCE_DOES_NOT_EXIST,
                )

            # Get version from alias
            if (hasattr(mr_model, 'custom_properties') and
                mr_model.custom_properties and
                f"mlflow.alias.{alias}" in mr_model.custom_properties):
                version = mr_model.custom_properties[f"mlflow.alias.{alias}"]
                return self.get_model_version(name, version)

            raise MlflowException(
                f"Alias '{alias}' not found for model '{name}'",
                RESOURCE_DOES_NOT_EXIST,
            )

        except StoreError as e:
            raise MlflowException(f"Failed to get model version by alias: {e}") from e

    def _await_model_version_creation(self, mv, await_creation_for):
        """
        Override MLflow's default polling behavior for model version creation.

        Since Model Registry creates versions synchronously and our plugin always
        returns status "READY", we don't need to poll. We just verify the version
        exists and is accessible.
        """
        # Try to get the model version to ensure it exists
        try:
            retrieved_mv = self.get_model_version(mv.name, mv.version)
            if retrieved_mv is None:
                raise MlflowException(
                    f"Model version {mv.version} for model {mv.name} was not found after creation",
                    RESOURCE_DOES_NOT_EXIST,
                )
            # If we got here, the version exists and is ready - no need to poll

        except Exception as e:
            raise MlflowException(
                f"Failed to verify model version creation for {mv.name} version {mv.version}: {e}"
            ) from e
