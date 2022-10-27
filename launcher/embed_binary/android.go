//go:build android
// +build android

package embed_binary

import (
	_ "embed"
)

//go:embed cqhttp-android
var embedding_cqhttp []byte
var PLANTFORM = Android_arm64
