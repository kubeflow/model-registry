package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// MODEL VERSION ARTIFACTS

func (suite *CoreTestSuite) TestCreateModelVersionArtifact() {
	// create model registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	createdArt, err := service.UpsertModelVersionArtifact(&openapi.Artifact{
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
	}, modelVersionId)
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

func (suite *CoreTestSuite) TestCreateDuplicateModelVersionArtifactFailure() {
	// create model registry service
	service := suite.setupModelRegistryService()
	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	artifact := &openapi.Artifact{
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
	}

	_, err := service.UpsertModelVersionArtifact(artifact, modelVersionId)
	suite.Nilf(err, "error creating new artifact: %v", err)

	// attempt to create dupliate version artifact
	_, err = service.UpsertModelVersionArtifact(artifact, modelVersionId)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "cannot register a duplicate version artifact")
	suite.Equal(409, statusResp, "duplicate version artifacts not allowed")
}

func (suite *CoreTestSuite) TestCreateModelVersionArtifactFailure() {
	// create model registry service
	service := suite.setupModelRegistryService()

	_, err := service.UpsertModelVersionArtifact(nil, "")
	suite.NotNil(err)
	suite.Equal("invalid artifact pointer, can't upsert nil: bad request", err.Error())

	modelVersionId := "9998"

	artifact := &openapi.Artifact{
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
	}

	_, err = service.UpsertModelVersionArtifact(artifact, "")
	suite.NotNil(err)
	suite.Equal("no model version found for id : not found", err.Error())

	_, err = service.UpsertModelVersionArtifact(artifact, modelVersionId)
	suite.NotNil(err)
	suite.Equal("no model version found for id 9998: not found", err.Error())
}

func (suite *CoreTestSuite) TestUpdateModelVersionArtifact() {
	// create model registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	createdArtifact, err := service.UpsertModelVersionArtifact(&openapi.Artifact{
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
	}, modelVersionId)
	suite.Nilf(err, "error creating new artifact: %v", err)

	newState := "MARKED_FOR_DELETION"
	createdArtifact.DocArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertModelVersionArtifact(createdArtifact, modelVersionId)
	suite.Nilf(err, "error updating artifact for %s: %v", modelVersionId, err)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.DocArtifact.Id)
	updatedArtifactId, _ := converter.StringToInt64(updatedArtifact.DocArtifact.Id)
	suite.Equal(createdArtifactId, updatedArtifactId)

	getById, err := suite.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	suite.Nilf(err, "error getting artifact by id %s", *createdArtifactId)

	suite.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	suite.Equal(fmt.Sprintf("%s:%s", modelVersionId, *createdArtifact.DocArtifact.Name), *getById.Artifacts[0].Name)
	suite.Equal(string(newState), getById.Artifacts[0].State.String())
	suite.Equal(*createdArtifact.DocArtifact.Uri, *getById.Artifacts[0].Uri)
	suite.Equal((*createdArtifact.DocArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["custom_string_prop"].GetStringValue())
}

func (suite *CoreTestSuite) TestUpdateModelVersionArtifactFailure() {
	// create model registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	createdArtifact, err := service.UpsertModelVersionArtifact(&openapi.Artifact{
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
	}, modelVersionId)
	suite.Nilf(err, "error creating new artifact for model version %s", modelVersionId)
	suite.NotNilf(createdArtifact.DocArtifact.Id, "created model artifact should not have nil Id")

	newState := "MARKED_FOR_DELETION"
	createdArtifact.DocArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertModelVersionArtifact(createdArtifact, modelVersionId)
	suite.Nilf(err, "error updating artifact for %s: %v", modelVersionId, err)

	wrongId := "5555"
	updatedArtifact.DocArtifact.Id = &wrongId
	_, err = service.UpsertModelVersionArtifact(updatedArtifact, modelVersionId)
	suite.NotNil(err)
	suite.Equal(fmt.Sprintf("no artifact found for id %s: not found", wrongId), err.Error())

	rmName := "x1"
	mvName := "x2"
	modelVersion2Id := suite.registerModelVersion(service, &rmName, &rmName, &mvName, &mvName)

	_, err = service.UpsertModelVersionArtifact(createdArtifact, modelVersion2Id)
	suite.Nilf(err, "error associating artifact %s to model version %s: %v",
		*createdArtifact.DocArtifact.Id, modelVersion2Id, err)
}

