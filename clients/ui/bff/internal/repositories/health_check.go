package repositories

import "github.com/kubeflow/model-registry/ui/bff/internal/models"

type HealthCheckRepository struct{}

func NewHealthCheckRepository() *HealthCheckRepository {
	return &HealthCheckRepository{}
}

func (r *HealthCheckRepository) HealthCheck(version string) (models.HealthCheckModel, error) {
	var res = models.HealthCheckModel{
		Status: "available",
		SystemInfo: models.SystemInfo{
			Version: version,
		},
	}

	return res, nil
}
