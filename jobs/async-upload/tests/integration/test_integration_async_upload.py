"""Integration tests for async-upload job functionality."""

import os
import time
import uuid
from pathlib import Path

import pytest
import requests
import yaml

# Import model registry client for API operations
try:
    from model_registry import ModelRegistry
    from model_registry.types import ArtifactState
except ImportError:
    pytest.skip("model-registry client not available", allow_module_level=True)

HTTP_SOURCE = "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/classification/mnist/model/mnist-8.onnx"


@pytest.fixture(scope="session")
def k8s_client():
    """Load Kubernetes config and return client."""
    kubernetes = pytest.importorskip("kubernetes")
    try:
        kubernetes.config.load_kube_config()
    except Exception:
        # If kubeconfig loading fails, try in-cluster config
        kubernetes.config.load_incluster_config()
    return kubernetes.client.ApiClient()


@pytest.fixture(scope="session")
def k8s_batch_client():
    """Return Kubernetes Batch API client."""
    kubernetes = pytest.importorskip("kubernetes")
    try:
        kubernetes.config.load_kube_config()
    except Exception:
        kubernetes.config.load_incluster_config()
    return kubernetes.client.BatchV1Api()


@pytest.fixture(scope="session")
def k8s_core_client():
    """Return Kubernetes Core API client."""
    kubernetes = pytest.importorskip("kubernetes")
    try:
        kubernetes.config.load_kube_config()
    except Exception:
        kubernetes.config.load_incluster_config()
    return kubernetes.client.CoreV1Api()


@pytest.fixture(scope="session")
def model_registry_client():
    """Create model registry client for integration tests."""
    mr_host_url = os.environ.get("MR_HOST_URL", "http://localhost:8080")
    # Parse URL to extract host and port
    from urllib.parse import urlparse

    parsed = urlparse(mr_host_url)
    host = f"{parsed.scheme}://{parsed.hostname}"
    port = parsed.port or (443 if parsed.scheme == "https" else 8080)

    return ModelRegistry(host, port, author="integration-test", is_secure=False)


def apply_job_with_strategic_merge(
    rm_id: str,
    mv_id: str,
    ma_id: str,
    container_image_uri: str,
    k8s_client,
    env,
    tmp_path,
) -> str:
    """Apply job using Kustomize strategic merge patches."""
    import subprocess
    import time

    patch_env = {
        "MODEL_SYNC_MODEL_ID": rm_id,
        "MODEL_SYNC_MODEL_VERSION_ID": mv_id,
        "MODEL_SYNC_MODEL_ARTIFACT_ID": ma_id,
        **env,
    }
    patch_env_list = [{"name": name, "value": value} for name, value in patch_env.items()]

    # Strategic merge patch template - only patch the image and env vars, keep original job name
    patch_obj = {
        "apiVersion": "batch/v1",
        "kind": "Job",
        "metadata": {"name": "my-async-upload-job"},
        "spec": {
            "template": {
                "spec": {
                    "containers": [
                        {
                            "name": "async-upload",
                            "image": container_image_uri,
                            "env": patch_env_list,
                        }
                    ]
                }
            }
        },
    }

    # patch_env = patch_obj['spec']['template']['spec']['containers'][0]['env']

    patch_template = yaml.safe_dump(patch_obj, sort_keys=False)

    # Get the path to the sample job file
    base_job_path = Path(__file__).parent.parent.parent / "samples" / "sample_job_s3_to_oci.yaml"

    # Kustomization template using relative path and modern patches syntax
    kustomization_obj = {
        "apiVersion": "kustomize.config.k8s.io/v1beta1",
        "kind": "Kustomization",
        "resources": ["sample_job_s3_to_oci.yaml"],
        "patches": [{"path": "patch.yaml", "target": {"kind": "Job", "name": "my-async-upload-job"}}],
    }

    kustomization_template = yaml.safe_dump(kustomization_obj)

    manifest_dir = tmp_path / "templates"
    manifest_dir.mkdir()

    # Copy the base job file into temp directory
    import shutil

    base_job_copy = manifest_dir / "sample_job_s3_to_oci.yaml"
    shutil.copy2(base_job_path, base_job_copy)

    # Write the patch file
    patch_file = manifest_dir / "patch.yaml"
    with open(patch_file, "w") as f:
        f.write(patch_template)

    # Write the kustomization file
    kustomize_file = manifest_dir / "kustomization.yaml"
    with open(kustomize_file, "w") as f:
        f.write(kustomization_template)

    # Delete existing job if it exists (Jobs are immutable)
    delete_result = subprocess.run(
        [
            "kubectl",
            "delete",
            "job",
            "my-async-upload-job",
            "-n",
            "default",
            "--ignore-not-found=true",
        ],
        capture_output=True,
        text=True,
        check=False,
    )
    if delete_result.returncode == 0:
        print(f"Deleted existing job: {delete_result.stdout.strip()}")
        # Wait a moment for deletion to complete
        time.sleep(3)

    # Apply resources using kubectl apply -k
    result = subprocess.run(
        ["kubectl", "apply", "-k", "."],
        capture_output=True,
        text=True,
        cwd=manifest_dir,
        check=False,
    )

    if result.returncode != 0:
        raise Exception(f"kubectl apply failed: {result.stderr}")

    # Describe job
    print("Applied Job:")
    result = subprocess.run(
        ["kubectl", "describe", "jobs/my-async-upload-job"],
        capture_output=True,
        text=True,
        cwd=manifest_dir,
        check=False,
    )
    print(result.stdout)

    # Return the original job name since we're not changing it
    return "my-async-upload-job"


