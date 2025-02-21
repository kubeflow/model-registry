package core

import (
	"context"
	"fmt"

	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// SERVE MODEL

func (suite *CoreTestSuite) TestCreateServeModel() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        modelVersionName,
		ExternalId:  &versionExternalId,
		Description: &modelVersionDescription,
		Author:      &author,
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %s", registeredModelId)
	createdVersionId := *createdVersion.Id
	createdVersionIdAsInt, _ := converter.StringToInt64(&createdVersionId)
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalId:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: createdVersionId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)
	suite.NotNil(createdEntity.Id, "created id should not be nil")

	state, _ := openapi.NewExecutionStateFromValue(executionState)
	suite.Equal(entityName, *createdEntity.Name)
	suite.Equal(*state, *createdEntity.LastKnownState)
	suite.Equal(createdVersionId, createdEntity.ModelVersionId)
	suite.Equal(entityDescription, *createdEntity.Description)
	suite.Equal(customString, (*createdEntity.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)
	getById, err := suite.mlmdClient.GetExecutionsByID(context.Background(), &proto.GetExecutionsByIDRequest{
		ExecutionIds: []int64{*createdEntityId},
	})
	suite.Nilf(err, "error getting Execution by id %d", createdEntityId)

	suite.Equal(*createdEntityId, *getById.Executions[0].Id)
	suite.Equal(fmt.Sprintf("%s:%s", inferenceServiceId, *createdEntity.Name), *getById.Executions[0].Name)
	suite.Equal(string(*createdEntity.LastKnownState), getById.Executions[0].LastKnownState.String())
	suite.Equal(*createdVersionIdAsInt, getById.Executions[0].Properties["model_version_id"].GetIntValue())
	suite.Equal(*createdEntity.Description, getById.Executions[0].Properties["description"].GetStringValue())
	suite.Equal((*createdEntity.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Executions[0].CustomProperties["custom_string_prop"].GetStringValue())

	inferenceServiceIdAsInt, _ := converter.StringToInt64(&inferenceServiceId)
	byCtx, _ := suite.mlmdClient.GetExecutionsByContext(context.Background(), &proto.GetExecutionsByContextRequest{
		ContextId: (*int64)(inferenceServiceIdAsInt),
	})
	suite.Equal(1, len(byCtx.Executions))
	suite.Equal(*createdEntityId, *byCtx.Executions[0].Id)
}

func (suite *CoreTestSuite) TestCreateDuplicateServeModelFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        modelVersionName,
		ExternalId:  &versionExternalId,
		Description: &modelVersionDescription,
		Author:      &author,
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %s", registeredModelId)
	createdVersionId := *createdVersion.Id
	//createdVersionIdAsInt, _ := converter.StringToInt64(&createdVersionId)
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalId:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: createdVersionId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)
	suite.NotNil(createdEntity.Id, "created id should not be nil")

	// attempt to create dupliate serve model
	_, err = service.UpsertServeModel(eut, &inferenceServiceId)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "cannot register a duplicate serve model")
	suite.Equal(409, statusResp, "duplicate serve models not allowed")
}

func (suite *CoreTestSuite) TestCreateServeModelFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)
	// end of data preparation

	_, err := service.UpsertServeModel(nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid serve model pointer, can't upsert nil: bad request", err.Error())

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalId:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: "9998",
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	_, err = service.UpsertServeModel(eut, nil)
	suite.NotNil(err)
	suite.Equal("missing inferenceServiceId, cannot create ServeModel without parent resource InferenceService: bad request", err.Error())

	_, err = service.UpsertServeModel(eut, &inferenceServiceId)
	suite.NotNil(err)
	suite.Equal("no model version found for id 9998: not found", err.Error())
}

func (suite *CoreTestSuite) TestUpdateServeModel() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        modelVersionName,
		ExternalId:  &versionExternalId,
		Description: &modelVersionDescription,
		Author:      &author,
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersionId := *createdVersion.Id
	createdVersionIdAsInt, _ := converter.StringToInt64(&createdVersionId)
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalId:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: createdVersionId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)

	newState := "UNKNOWN"
	createdEntity.LastKnownState = (*openapi.ExecutionState)(&newState)
	updatedEntity, err := service.UpsertServeModel(createdEntity, &inferenceServiceId)
	suite.Nilf(err, "error updating entity for %d: %v", inferenceServiceId, err)

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)
	updatedEntityId, _ := converter.StringToInt64(updatedEntity.Id)
	suite.Equal(createdEntityId, updatedEntityId)

	getById, err := suite.mlmdClient.GetExecutionsByID(context.Background(), &proto.GetExecutionsByIDRequest{
		ExecutionIds: []int64{*createdEntityId},
	})
	suite.Nilf(err, "error getting by id %d", createdEntityId)

	suite.Equal(*createdEntityId, *getById.Executions[0].Id)
	suite.Equal(fmt.Sprintf("%s:%s", inferenceServiceId, *createdEntity.Name), *getById.Executions[0].Name)
	suite.Equal(string(newState), getById.Executions[0].LastKnownState.String())
	suite.Equal(*createdVersionIdAsInt, getById.Executions[0].Properties["model_version_id"].GetIntValue())
	suite.Equal((*createdEntity.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Executions[0].CustomProperties["custom_string_prop"].GetStringValue())

	prevModelVersionId := updatedEntity.ModelVersionId
	updatedEntity.ModelVersionId = ""
	updatedEntity, err = service.UpsertServeModel(updatedEntity, &inferenceServiceId)
	suite.Nilf(err, "error updating entity for %d: %v", inferenceServiceId, err)
	suite.Equal(prevModelVersionId, updatedEntity.ModelVersionId)
}

