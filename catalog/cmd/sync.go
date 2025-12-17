package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/spf13/cobra"
)

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync Hugging Face model metadata",
	Long: `Sync metadata for Hugging Face models already in the catalog.
This command fetches the latest metadata from Hugging Face for models
already present in the catalog, compares it with existing data, and
updates models if changes are detected.`,
	RunE: runSync,
}

func runSync(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	ds, err := datastore.NewConnector("embedmd", &embedmd.EmbedMDConfig{
		DatabaseType: "postgres",
		DatabaseDSN:  "",
	})
	if err != nil {
		return fmt.Errorf("error creating datastore: %w", err)
	}

	repoSet, err := ds.Connect(service.DatastoreSpec())
	if err != nil {
		return fmt.Errorf("error initializing datastore: %v", err)
	}

	services := service.NewServices(
		getRepo[models.CatalogModelRepository](repoSet),
		getRepo[models.CatalogArtifactRepository](repoSet),
		getRepo[models.CatalogModelArtifactRepository](repoSet),
		getRepo[models.CatalogMetricsArtifactRepository](repoSet),
		getRepo[models.CatalogSourceRepository](repoSet),
		getRepo[models.PropertyOptionsRepository](repoSet),
	)

	// Sync all Hugging Face models in the catalog
	// The sync function will get all models and search them on Hugging Face
	glog.Infof("Starting sync for all models in catalog")

	err = catalog.SyncHuggingFaceModels(ctx, services)
	if err != nil {
		return fmt.Errorf("error syncing models: %w", err)
	}

	glog.Infof("Sync completed")
	return nil
}
