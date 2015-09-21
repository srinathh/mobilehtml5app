// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/srinathh/mobilehtml5app/example/todoapp"
)

func main() {
	tmpdir, err := ioutil.TempDir("", "todoapp")
	if err != nil {
		log.Fatalf("Could not create a temporary directory: %s", err)
	}

	settings := map[string]string{"persistdir": tmpdir}
	buf := bytes.Buffer{}
	json.NewEncoder(&buf).Encode(settings)

	appurl, err := todoapp.Start(buf.String())
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	log.Printf("App URL is: %s\nPersist Dir is: %s\n", appurl, tmpdir)
	select {}
}
