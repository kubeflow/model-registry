# ModelArtifactCreate

An ML model artifact.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**custom_properties** | [**Dict[str, MetadataValue]**](MetadataValue.md) | User provided custom properties which are not defined by its type. | [optional] 
**description** | **str** | An optional description about the resource. | [optional] 
**external_id** | **str** | The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance. | [optional] 
**uri** | **str** | The uniform resource identifier of the physical artifact. May be empty if there is no physical artifact. | [optional] 
**state** | [**ArtifactState**](ArtifactState.md) |  | [optional] 
**name** | **str** | The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set. | [optional] 
**model_format_name** | **str** | Name of the model format. | [optional] 
**storage_key** | **str** | Storage secret name. | [optional] 
**storage_path** | **str** | Path for model in storage provided by &#x60;storageKey&#x60;. | [optional] 
**model_format_version** | **str** | Version of the model format. | [optional] 
**service_account_name** | **str** | Name of the service account with storage secret. | [optional] 

## Example

```python
from mr_openapi.models.model_artifact_create import ModelArtifactCreate

# TODO update the JSON string below
json = "{}"
# create an instance of ModelArtifactCreate from a JSON string
model_artifact_create_instance = ModelArtifactCreate.from_json(json)
# print the JSON string representation of the object
print(ModelArtifactCreate.to_json())

# convert the object into a dict
model_artifact_create_dict = model_artifact_create_instance.to_dict()
# create an instance of ModelArtifactCreate from a dict
model_artifact_create_from_dict = ModelArtifactCreate.from_dict(model_artifact_create_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


