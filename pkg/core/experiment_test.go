package core

import (
	"context"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// EXPERIMENTS

func (suite *CoreTestSuite) TestCreateExperiment() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	state := openapi.EXPERIMENTSTATE_LIVE
	// create a new experiment
	experiment := &openapi.Experiment{
		Name:        experimentName,
		ExternalId:  &experimentExternalId,
		Description: &experimentDescription,
		Owner:       &experimentOwner,
		State:       &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	// test
	createdExperiment, err := service.UpsertExperiment(experiment)

	// checks
	suite.Nilf(err, "error creating experiment: %v", err)
	suite.NotNilf(createdExperiment.Id, "created experiment should not have nil Id")

	createdExperimentId, _ := converter.StringToInt64(createdExperiment.Id)
	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdExperimentId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	byId, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdExperimentId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")
	suite.Equalf(*experimentTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *experimentTypeName)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	suite.Equal(*createdExperiment.Id, *ctxId, "returned experiment id should match the mlmd one")
	suite.Equal(experimentName, *ctx.Name, "saved experiment name should match the provided one")
	suite.Equal(experimentExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(experimentDescription, ctx.Properties["description"].GetStringValue(), "saved description should match the provided one")
	suite.Equal(experimentOwner, ctx.Properties["owner"].GetStringValue(), "saved owner should match the provided one")
	suite.Equal(string(state), ctx.Properties["state"].GetStringValue(), "saved state should match the provided one")
	suite.Equal(myCustomProp, ctx.CustomProperties["myCustomProp"].GetStringValue(), "saved myCustomProp custom property should match the provided one")

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")
}

func (suite *CoreTestSuite) TestCreateDuplicateExperimentFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	state := openapi.EXPERIMENTSTATE_LIVE
	// create a new experiment
	experiment := &openapi.Experiment{
		Name:        experimentName,
		ExternalId:  &experimentExternalId,
		Description: &experimentDescription,
		Owner:       &experimentOwner,
		State:       &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	// create the first experiment
	_, err := service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	// attempt to create duplicate experiment
	_, err = service.UpsertExperiment(experiment)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "cannot register an experiment with duplicate names")
	suite.Equal(409, statusResp, "duplicate experiment names not allowed")
}

func (suite *CoreTestSuite) TestUpdateExperiment() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// create a new experiment
	experiment := &openapi.Experiment{
		Name:       experimentName,
		Owner:      &experimentOwner,
		ExternalId: &experimentExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	// test
	createdExperiment, err := service.UpsertExperiment(experiment)

	// checks
	suite.Nilf(err, "error creating experiment: %v", err)
	suite.NotNilf(createdExperiment.Id, "created experiment should not have nil Id")
	createdExperimentId, _ := converter.StringToInt64(createdExperiment.Id)

	// checks created experiment matches original one except for Id
	suite.Equal(experiment.Name, createdExperiment.Name, "returned experiment name should match the original one")
	suite.Equal(*experiment.ExternalId, *createdExperiment.ExternalId, "returned experiment external id should match the original one")
	suite.Equal(*experiment.CustomProperties, *createdExperiment.CustomProperties, "returned experiment custom props should match the original one")

	// update existing experiment
	newExperimentExternalId := "newExternalId"
	newOwner := "newOwner"
	newCustomProp := "updated myCustomProp"

	createdExperiment.ExternalId = &newExperimentExternalId
	createdExperiment.Owner = &newOwner
	(*createdExperiment.CustomProperties)["myCustomProp"] = openapi.MetadataValue{
		MetadataStringValue: converter.NewMetadataStringValue(newCustomProp),
	}
	// check can also define customProperty of name "owner", in addition to built-in property "owner"
	(*createdExperiment.CustomProperties)["owner"] = openapi.MetadataValue{
		MetadataStringValue: converter.NewMetadataStringValue(newCustomProp),
	}

	// update the experiment
	createdExperiment, err = service.UpsertExperiment(createdExperiment)
	suite.Nilf(err, "error updating experiment: %v", err)

	// still one experiment
	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdExperimentId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	suite.Equal(*createdExperiment.Id, *ctxId, "returned experiment id should match the mlmd one")
	suite.Equal(experimentName, *ctx.Name, "saved experiment name should match the provided one")
	suite.Equal(newExperimentExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(newOwner, ctx.Properties["owner"].GetStringValue(), "saved owner custom property should match the provided one")
	suite.Equal(newCustomProp, ctx.CustomProperties["myCustomProp"].GetStringValue(), "saved myCustomProp custom property should match the provided one")
	suite.Equal(newCustomProp, ctx.CustomProperties["owner"].GetStringValue(), "check can define custom property 'owner' and should match the provided one")

	// update the experiment keeping empty name
	newExperimentExternalId = "newNewExternalId"
	createdExperiment.ExternalId = &newExperimentExternalId
	createdExperiment.Name = ""
	createdExperiment, err = service.UpsertExperiment(createdExperiment)
	suite.Nilf(err, "error updating experiment: %v", err)

	// still one experiment
	getAllResp, err = suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err = suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdExperimentId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx = ctxById.Contexts[0]
	ctxId = converter.Int64ToString(ctx.Id)
	suite.Equal(*createdExperiment.Id, *ctxId, "returned experiment id should match the mlmd one")
	suite.Equal(experimentName, *ctx.Name, "saved experiment name should match the provided one")
	suite.Equal(newExperimentExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(newOwner, ctx.Properties["owner"].GetStringValue(), "saved owner custom property should match the provided one")
	suite.Equal(newCustomProp, ctx.CustomProperties["myCustomProp"].GetStringValue(), "saved myCustomProp custom property should match the provided one")
	suite.Equal(newCustomProp, ctx.CustomProperties["owner"].GetStringValue(), "check can define custom property 'owner' and should match the provided one")
}

func (suite *CoreTestSuite) TestGetExperimentById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	state := openapi.EXPERIMENTSTATE_LIVE
	// create a new experiment
	experiment := &openapi.Experiment{
		Name:       experimentName,
		ExternalId: &experimentExternalId,
		State:      &state,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	// test
	createdExperiment, err := service.UpsertExperiment(experiment)

	// checks
	suite.Nilf(err, "error creating experiment: %v", err)

	getExperimentById, err := service.GetExperimentById(*createdExperiment.Id)
	suite.Nilf(err, "error getting experiment by id %s: %v", *createdExperiment.Id, err)

	// checks created experiment matches original one except for Id
	suite.Equal(experiment.Name, getExperimentById.Name, "saved experiment name should match the original one")
	suite.Equal(*experiment.ExternalId, *getExperimentById.ExternalId, "saved experiment external id should match the original one")
	suite.Equal(*experiment.State, *getExperimentById.State, "saved experiment state should match the original one")
	suite.Equal(*experiment.CustomProperties, *getExperimentById.CustomProperties, "saved experiment custom props should match the original one")
}

func (suite *CoreTestSuite) TestGetExperimentByParamsWithNoResults() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	_, err := service.GetExperimentByParams(apiutils.Of("not-present"), nil)
	suite.NotNil(err)
	suite.Equal("no experiment found for provided parameters: not found", err.Error())
}

func (suite *CoreTestSuite) TestGetExperimentByParamsName() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// create a new experiment
	experiment := &openapi.Experiment{
		Name:       experimentName,
		ExternalId: &experimentExternalId,
	}

	createdExperiment, err := service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	byName, err := service.GetExperimentByParams(&experimentName, nil)
	suite.Nilf(err, "error getting experiment by name: %v", err)

	suite.Equalf(*createdExperiment.Id, *byName.Id, "the returned experiment id should match the retrieved by name")
}

func (suite *CoreTestSuite) TestGetExperimentByParamsInvalid() {
	// trigger a 400 bad request to test unallowed query params
	// create mode registry service
	service := suite.setupModelRegistryService()

	experiment := &openapi.Experiment{
		Name:       experimentName,
		ExternalId: &experimentExternalId,
	}

	// must register an experiment first, otherwise the http error will be a 404
	_, err := service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	invalidName := "\xFF"

	_, err = service.GetExperimentByParams(&invalidName, nil)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "invalid parameter used to retrieve experiment")
	suite.Equal(400, statusResp, "invalid parameter used to retrieve experiment")
}

