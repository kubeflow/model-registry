package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/db/models"
)

const (
	defaultHuggingFaceURL = "https://huggingface.co"
	apiKeyKey             = "apiKey"
	urlKey                = "url"
	modelLimitKey         = "modelLimit"
	includedModelsKey     = "includedModels"
	defaultModelLimit     = 100
)

// hfModel implements apimodels.CatalogModel and populates it from HuggingFace API data
type hfModel struct {
	apimodels.CatalogModel
}

type hfModelProvider struct {
	client         *http.Client
	sourceId       string
	apiKey         string
	baseURL        string
	modelLimit     int
	includedModels []string
	excludedModels []string
}

// hfModelInfo represents the structure of HuggingFace API model information
type hfModelInfo struct {
	ID          string    `json:"id"`
	Author      string    `json:"author,omitempty"`
	Sha         string    `json:"sha,omitempty"`
	CreatedAt   string    `json:"createdAt,omitempty"`
	UpdatedAt   string    `json:"updatedAt,omitempty"`
	Private     bool      `json:"private,omitempty"`
	Downloads   int       `json:"downloads,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	PipelineTag string    `json:"pipeline_tag,omitempty"`
	LibraryName string    `json:"library_name,omitempty"`
	ModelID     string    `json:"modelId,omitempty"`
	Task        string    `json:"task,omitempty"`
	Siblings    []hfFile  `json:"siblings,omitempty"`
	Config      *hfConfig `json:"config,omitempty"`
	CardData    *hfCard   `json:"cardData,omitempty"`
}

type hfFile struct {
	RFileName string `json:"rfilename"`
}

type hfConfig struct {
	Architectures []string `json:"architectures,omitempty"`
	ModelType     string   `json:"model_type,omitempty"`
}

type hfCard struct {
	Data map[string]interface{} `json:"data,omitempty"`
}

// populateFromHFInfo populates the hfModel's CatalogModel fields from HuggingFace API data
func (hfm *hfModel) populateFromHFInfo(hfInfo *hfModelInfo, sourceId string, originalModelName string) {
	// Set model name
	modelName := hfInfo.ID
	if modelName == "" {
		modelName = hfInfo.ModelID
	}
	if modelName == "" {
		modelName = originalModelName
	}
	hfm.Name = modelName

	// Set ExternalId
	if hfInfo.ID != "" {
		hfm.ExternalId = &hfInfo.ID
	}

	// Set SourceId
	if sourceId != "" {
		hfm.SourceId = &sourceId
	}

	// Convert timestamps
	if hfInfo.CreatedAt != "" {
		if createTime, err := parseHFTime(hfInfo.CreatedAt); err == nil {
			createTimeStr := strconv.FormatInt(createTime, 10)
			hfm.CreateTimeSinceEpoch = &createTimeStr
		}
	}
	if hfInfo.UpdatedAt != "" {
		if updateTime, err := parseHFTime(hfInfo.UpdatedAt); err == nil {
			updateTimeStr := strconv.FormatInt(updateTime, 10)
			hfm.LastUpdateTimeSinceEpoch = &updateTimeStr
		}
	}

	// Extract description from cardData if available
	if hfInfo.CardData != nil && hfInfo.CardData.Data != nil {
		if desc, ok := hfInfo.CardData.Data["description"].(string); ok && desc != "" {
			hfm.Description = &desc
		}
		if readme, ok := hfInfo.CardData.Data["readme"].(string); ok && readme != "" {
			hfm.Readme = &readme
		}
		if license, ok := hfInfo.CardData.Data["license"].(string); ok && license != "" {
			hfm.License = &license
		}
		if licenseLink, ok := hfInfo.CardData.Data["license_link"].(string); ok && licenseLink != "" {
			hfm.LicenseLink = &licenseLink
		}
	}

	// Set provider from author
	if hfInfo.Author != "" {
		hfm.Provider = &hfInfo.Author
	}

	// Set library name
	if hfInfo.LibraryName != "" {
		hfm.LibraryName = &hfInfo.LibraryName
	}

	// Convert tasks
	var tasks []string
	if hfInfo.Task != "" {
		tasks = append(tasks, hfInfo.Task)
	}
	if hfInfo.PipelineTag != "" && hfInfo.PipelineTag != hfInfo.Task {
		tasks = append(tasks, hfInfo.PipelineTag)
	}
	if len(tasks) > 0 {
		hfm.Tasks = tasks
	}

	// Convert tags and other metadata to custom properties
	customProps := make(map[string]apimodels.MetadataValue)
	if len(hfInfo.Tags) > 0 {
		if tagsJSON, err := json.Marshal(hfInfo.Tags); err == nil {
			customProps["hf_tags"] = apimodels.MetadataValue{
				MetadataStringValue: &apimodels.MetadataStringValue{
					StringValue: string(tagsJSON),
				},
			}
		}
	}

	if hfInfo.Downloads > 0 {
		customProps["hf_downloads"] = apimodels.MetadataValue{
			MetadataIntValue: &apimodels.MetadataIntValue{
				IntValue: strconv.FormatInt(int64(hfInfo.Downloads), 10),
			},
		}
	}

	if hfInfo.Sha != "" {
		customProps["hf_sha"] = apimodels.MetadataValue{
			MetadataStringValue: &apimodels.MetadataStringValue{
				StringValue: hfInfo.Sha,
			},
		}
	}

	if hfInfo.Config != nil {
		if len(hfInfo.Config.Architectures) > 0 {
			if archJSON, err := json.Marshal(hfInfo.Config.Architectures); err == nil {
				customProps["hf_architectures"] = apimodels.MetadataValue{
					MetadataStringValue: &apimodels.MetadataStringValue{
						StringValue: string(archJSON),
					},
				}
			}
		}
		if hfInfo.Config.ModelType != "" {
			customProps["hf_model_type"] = apimodels.MetadataValue{
				MetadataStringValue: &apimodels.MetadataStringValue{
					StringValue: hfInfo.Config.ModelType,
				},
			}
		}
	}

	if len(customProps) > 0 {
		hfm.SetCustomProperties(customProps)
	}
}

func (p *hfModelProvider) Models(ctx context.Context) (<-chan ModelProviderRecord, error) {
	// Read the catalog and report errors
	catalog, err := p.getModelsFromHF(ctx)
	if err != nil {
		return nil, err
	}

	ch := make(chan ModelProviderRecord)
	go func() {
		defer close(ch)

		// Send the initial list right away.
		p.emit(ctx, catalog, ch)
	}()

	return ch, nil
}

func (p *hfModelProvider) getModelsFromHF(ctx context.Context) ([]ModelProviderRecord, error) {
	var records []ModelProviderRecord

	for _, modelName := range p.includedModels {
		// Skip if excluded
		if isModelExcluded(modelName, p.excludedModels) {
			glog.V(2).Infof("Skipping excluded model: %s", modelName)
			continue
		}

		modelInfo, err := p.fetchModelInfo(ctx, modelName)
		if err != nil {
			glog.Errorf("Failed to fetch model info for %s: %v", modelName, err)
			continue
		}

		record := p.convertHFModelToRecord(modelInfo, modelName)
		records = append(records, record)
	}

	return records, nil
}

func (p *hfModelProvider) fetchModelInfo(ctx context.Context, modelName string) (*hfModelInfo, error) {
	// The HF API requires the full model identifier: org/model-name (aka repo/model-name)

	// Normalize the model name (remove any leading/trailing slashes)
	modelName = strings.Trim(modelName, "/")

	// Construct the API URL with the full model identifier
	apiURL := fmt.Sprintf("%s/api/models/%s", p.baseURL, modelName)

	glog.V(2).Infof("Fetching HuggingFace model info from: %s", apiURL)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent header (HuggingFace API expects this)
	req.Header.Set("User-Agent", "model-registry-catalog")

	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch model info for %s: %w", modelName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HuggingFace API returned status %d for model %s: %s", resp.StatusCode, modelName, string(bodyBytes))
	}

	var modelInfo hfModelInfo
	if err := json.NewDecoder(resp.Body).Decode(&modelInfo); err != nil {
		return nil, fmt.Errorf("failed to decode model info for %s: %w", modelName, err)
	}

	// Ensure ID is set from modelName if not present in API response
	if modelInfo.ID == "" {
		modelInfo.ID = modelName
	}

	return &modelInfo, nil
}

func (p *hfModelProvider) convertHFModelToRecord(hfInfo *hfModelInfo, originalModelName string) ModelProviderRecord {
	// Create hfModel and populate it from HF API data
	hfm := &hfModel{}
	hfm.populateFromHFInfo(hfInfo, p.sourceId, originalModelName)

	// Convert to database model
	model := dbmodels.CatalogModelImpl{}

	// Convert model attributes
	modelName := hfm.Name
	attrs := &dbmodels.CatalogModelAttributes{
		Name:       &modelName,
		ExternalID: hfm.ExternalId,
	}

	// Convert timestamps if available
	if hfm.CreateTimeSinceEpoch != nil {
		if createTime, err := strconv.ParseInt(*hfm.CreateTimeSinceEpoch, 10, 64); err == nil {
			attrs.CreateTimeSinceEpoch = &createTime
		}
	}
	if hfm.LastUpdateTimeSinceEpoch != nil {
		if updateTime, err := strconv.ParseInt(*hfm.LastUpdateTimeSinceEpoch, 10, 64); err == nil {
			attrs.LastUpdateTimeSinceEpoch = &updateTime
		}
	}

	model.Attributes = attrs

	// Convert model properties
	properties, customProperties := convertHFModelProperties(&hfm.CatalogModel)
	if len(properties) > 0 {
		model.Properties = &properties
	}
	if len(customProperties) > 0 {
		model.CustomProperties = &customProperties
	}

	return ModelProviderRecord{
		Model:     &model,
		Artifacts: []dbmodels.CatalogArtifact{}, // HF models don't have artifacts from the API
	}
}

// convertHFModelProperties converts CatalogModel properties to database format
func convertHFModelProperties(catalogModel *apimodels.CatalogModel) ([]models.Properties, []models.Properties) {
	var properties []models.Properties
	var customProperties []models.Properties

	// Regular properties
	if catalogModel.Description != nil {
		properties = append(properties, models.NewStringProperty("description", *catalogModel.Description, false))
	}
	if catalogModel.Readme != nil {
		properties = append(properties, models.NewStringProperty("readme", *catalogModel.Readme, false))
	}
	if catalogModel.Provider != nil {
		properties = append(properties, models.NewStringProperty("provider", *catalogModel.Provider, false))
	}
	if catalogModel.License != nil {
		properties = append(properties, models.NewStringProperty("license", *catalogModel.License, false))
	}
	if catalogModel.LicenseLink != nil {
		properties = append(properties, models.NewStringProperty("license_link", *catalogModel.LicenseLink, false))
	}
	if catalogModel.LibraryName != nil {
		properties = append(properties, models.NewStringProperty("library_name", *catalogModel.LibraryName, false))
	}
	if catalogModel.SourceId != nil {
		properties = append(properties, models.NewStringProperty("source_id", *catalogModel.SourceId, false))
	}

	// Convert array properties
	if len(catalogModel.Tasks) > 0 {
		if tasksJSON, err := json.Marshal(catalogModel.Tasks); err == nil {
			properties = append(properties, models.NewStringProperty("tasks", string(tasksJSON), false))
		}
	}

	// Convert custom properties from the CatalogModel
	if catalogModel.CustomProperties != nil {
		for key, value := range catalogModel.GetCustomProperties() {
			prop := convertMetadataValueToProperty(key, value)
			customProperties = append(customProperties, prop)
		}
	}

	return properties, customProperties
}

// parseHFTime parses HuggingFace timestamp format (ISO 8601)
func parseHFTime(timeStr string) (int64, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return 0, err
	}
	return t.UnixMilli(), nil
}

func (p *hfModelProvider) emit(ctx context.Context, models []ModelProviderRecord, out chan<- ModelProviderRecord) {
	done := ctx.Done()
	for _, model := range models {
		// Check if model should be excluded by name
		if model.Model.GetAttributes() != nil && model.Model.GetAttributes().Name != nil {
			modelName := *model.Model.GetAttributes().Name
			if isModelExcluded(modelName, p.excludedModels) {
				glog.V(2).Infof("Skipping excluded model in emit: %s", modelName)
				continue
			}
		}

		select {
		case out <- model:
		case <-done:
			return
		}
	}
}

func isModelExcluded(modelName string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.HasSuffix(pattern, "*") {
			if strings.HasPrefix(modelName, strings.TrimSuffix(pattern, "*")) {
				return true
			}
		} else if modelName == pattern {
			return true
		}
	}
	return false
}

// validateCredentials checks if the HuggingFace API key credentials are valid
func (p *hfModelProvider) validateCredentials(ctx context.Context) error {
	glog.Infof("Validating HuggingFace API credentials")

	// Make a simple API call to validate credentials
	apiURL := p.baseURL + "/api/whoami-v2"
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create validation request: %w", err)
	}

	req.Header.Set("User-Agent", "model-registry-catalog")

	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to validate HuggingFace credentials: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid HuggingFace API credentials")
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HuggingFace API validation failed with status: %d: %s", resp.StatusCode, string(bodyBytes))
	}

	glog.Infof("HuggingFace credentials validated successfully")
	return nil
}

func newHFModelProvider(ctx context.Context, source *Source, reldir string) (<-chan ModelProviderRecord, error) {
	p := &hfModelProvider{}
	p.client = &http.Client{Timeout: 30 * time.Second}

	// Parse Source ID
	sourceId := source.GetId()
	if sourceId == "" {
		return nil, fmt.Errorf("missing source ID for HuggingFace catalog")
	}
	p.sourceId = sourceId

	// Parse API key
	apiKey, ok := source.Properties[apiKeyKey].(string)
	if !ok || apiKey == "" {
		return nil, fmt.Errorf("missing or invalid '%s' property for HuggingFace catalog", apiKeyKey)
	}
	p.apiKey = apiKey

	// Parse base URL (optional, defaults to huggingface.co)
	// This allows tests to use mock servers by providing a custom URL
	p.baseURL = defaultHuggingFaceURL
	if url, ok := source.Properties[urlKey].(string); ok && url != "" {
		p.baseURL = strings.TrimSuffix(url, "/") // Remove trailing slash if present
	}

	// Validate credentials before proceeding
	if err := p.validateCredentials(ctx); err != nil {
		glog.Errorf("HuggingFace catalog credential validation failed: %v", err)
		return nil, fmt.Errorf("failed to validate HuggingFace catalog credentials: %w", err)
	}

	// Parse includedModels
	includedModels, ok := source.Properties[includedModelsKey].([]any)
	if !ok {
		return nil, fmt.Errorf("%q property should be a list", includedModelsKey)
	}

	if len(includedModels) == 0 {
		return nil, fmt.Errorf("%q property cannot be empty", includedModelsKey)
	}

	modelStrings := make([]string, 0, len(includedModels))
	for _, model := range includedModels {
		modelStr, ok := model.(string)
		if !ok {
			return nil, fmt.Errorf("%s: invalid list: expected string, got %T", includedModelsKey, model)
		}
		modelStrings = append(modelStrings, modelStr)
	}

	p.includedModels = modelStrings

	// Parse excludedModels
	p.excludedModels = []string{}
	if _, exists := source.Properties[excludedModelsKey]; exists {
		excludedModels, ok := source.Properties[excludedModelsKey].([]any)
		if !ok {
			return nil, fmt.Errorf("%q property should be a list", excludedModelsKey)
		}

		p.excludedModels = make([]string, 0, len(excludedModels))
		for _, name := range excludedModels {
			nameStr, ok := name.(string)
			if !ok {
				return nil, fmt.Errorf("%s: invalid list: wanted string, got %T", excludedModelsKey, name)
			}
			p.excludedModels = append(p.excludedModels, nameStr)
		}
	}

	// TODO: Implement model limit when organizations level includedModels are supported
	if limit, ok := source.Properties[modelLimitKey].(int); ok && limit > 0 {
		p.modelLimit = limit
		glog.Infof("Configuring HuggingFace catalog with URL: %s, modelLimit: %d", p.baseURL, limit)
		// Note: modelLimit is stored but not currently enforced in the current implementation
		// This could be used to limit the number of models fetched from includedModels
	} else {
		p.modelLimit = defaultModelLimit
		glog.Infof("Configuring HuggingFace catalog with URL: %s, %d included model(s)", p.baseURL, len(p.includedModels))
	}

	return p.Models(ctx)
}

func init() {
	if err := RegisterModelProvider("hf", newHFModelProvider); err != nil {
		panic(err)
	}
}
