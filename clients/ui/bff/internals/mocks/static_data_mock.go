package mocks

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
)

func GetRegisteredModelMocks() []openapi.RegisteredModel {
	model1 := openapi.RegisteredModel{
		CustomProperties: &map[string]openapi.MetadataValue{
			"my-label9": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue:  "property9",
					MetadataType: "string",
				},
			},
		},
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
		CustomProperties: &map[string]openapi.MetadataValue{
			"my-label9": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue:  "property9",
					MetadataType: "string",
				},
			},
		},
		Name:                     "Model Two",
		Description:              stringToPointer("This model does things and stuff"),
		ExternalId:               stringToPointer("345235987"),
		Id:                       stringToPointer("2"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		Owner:                    stringToPointer("John Watson"),
		State:                    stateToPointer(openapi.REGISTEREDMODELSTATE_LIVE),
	}

	return []openapi.RegisteredModel{model1, model2}
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
		CustomProperties: &map[string]openapi.MetadataValue{
			"my-label9": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue:  "property9",
					MetadataType: "string",
				},
			},
		},
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
		CustomProperties: &map[string]openapi.MetadataValue{
			"my-label9": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue:  "property9",
					MetadataType: "string",
				},
			},
		},
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

	return []openapi.ModelVersion{model1, model2}
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
		ArtifactType:             "TYPE_ONE",
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
		ArtifactType:             "TYPE_TWO",
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
		"my-label9": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "property9",
				MetadataType: "string",
			},
		},
	}

	return &result
}
