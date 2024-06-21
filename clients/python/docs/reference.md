# API Reference

## Client

```{eval-rst}
.. automodule:: model_registry
```

## Core

### Register objects

<!-- TODO: update section -->
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

another_version = registry.get_model_version_by_id("another-version-id", another_model.id)

another_trained_model = registry.get_model_artifact_by_id("another-model-artifact-id")
```

#### By parameters

External IDs can be used to query objects in the wild.
Note that external IDs must be unique among artifacts or contexts (see [Types](#types)).

```py
trained_model = ModelArtifact(
    name="my_model_name",
    uri="resource_URI",
    description="Model description",
    external_id="unique_reference"
)

# As a version is a context, we can have the same external_id as the above
version = ModelVersion(
    name="v1.0",
    author="model author",
    external_id="unique_reference"
)

# Registering this will cause an error!
# model = RegisteredModel(
#    name="my_model_name",
#    external_id="unique_reference",  # this wouldn't be unique
# )

model = RegisteredModel(
    name="my_model_name",
    external_id="another_unique_reference"
)
```

You can also perform queries by parameters:

```py
# We can get the model artifact associated to a version
another_trained_model = registry.get_model_artifact_by_params(model_version_id=another_version.id)

# Or by its unique identifier
trained_model = registry.get_model_artifact_by_params(external_id="unique_reference")

# Same thing for a version
version = registry.get_model_version_by_params(external_id="unique_reference")

# Or for a model
model = registry.get_registered_model_by_params(external_id="another_unique_reference")

# We can also get a version by its name and associated model id
version = registry.get_model_version_by_params(version="v1.0", registered_model_id="x")

# And we can get a model by simply calling its name
model = registry.get_registered_model_by_params(name="my_model_name")
```

### Query multiple objects

We can query all objects of a type

```py
models = registry.get_registered_models()

versions = registry.get_model_versions("registered_model_id")

# We can get a list of all model artifacts
all_model_artifacts = registry.get_model_artifacts()
```

To limit or order the query, provide a {py:class}`model_registry.types.ListOptions` object.

```py
from model_registry import ListOptions

options = ListOptions(limit=50)

first_50_models = registry.get_registered_models(options)

# By default we get ascending order
options = ListOptions.order_by_creation_time(is_asc=False)

last_50_models = registry.get_registered_models(options)
```

```{eval-rst}
.. automodule:: model_registry.core
```

## Types

### Create objects

Registry objects can be created by doing

```py
from model_registry.types import ModelArtifact, ModelVersion, RegisteredModel

trained_model = ModelArtifact(
    name="model-exec",
    uri="resource_URI",
    description="Model description",
)

version = ModelVersion(
    name="v1.0",
    author="model author",
)

model = RegisteredModel(
    name="model",
    owner="team",
)
```

```{eval-rst}
.. automodule:: model_registry.types
```

## Exceptions

```{eval-rst}
.. automodule:: model_registry.exceptions
```
