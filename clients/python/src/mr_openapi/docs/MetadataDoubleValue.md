# MetadataDoubleValue

A double property value.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**double_value** | **float** |  | 
**metadata_type** | **str** |  | [default to 'MetadataDoubleValue']

## Example

```python
from mr_openapi.models.metadata_double_value import MetadataDoubleValue

# TODO update the JSON string below
json = "{}"
# create an instance of MetadataDoubleValue from a JSON string
metadata_double_value_instance = MetadataDoubleValue.from_json(json)
# print the JSON string representation of the object
print(MetadataDoubleValue.to_json())

# convert the object into a dict
metadata_double_value_dict = metadata_double_value_instance.to_dict()
# create an instance of MetadataDoubleValue from a dict
metadata_double_value_from_dict = MetadataDoubleValue.from_dict(metadata_double_value_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


