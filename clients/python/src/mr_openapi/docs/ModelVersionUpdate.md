# ModelVersionUpdate

Represents a ModelVersion belonging to a RegisteredModel.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**state** | [**ModelVersionState**](ModelVersionState.md) |  | [optional] 
**author** | **str** | Name of the author. | [optional] 

## Example

```python
from mr_openapi.models.model_version_update import ModelVersionUpdate

# TODO update the JSON string below
json = "{}"
# create an instance of ModelVersionUpdate from a JSON string
model_version_update_instance = ModelVersionUpdate.from_json(json)
# print the JSON string representation of the object
print(ModelVersionUpdate.to_json())

# convert the object into a dict
model_version_update_dict = model_version_update_instance.to_dict()
# create an instance of ModelVersionUpdate from a dict
model_version_update_from_dict = ModelVersionUpdate.from_dict(model_version_update_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


