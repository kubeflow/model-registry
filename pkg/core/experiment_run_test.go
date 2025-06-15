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

// EXPERIMENT RUNS

func (suite *CoreTestSuite) TestCreateExperimentRun() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentId := suite.registerExperiment(service, nil, nil)

	state := openapi.EXPERIMENTRUNSTATE_LIVE
	status := openapi.EXPERIMENTRUNSTATUS_RUNNING
	experimentRun := &openapi.ExperimentRun{
		Name:        &experimentRunName,
		ExternalId:  &experimentRunExternalId,
		Description: &experimentRunDescription,
		State:       &state,
		Status:      &status,
		Owner:       &experimentRunOwner,
	}

	createdExperimentRun, err := service.UpsertExperimentRun(experimentRun, &experimentId)
	suite.Nilf(err, "error creating new experiment run for %s", experimentId)
	suite.Equal(createdExperimentRun.ExperimentId, experimentId, "ExperimentId should match the actual parent experiment")

	suite.NotNilf(createdExperimentRun.Id, "created experiment run should not have nil Id")

	createdExperimentRunId, _ := converter.StringToInt64(createdExperimentRun.Id)

	byId, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdExperimentRunId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	suite.Equal(*createdExperimentRunId, *byId.Contexts[0].Id, "returned experiment run id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", experimentId, experimentRunName), *byId.Contexts[0].Name, "saved experiment run name should match the provided one")
	suite.Equal(experimentRunExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(experimentRunOwner, byId.Contexts[0].Properties["owner"].GetStringValue(), "saved owner property should match the provided one")
	suite.Equal(experimentRunDescription, byId.Contexts[0].Properties["description"].GetStringValue(), "saved description should match the provided one")
	suite.Equal(string(state), byId.Contexts[0].Properties["state"].GetStringValue(), "saved state should match the provided one")
	suite.Equal(string(status), byId.Contexts[0].Properties["status"].GetStringValue(), "saved status should match the provided one")
	suite.Equalf(*experimentRunTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *experimentRunTypeName)

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(2, len(getAllResp.Contexts), "there should be two contexts saved in mlmd (experiment + experiment run)")
}

func (suite *CoreTestSuite) TestCreateDuplicateExperimentRunFailure() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentId := suite.registerExperiment(service, nil, nil)

	state := openapi.EXPERIMENTRUNSTATE_LIVE
	status := openapi.EXPERIMENTRUNSTATUS_RUNNING
	experimentRun := &openapi.ExperimentRun{
		Name:        &experimentRunName,
		ExternalId:  &experimentRunExternalId,
		Description: &experimentRunDescription,
		State:       &state,
		Status:      &status,
		Owner:       &experimentRunOwner,
	}

	createdExperimentRun, err := service.UpsertExperimentRun(experimentRun, &experimentId)
	suite.Nilf(err, "error creating new experiment run for %s", experimentId)
	suite.Equal(createdExperimentRun.ExperimentId, experimentId, "ExperimentId should match the actual parent experiment")

	// attempt to create duplicate experiment run
	_, err = service.UpsertExperimentRun(experimentRun, &experimentId)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "cannot register a duplicate experiment run")
	suite.Equal(409, statusResp, "duplicate experiment runs not allowed")
}

