package catalog

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golang/glog"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	graphQL "github.com/shurcooL/graphql"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type rhecModel struct {
	model.CatalogModel `yaml:",inline"`
	Artifacts          []*model.CatalogModelArtifact `yaml:"artifacts"`
}

// rhecCatalogConfig defines the structure of the RHEC catalog configuration file (ex. sample-rhec.yaml).
type rhecCatalogConfig struct {
	Source string `yaml:"source"`
	Models []struct {
		Repository string `yaml:"repository"`
	} `yaml:"models"`
}

type rhecCatalogImpl struct {
	modelsLock sync.RWMutex
	models     map[string]*rhecModel
}

var _ CatalogSourceProvider = &rhecCatalogImpl{}

func (r *rhecCatalogImpl) GetModel(ctx context.Context, name string) (*model.CatalogModel, error) {
	r.modelsLock.RLock()
	defer r.modelsLock.RUnlock()

	rm := r.models[name]
	if rm == nil {
		return nil, nil
	}
	cp := rm.CatalogModel
	return &cp, nil
}

func (r *rhecCatalogImpl) ListModels(ctx context.Context, params ListModelsParams) (model.CatalogModelList, error) {
	r.modelsLock.RLock()
	defer r.modelsLock.RUnlock()

	items := make([]model.CatalogModel, 0, len(r.models))
	for _, rm := range r.models {
		items = append(items, rm.CatalogModel)
	}

	count := len(items)
	if count > math.MaxInt32 {
		count = math.MaxInt32
	}

	return model.CatalogModelList{
		Items:         items,
		PageSize:      int32(count),
		Size:          int32(count),
		NextPageToken: "",
	}, nil
}

func (r *rhecCatalogImpl) GetArtifacts(ctx context.Context, name string) (*model.CatalogModelArtifactList, error) {
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

	list := model.CatalogModelArtifactList{
		Items:    make([]model.CatalogModelArtifact, count),
		PageSize: int32(count),
		Size:     int32(count),
	}
	for i := range list.Items {
		list.Items[i] = *rm.Artifacts[i]
	}
	return &list, nil
}

var getRepositoryQuery struct {
	GetRepositoryResponse struct {
		Error struct {
			Detail graphQL.String `graphql:"detail"`
			Status graphQL.String `graphql:"status"`
		} `graphql:"error"`
		Data struct {
			CreationDate      graphQL.String   `graphql:"creation_date"`
			LastUpdateDate    graphQL.String   `graphql:"last_update_date"`
			ReleaseCategories []graphQL.String `graphql:"release_categories"`
			VendorLabel       graphQL.String   `graphql:"vendor_label"`
			DisplayData       struct {
				ShortDescription graphQL.String `graphql:"short_description"`
				LongDescription  graphQL.String `graphql:"long_description"`
			} `graphql:"display_data"`
		} `graphql:"data"`
	} `graphql:"get_repository_by_registry_path(registry: $registry, repository: $repository)"`
}

// 		CreateTimeSinceEpoch:      			repo.creation_date
// 		LastUpdateTimeSinceEpoch: 			repo.last_update_date
// 		Description: 						repo.RepositoryDisplayData.short_description
// 		Readme:								repo.RepositoryDisplayData.long_description
// 		Maturity:							repo.release_categories (returns list, take first?)
// 		Language:							????????
// 		Tasks:								image.parseddata.labels[].name
// 		Provider:							repo.vendor_label
// 		Logo:								???????
// 		License:							???????
// 		LicenseLink:						???????
// 		LibraryName:						???????
// 		CustomProperties:					???????
// 		SourceId:							rhec
//      Name:								repository (this is the full rhelai1/model-name-blah) + tag

// type CatalogModelArtifact struct {
//     CreateTimeSinceEpoch: image.creation_date
//     LastUpdateTimeSinceEpoch: image.last_update_date
//     Uri: "registry.redhat.io" + "/" + repository (rhelai/model-name-blah) + ":" + image.tags[each].name
// 		oci://registry.redhat.io/rhelai1/modelcar-granite-7b-redhat-lab:1.4.0
// }

