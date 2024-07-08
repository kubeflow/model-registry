# ModelArtifactList

List of ModelArtifact entities.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**next_page_token** | **str** | Token to use to retrieve next page of results. | 
**page_size** | **int** | Maximum number of resources to return in the result. | 
**size** | **int** | Number of items in result list. | 
**items** | [**List[ModelArtifact]**](ModelArtifact.md) | Array of &#x60;ModelArtifact&#x60; entities. | [optional] 

## Example

```python
from mr_openapi.models.model_artifact_list import ModelArtifactList

# TODO update the JSON string below
json = "{}"
# create an instance of ModelArtifactList from a JSON string
model_artifact_list_instance = ModelArtifactList.from_json(json)
# print the JSON string representation of the object
print(ModelArtifactList.to_json())

# convert the object into a dict
model_artifact_list_dict = model_artifact_list_instance.to_dict()
# create an instance of ModelArtifactList from a dict
model_artifact_list_from_dict = ModelArtifactList.from_dict(model_artifact_list_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


