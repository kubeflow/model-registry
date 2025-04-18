package models

// ModelRegistrySettingsPayload defines the structure for creating a ModelRegistryKind with optional credentials.
type ModelRegistrySettingsPayload struct {
	ModelRegistry            ModelRegistryKind `json:"modelRegistry"`
	DatabasePassword         *string           `json:"databasePassword,omitempty"`         // Use pointer for optional field
	NewDatabaseCACertificate *string           `json:"newDatabaseCACertificate,omitempty"` // Use pointer for optional field
}

// ModelRegistryKindDetail includes ModelRegistryKind and potentially the database password.
type ModelRegistryKindDetail struct {
	ModelRegistryKind         // Embed ModelRegistryKind
	DatabasePassword  *string `json:"databasePassword,omitempty"` // Use pointer for optional field
}
