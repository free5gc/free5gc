#!/bin/bash
set -o xtrace

# # install all go dependency package
go get -u github.com/aead/cmac
go get -u github.com/antihax/optional
go get -u github.com/bronze1man/radius
go get -u github.com/cydev/zero
go get -u github.com/davecgh/go-spew
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/evanphx/json-patch
go get -u github.com/gin-contrib/sse
go get -u github.com/gin-gonic/gin
go get -u github.com/go-stack/stack
go get -u github.com/golang/protobuf
go get -u github.com/golang/snappy
go get -u github.com/google/go-cmp
go get -u github.com/google/gopacket
go get -u github.com/google/uuid
go get -u github.com/gorilla/mux
go get -u github.com/ishidawataru/sctp
go get -u github.com/jinzhu/copier
go get -u github.com/json-iterator/go
go get -u github.com/konsorten/go-windows-terminal-sequences
go get -u github.com/kr/pretty
go get -u github.com/mattn/go-isatty
go get -u github.com/mitchellh/mapstructure
go get -u github.com/modern-go/concurrent
go get -u github.com/modern-go/reflect2
go get -u github.com/mohae/deepcopy
go get -u github.com/pkg/errors
go get -u github.com/satori/go.uuid
go get -u github.com/sirupsen/logrus
go get -u github.com/stretchr/testify
go get -u github.com/tidwall/pretty
go get -u github.com/ugorji/go
go get -u github.com/urfave/cli
go get -u github.com/xdg/scram
go get -u github.com/xdg/stringprep
go get -u go.mongodb.org/mongo-driver
go get -u golang.org/x/crypto
go get -u golang.org/x/net
go get -u golang.org/x/oauth2
go get -u golang.org/x/sys
go get -u google.golang.org/appengine
go get -u gopkg.in/check.v1
go get -u gopkg.in/go-playground/assert.v1
go get -u gopkg.in/go-playground/validator.v8
go get -u gopkg.in/yaml.v2
go get -u gopkg.in/yaml.v3

# change the version of specific dependency package for some issues in PCF
go get -u github.com/satori/go.uuid
cd $GOPATH/src/github.com/satori/go.uuid
git checkout v1.2.0