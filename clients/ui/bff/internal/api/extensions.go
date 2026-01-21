package api

import (
	"log/slog"
	"sync"

	"github.com/julienschmidt/httprouter"
)

// HandlerID identifies an overridable HTTP handler.
// TODO(upstream): Keep this type exported so downstream code can reference arbitrary handler keys without
// requiring upstream to enumerate them.
type HandlerID string

// HandlerFactory builds a router handler that has access to the App instance.
// Implementations can opt to call buildDefault() to reuse the default upstream handler.
type HandlerFactory func(app *App, buildDefault func() httprouter.Handle) httprouter.Handle

var (
	handlerOverrideMu sync.RWMutex
	handlerOverrides  = map[HandlerID]HandlerFactory{}
)

// RegisterHandlerOverride allows downstream code to override a specific handler.
// TODO(downstream): Call this from vendor packages (e.g., internal/redhat/handlers) to inject custom behavior.
// Calling this function multiple times for the same id will replace the previous override.
func RegisterHandlerOverride(id HandlerID, factory HandlerFactory) {
	handlerOverrideMu.Lock()
	defer handlerOverrideMu.Unlock()
	handlerOverrides[id] = factory
}

//nolint:unused // Used by downstream implementations
func getHandlerOverride(id HandlerID) HandlerFactory {
	handlerOverrideMu.RLock()
	defer handlerOverrideMu.RUnlock()
	return handlerOverrides[id]
}

// handlerWithOverride returns the handler registered for the given id or builds the default one.
// TODO(upstream): This glue stays upstream so the router keeps working even when no downstream overrides exist.
//
//nolint:unused // Used by downstream implementations
func (app *App) handlerWithOverride(id HandlerID, buildDefault func() httprouter.Handle) httprouter.Handle {
	if factory := getHandlerOverride(id); factory != nil {
		app.logHandlerOverride(id)
		return factory(app, buildDefault)
	}
	return buildDefault()
}

//nolint:unused // Used by downstream implementations
func (app *App) logHandlerOverride(id HandlerID) {
	if app == nil || app.logger == nil {
		return
	}
	app.logger.Debug("Using handler override", slog.String("handlerID", string(id)))
}
