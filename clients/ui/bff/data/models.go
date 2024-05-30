package data

// Models struct is a single convenient container to hold and represent all our data.
type Models struct {
	HealthCheck HealthCheckModel
}

func NewModels() Models {
	return Models{
		HealthCheck: HealthCheckModel{},
	}
}
