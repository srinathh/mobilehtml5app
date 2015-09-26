package todoapp

// command to build static assets use the command -
//  go-bindata -nocompress -prefix res res/...

import (
	"bytes"
	"log"
	"net/http"

	"golang.org/x/net/context"
)

func serveIndex(_ context.Context, w http.ResponseWriter, r *http.Request) {
	serveAsset("app/index.html", w, r)
}

func serveComponents(_ context.Context, w http.ResponseWriter, r *http.Request) {
	serveAsset("app/components.js", w, r)
}

func serveRes(c context.Context, w http.ResponseWriter, r *http.Request) {
	respath := c.Value("respath").(string)
	if respath[0] == '/' {
		respath = respath[1:]
	}
	serveAsset(respath, w, r)
}

func serveAsset(fpath string, w http.ResponseWriter, r *http.Request) {
	data, err := Asset(fpath)
	if err != nil {
		log.Printf("serveAsset: could not load asset data %s: %s", fpath, err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	finfo, err := AssetInfo(fpath)
	if err != nil {
		log.Printf("serveAsset: could not load asset fileinfo %s: %s", fpath, err)
		http.Error(w, "", http.StatusNotFound)
		return
	}
	http.ServeContent(w, r, finfo.Name(), finfo.ModTime(), bytes.NewReader(data))
}
