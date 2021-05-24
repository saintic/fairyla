.PHONY: help

BINARY=fairyla
LINUX=$(BINARY).linux-amd64
MACOS=$(BINARY).darwin-amd64
WIN=$(BINARY).windows-amd64.exe
VERSION=$(shell grep "const version" server/main.go | tr -d '"' | awk '{print $$NF}')

help:
	@echo "  make docker  - build container image fastly"
	@echo "  make unpack  - format, compile go code and ui"
	@echo "  make release - unpack, update version and commit"

docker:
	docker build -t staugur/fairyla .

unpack:
	cd client && yarn build
	cp -f NOTICE server/
	cd server && make gotool && make build-linux && make build-macos && make build-windows
	cd server/bin && mv $(LINUX) $(BINARY)
	cd server && tar zcvf bin/$(BINARY).$(VERSION)-linux-amd64.tar.gz bin/$(BINARY) ui NOTICE && rm bin/$(BINARY)
	cd server/bin && mv $(MACOS) $(BINARY)
	cd server && tar zcvf bin/$(BINARY).$(VERSION)-darwin-amd64.tar.gz bin/$(BINARY) ui NOTICE && rm bin/$(BINARY)
	cd server/bin && mv $(WIN) $(BINARY).exe
	cd server && zip -r bin/$(BINARY).$(VERSION)-windows-amd64.zip bin/$(BINARY).exe ui NOTICE && rm bin/$(BINARY).exe
	rm -f server/NOTICE

release: unpack
	git add . && git ci -m "bump version $(VERSION)"

