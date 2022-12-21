//go:build (android && amd64) || tag_android_x86_64
// +build android,amd64 tag_android_x86_64

package embed_binary

import (
	_ "embed"
)

//go:embed go-cqhttp_linux_amd64.brotli
var embedding_cqhttp []byte
var PLANTFORM = Android_x86_64
