package todoapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/srinathh/mobilehtml5app/server"
)

var srv *server.Server

const (
	persistdir = "persistdir"
)

func logger(h server.ContextHandlerFunc) server.ContextHandlerFunc {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		log.Printf(r.URL.String())
		h.ServeHTTP(c, w, r)
	}
}

// Start is called by the native portion of the webapp to start the web server.
// It returns the server root URL (without the trailing slash) and any errors.
func Start(settings string) (string, error) {

	//initialization
	settingsMap := make(map[string]string)
	if err := json.NewDecoder(bytes.NewBufferString(settings)).Decode(&settingsMap); err != nil {
		return "", fmt.Errorf("could not decode settings: %s", err)
	}
	srv := server.NewServer()

	//setting up our handlers
	srv.HandleFunc(server.GET, "/", logger(serveIndex))
	srv.HandleFunc(server.GET, "/items", logger(fetchAll))
	srv.HandleFunc(server.POST, "/items/new", logger(createItem))
	srv.HandleFunc(server.GET, "/items/:itemid", logger(deleteItem))
	srv.HandleFunc(server.GET, "/res/*respath", logger(serveRes))
	//starting the server

	backend := newMapBackend()
	backend.createSample()
	bk = backend

	return srv.Start("127.0.0.1:0", settingsMap)
}

type backend interface {
	fetchAll() ([]item, error)
	create(item) error
	delete(id string) error
}

var bk backend

// Stop is called by the native portion of the webapp to stop the web server.
func Stop() {
	srv.Stop(time.Millisecond * 100)
}
