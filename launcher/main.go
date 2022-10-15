package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"omega_launcher/embed_binary"
	"omega_launcher/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pterm/pterm"
)

// 修改为官网地址
var STOARGE_REPO = "https://storage.fastbuilder.pro/"

type BotConfig struct {
	RentalCode       string `json:"租赁服号"`
	RentalPasswd     string `json:"租赁服密码"`
	FBToken          string `json:"FBToken"`
	QGroupLinkEnable bool   `json:"是否开启群服互通"`
	StartOmega       bool   `json:"是否启动Omega"`
	UpdateFB         bool   `json:"是否更新FB"`
}

func main() {
	// 添加启动信息
	pterm.Info.Printfln("Omega Launcher - Author: CMA2401PT")
	pterm.Info.Printfln("Modify By Liliya233")
	// 启动
	if err := os.Chdir(GetCurrentDir()); err != nil {
		panic(err)
	}
	StartOmegaHelper()
}

//go:embed config.yml
var defaultConfigBytes []byte

//go:embed 组件-群服互通-1.json
var defaultQGroupLinkConfigByte []byte

func CQHttpEnablerHelper() {
	if err := utils.WriteFileData(GetCqHttpExec(), GetCqHttpBinary()); err != nil {
		panic(err)
	}
	if !utils.MkDir(path.Join(GetCurrentDir(), "cqhttp_storage")) {
		panic("无法创建cqhttp_storage目录")
	}
	pterm.Warning.Println("请注意，你只能通过上传 config.yml 与 device.json 至 cqhttp_storage 目录的方式来为服务器配置群服互通")
	pterm.Info.Print("现在你可以进行文件上传的操作了，输入 y / 或者 n 继续配置群服互通: ")
	utils.GetInputYN()
	configFile := path.Join(GetCurrentDir(), "cqhttp_storage", "config.yml")
	omegaConfigFile := path.Join(GetCurrentDir(), "omega_storage", "配置", "群服互通", "组件-群服互通-1.json")
	if !utils.IsFile(configFile) {
		pterm.Info.Printf("请输入QQ账号: ")
		Code := utils.GetValidInput()
		pterm.Info.Printf("请输入QQ密码（想扫码登录则留空）: ")
		Passwd := utils.GetInput()
		if Passwd == "" {
			Passwd = "''"
		}
		defaultConfigStr := string(defaultConfigBytes)
		cfgStr := strings.ReplaceAll(defaultConfigStr, "[QQ账号]", Code)
		cfgStr = strings.ReplaceAll(cfgStr, "[QQ密码]", Passwd)
		utils.WriteFileData(configFile, []byte(cfgStr))
		pterm.Info.Printf("请输入想链接的群号: ")
		GroupCode := utils.GetValidInput()
		groupCfgStr := strings.ReplaceAll(string(defaultQGroupLinkConfigByte), "[群号]", GroupCode)
		utils.WriteFileData(omegaConfigFile, []byte(groupCfgStr))
	} else {
		pterm.Success.Println("尝试使用 cqhttp_storage 目录下的 config.yml 与 device.json 来配置群服互通")
	}
	RunCQHttp()
}

func WaitConnect() {
	for {
		u := url.URL{Scheme: "ws", Host: "127.0.0.1:6700"}
		var err error
		_, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			// time.Sleep(1)
			continue
		} else {
			return
		}
	}
}

func RunCQHttp() {
	args := []string{"-faststart"}
	cmd := exec.Command(GetCqHttpExec(), args...)
	cmd.Dir = path.Join(GetCurrentDir(), "cqhttp_storage")
	cqHttpOut, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		reader := bufio.NewReader(cqHttpOut)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				fmt.Println("reader exit")
				return
			}
			fmt.Print("[CQHTTP] " + readString)
		}
	}()
	go func() {
		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
	}()
	WaitConnect()
	pterm.Success.Println("CQ-Http已经成功启动了！")
	pterm.Info.Println("将 config.yml 与 device.json 上传至服务器 cqhttp_storage 目录下即可配置群服互通")
}

func LoadCurrentFBToken() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	fbconfigdir := filepath.Join(homedir, ".config/fastbuilder")
	token := filepath.Join(fbconfigdir, "fbtoken")
	if utils.IsFile(token) {
		if data, err := utils.GetFileData(token); err == nil {
			return string(data)
		}
	}
	return ""
}

func RequestToken() string {
	currentFbToken := LoadCurrentFBToken()
	if currentFbToken != "" && strings.HasPrefix(currentFbToken, "w9/BeLNV/9") {
		pterm.Info.Printf("要使用现有的 Fastbuilder 账户登录吗?  使用现有账户请输入 y , 使用新账户请输入 n: ")
		if utils.GetInputYN() {
			return currentFbToken
		}
	}
	pterm.Info.Printf("请输入 Fastbuilder 账号/或者输入 Token: ")
	Code := utils.GetValidInput()
	if strings.HasPrefix(Code, "w9/BeLNV/9") {
		// pterm.Success.Printfln("您输入的是 Token, 因此无需输入密码了")
		// time.Sleep(time.Second)
		return Code
	}
	pterm.Info.Printf("请输入 Fastbuilder 密码: ")
	Passwd := utils.GetValidInput()
	tokenstruct := &map[string]interface{}{
		"encrypt_token": true,
		"username":      Code,
		"password":      Passwd,
	}
	if token, err := json.Marshal(tokenstruct); err != nil {
		panic(err)
	} else {
		return string(token)
	}
}

