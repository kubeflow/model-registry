package mapper

import (
	"fmt"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/constants"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
)

const (
	invalidTypeId            = int64(9999)
	registeredModelTypeId    = int64(1)
	modelVersionTypeId       = int64(2)
	docArtifactTypeId        = int64(3)
	modelArtifactTypeId      = int64(4)
	servingEnvironmentTypeId = int64(5)
	inferenceServiceTypeId   = int64(6)
	serveModelTypeId         = int64(7)
)

var typesMap = map[string]int64{
	constants.RegisteredModelTypeName:    registeredModelTypeId,
	constants.ModelVersionTypeName:       modelVersionTypeId,
	constants.DocArtifactTypeName:        docArtifactTypeId,
	constants.ModelArtifactTypeName:      modelArtifactTypeId,
	constants.ServingEnvironmentTypeName: servingEnvironmentTypeId,
	constants.InferenceServiceTypeName:   inferenceServiceTypeId,
	constants.ServeModelTypeName:         serveModelTypeId,
}

func setup(t *testing.T) (*assert.Assertions, *Mapper) {
	return assert.New(t), NewMapper(
		typesMap,
	)
}

func TestMapFromRegisteredModel(t *testing.T) {
	assertion, m := setup(t)

	ctx, err := m.MapFromRegisteredModel(&openapi.RegisteredModel{Name: of("ModelName")})
	assertion.Nil(err)
	assertion.Equal("ModelName", ctx.GetName())
	assertion.Equal(registeredModelTypeId, ctx.GetTypeId())
}

func TestMapFromModelVersion(t *testing.T) {
	assertion, m := setup(t)

	ctx, err := m.MapFromModelVersion(&openapi.ModelVersion{Name: of("v1")}, "1", of("ModelName"))
	assertion.Nil(err)
	assertion.Equal("1:v1", ctx.GetName())
	assertion.Equal(modelVersionTypeId, ctx.GetTypeId())
}

func TestMapFromDocArtifact(t *testing.T) {
	assertion, m := setup(t)

	ctx, err := m.MapFromArtifact(&openapi.Artifact{
		DocArtifact: &openapi.DocArtifact{Name: of("DocArtifact")},
	}, of("2"))
	assertion.Nil(err)
	assertion.Equal("2:DocArtifact", ctx.GetName())
	assertion.Equal(docArtifactTypeId, ctx.GetTypeId())
}

func TestMapFromModelArtifact(t *testing.T) {
	assertion, m := setup(t)

	ctx, err := m.MapFromArtifact(&openapi.Artifact{
		ModelArtifact: &openapi.ModelArtifact{Name: of("ModelArtifact")},
	}, of("2"))
	assertion.Nil(err)
	assertion.Equal("2:ModelArtifact", ctx.GetName())
	assertion.Equal(modelArtifactTypeId, ctx.GetTypeId())
}

func TestMapFromModelArtifacts(t *testing.T) {
	assertion, m := setup(t)

	ctxList, err := m.MapFromModelArtifacts([]openapi.ModelArtifact{{Name: of("ModelArtifact1")}, {Name: of("ModelArtifact2")}}, of("2"))
	assertion.Nil(err)
	assertion.Equal(2, len(ctxList))
	assertion.Equal("2:ModelArtifact1", ctxList[0].GetName())
	assertion.Equal("2:ModelArtifact2", ctxList[1].GetName())
	assertion.Equal(modelArtifactTypeId, ctxList[0].GetTypeId())
	assertion.Equal(modelArtifactTypeId, ctxList[1].GetTypeId())
}

func TestMapFromModelArtifactsEmpty(t *testing.T) {
	assertion, m := setup(t)

	ctxList, err := m.MapFromModelArtifacts([]openapi.ModelArtifact{}, of("2"))
	assertion.Nil(err)
	assertion.Equal(0, len(ctxList))

	ctxList, err = m.MapFromModelArtifacts(nil, nil)
	assertion.Nil(err)
	assertion.Equal(0, len(ctxList))
}

func TestMapFromServingEnvironment(t *testing.T) {
	assertion, m := setup(t)

	ctx, err := m.MapFromServingEnvironment(&openapi.ServingEnvironment{Name: of("Env")})
	assertion.Nil(err)
	assertion.Equal("Env", ctx.GetName())
	assertion.Equal(servingEnvironmentTypeId, ctx.GetTypeId())
}

func TestMapFromInferenceService(t *testing.T) {
	assertion, m := setup(t)

	ctx, err := m.MapFromInferenceService(&openapi.InferenceService{Name: of("IS"), ServingEnvironmentId: "5", RegisteredModelId: "1"}, "5")
	assertion.Nil(err)
	assertion.Equal("5:IS", ctx.GetName())
	assertion.Equal(inferenceServiceTypeId, ctx.GetTypeId())
}

func TestMapFromInferenceServiceMissingRequiredIds(t *testing.T) {
	assertion, m := setup(t)

	_, err := m.MapFromInferenceService(&openapi.InferenceService{Name: of("IS")}, "5")
	assertion.NotNil(err)
	assertion.Equal("error setting field Properties: missing required RegisteredModelId field", err.Error())
}

