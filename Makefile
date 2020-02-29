version := $(shell git describe --abbrev=0 --tags)
LD_FLAGS := -w -X github.com/yottab/cli/cmd.version=$(version) -extldflags "-static"
define GOBUILD
	$(eval BIN := build/yb-$(version)-$(1)-$(2)$(3))
	CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build --tags netgo -a -ldflags '$(LD_FLAGS)' -o $(BIN) -i main.go
	tar -cjf $(BIN).tar.bz2 $(BIN)
	rm $(BIN)
	openssl dgst -sha256 $(BIN).tar.bz2 | cut -f2 -d' ' > $(BIN).sha256
endef


all: clean build

clean:
	rm -rf build/*

build: 
	CGO_ENABLED=0 go build -a --tags netgo -ldflags '$(LD_FLAGS)' -o build/yb -i main.go

run:
	./build/yb

build-linux:
	$(call GOBUILD,linux,amd64,)

build-linux-arm:
	$(call GOBUILD,linux,arm,)

build-windows:
	$(call GOBUILD,windows,amd64,.exe)

build-darwin:
	$(call GOBUILD,darwin,amd64,)

build-all: build-linux build-linux-arm build-windows build-darwin
	ls build

docker-build:
	docker build -t yb:$(version) . && docker tag yb:$(version) hub.yottab.io/yottab-library/cli:$(version)

docker-push:
	docker push hub.yottab.io/yottab-library/cli:$(version)


update-deps: download tidy vendor

download:
	GONOSUMDB="gitlab.com/projectxchange,github.com/btcsuite" HTTPS_PROXY=socks5://127.0.0.1:9150 GO111MODULE=on go mod download

tidy:
	GONOSUMDB="gitlab.com/projectxchange,github.com/btcsuite" HTTPS_PROXY=socks5://127.0.0.1:9150 GO111MODULE=on go mod tidy -v

vendor:
	GONOSUMDB="gitlab.com/projectxchange,github.com/btcsuite" HTTPS_PROXY=socks5://127.0.0.1:9150 GO111MODULE=on go mod vendor
