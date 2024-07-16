/*
 * Model Registry REST API
 *
 * REST API for Model Registry to create and manage ML model metadata
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 *
 */

// File generated by scripts/gen_type_assert.sh - DO NOT EDIT

package openapi

import (
	model "github.com/kubeflow/model-registry/pkg/openapi"
)

// AssertArtifactRequired checks if the required fields are not zero-ed
func AssertArtifactRequired(obj model.Artifact) error {
	return nil
}

// AssertArtifactConstraints checks if the values respects the defined constraints
func AssertArtifactConstraints(obj model.Artifact) error {
	return nil
}

// AssertArtifactListRequired checks if the required fields are not zero-ed
func AssertArtifactListRequired(obj model.ArtifactList) error {
	elements := map[string]interface{}{
		"nextPageToken": obj.NextPageToken,
		"pageSize":      obj.PageSize,
		"size":          obj.Size,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Items {
		if err := AssertArtifactRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertArtifactListConstraints checks if the values respects the defined constraints
func AssertArtifactListConstraints(obj model.ArtifactList) error {
	return nil
}

// AssertArtifactStateRequired checks if the required fields are not zero-ed
func AssertArtifactStateRequired(obj model.ArtifactState) error {
	return nil
}

// AssertArtifactStateConstraints checks if the values respects the defined constraints
func AssertArtifactStateConstraints(obj model.ArtifactState) error {
	return nil
}

// AssertBaseArtifactCreateRequired checks if the required fields are not zero-ed
func AssertBaseArtifactCreateRequired(obj model.BaseArtifactCreate) error {
	return nil
}

// AssertBaseArtifactCreateConstraints checks if the values respects the defined constraints
func AssertBaseArtifactCreateConstraints(obj model.BaseArtifactCreate) error {
	return nil
}

// AssertBaseArtifactRequired checks if the required fields are not zero-ed
func AssertBaseArtifactRequired(obj model.BaseArtifact) error {
	return nil
}

// AssertBaseArtifactConstraints checks if the values respects the defined constraints
func AssertBaseArtifactConstraints(obj model.BaseArtifact) error {
	return nil
}

// AssertBaseArtifactUpdateRequired checks if the required fields are not zero-ed
func AssertBaseArtifactUpdateRequired(obj model.BaseArtifactUpdate) error {
	return nil
}

// AssertBaseArtifactUpdateConstraints checks if the values respects the defined constraints
func AssertBaseArtifactUpdateConstraints(obj model.BaseArtifactUpdate) error {
	return nil
}

// AssertBaseExecutionCreateRequired checks if the required fields are not zero-ed
func AssertBaseExecutionCreateRequired(obj model.BaseExecutionCreate) error {
	return nil
}

// AssertBaseExecutionCreateConstraints checks if the values respects the defined constraints
func AssertBaseExecutionCreateConstraints(obj model.BaseExecutionCreate) error {
	return nil
}

// AssertBaseExecutionRequired checks if the required fields are not zero-ed
func AssertBaseExecutionRequired(obj model.BaseExecution) error {
	return nil
}

// AssertBaseExecutionConstraints checks if the values respects the defined constraints
func AssertBaseExecutionConstraints(obj model.BaseExecution) error {
	return nil
}

// AssertBaseExecutionUpdateRequired checks if the required fields are not zero-ed
func AssertBaseExecutionUpdateRequired(obj model.BaseExecutionUpdate) error {
	return nil
}

// AssertBaseExecutionUpdateConstraints checks if the values respects the defined constraints
func AssertBaseExecutionUpdateConstraints(obj model.BaseExecutionUpdate) error {
	return nil
}

// AssertBaseResourceCreateRequired checks if the required fields are not zero-ed
func AssertBaseResourceCreateRequired(obj model.BaseResourceCreate) error {
	return nil
}

// AssertBaseResourceCreateConstraints checks if the values respects the defined constraints
func AssertBaseResourceCreateConstraints(obj model.BaseResourceCreate) error {
	return nil
}

// AssertBaseResourceRequired checks if the required fields are not zero-ed
func AssertBaseResourceRequired(obj model.BaseResource) error {
	return nil
}

// AssertBaseResourceConstraints checks if the values respects the defined constraints
func AssertBaseResourceConstraints(obj model.BaseResource) error {
	return nil
}

// AssertBaseResourceListRequired checks if the required fields are not zero-ed
func AssertBaseResourceListRequired(obj model.BaseResourceList) error {
	elements := map[string]interface{}{
		"nextPageToken": obj.NextPageToken,
		"pageSize":      obj.PageSize,
		"size":          obj.Size,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertBaseResourceListConstraints checks if the values respects the defined constraints
func AssertBaseResourceListConstraints(obj model.BaseResourceList) error {
	return nil
}

// AssertBaseResourceUpdateRequired checks if the required fields are not zero-ed
func AssertBaseResourceUpdateRequired(obj model.BaseResourceUpdate) error {
	return nil
}

// AssertBaseResourceUpdateConstraints checks if the values respects the defined constraints
func AssertBaseResourceUpdateConstraints(obj model.BaseResourceUpdate) error {
	return nil
}

// AssertDocArtifactRequired checks if the required fields are not zero-ed
func AssertDocArtifactRequired(obj model.DocArtifact) error {
	elements := map[string]interface{}{
		"artifactType": obj.ArtifactType,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertDocArtifactConstraints checks if the values respects the defined constraints
func AssertDocArtifactConstraints(obj model.DocArtifact) error {
	return nil
}

// AssertErrorRequired checks if the required fields are not zero-ed
func AssertErrorRequired(obj model.Error) error {
	elements := map[string]interface{}{
		"code":    obj.Code,
		"message": obj.Message,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertErrorConstraints checks if the values respects the defined constraints
func AssertErrorConstraints(obj model.Error) error {
	return nil
}

// AssertExecutionStateRequired checks if the required fields are not zero-ed
func AssertExecutionStateRequired(obj model.ExecutionState) error {
	return nil
}

// AssertExecutionStateConstraints checks if the values respects the defined constraints
func AssertExecutionStateConstraints(obj model.ExecutionState) error {
	return nil
}

// AssertInferenceServiceCreateRequired checks if the required fields are not zero-ed
func AssertInferenceServiceCreateRequired(obj model.InferenceServiceCreate) error {
	elements := map[string]interface{}{
		"registeredModelId":    obj.RegisteredModelId,
		"servingEnvironmentId": obj.ServingEnvironmentId,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertInferenceServiceCreateConstraints checks if the values respects the defined constraints
func AssertInferenceServiceCreateConstraints(obj model.InferenceServiceCreate) error {
	return nil
}

// AssertInferenceServiceRequired checks if the required fields are not zero-ed
func AssertInferenceServiceRequired(obj model.InferenceService) error {
	elements := map[string]interface{}{
		"registeredModelId":    obj.RegisteredModelId,
		"servingEnvironmentId": obj.ServingEnvironmentId,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertInferenceServiceConstraints checks if the values respects the defined constraints
func AssertInferenceServiceConstraints(obj model.InferenceService) error {
	return nil
}

// AssertInferenceServiceListRequired checks if the required fields are not zero-ed
func AssertInferenceServiceListRequired(obj model.InferenceServiceList) error {
	elements := map[string]interface{}{
		"nextPageToken": obj.NextPageToken,
		"pageSize":      obj.PageSize,
		"size":          obj.Size,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Items {
		if err := AssertInferenceServiceRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertInferenceServiceListConstraints checks if the values respects the defined constraints
func AssertInferenceServiceListConstraints(obj model.InferenceServiceList) error {
	return nil
}

// AssertInferenceServiceStateRequired checks if the required fields are not zero-ed
func AssertInferenceServiceStateRequired(obj model.InferenceServiceState) error {
	return nil
}

// AssertInferenceServiceStateConstraints checks if the values respects the defined constraints
func AssertInferenceServiceStateConstraints(obj model.InferenceServiceState) error {
	return nil
}

// AssertInferenceServiceUpdateRequired checks if the required fields are not zero-ed
func AssertInferenceServiceUpdateRequired(obj model.InferenceServiceUpdate) error {
	return nil
}

// AssertInferenceServiceUpdateConstraints checks if the values respects the defined constraints
func AssertInferenceServiceUpdateConstraints(obj model.InferenceServiceUpdate) error {
	return nil
}

// AssertMetadataBoolValueRequired checks if the required fields are not zero-ed
func AssertMetadataBoolValueRequired(obj model.MetadataBoolValue) error {
	elements := map[string]interface{}{
		"bool_value":   obj.BoolValue,
		"metadataType": obj.MetadataType,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertMetadataBoolValueConstraints checks if the values respects the defined constraints
func AssertMetadataBoolValueConstraints(obj model.MetadataBoolValue) error {
	return nil
}

// AssertMetadataDoubleValueRequired checks if the required fields are not zero-ed
func AssertMetadataDoubleValueRequired(obj model.MetadataDoubleValue) error {
	elements := map[string]interface{}{
		"double_value": obj.DoubleValue,
		"metadataType": obj.MetadataType,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertMetadataDoubleValueConstraints checks if the values respects the defined constraints
func AssertMetadataDoubleValueConstraints(obj model.MetadataDoubleValue) error {
	return nil
}

// AssertMetadataIntValueRequired checks if the required fields are not zero-ed
func AssertMetadataIntValueRequired(obj model.MetadataIntValue) error {
	elements := map[string]interface{}{
		"int_value":    obj.IntValue,
		"metadataType": obj.MetadataType,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertMetadataIntValueConstraints checks if the values respects the defined constraints
func AssertMetadataIntValueConstraints(obj model.MetadataIntValue) error {
	return nil
}

// AssertMetadataProtoValueRequired checks if the required fields are not zero-ed
func AssertMetadataProtoValueRequired(obj model.MetadataProtoValue) error {
	elements := map[string]interface{}{
		"type":         obj.Type,
		"proto_value":  obj.ProtoValue,
		"metadataType": obj.MetadataType,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertMetadataProtoValueConstraints checks if the values respects the defined constraints
func AssertMetadataProtoValueConstraints(obj model.MetadataProtoValue) error {
	return nil
}

// AssertMetadataStringValueRequired checks if the required fields are not zero-ed
func AssertMetadataStringValueRequired(obj model.MetadataStringValue) error {
	elements := map[string]interface{}{
		"string_value": obj.StringValue,
		"metadataType": obj.MetadataType,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertMetadataStringValueConstraints checks if the values respects the defined constraints
func AssertMetadataStringValueConstraints(obj model.MetadataStringValue) error {
	return nil
}

// AssertMetadataStructValueRequired checks if the required fields are not zero-ed
func AssertMetadataStructValueRequired(obj model.MetadataStructValue) error {
	elements := map[string]interface{}{
		"struct_value": obj.StructValue,
		"metadataType": obj.MetadataType,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertMetadataStructValueConstraints checks if the values respects the defined constraints
func AssertMetadataStructValueConstraints(obj model.MetadataStructValue) error {
	return nil
}

// AssertMetadataValueRequired checks if the required fields are not zero-ed
func AssertMetadataValueRequired(obj model.MetadataValue) error {
	// FIXME(manual): Wrong autogenerated logic, just one elem should be non-zero.
	// elements := map[string]interface{}{
	// 	"int_value":    obj.IntValue,
	// 	"metadataType": obj.MetadataType,
	// 	"double_value": obj.DoubleValue,
	// 	"string_value": obj.StringValue,
	// 	"struct_value": obj.StructValue,
	// 	"type":         obj.Type,
	// 	"proto_value":  obj.ProtoValue,
	// 	"bool_value":   obj.BoolValue,
	// }
	// for name, el := range elements {
	// 	if isZero := IsZeroValue(el); isZero {
	// 		return &RequiredError{Field: name}
	// 	}
	// }

	return nil
}

// AssertMetadataValueConstraints checks if the values respects the defined constraints
func AssertMetadataValueConstraints(obj model.MetadataValue) error {
	return nil
}

// AssertModelArtifactCreateRequired checks if the required fields are not zero-ed
func AssertModelArtifactCreateRequired(obj model.ModelArtifactCreate) error {
	return nil
}

// AssertModelArtifactCreateConstraints checks if the values respects the defined constraints
func AssertModelArtifactCreateConstraints(obj model.ModelArtifactCreate) error {
	return nil
}

// AssertModelArtifactRequired checks if the required fields are not zero-ed
func AssertModelArtifactRequired(obj model.ModelArtifact) error {
	elements := map[string]interface{}{
		"artifactType": obj.ArtifactType,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertModelArtifactConstraints checks if the values respects the defined constraints
func AssertModelArtifactConstraints(obj model.ModelArtifact) error {
	return nil
}

// AssertModelArtifactListRequired checks if the required fields are not zero-ed
func AssertModelArtifactListRequired(obj model.ModelArtifactList) error {
	elements := map[string]interface{}{
		"nextPageToken": obj.NextPageToken,
		"pageSize":      obj.PageSize,
		"size":          obj.Size,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Items {
		if err := AssertModelArtifactRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertModelArtifactListConstraints checks if the values respects the defined constraints
func AssertModelArtifactListConstraints(obj model.ModelArtifactList) error {
	return nil
}

// AssertModelArtifactUpdateRequired checks if the required fields are not zero-ed
func AssertModelArtifactUpdateRequired(obj model.ModelArtifactUpdate) error {
	return nil
}

// AssertModelArtifactUpdateConstraints checks if the values respects the defined constraints
func AssertModelArtifactUpdateConstraints(obj model.ModelArtifactUpdate) error {
	return nil
}

// AssertModelVersionCreateRequired checks if the required fields are not zero-ed
func AssertModelVersionCreateRequired(obj model.ModelVersionCreate) error {
	elements := map[string]interface{}{
		"registeredModelId": obj.RegisteredModelId,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertModelVersionCreateConstraints checks if the values respects the defined constraints
func AssertModelVersionCreateConstraints(obj model.ModelVersionCreate) error {
	return nil
}

// AssertModelVersionRequired checks if the required fields are not zero-ed
func AssertModelVersionRequired(obj model.ModelVersion) error {
	elements := map[string]interface{}{
		"registeredModelId": obj.RegisteredModelId,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertModelVersionConstraints checks if the values respects the defined constraints
func AssertModelVersionConstraints(obj model.ModelVersion) error {
	return nil
}

// AssertModelVersionListRequired checks if the required fields are not zero-ed
func AssertModelVersionListRequired(obj model.ModelVersionList) error {
	elements := map[string]interface{}{
		"nextPageToken": obj.NextPageToken,
		"pageSize":      obj.PageSize,
		"size":          obj.Size,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Items {
		if err := AssertModelVersionRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertModelVersionListConstraints checks if the values respects the defined constraints
func AssertModelVersionListConstraints(obj model.ModelVersionList) error {
	return nil
}

// AssertModelVersionStateRequired checks if the required fields are not zero-ed
func AssertModelVersionStateRequired(obj model.ModelVersionState) error {
	return nil
}

// AssertModelVersionStateConstraints checks if the values respects the defined constraints
func AssertModelVersionStateConstraints(obj model.ModelVersionState) error {
	return nil
}

// AssertModelVersionUpdateRequired checks if the required fields are not zero-ed
func AssertModelVersionUpdateRequired(obj model.ModelVersionUpdate) error {
	return nil
}

// AssertModelVersionUpdateConstraints checks if the values respects the defined constraints
func AssertModelVersionUpdateConstraints(obj model.ModelVersionUpdate) error {
	return nil
}

// AssertOrderByFieldRequired checks if the required fields are not zero-ed
func AssertOrderByFieldRequired(obj model.OrderByField) error {
	return nil
}

// AssertOrderByFieldConstraints checks if the values respects the defined constraints
func AssertOrderByFieldConstraints(obj model.OrderByField) error {
	return nil
}

// AssertRegisteredModelCreateRequired checks if the required fields are not zero-ed
func AssertRegisteredModelCreateRequired(obj model.RegisteredModelCreate) error {
	return nil
}

// AssertRegisteredModelCreateConstraints checks if the values respects the defined constraints
func AssertRegisteredModelCreateConstraints(obj model.RegisteredModelCreate) error {
	return nil
}

// AssertRegisteredModelRequired checks if the required fields are not zero-ed
func AssertRegisteredModelRequired(obj model.RegisteredModel) error {
	return nil
}

// AssertRegisteredModelConstraints checks if the values respects the defined constraints
func AssertRegisteredModelConstraints(obj model.RegisteredModel) error {
	return nil
}

// AssertRegisteredModelListRequired checks if the required fields are not zero-ed
func AssertRegisteredModelListRequired(obj model.RegisteredModelList) error {
	elements := map[string]interface{}{
		"nextPageToken": obj.NextPageToken,
		"pageSize":      obj.PageSize,
		"size":          obj.Size,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Items {
		if err := AssertRegisteredModelRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRegisteredModelListConstraints checks if the values respects the defined constraints
func AssertRegisteredModelListConstraints(obj model.RegisteredModelList) error {
	return nil
}

// AssertRegisteredModelStateRequired checks if the required fields are not zero-ed
func AssertRegisteredModelStateRequired(obj model.RegisteredModelState) error {
	return nil
}

// AssertRegisteredModelStateConstraints checks if the values respects the defined constraints
func AssertRegisteredModelStateConstraints(obj model.RegisteredModelState) error {
	return nil
}

// AssertRegisteredModelUpdateRequired checks if the required fields are not zero-ed
func AssertRegisteredModelUpdateRequired(obj model.RegisteredModelUpdate) error {
	return nil
}

// AssertRegisteredModelUpdateConstraints checks if the values respects the defined constraints
func AssertRegisteredModelUpdateConstraints(obj model.RegisteredModelUpdate) error {
	return nil
}

// AssertServeModelCreateRequired checks if the required fields are not zero-ed
func AssertServeModelCreateRequired(obj model.ServeModelCreate) error {
	elements := map[string]interface{}{
		"modelVersionId": obj.ModelVersionId,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertServeModelCreateConstraints checks if the values respects the defined constraints
func AssertServeModelCreateConstraints(obj model.ServeModelCreate) error {
	return nil
}

// AssertServeModelRequired checks if the required fields are not zero-ed
func AssertServeModelRequired(obj model.ServeModel) error {
	elements := map[string]interface{}{
		"modelVersionId": obj.ModelVersionId,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertServeModelConstraints checks if the values respects the defined constraints
func AssertServeModelConstraints(obj model.ServeModel) error {
	return nil
}

// AssertServeModelListRequired checks if the required fields are not zero-ed
func AssertServeModelListRequired(obj model.ServeModelList) error {
	elements := map[string]interface{}{
		"nextPageToken": obj.NextPageToken,
		"pageSize":      obj.PageSize,
		"size":          obj.Size,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Items {
		if err := AssertServeModelRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertServeModelListConstraints checks if the values respects the defined constraints
func AssertServeModelListConstraints(obj model.ServeModelList) error {
	return nil
}

// AssertServeModelUpdateRequired checks if the required fields are not zero-ed
func AssertServeModelUpdateRequired(obj model.ServeModelUpdate) error {
	return nil
}

// AssertServeModelUpdateConstraints checks if the values respects the defined constraints
func AssertServeModelUpdateConstraints(obj model.ServeModelUpdate) error {
	return nil
}

// AssertServingEnvironmentCreateRequired checks if the required fields are not zero-ed
func AssertServingEnvironmentCreateRequired(obj model.ServingEnvironmentCreate) error {
	return nil
}

// AssertServingEnvironmentCreateConstraints checks if the values respects the defined constraints
func AssertServingEnvironmentCreateConstraints(obj model.ServingEnvironmentCreate) error {
	return nil
}

// AssertServingEnvironmentRequired checks if the required fields are not zero-ed
func AssertServingEnvironmentRequired(obj model.ServingEnvironment) error {
	return nil
}

// AssertServingEnvironmentConstraints checks if the values respects the defined constraints
func AssertServingEnvironmentConstraints(obj model.ServingEnvironment) error {
	return nil
}

// AssertServingEnvironmentListRequired checks if the required fields are not zero-ed
func AssertServingEnvironmentListRequired(obj model.ServingEnvironmentList) error {
	elements := map[string]interface{}{
		"nextPageToken": obj.NextPageToken,
		"pageSize":      obj.PageSize,
		"size":          obj.Size,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Items {
		if err := AssertServingEnvironmentRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertServingEnvironmentListConstraints checks if the values respects the defined constraints
func AssertServingEnvironmentListConstraints(obj model.ServingEnvironmentList) error {
	return nil
}

// AssertServingEnvironmentUpdateRequired checks if the required fields are not zero-ed
func AssertServingEnvironmentUpdateRequired(obj model.ServingEnvironmentUpdate) error {
	return nil
}

// AssertServingEnvironmentUpdateConstraints checks if the values respects the defined constraints
func AssertServingEnvironmentUpdateConstraints(obj model.ServingEnvironmentUpdate) error {
	return nil
}

// AssertSortOrderRequired checks if the required fields are not zero-ed
func AssertSortOrderRequired(obj model.SortOrder) error {
	return nil
}

// AssertSortOrderConstraints checks if the values respects the defined constraints
func AssertSortOrderConstraints(obj model.SortOrder) error {
	return nil
}
