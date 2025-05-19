package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/api"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestOverrideListReturnsRegistries(t *testing.T) {
	factory := &fakeKubeFactory{}
	app := newRedHatTestApp(factory)

	repo := &mockModelRegistrySettingsRepo{
		listFn: func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, selector string) ([]models.ModelRegistryKind, error) {
			if namespace != "test-ns" {
				t.Fatalf("unexpected namespace %s", namespace)
			}
			if selector != "env=prod" {
				t.Fatalf("unexpected label selector %s", selector)
			}
			if ctx == nil {
				t.Fatalf("expected context to be propagated")
			}
			return []models.ModelRegistryKind{{Metadata: models.Metadata{Name: "demo"}}}, nil
		},
	}

	withRepository(t, repo)

	handler := overrideModelRegistrySettingsList(app, failDefault(t))

	req := httptest.NewRequest(http.MethodGet, api.ModelRegistrySettingsListPath+"?namespace=test-ns&labelSelector=env%3Dprod", nil)
	rr := httptest.NewRecorder()

	handler(rr, req, nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp api.ModelRegistrySettingsListEnvelope
	decodeResponse(t, rr, &resp)

	if len(resp.Data) != 1 || resp.Data[0].Metadata.Name != "demo" {
		t.Fatalf("unexpected response payload: %+v", resp.Data)
	}
}

func TestOverrideCreateHonorsDryRun(t *testing.T) {
	factory := &fakeKubeFactory{}
	app := newRedHatTestApp(factory)

	repo := &mockModelRegistrySettingsRepo{
		createFn: func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, payload models.ModelRegistrySettingsPayload, dryRunOnly bool) (models.ModelRegistryKind, error) {
			if !dryRunOnly {
				t.Fatalf("expected dry run flag")
			}
			if namespace != "redhat" {
				t.Fatalf("unexpected namespace %s", namespace)
			}
			return models.ModelRegistryKind{Metadata: models.Metadata{Name: "mr"}}, nil
		},
	}

	withRepository(t, repo)

	body := bytes.NewBufferString(`{"data":{"modelRegistry":{"metadata":{"name":"mr"}}}}`)
	req := httptest.NewRequest(http.MethodPost, api.ModelRegistrySettingsListPath+"?namespace=redhat&dryRun=true", body)
	rr := httptest.NewRecorder()

	handler := overrideModelRegistrySettingsCreate(app, failDefault(t))
	handler(rr, req, nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 for dry run, got %d", rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "" {
		t.Fatalf("dry run should not set location header")
	}
}

func TestOverrideCreateSetsLocationOnSuccess(t *testing.T) {
	factory := &fakeKubeFactory{}
	app := newRedHatTestApp(factory)

	repo := &mockModelRegistrySettingsRepo{
		createFn: func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, payload models.ModelRegistrySettingsPayload, dryRunOnly bool) (models.ModelRegistryKind, error) {
			if dryRunOnly {
				t.Fatalf("unexpected dry run flag")
			}
			return models.ModelRegistryKind{Metadata: models.Metadata{Name: "fresh"}}, nil
		},
	}

	withRepository(t, repo)

	body := bytes.NewBufferString(`{"data":{"modelRegistry":{"metadata":{"name":"fresh"}}}}`)
	req := httptest.NewRequest(http.MethodPost, api.ModelRegistrySettingsListPath+"?namespace=redhat", body)
	rr := httptest.NewRecorder()

	handler := overrideModelRegistrySettingsCreate(app, failDefault(t))
	handler(rr, req, nil)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	expectedLocation := api.ModelRegistrySettingsListPath + "/fresh"
	if rr.Header().Get("Location") != expectedLocation {
		t.Fatalf("unexpected location header %s", rr.Header().Get("Location"))
	}
}

func TestOverrideGetReturnsCredentials(t *testing.T) {
	factory := &fakeKubeFactory{}
	app := newRedHatTestApp(factory)

	repo := &mockModelRegistrySettingsRepo{
		getFn: func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string) (models.ModelRegistryKind, *string, error) {
			if name != "target" {
				t.Fatalf("unexpected name %s", name)
			}
			pwd := "secret"
			return models.ModelRegistryKind{Metadata: models.Metadata{Name: name}}, &pwd, nil
		},
	}

	withRepository(t, repo)

	req := httptest.NewRequest(http.MethodGet, api.ModelRegistrySettingsPath+"?namespace=admin", nil)
	rr := httptest.NewRecorder()

	handler := overrideModelRegistrySettingsGet(app, failDefault(t))
	handler(rr, req, httprouter.Params{{Key: api.ModelRegistryId, Value: "target"}})

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp api.ModelRegistryAndCredentialsSettingsEnvelope
	decodeResponse(t, rr, &resp)

	if resp.Data.DatabasePassword != "secret" {
		t.Fatalf("expected password to be propagated, got %+v", resp.Data)
	}
}

