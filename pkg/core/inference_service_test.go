package core

import (
	"context"
	"fmt"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// INFERENCE SERVICE

func (suite *CoreTestSuite) TestCreateInferenceService() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)
	runtime := "model-server"
	desiredState := openapi.INFERENCESERVICESTATE_DEPLOYED

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		Runtime:              &runtime,
		DesiredState:         &desiredState,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %s: %v", parentResourceId, err)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	byId, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdEntityId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	suite.Equal(*createdEntityId, *byId.Contexts[0].Id, "returned id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", parentResourceId, entityName), *byId.Contexts[0].Name, "saved name should match the provided one")
	suite.Equal(entityExternalId2, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(customString, byId.Contexts[0].CustomProperties["custom_string_prop"].GetStringValue(), "saved custom_string_prop custom property should match the provided one")
	suite.Equal(entityDescription, byId.Contexts[0].Properties["description"].GetStringValue(), "saved description should match the provided one")
	suite.Equal(runtime, byId.Contexts[0].Properties["runtime"].GetStringValue(), "saved runtime should match the provided one")
	suite.Equal(string(desiredState), byId.Contexts[0].Properties["desired_state"].GetStringValue(), "saved state should match the provided one")
	suite.Equalf(*inferenceServiceTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *inferenceServiceTypeName)

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(3, len(getAllResp.Contexts), "there should be 3 contexts (RegisteredModel, ServingEnvironment, InferenceService) saved in mlmd")
}

func (suite *CoreTestSuite) TestCreateDuplicateInferenceServiceFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)
	runtime := "model-server"
	desiredState := openapi.INFERENCESERVICESTATE_DEPLOYED

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		Runtime:              &runtime,
		DesiredState:         &desiredState,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	_, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %s: %v", parentResourceId, err)

	// attempt to create dupliate inference service
	_, err = service.UpsertInferenceService(eut)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "cannot register a duplicate inference service")
	suite.Equal(409, statusResp, "duplicate inference services not allowed")

}

func (suite *CoreTestSuite) TestCreateInferenceServiceFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	_, err := service.UpsertInferenceService(nil)
	suite.NotNil(err)
	suite.Equal("invalid inference service pointer, can't upsert nil: bad request", err.Error())

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		ServingEnvironmentId: "9999",
		RegisteredModelId:    "9998",
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	_, err = service.UpsertInferenceService(eut)
	suite.NotNil(err)
	suite.Equal("no serving environment found for id 9999: not found", err.Error())

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	eut.ServingEnvironmentId = parentResourceId

	_, err = service.UpsertInferenceService(eut)
	suite.NotNil(err)
	suite.Equal("no registered model found for id 9998: not found", err.Error())
}

func (suite *CoreTestSuite) TestUpdateInferenceService() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %s", parentResourceId)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	newExternalId := "org.my_awesome_entity@v1"
	newScore := 0.95

	createdEntity.ExternalId = &newExternalId
	(*createdEntity.CustomProperties)["score"] = openapi.MetadataValue{
		MetadataDoubleValue: converter.NewMetadataDoubleValue(newScore),
	}

	updatedEntity, err := service.UpsertInferenceService(createdEntity)
	suite.Nilf(err, "error updating new entity for %s: %v", registeredModelId, err)

	updateEntityId, _ := converter.StringToInt64(updatedEntity.Id)
	suite.Equal(*createdEntityId, *updateEntityId, "created and updated should have same id")

	byId, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateEntityId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be 1 context saved in mlmd by id")

	suite.Equal(*updateEntityId, *byId.Contexts[0].Id, "returned id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", parentResourceId, *eut.Name), *byId.Contexts[0].Name, "saved name should match the provided one")
	suite.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(customString, byId.Contexts[0].CustomProperties["custom_string_prop"].GetStringValue(), "saved custom_string_prop custom property should match the provided one")
	suite.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	suite.Equalf(*inferenceServiceTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *inferenceServiceTypeName)

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(3, len(getAllResp.Contexts), "there should be 3 contexts saved in mlmd")

	// update with nil name
	newExternalId = "org.my_awesome_entity_@v1"
	updatedEntity.ExternalId = &newExternalId
	updatedEntity.Name = nil
	updatedEntity, err = service.UpsertInferenceService(updatedEntity)
	suite.Nilf(err, "error updating new model version for %s: %v", updateEntityId, err)

	updateEntityId, _ = converter.StringToInt64(updatedEntity.Id)
	suite.Equal(*createdEntityId, *updateEntityId, "created and updated should have same id")

	byId, err = suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateEntityId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be 1 context saved in mlmd by id")

	suite.Equal(*updateEntityId, *byId.Contexts[0].Id, "returned id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", parentResourceId, *eut.Name), *byId.Contexts[0].Name, "saved name should match the provided one")
	suite.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(customString, byId.Contexts[0].CustomProperties["custom_string_prop"].GetStringValue(), "saved custom_string_prop custom property should match the provided one")
	suite.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	suite.Equalf(*inferenceServiceTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *inferenceServiceTypeName)

	// update with empty registeredModelId
	newExternalId = "org.my_awesome_entity_@v1"
	prevRegModelId := updatedEntity.RegisteredModelId
	updatedEntity.RegisteredModelId = ""
	updatedEntity, err = service.UpsertInferenceService(updatedEntity)
	suite.Nil(err)
	suite.Equal(prevRegModelId, updatedEntity.RegisteredModelId)
}

