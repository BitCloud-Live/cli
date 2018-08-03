version := $(shell git describe --abbrev=0 --tags)

build-linux:
	 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/uvcli-$(version)-linux-amd64 -i main.go
build-linux-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/uvcli-$(version)-linux-arm -i main.go
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/uvcli-$(version)-windows-amd64.exe -i main.go
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/uvcli-$(version)-darwin-amd64 -i main.go
build-all: build-linux build-linux-arm build-windows build-darwin
	ls build
build:
	CGO_ENABLED=0 go build -o build/cli -i main.go

docker-build:
	docker build -t uv-cli:$(version) . && docker tag uv-cli:$(version) hub.uvcloud.ir/uvcloud/uv-cli:$(version)

docker-push:
	docker push hub.uvcloud.ir/uvcloud/uv-cli:$(version)

run:
	./build/cli

clean:
	echo "{}" > ~/.uv/config.json && rm -rf build/*
