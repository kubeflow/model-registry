package converter

import "github.com/opendatahub-io/model-registry/pkg/openapi"

type OpenAPIModel interface {
	openapi.RegisteredModel |
		openapi.ModelVersion |
		openapi.ModelArtifact |
		openapi.DocArtifact |
		openapi.ServingEnvironment |
		openapi.InferenceService |
		openapi.ServeModel
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

func InitRegisteredModelWithUpdate(source OpenapiUpdateWrapper[openapi.RegisteredModel]) openapi.RegisteredModel {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.RegisteredModel{}
}

func InitModelVersionWithUpdate(source OpenapiUpdateWrapper[openapi.ModelVersion]) openapi.ModelVersion {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.ModelVersion{}
}

func InitDocArtifactWithUpdate(source OpenapiUpdateWrapper[openapi.DocArtifact]) openapi.DocArtifact {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.DocArtifact{}
}

func InitModelArtifactWithUpdate(source OpenapiUpdateWrapper[openapi.ModelArtifact]) openapi.ModelArtifact {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.ModelArtifact{}
}

func InitServingEnvironmentWithUpdate(source OpenapiUpdateWrapper[openapi.ServingEnvironment]) openapi.ServingEnvironment {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.ServingEnvironment{}
}

func InitInferenceServiceWithUpdate(source OpenapiUpdateWrapper[openapi.InferenceService]) openapi.InferenceService {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.InferenceService{}
}

func InitServeModelWithUpdate(source OpenapiUpdateWrapper[openapi.ServeModel]) openapi.ServeModel {
	if source.Update != nil {
		return *source.Update
	}
	return openapi.ServeModel{}
}
