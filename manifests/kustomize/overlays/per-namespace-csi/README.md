# Per-Namespace CSI Overlay

This overlay deploys Model Registry with the CSI `ClusterStorageContainer` for a single Kubeflow profile namespace.

## Usage

Update the `namespace:` field in `kustomization.yaml` to the target profile namespace, then apply the overlay:

```sh
kubectl apply -k manifests/kustomize/overlays/per-namespace-csi
```

To deploy into multiple profiles, repeat the same step after changing the `namespace:` field each time, for example once for `profile-alpha` and again for `profile-beta`.

This overlay uses the `../db` backend by default. If you prefer PostgreSQL, replace `../db` with `../postgres` in `kustomization.yaml`.

## Multi-Profile Notes

`ClusterStorageContainer` is cluster-scoped, so Kubernetes can only have one `model-registry-storage-initializer` resource at a time. In a multi-profile deployment:

- `MODEL_REGISTRY_BASE_URL` is only a shared fallback default.
- The most recently applied overlay instance sets that fallback value.
- Per-profile isolation should use an embedded Model Registry service URL in each `InferenceService` `storageUri`.

Use this URI pattern:

```text
model-registry://model-registry-service.<PROFILE_NAMESPACE>.svc.cluster.local:8080/<MODEL_NAME>/<VERSION_NAME>
```

Example:

```text
model-registry://model-registry-service.profile-alpha.svc.cluster.local:8080/iris/v1
```

That fully qualified URI tells the CSI storage initializer which Model Registry service to query, even when the shared `ClusterStorageContainer` default points at a different namespace.

See [cmd/csi/GET_STARTED.md](../../../../cmd/csi/GET_STARTED.md) for a full `InferenceService` example and additional guidance.
