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
    map[string]jsonHttpHandler.GlobalsReceivingHandlerFunc{
      "/resources": func(g jsonHttpHandler.Globals) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
          w.WriteHeader(http.StatusOK)
          fmt.Fprintf(w, GetIndexBody)
        }
      },
      "/resource": func(g jsonHttpHandler.Globals) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
          w.WriteHeader(http.StatusOK)
          fmt.Fprintf(w, GetDetailBody)
        }
      },
      "/resources_create": func(g jsonHttpHandler.Globals) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
          w.WriteHeader(http.StatusCreated)
          fmt.Fprintf(w, PostCreateBody)
        }
      },
      "/resource_update": func(g jsonHttpHandler.Globals) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
          w.WriteHeader(http.StatusOK)
          fmt.Fprintf(w, PatchUpdateBody)
        }
      },
      "/resource_delete": func(g jsonHttpHandler.Globals) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
          w.WriteHeader(http.StatusOK)
          fmt.Fprintf(w, DeleteDestroyBody)
        }
      },
      "/error": func(g jsonHttpHandler.Globals) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
          panic("The impossible has happened")
        }
      },
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
  request, _ := http.NewRequest("GET", "/resource", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{JsonContentType}),
    recorder.Header()[JsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), GetDetailBody, recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestCreateRoute() {
  request, _ := http.NewRequest("POST", "/resources_create", nil)
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
  request, _ := http.NewRequest("PATCH", "/resource_update", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{JsonContentType}),
    recorder.Header()[JsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), PatchUpdateBody, recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestDestroyRoute() {
  request, _ := http.NewRequest("DELETE", "/resource_delete", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{JsonContentType}),
    recorder.Header()[JsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), DeleteDestroyBody, recorder.Body.String())
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
