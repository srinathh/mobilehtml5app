package todoapp

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/srinathh/mobilehtml5app/example/todoapp/data"
	"github.com/srinathh/mobilehtml5app/server"
)

var srv *server.Server
var bk backend

const (
	persistdir = "persistdir"
	backendptr = "backendptr"
	bgimg      = "bgimg"
)

func logger(h server.ContextHandlerFunc) server.ContextHandlerFunc {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		log.Printf(r.URL.String())
		h.ServeHTTP(c, w, r)
	}
}

func loadBG() (image.Image, error) {
	byt, err := data.Asset("img/bg.jpg")
	if err != nil {
		return nil, fmt.Errorf("error loading background image: %s", err)
	}
	return jpeg.Decode(bytes.NewReader(byt))
}

// Start is called by the native portion of the webapp to start the web server.
// It returns the server root URL (without the trailing slash) and any errors.
func Start(pdir string) (string, error) {

	srv = server.NewServer()

	//setting up our handlers
	srv.HandleFunc(server.GET, "/", logger(serveIndex))
	srv.HandleFunc(server.GET, "/components.js", logger(serveComponents))
	srv.HandleFunc(server.GET, "/items", logger(fetchAll))
	srv.HandleFunc(server.POST, "/items/new", logger(createItem))
	srv.HandleFunc(server.GET, "/items/:itemid", logger(deleteItem))
	srv.HandleFunc(server.GET, "/res/*respath", logger(serveRes))
	srv.HandleFunc(server.GET, "/bg/:width/:height", logger(serveBg))
	//starting the server

	var err error

	bk, err = newBoltBackend(pdir)
	if err != nil {
		return "", err
	}

	var bg image.Image

	bg, err = loadBG()
	if err != nil {
		return "", err
	}

	return srv.Start("127.0.0.1:0", map[string]interface{}{backendptr: bk, bgimg: bg})
}

type backend interface {
	fetchAll() ([]item, error)
	create(item) error
	delete(id string) error
	stop()
}

// Stop is called by the native portion of the webapp to stop the web server.
func Stop() {
	log.Println("calling backend.stop")
	bk.stop()
	log.Println("calling server.Stop")
	srv.Stop(time.Millisecond * 100)
	log.Println("finishing with Stop")
}
