# Model Registry Python Client

[![Python](https://img.shields.io/badge/python%20-3.9%7C3.10-blue)](https://github.com/kubeflow/model-registry)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](../../../LICENSE)

This library provides a high level interface for interacting with a model registry server.

## Basic usage

```py
from model_registry import ModelRegistry

registry = ModelRegistry("server-address", author="Ada Lovelace")  # Defaults to a secure connection via port 443

# registry = ModelRegistry("server-address", 1234, author="Ada Lovelace", is_secure=False)  # To use MR without TLS

model = registry.register_model(
    "my-model",  # model name
    "https://storage-place.my-company.com",  # model URI
    version="2.0.0",
    description="lorem ipsum",
    model_format_name="onnx",
    model_format_version="1",
    storage_key="my-data-connection",
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

### Importing from S3

When registering models stored on S3-compatible object storage, you should use `utils.s3_uri_from` to build an
unambiguous URI for your artifact.

```py
from model_registry import ModelRegistry, utils

registry = ModelRegistry(server_address="server-address", port=9090, author="author")

model = registry.register_model(
    "my-model",  # model name
    uri=utils.s3_uri_from("path/to/model", "my-bucket"),
    version="2.0.0",
    description="lorem ipsum",
    model_format_name="onnx",
    model_format_version="1",
    storage_key="my-data-connection",
    metadata={
        # can be one of the following types
        "int_key": 1,
        "bool_key": False,
        "float_key": 3.14,
        "str_key": "str_value",
    }
)
```

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

## Advanced use-cases

### Using Model Registry Python Client with newer Python versions (>=3.11)

> [!CAUTION]
> The mechanism described in this section is a temporary workaround and likely will never be supported.
> This workaround is ONLY applicable if your Python/Notebook project does NOT make use of MLMD directly or indirectly.

<!-- a longer-term plan to address this ties to the investigations to rebase this client on top of MR REST api directly,
so to avoid having to wrap the MLMD Wheel. See more: https://github.com/kubeflow/model-registry/pull/59 -->

This project _currently_ depends for internal implementations on the Google's [MLMD Python library](https://pypi.org/project/ml-metadata/).
Due to this dependency, this project supports [only the Python versions](https://github.com/kubeflow/model-registry/blob/8d77c13100c6cc5a9465d4293403114a3576fdd7/clients/python/pyproject.toml#L14) which are also available for the MLMD library (see more [here](https://pypi.org/project/ml-metadata/#files)).

As a workaround, **only IF your Python/Notebook project does NOT make use of MLMD directly or indirectly**,
you could opt-in to make use of a non-constrained variant of the MLMD dependency supporting _only_ remote gRPC calls (and not constrained by specific Python versions or architectures):

```
!pip install "https://github.com/opendatahub-io/ml-metadata/releases/download/v1.14.0%2Bremote.1/ml_metadata-1.14.0+remote.1-py3-none-any.whl" # need a Python 3.11 compatible version
!pip install --no-deps --ignore-requires-python --pre "model-registry" # ignore dependencies because of the above override
```

You can read more about this use-case, in the [Remote-only packaging of MLMD Python lib](https://github.com/kubeflow/model-registry/blob/main/docs/remote_only_packaging_of_MLMD_Python_lib.md) document.

## Development

Common tasks, such as building documentation and running tests, can be executed using [`nox`](https://github.com/wntrblm/nox) sessions.

Use `nox -l` to list sessions and execute them using `nox -s [session]`.

### Running Locally on Mac M1 or M2 (arm64 architecture)

If you want run tests locally you will need to set up a development environment, including docker engine; we recommend following the instructions [here](https://github.com/kubeflow/model-registry/blob/main/CONTRIBUTING.md#docker-engine).

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
