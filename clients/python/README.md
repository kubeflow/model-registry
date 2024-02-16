# Model Registry Python Client

[![Python](https://img.shields.io/badge/python%20-3.9%7C3.10-blue)](https://github.com/opendatahub-io/model-registry)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](../../../LICENSE)

This library provides a high level interface for interacting with a model registry server.

## Basic usage

```py
from model_registry import ModelRegistry

registry = ModelRegistry(server_address="server-address", port=9090, author="author")

model = registry.register_model(
    "my-model",  # model name
    "s3://path/to/model",  # model URI
    version="2.0.0",
    description="lorem ipsum",
    model_format_name="onnx",
    model_format_version="1",
    storage_key="aws-connection-path",
    storage_path="path/to/model",
    metadata={
        # can be one of the following types
        "int_key": 1,
        "bool_key": False,
        "float_key": 3.14,
        "str_key": "str_value",
    }
)

model = registry.get_registered_model("my-model")

version = registry.get_model_version("my-model", "v2.0")

experiment = registry.get_model_artifact("my-model", "v2.0")
```

### Default values for metadata

If not supplied, `metadata` values defaults to a predefined set of conventional values.
Reference the technical documentation in the pydoc of the client.

### Importing from Hugging Face Hub

To import models from Hugging Face Hub, start by installing the `huggingface-hub` package, either directly or as an
extra (available as `model-registry[hf]`).
Models can be imported with

```py
hf_model = registry.register_hf_model(
    "hf-namespace/hf-model",  # HF repo
    "relative/path/to/model/file.onnx",
    version="1.2.3",
    model_name="my-model",
    description="lorem ipsum",
    model_format_name="onnx",
    model_format_version="1",
)
```

There are caveats to be noted when using this method:

- It's only possible to import a single model file per Hugging Face Hub repo right now.
- If the model you want to import is in a global namespace, you should provide an author, e.g.

    ```py
    hf_model = registry.register_hf_model(
        "gpt2",  # this model implicitly has no author
        "onnx/decoder_model.onnx",
        author="OpenAI",  # Defaults to unknown in the absence of an author
        version="1.0.0",
        description="gpt-2 model",
        model_format_name="onnx",
        model_format_version="1",
    )
    ```

## Development

Common tasks, such as building documentation and running tests, can be executed using [`nox`](https://github.com/wntrblm/nox) sessions.

Use `nox -l` to list sessions and execute them using `nox -s [session]`.

### Running Locally on Mac M1 or M2 (arm64 architecture)

If you want run tests locally you will need to set up a colima develeopment environment using the instructions [here](https://github.com/opendatahub-io/model-registry/blob/main/CONTRIBUTING.md#colima)

You will also have to change the package source to one compatible with ARM64 architecture. This can be actioned by uncommenting lines 14 or 15 in the pyproject.toml file. Run the following command after you have uncommented the line.

```sh
poetry lock
```
Use the following commands to directly run the tests with individual test output. Alternatively you can use the nox session commands above.

```sh
poetry install
poetry run pytest -v
```

<!-- github-only -->
