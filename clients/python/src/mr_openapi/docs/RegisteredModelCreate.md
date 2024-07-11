# RegisteredModelCreate

A registered model in model registry. A registered model has ModelVersion children.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**name** | **str** | The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set. | [optional] 
**owner** | **str** |  | [optional] 
**state** | [**RegisteredModelState**](RegisteredModelState.md) |  | [optional] 

## Example

```python
from mr_openapi.models.registered_model_create import RegisteredModelCreate

# TODO update the JSON string below
json = "{}"
# create an instance of RegisteredModelCreate from a JSON string
registered_model_create_instance = RegisteredModelCreate.from_json(json)
# print the JSON string representation of the object
print(RegisteredModelCreate.to_json())

# convert the object into a dict
registered_model_create_dict = registered_model_create_instance.to_dict()
# create an instance of RegisteredModelCreate from a dict
registered_model_create_from_dict = RegisteredModelCreate.from_dict(registered_model_create_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


