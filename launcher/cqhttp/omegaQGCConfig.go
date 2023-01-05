package cqhttp

import (
	"encoding/json"
	"io/fs"
	"omega_launcher/defines"
	"omega_launcher/utils"
	"path"
	"path/filepath"
	"strings"
)

func getOmegaConfig() (string, *defines.ComponentConfig) {
	// 默认的空配置
	cfg := &defines.ComponentConfig{}
	// 默认配置文件路径
	fp := path.Join(utils.GetCurrentDataDir(), "omega_storage", "配置", "群服互通", "组件-群服互通.json")
	// 尝试从配置文件夹下寻找全部群服互通配置文件
	if err := filepath.Walk(path.Join(utils.GetCurrentDataDir(), "omega_storage", "配置"), func(filePath string, info fs.FileInfo, err error) error {
		// 跳过目录
		if info.IsDir() {
			return nil
		}
		// 识别非json组件文件并跳过
		fileName := info.Name()
		if !strings.HasPrefix(fileName, "组件") || !strings.HasSuffix(fileName, ".json") {
			return nil
		}
		// 对配置文件进行解析
		currentCfg := &defines.ComponentConfig{}
		if parseErr := utils.GetJsonData(filePath, currentCfg); parseErr != nil {
			return nil
		}
		// 如果不是群服互通组件, 则跳过
		if currentCfg.Name != "群服互通" {
			return nil
		}
		// 如果存在多个群服互通组件, 则报错
		if cfg.Configs != nil {
			panic("当前存在多个群服互通组件, 请自行删除多余的群服互通组件")
		}
		// 更新配置与路径信息
		cfg = currentCfg
		fp = filePath
		return nil
	}); err != nil {
		panic(err)
	}
	// 未找到配置时, 使用默认配置
	if cfg.Name != "群服互通" {
		err := json.Unmarshal(defaultQGroupLinkConfigByte, cfg)
		if err != nil {
			panic(err)
		}
	}
	// 将组件修改为开启状态
	cfg.Disabled = false
	return fp, cfg
}

func updateOmegaConfig(fp string, cfg *defines.ComponentConfig) {
	err := utils.WriteJsonData(fp, cfg)
	if err != nil {
		panic(err)
	}
}
