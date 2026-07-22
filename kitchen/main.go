package main

import (
    "net/http"

    "kitchen/handler"
    "kitchen/store"
)

func main() {
    h := &handler.Handler{Store: store.NewMemoryStore()}
    http.ListenAndServe(":8080", handler.WithCORS(h.Routes()))
}