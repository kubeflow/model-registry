package cmd

import (
	"github.com/kubeflow/model-registry/catalog/cmd"
)

func init() {
	rootCmd.AddCommand(cmd.CatalogCmd)
}
