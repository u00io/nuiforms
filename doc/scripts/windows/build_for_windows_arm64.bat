echo Building for Windows ARM64...
SET GOOS=windows
SET GOARCH=arm64
go build -ldflags -H=windowsgui -o ../../../bin/windows_arm64/nuiforms.exe ../../../main.go
