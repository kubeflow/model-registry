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

// MODEL VERSIONS

func (suite *CoreTestSuite) TestCreateModelVersion() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	state := openapi.MODELVERSIONSTATE_LIVE
	modelVersion := &openapi.ModelVersion{
		Name:        modelVersionName,
		ExternalId:  &versionExternalId,
		Description: &modelVersionDescription,
		State:       &state,
		Author:      &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	suite.Equal((*createdVersion).RegisteredModelId, registeredModelId, "RegisteredModelId should match the actual owner-entity")

	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")

	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	byId, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	suite.Equal(*createdVersionId, *byId.Contexts[0].Id, "returned model id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", registeredModelId, modelVersionName), *byId.Contexts[0].Name, "saved model name should match the provided one")
	suite.Equal(versionExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(author, byId.Contexts[0].Properties["author"].GetStringValue(), "saved author property should match the provided one")
	suite.Equal(modelVersionDescription, byId.Contexts[0].Properties["description"].GetStringValue(), "saved description should match the provided one")
	suite.Equal(string(state), byId.Contexts[0].Properties["state"].GetStringValue(), "saved state should match the provided one")
	suite.Equalf(*modelVersionTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *modelVersionTypeName)

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(2, len(getAllResp.Contexts), "there should be two contexts saved in mlmd")
}

func (suite *CoreTestSuite) TestCreateDuplicateModelVersionFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	state := openapi.MODELVERSIONSTATE_LIVE
	modelVersion := &openapi.ModelVersion{
		Name:        modelVersionName,
		ExternalId:  &versionExternalId,
		Description: &modelVersionDescription,
		State:       &state,
		Author:      &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	suite.Equal((*createdVersion).RegisteredModelId, registeredModelId, "RegisteredModelId should match the actual owner-entity")

	// attempt to create dupliate model version
	_, err = service.UpsertModelVersion(modelVersion, &registeredModelId)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "cannot register a duplicate model version")
	suite.Equal(409, statusResp, "duplicate model versions not allowed")
}

func (suite *CoreTestSuite) TestCreateModelVersionFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	_, err := service.UpsertModelVersion(nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid model version pointer, can't upsert nil: bad request", err.Error())

	registeredModelId := "9999"

	modelVersion := &openapi.ModelVersion{
		Name:       modelVersionName,
		ExternalId: &versionExternalId,
		Author:     &author,
	}

	_, err = service.UpsertModelVersion(modelVersion, nil)
	suite.NotNil(err)
	suite.Equal("missing registered model id, cannot create model version without registered model: bad request", err.Error())

	_, err = service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.NotNil(err)
	suite.Equal("no registered model found for id 9999: not found", err.Error())
}

func (suite *CoreTestSuite) TestUpdateModelVersion() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       modelVersionName,
		ExternalId: &versionExternalId,
		Author:     &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	newExternalId := "org.my_awesome_model@v1"
	newScore := 0.95

	createdVersion.ExternalId = &newExternalId
	(*createdVersion.CustomProperties)["score"] = openapi.MetadataValue{
		MetadataDoubleValue: converter.NewMetadataDoubleValue(newScore),
	}

	updatedVersion, err := service.UpsertModelVersion(createdVersion, &registeredModelId)
	suite.Nilf(err, "error updating new model version for %s: %v", registeredModelId, err)
	suite.Equal((*updatedVersion).RegisteredModelId, registeredModelId, "RegisteredModelId should match the actual owner-entity")

	updateVersionId, _ := converter.StringToInt64(updatedVersion.Id)
	suite.Equal(*createdVersionId, *updateVersionId, "created and updated model version should have same id")

	byId, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateVersionId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	suite.Equal(*updateVersionId, *byId.Contexts[0].Id, "returned model id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", registeredModelId, modelVersionName), *byId.Contexts[0].Name, "saved model name should match the provided one")
	suite.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(author, byId.Contexts[0].Properties["author"].GetStringValue(), "saved author property should match the provided one")
	suite.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	suite.Equalf(*modelVersionTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *modelVersionTypeName)

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(2, len(getAllResp.Contexts), "there should be two contexts saved in mlmd")

	// update with nil name
	newExternalId = "org.my_awesome_model_@v1"
	updatedVersion.ExternalId = &newExternalId
	updatedVersion.Name = ""
	updatedVersion, err = service.UpsertModelVersion(updatedVersion, &registeredModelId)
	suite.Nilf(err, "error updating new model version for %s: %v", registeredModelId, err)

	updateVersionId, _ = converter.StringToInt64(updatedVersion.Id)
	suite.Equal(*createdVersionId, *updateVersionId, "created and updated model version should have same id")

	byId, err = suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateVersionId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	suite.Equal(*updateVersionId, *byId.Contexts[0].Id, "returned model id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", registeredModelId, modelVersionName), *byId.Contexts[0].Name, "saved model name should match the provided one")
	suite.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(author, byId.Contexts[0].Properties["author"].GetStringValue(), "saved author property should match the provided one")
	suite.Equal(newScore, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	suite.Equalf(*modelVersionTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *modelVersionTypeName)
}

