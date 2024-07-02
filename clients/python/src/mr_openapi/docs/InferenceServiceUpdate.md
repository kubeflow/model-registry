# InferenceServiceUpdate

An `InferenceService` entity in a `ServingEnvironment` represents a deployed `ModelVersion` from a `RegisteredModel` created by Model Serving.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**model_version_id** | **str** | ID of the &#x60;ModelVersion&#x60; to serve. If it&#39;s unspecified, then the latest &#x60;ModelVersion&#x60; by creation order will be served. | [optional] 
**runtime** | **str** | Model runtime. | [optional] 
**desired_state** | [**InferenceServiceState**](InferenceServiceState.md) |  | [optional] 

## Example

```python
from mr_openapi.models.inference_service_update import InferenceServiceUpdate

# TODO update the JSON string below
json = "{}"
# create an instance of InferenceServiceUpdate from a JSON string
inference_service_update_instance = InferenceServiceUpdate.from_json(json)
# print the JSON string representation of the object
print(InferenceServiceUpdate.to_json())

# convert the object into a dict
inference_service_update_dict = inference_service_update_instance.to_dict()
# create an instance of InferenceServiceUpdate from a dict
inference_service_update_from_dict = InferenceServiceUpdate.from_dict(inference_service_update_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


