# Kubeflow Model Registry UI BFF
The Kubeflow Model Registry UI BFF is the _backend for frontend_ (BFF) used by the Kubeflow Model Registry UI.

# Building and Deploying
TBD

# Development
TBD

## Getting started

### Endpoints

| URL Pattern             | Handler              | Action                        |
|-------------------------|----------------------|-------------------------------|
| GET /v1/healthcheck     | HealthcheckHandler   | Show application information. |
| GET /v1/model-registry/ | ModelRegistryHandler | Get all model registries,     |

### Sample local calls
```
# GET /v1/healthcheck
curl -i localhost:4000/api/v1/healthcheck/
```