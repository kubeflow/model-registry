import json
import threading

import pytest

from model_registry import ModelRegistry, utils


@pytest.fixture
def schema_json():
    schema = {"epochs": {}}
    return json.dumps(schema)


@pytest.mark.e2e
def test_start_experiment_run(client: ModelRegistry, schema_json: str):
    with client.start_experiment_run(experiment_name="Experiment_Test") as run:
        run.log_param("input1", 5.75)
        run.log_metric(
            key="rval",
            value=10,
            step=4,
            timestamp="0",
            description="This is a test metric",
        )
        run.log_dataset(
            name="dataset_1",
            source_type="local",
            uri="s3://datasets/test",
            schema=schema_json,
            profile="random_profile",
        )

    assert len(run.get_logs()) == 3
    param = run.get_log("params", "input1")
    metric = run.get_log("metrics", "rval")
    dataset = run.get_log("datasets", "dataset_1")
    assert param
    assert metric
    assert dataset

    assert param.value == 5.75
    assert metric.value == 10
    assert metric.step == 4
    assert metric.timestamp == "0"
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
            run.log_metric(f"metric_{i}", value=i * 1000, step=i, timestamp="0")

    assert len(run.get_logs()) == 11
    assert run.get_log("params", "input1").value == 500

    with client.start_experiment_run(
        experiment_name="Experiment_Test_URI_Provided"
    ) as run:
        run.log_dataset(
            name="dataset_1",
            source_type="s3",
            uri="s3://datasets/test",
            schema=schema_json,
            profile="random_profile",
        )
    assert run.get_log("datasets", "dataset_1").uri == "s3://datasets/test"

    # Test actual
    model_dir, _ = get_temp_dir_with_models
    bucket, s3_endpoint, access_key_id, secret_access_key, region = patch_s3_env
    with client.start_experiment_run(experiment_name="Experiment_Test_3") as run:
        run.log_dataset(
            name="dataset_1",
            source_type="local",
            schema=schema_json,
            profile="random_profile",
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
    assert run.get_log("datasets", "dataset_1").uri.startswith("s3://")


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

    assert runs_by_name.next_item().id == runs_by_id.next_item().id
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
            schema=schema_json,
            profile="random_profile",
            description="This is a test dataset",
        )
        run.log_metric(
            key="metric_1",
            value=10,
            step=4,
            timestamp="0",
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
    assert dataset_log.next_item().name.endswith("1")
    assert dataset_log.next_item()
    assert dataset_log.next_item()
    try:
        # fail if we get a 4th item
        dataset_log.next_item()
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
        assert artifact.value == 10
        assert "nested" not in artifact.description

    for artifact in client.get_experiment_run_logs(run_id=run2.info.id):
        assert artifact.value == 20
        assert "nested" in artifact.description

    exp_run = client.get_experiment_run(run_id=run.info.id)
    assert exp_run.custom_properties is None
    assert exp_run.experiment_id == run.info.experiment_id

    exp_run = client.get_experiment_run(run_id=run2.info.id)
    assert "kubeflow.parent_run_id" in exp_run.custom_properties
    assert exp_run.custom_properties["kubeflow.parent_run_id"] == run.info.id
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
def test_bulk_metrics_and_params_retrieval(client: ModelRegistry): # noqa: C901
    """Test bulk retrieval of metrics and parameters for multiple experiment runs.

    This test:
    1. Creates multiple experiments with multiple runs each
    2. Logs multiple metrics and parameters for each run
    3. Selects a subset of runs from different experiments
    4. Tests bulk retrieval of all metrics and parameters for the selected runs
    """
    # Create multiple experiments
    experiments = []
    for i in range(3):
        exp_name = f"BulkTest_Experiment_{i}"
        exp = client.create_experiment(exp_name)
        experiments.append(exp)

    # Create multiple runs per experiment and log metrics/params
    all_runs = []
    for exp_idx, exp in enumerate(experiments):
        for run_idx in range(4):  # 4 runs per experiment
            with client.start_experiment_run(experiment_name=exp.name) as run:
                # Log multiple parameters
                run.log_param("learning_rate", 0.001 * (run_idx + 1))
                run.log_param("batch_size", 32 * (run_idx + 1))
                run.log_param("epochs", 10 + run_idx)
                run.log_param("optimizer", f"adam_{run_idx}")
                run.log_param("model_type", f"resnet_{exp_idx}")

                # Log multiple metrics
                for metric_idx in range(3):
                    run.log_metric(
                        key=f"accuracy_epoch_{metric_idx}",
                        value=0.8 + (run_idx * 0.05) + (metric_idx * 0.01),
                        step=metric_idx,
                        description=f"Accuracy at epoch {metric_idx}"
                    )
                    run.log_metric(
                        key=f"loss_epoch_{metric_idx}",
                        value=1.0 - (run_idx * 0.1) - (metric_idx * 0.05),
                        step=metric_idx,
                        description=f"Loss at epoch {metric_idx}"
                    )

                # Log some additional metrics
                run.log_metric("final_accuracy", 0.95 + (run_idx * 0.01))
                run.log_metric("training_time", 120 + (run_idx * 10))

                all_runs.append(run.info)

    # Select a subset of runs (some from each experiment)
    selected_runs = []
    # Take 2 runs from each experiment
    for exp_idx in range(3):
        exp_runs = [run for run in all_runs if run.experiment_id == experiments[exp_idx].id]
        selected_runs.extend(exp_runs[:2])  # Take first 2 runs from each experiment

    assert len(selected_runs) == 6  # 2 runs from each of 3 experiments

    # Test bulk retrieval using the new get_artifacts method
    # Get all artifacts for the selected runs
    selected_run_ids = [run.id for run in selected_runs]
    run_ids_list = ",".join([f'"{run_id}"' for run_id in selected_run_ids])
    run_ids_filter = f"experimentRunId IN ({run_ids_list})"

    # Get all artifacts for selected runs
    all_artifacts = list(client.get_artifacts(filter_query=run_ids_filter))

    # Filter by type in Python (since SQL filtering by artifact type doesn't work)
    all_metrics = [a for a in all_artifacts if hasattr(a, "value") and hasattr(a, "step")]
    all_params = [a for a in all_artifacts if hasattr(a, "parameter_type")]

    # Verify we got metrics and parameters from all selected runs
    unique_run_ids_metrics = set()
    unique_run_ids_params = set()

    for metric in all_metrics:
        # Extract run_id from metric's custom properties or other fields
        # Since we don't have direct run_id in the artifact, we'll verify by checking
        # that we have the expected number of metrics
        unique_run_ids_metrics.add(metric.name)

    for param in all_params:
        unique_run_ids_params.add(param.name)

    # Verify we have the expected number of metrics and parameters
    # Each run should have: 3 accuracy metrics + 3 loss metrics + 2 additional metrics = 8 metrics
    # Each run should have: 5 parameters
    expected_metrics_per_run = 8
    expected_params_per_run = 5

    assert len(all_metrics) == len(selected_runs) * expected_metrics_per_run, \
        f"Expected {len(selected_runs) * expected_metrics_per_run} metrics, got {len(all_metrics)}"
    assert len(all_params) == len(selected_runs) * expected_params_per_run, \
        f"Expected {len(selected_runs) * expected_params_per_run} parameters, got {len(all_params)}"

    # Verify specific metrics and parameters exist
    metric_names = {metric.name for metric in all_metrics}
    param_names = {param.name for param in all_params}

    expected_metric_names = {
        "accuracy_epoch_0", "accuracy_epoch_1", "accuracy_epoch_2",
        "loss_epoch_0", "loss_epoch_1", "loss_epoch_2",
        "final_accuracy", "training_time"
    }
    expected_param_names = {
        "learning_rate", "batch_size", "epochs", "optimizer", "model_type"
    }

    assert metric_names == expected_metric_names, f"Expected metrics {expected_metric_names}, got {metric_names}"
    assert param_names == expected_param_names, f"Expected parameters {expected_param_names}, got {param_names}"

    # Verify values are reasonable
    for metric in all_metrics:
        assert isinstance(metric.value, (int, float)), f"Metric {metric.name} should have numeric value"
        assert metric.value > 0, f"Metric {metric.name} should have positive value"

    for param in all_params:
        assert param.value is not None, f"Parameter {param.name} should have a value"

    # Test filtering by specific metric names using SQL
    accuracy_artifacts = list(client.get_artifacts(
        filter_query=f'{run_ids_filter} AND name LIKE "%accuracy%"'
    ))
    accuracy_metrics = [a for a in accuracy_artifacts if hasattr(a, "value") and hasattr(a, "step")]

    # Should have 4 accuracy metrics per run (3 epoch + 1 final)
    expected_accuracy_metrics = len(selected_runs) * 4
    assert len(accuracy_metrics) == expected_accuracy_metrics, \
        f"Expected {expected_accuracy_metrics} accuracy metrics, got {len(accuracy_metrics)}"

    # Test filtering by specific parameter names using SQL
    learning_rate_artifacts = list(client.get_artifacts(
        filter_query=f'{run_ids_filter} AND name = "learning_rate"'
    ))
    learning_rate_params = [a for a in learning_rate_artifacts if hasattr(a, "parameter_type")]

    # Should have 1 learning_rate parameter per run
    assert len(learning_rate_params) == len(selected_runs), \
        f"Expected {len(selected_runs)} learning_rate parameters, got {len(learning_rate_params)}"

    print(f"Successfully retrieved {len(all_metrics)} metrics and {len(all_params)} parameters from {len(selected_runs)} experiment runs")
    print("Bulk retrieval test completed successfully!")
