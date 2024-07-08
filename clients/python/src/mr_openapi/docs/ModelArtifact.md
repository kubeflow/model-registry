# ModelArtifact

An ML model artifact.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**artifact_type** | **str** |  | [default to 'model-artifact']
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**uri** | **str** | The uniform resource identifier of the physical artifact. May be empty if there is no physical artifact. | [optional] 
**state** | [**ArtifactState**](ArtifactState.md) |  | [optional] 
**name** | **str** | The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set. | [optional] 
**id** | **str** | Output only. The unique server generated id of the resource. | [optional] [readonly] 
**create_time_since_epoch** | **str** | Output only. Create time of the resource in millisecond since epoch. | [optional] [readonly] 
**last_update_time_since_epoch** | **str** | Output only. Last update time of the resource since epoch in millisecond since epoch. | [optional] [readonly] 
**model_format_name** | **str** | Name of the model format. | [optional] 
**storage_key** | **str** | Storage secret name. | [optional] 
**storage_path** | **str** | Path for model in storage provided by &#x60;storageKey&#x60;. | [optional] 
**model_format_version** | **str** | Version of the model format. | [optional] 
**service_account_name** | **str** | Name of the service account with storage secret. | [optional] 

## Example

```python
from mr_openapi.models.model_artifact import ModelArtifact

# TODO update the JSON string below
json = "{}"
# create an instance of ModelArtifact from a JSON string
model_artifact_instance = ModelArtifact.from_json(json)
# print the JSON string representation of the object
print(ModelArtifact.to_json())

# convert the object into a dict
model_artifact_dict = model_artifact_instance.to_dict()
# create an instance of ModelArtifact from a dict
model_artifact_from_dict = ModelArtifact.from_dict(model_artifact_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


