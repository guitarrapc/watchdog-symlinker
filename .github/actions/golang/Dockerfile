FROM golang:1.11.4
 
ENV LINT_VERSION="v1.12.5"
COPY entrypoint.sh /entrypoint.sh
 
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin ${LINT_VERSION} \
  && chmod +x /entrypoint.sh
RUN go get -u github.com/golang/dep/cmd/dep

ENTRYPOINT ["/entrypoint.sh"]