package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "techwave/routes"
)

func main() {
    r := mux.NewRouter()
    routes.RegisterRoutes(r)
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}