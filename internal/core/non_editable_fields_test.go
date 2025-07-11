package core_test

import (
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/core"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNonEditableFieldsProtection verifies that non-editable fields are protected for all entity types
func TestNonEditableFieldsProtection(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("RegisteredModel non-editable fields protection", func(t *testing.T) {
		testRegisteredModelNonEditableFields(t, service)
	})

	t.Run("ModelVersion non-editable fields protection", func(t *testing.T) {
		testModelVersionNonEditableFields(t, service)
	})

	t.Run("ServingEnvironment non-editable fields protection", func(t *testing.T) {
		testServingEnvironmentNonEditableFields(t, service)
	})

	t.Run("InferenceService non-editable fields protection", func(t *testing.T) {
		testInferenceServiceNonEditableFields(t, service)
	})

	t.Run("ServeModel non-editable fields protection", func(t *testing.T) {
		testServeModelNonEditableFields(t, service)
	})

	t.Run("ExperimentRun non-editable fields protection", func(t *testing.T) {
		testExperimentRunNonEditableFields(t, service)
	})

	t.Run("Artifact non-editable fields protection", func(t *testing.T) {
		testArtifactNonEditableFields(t, service)
	})
}

func testRegisteredModelNonEditableFields(t *testing.T, service *core.ModelRegistryService) {
	// Non-editable: Id, CreateTimeSinceEpoch, LastUpdateTimeSinceEpoch, Name
	// Editable: Description, ExternalId, CustomProperties, State, Owner, etc.

	// Create initial model
	model := &openapi.RegisteredModel{
		Name:        "test-rm-non-editable",
		Description: apiutils.Of("Original description"),
		Owner:       apiutils.Of("original-owner"),
		CustomProperties: &map[string]openapi.MetadataValue{
			"original": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue:  "value",
					MetadataType: "MetadataStringValue",
				},
			},
		},
	}

	created, err := service.UpsertRegisteredModel(model)
	require.NoError(t, err)

	// Store originals
	originalId := *created.Id
	originalName := created.Name
	originalCreateTime := *created.CreateTimeSinceEpoch

	time.Sleep(10 * time.Millisecond)

	// Attempt to hack non-editable fields
	updateRequest := &openapi.RegisteredModel{
		Id:                       created.Id,
		Name:                     "HACKED_NAME",                      // Should be ignored
		CreateTimeSinceEpoch:     apiutils.Of("9999999999"),          // Should be ignored
		LastUpdateTimeSinceEpoch: apiutils.Of("8888888888"),          // Should be ignored
		Description:              apiutils.Of("Updated description"), // Should work
		Owner:                    apiutils.Of("updated-owner"),       // Should work
	}

	updated, err := service.UpsertRegisteredModel(updateRequest)
	require.NoError(t, err)

	// Verify protection
	assert.Equal(t, originalId, *updated.Id, "RegisteredModel ID should not be changeable")
	assert.Equal(t, originalName, updated.Name, "RegisteredModel Name should not be changeable")
	assert.Equal(t, originalCreateTime, *updated.CreateTimeSinceEpoch, "RegisteredModel CreateTime should not be changeable")
	assert.Equal(t, "Updated description", *updated.Description, "RegisteredModel Description should be updatable")
	assert.Equal(t, "updated-owner", *updated.Owner, "RegisteredModel Owner should be updatable")
}

