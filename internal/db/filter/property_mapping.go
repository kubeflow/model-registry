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
	ValueType string // "int_value", "string_value", "double_value", "bool_value"
	Column    string // Database column name (for entity table properties)
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
		ValueType: "string_value", // Default to string, will be inferred at runtime
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
	"id":                       {EntityTable, "int_value", "id"},
	"name":                     {EntityTable, "string_value", "name"},
	"externalId":               {EntityTable, "string_value", "external_id"},
	"createTimeSinceEpoch":     {EntityTable, "int_value", "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {EntityTable, "int_value", "last_update_time_since_epoch"},

	// Properties that are stored in ContextProperty table but are "well-known" (not custom)
	// These are properties that the application manages, not user-defined custom properties
	"registeredModelId":    {PropertyTable, "int_value", ""},
	"modelVersionId":       {PropertyTable, "int_value", ""},
	"servingEnvironmentId": {PropertyTable, "int_value", ""},
	"experimentId":         {PropertyTable, "int_value", ""},
	"runtime":              {PropertyTable, "string_value", ""},
	"desiredState":         {PropertyTable, "string_value", ""},
	"state":                {PropertyTable, "string_value", ""},
	"owner":                {PropertyTable, "string_value", ""},
	"author":               {PropertyTable, "string_value", ""},
	"status":               {PropertyTable, "string_value", ""},
	"endTimeSinceEpoch":    {PropertyTable, "int_value", ""},
	"startTimeSinceEpoch":  {PropertyTable, "int_value", ""},
}
var artifactPropertyMap = EntityPropertyMap{
	// Entity table columns (Artifact table)
	"id":                       {EntityTable, "int_value", "id"},
	"name":                     {EntityTable, "string_value", "name"},
	"externalId":               {EntityTable, "string_value", "external_id"},
	"createTimeSinceEpoch":     {EntityTable, "int_value", "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {EntityTable, "int_value", "last_update_time_since_epoch"},
	"uri":                      {EntityTable, "string_value", "uri"},
	"state":                    {EntityTable, "int_value", "state"},

	// Properties that are stored in ArtifactProperty table but are "well-known"
	"modelFormatName":    {PropertyTable, "string_value", ""},
	"modelFormatVersion": {PropertyTable, "string_value", ""},
	"storageKey":         {PropertyTable, "string_value", ""},
	"storagePath":        {PropertyTable, "string_value", ""},
	"serviceAccountName": {PropertyTable, "string_value", ""},
	"modelSourceKind":    {PropertyTable, "string_value", ""},
	"modelSourceClass":   {PropertyTable, "string_value", ""},
	"modelSourceGroup":   {PropertyTable, "string_value", ""},
	"modelSourceId":      {PropertyTable, "string_value", ""},
	"modelSourceName":    {PropertyTable, "string_value", ""},
	"value":              {PropertyTable, "double_value", ""}, // For metrics/parameters
	"timestamp":          {PropertyTable, "int_value", ""},    // For metrics
	"step":               {PropertyTable, "int_value", ""},    // For metrics
	"parameterType":      {PropertyTable, "string_value", ""}, // For parameters
	"digest":             {PropertyTable, "string_value", ""}, // For datasets
	"sourceType":         {PropertyTable, "string_value", ""}, // For datasets
	"source":             {PropertyTable, "string_value", ""}, // For datasets
	"schema":             {PropertyTable, "string_value", ""}, // For datasets
	"profile":            {PropertyTable, "string_value", ""}, // For datasets
}

// executionPropertyMap defines properties for Execution entities
// Used by: ServeModel
var executionPropertyMap = EntityPropertyMap{
	// Entity table columns (Execution table)
	"id":                       {EntityTable, "int_value", "id"},
	"name":                     {EntityTable, "string_value", "name"},
	"externalId":               {EntityTable, "string_value", "external_id"},
	"createTimeSinceEpoch":     {EntityTable, "int_value", "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {EntityTable, "int_value", "last_update_time_since_epoch"},
	"lastKnownState":           {EntityTable, "int_value", "last_known_state"},

	// Properties that are stored in ExecutionProperty table but are "well-known"
	"modelVersionId":       {PropertyTable, "int_value", ""},
	"inferenceServiceId":   {PropertyTable, "int_value", ""},
	"registeredModelId":    {PropertyTable, "int_value", ""},
	"servingEnvironmentId": {PropertyTable, "int_value", ""},
}
