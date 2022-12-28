package fastbuilder

import (
	"bufio"
	"fmt"
	"io"
	"omega_launcher/cqhttp"
	"omega_launcher/embed_binary"
	"omega_launcher/utils"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/pterm/pterm"
	"golang.org/x/term"
)

// 启动器配置文件结构
type BotConfig struct {
	Repo             int    `json:"仓库序号"`
	RentalCode       string `json:"租赁服号"`
	RentalPasswd     string `json:"租赁服密码"`
	FBToken          string `json:"FBToken"`
	QGroupLinkEnable bool   `json:"是否开启群服互通"`
	StartOmega       bool   `json:"是否启动Omega"`
	UpdateFB         bool   `json:"是否更新FB"`
}

// 保存配置文件
func saveConfig(cfg *BotConfig) {
	if err := utils.WriteJsonData(path.Join(utils.GetCurrentDataDir(), "服务器登录配置.json"), cfg); err != nil {
		pterm.Error.Println("无法记录配置, 不过可能不是什么大问题")
	}
}

// 配置Token
func FBTokenSetup(cfg *BotConfig) *BotConfig {
	if cfg.FBToken != "" {
		pterm.Info.Printf("要使用上次的 Fastbuilder 账号登录吗? 要请输入 y, 需要修改请输入 n: ")
		if utils.GetInputYN() {
			return cfg
		}
	}
	cfg.FBToken = RequestToken()
	return cfg
}

// 配置租赁服信息
func RentalServerSetup(cfg *BotConfig) *BotConfig {
	pterm.Info.Printf("请输入租赁服号: ")
	cfg.RentalCode = utils.GetValidInput()
	pterm.Info.Printf("请输入租赁服密码 (没有则留空, 不会回显): ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	if err != nil {
		panic(err)
	}
	cfg.RentalPasswd = string(bytePassword)
	return cfg
}

func StartHelper() {
	// 读取配置出错则退出
	botConfig := &BotConfig{}
	if err := utils.GetJsonData(path.Join(utils.GetCurrentDataDir(), "服务器登录配置.json"), botConfig); err != nil {
		panic(err)
	}
	// 询问是否使用上一次的配置
	if botConfig.FBToken != "" && botConfig.RentalCode != "" {
		pterm.Info.Printf("要使用和上次完全相同的配置启动吗? 要请输入 y, 不要请输入 n: ")
		if utils.GetInputYN() {
			// 更新FB
			if botConfig.UpdateFB {
				UpdateFB(botConfig, false)
			}
			// 群服互通
			if botConfig.QGroupLinkEnable && botConfig.StartOmega {
				cqhttp.RunCQHttp()
			}
			// 启动Omega或者FB
			Run(botConfig)
			return
		}
	}
	// 配置FB更新
	pterm.Info.Printf("需要启动器帮忙下载或更新 Fastbuilder 吗? 要请输入 y, 不要请输入 n: ")
	botConfig.UpdateFB = false
	if utils.GetInputYN() {
		UpdateFB(botConfig, true)
		botConfig.UpdateFB = true
	} else {
		pterm.Warning.Println("将会使用该路径的 Fastbuilder: " + GetFBExecPath())
		time.Sleep(time.Second)
	}
	// 配置FB
	botConfig = FBTokenSetup(botConfig)
	// 配置租赁服登录
	if botConfig.RentalCode != "" {
		pterm.Info.Printf("要使用上次的租赁服配置吗? 要请输入 y, 不要请输入 n: ")
		if !utils.GetInputYN() {
			botConfig = RentalServerSetup(botConfig)
		}
	} else {
		botConfig = RentalServerSetup(botConfig)
	}
	// 询问是否使用Omega
	pterm.Info.Printf("要启动 Omega 还是 Fastbuilder? 启动 Omega 请输入 y, 启动 Fastbuilder 请输入 n: ")
	botConfig.StartOmega = false
	if utils.GetInputYN() {
		botConfig.StartOmega = true
		// 配置群服互通
		pterm.Info.Printf("需要启动器帮忙配置群服互通吗? 要请输入 y, 不要请输入 n: ")
		botConfig.QGroupLinkEnable = false
		if utils.GetInputYN() {
			botConfig.QGroupLinkEnable = true
			if !utils.IsDir(path.Join(utils.GetCurrentDataDir(), "omega_storage")) {
				pterm.Warning.Printf("首次启动时配置群服互通会导致新生成的组件均为非启用状态, 要继续吗? 要请输入 y, 不要请输入 n: ")
				if utils.GetInputYN() {
					cqhttp.CQHttpEnablerHelper()
				} else {
					botConfig.QGroupLinkEnable = false
				}
			} else {
				cqhttp.CQHttpEnablerHelper()
			}
		}
	}
	// 启动Omega或者FB
	Run(botConfig)
}

