//go:build android && arm64
// +build android,arm64

package embed_binary

import (
	_ "embed"
)

//go:embed assets/go-cqhttp_linux_arm64.brotli
var embedding_cqhttp []byte
var PLANTFORM = Android_arm64
