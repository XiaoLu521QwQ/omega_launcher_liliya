package cqhttp

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"omega_launcher/defines"
	"omega_launcher/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pterm/pterm"
	"golang.org/x/term"
	v2 "gopkg.in/yaml.v2"
)

//go:embed assets/组件-群服互通.json
var defaultQGroupLinkConfigByte []byte

//go:embed assets/config.yml
var defaultConfigBytes []byte

// 从cqhttp配置里读取QQ账密信息
func getInfoFormCQConfig(cqCfg string) *defines.CQHttpConfig {
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
func updateConfig(cqCfg, address, account, password string) {
	// 将获取的信息写入到cqhttp配置文件
	cfgStr := strings.ReplaceAll(string(defaultConfigBytes), "[地址]", address)
	cfgStr = strings.ReplaceAll(cfgStr, "[QQ账号]", account)
	cfgStr = strings.ReplaceAll(cfgStr, "[QQ密码]", password)
	err := utils.WriteFileData(cqCfg, []byte(cfgStr))
	if err != nil {
		panic(err)
	}
}

func cqhttpInit(cfg *defines.ComponentConfig, cqCfg string) {
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
	updateConfig(cqCfg, cfg.Configs.Address, Code, Passwd)
}

func getOmegaConfig() *defines.ComponentConfig {
	// 确保此路径可用
	utils.MkDir(path.Join(utils.GetCurrentDataDir(), "omega_storage", "配置", "群服互通"))
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
			panic("当前存在多个群服互通组件, 请手动删除多余的群服互通组件")
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
	// 更新Omega群服互通组件配置
	err := utils.WriteJsonData(fp, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func CQHttpEnablerHelper(botCfg *defines.LauncherConfig) {
	cfg := getOmegaConfig()
	// 解压go-cqhttp失败则退出
	if err := utils.WriteFileData(GetCqHttpExec(), GetCqHttpBinary()); err != nil {
		panic(err)
	}
	// 无法创建目录则退出
	if !utils.MkDir(path.Join(utils.GetCurrentDataDir(), "cqhttp_storage")) {
		panic("无法创建 cqhttp_storage 目录")
	}
	// 提示与确认
	pterm.Warning.Println("请注意, 非本地环境只能上传 config.yml, device.json 与 session.token 来配置 go-cqhttp")
	pterm.Info.Print("现在你可以进行文件上传的操作了, 输入 y 继续配置 go-cqhttp: ")
	utils.GetInputYN()
	// 配置文件路径
	cqCfg := path.Join(utils.GetCurrentDataDir(), "cqhttp_storage", "config.yml")
	// 如果go-cqhttp配置文件不存在, 则执行初始化操作
	if utils.IsFile(cqCfg) {
		pterm.Info.Print("已读取到 go-cqhttp 配置文件, 要使用吗? 使用请输入 y, 需要重新设置请输入 n: ")
		if utils.GetInputYN() {
			pterm.Warning.Println("尝试使用现有的 config.yml, device.json 与 session.token 来配置 go-cqhttp")
			// 在使用上次的配置时，将读取cqhttp配置文件的账密以及群服互通组件的地址，然后对cqhttp配置文件进行更新
			if re := getInfoFormCQConfig(cqCfg); re != nil {
				updateConfig(cqCfg, cfg.Configs.Address, re.Account.Uin, re.Account.Password)
			}
		} else {
			cqhttpInit(cfg, cqCfg)
		}
	} else {
		cqhttpInit(cfg, cqCfg)
	}
	// 运行cqhttp
	Run(botCfg)
}

func Run(botCfg *defines.LauncherConfig) {
	// 读取Omega配置
	cfg := getOmegaConfig()
	// 配置文件路径
	cqCfg := path.Join(utils.GetCurrentDataDir(), "cqhttp_storage", "config.yml")
	// 启动前, 将Omega配置内的IP地址同步到go-cqhttp配置文件
	if re := getInfoFormCQConfig(cqCfg); re != nil {
		updateConfig(cqCfg, cfg.Configs.Address, re.Account.Uin, re.Account.Password)
	}
	pterm.Warning.Println("如果长时间未启动 Omega, 请检查 config.yml 与 群服互通组件 设置的地址是否一致")
	// 配置启动参数
	args := []string{"-faststart"}
	// 配置执行目录
	cmd := exec.Command(GetCqHttpExec(), args...)
	cmd.Dir = path.Join(utils.GetCurrentDataDir(), "cqhttp_storage")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cqhttp_out, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	// 从管道中获取并打印cqhttp输出内容
	stopOutput := false
	go func() {
		reader := bufio.NewReader(cqhttp_out)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			if stopOutput {
				return
			}
			fmt.Print(readString + "\033[0m")
		}
	}()
	// 启动并持续运行cqhttp
	go func() {
		err := cmd.Start()
		if err != nil {
			pterm.Error.Println("go-cqhttp 启动时出现错误")
			panic(err)
		}
		err = cmd.Wait()
		if err != nil {
			pterm.Error.Println("go-cqhttp 运行时出现错误")
			panic(err)
		}
	}()
	// 等待cqhttp启动完成
	WaitConnect(cfg.Configs.Address)
	// 配置完成后, 根据设置决定是否关闭go-cqhttp输出
	if botCfg.BlockCQHttpOutput {
		pterm.Success.Println("将屏蔽 go-cqhttp 的输出内容")
		stopOutput = true
	}
	pterm.Success.Println("go-cqhttp 已经成功启动了, 可前往 omega_storage 文件夹对群服互通组件进行进一步配置")
	pterm.Info.Println("若要为服务器配置群服互通, 请执行以下的操作：")
	pterm.Info.Println("1. 在服务器使用启动器配置群服互通, 直至启动器要求进行文件上传操作")
	pterm.Info.Println("2. 将 cqhttp_storage 目录下的 config.yml, device.json 与 session.token 上传至服务器同样的目录下")
	pterm.Info.Println("3. 将 omega_storage/配置 目录下的 群服互通组件 上传至服务器同样的目录下")
	pterm.Info.Println("4. 在服务器上进行确认, 此时应该配置成功了")
	pterm.Info.Println("如果仍未配置成功, 请删除现有的 go-cqhttp 与 群服互通组件 后再重新进行配置")
	pterm.Info.Println("如果遇到关于 go-cqhttp 意料之外的问题, 可前往 https://docs.go-cqhttp.org/ 寻找可用的信息")
}
