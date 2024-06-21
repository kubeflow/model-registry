# ModelVersionList

List of ModelVersion entities.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**next_page_token** | **str** | Token to use to retrieve next page of results. | 
**page_size** | **int** | Maximum number of resources to return in the result. | 
**size** | **int** | Number of items in result list. | 
**items** | [**List[ModelVersion]**](ModelVersion.md) | Array of &#x60;ModelVersion&#x60; entities. | [optional] 

## Example

```python
from mr_openapi.models.model_version_list import ModelVersionList

# TODO update the JSON string below
json = "{}"
# create an instance of ModelVersionList from a JSON string
model_version_list_instance = ModelVersionList.from_json(json)
# print the JSON string representation of the object
print(ModelVersionList.to_json())

# convert the object into a dict
model_version_list_dict = model_version_list_instance.to_dict()
# create an instance of ModelVersionList from a dict
model_version_list_from_dict = ModelVersionList.from_dict(model_version_list_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