func (suite *CoreTestSuite) TestCreateExperimentRunFailure() {
	// create model registry service
	service := suite.setupModelRegistryService()

	_, err := service.UpsertExperimentRun(nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid experiment run pointer, can't upsert nil: bad request", err.Error())

	experimentId := "9999"

	experimentRun := &openapi.ExperimentRun{
		Name:       &experimentRunName,
		ExternalId: &experimentRunExternalId,
		Owner:      &experimentRunOwner,
	}

	_, err = service.UpsertExperimentRun(experimentRun, nil)
	suite.NotNil(err)
	suite.Equal("missing experiment id, cannot create experiment run without experiment: bad request", err.Error())

	_, err = service.UpsertExperimentRun(experimentRun, &experimentId)
	suite.NotNil(err)
	suite.Equal("no experiment found for id 9999: not found", err.Error())
}

func (suite *CoreTestSuite) TestUpdateExperimentRun() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentId := suite.registerExperiment(service, nil, nil)

	experimentRun := &openapi.ExperimentRun{
		Name:       &experimentRunName,
		ExternalId: &experimentRunExternalId,
		Owner:      &experimentRunOwner,
	}

	createdExperimentRun, err := service.UpsertExperimentRun(experimentRun, &experimentId)
	suite.Nilf(err, "error creating new experiment run for %s", experimentId)

	suite.NotNilf(createdExperimentRun.Id, "created experiment run should not have nil Id")
	createdExperimentRunId, _ := converter.StringToInt64(createdExperimentRun.Id)

	newExternalId := "org.myawesomeexperiment.run1.updated"
	newStatus := openapi.EXPERIMENTRUNSTATUS_FINISHED

	createdExperimentRun.ExternalId = &newExternalId
	createdExperimentRun.Status = &newStatus
	(*createdExperimentRun.CustomProperties)["score"] = openapi.MetadataValue{
		MetadataDoubleValue: converter.NewMetadataDoubleValue(0.95),
	}

	updatedExperimentRun, err := service.UpsertExperimentRun(createdExperimentRun, &experimentId)
	suite.Nilf(err, "error updating experiment run for %s: %v", experimentId, err)
	suite.Equal(updatedExperimentRun.ExperimentId, experimentId, "ExperimentId should match the actual parent experiment")

	updateExperimentRunId, _ := converter.StringToInt64(updatedExperimentRun.Id)
	suite.Equal(*createdExperimentRunId, *updateExperimentRunId, "created and updated experiment run should have same id")

	byId, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateExperimentRunId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	suite.Equal(*updateExperimentRunId, *byId.Contexts[0].Id, "returned experiment run id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", experimentId, experimentRunName), *byId.Contexts[0].Name, "saved experiment run name should match the provided one")
	suite.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(experimentRunOwner, byId.Contexts[0].Properties["owner"].GetStringValue(), "saved owner property should match the provided one")
	suite.Equal(string(newStatus), byId.Contexts[0].Properties["status"].GetStringValue(), "saved status should match the provided one")
	suite.Equal(0.95, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	suite.Equalf(*experimentRunTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *experimentRunTypeName)

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(2, len(getAllResp.Contexts), "there should be two contexts saved in mlmd")

	// update with nil name
	newExternalId = "org.myawesomeexperiment.run1.updated2"
	updatedExperimentRun.ExternalId = &newExternalId
	updatedExperimentRun.Name = nil
	updatedExperimentRun, err = service.UpsertExperimentRun(updatedExperimentRun, &experimentId)
	suite.Nilf(err, "error updating experiment run for %s: %v", experimentId, err)

	updateExperimentRunId, _ = converter.StringToInt64(updatedExperimentRun.Id)
	suite.Equal(*createdExperimentRunId, *updateExperimentRunId, "created and updated experiment run should have same id")

	byId, err = suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*updateExperimentRunId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)
	suite.Equal(1, len(byId.Contexts), "there should be just one context saved in mlmd")

	suite.Equal(*updateExperimentRunId, *byId.Contexts[0].Id, "returned experiment run id should match the mlmd one")
	suite.Equal(fmt.Sprintf("%s:%s", experimentId, experimentRunName), *byId.Contexts[0].Name, "saved experiment run name should match the provided one")
	suite.Equal(newExternalId, *byId.Contexts[0].ExternalId, "saved external id should match the provided one")
	suite.Equal(experimentRunOwner, byId.Contexts[0].Properties["owner"].GetStringValue(), "saved owner property should match the provided one")
	suite.Equal(string(newStatus), byId.Contexts[0].Properties["status"].GetStringValue(), "saved status should match the provided one")
	suite.Equal(0.95, byId.Contexts[0].CustomProperties["score"].GetDoubleValue(), "saved score custom property should match the provided one")
	suite.Equalf(*experimentRunTypeName, *byId.Contexts[0].Type, "saved context should be of type of %s", *experimentRunTypeName)
}

