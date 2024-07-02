# BaseResourceList


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**next_page_token** | **str** | Token to use to retrieve next page of results. | 
**page_size** | **int** | Maximum number of resources to return in the result. | 
**size** | **int** | Number of items in result list. | 

## Example

```python
from mr_openapi.models.base_resource_list import BaseResourceList

# TODO update the JSON string below
json = "{}"
# create an instance of BaseResourceList from a JSON string
base_resource_list_instance = BaseResourceList.from_json(json)
# print the JSON string representation of the object
print(BaseResourceList.to_json())

# convert the object into a dict
base_resource_list_dict = base_resource_list_instance.to_dict()
# create an instance of BaseResourceList from a dict
base_resource_list_from_dict = BaseResourceList.from_dict(base_resource_list_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