func (suite *CoreTestSuite) TestUpdateInferenceServiceFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %s", parentResourceId)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")

	newExternalId := "org.my_awesome_entity@v1"
	newScore := 0.95

	createdEntity.ExternalId = &newExternalId
	(*createdEntity.CustomProperties)["score"] = openapi.MetadataValue{
		MetadataDoubleValue: converter.NewMetadataDoubleValue(newScore),
	}

	wrongId := "9999"
	createdEntity.Id = &wrongId
	_, err = service.UpsertInferenceService(createdEntity)
	suite.NotNil(err)
	suite.Equal(fmt.Sprintf("no InferenceService found for id %s: not found", wrongId), err.Error())
}

func (suite *CoreTestSuite) TestGetInferenceServiceById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	state := openapi.INFERENCESERVICESTATE_UNDEPLOYED
	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		DesiredState:         &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %s", parentResourceId)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getById, err := service.GetInferenceServiceById(*createdEntity.Id)
	suite.Nilf(err, "error getting model version with id %d", *createdEntityId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdEntityId,
		},
	})
	suite.Nilf(err, "error retrieving context, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*getById.Id, *converter.Int64ToString(ctx.Id), "returned id should match the mlmd context one")
	suite.Equal(*eut.Name, *getById.Name, "saved name should match the provided one")
	suite.Equal(*eut.ExternalId, *getById.ExternalId, "saved external id should match the provided one")
	suite.Equal(*eut.DesiredState, *getById.DesiredState, "saved state should match the provided one")
	suite.Equal((*getById.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, customString, "saved custom_string_prop custom property should match the provided one")
}

func (suite *CoreTestSuite) TestGetRegisteredModelByInferenceServiceId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}
	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %s", parentResourceId)
	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")

	getRM, err := service.GetRegisteredModelByInferenceService(*createdEntity.Id)
	suite.Nilf(err, "error getting using id %s", *createdEntity.Id)

	suite.Equal(registeredModelId, *getRM.Id, "returned id should match the original registeredModelId")
}

func (suite *CoreTestSuite) TestGetModelVersionByInferenceServiceId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion1Name := "v1"
	modelVersion1 := &openapi.ModelVersion{Name: modelVersion1Name, Description: &modelVersionDescription}
	createdVersion1, err := service.UpsertModelVersion(modelVersion1, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersion1Id := *createdVersion1.Id

	modelVersion2Name := "v2"
	modelVersion2 := &openapi.ModelVersion{Name: modelVersion2Name, Description: &modelVersionDescription}
	createdVersion2, err := service.UpsertModelVersion(modelVersion2, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	createdVersion2Id := *createdVersion2.Id
	// end of data preparation

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		ModelVersionId:       nil, // first we test by unspecified
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}
	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %s", parentResourceId)

	getVModel, err := service.GetModelVersionByInferenceService(*createdEntity.Id)
	suite.Nilf(err, "error getting using id %s", *createdEntity.Id)
	suite.Equal(createdVersion2Id, *getVModel.Id, "returned id shall be the latest ModelVersion by creation order")

	// here we used the returned entity (so ID is populated), and we update to specify the "ID of the ModelVersion to serve"
	createdEntity.ModelVersionId = &createdVersion1Id
	_, err = service.UpsertInferenceService(createdEntity)
	suite.Nilf(err, "error updating eut for %s", parentResourceId)

	getVModel, err = service.GetModelVersionByInferenceService(*createdEntity.Id)
	suite.Nilf(err, "error getting using id %s", *createdEntity.Id)
	suite.Equal(createdVersion1Id, *getVModel.Id, "returned id shall be the specified one")
}

