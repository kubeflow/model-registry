# Model Registry Python Client

[![Python](https://img.shields.io/badge/python%20-3.9%7C3.10%7C3.11%7C3.12-blue)](https://github.com/kubeflow/model-registry)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](../../../LICENSE)
[![Read the Docs](https://img.shields.io/readthedocs/model-registry)](https://model-registry.readthedocs.io/en/latest/)
[![Tutorial Website](https://img.shields.io/badge/Website-green?style=plastic&label=Tutorial&labelColor=blue)](https://www.kubeflow.org/docs/components/model-registry/getting-started/)

This library provides a high level interface for interacting with a model registry server.

## Installation

In your Python environment, you can install the latest version of the Model Registry Python client with:

```
pip install --pre model-registry
```

### Installing extras

Some capabilities of this Model Registry Python client, such as [importing model from Hugging Face](#importing-from-hugging-face-hub),
require additional dependencies.

By [installing an extra variant](https://packaging.python.org/en/latest/tutorials/installing-packages/#installing-extras) of this package
the additional dependencies will be managed for you automatically, for instance with:

```
pip install --pre "model-registry[hf]"
```

This step is not required if you already installed the additional dependencies already, for instance with:

```
pip install huggingface-hub
```

## Basic usage

```py
from model_registry import ModelRegistry

registry = ModelRegistry("https://server-address", author="Ada Lovelace")  # Defaults to a secure connection via port 443

# registry = ModelRegistry("http://server-address", 1234, author="Ada Lovelace", is_secure=False)  # To use MR without TLS

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

version = registry.get_model_version("my-model", "2.0.0")

experiment = registry.get_model_artifact("my-model", "2.0.0")

# change is not reflected on pushed model version
version.description = "Updated model version"

# you can update it using
registry.update(version)
```

### Importing from S3

When registering models stored on S3-compatible object storage, you should use `utils.s3_uri_from` to build an
unambiguous URI for your artifact.

```py
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
extra (available as `model-registry[hf]`). Reference section "[installing extras](#installing-extras)" above for
more information.

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

### Listing models

To list models you can use
```py
for model in registry.get_registered_models():
    ...

# and versions associated with a model
for version in registry.get_model_versions("my-model"):
    ...
```

You can also use `order_by_creation_time`, `order_by_update_time`, or `order_by_id` to change the sorting order

```py
latest_updates = registry.get_model_versions("my-model").order_by_update_time().descending()
for version in latest_updates:
    ...
```

By default, all queries will be `ascending`, but this method is also available for explicitness.

> Note: You can also set the `page_size()` that you want the Pager to use when invoking the Model Registry backend.
> When using it as an iterator, it will automatically manage pages for you.

#### Implementation notes

The pager will manage pages for you in order to prevent infinite looping.
Currently, the Model Registry backend treats model lists as a circular buffer, and **will not end iteration** for you.

## Development

Common tasks, such as building documentation and running tests, can be executed using [`nox`](https://github.com/wntrblm/nox) sessions.

Use `nox -l` to list sessions and execute them using `nox -s [session]`.

Alternatively, use `make install` to setup a local Python virtual environment with `poetry`.

To run the tests you will need `docker` (or equivalent) and the `compose` extension command.
This is necessary as the test suite will manage a Model Registry server and an MLMD instance to ensure a clean state on
each run.
You can use `make test` to execute `pytest`.

### Running Locally on Mac M1 or M2 (arm64 architecture)

Check out our [recommendations on setting up your docker engine](https://github.com/kubeflow/model-registry/blob/main/CONTRIBUTING.md#docker-engine) on an ARM processor.

<!-- github-only -->
