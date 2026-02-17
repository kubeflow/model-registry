package models

// ModelTransferJobSourceType represents the type of source for a transfer job
type ModelTransferJobSourceType string

const (
	ModelTransferJobSourceTypeS3  ModelTransferJobSourceType = "s3"
	ModelTransferJobSourceTypeOCI ModelTransferJobSourceType = "oci"
	ModelTransferJobSourceTypeURI ModelTransferJobSourceType = "uri"
)

// ModelTransferJobDestinationType represents the type of destination for a transfer job
type ModelTransferJobDestinationType string

const (
	ModelTransferJobDestinationTypeS3  ModelTransferJobDestinationType = "s3"
	ModelTransferJobDestinationTypeOCI ModelTransferJobDestinationType = "oci"
)

// ModelTransferJobStatus represents the status of a transfer job
type ModelTransferJobStatus string

const (
	ModelTransferJobStatusPending   ModelTransferJobStatus = "PENDING"
	ModelTransferJobStatusRunning   ModelTransferJobStatus = "RUNNING"
	ModelTransferJobStatusCompleted ModelTransferJobStatus = "COMPLETED"
	ModelTransferJobStatusFailed    ModelTransferJobStatus = "FAILED"
	ModelTransferJobStatusCancelled ModelTransferJobStatus = "CANCELLED"
)

// ModelTransferJobUploadIntent represents the intent of the upload
type ModelTransferJobUploadIntent string

const (
	ModelTransferJobUploadIntentCreateModel    ModelTransferJobUploadIntent = "create_model"
	ModelTransferJobUploadIntentCreateVersion  ModelTransferJobUploadIntent = "create_version"
	ModelTransferJobUploadIntentUpdateArtifact ModelTransferJobUploadIntent = "update_artifact"
)

// ModelTransferJobSource represents the source configuration for a transfer job
type ModelTransferJobSource struct {
	Type     ModelTransferJobSourceType `json:"type"`
	Bucket   string                     `json:"bucket,omitempty"`
	Key      string                     `json:"key,omitempty"`
	Region   string                     `json:"region,omitempty"`
	Endpoint string                     `json:"endpoint,omitempty"`
	URI      string                     `json:"uri,omitempty"`
	Registry string                     `json:"registry,omitempty"`
}

// ModelTransferJobDestination represents the destination configuration for a transfer job
type ModelTransferJobDestination struct {
	Type     ModelTransferJobDestinationType `json:"type"`
	Bucket   string                          `json:"bucket,omitempty"`
	Key      string                          `json:"key,omitempty"`
	Region   string                          `json:"region,omitempty"`
	Endpoint string                          `json:"endpoint,omitempty"`
	URI      string                          `json:"uri,omitempty"`
	Registry string                          `json:"registry,omitempty"`
}

// ModelTransferJob represents a model transfer job
type ModelTransferJob struct {
	Id                       string                       `json:"id"`
	Name                     string                       `json:"name"`
	Description              string                       `json:"description,omitempty"`
	Source                   ModelTransferJobSource       `json:"source"`
	Destination              ModelTransferJobDestination  `json:"destination"`
	UploadIntent             ModelTransferJobUploadIntent `json:"uploadIntent"`
	RegisteredModelId        string                       `json:"registeredModelId,omitempty"`
	RegisteredModelName      string                       `json:"registeredModelName,omitempty"`
	ModelVersionId           string                       `json:"modelVersionId,omitempty"`
	ModelVersionName         string                       `json:"modelVersionName,omitempty"`
	ModelArtifactId          string                       `json:"modelArtifactId,omitempty"`
	ModelArtifactName        string                       `json:"modelArtifactName,omitempty"`
	Namespace                string                       `json:"namespace,omitempty"`
	Author                   string                       `json:"author,omitempty"`
	Status                   ModelTransferJobStatus       `json:"status"`
	CreateTimeSinceEpoch     string                       `json:"createTimeSinceEpoch"`
	LastUpdateTimeSinceEpoch string                       `json:"lastUpdateTimeSinceEpoch"`
	ErrorMessage             string                       `json:"errorMessage,omitempty"`
}

// ModelTransferJobList represents a list of model transfer jobs
type ModelTransferJobList struct {
	Items         []ModelTransferJob `json:"items"`
	Size          int                `json:"size"`
	PageSize      int                `json:"pageSize"`
	NextPageToken string             `json:"nextPageToken"`
}

// ModelTransferJobOperationStatus is the response body for delete and update operations.
// TODO: Remove this type when the actual implementation returns the real resource in the response.
type ModelTransferJobOperationStatus struct {
	Status string `json:"status"`
}
