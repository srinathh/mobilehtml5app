// Package server provides an integrated http server with graceful shutdown
// and parameterized routing capabilities for serving webapps locally on mobile.
//
// Create an instance of the Server by calling Server.NewServer(), register your
// handlers using Server.Router.Handle() or Server.RouterHandleFunc() and start the
// server with Server.Start(). The server expects handlers to be of the
// contextrouter.ContextHandler type which are similar to http.Handler but additionally
// take a context.Context as the first parameter in ServeHTTP()
//
// Unlike long-running servers, mobile apps should expect to be closed by
// the operating system at any time. Call Server.Stop() to shut down the server
// gracefully. When Stop() is called, the server instance closes the Done()
// channel on the router Context  to signal Handlers and blocks upto a timeOut
// duration for them to finish before shutting down the server. Android
// documentation suggests apps should not block the UI thread for more than
// 100 to 200 milliseconds. Handlers that might spawn long-running functions
// or computation should check for Done channel closure and abandon or finish
// work if closed. See https://blog.golang.org/context for an illustration. Server
// uses github.com/tylerb/graceful package for the shutdown functionality.
package server

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/srinathh/mobilehtml5app/contextrouter"
	"github.com/tylerb/graceful"
)

// maxVerify and verifyDelay repesent the number of attempts and delay between
// attempts that Server.Start() should use when trying to verify that the server
// is actually running
const maxVerify = 5
const verifyDelay = time.Millisecond * 30

// Server is an integrated http server with graceful shutdown and parameterized
// routing capabilities
type Server struct {
	Router *contextrouter.ContextRouter
	server *graceful.Server
	sync.RWMutex
}

// NewServer initializes and returns a new Server. Call Start() to start the
// server and Stop() to shut it down.
func NewServer() *Server {
	return &Server{
		Router: contextrouter.New(),
		server: nil,
	}
}

// Start creates and starts a graceful HTTP server listening on the specified
// address and verifies it is running. It creates a new root context to
// use with the server instance and will also copy any settings passed as key/value
// pairs in ctxValues to the context. If a server is already running,
// Start will call Stop() first to close it. Start will return the root url
// of the server (without the trailing slash) if successfully started. This
// could be useful if you have requested for a system chosen port
func (s *Server) Start(addr string) (string, error) {
	if s.server != nil {
		s.Stop(time.Millisecond * 100)
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("could not listen on %s: %s", addr, err)
	}
	s.server = &graceful.Server{
		Server: &http.Server{
			Addr:    l.Addr().String(),
			Handler: s.Router},
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
		s.Router.Stop()
		s.server.Stop(timeOut)
		select {
		case <-s.server.StopChan():
		}
		s.server = nil
	}
}
