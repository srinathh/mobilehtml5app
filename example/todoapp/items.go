package todoapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"golang.org/x/net/context"
)

const (
	maxTextSize = 65536
	maxPriority = 1
)

type item struct {
	ID       string
	time     time.Time
	Text     string //Title of the item
	Priority int
}

func (i item) check() error {
	if l := len(i.Text); l > maxTextSize {
		return fmt.Errorf("checkInputs: text is too long: %d", l)
	}

	if i.Priority < 0 || i.Priority > maxPriority {
		return fmt.Errorf("checkInputs: bad priority: %d", i.Priority)
	}
	return nil
}

const timestamp = "2006-01-02T15:04:05.000Z"

type itemSorter []item

func (s itemSorter) Len() int { return len(s) }

func (s itemSorter) Less(i, j int) bool {
	if s[i].Priority > s[j].Priority {
		return true
	} else if s[i].Priority < s[j].Priority {
		return false
	}

	if s[i].time.Before(s[j].time) {
		return true
	}
	return false
}

func (s itemSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

const (
	priorityNormal    int = 0
	priorityImportant int = 1
)

func fetchAll(c context.Context, w http.ResponseWriter, r *http.Request) {
	items, err := bk.fetchAll()
	sort.Sort(itemSorter(items))

	if err != nil {
		log.Printf("fetchAll: error fetching items: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("fetchAll: error encoding items: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func createItem(c context.Context, w http.ResponseWriter, r *http.Request) {
	var i item
	if err := json.Unmarshal([]byte(r.PostFormValue("data")), &i); err != nil {
		log.Printf("createItem: error decoding item: %s\n", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := i.check(); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	i.time = time.Now()
	i.ID = i.time.Format(timestamp)

	if err := bk.create(i); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func deleteItem(c context.Context, w http.ResponseWriter, r *http.Request) {
	id := c.Value("itemid").(string)
	if err := bk.delete(id); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}
