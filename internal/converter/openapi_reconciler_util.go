package converter

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
)

func UpdateExistingArtifact(genc OpenAPIReconciler, source OpenapiUpdateWrapper[openapi.Artifact]) (openapi.Artifact, error) {
	art := InitWithExisting(source)
	if source.Update == nil {
		return art, nil
	}
	ma, err := genc.UpdateExistingModelArtifact(OpenapiUpdateWrapper[openapi.ModelArtifact]{Existing: art.ModelArtifact, Update: source.Update.ModelArtifact})
	if err != nil {
		return art, err
	}
	da, err := genc.UpdateExistingDocArtifact(OpenapiUpdateWrapper[openapi.DocArtifact]{Existing: art.DocArtifact, Update: source.Update.DocArtifact})
	if err != nil {
		return art, err
	}
	art.DocArtifact = &da
	art.ModelArtifact = &ma
	return art, nil
}
