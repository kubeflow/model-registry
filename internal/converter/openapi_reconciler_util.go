package converter

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
)

func UpdateExistingArtifact(genc OpenAPIReconciler, source OpenapiUpdateWrapper[openapi.Artifact]) (openapi.Artifact, error) {
	art := InitWithExisting(source)

	if source.Update == nil {
		return art, nil
	}

	if source.Update.ModelArtifact != nil {
		ma, err := genc.UpdateExistingModelArtifact(OpenapiUpdateWrapper[openapi.ModelArtifact]{Existing: art.ModelArtifact, Update: source.Update.ModelArtifact})
		if err != nil {
			return art, err
		}

		art.ModelArtifact = &ma
	}

	if source.Update.DocArtifact != nil {
		da, err := genc.UpdateExistingDocArtifact(OpenapiUpdateWrapper[openapi.DocArtifact]{Existing: art.DocArtifact, Update: source.Update.DocArtifact})
		if err != nil {
			return art, err
		}

		art.DocArtifact = &da
	}

	return art, nil
}
