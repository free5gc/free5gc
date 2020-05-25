#!/bin/bash
set -o xtrace

# Install go packages which are not required to switch version
go get -u github.com/aead/cmac
go get -u github.com/bronze1man/radius
go get -u github.com/cydev/zero
go get -u github.com/jinzhu/copier
go get -u github.com/mohae/deepcopy
go get -u github.com/sparrc/go-ping
go get -u github.com/xdg/scram
go get -u golang.org/x/crypto
go get -u golang.org/x/net
go get -u golang.org/x/oauth2
go get -u golang.org/x/sys
go get -u golang.org/x/tools
go get -u google.golang.org/genproto
go get -u gopkg.in/yaml.v3

# Install go packages which are required to switch version
go get -u git.cs.nctu.edu.tw/calee/sctp
cd $GOPATH/src/git.cs.nctu.edu.tw/calee/sctp
git checkout v1.0.0
go get -u github.com/antihax/optional
cd $GOPATH/src/github.com/antihax/optional
git checkout v1.0.0
go get -u github.com/davecgh/go-spew
cd $GOPATH/src/github.com/davecgh/go-spew
git checkout v1.1.1
go get -u github.com/dgrijalva/jwt-go
cd $GOPATH/src/github.com/dgrijalva/jwt-go
git checkout v3.2.0
go get -u github.com/evanphx/json-patch
cd $GOPATH/src/github.com/evanphx/json-patch
git checkout v4.5.0
go get -u github.com/gin-contrib/cors
cd $GOPATH/src/github.com/gin-contrib/cors
git checkout v1.3.0
go get -u github.com/gin-contrib/sse
cd $GOPATH/src/github.com/gin-contrib/sse
git checkout v0.1.0
go get -u github.com/gin-gonic/gin
cd $GOPATH/src/github.com/gin-gonic/gin
git checkout v1.4.0
go get -u github.com/go-stack/stack
cd $GOPATH/src/github.com/go-stack/stack
git checkout v1.8.0
go get -u github.com/golang/snappy
cd $GOPATH/src/github.com/golang/snappy
git checkout v0.0.1
go get -u github.com/google/go-cmp
cd $GOPATH/src/github.com/google/go-cmp
git checkout v0.4.0
go get -u github.com/google/gopacket
cd $GOPATH/src/github.com/google/gopacket
git checkout v1.1.17
go get -u github.com/google/uuid
cd $GOPATH/src/github.com/google/uuid
git checkout v1.1.1
go get -u github.com/konsorten/go-windows-terminal-sequences
cd $GOPATH/src/github.com/konsorten/go-windows-terminal-sequences
git checkout v1.0.2
go get -u github.com/mattn/go-isatty
cd $GOPATH/src/github.com/mattn/go-isatty
git checkout v0.0.8
go get -u github.com/mitchellh/mapstructure
cd $GOPATH/src/github.com/mitchellh/mapstructure
git checkout v1.1.2
go get -u github.com/pkg/errors
cd $GOPATH/src/github.com/pkg/errors
git checkout v0.9.1
go get -u github.com/sirupsen/logrus
cd $GOPATH/src/github.com/sirupsen/logrus
git checkout v1.4.2
go get -u github.com/stretchr/testify
cd $GOPATH/src/github.com/stretchr/testify
git checkout v1.4.0
go get -u github.com/tidwall/pretty
cd $GOPATH/src/github.com/tidwall/pretty
git checkout v1.0.0
go get -u github.com/ugorji/go
cd $GOPATH/src/github.com/ugorji/go
git checkout v1.1.7
go get -u github.com/urfave/cli
cd $GOPATH/src/github.com/urfave/cli
git checkout v1.20.0
go get -u github.com/vishvananda/netlink
cd $GOPATH/src/github.com/vishvananda/netlink
git checkout v1.1.0
go get -u github.com/wmnsk/go-gtp
cd $GOPATH/src/github.com/wmnsk/go-gtp
git checkout v0.7.4
go get -u github.com/xdg/stringprep
cd $GOPATH/src/github.com/xdg/stringprep
git checkout v1.0.0
go get -u go.mongodb.org/mongo-driver
cd $GOPATH/src/go.mongodb.org/mongo-driver
git checkout v1.0.2
go get -u golang.org/x/mod
cd $GOPATH/src/golang.org/x/mod
git checkout v0.2.0
go get -u google.golang.org/appengine
cd $GOPATH/src/google.golang.org/appengine
git checkout v1.6.1
go get -u gopkg.in/go-playground/validator.v8
cd $GOPATH/src/gopkg.in/go-playground/validator.v8
git checkout v8.18.2
go get -u gopkg.in/yaml.v2
cd $GOPATH/src/gopkg.in/yaml.v2
git checkout v2.2.8
