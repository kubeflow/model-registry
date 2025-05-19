package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/kubeflow/model-registry/ui/bff/internal/api"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	redhatrepos "github.com/kubeflow/model-registry/ui/bff/internal/redhat/repositories"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// TODO(downstream): Keep handler IDs in sync with the upstream router keys until the API exposes a catalog.
	modelRegistrySettingsListHandlerID   = api.HandlerID("modelRegistrySettings:list")
	modelRegistrySettingsCreateHandlerID = api.HandlerID("modelRegistrySettings:create")
	modelRegistrySettingsGetHandlerID    = api.HandlerID("modelRegistrySettings:get")
	modelRegistrySettingsUpdateHandlerID = api.HandlerID("modelRegistrySettings:update")
	modelRegistrySettingsDeleteHandlerID = api.HandlerID("modelRegistrySettings:delete")
)

func init() {
	registerModelRegistrySettingsOverride(modelRegistrySettingsListHandlerID, overrideModelRegistrySettingsList)
	registerModelRegistrySettingsOverride(modelRegistrySettingsCreateHandlerID, overrideModelRegistrySettingsCreate)
	registerModelRegistrySettingsOverride(modelRegistrySettingsGetHandlerID, overrideModelRegistrySettingsGet)
	registerModelRegistrySettingsOverride(modelRegistrySettingsUpdateHandlerID, overrideModelRegistrySettingsUpdate)
	registerModelRegistrySettingsOverride(modelRegistrySettingsDeleteHandlerID, overrideModelRegistrySettingsDelete)
}

type modelRegistrySettingsRepository interface {
	List(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, labelSelector string) ([]models.ModelRegistryKind, error)
	Create(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, payload models.ModelRegistrySettingsPayload, dryRunOnly bool) (models.ModelRegistryKind, error)
	Get(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string) (models.ModelRegistryKind, *string, error)
	Patch(
		ctx context.Context,
		client k8s.KubernetesClientInterface,
		namespace string,
		name string,
		patch map[string]interface{},
		databasePassword *string,
		newCACertificate *string,
		dryRunOnly bool,
	) (models.ModelRegistryKind, *string, error)
	Delete(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, name string, dryRunOnly bool) (*metav1.Status, error)
}

var (
	newModelRegistrySettingsRepository = func(app *api.App) modelRegistrySettingsRepository {
		if app == nil {
			return redhatrepos.NewModelRegistrySettingsRepository(nil)
		}
		return redhatrepos.NewModelRegistrySettingsRepository(app.Logger())
	}

	buildModelRegistryPatch = redhatrepos.BuildModelRegistryPatch
)

func overrideModelRegistrySettingsList(app *api.App, buildDefault func() httprouter.Handle) httprouter.Handle {
	if !shouldUseRedHatOverrides(app) {
		return buildDefault()
	}

	repo := newModelRegistrySettingsRepository(app)

	return app.AttachNamespace(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		namespace, ok := namespaceFromContext(app, w, r)
		if !ok {
			return
		}

		client, ok := getKubernetesClient(app, w, r)
		if !ok {
			return
		}

		labelSelector := getLabelSelector(r)

		registries, err := repo.List(r.Context(), client, namespace, labelSelector)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		resp := api.ModelRegistrySettingsListEnvelope{Data: registries}
		if err := app.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
			app.ServerError(w, r, err)
		}
	})
}

func overrideModelRegistrySettingsCreate(app *api.App, buildDefault func() httprouter.Handle) httprouter.Handle {
	if !shouldUseRedHatOverrides(app) {
		return buildDefault()
	}

	repo := newModelRegistrySettingsRepository(app)

	return app.AttachNamespace(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		namespace, ok := namespaceFromContext(app, w, r)
		if !ok {
			return
		}

		client, ok := getKubernetesClient(app, w, r)
		if !ok {
			return
		}

		dryRunOnly := isDryRunRequest(r)

		var envelope api.ModelRegistrySettingsPayloadEnvelope
		if err := app.ReadJSON(w, r, &envelope); err != nil {
			app.BadRequest(w, r, err)
			return
		}

		result, err := repo.Create(r.Context(), client, namespace, envelope.Data, dryRunOnly)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		headers := http.Header{}
		status := http.StatusCreated
		if dryRunOnly {
			status = http.StatusOK
		} else if location := buildResourceLocation(r, result.Metadata.Name); location != "" {
			headers.Set("Location", location)
		}

		resp := api.ModelRegistrySettingsEnvelope{Data: result}
		if err := app.WriteJSON(w, status, resp, headers); err != nil {
			app.ServerError(w, r, err)
		}
	})
}

func overrideModelRegistrySettingsGet(app *api.App, buildDefault func() httprouter.Handle) httprouter.Handle {
	if !shouldUseRedHatOverrides(app) {
		return buildDefault()
	}

	repo := newModelRegistrySettingsRepository(app)

	return app.AttachNamespace(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		namespace, ok := namespaceFromContext(app, w, r)
		if !ok {
			return
		}

		name := strings.TrimSpace(ps.ByName(api.ModelRegistryId))
		if name == "" {
			app.BadRequest(w, r, fmt.Errorf("missing model registry identifier"))
			return
		}

		client, ok := getKubernetesClient(app, w, r)
		if !ok {
			return
		}

		registry, password, err := repo.Get(r.Context(), client, namespace, name)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		payload := models.ModelRegistryAndCredentials{ModelRegistry: registry}
		if password != nil {
			payload.DatabasePassword = *password
		}

		resp := api.ModelRegistryAndCredentialsSettingsEnvelope{Data: payload}
		if err := app.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
			app.ServerError(w, r, err)
		}
	})
}

