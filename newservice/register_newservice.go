package main

import "net/http"

func init() {
    http.HandleFunc("/new/get", GetHello)
    http.HandleFunc("/new/post", PostEcho)
}
