package modelcatalog

import (
	"strings"

	"github.com/kubeflow/model-registry/catalog/internal/catalog/modelcatalog/models"
)

// DisplayNameFromStoredName returns the model display name from the stored (namespaced) name.
// Stored names use the format "sourceId:modelName" for DB uniqueness; this strips the prefix
// so callers get the name without the source id prepended.
func DisplayNameFromStoredName(storedName string) string {
	if idx := strings.Index(storedName, ":"); idx >= 0 {
		return storedName[idx+1:]
	}
	return storedName
}

// DisplayName returns the display name for a catalog model (name without "sourceId:" prefix).
// Use this when reading entities from the DB to get the user-facing model name.
func DisplayName(m models.CatalogModel) string {
	if m == nil || m.GetAttributes() == nil || m.GetAttributes().Name == nil {
		return ""
	}
	return DisplayNameFromStoredName(*m.GetAttributes().Name)
}
