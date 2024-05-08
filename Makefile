NAMES=iamtf josso
VERSION=$(shell git describe --tags --always --dirty)

PLATFORMS=darwin linux windows openbsd
ARCHITECTURES=amd64 386 arm64 arm

default: build

dep: # Download required dependencies
	go mod tidy
	go mod vendor

build:
	CGO_ENABLED=0 go build -o ./tmp/wazuh-jumpcloud-integration -a -ldflags '-X main.version=$(VERSION) -extldflags "-static"' ./cmd/service/main.go

#CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' .


run: build
	~/go/bin/custom-jumpcloud

test: build
	go test ./... -v

dist:
	$(foreach NAME,$(NAMES),\
		$(foreach GOOS,$(PLATFORMS),\
			$(foreach GOARCH,$(ARCHITECTURES),\
				$(shell export GOOS=$(GOOS);\
					export BINARY=$(NAME)ctl;\
					export GOARCH=$(GOARCH);\
					OUT_DIR='./.tmp/$(NAME)/$(GOOS)/$(GOARCH)/$(VERSION)';\
					go build -ldflags "-X main.version=$(VERSION)" -v -o $${OUT_DIR%.}/$(BINARY) ./$(NAME)ctl; \
					if test -f $${OUT_DIR}/$${BINARY} ; then cd $${OUT_DIR} ; zip -q ../../../$${BINARY}-$(GOOS)-$(GOARCH)-$(VERSION).zip $${BINARY} ; fi; \
					if test -f $${OUT_DIR}/$${BINARY}.exe ; then cd $${OUT_DIR} ; zip -q ../../../$${BINARY}-$(GOOS)-$(GOARCH)-$(VERSION).zip $${BINARY}.exe ; fi \
					))))