func testModelVersionNonEditableFields(t *testing.T, service *core.ModelRegistryService) {
	// Non-editable: Id, CreateTimeSinceEpoch, LastUpdateTimeSinceEpoch, Name, RegisteredModelId
	// Editable: Description, ExternalId, CustomProperties, State, Author

	// Create parent registered model first
	rm := &openapi.RegisteredModel{Name: "test-mv-parent"}
	createdRM, err := service.UpsertRegisteredModel(rm)
	require.NoError(t, err)

	// Create model version
	mv := &openapi.ModelVersion{
		Name:        "v1.0",
		Description: apiutils.Of("Original version"),
		Author:      apiutils.Of("original-author"),
	}

	created, err := service.UpsertModelVersion(mv, createdRM.Id)
	require.NoError(t, err)

	// Store originals
	originalId := *created.Id
	originalName := created.Name
	originalRegisteredModelId := created.RegisteredModelId
	originalCreateTime := *created.CreateTimeSinceEpoch

	time.Sleep(10 * time.Millisecond)

	// Attempt to hack non-editable fields
	updateRequest := &openapi.ModelVersion{
		Id:                       created.Id,
		Name:                     "HACKED_NAME",                      // Should be ignored
		RegisteredModelId:        "999999",                           // Should be ignored
		CreateTimeSinceEpoch:     apiutils.Of("9999999999"),          // Should be ignored
		LastUpdateTimeSinceEpoch: apiutils.Of("8888888888"),          // Should be ignored
		Description:              apiutils.Of("Updated description"), // Should work
		Author:                   apiutils.Of("updated-author"),      // Should work
	}

	updated, err := service.UpsertModelVersion(updateRequest, nil)
	require.NoError(t, err)

	// Verify protection
	assert.Equal(t, originalId, *updated.Id, "ModelVersion ID should not be changeable")
	assert.Equal(t, originalName, updated.Name, "ModelVersion Name should not be changeable")
	assert.Equal(t, originalRegisteredModelId, updated.RegisteredModelId, "ModelVersion RegisteredModelId should not be changeable")
	assert.Equal(t, originalCreateTime, *updated.CreateTimeSinceEpoch, "ModelVersion CreateTime should not be changeable")
	assert.Equal(t, "Updated description", *updated.Description, "ModelVersion Description should be updatable")
	assert.Equal(t, "updated-author", *updated.Author, "ModelVersion Author should be updatable")
}

func testServingEnvironmentNonEditableFields(t *testing.T, service *core.ModelRegistryService) {
	// Non-editable: Id, CreateTimeSinceEpoch, LastUpdateTimeSinceEpoch, Name
	// Editable: Description, ExternalId, CustomProperties

	se := &openapi.ServingEnvironment{
		Name:        "test-se-non-editable",
		Description: apiutils.Of("Original description"),
	}

	created, err := service.UpsertServingEnvironment(se)
	require.NoError(t, err)

	// Store originals
	originalId := *created.Id
	originalName := created.Name
	originalCreateTime := *created.CreateTimeSinceEpoch

	time.Sleep(10 * time.Millisecond)

	// Attempt to hack non-editable fields
	updateRequest := &openapi.ServingEnvironment{
		Id:                       created.Id,
		Name:                     "HACKED_NAME",                      // Should be ignored
		CreateTimeSinceEpoch:     apiutils.Of("9999999999"),          // Should be ignored
		LastUpdateTimeSinceEpoch: apiutils.Of("8888888888"),          // Should be ignored
		Description:              apiutils.Of("Updated description"), // Should work
	}

	updated, err := service.UpsertServingEnvironment(updateRequest)
	require.NoError(t, err)

	// Verify protection
	assert.Equal(t, originalId, *updated.Id, "ServingEnvironment ID should not be changeable")
	assert.Equal(t, originalName, updated.Name, "ServingEnvironment Name should not be changeable")
	assert.Equal(t, originalCreateTime, *updated.CreateTimeSinceEpoch, "ServingEnvironment CreateTime should not be changeable")
	assert.Equal(t, "Updated description", *updated.Description, "ServingEnvironment Description should be updatable")
}

