package main

import (
	_ "embed"
	"omega_launcher/cqhttp"
	"omega_launcher/defines"
	"omega_launcher/fastbuilder"
	"omega_launcher/utils"
	"os"
	"path"

	"github.com/pterm/pterm"
)

//go:embed VERSION
var version []byte

func main() {
	// 添加启动信息
	pterm.DefaultBox.Println("https://github.com/Liliya233/omega_launcher")
	pterm.Info.Println("Omega Launcher", pterm.Yellow("(", string(version), ")"))
	pterm.Info.Println("Author: CMA2401PT, Modified: Liliya233")
	// 确保目录可用
	if err := os.Chdir(utils.GetCurrentDir()); err != nil {
		pterm.Error.Printf("读取当前目录时出现问题")
		panic(err)
	}
	// 启动
	// 读取配置出错则退出
	launcherConfig := &defines.LauncherConfig{}
	if err := utils.GetJsonData(path.Join(utils.GetCurrentDataDir(), "服务器登录配置.json"), launcherConfig); err != nil {
		panic(err)
	}
	// 询问是否使用上一次的配置
	if launcherConfig.FBToken != "" && launcherConfig.RentalCode != "" {
		pterm.Info.Printf("要使用和上次完全相同的配置启动吗? 要请输入 y, 不要请输入 n: ")
		if utils.GetInputYN() {
			// 更新FB
			if launcherConfig.UpdateFB {
				fastbuilder.Update(launcherConfig, false)
			} else {
				fastbuilder.CheckExecFile()
			}
			// 群服互通
			if launcherConfig.QGroupLinkEnable && launcherConfig.StartOmega {
				cqhttp.Run(launcherConfig)
			}
			// 启动Omega或者FB
			fastbuilder.Run(launcherConfig)
			return
		}
	}
	// 配置FB更新
	pterm.Info.Printf("需要启动器帮忙下载或更新 Fastbuilder 吗? 要请输入 y, 不要请输入 n: ")
	launcherConfig.UpdateFB = false
	if utils.GetInputYN() {
		fastbuilder.Update(launcherConfig, true)
		launcherConfig.UpdateFB = true
	} else {
		fastbuilder.CheckExecFile()
	}
	// 配置FB
	fastbuilder.FBTokenSetup(launcherConfig)
	// 配置租赁服登录
	if launcherConfig.RentalCode != "" {
		pterm.Info.Printf("要使用上次的租赁服配置吗? 要请输入 y, 不要请输入 n: ")
		if !utils.GetInputYN() {
			fastbuilder.RentalServerSetup(launcherConfig)
		}
	} else {
		fastbuilder.RentalServerSetup(launcherConfig)
	}
	// 询问是否使用Omega
	pterm.Info.Printf("要启动 Omega 还是 Fastbuilder? 启动 Omega 请输入 y, 启动 Fastbuilder 请输入 n: ")
	launcherConfig.StartOmega = false
	if utils.GetInputYN() {
		launcherConfig.StartOmega = true
		// 配置群服互通
		pterm.Info.Printf("需要启动器帮忙配置群服互通吗? 要请输入 y, 不要请输入 n: ")
		launcherConfig.QGroupLinkEnable = false
		if utils.GetInputYN() {
			launcherConfig.QGroupLinkEnable = true
			pterm.Info.Printf("需要在配置完成后屏蔽 go-cqhttp 的输出吗? 要请输入 y, 不要请输入 n: ")
			launcherConfig.BlockCQHttpOutput = false
			if utils.GetInputYN() {
				launcherConfig.BlockCQHttpOutput = true
			}
			if !utils.IsDir(path.Join(utils.GetCurrentDataDir(), "omega_storage", "配置")) {
				pterm.Warning.Printf("首次启动时配置群服互通会导致新生成的组件均为非启用状态, 要继续吗? 要请输入 y, 不要请输入 n: ")
				if utils.GetInputYN() {
					utils.MkDir(path.Join(utils.GetCurrentDataDir(), "omega_storage", "配置", "群服互通"))
					cqhttp.CQHttpEnablerHelper(launcherConfig)
				} else {
					launcherConfig.QGroupLinkEnable = false
				}
			} else {
				cqhttp.CQHttpEnablerHelper(launcherConfig)
			}
		}
	}
	// 启动Omega或者FB
	fastbuilder.Run(launcherConfig)
}
