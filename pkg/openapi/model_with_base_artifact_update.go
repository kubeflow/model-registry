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

// checks if the WithBaseArtifactUpdate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WithBaseArtifactUpdate{}

// WithBaseArtifactUpdate struct for WithBaseArtifactUpdate
type WithBaseArtifactUpdate struct {
	// The uniform resource identifier of the physical artifact. May be empty if there is no physical artifact.
	Uri   *string        `json:"uri,omitempty"`
	State *ArtifactState `json:"state,omitempty"`
}

// NewWithBaseArtifactUpdate instantiates a new WithBaseArtifactUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWithBaseArtifactUpdate() *WithBaseArtifactUpdate {
	this := WithBaseArtifactUpdate{}
	var state ArtifactState = ARTIFACTSTATE_UNKNOWN
	this.State = &state
	return &this
}

// NewWithBaseArtifactUpdateWithDefaults instantiates a new WithBaseArtifactUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWithBaseArtifactUpdateWithDefaults() *WithBaseArtifactUpdate {
	this := WithBaseArtifactUpdate{}
	var state ArtifactState = ARTIFACTSTATE_UNKNOWN
	this.State = &state
	return &this
}

// GetUri returns the Uri field value if set, zero value otherwise.
func (o *WithBaseArtifactUpdate) GetUri() string {
	if o == nil || IsNil(o.Uri) {
		var ret string
		return ret
	}
	return *o.Uri
}

// GetUriOk returns a tuple with the Uri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WithBaseArtifactUpdate) GetUriOk() (*string, bool) {
	if o == nil || IsNil(o.Uri) {
		return nil, false
	}
	return o.Uri, true
}

// HasUri returns a boolean if a field has been set.
func (o *WithBaseArtifactUpdate) HasUri() bool {
	if o != nil && !IsNil(o.Uri) {
		return true
	}

	return false
}

// SetUri gets a reference to the given string and assigns it to the Uri field.
func (o *WithBaseArtifactUpdate) SetUri(v string) {
	o.Uri = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *WithBaseArtifactUpdate) GetState() ArtifactState {
	if o == nil || IsNil(o.State) {
		var ret ArtifactState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WithBaseArtifactUpdate) GetStateOk() (*ArtifactState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *WithBaseArtifactUpdate) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given ArtifactState and assigns it to the State field.
func (o *WithBaseArtifactUpdate) SetState(v ArtifactState) {
	o.State = &v
}

func (o WithBaseArtifactUpdate) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WithBaseArtifactUpdate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Uri) {
		toSerialize["uri"] = o.Uri
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	return toSerialize, nil
}

type NullableWithBaseArtifactUpdate struct {
	value *WithBaseArtifactUpdate
	isSet bool
}

func (v NullableWithBaseArtifactUpdate) Get() *WithBaseArtifactUpdate {
	return v.value
}

func (v *NullableWithBaseArtifactUpdate) Set(val *WithBaseArtifactUpdate) {
	v.value = val
	v.isSet = true
}

func (v NullableWithBaseArtifactUpdate) IsSet() bool {
	return v.isSet
}

func (v *NullableWithBaseArtifactUpdate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWithBaseArtifactUpdate(val *WithBaseArtifactUpdate) *NullableWithBaseArtifactUpdate {
	return &NullableWithBaseArtifactUpdate{value: val, isSet: true}
}

func (v NullableWithBaseArtifactUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWithBaseArtifactUpdate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
