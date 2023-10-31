# Model Registry Python Client

[![Python](https://img.shields.io/badge/python%20-3.9%7C3.10-blue)](https://github.com/opendatahub-io/model-registry)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](../../../LICENSE)

This library provides a low level interface for interacting with a model registry server.

Types are based on [ML Metadata](https://github.com/google/ml-metadata), with Pythonic class wrappers.

## Basic usage

### Create objects

Registry objects can be created by doing

<!-- TODO: #120 Refer to types documentation -->


```py
from model_registry.types import ModelArtifact, ModelVersion, RegisteredModel

trained_model = ModelArtifact("my_model_name", "resource_URI",
                              description="Model description")

version = ModelVersion(trained_model, "v1.0", "model author")

model = RegisteredModel("my_model_name")
```

### Register objects

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

### Query objects

There are several ways to get previously registered objects from the registry.

#### By ID

IDs are created once the object is registered, you can either keep the string returned by the
`upsert_*` functions, or access the `id` property of the objects.

```py
new_model = RegisteredModel("new_model")

new_model_id = registry.upsert_registered_model(new_model)

assert new_model_id == new_model.id
```

To query objects using IDs, do

```py
another_model = registry.get_registered_model_by_id("another-model-id")

# fetching a registered_model will also fetch its associated versions
versions = another_model.versions

another_version = registry.get_model_version_by_id("another-version-id", another_model.id)

# fetching a version will also fetch its associated model artifact
model_artifact = another_version.model

another_trained_model = registry.get_model_artifact_by_id("another-model-artifact-id")
```

#### By parameters

<!-- TODO: #120 provide a link to the reference docs instead of code -->
External IDs can be used to query objects in the wild.
Note that external IDs must be unique among [artifacts](src/model_registry/types/artifacts.py) or
[contexts](src/model_registry/types/contexts.py).

```py
trained_model = ModelArtifact("my_model_name", "resource_URI",
                              description="Model description",
                              external_id="unique_reference")

# As a version is a context, we can have the same external_id as the above
version = ModelVersion(trained_model, "v1.0", "model author",
                       external_id="unique_reference")

# Registering this will cause an error!
# model = RegisteredModel("my_model_name",
#                         external_id="unique_reference")

model = RegisteredModel("my_model_name",
                        external_id="another_unique_reference")
```

You can also perform queries by parameters to get model artifacts:

```py
# We can get the model artifact associated to a version
another_trained_model = registry.get_model_artifact_by_params(model_version_id=another_version.id)

# Or by its unique identifier
trained_model = registry.get_model_artifact_by_params(external_id="unique_reference")
```

### Query multiple objects

We can query all objects of a type

```py
models = registry.get_registered_models()

versions = registry.get_model_versions("registered_model_id")

# We can get associated model artifacts with the versions
model_artifacts = [version.model for version in versions]
```

To limit or order the query, provide a [`ListOptions`](src/model_registry/types/options.py) object

```py
from model_registry import ListOptions, OrderByField

options = ListOptions(limit=50)

first_50_models = registry.get_registered_models(options)

# By default we get ascending order
options = ListOptions(order_by=OrderByField.CREATE_TIME, is_asc=False)

last_50_models = registry.get_registered_models(options)
```

## Development

Common tasks, such as building documentation and running tests, can be executed using [`nox`](https://github.com/wntrblm/nox) sessions.

Use `nox -l` to list sessions and execute them using `nox -s [session]`.

<!-- github-only -->
