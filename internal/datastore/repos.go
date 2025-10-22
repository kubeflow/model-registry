package datastore

import "reflect"

// Spec is the specification for the datastore.
// Each entry
type Spec struct {
	// Maps artifact type names to an initializer.
	ArtifactTypes map[string]*SpecType

	// Maps context type names to an initializer.
	ContextTypes map[string]*SpecType

	// Maps execution type names to an initializer.
	ExecutionTypes map[string]*SpecType

	// Any repo initialization functions that don't map to a type can be
	// added here.
	Others []any
}

// NewSpec returns an empty Spec instance.
func NewSpec() *Spec {
	return &Spec{
		ArtifactTypes:  map[string]*SpecType{},
		ContextTypes:   map[string]*SpecType{},
		ExecutionTypes: map[string]*SpecType{},
		Others:         []any{},
	}
}

// AddArtifact adds an artifact type to the spec.
func (s *Spec) AddArtifact(name string, t *SpecType) *Spec {
	s.ArtifactTypes[name] = t
	return s
}

// AddContext adds a context type to the spec.
func (s *Spec) AddContext(name string, t *SpecType) *Spec {
	s.ContextTypes[name] = t
	return s
}

// AddExecution adds an execution type to the spec.
func (s *Spec) AddExecution(name string, t *SpecType) *Spec {
	s.ExecutionTypes[name] = t
	return s
}

// AddOther adds a repo initializer to the spec.
func (s *Spec) AddOther(initFn any) *Spec {
	s.Others = append(s.Others, initFn)
	return s
}

// AllNames returns all the type names in the spec.
func (s *Spec) AllNames() []string {
	names := make([]string, 0, len(s.ArtifactTypes)+len(s.ContextTypes)+len(s.ExecutionTypes))
	for n := range s.ArtifactTypes {
		names = append(names, n)
	}
	for n := range s.ContextTypes {
		names = append(names, n)
	}
	for n := range s.ExecutionTypes {
		names = append(names, n)
	}
	return names
}

// PropertyType is the data type of a property's value.
type PropertyType int32

const (
	PropertyTypeUnknown = iota
	PropertyTypeInt
	PropertyTypeDouble
	PropertyTypeString
	PropertyTypeStruct
	PropertyTypeProto
	PropertyTypeBoolean
)

// SpecType is a single type in the spec.
type SpecType struct {
	// InitFn is a pointer to an initialization function that returns the
	// repository instance and an optional error.
	//
	// Data Store implementations must pass arguments to the initialization
	// functions that are required by the repository (database handles,
	// HTTP clients, etc.) and should also provide values for the following
	// types:
	//
	//   - int64: type id
	//   - ArtifactTypeMap: Map of all artifact names to type IDs
	//   - ContextTypeMap: Map of all context names to type IDs
	//   - ExecutionTypeMap: Map of all execution names to type IDs
	InitFn any

	// Defined (non-custom) properties of the type.
	Properties map[string]PropertyType
}

// NewSpecType creates a SpecType instance.
func NewSpecType(initFn any) *SpecType {
	return &SpecType{
		InitFn:     initFn,
		Properties: map[string]PropertyType{},
	}
}

// AddInt adds an int property to the type spec.
func (st *SpecType) AddInt(name string) *SpecType {
	st.Properties[name] = PropertyTypeInt
	return st
}

// AddDouble adds a double property to the type spec.
func (st *SpecType) AddDouble(name string) *SpecType {
	st.Properties[name] = PropertyTypeDouble
	return st
}

// AddString adds a string property to the type spec.
func (st *SpecType) AddString(name string) *SpecType {
	st.Properties[name] = PropertyTypeString
	return st
}

// AddStruct adds a struct property to the type spec.
func (st *SpecType) AddStruct(name string) *SpecType {
	st.Properties[name] = PropertyTypeStruct
	return st
}

// AddProto adds a proto property to the type spec.
func (st *SpecType) AddProto(name string) *SpecType {
	st.Properties[name] = PropertyTypeProto
	return st
}

// AddBoolean adds a boolean property to the type spec.
func (st *SpecType) AddBoolean(name string) *SpecType {
	st.Properties[name] = PropertyTypeBoolean
	return st
}

// RepoSet holds repository implementions.
type RepoSet interface {
	// TypeMap returns a map of type names to IDs
	TypeMap() map[string]int32

	// Repository returns a repository instance of the specified type.
	Repository(t reflect.Type) (any, error)
}

// ArtifactTypeMap maps artifact type names to IDs
type ArtifactTypeMap map[string]int32

// ContextTypeMap maps context type names to IDs
type ContextTypeMap map[string]int32

// ExecutionTypeMap maps execution type names to IDs
type ExecutionTypeMap map[string]int32