func TestMapFromServeModel(t *testing.T) {
	assertion, m := setup(t)

	ctx, err := m.MapFromServeModel(&openapi.ServeModel{Name: of("Serve"), ModelVersionId: "1"}, "10")
	assertion.Nil(err)
	assertion.Equal("10:Serve", ctx.GetName())
	assertion.Equal(serveModelTypeId, ctx.GetTypeId())
}

func TestMapFromServeModelMissingRequiredId(t *testing.T) {
	assertion, m := setup(t)

	_, err := m.MapFromServeModel(&openapi.ServeModel{Name: of("Serve")}, "10")
	assertion.NotNil(err)
	assertion.Equal("error setting field Properties: missing required ModelVersionId field", err.Error())
}

func TestMapToRegisteredModel(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToRegisteredModel(&proto.Context{
		TypeId: of(registeredModelTypeId),
		Type:   of(constants.RegisteredModelTypeName),
	})
	assertion.Nil(err)
}

func TestMapToRegisteredModelInvalid(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToRegisteredModel(&proto.Context{
		TypeId: of(invalidTypeId),
		Type:   of("odh.OtherEntity"),
	})
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("invalid entity: expected %s but received odh.OtherEntity, please check the provided id", constants.RegisteredModelTypeName), err.Error())
}

func TestMapToModelVersion(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToModelVersion(&proto.Context{
		TypeId: of(modelVersionTypeId),
		Type:   of(constants.ModelVersionTypeName),
	})
	assertion.Nil(err)
}

func TestMapToModelVersionInvalid(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToModelVersion(&proto.Context{
		TypeId: of(invalidTypeId),
		Type:   of("odh.OtherEntity"),
	})
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("invalid entity: expected %s but received odh.OtherEntity, please check the provided id", constants.ModelVersionTypeName), err.Error())
}

func TestMapToDocArtifact(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToArtifact(&proto.Artifact{
		TypeId: of(docArtifactTypeId),
		Type:   of(constants.DocArtifactTypeName),
	})
	assertion.Nil(err)
}

func TestMapToModelArtifact(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToArtifact(&proto.Artifact{
		TypeId: of(modelArtifactTypeId),
		Type:   of(constants.ModelArtifactTypeName),
	})
	assertion.Nil(err)
}

func TestMapToModelArtifactMissingType(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToArtifact(&proto.Artifact{
		TypeId: of(modelArtifactTypeId),
	})
	assertion.NotNil(err)
	assertion.Equal("invalid artifact type, can't map from nil", err.Error())
}

func TestMapToArtifactInvalid(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToArtifact(&proto.Artifact{
		TypeId: of(invalidTypeId),
		Type:   of("odh.OtherEntity"),
	})
	assertion.NotNil(err)
	assertion.Equal("unknown artifact type: odh.OtherEntity", err.Error())
}

func TestMapToServingEnvironment(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToServingEnvironment(&proto.Context{
		TypeId: of(servingEnvironmentTypeId),
		Type:   of(constants.ServingEnvironmentTypeName),
	})
	assertion.Nil(err)
}

func TestMapToServingEnvironmentInvalid(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToServingEnvironment(&proto.Context{
		TypeId: of(invalidTypeId),
		Type:   of("odh.OtherEntity"),
	})
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("invalid entity: expected %s but received odh.OtherEntity, please check the provided id", constants.ServingEnvironmentTypeName), err.Error())
}

func TestMapToInferenceService(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToInferenceService(&proto.Context{
		TypeId: of(inferenceServiceTypeId),
		Type:   of(constants.InferenceServiceTypeName),
	})
	assertion.Nil(err)
}

func TestMapToInferenceServiceInvalid(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToInferenceService(&proto.Context{
		TypeId: of(invalidTypeId),
		Type:   of("odh.OtherEntity"),
	})
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("invalid entity: expected %s but received odh.OtherEntity, please check the provided id", constants.InferenceServiceTypeName), err.Error())
}

func TestMapToServeModel(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToServeModel(&proto.Execution{
		TypeId: of(serveModelTypeId),
		Type:   of(constants.ServeModelTypeName),
	})
	assertion.Nil(err)
}

func TestMapToServeModelInvalid(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToServeModel(&proto.Execution{
		TypeId: of(invalidTypeId),
		Type:   of("odh.OtherEntity"),
	})
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("invalid entity: expected %s but received odh.OtherEntity, please check the provided id", constants.ServeModelTypeName), err.Error())
}

func TestMapTo(t *testing.T) {
	_, err := mapTo[*proto.Execution, any](&proto.Execution{TypeId: of(registeredModelTypeId)}, typesMap, "notExisitingTypeName", func(e *proto.Execution) (*any, error) { return nil, nil })
	assert.NotNil(t, err)
	assert.Equal(t, "unknown type name provided: notExisitingTypeName", err.Error())
}

// of returns a pointer to the provided literal/const input
func of[E any](e E) *E {
	return &e
}
