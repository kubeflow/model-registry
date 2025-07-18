# Model Catalog Manifests

To deploy the model catalog:

```sh
kubectl apply -k . -n NAMESPACE
```

Replace `NAMESPACE` with your desired Kubernetes namespace.

To configure customized models, update `sources.yaml` and `sample-catalog.yaml`.

To add Red Hat Ecosystem Catalog models, update `sources.yaml` and add each path as a separate repository under `models`.
