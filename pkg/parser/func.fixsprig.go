package parser

import (
	"encoding/base32"
	"encoding/base64"
)

func base64decode(v string) string {
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		// NOTE: do nothing
	}
	return string(data)
}

func base32decode(v string) string {
	data, err := base32.StdEncoding.DecodeString(v)
	if err != nil {
		// NOTE: do nothing
	}
	return string(data)
}
