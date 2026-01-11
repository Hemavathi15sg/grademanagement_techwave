package main

import (
	"log"
	"net/http"
	"techwave/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterRoutes(r)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
