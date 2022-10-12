package utils

import (
	"crypto/sha256"
	"fmt"
)

// 修改hash计算方式以适配官网
func GetBinaryHash(fileData []byte) string {
	cvt := func(in [32]byte) []byte {
		return in[:32]
	}
	hashedBytes := cvt(sha256.Sum256(fileData))
	return fmt.Sprintf("%x", hashedBytes)
}
