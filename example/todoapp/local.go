// +build ignore

package main

import (
	"log"

	"github.com/srinathh/mobilehtml5app/example/todoapp"
)

func main() {
	tmpdir := "/tmp"
	appurl, err := todoapp.Start(tmpdir)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	log.Printf("App URL is: %s\nPersist Dir is: %s\n", appurl, tmpdir)
	select {}
}
