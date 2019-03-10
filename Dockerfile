FROM golang:1.12 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
WORKDIR /go/src/github.com/guitarrapc/watchdog-symlinker
RUN go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
COPY . .
RUN make