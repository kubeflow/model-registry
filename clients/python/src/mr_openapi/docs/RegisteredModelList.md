# RegisteredModelList

List of RegisteredModels.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**next_page_token** | **str** | Token to use to retrieve next page of results. | 
**page_size** | **int** | Maximum number of resources to return in the result. | 
**size** | **int** | Number of items in result list. | 
**items** | [**List[RegisteredModel]**](RegisteredModel.md) |  | [optional] 

## Example

```python
from mr_openapi.models.registered_model_list import RegisteredModelList

# TODO update the JSON string below
json = "{}"
# create an instance of RegisteredModelList from a JSON string
registered_model_list_instance = RegisteredModelList.from_json(json)
# print the JSON string representation of the object
print(RegisteredModelList.to_json())

# convert the object into a dict
registered_model_list_dict = registered_model_list_instance.to_dict()
# create an instance of RegisteredModelList from a dict
registered_model_list_from_dict = RegisteredModelList.from_dict(registered_model_list_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