func TestOverrideUpdateUsesProvidedPatch(t *testing.T) {
	factory := &fakeKubeFactory{}
	app := newRedHatTestApp(factory)

	repo := &mockModelRegistrySettingsRepo{
		patchFn: func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string, patch map[string]interface{}, dbPassword *string, ca *string, dryRunOnly bool) (models.ModelRegistryKind, *string, error) {
			if patch["spec"] == nil {
				t.Fatalf("expected spec in patch")
			}
			if dbPassword == nil || *dbPassword != "updated" {
				t.Fatalf("missing password override")
			}
			return models.ModelRegistryKind{Metadata: models.Metadata{Name: name}}, dbPassword, nil
		},
	}

	withRepository(t, repo)

	body := bytes.NewBufferString(`{"data":{"patch":{"spec":{"example":true}},"databasePassword":"updated"}}`)
	req := httptest.NewRequest(http.MethodPatch, api.ModelRegistrySettingsPath+"?namespace=admin", body)
	rr := httptest.NewRecorder()

	handler := overrideModelRegistrySettingsUpdate(app, failDefault(t))
	handler(rr, req, httprouter.Params{{Key: api.ModelRegistryId, Value: "mr"}})

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp api.ModelRegistrySettingsEnvelope
	decodeResponse(t, rr, &resp)
	if resp.Data.Metadata.Name != "mr" {
		t.Fatalf("unexpected response %+v", resp.Data)
	}
}

func TestOverrideUpdateBuildsPatchFromModel(t *testing.T) {
	factory := &fakeKubeFactory{}
	app := newRedHatTestApp(factory)

	repo := &mockModelRegistrySettingsRepo{
		patchFn: func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string, patch map[string]interface{}, dbPassword *string, ca *string, dryRunOnly bool) (models.ModelRegistryKind, *string, error) {
			if patch["built"] != true {
				t.Fatalf("expected custom patch map, got %#v", patch)
			}
			return models.ModelRegistryKind{Metadata: models.Metadata{Name: name}}, nil, nil
		},
	}

	withRepository(t, repo)
	withPatchBuilder(t, func(model models.ModelRegistryKind) (map[string]interface{}, error) {
		if model.Metadata.Name != "mr" {
			t.Fatalf("expected handler to set metadata name")
		}
		return map[string]interface{}{"built": true}, nil
	})

	body := bytes.NewBufferString(`{"data":{"modelRegistry":{"metadata":{}}}}`)
	req := httptest.NewRequest(http.MethodPatch, api.ModelRegistrySettingsPath+"?namespace=admin", body)
	rr := httptest.NewRecorder()

	handler := overrideModelRegistrySettingsUpdate(app, failDefault(t))
	handler(rr, req, httprouter.Params{{Key: api.ModelRegistryId, Value: "mr"}})

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestOverrideDeleteDryRunReturnsStatus(t *testing.T) {
	factory := &fakeKubeFactory{}
	app := newRedHatTestApp(factory)

	repo := &mockModelRegistrySettingsRepo{
		deleteFn: func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string, dryRunOnly bool) (*metav1.Status, error) {
			if !dryRunOnly {
				t.Fatalf("expected dry run delete")
			}
			return &metav1.Status{Status: "Success"}, nil
		},
	}

	withRepository(t, repo)

	req := httptest.NewRequest(http.MethodDelete, api.ModelRegistrySettingsPath+"?namespace=ns&dryRun=true", nil)
	rr := httptest.NewRecorder()

	handler := overrideModelRegistrySettingsDelete(app, failDefault(t))
	handler(rr, req, httprouter.Params{{Key: api.ModelRegistryId, Value: "mr"}})

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var status metav1.Status
	decodeResponse(t, rr, &status)
	if status.Status != "Success" {
		t.Fatalf("unexpected status response %#v", status)
	}
}

