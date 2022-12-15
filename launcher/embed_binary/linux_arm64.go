//go:build linux && !android && arm64
// +build linux,!android,arm64

package embed_binary

import (
	_ "embed"
)

//go:embed go-cqhttp_linux_arm64.brotli
var embedding_cqhttp []byte
var PLANTFORM = Linux_arm64
