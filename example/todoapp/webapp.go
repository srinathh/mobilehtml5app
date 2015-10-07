package todoapp

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"time"

	"github.com/srinathh/mobilehtml5app/contextrouter"
	"github.com/srinathh/mobilehtml5app/example/todoapp/data"
	"github.com/srinathh/mobilehtml5app/server"
	"golang.org/x/net/context"
)

// App implements a web server backend for an android app
type App struct {
	srv *server.Server
	bk  backend
	bg  image.Image
}

// NewApp returns an App
func NewApp(pdir string) (*App, error) {
	srv := server.NewServer()
	bk, err := newBoltBackend(pdir)
	if err != nil {
		return nil, err
	}
	bg, err := loadBG()
	if err != nil {
		return nil, err
	}

	app := &App{
		srv: srv,
		bk:  bk,
		bg:  bg,
	}

	srv.Router.HandleFunc(contextrouter.GET, "/", logger(serveIndex))
	srv.Router.HandleFunc(contextrouter.GET, "/items", logger(app.fetchAll))
	srv.Router.HandleFunc(contextrouter.POST, "/items/new", logger(app.createItem))
	srv.Router.HandleFunc(contextrouter.GET, "/items/:itemid", logger(app.deleteItem))
	srv.Router.HandleFunc(contextrouter.GET, "/res/*respath", logger(serveRes))
	srv.Router.HandleFunc(contextrouter.GET, "/bg/:width/:height", logger(app.serveBg))

	return app, nil
}

// Start is called by the native portion of the webapp to start the web server.
// It returns the server root URL (without the trailing slash) and any errors.
func (app *App) Start() (string, error) {
	return app.srv.Start("127.0.0.1:0")
}

// Stop is called by the native portion of the webapp to stop the web server.
func (app *App) Stop() {
	app.srv.Stop(time.Millisecond * 100)
}

func logger(h contextrouter.ContextHandlerFunc) contextrouter.ContextHandlerFunc {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		log.Printf(r.URL.String())
		h.ServeHTTP(c, w, r)
	}
}

type backend interface {
	fetchAll() ([]item, error)
	create(item) error
	delete(id string) error
	stop()
}

func loadBG() (image.Image, error) {
	byt, err := data.Asset("img/bg.jpg")
	if err != nil {
		return nil, fmt.Errorf("error loading background image: %s", err)
	}
	return jpeg.Decode(bytes.NewReader(byt))
}
