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

// checks if the RegisteredModelList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RegisteredModelList{}

// RegisteredModelList List of RegisteredModels.
type RegisteredModelList struct {
	// Token to use to retrieve next page of results.
	NextPageToken string `json:"nextPageToken"`
	// Maximum number of resources to return in the result.
	PageSize int32 `json:"pageSize"`
	// Number of items in result list.
	Size int32 `json:"size"`
	//
	Items []RegisteredModel `json:"items"`
}

// NewRegisteredModelList instantiates a new RegisteredModelList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRegisteredModelList(nextPageToken string, pageSize int32, size int32, items []RegisteredModel) *RegisteredModelList {
	this := RegisteredModelList{}
	this.NextPageToken = nextPageToken
	this.PageSize = pageSize
	this.Size = size
	this.Items = items
	return &this
}

// NewRegisteredModelListWithDefaults instantiates a new RegisteredModelList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRegisteredModelListWithDefaults() *RegisteredModelList {
	this := RegisteredModelList{}
	return &this
}

// GetNextPageToken returns the NextPageToken field value
func (o *RegisteredModelList) GetNextPageToken() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value
// and a boolean to check if the value has been set.
func (o *RegisteredModelList) GetNextPageTokenOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NextPageToken, true
}

// SetNextPageToken sets field value
func (o *RegisteredModelList) SetNextPageToken(v string) {
	o.NextPageToken = v
}

// GetPageSize returns the PageSize field value
func (o *RegisteredModelList) GetPageSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.PageSize
}

// GetPageSizeOk returns a tuple with the PageSize field value
// and a boolean to check if the value has been set.
func (o *RegisteredModelList) GetPageSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PageSize, true
}

// SetPageSize sets field value
func (o *RegisteredModelList) SetPageSize(v int32) {
	o.PageSize = v
}

// GetSize returns the Size field value
func (o *RegisteredModelList) GetSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Size
}

// GetSizeOk returns a tuple with the Size field value
// and a boolean to check if the value has been set.
func (o *RegisteredModelList) GetSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Size, true
}

// SetSize sets field value
func (o *RegisteredModelList) SetSize(v int32) {
	o.Size = v
}

// GetItems returns the Items field value
func (o *RegisteredModelList) GetItems() []RegisteredModel {
	if o == nil {
		var ret []RegisteredModel
		return ret
	}

	return o.Items
}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
func (o *RegisteredModelList) GetItemsOk() ([]RegisteredModel, bool) {
	if o == nil {
		return nil, false
	}
	return o.Items, true
}

// SetItems sets field value
func (o *RegisteredModelList) SetItems(v []RegisteredModel) {
	o.Items = v
}

func (o RegisteredModelList) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RegisteredModelList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["nextPageToken"] = o.NextPageToken
	toSerialize["pageSize"] = o.PageSize
	toSerialize["size"] = o.Size
	toSerialize["items"] = o.Items
	return toSerialize, nil
}

type NullableRegisteredModelList struct {
	value *RegisteredModelList
	isSet bool
}

func (v NullableRegisteredModelList) Get() *RegisteredModelList {
	return v.value
}

func (v *NullableRegisteredModelList) Set(val *RegisteredModelList) {
	v.value = val
	v.isSet = true
}

func (v NullableRegisteredModelList) IsSet() bool {
	return v.isSet
}

func (v *NullableRegisteredModelList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRegisteredModelList(val *RegisteredModelList) *NullableRegisteredModelList {
	return &NullableRegisteredModelList{value: val, isSet: true}
}

func (v NullableRegisteredModelList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRegisteredModelList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
