#!/bin/bash

GO_VERSION="1.26.2"
GOLANGCI_LINT_VERSION="2.11.4"

# this assumes your current version of Go is in the default location:
sudo rm -rf /usr/local/go
wget https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz
sudo tar -C /usr/local -zxvf go${GO_VERSION}.linux-amd64.tar.gz
rm go${GO_VERSION}.linux-amd64.tar.gz

# this assumes golangci-lint is in your PATH:
if golangci-lint version > /dev/null 2>&1; then
    curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v$GOLANGCI_LINT_VERSION
fi
