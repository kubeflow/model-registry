# Kubeflow Model Registry – Async Upload Job

> Lightweight, non‑root Python Job containerised for **transferring model artefacts** between storage back‑ends (S3, OCI, PVC, …) and registering them in **Kubeflow Model Registry**. Born out of the discussion in [kubeflow/model-registry #1108](https://github.com/kubeflow/model-registry/issues/1108).

---

## Quick start

```bash
# 1 – Build & tag ( reproducible build args are optional )
docker build \
  --build-arg VCS_REF=$(git rev-parse --short HEAD) \
  --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  -t quay.io/<org>/async-upload-job:$(git rev-parse --short HEAD) .

# 2 – Push to your registry
docker push quay.io/<org>/async-upload-job:<tag>

# TODO: Run locally...
```

## References

- Issue thread : [https://github.com/kubeflow/model-registry/issues/1108](https://github.com/kubeflow/model-registry/issues/1108)
- OCI Image Spec : [https://github.com/opencontainers/image-spec](https://github.com/opencontainers/image-spec)
- Kubernetes Pod Security : [https://kubernetes.io/docs/concepts/security/pod-security-standards/](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
