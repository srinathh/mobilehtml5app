package todoapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"
)

const (
	maxTextSize = 65536
	maxStatus   = 1
	maxPriority = 1
)

type item struct {
	Text     string //Title of the item
	Status   int
	Priority int
}

const (
	statusOpen   = 0
	statusClosed = 1
)

const (
	priorityNormal    int = 0
	priorityImportant int = 1
)

/*
srv.HandleFunc(server.GET, "/items", logger(serveItems))
srv.HandleFunc(server.PUT, "/items/:itemid", logger(putEditItem))
*/

func serveItems(c context.Context, w http.ResponseWriter, r *http.Request) {
	items, err := bk.fetchAll()
	if err != nil {
		log.Printf("serveItems: error fetching items: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("serveItems: error encoding items: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func checkInputs(i item) error {
	// impose a hard limit on text size
	if l := len(i.Text); l > maxTextSize {
		return fmt.Errorf("checkInputs: text is too long: %d", l)
	}

	if i.Status < 0 || i.Status > maxStatus {
		return fmt.Errorf("checkInputs: bad status: %d", i.Status)
	}

	if i.Priority < 0 || i.Priority > maxPriority {
		return fmt.Errorf("checkInputs: bad priority: %d", i.Priority)
	}
	return nil
}

func putEditItem(c context.Context, w http.ResponseWriter, r *http.Request) {
	itemid := c.Value("itemid").(string)
	var i item
	if err := json.Unmarshal([]byte(r.PostFormValue("data")), &i); err != nil {
		log.Printf("putEditItems: error decoding item: %s\n", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	switch itemid {
	case "new":
		itemid = strconv.FormatInt(time.Now().UnixNano(), 10)
	default:
		if !bk.exists(itemid) {
			log.Printf("putEditItems: error looking up itemid: %s\n", itemid)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}

	if err := checkInputs(i); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := bk.edit(itemid, i); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}
