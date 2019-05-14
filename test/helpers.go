package test

import (
  "fmt"
  "net/http"
  "app/jsonHttpHandler"
)

func getBasicJsonHttpHandler() *jsonHttpHandler.JsonHttpHandler {
  return jsonHttpHandler.New(
    &Globals{},
    []jsonHttpHandler.RouteData{
      jsonHttpHandler.NewRouteData(
        http.MethodGet,
        "/resources",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, getIndexBody)
          }
        },
        []jsonHttpHandler.Middleware{},
      ),
      jsonHttpHandler.NewRouteData(
        http.MethodPost,
        "/resources",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusCreated)
            fmt.Fprintf(w, postCreateBody)
          }
        },
        []jsonHttpHandler.Middleware{},
      ),
      jsonHttpHandler.NewRouteData(
        http.MethodGet,
        "/resources/:id",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            id := jsonHttpHandler.GetPathParamP(r.Context(), "id")

            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, fmt.Sprintf("%s %s", getDetailBody, id))
          }
        },
        []jsonHttpHandler.Middleware{},
      ),
      jsonHttpHandler.NewRouteData(
        http.MethodPatch,
        "/resources/:id",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            id := jsonHttpHandler.GetPathParamP(r.Context(), "id")

            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, fmt.Sprintf("%s %s", patchUpdateBody, id))
          }
        },
        []jsonHttpHandler.Middleware{},
      ),
      jsonHttpHandler.NewRouteData(
        http.MethodDelete,
        "/resources/:id",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            id := jsonHttpHandler.GetPathParamP(r.Context(), "id")

            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, fmt.Sprintf("%s %s", deleteDestroyBody, id))
          }
        },
        []jsonHttpHandler.Middleware{},
      ),
      jsonHttpHandler.NewRouteData(
        http.MethodGet,
        "/error",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            panic("The impossible has happened")
          }
        },
        []jsonHttpHandler.Middleware{},
      ),
      jsonHttpHandler.NewRouteData(
        http.MethodPost,
        "/middlewares",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, middlewaresBody)
          }
        },
        []jsonHttpHandler.Middleware{
          requirePayload,
          requireExactPayload(wantedPayload),
        },
      ),
    },
  )
}

func getJsonHttpHandlerWithCors(corsHandler jsonHttpHandler.GlobalsReceivingHandlerFunc) *jsonHttpHandler.JsonHttpHandler {
  return jsonHttpHandler.NewWithCors(
    &Globals{},
    corsHandler,
    []jsonHttpHandler.RouteData{
      jsonHttpHandler.NewRouteData(
        http.MethodGet,
        "/resources",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, getIndexBody)
          }
        },
        []jsonHttpHandler.Middleware{},
      ),
      jsonHttpHandler.NewRouteData(
        http.MethodPost,
        "/resources",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusCreated)
            fmt.Fprintf(w, postCreateBody)
          }
        },
        []jsonHttpHandler.Middleware{},
      ),
      jsonHttpHandler.NewRouteData(
        http.MethodGet,
        "/resources/:id",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            id := jsonHttpHandler.GetPathParamP(r.Context(), "id")

            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, fmt.Sprintf("%s %s", getDetailBody, id))
          }
        },
        []jsonHttpHandler.Middleware{},
      ),
      jsonHttpHandler.NewRouteData(
        http.MethodPatch,
        "/resources/:id",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            id := jsonHttpHandler.GetPathParamP(r.Context(), "id")

            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, fmt.Sprintf("%s %s", patchUpdateBody, id))
          }
        },
        []jsonHttpHandler.Middleware{},
      ),
      jsonHttpHandler.NewRouteData(
        http.MethodDelete,
        "/resources/:id",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            id := jsonHttpHandler.GetPathParamP(r.Context(), "id")

            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, fmt.Sprintf("%s %s", deleteDestroyBody, id))
          }
        },
        []jsonHttpHandler.Middleware{},
      ),
      jsonHttpHandler.NewRouteData(
        http.MethodGet,
        "/error",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            panic("The impossible has happened")
          }
        },
        []jsonHttpHandler.Middleware{},
      ),
      jsonHttpHandler.NewRouteData(
        http.MethodPost,
        "/middlewares",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, middlewaresBody)
          }
        },
        []jsonHttpHandler.Middleware{
          requirePayload,
          requireExactPayload(wantedPayload),
        },
      ),
    },
  )
}
