# Kubeflow Model Registry UI BFF
The Kubeflow Model Registry UI BFF is the _backend for frontend_ (BFF) used by the Kubeflow Model Registry UI.

# Building and Deploying
TBD

# Development
TBD

## Getting started

### Endpoints

| URL Pattern                                                   | Handler                 | Action                                       |
|---------------------------------------------------------------|-------------------------|----------------------------------------------|
| GET /v1/healthcheck                                           | HealthcheckHandler      | Show application information.                |
| GET /v1/model-registry/                                       | ModelRegistryHandler    | Get all model registries,                    |
| GET /v1/model-registry/{model_registry_id}/registered_models  | RegisteredModelsHandler | Gets a list of all RegisteredModel entities. |
| POST /v1/model-registry/{model_registry_id}/registered_models | RegisteredModelsHandler | Create a RegisteredModel entity.             |

### Sample local calls
```
# GET /v1/healthcheck
curl -i localhost:4000/api/v1/healthcheck/
```
```
# GET /v1/model-registry/ 
curl -i localhost:4000/api/v1/model-registry/
```
```
# GET /v1/model-registry/{model_registry_id}/registered_models
curl -i localhost:4000/api/v1/model-registry/model-registry/registered_models
```
```
#POST /v1/model-registry/{model_registry_id}/registered_models
curl -i -X POST "http://localhost:4000/api/v1/model-registry/model-registry/registered_models" \
     -H "Content-Type: application/json" \
     -d '{
  "customProperties": {
    "my-label9": {
      "metadataType": "MetadataStringValue",
      "string_value": "val"
    }
  },
  "description": "bella description",
  "externalId": "9927",
  "name": "bella",
  "owner": "eder",
  "state": "LIVE"
}'
```