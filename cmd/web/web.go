package main

import (
	"fmt"
	"log"
	"net/http"

	route "github.com/iineva/templates/cmd/web/route"
	"github.com/iineva/templates/pkg/env"
	"github.com/iineva/templates/templates"
)

func main() {
	log.Println("SERVER STARTING...")

	host :=  fmt.Sprintf("%s:%s", env.Get("ADDRESS", "127.0.0.1"), env.Get("PORT", "8080"))

	// parser API
	http.HandleFunc("/v1/", route.HandlerTemplate)
	// healthy check
	http.HandleFunc("/",  func(w http.ResponseWriter, r *http.Request) {
		 fmt.Fprint(w, "")
	})
	// static files
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(templates.FS))))

	log.Printf("SERVER LISTEN ON: %v", host)
	log.Fatal(http.ListenAndServe(host, nil))
}