func (suite *CoreTestSuite) TestUpdateModelVersionFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       modelVersionName,
		ExternalId: &versionExternalId,
		Author:     &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %s", registeredModelId)
	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")

	newExternalId := "org.my_awesome_model@v1"
	newScore := 0.95

	createdVersion.ExternalId = &newExternalId
	(*createdVersion.CustomProperties)["score"] = openapi.MetadataValue{
		MetadataDoubleValue: converter.NewMetadataDoubleValue(newScore),
	}

	wrongId := "9999"
	createdVersion.Id = &wrongId
	_, err = service.UpsertModelVersion(createdVersion, &registeredModelId)
	suite.NotNil(err)
	suite.Equal(fmt.Sprintf("no model version found for id %s: not found", wrongId), err.Error())
}

func (suite *CoreTestSuite) TestGetModelVersionById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	state := openapi.MODELVERSIONSTATE_ARCHIVED
	modelVersion := &openapi.ModelVersion{
		Name:       modelVersionName,
		ExternalId: &versionExternalId,
		State:      &state,
		Author:     &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	getById, err := service.GetModelVersionById(*createdVersion.Id)
	suite.Nilf(err, "error getting model version with id %d", *createdVersionId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getById.Id, "returned model version id should match the mlmd context one")
	suite.Equal(modelVersion.Name, getById.Name, "saved model name should match the provided one")
	suite.Equal(*modelVersion.ExternalId, *getById.ExternalId, "saved external id should match the provided one")
	suite.Equal(*modelVersion.State, *getById.State, "saved model state should match the original one")
	suite.Equal(*getById.Author, author, "saved author property should match the provided one")
}

func (suite *CoreTestSuite) TestGetModelVersionByParamsWithNoResults() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	_, err := service.GetModelVersionByParams(apiutils.Of("not-present"), &registeredModelId, nil)
	suite.NotNil(err)
	suite.Equal("no model versions found for versionName=not-present, registeredModelId=1, externalId=: not found", err.Error())
}

func (suite *CoreTestSuite) TestGetModelVersionByParamsName() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       modelVersionName,
		ExternalId: &versionExternalId,
		Author:     &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	getByName, err := service.GetModelVersionByParams(&modelVersionName, &registeredModelId, nil)
	suite.Nilf(err, "error getting model version by name %d", *createdVersionId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getByName.Id, "returned model version id should match the mlmd context one")
	suite.Equal(fmt.Sprintf("%s:%s", registeredModelId, getByName.Name), *ctx.Name, "saved model name should match the provided one")
	suite.Equal(*ctx.ExternalId, *getByName.ExternalId, "saved external id should match the provided one")
	suite.Equal(ctx.Properties["author"].GetStringValue(), *getByName.Author, "saved author property should match the provided one")
}

func (suite *CoreTestSuite) TestGetModelVersionByParamsInvalid() {
	// trigger a 400 bad request to test unallowed query params
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       modelVersionName,
		ExternalId: &versionExternalId,
		Author:     &author,
	}

	// must register a model version first, otherwise the http error will be a 404
	_, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	invalidName := "\xFF"

	_, err = service.GetModelVersionByParams(&invalidName, &registeredModelId, nil)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "invalid parameter used to retreive model version")
	suite.Equal(400, statusResp, "invalid parameter used to retreive model version")
}

func (suite *CoreTestSuite) TestGetModelVersionByParamsExternalId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       modelVersionName,
		ExternalId: &versionExternalId,
		Author:     &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")
	createdVersionId, _ := converter.StringToInt64(createdVersion.Id)

	getByExternalId, err := service.GetModelVersionByParams(nil, nil, modelVersion.ExternalId)
	suite.Nilf(err, "error getting model version by external id %d", *modelVersion.ExternalId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdVersionId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getByExternalId.Id, "returned model version id should match the mlmd context one")
	suite.Equal(fmt.Sprintf("%s:%s", registeredModelId, getByExternalId.Name), *ctx.Name, "saved model name should match the provided one")
	suite.Equal(*ctx.ExternalId, *getByExternalId.ExternalId, "saved external id should match the provided one")
	suite.Equal(ctx.Properties["author"].GetStringValue(), *getByExternalId.Author, "saved author property should match the provided one")
}

