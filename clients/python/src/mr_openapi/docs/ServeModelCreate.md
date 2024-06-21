# ServeModelCreate

An ML model serving action.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**last_known_state** | [**ExecutionState**](ExecutionState.md) |  | [optional] 
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**name** | **str** | The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set. | [optional] 
**model_version_id** | **str** | ID of the &#x60;ModelVersion&#x60; that was served in &#x60;InferenceService&#x60;. | 

## Example

```python
from mr_openapi.models.serve_model_create import ServeModelCreate

# TODO update the JSON string below
json = "{}"
# create an instance of ServeModelCreate from a JSON string
serve_model_create_instance = ServeModelCreate.from_json(json)
# print the JSON string representation of the object
print(ServeModelCreate.to_json())

# convert the object into a dict
serve_model_create_dict = serve_model_create_instance.to_dict()
# create an instance of ServeModelCreate from a dict
serve_model_create_from_dict = ServeModelCreate.from_dict(serve_model_create_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


