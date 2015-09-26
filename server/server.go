// Package server provides an integrated http server with graceful shutdown
// and parameterized routing capabilities for serving webapps locally on mobile.
//
// Create an instance of the Server by calling Server.NewServer(), register your
// handlers using Server.Handle() or Server.HandleFunc() and start the
// server with Server.Start(). The server expects handlers to be of the
// ContextHandler type which are similar to http.Handler but additionally
// take a context.Context as the first parameter in ServeHTTP()
//
// Unlike long-running servers, mobile apps should expect to be closed by
// the operating system at any time. Call Server.Stop() to shut down the server
// gracefully. When Stop() is called, the server instance closes the Done()
// channel on the Context passed to signal ContextHandlers and blocks upto a timeOut
// duration for them to finish before shutting down the server. Android
// documentation suggests apps should not block the UI thread for more than
// 100 to 200 milliseconds. Handlers that might spawn long-running functions
// or computation should check for Done channel closure and abandon or finish
// work if closed. See https://blog.golang.org/context for an illustration. Server
// uses github.com/tylerb/graceful package for the shutdown functionality.
//
// The Context passed to ContextHandlers is also used to pass named routing
// parameters which can be retreived using Context.Value(). Server uses
// github.com/julienschmidt/httprouter as the integrated router. For details
// on named parameters, see http://godoc.org/github.com/julienschmidt/httprouter
//
// The caller can also pass arbitrary instance specific settings to Handlers as
// string key-value pairs in Server.Start(). These are added to the context of
// the active server instance and passed to all ContextHandlers avoiding the
// need for managing global state variables. Possible uses could include passing
// settings form the native portion of the mobile app such as primary user account
// or persistant writeable directory path.
package server

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/tylerb/graceful"
	"golang.org/x/net/context"
)

// maxVerify and verifyDelay repesent the number of attempts and delay between
// attempts that Server.Start() should use when trying to verify that the server
// is actually running
const maxVerify = 5
const verifyDelay = time.Millisecond * 30

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

// Server is an integrated http server with graceful shutdown and parameterized
// routing capabilities
type Server struct {
	router     *httprouter.Router
	context    context.Context
	cancelfunc context.CancelFunc
	server     *graceful.Server
	mtx        sync.RWMutex
}

// NewServer initializes and returns a new Server. Call Start() to start the
// server and Stop() to shut it down.
func NewServer() *Server {
	// here we initialize only the router. The graceful.Server, Context & its
	// corresponding cancelfunc will ne created on Start()
	rtr := httprouter.New()
	return &Server{
		context:    nil,
		cancelfunc: nil,
		router:     rtr,
		server:     nil,
	}
}

// Handle registers a ContextHandler for the required method and route.
// See https://github.com/julienschmidt/httprouter for details on named parameters.
func (s *Server) Handle(method Method, path string, handler ContextHandler) {
	s.mtx.Lock()
	s.router.Handle(string(method), path, s.wrapToHandle(handler))
	s.mtx.Unlock()
}

// HandleFunc registers a ContextHandlerFunc for the required method and route.
// See https://github.com/julienschmidt/httprouter for details on named parameters.
func (s *Server) HandleFunc(method Method, path string, handler func(context.Context, http.ResponseWriter, *http.Request)) {
	s.Handle(method, path, ContextHandlerFunc(handler))
}

// wrapToHandle wraps ContextHandlers to the httprouter.Handle type using a
// function closure which passes httprouter.Params as Context.Values to the
// registered ContextHandlers
func (s *Server) wrapToHandle(handler ContextHandler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// we set a readlock to prevent the root context from being changed by
		// AddContextValue while anyu handlers are running.
		s.mtx.RLock()
		c := s.context
		for _, p := range params {
			c = context.WithValue(c, p.Key, p.Value)
		}
		handler.ServeHTTP(c, w, req)
		s.mtx.RUnlock()
	})
}

// Start creates and starts a graceful HTTP server listening on the specified
// address and verifies it is running. It creates a new root context to
// use with the server instance and will also copy any settings passed as key/value
// pairs in ctxValues to the context. If a server is already running,
// Start will call Stop() first to close it. Start will return the root url
// of the server (without the trailing slash) if successfully started. This
// could be useful if you have requested for a system chosen port
func (s *Server) Start(addr string, ctxValues map[string]interface{}) (string, error) {
	if s.server != nil {
		s.Stop(time.Millisecond * 100)
	}
	s.context, s.cancelfunc = context.WithCancel(context.Background())
	for k, v := range ctxValues {
		s.context = context.WithValue(s.context, k, v)
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("could not listen on %s: %s", addr, err)
	}

	s.server = &graceful.Server{
		Server: &http.Server{
			Addr:    l.Addr().String(),
			Handler: s.router},
		Timeout: time.Millisecond * 100,
	}
	go s.server.Serve(l)

	for j := 0; j < maxVerify; j++ {
		select {
		case <-time.After(verifyDelay):
			conn, err := net.Dial("tcp", l.Addr().String())
			if err != nil {
				continue
			}
			conn.Close()
			return "http://" + l.Addr().String(), nil
		}
	}
	s.Stop(time.Millisecond * 100)
	return "", fmt.Errorf("could not verify that the server is started")
}

// Stop closes the done channel of the root Context of the server to signal
// any open handlers to terminate and shuts down the server after
// waiting for upto the TimeOut period for any handlers to close. Stop blocks
// until the server closes
func (s *Server) Stop(timeOut time.Duration) {
	if s.server != nil {
		s.cancelfunc()
		s.server.Stop(timeOut)
		select {
		case <-s.server.StopChan():
		}
		s.server = nil
		s.cancelfunc = nil
		s.context = nil
	}
}
