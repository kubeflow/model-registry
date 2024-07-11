# MetadataBoolValue

A bool property value.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**bool_value** | **bool** |  | 
**metadata_type** | **str** |  | [default to 'MetadataBoolValue']

## Example

```python
from mr_openapi.models.metadata_bool_value import MetadataBoolValue

# TODO update the JSON string below
json = "{}"
# create an instance of MetadataBoolValue from a JSON string
metadata_bool_value_instance = MetadataBoolValue.from_json(json)
# print the JSON string representation of the object
print(MetadataBoolValue.to_json())

# convert the object into a dict
metadata_bool_value_dict = metadata_bool_value_instance.to_dict()
# create an instance of MetadataBoolValue from a dict
metadata_bool_value_from_dict = MetadataBoolValue.from_dict(metadata_bool_value_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


