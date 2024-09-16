package mocks

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

func GenerateMockRegisteredModelList() openapi.RegisteredModelList {
	var models []openapi.RegisteredModel
	for i := 0; i < 2; i++ {
		model := GenerateMockRegisteredModel()
		models = append(models, model)
	}

	return openapi.RegisteredModelList{
		NextPageToken: gofakeit.UUID(),
		PageSize:      int32(gofakeit.Number(1, 20)),
		Size:          int32(len(models)),
		Items:         models,
	}
}

func GenerateMockRegisteredModel() openapi.RegisteredModel {
	model := openapi.RegisteredModel{
		CustomProperties: &map[string]openapi.MetadataValue{
			"example_key": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue:  gofakeit.Sentence(3),
					MetadataType: "string",
				},
			},
		},
		Description:              stringToPointer(gofakeit.Sentence(5)),
		ExternalId:               stringToPointer(gofakeit.UUID()),
		Name:                     gofakeit.Name(),
		Id:                       stringToPointer(gofakeit.UUID()),
		CreateTimeSinceEpoch:     stringToPointer(fmt.Sprintf("%d", gofakeit.Date().UnixMilli())),
		LastUpdateTimeSinceEpoch: stringToPointer(fmt.Sprintf("%d", gofakeit.Date().UnixMilli())),
		Owner:                    stringToPointer(gofakeit.Name()),
		State:                    stateToPointer(openapi.RegisteredModelState(gofakeit.RandomString([]string{string(openapi.REGISTEREDMODELSTATE_LIVE), string(openapi.REGISTEREDMODELSTATE_ARCHIVED)}))),
	}
	return model
}

func GenerateMockModelVersion() openapi.ModelVersion {
	model := openapi.ModelVersion{
		CustomProperties: &map[string]openapi.MetadataValue{
			"example_key": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue:  gofakeit.Sentence(3),
					MetadataType: "string",
				},
			},
		},
		Description:              stringToPointer(gofakeit.Sentence(5)),
		ExternalId:               stringToPointer(gofakeit.UUID()),
		Name:                     gofakeit.Name(),
		Id:                       stringToPointer(gofakeit.UUID()),
		CreateTimeSinceEpoch:     stringToPointer(fmt.Sprintf("%d", gofakeit.Date().UnixMilli())),
		LastUpdateTimeSinceEpoch: stringToPointer(fmt.Sprintf("%d", gofakeit.Date().UnixMilli())),
		Author:                   stringToPointer(gofakeit.Name()),
		State:                    stateToPointer(openapi.ModelVersionState(gofakeit.RandomString([]string{string(openapi.MODELVERSIONSTATE_LIVE), string(openapi.MODELVERSIONSTATE_ARCHIVED)}))),
	}
	return model
}

func GenerateMockModelVersionList() openapi.ModelVersionList {
	var versions []openapi.ModelVersion

	for i := 0; i < 2; i++ {
		version := GenerateMockModelVersion()
		versions = append(versions, version)
	}

	return openapi.ModelVersionList{
		NextPageToken: gofakeit.UUID(),
		PageSize:      int32(gofakeit.Number(1, 20)),
		Size:          int32(len(versions)),
		Items:         versions,
	}
}

func stateToPointer[T any](s T) *T {
	return &s
}

func stringToPointer(s string) *string {
	return &s
}
