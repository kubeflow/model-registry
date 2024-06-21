# ServingEnvironmentList

List of ServingEnvironments.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**next_page_token** | **str** | Token to use to retrieve next page of results. | 
**page_size** | **int** | Maximum number of resources to return in the result. | 
**size** | **int** | Number of items in result list. | 
**items** | [**List[ServingEnvironment]**](ServingEnvironment.md) |  | [optional] 

## Example

```python
from mr_openapi.models.serving_environment_list import ServingEnvironmentList

# TODO update the JSON string below
json = "{}"
# create an instance of ServingEnvironmentList from a JSON string
serving_environment_list_instance = ServingEnvironmentList.from_json(json)
# print the JSON string representation of the object
print(ServingEnvironmentList.to_json())

# convert the object into a dict
serving_environment_list_dict = serving_environment_list_instance.to_dict()
# create an instance of ServingEnvironmentList from a dict
serving_environment_list_from_dict = ServingEnvironmentList.from_dict(serving_environment_list_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


