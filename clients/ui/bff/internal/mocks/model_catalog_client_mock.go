package mocks

import (
	"log/slog"
	"net/url"
	"strings"

	"github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/stretchr/testify/mock"
)

type ModelCatalogClientMock struct {
	mock.Mock
}

func NewModelCatalogClientMock(logger *slog.Logger) (*ModelCatalogClientMock, error) {
	return &ModelCatalogClientMock{}, nil
}

func (m *ModelCatalogClientMock) GetAllCatalogModelsAcrossSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*openapi.CatalogModelList, error) {
	allMockSources := getAllCatalogModelMocks()

	return &allMockSources, nil
}

func (m *ModelCatalogClientMock) GetCatalogSourceModel(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*openapi.CatalogModel, error) {
	catalogModel := GetCatalogModelMocks()[0]

	return &catalogModel, nil
}

func (m *ModelCatalogClientMock) GetAllCatalogSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*openapi.CatalogSourceList, error) {
	allMockSources := GetCatalogSourceListMock()
	var filteredMockSources []openapi.CatalogSource

	name := pageValues.Get("name")

	if name != "" {
		nameFilterLower := strings.ToLower(name)
		for _, source := range allMockSources.Items {
			if strings.ToLower(source.Id) == nameFilterLower || strings.ToLower(source.Name) == nameFilterLower {
				filteredMockSources = append(filteredMockSources, source)
			}
		}
	} else {
		filteredMockSources = allMockSources.Items
	}
	catalogSourceList := openapi.CatalogSourceList{
		Items:         filteredMockSources,
		PageSize:      int32(10),
		NextPageToken: "",
		Size:          int32(len(filteredMockSources)),
	}

	return &catalogSourceList, nil
}

func (m *ModelCatalogClientMock) GetCatalogModelArtifacts(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*openapi.CatalogModelArtifactList, error) {
	allMockModelArtifacts := GetCatalogModelArtifactListMock()
	return &allMockModelArtifacts, nil
}
