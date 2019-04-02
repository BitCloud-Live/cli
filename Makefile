version := $(shell git describe --abbrev=0 --tags)
LD_FLAGS := -w -X github.com/yottab/cli/cmd.version=$(version) -extldflags "-static"
define GOBUILD
	CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build --tags netgo -a -ldflags '$(LD_FLAGS)' -o build/yb-$(version)-$(1)-$(2) -i main.go
endef


all: clean build

clean:
	rm -rf build/*

build: 
	CGO_ENABLED=0 go build --tags netgo -ldflags '$(LD_FLAGS)' -o build/yb -i main.go

run:
	./build/yb

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
	docker build -t yb:$(version) . && docker tag yb:$(version) hub.yottab.io/yottab/cli:$(version)

docker-push:
	docker push hub.yottab.io/yottab/cli:$(version)
