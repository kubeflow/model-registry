package catalog

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/genqlient"
	"github.com/kubeflow/model-registry/catalog/pkg/openapi"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	models "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

type rhecModel struct {
	models.CatalogModel `yaml:",inline"`
	Artifacts           []*openapi.CatalogModelArtifact `yaml:"artifacts"`
}

// rhecCatalogConfig defines the structure of the RHEC catalog configuration.
type rhecCatalogConfig struct {
	Models []struct {
		Repository string `yaml:"repository"`
	} `yaml:"models"`
}

type rhecCatalogImpl struct {
	modelsLock sync.RWMutex
	models     map[string]*rhecModel
}

var _ CatalogSourceProvider = &rhecCatalogImpl{}

func (r *rhecCatalogImpl) GetModel(ctx context.Context, name string) (*openapi.CatalogModel, error) {
	r.modelsLock.RLock()
	defer r.modelsLock.RUnlock()

	rm := r.models[name]
	if rm == nil {
		return nil, nil
	}
	cp := rm.CatalogModel
	return &cp, nil
}

func (r *rhecCatalogImpl) ListModels(ctx context.Context, params ListModelsParams) (openapi.CatalogModelList, error) {
	r.modelsLock.RLock()
	defer r.modelsLock.RUnlock()

	var filteredModels []*model.CatalogModel
	for _, rm := range r.models {
		cm := rm.CatalogModel
		if params.Query != "" {
			query := strings.ToLower(params.Query)
			// Check if query matches name, description, tasks, provider, or libraryName
			if !strings.Contains(strings.ToLower(cm.Name), query) &&
				!strings.Contains(strings.ToLower(cm.GetDescription()), query) &&
				!strings.Contains(strings.ToLower(cm.GetProvider()), query) &&
				!strings.Contains(strings.ToLower(cm.GetLibraryName()), query) {

				// Check tasks
				foundInTasks := false
				for _, task := range cm.GetTasks() { // Use GetTasks() for nil safety
					if strings.Contains(strings.ToLower(task), query) {
						foundInTasks = true
						break
					}
				}
				if !foundInTasks {
					continue // Skip if no match in any searchable field
				}
			}
		}
		filteredModels = append(filteredModels, &cm)
	}

	// Sort the filtered models
	sort.Slice(filteredModels, func(i, j int) bool {
		a := filteredModels[i]
		b := filteredModels[j]

		var less bool
		switch params.OrderBy {
		case model.ORDERBYFIELD_CREATE_TIME:
			// Convert CreateTimeSinceEpoch (string) to int64 for comparison
			// Handle potential nil or conversion errors by treating as 0
			aTime, _ := strconv.ParseInt(a.GetCreateTimeSinceEpoch(), 10, 64)
			bTime, _ := strconv.ParseInt(b.GetCreateTimeSinceEpoch(), 10, 64)
			less = aTime < bTime
		case model.ORDERBYFIELD_LAST_UPDATE_TIME:
			// Convert LastUpdateTimeSinceEpoch (string) to int64 for comparison
			// Handle potential nil or conversion errors by treating as 0
			aTime, _ := strconv.ParseInt(a.GetLastUpdateTimeSinceEpoch(), 10, 64)
			bTime, _ := strconv.ParseInt(b.GetLastUpdateTimeSinceEpoch(), 10, 64)
			less = aTime < bTime
		case model.ORDERBYFIELD_NAME:
			fallthrough
		default:
			// Fallback to name sort if an unknown sort field is provided
			less = strings.Compare(a.Name, b.Name) < 0
		}

		if params.SortOrder == model.SORTORDER_DESC {
			return !less
		}
		return less
	})

	count := len(filteredModels)
	if count > math.MaxInt32 {
		count = math.MaxInt32
	}

	list := model.CatalogModelList{
		Items:    make([]model.CatalogModel, count),
		PageSize: int32(count),
		Size:     int32(count),
	}
	for i := range list.Items {
		list.Items[i] = *filteredModels[i]
	}
	return list, nil // Return the struct value directly
}

func (r *rhecCatalogImpl) GetArtifacts(ctx context.Context, name string) (*openapi.CatalogModelArtifactList, error) {
	r.modelsLock.RLock()
	defer r.modelsLock.RUnlock()

	rm := r.models[name]
	if rm == nil {
		return nil, nil
	}

	count := len(rm.Artifacts)
	if count > math.MaxInt32 {
		count = math.MaxInt32
	}

	list := openapi.CatalogModelArtifactList{
		Items:    make([]openapi.CatalogModelArtifact, count),
		PageSize: int32(count),
		Size:     int32(count),
	}
	for i := range list.Items {
		list.Items[i] = *rm.Artifacts[i]
	}
	return &list, nil
}

func fetchRepository(ctx context.Context, client graphql.Client, repository string) (*genqlient.GetRepositoryResponse, error) {
	resp, err := genqlient.GetRepository(ctx, client, "registry.access.redhat.com", repository)
	if err != nil {
		return nil, fmt.Errorf("failed to query rhec repository: %w", err)
	}

	if err := resp.Get_repository_by_registry_path.Error; err.Detail != "" || err.Status != 0 {
		return nil, fmt.Errorf("rhec repository query error: detail: %s, status: %d", err.Detail, err.Status)
	}
	return resp, nil
}

