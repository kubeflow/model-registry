# ServeModelUpdate

An ML model serving action.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**last_known_state** | [**ExecutionState**](ExecutionState.md) |  | [optional] 
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 

## Example

```python
from mr_openapi.models.serve_model_update import ServeModelUpdate

# TODO update the JSON string below
json = "{}"
# create an instance of ServeModelUpdate from a JSON string
serve_model_update_instance = ServeModelUpdate.from_json(json)
# print the JSON string representation of the object
print(ServeModelUpdate.to_json())

# convert the object into a dict
serve_model_update_dict = serve_model_update_instance.to_dict()
# create an instance of ServeModelUpdate from a dict
serve_model_update_from_dict = ServeModelUpdate.from_dict(serve_model_update_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


