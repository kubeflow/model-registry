package converter

import (
	"fmt"

	"github.com/kubeflow/model-registry/pkg/openapi"
)

// getArtifactTypeName returns the type name of an artifact.
func getArtifactTypeName(art *openapi.Artifact) string {
	if art == nil {
		return "unknown"
	}
	if art.ModelArtifact != nil {
		return string(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT)
	}
	if art.DocArtifact != nil {
		return string(openapi.ARTIFACTTYPEQUERYPARAM_DOC_ARTIFACT)
	}
	if art.DataSet != nil {
		return string(openapi.ARTIFACTTYPEQUERYPARAM_DATASET_ARTIFACT)
	}
	if art.Metric != nil {
		return string(openapi.ARTIFACTTYPEQUERYPARAM_METRIC)
	}
	if art.Parameter != nil {
		return string(openapi.ARTIFACTTYPEQUERYPARAM_PARAMETER)
	}
	return "unknown"
}

func UpdateExistingArtifact(genc OpenAPIReconciler, source OpenapiUpdateWrapper[openapi.Artifact]) (openapi.Artifact, error) {
	art := InitWithExisting(source)

	if source.Update == nil {
		return art, nil
	}

	// Validate that the artifact type in the update matches the existing artifact type
	// Changing artifact type is not allowed
	existingType := getArtifactTypeName(source.Existing)
	updateType := getArtifactTypeName(source.Update)

	if existingType != updateType {
		return art, fmt.Errorf("cannot change artifact type from '%s' to '%s': artifact type is immutable", existingType, updateType)
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

	if source.Update.DataSet != nil {
		ds, err := genc.UpdateExistingDataSet(OpenapiUpdateWrapper[openapi.DataSet]{Existing: art.DataSet, Update: source.Update.DataSet})
		if err != nil {
			return art, err
		}
		art.DataSet = &ds
	}

	if source.Update.Metric != nil {
		mt, err := genc.UpdateExistingMetric(OpenapiUpdateWrapper[openapi.Metric]{Existing: art.Metric, Update: source.Update.Metric})
		if err != nil {
			return art, err
		}
		art.Metric = &mt
	}

	if source.Update.Parameter != nil {
		pa, err := genc.UpdateExistingParameter(OpenapiUpdateWrapper[openapi.Parameter]{Existing: art.Parameter, Update: source.Update.Parameter})
		if err != nil {
			return art, err
		}
		art.Parameter = &pa
	}

	return art, nil
}
