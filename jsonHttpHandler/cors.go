package jsonHttpHandler

import (
  "fmt"
  "net/url"
  "net/http"
)

const (
  OriginHeader = "Origin"
  AllowOriginHeader = "Access-Control-Allow-Origin"
  AllowMethodsHeader = "Access-Control-Allow-Methods"
  AllowHeadersHeader = "Access-Control-Allow-Headers"
  AllowCredentialsHeader = "Access-Control-Allow-Credentials"
  AccessControlMaxAgeHeader = "Access-Control-Max-Age"
  VaryHeader = "Vary"

  AllMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD"
  AccessControlMaxAgeHeaderValue = "1728000"
  VaryHeaderValue = "Origin"
)

func corsNoop(g Globals) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "")
  }
}

func AddCorsHeaders(w http.ResponseWriter, origin string) {
  w.Header().Set(AllowOriginHeader, origin)
  w.Header().Set(AllowMethodsHeader, AllMethods)
  w.Header().Set(AccessControlMaxAgeHeader, AccessControlMaxAgeHeaderValue)
  w.Header().Set(AllowCredentialsHeader, "true")
}

func isLocalhost(origin string) bool {
  maybeUrl, err := url.Parse(origin)

  if err != nil {
    fmt.Println(origin)
    fmt.Println(fmt.Sprintf("%v", err))
    return false
  }

  return maybeUrl.Hostname() == "localhost"
}

func ListBasedCorsHandler(allowedOrigins []string) GlobalsReceivingHandlerFunc {
  return func (g Globals) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set(VaryHeader, VaryHeaderValue)

      sliceWithOrigin, ok := r.Header[OriginHeader]

      if ok {
        origin := sliceWithOrigin[0]

        if isLocalhost(origin) {
          AddCorsHeaders(w, origin)
        } else {
          for _, allowedOrigin := range allowedOrigins {
            if allowedOrigin == origin {
              AddCorsHeaders(w, origin)
              break
            }
          }
        }

        w.WriteHeader(http.StatusNoContent)
        fmt.Fprintf(w, "")
      } else {
        w.WriteHeader(http.StatusNoContent)
        fmt.Fprintf(w, "")
      }
    }
  }
}

func ListBasedCorsHandlerWithLocalhost(allowedOrigins []string) GlobalsReceivingHandlerFunc {
  return func (g Globals) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set(VaryHeader, VaryHeaderValue)

      sliceWithOrigin, ok := r.Header[OriginHeader]

      if ok {
        origin := sliceWithOrigin[0]

        if isLocalhost(origin) {
          AddCorsHeaders(w, origin)
        } else {
          for _, allowedOrigin := range allowedOrigins {
            if allowedOrigin == origin {
              AddCorsHeaders(w, origin)
              break
            }
          }
        }

        w.WriteHeader(http.StatusNoContent)
        fmt.Fprintf(w, "")
      } else {
        w.WriteHeader(http.StatusNoContent)
        fmt.Fprintf(w, "")
      }
    }
  }
}
