package todoapp

// command to build static assets use the command -
//  go-bindata -nocompress -prefix res res/...

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/srinathh/mobilehtml5app/example/todoapp/data"
	"golang.org/x/net/context"
)

func serveIndex(_ context.Context, w http.ResponseWriter, r *http.Request) {
	serveAsset("app/index.html", w, r)
}

func serveRes(c context.Context, w http.ResponseWriter, r *http.Request) {
	respath := c.Value("respath").(string)
	if respath[0] == '/' {
		respath = respath[1:]
	}
	serveAsset(respath, w, r)
}

func serveAsset(fpath string, w http.ResponseWriter, r *http.Request) {
	b, err := data.Asset(fpath)
	if err != nil {
		log.Printf("serveAsset: could not load asset data %s: %s", fpath, err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	finfo, err := data.AssetInfo(fpath)
	if err != nil {
		log.Printf("serveAsset: could not load asset fileinfo %s: %s", fpath, err)
		http.Error(w, "", http.StatusNotFound)
		return
	}
	http.ServeContent(w, r, finfo.Name(), finfo.ModTime(), bytes.NewReader(b))
}

func fitCropScale(i image.Image, r image.Rectangle) image.Image {
	wantRatio := float64(r.Bounds().Dx()) / float64(r.Bounds().Dy())
	haveRatio := float64(i.Bounds().Dx()) / float64(i.Bounds().Dy())

	sliceRect := image.Rectangle{}

	if haveRatio > wantRatio {
		wantwidth := wantRatio * float64(i.Bounds().Dy())
		sliceRect = image.Rect(i.Bounds().Dx()/2-int(wantwidth/2), 0, i.Bounds().Dx()/2+int(wantwidth/2), i.Bounds().Dy())
	} else {
		wantheight := float64(i.Bounds().Dx()) / wantRatio
		sliceRect = image.Rect(0, i.Bounds().Dy()/2-int(wantheight/2), i.Bounds().Dx(), i.Bounds().Dy()/2+int(wantheight/2))
	}

	return imaging.Resize(imaging.Crop(i, sliceRect), r.Dx(), r.Dy(), imaging.Lanczos)
}

func (a *App) serveBg(c context.Context, w http.ResponseWriter, r *http.Request) {
	bg := a.bg

	width, err := strconv.Atoi(c.Value("width").(string))
	if err != nil {
		log.Printf("serveBG: error decoding width %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	height, err := strconv.Atoi(c.Value("height").(string))
	if err != nil {
		log.Printf("serveBG: error decoding height %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	buf := bytes.Buffer{}

	if err := jpeg.Encode(&buf, fitCropScale(bg, image.Rect(0, 0, width, height)), &jpeg.Options{Quality: 80}); err != nil {
		log.Printf("serveBG: error encoding background :%s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, "bg.jpg", time.Now(), bytes.NewReader(buf.Bytes()))
}
