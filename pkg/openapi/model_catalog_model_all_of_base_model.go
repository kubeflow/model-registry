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

// checks if the CatalogModelAllOfBaseModel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CatalogModelAllOfBaseModel{}

// CatalogModelAllOfBaseModel struct for CatalogModelAllOfBaseModel
type CatalogModelAllOfBaseModel struct {
	// Name of the catalog for an external base model. Omit for models in the same catalog.
	Catalog *string `json:"catalog,omitempty"`
	// Name of the repository in an external catalog where the base model exists. Omit for models in the same catalog.
	Repository *string `json:"repository,omitempty"`
	Name       *string `json:"name,omitempty"`
}

// NewCatalogModelAllOfBaseModel instantiates a new CatalogModelAllOfBaseModel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCatalogModelAllOfBaseModel() *CatalogModelAllOfBaseModel {
	this := CatalogModelAllOfBaseModel{}
	return &this
}

// NewCatalogModelAllOfBaseModelWithDefaults instantiates a new CatalogModelAllOfBaseModel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCatalogModelAllOfBaseModelWithDefaults() *CatalogModelAllOfBaseModel {
	this := CatalogModelAllOfBaseModel{}
	return &this
}

// GetCatalog returns the Catalog field value if set, zero value otherwise.
func (o *CatalogModelAllOfBaseModel) GetCatalog() string {
	if o == nil || IsNil(o.Catalog) {
		var ret string
		return ret
	}
	return *o.Catalog
}

// GetCatalogOk returns a tuple with the Catalog field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CatalogModelAllOfBaseModel) GetCatalogOk() (*string, bool) {
	if o == nil || IsNil(o.Catalog) {
		return nil, false
	}
	return o.Catalog, true
}

// HasCatalog returns a boolean if a field has been set.
func (o *CatalogModelAllOfBaseModel) HasCatalog() bool {
	if o != nil && !IsNil(o.Catalog) {
		return true
	}

	return false
}

// SetCatalog gets a reference to the given string and assigns it to the Catalog field.
func (o *CatalogModelAllOfBaseModel) SetCatalog(v string) {
	o.Catalog = &v
}

// GetRepository returns the Repository field value if set, zero value otherwise.
func (o *CatalogModelAllOfBaseModel) GetRepository() string {
	if o == nil || IsNil(o.Repository) {
		var ret string
		return ret
	}
	return *o.Repository
}

// GetRepositoryOk returns a tuple with the Repository field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CatalogModelAllOfBaseModel) GetRepositoryOk() (*string, bool) {
	if o == nil || IsNil(o.Repository) {
		return nil, false
	}
	return o.Repository, true
}

// HasRepository returns a boolean if a field has been set.
func (o *CatalogModelAllOfBaseModel) HasRepository() bool {
	if o != nil && !IsNil(o.Repository) {
		return true
	}

	return false
}

// SetRepository gets a reference to the given string and assigns it to the Repository field.
func (o *CatalogModelAllOfBaseModel) SetRepository(v string) {
	o.Repository = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *CatalogModelAllOfBaseModel) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CatalogModelAllOfBaseModel) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *CatalogModelAllOfBaseModel) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *CatalogModelAllOfBaseModel) SetName(v string) {
	o.Name = &v
}

func (o CatalogModelAllOfBaseModel) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CatalogModelAllOfBaseModel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Catalog) {
		toSerialize["catalog"] = o.Catalog
	}
	if !IsNil(o.Repository) {
		toSerialize["repository"] = o.Repository
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	return toSerialize, nil
}

type NullableCatalogModelAllOfBaseModel struct {
	value *CatalogModelAllOfBaseModel
	isSet bool
}

func (v NullableCatalogModelAllOfBaseModel) Get() *CatalogModelAllOfBaseModel {
	return v.value
}

func (v *NullableCatalogModelAllOfBaseModel) Set(val *CatalogModelAllOfBaseModel) {
	v.value = val
	v.isSet = true
}

func (v NullableCatalogModelAllOfBaseModel) IsSet() bool {
	return v.isSet
}

func (v *NullableCatalogModelAllOfBaseModel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCatalogModelAllOfBaseModel(val *CatalogModelAllOfBaseModel) *NullableCatalogModelAllOfBaseModel {
	return &NullableCatalogModelAllOfBaseModel{value: val, isSet: true}
}

func (v NullableCatalogModelAllOfBaseModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCatalogModelAllOfBaseModel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
