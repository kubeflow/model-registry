# ServeModelList

List of ServeModel entities.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**next_page_token** | **str** | Token to use to retrieve next page of results. | 
**page_size** | **int** | Maximum number of resources to return in the result. | 
**size** | **int** | Number of items in result list. | 
**items** | [**List[ServeModel]**](ServeModel.md) | Array of &#x60;ModelArtifact&#x60; entities. | [optional] 

## Example

```python
from mr_openapi.models.serve_model_list import ServeModelList

# TODO update the JSON string below
json = "{}"
# create an instance of ServeModelList from a JSON string
serve_model_list_instance = ServeModelList.from_json(json)
# print the JSON string representation of the object
print(ServeModelList.to_json())

# convert the object into a dict
serve_model_list_dict = serve_model_list_instance.to_dict()
# create an instance of ServeModelList from a dict
serve_model_list_from_dict = ServeModelList.from_dict(serve_model_list_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


