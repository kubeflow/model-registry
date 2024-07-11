# InferenceService

An `InferenceService` entity in a `ServingEnvironment` represents a deployed `ModelVersion` from a `RegisteredModel` created by Model Serving.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**name** | **str** | The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set. | [optional] 
**id** | **str** | Output only. The unique server generated id of the resource. | [optional] [readonly] 
**create_time_since_epoch** | **str** | Output only. Create time of the resource in millisecond since epoch. | [optional] [readonly] 
**last_update_time_since_epoch** | **str** | Output only. Last update time of the resource since epoch in millisecond since epoch. | [optional] [readonly] 
**model_version_id** | **str** | ID of the &#x60;ModelVersion&#x60; to serve. If it&#39;s unspecified, then the latest &#x60;ModelVersion&#x60; by creation order will be served. | [optional] 
**runtime** | **str** | Model runtime. | [optional] 
**desired_state** | [**InferenceServiceState**](InferenceServiceState.md) |  | [optional] 
**registered_model_id** | **str** | ID of the &#x60;RegisteredModel&#x60; to serve. | 
**serving_environment_id** | **str** | ID of the parent &#x60;ServingEnvironment&#x60; for this &#x60;InferenceService&#x60; entity. | 

## Example

```python
from mr_openapi.models.inference_service import InferenceService

# TODO update the JSON string below
json = "{}"
# create an instance of InferenceService from a JSON string
inference_service_instance = InferenceService.from_json(json)
# print the JSON string representation of the object
print(InferenceService.to_json())

# convert the object into a dict
inference_service_dict = inference_service_instance.to_dict()
# create an instance of InferenceService from a dict
inference_service_from_dict = InferenceService.from_dict(inference_service_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


