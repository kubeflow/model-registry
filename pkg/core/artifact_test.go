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

// ARTIFACTS

func (suite *CoreTestSuite) TestCreateArtifact() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	createdArt, err := service.UpsertArtifact(&openapi.Artifact{
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
	}, &modelVersionId)
	suite.Nilf(err, "error creating new artifact for %d: %v", modelVersionId, err)

	docArtifact := createdArt.DocArtifact
	suite.NotNilf(docArtifact, "error creating new artifact for %d", modelVersionId)
	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	suite.NotNil(docArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *docArtifact.Name)
	suite.Equal(*state, *docArtifact.State)
	suite.Equal(artifactUri, *docArtifact.Uri)
	suite.Equal(artifactDescription, *docArtifact.Description)
	suite.Equal(customString, (*docArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)
}

func (suite *CoreTestSuite) TestCreateArtifactFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := "9998"

	var artifact openapi.Artifact
	artifact.DocArtifact = &openapi.DocArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	_, err := service.UpsertArtifact(&artifact, nil)
	suite.NotNil(err)
	suite.Equal("missing model version id, cannot create artifact without model version: bad request", err.Error())

	_, err = service.UpsertArtifact(&artifact, &modelVersionId)
	suite.NotNil(err)
	suite.Equal("no model version found for id 9998: not found", err.Error())
}

func (suite *CoreTestSuite) TestUpdateArtifact() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	createdArtifact, err := service.UpsertArtifact(&openapi.Artifact{
		DocArtifact: &openapi.DocArtifact{
			Name:  &artifactName,
			State: (*openapi.ArtifactState)(&artifactState),
			Uri:   &artifactUri,
			CustomProperties: &map[string]openapi.MetadataValue{
				"custom_string_prop": {
					MetadataStringValue: converter.NewMetadataStringValue(customString),
				},
			},
		},
	}, &modelVersionId)
	suite.Nilf(err, "error creating new artifact for %d", modelVersionId)

	newState := "MARKED_FOR_DELETION"
	createdArtifact.DocArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertArtifact(createdArtifact, &modelVersionId)
	suite.Nilf(err, "error updating artifact for %d: %v", modelVersionId, err)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.DocArtifact.Id)
	updatedArtifactId, _ := converter.StringToInt64(updatedArtifact.DocArtifact.Id)
	suite.Equal(createdArtifactId, updatedArtifactId)

	getById, err := suite.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	suite.Nilf(err, "error getting artifact by id %d", createdArtifactId)

	suite.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	suite.Equal(fmt.Sprintf("%s:%s", modelVersionId, *createdArtifact.DocArtifact.Name), *getById.Artifacts[0].Name)
	suite.Equal(string(newState), getById.Artifacts[0].State.String())
	suite.Equal(*createdArtifact.DocArtifact.Uri, *getById.Artifacts[0].Uri)
	suite.Equal((*createdArtifact.DocArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["custom_string_prop"].GetStringValue())
}

func (suite *CoreTestSuite) TestUpdateArtifactFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	createdArtifact, err := service.UpsertArtifact(&openapi.Artifact{
		DocArtifact: &openapi.DocArtifact{
			Name:  &artifactName,
			State: (*openapi.ArtifactState)(&artifactState),
			Uri:   &artifactUri,
			CustomProperties: &map[string]openapi.MetadataValue{
				"custom_string_prop": {
					MetadataStringValue: converter.NewMetadataStringValue(customString),
				},
			},
		},
	}, &modelVersionId)
	suite.Nilf(err, "error creating new artifact for model version %s", modelVersionId)
	suite.NotNilf(createdArtifact.DocArtifact.Id, "created model artifact should not have nil Id")

	newState := "MARKED_FOR_DELETION"
	createdArtifact.DocArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertArtifact(createdArtifact, &modelVersionId)
	suite.Nilf(err, "error updating artifact for %d: %v", modelVersionId, err)

	wrongId := "5555"
	updatedArtifact.DocArtifact.Id = &wrongId
	_, err = service.UpsertArtifact(updatedArtifact, &modelVersionId)
	suite.NotNil(err)
	suite.Equal(fmt.Sprintf("no artifact found for id %s: not found", wrongId), err.Error())
}

func (suite *CoreTestSuite) TestGetArtifactById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	createdArtifact, err := service.UpsertArtifact(&openapi.Artifact{
		DocArtifact: &openapi.DocArtifact{
			Name:  &artifactName,
			State: (*openapi.ArtifactState)(&artifactState),
			Uri:   &artifactUri,
			CustomProperties: &map[string]openapi.MetadataValue{
				"custom_string_prop": {
					MetadataStringValue: converter.NewMetadataStringValue(customString),
				},
			},
		},
	}, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.DocArtifact.Id)

	getById, err := service.GetArtifactById(*createdArtifact.DocArtifact.Id)
	suite.Nilf(err, "error getting artifact by id %d", createdArtifactId)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	suite.NotNil(createdArtifact.DocArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *getById.DocArtifact.Name)
	suite.Equal(*state, *getById.DocArtifact.State)
	suite.Equal(artifactUri, *getById.DocArtifact.Uri)
	suite.Equal(customString, (*getById.DocArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*createdArtifact, *getById, "artifacts returned during creation and on get by id should be equal")
}

