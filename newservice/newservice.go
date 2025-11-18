package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// GET /new/get
func GetHello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from new API service!")
}

// POST /new/post
func PostEcho(w http.ResponseWriter, r *http.Request) {
    var body map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(body)
}
