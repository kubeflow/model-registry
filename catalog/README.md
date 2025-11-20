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
