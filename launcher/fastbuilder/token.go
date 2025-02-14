package fastbuilder

import (
	"encoding/json"
	"fmt"
	"omega_launcher/utils"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pterm/pterm"
	"golang.org/x/term"
)

// 加载现有的token
func loadCurrentFBToken() string {
	// 获取目录
	homedir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	fbconfigdir := filepath.Join(homedir, ".config/fastbuilder")
	token := filepath.Join(fbconfigdir, "fbtoken")
	// 尝试读取token文件
	if utils.IsFile(token) {
		if data, err := utils.GetFileData(token); err == nil {
			return string(data)
		}
	}
	return ""
}

// 请求token
func requestToken() string {
	// 尝试加载现有的token
	currentFbToken := loadCurrentFBToken()
	// 读取成功, 提示是否使用
	if currentFbToken != "" && strings.HasPrefix(currentFbToken, "w9/BeLNV/9") {
		pterm.Info.Printf("要使用现有的 Fastbuilder 账户登录吗? 使用现有账户请输入 y, 使用新账户请输入 n: ")
		if utils.GetInputYN() {
			return currentFbToken
		}
	}
	// 获取新的token
	pterm.Info.Printf("请输入 Fastbuilder 账号, 或者输入 Token: ")
	Code := utils.GetValidInput()
	// 输入token则直接返回
	if strings.HasPrefix(Code, "w9/BeLNV/9") {
		return Code
	}
	pterm.Info.Printf("请输入 Fastbuilder 密码 (不会回显): ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	if err != nil {
		panic(err)
	}
	Passwd := string(bytePassword)
	// 根据输入信息构建新token
	tokenstruct := &map[string]interface{}{
		"encrypt_token": true,
		"username":      Code,
		"password":      Passwd,
	}
	token, err := json.Marshal(tokenstruct)
	if err != nil {
		panic(err)
	}
	return string(token)
}
