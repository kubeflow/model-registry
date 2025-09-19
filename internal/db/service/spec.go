package service

import (
	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/defaults"
)

func DatastoreSpec() *datastore.Spec {
	return datastore.NewSpec().
		AddArtifact(defaults.ModelArtifactTypeName, datastore.NewSpecType(NewModelArtifactRepository).
			AddString("description").
			AddString("model_format_name").
			AddString("model_format_version").
			AddString("service_account_name").
			AddString("storage_key").
			AddString("storage_path"),
		).
		AddArtifact(defaults.DocArtifactTypeName, datastore.NewSpecType(NewDocArtifactRepository).
			AddString("description"),
		).
		AddArtifact(defaults.DataSetTypeName, datastore.NewSpecType(NewDataSetRepository).
			AddString("description").
			AddString("digest").
			AddString("source_type").
			AddString("source").
			AddString("schema").
			AddString("profile"),
		).
		AddArtifact(defaults.MetricTypeName, datastore.NewSpecType(NewMetricRepository).
			AddString("description").
			AddProto("value").
			AddString("timestamp").
			AddInt("step"),
		).
		AddArtifact(defaults.ParameterTypeName, datastore.NewSpecType(NewParameterRepository).
			AddString("description").
			AddString("value").
			AddString("parameter_type"),
		).
		AddArtifact(defaults.MetricHistoryTypeName, datastore.NewSpecType(NewMetricHistoryRepository).
			AddString("description").
			AddProto("value").
			AddString("timestamp").
			AddInt("step"),
		).
		AddContext(defaults.RegisteredModelTypeName, datastore.NewSpecType(NewRegisteredModelRepository).
			AddString("description").
			AddString("owner").
			AddString("state").
			AddStruct("language").
			AddString("library_name").
			AddString("license_link").
			AddString("license").
			AddString("logo").
			AddString("maturity").
			AddString("provider").
			AddString("readme").
			AddStruct("tasks"),
		).
		AddContext(defaults.ModelVersionTypeName, datastore.NewSpecType(NewModelVersionRepository).
			AddString("author").
			AddString("description").
			AddString("model_name").
			AddString("state").
			AddString("version"),
		).
		AddContext(defaults.ServingEnvironmentTypeName, datastore.NewSpecType(NewServingEnvironmentRepository).
			AddString("description"),
		).
		AddContext(defaults.InferenceServiceTypeName, datastore.NewSpecType(NewInferenceServiceRepository).
			AddString("description").
			AddString("desired_state").
			AddInt("model_version_id").
			AddInt("registered_model_id").
			AddString("runtime").
			AddInt("serving_environment_id"),
		).
		AddContext(defaults.ExperimentTypeName, datastore.NewSpecType(NewExperimentRepository).
			AddString("description").
			AddString("owner").
			AddString("state"),
		).
		AddContext(defaults.ExperimentRunTypeName, datastore.NewSpecType(NewExperimentRunRepository).
			AddString("description").
			AddString("owner").
			AddString("state").
			AddString("status").
			AddInt("start_time_since_epoch").
			AddInt("end_time_since_epoch").
			AddInt("experiment_id"),
		).
		AddExecution(defaults.ServeModelTypeName, datastore.NewSpecType(NewServeModelRepository).
			AddString("description").
			AddInt("model_version_id"),
		).
		AddOther(NewArtifactRepository)
}
