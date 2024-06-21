# InferenceServiceList

List of InferenceServices.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**next_page_token** | **str** | Token to use to retrieve next page of results. | 
**page_size** | **int** | Maximum number of resources to return in the result. | 
**size** | **int** | Number of items in result list. | 
**items** | [**List[InferenceService]**](InferenceService.md) |  | [optional] 

## Example

```python
from mr_openapi.models.inference_service_list import InferenceServiceList

# TODO update the JSON string below
json = "{}"
# create an instance of InferenceServiceList from a JSON string
inference_service_list_instance = InferenceServiceList.from_json(json)
# print the JSON string representation of the object
print(InferenceServiceList.to_json())

# convert the object into a dict
inference_service_list_dict = inference_service_list_instance.to_dict()
# create an instance of InferenceServiceList from a dict
inference_service_list_from_dict = InferenceServiceList.from_dict(inference_service_list_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


