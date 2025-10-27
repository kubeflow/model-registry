package filter

// PropertyLocation indicates where a property is stored
type PropertyLocation int

const (
	EntityTable   PropertyLocation = iota // Property is a column in the main entity table
	PropertyTable                         // Property is stored in the entity's property table
	Custom                                // Property is a custom property in the property table
	RelatedEntity                         // Property is in a related entity (requires JOIN)
)

// RelatedEntityType indicates the type of related entity
type RelatedEntityType string

const (
	RelatedEntityArtifact  RelatedEntityType = "artifact"
	RelatedEntityContext   RelatedEntityType = "context"
	RelatedEntityExecution RelatedEntityType = "execution"
)

// PropertyDefinition defines how a property should be handled
type PropertyDefinition struct {
	Location  PropertyLocation
	ValueType string // IntValueType, StringValueType, DoubleValueType, BoolValueType
	Column    string // Database column name (for entity table) or property name (for property table)

	// Fields for related entity properties
	RelatedEntityType RelatedEntityType // Type of related entity (artifact, context, execution)
	RelatedProperty   string            // Property name in the related entity
	JoinTable         string            // Table to join through (e.g., "Attribution", "ParentContext")
}

// EntityPropertyMap maps property names to their definitions for each entity type
type EntityPropertyMap map[string]PropertyDefinition

// GetPropertyDefinition returns the property definition for a given entity type and property name
func GetPropertyDefinition(entityType EntityType, propertyName string) PropertyDefinition {
	entityMap := getEntityPropertyMap(entityType)

	if def, exists := entityMap[propertyName]; exists {
		return def
	}

	// If not found in the map, assume it's a custom property
	return PropertyDefinition{
		Location:  Custom,
		ValueType: StringValueType, // Default to string, will be inferred at runtime
		Column:    propertyName,    // Use the property name as-is for custom properties
	}
}

// getEntityPropertyMap returns the property mapping for a specific entity type
func getEntityPropertyMap(entityType EntityType) EntityPropertyMap {
	switch entityType {
	case EntityTypeContext:
		return contextPropertyMap
	case EntityTypeArtifact:
		return artifactPropertyMap
	case EntityTypeExecution:
		return executionPropertyMap
	default:
		return contextPropertyMap // Default fallback
	}
}

