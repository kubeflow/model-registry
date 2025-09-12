package datastore

import "reflect"

// RepoSetSpec is the specification to construct a RepoSet. Each `any` is
// expected to be a function pointer to an initialization function that returns
// a repository instance and an optional error.
//
// RepoSet implementations must pass arguments to the initialization functions
// that are required by the repository (database handles, HTTP clients, etc.)
// and should also provide values for the following types:
//
//   - int64: type id
//   - ArtifactTypeMap: Map of all artifact names to type IDs
//   - ContextTypeMap: Map of all context names to type IDs
//   - ExecutionTypeMap: Map of all execution names to type IDs
type RepoSetSpec struct {
	// Maps artifact type names to an initializer.
	ArtifactTypes map[string]any

	// Maps context type names to an initializer.
	ContextTypes map[string]any

	// Maps execution type names to an initializer.
	ExecutionTypes map[string]any

	// Any repos initialization functions that don't map to a type can be
	// added here.
	Others []any
}

// AllNames returns all the type names in the spec.
func (s *RepoSetSpec) AllNames() []string {
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

// RepoSet holds repository implementions.
type RepoSet interface {
	// TypeMap returns a map of type names to IDs
	TypeMap() map[string]int64

	// Repository returns a repository instance of the specified type.
	Repository(t reflect.Type) (any, error)
}

// ArtifactTypeMap maps artifact type names to IDs
type ArtifactTypeMap map[string]int64

// ContextTypeMap maps context type names to IDs
type ContextTypeMap map[string]int64

// ExecutionTypeMap maps execution type names to IDs
type ExecutionTypeMap map[string]int64
