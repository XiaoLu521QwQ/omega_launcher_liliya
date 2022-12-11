SET CGO_ENABLED=0
SET GOARCH=amd64
SET GOOS=windows
go build -ldflags="-w -s" -o ../build/launcher.exe
SET GOOS=linux
go build -ldflags="-w -s" -o ../build/launcher_linux
SET GOOS=darwin
go build -ldflags="-w -s" -o ../build/launcher_darwin
SET GOOS=android
SET GOARCH=arm64
go build -ldflags="-w -s" -o ../build/launcher_android

:: upx.exe -9 ../build/launcher.exe
:: upx.exe -9 ../build/launcher_linux
:: upx.exe -9 ../build/launcher_darwin
:: upx.exe -9 ../build/launcher_android
