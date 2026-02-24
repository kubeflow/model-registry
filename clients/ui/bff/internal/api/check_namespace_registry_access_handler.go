package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
)

const CheckNamespaceRegistryAccessPath = ApiPathPrefix + "/check-namespace-registry-access"

type CheckNamespaceRegistryAccessRequest struct {
	Namespace         string `json:"namespace"`
	RegistryName      string `json:"registryName"`
	RegistryNamespace string `json:"registryNamespace"`
}

type CheckNamespaceRegistryAccessRequestEnvelope Envelope[CheckNamespaceRegistryAccessRequest, None]

type CheckNamespaceRegistryAccessResponse struct {
	HasAccess bool `json:"hasAccess"`
}

type CheckNamespaceRegistryAccessEnvelope Envelope[CheckNamespaceRegistryAccessResponse, None]

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
