.PHONY: help

BINARY=fairyla
LINUX=$(BINARY).linux-amd64
MACOS=$(BINARY).darwin-amd64
WIN=$(BINARY).windows-amd64.exe
VERSION=$(shell grep "const version" main.go | tr -d '"' | awk '{print $$NF}')
LDFLAGS=-ldflags "-s -w"

help:
	@echo "  make clean  - Remove binaries and vim swap files"
	@echo "  make gotool - Run go tool 'fmt' and 'vet'"
	@echo "  make build  - Compile go code and generate binary file"
	@echo "  make release- Format go code and compile to generate binary release"

gotool:
	go fmt ./
	go vet ./

build:
	go build ${LDFLAGS} -o bin/$(BINARY) && chmod +x bin/$(BINARY)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/$(LINUX) && chmod +x bin/$(LINUX)

build-macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/$(MACOS) && chmod +x bin/$(MACOS)

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/$(WIN) && chmod +x bin/$(WIN)

docker:
	docker build -t staugur/fairyla .

clean:
	find bin -name '*.tar.gz' -exec rm -f {} +
	find bin -name '*.zip' -exec rm -f {} +
	find bin -name 'fairyla.*-amd64*' -exec rm -f {} +
