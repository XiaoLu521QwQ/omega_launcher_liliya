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

// 获取文件Hash
func GetFileHash(fname string) string {
	if IsFile(fname) {
		fileData, err := GetFileData(fname)
		if err != nil {
			panic(err)
		}
		return GetBinaryHash(fileData)
	}
	return ""
}
