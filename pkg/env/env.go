package env

import (
	"fmt"
	"net/http"
	"os"
)

var (
	publicURL = os.Getenv("PUBLIC_URL")
)

func PublicURL(r *http.Request) string {
	if publicURL != "" {
		return publicURL
	}

	// auto check public URL
	protocol := "http"
	if r.Header.Get("x-forwarded-proto") != "" {
		protocol = r.Header.Get("x-forwarded-proto")
	} else	if r.TLS != nil {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s", protocol, r.Host)
}

func Get(key string, def ...string) string {
	v := os.Getenv(key)
	if v == "" && len(def) != 0 {
		v = def[0]
	}
	return v
}
