package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from GET method!")
    })

    http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            fmt.Fprintf(w, "Data received via POST!")
        } else {
            fmt.Fprintf(w, "Please use POST method.")
        }
    })

    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}
