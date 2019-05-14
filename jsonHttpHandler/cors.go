package jsonHttpHandler

import (
  "fmt"
  "net/http"
)

func corsNoop(g Globals) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "{}")
  }
}

