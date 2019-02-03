FROM golang:1.11 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src/github.com/watchdog-symlinker
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
COPY . .
RUN make