func fetchRepositoryImages(ctx context.Context, client graphql.Client, repository string) ([]genqlient.FindRepositoryImagesFind_repository_images_by_registry_pathContainerImagePaginatedResponseDataContainerImage, error) {
	resp, err := genqlient.FindRepositoryImages(ctx, client, "registry.access.redhat.com", repository)
	if err != nil {
		return nil, fmt.Errorf("failed to query rhec images: %w", err)
	}

	if err := resp.Find_repository_images_by_registry_path.Error; err.Detail != "" || err.Status != 0 {
		return nil, fmt.Errorf("rhec images query error: detail: %s, status: %d", err.Detail, err.Status)
	}
	return resp.Find_repository_images_by_registry_path.Data, nil
}

func newRhecModel(repoData *genqlient.GetRepositoryResponse, imageData genqlient.FindRepositoryImagesFind_repository_images_by_registry_pathContainerImagePaginatedResponseDataContainerImage, imageTagName, repositoryName string) *rhecModel {

	sourceId := "rhec"
	createTime := repoData.Get_repository_by_registry_path.Data.Creation_date.Format(time.RFC3339)
	lastUpdateTime := repoData.Get_repository_by_registry_path.Data.Last_update_date.Format(time.RFC3339)
	description := repoData.Get_repository_by_registry_path.Data.Display_data.Short_description
	readme := repoData.Get_repository_by_registry_path.Data.Display_data.Long_description
	provider := repoData.Get_repository_by_registry_path.Data.Vendor_label

	var maturity *string
	if len(repoData.Get_repository_by_registry_path.Data.Release_categories) > 0 {
		maturityStr := repoData.Get_repository_by_registry_path.Data.Release_categories[0]
		maturity = &maturityStr
	}

	var tasks []string
	for _, label := range imageData.Parsed_data.Labels {
		tasks = append(tasks, label.Value)
	}
	imageCreationDate := imageData.Creation_date.Format(time.RFC3339)
	imageLastUpdateDate := imageData.Last_update_date.Format(time.RFC3339)

	modelName := repositoryName + ":" + imageTagName

	return &rhecModel{
		CatalogModel: openapi.CatalogModel{
			Name:                     modelName,
			CreateTimeSinceEpoch:     &createTime,
			LastUpdateTimeSinceEpoch: &lastUpdateTime,
			Description:              &description,
			Readme:                   &readme,
			Maturity:                 maturity,
			Language:                 []string{},
			Tasks:                    tasks,
			Provider:                 &provider,
			Logo:                     nil,
			License:                  nil,
			LicenseLink:              nil,
			LibraryName:              nil,
			SourceId:                 &sourceId,
		},
		Artifacts: []*openapi.CatalogModelArtifact{
			{
				Uri:                      "oci://registry.redhat.io/" + repositoryName + ":" + imageTagName,
				CreateTimeSinceEpoch:     &imageCreationDate,
				LastUpdateTimeSinceEpoch: &imageLastUpdateDate,
			},
		},
	}
}

func (r *rhecCatalogImpl) load(modelsList []any) error {
	graphqlClient := graphql.NewClient("https://catalog.redhat.com/api/containers/graphql/", http.DefaultClient)
	ctx := context.Background()

	models := make(map[string]*rhecModel)
	for _, modelEntry := range modelsList {
		modelMap, ok := modelEntry.(map[string]any)
		if !ok {
			glog.Warningf("skipping invalid entry in 'models' list")
			continue
		}
		repo, ok := modelMap["repository"].(string)
		if !ok {
			glog.Warningf("skipping model with missing or invalid 'repository'")
			continue
		}

		repoData, err := fetchRepository(ctx, graphqlClient, repo)
		if err != nil {
			return err
		}

		imagesData, err := fetchRepositoryImages(ctx, graphqlClient, repo)
		if err != nil {
			return err
		}

		for _, image := range imagesData {
			for _, imageRepository := range image.Repositories {
				for _, imageTag := range imageRepository.Tags {
					tagName := imageTag.Name
					fullModelName := repo + ":" + tagName
					model := newRhecModel(repoData, image, tagName, repo)
					models[fullModelName] = model
				}
			}
		}
	}

	r.modelsLock.Lock()
	defer r.modelsLock.Unlock()
	r.models = models

	return nil
}

func newRhecCatalog(source *CatalogSourceConfig) (CatalogSourceProvider, error) {
	modelsData, ok := source.Properties["models"]
	if !ok {
		return nil, fmt.Errorf("missing 'models' property for rhec catalog")
	}

	modelsList, ok := modelsData.([]any)
	if !ok {
		return nil, fmt.Errorf("'models' property should be a list")
	}

	r := &rhecCatalogImpl{
		models: make(map[string]*rhecModel),
	}

	err := r.load(modelsList)
	if err != nil {
		return nil, fmt.Errorf("error loading rhec catalog: %w", err)
	}

	return r, nil
}

func init() {
	if err := RegisterCatalogType("rhec", newRhecCatalog); err != nil {
		panic(err)
	}
}
