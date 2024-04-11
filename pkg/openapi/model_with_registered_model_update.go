/*
Model Registry REST API

REST API for Model Registry to create and manage ML model metadata

API version: v1alpha3
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the WithRegisteredModelUpdate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WithRegisteredModelUpdate{}

// WithRegisteredModelUpdate A registered model in model registry. A registered model has ModelVersion children.
type WithRegisteredModelUpdate struct {
	State *RegisteredModelState `json:"state,omitempty"`
}

// NewWithRegisteredModelUpdate instantiates a new WithRegisteredModelUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWithRegisteredModelUpdate() *WithRegisteredModelUpdate {
	this := WithRegisteredModelUpdate{}
	var state RegisteredModelState = REGISTEREDMODELSTATE_LIVE
	this.State = &state
	return &this
}

// NewWithRegisteredModelUpdateWithDefaults instantiates a new WithRegisteredModelUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWithRegisteredModelUpdateWithDefaults() *WithRegisteredModelUpdate {
	this := WithRegisteredModelUpdate{}
	var state RegisteredModelState = REGISTEREDMODELSTATE_LIVE
	this.State = &state
	return &this
}

// GetState returns the State field value if set, zero value otherwise.
func (o *WithRegisteredModelUpdate) GetState() RegisteredModelState {
	if o == nil || IsNil(o.State) {
		var ret RegisteredModelState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WithRegisteredModelUpdate) GetStateOk() (*RegisteredModelState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *WithRegisteredModelUpdate) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given RegisteredModelState and assigns it to the State field.
func (o *WithRegisteredModelUpdate) SetState(v RegisteredModelState) {
	o.State = &v
}

func (o WithRegisteredModelUpdate) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WithRegisteredModelUpdate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	return toSerialize, nil
}

type NullableWithRegisteredModelUpdate struct {
	value *WithRegisteredModelUpdate
	isSet bool
}

func (v NullableWithRegisteredModelUpdate) Get() *WithRegisteredModelUpdate {
	return v.value
}

func (v *NullableWithRegisteredModelUpdate) Set(val *WithRegisteredModelUpdate) {
	v.value = val
	v.isSet = true
}

func (v NullableWithRegisteredModelUpdate) IsSet() bool {
	return v.isSet
}

func (v *NullableWithRegisteredModelUpdate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWithRegisteredModelUpdate(val *WithRegisteredModelUpdate) *NullableWithRegisteredModelUpdate {
	return &NullableWithRegisteredModelUpdate{value: val, isSet: true}
}

func (v NullableWithRegisteredModelUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWithRegisteredModelUpdate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}