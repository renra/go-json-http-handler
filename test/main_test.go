package test

import (
  "fmt"
  "net/http"
  "strings"
  "testing"
  "net/http/httptest"
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
  getIndexBody = "GET resources"
  postCreateBody = "POST resources"
  getDetailBody = "GET resource"
  patchUpdateBody = "PATCH resource"
  deleteDestroyBody = "DELETE resource"
  emptyBody = "{}"
  jsonContentTypeHeader = "Content-Type"
  jsonContentType = "application/json; charset=utf-8"
  wantedPayload = "foo"
  middlewaresBody = "middlewares"
)

type JsonHttpApiSuite struct {
  suite.Suite
  handler *jsonHttpHandler.JsonHttpHandler
}

func requirePayload(g jsonHttpHandler.Globals, next http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    payload, _ := r.Context().Value(jsonHttpHandler.PayloadKey).(*string)

    if payload == nil {
      w.WriteHeader(http.StatusUnauthorized)
      fmt.Fprintf(w, "")
    } else {
      if *payload != "" {
        next(w, r)
      } else {
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "")
      }
    }
  }
}

func requireExactPayload(pattern string) jsonHttpHandler.Middleware {
  return func(g jsonHttpHandler.Globals, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      payload := r.Context().Value(jsonHttpHandler.PayloadKey).(*string)

      if payload == nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "")
      } else {
        if *payload == pattern {
          next(w, r)
        } else {
          w.WriteHeader(http.StatusBadRequest)
          fmt.Fprintf(w, "")
        }
      }
    }
  }
}

func (suite *JsonHttpApiSuite) SetupTest() {
  suite.handler = getBasicJsonHttpHandler()
}

