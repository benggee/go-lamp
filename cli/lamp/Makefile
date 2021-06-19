version := $(shell /bin/date "+%Y-%m-%d %H:%M")

build:
	go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" lamp.go
	$(if $(shell command -v upx), upx lamp)
mac:
	GOOS=darwin go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o lamp-darwin lamp.go
	$(if $(shell command -v upx), upx lamp-darwin)
win:
	GOOS=windows go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o lamp.exe lamp.go
	$(if $(shell command -v upx), upx lamp.exe)
linux:
	GOOS=linux go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o lamp-linux lamp.go
	$(if $(shell command -v upx), upx lamp-linux)