func (suite *CoreTestSuite) TestUpdateExperimentRunFailure() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentId := suite.registerExperiment(service, nil, nil)

	experimentRun := &openapi.ExperimentRun{
		Name:       &experimentRunName,
		ExternalId: &experimentRunExternalId,
		Owner:      &experimentRunOwner,
	}

	createdExperimentRun, err := service.UpsertExperimentRun(experimentRun, &experimentId)
	suite.Nilf(err, "error creating new experiment run for %s", experimentId)
	suite.NotNilf(createdExperimentRun.Id, "created experiment run should not have nil Id")

	newExternalId := "org.myawesomeexperiment.run1.updated"
	newStatus := openapi.EXPERIMENTRUNSTATUS_FINISHED

	createdExperimentRun.ExternalId = &newExternalId
	createdExperimentRun.Status = &newStatus
	(*createdExperimentRun.CustomProperties)["score"] = openapi.MetadataValue{
		MetadataDoubleValue: converter.NewMetadataDoubleValue(0.95),
	}

	wrongId := "9999"
	createdExperimentRun.Id = &wrongId
	_, err = service.UpsertExperimentRun(createdExperimentRun, &experimentId)
	suite.NotNil(err)
	suite.Equal(fmt.Sprintf("no experiment run found for id %s: not found", wrongId), err.Error())
}

func (suite *CoreTestSuite) TestGetExperimentRunById() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentId := suite.registerExperiment(service, nil, nil)

	state := openapi.EXPERIMENTRUNSTATE_ARCHIVED
	status := openapi.EXPERIMENTRUNSTATUS_FINISHED
	experimentRun := &openapi.ExperimentRun{
		Name:       &experimentRunName,
		ExternalId: &experimentRunExternalId,
		State:      &state,
		Status:     &status,
		Owner:      &experimentRunOwner,
	}

	createdExperimentRun, err := service.UpsertExperimentRun(experimentRun, &experimentId)
	suite.Nilf(err, "error creating new experiment run for %s", experimentId)

	suite.NotNilf(createdExperimentRun.Id, "created experiment run should not have nil Id")
	createdExperimentRunId, _ := converter.StringToInt64(createdExperimentRun.Id)

	getById, err := service.GetExperimentRunById(*createdExperimentRun.Id)
	suite.Nilf(err, "error getting experiment run with id %s", *createdExperimentRunId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdExperimentRunId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getById.Id, "returned experiment run id should match the mlmd context one")
	suite.Equal(*experimentRun.Name, *getById.Name, "saved experiment run name should match the provided one")
	suite.Equal(*experimentRun.ExternalId, *getById.ExternalId, "saved external id should match the provided one")
	suite.Equal(*experimentRun.State, *getById.State, "saved experiment run state should match the original one")
	suite.Equal(*experimentRun.Status, *getById.Status, "saved experiment run status should match the original one")
	suite.Equal(*getById.Owner, experimentRunOwner, "saved owner property should match the provided one")
}

func (suite *CoreTestSuite) TestGetExperimentRunByParamsWithNoResults() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentId := suite.registerExperiment(service, nil, nil)

	_, err := service.GetExperimentRunByParams(apiutils.Of("not-present"), &experimentId, nil)
	suite.NotNil(err)
	suite.Equal("no experiment run found for provided parameters: not found", err.Error())
}

func (suite *CoreTestSuite) TestGetExperimentRunByParamsName() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentId := suite.registerExperiment(service, nil, nil)

	experimentRun := &openapi.ExperimentRun{
		Name:       &experimentRunName,
		ExternalId: &experimentRunExternalId,
		Owner:      &experimentRunOwner,
	}

	createdExperimentRun, err := service.UpsertExperimentRun(experimentRun, &experimentId)
	suite.Nilf(err, "error creating new experiment run for %s", experimentId)

	suite.NotNilf(createdExperimentRun.Id, "created experiment run should not have nil Id")
	createdExperimentRunId, _ := converter.StringToInt64(createdExperimentRun.Id)

	getByName, err := service.GetExperimentRunByParams(&experimentRunName, &experimentId, nil)
	suite.Nilf(err, "error getting experiment run by name %s", *createdExperimentRunId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdExperimentRunId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getByName.Id, "returned experiment run id should match the mlmd context one")
	suite.Equal(fmt.Sprintf("%s:%s", experimentId, *getByName.Name), *ctx.Name, "saved experiment run name should match the provided one")
	suite.Equal(*ctx.ExternalId, *getByName.ExternalId, "saved external id should match the provided one")
	suite.Equal(ctx.Properties["owner"].GetStringValue(), *getByName.Owner, "saved owner property should match the provided one")
}

