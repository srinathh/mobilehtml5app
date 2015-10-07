// +build ignore

package main

import (
	"log"

	"github.com/srinathh/mobilehtml5app/example/todoapp"
)

func main() {
	tmpdir := "/tmp"
	app, err := todoapp.NewApp(tmpdir)
	if err != nil {
		log.Fatal("Error creating app: %s", err)
	}

	appurl, err := app.Start()
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	log.Printf("App URL is: %s\nPersist Dir is: %s\n", appurl, tmpdir)
	select {}
}
