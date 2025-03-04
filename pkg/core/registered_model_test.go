package core

import (
	"context"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// REGISTERED MODELS

func (suite *CoreTestSuite) TestCreateRegisteredModel() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	state := openapi.REGISTEREDMODELSTATE_ARCHIVED
	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:        modelName,
		ExternalId:  &modelExternalId,
		Description: &modelDescription,
		Owner:       &modelOwner,
		State:       &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)

	// checks
	suite.Nilf(err, "error creating registered model: %v", err)
	suite.NotNilf(createdModel.Id, "created registered model should not have nil Id")

	createdModelId, _ := converter.StringToInt64(createdModel.Id)
	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdModelId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	suite.Equal(*createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	suite.Equal(modelName, *ctx.Name, "saved model name should match the provided one")
	suite.Equal(modelExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(modelDescription, ctx.Properties["description"].GetStringValue(), "saved description should match the provided one")
	suite.Equal(modelOwner, ctx.Properties["owner"].GetStringValue(), "saved owner should match the provided one")
	suite.Equal(string(state), ctx.Properties["state"].GetStringValue(), "saved state should match the provided one")
	suite.Equal(myCustomProp, ctx.CustomProperties["myCustomProp"].GetStringValue(), "saved myCustomProp custom property should match the provided one")

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")
}

func (suite *CoreTestSuite) TestCreateDuplicateRegisteredModelFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	state := openapi.REGISTEREDMODELSTATE_ARCHIVED
	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:        modelName,
		ExternalId:  &modelExternalId,
		Description: &modelDescription,
		Owner:       &modelOwner,
		State:       &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	// create the first model
	_, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	// attempt to create dupliate model
	_, err = service.UpsertRegisteredModel(registeredModel)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "cannot register a model with duplicate names")
	suite.Equal(409, statusResp, "duplicate model names not allowed")
}

