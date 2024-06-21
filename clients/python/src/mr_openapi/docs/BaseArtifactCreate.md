# BaseArtifactCreate


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**uri** | **str** | The uniform resource identifier of the physical artifact. May be empty if there is no physical artifact. | [optional] 
**state** | [**ArtifactState**](ArtifactState.md) |  | [optional] 
**name** | **str** | The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set. | [optional] 

## Example

```python
from mr_openapi.models.base_artifact_create import BaseArtifactCreate

# TODO update the JSON string below
json = "{}"
# create an instance of BaseArtifactCreate from a JSON string
base_artifact_create_instance = BaseArtifactCreate.from_json(json)
# print the JSON string representation of the object
print(BaseArtifactCreate.to_json())

# convert the object into a dict
base_artifact_create_dict = base_artifact_create_instance.to_dict()
# create an instance of BaseArtifactCreate from a dict
base_artifact_create_from_dict = BaseArtifactCreate.from_dict(base_artifact_create_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


