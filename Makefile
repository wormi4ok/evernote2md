LATEST_TAG := $(shell git describe --abbrev=0 --tags)
SUFFIX ?= rc1

.PHONY: build
build:
	 go build -trimpath -o bin/evernote2md -ldflags "-s -w -X main.version=$(LATEST_TAG)-${SUFFIX}"

.PHONY: test
test:
	go test ./...