func FBTokenSetup(cfg *BotConfig) {
	if cfg.FBToken != "" {
		pterm.Info.Printf("要使用上次的 Fastbuilder 账号登录吗?  要请输入 y , 需要修改请输入 n: ")
		if utils.GetInputYN() {
			return
		}
	}
	newToken := RequestToken()
	cfg.FBToken = newToken
}

func RentalServerSetup(cfg *BotConfig) {
	pterm.Info.Printf("请输入租赁服账号: ")
	cfg.RentalCode = utils.GetValidInput()
	pterm.Info.Printf("请输入租赁服密码（没有则留空）: ")
	cfg.RentalPasswd = utils.GetInput()
}

func UpdateFB() {
	pterm.Warning.Println("正在从官网获取更新信息...")
	targetHash := GetRemoteOmegaHash()
	currentHash := GetCurrentOmegaHash()
	//fmt.Println(targetHash)
	//fmt.Println(currentHash)
	if targetHash == currentHash {
		pterm.Success.Println("太好了，你的 Fastbuilder 已经是最新的了!")
	} else {
		pterm.Warning.Println("正在为你下载最新的 Fastbuilder, 请保持耐心...")
		DownloadOmega()
	}
}

func StartOmegaHelper() {
	// 读取配置出错则退出
	botConfig := &BotConfig{}
	if err := utils.GetJsonData(path.Join(GetCurrentDir(), "服务器登录配置.json"), botConfig); err != nil {
		panic(err)
	}
	// 询问是否使用上一次的配置
	if botConfig.FBToken != "" && botConfig.RentalCode != "" {
		pterm.Info.Printf("要使用和上次完全相同的配置启动吗? 要请输入 y, 不要请输入 n (10秒后会自动确认): ")
		if utils.GetInputYNTimeLimit(10) {
			// 更新FB
			if botConfig.UpdateFB {
				UpdateFB()
			}
			// 群服互通
			if botConfig.QGroupLinkEnable {
				if utils.IsDir(path.Join(GetCurrentDir(), "omega_storage")) {
					RunCQHttp()
				} else {
					pterm.Warning.Println("在Omega完全启动前，将不会进行群服互通的配置")
				}
			}
			Run(botConfig)
			return
		}
	}
	// 配置FB更新
	pterm.Info.Printf("需要从官网下载或更新 Fastbuilder 吗?  要请输入 y, 不要请输入 n: ")
	if utils.GetInputYN() {
		UpdateFB()
		botConfig.UpdateFB = true
	} else {
		pterm.Warning.Println("将会使用该路径的 Fastbuilder：" + GetOmegaExecName())
		botConfig.UpdateFB = false
		time.Sleep(time.Second)
	}
	// 配置FB
	FBTokenSetup(botConfig)
	// 配置租赁服登录
	if botConfig.RentalCode != "" {
		pterm.Info.Printf("要使用上次的租赁服配置吗?  要请输入 y, 不要请输入 n : ")
		if !utils.GetInputYN() {
			RentalServerSetup(botConfig)
		}
	} else {
		RentalServerSetup(botConfig)
	}
	// 询问是否使用Omega
	pterm.Info.Printf("要启动 Omega 还是 Fastbuilder?  启动 Omega 请输入 y, 启动 Fastbuilder 请输入 n: ")
	if utils.GetInputYN() {
		botConfig.StartOmega = true
		// 配置群服互通
		if utils.IsDir(path.Join(GetCurrentDir(), "omega_storage")) {
			pterm.Info.Printf("要启用群服互通吗?  要请输入 y, 不要请输入 n: ")
			if utils.GetInputYN() {
				CQHttpEnablerHelper()
				botConfig.QGroupLinkEnable = true
			}
		} else {
			pterm.Warning.Println("在Omega完全启动前，将不会进行群服互通的配置")
			botConfig.QGroupLinkEnable = false
		}
	}
	// 将本次配置写入文件
	if err := utils.WriteJsonData(path.Join(GetCurrentDir(), "服务器登录配置.json"), botConfig); err != nil {
		pterm.Error.Println("无法记录配置，不过可能不是什么大问题")
	}
	// 启动Omega或者FB
	Run(botConfig)
}

