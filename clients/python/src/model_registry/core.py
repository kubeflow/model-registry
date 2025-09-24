"""Client for the model registry."""

from __future__ import annotations

from collections.abc import AsyncIterator
from contextlib import asynccontextmanager
from dataclasses import dataclass
from typing import TypeVar, cast, overload

from mr_openapi import (
    ApiClient,
    Configuration,
    ModelRegistryServiceApi,
)
from mr_openapi import (
    exceptions as mr_exceptions,
)

from ._utils import required_args
from .types import (
    Artifact,
    Experiment,
    ExperimentRun,
    ExperimentRunArtifact,
    ListOptions,
    ModelArtifact,
    ModelVersion,
    RegisteredModel,
)

ArtifactT = TypeVar("ArtifactT", bound=Artifact)


@dataclass
class ModelRegistryAPIClient:
    """Model registry API."""

    config: Configuration

    @classmethod
    def secure_connection(
        cls,
        server_address: str,
        port: int = 443,
        *,
        user_token: str,
        custom_ca: str | None = None,
    ) -> ModelRegistryAPIClient:
        """Constructor.

        Args:
            server_address: Server address.
            port: Server port. Defaults to 443.

        Keyword Args:
            user_token: The PEM-encoded user token as a string.
            custom_ca: The path to a PEM-
        """
        return cls(
            Configuration(
                f"{server_address}:{port}",
                access_token=user_token,
                ssl_ca_cert=custom_ca,
            )
        )

    @classmethod
    def insecure_connection(
        cls,
        server_address: str,
        port: int,
        user_token: str | None = None,
    ) -> ModelRegistryAPIClient:
        """Constructor.

        Args:
            server_address: Server address.
            port: Server port.
            user_token: The PEM-encoded user token as a string.
        """
        config = Configuration(
            host=f"{server_address}:{port}",
            access_token=user_token,
            verify_ssl=False,
        )
        return cls(config)

    @asynccontextmanager
    async def get_client(self) -> AsyncIterator[ModelRegistryServiceApi]:
        """Get a client for the model registry."""
        api_client = ApiClient(self.config)
        client = ModelRegistryServiceApi(api_client)

        try:
            yield client
        finally:
            await api_client.close()

    async def upsert_registered_model(self, registered_model: RegisteredModel) -> RegisteredModel:
        """Upsert a registered model.

        Updates or creates a registered model on the server.

        Args:
            registered_model: Registered model.

        Returns:
            New registered model.
        """
        async with self.get_client() as client:
            if registered_model.id:
                rm = await client.update_registered_model(registered_model.id, registered_model.update())
            else:
                rm = await client.create_registered_model(registered_model.create())

        return RegisteredModel.from_basemodel(rm)

    async def get_registered_model_by_id(self, id: str) -> RegisteredModel | None:
        """Fetch a registered model by its ID.

        Args:
            id: Registered model ID.

        Returns:
            Registered model.
        """
        async with self.get_client() as client:
            try:
                rm = await client.get_registered_model(id)
            except mr_exceptions.NotFoundException:
                return None

        return RegisteredModel.from_basemodel(rm)

    @overload
    async def get_registered_model_by_params(self, name: str): ...

    @overload
    async def get_registered_model_by_params(self, *, external_id: str): ...

    @required_args(("name",), ("external_id",))
    async def get_registered_model_by_params(
        self, name: str | None = None, external_id: str | None = None
    ) -> RegisteredModel | None:
        """Fetch a registered model by its name or external ID.

        Args:
            name: Registered model name.
            external_id: Registered model external ID.

        Returns:
            Registered model.
        """
        async with self.get_client() as client:
            try:
                rm = await client.find_registered_model(name=name, external_id=external_id)
            except mr_exceptions.NotFoundException:
                return None

        return RegisteredModel.from_basemodel(rm)

    async def get_registered_models(self, options: ListOptions | None = None) -> list[RegisteredModel]:
        """Fetch registered models.

        Args:
            options: Options for listing registered models.

        Returns:
            Registered models.
        """
        async with self.get_client() as client:
            rm_list = await client.get_registered_models(**(options or ListOptions()).as_options())

        if options:
            options.next_page_token = rm_list.next_page_token

        return [RegisteredModel.from_basemodel(rm) for rm in rm_list.items or []]

    async def upsert_model_version(
        self, model_version: ModelVersion, registered_model_id: str | None = None
    ) -> ModelVersion:
        """Upsert a model version.

        Updates or creates a model version on the server.

        Args:
            model_version: Model version to upsert.
            registered_model_id: ID of the registered model this version will be associated to. Can be None when updating an existing model version.

        Returns:
            New model version.
        """
        async with self.get_client() as client:
            if model_version.id:
                mv = await client.update_model_version(model_version.id, model_version.update())
            elif registered_model_id:
                mv = await client.create_model_version(model_version.create(registered_model_id=registered_model_id))
            else:
                msg = f"Registered model ID required for creating a new model version: {model_version}"
                raise ValueError(msg)

        return ModelVersion.from_basemodel(mv)

    async def get_model_version_by_id(self, model_version_id: str) -> ModelVersion | None:
        """Fetch a model version by its ID.

        Args:
            model_version_id: Model version ID.

        Returns:
            Model version.
        """
        async with self.get_client() as client:
            try:
                mv = await client.get_model_version(model_version_id)
            except mr_exceptions.NotFoundException:
                return None

        return ModelVersion.from_basemodel(mv)

    async def get_model_versions(
        self, registered_model_id: str, options: ListOptions | None = None
    ) -> list[ModelVersion]:
        """Fetch model versions by registered model ID.

        Args:
            registered_model_id: Registered model ID.
            options: Options for listing model versions.

        Returns:
            Model versions.
        """
        async with self.get_client() as client:
            mv_list = await client.get_registered_model_versions(
                registered_model_id, **(options or ListOptions()).as_options()
            )

        if options:
            options.next_page_token = mv_list.next_page_token

        return [ModelVersion.from_basemodel(mv) for mv in mv_list.items or []]

    @overload
    async def get_model_version_by_params(self, registered_model_id: str, name: str): ...

    @overload
    async def get_model_version_by_params(self, *, external_id: str): ...

    @required_args(
        (
            "registered_model_id",
            "name",
        ),
        ("external_id",),
    )
    async def get_model_version_by_params(
        self,
        registered_model_id: str | None = None,
        name: str | None = None,
        external_id: str | None = None,
    ) -> ModelVersion | None:
        """Fetch a model version by associated parameters.

        Either fetches by using external ID or by using registered model ID and version name.

        Args:
            registered_model_id: Registered model ID.
            name: Model version.
            external_id: Model version external ID.

        Returns:
            Model version.
        """
        async with self.get_client() as client:
            try:
                mv = await client.find_model_version(
                    name=name,
                    external_id=external_id,
                    parent_resource_id=registered_model_id,
                )
            except mr_exceptions.NotFoundException:
                return None

        return ModelVersion.from_basemodel(mv)

    async def upsert_model_artifact(self, model_artifact: ModelArtifact) -> ModelArtifact:
        """Upsert a model artifact.

        Updates or creates a model artifact on the server.

        Args:
            model_artifact: Model artifact to upsert.
            model_version_id: ID of the model version this artifact will be associated to.

        Returns:
            New model artifact.
        """
        async with self.get_client() as client:
            if not model_artifact.id:
                ma = await client.create_model_artifact(model_artifact.create())
            else:
                ma = await client.update_model_artifact(model_artifact.id, model_artifact.update())
        return ModelArtifact.from_basemodel(ma)

    async def upsert_model_version_artifact(self, artifact: ArtifactT, model_version_id: str) -> ArtifactT:
        """Creates a model version artifact.

        Creates a model version artifact on the server.

        Args:
            artifact: Model version artifact to upsert.
            model_version_id: ID of the model version this artifact will be associated to.

        Returns:
            New model version artifact.
        """
        async with self.get_client() as client:
            return cast(
                ArtifactT,
                Artifact.validate_artifact(
                    await client.upsert_model_version_artifact(model_version_id, artifact.wrap())
                ),
            )

    async def get_model_artifact_by_id(self, id: str) -> ModelArtifact | None:
        """Fetch a model artifact by its ID.

        Args:
            id: Model artifact ID.

        Returns:
            Model artifact.
        """
        async with self.get_client() as client:
            try:
                ma = await client.get_model_artifact(id)
            except mr_exceptions.NotFoundException:
                return None

        return ModelArtifact.from_basemodel(ma)

    @overload
    async def get_model_artifact_by_params(
        self,
        name: str,
        model_version_id: str,
    ): ...

    @overload
    async def get_model_artifact_by_params(self, *, external_id: str): ...

    @required_args(
        (
            "name",
            "model_version_id",
        ),
        ("external_id",),
    )
    async def get_model_artifact_by_params(
        self,
        name: str | None = None,
        model_version_id: str | None = None,
        external_id: str | None = None,
    ) -> ModelArtifact | None:
        """Fetch a model artifact either by external ID or by its name and the ID of its associated model version.

        Args:
            name: Model artifact name.
            model_version_id: ID of the associated model version.
            external_id: Model artifact external ID.

        Returns:
            Model artifact.
        """
        async with self.get_client() as client:
            try:
                ma = await client.find_model_artifact(
                    name=name,
                    parent_resource_id=model_version_id,
                    external_id=external_id,
                )
            except mr_exceptions.NotFoundException:
                return None

        return ModelArtifact.from_basemodel(ma)

    async def get_model_artifacts(
        self,
        model_version_id: str | None = None,
        options: ListOptions | None = None,
    ) -> list[ModelArtifact]:
        """Fetches model artifacts.

        Args:
            model_version_id: ID of the associated model version.
            options: Options for listing model artifacts.

        Returns:
            Model artifacts.
        """
        async with self.get_client() as client:
            if model_version_id:
                art_list = await client.get_model_version_artifacts(
                    model_version_id, **(options or ListOptions()).as_options()
                )
                if options:
                    options.next_page_token = art_list.next_page_token
                models = []
                for art in art_list.items or []:
                    converted = Artifact.validate_artifact(art)
                    if isinstance(converted, ModelArtifact):
                        models.append(converted)
                return models

            ma_list = await client.get_model_artifacts(**(options or ListOptions()).as_options())
            if options:
                options.next_page_token = ma_list.next_page_token
            return [ModelArtifact.from_basemodel(ma) for ma in ma_list.items or []]

    async def get_model_version_artifacts(
        self,
        model_version_id: str,
        options: ListOptions | None = None,
    ) -> list[Artifact]:
        """Fetches model artifacts.

        Args:
            model_version_id: ID of the associated model version.
            options: Options for listing model artifacts.

        Returns:
            Model artifacts.
        """
        async with self.get_client() as client:
            art_list = await client.get_model_version_artifacts(
                model_version_id, **(options or ListOptions()).as_options()
            )
        if options:
            options.next_page_token = art_list.next_page_token
        return [Artifact.validate_artifact(art) for art in art_list.items or []]

    async def upsert_experiment(self, experiment: Experiment) -> Experiment:
        """Upsert an experiment.

        Updates or creates an experiment on the server.

        Args:
            experiment: Experiment to upsert.
        """
        async with self.get_client() as client:
            if experiment.id:
                exp = await client.update_experiment(experiment.id, experiment.update())
            elif experiment.name:
                if exp := await self.get_experiment_by_name(experiment.name):
                    exp = await client.update_experiment(exp.id, experiment.update())
                else:
                    exp = await client.create_experiment(experiment.create())
        return Experiment.from_basemodel(exp)

    async def get_experiment_by_name(self, name: str) -> Experiment | None:
        """Fetch an experiment by its name.

        Args:
            name: Experiment name.
        """
        async with self.get_client() as client:
            try:
                exp = await client.find_experiment(name=name)
            except mr_exceptions.NotFoundException:
                return None

        return Experiment.from_basemodel(exp)

    async def get_experiment_by_id(self, id: str | int) -> Experiment | None:
        """Fetch an experiment by its ID.

        Args:
            id: Experiment ID.
        """
        async with self.get_client() as client:
            try:
                exp = await client.get_experiment(str(id))
            except mr_exceptions.NotFoundException:
                return None

        return RegisteredModel.from_basemodel(exp)

    async def get_experiments(self, options: ListOptions | None = None) -> list[Experiment]:
        """Fetch experiments.

        Args:
            options: Options for listing experiments.
        """
        async with self.get_client() as client:
            exp_list = await client.get_experiments(**(options or ListOptions()).as_options())
        if options:
            options.next_page_token = exp_list.next_page_token
        return [Experiment.from_basemodel(exp) for exp in exp_list.items or []]

    async def upsert_experiment_run(self, experiment_run: ExperimentRun) -> ExperimentRun:
        """Upsert an experiment run.

        Updates or creates an experiment run on the server.

        Args:
            experiment_run: Experiment run to upsert.
        """
        async with self.get_client() as client:
            if experiment_run.id:
                exp_run = await client.create_experiment_run(experiment_run.id, experiment_run.update())
            else:
                exp_run = await client.create_experiment_run(experiment_run.create())

        return ExperimentRun.from_basemodel(exp_run)

    async def get_experiment_runs_by_experiment_id(
        self, experiment_id: str | int, options: ListOptions | None = None
    ) -> list[ExperimentRun]:
        """Fetch experiment runs by experiment ID.

        Args:
            experiment_id: Experiment ID.
            options: Options for listing experiment runs.
        """
        async with self.get_client() as client:
            try:
                exp_runs = await client.get_experiment_experiment_runs(
                    str(experiment_id), **(options or ListOptions()).as_options()
                )
            except mr_exceptions.NotFoundException:
                return []

        if options:
            options.next_page_token = exp_runs.next_page_token

        return [ExperimentRun.from_basemodel(exp_run) for exp_run in exp_runs.items or []]

    async def get_experiment_runs_by_experiment_name(
        self, experiment_name: str, options: ListOptions | None = None
    ) -> list[ExperimentRun]:
        """Fetch experiment runs by experiment name.

        Args:
            experiment_name: Experiment run to upsert.
            options: Options for listing experiment runs.

        """
        async with self.get_client() as client:
            try:
                exp = await self.get_experiment_by_name(experiment_name)
                if not exp:
                    return []
                exp_runs = await client.get_experiment_experiment_runs(
                    str(exp.id), **(options or ListOptions()).as_options()
                )
            except mr_exceptions.NotFoundException:
                return []

        if options:
            options.next_page_token = exp_runs.next_page_token

        return [ExperimentRun.from_basemodel(exp_run) for exp_run in exp_runs.items or []]

    async def get_experiment_run_by_experiment_and_run_id(
        self,
        run_id: str | int,
        experiment_name: str | None = None,
        experiment_id: str | int | None = None,
    ) -> ExperimentRun:
        """Fetch experiment run by experiment name / ID and the run ID.

        Args:
            run_id: Run ID.
            experiment_name: Experiment name.
            experiment_id: Experiment ID.
            options: Options for listing experiment runs.

        Returns:
            Experiment run.
        """
        async with self.get_client() as client:
            try:
                if experiment_name:
                    exp = await self.get_experiment_by_name(experiment_name)
                elif experiment_id:
                    exp = await self.get_experiment_by_id(str(experiment_id))
                else:
                    msg = "Either experiment_name or experiment_id must be provided"
                    raise ValueError(msg)
                if not exp:
                    return None

                exp_run = await client.get_experiment_run(str(run_id))
            except mr_exceptions.NotFoundException:
                return None

        return ExperimentRun.from_basemodel(exp_run)

    async def get_experiment_run_by_experiment_and_run_name(
        self,
        run_name: str,
        experiment_name: str | None = None,
        experiment_id: str | int | None = None,
        options: ListOptions | None = None,
    ) -> ExperimentRun:
        """Fetch experiment runs by experiment name / ID and the run ID.

        Args:
            run_name: Run name.
            experiment_name: Experiment name.
            experiment_id: Experiment ID.
            options: Options for listing experiment runs.

        Returns:
            Experiment run.
        """
        async with self.get_client() as client:
            exp = None
            try:
                if experiment_name:
                    exp = await self.get_experiment_by_name(experiment_name)
                elif experiment_id:
                    exp = await self.get_experiment_by_id(str(experiment_id))

                if not exp:
                    return None

                exp_run = await client.get_experiment_run(exp.id)
            except mr_exceptions.NotFoundException:
                return None

        return ExperimentRun.from_basemodel(exp_run)

    async def get_experiment_run_by_id(self, id: str) -> ExperimentRun:
        """Fetch an experiment run by its ID.

        Args:
            id: Experiment run ID.
        """
        async with self.get_client() as client:
            try:
                exp_run = await client.get_experiment_run(id)
            except mr_exceptions.NotFoundException:
                return None

        return ExperimentRun.from_basemodel(exp_run)

    async def upsert_experiment_run_artifact(
        self, experiment_run_id: str, artifact: ExperimentRunArtifact
    ) -> ExperimentRunArtifact:
        """Upsert an experiment run artifact (parameter, metric, or dataset).

        Updates or creates an experiment run on the server.

        Args:
            experiment_run_id: Experiment run ID.
            artifact: Artifact to upsert.
        """
        async with self.get_client() as client:
            return Artifact.validate_artifact(
                await client.upsert_experiment_run_artifact(
                    experimentrun_id=experiment_run_id, artifact=artifact.wrap()
                )
            )

    @overload
    async def get_artifacts_by_experiment_run_params(self, run_id: str | int, options: ListOptions | None = None): ...

    @overload
    async def get_artifacts_by_experiment_run_params(
        self,
        run_name: str,
        experiment_name: str | None = None,
        options: ListOptions | None = None,
    ): ...

    @overload
    async def get_artifacts_by_experiment_run_params(
        self,
        run_name: str,
        experiment_id: str | int | None = None,
        options: ListOptions | None = None,
    ): ...

    @required_args(("run_id",), ("run_name", "experiment_name"), ("run_name", "experiment_id"))
    async def get_artifacts_by_experiment_run_params(
        self,
        run_id: str | int | None = None,
        run_name: str | None = None,
        experiment_name: str | None = None,
        experiment_id: str | int | None = None,
        *,
        options: ListOptions | None = None,
    ) -> ExperimentRunArtifact:
        """Fetch a log by experiment run ID and name.

        Args:
            run_id: Experiment run ID.
            run_name: Experiment run name.
            experiment_name: Experiment name.
            experiment_id: Experiment ID.

        Keyword Args:
            options: Options for listing experiment run artifacts.
        """
        async with self.get_client() as client:
            try:
                if not run_id and run_name:
                    if experiment_name:
                        exp_runs = await self.get_experiment_runs_by_experiment_name(
                            experiment_name=experiment_name,
                            options=ListOptions(limit=100),
                        )
                    elif experiment_id:
                        exp_runs = await self.get_experiment_runs_by_experiment_id(
                            experiment_id=experiment_id,
                            options=ListOptions(limit=100),
                        )
                    else:
                        msg = "Either experiment_name or experiment_id must be provided"
                        raise ValueError(msg)

                    run = next((r for r in exp_runs if r.name == run_name), None)
                    if not run:
                        print(
                            f"Could not find run {run_name} "
                            f"in experiment {experiment_name} within the first 100 runs. "
                            "Please narrow your search by run id."
                        )
                        return []
                    run_id = run.id

                logs = await client.get_experiment_run_artifacts(
                    str(run_id), **(options or ListOptions()).as_options()
                )
            except mr_exceptions.NotFoundException:
                return []

        if options:
            options.next_page_token = logs.next_page_token
        return [Artifact.validate_artifact(log) for log in logs.items or []]

    async def get_artifacts(
        self,
        filter_query: str | None = None,
        artifact_type: str | None = None,
        options: ListOptions | None = None,
    ) -> list[Artifact]:
        """Get artifacts with filtering capabilities.

        Args:
            filter_query: A SQL-like query string to filter the list of entities.
            artifact_type: Specifies the artifact type for listing artifacts.
            options: Options for listing artifacts.

        Returns:
            List of artifacts matching the filter criteria.
        """
        # Create options dict with additional filters
        options_dict = (options or ListOptions()).as_options()
        if filter_query is not None:
            options_dict["filter_query"] = filter_query
        if artifact_type is not None:
            options_dict["artifact_type"] = artifact_type

        async with self.get_client() as client:
            try:
                response = await client.get_artifacts(**options_dict)
            except mr_exceptions.NotFoundException:
                return []

        if options:
            options.next_page_token = response.next_page_token
        return [Artifact.validate_artifact(artifact) for artifact in response.items or []]