func (suite *CoreTestSuite) TestGetExperimentRunByParamsInvalid() {
	// trigger a 400 bad request to test unallowed query params
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentId := suite.registerExperiment(service, nil, nil)

	experimentRun := &openapi.ExperimentRun{
		Name:       &experimentRunName,
		ExternalId: &experimentRunExternalId,
		Owner:      &experimentRunOwner,
	}

	// must register an experiment run first, otherwise the http error will be a 404
	_, err := service.UpsertExperimentRun(experimentRun, &experimentId)
	suite.Nilf(err, "error creating new experiment run for %s", experimentId)

	invalidName := "\xFF"

	_, err = service.GetExperimentRunByParams(&invalidName, &experimentId, nil)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "invalid parameter used to retrieve experiment run")
	suite.Equal(400, statusResp, "invalid parameter used to retrieve experiment run")
}

func (suite *CoreTestSuite) TestGetExperimentRunByParamsExternalId() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentId := suite.registerExperiment(service, nil, nil)

	experimentRun := &openapi.ExperimentRun{
		Name:       &experimentRunName,
		ExternalId: &experimentRunExternalId,
		Owner:      &experimentRunOwner,
	}

	createdExperimentRun, err := service.UpsertExperimentRun(experimentRun, &experimentId)
	suite.Nilf(err, "error creating new experiment run for %s", experimentId)

	suite.NotNilf(createdExperimentRun.Id, "created experiment run should not have nil Id")
	createdExperimentRunId, _ := converter.StringToInt64(createdExperimentRun.Id)

	getByExternalId, err := service.GetExperimentRunByParams(nil, nil, experimentRun.ExternalId)
	suite.Nilf(err, "error getting experiment run by external id %s", *experimentRun.ExternalId)

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{
			*createdExperimentRunId,
		},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	suite.Equal(*converter.Int64ToString(ctx.Id), *getByExternalId.Id, "returned experiment run id should match the mlmd context one")
	suite.Equal(fmt.Sprintf("%s:%s", experimentId, *getByExternalId.Name), *ctx.Name, "saved experiment run name should match the provided one")
	suite.Equal(*ctx.ExternalId, *getByExternalId.ExternalId, "saved external id should match the provided one")
	suite.Equal(ctx.Properties["owner"].GetStringValue(), *getByExternalId.Owner, "saved owner property should match the provided one")
}

func (suite *CoreTestSuite) TestGetExperimentRunByEmptyParams() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentId := suite.registerExperiment(service, nil, nil)

	experimentRun := &openapi.ExperimentRun{
		Name:       &experimentRunName,
		ExternalId: &experimentRunExternalId,
		Owner:      &experimentRunOwner,
	}

	createdExperimentRun, err := service.UpsertExperimentRun(experimentRun, &experimentId)
	suite.Nilf(err, "error creating new experiment run for %s", experimentId)
	suite.NotNilf(createdExperimentRun.Id, "created experiment run should not have nil Id")

	_, err = service.GetExperimentRunByParams(nil, nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid parameters call, supply either (name and experimentId), or externalId: bad request", err.Error())
}

