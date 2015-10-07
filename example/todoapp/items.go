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
	maxPriority       = 1
	priorityNormal    = 0
	priorityImportant = 1
)

type item struct {
	ID       string
	time     time.Time
	Text     string //Title of the item
	Priority int
}

func (i item) check() error {
	if i.time == time.Unix(0, 0) {
		return fmt.Errorf("zero time not allowed")
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
	//higher priority item takes precedence
	if s[i].Priority > s[j].Priority {
		return true
	} else if s[i].Priority < s[j].Priority {
		return false
	}

	//earlier item takes precedence
	if s[i].time.Before(s[j].time) {
		return true
	}
	return false
}

func (s itemSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (a *App) fetchAll(c context.Context, w http.ResponseWriter, r *http.Request) {
	items, err := a.bk.fetchAll()
	if err != nil {
		log.Printf("fetchAll: error fetching items: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	sort.Sort(itemSorter(items))

	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("fetchAll: error encoding items: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (a *App) createItem(c context.Context, w http.ResponseWriter, r *http.Request) {
	var i item
	if err := json.Unmarshal([]byte(r.PostFormValue("data")), &i); err != nil {
		log.Printf("createItem: error decoding item: %s\n", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	i.time = time.Now()
	i.ID = i.time.Format(timestamp)

	if err := i.check(); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := a.bk.create(i); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func (a *App) deleteItem(c context.Context, w http.ResponseWriter, r *http.Request) {
	id := c.Value("itemid").(string)
	if err := a.bk.delete(id); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}
