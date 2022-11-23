package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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
		if strings.HasPrefix(s, "y") || strings.HasPrefix(s, "Y") || s == "" {
			return true
		} else if strings.HasPrefix(s, "n") || strings.HasPrefix(s, "N") {
			return false
		}
		pterm.Error.Println("无效输入，输入应该为 y 或者 n")
	}
}

func GetIntInputInScope(a, b int) int {
	for {
		s := GetInput()
		num, err := strconv.Atoi(s)
		if err != nil {
			pterm.Error.Println("无效输入，请重新输入")
			continue
		}
		if num < a || num > b {
			pterm.Error.Println(fmt.Sprintf("只能输入%d到%d之间的整数，请重新输入", a, b))
			continue
		}
		return num
	}
}

// 错误示范，函数退出后协程仍然在等待输入(((
/*
func GetInputYNTimeLimit(sec int) bool {
	// 建立频道
	chn := make(chan bool)
	// 计时器
	timeout, cancel := context.WithTimeout(context.Background(), time.Duration(sec)*time.Second)
	defer cancel()
	// 将用户输入结果传入chn
	go func() {
		chn <- GetInputYN()
	}()
	// 返回输入结果或超时处理
	select {
	case chn := <-chn:
		return chn
	case <-timeout.Done():
		fmt.Println("<超时自动确认>")
		return true
	}
}
*/
