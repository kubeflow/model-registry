# ModelVersionCreate

Represents a ModelVersion belonging to a RegisteredModel.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**name** | **str** | The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set. | [optional] 
**state** | [**ModelVersionState**](ModelVersionState.md) |  | [optional] 
**author** | **str** | Name of the author. | [optional] 
**registered_model_id** | **str** | ID of the &#x60;RegisteredModel&#x60; to which this version belongs. | 

## Example

```python
from mr_openapi.models.model_version_create import ModelVersionCreate

# TODO update the JSON string below
json = "{}"
# create an instance of ModelVersionCreate from a JSON string
model_version_create_instance = ModelVersionCreate.from_json(json)
# print the JSON string representation of the object
print(ModelVersionCreate.to_json())

# convert the object into a dict
model_version_create_dict = model_version_create_instance.to_dict()
# create an instance of ModelVersionCreate from a dict
model_version_create_from_dict = ModelVersionCreate.from_dict(model_version_create_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