func (suite *CoreTestSuite) TestGetModelArtifactByInferenceServiceId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion1Name := "v1"
	modelVersion1 := &openapi.ModelVersion{Name: modelVersion1Name, Description: &modelVersionDescription}
	createdVersion1, err := service.UpsertModelVersion(modelVersion1, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %s", registeredModelId)
	modelArtifact1Name := "v1-artifact"
	modelArtifact1 := &openapi.ModelArtifact{Name: &modelArtifact1Name}
	art1, err := service.UpsertModelVersionArtifact(&openapi.Artifact{ModelArtifact: modelArtifact1}, *createdVersion1.Id)
	suite.Nilf(err, "error creating new model artifact for %s", *createdVersion1.Id)
	ma1 := art1.ModelArtifact

	modelVersion2Name := "v2"
	modelVersion2 := &openapi.ModelVersion{Name: modelVersion2Name, Description: &modelVersionDescription}
	createdVersion2, err := service.UpsertModelVersion(modelVersion2, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %s", registeredModelId)
	modelArtifact2Name := "v2-artifact"
	modelArtifact2 := &openapi.ModelArtifact{Name: &modelArtifact2Name}
	art2, err := service.UpsertModelVersionArtifact(&openapi.Artifact{ModelArtifact: modelArtifact2}, *createdVersion2.Id)
	suite.Nilf(err, "error creating new model artifact for %s", *createdVersion2.Id)
	ma2 := art2.ModelArtifact
	// end of data preparation

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		ModelVersionId:       nil, // first we test by unspecified
	}
	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %s", parentResourceId)

	getModelArt, err := service.GetModelArtifactByInferenceService(*createdEntity.Id)
	suite.Nilf(err, "error getting using id %s", *createdEntity.Id)
	suite.Equal(*ma2.Id, *getModelArt.Id, "returned id shall be the latest ModelVersion by creation order")

	// here we used the returned entity (so ID is populated), and we update to specify the "ID of the ModelVersion to serve"
	createdEntity.ModelVersionId = createdVersion1.Id
	_, err = service.UpsertInferenceService(createdEntity)
	suite.Nilf(err, "error updating eut for %s", parentResourceId)

	getModelArt, err = service.GetModelArtifactByInferenceService(*createdEntity.Id)
	suite.Nilf(err, "error getting using id %s", *createdEntity.Id)
	suite.Equal(*ma1.Id, *getModelArt.Id, "returned id shall be the specified one")
}

func (suite *CoreTestSuite) TestGetInferenceServiceByParamsWithNoResults() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)

	_, err := service.GetInferenceServiceByParams(apiutils.Of("not-present"), &parentResourceId, nil)
	suite.NotNil(err)
	suite.Equal("no inference services found for name=not-present, servingEnvironmentId=1, externalId=: not found", err.Error())
}

func (suite *CoreTestSuite) TestGetInferenceServiceByParamsName() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %s", parentResourceId)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getByName, err := service.GetInferenceServiceByParams(&entityName, &parentResourceId, nil)
	suite.Nilf(err, "error getting model version by name %d", *createdEntityId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdEntityId,
		},
	})
	suite.Nilf(err, "error retrieving context, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getByName.Id, "returned id should match the mlmd context one")
	suite.Equal(fmt.Sprintf("%s:%s", parentResourceId, *getByName.Name), *ctx.Name, "saved name should match the provided one")
	suite.Equal(*ctx.ExternalId, *getByName.ExternalId, "saved external id should match the provided one")
	suite.Equal(ctx.CustomProperties["custom_string_prop"].GetStringValue(), (*getByName.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, "saved custom_string_prop custom property should match the provided one")
}

func (suite *CoreTestSuite) TestGetInferenceServiceByParamInvalid() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	// must register an inference service first, otherwise the http error will be a 404
	_, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %s", parentResourceId)

	invalidName := "\xFF"

	_, err = service.GetInferenceServiceByParams(&invalidName, &parentResourceId, nil)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "invalid parameter used to retreive inference service")
	suite.Equal(400, statusResp, "invalid parameter used to retreive inference service")
}

