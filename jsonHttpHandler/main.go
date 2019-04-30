package jsonHttpHandler

import (
  "fmt"
  "context"
  "strings"
  "net/http"
  "github.com/renra/go-errtrace/errtrace"
)

type Config interface {
  Get(string) (interface{}, *errtrace.Error)
  GetP(string) interface{}
  GetString(string) (string, *errtrace.Error)
  GetStringP(string) string
  GetInt(string) (int, *errtrace.Error)
  GetIntP(string) int
  GetFloat(string) (float64, *errtrace.Error)
  GetFloatP(string) float64
  GetBool(string) (bool, *errtrace.Error)
  GetBoolP(string) bool
}

type Logger interface {
  LogWithSeverity(map[string]string, int)
}

type Globals interface {
  Config() Config
  Logger() Logger
  Log(string)
  LogErrorWithTrace(*errtrace.Error)
  Clients() map[string]interface{}
}

type GlobalsReceivingHandlerFunc func(Globals) http.HandlerFunc

type Middleware func(http.HandlerFunc) http.HandlerFunc

type RouteData struct {
  Verb string
  Middlewares []Middleware
  Handler GlobalsReceivingHandlerFunc
  Pattern string
}

func NewRouteData(verb string, pattern string, handler GlobalsReceivingHandlerFunc, middlewares []Middleware) RouteData {
  return RouteData{
    Verb: verb,
    Pattern: pattern,
    Handler: handler,
    Middlewares: middlewares,
  }
}

type JsonHttpHandler struct {
  globals Globals
  routeMap []RouteData
}

func (h JsonHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  h.globals.Log(fmt.Sprintf("%s", r.URL))

  defer func() {
    if r := recover(); r != nil {
      h.globals.LogErrorWithTrace(errtrace.Wrap(r))

      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "{}")
    }
  }()

  for _, routeData := range h.routeMap {
    if routeData.Verb == r.Method {
      doesRouteMatch, pathParams := GetMatchAndPathParams(routeData.Pattern, r.URL.Path)

      if doesRouteMatch {
        ctx := r.Context()
        r = r.WithContext(context.WithValue(ctx, PathParamsKey, *pathParams))

        routeData.Handler(h.globals)(w, r)
        return
      }
    }
  }

  w.WriteHeader(http.StatusNotFound)
  fmt.Fprintf(w, "{}")
}

func GetMatchAndPathParams(pattern string, urlPath string) (bool, *map[string]string) {
  patternShards := strings.Split(pattern, "/")
  pathShards := strings.Split(urlPath, "/")

  if len(patternShards) != len(pathShards) {
    return false, nil
  }

  params := make(map[string]string, 0)

  for i, _ := range patternShards {
    shards := strings.Split(patternShards[i], ":")

    if len(shards) == 1 {
      if patternShards[i] != pathShards[i] {
        return false, nil
      }
    } else {
      params[shards[1]] = pathShards[i]
    }
  }

  return true, &params
}

func New(g Globals, routeMap []RouteData) *JsonHttpHandler {
  return &JsonHttpHandler{globals: g, routeMap: routeMap}
}
