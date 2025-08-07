package models

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
)

type CatalogModel struct {
	CreateTimeSinceEpoch     *string                            `json:"createTimeSinceEpoch,omitempty"`
	CustomProperties         *map[string]*openapi.MetadataValue `json:"customProperties,omitempty"`
	Description              *string                            `json:"description,omitempty"`
	Language                 []string                           `json:"language,omitempty"`
	LastUpdateTimeSinceEpoch *string                            `json:"lastUpdateTimeSinceEpoch,omitempty"`
	LibraryName              *string                            `json:"libraryName,omitempty"`
	License                  *string                            `json:"license,omitempty"`
	LicenseLink              *string                            `json:"licenseLink,omitempty"`
	Logo                     *string                            `json:"logo,omitempty"`
	Maturity                 *string                            `json:"maturity,omitempty"`
	Name                     string                             `json:"name"`
	Provider                 *string                            `json:"provider,omitempty"`
	Readme                   *string                            `json:"readme,omitempty"`
	SourceId                 *string                            `json:"source_id,omitempty"`
	Tasks                    []string                           `json:"tasks,omitempty"`
}

type CatalogModelList struct {
	NextPageToken string         `json:"nextPageToken"`
	PageSize      int32          `json:"pageSize"`
	Size          int32          `json:"size"`
	Items         []CatalogModel `json:"items"`
}
