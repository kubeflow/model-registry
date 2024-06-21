# ServingEnvironmentCreate

A Model Serving environment for serving `RegisteredModels`.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**name** | **str** | The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set. | [optional] 

## Example

```python
from mr_openapi.models.serving_environment_create import ServingEnvironmentCreate

# TODO update the JSON string below
json = "{}"
# create an instance of ServingEnvironmentCreate from a JSON string
serving_environment_create_instance = ServingEnvironmentCreate.from_json(json)
# print the JSON string representation of the object
print(ServingEnvironmentCreate.to_json())

# convert the object into a dict
serving_environment_create_dict = serving_environment_create_instance.to_dict()
# create an instance of ServingEnvironmentCreate from a dict
serving_environment_create_from_dict = ServingEnvironmentCreate.from_dict(serving_environment_create_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