func (suite *CoreTestSuite) TestGetExperimentRuns() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentId := suite.registerExperiment(service, nil, nil)

	experimentRun1 := &openapi.ExperimentRun{
		Name:       &experimentRunName,
		ExternalId: &experimentRunExternalId,
	}

	secondExperimentRunName := "run2"
	secondExperimentRunExtId := "org.myawesomeexperiment.run2"
	experimentRun2 := &openapi.ExperimentRun{
		Name:       &secondExperimentRunName,
		ExternalId: &secondExperimentRunExtId,
	}

	thirdExperimentRunName := "run3"
	thirdExperimentRunExtId := "org.myawesomeexperiment.run3"
	experimentRun3 := &openapi.ExperimentRun{
		Name:       &thirdExperimentRunName,
		ExternalId: &thirdExperimentRunExtId,
	}

	createdExperimentRun1, err := service.UpsertExperimentRun(experimentRun1, &experimentId)
	suite.Nilf(err, "error creating new experiment run for %s", experimentId)

	createdExperimentRun2, err := service.UpsertExperimentRun(experimentRun2, &experimentId)
	suite.Nilf(err, "error creating new experiment run for %s", experimentId)

	createdExperimentRun3, err := service.UpsertExperimentRun(experimentRun3, &experimentId)
	suite.Nilf(err, "error creating new experiment run for %s", experimentId)

	anotherExperimentName := "AnotherExperiment"
	anotherExperimentExtId := "org.another"
	anotherExperimentId := suite.registerExperiment(service, &anotherExperimentName, &anotherExperimentExtId)

	anotherExperimentRunName := "another-run"
	anotherExperimentRunExtId := "org.another.run1"
	experimentRunAnother := &openapi.ExperimentRun{
		Name:       &anotherExperimentRunName,
		ExternalId: &anotherExperimentRunExtId,
	}

	_, err = service.UpsertExperimentRun(experimentRunAnother, &anotherExperimentId)
	suite.Nilf(err, "error creating new experiment run for %s", anotherExperimentId)

	createdExperimentRunId1, _ := converter.StringToInt64(createdExperimentRun1.Id)
	createdExperimentRunId2, _ := converter.StringToInt64(createdExperimentRun2.Id)
	createdExperimentRunId3, _ := converter.StringToInt64(createdExperimentRun3.Id)

	getAll, err := service.GetExperimentRuns(api.ListOptions{}, nil)
	suite.Nilf(err, "error getting all experiment runs")
	suite.Equal(int32(4), getAll.Size, "expected four experiment runs across all experiments")

	getAllByExperiment, err := service.GetExperimentRuns(api.ListOptions{}, &experimentId)
	suite.Nilf(err, "error getting all experiment runs")
	suite.Equalf(int32(3), getAllByExperiment.Size, "expected three experiment runs for experiment %s", experimentId)

	suite.Equal(*converter.Int64ToString(createdExperimentRunId1), *getAllByExperiment.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdExperimentRunId2), *getAllByExperiment.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdExperimentRunId3), *getAllByExperiment.Items[2].Id)

	// order by last update time, expecting last created as first
	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByExperiment, err = service.GetExperimentRuns(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &experimentId)
	suite.Nilf(err, "error getting all experiment runs")
	suite.Equalf(int32(3), getAllByExperiment.Size, "expected three experiment runs for experiment %s", experimentId)

	suite.Equal(*converter.Int64ToString(createdExperimentRunId1), *getAllByExperiment.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdExperimentRunId2), *getAllByExperiment.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdExperimentRunId3), *getAllByExperiment.Items[0].Id)

	// update the second experiment run
	newExperimentRunExternalId := "updated.org:run2"
	createdExperimentRun2.ExternalId = &newExperimentRunExternalId
	createdExperimentRun2, err = service.UpsertExperimentRun(createdExperimentRun2, &experimentId)
	suite.Nilf(err, "error updating experiment run for %s", experimentId)

	suite.Equal(newExperimentRunExternalId, *createdExperimentRun2.ExternalId)

	getAllByExperiment, err = service.GetExperimentRuns(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &experimentId)
	suite.Nilf(err, "error getting all experiment runs")
	suite.Equalf(int32(3), getAllByExperiment.Size, "expected three experiment runs for experiment %s", experimentId)

	suite.Equal(*converter.Int64ToString(createdExperimentRunId1), *getAllByExperiment.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdExperimentRunId2), *getAllByExperiment.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdExperimentRunId3), *getAllByExperiment.Items[1].Id)
}

func (suite *CoreTestSuite) TestCreateExperimentRunWithCustomPropFailure() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentId := suite.registerExperiment(service, nil, nil)

	// create an experiment run with incomplete customProperty fields
	experimentRun := &openapi.ExperimentRun{
		Name: &experimentRunName,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp1": {},
		},
	}

	// test
	_, err := service.UpsertExperimentRun(experimentRun, &experimentId)

	// checks
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "error creating experiment run: %v", err)
	suite.Equal(400, statusResp, "customProperties must include metadataType")
}

