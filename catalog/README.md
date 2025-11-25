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

The HuggingFace provider requires an API key for authentication. Set the `HF_API_KEY` environment variable:

```bash
export HF_API_KEY="your-huggingface-api-key-here"
```

**Getting a HuggingFace API Key:**
1. Sign up or log in to [HuggingFace](https://huggingface.co)
2. Go to your [Settings > Access Tokens](https://huggingface.co/settings/tokens)
3. Create a new token with "Read" permissions
4. Copy the token and set it as the `HF_API_KEY` environment variable

**For Kubernetes deployments:**
- Store the API key in a Kubernetes Secret
- Reference it in your deployment configuration
- The catalog service will read it from the `HF_API_KEY` environment variable

#### 2. Configure the Source

Add a HuggingFace source to your `catalog-sources.yaml`:

```yaml
catalogs:
  - name: "HuggingFace Hub"
    id: "huggingface"
    type: "hf"
    enabled: true
    properties:
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
```

#### Excluded Models

The `excludedModels` property supports:
- **Exact matches**: `"meta-llama/Llama-3.1-8B-Instruct"` - excludes this specific model
- **Pattern matching**: `"test-*"` - excludes all models starting with "test-"

## Development

### Prerequisites
- Go >= 1.24
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
