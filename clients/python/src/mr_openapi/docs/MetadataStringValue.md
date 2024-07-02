# MetadataStringValue

A string property value.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**string_value** | **str** |  | 
**metadata_type** | **str** |  | [default to 'MetadataStringValue']

## Example

```python
from mr_openapi.models.metadata_string_value import MetadataStringValue

# TODO update the JSON string below
json = "{}"
# create an instance of MetadataStringValue from a JSON string
metadata_string_value_instance = MetadataStringValue.from_json(json)
# print the JSON string representation of the object
print(MetadataStringValue.to_json())

# convert the object into a dict
metadata_string_value_dict = metadata_string_value_instance.to_dict()
# create an instance of MetadataStringValue from a dict
metadata_string_value_from_dict = MetadataStringValue.from_dict(metadata_string_value_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


