# Model Catalog Manifests

To deploy the model catalog:

```sh
kubectl apply -k . -n NAMESPACE
```

Replace `NAMESPACE` with your desired Kubernetes namespace.

## sources.yaml Configuration

The `sources.yaml` file configures the model catalog sources. It contains a top-level `catalogs` list, where each entry defines a single catalog source.

### Common Properties

Each catalog source entry supports the following common properties:

- **`name`** (*string*, required): A user-friendly name for the catalog source.
- **`id`** (*string*, required): A unique identifier for the catalog source.
- **`type`** (*string*, required): The type of catalog source. There is currently one supported type: `yaml`.
- **`enabled`** (*boolean*, optional): Whether the catalog source is enabled. Defaults to `true` if not specified.

### Catalog Source Types

Below are the `properties` supported for the yaml type source.

#### `yaml`

The `yaml` type sources model metadata from a local YAML file.

##### Properties

- **`yamlCatalogPath`** (*string*, required): The path to the YAML file containing the model definitions. This path is relative to the directory where the `sources.yaml` file is located.
- **`excludedModels`** (*string list*, optional): A list of models to exclude from the catalog. These can be an exact name with a tag (e.g., `model-a:1.0`) or a pattern ending with `*` to exclude all tags for a repository (e.g., `model-b:*`).

##### Example

```yaml
catalogs:
  - name: Sample Catalog
    id: sample_custom_catalog
    type: yaml
    enabled: true
    properties:
      yamlCatalogPath: sample-catalog.yaml
      excludedModels:
      - model-a:1.0
      - model-b:*
```
