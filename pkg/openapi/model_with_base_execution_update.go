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

// checks if the WithBaseExecutionUpdate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WithBaseExecutionUpdate{}

// WithBaseExecutionUpdate struct for WithBaseExecutionUpdate
type WithBaseExecutionUpdate struct {
	LastKnownState *ExecutionState `json:"lastKnownState,omitempty"`
}

// NewWithBaseExecutionUpdate instantiates a new WithBaseExecutionUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWithBaseExecutionUpdate() *WithBaseExecutionUpdate {
	this := WithBaseExecutionUpdate{}
	var lastKnownState ExecutionState = EXECUTIONSTATE_UNKNOWN
	this.LastKnownState = &lastKnownState
	return &this
}

// NewWithBaseExecutionUpdateWithDefaults instantiates a new WithBaseExecutionUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWithBaseExecutionUpdateWithDefaults() *WithBaseExecutionUpdate {
	this := WithBaseExecutionUpdate{}
	var lastKnownState ExecutionState = EXECUTIONSTATE_UNKNOWN
	this.LastKnownState = &lastKnownState
	return &this
}

// GetLastKnownState returns the LastKnownState field value if set, zero value otherwise.
func (o *WithBaseExecutionUpdate) GetLastKnownState() ExecutionState {
	if o == nil || IsNil(o.LastKnownState) {
		var ret ExecutionState
		return ret
	}
	return *o.LastKnownState
}

// GetLastKnownStateOk returns a tuple with the LastKnownState field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WithBaseExecutionUpdate) GetLastKnownStateOk() (*ExecutionState, bool) {
	if o == nil || IsNil(o.LastKnownState) {
		return nil, false
	}
	return o.LastKnownState, true
}

// HasLastKnownState returns a boolean if a field has been set.
func (o *WithBaseExecutionUpdate) HasLastKnownState() bool {
	if o != nil && !IsNil(o.LastKnownState) {
		return true
	}

	return false
}

// SetLastKnownState gets a reference to the given ExecutionState and assigns it to the LastKnownState field.
func (o *WithBaseExecutionUpdate) SetLastKnownState(v ExecutionState) {
	o.LastKnownState = &v
}

func (o WithBaseExecutionUpdate) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WithBaseExecutionUpdate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.LastKnownState) {
		toSerialize["lastKnownState"] = o.LastKnownState
	}
	return toSerialize, nil
}

type NullableWithBaseExecutionUpdate struct {
	value *WithBaseExecutionUpdate
	isSet bool
}

func (v NullableWithBaseExecutionUpdate) Get() *WithBaseExecutionUpdate {
	return v.value
}

func (v *NullableWithBaseExecutionUpdate) Set(val *WithBaseExecutionUpdate) {
	v.value = val
	v.isSet = true
}

func (v NullableWithBaseExecutionUpdate) IsSet() bool {
	return v.isSet
}

func (v *NullableWithBaseExecutionUpdate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWithBaseExecutionUpdate(val *WithBaseExecutionUpdate) *NullableWithBaseExecutionUpdate {
	return &NullableWithBaseExecutionUpdate{value: val, isSet: true}
}

func (v NullableWithBaseExecutionUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWithBaseExecutionUpdate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}