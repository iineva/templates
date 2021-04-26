package templates

import "embed"

//go:embed config
//go:embed define
//go:embed rules
var FS embed.FS
