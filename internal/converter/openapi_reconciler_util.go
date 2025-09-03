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
