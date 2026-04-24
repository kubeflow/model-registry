package cmd

import (
	"github.com/kubeflow/hub/catalog/cmd"
)

func init() {
	rootCmd.AddCommand(cmd.CatalogCmd)
}
