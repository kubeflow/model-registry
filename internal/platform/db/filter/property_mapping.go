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
	ValueType string
	Column    string

	RelatedEntityType RelatedEntityType
	RelatedProperty   string
	JoinTable         string
}

// EntityPropertyMap maps property names to their definitions for each entity type
type EntityPropertyMap map[string]PropertyDefinition

// GetPropertyDefinition returns the property definition for a given entity type and property name
func GetPropertyDefinition(entityType EntityType, propertyName string) PropertyDefinition {
	entityMap := getEntityPropertyMap(entityType)

	if def, exists := entityMap[propertyName]; exists {
		return def
	}

	return PropertyDefinition{
		Location:  Custom,
		ValueType: StringValueType,
		Column:    propertyName,
	}
}

func getEntityPropertyMap(entityType EntityType) EntityPropertyMap {
	switch entityType {
	case EntityTypeContext:
		return contextPropertyMap
	case EntityTypeArtifact:
		return artifactPropertyMap
	case EntityTypeExecution:
		return executionPropertyMap
	default:
		return contextPropertyMap
	}
}

var contextPropertyMap = EntityPropertyMap{
	"id":                       {Location: EntityTable, ValueType: IntValueType, Column: "id"},
	"name":                     {Location: EntityTable, ValueType: StringValueType, Column: "name"},
	"externalId":               {Location: EntityTable, ValueType: StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: EntityTable, ValueType: IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: EntityTable, ValueType: IntValueType, Column: "last_update_time_since_epoch"},

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
	"id":                       {Location: EntityTable, ValueType: IntValueType, Column: "id"},
	"name":                     {Location: EntityTable, ValueType: StringValueType, Column: "name"},
	"externalId":               {Location: EntityTable, ValueType: StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: EntityTable, ValueType: IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: EntityTable, ValueType: IntValueType, Column: "last_update_time_since_epoch"},
	"uri":                      {Location: EntityTable, ValueType: StringValueType, Column: "uri"},
	"state":                    {Location: EntityTable, ValueType: IntValueType, Column: "state"},

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
	"value":              {Location: PropertyTable, ValueType: DoubleValueType, Column: "value"},
	"timestamp":          {Location: PropertyTable, ValueType: IntValueType, Column: "timestamp"},
	"step":               {Location: PropertyTable, ValueType: IntValueType, Column: "step"},
	"parameterType":      {Location: PropertyTable, ValueType: StringValueType, Column: "parameter_type"},
	"digest":             {Location: PropertyTable, ValueType: StringValueType, Column: "digest"},
	"sourceType":         {Location: PropertyTable, ValueType: StringValueType, Column: "source_type"},
	"source":             {Location: PropertyTable, ValueType: StringValueType, Column: "source"},
	"schema":             {Location: PropertyTable, ValueType: StringValueType, Column: "schema"},
	"profile":            {Location: PropertyTable, ValueType: StringValueType, Column: "profile"},
	"experimentId":       {Location: PropertyTable, ValueType: IntValueType, Column: "experiment_id"},
	"experimentRunId":    {Location: PropertyTable, ValueType: IntValueType, Column: "experiment_run_id"},
}

var executionPropertyMap = EntityPropertyMap{
	"id":                       {Location: EntityTable, ValueType: IntValueType, Column: "id"},
	"name":                     {Location: EntityTable, ValueType: StringValueType, Column: "name"},
	"externalId":               {Location: EntityTable, ValueType: StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: EntityTable, ValueType: IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: EntityTable, ValueType: IntValueType, Column: "last_update_time_since_epoch"},
	"lastKnownState":           {Location: EntityTable, ValueType: IntValueType, Column: "last_known_state"},

	"modelVersionId":       {Location: PropertyTable, ValueType: IntValueType, Column: "model_version_id"},
	"inferenceServiceId":   {Location: PropertyTable, ValueType: IntValueType, Column: "inference_service_id"},
	"registeredModelId":    {Location: PropertyTable, ValueType: IntValueType, Column: "registered_model_id"},
	"servingEnvironmentId": {Location: PropertyTable, ValueType: IntValueType, Column: "serving_environment_id"},
}
