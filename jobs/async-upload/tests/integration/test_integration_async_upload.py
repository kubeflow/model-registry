"""Integration tests for async-upload job functionality."""

import json
import os
import tempfile
import time
import uuid
from pathlib import Path
from typing import Dict, Any

import pytest
import requests
import yaml

# Import model registry client for API operations
try:
    from model_registry import ModelRegistry
    from model_registry.types import ArtifactState
except ImportError:
    pytest.skip("model-registry client not available", allow_module_level=True)


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


def apply_job_with_patches(base_job_path: str, patches: Dict[str, Any], k8s_client) -> None:
    """Apply job with patches using pure Python approach."""
    # Load the base job YAML
    with open(base_job_path) as f:
        job_docs = list(yaml.safe_load_all(f))
    
    # Apply patches to each document
    for doc in job_docs:
        if doc and doc.get("kind"):
            # Only apply patches to Job resources
            if doc.get("kind") == "Job":
                # Apply patches using Python dict operations
                for patch_path, patch_value in patches.items():
                    apply_patch_to_dict(doc, patch_path, patch_value)
            
            # Apply the resource using kubernetes client
            apply_resource(doc, k8s_client)


def apply_patch_to_dict(doc: Dict[str, Any], path: str, value: Any) -> None:
    """Apply a single patch to a dictionary using JSONPath-style path."""
    # Convert JSONPath to Python dict navigation
    # Example: "/metadata/name" -> ["metadata", "name"]
    # Example: "/spec/template/spec/containers/0/image" -> navigate to containers[0].image
    path_parts = [part for part in path.split("/") if part]
    
    # Navigate to the parent of the target
    current: Any = doc
    for part in path_parts[:-1]:
        # Check if part is a numeric index (for arrays)
        if part.isdigit():
            # It's an array index
            index = int(part)
            current = current[index]
        else:
            # It's a dictionary key
            if part not in current:
                current[part] = {}
            current = current[part]
    
    # Set the final value
    final_key = path_parts[-1]
    
    # Check if final key is a numeric index
    if final_key.isdigit():
        index = int(final_key)
        current[index] = value
    else:
        current[final_key] = value


def apply_resource(resource: Dict[str, Any], k8s_client) -> None:
    """Apply a single Kubernetes resource."""
    try:
        from kubernetes.dynamic import DynamicClient
    except ImportError:
        pytest.skip("kubernetes.dynamic not available")
    
    dyn_client = DynamicClient(k8s_client)
    
    # Get the API version and kind
    api_version = resource.get("apiVersion", "v1")
    kind = resource.get("kind")
    name = resource["metadata"]["name"]
    namespace = resource["metadata"].get("namespace", "default")
    
    # Create the resource
    api_resource = dyn_client.resources.get(api_version=api_version, kind=kind)
    
    # For Jobs, delete existing one first to avoid conflicts
    if kind == "Job":
        try:
            api_resource.delete(name=name, namespace=namespace)
            time.sleep(2)  # Wait for deletion to complete
        except Exception:
            pass  # Ignore errors if job doesn't exist
    
    # Apply the resource
    try:
        api_resource.create(body=resource)
    except Exception as e:
        if hasattr(e, 'status') and getattr(e, 'status', None) == 409:  # Conflict - resource already exists
            # For non-Job resources, just ignore if they already exist
            if kind != "Job":
                print(f"Resource {kind}/{name} already exists, skipping...")
                return
            else:
                raise  # Jobs should have been deleted above, so this is unexpected
        else:
            raise


