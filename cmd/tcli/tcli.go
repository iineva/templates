package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/iineva/templates/pkg/parser"
)

var (
	ErrMethodNotFound = errors.New("method not found")
)

func main() {

	m := flag.String("m", "", "call spring func")
	f := flag.String("f", "", "parse file")
	a := flag.String("a", "", "parse file's args, like url query")
	flag.Usage = usage
	flag.Parse()

	output := bytes.NewBuffer(make([]byte, 0))
	if *m != "" {
		p, err := parser.New(&url.URL{}, "")
		if err != nil {
			panic(err)
		}
		tpl := fmt.Sprintf("{{ %s }}", *m)
		err = p.ParseAndWriteFromReader(bytes.NewBufferString(tpl), output)
		if err != nil {
			panic(err)
		}
	} else if *f != "" {
		fp, err := filepath.Abs(*f)
		fp = "file://" + fp

		p, err := parser.New(&url.URL{
			RawQuery: *a,
		}, fp)
		if err != nil {
			panic(err)
		}

		err = p.ParseAndWrite(output)
		if err != nil {
			panic(err)
		}
	}

	log.Print("parse done!")

	io.Copy(output, bytes.NewBufferString("\n"))

	io.Copy(os.Stdout, output)
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: tcli [options]
Options:
`)
	flag.PrintDefaults()
}
