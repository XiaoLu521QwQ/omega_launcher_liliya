package fastbuilder

import (
	"encoding/json"
	"omega_launcher/embed_binary"
	"omega_launcher/utils"
	"path"
	"path/filepath"

	"github.com/pterm/pterm"
)

// 获取FB文件名
func GetFBExecName() string {
	name := ""
	switch embed_binary.GetPlantform() {
	case embed_binary.WINDOWS_arm64:
		// 不存在该构建
		name = ""
	case embed_binary.WINDOWS_x86_64:
		name = "phoenixbuilder-windows-executable-x86_64.exe"
	case embed_binary.Linux_arm64:
		name = "phoenixbuilder-aarch64"
	case embed_binary.Linux_x86_64:
		name = "phoenixbuilder"
	case embed_binary.MACOS_arm64:
		name = "phoenixbuilder-macos-arm64"
	case embed_binary.MACOS_x86_64:
		name = "phoenixbuilder-macos-x86_64"
	case embed_binary.Android_arm64:
		name = "phoenixbuilder-android-termux-shared-executable-arm64"
	case embed_binary.Android_x86_64:
		name = "phoenixbuilder-android-termux-shared-executable-x86_64"
	}
	if name == "" {
		panic("尚未支持该平台" + embed_binary.GetPlantform())
	}
	return name
}

// 获取FB文件路径
func getFBExecPath() string {
	path := path.Join(utils.GetCurrentDir(), GetFBExecName())
	result, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return result
}

// 获取本地FB文件Hash
func getCurrentFBHash() string {
	exec := getFBExecPath()
	return utils.GetFileHash(exec)
}

// 获取远程仓库的Hash
func getRemoteFBHash(url string) string {
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
		pterm.Error.Printfln("未能从远程仓库获取 Hash")
	}
	return hash
}

// 检查当前目录是否存在FB执行文件, 不存在时将会panic
func CheckExecFile() {
	if utils.IsFile(getFBExecPath()) {
		pterm.Success.Println("已从当前目录读取到 Fastbuilder")
	} else {
		pterm.Error.Println("请先下载 Fastbuilder 至当前目录")
		pterm.Error.Println("所需 Fastbuilder 的文件名为: " + GetFBExecName())
		panic("当前目录不存在 Fastbuilder")
	}
}
