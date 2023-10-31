# Model Registry Python Client

[![Python](https://img.shields.io/badge/python%20-3.9%7C3.10-blue)](https://github.com/opendatahub-io/model-registry)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](../../../LICENSE)

This library provides a low level interface for interacting with a model registry server.

Types are based on [ML Metadata](https://github.com/google/ml-metadata), with Pythonic class wrappers.

## Basic usage

Registry objects can be created by doing

<!-- TODO: #120 Refer to types documentation -->


```py
from model_registry.types import ModelArtifact, ModelVersion, RegisteredModel

trained_model = ModelArtifact("my_model_name", "resource_URI",
                              description="Model description")

version = ModelVersion(trained_model, "v1.0", "model author")

model = RegisteredModel("my_model_name")
```

<!-- TODO: #120 provide a link to the reference docs instead of code -->
To register those objects, you can use the [`model_registry.ModelRegistry` class](src/model_registry/registry/client.py):

```py
from model_registry import ModelRegistry

registry = ModelRegistry("server-address", "port")

model_id = registry.upsert_registered_model(model)

# we need a model to associate the version to
version_id = registry.upsert_model_version(version, model_id)

# we need a version to associate an trained model to
experiment_id = registry.upsert_model_artifact(trained_model, version_id)
```

To get previously registered objects from the registry, use
```py
another_model = registry.get_registered_model_by_id("another-model-id")

another_version = registry.get_model_version_by_id("another-version-id", another_model.id)

another_experiment = registry.get_model_artifact_by_id("another-model-artifact-id")
```

## Development

Common tasks, such as building documentation and running tests, can be executed using [`nox`](https://github.com/wntrblm/nox) sessions.

Use `nox -l` to list sessions and execute them using `nox -s [session]`.

<!-- github-only -->
