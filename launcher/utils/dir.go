package utils

import (
	"omega_launcher/embed_binary"
	"os"
	"path"
	"path/filepath"
)

func GetCurrentDir() string {
	// 兼容配套的Docker
	if IsFile(path.Join("/ome", "launcher_liliya")) {
		return path.Join("/workspace")
	}
	pathExecutable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(pathExecutable)
}

func GetCurrentDataDir() string {
	// 兼容配套的Docker
	if IsFile(path.Join("/ome", "launcher_liliya")) {
		return path.Join("/workspace")
	}
	// 兼容Android
	plantform := embed_binary.GetPlantform()
	if plantform == embed_binary.Android_arm64 || plantform == embed_binary.Android_x86_64 {
		return path.Join("/sdcard", "Download")
	}
	pathExecutable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(pathExecutable)
}
