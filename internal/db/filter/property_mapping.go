package filter

import (
	platformfilter "github.com/kubeflow/hub/internal/platform/db/filter"
)

// Property location types
type PropertyLocation = platformfilter.PropertyLocation

const (
	EntityTable   = platformfilter.EntityTable
	PropertyTable = platformfilter.PropertyTable
	Custom        = platformfilter.Custom
	RelatedEntity = platformfilter.RelatedEntity
)

// Related entity types
type RelatedEntityType = platformfilter.RelatedEntityType

const (
	RelatedEntityArtifact  = platformfilter.RelatedEntityArtifact
	RelatedEntityContext   = platformfilter.RelatedEntityContext
	RelatedEntityExecution = platformfilter.RelatedEntityExecution
)

// Property definition types
type PropertyDefinition = platformfilter.PropertyDefinition
type EntityPropertyMap = platformfilter.EntityPropertyMap

// GetPropertyDefinition returns the property definition for a given entity type and property name
var GetPropertyDefinition = platformfilter.GetPropertyDefinition
