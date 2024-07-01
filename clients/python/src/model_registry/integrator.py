from abc import ABC, abstractmethod
from enum import Enum
import re
from typing import Any, Dict
from .exceptions import StoreException
from warnings import warn


class supportedPlatform(Enum):
    HUGGING_FACE = "HuggingFace"
    MLFLOW = "MlFlow"

    def __contains__(cls, item: str):
        try:
            cls(item)
        except ValueError:
            return False
        return True


class ModelInfoProvider(ABC):
    @property
    def passByInfo(self):
        return ["version", "model_format_name", "model_format_version", "description"]

    @property
    def modelInfo(self):
        raise NotImplementedError("This is an abstract class")

    @abstractmethod
    def get_model_info(self) -> dict:
        pass

    def _extractDirectPassByInfo(self, kwargs: dict) -> dict:
        for keyword in self.passByInfo:
            self._modelInfo[keyword] = kwargs[keyword]
            del kwargs[keyword]
        return kwargs


class HuggingFaceModelInfoProvider(ModelInfoProvider):
    _modelInfo = {}

    def get_model_info(self, **kwargs) -> dict:
        kwargs = self._extractDirectPassByInfo(kwargs)
        return self._getHfModelInfo(**kwargs)

    def _getHfModelInfo(
        self, repo: str, path: str, *, author: str = None, model_name: str = None, git_ref: str = "main"
    ) -> dict:
        try:
            from huggingface_hub import HfApi, hf_hub_url, utils
        except ImportError as e:
            msg = "huggingface_hub is not installed"
            raise StoreException(msg) from e

        api = HfApi()
        try:
            model_info = api.model_info(repo, revision=git_ref)
        except utils.RepositoryNotFoundError as e:
            msg = f"Repository {repo} does not exist"
            raise StoreException(msg) from e

        if not author:
            author = model_info.author or "unknown"
            
        source_uri = hf_hub_url(repo, path, revision=git_ref)
        metadata = {
            "repo": repo,
            "source_uri": source_uri,
            "model_origin": "huggingface_hub",
            "model_author": author,
        }
        
        self._modelInfo.update(
            {
                "repo": repo,
                "path": path,
                "author": author,
                "model_name": model_name or repo,
                "git_ref": git_ref,
                "source_uri": hf_hub_url(repo, path, revision=git_ref),
                "metadata": metadata,
                "storage_path": path
            }
        )

        return self._modelInfo


class MLflowModelInfoProvider(ModelInfoProvider):
    _modelInfo = {}

    def get_model_info(self, **kwargs) -> dict:
        kwargs = self._extractDirectPassByInfo(kwargs)
        return self._getMlflowModelInfo(**kwargs)

    def _validateArtifactUri(self, mlflowSourceUri: str) -> str:
        prefixes = ["gs://", "s3://", "hdfs://"]

        regex_patterns = [
            r"https://(.+?).blob.core.windows.net/(.+)",
            r"https://(.+?).file.core.windows.net/(.+)",
            r"https?://(.+)/(.+)",
        ]

        for prefix in prefixes:
            if mlflowSourceUri.startswith(prefix):
                return mlflowSourceUri

        for pattern in regex_patterns:
            if re.match(pattern, mlflowSourceUri):
                return mlflowSourceUri

        warn("Unsupported artifact uri format, may not accessible by downstream application", stacklevel=2)
        return mlflowSourceUri

    def _getMlflowModelInfo(self, tracking_uri: str, registered_name: str, registered_version: str, model_name: str, author: str = None) -> dict:
        try:
            import mlflow
            from mlflow.exceptions import RestException
        except ImportError as e:
            msg = "mlflow is not installed"
            raise StoreException(msg) from e
        mlflow.set_tracking_uri(tracking_uri)
        client = mlflow.tracking.MlflowClient()
        try:
            if registered_version:
                model_version_details = client.get_model_version(name=registered_name, version=registered_version)
            else:
                model_version_details = client.get_registered_model(registered_name).latest_versions[0]
        except RestException as e:
            msg = f"Mlflow API error: {e.message}"
            raise StoreException(msg)

        if not author:
            if "author" in model_version_details.tags:
                self._modelInfo["author"] = model_version_details.tags["author"]
            else:
                self._modelInfo["author"] = "unknown"
        else:
            self._modelInfo["author"] = author

        sourceUri = self._validateArtifactUri(model_version_details.source)
        model_file_name = mlflow.pyfunc.load_model(sourceUri).metadata.flavors["python_function"]["data"]
        self._modelInfo.update(
            {
                "uri": f"{sourceUri}/{model_file_name}",
                "metadata": {
                    "name_in_origin_mr": model_version_details.name,
                    "version_in_origin_mr": model_version_details.version
                    "tracking_uri": tracking_uri,
                    "model_origin": "MlFlow",
                    "model_author": self._modelInfo["author"],
                },
            }
        )

        self._modelInfo.update(model_version_details.tags)
        self._modelInfo["name"] = model_name or model_version_details.name
        return self._modelInfo


class ModelInfoManager:

    @classmethod
    def providers(cls):
        return {
            supportedPlatform.HUGGING_FACE: HuggingFaceModelInfoProvider(),
            supportedPlatform.MLFLOW: MLflowModelInfoProvider(),
        }

    @classmethod
    def get_model_info(cls, platform: str, kwargs: dict) -> dict:
        assert platform in [p.value for p in supportedPlatform]
        provider = cls.providers().get(supportedPlatform(platform))
        if provider:
            return provider.get_model_info(**kwargs)
        else:
            return {"error": "Platform not supported"}
