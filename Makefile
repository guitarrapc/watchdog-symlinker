# NAME     := watchdog-symlinker
# VERSION  := v0.1
# ...

VERSION  := v0.0.3
REVISION := $(shell git rev-parse --short HEAD)
SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

all: setup test build

setup:
	go mod download

test:
    # go test -v ./..
	go test
	golangci-lint run

build:
	go build -o app
