package catalog

import (
	"context"
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
	Gated       string    `json:"gated,omitempty"`
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

var (
	catalogModelLogo = "data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48c3ZnIGlkPSJ1dWlkLWE2ODNjMGQyLWViOTAtNGI0Yi1hOWE0LTVlMzI2NWYwNzVjNSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB3aWR0aD0iMTc1IiBoZWlnaHQ9IjE3NSIgdmlld0JveD0iMCAwIDE3NSAxNzUiPjxkZWZzPjxzdHlsZT4udXVpZC00MmIxYTFkMi05NWQ4LTQ3OTgtYmY5ZS04YTc1NzAyOGM0NWJ7ZmlsbDpub25lO30udXVpZC01MjRiZWEwYi1mMTg2LTQ1ZjAtOWQ4Ny01ZTBhNjlkZjZhYTN7ZmlsbDojZTBlMGUwO30udXVpZC02Mzc3YmY0MS04OTU5LTQyN2ItYWFlZC0xOWRjMDlmNzM1MGV7ZmlsbDojZmZmO30udXVpZC0xNTQyMDVmNy02YzY5LTQyODYtYjllZC1hZDcwMjI0YzdiODR7ZmlsbDojZTAwO308L3N0eWxlPjwvZGVmcz48cmVjdCBjbGFzcz0idXVpZC00MmIxYTFkMi05NWQ4LTQ3OTgtYmY5ZS04YTc1NzAyOGM0NWIiIHdpZHRoPSIxNzUiIGhlaWdodD0iMTc1Ii8+PHJlY3QgY2xhc3M9InV1aWQtNjM3N2JmNDEtODk1OS00MjdiLWFhZWQtMTlkYzA5ZjczNTBlIiB4PSIxMi41IiB5PSIxMi41IiB3aWR0aD0iMTUwIiBoZWlnaHQ9IjE1MCIgcng9IjM3LjUiIHJ5PSIzNy41Ii8+PHBhdGggY2xhc3M9InV1aWQtNTI0YmVhMGItZjE4Ni00NWYwLTlkODctNWUwYTY5ZGY2YWEzIiBkPSJNMTI1LDE3LjcwODMzYzE3LjgwNTY3LDAsMzIuMjkxNjcsMTQuNDg2LDMyLjI5MTY3LDMyLjI5MTY3djc1YzAsMTcuODA1NjYtMTQuNDg1OTksMzIuMjkxNjctMzIuMjkxNjcsMzIuMjkxNjdINTBjLTE3LjgwNTY3LDAtMzIuMjkxNjctMTQuNDg2LTMyLjI5MTY3LTMyLjI5MTY3VjUwYzAtMTcuODA1NjcsMTQuNDg1OTktMzIuMjkxNjcsMzIuMjkxNjctMzIuMjkxNjdoNzVNMTI1LDEyLjVINTBjLTIwLjcxMDY3LDAtMzcuNSwxNi43ODkzMy0zNy41LDM3LjV2NzVjMCwyMC43MTA2NywxNi43ODkzMywzNy41LDM3LjUsMzcuNWg3NWMyMC43MTA2NiwwLDM3LjUtMTYuNzg5MzMsMzcuNS0zNy41VjUwYzAtMjAuNzEwNjctMTYuNzg5MzQtMzcuNS0zNy41LTM3LjVoMFoiLz48cGF0aCBjbGFzcz0idXVpZC0xNTQyMDVmNy02YzY5LTQyODYtYjllZC1hZDcwMjI0YzdiODQiIGQ9Ik0xMTYuNjY2NjIsMTE5LjI3MDgzaC01OC4zMzMzM2MtMS4wNTM4NywwLTIuMDAxOTUtLjYzNDc3LTIuNDA2ODItMS42MDcyNi0uNDAyODMtLjk3MjQ5LS4xNzkwNC0yLjA5MzUxLjU2NTU5LTIuODM4MTNsMjcuMzI1NDQtMjcuMzI1NDQtMjcuMzI1NDQtMjcuMzI1NDRjLS43NDQ2My0uNzQ0NjMtLjk2ODQyLTEuODY1NjQtLjU2NTU5LTIuODM4MTMuNDA0ODctLjk3MjQ5LDEuMzUyOTUtMS42MDcyNiwyLjQwNjgyLTEuNjA3MjZoNTguMzMzMzNjMS4wNTM4NywwLDIuMDAxOTUuNjM0NzcsMi40MDY4MiwxLjYwNzI2LjQwMjgzLjk3MjQ5LjE3OTA0LDIuMDkzNTEtLjU2NTU5LDIuODM4MTNsLTI3LjMyNTQ0LDI3LjMyNTQ0LDI3LjMyNTQ0LDI3LjMyNTQ0Yy43NDQ2My43NDQ2My45Njg0MiwxLjg2NTY0LjU2NTU5LDIuODM4MTMtLjQwNDg3Ljk3MjQ5LTEuMzUyOTUsMS42MDcyNi0yLjQwNjgyLDEuNjA3MjZaTTY0LjYxOTkxLDExNC4wNjI1aDQ1Ljc2MDA5bC0yMi44ODAwNS0yMi44ODAwNS0yMi44ODAwNSwyMi44ODAwNVpNNjQuNjE5OTEsNjAuOTM3NWwyMi44ODAwNSwyMi44ODAwNSwyMi44ODAwNS0yMi44ODAwNWgtNDUuNzYwMDlaIi8+PHBhdGggY2xhc3M9InV1aWQtMTU0MjA1ZjctNmM2OS00Mjg2LWI5ZWQtYWQ3MDIyNGM3Yjg0IiBkPSJNMTE2LjY2NjYyLDkwLjEwNDE3aC01OC4zMzMzM2MtMS40Mzg0LDAtMi42MDQxNy0xLjE2NTc3LTIuNjA0MTctMi42MDQxN3MxLjE2NTc3LTIuNjA0MTcsMi42MDQxNy0yLjYwNDE3aDU4LjMzMzMzYzEuNDM4NCwwLDIuNjA0MTcsMS4xNjU3NywyLjYwNDE3LDIuNjA0MTdzLTEuMTY1NzcsMi42MDQxNy0yLjYwNDE3LDIuNjA0MTdaIi8+PHJlY3QgY2xhc3M9InV1aWQtNjM3N2JmNDEtODk1OS00MjdiLWFhZWQtMTlkYzA5ZjczNTBlIiB4PSIxMTIuNSIgeT0iMTEyLjUiIHdpZHRoPSIxMi41IiBoZWlnaHQ9IjEyLjUiLz48cGF0aCBkPSJNMTI0Ljk5OTk1LDEyNy42MDQxN2gtMTIuNWMtMS40Mzg0LDAtMi42MDQxNy0xLjE2NTc3LTIuNjA0MTctMi42MDQxN3YtMTIuNWMwLTEuNDM4NCwxLjE2NTc3LTIuNjA0MTcsMi42MDQxNy0yLjYwNDE3aDEyLjVjMS40Mzg0LDAsMi42MDQxNywxLjE2NTc3LDIuNjA0MTcsMi42MDQxN3YxMi41YzAsMS40Mzg0LTEuMTY1NzcsMi42MDQxNy0yLjYwNDE3LDIuNjA0MTdaTTExNS4xMDQxMiwxMjIuMzk1ODNoNy4yOTE2N3YtNy4yOTE2N2gtNy4yOTE2N3Y3LjI5MTY3WiIvPjxyZWN0IGNsYXNzPSJ1dWlkLTYzNzdiZjQxLTg5NTktNDI3Yi1hYWVkLTE5ZGMwOWY3MzUwZSIgeD0iNTAiIHk9IjExMi41IiB3aWR0aD0iMTIuNSIgaGVpZ2h0PSIxMi41Ii8+PHBhdGggZD0iTTYyLjQ5OTk1LDEyNy42MDQxN2gtMTIuNWMtMS40Mzg0LDAtMi42MDQxNy0xLjE2NTc3LTIuNjA0MTctMi42MDQxN3YtMTIuNWMwLTEuNDM4NCwxLjE2NTc3LTIuNjA0MTcsMi42MDQxNy0yLjYwNDE3aDEyLjVjMS40Mzg0LDAsMi42MDQxNywxLjE2NTc3LDIuNjA0MTcsMi42MDQxN3YxMi41YzAsMS40Mzg0LTEuMTY1NzcsMi42MDQxNy0yLjYwNDE3LDIuNjA0MTdaTTUyLjYwNDEyLDEyMi4zOTU4M2g3LjI5MTY3di03LjI5MTY3aC03LjI5MTY3djcuMjkxNjdaIi8+PHJlY3QgY2xhc3M9InV1aWQtNjM3N2JmNDEtODk1OS00MjdiLWFhZWQtMTlkYzA5ZjczNTBlIiB4PSIxMTIuNSIgeT0iNTAiIHdpZHRoPSIxMi41IiBoZWlnaHQ9IjEyLjUiLz48cGF0aCBkPSJNMTI0Ljk5OTk1LDY1LjEwNDE3aC0xMi41Yy0xLjQzODQsMC0yLjYwNDE3LTEuMTY1NzctMi42MDQxNy0yLjYwNDE3di0xMi41YzAtMS40Mzg0LDEuMTY1NzctMi42MDQxNywyLjYwNDE3LTIuNjA0MTdoMTIuNWMxLjQzODQsMCwyLjYwNDE3LDEuMTY1NzcsMi42MDQxNywyLjYwNDE3djEyLjVjMCwxLjQzODQtMS4xNjU3NywyLjYwNDE3LTIuNjA0MTcsMi42MDQxN1pNMTE1LjEwNDEyLDU5Ljg5NTgzaDcuMjkxNjd2LTcuMjkxNjdoLTcuMjkxNjd2Ny4yOTE2N1oiLz48cmVjdCBjbGFzcz0idXVpZC02Mzc3YmY0MS04OTU5LTQyN2ItYWFlZC0xOWRjMDlmNzM1MGUiIHg9IjUwIiB5PSI1MCIgd2lkdGg9IjEyLjUiIGhlaWdodD0iMTIuNSIvPjxwYXRoIGQ9Ik02Mi40OTk5NSw2NS4xMDQxN2gtMTIuNWMtMS40Mzg0LDAtMi42MDQxNy0xLjE2NTc3LTIuNjA0MTctMi42MDQxN3YtMTIuNWMwLTEuNDM4NCwxLjE2NTc3LTIuNjA0MTcsMi42MDQxNy0yLjYwNDE3aDEyLjVjMS40Mzg0LDAsMi42MDQxNywxLjE2NTc3LDIuNjA0MTcsMi42MDQxN3YxMi41YzAsMS40Mzg0LTEuMTY1NzcsMi42MDQxNy0yLjYwNDE3LDIuNjA0MTdaTTUyLjYwNDEyLDU5Ljg5NTgzaDcuMjkxNjd2LTcuMjkxNjdoLTcuMjkxNjd2Ny4yOTE2N1oiLz48cmVjdCBjbGFzcz0idXVpZC02Mzc3YmY0MS04OTU5LTQyN2ItYWFlZC0xOWRjMDlmNzM1MGUiIHg9IjUwIiB5PSI4MS4yNSIgd2lkdGg9IjEyLjUiIGhlaWdodD0iMTIuNSIvPjxwYXRoIGQ9Ik02Mi40OTk5NSw5Ni4zNTQxN2gtMTIuNWMtMS40Mzg0LDAtMi42MDQxNy0xLjE2NTc3LTIuNjA0MTctMi42MDQxN3YtMTIuNWMwLTEuNDM4NCwxLjE2NTc3LTIuNjA0MTcsMi42MDQxNy0yLjYwNDE3aDEyLjVjMS40Mzg0LDAsMi42MDQxNywxLjE2NTc3LDIuNjA0MTcsMi42MDQxN3YxMi41YzAsMS40Mzg0LTEuMTY1NzcsMi42MDQxNy0yLjYwNDE3LDIuNjA0MTdaTTUyLjYwNDEyLDkxLjE0NTgzaDcuMjkxNjd2LTcuMjkxNjdoLTcuMjkxNjd2Ny4yOTE2N1oiLz48cmVjdCBjbGFzcz0idXVpZC02Mzc3YmY0MS04OTU5LTQyN2ItYWFlZC0xOWRjMDlmNzM1MGUiIHg9IjExMi41IiB5PSI4MS4yNSIgd2lkdGg9IjEyLjUiIGhlaWdodD0iMTIuNSIvPjxwYXRoIGQ9Ik0xMjQuOTk5OTUsOTYuMzU0MTdoLTEyLjVjLTEuNDM4NCwwLTIuNjA0MTctMS4xNjU3Ny0yLjYwNDE3LTIuNjA0MTd2LTEyLjVjMC0xLjQzODQsMS4xNjU3Ny0yLjYwNDE3LDIuNjA0MTctMi42MDQxN2gxMi41YzEuNDM4NCwwLDIuNjA0MTcsMS4xNjU3NywyLjYwNDE3LDIuNjA0MTd2MTIuNWMwLDEuNDM4NC0xLjE2NTc3LDIuNjA0MTctMi42MDQxNywyLjYwNDE3Wk0xMTUuMTA0MTIsOTEuMTQ1ODNoNy4yOTE2N3YtNy4yOTE2N2gtNy4yOTE2N3Y3LjI5MTY3WiIvPjxyZWN0IGNsYXNzPSJ1dWlkLTYzNzdiZjQxLTg5NTktNDI3Yi1hYWVkLTE5ZGMwOWY3MzUwZSIgeD0iODEuMjUiIHk9IjgxLjI1IiB3aWR0aD0iMTIuNSIgaGVpZ2h0PSIxMi41Ii8+PHBhdGggZD0iTTkzLjc0OTk1LDk2LjM1NDE3aC0xMi41Yy0xLjQzODQsMC0yLjYwNDE3LTEuMTY1NzctMi42MDQxNy0yLjYwNDE3di0xMi41YzAtMS40Mzg0LDEuMTY1NzctMi42MDQxNywyLjYwNDE3LTIuNjA0MTdoMTIuNWMxLjQzODQsMCwyLjYwNDE3LDEuMTY1NzcsMi42MDQxNywyLjYwNDE3djEyLjVjMCwxLjQzODQtMS4xNjU3NywyLjYwNDE3LTIuNjA0MTcsMi42MDQxN1pNODMuODU0MTIsOTEuMTQ1ODNoNy4yOTE2N3YtNy4yOTE2N2gtNy4yOTE2N3Y3LjI5MTY3WiIvPjwvc3ZnPg=="
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
			StringValue: hfInfo.Gated,
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
		if isModelExcluded(modelName, p.excludedModels) {
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
			if isModelExcluded(finalModelName, p.excludedModels) {
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

	// Parse API key from environment variable
	apiKey := os.Getenv("HF_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing HF_API_KEY environment variable for HuggingFace catalog")
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
		if len(p.excludedModels) > 0 {
			glog.Infof("Configured %d excluded model pattern(s) for HuggingFace catalog", len(p.excludedModels))
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
