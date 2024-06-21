# MetadataIntValue

An integer (int64) property value.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**int_value** | **str** |  | 
**metadata_type** | **str** |  | [default to 'MetadataIntValue']

## Example

```python
from mr_openapi.models.metadata_int_value import MetadataIntValue

# TODO update the JSON string below
json = "{}"
# create an instance of MetadataIntValue from a JSON string
metadata_int_value_instance = MetadataIntValue.from_json(json)
# print the JSON string representation of the object
print(MetadataIntValue.to_json())

# convert the object into a dict
metadata_int_value_dict = metadata_int_value_instance.to_dict()
# create an instance of MetadataIntValue from a dict
metadata_int_value_from_dict = MetadataIntValue.from_dict(metadata_int_value_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


