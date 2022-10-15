package utils

import (
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
