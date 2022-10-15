package fastbuilder

import (
	"encoding/json"
	"omega_launcher/embed_binary"
	"omega_launcher/utils"
	"path"
	"path/filepath"
)

// 获取FB文件名
func GetFBExecName() string {
	name := ""
	switch embed_binary.GetPlantform() {
	case embed_binary.WINDOWS_x86_64:
		name = "phoenixbuilder-windows-executable-x86_64.exe"
	case embed_binary.Linux_x86_64:
		name = "phoenixbuilder"
	case embed_binary.MACOS_x86_64:
		name = "phoenixbuilder-macos-x86_64"
	}
	return name
}

// 获取FB文件路径
func GetFBExecPath() string {
	path := path.Join(utils.GetCurrentDir(), GetFBExecName())
	result, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return result
}

// 获取本地FB文件Hash
func GetCurrentFBHash() string {
	exec := GetFBExecPath()
	return utils.GetFileHash(exec)
}

// 获取远程仓库的Hash
func GetRemoteFBHash(url string) string {
	// 获取文件内容
	jsonData := utils.DownloadSmallContent(url + "hashes.json")
	// 解析文件内容
	var hash string
	hashMap := make(map[string]string, 0)
	if err := json.Unmarshal([]byte(jsonData), &hashMap); err != nil {
		panic(err)
	}
	hash = hashMap[GetFBExecName()]
	if hash == "" {
		panic("未知平台" + embed_binary.GetPlantform())
	}
	return hash
}