func overrideModelRegistrySettingsUpdate(app *api.App, buildDefault func() httprouter.Handle) httprouter.Handle {
	if !shouldUseRedHatOverrides(app) {
		return buildDefault()
	}

	repo := newModelRegistrySettingsRepository(app)

	return app.AttachNamespace(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		namespace, ok := namespaceFromContext(app, w, r)
		if !ok {
			return
		}

		name := strings.TrimSpace(ps.ByName(api.ModelRegistryId))
		if name == "" {
			app.BadRequest(w, r, fmt.Errorf("missing model registry identifier"))
			return
		}

		client, ok := getKubernetesClient(app, w, r)
		if !ok {
			return
		}

		var envelope modelRegistrySettingsPatchEnvelope
		if err := app.ReadJSON(w, r, &envelope); err != nil {
			app.BadRequest(w, r, err)
			return
		}

		patchBody := envelope.Data.Patch
		if len(patchBody) == 0 {
			if envelope.Data.ModelRegistry == nil {
				app.BadRequest(w, r, fmt.Errorf("patch body or modelRegistry payload is required"))
				return
			}
			modelCopy := *envelope.Data.ModelRegistry
			if modelCopy.Metadata.Name == "" {
				modelCopy.Metadata.Name = name
			}
			builtPatch, err := buildModelRegistryPatch(modelCopy)
			if err != nil {
				app.ServerError(w, r, err)
				return
			}
			patchBody = builtPatch
		}

		dryRunOnly := isDryRunRequest(r)

		result, _, err := repo.Patch(
			r.Context(),
			client,
			namespace,
			name,
			patchBody,
			envelope.Data.DatabasePassword,
			envelope.Data.NewDatabaseCACertificate,
			dryRunOnly,
		)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		resp := api.ModelRegistrySettingsEnvelope{Data: result}
		if err := app.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
			app.ServerError(w, r, err)
		}
	})
}

func overrideModelRegistrySettingsDelete(app *api.App, buildDefault func() httprouter.Handle) httprouter.Handle {
	if !shouldUseRedHatOverrides(app) {
		return buildDefault()
	}

	repo := newModelRegistrySettingsRepository(app)

	return app.AttachNamespace(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		namespace, ok := namespaceFromContext(app, w, r)
		if !ok {
			return
		}

		name := strings.TrimSpace(ps.ByName(api.ModelRegistryId))
		if name == "" {
			app.BadRequest(w, r, fmt.Errorf("missing model registry identifier"))
			return
		}

		client, ok := getKubernetesClient(app, w, r)
		if !ok {
			return
		}

		dryRunOnly := isDryRunRequest(r)

		statusObj, err := repo.Delete(r.Context(), client, namespace, name, dryRunOnly)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		if dryRunOnly {
			if err := app.WriteJSON(w, http.StatusOK, statusObj, nil); err != nil {
				app.ServerError(w, r, err)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

type modelRegistrySettingsPatchEnvelope struct {
	Data modelRegistrySettingsPatchPayload `json:"data"`
}

type modelRegistrySettingsPatchPayload struct {
	ModelRegistry            *models.ModelRegistryKind `json:"modelRegistry,omitempty"`
	Patch                    map[string]interface{}    `json:"patch,omitempty"`
	DatabasePassword         *string                   `json:"databasePassword,omitempty"`
	NewDatabaseCACertificate *string                   `json:"newDatabaseCACertificate,omitempty"`
}

func registerModelRegistrySettingsOverride(id api.HandlerID, factory api.HandlerFactory) {
	api.RegisterHandlerOverride(id, factory)
}

// shouldUseRedHatOverrides returns true when the downstream override should run.
// The override is disabled when running in mock mode (MockK8Client=true) so that
// stub endpoints are used instead for local development without a real cluster.
func shouldUseRedHatOverrides(app *api.App) bool {
	if app == nil {
		return false
	}
	// When K8s client is mocked, use the upstream stub handlers instead of real implementations
	return !app.Config().MockK8Client
}

func namespaceFromContext(app *api.App, w http.ResponseWriter, r *http.Request) (string, bool) {
	if r == nil {
		app.BadRequest(w, r, fmt.Errorf("missing request context"))
		return "", false
	}
	namespace, _ := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if namespace == "" {
		app.BadRequest(w, r, fmt.Errorf("missing namespace in the context"))
		return "", false
	}
	return namespace, true
}

func getKubernetesClient(app *api.App, w http.ResponseWriter, r *http.Request) (k8s.KubernetesClientInterface, bool) {
	client, err := app.KubernetesClientFactory().GetClient(r.Context())
	if err != nil {
		app.ServerError(w, r, fmt.Errorf("failed to get Kubernetes client: %w", err))
		return nil, false
	}
	return client, true
}

func getLabelSelector(r *http.Request) string {
	if r == nil {
		return ""
	}
	return r.URL.Query().Get("labelSelector")
}

func isDryRunRequest(r *http.Request) bool {
	if r == nil {
		return false
	}
	value := r.URL.Query().Get("dryRun")
	if value == "" {
		return false
	}
	dryRun, err := strconv.ParseBool(value)
	return err == nil && dryRun
}

func buildResourceLocation(r *http.Request, resourceName string) string {
	if r == nil || resourceName == "" {
		return ""
	}
	base := strings.TrimSuffix(r.URL.Path, "/")
	if base == "" {
		base = "/"
	}
	return fmt.Sprintf("%s/%s", base, resourceName)
}