func (suite *CoreTestSuite) TestGetModelVersionByEmptyParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion := &openapi.ModelVersion{
		Name:       modelVersionName,
		ExternalId: &versionExternalId,
		Author:     &author,
	}

	createdVersion, err := service.UpsertModelVersion(modelVersion, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)
	suite.NotNilf(createdVersion.Id, "created model version should not have nil Id")

	_, err = service.GetModelVersionByParams(nil, nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid parameters call, supply either (versionName and registeredModelId), or externalId: bad request", err.Error())
}

func (suite *CoreTestSuite) TestGetModelVersions() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModelId := suite.registerModel(service, nil, nil)

	modelVersion1 := &openapi.ModelVersion{
		Name:       modelVersionName,
		ExternalId: &versionExternalId,
	}

	secondModelVersionName := "v2"
	secondModelVersionExtId := "org.myawesomemodel@v2"
	modelVersion2 := &openapi.ModelVersion{
		Name:       secondModelVersionName,
		ExternalId: &secondModelVersionExtId,
	}

	thirdModelVersionName := "v3"
	thirdModelVersionExtId := "org.myawesomemodel@v3"
	modelVersion3 := &openapi.ModelVersion{
		Name:       thirdModelVersionName,
		ExternalId: &thirdModelVersionExtId,
	}

	createdVersion1, err := service.UpsertModelVersion(modelVersion1, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	createdVersion2, err := service.UpsertModelVersion(modelVersion2, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	createdVersion3, err := service.UpsertModelVersion(modelVersion3, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	anotherRegModelName := "AnotherModel"
	anotherRegModelExtId := "org.another"
	anotherRegisteredModelId := suite.registerModel(service, &anotherRegModelName, &anotherRegModelExtId)

	anotherModelVersionName := "v1.0"
	anotherModelVersionExtId := "org.another@v1.0"
	modelVersionAnother := &openapi.ModelVersion{
		Name:       anotherModelVersionName,
		ExternalId: &anotherModelVersionExtId,
	}

	_, err = service.UpsertModelVersion(modelVersionAnother, &anotherRegisteredModelId)
	suite.Nilf(err, "error creating new model version for %d", anotherRegisteredModelId)

	createdVersionId1, _ := converter.StringToInt64(createdVersion1.Id)
	createdVersionId2, _ := converter.StringToInt64(createdVersion2.Id)
	createdVersionId3, _ := converter.StringToInt64(createdVersion3.Id)

	getAll, err := service.GetModelVersions(api.ListOptions{}, nil)
	suite.Nilf(err, "error getting all model versions")
	suite.Equal(int32(4), getAll.Size, "expected four model versions across all registered models")

	getAllByRegModel, err := service.GetModelVersions(api.ListOptions{}, &registeredModelId)
	suite.Nilf(err, "error getting all model versions")
	suite.Equalf(int32(3), getAllByRegModel.Size, "expected three model versions for registered model %d", registeredModelId)

	suite.Equal(*converter.Int64ToString(createdVersionId1), *getAllByRegModel.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdVersionId2), *getAllByRegModel.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdVersionId3), *getAllByRegModel.Items[2].Id)

	// order by last update time, expecting last created as first
	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByRegModel, err = service.GetModelVersions(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &registeredModelId)
	suite.Nilf(err, "error getting all model versions")
	suite.Equalf(int32(3), getAllByRegModel.Size, "expected three model versions for registered model %d", registeredModelId)

	suite.Equal(*converter.Int64ToString(createdVersionId1), *getAllByRegModel.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdVersionId2), *getAllByRegModel.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdVersionId3), *getAllByRegModel.Items[0].Id)

	// update the second version
	newVersionExternalId := "updated.org:v2"
	createdVersion2.ExternalId = &newVersionExternalId
	createdVersion2, err = service.UpsertModelVersion(createdVersion2, &registeredModelId)
	suite.Nilf(err, "error creating new model version for %d", registeredModelId)

	suite.Equal(newVersionExternalId, *createdVersion2.ExternalId)

	getAllByRegModel, err = service.GetModelVersions(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &registeredModelId)
	suite.Nilf(err, "error getting all model versions")
	suite.Equalf(int32(3), getAllByRegModel.Size, "expected three model versions for registered model %d", registeredModelId)

	suite.Equal(*converter.Int64ToString(createdVersionId1), *getAllByRegModel.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdVersionId2), *getAllByRegModel.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdVersionId3), *getAllByRegModel.Items[1].Id)
}
