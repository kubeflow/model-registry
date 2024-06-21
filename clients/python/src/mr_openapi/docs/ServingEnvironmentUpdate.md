# ServingEnvironmentUpdate

A Model Serving environment for serving `RegisteredModels`.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 

## Example

```python
from mr_openapi.models.serving_environment_update import ServingEnvironmentUpdate

# TODO update the JSON string below
json = "{}"
# create an instance of ServingEnvironmentUpdate from a JSON string
serving_environment_update_instance = ServingEnvironmentUpdate.from_json(json)
# print the JSON string representation of the object
print(ServingEnvironmentUpdate.to_json())

# convert the object into a dict
serving_environment_update_dict = serving_environment_update_instance.to_dict()
# create an instance of ServingEnvironmentUpdate from a dict
serving_environment_update_from_dict = ServingEnvironmentUpdate.from_dict(serving_environment_update_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


