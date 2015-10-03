package contextrouter

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

// Method denotes a HTTP method to be specified to the Router
type Method string

// Supported HTTP Methods for use with Server.Handle and Server.HandlerFunc
const (
	GET     Method = "GET"
	DELETE         = "DELETE"
	HEAD           = "HEAD"
	OPTIONS        = "OPTIONS"
	PATCH          = "PATHCH"
	POST           = "POST"
	PUT            = "PUT"
)

// ContextHandler is analogous to http.Handler but takes a Context as the
// first parameter. Shutdown signals, named routing parameters and any global
// key value settings are passed via context. Handlers can pass the context
// to other go functions across API bounndries
type ContextHandler interface {
	ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request)
}

// ContextHandlerFunc is analogous to http.HandlerFunc similar to ContextHandler
type ContextHandlerFunc func(c context.Context, w http.ResponseWriter, r *http.Request)

// ServeHTTP enables ContextHandlerFunc to satisfy the ContextHandler interface
func (f ContextHandlerFunc) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) {
	f(c, w, r)
}

// ContextWrapper is a convenience wrapper for http.Handler into ContextHandler.
// Do not use this if you have long-running routines.
func ContextWrapper(h http.Handler) ContextHandler {
	return ContextHandlerFunc(func(_ context.Context, w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

// ContextRouter is an http router integrating a context.
type ContextRouter struct {
	router     *httprouter.Router
	context    context.Context
	cancelfunc context.CancelFunc
	sync.RWMutex
}

// New initializes and returns a new router.
func New() *ContextRouter {
	ctx, cfunc := context.WithCancel(context.Background())
	// corresponding cancelfunc will ne created on Start()
	return &ContextRouter{
		context:    ctx,
		cancelfunc: cfunc,
		router:     httprouter.New(),
	}
}

// Handle registers a ContextHandler for the required method and route.
// See https://github.com/julienschmidt/httprouter for details on named parameters.
func (s *ContextRouter) Handle(method Method, path string, handler ContextHandler) {
	s.Lock()
	s.router.Handle(string(method), path, s.wrapToHandle(handler))
	s.Unlock()
}

// HandleFunc registers a ContextHandlerFunc for the required method and route.
// See https://github.com/julienschmidt/httprouter for details on named parameters.
func (s *ContextRouter) HandleFunc(method Method, path string, handler func(context.Context, http.ResponseWriter, *http.Request)) {
	s.Handle(method, path, ContextHandlerFunc(handler))
}

// wrapToHandle wraps ContextHandlers to the httprouter.Handle type using a
// function closure which passes httprouter.Params as Context.Values to the
// registered ContextHandlers
func (s *ContextRouter) wrapToHandle(handler ContextHandler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// readlock to prevent the context from being changed by when a request is in flight
		s.RLock()
		c := s.context
		for _, p := range params {
			c = context.WithValue(c, p.Key, p.Value)
		}
		handler.ServeHTTP(c, w, req)
		s.RUnlock()
	})
}

// ServeHTTP routes requests to the appropriate handlers
func (s *ContextRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Stop closes the done channel of the root Context of the server to signal to
// any long running handlers to stop their work.
func (s *ContextRouter) Stop() {
	if s.router != nil {
		s.cancelfunc()
	}
	s.Lock()
	s.context, s.cancelfunc = context.WithCancel(context.Background())
	s.Unlock()
}
