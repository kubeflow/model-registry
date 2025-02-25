package mocks

import (
	"context"
	"log/slog"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
)

func GetRegisteredModelMocks() []openapi.RegisteredModel {
	model1 := openapi.RegisteredModel{
		CustomProperties:         newCustomProperties(),
		Name:                     "Model One",
		Description:              stringToPointer("This model does things and stuff"),
		ExternalId:               stringToPointer("934589798"),
		Id:                       stringToPointer("1"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		Owner:                    stringToPointer("Sherlock Holmes"),
		State:                    stateToPointer(openapi.REGISTEREDMODELSTATE_LIVE),
	}

	model2 := openapi.RegisteredModel{
		CustomProperties:         newCustomProperties(),
		Name:                     "Model Two",
		Description:              stringToPointer("This model does things and stuff"),
		ExternalId:               stringToPointer("345235987"),
		Id:                       stringToPointer("2"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		Owner:                    stringToPointer("John Watson"),
		State:                    stateToPointer(openapi.REGISTEREDMODELSTATE_LIVE),
	}

	model3 := openapi.RegisteredModel{
		CustomProperties:         newCustomProperties(),
		Name:                     "Model Three",
		Description:              stringToPointer("This model does things and stuff"),
		ExternalId:               stringToPointer("345235989"),
		Id:                       stringToPointer("3"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249933"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249933"),
		Owner:                    stringToPointer("M. Oriarty"),
		State:                    stateToPointer(openapi.REGISTEREDMODELSTATE_ARCHIVED),
	}

	return []openapi.RegisteredModel{model1, model2, model3}
}

func GetRegisteredModelListMock() openapi.RegisteredModelList {
	models := GetRegisteredModelMocks()

	return openapi.RegisteredModelList{
		NextPageToken: "abcdefgh",
		PageSize:      2,
		Size:          int32(len(models)),
		Items:         models,
	}
}

func GetModelVersionMocks() []openapi.ModelVersion {
	model1 := openapi.ModelVersion{
		CustomProperties:         newCustomProperties(),
		Name:                     "Version One",
		Description:              stringToPointer("This version improves stuff and things"),
		ExternalId:               stringToPointer("934589798"),
		Id:                       stringToPointer("1"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		RegisteredModelId:        "1",
		Author:                   stringToPointer("Sherlock Holmes"),
		State:                    stateToPointer(openapi.MODELVERSIONSTATE_LIVE),
	}

	model2 := openapi.ModelVersion{
		CustomProperties:         newCustomProperties(),
		Name:                     "Version Two",
		Description:              stringToPointer("This version improves stuff and things"),
		ExternalId:               stringToPointer("934589799"),
		Id:                       stringToPointer("2"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		RegisteredModelId:        "2",
		Author:                   stringToPointer("Sherlock Holmes"),
		State:                    stateToPointer(openapi.MODELVERSIONSTATE_LIVE),
	}

	model3 := openapi.ModelVersion{
		CustomProperties:         newCustomProperties(),
		Name:                     "Version Three",
		Description:              stringToPointer("This version didn't improve stuff and things"),
		ExternalId:               stringToPointer("934589791"),
		Id:                       stringToPointer("3"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		RegisteredModelId:        "3",
		Author:                   stringToPointer("Sherlock Holmes"),
		State:                    stateToPointer(openapi.MODELVERSIONSTATE_ARCHIVED),
	}

	return []openapi.ModelVersion{model1, model2, model3}
}

func GetModelVersionListMock() openapi.ModelVersionList {
	versions := GetModelVersionMocks()

	return openapi.ModelVersionList{
		NextPageToken: "abcdefgh",
		PageSize:      2,
		Items:         versions,
		Size:          2,
	}
}

func GetModelArtifactMocks() []openapi.ModelArtifact {
	artifact1 := openapi.ModelArtifact{
		ArtifactType:             stringToPointer("TYPE_ONE"),
		CustomProperties:         newCustomProperties(),
		Description:              stringToPointer("This artifact can do more than you would expect"),
		ExternalId:               stringToPointer("1000001"),
		Uri:                      stringToPointer("http://localhost/artifacts/1"),
		State:                    stateToPointer(openapi.ARTIFACTSTATE_LIVE),
		Name:                     stringToPointer("Artifact One"),
		Id:                       stringToPointer("1"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		ModelFormatName:          stringToPointer("ONNX"),
		StorageKey:               stringToPointer("key1"),
		StoragePath:              stringToPointer("/artifacts/1"),
		ModelFormatVersion:       stringToPointer("1.0.0"),
		ServiceAccountName:       stringToPointer("service-1"),
	}

	artifact2 := openapi.ModelArtifact{
		ArtifactType:             stringToPointer("TYPE_TWO"),
		CustomProperties:         newCustomProperties(),
		Description:              stringToPointer("This artifact can do more than you would expect, but less than you would hope"),
		ExternalId:               stringToPointer("1000002"),
		Uri:                      stringToPointer("http://localhost/artifacts/2"),
		State:                    stateToPointer(openapi.ARTIFACTSTATE_PENDING),
		Name:                     stringToPointer("Artifact Two"),
		Id:                       stringToPointer("2"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		ModelFormatName:          stringToPointer("TensorFlow"),
		StorageKey:               stringToPointer("key2"),
		StoragePath:              stringToPointer("/artifacts/2"),
		ModelFormatVersion:       stringToPointer("1.0.0"),
		ServiceAccountName:       stringToPointer("service-2"),
	}

	return []openapi.ModelArtifact{artifact1, artifact2}
}

func GetModelArtifactListMock() openapi.ModelArtifactList {
	return openapi.ModelArtifactList{
		NextPageToken: "abcdefgh",
		PageSize:      2,
		Items:         GetModelArtifactMocks(),
		Size:          2,
	}
}

func newCustomProperties() *map[string]openapi.MetadataValue {
	result := map[string]openapi.MetadataValue{
		"tensorflow": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"pytorch": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"mll": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"rnn": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"AWS_KEY": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "asdf89asdf098asdfa",
				MetadataType: "MetadataStringValue",
			},
		},
		"AWS_PASSWORD": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "*AadfeDs34adf",
				MetadataType: "MetadataStringValue",
			},
		},
	}

	return &result
}

func NewMockSessionContext(parent context.Context) context.Context {
	if parent == nil {
		parent = context.TODO()
	}
	traceId := uuid.NewString()
	ctx := context.WithValue(parent, constants.TraceIdKey, traceId)

	traceLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx = context.WithValue(ctx, constants.TraceLoggerKey, traceLogger)
	return ctx
}

func NewMockSessionContextNoParent() context.Context {
	return NewMockSessionContext(context.TODO())
}

func GenerateMockArtifactList() openapi.ArtifactList {
	var artifacts []openapi.Artifact
	for i := 0; i < 2; i++ {
		artifact := GenerateMockArtifact()
		artifacts = append(artifacts, artifact)
	}

	return openapi.ArtifactList{
		NextPageToken: gofakeit.UUID(),
		PageSize:      int32(gofakeit.Number(1, 20)),
		Size:          int32(len(artifacts)),
		Items:         artifacts,
	}
}

func GenerateMockArtifact() openapi.Artifact {
	modelArtifact := GenerateMockModelArtifact()

	mockData := openapi.Artifact{
		ModelArtifact: &modelArtifact,
	}
	return mockData
}
