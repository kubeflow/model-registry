# BFF Handler Extensions

Some downstream builds need to add behavior that does not belong upstream yet. Instead of maintaining long-lived forks, the BFF exposes a simple extension registry so downstream code can override individual handlers while reusing the rest of the stack.

## Table of Contents

- [Core Concepts](#core-concepts)
- [Ownership Boundaries](#ownership-boundaries)
- [Available App Methods](#available-app-methods)
- [Step-by-Step: Overriding an Existing Handler](#step-by-step-overriding-an-existing-handler)
- [Step-by-Step: Adding a Downstream-Only Route](#step-by-step-adding-a-downstream-only-route)
- [Accessing the Kubernetes Client](#accessing-the-kubernetes-client)
- [Importing Upstream Types](#importing-upstream-types)
- [Testing Downstream Overrides](#testing-downstream-overrides)
- [Current Handler ID Coverage](#current-handler-id-coverage)

---

## Core Concepts

- **Handler IDs** – Each overridable endpoint has a stable `HandlerID` string constant defined upstream in `internal/api/app.go`. Downstream code references these IDs to register overrides.
- **Handler Factories** – Downstream packages call `api.RegisterHandlerOverride(id, factory)` inside `init()`. Factories receive the `*api.App` plus a `buildDefault` function, so you can either replace the handler entirely or fall back to upstream logic conditionally.
- **Downstream-Only Routes** – Some endpoints exist upstream only as stubs returning 501 Not Implemented. These are designed to be fully implemented downstream. The route is wired upstream, but the real logic lives in the override.
- **Dependencies** – The `App` exposes read-only accessors for configuration, repositories, and the Kubernetes client factory. It also exposes helper methods such as `BadRequest`, `ServerError`, and `WriteJSON` so overrides can follow the same conventions as upstream handlers.

---

## Ownership Boundaries

- **Upstream-only artifacts** live under `internal/api`, `internal/repositories`, and the default handler tree. These packages must remain vendor-neutral and keep their existing contracts intact so downstream imports keep compiling.
- **Downstream-only artifacts** live under `clients/ui/bff/internal/<vendor>/` (e.g., `internal/redhat/`). Any logic that assumes vendor-specific credentials, namespaces, or controllers must stay here so other distributions do not pick it up accidentally.
- **Shared interfaces** (for example, repository interfaces or the handler override registry itself) stay upstream. Only implementations specific to a vendor move downstream.

**Rule of thumb:** If a change requires vendor-only RBAC, Kubernetes resources, or APIs invisible to open-source users, keep it downstream. Everything else should be proposed upstream.

---

## Available App Methods

The `*api.App` instance provides these exported methods for use in downstream handlers:

| Method | Description |
|--------|-------------|
| `app.Config()` | Returns `config.EnvConfig` with deployment settings |
| `app.Logger()` | Returns `*slog.Logger` for structured logging |
| `app.KubernetesClientFactory()` | Returns factory to build Kubernetes clients |
| `app.Repositories()` | Returns `*repositories.Repositories` for data access |
| `app.WriteJSON(w, status, data, headers)` | Writes JSON response with proper content-type |
| `app.ReadJSON(w, r, dst)` | Parses JSON request body into destination struct |
| `app.BadRequest(w, r, err)` | Returns HTTP 400 with error message |
| `app.ServerError(w, r, err)` | Returns HTTP 500 with error message |
| `app.NotImplemented(w, r, feature)` | Returns HTTP 501 for unimplemented features |
| `app.AttachNamespace(handler)` | Middleware that extracts `namespace` query param into context |

---

## Step-by-Step: Overriding an Existing Handler

This workflow overrides an upstream handler with downstream-specific logic.

### 1. Add the blank import in `main.go`

In `clients/ui/bff/cmd/main.go`, add a blank import for your handlers package:

```go
import (
    // ... other imports ...

    // Import downstream handlers to register overrides via init()
    _ "github.com/kubeflow/model-registry/ui/bff/internal/redhat/handlers"
)
```

### 2. Locate the Handler ID upstream

Handler IDs are defined in `internal/api/app.go`. For example, Model Registry Settings endpoints:

```go
const (
    handlerModelRegistrySettingsListID   HandlerID = "modelRegistrySettings:list"
    handlerModelRegistrySettingsCreateID HandlerID = "modelRegistrySettings:create"
    handlerModelRegistrySettingsGetID    HandlerID = "modelRegistrySettings:get"
    handlerModelRegistrySettingsUpdateID HandlerID = "modelRegistrySettings:update"
    handlerModelRegistrySettingsDeleteID HandlerID = "modelRegistrySettings:delete"
)
```

### 3. Create the downstream handler file

Create a file in `internal/<vendor>/handlers/` (e.g., `internal/redhat/handlers/my_handler.go`):

```go
package handlers

import (
    "fmt"
    "net/http"

    "github.com/julienschmidt/httprouter"

    "github.com/kubeflow/model-registry/ui/bff/internal/api"
    "github.com/kubeflow/model-registry/ui/bff/internal/constants"
)

// Mirror the upstream handler ID string
const myHandlerID = api.HandlerID("modelRegistrySettings:list")

func init() {
    api.RegisterHandlerOverride(myHandlerID, myOverrideFactory)
}

func myOverrideFactory(app *api.App, buildDefault func() httprouter.Handle) httprouter.Handle {
    // Optionally fall back to upstream when mocking
    if app.Config().MockK8Client {
        return buildDefault()
    }

    return app.AttachNamespace(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
        if !ok || namespace == "" {
            app.BadRequest(w, r, fmt.Errorf("missing namespace"))
            return
        }

        // Your downstream implementation here
        resp := api.ModelRegistrySettingsListEnvelope{Data: nil}
        if err := app.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
            app.ServerError(w, r, err)
        }
    })
}
```

---

## Step-by-Step: Adding a Downstream-Only Route

For routes that have **no real upstream implementation** (upstream returns 501 Not Implemented), use this pattern.

### 1. Upstream defines the route and stub handler

In `internal/api/app.go`, the route is defined with a handler that returns 501:

```go
// Path constant
const KubernetesServicesListPath = SettingsPath + "/services"

// Handler ID
const handlerKubernetesServicesListID HandlerID = "kubernetes:services:list"

// Route registration in Routes()
apiRouter.GET(
    KubernetesServicesListPath,
    app.handlerWithOverride(handlerKubernetesServicesListID, func() httprouter.Handle {
        return app.AttachNamespace(app.kubernetesServicesNotImplementedHandler)
    }),
)
```

The stub handler in `internal/api/kubernetes_resources_handler.go`:

```go
// Generic handler for endpoints not implemented upstream
func (app *App) endpointNotImplementedHandler(feature string) func(http.ResponseWriter, *http.Request, httprouter.Params) {
    return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        app.NotImplemented(w, r, feature)
    }
}
```

### 2. Downstream provides the real implementation

In `internal/redhat/handlers/kubernetes_services.go`:

```go
package handlers

import (
    "fmt"
    "net/http"

    "github.com/julienschmidt/httprouter"
    "github.com/kubeflow/model-registry/ui/bff/internal/api"
    "github.com/kubeflow/model-registry/ui/bff/internal/constants"
)

const kubernetesServicesListHandlerID = api.HandlerID("kubernetes:services:list")

func init() {
    api.RegisterHandlerOverride(kubernetesServicesListHandlerID, overrideKubernetesServicesList)
}

func overrideKubernetesServicesList(app *api.App, buildDefault func() httprouter.Handle) httprouter.Handle {
    // Fall back to stub when K8s client is mocked
    if app.Config().MockK8Client {
        return buildDefault()
    }

    return app.AttachNamespace(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
        if !ok || namespace == "" {
            app.BadRequest(w, r, fmt.Errorf("missing namespace"))
            return
        }

        client, err := app.KubernetesClientFactory().GetClient(r.Context())
        if err != nil {
            app.ServerError(w, r, fmt.Errorf("failed to get client: %w", err))
            return
        }

        // Real implementation using the Kubernetes client
        services, err := client.GetServiceDetails(r.Context(), namespace)
        if err != nil {
            app.ServerError(w, r, err)
            return
        }

        // Build and return response
        items := make([]api.KubernetesServiceItem, 0, len(services))
        for _, svc := range services {
            items = append(items, api.KubernetesServiceItem{
                Name:      svc.Name,
                Namespace: namespace,
            })
        }

        resp := api.KubernetesServicesListEnvelope{Data: items}
        if err := app.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
            app.ServerError(w, r, err)
        }
    })
}
```

---

## Accessing the Kubernetes Client

To interact with Kubernetes resources in your handler:

```go
func myHandler(app *api.App) httprouter.Handle {
    return app.AttachNamespace(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        // Get the Kubernetes client from the factory
        client, err := app.KubernetesClientFactory().GetClient(r.Context())
        if err != nil {
            app.ServerError(w, r, fmt.Errorf("failed to get Kubernetes client: %w", err))
            return
        }

        // Use the client to interact with Kubernetes
        // The client interface is defined in internal/integrations/kubernetes
        namespace := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
        
        // Example: Get service details
        services, err := client.GetServiceDetails(r.Context(), namespace)
        if err != nil {
            app.ServerError(w, r, err)
            return
        }
        
        // Process services...
    })
}
```

The `KubernetesClientInterface` provides methods for common operations. Check `internal/integrations/kubernetes/client.go` for available methods.

---

## Importing Upstream Types

When building downstream handlers, you'll commonly import these packages:

```go
import (
    // Core API types and App
    "github.com/kubeflow/model-registry/ui/bff/internal/api"
    
    // Configuration types
    "github.com/kubeflow/model-registry/ui/bff/internal/config"
    
    // Context keys for namespace, identity, etc.
    "github.com/kubeflow/model-registry/ui/bff/internal/constants"
    
    // Kubernetes client interface
    k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
    
    // Data models
    "github.com/kubeflow/model-registry/ui/bff/internal/models"
    
    // Upstream repositories (if needed)
    "github.com/kubeflow/model-registry/ui/bff/internal/repositories"
    
    // Router
    "github.com/julienschmidt/httprouter"
)
```

**Commonly used types from `api` package:**

- `api.App` – Main application instance
- `api.HandlerID` – Type for handler identifiers
- `api.HandlerFactory` – Signature for override factories
- `api.RegisterHandlerOverride()` – Register an override
- Envelope types like `api.ModelRegistrySettingsListEnvelope`, `api.KubernetesServicesListEnvelope`

---

## Testing Downstream Overrides

Keep unit tests next to your override implementations in the downstream package.

### Testing with mocked Kubernetes client

When `MockK8Client=true`, your handlers should fall back to upstream stubs or return mock data:

```go
func TestMyHandler_MockMode(t *testing.T) {
    // Create app with mocked K8s client
    cfg := config.EnvConfig{
        MockK8Client: true,
    }
    
    // Your handler factory should return buildDefault() or mock behavior
    // when app.Config().MockK8Client is true
}
```

### Pattern for conditional override activation

```go
func shouldUseDownstreamOverrides(app *api.App) bool {
    if app == nil {
        return false
    }
    // When K8s client is mocked, use upstream stub handlers
    return !app.Config().MockK8Client
}

func myOverrideFactory(app *api.App, buildDefault func() httprouter.Handle) httprouter.Handle {
    // Fall back to upstream stub when K8s is mocked
    if !shouldUseDownstreamOverrides(app) {
        return buildDefault()
    }
    
    // Real implementation
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        // ...
    }
}
```

### Test file structure

```
internal/redhat/handlers/
├── model_registry_settings.go
├── model_registry_settings_test.go
├── kubernetes_services.go
└── kubernetes_services_test.go
```

---

## Current Handler ID Coverage

These handler IDs are currently wired with `handlerWithOverride()` upstream:

| Handler ID | HTTP Method | Path | Description |
|------------|-------------|------|-------------|
| `modelRegistrySettings:list` | GET | `/api/v1/settings/model_registry` | List all model registries |
| `modelRegistrySettings:create` | POST | `/api/v1/settings/model_registry` | Create a model registry |
| `modelRegistrySettings:get` | GET | `/api/v1/settings/model_registry/:model_registry_id` | Get a model registry |
| `modelRegistrySettings:update` | PATCH | `/api/v1/settings/model_registry/:model_registry_id` | Update a model registry |
| `modelRegistrySettings:delete` | DELETE | `/api/v1/settings/model_registry/:model_registry_id` | Delete a model registry |
| `kubernetes:services:list` | GET | `/api/v1/settings/services` | List Kubernetes services (downstream-only) |

---

## Change Workflow Summary

### To override an existing handler:

1. Find the handler ID in `internal/api/app.go`
2. Create handler file in `internal/<vendor>/handlers/`
3. Register override in `init()` with `api.RegisterHandlerOverride()`
4. Add blank import in `cmd/main.go` (if not already present)

### To add a downstream-only route:

1. Define route path constant and handler ID in `internal/api/app.go`
2. Create stub handler returning 501 in `internal/api/`
3. Wire route with `handlerWithOverride()` in `Routes()`
4. Implement real handler in `internal/<vendor>/handlers/`
5. Register override in `init()` with `api.RegisterHandlerOverride()`