func (suite *JsonHttpApiSuite) TestIndexRoute() {
  request, _ := http.NewRequest(http.MethodGet, "/resources", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{jsonContentType}),
    recorder.Header()[jsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), getIndexBody, recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestDetailRoute() {
  id := 12
  request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/resources/%d", id), nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{jsonContentType}),
    recorder.Header()[jsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), fmt.Sprintf("%s %d", getDetailBody, id), recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestCreateRoute() {
  request, _ := http.NewRequest(http.MethodPost, "/resources", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{jsonContentType}),
    recorder.Header()[jsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusCreated, recorder.Code)
  assert.Equal(suite.T(), postCreateBody, recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestUpdateRoute() {
  id := 12
  request, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/resources/%d", id), nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{jsonContentType}),
    recorder.Header()[jsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), fmt.Sprintf("%s %d", patchUpdateBody, id), recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestDestroyRoute() {
  id := 12
  request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/resources/%d", id), nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{jsonContentType}),
    recorder.Header()[jsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), fmt.Sprintf("%s %d", deleteDestroyBody, id), recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestNotFound() {
  request, _ := http.NewRequest(http.MethodGet, "/unknown_url", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{jsonContentType}),
    recorder.Header()[jsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusNotFound, recorder.Code)
  assert.Equal(suite.T(), emptyBody, recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestInternalServerError() {
  request, _ := http.NewRequest(http.MethodGet, "/error", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{jsonContentType}),
    recorder.Header()[jsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
  assert.Equal(suite.T(), emptyBody, recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestMiddlewares() {
  request, _ := http.NewRequest(http.MethodPost, "/middlewares", strings.NewReader(wantedPayload))
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{jsonContentType}),
    recorder.Header()[jsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), middlewaresBody, recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestMiddlewares_PayloadMismatch() {
  request, _ := http.NewRequest(http.MethodPost, "/middlewares", strings.NewReader("whatever"))
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{jsonContentType}),
    recorder.Header()[jsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
  assert.Equal(suite.T(), "", recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestMiddlewares_NoPayload() {
  request, _ := http.NewRequest(http.MethodPost, "/middlewares", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    []string([]string{jsonContentType}),
    recorder.Header()[jsonContentTypeHeader],
  )

  assert.Equal(suite.T(), http.StatusUnauthorized, recorder.Code)
  assert.Equal(suite.T(), "", recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestCORS_WithoutAnyHandler() {
  request, _ := http.NewRequest(http.MethodOptions, "/resources", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    0,
    len(recorder.Header()[jsonContentTypeHeader]),
  )

  assert.Equal(suite.T(), http.StatusNotFound, recorder.Code)
  assert.Equal(suite.T(), "", recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestCors_WithCustomHandler() {
  suite.handler = getJsonHttpHandlerWithCors(
    func (g jsonHttpHandler.Globals) http.HandlerFunc {
      return func (w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "")
      }
    },
  )

  request, _ := http.NewRequest(http.MethodOptions, "/resources", nil)
  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    0,
    len(recorder.Header()[jsonContentTypeHeader]),
  )

  assert.Equal(suite.T(), http.StatusOK, recorder.Code)
  assert.Equal(suite.T(), "", recorder.Body.String())
}

func (suite *JsonHttpApiSuite) TestCors_WithProvidedHandler() {
  origin1 := "https://server.foo.bar"
  origin2 := "https://server.whatever"

  allowedOrigins := []string{origin1, origin2}

  suite.handler = getJsonHttpHandlerWithCors(
    jsonHttpHandler.ListBasedCorsHandler(allowedOrigins),
  )

  request, _ := http.NewRequest(http.MethodOptions, "/resources", nil)
  request.Header.Add(jsonHttpHandler.OriginHeader, origin1)

  recorder := httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    0,
    len(recorder.Header()[jsonContentTypeHeader]),
  )

  assert.Equal(
    suite.T(),
    []string{origin1},
    recorder.Header()[jsonHttpHandler.AllowOriginHeader],
  )

  assert.Equal(
    suite.T(),
    []string{jsonHttpHandler.AllMethods},
    recorder.Header()[jsonHttpHandler.AllowMethodsHeader],
  )

  assert.Equal(
    suite.T(),
    []string{"true"},
    recorder.Header()[jsonHttpHandler.AllowCredentialsHeader],
  )

  assert.Equal(
    suite.T(),
    []string{jsonHttpHandler.AccessControlMaxAgeHeaderValue},
    recorder.Header()[jsonHttpHandler.AccessControlMaxAgeHeader],
  )

  assert.Equal(
    suite.T(),
    []string{jsonHttpHandler.VaryHeaderValue},
    recorder.Header()[jsonHttpHandler.VaryHeader],
  )

  assert.Equal(suite.T(), http.StatusNoContent, recorder.Code)
  assert.Equal(suite.T(), "", recorder.Body.String())

  request, _ = http.NewRequest(http.MethodOptions, "/resources", nil)
  request.Header.Add(jsonHttpHandler.OriginHeader, origin2)

  recorder = httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    0,
    len(recorder.Header()[jsonContentTypeHeader]),
  )

  assert.Equal(
    suite.T(),
    []string{origin2},
    recorder.Header()[jsonHttpHandler.AllowOriginHeader],
  )

  assert.Equal(
    suite.T(),
    []string{jsonHttpHandler.AllMethods},
    recorder.Header()[jsonHttpHandler.AllowMethodsHeader],
  )

  assert.Equal(
    suite.T(),
    []string{"true"},
    recorder.Header()[jsonHttpHandler.AllowCredentialsHeader],
  )

  assert.Equal(
    suite.T(),
    []string{jsonHttpHandler.AccessControlMaxAgeHeaderValue},
    recorder.Header()[jsonHttpHandler.AccessControlMaxAgeHeader],
  )

  assert.Equal(
    suite.T(),
    []string{jsonHttpHandler.VaryHeaderValue},
    recorder.Header()[jsonHttpHandler.VaryHeader],
  )

  assert.Equal(suite.T(), http.StatusNoContent, recorder.Code)
  assert.Equal(suite.T(), "", recorder.Body.String())

  request, _ = http.NewRequest(http.MethodOptions, "/resources", nil)
  request.Header.Add(jsonHttpHandler.OriginHeader, "https://some-totally-different.domain.net")

  recorder = httptest.NewRecorder()

  suite.handler.ServeHTTP(recorder, request)

  assert.Equal(
    suite.T(),
    0,
    len(recorder.Header()[jsonContentTypeHeader]),
  )

  assert.Equal(
    suite.T(),
    0,
    len(recorder.Header()[jsonHttpHandler.AllowOriginHeader]),
  )

  assert.Equal(
    suite.T(),
    0,
    len(recorder.Header()[jsonHttpHandler.AllowMethodsHeader]),
  )

  assert.Equal(
    suite.T(),
    0,
    len(recorder.Header()[jsonHttpHandler.AllowCredentialsHeader]),
  )

  assert.Equal(
    suite.T(),
    0,
    len(recorder.Header()[jsonHttpHandler.AccessControlMaxAgeHeader]),
  )

  assert.Equal(
    suite.T(),
    []string{jsonHttpHandler.VaryHeaderValue},
    recorder.Header()[jsonHttpHandler.VaryHeader],
  )

  assert.Equal(suite.T(), http.StatusNoContent, recorder.Code)
  assert.Equal(suite.T(), "", recorder.Body.String())
}
