package filter

// RestEntityType represents the specific REST API entity type
type RestEntityType string

const (
	// Context-based REST entities
	RestEntityRegisteredModel    RestEntityType = "RegisteredModel"
	RestEntityModelVersion       RestEntityType = "ModelVersion"
	RestEntityInferenceService   RestEntityType = "InferenceService"
	RestEntityServingEnvironment RestEntityType = "ServingEnvironment"
	RestEntityExperiment         RestEntityType = "Experiment"
	RestEntityExperimentRun      RestEntityType = "ExperimentRun"

	// Artifact-based REST entities
	RestEntityModelArtifact RestEntityType = "ModelArtifact"
	RestEntityDocArtifact   RestEntityType = "DocArtifact"
	RestEntityDataSet       RestEntityType = "DataSet"
	RestEntityMetric        RestEntityType = "Metric"
	RestEntityParameter     RestEntityType = "Parameter"

	// Execution-based REST entities
	RestEntityServeModel RestEntityType = "ServeModel"
)

// isChildEntity returns true if the REST entity type uses prefixed names (parentId:name)
func isChildEntity(entityType RestEntityType) bool {
	// Only top-level entities don't use prefixed names
	switch entityType {
	case RestEntityRegisteredModel, RestEntityExperiment, RestEntityServingEnvironment:
		return false
	default:
		// All other entities are child entities that use prefixed names
		return true
	}
}

// RestEntityPropertyMap maps REST entity types to their allowed properties
var RestEntityPropertyMap = map[RestEntityType]map[string]bool{
	// Context-based entities
	RestEntityRegisteredModel: {
		// Common Context properties
		"id": true, "name": true, "externalId": true,
		"createTimeSinceEpoch": true, "lastUpdateTimeSinceEpoch": true,
		// RegisteredModel-specific properties
		"state": true, "owner": true,
		// No experiment or serving-specific properties allowed
	},

	RestEntityModelVersion: {
		// Common Context properties
		"id": true, "name": true, "externalId": true,
		"createTimeSinceEpoch": true, "lastUpdateTimeSinceEpoch": true,
		// ModelVersion-specific properties
		"registeredModelId": true, "state": true, "author": true,
		// No experiment or serving-specific properties allowed
	},

	RestEntityInferenceService: {
		// Common Context properties
		"id": true, "name": true, "externalId": true,
		"createTimeSinceEpoch": true, "lastUpdateTimeSinceEpoch": true,
		// InferenceService-specific properties
		"registeredModelId": true, "modelVersionId": true, "servingEnvironmentId": true,
		"runtime": true, "desiredState": true,
		// No experiment-specific properties allowed
	},

	RestEntityServingEnvironment: {
		// Common Context properties
		"id": true, "name": true, "externalId": true,
		"createTimeSinceEpoch": true, "lastUpdateTimeSinceEpoch": true,
		// ServingEnvironment-specific properties (minimal)
		// No inference or experiment-specific properties allowed
	},

	RestEntityExperiment: {
		// Common Context properties
		"id": true, "name": true, "externalId": true,
		"createTimeSinceEpoch": true, "lastUpdateTimeSinceEpoch": true,
		// Experiment-specific properties
		"state": true, "owner": true,
		// No serving or model-specific properties allowed
	},

	RestEntityExperimentRun: {
		// Common Context properties
		"id": true, "name": true, "externalId": true,
		"createTimeSinceEpoch": true, "lastUpdateTimeSinceEpoch": true,
		// ExperimentRun-specific properties
		"experimentId": true, "startTimeSinceEpoch": true, "endTimeSinceEpoch": true,
		"status": true, "state": true, "owner": true,
		// No serving or model-specific properties allowed
	},

	// Artifact-based entities
	RestEntityModelArtifact: {
		// Common Artifact properties
		"id": true, "name": true, "externalId": true,
		"createTimeSinceEpoch": true, "lastUpdateTimeSinceEpoch": true,
		"uri": true, "state": true,
		// ModelArtifact-specific properties
		"modelFormatName": true, "modelFormatVersion": true,
		"storageKey": true, "storagePath": true, "serviceAccountName": true,
		"modelSourceKind": true, "modelSourceClass": true, "modelSourceGroup": true,
		"modelSourceId": true, "modelSourceName": true,
		// No metric/parameter/dataset-specific properties allowed
	},

	RestEntityDocArtifact: {
		// Common Artifact properties
		"id": true, "name": true, "externalId": true,
		"createTimeSinceEpoch": true, "lastUpdateTimeSinceEpoch": true,
		"uri": true, "state": true,
		// DocArtifact has minimal additional properties
		// No metric/parameter/dataset-specific properties allowed
	},

	RestEntityDataSet: {
		// Common Artifact properties
		"id": true, "name": true, "externalId": true,
		"createTimeSinceEpoch": true, "lastUpdateTimeSinceEpoch": true,
		"uri": true, "state": true,
		// DataSet-specific properties
		"digest": true, "sourceType": true, "source": true, "schema": true, "profile": true,
		// No metric/parameter/model-specific properties allowed
	},

	RestEntityMetric: {
		// Common Artifact properties
		"id": true, "name": true, "externalId": true,
		"createTimeSinceEpoch": true, "lastUpdateTimeSinceEpoch": true,
		"uri": true, "state": true,
		// Metric-specific properties
		"value": true, "timestamp": true, "step": true,
		// No parameter/dataset/model-specific properties allowed
	},

	RestEntityParameter: {
		// Common Artifact properties
		"id": true, "name": true, "externalId": true,
		"createTimeSinceEpoch": true, "lastUpdateTimeSinceEpoch": true,
		"uri": true, "state": true,
		// Parameter-specific properties
		"value": true, "parameterType": true,
		// No metric/dataset/model-specific properties allowed
	},

	// Execution-based entities
	RestEntityServeModel: {
		// Common Execution properties
		"id": true, "name": true, "externalId": true,
		"createTimeSinceEpoch": true, "lastUpdateTimeSinceEpoch": true,
		"lastKnownState": true,
		// ServeModel-specific properties
		"modelVersionId": true, "inferenceServiceId": true,
		"registeredModelId": true, "servingEnvironmentId": true,
	},
}

