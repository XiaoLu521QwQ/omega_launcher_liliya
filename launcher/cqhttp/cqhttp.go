package cqhttp

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"omega_launcher/utils"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/pterm/pterm"
	"golang.org/x/term"
)

//go:embed raw/组件-群服互通-1.json
var defaultQGroupLinkConfigByte []byte

//go:embed raw/config.yml
var defaultConfigBytes []byte

type QGroupLink struct {
	Address                   string                        `json:"CQHTTP正向Websocket代理地址"`
	GameMessageFormat         string                        `json:"游戏消息格式化模版"`
	QQMessageFormat           string                        `json:"Q群消息格式化模版"`
	Groups                    map[string]int64              `json:"链接的Q群"`
	Selector                  string                        `json:"游戏内可以听到QQ消息的玩家的选择器"`
	NoBotMsg                  bool                          `json:"不要转发机器人的消息"`
	ChatOnly                  bool                          `json:"只转发聊天消息"`
	MuteIgnored               bool                          `json:"屏蔽其他群的消息"`
	FilterQQToServerMsgByHead string                        `json:"仅仅转发开头为以下特定字符的消息到服务器"`
	FilterServerToQQMsgByHead string                        `json:"仅仅转发开头为以下特定字符的消息到QQ"`
	AllowedCmdExecutor        map[int64]bool                `json:"允许这些人透过QQ执行命令"`
	AllowdFakeCmdExecutor     map[int64]map[string][]string `json:"允许这些人透过QQ执行伪命令"`
	DenyCmds                  map[string]string             `json:"屏蔽这些指令"`
}

type ComponentConfig struct {
	Name        string      `json:"名称"`
	Description string      `json:"描述"`
	Disabled    bool        `json:"是否禁用"`
	Version     string      `json:"版本"`
	Source      string      `json:"来源"`
	Configs     *QGroupLink `json:"配置"`
}

func CQHttpEnablerHelper() {
	// 尝试读取Omega配置, 读取出错时使用默认配置
	cfg := &ComponentConfig{}
	if err := utils.GetJsonData(path.Join(utils.GetCurrentDir(), "omega_storage", "配置", "群服互通", "组件-群服互通-1.json"), cfg); err != nil {
		err := json.Unmarshal(defaultQGroupLinkConfigByte, cfg)
		if err != nil {
			panic(err)
		}
	}
	// 解压go-cqhttp失败则退出
	if err := utils.WriteFileData(GetCqHttpExec(), GetCqHttpBinary()); err != nil {
		panic(err)
	}
	// 无法创建目录则退出
	if !utils.MkDir(path.Join(utils.GetCurrentDir(), "cqhttp_storage")) {
		panic("无法创建cqhttp_storage目录")
	}
	// 提示与确认
	pterm.Warning.Println("请注意, 非本地环境只能通过上传 config.yml, device.json 与 session.token 的方式来配置 go-cqhttp")
	pterm.Info.Print("现在你可以进行文件上传的操作了, 输入 y 继续配置 go-cqhttp: ")
	utils.GetInputYN()
	// 配置文件路径
	configFile := path.Join(utils.GetCurrentDir(), "cqhttp_storage", "config.yml")
	omegaConfigFile := path.Join(utils.GetCurrentDir(), "omega_storage", "配置", "群服互通", "组件-群服互通-1.json")
	// 如果go-cqhttp配置文件不存在, 则执行初始化操作
	if !utils.IsFile(configFile) {
		// 获取cqhttp配置信息
		pterm.Info.Printf("请输入QQ账号: ")
		Code := utils.GetValidInput()
		pterm.Info.Printf("请输入QQ密码 (不会回显): ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Print("\n")
		if err != nil {
			panic(err)
		}
		Passwd := string(bytePassword)
		// 检查Omega配置文件的地址是否可用
		if !utils.IsAddressAvailable(cfg.Configs.Address) {
			port, err := utils.GetAvailablePort()
			if err != nil {
				pterm.Error.Println("无法为 go-cqhttp 获取可用端口")
				panic(err)
			}
			cfg.Configs.Address = fmt.Sprintf("127.0.0.1:%d", port)
		}
		// 将获取的信息写入到cqhttp配置文件
		cfgStr := strings.ReplaceAll(string(defaultConfigBytes), "[地址]", cfg.Configs.Address)
		cfgStr = strings.ReplaceAll(cfgStr, "[QQ账号]", Code)
		cfgStr = strings.ReplaceAll(cfgStr, "[QQ密码]", fmt.Sprintf("'%s'", Passwd))
		err = utils.WriteFileData(configFile, []byte(cfgStr))
		if err != nil {
			panic(err)
		}
		// 更新Omega群服互通组件配置
		err = utils.WriteJsonData(omegaConfigFile, cfg)
		if err != nil {
			panic(err)
		}
	} else {
		pterm.Success.Println("尝试使用 cqhttp_storage 目录下的 config.yml, device.json 与 session.token 来配置 go-cqhttp")
	}
	// 运行cqhttp
	RunCQHttp()
}

func RunCQHttp() {
	pterm.Warning.Println("如果长时间未启动Omega, 请检查 config.yml 与 组件-群服互通-1.json 设置的地址是否一致")
	// 尝试读取Omega配置, 读取出错时使用默认配置
	cfg := &ComponentConfig{}
	if err := utils.GetJsonData(path.Join(utils.GetCurrentDir(), "omega_storage", "配置", "群服互通", "组件-群服互通-1.json"), cfg); err != nil {
		err := json.Unmarshal(defaultQGroupLinkConfigByte, cfg)
		if err != nil {
			panic(err)
		}
	}
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
	WaitConnect(cfg.Configs.Address)
	pterm.Success.Println("go-cqhttp 已经成功启动了, 可前往 omega_storage 文件夹对群服互通组件进行进一步配置")
	pterm.Info.Println("若要为服务器配置群服互通, 请执行以下的操作：")
	pterm.Info.Println("1. 在服务器使用启动器配置群服互通, 直至看到\"现在你可以进行文件上传..\"的提示")
	pterm.Info.Println("2. 将本地 cqhttp_storage 目录下的 config.yml, device.json 与 session.token 上传至服务器同样的目录下")
	pterm.Info.Println("3. 将本地 omega_storage/配置/群服互通 目录下的 组件-群服互通-1.json 上传至服务器同样的目录下")
	pterm.Info.Println("4. 在服务器上进行确认, 此时应该配置成功了")
	pterm.Info.Println("如果遇到意料之外的问题, 请重新操作或前往 https://docs.go-cqhttp.org/ 寻找可用的信息")
}