// EXPERIMENT RUN ARTIFACTS

func (suite *CoreTestSuite) TestUpsertExperimentRunArtifact() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentRunId := suite.registerExperimentRun(service, nil, nil, nil, nil)

	createdArt, err := service.UpsertExperimentRunArtifact(&openapi.Artifact{
		DocArtifact: &openapi.DocArtifact{
			Name:        &artifactName,
			State:       (*openapi.ArtifactState)(&artifactState),
			Uri:         &artifactUri,
			Description: &artifactDescription,
			CustomProperties: &map[string]openapi.MetadataValue{
				"custom_string_prop": {
					MetadataStringValue: converter.NewMetadataStringValue(customString),
				},
			},
		},
	}, experimentRunId)
	suite.Nilf(err, "error creating new artifact: %v", err)

	docArtifact := createdArt.DocArtifact
	suite.NotNil(docArtifact, "error creating new artifact")
	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	suite.NotNil(docArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *docArtifact.Name)
	suite.Equal(*state, *docArtifact.State)
	suite.Equal(artifactUri, *docArtifact.Uri)
	suite.Equal(artifactDescription, *docArtifact.Description)
	suite.Equal(customString, (*docArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)
}

func (suite *CoreTestSuite) TestGetExperimentRunArtifacts() {
	// create model registry service
	service := suite.setupModelRegistryService()

	experimentRunId := suite.registerExperimentRun(service, nil, nil, nil, nil)

	secondArtifactName := "second-name"
	secondArtifactExtId := "second-ext-id"
	secondArtifactUri := "second-uri"

	createdArtifact1, err := service.UpsertExperimentRunArtifact(&openapi.Artifact{
		ModelArtifact: &openapi.ModelArtifact{
			Name:       &artifactName,
			State:      (*openapi.ArtifactState)(&artifactState),
			Uri:        &artifactUri,
			ExternalId: &artifactExtId,
			CustomProperties: &map[string]openapi.MetadataValue{
				"custom_string_prop": {
					MetadataStringValue: converter.NewMetadataStringValue(customString),
				},
			},
		},
	}, experimentRunId)
	suite.Nilf(err, "error creating new artifact: %v", err)
	createdArtifact2, err := service.UpsertExperimentRunArtifact(&openapi.Artifact{
		DocArtifact: &openapi.DocArtifact{
			Name:       &secondArtifactName,
			State:      (*openapi.ArtifactState)(&artifactState),
			Uri:        &secondArtifactUri,
			ExternalId: &secondArtifactExtId,
			CustomProperties: &map[string]openapi.MetadataValue{
				"custom_string_prop": {
					MetadataStringValue: converter.NewMetadataStringValue(customString),
				},
			},
		},
	}, experimentRunId)
	suite.Nilf(err, "error creating new artifact: %v", err)

	createdArtifactId1, _ := converter.StringToInt64(createdArtifact1.ModelArtifact.Id)
	createdArtifactId2, _ := converter.StringToInt64(createdArtifact2.DocArtifact.Id)

	getAll, err := service.GetExperimentRunArtifacts(api.ListOptions{}, &experimentRunId)
	suite.Nilf(err, "error getting all experiment run artifacts")
	suite.Equalf(int32(2), getAll.Size, "expected two artifacts")

	suite.Equal(*converter.Int64ToString(createdArtifactId1), *getAll.Items[0].ModelArtifact.Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId2), *getAll.Items[1].DocArtifact.Id)

	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByExperimentRun, err := service.GetExperimentRunArtifacts(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &experimentRunId)
	suite.Nilf(err, "error getting all experiment run artifacts: %v", err)
	suite.Equalf(int32(2), getAllByExperimentRun.Size, "expected 2 artifacts for experiment run %v", experimentRunId)

	suite.Equal(*converter.Int64ToString(createdArtifactId1), *getAllByExperimentRun.Items[1].ModelArtifact.Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId2), *getAllByExperimentRun.Items[0].DocArtifact.Id)
}