func (suite *CoreTestSuite) TestGetInfernenceServiceByParamsExternalId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %s", parentResourceId)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	getByExternalId, err := service.GetInferenceServiceByParams(nil, nil, eut.ExternalId)
	suite.Nilf(err, "error getting by external id %d", *eut.ExternalId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdEntityId,
		},
	})
	suite.Nilf(err, "error retrieving context, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getByExternalId.Id, "returned id should match the mlmd context one")
	suite.Equal(fmt.Sprintf("%s:%s", parentResourceId, *getByExternalId.Name), *ctx.Name, "saved name should match the provided one")
	suite.Equal(*ctx.ExternalId, *getByExternalId.ExternalId, "saved external id should match the provided one")
	suite.Equal(ctx.CustomProperties["custom_string_prop"].GetStringValue(), (*getByExternalId.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, "saved custom_string_prop custom property should match the provided one")
}

func (suite *CoreTestSuite) TestGetInferenceServiceByEmptyParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		Description:          &entityDescription,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdEntity, err := service.UpsertInferenceService(eut)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	suite.NotNilf(createdEntity.Id, "created eut should not have nil Id")

	_, err = service.GetInferenceServiceByParams(nil, nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid parameters call, supply either (name and servingEnvironmentId), or externalId: bad request", err.Error())
}

func (suite *CoreTestSuite) TestGetInferenceServices() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	parentResourceId := suite.registerServingEnvironment(service, "", nil)
	registeredModelId := suite.registerModel(service, nil, nil)

	eut1 := &openapi.InferenceService{
		Name:                 &entityName,
		ExternalId:           &entityExternalId2,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		Runtime:              apiutils.Of("model-server0"),
	}

	secondName := "v2"
	secondExtId := "org.myawesomeentity@v2"
	eut2 := &openapi.InferenceService{
		Name:                 &secondName,
		ExternalId:           &secondExtId,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		Runtime:              apiutils.Of("model-server1"),
	}

	thirdName := "v3"
	thirdExtId := "org.myawesomeentity@v3"
	eut3 := &openapi.InferenceService{
		Name:                 &thirdName,
		ExternalId:           &thirdExtId,
		ServingEnvironmentId: parentResourceId,
		RegisteredModelId:    registeredModelId,
		Runtime:              apiutils.Of("model-server2"),
	}

	createdEntity1, err := service.UpsertInferenceService(eut1)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	createdEntity2, err := service.UpsertInferenceService(eut2)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	createdEntity3, err := service.UpsertInferenceService(eut3)
	suite.Nilf(err, "error creating new eut for %v", parentResourceId)

	anotherParentResourceName := "AnotherModel"
	anotherParentResourceExtId := "org.another"
	anotherParentResourceId := suite.registerServingEnvironment(service, anotherParentResourceName, &anotherParentResourceExtId)

	anotherName := "v1.0"
	anotherExtId := "org.another@v1.0"
	eutAnother := &openapi.InferenceService{
		Name:                 &anotherName,
		ExternalId:           &anotherExtId,
		ServingEnvironmentId: anotherParentResourceId,
		RegisteredModelId:    registeredModelId,
		Runtime:              apiutils.Of("model-server3"),
	}

	_, err = service.UpsertInferenceService(eutAnother)
	suite.Nilf(err, "error creating new model version for %d", anotherParentResourceId)

	createdId1, _ := converter.StringToInt64(createdEntity1.Id)
	createdId2, _ := converter.StringToInt64(createdEntity2.Id)
	createdId3, _ := converter.StringToInt64(createdEntity3.Id)

	getAll, err := service.GetInferenceServices(api.ListOptions{}, nil, nil)
	suite.Nilf(err, "error getting all")
	suite.Equal(int32(4), getAll.Size, "expected 4 across all parent resources")

	getAllByParentResource, err := service.GetInferenceServices(api.ListOptions{}, &parentResourceId, nil)
	suite.Nilf(err, "error getting all")
	suite.Equalf(int32(3), getAllByParentResource.Size, "expected 3 for parent resource %d", parentResourceId)

	suite.Equal(*converter.Int64ToString(createdId1), *getAllByParentResource.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdId2), *getAllByParentResource.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdId3), *getAllByParentResource.Items[2].Id)

	modelServer := "model-server1"
	getAllByParentResourceAndRuntime, err := service.GetInferenceServices(api.ListOptions{}, &parentResourceId, &modelServer)
	suite.Nilf(err, "error getting all")
	suite.Equalf(int32(1), getAllByParentResourceAndRuntime.Size, "expected 1 for parent resource %s and runtime %s", parentResourceId, modelServer)

	suite.Equal(*converter.Int64ToString(createdId1), *getAllByParentResource.Items[0].Id)

	// order by last update time, expecting last created as first
	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByParentResource, err = service.GetInferenceServices(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &parentResourceId, nil)
	suite.Nilf(err, "error getting all")
	suite.Equalf(int32(3), getAllByParentResource.Size, "expected 3 for parent resource %d", parentResourceId)

	suite.Equal(*converter.Int64ToString(createdId1), *getAllByParentResource.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdId2), *getAllByParentResource.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdId3), *getAllByParentResource.Items[0].Id)

	// update the second entity
	newExternalId := "updated.org:v2"
	createdEntity2.ExternalId = &newExternalId
	createdEntity2, err = service.UpsertInferenceService(createdEntity2)
	suite.Nilf(err, "error creating new eut2 for %d", parentResourceId)

	suite.Equal(newExternalId, *createdEntity2.ExternalId)

	getAllByParentResource, err = service.GetInferenceServices(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &parentResourceId, nil)
	suite.Nilf(err, "error getting all")
	suite.Equalf(int32(3), getAllByParentResource.Size, "expected 3 for parent resource %d", parentResourceId)

	suite.Equal(*converter.Int64ToString(createdId1), *getAllByParentResource.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdId2), *getAllByParentResource.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdId3), *getAllByParentResource.Items[1].Id)
}
