package main

import (
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
    "github.com/go-redis/redis/v8"
    "techwave/routes"
)

func main() {
    // Initialize Redis client
    // Use environment variable REDIS_ADDR or default to localhost:6379
    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "localhost:6379"
    }
    
    redisClient := redis.NewClient(&redis.Options{
        Addr: redisAddr,
        Password: os.Getenv("REDIS_PASSWORD"), // Optional password
        DB: 0, // Default DB
    })
    
    r := mux.NewRouter()
    routes.RegisterRoutes(r, redisClient)
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}