def wait_for_job_completion(
    job_name: str,
    namespace: str,
    batch_client,
    timeout_seconds: int = 600
) -> bool:
    """Wait for job completion and return success status."""
    start_time = time.time()
    
    while time.time() - start_time < timeout_seconds:
        try:
            job = batch_client.read_namespaced_job(name=job_name, namespace=namespace)
            
            # Check if job is complete
            if job.status.conditions:
                for condition in job.status.conditions:
                    if condition.type == "Complete" and condition.status == "True":
                        return True
                    elif condition.type == "Failed" and condition.status == "True":
                        return False
            
            time.sleep(10)  # Wait 10 seconds before checking again
            
        except Exception as e:
            if hasattr(e, 'status') and getattr(e, 'status', None) == 404:
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
        's3',
        endpoint_url='http://localhost:9000',
        aws_access_key_id=access_key,
        aws_secret_access_key=secret_key,
        region_name='us-east-1'  # MinIO doesn't care about region but boto3 requires it
    )
    
    try:
        # Upload the file
        s3_client.upload_file(file_path, bucket, key)
    except ClientError as e:
        raise Exception(f"Failed to upload to MinIO: {e}")


@pytest.mark.integration
def test_async_upload_integration(
    k8s_client,
    k8s_batch_client,
    k8s_core_client,
    model_registry_client: ModelRegistry
):
    """Test the complete async-upload job integration.
    
    This test:
    1. Creates a RegisteredModel, ModelVersion, and placeholder ModelArtifact
    2. Downloads and uploads an ONNX model to MinIO
    3. Creates and applies a Kubernetes job using kustomize
    4. Waits for job completion
    5. Validates the final artifact state
    """
    
    # Configuration
    mr_host_url = os.environ.get("MR_HOST_URL", "http://localhost:8080")
    container_image_uri = os.environ.get(
        "CONTAINER_IMAGE_URI",
        "ghcr.io/kubeflow/model-registry/job/async-upload:latest"
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
        description="Test model for async upload"
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
    with tempfile.TemporaryDirectory() as temp_dir:
        model_file = Path(temp_dir) / "mnist-8.onnx"
        
        # Download the model
        response = requests.get(
            "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/classification/mnist/model/mnist-8.onnx"
        )
        response.raise_for_status()
        
        with open(model_file, "wb") as f:
            f.write(response.content)
        
        print("Uploading to MinIO...")
        bucket = "default"
        key = f"my-model/mnist-8.onnx"
        upload_to_minio(str(model_file), bucket, key)
        
        # Step 3: Prepare job patches
        print("Preparing job patches...")
        base_job_path = Path(__file__).parent.parent.parent / "samples" / "sample_job_s3_to_oci.yaml"
        
        # Define patches to apply to the job
        patches = {
            "/metadata/name": job_name,
            "/spec/template/spec/containers/0/image": container_image_uri,
            "/spec/template/spec/containers/0/env/12/value": rm.id,
            "/spec/template/spec/containers/0/env/13/value": mv.id,
            "/spec/template/spec/containers/0/env/14/value": ma.id,
        }
        
        # Step 4: Apply resources (job cleanup is handled automatically)
        print("Applying resources...")
        
        # Step 5: Apply the job with patches
        apply_job_with_patches(str(base_job_path), patches, k8s_client)
        
        # Step 6: Wait for job completion
        print("Waiting for job completion...")
        success = wait_for_job_completion(job_name, namespace, k8s_batch_client, timeout_seconds=600)
        
        if not success:
            # Get job logs for debugging
            try:
                pods = k8s_core_client.list_namespaced_pod(
                    namespace=namespace,
                    label_selector=f"job-name={job_name}"
                )
                for pod in pods.items:
                    logs = k8s_core_client.read_namespaced_pod_log(
                        name=pod.metadata.name,
                        namespace=namespace
                    )
                    print(f"Pod {pod.metadata.name} logs:\n{logs}")
            except Exception as e:
                print(f"Could not get pod logs: {e}")
            
            pytest.fail("Job did not complete successfully")
        
        print("Job completed successfully!")
        
        # Step 7: Validate the final artifact state
        print("Validating final artifact state...")
        
        # Wait a bit for the artifact to be updated
        time.sleep(5)
        
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
            k8s_batch_client.delete_namespaced_job(name=job_name, namespace=namespace)
        except Exception:
            pass  # Ignore cleanup errors
        
        print("Integration test completed successfully!") 