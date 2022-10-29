package utils

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func DownloadSmallContent(sourceUrl string) []byte {
	// Get the data
	resp, err := http.Get(sourceUrl)
	if err != nil {
		panic(err)
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