func Run(cfg *BotConfig) {
	// 打印警告信息, Windows新版终端存在此问题，暂时没找到解决方法（
	plantform := embed_binary.GetPlantform()
	if plantform == embed_binary.WINDOWS_arm64 || plantform == embed_binary.WINDOWS_x86_64 {
		pterm.Warning.Println("对于Windows新版终端, 直接点击关闭按钮会导致程序在后台持续运行")
	}
	// 配置启动参数
	args := []string{"-M", "--plain-token", cfg.FBToken, "--no-update-check", "-c", cfg.RentalCode}
	// 是否需要租赁服密码
	if cfg.RentalPasswd != "" {
		args = append(args, "-p")
		args = append(args, cfg.RentalPasswd)
	}
	// 是否启动Omega
	if cfg.StartOmega {
		args = append(args, "-O")
		pterm.Warning.Println("请使用 stop 命令来正确的退出程序")
	} else {
		pterm.Warning.Println("请使用 exit / fbexit 命令来正确的退出程序")
	}
	// 建立频道
	readC := make(chan string)
	stop := make(chan string)
	// 持续将输入信息输入到频道中
	go func() {
		for {
			s := utils.GetInput()
			readC <- s
		}
	}()
	// 读取验证服务器返回的Token并保存
	go func() {
		for {
			if strings.HasPrefix(cfg.FBToken, "w9/BeLNV/9") {
				pterm.Success.Println("成功获取到Token")
				saveConfig(cfg)
				return
			}
			cfg.FBToken = LoadCurrentFBToken()
		}
	}()
	// 重启间隔
	restartTime := 0
	for {
		// 记录启动时间
		startTime := time.Now()
		// 是否停止
		isStopped := false
		// 最近一次输入, 用于忽略对输入内容的重复输出
		lastInput := ""
		// 启动时提示信息
		pterm.Success.Println("如果 Omega/Fastbuilder 崩溃了, 它将在一段时间后自动重启")
		// 启动命令
		cmd := exec.Command(GetFBExecPath(), args...)
		cmd.Dir = path.Join(utils.GetCurrentDataDir())
		// 建立从Fastbuilder到控制台的输出管道
		omega_out, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		// 建立从控制台到Fastbuilder的输入管道
		omega_in, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		// 建立从Fastbuilder到控制台的错误管道
		omega_err, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		// 从管道中获取并打印Fastbuilder输出内容
		go func() {
			reader := bufio.NewReader(omega_out)
			for {
				readString, err := reader.ReadString('\n')
				readString = strings.TrimPrefix(readString, "> ")
				if lastInput != "" && strings.HasPrefix(readString, lastInput) {
					lastInput = ""
					continue
				}
				lastInput = ""
				if readString == "\n" {
					continue
				}
				if err != nil || err == io.EOF {
					//pterm.Error.Println("读取 Omega/Fastbuilder 输出内容时出现错误")
					return
				}
				fmt.Print(readString + "\033[0m")
			}
		}()
		// 从管道中获取并打印Fastbuilder错误内容
		go func() {
			reader := bufio.NewReader(omega_err)
			for {
				readString, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					//pterm.Error.Println("读取 Omega/Fastbuilder 错误内容时出现错误")
					return
				}
				pterm.Error.Println(readString)
			}
		}()
		// 在未收到停止信号前, 启动器会一直将控制台输入的内容通过管道发送给Fastbuilder
		go func() {
			for {
				select {
				case <-stop:
					return
				case s := <-readC:
					lastInput = s
					// 接收到停止命令时处理
					if (cfg.StartOmega && s == "stop") || s == "exit" || s == "fbexit" {
						// 关闭重启
						isStopped = true
						// 发出停止命令
						omega_in.Write([]byte(s + "\n"))
						// 输出信息
						pterm.Success.Println("正在等待 Omega/Fastbuilder 处理退出命令")
						// 停止接收输入
						return
					} else {
						omega_in.Write([]byte(s + "\n"))
					}
				}
			}
		}()
		// 启动并持续运行Fastbuilder
		err = cmd.Start()
		if err != nil {
			pterm.Error.Println("Omega/Fastbuilder 启动时出现错误")
			pterm.Error.Println(err)
		}
		err = cmd.Wait()
		if err != nil {
			pterm.Error.Println("Omega/Fastbuilder 运行时出现错误")
			pterm.Error.Println(err)
		}
		// 如果运行到这里, 说明Fastbuilder出现错误或退出运行了
		cmd.Process.Kill()
		// 判断是否正常退出
		if isStopped {
			pterm.Success.Println("Omega/Fastbuilder 已正常退出, 启动器将结束运行")
			time.Sleep(3 * time.Second)
			break
		} else {
			stop <- "stop!!"
			pterm.Error.Println("Oh no! Fastbuilder crashed!") // ?
		}
		// 为了避免频繁请求, 崩溃后将等待一段时间后重启, 可手动跳过等待
		if time.Since(startTime) < time.Minute {
			if restartTime < 1<<20 {
				restartTime = restartTime<<1 + 1
			}
		} else {
			restartTime = 0
		}
		pterm.Warning.Printf("似乎发生了错误, %d秒后会重新启动 Omega/Fastbuilder (按回车立即重启)", restartTime)
		// 等待输入或计时结束
		select {
		case <-readC:
			restartTime = 0
		case <-time.After(time.Second * time.Duration(restartTime)):
			fmt.Println("")
		}
	}
}
