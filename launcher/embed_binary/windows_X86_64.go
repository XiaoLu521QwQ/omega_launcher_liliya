//go:build windows && amd64
// +build windows,amd64

package embed_binary

import (
	_ "embed"
)

//go:embed go-cqhttp_windows_amd64.exe.brotli
var embedding_cqhttp []byte
var PLANTFORM = WINDOWS_x86_64
