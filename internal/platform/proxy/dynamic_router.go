package proxy

import (
	"net/http"
	"sync"
)

type DynamicRouter struct {
	mu     sync.RWMutex
	router http.Handler
}

func NewDynamicRouter() *DynamicRouter {
	return &DynamicRouter{
		router: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "router not configured", http.StatusServiceUnavailable)
		}),
	}
}

func (d *DynamicRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()

	router := d.router

	d.mu.RUnlock()

	router.ServeHTTP(w, r)
}

func (d *DynamicRouter) SetRouter(router http.Handler) {
	d.mu.Lock()

	d.router = router

	d.mu.Unlock()
}