func (suite *CoreTestSuite) TestUpsertModelVersionStandaloneArtifact() {
	// create model registry service
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
	})
	suite.Nilf(err, "error creating new artifact: %v", err)

	upsertedMVArtifact, err := service.UpsertModelVersionArtifact(createdArtifact, modelVersionId)
	suite.Nilf(err, "error upserting standalone artifact: %v", err)
	suite.Equal(*createdArtifact.DocArtifact.Id, *upsertedMVArtifact.DocArtifact.Id)

	associatedMV, err := service.getModelVersionByArtifactId(*createdArtifact.DocArtifact.Id)
	suite.Nilf(err, "error getting model version by artifact id: %v", err)
	suite.Equal(modelVersionId, *associatedMV.Id)
}

func (suite *CoreTestSuite) TestGetModelVersionArtifacts() {
	// create model registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	secondArtifactName := "second-name"
	secondArtifactExtId := "second-ext-id"
	secondArtifactUri := "second-uri"

	createdArtifact1, err := service.UpsertModelVersionArtifact(&openapi.Artifact{
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
	}, modelVersionId)
	suite.Nilf(err, "error creating new artifact: %v", err)
	createdArtifact2, err := service.UpsertModelVersionArtifact(&openapi.Artifact{
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
	}, modelVersionId)
	suite.Nilf(err, "error creating new artifact: %v", err)

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
	suite.Nilf(err, "error getting all model artifacts: %v", err)
	suite.Equalf(int32(2), getAllByModelVersion.Size, "expected 2 artifacts for model version %v", modelVersionId)

	suite.Equal(*converter.Int64ToString(createdArtifactId1), *getAllByModelVersion.Items[1].ModelArtifact.Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId2), *getAllByModelVersion.Items[0].DocArtifact.Id)
}

// ARTIFACTS

func (suite *CoreTestSuite) TestCreateArtifact() {
	// create model registry service
	service := suite.setupModelRegistryService()

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
	})
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

func (suite *CoreTestSuite) TestCreateArtifactFailure() {
	// create model registry service
	service := suite.setupModelRegistryService()

	_, err := service.UpsertArtifact(nil)
	suite.NotNil(err)
	suite.Equal("invalid artifact pointer, can't upsert nil: bad request", err.Error())

	artifact := &openapi.Artifact{}

	_, err = service.UpsertArtifact(artifact)
	suite.NotNil(err)
	suite.Equal("invalid artifact type, must be either ModelArtifact or DocArtifact: bad request", err.Error())
}

func (suite *CoreTestSuite) TestUpdateArtifact() {
	// create model registry service
	service := suite.setupModelRegistryService()

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
	})
	suite.Nilf(err, "error creating new artifact: %v", err)

	newState := "MARKED_FOR_DELETION"
	createdArtifact.DocArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertArtifact(createdArtifact)
	suite.Nilf(err, "error updating artifact: %v", err)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.DocArtifact.Id)
	updatedArtifactId, _ := converter.StringToInt64(updatedArtifact.DocArtifact.Id)
	suite.Equal(createdArtifactId, updatedArtifactId)

	getById, err := suite.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	suite.Nilf(err, "error getting artifact by id %s: %v", *createdArtifactId, err)

	suite.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	fmt.Printf("da name: %s, db name: %s", *createdArtifact.DocArtifact.Name, *getById.Artifacts[0].Name)
	exploded := strings.Split(*getById.Artifacts[0].Name, ":")
	suite.NotZero(exploded[0], "prefix should not be empty")
	suite.Equal(exploded[1], *createdArtifact.DocArtifact.Name)
	suite.Equal(string(newState), getById.Artifacts[0].State.String())
	suite.Equal(*createdArtifact.DocArtifact.Uri, *getById.Artifacts[0].Uri)
	suite.Equal((*createdArtifact.DocArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["custom_string_prop"].GetStringValue())
}

