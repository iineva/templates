package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const recursionMaxNums = 100

func (p *Parser)funcImport(target string) (string, error) {

	rootURL := p.rootURL

	newParser, err := New(rootURL, target)
	if err != nil {
		return "", err
	}
	newParser.parent = p

	if newParser.rootParser().importCount > recursionMaxNums {
		return "", errors.New("import count out of limit, may be dead cycle or import too much URLs")
	}
	newParser.rootParser().importCount++

	if target == rootURL.String() || p.isDeadCycle(target) {
		return "", fmt.Errorf("dead cycle %v to import %v", p.targetURL, target)
	}

	b := &bytes.Buffer{}
	if err := newParser.ParseAndWrite(b); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (p *Parser)funcInclude(name string, data... interface{}) (string, error) {
	var buf strings.Builder
	var param interface{}
	if len(data) == 1 {
		param = data[0]
	} else {
		param = data
	}
	err := p.template.ExecuteTemplate(&buf, name, param)
	return buf.String(), err
}

func (p *Parser)funcLoadDefine(u string) (string, error) {
	targetURL, err := p.getAbsoluteURL(u)
	if err != nil {
		return "", err
	}
	content, err := p.httpGet(targetURL.String())
	if err != nil {
		return "", err
	}
	defTpl := p.template.New(u)
	if _, err := defTpl.Parse(string(content)); err != nil {
		return "", err
	}
	return "", nil
}
