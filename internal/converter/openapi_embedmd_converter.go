package converter

import (
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// goverter:converter
// goverter:output:file ./generated/openapi_embedmd_converter.gen.go
// goverter:wrapErrors
// goverter:matchIgnoreCase
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend Int64ToString
// goverter:extend StringToInt32
// goverter:extend MapOpenAPICustomPropertiesEmbedMD
type OpenAPIToEmbedMDConverter interface {
	// goverter:autoMap Model
	// goverter:map Model Properties | MapRegisteredModelPropertiesEmbedMD
	// goverter:map Model Attributes | MapRegisteredModelAttributesEmbedMD
	// goverter:map . TypeID | MapRegisteredModelTypeIDEmbedMD
	ConvertRegisteredModel(source *OpenAPIModelWrapper[openapi.RegisteredModel]) (*models.RegisteredModelImpl, error)

	// goverter:autoMap Model
	// goverter:map Model Properties | MapModelVersionPropertiesEmbedMD
	// goverter:map Model Attributes | MapModelVersionAttributesEmbedMD
	// goverter:map . TypeID | MapModelVersionTypeIDEmbedMD
	ConvertModelVersion(source *OpenAPIModelWrapper[openapi.ModelVersion]) (*models.ModelVersionImpl, error)

	// goverter:autoMap Model
	// goverter:map Model Properties | MapModelArtifactPropertiesEmbedMD
	// goverter:map Model Attributes | MapModelArtifactAttributesEmbedMD
	// goverter:map . TypeID | MapModelArtifactTypeIDEmbedMD
	ConvertModelArtifact(source *OpenAPIModelWrapper[openapi.ModelArtifact]) (*models.ModelArtifactImpl, error)

	// goverter:autoMap Model
	// goverter:map Model Properties | MapDocArtifactPropertiesEmbedMD
	// goverter:map Model Attributes | MapDocArtifactAttributesEmbedMD
	// goverter:map . TypeID | MapDocArtifactTypeIDEmbedMD
	ConvertDocArtifact(source *OpenAPIModelWrapper[openapi.DocArtifact]) (*models.DocArtifactImpl, error)

	// goverter:autoMap Model
	// goverter:map Model Properties | MapServingEnvironmentPropertiesEmbedMD
	// goverter:map Model Attributes | MapServingEnvironmentAttributesEmbedMD
	// goverter:map . TypeID | MapServingEnvironmentTypeIDEmbedMD
	ConvertServingEnvironment(source *OpenAPIModelWrapper[openapi.ServingEnvironment]) (*models.ServingEnvironmentImpl, error)

	// goverter:autoMap Model
	// goverter:map Model Properties | MapInferenceServicePropertiesEmbedMD
	// goverter:map Model Attributes | MapInferenceServiceAttributesEmbedMD
	// goverter:map . TypeID | MapInferenceServiceTypeIDEmbedMD
	ConvertInferenceService(source *OpenAPIModelWrapper[openapi.InferenceService]) (*models.InferenceServiceImpl, error)

	// goverter:autoMap Model
	// goverter:map Model Properties | MapServeModelPropertiesEmbedMD
	// goverter:map Model Attributes | MapServeModelAttributesEmbedMD
	// goverter:map . TypeID | MapServeModelTypeIDEmbedMD
	ConvertServeModel(source *OpenAPIModelWrapper[openapi.ServeModel]) (*models.ServeModelImpl, error)
}
