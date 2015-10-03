package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/srinathh/mobilehtml5app/contextrouter"

	"golang.org/x/net/context"
)

func initServer() *Server {
	srv := NewServer()

	srv.Router.HandleFunc(contextrouter.GET, "/:hellostring/:name", func(c context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s, %s", c.Value("hellostring").(string), c.Value("name").(string))
	})

	srv.Router.Handle(contextrouter.GET, "/", contextrouter.ContextWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, Stranger")
	})))
	return srv
}

func checkResponse(url, want string) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not fetch from %s: %s", url+"/Alice", err)
	}
	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("could not read response body from %s: %s", url+"/Alice", err)
	}
	res.Body.Close()
	if string(got) != want {
		return fmt.Errorf("want: %s got: %s", want, got)
	}
	return nil
}

func TestBasic(t *testing.T) {
	srv := initServer()

	rooturl, err := srv.Start("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	if err := checkResponse(rooturl+"/Namaste/Alice", "Namaste, Alice"); err != nil {
		t.Error(err)
	}
	if err := checkResponse(rooturl, "Hello, Stranger"); err != nil {
		t.Error(err)
	}
	srv.Stop(time.Millisecond * 100)
}

func TestStartStop(t *testing.T) {
	srv := initServer()

	for j := 0; j < 5; j++ {
		//start the server
		rooturl, err := srv.Start("127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		if err := checkResponse(rooturl+"/Namaste/Alice", "Namaste, Alice"); err != nil {
			t.Error(err)
		}
		srv.Stop(time.Millisecond * 100)
		if err := checkResponse(rooturl+"/Namaste/Alice", "Namaste, Alice"); err == nil {
			t.Errorf("got a valid response after server close")
		}
	}
}

func TestIllegalStart(t *testing.T) {
	srv := initServer()
	// in most setups, this attempt to connect to an arbitrary low number port should fail
	// and return server unable to start type of error
	if _, err := srv.Start("127.0.0.1:1"); err == nil {
		t.Errorf("started on a low numbered port!!!")
		srv.Stop(time.Millisecond * 100)
	}
}

func TestUnStoppedStart(t *testing.T) {
	srv := initServer()
	_, err := srv.Start("127.0.0.1:9999")
	if err != nil {
		t.Fatal(err)
	}
	_, err = srv.Start("127.0.0.1:9999")
	if err != nil {
		t.Fatal(err)
	}
	srv.Stop(time.Millisecond * 100)
}
