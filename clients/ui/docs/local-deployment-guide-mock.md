# Local Deployment Guide (Mock BFF)

## Overview

This guide walks you through setting up a complete local development environment for the Model Registry UI using Docker/Podman Compose. This setup includes:

* PostgreSQL database for model registry storage
* Model Registry server with direct database connection
* Model Catalog service for browsing community models
* Mock BFF service for UI backend functionality
* React frontend with development server

This approach provides a fully functional environment without requiring a Kubernetes cluster.

## Prerequisites

The following tools need to be installed in your local environment:

* Docker or Podman - [Docker Instructions](https://www.docker.com) or [Podman Instructions](https://podman.io)
* jq - JSON command line processor for API testing
* curl - for making HTTP requests

Note: All tools can be installed using your OS package manager.

## Setup

### 1. Build and Start Services

From the root directory of the model-registry repository, run:

```shell
docker-compose -f docker-compose-dev.yaml up --build
```

Or if using Podman:

```shell
podman-compose -f docker-compose-dev.yaml up --build
```

This will start the following services:
- **postgres**: PostgreSQL database on port 5432
- **model-registry**: Main model registry server on port 8080
- **model-catalog**: Catalog service on port 8081
- **mock-bff**: Mock backend-for-frontend service on port 4000
- **ui-frontend**: React development server on port 3000

Wait for all services to start. You should see logs indicating that each service is ready.

### 2. Verify Services

Check that all services are running:

```shell
# Check model registry API
curl http://localhost:8080/api/model_registry/v1alpha3/registered_models

# Check model catalog
curl http://localhost:8081/api/model_catalog/v1alpha1/sources

# Check mock BFF
curl http://localhost:4000/api/v1/user

# Check UI (should redirect to login or show the interface)
curl -I http://localhost:3000
```

### 3. Access the UI

Open your web browser and navigate to:

```
http://localhost:3000/model-registry
```

You should see the Model Registry UI interface.

## Testing with Sample Data

### Creating a Registered Model

First, create a registered model and capture its ID:

```shell
MODEL_ID=$(curl -s -X POST http://localhost:8080/api/model_registry/v1alpha3/registered_models \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-example-model",
    "description": "A test model for demonstration"
  }' | jq -r '.id')

echo "Created model with ID: $MODEL_ID"
```

### Creating a Model Version

Using the captured model ID, create a version:

```shell
VERSION_ID=$(curl -s -X POST http://localhost:8080/api/model_registry/v1alpha3/registered_models/$MODEL_ID/versions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "v1.0.0",
    "description": "Initial version of the test model",
    "author": "Model Developer",
    "registeredModelId": "'$MODEL_ID'"
  }' | jq -r '.id')

echo "Created version with ID: $VERSION_ID"
```

### Creating a Model Artifact

Add an artifact to the model version:

```shell
ARTIFACT_ID=$(curl -s -X POST http://localhost:8080/api/model_registry/v1alpha3/model_versions/$VERSION_ID/artifacts \
  -H "Content-Type: application/json" \
  -d '{
    "name": "model-artifact",
    "description": "Trained model artifact",
    "uri": "s3://my-bucket/models/my-test-model/v1.0.0/model.tar.gz",
    "artifactType": "model-artifact"
  }' | jq -r '.id')

echo "Created artifact with ID: $ARTIFACT_ID"
```

### Verify Data in UI

1. Refresh the UI at `http://localhost:3000/model-registry`
2. You should now see your test model "my-test-model" in the models list
3. Click on the model to view its details and versions
4. Navigate to the version to see the associated artifacts

## Troubleshooting

### Services Not Starting

If services fail to start, check the logs:

```shell
docker-compose -f docker-compose-dev.yaml logs [service-name]
```

Common issues:
- **Port conflicts**: Ensure ports 3000, 4000, 5432, 8080, and 8081 are available
- **Docker/Podman issues**: Restart Docker/Podman daemon

### UI Not Loading

1. Verify the mock-bff service is running:
   ```shell
   curl http://localhost:4000/api/v1/user
   ```

2. Check browser console for JavaScript errors
3. Ensure all environment variables are properly set in the docker-compose file

### API Calls Failing

1. Verify model-registry service is responding:
   ```shell
   curl http://localhost:8080/api/model_registry/v1alpha3/registered_models
   ```

2. Check that the mock-bff is properly proxying requests:
   ```shell
   curl http://localhost:4000/model-registry/api/v1/model_registry/model-registry/registered_models
   ```

### Database Issues

1. Check PostgreSQL logs:
   ```shell
   docker-compose -f docker-compose-dev.yaml logs postgres
   ```

2. Verify database connectivity:
   ```shell
   docker-compose -f docker-compose-dev.yaml exec postgres psql -U modelregistry -d modelregistry -c "\dt"
   ```

## Stopping the Environment

To stop all services:

```shell
docker-compose -f docker-compose-dev.yaml down
```

To also remove volumes (this will delete all data):

```shell
docker-compose -f docker-compose-dev.yaml down -v
```

## Development Workflow

This setup is ideal for:
- Frontend development testing
- API integration testing
- Full-stack feature development
- UI/UX testing with real data

The mock-bff service handles the translation between the UI's expected BFF format and the model-registry's native API format, eliminating the need for a Kubernetes cluster during development.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   UI Frontend   │───▶│    Mock BFF     │───▶│ Model Registry  │
│   (Port 3000)   │    │   (Port 4000)   │    │   (Port 8080)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                 │                       │
                                 │                       ▼
                                 │              ┌─────────────────┐
                                 │              │   PostgreSQL    │
                                 │              │   (Port 5432)   │
                                 │              └─────────────────┘
                                 ▼
                        ┌─────────────────┐
                        │ Model Catalog   │
                        │   (Port 8081)   │
                        └─────────────────┘
```
