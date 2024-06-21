# InferenceServiceCreate

An `InferenceService` entity in a `ServingEnvironment` represents a deployed `ModelVersion` from a `RegisteredModel` created by Model Serving.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**name** | **str** | The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set. | [optional] 
**model_version_id** | **str** | ID of the &#x60;ModelVersion&#x60; to serve. If it&#39;s unspecified, then the latest &#x60;ModelVersion&#x60; by creation order will be served. | [optional] 
**runtime** | **str** | Model runtime. | [optional] 
**desired_state** | [**InferenceServiceState**](InferenceServiceState.md) |  | [optional] 
**registered_model_id** | **str** | ID of the &#x60;RegisteredModel&#x60; to serve. | 
**serving_environment_id** | **str** | ID of the parent &#x60;ServingEnvironment&#x60; for this &#x60;InferenceService&#x60; entity. | 

## Example

```python
from mr_openapi.models.inference_service_create import InferenceServiceCreate

# TODO update the JSON string below
json = "{}"
# create an instance of InferenceServiceCreate from a JSON string
inference_service_create_instance = InferenceServiceCreate.from_json(json)
# print the JSON string representation of the object
print(InferenceServiceCreate.to_json())

# convert the object into a dict
inference_service_create_dict = inference_service_create_instance.to_dict()
# create an instance of InferenceServiceCreate from a dict
inference_service_create_from_dict = InferenceServiceCreate.from_dict(inference_service_create_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


