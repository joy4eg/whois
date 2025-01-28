package data

import "embed"

// Files is an embedded filesystem with configuration files.
//
//go:embed *.json
var Files embed.FS
