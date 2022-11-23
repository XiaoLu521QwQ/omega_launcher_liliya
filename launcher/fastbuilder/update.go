package fastbuilder

import (
	"bytes"
	"io"
	"omega_launcher/utils"

	"github.com/pterm/pterm"
)

// Fastbuilder远程仓库地址
var STOARGE_REPO = ""
var REMOTE_REPO = "https://github.com/LNSSPsd/PhoenixBuilder/releases/latest/download/"
var MIRROR_REPO = "https://hub.fgit.ml/LNSSPsd/PhoenixBuilder/releases/latest/download/"
var LOCAL_REPO = "http://fileserver:12333/res/"

// 仓库选择
func selectRepo(cfg *BotConfig, reselect bool) *BotConfig {
	if reselect || cfg.Repo < 1 || cfg.Repo > 3 {
		pterm.Info.Printf("当前可选择的仓库有：\n1. Github仓库\n2. Github镜像仓库\n3. 本地仓库(自用)\n")
		pterm.Info.Printf("请输入序号来选择一个仓库: ")
		cfg.Repo = utils.GetIntInputInScope(1, 3)
	}
	if cfg.Repo == 1 {
		pterm.Info.Println("将使用Github仓库进行更新")
		STOARGE_REPO = REMOTE_REPO
	} else if cfg.Repo == 2 {
		pterm.Info.Println("将使用Github镜像仓库进行更新")
		STOARGE_REPO = MIRROR_REPO
	} else if cfg.Repo == 3 {
		pterm.Info.Println("将使用本地仓库进行更新")
		STOARGE_REPO = LOCAL_REPO
	} else {
		panic("无效的仓库值，请重新进行选择")
	}
	return cfg
}

// 下载FB
func downloadFB() {
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
func UpdateFB(cfg *BotConfig, reselect bool) *BotConfig {
	cfg = selectRepo(cfg, reselect)
	pterm.Warning.Println("正在从指定仓库获取更新信息..")
	targetHash := GetRemoteFBHash(STOARGE_REPO)
	currentHash := GetCurrentFBHash()
	//fmt.Println(targetHash)
	//fmt.Println(currentHash)
	if targetHash == currentHash {
		pterm.Success.Println("太好了，你的 Fastbuilder 已经是最新的了!")
	} else {
		pterm.Warning.Println("正在为你下载最新的 Fastbuilder, 请保持耐心..")
		downloadFB()
	}
	return cfg
}
