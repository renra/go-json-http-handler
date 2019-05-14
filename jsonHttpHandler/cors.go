package jsonHttpHandler

import (
  "fmt"
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

func ListBasedCorsHandler(allowedOrigins []string) GlobalsReceivingHandlerFunc {
  return func (g Globals) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set(VaryHeader, VaryHeaderValue)

      sliceWithOrigin, ok := r.Header[OriginHeader]

      if ok {
        origin := sliceWithOrigin[0]

        for _, allowedOrigin := range allowedOrigins {
          if allowedOrigin == origin {
            w.Header().Set(AllowOriginHeader, origin)
            w.Header().Set(AllowMethodsHeader, AllMethods)
            w.Header().Set(AccessControlMaxAgeHeader, AccessControlMaxAgeHeaderValue)
            w.Header().Set(AllowCredentialsHeader, "true")
            break
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


    //w.Header().Set(AllowOriginHeader, "*")
    //w.Header().Set(AllowMethodsHeader, "POST, GET, OPTIONS, PUT, DELETE")
    //w.Header().Set(AllowHeadersHeader, "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    //w.Header().Set(AllowCredentialsHeader, true)
