package catalog

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	defaultAPIKeyEnvVar   = "HF_API_KEY"
	urlKey                = "url"
	apiKeyEnvVarKey       = "apiKeyEnvVar"
)

// gatedString is a custom type that can unmarshal both boolean and string values from JSON
// It converts booleans to strings (false -> "false", true -> "true")
type gatedString string

// UnmarshalJSON implements json.Unmarshaler to handle both boolean and string values
func (g *gatedString) UnmarshalJSON(data []byte) error {
	// Handle null/empty
	if len(data) == 0 || string(data) == "null" {
		*g = gatedString("")
		return nil
	}

	// Try to unmarshal as boolean first (handles true/false)
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		*g = gatedString(strconv.FormatBool(b))
		return nil
	}

	// If not a boolean, try as string (handles quoted strings)
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("gated field must be boolean or string, got: %s", string(data))
	}
	*g = gatedString(s)
	return nil
}

// String returns the string value
func (g gatedString) String() string {
	return string(g)
}

// hfModel implements apimodels.CatalogModel and populates it from HuggingFace API data
type hfModel struct {
	apimodels.CatalogModel
}

type hfModelProvider struct {
	client         *http.Client
	sourceId       string
	apiKey         string
	baseURL        string
	includedModels []string
	filter         *ModelFilter
}

