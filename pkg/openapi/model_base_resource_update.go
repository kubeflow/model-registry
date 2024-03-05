/*
Model Registry REST API

REST API for Model Registry to create and manage ML model metadata

API version: v1alpha2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the BaseResourceUpdate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BaseResourceUpdate{}

// BaseResourceUpdate struct for BaseResourceUpdate
type BaseResourceUpdate struct {
	// User provided custom properties which are not defined by its type.
	CustomProperties *map[string]MetadataValue `json:"customProperties,omitempty"`
	// An optional description about the resource.
	Description *string `json:"description,omitempty"`
	// The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance.
	ExternalId *string `json:"externalId,omitempty"`
}

// NewBaseResourceUpdate instantiates a new BaseResourceUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBaseResourceUpdate() *BaseResourceUpdate {
	this := BaseResourceUpdate{}
	return &this
}

// NewBaseResourceUpdateWithDefaults instantiates a new BaseResourceUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBaseResourceUpdateWithDefaults() *BaseResourceUpdate {
	this := BaseResourceUpdate{}
	return &this
}

// GetCustomProperties returns the CustomProperties field value if set, zero value otherwise.
func (o *BaseResourceUpdate) GetCustomProperties() map[string]MetadataValue {
	if o == nil || IsNil(o.CustomProperties) {
		var ret map[string]MetadataValue
		return ret
	}
	return *o.CustomProperties
}

// GetCustomPropertiesOk returns a tuple with the CustomProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseResourceUpdate) GetCustomPropertiesOk() (*map[string]MetadataValue, bool) {
	if o == nil || IsNil(o.CustomProperties) {
		return nil, false
	}
	return o.CustomProperties, true
}

// HasCustomProperties returns a boolean if a field has been set.
func (o *BaseResourceUpdate) HasCustomProperties() bool {
	if o != nil && !IsNil(o.CustomProperties) {
		return true
	}

	return false
}

// SetCustomProperties gets a reference to the given map[string]MetadataValue and assigns it to the CustomProperties field.
func (o *BaseResourceUpdate) SetCustomProperties(v map[string]MetadataValue) {
	o.CustomProperties = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *BaseResourceUpdate) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseResourceUpdate) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *BaseResourceUpdate) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *BaseResourceUpdate) SetDescription(v string) {
	o.Description = &v
}

// GetExternalId returns the ExternalId field value if set, zero value otherwise.
func (o *BaseResourceUpdate) GetExternalId() string {
	if o == nil || IsNil(o.ExternalId) {
		var ret string
		return ret
	}
	return *o.ExternalId
}

// GetExternalIdOk returns a tuple with the ExternalId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseResourceUpdate) GetExternalIdOk() (*string, bool) {
	if o == nil || IsNil(o.ExternalId) {
		return nil, false
	}
	return o.ExternalId, true
}

// HasExternalId returns a boolean if a field has been set.
func (o *BaseResourceUpdate) HasExternalId() bool {
	if o != nil && !IsNil(o.ExternalId) {
		return true
	}

	return false
}

// SetExternalId gets a reference to the given string and assigns it to the ExternalId field.
func (o *BaseResourceUpdate) SetExternalId(v string) {
	o.ExternalId = &v
}

func (o BaseResourceUpdate) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BaseResourceUpdate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.CustomProperties) {
		toSerialize["customProperties"] = o.CustomProperties
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.ExternalId) {
		toSerialize["externalId"] = o.ExternalId
	}
	return toSerialize, nil
}

type NullableBaseResourceUpdate struct {
	value *BaseResourceUpdate
	isSet bool
}

func (v NullableBaseResourceUpdate) Get() *BaseResourceUpdate {
	return v.value
}

func (v *NullableBaseResourceUpdate) Set(val *BaseResourceUpdate) {
	v.value = val
	v.isSet = true
}

func (v NullableBaseResourceUpdate) IsSet() bool {
	return v.isSet
}

func (v *NullableBaseResourceUpdate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBaseResourceUpdate(val *BaseResourceUpdate) *NullableBaseResourceUpdate {
	return &NullableBaseResourceUpdate{value: val, isSet: true}
}

func (v NullableBaseResourceUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBaseResourceUpdate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