func testInferenceServiceNonEditableFields(t *testing.T, service *core.ModelRegistryService) {
	// Non-editable: Id, CreateTimeSinceEpoch, LastUpdateTimeSinceEpoch, Name, RegisteredModelId, ServingEnvironmentId
	// Editable: Description, ExternalId, CustomProperties, ModelVersionId, Runtime, DesiredState

	// Create dependencies
	rm := &openapi.RegisteredModel{Name: "test-is-rm"}
	createdRM, err := service.UpsertRegisteredModel(rm)
	require.NoError(t, err)

	mv := &openapi.ModelVersion{Name: "v1"}
	createdMV, err := service.UpsertModelVersion(mv, createdRM.Id)
	require.NoError(t, err)

	se := &openapi.ServingEnvironment{Name: "test-is-se"}
	createdSE, err := service.UpsertServingEnvironment(se)
	require.NoError(t, err)

	// Create inference service
	is := &openapi.InferenceService{
		Name:                 apiutils.Of("test-is"),
		ServingEnvironmentId: *createdSE.Id,
		RegisteredModelId:    *createdRM.Id,
		ModelVersionId:       createdMV.Id,
		Description:          apiutils.Of("Original description"),
		Runtime:              apiutils.Of("original-runtime"),
	}

	created, err := service.UpsertInferenceService(is)
	require.NoError(t, err)

	// Store originals
	originalId := *created.Id
	originalName := *created.Name
	originalServingEnvId := created.ServingEnvironmentId
	originalRegModelId := created.RegisteredModelId

	time.Sleep(10 * time.Millisecond)

	// Attempt to hack non-editable fields
	updateRequest := &openapi.InferenceService{
		Id:                   created.Id,
		Name:                 apiutils.Of("HACKED_NAME"),         // Should be ignored
		ServingEnvironmentId: "999999",                           // Should be ignored
		RegisteredModelId:    "888888",                           // Should be ignored
		ModelVersionId:       apiutils.Of("777777"),              // Should work (editable)
		Description:          apiutils.Of("Updated description"), // Should work
		Runtime:              apiutils.Of("updated-runtime"),     // Should work
	}

	updated, err := service.UpsertInferenceService(updateRequest)
	require.NoError(t, err)

	// Verify protection
	assert.Equal(t, originalId, *updated.Id, "InferenceService ID should not be changeable")
	assert.Equal(t, originalName, *updated.Name, "InferenceService Name should not be changeable")
	assert.Equal(t, originalServingEnvId, updated.ServingEnvironmentId, "InferenceService ServingEnvironmentId should not be changeable")
	assert.Equal(t, originalRegModelId, updated.RegisteredModelId, "InferenceService RegisteredModelId should not be changeable")
	assert.Equal(t, "Updated description", *updated.Description, "InferenceService Description should be updatable")
	assert.Equal(t, "updated-runtime", *updated.Runtime, "InferenceService Runtime should be updatable")
}

func testServeModelNonEditableFields(t *testing.T, service *core.ModelRegistryService) {
	// Non-editable: Id, CreateTimeSinceEpoch, LastUpdateTimeSinceEpoch, Name, ModelVersionId
	// Editable: Description, ExternalId, CustomProperties, LastKnownState

	// Create dependencies
	rm := &openapi.RegisteredModel{Name: "test-sm-rm"}
	createdRM, err := service.UpsertRegisteredModel(rm)
	require.NoError(t, err)

	mv := &openapi.ModelVersion{Name: "v1"}
	createdMV, err := service.UpsertModelVersion(mv, createdRM.Id)
	require.NoError(t, err)

	se := &openapi.ServingEnvironment{Name: "test-sm-se"}
	createdSE, err := service.UpsertServingEnvironment(se)
	require.NoError(t, err)

	is := &openapi.InferenceService{
		Name:                 apiutils.Of("test-sm-is"),
		ServingEnvironmentId: *createdSE.Id,
		RegisteredModelId:    *createdRM.Id,
		ModelVersionId:       createdMV.Id,
	}
	createdIS, err := service.UpsertInferenceService(is)
	require.NoError(t, err)

	// Create serve model
	sm := &openapi.ServeModel{
		Name:           apiutils.Of("test-sm"),
		ModelVersionId: *createdMV.Id,
		Description:    apiutils.Of("Original description"),
	}

	created, err := service.UpsertServeModel(sm, createdIS.Id)
	require.NoError(t, err)

	// Store originals
	originalId := *created.Id
	originalName := *created.Name
	originalModelVersionId := created.ModelVersionId

	time.Sleep(10 * time.Millisecond)

	// Attempt to hack non-editable fields
	updateRequest := &openapi.ServeModel{
		Id:             created.Id,
		Name:           apiutils.Of("HACKED_NAME"),         // Should be ignored
		ModelVersionId: "999999",                           // Should be ignored
		Description:    apiutils.Of("Updated description"), // Should work
	}

	updated, err := service.UpsertServeModel(updateRequest, nil)
	require.NoError(t, err)

	// Verify protection
	assert.Equal(t, originalId, *updated.Id, "ServeModel ID should not be changeable")
	assert.Equal(t, originalName, *updated.Name, "ServeModel Name should not be changeable")
	assert.Equal(t, originalModelVersionId, updated.ModelVersionId, "ServeModel ModelVersionId should not be changeable")
	assert.Equal(t, "Updated description", *updated.Description, "ServeModel Description should be updatable")
}

