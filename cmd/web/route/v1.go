package route

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/iineva/templates/pkg/env"
	"github.com/iineva/templates/pkg/parser"
)

func checkIsSelfRequest(r *http.Request) bool {
	u := r.Header.Get("user-agent")
	return strings.Contains(u, "templates/1.0")
}

func HandlerTemplate(w http.ResponseWriter, r *http.Request) {

	var requestURL *url.URL
	var targetURL string
	var p *parser.Parser
	var err error

	if checkIsSelfRequest(r) {
		err = errors.New("cant not call server it self")
		goto error_block
	}

	requestURL, _ = url.ParseRequestURI(env.PublicURL(r) + r.URL.String())
	targetURL = r.URL.Query().Get("url")
	if targetURL == "" {
		err = errors.New("url = null")
		goto error_block
	}
	p, err = parser.New(requestURL, targetURL)

	if err != nil {
		goto error_block
	}

	err = p.ParseAndWrite(w)
	if err != nil {
		goto error_block
	}
	return

error_block:
	log.Printf("error: %v", err.Error())
  http.Error(w, err.Error(), 500)
}