func (suite *CoreTestSuite) TestGetExperimentByParamsExternalId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// create a new experiment
	experiment := &openapi.Experiment{
		Name:       experimentName,
		ExternalId: &experimentExternalId,
	}

	createdExperiment, err := service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	byExtId, err := service.GetExperimentByParams(nil, &experimentExternalId)
	suite.Nilf(err, "error getting experiment by external id: %v", err)

	suite.Equalf(*createdExperiment.Id, *byExtId.Id, "the returned experiment id should match the retrieved by external id")
}

func (suite *CoreTestSuite) TestGetExperimentByEmptyParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// create a new experiment
	experiment := &openapi.Experiment{
		Name:       experimentName,
		ExternalId: &experimentExternalId,
	}

	_, err := service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	_, err = service.GetExperimentByParams(nil, nil)
	suite.NotNil(err)
	suite.Equal("at least one parameter (name or externalId) must be provided", err.Error())
}

func (suite *CoreTestSuite) TestGetExperimentsOrderedById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	orderBy := "ID"

	// create a new experiment
	experiment := &openapi.Experiment{
		Name:       experimentName,
		ExternalId: &experimentExternalId,
	}

	_, err := service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	newExperimentName := "PricingExperiment2"
	newExperimentExternalId := "myExternalId2"
	experiment.Name = newExperimentName
	experiment.ExternalId = &newExperimentExternalId
	_, err = service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	newExperimentName = "PricingExperiment3"
	newExperimentExternalId = "myExternalId3"
	experiment.Name = newExperimentName
	experiment.ExternalId = &newExperimentExternalId
	_, err = service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	orderedById, err := service.GetExperiments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	suite.Nilf(err, "error getting experiments: %v", err)

	suite.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		suite.Less(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}

	orderedById, err = service.GetExperiments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	suite.Nilf(err, "error getting experiments: %v", err)

	suite.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		suite.Greater(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}
}

