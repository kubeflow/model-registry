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
