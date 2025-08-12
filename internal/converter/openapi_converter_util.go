package converter

import "github.com/kubeflow/model-registry/pkg/openapi"

type OpenAPIModel interface {
	openapi.Artifact |
		openapi.RegisteredModel |
		openapi.ModelVersion |
		openapi.ModelArtifact |
		openapi.DocArtifact |
		openapi.DataSet |
		openapi.Metric |
		openapi.Parameter |
		openapi.ServingEnvironment |
		openapi.InferenceService |
		openapi.ServeModel |
		openapi.Experiment |
		openapi.ExperimentRun
}

type OpenapiUpdateWrapper[
	M OpenAPIModel,
] struct {
	Existing *M
	Update   *M
}

func NewOpenapiUpdateWrapper[
	M OpenAPIModel,
](existing *M, update *M) OpenapiUpdateWrapper[M] {
	return OpenapiUpdateWrapper[M]{
		Existing: existing,
		Update:   update,
	}
}

func InitWithExisting[M OpenAPIModel](source OpenapiUpdateWrapper[M]) M {
	return *source.Existing
}

func InitWithUpdate[M OpenAPIModel](source OpenapiUpdateWrapper[M]) M {
	if source.Update != nil {
		return *source.Update
	}
	var m M
	return m
}
