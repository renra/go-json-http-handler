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

  AnyOrigin = "*"
  AnyMethod = "*"
  AllowedHeaders = "*"
  AccessControlMaxAgeHeaderValue = "86400"
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
  w.Header().Set(AllowMethodsHeader, AnyMethod)
  w.Header().Set(AllowHeadersHeader, AllowedHeaders)
  w.Header().Set(AccessControlMaxAgeHeader, AccessControlMaxAgeHeaderValue)
  w.Header().Set(AllowCredentialsHeader, "true")
}

func AddCorsHeadersForAnyOrigin(w http.ResponseWriter) {
  w.Header().Set(AllowOriginHeader, AnyOrigin)
}

func isLocalhost(origin string) bool {
  maybeUrl, err := url.Parse(origin)

  if err != nil {
    return false
  }

  return maybeUrl.Hostname() == "localhost"
}

type CorsAllower func([]string, string) bool

func CorsHandler(allowedOrigins []string, isCorsAllowed CorsAllower) GlobalsReceivingHandlerFunc {
  return func (g Globals) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set(VaryHeader, VaryHeaderValue)

      sliceWithOrigin, ok := r.Header[OriginHeader]

      if ok {
        origin := sliceWithOrigin[0]

        if isCorsAllowed(allowedOrigins, origin) {
          AddCorsHeaders(w, origin)
        }
      }

      w.WriteHeader(http.StatusNoContent)
      fmt.Fprintf(w, "")
    }
  }
}

func ListBasedCorsHandler(allowedOrigins []string) GlobalsReceivingHandlerFunc {
  return CorsHandler(allowedOrigins, func(origins []string, origin string) bool {
    for _, allowedOrigin := range origins {
      if allowedOrigin == origin {
        return true
      }
    }

    return false
  })
}

func ListBasedCorsHandlerWithLocalhost(allowedOrigins []string) GlobalsReceivingHandlerFunc {
  return CorsHandler(allowedOrigins, func(origins []string, origin string) bool {
    if isLocalhost(origin) {
      return true
    } else {
      for _, allowedOrigin := range origins {
        if allowedOrigin == origin {
          return true
        }
      }

      return false
    }
  })
}
