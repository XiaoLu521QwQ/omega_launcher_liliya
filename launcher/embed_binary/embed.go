package embed_binary

const (
	Android_arm64  = "android_arm64"
	Android_x86_64 = "android_x86_64"
	Linux_arm64    = "linux_arm64"
	Linux_x86_64   = "linux_x86_64"
	MACOS_arm64    = "macos_arm64"
	MACOS_x86_64   = "macos_x86_64"
	WINDOWS_arm64  = "windows_arm64"
	WINDOWS_x86_64 = "windows_x86_64"
)

func GetCqHttpBinary() []byte {
	return embedding_cqhttp
}

func GetPlantform() string {
	return PLANTFORM
}