def wait_for_job_completion(job_name: str, namespace: str, batch_client, timeout_seconds: int = 60) -> bool:
    """Wait for job completion and return success status."""
    start_time = time.time()

    while time.time() - start_time < timeout_seconds:
        try:
            job = batch_client.read_namespaced_job(name=job_name, namespace=namespace)

            # Check if job is complete
            if job.status.conditions:
                for condition in job.status.conditions:
                    print(condition)
                    if condition.type == "Complete" and condition.status == "True":
                        return True
                    elif condition.type == "Failed" and condition.status == "True":
                        return False

            time.sleep(10)  # Wait 10 seconds before checking again

        except Exception as e:
            if hasattr(e, "status") and getattr(e, "status", None) == 404:
                # Job doesn't exist yet
                time.sleep(5)
                continue
            else:
                raise

    return False


def upload_to_minio(file_path: str, bucket: str, key: str) -> None:
    """Upload file to MinIO using boto3."""
    import boto3
    from botocore.exceptions import ClientError

    # MinIO credentials (hardcoded for this test)
    access_key = "minioadmin"
    secret_key = "minioadmin"

    # Create S3 client configured for MinIO
    s3_client = boto3.client(
        "s3",
        endpoint_url="http://localhost:9000",
        aws_access_key_id=access_key,
        aws_secret_access_key=secret_key,
        region_name="us-east-1",  # MinIO doesn't care about region but boto3 requires it
    )

    try:
        # Upload the file
        s3_client.upload_file(file_path, bucket, key)
    except ClientError as e:
        raise Exception(f"Failed to upload to MinIO: {e}")


def _setup_s3(tmp_path):
    model_dirpath = tmp_path / "model"
    model_filepath = model_dirpath / "mnist-8.onnx"

    # Download the model
    response = requests.get(HTTP_SOURCE)
    response.raise_for_status()

    model_dirpath.mkdir()
    with open(model_filepath, "wb") as f:
        f.write(response.content)

    print("Uploading to MinIO...")
    bucket = "default"
    key = "my-model/mnist-8.onnx"
    upload_to_minio(str(model_filepath), bucket, key)