func TestOverrideDeleteWritesNoContentOnSuccess(t *testing.T) {
	factory := &fakeKubeFactory{}
	app := newRedHatTestApp(factory)

	repo := &mockModelRegistrySettingsRepo{
		deleteFn: func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string, dryRunOnly bool) (*metav1.Status, error) {
			if dryRunOnly {
				t.Fatalf("unexpected dry run flag")
			}
			return &metav1.Status{Status: "Success"}, nil
		},
	}

	withRepository(t, repo)

	req := httptest.NewRequest(http.MethodDelete, api.ModelRegistrySettingsPath+"?namespace=ns", nil)
	rr := httptest.NewRecorder()

	handler := overrideModelRegistrySettingsDelete(app, failDefault(t))
	handler(rr, req, httprouter.Params{{Key: api.ModelRegistryId, Value: "mr"}})

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func newRedHatTestApp(factory k8s.KubernetesClientFactory) *api.App {
	cfg := config.EnvConfig{AuthMethod: config.AuthMethodUser}
	return api.NewTestApp(cfg, noopLogger(), factory, nil)
}

func noopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func failDefault(t *testing.T) func() httprouter.Handle {
	return func() httprouter.Handle {
		return func(http.ResponseWriter, *http.Request, httprouter.Params) {
			t.Fatalf("default handler should not be invoked")
		}
	}
}

func decodeResponse(t *testing.T, rr *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.Unmarshal(rr.Body.Bytes(), v); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func withRepository(t *testing.T, repo modelRegistrySettingsRepository) {
	t.Helper()
	original := newModelRegistrySettingsRepository
	newModelRegistrySettingsRepository = func(*api.App) modelRegistrySettingsRepository {
		return repo
	}
	t.Cleanup(func() {
		newModelRegistrySettingsRepository = original
	})
}

func withPatchBuilder(t *testing.T, builder func(models.ModelRegistryKind) (map[string]interface{}, error)) {
	t.Helper()
	original := buildModelRegistryPatch
	buildModelRegistryPatch = builder
	t.Cleanup(func() {
		buildModelRegistryPatch = original
	})
}

type fakeKubeFactory struct{}

func (f *fakeKubeFactory) GetClient(ctx context.Context) (k8s.KubernetesClientInterface, error) {
	return nil, nil
}

func (f *fakeKubeFactory) ExtractRequestIdentity(http.Header) (*k8s.RequestIdentity, error) {
	return nil, nil
}

func (f *fakeKubeFactory) ValidateRequestIdentity(*k8s.RequestIdentity) error {
	return nil
}

type mockModelRegistrySettingsRepo struct {
	listFn   func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, selector string) ([]models.ModelRegistryKind, error)
	createFn func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, payload models.ModelRegistrySettingsPayload, dryRunOnly bool) (models.ModelRegistryKind, error)
	getFn    func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string) (models.ModelRegistryKind, *string, error)
	patchFn  func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string, patch map[string]interface{}, dbPassword *string, ca *string, dryRunOnly bool) (models.ModelRegistryKind, *string, error)
	deleteFn func(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string, dryRunOnly bool) (*metav1.Status, error)
}

func (m *mockModelRegistrySettingsRepo) List(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, selector string) ([]models.ModelRegistryKind, error) {
	if m.listFn == nil {
		return nil, errors.New("list not implemented")
	}
	return m.listFn(ctx, client, namespace, selector)
}

func (m *mockModelRegistrySettingsRepo) Create(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, payload models.ModelRegistrySettingsPayload, dryRunOnly bool) (models.ModelRegistryKind, error) {
	if m.createFn == nil {
		return models.ModelRegistryKind{}, errors.New("create not implemented")
	}
	return m.createFn(ctx, client, namespace, payload, dryRunOnly)
}

func (m *mockModelRegistrySettingsRepo) Get(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string) (models.ModelRegistryKind, *string, error) {
	if m.getFn == nil {
		return models.ModelRegistryKind{}, nil, errors.New("get not implemented")
	}
	return m.getFn(ctx, client, namespace, name)
}

func (m *mockModelRegistrySettingsRepo) Patch(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string, patch map[string]interface{}, dbPassword *string, ca *string, dryRunOnly bool) (models.ModelRegistryKind, *string, error) {
	if m.patchFn == nil {
		return models.ModelRegistryKind{}, nil, errors.New("patch not implemented")
	}
	return m.patchFn(ctx, client, namespace, name, patch, dbPassword, ca, dryRunOnly)
}

func (m *mockModelRegistrySettingsRepo) Delete(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string, dryRunOnly bool) (*metav1.Status, error) {
	if m.deleteFn == nil {
		return nil, errors.New("delete not implemented")
	}
	return m.deleteFn(ctx, client, namespace, name, dryRunOnly)
}
