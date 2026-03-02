package basecatalog

import (
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/golang/glog"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	mrmodels "github.com/kubeflow/model-registry/internal/db/models"
)

// CleanupOrphanedCatalogSources removes CatalogSource records for sources that are no
// longer in the config. currentSourceIDs is the union of all active source IDs from all
// loaders; any source not in this set is deleted.
func CleanupOrphanedCatalogSources(repo dbmodels.CatalogSourceRepository, currentSourceIDs mapset.Set[string]) error {
	existingSources, err := repo.GetAll()
	if err != nil {
		return fmt.Errorf("unable to get existing catalog sources: %w", err)
	}

	for _, source := range existingSources {
		attrs := source.GetAttributes()
		if attrs == nil || attrs.Name == nil {
			continue
		}

		sourceID := *attrs.Name
		if !currentSourceIDs.Contains(sourceID) {
			glog.Infof("Removing orphaned catalog source %s (no longer in config)", sourceID)
			if delErr := repo.Delete(sourceID); delErr != nil {
				glog.Errorf("failed to delete orphaned catalog source %s: %v", sourceID, delErr)
			}
		}
	}

	return nil
}

// SaveSourceStatus persists the operational status of a catalog source to the database.
// Valid status values are the SourceStatus* constants defined in this package.
func SaveSourceStatus(repo dbmodels.CatalogSourceRepository, sourceID, status, errorMsg string) {
	switch status {
	case SourceStatusAvailable, SourceStatusPartiallyAvailable, SourceStatusError, SourceStatusDisabled:
		// valid
	default:
		glog.Errorf("invalid status %q for source %s", status, sourceID)
		return
	}

	source := &dbmodels.CatalogSourceImpl{
		Attributes: &dbmodels.CatalogSourceAttributes{
			Name: &sourceID,
		},
	}

	props := []mrmodels.Properties{
		mrmodels.NewStringProperty("status", status, false),
	}

	if errorMsg != "" {
		props = append(props, mrmodels.NewStringProperty("error", errorMsg, false))
	}

	source.Properties = &props

	if _, err := repo.Save(source); err != nil {
		glog.Errorf("failed to save status for source %s: %v", sourceID, err)
	}
}