func (suite *CoreTestSuite) TestUpdateRegisteredModel() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       modelName,
		Owner:      &modelOwner,
		ExternalId: &modelExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)

	// checks
	suite.Nilf(err, "error creating registered model: %v", err)
	suite.NotNilf(createdModel.Id, "created registered model should not have nil Id")
	createdModelId, _ := converter.StringToInt64(createdModel.Id)

	// checks created model matches original one except for Id
	suite.Equal(registeredModel.Name, createdModel.Name, "returned model name should match the original one")
	suite.Equal(*registeredModel.ExternalId, *createdModel.ExternalId, "returned model external id should match the original one")
	suite.Equal(*registeredModel.CustomProperties, *createdModel.CustomProperties, "returned model custom props should match the original one")

	// update existing model
	newModelExternalId := "newExternalId"
	newOwner := "newOwner"
	newCustomProp := "updated myCustomProp"

	createdModel.ExternalId = &newModelExternalId
	createdModel.Owner = &newOwner
	(*createdModel.CustomProperties)["myCustomProp"] = openapi.MetadataValue{
		MetadataStringValue: converter.NewMetadataStringValue(newCustomProp),
	}
	// check can also define customProperty of name "owner", in addition to built-in property "owner"
	(*createdModel.CustomProperties)["owner"] = openapi.MetadataValue{
		MetadataStringValue: converter.NewMetadataStringValue(newCustomProp),
	}

	// update the model
	createdModel, err = service.UpsertRegisteredModel(createdModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	// still one registered model
	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdModelId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	suite.Equal(*createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	suite.Equal(modelName, *ctx.Name, "saved model name should match the provided one")
	suite.Equal(newModelExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(newOwner, ctx.Properties["owner"].GetStringValue(), "saved owner custom property should match the provided one")
	suite.Equal(newCustomProp, ctx.CustomProperties["myCustomProp"].GetStringValue(), "saved myCustomProp custom property should match the provided one")
	suite.Equal(newCustomProp, ctx.CustomProperties["owner"].GetStringValue(), "check can define custom property 'onwer' and should match the provided one")

	// update the model keeping nil name
	newModelExternalId = "newNewExternalId"
	createdModel.ExternalId = &newModelExternalId
	createdModel.Name = ""
	createdModel, err = service.UpsertRegisteredModel(createdModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	// still one registered model
	getAllResp, err = suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err = suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdModelId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx = ctxById.Contexts[0]
	ctxId = converter.Int64ToString(ctx.Id)
	suite.Equal(*createdModel.Id, *ctxId, "returned model id should match the mlmd one")
	suite.Equal(modelName, *ctx.Name, "saved model name should match the provided one")
	suite.Equal(newModelExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(newOwner, ctx.Properties["owner"].GetStringValue(), "saved owner custom property should match the provided one")
	suite.Equal(newCustomProp, ctx.CustomProperties["myCustomProp"].GetStringValue(), "saved myCustomProp custom property should match the provided one")
	suite.Equal(newCustomProp, ctx.CustomProperties["owner"].GetStringValue(), "check can define custom property 'onwer' and should match the provided one")
}

func (suite *CoreTestSuite) TestGetRegisteredModelById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	state := openapi.REGISTEREDMODELSTATE_LIVE
	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       modelName,
		ExternalId: &modelExternalId,
		State:      &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	// test
	createdModel, err := service.UpsertRegisteredModel(registeredModel)

	// checks
	suite.Nilf(err, "error creating registered model: %v", err)

	getModelById, err := service.GetRegisteredModelById(*createdModel.Id)
	suite.Nilf(err, "error getting registered model by id %s: %v", *createdModel.Id, err)

	// checks created model matches original one except for Id
	suite.Equal(registeredModel.Name, getModelById.Name, "saved model name should match the original one")
	suite.Equal(*registeredModel.ExternalId, *getModelById.ExternalId, "saved model external id should match the original one")
	suite.Equal(*registeredModel.State, *getModelById.State, "saved model state should match the original one")
	suite.Equal(*registeredModel.CustomProperties, *getModelById.CustomProperties, "saved model custom props should match the original one")
}

func (suite *CoreTestSuite) TestGetRegisteredModelByParamsWithNoResults() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	_, err := service.GetRegisteredModelByParams(apiutils.Of("not-present"), nil)
	suite.NotNil(err)
	suite.Equal("no registered models found for name=not-present, externalId=: not found", err.Error())
}

func (suite *CoreTestSuite) TestGetRegisteredModelByParamsName() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       modelName,
		ExternalId: &modelExternalId,
	}

	createdModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	byName, err := service.GetRegisteredModelByParams(&modelName, nil)
	suite.Nilf(err, "error getting registered model by name: %v", err)

	suite.Equalf(*createdModel.Id, *byName.Id, "the returned model id should match the retrieved by name")
}

func (suite *CoreTestSuite) TestGetRegisteredModelByParamsInvalid() {
	// trigger a 400 bad request to test unallowed query params
	// create mode registry service
	service := suite.setupModelRegistryService()

	registeredModel := &openapi.RegisteredModel{
		Name:       modelName,
		ExternalId: &modelExternalId,
	}

	// must register a model first, otherwise the http error will be a 404
	_, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	invalidName := "\xFF"

	_, err = service.GetRegisteredModelByParams(&invalidName, nil)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "invalid parameter used to retreive registered model")
	suite.Equal(400, statusResp, "invalid parameter used to retreive registered model")
}

func (suite *CoreTestSuite) TestGetRegisteredModelByParamsExternalId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       modelName,
		ExternalId: &modelExternalId,
	}

	createdModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	byName, err := service.GetRegisteredModelByParams(nil, &modelExternalId)
	suite.Nilf(err, "error getting registered model by external id: %v", err)

	suite.Equalf(*createdModel.Id, *byName.Id, "the returned model id should match the retrieved by name")
}

func (suite *CoreTestSuite) TestGetRegisteredModelByEmptyParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       modelName,
		ExternalId: &modelExternalId,
	}

	_, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	_, err = service.GetRegisteredModelByParams(nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid parameters call, supply either name or externalId: bad request", err.Error())
}

