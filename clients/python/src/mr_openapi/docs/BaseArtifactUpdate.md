# BaseArtifactUpdate


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**uri** | **str** | The uniform resource identifier of the physical artifact. May be empty if there is no physical artifact. | [optional] 
**state** | [**ArtifactState**](ArtifactState.md) |  | [optional] 

## Example

```python
from mr_openapi.models.base_artifact_update import BaseArtifactUpdate

# TODO update the JSON string below
json = "{}"
# create an instance of BaseArtifactUpdate from a JSON string
base_artifact_update_instance = BaseArtifactUpdate.from_json(json)
# print the JSON string representation of the object
print(BaseArtifactUpdate.to_json())

# convert the object into a dict
base_artifact_update_dict = base_artifact_update_instance.to_dict()
# create an instance of BaseArtifactUpdate from a dict
base_artifact_update_from_dict = BaseArtifactUpdate.from_dict(base_artifact_update_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


