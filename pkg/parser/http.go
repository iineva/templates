package parser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/iineva/templates/pkg/random"
)

const cacheExpire = time.Second * 30

type cache struct {
	data []byte
	date time.Time
}

var (
	cacheMap = sync.Map{}
	clientID = random.RandString(8)
)

func (p *Parser) loadFile(targetURL string) ([]byte, error) {
	log.Printf("FILE GET: %v", targetURL)

	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	d, err := os.ReadFile(u.Path)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (p *Parser) loadHttp(targetURL string) ([]byte, error) {
	log.Printf("HTTP GET: %v", targetURL)

	req, _ := http.NewRequest("GET", targetURL, nil)
	req.Header.Add("user-agent", fmt.Sprintf("templates/1.0 (id %s)", clientID))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("http status code: " + resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *Parser) httpGet(targetURL string) (string, error) {
	// load cache
	if c, ok := cacheMap.Load(targetURL); ok {
		ca := c.(cache)
		if time.Since(ca.date) < cacheExpire {
			return string(ca.data), nil
		}
	}

	// delete cache
	cacheMap.Delete(targetURL)

	var data []byte
	if strings.HasPrefix(targetURL, "file://") {
		d, err := p.loadFile(targetURL)
		if err != nil {
			return "", err
		}
		data = d
	} else if strings.HasPrefix(targetURL, "http") {
		d, err := p.loadHttp(targetURL)
		if err != nil {
			return "", err
		}
		data = d
	}

	// cache
	cacheMap.Store(targetURL, cache{
		data: data,
		date: time.Now(),
	})

	return string(data), nil
}

func (p *Parser) getTargetURL() (*url.URL, error) {
	targetURL, err := url.Parse(p.targetURL)
	if err != nil {
		return nil, err
	}

	if targetURL.Scheme != "" {
		return targetURL, nil
	}

	// get parent full url
	var u *url.URL
	if p.parent != nil {
		pURL, err := p.parent.getTargetURL()
		if err != nil {
			return nil, err
		}
		u = pURL
	}

	targetURL.Scheme = u.Scheme
	targetURL.Host = u.Host
	if !strings.HasPrefix(targetURL.Path, "/") {
		// ./file.conf or file.conf to http://exemple.com/path/to/file.conf
		targetURL.Path = path.Join(path.Dir(u.Path), targetURL.Path)
	}

	return targetURL, nil
}

func (p *Parser) getAbsoluteURL(u string) (*url.URL, error) {
	absoluteURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(absoluteURL.Scheme, "http") {
		return absoluteURL, nil
	}

	targetURL, err := p.getTargetURL()
	if err != nil {
		return nil, err
	}

	absoluteURL.Scheme = targetURL.Scheme
	absoluteURL.Host = targetURL.Host
	if !strings.HasPrefix(absoluteURL.Path, "/") {
		// ./file.conf or file.conf to http://exemple.com/path/to/file.conf
		absoluteURL.Path = path.Join(path.Dir(targetURL.Path), absoluteURL.Path)
	}

	return absoluteURL, nil
}