// hfModelInfo represents the structure of HuggingFace API model information
type hfModelInfo struct {
	ID          string      `json:"id"`
	Author      string      `json:"author,omitempty"`
	Sha         string      `json:"sha,omitempty"`
	CreatedAt   string      `json:"createdAt,omitempty"`
	UpdatedAt   string      `json:"updatedAt,omitempty"`
	Private     bool        `json:"private,omitempty"`
	Gated       gatedString `json:"gated,omitempty"`
	Downloads   int         `json:"downloads,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	PipelineTag string      `json:"pipeline_tag,omitempty"`
	LibraryName string      `json:"library_name,omitempty"`
	ModelID     string      `json:"modelId,omitempty"`
	Task        string      `json:"task,omitempty"`
	Siblings    []hfFile    `json:"siblings,omitempty"`
	Config      *hfConfig   `json:"config,omitempty"`
	CardData    *hfCard     `json:"cardData,omitempty"`
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

//go:embed assets/catalog_logo.svg
var catalogLogoSVG []byte

var (
	catalogModelLogo = "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(catalogLogoSVG)
)

// populateFromHFInfo populates the hfModel's CatalogModel fields from HuggingFace API data
func (hfm *hfModel) populateFromHFInfo(ctx context.Context, provider *hfModelProvider, hfInfo *hfModelInfo, sourceId string, originalModelName string) {
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

	// Extract license from tags
	// Skip license tags in custom properties to avoid duplication
	var filteredTags []string
	if len(hfInfo.Tags) > 0 {
		filteredTags = make([]string, 0, len(hfInfo.Tags))
		for _, tag := range hfInfo.Tags {
			if strings.HasPrefix(tag, "license:") {
				// Extract license (only first one)
				if hfm.License == nil {
					license := strings.TrimPrefix(tag, "license:")
					if license != "" {
						hfm.License = &license
					}
				}
			} else {
				filteredTags = append(filteredTags, tag)
			}
		}
	}

	// Extract README from sibling files first (preferred source)
	// Check for common README filenames
	readmeFilenames := []string{"README.md", "readme.md", "Readme.md", "README", "readme"}

	for _, sibling := range hfInfo.Siblings {
		for _, readmeFilename := range readmeFilenames {
			if sibling.RFileName == readmeFilename {
				if readmeContent, err := provider.fetchFileContent(ctx, modelName, readmeFilename); err == nil {
					hfm.Readme = &readmeContent
					break
				} else {
					glog.V(2).Infof("Failed to fetch README from sibling file %s for model %s: %v", readmeFilename, modelName, err)
				}
			}
		}
		if hfm.Readme != nil {
			break
		}
	}

	// Extract description from cardData if available
	if hfInfo.CardData != nil && hfInfo.CardData.Data != nil {
		// Extract description from cardData if available
		if desc, ok := hfInfo.CardData.Data["description"].(string); ok && desc != "" {
			hfm.Description = &desc
		}

		// Extract language from cardData if available
		if langData, ok := hfInfo.CardData.Data["language"].([]interface{}); ok && len(langData) > 0 {
			languages := make([]string, 0, len(langData))
			for _, lang := range langData {
				if langStr, ok := lang.(string); ok && langStr != "" {
					languages = append(languages, langStr)
				}
			}
			if len(languages) > 0 {
				hfm.Language = languages
			}
		}

		// Extract license link from cardData if available
		// Check common field names for license link/URL
		if hfm.LicenseLink == nil {
			licenseLinkFields := []string{"license_link", "licenseLink", "license_url", "licenseUrl", "license"}
			for _, field := range licenseLinkFields {
				if link, ok := hfInfo.CardData.Data[field].(string); ok && link != "" {
					if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
						hfm.LicenseLink = &link
						break
					}
				}
			}
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

	// Set logo
	hfm.Logo = &catalogModelLogo

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

	customProps["hf_private"] = apimodels.MetadataValue{
		MetadataStringValue: &apimodels.MetadataStringValue{
			StringValue: strconv.FormatBool(hfInfo.Private),
		},
	}

	customProps["hf_gated"] = apimodels.MetadataValue{
		MetadataStringValue: &apimodels.MetadataStringValue{
			StringValue: hfInfo.Gated.String(),
		},
	}

	if len(filteredTags) > 0 {
		if tagsJSON, err := json.Marshal(filteredTags); err == nil {
			customProps["hf_tags"] = apimodels.MetadataValue{
				MetadataStringValue: &apimodels.MetadataStringValue{
					StringValue: string(tagsJSON),
				},
			}
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
		// Skip if excluded - check before fetching to avoid unnecessary API calls
		if !p.filter.Allows(modelName) {
			glog.V(2).Infof("Skipping excluded model: %s", modelName)
			continue
		}

		modelInfo, err := p.fetchModelInfo(ctx, modelName)
		if err != nil {
			glog.Errorf("Failed to fetch model info for %s: %v", modelName, err)
			continue
		}

		record := p.convertHFModelToRecord(ctx, modelInfo, modelName)

		// Additional safety check: verify the final model name is not excluded
		// (in case the model name changed during conversion, e.g., from hfInfo.ID)
		if record.Model.GetAttributes() != nil && record.Model.GetAttributes().Name != nil {
			finalModelName := *record.Model.GetAttributes().Name
			if !p.filter.Allows(finalModelName) {
				glog.V(2).Infof("Skipping excluded model (after conversion): %s", finalModelName)
				continue
			}
		}

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

// fetchFileContent fetches the content of a file from HuggingFace repository
func (p *hfModelProvider) fetchFileContent(ctx context.Context, modelName string, filename string) (string, error) {
	// Normalize the model name (remove any leading/trailing slashes)
	modelName = strings.Trim(modelName, "/")

	// Construct the API URL for raw file content
	// HuggingFace API endpoint: {baseURL}/{model_id}/raw/main/{filename}
	apiURL := fmt.Sprintf("%s/%s/raw/main/%s", p.baseURL, modelName, filename)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent header
	req.Header.Set("User-Agent", "model-registry-catalog")

	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file %s for model %s: %w", filename, modelName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HuggingFace API returned status %d for file %s in model %s: %s", resp.StatusCode, filename, modelName, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read file content for %s in model %s: %w", filename, modelName, err)
	}

	return string(bodyBytes), nil
}

func (p *hfModelProvider) convertHFModelToRecord(ctx context.Context, hfInfo *hfModelInfo, originalModelName string) ModelProviderRecord {
	// Create hfModel and populate it from HF API data
	hfm := &hfModel{}
	hfm.populateFromHFInfo(ctx, p, hfInfo, p.sourceId, originalModelName)

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
	if catalogModel.Logo != nil {
		properties = append(properties, models.NewStringProperty("logo", *catalogModel.Logo, false))
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
	if len(catalogModel.Language) > 0 {
		if languageJSON, err := json.Marshal(catalogModel.Language); err == nil {
			properties = append(properties, models.NewStringProperty("language", string(languageJSON), false))
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
			if !p.filter.Allows(modelName) {
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

	// Parse API key from environment variable
	// Allow the environment variable name to be configured via properties, defaulting to HF_API_KEY
	apiKeyEnvVar := defaultAPIKeyEnvVar
	if envVar, ok := source.Properties[apiKeyEnvVarKey].(string); ok && envVar != "" {
		apiKeyEnvVar = envVar
	}
	apiKey := os.Getenv(apiKeyEnvVar)
	if apiKey == "" {
		return nil, fmt.Errorf("missing %s environment variable for HuggingFace catalog", apiKeyEnvVar)
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

	// Use top-level IncludedModels from Source as the list of models to fetch
	// These can be specific model names (required for HF API) or patterns
	if len(source.IncludedModels) == 0 {
		return nil, fmt.Errorf("includedModels cannot be empty for HuggingFace catalog")
	}

	p.includedModels = source.IncludedModels

	// Create ModelFilter from source configuration (handles IncludedModels/ExcludedModels from Source)
	// Note: IncludedModels are used both for fetching and filtering
	filter, err := NewModelFilterFromSource(source, nil, nil)
	if err != nil {
		return nil, err
	}
	p.filter = filter

	return p.Models(ctx)
}

func init() {
	if err := RegisterModelProvider("hf", newHFModelProvider); err != nil {
		panic(err)
	}
}