func (suite *CoreTestSuite) TestUpdateArtifactFailure() {
	// create model registry service
	service := suite.setupModelRegistryService()

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
	})
	suite.Nilf(err, "error creating new artifact for model version: %v", err)
	suite.NotNilf(createdArtifact.DocArtifact.Id, "created model artifact should not have nil Id")

	newState := "MARKED_FOR_DELETION"
	createdArtifact.DocArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertArtifact(createdArtifact)
	suite.Nilf(err, "error updating artifact: %v", err)

	wrongId := "5555"
	updatedArtifact.DocArtifact.Id = &wrongId
	_, err = service.UpsertArtifact(updatedArtifact)
	suite.NotNil(err)
	suite.Equal(fmt.Sprintf("no artifact found for id %s: not found", wrongId), err.Error())

	ma := &openapi.Artifact{
		ModelArtifact: &openapi.ModelArtifact{
			Id: createdArtifact.DocArtifact.Id,
		},
	}
	_, err = service.UpsertArtifact(ma)
	suite.NotNil(err)
	suite.Equal(fmt.Sprintf("mismatched types, artifact with id %s is not a model artifact: bad request", *createdArtifact.DocArtifact.Id), err.Error())
}

func (suite *CoreTestSuite) TestGetArtifactById() {
	// create model registry service
	service := suite.setupModelRegistryService()

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
	})
	suite.Nilf(err, "error creating new model artifact: %v", err)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.DocArtifact.Id)

	getById, err := service.GetArtifactById(*createdArtifact.DocArtifact.Id)
	suite.Nilf(err, "error getting artifact by id %s: %v", *createdArtifactId, err)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	suite.NotNil(createdArtifact.DocArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *getById.DocArtifact.Name)
	suite.Equal(*state, *getById.DocArtifact.State)
	suite.Equal(artifactUri, *getById.DocArtifact.Uri)
	suite.Equal(customString, (*getById.DocArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*createdArtifact, *getById, "artifacts returned during creation and on get by id should be equal")
}

func (suite *CoreTestSuite) TestGetArtifactByParams() {
	// create model registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	docArtifact := &openapi.DocArtifact{
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

	art, err := service.UpsertModelVersionArtifact(&openapi.Artifact{DocArtifact: docArtifact}, modelVersionId)
	suite.Nilf(err, "error creating new model artifact: %v", err)
	da := art.DocArtifact

	createdArtifactId, _ := converter.StringToInt64(da.Id)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)

	artByName, err := service.GetArtifactByParams(&artifactName, &modelVersionId, nil)
	suite.Nilf(err, "error getting model artifact by id %s: %v", *createdArtifactId, err)
	daByName := artByName.DocArtifact

	suite.NotNil(da.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *daByName.Name)
	suite.Equal(artifactExtId, *daByName.ExternalId)
	suite.Equal(*state, *daByName.State)
	suite.Equal(artifactUri, *daByName.Uri)
	suite.Equal(customString, (*daByName.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*da, *daByName, "artifacts returned during creation and on get by name should be equal")

	getByExtId, err := service.GetArtifactByParams(nil, nil, &artifactExtId)
	suite.Nilf(err, "error getting model artifact by id %s: %v", *createdArtifactId, err)
	daByExtId := getByExtId.DocArtifact

	suite.NotNil(da.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *daByExtId.Name)
	suite.Equal(artifactExtId, *daByExtId.ExternalId)
	suite.Equal(*state, *daByExtId.State)
	suite.Equal(artifactUri, *daByExtId.Uri)
	suite.Equal(customString, (*daByExtId.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*da, *daByExtId, "artifacts returned during creation and on get by ext id should be equal")
}

func (suite *CoreTestSuite) TestGetArtifactByParamsInvalid() {
	// trigger a 400 bad request to test unallowed query params
	// create model registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	invalidName := "\xFF"

	_, err := service.GetArtifactByParams(&invalidName, &modelVersionId, nil)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "invalid parameter used to retreive artifact")
	suite.Equal(400, statusResp, "invalid parameter used to retreive artifact")
}