func testExperimentRunNonEditableFields(t *testing.T, service *core.ModelRegistryService) {
	// Non-editable: Id, CreateTimeSinceEpoch, LastUpdateTimeSinceEpoch, Name, ExperimentId
	// Editable: Description, ExternalId, CustomProperties, State, Owner, Status, EndTimeSinceEpoch

	// Create parent experiment
	exp := &openapi.Experiment{Name: "test-er-exp"}
	createdExp, err := service.UpsertExperiment(exp)
	require.NoError(t, err)

	// Create experiment run
	er := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-er"),
		Description: apiutils.Of("Original description"),
		Owner:       apiutils.Of("original-owner"),
	}

	created, err := service.UpsertExperimentRun(er, createdExp.Id)
	require.NoError(t, err)

	// Store originals
	originalId := *created.Id
	originalName := *created.Name
	originalExperimentId := created.ExperimentId

	time.Sleep(10 * time.Millisecond)

	// Attempt to hack non-editable fields
	updateRequest := &openapi.ExperimentRun{
		Id:           created.Id,
		Name:         apiutils.Of("HACKED_NAME"),         // Should be ignored
		ExperimentId: "999999",                           // Should be ignored
		Description:  apiutils.Of("Updated description"), // Should work
		Owner:        apiutils.Of("updated-owner"),       // Should work
	}

	updated, err := service.UpsertExperimentRun(updateRequest, createdExp.Id)
	require.NoError(t, err)

	// Verify protection
	assert.Equal(t, originalId, *updated.Id, "ExperimentRun ID should not be changeable")
	assert.Equal(t, originalName, *updated.Name, "ExperimentRun Name should not be changeable")
	assert.Equal(t, originalExperimentId, updated.ExperimentId, "ExperimentRun ExperimentId should not be changeable")
	assert.Equal(t, "Updated description", *updated.Description, "ExperimentRun Description should be updatable")
	assert.Equal(t, "updated-owner", *updated.Owner, "ExperimentRun Owner should be updatable")
}

func testArtifactNonEditableFields(t *testing.T, service *core.ModelRegistryService) {
	// Test ModelArtifact, DocArtifact, DataSet, Metric, Parameter
	// Non-editable for all: Id, CreateTimeSinceEpoch, LastUpdateTimeSinceEpoch, Name, ArtifactType
	// Editable: Description, ExternalId, CustomProperties, Uri, State, etc.

	t.Run("ModelArtifact non-editable fields", func(t *testing.T) {
		// Create model artifact
		artifact := &openapi.Artifact{
			ModelArtifact: &openapi.ModelArtifact{
				Name:        apiutils.Of("test-ma"),
				Description: apiutils.Of("Original description"),
				Uri:         apiutils.Of("s3://original/path"),
			},
		}

		created, err := service.UpsertArtifact(artifact)
		require.NoError(t, err)

		// Store originals
		originalId := *created.ModelArtifact.Id
		originalName := *created.ModelArtifact.Name
		originalArtifactType := *created.ModelArtifact.ArtifactType

		time.Sleep(10 * time.Millisecond)

		// Attempt to hack non-editable fields
		updateRequest := &openapi.Artifact{
			ModelArtifact: &openapi.ModelArtifact{
				Id:           created.ModelArtifact.Id,
				Name:         apiutils.Of("HACKED_NAME"),         // Should be ignored
				ArtifactType: apiutils.Of("HACKED_TYPE"),         // Should be ignored
				Description:  apiutils.Of("Updated description"), // Should work
				Uri:          apiutils.Of("s3://updated/path"),   // Should work
			},
		}

		updated, err := service.UpsertArtifact(updateRequest)
		require.NoError(t, err)

		// Verify protection
		assert.Equal(t, originalId, *updated.ModelArtifact.Id, "ModelArtifact ID should not be changeable")
		assert.Equal(t, originalName, *updated.ModelArtifact.Name, "ModelArtifact Name should not be changeable")
		assert.Equal(t, originalArtifactType, *updated.ModelArtifact.ArtifactType, "ModelArtifact ArtifactType should not be changeable")
		assert.Equal(t, "Updated description", *updated.ModelArtifact.Description, "ModelArtifact Description should be updatable")
		assert.Equal(t, "s3://updated/path", *updated.ModelArtifact.Uri, "ModelArtifact Uri should be updatable")
	})

	// Similar tests can be added for DocArtifact, DataSet, Metric, Parameter
}
