from __future__ import annotations

import json
import threading

import pytest

from model_registry import ModelRegistry, utils
from model_registry.types import Artifact, ListOptions, Pager


def get_artifacts(
    model_registry: ModelRegistry,
    filter_query: str | None = None,
    artifact_type: str | None = None,
) -> Pager[Artifact]:
    def artifacts_list(options: ListOptions) -> list[Artifact]:
        return model_registry.async_runner(
            model_registry._api.get_artifacts(
                filter_query=filter_query,
                artifact_type=artifact_type,
                options=options,
            )
        )

    return Pager[Artifact](artifacts_list)


@pytest.fixture
def schema_json():
    schema: dict[str, dict] = {"epochs": {}}
    return json.dumps(schema)


@pytest.mark.e2e
def test_start_experiment_run(client: ModelRegistry, schema_json: str):
    with client.start_experiment_run(experiment_name="Experiment_Test") as run:
        run.log_param("input1", 5.75)
        run.log_metric(
            key="rval",
            value=10,
            step=4,
            timestamp="0",  # type: ignore[arg-type]
            description="This is a test metric",
        )
        run.log_dataset(
            name="dataset_1",
            source_type="local",
            uri="s3://datasets/test",
            schema=schema_json,  # type: ignore[arg-type]
            profile="random_profile",  # type: ignore[arg-type]
        )

    assert len(run.get_logs()) == 3
    param = run.get_log("params", "input1")
    metric = run.get_log("metrics", "rval")
    dataset = run.get_log("datasets", "dataset_1")
    assert param
    assert metric
    assert dataset

    assert param.value == 5.75  # type: ignore[union-attr]
    assert metric.value == 10  # type: ignore[union-attr]
    assert metric.step == 4  # type: ignore[union-attr]
    assert metric.timestamp == "0"  # type: ignore[union-attr]
    assert metric.description == "This is a test metric"
    assert metric.name == "rval"


@pytest.mark.skip(reason="Skipping test_start_experiment_run_with_advanced_scenarios")
@pytest.mark.e2e
def test_start_experiment_run_with_advanced_scenarios(
    client: ModelRegistry, get_temp_dir_with_models, patch_s3_env, schema_json: str
):
    with client.start_experiment_run(experiment_name="Experiment_Test") as run:
        run.log_param("input1", 5.75)
        run.log_param("input1", 500)
        for i in range(10):
            run.log_metric(f"metric_{i}", value=i * 1000, step=i, timestamp="0")  # type: ignore[arg-type]

    assert len(run.get_logs()) == 11
    assert run.get_log("params", "input1").value == 500  # type: ignore[union-attr]

    with client.start_experiment_run(
        experiment_name="Experiment_Test_URI_Provided"
    ) as run:
        run.log_dataset(
            name="dataset_1",
            source_type="s3",
            uri="s3://datasets/test",
            schema=schema_json,  # type: ignore[arg-type]
            profile="random_profile",  # type: ignore[arg-type]
        )
    assert run.get_log("datasets", "dataset_1").uri == "s3://datasets/test"  # type: ignore[union-attr]

    # Test actual
    model_dir, _ = get_temp_dir_with_models
    bucket, s3_endpoint, access_key_id, secret_access_key, region = patch_s3_env
    with client.start_experiment_run(experiment_name="Experiment_Test_3") as run:
        run.log_dataset(
            name="dataset_1",
            source_type="local",
            schema=schema_json,  # type: ignore[arg-type]
            profile="random_profile",  # type: ignore[arg-type]
            file_path=model_dir,
            s3_auth=utils.S3Params(
                endpoint_url=s3_endpoint,
                bucket_name=bucket,
                s3_prefix="datasets",
                access_key_id=access_key_id,
                secret_access_key=secret_access_key,
                region=region,
            ),
        )
    assert run.get_log("datasets", "dataset_1").uri.startswith("s3://")  # type: ignore[union-attr]


@pytest.mark.e2e
def test_experiments(client: ModelRegistry):
    with client.start_experiment_run(experiment_name="Experiment_Test_2"):
        pass
    found_exp = False
    for experiment in client.get_experiments():
        if experiment.name == "Experiment_Test_2":
            found_exp = True
    assert found_exp


@pytest.mark.e2e
def test_get_experiment_runs(client: ModelRegistry):
    """This tests:
    - get_experiment_runs(name)
    - get_experiment_runs(id)
    """
    with client.start_experiment_run(experiment_name="Experiment_Test_3") as run:
        run.log_param("input1", 5.75)
    runs_by_name = client.get_experiment_runs(experiment_name="Experiment_Test_3")
    runs_by_id = client.get_experiment_runs(experiment_id=run.info.experiment_id)

    assert runs_by_name.next_item().id == runs_by_id.next_item().id  # type: ignore[attr-defined]
    runs_by_name.restart()
    runs_by_id.restart()

    found_exp_run_by_id = False
    found_exp_run_by_name = False

    for r in runs_by_name:
        if r.name == run.info.name:
            found_exp_run_by_name = True

    for r in runs_by_id:
        if r.id == run.info.id:
            found_exp_run_by_id = True

    assert found_exp_run_by_id
    assert found_exp_run_by_name


