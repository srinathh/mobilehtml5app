package todoapp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"golang.org/x/net/context"
)

var wants = map[string]item{
	"1442717418486489952": item{
		Text:     "Test Item 1",
		Status:   0,
		Priority: 0,
	},
	"1442717451419839383": item{
		Text:     "Test Item 2",
		Status:   1,
		Priority: 0,
	},
	"1442717480424770931": item{
		Text:     "Test Item 3",
		Status:   0,
		Priority: 1,
	},
	"1442717575928777740": item{
		Text:     "Test Item 4",
		Status:   1,
		Priority: 1,
	},
}

func TestFetchItems(t *testing.T) {
	bk = &mapBackend{
		Items: wants,
	}
	w := httptest.NewRecorder()
	serveItems(context.Background(), w, nil)
	got := map[string]item{}
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if len(got) != len(wants) {
		t.Errorf("length mismatch\ngot: %v\nwant: %v\n", got, wants)
	}

	for k, v := range got {
		if v != wants[k] {
			t.Errorf("item mismatch\ngot: %v\nwant: %v\n", v, wants[k])
		}
	}

}

func TestEditItem(t *testing.T) {
	backend := newMapBackend()
	bk = backend
	w := httptest.NewRecorder()
	want := item{
		Text:     "Test Item 1",
		Status:   0,
		Priority: 0,
	}
	b, _ := json.Marshal(want)
	form := url.Values{}
	form.Add("data", string(b))
	r, err := http.NewRequest("POST", "127.0.0.1:9898/items/new", nil)
	r.PostForm = form
	putEditItem(context.WithValue(context.Background(), "itemid", "new"), w, r)
	if w.Code != http.StatusOK {
		t.Fatal(err)
	}
	if len(backend.Items) != 1 {
		t.Fatal("Unable to set item")
	}

	for _, v := range backend.Items {
		if v != want {
			t.Error("mismatch")
		}
	}
}
