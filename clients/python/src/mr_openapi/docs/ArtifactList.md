# ArtifactList

A list of Artifact entities.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**next_page_token** | **str** | Token to use to retrieve next page of results. | 
**page_size** | **int** | Maximum number of resources to return in the result. | 
**size** | **int** | Number of items in result list. | 
**items** | [**List[Artifact]**](Artifact.md) | Array of &#x60;Artifact&#x60; entities. | [optional] 

## Example

```python
from mr_openapi.models.artifact_list import ArtifactList

# TODO update the JSON string below
json = "{}"
# create an instance of ArtifactList from a JSON string
artifact_list_instance = ArtifactList.from_json(json)
# print the JSON string representation of the object
print(ArtifactList.to_json())

# convert the object into a dict
artifact_list_dict = artifact_list_instance.to_dict()
# create an instance of ArtifactList from a dict
artifact_list_from_dict = ArtifactList.from_dict(artifact_list_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


