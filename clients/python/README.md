# Model Registry Python Client

[![Python](https://img.shields.io/badge/python%20-3.9%7C3.10-blue)](https://github.com/opendatahub-io/model-registry)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](../../../LICENSE)

This library provides a high level interface for interacting with a model registry server.

## Basic usage

```py
from model_registry import ModelRegistry

registry = ModelRegistry(address="server-address", port="port", author="author")

model = registry.register_model("my-model",
                                "s3://path/to/model",
                                format_name="onnx",
                                format_version="v1",
                                secret_account_name="access",
                                version="v2.0",
                                description="lorem ipsum",
                                )

model = registry.get_registered_model("my-model")

version = registry.get_model_version("my-model", "v2.0")

experiment = registry.get_model_artifact("my-model", "v2.0")
```

## Development

Common tasks, such as building documentation and running tests, can be executed using [`nox`](https://github.com/wntrblm/nox) sessions.

Use `nox -l` to list sessions and execute them using `nox -s [session]`.

<!-- github-only -->
