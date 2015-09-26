package todoapp

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/srinathh/mobilehtml5app/server"
)

var srv *server.Server
var bk backend

const (
	persistdir = "persistdir"
	backendptr = "backendptr"
)

func logger(h server.ContextHandlerFunc) server.ContextHandlerFunc {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		log.Printf(r.URL.String())
		h.ServeHTTP(c, w, r)
	}
}

// Start is called by the native portion of the webapp to start the web server.
// It returns the server root URL (without the trailing slash) and any errors.
func Start(pdir string) (string, error) {

	srv := server.NewServer()

	//setting up our handlers
	srv.HandleFunc(server.GET, "/", logger(serveIndex))
	srv.HandleFunc(server.GET, "/components.js", logger(serveComponents))
	srv.HandleFunc(server.GET, "/items", logger(fetchAll))
	srv.HandleFunc(server.POST, "/items/new", logger(createItem))
	srv.HandleFunc(server.GET, "/items/:itemid", logger(deleteItem))
	srv.HandleFunc(server.GET, "/res/*respath", logger(serveRes))
	//starting the server

	bk, err := newBoltBackend(pdir)
	if err != nil {
		return "", err
	}

	return srv.Start("127.0.0.1:0", map[string]interface{}{backendptr: bk})
}

type backend interface {
	FetchAll() ([]item, error)
	Create(item) error
	Delete(id string) error
	Stop()
}

// Stop is called by the native portion of the webapp to stop the web server.
func Stop() {
	bk.Stop()
	srv.Stop(time.Millisecond * 100)
}