func (suite *CoreTestSuite) TestGetExperimentsOrderedByLastUpdate() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	orderBy := "LAST_UPDATE_TIME"

	// create a new experiment
	experiment := &openapi.Experiment{
		Name:       experimentName,
		ExternalId: &experimentExternalId,
	}

	firstExperiment, err := service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	newExperimentName := "PricingExperiment2"
	newExperimentExternalId := "myExternalId2"
	experiment.Name = newExperimentName
	experiment.ExternalId = &newExperimentExternalId
	secondExperiment, err := service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	newExperimentName = "PricingExperiment3"
	newExperimentExternalId = "myExternalId3"
	experiment.Name = newExperimentName
	experiment.ExternalId = &newExperimentExternalId
	thirdExperiment, err := service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	// update second experiment
	secondExperiment.ExternalId = nil
	_, err = service.UpsertExperiment(secondExperiment)
	suite.Nilf(err, "error updating experiment: %v", err)

	orderedById, err := service.GetExperiments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	suite.Nilf(err, "error getting experiments: %v", err)

	suite.Equal(3, int(orderedById.Size))
	suite.Equal(*firstExperiment.Id, *orderedById.Items[0].Id)
	suite.Equal(*thirdExperiment.Id, *orderedById.Items[1].Id)
	suite.Equal(*secondExperiment.Id, *orderedById.Items[2].Id)

	orderedById, err = service.GetExperiments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	suite.Nilf(err, "error getting experiments: %v", err)

	suite.Equal(3, int(orderedById.Size))
	suite.Equal(*secondExperiment.Id, *orderedById.Items[0].Id)
	suite.Equal(*thirdExperiment.Id, *orderedById.Items[1].Id)
	suite.Equal(*firstExperiment.Id, *orderedById.Items[2].Id)
}

func (suite *CoreTestSuite) TestGetExperimentsWithPageSize() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	pageSize := int32(1)
	pageSize2 := int32(2)
	experimentName := "PricingExperiment1"
	experimentExternalId := "myExternalId1"

	// create a new experiment
	experiment := &openapi.Experiment{
		Name:       experimentName,
		ExternalId: &experimentExternalId,
	}

	firstExperiment, err := service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	newExperimentName := "PricingExperiment2"
	newExperimentExternalId := "myExternalId2"
	experiment.Name = newExperimentName
	experiment.ExternalId = &newExperimentExternalId
	secondExperiment, err := service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	newExperimentName = "PricingExperiment3"
	newExperimentExternalId = "myExternalId3"
	experiment.Name = newExperimentName
	experiment.ExternalId = &newExperimentExternalId
	thirdExperiment, err := service.UpsertExperiment(experiment)
	suite.Nilf(err, "error creating experiment: %v", err)

	truncatedList, err := service.GetExperiments(api.ListOptions{
		PageSize: &pageSize,
	})
	suite.Nilf(err, "error getting experiments: %v", err)

	suite.Equal(1, int(truncatedList.Size))
	suite.NotEqual("", truncatedList.NextPageToken, "next page token should not be empty")
	suite.Equal(*firstExperiment.Id, *truncatedList.Items[0].Id)

	truncatedList, err = service.GetExperiments(api.ListOptions{
		PageSize:      &pageSize2,
		NextPageToken: &truncatedList.NextPageToken,
	})
	suite.Nilf(err, "error getting experiments: %v", err)

	suite.Equal(2, int(truncatedList.Size))
	suite.Equal("", truncatedList.NextPageToken, "next page token should be empty as list item returned")
	suite.Equal(*secondExperiment.Id, *truncatedList.Items[0].Id)
	suite.Equal(*thirdExperiment.Id, *truncatedList.Items[1].Id)
}

func (suite *CoreTestSuite) TestCreateExperimentWithCustomPropFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// create a new experiment with incomplete customProperty fields
	experiment := &openapi.Experiment{
		Name: experimentName,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp1": {},
		},
	}

	// test
	_, err := service.UpsertExperiment(experiment)

	// checks
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "error creating experiment: %v", err)
	suite.Equal(400, statusResp, "customProperties must include metadataType")
}
