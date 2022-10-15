package main

import (
	_ "embed"
	"omega_launcher/fastbuilder"
	"omega_launcher/utils"
	"os"

	"github.com/pterm/pterm"
)

func main() {
	// 添加启动信息
	pterm.Info.Printfln("Omega Launcher - Author: CMA2401PT")
	pterm.Info.Printfln("Modify By Liliya233")
	// 启动
	if err := os.Chdir(utils.GetCurrentDir()); err != nil {
		panic(err)
	}
	fastbuilder.StartHelper()
}
