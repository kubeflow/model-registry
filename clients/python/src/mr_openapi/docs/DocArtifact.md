# DocArtifact

A document.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**artifact_type** | **str** |  | [default to 'doc-artifact']
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**uri** | **str** | The uniform resource identifier of the physical artifact. May be empty if there is no physical artifact. | [optional] 
**state** | [**ArtifactState**](ArtifactState.md) |  | [optional] 
**name** | **str** | The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set. | [optional] 
**id** | **str** | Output only. The unique server generated id of the resource. | [optional] [readonly] 
**create_time_since_epoch** | **str** | Output only. Create time of the resource in millisecond since epoch. | [optional] [readonly] 
**last_update_time_since_epoch** | **str** | Output only. Last update time of the resource since epoch in millisecond since epoch. | [optional] [readonly] 

## Example

```python
from mr_openapi.models.doc_artifact import DocArtifact

# TODO update the JSON string below
json = "{}"
# create an instance of DocArtifact from a JSON string
doc_artifact_instance = DocArtifact.from_json(json)
# print the JSON string representation of the object
print(DocArtifact.to_json())

# convert the object into a dict
doc_artifact_dict = doc_artifact_instance.to_dict()
# create an instance of DocArtifact from a dict
doc_artifact_from_dict = DocArtifact.from_dict(doc_artifact_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


