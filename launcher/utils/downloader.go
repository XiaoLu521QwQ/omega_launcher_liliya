package utils

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/pterm/pterm"
)

func DownloadSmallContent(sourceUrl string) []byte {
	// Get the data
	resp, err := http.Get(sourceUrl)
	if err != nil {
		if strings.Contains(err.Error(), "->[::1]:53") {
			pterm.Warning.Println("域名解析出错，建议为终端配置DNS")
			net.DefaultResolver.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
				dialer := &net.Dialer{
					Timeout: 10 * time.Second,
				}
				return dialer.DialContext(ctx, "tcp", "8.8.8.8:53")
			}
			resp, err = http.Get(sourceUrl)
			if err != nil {
				pterm.Error.Println("从指定仓库下载资源时出现错误，请重试或更换仓库")
				panic(err)
			}
		} else {
			pterm.Error.Println("从指定仓库下载资源时出现错误，请重试或更换仓库")
			panic(err)
		}
	}
	defer resp.Body.Close()
	// Size
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	downloadSize := int64(size)

	// Progress Bar
	bar := pb.Full.Start64(downloadSize)
	bar.SetWidth(-1)
	bar.SetMaxWidth(100)
	bar.SetRefreshRate(time.Millisecond)
	defer bar.Finish()

	// Reader
	barReader := bar.NewProxyReader(resp.Body)

	// Buffer
	contents := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(contents, barReader); err == nil {
		return contents.Bytes()
	} else {
		panic(err)
	}
}
