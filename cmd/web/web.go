package main

import (
	"log"
	"net/http"

	route "github.com/iineva/templates/cmd/web/route"
	"github.com/iineva/templates/templates"
)

func main() {
	log.Println("SERVER STARTING...")

	host := "0.0.0.0:8080"

	// parser API
	http.HandleFunc("/v1/", route.HandlerTemplate)
	// static files
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(templates.FS))))

	log.Printf("SERVER LISTEN ON: %v", host)
	log.Fatal(http.ListenAndServe(host, nil))
}
