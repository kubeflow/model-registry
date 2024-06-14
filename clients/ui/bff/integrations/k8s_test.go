package integrations

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log/slog"
	"os"
	"testing"
)

func TestBuildModelRegistryServiceCache(t *testing.T) {
	services := v1.ServiceList{
		Items: []v1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "service-dora",
				},
				Spec: v1.ServiceSpec{
					ClusterIP: "10.0.0.1",
					Ports: []v1.ServicePort{
						{
							Name: "http-api",
							Port: 80,
						},
					},
					Selector: map[string]string{
						"component": "model-registry-server",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "service-bella",
				},
				Spec: v1.ServiceSpec{
					ClusterIP: "10.0.0.2",
					Ports: []v1.ServicePort{
						{
							Name: "http-api",
							Port: 8080,
						},
					},
					Selector: map[string]string{
						"component": "model-registry-server",
					},
				},
			},
		},
	}
	expectedServiceCache := map[string]ServiceDetails{
		"service-dora": {
			Name:      "service-dora",
			ClusterIP: "10.0.0.1",
			HTTPPort:  80,
		},
		"service-bella": {
			Name:      "service-bella",
			ClusterIP: "10.0.0.2",
			HTTPPort:  8080,
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	serviceCache, err := buildModelRegistryServiceCache(logger, services)
	assert.NoError(t, err, "unexpected error while building service cache")
	assert.Equal(t, expectedServiceCache, serviceCache, "serviceCache does not match expected value")
}