func Run(cfg *BotConfig) {
	// 敏感信息不应进行打印
	// fmt.Println(cfg.Token)
	args := []string{"-M", "--plain-token", cfg.FBToken, "--no-update-check", "-c", cfg.RentalCode}
	if cfg.RentalPasswd != "" {
		args = append(args, "-p")
		args = append(args, cfg.RentalPasswd)
	}
	// 是否启动Omega
	if cfg.StartOmega {
		args = append(args, "-O")
	}
	readC := make(chan string)
	go func() {
		for {
			s := utils.GetInput()
			readC <- s
		}
	}()
	// t := time.NewTicker(10 * time.Second)
	for {
		cmd := exec.Command(GetOmegaExecName(), args...)
		omega_out, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		omega_in, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		pterm.Success.Println("如果 Omega/Fastbuilder 崩溃了，它会在最长 10 秒后自动重启")
		stopped := false
		go func() {
			reader := bufio.NewReader(omega_out)
			for {
				readString, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					fmt.Println("reader exit")
					return
				}
				fmt.Print(readString)
			}
		}()

		go func() {
			for {
				s := <-readC
				if stopped {
					readC <- s
					return
				}
				omega_in.Write([]byte(s + "\n"))
			}
		}()

		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
		stopped = true
		// 为了避免频繁请求，崩溃后将等待10秒后重启，可手动跳过等待
		pterm.Error.Println("Oh no! Fastbuilder crashed!") // ?
		pterm.Warning.Print("似乎发生了错误，要重启 Fastbuilder 吗? 输入 y / 或者 n 确认(10秒后会自动确认): ")
		utils.GetInputYNTimeLimit(10)
	}
}

func GetCqHttpBinary() []byte {
	return embed_binary.GetCqHttpBinary()
}

func GetCurrentDir() string {
	// 兼容配套的Dockerfile
	if utils.IsFile(path.Join("/ome", "launcher_liliya")) {
		return path.Join("/workspace")
	}
	pathExecutable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dirPathExecutable := filepath.Dir(pathExecutable)
	return dirPathExecutable
}

func GetOmegaExecName() string {
	omega := "fastbuilder"
	switch GetPlantform() {
	case embed_binary.WINDOWS_x86_64:
		omega = "phoenixbuilder-windows-executable-x86_64.exe"
	case embed_binary.Linux_x86_64:
		omega = "phoenixbuilder"
	case embed_binary.MACOS_x86_64:
		omega = "phoenixbuilder-macos-x86_64"
	}
	omega = path.Join(GetCurrentDir(), omega)
	p, err := filepath.Abs(omega)
	if err != nil {
		panic(err)
	}
	return p
}

func GetCqHttpExec() string {
	cqhttp := "cqhttp"
	if GetPlantform() == embed_binary.WINDOWS_x86_64 {
		cqhttp = "cqhttp.exe"
	}
	cqhttp = path.Join(GetCurrentDir(), cqhttp)
	p, err := filepath.Abs(cqhttp)
	if err != nil {
		panic(err)
	}
	return p
}

func GetPlantform() string {
	return embed_binary.GetPlantform()
}

// 添加解析JSON的操作以适配官网
func GetRemoteOmegaHash() string {
	jsonData := utils.DownloadSmallContent(STOARGE_REPO + "hashes.json")
	var hash string
	hashMap := make(map[string]string, 0)
	if err := json.Unmarshal([]byte(jsonData), &hashMap); err != nil {
		panic(err)
	}
	switch GetPlantform() {
	case embed_binary.WINDOWS_x86_64:
		hash = hashMap["phoenixbuilder-windows-executable-x86_64.exe"]
	case embed_binary.Linux_x86_64:
		hash = hashMap["phoenixbuilder"]
	case embed_binary.MACOS_x86_64:
		hash = hashMap["phoenixbuilder-macos-x86_64"]
	default:
		panic("未知平台" + GetPlantform())
	}
	return hash
}

func GetFileHash(fname string) string {
	if utils.IsFile(fname) {
		fileData, err := utils.GetFileData(fname)
		if err != nil {
			panic(err)
		}
		return utils.GetBinaryHash(fileData)
	}
	return ""
}

func GetCurrentOmegaHash() string {
	exec := GetOmegaExecName()
	return GetFileHash(exec)
}

func GetCQHttpHash() string {
	exec := GetCqHttpExec()
	return GetFileHash(exec)
}

func GetEmbeddedCQHttpHash() string {
	return utils.GetBinaryHash(GetCqHttpBinary())
}

func DownloadOmega() {
	exec := GetOmegaExecName()
	url := ""
	switch GetPlantform() {
	case embed_binary.WINDOWS_x86_64:
		url = STOARGE_REPO + "phoenixbuilder-windows-executable-x86_64.exe"
	case embed_binary.Linux_x86_64:
		url = STOARGE_REPO + "phoenixbuilder"
	case embed_binary.MACOS_x86_64:
		url = STOARGE_REPO + "phoenixbuilder-macos-x86_64"
	default:
		panic("未知平台" + GetPlantform())
	}
	compressedData := utils.DownloadSmallContent(url)
	var execBytes []byte
	var err error
	// 官网并没有提供brotli，所以对读取操作进行修改
	if execBytes, err = io.ReadAll(bytes.NewReader(compressedData)); err != nil {
		panic(err)
	}
	if err := utils.WriteFileData(exec, execBytes); err != nil {
		panic(err)
	}
}
