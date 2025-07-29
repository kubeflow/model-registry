# Model Catalog Manifests

To deploy the model catalog:

```sh
kubectl apply -k . -n NAMESPACE
```

Replace `NAMESPACE` with your desired Kubernetes namespace.

## Configure Custom Models
To configure customized models, update `sources.yaml` and `sample-catalog.yaml`.

## [Adding Red Hat Ecosystem Catalog models](https://github.com/kubeflow/model-registry/blob/679114f0e9cca631e4c16166affa0966c6a371ff/catalog/internal/catalog/genqlient/README.md?plain=1#L1)
To add Red Hat Ecosystem Catalog models, update `sources.yaml` and add each path as a separate repository under `models`.

