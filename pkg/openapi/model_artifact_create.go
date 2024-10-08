/*
Model Registry REST API

REST API for Model Registry to create and manage ML model metadata

API version: v1alpha3
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"fmt"
)

// ArtifactCreate - An Artifact to be created.
type ArtifactCreate struct {
	DocArtifactCreate   *DocArtifactCreate
	ModelArtifactCreate *ModelArtifactCreate
}

// DocArtifactCreateAsArtifactCreate is a convenience function that returns DocArtifactCreate wrapped in ArtifactCreate
func DocArtifactCreateAsArtifactCreate(v *DocArtifactCreate) ArtifactCreate {
	return ArtifactCreate{
		DocArtifactCreate: v,
	}
}

// ModelArtifactCreateAsArtifactCreate is a convenience function that returns ModelArtifactCreate wrapped in ArtifactCreate
func ModelArtifactCreateAsArtifactCreate(v *ModelArtifactCreate) ArtifactCreate {
	return ArtifactCreate{
		ModelArtifactCreate: v,
	}
}

// Unmarshal JSON data into one of the pointers in the struct
func (dst *ArtifactCreate) UnmarshalJSON(data []byte) error {
	var err error
	// use discriminator value to speed up the lookup
	var jsonDict map[string]interface{}
	err = newStrictDecoder(data).Decode(&jsonDict)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON into map for the discriminator lookup")
	}

	// check if the discriminator value is 'DocArtifactCreate'
	if jsonDict["artifactType"] == "DocArtifactCreate" {
		// try to unmarshal JSON data into DocArtifactCreate
		err = json.Unmarshal(data, &dst.DocArtifactCreate)
		if err == nil {
			return nil // data stored in dst.DocArtifactCreate, return on the first match
		} else {
			dst.DocArtifactCreate = nil
			return fmt.Errorf("failed to unmarshal ArtifactCreate as DocArtifactCreate: %s", err.Error())
		}
	}

	// check if the discriminator value is 'ModelArtifactCreate'
	if jsonDict["artifactType"] == "ModelArtifactCreate" {
		// try to unmarshal JSON data into ModelArtifactCreate
		err = json.Unmarshal(data, &dst.ModelArtifactCreate)
		if err == nil {
			return nil // data stored in dst.ModelArtifactCreate, return on the first match
		} else {
			dst.ModelArtifactCreate = nil
			return fmt.Errorf("failed to unmarshal ArtifactCreate as ModelArtifactCreate: %s", err.Error())
		}
	}

	// check if the discriminator value is 'doc-artifact'
	if jsonDict["artifactType"] == "doc-artifact" {
		// try to unmarshal JSON data into DocArtifactCreate
		err = json.Unmarshal(data, &dst.DocArtifactCreate)
		if err == nil {
			return nil // data stored in dst.DocArtifactCreate, return on the first match
		} else {
			dst.DocArtifactCreate = nil
			return fmt.Errorf("failed to unmarshal ArtifactCreate as DocArtifactCreate: %s", err.Error())
		}
	}

	// check if the discriminator value is 'model-artifact'
	if jsonDict["artifactType"] == "model-artifact" {
		// try to unmarshal JSON data into ModelArtifactCreate
		err = json.Unmarshal(data, &dst.ModelArtifactCreate)
		if err == nil {
			return nil // data stored in dst.ModelArtifactCreate, return on the first match
		} else {
			dst.ModelArtifactCreate = nil
			return fmt.Errorf("failed to unmarshal ArtifactCreate as ModelArtifactCreate: %s", err.Error())
		}
	}

	return nil
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src ArtifactCreate) MarshalJSON() ([]byte, error) {
	if src.DocArtifactCreate != nil {
		return json.Marshal(&src.DocArtifactCreate)
	}

	if src.ModelArtifactCreate != nil {
		return json.Marshal(&src.ModelArtifactCreate)
	}

	return nil, nil // no data in oneOf schemas
}

// Get the actual instance
func (obj *ArtifactCreate) GetActualInstance() interface{} {
	if obj == nil {
		return nil
	}
	if obj.DocArtifactCreate != nil {
		return obj.DocArtifactCreate
	}

	if obj.ModelArtifactCreate != nil {
		return obj.ModelArtifactCreate
	}

	// all schemas are nil
	return nil
}

type NullableArtifactCreate struct {
	value *ArtifactCreate
	isSet bool
}

func (v NullableArtifactCreate) Get() *ArtifactCreate {
	return v.value
}

func (v *NullableArtifactCreate) Set(val *ArtifactCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableArtifactCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableArtifactCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableArtifactCreate(val *ArtifactCreate) *NullableArtifactCreate {
	return &NullableArtifactCreate{value: val, isSet: true}
}

func (v NullableArtifactCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableArtifactCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
