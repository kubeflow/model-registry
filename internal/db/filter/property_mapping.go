package filter

// PropertyLocation indicates where a property is stored
type PropertyLocation int

const (
	EntityTable   PropertyLocation = iota // Property is a column in the main entity table
	PropertyTable                         // Property is stored in the entity's property table
	Custom                                // Property is a custom property in the property table
)

// PropertyDefinition defines how a property should be handled
type PropertyDefinition struct {
	Location  PropertyLocation
	ValueType string // IntValueType, StringValueType, DoubleValueType, BoolValueType
	Column    string // Database column name (for entity table) or property name (for property table)
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
	"id":                       {EntityTable, IntValueType, "id"},
	"name":                     {EntityTable, StringValueType, "name"},
	"externalId":               {EntityTable, StringValueType, "external_id"},
	"createTimeSinceEpoch":     {EntityTable, IntValueType, "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {EntityTable, IntValueType, "last_update_time_since_epoch"},

	// Properties that are stored in ContextProperty table but are "well-known" (not custom)
	// These are properties that the application manages, not user-defined custom properties
	"registeredModelId":    {PropertyTable, IntValueType, "registered_model_id"},
	"modelVersionId":       {PropertyTable, IntValueType, "model_version_id"},
	"servingEnvironmentId": {PropertyTable, IntValueType, "serving_environment_id"},
	"experimentId":         {PropertyTable, IntValueType, "experiment_id"},
	"runtime":              {PropertyTable, StringValueType, "runtime"},
	"desiredState":         {PropertyTable, StringValueType, "desired_state"},
	"state":                {PropertyTable, StringValueType, "state"},
	"owner":                {PropertyTable, StringValueType, "owner"},
	"author":               {PropertyTable, StringValueType, "author"},
	"status":               {PropertyTable, StringValueType, "status"},
	"endTimeSinceEpoch":    {PropertyTable, StringValueType, "end_time_since_epoch"},
	"startTimeSinceEpoch":  {PropertyTable, StringValueType, "start_time_since_epoch"},
}
var artifactPropertyMap = EntityPropertyMap{
	// Entity table columns (Artifact table)
	"id":                       {EntityTable, IntValueType, "id"},
	"name":                     {EntityTable, StringValueType, "name"},
	"externalId":               {EntityTable, StringValueType, "external_id"},
	"createTimeSinceEpoch":     {EntityTable, IntValueType, "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {EntityTable, IntValueType, "last_update_time_since_epoch"},
	"uri":                      {EntityTable, StringValueType, "uri"},
	"state":                    {EntityTable, IntValueType, "state"},

	// Properties that are stored in ArtifactProperty table but are "well-known"
	"modelFormatName":    {PropertyTable, StringValueType, "model_format_name"},
	"modelFormatVersion": {PropertyTable, StringValueType, "model_format_version"},
	"storageKey":         {PropertyTable, StringValueType, "storage_key"},
	"storagePath":        {PropertyTable, StringValueType, "storage_path"},
	"serviceAccountName": {PropertyTable, StringValueType, "service_account_name"},
	"modelSourceKind":    {PropertyTable, StringValueType, "model_source_kind"},
	"modelSourceClass":   {PropertyTable, StringValueType, "model_source_class"},
	"modelSourceGroup":   {PropertyTable, StringValueType, "model_source_group"},
	"modelSourceId":      {PropertyTable, StringValueType, "model_source_id"},
	"modelSourceName":    {PropertyTable, StringValueType, "model_source_name"},
	"value":              {PropertyTable, DoubleValueType, "value"},          // For metrics/parameters
	"timestamp":          {PropertyTable, IntValueType, "timestamp"},         // For metrics
	"step":               {PropertyTable, IntValueType, "step"},              // For metrics
	"parameterType":      {PropertyTable, StringValueType, "parameter_type"}, // For parameters
	"digest":             {PropertyTable, StringValueType, "digest"},         // For datasets
	"sourceType":         {PropertyTable, StringValueType, "source_type"},    // For datasets
	"source":             {PropertyTable, StringValueType, "source"},         // For datasets
	"schema":             {PropertyTable, StringValueType, "schema"},         // For datasets
	"profile":            {PropertyTable, StringValueType, "profile"},        // For datasets
}

// executionPropertyMap defines properties for Execution entities
// Used by: ServeModel
var executionPropertyMap = EntityPropertyMap{
	// Entity table columns (Execution table)
	"id":                       {EntityTable, IntValueType, "id"},
	"name":                     {EntityTable, StringValueType, "name"},
	"externalId":               {EntityTable, StringValueType, "external_id"},
	"createTimeSinceEpoch":     {EntityTable, IntValueType, "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {EntityTable, IntValueType, "last_update_time_since_epoch"},
	"lastKnownState":           {EntityTable, IntValueType, "last_known_state"},

	// Properties that are stored in ExecutionProperty table but are "well-known"
	"modelVersionId":       {PropertyTable, IntValueType, "model_version_id"},
	"inferenceServiceId":   {PropertyTable, IntValueType, "inference_service_id"},
	"registeredModelId":    {PropertyTable, IntValueType, "registered_model_id"},
	"servingEnvironmentId": {PropertyTable, IntValueType, "serving_environment_id"},
}
