package data

type SystemInfo struct {
	Version string `json:"version"`
}

type HealthCheckModel struct {
	Status     string     `json:"status"`
	SystemInfo SystemInfo `json:"system_info"`
}

func (m HealthCheckModel) HealthCheck(version string) (HealthCheckModel, error) {

	var res = HealthCheckModel{
		Status: "available",
		SystemInfo: SystemInfo{
			Version: version,
		},
	}

	return res, nil
}
