# NAME     := watchdog-symlinker
# VERSION  := v0.1
# ...

all: setup build

setup:
	which dep
	dep ensure

build:
	go build -o app
