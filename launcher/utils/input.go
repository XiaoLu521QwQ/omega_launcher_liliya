package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

func GetInput() string {
	buf := bufio.NewReader(os.Stdin)
	l, _, _ := buf.ReadLine()
	return string(strings.TrimSpace(string(l)))
}

func GetValidInput() string {
	for {
		s := GetInput()
		if s == "" {
			pterm.Error.Println("无效输入，输入不能为空")
			continue
		}
		return s
	}
}

func GetInputYN() bool {
	for {
		s := GetInput()
		if strings.HasPrefix(s, "y") || strings.HasPrefix(s, "Y") {
			return true
		} else if strings.HasPrefix(s, "n") || strings.HasPrefix(s, "N") {
			return false
		}
		pterm.Error.Println("无效输入，输入应该为 y 或者 n")
	}
}

func GetInputYNTimeLimit(sec int) bool {
	chn := make(chan bool)
	// 计时器
	t := time.NewTimer(time.Duration(sec) * time.Second)
	// 将用户输入结果传入chn
	go func() {
		for {
			s := GetInput()
			if strings.HasPrefix(s, "y") || strings.HasPrefix(s, "Y") {
				chn <- true
			} else if strings.HasPrefix(s, "n") || strings.HasPrefix(s, "N") {
				chn <- false
			} else {
				pterm.Error.Println("无效输入，输入应该为 y 或者 n")
				t.Reset(time.Duration(sec) * time.Second)
			}
		}
	}()
	// 返回输入结果或超时处理
	select {
	case chn := <-chn:
		return chn
	case <-t.C:
		fmt.Println("<超时自动确认>")
		return true
	}
}