func (suite *CoreTestSuite) TestGetArtifacts() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	secondArtifactName := "second-name"
	secondArtifactExtId := "second-ext-id"
	secondArtifactUri := "second-uri"

	createdArtifact1, err := service.UpsertArtifact(&openapi.Artifact{
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
	}, &modelVersionId)
	suite.Nilf(err, "error creating new artifact for %d", modelVersionId)
	createdArtifact2, err := service.UpsertArtifact(&openapi.Artifact{
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
	}, &modelVersionId)
	suite.Nilf(err, "error creating new artifact for %d", modelVersionId)

	createdArtifactId1, _ := converter.StringToInt64(createdArtifact1.ModelArtifact.Id)
	createdArtifactId2, _ := converter.StringToInt64(createdArtifact2.DocArtifact.Id)

	getAll, err := service.GetArtifacts(api.ListOptions{}, &modelVersionId)
	suite.Nilf(err, "error getting all model artifacts")
	suite.Equalf(int32(2), getAll.Size, "expected two artifacts")

	suite.Equal(*converter.Int64ToString(createdArtifactId1), *getAll.Items[0].ModelArtifact.Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId2), *getAll.Items[1].DocArtifact.Id)

	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByModelVersion, err := service.GetArtifacts(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &modelVersionId)
	suite.Nilf(err, "error getting all model artifacts for %d", modelVersionId)
	suite.Equalf(int32(2), getAllByModelVersion.Size, "expected 2 artifacts for model version %d", modelVersionId)

	suite.Equal(*converter.Int64ToString(createdArtifactId1), *getAllByModelVersion.Items[1].ModelArtifact.Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId2), *getAllByModelVersion.Items[0].DocArtifact.Id)
}

// MODEL ARTIFACTS

func (suite *CoreTestSuite) TestCreateModelArtifact() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact, err := service.UpsertModelArtifact(&openapi.ModelArtifact{
		Name:               &artifactName,
		State:              (*openapi.ArtifactState)(&artifactState),
		Uri:                &artifactUri,
		Description:        &artifactDescription,
		ModelFormatName:    apiutils.Of("onnx"),
		ModelFormatVersion: apiutils.Of("1"),
		StorageKey:         apiutils.Of("aws-connection-models"),
		StoragePath:        apiutils.Of("bucket"),
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	suite.NotNil(modelArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *modelArtifact.Name)
	suite.Equal(*state, *modelArtifact.State)
	suite.Equal(artifactUri, *modelArtifact.Uri)
	suite.Equal(artifactDescription, *modelArtifact.Description)
	suite.Equal("onnx", *modelArtifact.ModelFormatName)
	suite.Equal("1", *modelArtifact.ModelFormatVersion)
	suite.Equal("aws-connection-models", *modelArtifact.StorageKey)
	suite.Equal("bucket", *modelArtifact.StoragePath)
	suite.Equal(customString, (*modelArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)
}

func (suite *CoreTestSuite) TestCreateModelArtifactFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := "9998"

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	_, err := service.UpsertModelArtifact(modelArtifact, nil)
	suite.NotNil(err)
	suite.Equal("missing model version id, cannot create artifact without model version: bad request", err.Error())

	_, err = service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.NotNil(err)
	suite.Equal("no model version found for id 9998: not found", err.Error())
}

func (suite *CoreTestSuite) TestUpdateModelArtifact() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	newState := "MARKED_FOR_DELETION"
	createdArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertModelArtifact(createdArtifact, &modelVersionId)
	suite.Nilf(err, "error updating model artifact for %d: %v", modelVersionId, err)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)
	updatedArtifactId, _ := converter.StringToInt64(updatedArtifact.Id)
	suite.Equal(createdArtifactId, updatedArtifactId)

	getById, err := suite.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	suite.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	suite.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	suite.Equal(fmt.Sprintf("%s:%s", modelVersionId, *createdArtifact.Name), *getById.Artifacts[0].Name)
	suite.Equal(string(newState), getById.Artifacts[0].State.String())
	suite.Equal(*createdArtifact.Uri, *getById.Artifacts[0].Uri)
	suite.Equal((*createdArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["custom_string_prop"].GetStringValue())
}

func (suite *CoreTestSuite) TestUpdateModelArtifactFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for model version %s", modelVersionId)
	suite.NotNilf(createdArtifact.Id, "created model artifact should not have nil Id")
}

