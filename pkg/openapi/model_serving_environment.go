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

// checks if the ServingEnvironment type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ServingEnvironment{}

// ServingEnvironment A Model Serving environment for serving `RegisteredModels`.
type ServingEnvironment struct {
	// User provided custom properties which are not defined by its type.
	CustomProperties *map[string]MetadataValue `json:"customProperties,omitempty"`
	// An optional description about the resource.
	Description *string `json:"description,omitempty"`
	// The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance.
	ExternalId *string `json:"externalId,omitempty"`
	// The name of the ServingEnvironment.
	Name string `json:"name"`
	// The unique server generated id of the resource.
	Id *string `json:"id,omitempty"`
	// Output only. Create time of the resource in millisecond since epoch.
	CreateTimeSinceEpoch *string `json:"createTimeSinceEpoch,omitempty"`
	// Output only. Last update time of the resource since epoch in millisecond since epoch.
	LastUpdateTimeSinceEpoch *string `json:"lastUpdateTimeSinceEpoch,omitempty"`
}

// NewServingEnvironment instantiates a new ServingEnvironment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServingEnvironment(name string) *ServingEnvironment {
	this := ServingEnvironment{}
	this.Name = name
	return &this
}

// NewServingEnvironmentWithDefaults instantiates a new ServingEnvironment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServingEnvironmentWithDefaults() *ServingEnvironment {
	this := ServingEnvironment{}
	return &this
}

// GetCustomProperties returns the CustomProperties field value if set, zero value otherwise.
func (o *ServingEnvironment) GetCustomProperties() map[string]MetadataValue {
	if o == nil || IsNil(o.CustomProperties) {
		var ret map[string]MetadataValue
		return ret
	}
	return *o.CustomProperties
}

// GetCustomPropertiesOk returns a tuple with the CustomProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServingEnvironment) GetCustomPropertiesOk() (*map[string]MetadataValue, bool) {
	if o == nil || IsNil(o.CustomProperties) {
		return nil, false
	}
	return o.CustomProperties, true
}

// HasCustomProperties returns a boolean if a field has been set.
func (o *ServingEnvironment) HasCustomProperties() bool {
	if o != nil && !IsNil(o.CustomProperties) {
		return true
	}

	return false
}

// SetCustomProperties gets a reference to the given map[string]MetadataValue and assigns it to the CustomProperties field.
func (o *ServingEnvironment) SetCustomProperties(v map[string]MetadataValue) {
	o.CustomProperties = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *ServingEnvironment) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServingEnvironment) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *ServingEnvironment) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *ServingEnvironment) SetDescription(v string) {
	o.Description = &v
}

// GetExternalId returns the ExternalId field value if set, zero value otherwise.
func (o *ServingEnvironment) GetExternalId() string {
	if o == nil || IsNil(o.ExternalId) {
		var ret string
		return ret
	}
	return *o.ExternalId
}

// GetExternalIdOk returns a tuple with the ExternalId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServingEnvironment) GetExternalIdOk() (*string, bool) {
	if o == nil || IsNil(o.ExternalId) {
		return nil, false
	}
	return o.ExternalId, true
}

// HasExternalId returns a boolean if a field has been set.
func (o *ServingEnvironment) HasExternalId() bool {
	if o != nil && !IsNil(o.ExternalId) {
		return true
	}

	return false
}

// SetExternalId gets a reference to the given string and assigns it to the ExternalId field.
func (o *ServingEnvironment) SetExternalId(v string) {
	o.ExternalId = &v
}

// GetName returns the Name field value
func (o *ServingEnvironment) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *ServingEnvironment) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *ServingEnvironment) SetName(v string) {
	o.Name = v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *ServingEnvironment) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServingEnvironment) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ServingEnvironment) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ServingEnvironment) SetId(v string) {
	o.Id = &v
}

// GetCreateTimeSinceEpoch returns the CreateTimeSinceEpoch field value if set, zero value otherwise.
func (o *ServingEnvironment) GetCreateTimeSinceEpoch() string {
	if o == nil || IsNil(o.CreateTimeSinceEpoch) {
		var ret string
		return ret
	}
	return *o.CreateTimeSinceEpoch
}

// GetCreateTimeSinceEpochOk returns a tuple with the CreateTimeSinceEpoch field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServingEnvironment) GetCreateTimeSinceEpochOk() (*string, bool) {
	if o == nil || IsNil(o.CreateTimeSinceEpoch) {
		return nil, false
	}
	return o.CreateTimeSinceEpoch, true
}

// HasCreateTimeSinceEpoch returns a boolean if a field has been set.
func (o *ServingEnvironment) HasCreateTimeSinceEpoch() bool {
	if o != nil && !IsNil(o.CreateTimeSinceEpoch) {
		return true
	}

	return false
}

// SetCreateTimeSinceEpoch gets a reference to the given string and assigns it to the CreateTimeSinceEpoch field.
func (o *ServingEnvironment) SetCreateTimeSinceEpoch(v string) {
	o.CreateTimeSinceEpoch = &v
}

// GetLastUpdateTimeSinceEpoch returns the LastUpdateTimeSinceEpoch field value if set, zero value otherwise.
func (o *ServingEnvironment) GetLastUpdateTimeSinceEpoch() string {
	if o == nil || IsNil(o.LastUpdateTimeSinceEpoch) {
		var ret string
		return ret
	}
	return *o.LastUpdateTimeSinceEpoch
}

// GetLastUpdateTimeSinceEpochOk returns a tuple with the LastUpdateTimeSinceEpoch field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServingEnvironment) GetLastUpdateTimeSinceEpochOk() (*string, bool) {
	if o == nil || IsNil(o.LastUpdateTimeSinceEpoch) {
		return nil, false
	}
	return o.LastUpdateTimeSinceEpoch, true
}

// HasLastUpdateTimeSinceEpoch returns a boolean if a field has been set.
func (o *ServingEnvironment) HasLastUpdateTimeSinceEpoch() bool {
	if o != nil && !IsNil(o.LastUpdateTimeSinceEpoch) {
		return true
	}

	return false
}

// SetLastUpdateTimeSinceEpoch gets a reference to the given string and assigns it to the LastUpdateTimeSinceEpoch field.
func (o *ServingEnvironment) SetLastUpdateTimeSinceEpoch(v string) {
	o.LastUpdateTimeSinceEpoch = &v
}

func (o ServingEnvironment) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ServingEnvironment) ToMap() (map[string]interface{}, error) {
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
	toSerialize["name"] = o.Name
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.CreateTimeSinceEpoch) {
		toSerialize["createTimeSinceEpoch"] = o.CreateTimeSinceEpoch
	}
	if !IsNil(o.LastUpdateTimeSinceEpoch) {
		toSerialize["lastUpdateTimeSinceEpoch"] = o.LastUpdateTimeSinceEpoch
	}
	return toSerialize, nil
}

type NullableServingEnvironment struct {
	value *ServingEnvironment
	isSet bool
}

func (v NullableServingEnvironment) Get() *ServingEnvironment {
	return v.value
}

func (v *NullableServingEnvironment) Set(val *ServingEnvironment) {
	v.value = val
	v.isSet = true
}

func (v NullableServingEnvironment) IsSet() bool {
	return v.isSet
}

func (v *NullableServingEnvironment) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableServingEnvironment(val *ServingEnvironment) *NullableServingEnvironment {
	return &NullableServingEnvironment{value: val, isSet: true}
}

func (v NullableServingEnvironment) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableServingEnvironment) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
