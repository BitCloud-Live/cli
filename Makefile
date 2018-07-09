build-linux:
	 GOOS=linux GOARCH=amd64 go build -o build/cli-linux-amd64 -i main.go
build-linux-arm:
	GOOS=linux GOARCH=amd64 go build -o build/cli-linux-arm -i main.go
build-windows:
	GOOS=windows GOARCH=amd64 go build -o build/uvcli-windows-amd64.exe -i main.go
build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o build/uvcli-osx-amd64 -i main.go
build-all: build-linux build-linux-arm build-windows build-darwin
	ls build

build:
	go build -o build/cli -i main.go
run:
	./build/cli

clean:
	echo "{}" > ~/.uv/config.json