func (suite *CoreTestSuite) TestGetModelArtifactById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:  &artifactName,
		State: (*openapi.ArtifactState)(&artifactState),
		Uri:   &artifactUri,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)

	getById, err := service.GetModelArtifactById(*createdArtifact.Id)
	suite.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	suite.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *getById.Name)
	suite.Equal(*state, *getById.State)
	suite.Equal(artifactUri, *getById.Uri)
	suite.Equal(customString, (*getById.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*createdArtifact, *getById, "artifacts returned during creation and on get by id should be equal")
}

func (suite *CoreTestSuite) TestGetModelArtifactByParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:       &artifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &artifactUri,
		ExternalId: &artifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)

	getByName, err := service.GetModelArtifactByParams(&artifactName, &modelVersionId, nil)
	suite.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	suite.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *getByName.Name)
	suite.Equal(artifactExtId, *getByName.ExternalId)
	suite.Equal(*state, *getByName.State)
	suite.Equal(artifactUri, *getByName.Uri)
	suite.Equal(customString, (*getByName.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*createdArtifact, *getByName, "artifacts returned during creation and on get by name should be equal")

	getByExtId, err := service.GetModelArtifactByParams(nil, nil, &artifactExtId)
	suite.Nilf(err, "error getting model artifact by id %d", createdArtifactId)

	suite.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *getByExtId.Name)
	suite.Equal(artifactExtId, *getByExtId.ExternalId)
	suite.Equal(*state, *getByExtId.State)
	suite.Equal(artifactUri, *getByExtId.Uri)
	suite.Equal(customString, (*getByExtId.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*createdArtifact, *getByExtId, "artifacts returned during creation and on get by ext id should be equal")
}

func (suite *CoreTestSuite) TestGetModelArtifactByEmptyParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact := &openapi.ModelArtifact{
		Name:       &artifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &artifactUri,
		ExternalId: &artifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	_, err := service.UpsertModelArtifact(modelArtifact, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	_, err = service.GetModelArtifactByParams(nil, nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid parameters call, supply either (artifactName and modelVersionId), or externalId: bad request", err.Error())
}

func (suite *CoreTestSuite) TestGetModelArtifactByParamsWithNoResults() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	_, err := service.GetModelArtifactByParams(apiutils.Of("not-present"), &modelVersionId, nil)
	suite.NotNil(err)
	suite.Equal("no model artifacts found for artifactName=not-present, modelVersionId=2, externalId=: not found", err.Error())
}

func (suite *CoreTestSuite) TestGetModelArtifacts() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	modelArtifact1 := &openapi.ModelArtifact{
		Name:       &artifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &artifactUri,
		ExternalId: &artifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	secondArtifactName := "second-name"
	secondArtifactExtId := "second-ext-id"
	secondArtifactUri := "second-uri"
	modelArtifact2 := &openapi.ModelArtifact{
		Name:       &secondArtifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &secondArtifactUri,
		ExternalId: &secondArtifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	thirdArtifactName := "third-name"
	thirdArtifactExtId := "third-ext-id"
	thirdArtifactUri := "third-uri"
	modelArtifact3 := &openapi.ModelArtifact{
		Name:       &thirdArtifactName,
		State:      (*openapi.ArtifactState)(&artifactState),
		Uri:        &thirdArtifactUri,
		ExternalId: &thirdArtifactExtId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_string_prop": {
				MetadataStringValue: converter.NewMetadataStringValue(customString),
			},
		},
	}

	createdArtifact1, err := service.UpsertModelArtifact(modelArtifact1, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)
	createdArtifact2, err := service.UpsertModelArtifact(modelArtifact2, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)
	createdArtifact3, err := service.UpsertModelArtifact(modelArtifact3, &modelVersionId)
	suite.Nilf(err, "error creating new model artifact for %d", modelVersionId)

	createdArtifactId1, _ := converter.StringToInt64(createdArtifact1.Id)
	createdArtifactId2, _ := converter.StringToInt64(createdArtifact2.Id)
	createdArtifactId3, _ := converter.StringToInt64(createdArtifact3.Id)

	getAll, err := service.GetModelArtifacts(api.ListOptions{}, nil)
	suite.Nilf(err, "error getting all model artifacts")
	suite.Equalf(int32(3), getAll.Size, "expected three model artifacts")

	suite.Equal(*converter.Int64ToString(createdArtifactId1), *getAll.Items[0].Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId2), *getAll.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId3), *getAll.Items[2].Id)

	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByModelVersion, err := service.GetModelArtifacts(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, &modelVersionId)
	suite.Nilf(err, "error getting all model artifacts for %d", modelVersionId)
	suite.Equalf(int32(3), getAllByModelVersion.Size, "expected three model artifacts for model version %d", modelVersionId)

	suite.Equal(*converter.Int64ToString(createdArtifactId1), *getAllByModelVersion.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId2), *getAllByModelVersion.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId3), *getAllByModelVersion.Items[0].Id)
}
