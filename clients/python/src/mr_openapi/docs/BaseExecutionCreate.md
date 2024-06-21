# BaseExecutionCreate


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**last_known_state** | [**ExecutionState**](ExecutionState.md) |  | [optional] 
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**name** | **str** | The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set. | [optional] 

## Example

```python
from mr_openapi.models.base_execution_create import BaseExecutionCreate

# TODO update the JSON string below
json = "{}"
# create an instance of BaseExecutionCreate from a JSON string
base_execution_create_instance = BaseExecutionCreate.from_json(json)
# print the JSON string representation of the object
print(BaseExecutionCreate.to_json())

# convert the object into a dict
base_execution_create_dict = base_execution_create_instance.to_dict()
# create an instance of BaseExecutionCreate from a dict
base_execution_create_from_dict = BaseExecutionCreate.from_dict(base_execution_create_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


