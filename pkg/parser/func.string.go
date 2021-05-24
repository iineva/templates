package parser

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
)

func urlDecode(s string) (string, error) {
	return url.QueryUnescape(s)
}

func urlEncode(s string) string {
	return url.QueryEscape(s)
}

// Javascript URL
type JavascriptURL struct {
	GoURL    *url.URL // golang url.URL
	Hash     string   // "#xxx"
	Host     string   // "github.com:8888"
	Hostname string   // "github.com"
	Href     string   // "https://u:p@github.com:8888/search?l=Go&q=url&type=Repositories#xxx"
	Origin   string   // "https://github.com:8888"
	Password string   // "p"
	Pathname string   // "/search"
	Port     string   // "8888"
	Protocol string   // "https:"
	Search   string   // "?l=Go&q=url&type=Repositories"
	Username string   // "u"
}

func fromUrl(s string) (*JavascriptURL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	password, _ := u.User.Password()
	pu := &JavascriptURL{
		GoURL:    u, // golang url.URL
		Hash:     "#" + u.Fragment,
		Host:     u.Host,
		Hostname: u.Hostname(),
		Href:     s,
		Origin:   fmt.Sprintf("%v://%v", u.Scheme, u.Host),
		Password: password,
		Pathname: u.Path,
		Port:     u.Port(),
		Protocol: u.Scheme + ":",
		Search:   "?" + u.RawQuery,
		Username: u.User.Username(),
	}
	return pu, nil
}

func filterList(regex string, list []interface{}) ([]interface{}, error) {

	newList := []interface{}{}

	if len(list) == 0 {
		return newList, nil
	}

	switch list[0].(type) {
	case string:
		for _, item := range list {
			ok, err := regexp.MatchString(regex, item.(string))
			if err != nil {
				log.Print(err.Error())
				return newList, err
			}
			if ok {
				newList = append(newList, item)
			}
		}
	default:
		return newList, errors.New("unkonw type")
	}

	return newList, nil
}
