FROM golang:1.11 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src/github.com/watchdog-symlinker
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep
RUN make