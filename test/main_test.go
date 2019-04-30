package test

import (
  "fmt"
  "net/http"
  "net/http/httptest"
  "testing"
  "app/jsonHttpHandler"
  "github.com/stretchr/testify/suite"
  "github.com/stretchr/testify/assert"
  "github.com/renra/go-errtrace/errtrace"
)

type Logger struct {
}

func (l *Logger) LogWithSeverity(data map[string]string, severity int) {
}

type Config struct {
}

func (c *Config) Get(key string) (interface {}, *errtrace.Error) {
  return key, nil
}

func (c *Config) GetString(key string) (string, *errtrace.Error) {
  return key, nil
}

func (c *Config) GetInt(key string) (int, *errtrace.Error) {
  return 0, nil
}

func (c *Config) GetFloat(key string) (float64, *errtrace.Error) {
  return 0.0, nil
}

func (c *Config) GetBool(key string) (bool, *errtrace.Error) {
  return true, nil
}

func (c *Config) GetP(key string) interface {} {
  return key
}

func (c *Config) GetStringP(key string) string {
  return key
}

func (c *Config) GetIntP(key string) int {
  return 0
}

func (c *Config) GetFloatP(key string) float64 {
  return 0.0
}

func (c *Config) GetBoolP(key string) bool {
  return true
}

type Globals struct {
}

func (g *Globals) Config() jsonHttpHandler.Config {
  return &Config{}
}

func (g *Globals) Logger() jsonHttpHandler.Logger {
  return &Logger{}
}

func (g *Globals) Log(msg string) {
  g.Logger().LogWithSeverity(map[string]string{"msg": msg}, 1)
}

func (g *Globals) LogErrorWithTrace(err *errtrace.Error) {
  g.Logger().LogWithSeverity(map[string]string{"msg": err.Error(), "trace": err.StringStack()}, 0)
}

func (g *Globals) Clients() map[string]interface{} {
  return map[string]interface{}{}
}

func TestJsonHttpApi(t *testing.T) {
  suite.Run(t, new(JsonHttpApiSuite))
}

const (
  GetIndexBody = "GET resources"
  PostCreateBody = "POST resources"
  GetDetailBody = "GET resource"
  PatchUpdateBody = "PATCH resource"
  DeleteDestroyBody = "DELETE resource"
  EmptyBody = "{}"
  JsonContentTypeHeader = "Content-Type"
  JsonContentType = "application/json; charset=utf-8"
)

type JsonHttpApiSuite struct {
  suite.Suite
  handler *jsonHttpHandler.JsonHttpHandler
}

func (suite *JsonHttpApiSuite) SetupSuite() {
  suite.handler = jsonHttpHandler.New(
    &Globals{},
    []jsonHttpHandler.RouteData{
      jsonHttpHandler.NewRouteData(
        http.MethodGet,
        "/resources",
        func(g jsonHttpHandler.Globals) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, GetIndexBody)
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
            fmt.Fprintf(w, PostCreateBody)
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
            fmt.Fprintf(w, fmt.Sprintf("%s %s", GetDetailBody, id))
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
            fmt.Fprintf(w, fmt.Sprintf("%s %s", PatchUpdateBody, id))
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
            fmt.Fprintf(w, fmt.Sprintf("%s %s", DeleteDestroyBody, id))
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
    },
  )
}

func (suite *JsonHttpApiSuite) TestIndexRoute() {
  request, _ := http.NewRequest("GET", "/resources", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{JsonContentType}),
    recorder.Header()[JsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), GetIndexBody, recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestDetailRoute() {
  id := 12
  request, _ := http.NewRequest("GET", fmt.Sprintf("/resources/%d", id), nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{JsonContentType}),
    recorder.Header()[JsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), fmt.Sprintf("%s %d", GetDetailBody, id), recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestCreateRoute() {
  request, _ := http.NewRequest("POST", "/resources", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{JsonContentType}),
    recorder.Header()[JsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusCreated, recorder.Code)
  assert.Equal(suite.T(), PostCreateBody, recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestUpdateRoute() {
  id := 12
  request, _ := http.NewRequest("PATCH", fmt.Sprintf("/resources/%d", id), nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{JsonContentType}),
    recorder.Header()[JsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), fmt.Sprintf("%s %d", PatchUpdateBody, id), recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestDestroyRoute() {
  id := 12
  request, _ := http.NewRequest("DELETE", fmt.Sprintf("/resources/%d", id), nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{JsonContentType}),
    recorder.Header()[JsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), fmt.Sprintf("%s %d", DeleteDestroyBody, id), recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestNotFound() {
  request, _ := http.NewRequest("GET", "/unknown_url", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{JsonContentType}),
    recorder.Header()[JsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusNotFound, recorder.Code)
  assert.Equal(suite.T(), EmptyBody, recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestInternalServerError() {
  request, _ := http.NewRequest("GET", "/error", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{JsonContentType}),
    recorder.Header()[JsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
  assert.Equal(suite.T(), EmptyBody, recorder.Body.String())
}
