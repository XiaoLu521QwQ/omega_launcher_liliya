//go:build darwin && amd64
// +build darwin,amd64

package embed_binary

import (
	_ "embed"
)

//go:embed assets/go-cqhttp_darwin_amd64.brotli
var embedding_cqhttp []byte
var PLANTFORM = MACOS_x86_64
