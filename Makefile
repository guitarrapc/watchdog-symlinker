# NAME     := watchdog-symlinker
# VERSION  := v0.1
# ...

all: setup test build

setup:
	which dep
	dep ensure

test:
    # go test -v ./..
	golangci-lint run

build:
	go build -o app