var findRepositoryImagesQuery struct {
	FindRepositoryImagesResponse struct {
		Error []struct {
			Detail graphQL.String `graphql:"detail"`
			Status graphQL.String `graphql:"status"`
		} `graphql:"error"`
		Total graphQL.Int `graphql:"total"`
		Data  []struct {
			CreationDate   graphQL.String `graphql:"creation_date"`
			LastUpdateDate graphQL.String `graphql:"last_update_date"`
			Repositories   []struct {
				Registry graphQL.String `graphql:"registry"`
				Tags     []struct {
					Name graphQL.String `graphql:"name"`
				} `graphql:"tags"`
			} `graphql:"repositories"`
			ParsedData struct {
				Labels []struct {
					Name  graphQL.String `graphql:"name"`
					Value graphQL.String `graphql:"value"`
				} `graphql:"labels"`
			} `graphql:"parsed_data"`
		} `graphql:"data"`
	} `graphql:"find_repository_images_by_registry_path(registry: $registry, repository: $repository, sort_by: [{ field: \"creation_date\", order: DESC }])"`
}

func (r *rhecCatalogImpl) load(path string) error {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s file: %w", path, err)
	}

	var contents rhecCatalogConfig
	if err = yaml.UnmarshalStrict(fileBytes, &contents); err != nil {
		return fmt.Errorf("failed to parse %s file: %w", path, err)
	}

	graphQLClient := graphQL.NewClient("https://catalog.redhat.com/api/containers/graphql/", nil)

	models := make(map[string]*rhecModel, len(contents.Models))
	for _, m := range contents.Models {
		queryVariables := map[string]any{
			"registry":   graphQL.String("registry.access.redhat.com"),
			"repository": graphQL.String(m.Repository),
		}
		err := graphQLClient.Query(context.Background(), &getRepositoryQuery, queryVariables)
		if err != nil {
			return fmt.Errorf("failed to query rhec repository: %w", err)
		}

		if getRepositoryQuery.GetRepositoryResponse.Error.Detail != "" || getRepositoryQuery.GetRepositoryResponse.Error.Status != "" {
			return fmt.Errorf("rhec repository query error: detail: %s, status: %s", getRepositoryQuery.GetRepositoryResponse.Error.Detail, getRepositoryQuery.GetRepositoryResponse.Error.Status)
		}

		sourceId := "rhec"
		createTime := string(getRepositoryQuery.GetRepositoryResponse.Data.CreationDate)
		lastUpdateTime := string(getRepositoryQuery.GetRepositoryResponse.Data.LastUpdateDate)
		description := string(getRepositoryQuery.GetRepositoryResponse.Data.DisplayData.ShortDescription)
		readme := string(getRepositoryQuery.GetRepositoryResponse.Data.DisplayData.LongDescription)
		provider := string(getRepositoryQuery.GetRepositoryResponse.Data.VendorLabel)

		TEMPPOINTERPLACEHOLDER := " "

		var maturity *string
		if len(getRepositoryQuery.GetRepositoryResponse.Data.ReleaseCategories) > 0 {
			maturityStr := string(getRepositoryQuery.GetRepositoryResponse.Data.ReleaseCategories[0])
			maturity = &maturityStr
		}

		err = graphQLClient.Query(context.Background(), &findRepositoryImagesQuery, queryVariables)
		if err != nil {
			return fmt.Errorf("failed to query rhec images: %w", err)
		}

		if len(findRepositoryImagesQuery.FindRepositoryImagesResponse.Error) > 0 {
			var errorStrings []string
			for _, err := range findRepositoryImagesQuery.FindRepositoryImagesResponse.Error {
				errorStrings = append(errorStrings, fmt.Sprintf("detail: %s, status: %s", err.Detail, err.Status))
			}
			return fmt.Errorf("rhec images query errors: %s", strings.Join(errorStrings, "; "))
		}

		for _, image := range findRepositoryImagesQuery.FindRepositoryImagesResponse.Data {
			var tasks []string
			for _, label := range image.ParsedData.Labels {
				tasks = append(tasks, string(label.Value))
			}
			imageCreationDate := string(image.CreationDate)
			imageLastUpdateDate := string(image.LastUpdateDate)

			for _, imageRepository := range image.Repositories {
				for _, imageTag := range imageRepository.Tags {

					models[m.Repository+":"+string(imageTag.Name)] = &rhecModel{
						CatalogModel: model.CatalogModel{
							Name:                     m.Repository + ":" + string(imageTag.Name),
							CreateTimeSinceEpoch:     &createTime,
							LastUpdateTimeSinceEpoch: &lastUpdateTime,
							Description:              &description,
							Readme:                   &readme,
							Maturity:                 maturity,
							Language:                 []string{},
							Tasks:                    tasks,
							Provider:                 &provider,
							Logo:                     &TEMPPOINTERPLACEHOLDER,
							License:                  &TEMPPOINTERPLACEHOLDER,
							LicenseLink:              &TEMPPOINTERPLACEHOLDER,
							LibraryName:              &TEMPPOINTERPLACEHOLDER,
							SourceId:                 &sourceId,
						},
						Artifacts: []*model.CatalogModelArtifact{
							{
								Uri:                      "registry.redhat.io" + "/" + m.Repository + ":" + string(imageTag.Name),
								CreateTimeSinceEpoch:     &imageCreationDate,
								LastUpdateTimeSinceEpoch: &imageLastUpdateDate,
							},
						},
					}

					//todo: remove logging for test
					glog.Infof("RHEC response: %+v", getRepositoryQuery)
					updatedModel := models[m.Repository+":"+string(imageTag.Name)]
					maturityVal := "nil"
					if updatedModel.Maturity != nil {
						maturityVal = *updatedModel.Maturity
					}

					glog.Infof("updated model: Name=%s, CreateTime=%s, LastUpdate=%s, Description=%s, Readme=%s, Provider=%s, Maturity=%s",
						updatedModel.Name,
						*updatedModel.CreateTimeSinceEpoch,
						*updatedModel.LastUpdateTimeSinceEpoch,
						*updatedModel.Description,
						*updatedModel.Readme,
						*updatedModel.Provider,
						maturityVal,
					)
				}
			}

		}
	}

	r.modelsLock.Lock()
	defer r.modelsLock.Unlock()
	r.models = models

	return nil
}

const rhecCatalogPath = "yamlCatalogPath"

func newRhecCatalog(source *CatalogSourceConfig) (CatalogSourceProvider, error) {
	rhecModelFile, exists := source.Properties[rhecCatalogPath].(string)
	if !exists || rhecModelFile == "" {
		return nil, fmt.Errorf("missing %s string property", rhecCatalogPath)
	}

	rhecModelFile, err := filepath.Abs(rhecModelFile)
	if err != nil {
		return nil, fmt.Errorf("abs: %w", err)
	}

	p := &rhecCatalogImpl{}
	err = p.load(rhecModelFile)
	if err != nil {
		return nil, err
	}

	go func() {
		changes, err := getMonitor().Path(rhecModelFile)
		if err != nil {
			glog.Errorf("unable to watch RHEC catalog file: %v", err)
			return
		}

		for range changes {
			glog.Infof("Reloading RHEC catalog %s", rhecModelFile)

			err = p.load(rhecModelFile)
			if err != nil {
				glog.Errorf("unable to load RHEC catalog: %v", err)
			}
		}
	}()

	return p, nil
}

func init() {
	if err := RegisterCatalogType("rhec", newRhecCatalog); err != nil {
		panic(err)
	}
}