func (suite *CoreTestSuite) TestGetArtifacts() {
	// create model registry service
	service := suite.setupModelRegistryService()

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
	})
	suite.Nilf(err, "error creating new artifact: %v", err)
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
	})
	suite.Nilf(err, "error creating new artifact: %v", err)

	createdArtifactId1, _ := converter.StringToInt64(createdArtifact1.ModelArtifact.Id)
	createdArtifactId2, _ := converter.StringToInt64(createdArtifact2.DocArtifact.Id)

	getAll, err := service.GetArtifacts(api.ListOptions{}, nil)
	suite.Nilf(err, "error getting all model artifacts")
	suite.Equalf(int32(2), getAll.Size, "expected two artifacts")

	suite.Equal(*converter.Int64ToString(createdArtifactId1), *getAll.Items[0].ModelArtifact.Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId2), *getAll.Items[1].DocArtifact.Id)

	orderByLastUpdate := "LAST_UPDATE_TIME"
	getAllByModelVersion, err := service.GetArtifacts(api.ListOptions{
		OrderBy:   &orderByLastUpdate,
		SortOrder: &descOrderDirection,
	}, nil)
	suite.Nilf(err, "error getting all model artifacts: %v", err)
	suite.Equalf(int32(2), getAllByModelVersion.Size, "expected 2 artifacts: %v", err)

	suite.Equal(*converter.Int64ToString(createdArtifactId1), *getAllByModelVersion.Items[1].ModelArtifact.Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId2), *getAllByModelVersion.Items[0].DocArtifact.Id)
}

// MODEL ARTIFACTS

func (suite *CoreTestSuite) TestCreateModelArtifact() {
	// create model registry service
	service := suite.setupModelRegistryService()

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
	})
	suite.Nilf(err, "error creating new model artifact: %v", err)

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
	// create model registry service
	service := suite.setupModelRegistryService()

	_, err := service.UpsertModelArtifact(nil)
	suite.NotNil(err)
	suite.Equal("invalid artifact pointer, can't upsert nil: bad request", err.Error())
}

func (suite *CoreTestSuite) TestUpdateModelArtifact() {
	// create model registry service
	service := suite.setupModelRegistryService()

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

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact)
	suite.Nilf(err, "error creating new model artifact: %v", err)

	newState := "MARKED_FOR_DELETION"
	createdArtifact.State = (*openapi.ArtifactState)(&newState)
	updatedArtifact, err := service.UpsertModelArtifact(createdArtifact)
	suite.Nilf(err, "error updating model artifact: %v", err)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)
	updatedArtifactId, _ := converter.StringToInt64(updatedArtifact.Id)
	suite.Equal(createdArtifactId, updatedArtifactId)

	getById, err := suite.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	suite.Nilf(err, "error getting model artifact by id %s: %v", *createdArtifactId, err)

	suite.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	exploded := strings.Split(*getById.Artifacts[0].Name, ":")
	suite.NotZero(exploded[0], "prefix should not be empty")
	suite.Equal(exploded[1], *createdArtifact.Name)
	suite.Equal(string(newState), getById.Artifacts[0].State.String())
	suite.Equal(*createdArtifact.Uri, *getById.Artifacts[0].Uri)
	suite.Equal((*createdArtifact.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["custom_string_prop"].GetStringValue())
}

