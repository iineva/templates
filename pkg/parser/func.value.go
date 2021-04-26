package parser

import (
	"net/url"
)

func (p *Parser)funcRequestURL() string {
	return p.rootURL.String()
}

func (p *Parser)funcQuery(key string, query string) (string, error) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return "", err
	}
	return values.Get(key), nil
}

func (p *Parser)funcContextQuery(key string) (string, error) {
	targetURLObj, err := url.Parse(p.targetURL)
	if err != nil {
		return "", err
	}
	v := targetURLObj.Query().Get(key)
	if v != "" {
		return v, nil
	}
	return p.params.Get(key), nil
}

func (p *Parser)funcGetList(name string) []interface{} {
	return p.parserListMap[name]
}
// all always return ""
func (p *Parser)funcAddList(name string, item interface{}) string {
	p.parserListMap[name] = append(p.parserListMap[name], item)
	return ""
}
// all always return ""
func (p *Parser)funcDelList(name string) string {
	delete(p.parserListMap, name)
	return ""
}

// all always return ""
func (p *Parser)FuncSetValue(name string, value interface{}) string {
	p.parserValueMap[name] = value
	return ""
}
func (p *Parser)FuncGetValue(name string) interface{} {
	return p.parserValueMap[name]
}
// all always return ""
func (p *Parser)FuncDelValue(name string, value interface{}) string {
	delete(p.parserValueMap, name)
	return ""
}
