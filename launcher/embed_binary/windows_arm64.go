//go:build windows && arm64
// +build windows,arm64

package embed_binary

import (
	_ "embed"
)

//go:embed assets/go-cqhttp_windows_arm64.exe.brotli
var embedding_cqhttp []byte
var PLANTFORM = WINDOWS_arm64
