package route

import (
	"io"
	"net/http"
)

func Healthy(p string, next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == p {
			io.WriteString(w, "")
			return
		}
		next.ServeHTTP(w, r)
  })
}
