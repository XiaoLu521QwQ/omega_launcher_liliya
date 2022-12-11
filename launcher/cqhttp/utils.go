package cqhttp

import (
	"bytes"
	"io"
	"net/url"
	"omega_launcher/embed_binary"
	"omega_launcher/utils"
	"path"
	"path/filepath"

	"github.com/andybalholm/brotli"
	"github.com/gorilla/websocket"
)

func GetCqHttpExec() string {
	cqhttp := "cqhttp"
	if embed_binary.GetPlantform() == embed_binary.WINDOWS_x86_64 {
		cqhttp = "cqhttp.exe"
	}
	cqhttp = path.Join(utils.GetCurrentDir(), cqhttp)
	p, err := filepath.Abs(cqhttp)
	if err != nil {
		panic(err)
	}
	return p
}

func GetCqHttpBinary() []byte {
	compressedData := embed_binary.GetCqHttpBinary()
	var execBytes []byte
	var err error
	if execBytes, err = io.ReadAll(brotli.NewReader(bytes.NewReader(compressedData))); err != nil {
		panic(err)
	}
	return execBytes
}

func GetCQHttpHash() string {
	exec := GetCqHttpExec()
	return utils.GetFileHash(exec)
}

func GetEmbeddedCQHttpHash() string {
	return utils.GetBinaryHash(GetCqHttpBinary())
}

func WaitConnect(addr string) {
	for {
		u := url.URL{Scheme: "ws", Host: addr}
		if _, _, err := websocket.DefaultDialer.Dial(u.String(), nil); err != nil {
			// time.Sleep(1)
			continue
		} else {
			return
		}
	}
}
