package todoapp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"golang.org/x/net/context"
)

func TestFetchItems(t *testing.T) {
	backend := newMapBackend()
	backend.createSample()
	bk = backend
	w := httptest.NewRecorder()
	fetchAll(context.Background(), w, nil)

	got := []item{}
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if len(got) != len(backend.Items) {
		t.Errorf("length mismatch\ngot: %v\nwant: %v\n", got, backend.Items)
	}

	for _, item := range got {
		found, ok := backend.Items[item.ID]
		if !ok || item.Text != found.Text || item.Priority != found.Priority {
			t.Errorf("item mismatch\ngot: %v\nwant: %v\n", item, backend.Items[item.ID])
		}
	}
}

func TestDeleteItem(t *testing.T) {
	backend := newMapBackend()
	backend.createSample()
	bk = backend

	keys := []string{}
	for k, _ := range backend.Items {
		keys = append(keys, k)
	}

	w := httptest.NewRecorder()
	for _, k := range keys {
		c := context.WithValue(context.Background(), "itemid", k)
		deleteItem(c, w, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Did not get OK: %v: %v", w.Code, k)
		}
	}

	if len(backend.Items) != 0 {
		t.Fatalf("error deleting all entries")
	}
}

func TestCreateItem(t *testing.T) {
	backend := newMapBackend()
	bk = backend
	w := httptest.NewRecorder()
	want := item{
		ID:       "new",
		Text:     "Test Item 1",
		Priority: 0,
	}
	b, _ := json.Marshal(want)
	form := url.Values{}
	form.Add("data", string(b))
	r, err := http.NewRequest("POST", "127.0.0.1:9898/items/new", nil)
	r.PostForm = form

	createItem(context.Background(), w, r)
	if w.Code != http.StatusOK {
		t.Fatal(err)
	}
	if len(backend.Items) != 1 {
		t.Fatal("Unable to set item")
	}

	for _, item := range backend.Items {
		if item.Priority != want.Priority || item.Text != want.Text {
			t.Errorf("item mismatch\ngot: %v\nwant: %v\n", item, want)
		}
	}
}
