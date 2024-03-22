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

// checks if the ModelArtifactList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ModelArtifactList{}

// ModelArtifactList List of ModelArtifact entities.
type ModelArtifactList struct {
	// Array of `ModelArtifact` entities.
	Items []ModelArtifact `json:"items,omitempty"`
	// Token to use to retrieve next page of results.
	NextPageToken string `json:"nextPageToken"`
	// Maximum number of resources to return in the result.
	PageSize int32 `json:"pageSize"`
	// Number of items in result list.
	Size int32 `json:"size"`
}

// NewModelArtifactList instantiates a new ModelArtifactList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewModelArtifactList(nextPageToken string, pageSize int32, size int32) *ModelArtifactList {
	this := ModelArtifactList{}
	this.NextPageToken = nextPageToken
	this.PageSize = pageSize
	this.Size = size
	return &this
}

// NewModelArtifactListWithDefaults instantiates a new ModelArtifactList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewModelArtifactListWithDefaults() *ModelArtifactList {
	this := ModelArtifactList{}
	return &this
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *ModelArtifactList) GetItems() []ModelArtifact {
	if o == nil || IsNil(o.Items) {
		var ret []ModelArtifact
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ModelArtifactList) GetItemsOk() ([]ModelArtifact, bool) {
	if o == nil || IsNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *ModelArtifactList) HasItems() bool {
	if o != nil && !IsNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []ModelArtifact and assigns it to the Items field.
func (o *ModelArtifactList) SetItems(v []ModelArtifact) {
	o.Items = v
}

// GetNextPageToken returns the NextPageToken field value
func (o *ModelArtifactList) GetNextPageToken() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value
// and a boolean to check if the value has been set.
func (o *ModelArtifactList) GetNextPageTokenOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NextPageToken, true
}

// SetNextPageToken sets field value
func (o *ModelArtifactList) SetNextPageToken(v string) {
	o.NextPageToken = v
}

// GetPageSize returns the PageSize field value
func (o *ModelArtifactList) GetPageSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.PageSize
}

// GetPageSizeOk returns a tuple with the PageSize field value
// and a boolean to check if the value has been set.
func (o *ModelArtifactList) GetPageSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PageSize, true
}

// SetPageSize sets field value
func (o *ModelArtifactList) SetPageSize(v int32) {
	o.PageSize = v
}

// GetSize returns the Size field value
func (o *ModelArtifactList) GetSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Size
}

// GetSizeOk returns a tuple with the Size field value
// and a boolean to check if the value has been set.
func (o *ModelArtifactList) GetSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Size, true
}

// SetSize sets field value
func (o *ModelArtifactList) SetSize(v int32) {
	o.Size = v
}

func (o ModelArtifactList) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ModelArtifactList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	toSerialize["nextPageToken"] = o.NextPageToken
	toSerialize["pageSize"] = o.PageSize
	toSerialize["size"] = o.Size
	return toSerialize, nil
}

type NullableModelArtifactList struct {
	value *ModelArtifactList
	isSet bool
}

func (v NullableModelArtifactList) Get() *ModelArtifactList {
	return v.value
}

func (v *NullableModelArtifactList) Set(val *ModelArtifactList) {
	v.value = val
	v.isSet = true
}

func (v NullableModelArtifactList) IsSet() bool {
	return v.isSet
}

func (v *NullableModelArtifactList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableModelArtifactList(val *ModelArtifactList) *NullableModelArtifactList {
	return &NullableModelArtifactList{value: val, isSet: true}
}

func (v NullableModelArtifactList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableModelArtifactList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