@pytest.mark.e2e
def test_get_experiment_run_with_artifact_types(
    client: ModelRegistry, schema_json: str
):
    with client.start_experiment_run(experiment_name="Experiment_Test_4") as run:
        run.log_dataset(
            name="dataset_1",
            source_type="local",
            uri="s3://datasets/test",
            schema=schema_json,  # type: ignore[arg-type]
            profile="random_profile",  # type: ignore[arg-type]
            description="This is a test dataset",
        )
        run.log_metric(
            key="metric_1",
            value=10,
            step=4,
            timestamp="0",  # type: ignore[arg-type]
            description="This is a test metric",
        )
        run.log_param(
            key="param_1",
            value=10,
            description="This is a test param",
        )

    dataset_log = client.get_experiment_run_logs(
        run_id=run.info.id,
    )
    assert dataset_log.next_item().name.endswith("1")  # type: ignore[attr-defined]
    assert dataset_log.next_item()
    assert dataset_log.next_item()
    try:
        # fail if we get a 4th item
        dataset_log.next_item()  # type: ignore[unused-coroutine]
        pytest.fail("Expected StopIteration")
    except StopIteration:
        assert True


@pytest.mark.e2e
def test_start_experiment_run_nested(client: ModelRegistry):
    """Test that starting a nested experiment run without nested=True raises an error."""
    with client.start_experiment_run(experiment_name="Experiment_Test_6") as run:
        run.log_metric(
            key="rval",
            value=10,
            step=4,
            description="This is a test metric",
        )
        # This should fail because nested=True is not specified
        with pytest.raises(ValueError, match="Experiment run is already active"):  # noqa: SIM117
            with client.start_experiment_run():
                pass

        with client.start_experiment_run(nested=True) as run2:
            run2.log_metric(
                key="nested_metric",
                value=20,
                step=1,
                description="This is a nested test metric",
            )
    # Assert logs are correct
    for artifact in client.get_experiment_run_logs(run_id=run.info.id):
        assert artifact.value == 10  # type: ignore[union-attr]
        assert "nested" not in artifact.description  # type: ignore[operator]

    for artifact in client.get_experiment_run_logs(run_id=run2.info.id):
        assert artifact.value == 20  # type: ignore[union-attr]
        assert "nested" in artifact.description  # type: ignore[operator]

    exp_run = client.get_experiment_run(run_id=run.info.id)
    assert exp_run.custom_properties is None
    assert exp_run.experiment_id == run.info.experiment_id

    exp_run = client.get_experiment_run(run_id=run2.info.id)
    assert "kubeflow.parent_run_id" in exp_run.custom_properties  # type: ignore[operator]
    assert exp_run.custom_properties["kubeflow.parent_run_id"] == run.info.id  # type: ignore[index]
    assert exp_run.experiment_id == run.info.experiment_id
    assert exp_run.name == run2.info.name
    assert exp_run.id == run2.info.id


@pytest.mark.e2e
async def test_start_experiment_run_thread_safety(client: ModelRegistry):
    gen_name = utils.generate_name("Experiment_Test")
    exp = client.create_experiment(gen_name)
    assert exp.id

    def run_experiment(client: ModelRegistry):
        with client.start_experiment_run(experiment_name=gen_name) as run:
            run.log_metric(
                key="rval",
                value=10,
                step=4,
                description="This is a test metric",
            )

    threads = []
    for _i in range(5):
        thread = threading.Thread(target=run_experiment, args=(client,))
        threads.append(thread)
        thread.start()

    for thread in threads:
        thread.join()

    # Check that all runs are created
    exp_run = client.get_experiment_runs(experiment_name=gen_name)
    ctr = 0
    for run in exp_run:
        ctr += 1
        assert run.custom_properties is None
    assert ctr == 5


@pytest.mark.e2e
async def test_start_experiment_run_nested_thread_safety(client: ModelRegistry):
    gen_name = utils.generate_name("Experiment_Test")
    exp = client.create_experiment(gen_name)
    assert exp.id

    def run_experiment(client: ModelRegistry):
        with client.start_experiment_run(experiment_name=gen_name) as run:
            run.log_metric(
                key="rval",
                value=10,
                step=4,
                description="This is a test metric",
            )
            with client.start_experiment_run(nested=True) as run2:
                run2.log_metric(
                    key="rval",
                    value=20,
                    step=55,
                    description="This is a nested run test metric",
                )

    threads = []
    for _i in range(5):
        thread = threading.Thread(target=run_experiment, args=(client,))
        threads.append(thread)
        thread.start()

    for thread in threads:
        thread.join()

    # Check that all runs are created
    exp_run = client.get_experiment_runs(experiment_name=gen_name)
    ctr = 0
    for r in exp_run:
        ctr += 1
        if r.custom_properties:
            if parent_id := r.custom_properties.get("kubeflow.parent_run_id"):
                assert parent_id
        else:
            assert r.custom_properties is None

    assert ctr == 10


