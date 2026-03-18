package datastore

import "reflect"

// Spec is the specification for the datastore.
type Spec struct {
	ArtifactTypes  map[string]*SpecType
	ContextTypes   map[string]*SpecType
	ExecutionTypes map[string]*SpecType
	Others         []any
}

func NewSpec() *Spec {
	return &Spec{
		ArtifactTypes:  map[string]*SpecType{},
		ContextTypes:   map[string]*SpecType{},
		ExecutionTypes: map[string]*SpecType{},
		Others:         []any{},
	}
}

func (s *Spec) AddArtifact(name string, t *SpecType) *Spec {
	s.ArtifactTypes[name] = t
	return s
}

func (s *Spec) AddContext(name string, t *SpecType) *Spec {
	s.ContextTypes[name] = t
	return s
}

func (s *Spec) AddExecution(name string, t *SpecType) *Spec {
	s.ExecutionTypes[name] = t
	return s
}

func (s *Spec) AddOther(initFn any) *Spec {
	s.Others = append(s.Others, initFn)
	return s
}

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

type SpecType struct {
	InitFn     any
	Properties map[string]PropertyType
}

func NewSpecType(initFn any) *SpecType {
	return &SpecType{
		InitFn:     initFn,
		Properties: map[string]PropertyType{},
	}
}

func (st *SpecType) AddInt(name string) *SpecType {
	st.Properties[name] = PropertyTypeInt
	return st
}

func (st *SpecType) AddDouble(name string) *SpecType {
	st.Properties[name] = PropertyTypeDouble
	return st
}

func (st *SpecType) AddString(name string) *SpecType {
	st.Properties[name] = PropertyTypeString
	return st
}

func (st *SpecType) AddStruct(name string) *SpecType {
	st.Properties[name] = PropertyTypeStruct
	return st
}

func (st *SpecType) AddProto(name string) *SpecType {
	st.Properties[name] = PropertyTypeProto
	return st
}

func (st *SpecType) AddBoolean(name string) *SpecType {
	st.Properties[name] = PropertyTypeBoolean
	return st
}

type RepoSet interface {
	TypeMap() map[string]int32
	Repository(t reflect.Type) (any, error)
}

type ArtifactTypeMap map[string]int32
type ContextTypeMap map[string]int32
type ExecutionTypeMap map[string]int32
