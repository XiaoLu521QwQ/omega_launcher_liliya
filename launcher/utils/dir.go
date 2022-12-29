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
	// Android环境下, 尝试将数据文件放在 /sdcard/Download
	plantform := embed_binary.GetPlantform()
	if plantform == embed_binary.Android_arm64 || plantform == embed_binary.Android_x86_64 {
		if IsDir("/sdcard/Download/omega_storage") {
			return path.Join("/sdcard/Download")
		} else {
			if IsDir("/sdcard") {
				if MkDir("/sdcard/Download/omega_storage") {
					return path.Join("/sdcard/Download")
				}
			}
		}
	}
	return GetCurrentDir()
}

/*
// From PhoenixBuilder\omega\mainframe\bootstrap.go:227

o.storageRoot = "omega_storage"
// android
if utils.IsDir("/sdcard/Download/omega_storage") {
	o.storageRoot = "/sdcard/Download/omega_storage"
} else {
	if utils.IsDir("/sdcard") {
		if err := utils.MakeDirP("/sdcard/Download/omega_storage"); err == nil {
			o.storageRoot = "/sdcard/Download/omega_storage"
		}
	}
}
*/
