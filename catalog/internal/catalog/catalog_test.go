package catalog

import (
	"reflect"
	"sort"
	"testing"
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
			want:    []string{"catalog1", "catalog2"},
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
			gotKeys := make([]string, 0, len(got))
			for k := range got {
				gotKeys = append(gotKeys, k)
			}
			sort.Strings(gotKeys)
			if !reflect.DeepEqual(gotKeys, tt.want) {
				t.Errorf("LoadCatalogSources() got = %v, want %v", got, tt.want)
			}
		})
	}
}
