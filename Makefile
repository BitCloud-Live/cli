version := $(shell git describe --abbrev=0 --tags)
LD_FLAGS := -w -X github.com/uvcloud/uv-cli/cmd.version=$(version) -extldflags "-static"
define GOBUILD
	CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build --tags netgo -a -ldflags '$(LD_FLAGS)' -o build/uv-$(version)-$(1)-$(2) -i main.go
endef


all: clean build

clean:
	rm -rf build/*

build: 
	CGO_ENABLED=0 go build --tags netgo -ldflags '$(LD_FLAGS)' -o build/cli -i main.go

run:
	./build/cli

build-linux:
	$(call GOBUILD,linux,amd64)

build-linux-arm:
	$(call GOBUILD,linux,arm)

build-windows:
	$(call GOBUILD,windows,amd64)

build-darwin:
	$(call GOBUILD,darwin,amd64)

build-all: build-linux build-linux-arm build-windows build-darwin
	ls build

docker-build:
	docker build -t uv-cli:$(version) . && docker tag uv-cli:$(version) hub.uvcloud.ir/uvcloud/uv-cli:$(version)

docker-push:
	docker push hub.uvcloud.ir/uvcloud/uv-cli:$(version)