@pytest.mark.parametrize(
    "setup,env",
    [
        (
            None,
            {
                "MODEL_SYNC_SOURCE_TYPE": "uri",
                "MODEL_SYNC_SOURCE_URI": HTTP_SOURCE,
            },
        ),
        (
            _setup_s3,
            {
                "MODEL_SYNC_SOURCE_TYPE": "s3",
                "MODEL_SYNC_SOURCE_AWS_KEY": "my-model",
                "MODEL_SYNC_SOURCE_S3_CREDENTIALS_PATH": "/opt/creds/source",
            },
        ),
    ],
)
@pytest.mark.integration
def test_async_upload_integration(
    setup,
    env,
    tmp_path,
    k8s_client,
    k8s_batch_client,
    k8s_core_client,
    model_registry_client: ModelRegistry,
):
    """Test the complete async-upload job integration.

    This test:
    1. Creates a RegisteredModel, ModelVersion, and placeholder ModelArtifact
    2. Downloads and uploads an ONNX model to MinIO
    3. Creates and applies a Kubernetes job using kustomize
    4. Waits for job completion
    5. Validates the final artifact state
    """
    env = env.copy()

    if "MODEL_SYNC_DESTINATION_OCI_BASE_IMAGE" in os.environ:
        env["MODEL_SYNC_DESTINATION_OCI_BASE_IMAGE"] = os.environ["MODEL_SYNC_DESTINATION_OCI_BASE_IMAGE"]

    # TODO use yield fixtures for setup/teardown
    # https://docs.pytest.org/en/6.2.x/fixture.html#yield-fixtures-recommended
    if setup is not None:
        setup(tmp_path)

    # Configuration
    container_image_uri = os.environ.get(
        "CONTAINER_IMAGE_URI", "ghcr.io/kubeflow/model-registry/job/async-upload:latest"
    )
    job_name = f"test-async-upload-job-{uuid.uuid4().hex[:8]}"
    namespace = "default"

    # Step 1: Create RegisteredModel
    print("Creating RegisteredModel...")
    model_name = f"test-model-{uuid.uuid4().hex[:8]}"
    rm = model_registry_client.register_model(
        name=model_name,
        uri="PLACEHOLDER",  # Will be updated by the job
        version="v1.0.0",
        model_format_name="onnx",
        model_format_version="1.0",
        description="Test model for async upload",
    )
    assert rm.id
    print(f"  Created RegisteredModel with ID: {rm.id}")

    # Get the created model version and artifact
    mv = model_registry_client.get_model_version(model_name, "v1.0.0")
    assert mv and mv.id
    print(f"  Created ModelVersion with ID: {mv.id}")

    ma = model_registry_client.get_model_artifact(model_name, "v1.0.0")
    assert ma and ma.id
    print(f"  Created ModelArtifact with ID: {ma.id}")

    # Verify initial state
    assert ma.uri == "PLACEHOLDER"
    assert ma.state == ArtifactState.UNKNOWN

    # Step 2: Download and upload ONNX model
    print("Downloading ONNX model...")

    # Step 3: Apply the job with patches
    print("Applying resources...")
    actual_job_name = apply_job_with_strategic_merge(
        rm_id=rm.id,
        mv_id=mv.id,
        ma_id=ma.id,
        container_image_uri=container_image_uri,
        k8s_client=k8s_client,
        env=env,
        tmp_path=tmp_path,
    )

    # Use the actual job name returned from kustomize
    if actual_job_name:
        job_name = actual_job_name

    print(f"Waiting for job completion: {job_name}")
    success = wait_for_job_completion(job_name, namespace, k8s_batch_client)

    if not success:
        # Get job logs for debugging
        try:
            pods = k8s_core_client.list_namespaced_pod(namespace=namespace, label_selector=f"job-name={job_name}")
            for pod in pods.items:
                logs = k8s_core_client.read_namespaced_pod_log(name=pod.metadata.name, namespace=namespace)
                print(f"Pod {pod.metadata.name} logs:\n{logs}")
        except Exception as e:
            print(f"Could not get pod logs: {e}")

        pytest.fail("Job did not complete successfully")

    print("Job completed successfully!")

    # Step 7: Validate the final artifact state
    print("Validating final artifact state...")
    # TODO compare resulting OCI image model files with expected

    # Wait a bit for the artifact to be updated
    time.sleep(2)

    # Fetch the updated artifact
    updated_ma = model_registry_client.get_model_artifact(model_name, "v1.0.0")
    assert updated_ma

    # Validate the artifact was updated correctly
    assert updated_ma.uri != "PLACEHOLDER", f"URI was not updated: {updated_ma.uri}"
    assert updated_ma.state == ArtifactState.LIVE, f"State was not updated to LIVE: {updated_ma.state}"

    print(f"✅ Artifact URI updated to: {updated_ma.uri}")
    print(f"✅ Artifact state updated to: {updated_ma.state}")

    # Clean up the job
    try:
        import subprocess

        subprocess.run(
            ["kubectl", "delete", "job", job_name, "-n", namespace],
            capture_output=True,
            check=False,
        )
    except Exception:
        pass  # Ignore cleanup errors

    print("Integration test completed successfully!")