func (suite *CoreTestSuite) TestUpdateServeModelFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        modelVersionName,
		ExternalId:  &versionExternalId,
		Description: &modelVersionDescription,
		Author:      &author,
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersionId := *createdVersion.Id
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalId:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: createdVersionId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %d", inferenceServiceId)
	suite.NotNil(createdEntity.Id, "created id should not be nil")

	newState := "UNKNOWN"
	createdEntity.LastKnownState = (*openapi.ExecutionState)(&newState)
	updatedEntity, err := service.UpsertServeModel(createdEntity, &inferenceServiceId)
	suite.Nilf(err, "error updating entity for %s: %v", inferenceServiceId, err)

	wrongId := "9998"
	updatedEntity.Id = &wrongId
	_, err = service.UpsertServeModel(updatedEntity, &inferenceServiceId)
	suite.NotNil(err)
	suite.Equal(fmt.Sprintf("no ServeModel found for id %s: not found", wrongId), err.Error())
}

func (suite *CoreTestSuite) TestGetServeModelById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:        modelVersionName,
		ExternalId:  &versionExternalId,
		Description: &modelVersionDescription,
		Author:      &author,
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %s", registeredModelId)
	createdVersionId := *createdVersion.Id
	// end of data preparation

	eut := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		ExternalId:     &entityExternalId2,
		Description:    &entityDescription,
		Name:           &entityName,
		ModelVersionId: createdVersionId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity, err := service.UpsertServeModel(eut, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %s", inferenceServiceId)

	getById, err := service.GetServeModelById(*createdEntity.Id)
	suite.Nilf(err, "error getting entity by id %s", *createdEntity.Id)

	state, _ := openapi.NewExecutionStateFromValue(executionState)
	suite.NotNil(createdEntity.Id, "created artifact id should not be nil")
	suite.Equal(entityName, *getById.Name)
	suite.Equal(*state, *getById.LastKnownState)
	suite.Equal(createdVersionId, getById.ModelVersionId)
	suite.Equal(customString, (*getById.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*createdEntity, *getById, "artifacts returned during creation and on get by id should be equal")
}

func (suite *CoreTestSuite) TestGetServeModels() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)
	inferenceServiceId := suite.registerInferenceService(service, registeredModelId, nil, nil, nil, nil)

	modelVersion1Name := "v1"
	modelVersion1 := &openapi.ModelVersion{Name: modelVersion1Name, Description: &modelVersionDescription}
	createdVersion1, err := service.UpsertModelVersion(modelVersion1, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %s", registeredModelId)
	createdVersion1Id := *createdVersion1.Id

	modelVersion2Name := "v2"
	modelVersion2 := &openapi.ModelVersion{Name: modelVersion2Name, Description: &modelVersionDescription}
	createdVersion2, err := service.UpsertModelVersion(modelVersion2, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %s", registeredModelId)
	createdVersion2Id := *createdVersion2.Id

	modelVersion3Name := "v3"
	modelVersion3 := &openapi.ModelVersion{Name: modelVersion3Name, Description: &modelVersionDescription}
	createdVersion3, err := service.UpsertModelVersion(modelVersion3, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %s", registeredModelId)
	createdVersion3Id := *createdVersion3.Id
	// end of data preparation

	eut1Name := "sm1"
	eut1 := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		Description:    &entityDescription,
		Name:           &eut1Name,
		ModelVersionId: createdVersion1Id,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	eut2Name := "sm2"
	eut2 := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		Description:    &entityDescription,
		Name:           &eut2Name,
		ModelVersionId: createdVersion2Id,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	eut3Name := "sm3"
	eut3 := &openapi.ServeModel{
		LastKnownState: (*openapi.ExecutionState)(&executionState),
		Description:    &entityDescription,
		Name:           &eut3Name,
		ModelVersionId: createdVersion3Id,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity1, err := service.UpsertServeModel(eut1, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %s", inferenceServiceId)
	createdEntity2, err := service.UpsertServeModel(eut2, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %s", inferenceServiceId)
	createdEntity3, err := service.UpsertServeModel(eut3, &inferenceServiceId)
	suite.Nilf(err, "error creating new ServeModel for %s", inferenceServiceId)

	createdEntityId1, _ := converter.StringToInt64(createdEntity1.Id)
	createdEntityId2, _ := converter.StringToInt64(createdEntity2.Id)
	createdEntityId3, _ := converter.StringToInt64(createdEntity3.Id)

	getAll, err := service.GetServeModels(api.ListOptions{}, nil)
	suite.Nilf(err, "error getting all ServeModel")
	suite.Equalf(int32(3), getAll.Size, "expected three ServeModel")

	suite.Equal(*converter.Int64ToString(createdEntityId1), *getAll.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdEntityId2), *getAll.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdEntityId3), *getAll.Items[2].Id)

	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByInferenceService, err := service.GetServeModels(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &inferenceServiceId)
	suite.Nilf(err, "error getting all ServeModels for %s", inferenceServiceId)
	suite.Equalf(int32(3), getAllByInferenceService.Size, "expected three ServeModels for InferenceServiceId %s", inferenceServiceId)

	suite.Equal(*converter.Int64ToString(createdEntityId1), *getAllByInferenceService.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdEntityId2), *getAllByInferenceService.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdEntityId3), *getAllByInferenceService.Items[0].Id)
}