@pytest.mark.e2e
def test_bulk_metrics_and_params_retrieval(client: ModelRegistry):  # noqa: C901
    """Test bulk retrieval of metrics and parameters for multiple experiment runs.

    This test:
    1. Creates multiple experiments with multiple runs each
    2. Logs multiple metrics and parameters for each run
    3. Selects a subset of runs from different experiments
    4. Tests bulk retrieval of all metrics and parameters for the selected runs
    """
    experiments = []
    for i in range(4):
        exp_name = f"BulkTest_Experiment_{i}"
        exp = client.create_experiment(exp_name)
        experiments.append(exp)

    all_runs = []
    for exp_idx, exp in enumerate(experiments):
        exp_runs = []
        for run_idx in range(4):
            with client.start_experiment_run(experiment_name=exp.name) as run:
                run.log_param("learning_rate", 0.001 * (run_idx + 1))
                run.log_param("batch_size", 32 * (run_idx + 1))
                run.log_param("epochs", 10 + run_idx)
                run.log_param("optimizer", f"adam_{run_idx}")
                run.log_param("model_type", f"resnet_{exp_idx}")

                for metric_idx in range(3):
                    run.log_metric(
                        key=f"accuracy_epoch_{metric_idx}",
                        value=0.8 + (run_idx * 0.05) + (metric_idx * 0.01),
                        step=metric_idx,
                        description=f"Accuracy at epoch {metric_idx}",
                    )
                    run.log_metric(
                        key=f"loss_epoch_{metric_idx}",
                        value=1.0 - (run_idx * 0.1) - (metric_idx * 0.05),
                        step=metric_idx,
                        description=f"Loss at epoch {metric_idx}",
                    )
                run.log_metric("final_accuracy", 0.95 + (run_idx * 0.01))
                run.log_metric("training_time", 120 + (run_idx * 10))

                exp_runs.append(run.info)
        all_runs.append(exp_runs)

    selected_runs = [
        all_runs[0][0],
        all_runs[0][2],
        all_runs[1][1],
        all_runs[1][2],
        all_runs[1][3],
        all_runs[3][1],
    ]

    def sql_quote(s):
        escaped = s.replace("'", "''")
        return f"'{escaped}'"

    run_ids_list = ",".join([sql_quote(run.id) for run in selected_runs])
    run_ids_filter = f"experimentRunId IN ({run_ids_list})"

    metrics = list(get_artifacts(client, artifact_type="metric", filter_query=run_ids_filter))
    params = list(get_artifacts(client, artifact_type="parameter", filter_query=run_ids_filter))

    expected_metrics_per_run = 8
    expected_params_per_run = 5

    assert len(metrics) == len(selected_runs) * expected_metrics_per_run
    assert len(params) == len(selected_runs) * expected_params_per_run

    expected_metric_names = {
        "accuracy_epoch_0",
        "accuracy_epoch_1",
        "accuracy_epoch_2",
        "loss_epoch_0",
        "loss_epoch_1",
        "loss_epoch_2",
        "final_accuracy",
        "training_time",
    }
    expected_param_names = {"learning_rate", "batch_size", "epochs", "optimizer", "model_type"}

    expected_metric_attributes = {(r.experiment_id, r.id, n) for r in selected_runs for n in expected_metric_names}
    metric_attributes = {(m.experiment_id, m.experiment_run_id, m.name) for m in metrics}
    assert metric_attributes == expected_metric_attributes

    expected_param_attributes = {(r.experiment_id, r.id, n) for r in selected_runs for n in expected_param_names}
    param_attributes = {(m.experiment_id, m.experiment_run_id, m.name) for m in params}

    assert param_attributes == expected_param_attributes
    # Verify values are reasonable
    for metric in metrics:
        assert isinstance(metric.value, (int, float))  # type: ignore[attr-defined]
        assert metric.value > 0  # type: ignore[attr-defined]

    for param in params:
        assert param.value is not None  # type: ignore[attr-defined]

    # Test filtering by specific metric names using SQL
    accuracy_metrics = list(get_artifacts(client, artifact_type="metric", filter_query=f'{run_ids_filter} AND name LIKE "%accuracy%"'))

    # Should have 4 accuracy metrics per run (3 epoch + 1 final)
    expected_accuracy_metrics = len(selected_runs) * 4
    assert len(accuracy_metrics) == expected_accuracy_metrics

    # Test filtering by specific parameter names using SQL
    learning_rate_params  = list(get_artifacts(client, artifact_type="parameter", filter_query=f'{run_ids_filter} AND name = "learning_rate"'))

    # Should have 1 learning_rate parameter per run
    assert len(learning_rate_params) == len(selected_runs)
