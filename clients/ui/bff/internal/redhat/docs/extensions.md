# BFF handler extensions

Some downstream builds (for example, the RHOAI dashboard) need to add behavior that does not belong upstream yet. Instead of maintaining long-lived forks, the BFF now exposes a simple extension registry so downstream code can override individual handlers while reusing the rest of the stack.

## Core concepts

- **Handler IDs** – Each overridable endpoint exposes a stable `HandlerID` constant (see `internal/api/extensions.go`). Model Registry Settings routes are wired first.
- **Factories** – Downstream packages call `api.RegisterHandlerOverride(id, factory)` inside `init()`. Factories receive the `*api.App` plus the default handler builder, so you can either replace the handler entirely or fall back to upstream logic.
- **Dependencies** – The `App` exposes read-only accessors for configuration, repositories, and the Kubernetes client factory. It also exposes helper methods such as `BadRequest`, `ServerError`, and `WriteJSON` so overrides can follow the same conventions as upstream handlers.

## Ownership boundaries

- **Upstream-only artifacts** live under `internal/api`, `internal/repositories`, and the default handler tree. These packages must remain vendor-neutral and keep their existing contracts intact so downstream imports keep compiling.
- **Downstream-only artifacts** live under `clients/ui/bff/internal/redhat` (and sibling vendor folders, if they are ever added). Any logic that assumes Red Hat credentials, namespaces, or controllers must stay here so other distributions do not pick it up accidentally.
- **Shared interfaces** (for example repository interfaces or the handler override registry itself) stay upstream. Only implementers that are specific to a vendor move downstream.

Use this rule of thumb: if a change requires Red Hat-only RBAC, Kubernetes resources, or APIs that are invisible to open-source users, keep it downstream. Everything else should be proposed upstream.

## Minimal example

```go
package handlers

import (
    "fmt"
    "net/http"

    "github.com/julienschmidt/httprouter"

    "github.com/kubeflow/model-registry/ui/bff/internal/api"
    "github.com/kubeflow/model-registry/ui/bff/internal/constants"
)

const (
    modelRegistrySettingsListHandlerID = api.HandlerID("modelRegistrySettings:list")
)

func init() {
    api.RegisterHandlerOverride(modelRegistrySettingsListHandlerID, overrideModelRegistrySettingsList)
}

func overrideModelRegistrySettingsList(app *api.App, _ func() httprouter.Handle) httprouter.Handle {
    return app.AttachNamespace(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
        if !ok || namespace == "" {
            app.BadRequest(w, r, fmt.Errorf("missing namespace"))
            return
        }

        client, err := app.KubernetesClientFactory().GetClient(r.Context())
        if err != nil {
            app.ServerError(w, r, fmt.Errorf("failed to build client: %w", err))
            return
        }

        // Use the client to fetch data from Kubernetes and build the response.
        // Downstream repositories live in internal/redhat/repositories.
        _ = client // placeholder for actual implementation

        resp := api.ModelRegistrySettingsListEnvelope{Data: nil}
        if err := app.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
            app.ServerError(w, r, err)
        }
    })
}
```

Package registration is purely compile-time: add a blank import in `clients/ui/bff/cmd/main.go` for the downstream handlers package. When the package is imported, overrides registered in `init()` are automatically active—no configuration flags required. The `buildDefault` parameter is available if an override needs to delegate to upstream logic conditionally.

## Managing downstream overrides

- **Structure** – Place handler factories below `clients/ui/bff/internal/redhat/handlers`, and keep any repository or helper implementations under `clients/ui/bff/internal/redhat/repositories`. This mirrors the upstream layout so the APIs remain familiar.
- **Activation** – Overrides are active whenever their package is imported. Use a blank import (e.g., `_ "github.com/.../internal/redhat/handlers"`) in the main entry point to enable them. No configuration flags are needed.
- **Conditional delegation** – If an override needs to fall back to upstream logic under certain conditions, call `buildDefault()`. Otherwise, the downstream handler runs unconditionally.
- **Shared clients** – Build Kubernetes or database clients via `app.KubernetesClientFactory()` or other upstream factories. Never duplicate client configuration downstream; add capabilities to the upstream factory instead when needed.
- **Testing** – Keep unit and integration tests downstream next to the overrides. Use the upstream interfaces to mock dependencies the same way default handlers do.

## Change workflow

1. **Add the handler ID upstream** – introduce a new `HandlerID` constant and wrap the router registration with `app.handlerWithOverride`. Document the ID in this file under *Current coverage*.
2. **Introduce downstream logic** – implement handler factories (and any supporting repositories) under `clients/ui/bff/internal/redhat`. Register them in the package `init()` by calling `api.RegisterHandlerOverride`.
3. **Wire repositories** – if the downstream handler needs bespoke storage logic, implement a downstream repository that satisfies the upstream interface and expose it via `app.Repositories()` overrides.
4. **Document and test** – update this guide when extending coverage, and add downstream tests to catch regressions before shipping.

## Current coverage

| Handler ID | Path |
|------------|------|
| `modelRegistrySettings:list` | `GET /api/v1/settings/model_registry` |
| `modelRegistrySettings:create` | `POST /api/v1/settings/model_registry` |
| `modelRegistrySettings:get` | `GET /api/v1/settings/model_registry/:model_registry_id` |
| `modelRegistrySettings:update` | `PATCH /api/v1/settings/model_registry/:model_registry_id` |
| `modelRegistrySettings:delete` | `DELETE /api/v1/settings/model_registry/:model_registry_id` |

Additions follow the same pattern: declare a `HandlerID` constant, wrap the router registration with `app.handlerWithOverride`, and document the new ID here so downstream authors know it is available.
