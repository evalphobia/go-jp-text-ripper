.PHONY: init build _build clean lint test

GO111MODULE=off
VERSION := $(shell git tag --points-at HEAD --sort=-v:refname | head -n 1)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'

init:
	go mod download

# build binary
build:
	go build -ldflags "$(LDFLAGS)" -o bin/go-jp-text-ripper .

build-macos:
	@make _build BUILD_OS=darwin BUILD_ARCH=amd64

build-linux:
	@make _build BUILD_OS=linux BUILD_ARCH=amd64

build-windows:
	@make _build BUILD_OS=windows BUILD_ARCH=amd64

_build:
	@mkdir -p bin/release
	$(eval BUILD_OUTPUT := go-jp-text-ripper_${BUILD_OS}_${BUILD_ARCH}${BUILD_ARM})
	GOOS=${BUILD_OS} \
	GOARCH=${BUILD_ARCH} \
	go build -o bin/${BUILD_OUTPUT} .
	@if [ "${USE_ARCHIVE}" = "1" ]; then \
		gzip -k -f bin/${BUILD_OUTPUT} ;\
		mv bin/${BUILD_OUTPUT}.gz bin/release/ ;\
	fi

build-all: clean
	@make build-macos build-linux build-windows USE_ARCHIVE=1

clean:
	rm -f bin/go-jp-text-ripper_*
	rm -f bin/release/*


# Exec golint, vet, gofmt
lint:
	@type golangci-lint > /dev/null || go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

test:
	@type gosec > /dev/null || go get github.com/securego/gosec/cmd/gosec
	gosec -quiet ./...
	go test ./... -count=1;
