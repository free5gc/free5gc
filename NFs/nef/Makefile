BUILD_PATH = build
BIN_PATH = bin
CFG_PATH = config
PWD_PATH = $(shell pwd)

NF = nef

NF_GO_FILES = $(shell find . -name "*.go" ! -name "*_test.go")
NF_CFG_FILE = nefcfg.yaml

VERSION = $(shell git describe --tags)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_HASH = $(shell git log --pretty="%H" -1 | cut -c1-8)
COMMIT_TIME = $(shell git log --pretty="%ai" -1 | awk '{time=$$(1)"T"$$(2)"Z"; print time}')
LDFLAGS = -X github.com/free5gc/version.VERSION=$(VERSION) \
          -X github.com/free5gc/version.BUILD_TIME=$(BUILD_TIME) \
          -X github.com/free5gc/version.COMMIT_HASH=$(COMMIT_HASH) \
          -X github.com/free5gc/version.COMMIT_TIME=$(COMMIT_TIME)

.PHONY: $(NF) clean

.DEFAULT_GOAL: nf

nf: $(NF)

all: $(NF) config

$(NF): $(BUILD_PATH)/$(BIN_PATH)/$(NF)

$(BUILD_PATH)/$(BIN_PATH)/$(NF): cmd/main.go $(NF_GO_FILES)
	@echo "Start building $(NF)...."
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $@ cmd/main.go

config: $(BUILD_PATH)/$(CFG_PATH)/$(NF_CFG_FILE)

$(BUILD_PATH)/$(CFG_PATH)/$(NF_CFG_FILE): $(CFG_PATH)/$(NF_CFG_FILE)
	@echo "Start building $(NF_CFG_FILE)...."
	mkdir -p $(BUILD_PATH)/$(CFG_PATH)
	cp $(CFG_PATH)/$(NF_CFG_FILE) $(BUILD_PATH)/$(CFG_PATH)/$(NF_CFG_FILE)

clean:
	rm -rf $(BUILD_PATH)

