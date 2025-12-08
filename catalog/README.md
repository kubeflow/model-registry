# Model Catalog Service

The Model Catalog Service provides a **read-only discovery service** for ML models across multiple catalog sources. It acts as a federated metadata aggregation layer, allowing users to search and discover models from various external catalogs through a unified REST API.

## Architecture Overview

The catalog service operates as a **metadata aggregation layer** that:
- Federates model discovery across different external catalogs
- Provides a unified REST API for model search and discovery
- Uses pluggable source providers for extensibility
- Operates without traditional database storage (file-based configuration)

### Supported Catalog Sources

- **YAML Catalog** - Static YAML files containing model metadata
- **HuggingFace Hub** - Discover models from HuggingFace's model repository

## REST API

### Base URL
`/api/model_catalog/v1alpha1`

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/sources` | List all catalog sources with pagination |
| `GET` | `/models` | Search models across sources (requires `source` parameter) |
| `GET` | `/sources/{source_id}/models/{model_name+}` | Get specific model details |
| `GET` | `/sources/{source_id}/models/{model_name}/artifacts` | List model artifacts |

### OpenAPI Specification

View the complete API specification:
- [Swagger UI](https://www.kubeflow.org/docs/components/model-registry/reference/model-catalog-rest-api/#swagger-ui)
- [Swagger Playground](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/kubeflow/model-registry/main/api/openapi/catalog.yaml)

## Data Models

### CatalogSource
Simple source metadata:
```json
{
  "id": "string",
  "name": "string"
}
```

### CatalogModel
Rich model metadata including:
- Basic info: `name`, `description`, `readme`, `maturity`
- Technical: `language[]`, `tasks[]`, `libraryName`
- Legal: `license`, `licenseLink`, `provider`
- Extensible: `customProperties` (key-value metadata)

### CatalogModelArtifact
Artifact references:
```json
{
  "uri": "string",
  "customProperties": {}
}
```

## Configuration

The catalog service uses **file-based configuration** instead of traditional databases:

```yaml
# catalog-sources.yaml
catalogs:
  - id: "yaml-catalog"
    name: "Local YAML Catalog"
    type: "yaml"
    properties:
      path: "./models"
```

### HuggingFace Source Configuration

The HuggingFace catalog source allows you to discover and import models from the HuggingFace Hub. To configure a HuggingFace source:

#### 1. Set Your API Key

Setting a Hugging Face API key is optional. Hugging Face  requires an API key for authentication for full access to data of models that are private and/or gated. If an API key is NOT set, private models will be entirely unavailable and gated models will have limited metadata. By default, the service reads the API key from the `HF_API_KEY` environment variable:

```bash
export HF_API_KEY="your-huggingface-api-key-here"
```

**Getting a HuggingFace API Key:**
1. Sign up or log in to [HuggingFace](https://huggingface.co)
2. Go to your [Settings > Access Tokens](https://huggingface.co/settings/tokens)
3. Create a new token with "Read" permissions
4. Copy the token and set it as an environment variable

**For Kubernetes deployments:**
- Store the API key in a Kubernetes Secret
- Reference it in your deployment configuration
- The catalog service will read it from the configured environment variable (defaults to `HF_API_KEY`)

**Custom Environment Variable Name:**
You can configure a custom environment variable name per source by setting the `apiKeyEnvVar` property in your source configuration (see below). This is useful when you need different API keys for different sources.

**Important Notes:**
- **Private Models**: For private models, the API key must belong to an account that has been granted access to the model. Without proper access, the catalog service will not be able to retrieve model information.
- **Gated Models**: For gated models (models with usage restrictions), you must accept the model's terms of service on HuggingFace before the catalog service can access all available model information. Visit the model's page on HuggingFace and accept the terms to ensure full metadata is available.

#### 2. Configure the Source

Add a HuggingFace source to your `catalog-sources.yaml`:

```yaml
catalogs:
  - name: "HuggingFace Hub"
    id: "huggingface"
    type: "hf"
    enabled: true
    # Required: List of model identifiers to include
    # Format: "organization/model-name" or "username/model-name"
    includedModels:
      - "meta-llama/Llama-3.1-8B-Instruct"
      - "ibm-granite/granite-4.0-h-small"
      - "microsoft/phi-2"
    
    # Optional: Exclude specific models or patterns
    # Supports exact matches or patterns ending with "*"
    excludedModels:
      - "some-org/unwanted-model"
      - "another-org/test-*"  # Excludes all models starting with "test-"
    
    # Optional: Configure a custom environment variable name for the API key
    # Defaults to "HF_API_KEY" if not specified
    properties:
      apiKeyEnvVar: "MY_CUSTOM_API_KEY_VAR"
```

#### Model Filtering

Both `includedModels` and `excludedModels` are top-level properties (not nested under `properties`):

- **`includedModels`** (required): List of model identifiers to fetch from HuggingFace. Format: `"organization/model-name"` or `"username/model-name"`
- **`excludedModels`** (optional): List of models or patterns to exclude from the results

The `excludedModels` property supports:
- **Exact matches**: `"meta-llama/Llama-3.1-8B-Instruct"` - excludes this specific model
- **Pattern matching**: `"test-*"` - excludes all models starting with "test-"

## Development

### Prerequisites
- Go >= 1.25
- Java >= 11.0 (for OpenAPI generation)
- Node.js >= 20.0.0 (for GraphQL schema downloads)

### Building

Generate OpenAPI server code:
```bash
make gen/openapi-server
```

Generate OpenAPI client code:
```bash
make gen/openapi
```

### Project Structure

```
catalog/
├── cmd/                    # Main application entry point
├── internal/
│   ├── catalog/           # Core catalog logic and providers
│   │   ├── genqlient/     # GraphQL client generation
│   │   └── testdata/      # Test fixtures
│   └── server/openapi/    # REST API implementation
├── pkg/openapi/           # Generated OpenAPI client
├── scripts/               # Build and generation scripts
└── Makefile              # Build targets
```

### Adding New Catalog Providers

1. Implement the `CatalogSourceProvider` interface:
```go
type CatalogSourceProvider interface {
    GetModel(ctx context.Context, name string) (*model.CatalogModel, error)
    ListModels(ctx context.Context, params ListModelsParams) (model.CatalogModelList, error)
    GetArtifacts(ctx context.Context, name string) (*model.CatalogModelArtifactList, error)
}
```

2. Register your provider:
```go
catalog.RegisterCatalogType("my-catalog", func(source *Source) (CatalogSourceProvider, error) {
    return NewMyCatalogProvider(source)
})
```

### Testing

The catalog service includes comprehensive testing:
- Unit tests for core catalog logic
- Integration tests for provider implementations
- OpenAPI contract validation

### Configuration Hot Reloading

The service automatically reloads configuration when the catalog sources file changes, enabling dynamic catalog updates without service restarts.

## Integration

The catalog service is designed to complement the main Model Registry service by providing:
- External model discovery capabilities
- Unified metadata aggregation
- Read-only access to distributed model catalogs

For complete Model Registry documentation, see the [main README](../README.md).