// GetMLMDEntityType maps REST entity types to their underlying MLMD entity type
func GetMLMDEntityType(restEntityType RestEntityType) EntityType {
	switch restEntityType {
	case RestEntityRegisteredModel, RestEntityModelVersion, RestEntityInferenceService,
		RestEntityServingEnvironment, RestEntityExperiment, RestEntityExperimentRun:
		return EntityTypeContext

	case RestEntityModelArtifact, RestEntityDocArtifact, RestEntityDataSet,
		RestEntityMetric, RestEntityParameter:
		return EntityTypeArtifact

	case RestEntityServeModel:
		return EntityTypeExecution

	default:
		return EntityTypeContext // Default fallback
	}
}

// GetPropertyDefinitionForRestEntity returns property definition for a REST entity type
// This function determines the correct data type and storage location for properties
func GetPropertyDefinitionForRestEntity(restEntityType RestEntityType, propertyName string) PropertyDefinition {
	// Check if this is a well-known property for this specific REST entity type
	allowedProperties, exists := RestEntityPropertyMap[restEntityType]
	if exists {
		if _, isWellKnown := allowedProperties[propertyName]; isWellKnown {
			// Use the well-known property definition
			mlmdEntityType := GetMLMDEntityType(restEntityType)
			return GetPropertyDefinition(mlmdEntityType, propertyName)
		}
	}

	// Not a well-known property for this entity type, treat as custom
	return PropertyDefinition{
		Location:  Custom,
		ValueType: StringValueType, // Default, will be inferred at runtime
		Column:    propertyName,    // Use the property name as-is for custom properties
	}
}