func (suite *CoreTestSuite) TestGetModelArtifactById() {
	// create model registry service
	service := suite.setupModelRegistryService()

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

	createdArtifact, err := service.UpsertModelArtifact(modelArtifact)
	suite.Nilf(err, "error creating new model artifact: %v", err)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Id)

	getById, err := service.GetModelArtifactById(*createdArtifact.Id)
	suite.Nilf(err, "error getting model artifact by id %s: %v", *createdArtifactId, err)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	suite.NotNil(createdArtifact.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *getById.Name)
	suite.Equal(*state, *getById.State)
	suite.Equal(artifactUri, *getById.Uri)
	suite.Equal(customString, (*getById.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*createdArtifact, *getById, "artifacts returned during creation and on get by id should be equal")
}

func (suite *CoreTestSuite) TestGetModelArtifactByParams() {
	// create model registry service
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

	art, err := service.UpsertModelVersionArtifact(&openapi.Artifact{ModelArtifact: modelArtifact}, modelVersionId)
	suite.Nilf(err, "error creating new model artifact: %v", err)
	ma := art.ModelArtifact

	createdArtifactId, _ := converter.StringToInt64(ma.Id)

	state, _ := openapi.NewArtifactStateFromValue(artifactState)

	getByName, err := service.GetModelArtifactByParams(&artifactName, &modelVersionId, nil)
	suite.Nilf(err, "error getting model artifact by id %s: %v", *createdArtifactId, err)

	suite.NotNil(ma.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *getByName.Name)
	suite.Equal(artifactExtId, *getByName.ExternalId)
	suite.Equal(*state, *getByName.State)
	suite.Equal(artifactUri, *getByName.Uri)
	suite.Equal(customString, (*getByName.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*ma, *getByName, "artifacts returned during creation and on get by name should be equal")

	getByExtId, err := service.GetModelArtifactByParams(nil, nil, &artifactExtId)
	suite.Nilf(err, "error getting model artifact by id %s: %v", *createdArtifactId, err)

	suite.NotNil(ma.Id, "created artifact id should not be nil")
	suite.Equal(artifactName, *getByExtId.Name)
	suite.Equal(artifactExtId, *getByExtId.ExternalId)
	suite.Equal(*state, *getByExtId.State)
	suite.Equal(artifactUri, *getByExtId.Uri)
	suite.Equal(customString, (*getByExtId.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)

	suite.Equal(*ma, *getByExtId, "artifacts returned during creation and on get by ext id should be equal")
}

func (suite *CoreTestSuite) TestGetModelArtifactByEmptyParams() {
	// create model registry service
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

	_, err := service.UpsertModelVersionArtifact(&openapi.Artifact{ModelArtifact: modelArtifact}, modelVersionId)
	suite.Nilf(err, "error creating new model artifact: %v", err)

	_, err = service.GetModelArtifactByParams(nil, nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid parameters call, supply either (artifactName and parentResourceId), or externalId: bad request", err.Error())
}

func (suite *CoreTestSuite) TestGetModelArtifactByParamsWithNoResults() {
	// create model registry service
	service := suite.setupModelRegistryService()

	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	_, err := service.GetModelArtifactByParams(apiutils.Of("not-present"), &modelVersionId, nil)
	suite.NotNil(err)
	suite.Equal("no model artifacts found for artifactName=not-present, parentResourceId=2, externalId=: not found", err.Error())
}

func (suite *CoreTestSuite) TestGetModelArtifacts() {
	// create model registry service
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

	art1, err := service.UpsertModelVersionArtifact(&openapi.Artifact{ModelArtifact: modelArtifact1}, modelVersionId)
	suite.Nilf(err, "error creating new model artifact: %v", err)
	ma1 := art1.ModelArtifact
	art2, err := service.UpsertModelVersionArtifact(&openapi.Artifact{ModelArtifact: modelArtifact2}, modelVersionId)
	suite.Nilf(err, "error creating new model artifact: %v", err)
	ma2 := art2.ModelArtifact
	art3, err := service.UpsertModelVersionArtifact(&openapi.Artifact{ModelArtifact: modelArtifact3}, modelVersionId)
	suite.Nilf(err, "error creating new model artifact: %v", err)
	ma3 := art3.ModelArtifact

	createdArtifactId1, _ := converter.StringToInt64(ma1.Id)
	createdArtifactId2, _ := converter.StringToInt64(ma2.Id)
	createdArtifactId3, _ := converter.StringToInt64(ma3.Id)

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
	suite.Nilf(err, "error getting all model artifacts: %v", err)
	suite.Equalf(int32(3), getAllByModelVersion.Size, "expected three model artifacts for model version %v", modelVersionId)

	suite.Equal(*converter.Int64ToString(createdArtifactId1), *getAllByModelVersion.Items[2].Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId2), *getAllByModelVersion.Items[1].Id)
	suite.Equal(*converter.Int64ToString(createdArtifactId3), *getAllByModelVersion.Items[0].Id)
}

// DATASET ARTIFACTS

func (suite *CoreTestSuite) TestCreateDataSet() {
	// create model registry service
	service := suite.setupModelRegistryService()

	createdArt, err := service.UpsertArtifact(&openapi.Artifact{
		DataSet: &openapi.DataSet{
			Name:        &artifactName,
			State:       (*openapi.ArtifactState)(&artifactState),
			Uri:         &artifactUri,
			Description: &artifactDescription,
			Digest:      apiutils.Of("test-digest"),
			SourceType:  apiutils.Of("test-source-type"),
			Source:      apiutils.Of("test-source"),
			Schema:      apiutils.Of("test-schema"),
			Profile:     apiutils.Of("test-profile"),
			CustomProperties: &map[string]openapi.MetadataValue{
				"custom_string_prop": {
					MetadataStringValue: converter.NewMetadataStringValue(customString),
				},
			},
		},
	})
	suite.Nilf(err, "error creating new dataset: %v", err)

	dataset := createdArt.DataSet
	suite.NotNil(dataset, "error creating new dataset")
	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	suite.NotNil(dataset.Id, "created dataset id should not be nil")
	suite.Equal(artifactName, *dataset.Name)
	suite.Equal(*state, *dataset.State)
	suite.Equal(artifactUri, *dataset.Uri)
	suite.Equal(artifactDescription, *dataset.Description)
	suite.Equal("test-digest", *dataset.Digest)
	suite.Equal("test-source-type", *dataset.SourceType)
	suite.Equal("test-source", *dataset.Source)
	suite.Equal("test-schema", *dataset.Schema)
	suite.Equal("test-profile", *dataset.Profile)
	suite.Equal(customString, (*dataset.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)
}

func (suite *CoreTestSuite) TestUpdateDataSet() {
	// create model registry service
	service := suite.setupModelRegistryService()
	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	createdArtifact, err := service.UpsertModelVersionArtifact(&openapi.Artifact{
		DataSet: &openapi.DataSet{
			Name:  &artifactName,
			State: (*openapi.ArtifactState)(&artifactState),
			Uri:   &artifactUri,
			CustomProperties: &map[string]openapi.MetadataValue{
				"custom_string_prop": {
					MetadataStringValue: converter.NewMetadataStringValue(customString),
				},
			},
		},
	}, modelVersionId)
	suite.Nilf(err, "error creating new dataset: %v", err)

	// update with only id and updated state
	createdArtifact = &openapi.Artifact{
		DataSet: &openapi.DataSet{
			Id:    createdArtifact.DataSet.Id,
			State: apiutils.Of(openapi.ARTIFACTSTATE_MARKED_FOR_DELETION),
		},
	}
	updatedArtifact, err := service.UpsertArtifact(createdArtifact)
	suite.Nilf(err, "error updating dataset: %v", err)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.DataSet.Id)
	updatedArtifactId, _ := converter.StringToInt64(updatedArtifact.DataSet.Id)
	suite.Equal(createdArtifactId, updatedArtifactId)
	suite.Equal(artifactName, *updatedArtifact.DataSet.Name)
	suite.Equal(*createdArtifact.DataSet.State, *updatedArtifact.DataSet.State)

	getById, err := suite.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	suite.Nilf(err, "error getting dataset by id %s: %v", *createdArtifactId, err)

	suite.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	exploded := strings.Split(*getById.Artifacts[0].Name, ":")
	suite.NotZero(exploded[0], "prefix should not be empty")
	suite.Equal(exploded[1], *createdArtifact.DataSet.Name)
	suite.Equal(*updatedArtifact.DataSet.Uri, *getById.Artifacts[0].Uri)
	suite.Equal((*updatedArtifact.DataSet.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["custom_string_prop"].GetStringValue())
}

// METRIC ARTIFACTS

func (suite *CoreTestSuite) TestCreateMetric() {
	// create model registry service
	service := suite.setupModelRegistryService()
	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	createdArt, err := service.UpsertModelVersionArtifact(&openapi.Artifact{
		Metric: &openapi.Metric{
			Name:        &artifactName,
			State:       (*openapi.ArtifactState)(&artifactState),
			Description: &artifactDescription,
			Value:       apiutils.Of(42.5),
			Timestamp:   apiutils.Of("1234567890"),
			Step:        apiutils.Of(int64(1)),
			CustomProperties: &map[string]openapi.MetadataValue{
				"custom_string_prop": {
					MetadataStringValue: converter.NewMetadataStringValue(customString),
				},
			},
		},
	}, modelVersionId)
	suite.Nilf(err, "error creating new metric: %v", err)

	metric := createdArt.Metric
	suite.NotNil(metric, "error creating new metric")
	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	suite.NotNil(metric.Id, "created metric id should not be nil")
	suite.Equal(artifactName, *metric.Name)
	suite.Equal(*state, *metric.State)
	suite.Equal(artifactDescription, *metric.Description)
	suite.Equal(42.5, *metric.Value)
	suite.Equal("1234567890", *metric.Timestamp)
	suite.Equal(int64(1), *metric.Step)
	suite.Equal(customString, (*metric.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)
}

func (suite *CoreTestSuite) TestUpdateMetric() {
	// create model registry service
	service := suite.setupModelRegistryService()
	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	createdArtifact, err := service.UpsertModelVersionArtifact(&openapi.Artifact{
		Metric: &openapi.Metric{
			Name:  &artifactName,
			State: (*openapi.ArtifactState)(&artifactState),
			Value: apiutils.Of(42.5),
			CustomProperties: &map[string]openapi.MetadataValue{
				"custom_string_prop": {
					MetadataStringValue: converter.NewMetadataStringValue(customString),
				},
			},
		},
	}, modelVersionId)
	suite.Nilf(err, "error creating new metric: %v", err)

	newValue := 52.4
	// update with only id and updated value
	createdArtifact = &openapi.Artifact{
		Metric: &openapi.Metric{
			Id:    createdArtifact.Metric.Id,
			Value: &newValue,
		},
	}
	updatedArtifact, err := service.UpsertArtifact(createdArtifact)
	suite.Nilf(err, "error updating metric: %v", err)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Metric.Id)
	updatedArtifactId, _ := converter.StringToInt64(updatedArtifact.Metric.Id)
	suite.Equal(createdArtifactId, updatedArtifactId)
	suite.Equal(artifactName, *updatedArtifact.Metric.Name)
	suite.Equal(newValue, *updatedArtifact.Metric.Value)

	getById, err := suite.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	suite.Nilf(err, "error getting metric by id %s: %v", *createdArtifactId, err)

	suite.Equal(*createdArtifactId, *getById.Artifacts[0].Id)
	exploded := strings.Split(*getById.Artifacts[0].Name, ":")
	suite.NotZero(exploded[0], "prefix should not be empty")
	suite.Equal(exploded[1], *createdArtifact.Metric.Name)
}

// PARAMETER ARTIFACTS

func (suite *CoreTestSuite) TestCreateParameter() {
	// create model registry service
	service := suite.setupModelRegistryService()
	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	createdArt, err := service.UpsertModelVersionArtifact(&openapi.Artifact{
		Parameter: &openapi.Parameter{
			Name:          &artifactName,
			State:         (*openapi.ArtifactState)(&artifactState),
			Description:   &artifactDescription,
			Value:         apiutils.Of("test-value"),
			ParameterType: apiutils.Of(openapi.PARAMETERTYPE_STRING),
			CustomProperties: &map[string]openapi.MetadataValue{
				"custom_string_prop": {
					MetadataStringValue: converter.NewMetadataStringValue(customString),
				},
			},
		},
	}, modelVersionId)
	suite.Nilf(err, "error creating new parameter: %v", err)

	parameter := createdArt.Parameter
	suite.NotNil(parameter, "error creating new parameter")
	state, _ := openapi.NewArtifactStateFromValue(artifactState)
	suite.NotNil(parameter.Id, "created parameter id should not be nil")
	suite.Equal(artifactName, *parameter.Name)
	suite.Equal(*state, *parameter.State)
	suite.Equal(artifactDescription, *parameter.Description)
	suite.Equal("test-value", *parameter.Value)
	suite.Equal(openapi.PARAMETERTYPE_STRING, *parameter.ParameterType)
	suite.Equal(customString, (*parameter.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue)
}

func (suite *CoreTestSuite) TestUpdateParameter() {
	// create model registry service
	service := suite.setupModelRegistryService()
	modelVersionId := suite.registerModelVersion(service, nil, nil, nil, nil)

	createdArtifact, err := service.UpsertModelVersionArtifact(&openapi.Artifact{
		Parameter: &openapi.Parameter{
			Name:          &artifactName,
			State:         (*openapi.ArtifactState)(&artifactState),
			Value:         apiutils.Of("test-value"),
			ParameterType: apiutils.Of(openapi.PARAMETERTYPE_STRING),
			CustomProperties: &map[string]openapi.MetadataValue{
				"custom_string_prop": {
					MetadataStringValue: converter.NewMetadataStringValue(customString),
				},
			},
		},
	}, modelVersionId)
	suite.Nilf(err, "error creating new parameter: %v", err)

	// update with only id and updated value
	newValue := "new-test-value"
	createdArtifact = &openapi.Artifact{
		Parameter: &openapi.Parameter{
			Id:    createdArtifact.Parameter.Id,
			Value: apiutils.Of(newValue),
		},
	}
	updatedArtifact, err := service.UpsertArtifact(createdArtifact)
	suite.Nilf(err, "error updating parameter: %v", err)

	createdArtifactId, _ := converter.StringToInt64(createdArtifact.Parameter.Id)
	updatedArtifactId, _ := converter.StringToInt64(updatedArtifact.Parameter.Id)
	suite.Equal(createdArtifactId, updatedArtifactId)
	suite.Equal(artifactName, *updatedArtifact.Parameter.Name)
	suite.Equal(newValue, *updatedArtifact.Parameter.Value)

	getById, err := suite.mlmdClient.GetArtifactsByID(context.Background(), &proto.GetArtifactsByIDRequest{
		ArtifactIds: []int64{*createdArtifactId},
	})
	suite.Nilf(err, "error getting parameter by id %s: %v", *createdArtifactId, err)

	suite.Equal(*updatedArtifactId, *getById.Artifacts[0].Id)
	exploded := strings.Split(*getById.Artifacts[0].Name, ":")
	suite.NotZero(exploded[0], "prefix should not be empty")
	suite.Equal(exploded[1], *updatedArtifact.Parameter.Name)
	suite.Equal((*updatedArtifact.Parameter.CustomProperties)["custom_string_prop"].MetadataStringValue.StringValue, getById.Artifacts[0].CustomProperties["custom_string_prop"].GetStringValue())
}
