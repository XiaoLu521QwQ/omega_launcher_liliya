package cqhttp

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"omega_launcher/utils"
	"os/exec"
	"path"
	"strings"

	"github.com/pterm/pterm"
)

var ADDRESS = "127.0.0.1:6700"

//go:embed raw/组件-群服互通-1.json
var defaultQGroupLinkConfigByte []byte

//go:embed raw/config.yml
var defaultConfigBytes []byte

func CQHttpEnablerHelper() {
	// 无法读写文件则退出
	if err := utils.WriteFileData(GetCqHttpExec(), GetCqHttpBinary()); err != nil {
		panic(err)
	}
	// 无法创建目录则退出
	if !utils.MkDir(path.Join(utils.GetCurrentDir(), "cqhttp_storage")) {
		panic("无法创建cqhttp_storage目录")
	}
	// 提示与确认
	pterm.Warning.Println("请注意，你只能通过上传 config.yml 与 device.json 至 cqhttp_storage 目录的方式来为服务器配置群服互通")
	pterm.Info.Print("现在你可以进行文件上传的操作了，输入 y / 或者 n 继续配置群服互通: ")
	utils.GetInputYN()
	// 配置文件所在路径
	configFile := path.Join(utils.GetCurrentDir(), "cqhttp_storage", "config.yml")
	omegaConfigFile := path.Join(utils.GetCurrentDir(), "omega_storage", "配置", "群服互通", "组件-群服互通-1.json")
	// 如果配置文件不存在，则执行初始化操作
	if !utils.IsFile(configFile) {
		// 获取cqhttp配置信息
		pterm.Info.Printf("请输入QQ账号: ")
		Code := utils.GetValidInput()
		pterm.Info.Printf("请输入QQ密码（想扫码登录则留空）: ")
		Passwd := utils.GetInput()
		if Passwd == "" {
			Passwd = "''"
		}
		// 将获取的信息写入到cqhttp配置文件
		cfgStr := strings.ReplaceAll(string(defaultConfigBytes), "[地址]", ADDRESS)
		cfgStr = strings.ReplaceAll(cfgStr, "[QQ账号]", Code)
		cfgStr = strings.ReplaceAll(cfgStr, "[QQ密码]", Passwd)
		utils.WriteFileData(configFile, []byte(cfgStr))
		// 获取Omega组件配置信息
		pterm.Info.Printf("请输入想链接的群号: ")
		GroupCode := utils.GetValidInput()
		// 将获取的信息写入到Omega组件文件
		groupCfgStr := strings.ReplaceAll(string(defaultQGroupLinkConfigByte), "[地址]", ADDRESS)
		groupCfgStr = strings.ReplaceAll(groupCfgStr, "[群号]", GroupCode)
		utils.WriteFileData(omegaConfigFile, []byte(groupCfgStr))
	} else {
		pterm.Success.Println("尝试使用 cqhttp_storage 目录下的 config.yml 与 device.json 来配置群服互通")
	}
	// 运行cqhttp
	RunCQHttp()
}

func RunCQHttp() {
	// 配置启动参数
	args := []string{"-faststart"}
	// 配置执行目录
	cmd := exec.Command(GetCqHttpExec(), args...)
	cmd.Dir = path.Join(utils.GetCurrentDir(), "cqhttp_storage")
	// 建立cqhttp的输出管道
	cqHttpOut, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	// 从管道中获取并打印cqhttp输出内容
	go func() {
		reader := bufio.NewReader(cqHttpOut)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				pterm.Error.Println("读取CQHTTP输出内容时出现错误")
				return
			}
			fmt.Print("[CQHTTP] " + readString)
		}
	}()
	// 启动并持续运行cqhttp
	go func() {
		err = cmd.Start()
		if err != nil {
			pterm.Error.Println("CQHTTP启动时出现错误")
		}
		err = cmd.Wait()
		if err != nil {
			pterm.Error.Println("CQHTTP运行时出现错误")
		}
	}()
	// 等待cqhttp启动完成
	WaitConnect(ADDRESS)
	pterm.Success.Println("CQ-Http已经成功启动了！")
	pterm.Info.Println("将 config.yml 与 device.json 上传至服务器 cqhttp_storage 目录下即可配置群服互通")
}
