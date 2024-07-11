# ServeModel

An ML model serving action.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**last_known_state** | [**ExecutionState**](ExecutionState.md) |  | [optional] 
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**name** | **str** | The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set. | [optional] 
**id** | **str** | Output only. The unique server generated id of the resource. | [optional] [readonly] 
**create_time_since_epoch** | **str** | Output only. Create time of the resource in millisecond since epoch. | [optional] [readonly] 
**last_update_time_since_epoch** | **str** | Output only. Last update time of the resource since epoch in millisecond since epoch. | [optional] [readonly] 
**model_version_id** | **str** | ID of the &#x60;ModelVersion&#x60; that was served in &#x60;InferenceService&#x60;. | 

## Example

```python
from mr_openapi.models.serve_model import ServeModel

# TODO update the JSON string below
json = "{}"
# create an instance of ServeModel from a JSON string
serve_model_instance = ServeModel.from_json(json)
# print the JSON string representation of the object
print(ServeModel.to_json())

# convert the object into a dict
serve_model_dict = serve_model_instance.to_dict()
# create an instance of ServeModel from a dict
serve_model_from_dict = ServeModel.from_dict(serve_model_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

