package catalog

import (
	"reflect"
	"sort"
	"testing"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

func TestLoadCatalogSources(t *testing.T) {
	type args struct {
		catalogsPath string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "test-catalog-sources",
			args:    args{catalogsPath: "testdata/test-catalog-sources.yaml"},
			want:    []string{"catalog1", "catalog3", "catalog4"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadCatalogSources(tt.args.catalogsPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadCatalogSources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotKeys := make([]string, 0, len(got.All()))
			for k := range got.All() {
				gotKeys = append(gotKeys, k)
			}
			sort.Strings(gotKeys)
			if !reflect.DeepEqual(gotKeys, tt.want) {
				t.Errorf("LoadCatalogSources() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadCatalogSourcesEnabledDisabled(t *testing.T) {
	trueValue := true
	type args struct {
		catalogsPath string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]model.CatalogSource
		wantErr bool
	}{
		{
			name: "test-catalog-sources-enabled-and-disabled",
			args: args{catalogsPath: "testdata/test-catalog-sources.yaml"},
			want: map[string]model.CatalogSource{
				"catalog1": {
					Id:      "catalog1",
					Name:    "Catalog 1",
					Enabled: &trueValue,
				},
				"catalog3": {
					Id:      "catalog3",
					Name:    "Catalog 3",
					Enabled: &trueValue,
				},
				"catalog4": {
					Id:      "catalog4",
					Name:    "Catalog 4",
					Enabled: &trueValue,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadCatalogSources(tt.args.catalogsPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadCatalogSources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			gotMetadata := make(map[string]model.CatalogSource)
			for id, source := range got.All() {
				gotMetadata[id] = source.Metadata
			}

			if !reflect.DeepEqual(gotMetadata, tt.want) {
				t.Errorf("LoadCatalogSources() got metadata = %#v, want %#v", gotMetadata, tt.want)
			}
		})
	}
}
