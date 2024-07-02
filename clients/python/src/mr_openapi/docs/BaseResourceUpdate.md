# BaseResourceUpdate


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 

## Example

```python
from mr_openapi.models.base_resource_update import BaseResourceUpdate

# TODO update the JSON string below
json = "{}"
# create an instance of BaseResourceUpdate from a JSON string
base_resource_update_instance = BaseResourceUpdate.from_json(json)
# print the JSON string representation of the object
print(BaseResourceUpdate.to_json())

# convert the object into a dict
base_resource_update_dict = base_resource_update_instance.to_dict()
# create an instance of BaseResourceUpdate from a dict
base_resource_update_from_dict = BaseResourceUpdate.from_dict(base_resource_update_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


