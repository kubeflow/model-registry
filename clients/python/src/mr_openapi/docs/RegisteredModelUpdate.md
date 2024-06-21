# RegisteredModelUpdate

A registered model in model registry. A registered model has ModelVersion children.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**owner** | **str** |  | [optional] 
**state** | [**RegisteredModelState**](RegisteredModelState.md) |  | [optional] 

## Example

```python
from mr_openapi.models.registered_model_update import RegisteredModelUpdate

# TODO update the JSON string below
json = "{}"
# create an instance of RegisteredModelUpdate from a JSON string
registered_model_update_instance = RegisteredModelUpdate.from_json(json)
# print the JSON string representation of the object
print(RegisteredModelUpdate.to_json())

# convert the object into a dict
registered_model_update_dict = registered_model_update_instance.to_dict()
# create an instance of RegisteredModelUpdate from a dict
registered_model_update_from_dict = RegisteredModelUpdate.from_dict(registered_model_update_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


