SET CGO_ENABLED=0

SET GOARCH=amd64
SET GOOS=windows
go build -ldflags="-w -s" -trimpath -o ../build/omega_launcher_windows_amd64.exe
SET GOOS=linux
go build -ldflags="-w -s" -trimpath -o ../build/omega_launcher_linux_amd64
SET GOOS=darwin
go build -ldflags="-w -s" -trimpath -o ../build/omega_launcher_darwin_amd64
SET GOOS=android
go build -ldflags="-w -s" -trimpath -o ../build/omega_launcher_android_amd64

SET GOARCH=arm64
SET GOOS=windows
go build -ldflags="-w -s" -trimpath -o ../build/omega_launcher_windows_arm64.exe
SET GOOS=linux
go build -ldflags="-w -s" -trimpath -o ../build/omega_launcher_linux_arm64
SET GOOS=darwin
go build -ldflags="-w -s" -trimpath -o ../build/omega_launcher_darwin_arm64
SET GOOS=android
go build -ldflags="-w -s" -trimpath -o ../build/omega_launcher_android_arm64