func (suite *CoreTestSuite) TestGetRegisteredModelsOrderedById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	orderBy := "ID"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       modelName,
		ExternalId: &modelExternalId,
	}

	_, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	newModelName := "PricingModel2"
	newModelExternalId := "myExternalId2"
	registeredModel.Name = newModelName
	registeredModel.ExternalId = &newModelExternalId
	_, err = service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	newModelName = "PricingModel3"
	newModelExternalId = "myExternalId3"
	registeredModel.Name = newModelName
	registeredModel.ExternalId = &newModelExternalId
	_, err = service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	orderedById, err := service.GetRegisteredModels(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	suite.Nilf(err, "error getting registered models: %v", err)

	suite.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		suite.Less(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}

	orderedById, err = service.GetRegisteredModels(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	suite.Nilf(err, "error getting registered models: %v", err)

	suite.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		suite.Greater(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}
}

func (suite *CoreTestSuite) TestGetRegisteredModelsOrderedByLastUpdate() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	orderBy := "LAST_UPDATE_TIME"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       modelName,
		ExternalId: &modelExternalId,
	}

	firstModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	newModelName := "PricingModel2"
	newModelExternalId := "myExternalId2"
	registeredModel.Name = newModelName
	registeredModel.ExternalId = &newModelExternalId
	secondModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	newModelName = "PricingModel3"
	newModelExternalId = "myExternalId3"
	registeredModel.Name = newModelName
	registeredModel.ExternalId = &newModelExternalId
	thirdModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	// update second model
	secondModel.ExternalId = nil
	_, err = service.UpsertRegisteredModel(secondModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	orderedById, err := service.GetRegisteredModels(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	suite.Nilf(err, "error getting registered models: %v", err)

	suite.Equal(3, int(orderedById.Size))
	suite.Equal(*firstModel.Id, *orderedById.Items[0].Id)
	suite.Equal(*thirdModel.Id, *orderedById.Items[1].Id)
	suite.Equal(*secondModel.Id, *orderedById.Items[2].Id)

	orderedById, err = service.GetRegisteredModels(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	suite.Nilf(err, "error getting registered models: %v", err)

	suite.Equal(3, int(orderedById.Size))
	suite.Equal(*secondModel.Id, *orderedById.Items[0].Id)
	suite.Equal(*thirdModel.Id, *orderedById.Items[1].Id)
	suite.Equal(*firstModel.Id, *orderedById.Items[2].Id)
}

func (suite *CoreTestSuite) TestGetRegisteredModelsWithPageSize() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	pageSize := int32(1)
	pageSize2 := int32(2)
	modelName := "PricingModel1"
	modelExternalId := "myExternalId1"

	// register a new model
	registeredModel := &openapi.RegisteredModel{
		Name:       modelName,
		ExternalId: &modelExternalId,
	}

	firstModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	newModelName := "PricingModel2"
	newModelExternalId := "myExternalId2"
	registeredModel.Name = newModelName
	registeredModel.ExternalId = &newModelExternalId
	secondModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	newModelName = "PricingModel3"
	newModelExternalId = "myExternalId3"
	registeredModel.Name = newModelName
	registeredModel.ExternalId = &newModelExternalId
	thirdModel, err := service.UpsertRegisteredModel(registeredModel)
	suite.Nilf(err, "error creating registered model: %v", err)

	truncatedList, err := service.GetRegisteredModels(api.ListOptions{
		PageSize: &pageSize,
	})
	suite.Nilf(err, "error getting registered models: %v", err)

	suite.Equal(1, int(truncatedList.Size))
	suite.NotEqual("", truncatedList.NextPageToken, "next page token should not be empty")
	suite.Equal(*firstModel.Id, *truncatedList.Items[0].Id)

	truncatedList, err = service.GetRegisteredModels(api.ListOptions{
		PageSize:      &pageSize2,
		NextPageToken: &truncatedList.NextPageToken,
	})
	suite.Nilf(err, "error getting registered models: %v", err)

	suite.Equal(2, int(truncatedList.Size))
	suite.Equal("", truncatedList.NextPageToken, "next page token should be empty as list item returned")
	suite.Equal(*secondModel.Id, *truncatedList.Items[0].Id)
	suite.Equal(*thirdModel.Id, *truncatedList.Items[1].Id)
}
