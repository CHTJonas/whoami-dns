SHELL := bash
.ONESHELL:

VER=$(shell git describe --tags --always --dirty)
GO=$(shell which go)
GOOPTS=-trimpath -mod=readonly -ldflags "-X main.version=$(VER:v%=%) -s -w -buildid="
GOMOD=$(GO) mod
GOFMT=$(GO) fmt
GOTEST=$(GO) test $(GOOPTS)
GOBUILD=$(GO) build $(GOOPTS)

all: format build license

clean:
	@rm -rf bin

license: dir
	cp NOTICE bin/NOTICE
	cp LICENSE bin/LICENSE

dir:
	@if [ ! -d bin ]; then mkdir -p bin; fi

mod:
	@$(GOMOD) download

format:
	@$(GOFMT) ./...

test:
	$(GOTEST) -race -cover ./...

build/linux/armv7: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=arm
	export GOARM=7
	$(GOBUILD) -o bin/whoami-dns-linux-$(VER:v%=%)-armv7 cmd/main.go

build/linux/arm64: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=arm64
	$(GOBUILD) -o bin/whoami-dns-linux-$(VER:v%=%)-arm64 cmd/main.go

build/linux/386: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=386
	$(GOBUILD) -o bin/whoami-dns-linux-$(VER:v%=%)-386 cmd/main.go

build/linux/amd64: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=amd64
	$(GOBUILD) -o bin/whoami-dns-linux-$(VER:v%=%)-amd64 cmd/main.go

build/linux: build/linux/armv7 build/linux/arm64 build/linux/386 build/linux/amd64

build/freebsd/armv7: dir mod
	export CGO_ENABLED=0
	export GOOS=freebsd
	export GOARCH=arm
	export GOARM=7
	$(GOBUILD) -o bin/whoami-dns-freebsd-$(VER:v%=%)-armv7 cmd/main.go

build/freebsd/arm64: dir mod
	export CGO_ENABLED=0
	export GOOS=freebsd
	export GOARCH=arm64
	$(GOBUILD) -o bin/whoami-dns-freebsd-$(VER:v%=%)-arm64 cmd/main.go

build/freebsd/386: dir mod
	export CGO_ENABLED=0
	export GOOS=freebsd
	export GOARCH=386
	$(GOBUILD) -o bin/whoami-dns-freebsd-$(VER:v%=%)-386 cmd/main.go

build/freebsd/amd64: dir mod
	export CGO_ENABLED=0
	export GOOS=freebsd
	export GOARCH=amd64
	$(GOBUILD) -o bin/whoami-dns-freebsd-$(VER:v%=%)-amd64 cmd/main.go

build/freebsd: build/freebsd/armv7 build/freebsd/arm64 build/freebsd/386 build/freebsd/amd64

build/openbsd/armv7: dir mod
	export CGO_ENABLED=0
	export GOOS=openbsd
	export GOARCH=arm
	export GOARM=7
	$(GOBUILD) -o bin/whoami-dns-openbsd-$(VER:v%=%)-armv7 cmd/main.go

build/openbsd/arm64: dir mod
	export CGO_ENABLED=0
	export GOOS=openbsd
	export GOARCH=arm64
	$(GOBUILD) -o bin/whoami-dns-openbsd-$(VER:v%=%)-arm64 cmd/main.go

build/openbsd/386: dir mod
	export CGO_ENABLED=0
	export GOOS=openbsd
	export GOARCH=386
	$(GOBUILD) -o bin/whoami-dns-openbsd-$(VER:v%=%)-386 cmd/main.go

build/openbsd/amd64: dir mod
	export CGO_ENABLED=0
	export GOOS=openbsd
	export GOARCH=amd64
	$(GOBUILD) -o bin/whoami-dns-openbsd-$(VER:v%=%)-amd64 cmd/main.go

build/openbsd: build/openbsd/armv7 build/openbsd/arm64 build/openbsd/386 build/openbsd/amd64

build/darwin/arm64: dir mod
	export CGO_ENABLED=0
	export GOOS=darwin
	export GOARCH=arm64
	$(GOBUILD) -o bin/whoami-dns-darwin-$(VER:v%=%)-arm64 cmd/main.go

build/darwin/amd64: dir mod
	export CGO_ENABLED=0
	export GOOS=darwin
	export GOARCH=amd64
	$(GOBUILD) -o bin/whoami-dns-darwin-$(VER:v%=%)-amd64 cmd/main.go

build/darwin: build/darwin/arm64 build/darwin/amd64

build/windows/386: dir mod
	export CGO_ENABLED=0
	export GOOS=windows
	export GOARCH=386
	$(GOBUILD) -o bin/whoami-dns-windows-$(VER:v%=%)-386 cmd/main.go

build/windows/amd64: dir mod
	export CGO_ENABLED=0
	export GOOS=windows
	export GOARCH=amd64
	$(GOBUILD) -o bin/whoami-dns-windows-$(VER:v%=%)-amd64 cmd/main.go

build/windows: build/windows/386 build/windows/amd64

build: build/linux build/freebsd build/openbsd build/darwin build/windows