// contextPropertyMap defines properties for Context entities
// Used by: RegisteredModel, ModelVersion, InferenceService, ServingEnvironment, Experiment, ExperimentRun
var contextPropertyMap = EntityPropertyMap{
	// Entity table columns (Context table)
	"id":                       {Location: EntityTable, ValueType: IntValueType, Column: "id"},
	"name":                     {Location: EntityTable, ValueType: StringValueType, Column: "name"},
	"externalId":               {Location: EntityTable, ValueType: StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: EntityTable, ValueType: IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: EntityTable, ValueType: IntValueType, Column: "last_update_time_since_epoch"},

	// Properties that are stored in ContextProperty table but are "well-known" (not custom)
	// These are properties that the application manages, not user-defined custom properties
	"registeredModelId":    {Location: PropertyTable, ValueType: IntValueType, Column: "registered_model_id"},
	"modelVersionId":       {Location: PropertyTable, ValueType: IntValueType, Column: "model_version_id"},
	"servingEnvironmentId": {Location: PropertyTable, ValueType: IntValueType, Column: "serving_environment_id"},
	"experimentId":         {Location: PropertyTable, ValueType: IntValueType, Column: "experiment_id"},
	"runtime":              {Location: PropertyTable, ValueType: StringValueType, Column: "runtime"},
	"desiredState":         {Location: PropertyTable, ValueType: StringValueType, Column: "desired_state"},
	"state":                {Location: PropertyTable, ValueType: StringValueType, Column: "state"},
	"owner":                {Location: PropertyTable, ValueType: StringValueType, Column: "owner"},
	"author":               {Location: PropertyTable, ValueType: StringValueType, Column: "author"},
	"status":               {Location: PropertyTable, ValueType: StringValueType, Column: "status"},
	"endTimeSinceEpoch":    {Location: PropertyTable, ValueType: StringValueType, Column: "end_time_since_epoch"},
	"startTimeSinceEpoch":  {Location: PropertyTable, ValueType: StringValueType, Column: "start_time_since_epoch"},
}
var artifactPropertyMap = EntityPropertyMap{
	// Entity table columns (Artifact table)
	"id":                       {Location: EntityTable, ValueType: IntValueType, Column: "id"},
	"name":                     {Location: EntityTable, ValueType: StringValueType, Column: "name"},
	"externalId":               {Location: EntityTable, ValueType: StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: EntityTable, ValueType: IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: EntityTable, ValueType: IntValueType, Column: "last_update_time_since_epoch"},
	"uri":                      {Location: EntityTable, ValueType: StringValueType, Column: "uri"},
	"state":                    {Location: EntityTable, ValueType: IntValueType, Column: "state"},

	// Properties that are stored in ArtifactProperty table but are "well-known"
	"modelFormatName":    {Location: PropertyTable, ValueType: StringValueType, Column: "model_format_name"},
	"modelFormatVersion": {Location: PropertyTable, ValueType: StringValueType, Column: "model_format_version"},
	"storageKey":         {Location: PropertyTable, ValueType: StringValueType, Column: "storage_key"},
	"storagePath":        {Location: PropertyTable, ValueType: StringValueType, Column: "storage_path"},
	"serviceAccountName": {Location: PropertyTable, ValueType: StringValueType, Column: "service_account_name"},
	"modelSourceKind":    {Location: PropertyTable, ValueType: StringValueType, Column: "model_source_kind"},
	"modelSourceClass":   {Location: PropertyTable, ValueType: StringValueType, Column: "model_source_class"},
	"modelSourceGroup":   {Location: PropertyTable, ValueType: StringValueType, Column: "model_source_group"},
	"modelSourceId":      {Location: PropertyTable, ValueType: StringValueType, Column: "model_source_id"},
	"modelSourceName":    {Location: PropertyTable, ValueType: StringValueType, Column: "model_source_name"},
	"value":              {Location: PropertyTable, ValueType: DoubleValueType, Column: "value"},          // For metrics/parameters
	"timestamp":          {Location: PropertyTable, ValueType: IntValueType, Column: "timestamp"},         // For metrics
	"step":               {Location: PropertyTable, ValueType: IntValueType, Column: "step"},              // For metrics
	"parameterType":      {Location: PropertyTable, ValueType: StringValueType, Column: "parameter_type"}, // For parameters
	"digest":             {Location: PropertyTable, ValueType: StringValueType, Column: "digest"},         // For datasets
	"sourceType":         {Location: PropertyTable, ValueType: StringValueType, Column: "source_type"},    // For datasets
	"source":             {Location: PropertyTable, ValueType: StringValueType, Column: "source"},         // For datasets
	"schema":             {Location: PropertyTable, ValueType: StringValueType, Column: "schema"},         // For datasets
	"profile":            {Location: PropertyTable, ValueType: StringValueType, Column: "profile"},        // For datasets
	"experimentId":       {Location: PropertyTable, ValueType: IntValueType, Column: "experiment_id"},     // For all artifacts
	"experimentRunId":    {Location: PropertyTable, ValueType: IntValueType, Column: "experiment_run_id"}, // For all artifacts
}

// executionPropertyMap defines properties for Execution entities
// Used by: ServeModel
var executionPropertyMap = EntityPropertyMap{
	// Entity table columns (Execution table)
	"id":                       {Location: EntityTable, ValueType: IntValueType, Column: "id"},
	"name":                     {Location: EntityTable, ValueType: StringValueType, Column: "name"},
	"externalId":               {Location: EntityTable, ValueType: StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: EntityTable, ValueType: IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: EntityTable, ValueType: IntValueType, Column: "last_update_time_since_epoch"},
	"lastKnownState":           {Location: EntityTable, ValueType: IntValueType, Column: "last_known_state"},

	// Properties that are stored in ExecutionProperty table but are "well-known"
	"modelVersionId":       {Location: PropertyTable, ValueType: IntValueType, Column: "model_version_id"},
	"inferenceServiceId":   {Location: PropertyTable, ValueType: IntValueType, Column: "inference_service_id"},
	"registeredModelId":    {Location: PropertyTable, ValueType: IntValueType, Column: "registered_model_id"},
	"servingEnvironmentId": {Location: PropertyTable, ValueType: IntValueType, Column: "serving_environment_id"},
}
