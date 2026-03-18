package repository

import (
	"github.com/kubeflow/model-registry/internal/platform/db/entity"
	"github.com/kubeflow/model-registry/internal/platform/db/schema"
)

// MapPropertiesToArtifactProperty converts entity.Properties to schema.ArtifactProperty
func MapPropertiesToArtifactProperty(prop entity.Properties, artifactID int32, isCustomProperty bool) schema.ArtifactProperty {
	return schema.ArtifactProperty{
		ArtifactID:       artifactID,
		Name:             prop.Name,
		IsCustomProperty: isCustomProperty,
		IntValue:         prop.IntValue,
		DoubleValue:      prop.DoubleValue,
		StringValue:      prop.StringValue,
		BoolValue:        prop.BoolValue,
		ByteValue:        prop.ByteValue,
		ProtoValue:       prop.ProtoValue,
	}
}

// MapPropertiesToContextProperty converts entity.Properties to schema.ContextProperty
func MapPropertiesToContextProperty(prop entity.Properties, contextID int32, isCustomProperty bool) schema.ContextProperty {
	return schema.ContextProperty{
		ContextID:        contextID,
		Name:             prop.Name,
		IsCustomProperty: isCustomProperty,
		IntValue:         prop.IntValue,
		DoubleValue:      prop.DoubleValue,
		StringValue:      prop.StringValue,
		BoolValue:        prop.BoolValue,
		ByteValue:        prop.ByteValue,
		ProtoValue:       prop.ProtoValue,
	}
}

// MapPropertiesToExecutionProperty converts entity.Properties to schema.ExecutionProperty
func MapPropertiesToExecutionProperty(prop entity.Properties, executionID int32, isCustomProperty bool) schema.ExecutionProperty {
	return schema.ExecutionProperty{
		ExecutionID:      executionID,
		Name:             prop.Name,
		IsCustomProperty: isCustomProperty,
		IntValue:         prop.IntValue,
		DoubleValue:      prop.DoubleValue,
		StringValue:      prop.StringValue,
		BoolValue:        prop.BoolValue,
		ByteValue:        prop.ByteValue,
		ProtoValue:       prop.ProtoValue,
	}
}

// MapArtifactPropertyToProperties converts schema.ArtifactProperty to entity.Properties
func MapArtifactPropertyToProperties(artProperty schema.ArtifactProperty) entity.Properties {
	return entity.Properties{
		Name:             artProperty.Name,
		IsCustomProperty: artProperty.IsCustomProperty,
		IntValue:         artProperty.IntValue,
		DoubleValue:      artProperty.DoubleValue,
		StringValue:      artProperty.StringValue,
		BoolValue:        artProperty.BoolValue,
		ByteValue:        artProperty.ByteValue,
		ProtoValue:       artProperty.ProtoValue,
	}
}

// MapContextPropertyToProperties converts schema.ContextProperty to entity.Properties
func MapContextPropertyToProperties(contextProperty schema.ContextProperty) entity.Properties {
	return entity.Properties{
		Name:             contextProperty.Name,
		IsCustomProperty: contextProperty.IsCustomProperty,
		IntValue:         contextProperty.IntValue,
		DoubleValue:      contextProperty.DoubleValue,
		StringValue:      contextProperty.StringValue,
		BoolValue:        contextProperty.BoolValue,
		ByteValue:        contextProperty.ByteValue,
		ProtoValue:       contextProperty.ProtoValue,
	}
}

// MapExecutionPropertyToProperties converts schema.ExecutionProperty to entity.Properties
func MapExecutionPropertyToProperties(executionProperty schema.ExecutionProperty) entity.Properties {
	return entity.Properties{
		Name:             executionProperty.Name,
		IsCustomProperty: executionProperty.IsCustomProperty,
		IntValue:         executionProperty.IntValue,
		DoubleValue:      executionProperty.DoubleValue,
		StringValue:      executionProperty.StringValue,
		BoolValue:        executionProperty.BoolValue,
		ByteValue:        executionProperty.ByteValue,
		ProtoValue:       executionProperty.ProtoValue,
	}
}
