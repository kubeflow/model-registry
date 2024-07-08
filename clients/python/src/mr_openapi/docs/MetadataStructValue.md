# MetadataStructValue

A struct property value.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**struct_value** | **str** | Base64 encoded bytes for struct value | 
**metadata_type** | **str** |  | [default to 'MetadataStructValue']

## Example

```python
from mr_openapi.models.metadata_struct_value import MetadataStructValue

# TODO update the JSON string below
json = "{}"
# create an instance of MetadataStructValue from a JSON string
metadata_struct_value_instance = MetadataStructValue.from_json(json)
# print the JSON string representation of the object
print(MetadataStructValue.to_json())

# convert the object into a dict
metadata_struct_value_dict = metadata_struct_value_instance.to_dict()
# create an instance of MetadataStructValue from a dict
metadata_struct_value_from_dict = MetadataStructValue.from_dict(metadata_struct_value_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


