package parser

import (
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type Parser struct {
	targetURL   string
	rootURL     *url.URL
	params      url.Values
	template    *template.Template
	parent      *Parser
	importCount int

	parserListMap  map[string][]interface{}
	parserValueMap map[string]interface{}

	localMode bool
}

type Values struct {
	Search url.Values
}

func New(rootURL *url.URL, targetURL string) (*Parser, error) {

	rootURLParam := rootURL.Query()

	p := &Parser{
		rootURL:        rootURL,
		targetURL:      targetURL,
		params:         rootURLParam,
		parserListMap:  make(map[string][]interface{}),
		parserValueMap: make(map[string]interface{}),
		localMode:      strings.HasPrefix(targetURL, "file://"),
	}

	if err := p.init(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Parser) init() error {
	// block env functions
	sprigFuncMap := sprig.TxtFuncMap()
	delete(sprigFuncMap, "env")
	delete(sprigFuncMap, "expandenv")

	// 修复非标准base64编码的解码（缺少等号）
	sprigFuncMap["b64dec"] = base64decode
	sprigFuncMap["b32dec"] = base32decode

	p.template = template.New(p.targetURL).Funcs(sprigFuncMap).Funcs(template.FuncMap{
		// 获取参数，优先读取当前链接参数
		"queryContext": p.funcContextQuery,
		// 生成外网链接
		"requestURL": p.funcRequestURL,
		// import from URL
		"import": p.funcImport,
		// include like helm do
		"include": p.funcInclude,
		// load define
		"loadDefine": p.funcLoadDefine,
		// url encode/decode
		"urlEncode": urlEncode,
		"urlDecode": urlDecode,
		// parse url
		"fromUrl": fromUrl,
		// get list value
		"getList": p.funcGetList,
		"addList": p.funcAddList,
		"delList": p.funcDelList,
		// get value
		"getValue": p.FuncGetValue,
		"setValue": p.FuncSetValue,
		"delValue": p.FuncDelValue,
		// query
		"query": query,
		// filter list with regex
		"filterList": filterList,

		"toToml":        toTOML,
		"toYaml":        toYAML,
		"fromYaml":      fromYAML,
		"fromYamlArray": fromYAMLArray,
		"toJson":        toJSON,
		"fromJson":      fromJSON,
		"fromJsonArray": fromJSONArray,
	})

	// TODO: for test
	//p.addDefine("test/surge.rule.tpl")
	//p.addDefine("test/subscribe.tpl")
	return nil
}

func (p *Parser) addDefine(fileName string) error {
	def, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	defTpl := p.template.New(fileName)
	if _, err := defTpl.Parse(string(def)); err != nil {
		return err
	}
	return nil
}

func (p *Parser) isDeadCycle(u string) bool {
	if u == p.targetURL {
		return true
	}
	if p.parent != nil {
		return p.parent.isDeadCycle(u)
	}
	return false
}

func (p *Parser) rootParser() *Parser {
	if p.parent != nil {
		return p.parent.rootParser()
	}
	return p
}

func (p *Parser) loadTarget() (string, error) {
	targetURL, err := p.getTargetURL()
	if err != nil {
		return "", err
	}

	return p.httpGet(targetURL.String())
}

func (p *Parser) ParseAndWrite(w io.Writer) error {

	tpl, err := p.loadTarget()
	if err != nil {
		return err
	}

	if _, err := p.template.Parse(string(tpl)); err != nil {
		return err
	}
	return p.template.Execute(w, Values{
		Search: p.params,
	})
}

func (p *Parser) ParseAndWriteFromReader(r io.Reader, w io.Writer) error {

	tpl, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if _, err := p.template.Parse(string(tpl)); err != nil {
		return err
	}
	return p.template.Execute(w, Values{
		Search: p.params,
	})
}
