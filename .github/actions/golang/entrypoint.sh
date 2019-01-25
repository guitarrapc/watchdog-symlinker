#!/bin/bash
 
APP_DIR="/go/src/github.com/${GITHUB_REPOSITORY}/"
 
mkdir -p ${APP_DIR} && cp -r ./ ${APP_DIR} && cd ${APP_DIR}
 
export GO111MODULE=on
go mod tidy
go mod verify
 
if [[ "$1" == "lint" ]]; then
    echo "# Running GolangCI-Lint"
    golangci-lint --version
    golangci-lint run
fi
 
if [[ "$1" == "build" ]]; then
    echo "# Running Build"
    go build
fi