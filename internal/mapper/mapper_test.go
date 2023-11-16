package mapper

import (
	"fmt"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/stretchr/testify/assert"
)

const (
	invalidTypeId            = int64(9999)
	registeredModelTypeId    = int64(1)
	modelVersionTypeId       = int64(2)
	modelArtifactTypeId      = int64(3)
	servingEnvironmentTypeId = int64(4)
	inferenceServiceTypeId   = int64(5)
	serveModelTypeId         = int64(6)
)

func setup(t *testing.T) (*assert.Assertions, *Mapper) {
	return assert.New(t), NewMapper(
		registeredModelTypeId,
		modelVersionTypeId,
		modelArtifactTypeId,
		servingEnvironmentTypeId,
		inferenceServiceTypeId,
		serveModelTypeId,
	)
}

func TestMapToRegisteredModel(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToRegisteredModel(&proto.Context{
		TypeId: of(registeredModelTypeId),
	})
	assertion.Nil(err)
}

func TestMapToRegisteredModelInvalid(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToRegisteredModel(&proto.Context{
		TypeId: of(invalidTypeId),
	})
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("invalid TypeId, expected %d but received %d", registeredModelTypeId, invalidTypeId), err.Error())
}

func TestMapToModelVersion(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToModelVersion(&proto.Context{
		TypeId: of(modelVersionTypeId),
	})
	assertion.Nil(err)
}

func TestMapToModelVersionInvalid(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToModelVersion(&proto.Context{
		TypeId: of(invalidTypeId),
	})
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("invalid TypeId, expected %d but received %d", modelVersionTypeId, invalidTypeId), err.Error())
}

func TestMapToModelArtifact(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToModelArtifact(&proto.Artifact{
		TypeId: of(modelArtifactTypeId),
		Type:   of("odh.ModelArtifact"),
	})
	assertion.Nil(err)
}

func TestMapToModelArtifactMissingType(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToModelArtifact(&proto.Artifact{
		TypeId: of(modelArtifactTypeId),
	})
	assertion.NotNil(err)
	assertion.Equal("error setting field ArtifactType: invalid artifact type found: <nil>", err.Error())
}

func TestMapToModelArtifactInvalid(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToModelArtifact(&proto.Artifact{
		TypeId: of(invalidTypeId),
	})
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("invalid TypeId, expected %d but received %d", modelArtifactTypeId, invalidTypeId), err.Error())
}

func TestMapToServingEnvironment(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToServingEnvironment(&proto.Context{
		TypeId: of(servingEnvironmentTypeId),
	})
	assertion.Nil(err)
}

func TestMapToServingEnvironmentInvalid(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToServingEnvironment(&proto.Context{
		TypeId: of(invalidTypeId),
	})
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("invalid TypeId, expected %d but received %d", servingEnvironmentTypeId, invalidTypeId), err.Error())
}

func TestMapToInferenceService(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToInferenceService(&proto.Context{
		TypeId: of(inferenceServiceTypeId),
	})
	assertion.Nil(err)
}

func TestMapToInferenceServiceInvalid(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToInferenceService(&proto.Context{
		TypeId: of(invalidTypeId),
	})
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("invalid TypeId, expected %d but received %d", inferenceServiceTypeId, invalidTypeId), err.Error())
}

func TestMapToServeModel(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToServeModel(&proto.Execution{
		TypeId: of(serveModelTypeId),
	})
	assertion.Nil(err)
}

func TestMapToServeModelInvalid(t *testing.T) {
	assertion, m := setup(t)
	_, err := m.MapToServeModel(&proto.Execution{
		TypeId: of(invalidTypeId),
	})
	assertion.NotNil(err)
	assertion.Equal(fmt.Sprintf("invalid TypeId, expected %d but received %d", serveModelTypeId, invalidTypeId), err.Error())
}

// of returns a pointer to the provided literal/const input
func of[E any](e E) *E {
	return &e
}
