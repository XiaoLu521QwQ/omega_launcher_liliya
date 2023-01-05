package cqhttp

import (
	"fmt"
	"omega_launcher/defines"
	"omega_launcher/utils"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/pterm/pterm"
	"golang.org/x/term"
	v2 "gopkg.in/yaml.v2"
)

// 从cqhttp配置里读取QQ账密信息
func getCQConfig(cqCfg string) *defines.CQHttpConfig {
	cfg := &defines.CQHttpConfig{}
	data, err := os.ReadFile(cqCfg)
	if err != nil {
		return nil
	}
	if err := v2.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	if cfg.Account.Uin == "" {
		return nil
	}
	return cfg
}

// 将信息写入cqhttp配置文件
func updateCQConfig(cfg *defines.CQHttpConfig, address string) {
	cfgStr := strings.ReplaceAll(string(defaultConfigBytes), "[地址]", address)
	if cfg != nil {
		cfgStr = strings.ReplaceAll(cfgStr, "[QQ账号]", cfg.Account.Uin)
		cfgStr = strings.ReplaceAll(cfgStr, "[QQ密码]", cfg.Account.Password)
	} else {
		cfgStr = strings.ReplaceAll(cfgStr, "[QQ账号]", "null")
		cfgStr = strings.ReplaceAll(cfgStr, "[QQ密码]", "null")
	}
	err := utils.WriteFileData(path.Join(utils.GetCurrentDir(), "cqhttp_storage", "config.yml"), []byte(cfgStr))
	if err != nil {
		pterm.Error.Println("更新 go-cqhttp 配置文件时遇到问题")
		panic(err)
	}
}

// 初始化cqhttp配置文件
func cqhttpInit() {
	cfg := &defines.CQHttpConfig{}
	// 要求输入cqhttp配置信息
	pterm.Info.Printf("请输入QQ账号: ")
	cfg.Account.Uin = utils.GetValidInput()
	pterm.Info.Printf("请输入QQ密码 (不会回显): ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	if err != nil {
		panic(err)
	}
	cfg.Account.Password = string(bytePassword)
	updateCQConfig(cfg, "null")
}
