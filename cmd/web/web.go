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

	host :=  fmt.Sprintf("%s:%s", env.Get("ADDRESS", "0.0.0.0"), env.Get("PORT", "8080"))

	// parser API
	http.HandleFunc("/v1/", route.HandlerTemplate)
	// static files
	http.Handle("/", route.Healthy("/", http.FileServer(http.FS(templates.FS))))

	log.Printf("SERVER LISTEN ON: %v", host)
	log.Fatal(http.ListenAndServe(host, nil))
}
