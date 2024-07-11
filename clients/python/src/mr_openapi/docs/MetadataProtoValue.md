# MetadataProtoValue

A proto property value.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**type** | **str** | url describing proto value | 
**proto_value** | **str** | Base64 encoded bytes for proto value | 
**metadata_type** | **str** |  | [default to 'MetadataProtoValue']

## Example

```python
from mr_openapi.models.metadata_proto_value import MetadataProtoValue

# TODO update the JSON string below
json = "{}"
# create an instance of MetadataProtoValue from a JSON string
metadata_proto_value_instance = MetadataProtoValue.from_json(json)
# print the JSON string representation of the object
print(MetadataProtoValue.to_json())

# convert the object into a dict
metadata_proto_value_dict = metadata_proto_value_instance.to_dict()
# create an instance of MetadataProtoValue from a dict
metadata_proto_value_from_dict = MetadataProtoValue.from_dict(metadata_proto_value_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


