# Model Catalog stub implementation resources

This directory contains metadata for the initial stub iteration of the model catalog in the RHOAI dashboard. This stub implementation allows us to develop and demonstrate model catalog UI functionality without a real backend service. The model source files and scripts here will become obsolete once the upcoming model catalog backend service is implemented.

For the initial stub implementation, each of the YAML files in `./models` represents a model catalog source object. To include these models in the model catalog UI, this data is stored in a ConfigMap on a cluster. The YAML for each source is converted to JSON and inserted as an element in the `sources` array within the JSON blob `data.modelCatalogSources` in the ConfigMap. Multiple sources will appear in the UI under section headers.

There are two such ConfigMaps automatically created by the [manifests in the odh-dashboard repository](https://github.com/opendatahub-io/odh-dashboard/blob/main/manifests/rhoai/shared/apps/model-catalog). The manifests are specific to RHOAI and only automatically create these ConfigMaps on RHOAI clusters and not ODH. They are created in the application namespace (`redhat-ods-applications`). To test this functionality in an ODH cluster, you can create these configmaps manually in the `opendatahub` namespace based on these manifests.

- `model-catalog-sources`

  This is a managed resource containing the model sources shipped in product releases. [The manifests](https://github.com/opendatahub-io/odh-dashboard/blob/main/manifests/rhoai/shared/apps/model-catalog/model-catalog-configmap.yaml) contain the metadata content for these sources and their models. Edits to this ConfigMap on a RHOAI cluster will not persist, and upgrading the platform will replace its contents if they have changed.

- `model-catalog-unmanaged-sources`

  This is an unmanaged resource which by default contains an empty `sources` array in its JSON. Sources can be added to the model catalog on a cluster by editing this ConfigMap. Its contents will be combined with the managed ConfigMap (the lists of sources are appended together).

## Updating unmanaged sources in a cluster

The `scripts/update-unmanaged-sources-configmap.sh` script updates the `model-catalog-unmanaged-sources` ConfigMap with a new model catalog source. It takes one of the `models/*.yaml` files as input, converts it to JSON, and updates the ConfigMap in the cluster.

### Usage

```bash
./scripts/update-unmanaged-sources-configmap.sh <input-yaml-file> [namespace]
```

- `input-yaml-file`: Path to the YAML file containing the source configuration
- `namespace`: (Optional) The namespace where the ConfigMap exists. Defaults to `redhat-ods-applications`. Can be changed to `opendatahub` if needed.

### Examples

Using the default namespace (redhat-ods-applications):

```bash
./scripts/update-unmanaged-sources-configmap.sh models/neural-magic-models.yaml
```

Using the opendatahub namespace:

```bash
./scripts/update-unmanaged-sources-configmap.sh models/neural-magic-models.yaml opendatahub
```

### Prerequisites

- `yq` and `jq` must be installed
- Must be logged into an OpenShift cluster using `oc login`
- The ConfigMap `model-catalog-unmanaged-sources` must exist in the specified namespace
