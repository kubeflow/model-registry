// Package proxy provides dynamic routing capabilities for HTTP servers.
//
// This file contains the implementation of a dynamic router that allows
// changing the HTTP handler at runtime in a thread-safe manner. It is
// particularly useful for proxy servers that need to update their routing
// logic wihtout restarting the server.
package proxy

import (
	"net/http"
	"sync"
)

type dynamicRouter struct {
	mu     sync.RWMutex
	router http.Handler
}

func NewDynamicRouter() *dynamicRouter {
	return &dynamicRouter{}
}

func (d *dynamicRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()

	router := d.router

	d.mu.RUnlock()

	router.ServeHTTP(w, r)
}

func (d *dynamicRouter) SetRouter(router http.Handler) {
	d.mu.Lock()

	d.router = router

	d.mu.Unlock()
}
