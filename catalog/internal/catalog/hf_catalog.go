package catalog

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/pkg/openapi"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

type hfCatalogImpl struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

var _ CatalogSourceProvider = &hfCatalogImpl{}

const (
	defaultHuggingFaceURL = "https://huggingface.co"
)

func (h *hfCatalogImpl) GetModel(ctx context.Context, name string) (*openapi.CatalogModel, error) {
	// TODO: Implement HuggingFace model retrieval
	return nil, fmt.Errorf("HuggingFace model retrieval not yet implemented")
}

func (h *hfCatalogImpl) ListModels(ctx context.Context, params ListModelsParams) (model.CatalogModelList, error) {
	// TODO: Implement HuggingFace model listing
	// For now, return empty list to satisfy interface
	return model.CatalogModelList{
		Items:    []model.CatalogModel{},
		PageSize: 0,
		Size:     0,
	}, nil
}

func (h *hfCatalogImpl) GetArtifacts(ctx context.Context, name string) (*openapi.CatalogModelArtifactList, error) {
	// TODO: Implement HuggingFace model artifacts retrieval
	// For now, return empty list to satisfy interface
	return &openapi.CatalogModelArtifactList{
		Items:    []openapi.CatalogModelArtifact{},
		PageSize: 0,
		Size:     0,
	}, nil
}

// validateCredentials checks if the HuggingFace API credentials are valid
func (h *hfCatalogImpl) validateCredentials(ctx context.Context) error {
	glog.Infof("Validating HuggingFace API credentials")

	// Make a simple API call to validate credentials
	apiURL := h.baseURL + "/api/whoami-v2"
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create validation request: %w", err)
	}

	if h.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+h.apiKey)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to validate HuggingFace credentials: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid HuggingFace API credentials")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HuggingFace API validation failed with status: %d", resp.StatusCode)
	}

	glog.Infof("HuggingFace credentials validated successfully")
	return nil
}

// newHfCatalog creates a new HuggingFace catalog source
func newHfCatalog(source *CatalogSourceConfig) (CatalogSourceProvider, error) {
	apiKey, ok := source.Properties["apiKey"].(string)
	if !ok || apiKey == "" {
		return nil, fmt.Errorf("missing or invalid 'apiKey' property for HuggingFace catalog")
	}

	baseURL := defaultHuggingFaceURL
	if url, ok := source.Properties["url"].(string); ok && url != "" {
		baseURL = strings.TrimSuffix(url, "/")
	}

	// Optional model limit for future implementation
	modelLimit := 100
	if limit, ok := source.Properties["modelLimit"].(int); ok && limit > 0 {
		modelLimit = limit
	}

	glog.Infof("Configuring HuggingFace catalog with URL: %s, modelLimit: %d", baseURL, modelLimit)

	h := &hfCatalogImpl{
		client:  &http.Client{Timeout: 30 * time.Second},
		apiKey:  apiKey,
		baseURL: baseURL,
	}

	// Validate credentials during initialization (as required by Jira ticket)
	ctx := context.Background()
	if err := h.validateCredentials(ctx); err != nil {
		glog.Errorf("HuggingFace catalog credential validation failed: %v", err)
		return nil, fmt.Errorf("failed to validate HuggingFace catalog credentials: %w", err)
	}

	glog.Infof("HuggingFace catalog source configured successfully")
	return h, nil
}

func init() {
	if err := RegisterCatalogType("hf", newHfCatalog); err != nil {
		panic(err)
	}
}
