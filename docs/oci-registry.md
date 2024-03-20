# OCI Registry as a Kubeflow Model Registry

## Authors

- Ramkumar Chinchani (Cisco)
- _TBD_

## Maintainers

- Ramkumar Chinchani (Cisco)
- _TBD_

## Motivation

According to the [Kubeflow 2023
survey](https://blog.kubeflow.org/kubeflow-user-survey-2023/), 44% of users
identified Model Registry as one of the big gaps in the user’s ML Lifecycle
missing from the Kubeflow offering. 

![Kubeflow survey](diagrams/model-registry-kubeflowsurvey.png "Kubeflow survey")

## Solution Overview

[Open Container Initiative](https://opencontainers.org/) is a sibling (to CNCF)
organization under [The Linux Foundation](https://www.linuxfoundation.org/)
which has the container
[runtime](https://github.com/opencontainers/runtime-spec),
[image](https://github.com/opencontainers/image-spec) and
[distribution](https://github.com/opencontainers/distribution-spec)
specifications under its purvey which are vendor-neutral contracts that the Kubernetes
ecosystem relies on for running, filesystem layout, and pushing and pulling of
container images.

However, recent developments in the OCI, specifically
[_image_](https://github.com/opencontainers/image-spec/releases/tag/v1.1.0) and
[_distribution_](https://github.com/opencontainers/distribution-spec/releases/tag/v1.1.0)
spec **v1.1.0**, have included support for pushing arbitrary artifacts along
with support for relationships between artifacts.

## OCI v1.1.0 Conformant Registries

The following are the highlights about OCI artifact registries.

- Container images: these represent workloads and have been the traditional use
  case for an OCI conformant registry.

- Artifacts: these represent arbitrary data (ML model data or additional
  metadata in this context) that can also be pushed and pulled from an OCI
  conformant registry.

- Content-addressable: all data is organized as a Merkle DAG with SHA256-hashed
  blobs. This bodes well for reproducibility.

- Versioning: apart from the SHA256 hash, all data can be tagged with a human-readable version.

- Annotations: there is provision to append arbitrary annotations to any artifact.

- References: an artifact can be pushed along with a reference to another
  artifact (via the `Subject` field) which can be leveraged to address the data
  lineage use case.

- Provenance: each artifact can be cryptographically signed, with the signature
  as its own separate artifact "referring" to the signed artifact.

- Ecosystem tooling: there are OCI v1.1.0 conformant registries and clients
  already available which can be leveraged.

- Infrastructure reuse: a container image registry is already a critical piece of
  infrastructure which can now be reused.


## References

_TBD_

# Appendix

The following section demonstrates the supported workflow.

NOTE: This section is not an endorsement of all of the tools used but merely
represents a demonstration that functioning tools already exist. Readers are
free to pick and choose any tool as they see fit with the requirement that the
choice should be OCI v1.1.0 conformant.

[`zot`](https://zotregistry.dev) is chosen as the registry and
[`regctl`](https://github.com/regclient/regclient) as the client.

## Start a registry

```bash
podman run -p 5000:5000 ghcr.io/project-zot/zot-linux-amd64:latest
```

## Download model data

```bash
curl -v -L0 https://github.com/tarilabs/demo20231212/raw/main/v1.nb20231206162408/mnist.onnx -o mnist.onnx
```

## Upload model data with annotations

```bash
regctl artifact put \
  --annotation description="used for demo purposes" \
  --annotation model_format_name="onnx" \
  --annotation model_format_version="1" \
  --artifact-type "application/vnd.model.type" \
  localhost:5000/models/my-model-from-gh:v1 \
  -f mnist.onnx
```

## List all artifacts

```bash
regctl artifact list localhost:5000/models/my-model-from-gh:v1 --format '{{jsonPretty .}}
```

## Filter by artifact type

```bash
regctl artifact list --filter-artifact-type "application/vnd.model.type" localhost:5000/models/my-model-from-gh:v1 --format '{{jsonPretty .}}'
```

## Download model data

```bash
regctl artifact get localhost:5000/models/my-model-from-gh:v1 > mnist.onnx.copy
```
