package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
)

const CheckNamespaceRegistryAccessPath = ApiPathPrefix + "/check-namespace-registry-access"

// CheckNamespaceRegistryAccessRequest is the request body for the namespace registry access check.
type CheckNamespaceRegistryAccessRequest struct {
	Namespace         string `json:"namespace"`
	RegistryName      string `json:"registryName"`
	RegistryNamespace string `json:"registryNamespace"`
}

type CheckNamespaceRegistryAccessRequestEnvelope Envelope[CheckNamespaceRegistryAccessRequest, None]

// CheckNamespaceRegistryAccessResponse is the response body.
type CheckNamespaceRegistryAccessResponse struct {
	HasAccess bool `json:"hasAccess"`
}

type CheckNamespaceRegistryAccessEnvelope Envelope[CheckNamespaceRegistryAccessResponse, None]

// CheckNamespaceRegistryAccessHandler checks if the default SA in the given namespace
// can get the model registry service in the registry namespace (SubjectAccessReview).
// Uses the logged-in user's token (no privilege elevation).
func (app *App) CheckNamespaceRegistryAccessHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	identity, ok := ctx.Value(constants.RequestIdentityKey).(*kubernetes.RequestIdentity)
	if !ok || identity == nil {
		app.badRequestResponse(w, r, fmt.Errorf("missing RequestIdentity in context"))
		return
	}

	if err := app.kubernetesClientFactory.ValidateRequestIdentity(identity); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var envelope CheckNamespaceRegistryAccessRequestEnvelope
	if err := app.ReadJSON(w, r, &envelope); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	req := envelope.Data

	if req.Namespace == "" || req.RegistryName == "" || req.RegistryNamespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("namespace, registryName and registryNamespace are required"))
		return
	}

	client, err := app.kubernetesClientFactory.GetClient(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("failed to get Kubernetes client: %w", err))
		return
	}

	hasAccess, err := client.CanNamespaceAccessRegistry(ctx, identity, req.Namespace, req.RegistryName, req.RegistryNamespace)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("namespace registry access check failed: %w", err))
		return
	}

	resp := CheckNamespaceRegistryAccessEnvelope{
		Data: CheckNamespaceRegistryAccessResponse{HasAccess: hasAccess},
	}
	if err := app.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
