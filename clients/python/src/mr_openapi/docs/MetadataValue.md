# MetadataValue

A value in properties.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**int_value** | **str** |  | 
**metadata_type** | **str** |  | [default to 'MetadataBoolValue']
**double_value** | **float** |  | 
**string_value** | **str** |  | 
**struct_value** | **str** | Base64 encoded bytes for struct value | 
**type** | **str** | url describing proto value | 
**proto_value** | **str** | Base64 encoded bytes for proto value | 
**bool_value** | **bool** |  | 

## Example

```python
from mr_openapi.models.metadata_value import MetadataValue

# TODO update the JSON string below
json = "{}"
# create an instance of MetadataValue from a JSON string
metadata_value_instance = MetadataValue.from_json(json)
# print the JSON string representation of the object
print(MetadataValue.to_json())

# convert the object into a dict
metadata_value_dict = metadata_value_instance.to_dict()
# create an instance of MetadataValue from a dict
metadata_value_from_dict = MetadataValue.from_dict(metadata_value_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


