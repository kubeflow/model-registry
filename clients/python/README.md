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

## Development

To build the documentation, run `nox -s docs-build`.

<!-- github-only -->
