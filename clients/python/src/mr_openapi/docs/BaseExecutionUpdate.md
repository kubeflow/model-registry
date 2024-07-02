# BaseExecutionUpdate


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**last_known_state** | [**ExecutionState**](ExecutionState.md) |  | [optional] 

## Example

```python
from mr_openapi.models.base_execution_update import BaseExecutionUpdate

# TODO update the JSON string below
json = "{}"
# create an instance of BaseExecutionUpdate from a JSON string
base_execution_update_instance = BaseExecutionUpdate.from_json(json)
# print the JSON string representation of the object
print(BaseExecutionUpdate.to_json())

# convert the object into a dict
base_execution_update_dict = base_execution_update_instance.to_dict()
# create an instance of BaseExecutionUpdate from a dict
base_execution_update_from_dict = BaseExecutionUpdate.from_dict(base_execution_update_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


