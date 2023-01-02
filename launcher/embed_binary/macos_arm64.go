//go:build darwin && arm64
// +build darwin,arm64

package embed_binary

import (
	_ "embed"
)

//go:embed assets/go-cqhttp_darwin_arm64.brotli
var embedding_cqhttp []byte
var PLANTFORM = MACOS_arm64
