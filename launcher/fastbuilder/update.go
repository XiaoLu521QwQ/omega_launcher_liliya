package fastbuilder

import (
	"bytes"
	"io"
	"omega_launcher/utils"

	"github.com/pterm/pterm"
)

// Fastbuilder远程仓库地址
var STOARGE_REPO = "https://storage.fastbuilder.pro/"

// 下载FB
func DownloadFB() {
	var execBytes []byte
	var err error
	// 获取写入路径与远程仓库url
	path := GetFBExecPath()
	url := STOARGE_REPO + GetFBExecName()
	// 下载
	compressedData := utils.DownloadSmallContent(url)
	// 官网并没有提供brotli，所以对读取操作进行修改
	if execBytes, err = io.ReadAll(bytes.NewReader(compressedData)); err != nil {
		panic(err)
	}
	// 写入文件
	if err := utils.WriteFileData(path, execBytes); err != nil {
		panic(err)
	}
}

// 升级FB
func UpdateFB() {
	pterm.Warning.Println("正在从官网获取更新信息...")
	targetHash := GetRemoteFBHash(STOARGE_REPO)
	currentHash := GetCurrentFBHash()
	//fmt.Println(targetHash)
	//fmt.Println(currentHash)
	if targetHash == currentHash {
		pterm.Success.Println("太好了，你的 Fastbuilder 已经是最新的了!")
	} else {
		pterm.Warning.Println("正在为你下载最新的 Fastbuilder, 请保持耐心...")
		DownloadFB()
	}
}
