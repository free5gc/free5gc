GO_BIN_PATH = bin
GO_SRC_PATH = NFs
C_BUILD_PATH = build
ROOT_PATH = $(shell pwd)

NF = $(GO_NF)
GO_NF = amf ausf bsf nrf nssf pcf smf udm udr n3iwf upf chf tngf nef scp

WEBCONSOLE = webconsole

CFG_GENERATOR = cfg-generator
CFG_GENERATOR_PATH = config/cfg-generator
CFG_GENERATOR_GO_FILES = $(shell find $(CFG_GENERATOR_PATH) -name "*.go" ! -name "*_test.go") $(CFG_GENERATOR_PATH)/go.mod

NF_GO_FILES = $(shell find $(GO_SRC_PATH)/$(%) -name "*.go" ! -name "*_test.go")
WEBCONSOLE_GO_FILES = $(shell find $(WEBCONSOLE) -name "*.go" ! -name "*_test.go")
WEBCONSOLE_JS_FILES = $(shell find $(WEBCONSOLE)/frontend -name '*.tsx' ! -path "*/node_modules/*")
WEBCONSOLE_FRONTEND = $(WEBCONSOLE)/public

VERSION = $(shell git describe --tags)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_HASH = $(shell git submodule status | grep $(GO_SRC_PATH)/$(@F) | awk '{print $$(1)}' | cut -c1-8)
COMMIT_TIME = $(shell cd $(GO_SRC_PATH)/$(@F) && git log --pretty="@%at" -1 | xargs date -u +"%Y-%m-%dT%H:%M:%SZ" -d)
LDFLAGS = -X github.com/free5gc/util/version.VERSION=$(VERSION) \
          -X github.com/free5gc/util/version.BUILD_TIME=$(BUILD_TIME) \
          -X github.com/free5gc/util/version.COMMIT_HASH=$(COMMIT_HASH) \
          -X github.com/free5gc/util/version.COMMIT_TIME=$(COMMIT_TIME)

WEBCONSOLE_COMMIT_HASH = $(shell git submodule status | grep $(WEBCONSOLE) | awk '{print $$(1)}' | cut -c1-8)
WEBCONSOLE_COMMIT_TIME = $(shell cd $(WEBCONSOLE) && git log --pretty="@%at" -1 | xargs date -u +"%Y-%m-%dT%H:%M:%SZ" -d)
WEBCONSOLE_LDFLAGS = -X github.com/free5gc/util/version.VERSION=$(VERSION) \
                     -X github.com/free5gc/util/version.BUILD_TIME=$(BUILD_TIME) \
                     -X github.com/free5gc/util/version.COMMIT_HASH=$(WEBCONSOLE_COMMIT_HASH) \
                     -X github.com/free5gc/util/version.COMMIT_TIME=$(WEBCONSOLE_COMMIT_TIME)

.PHONY: $(NF) $(WEBCONSOLE) $(CFG_GENERATOR) clean

.DEFAULT_GOAL: nfs

nfs: $(NF) $(CFG_GENERATOR)

all: $(NF) $(CFG_GENERATOR) $(WEBCONSOLE)

debug: GCFLAGS += -N -l
debug: all

$(GO_NF): % : $(GO_BIN_PATH)/%

$(GO_BIN_PATH)/%: $(NF_GO_FILES)
# $(@F): The file-within-directory part of the file name of the target.
	@echo "Start building $(@F)...."
	cd $(GO_SRC_PATH)/$(@F)/cmd && \
	CGO_ENABLED=0 go build -gcflags "$(GCFLAGS)" -ldflags "$(LDFLAGS)" -o $(ROOT_PATH)/$@ main.go

vpath %.go $(addprefix $(GO_SRC_PATH)/, $(GO_NF))

# cfg-generator renders config/<nf>cfg.yaml from config/global-config.yaml;
# run.sh calls it before starting the NFs. The explicit rule below takes
# precedence over the $(GO_BIN_PATH)/% pattern rule used for the NFs.
$(CFG_GENERATOR): $(GO_BIN_PATH)/$(CFG_GENERATOR)

$(GO_BIN_PATH)/$(CFG_GENERATOR): $(CFG_GENERATOR_GO_FILES)
	@echo "Start building $(@F)...."
	cd $(CFG_GENERATOR_PATH) && \
	CGO_ENABLED=0 go build -o $(ROOT_PATH)/$@ .

$(C_NF): % :
	@echo "Start building $@...."
	cd $(GO_SRC_PATH)/$@ && \
	rm -rf $(C_BUILD_PATH) && \
	mkdir -p $(C_BUILD_PATH) && \
	cd ./$(C_BUILD_PATH) && \
	cmake .. && \
	make -j$(nproc)

$(WEBCONSOLE): $(WEBCONSOLE)/$(GO_BIN_PATH)/$(WEBCONSOLE) $(WEBCONSOLE_FRONTEND)

$(WEBCONSOLE)/$(GO_BIN_PATH)/$(WEBCONSOLE): $(WEBCONSOLE)/server.go $(WEBCONSOLE_GO_FILES)
	@echo "Start building $(@F)...."
	cd $(WEBCONSOLE) && \
	CGO_ENABLED=0 go build -ldflags "$(WEBCONSOLE_LDFLAGS)" -o $(ROOT_PATH)/$@ ./server.go

$(WEBCONSOLE_FRONTEND): $(WEBCONSOLE_JS_FILES)
	@echo "Start building $(@F) frontend...."
	cd $(WEBCONSOLE)/frontend && \
	sudo corepack enable && \
	yarn install && \
	yarn build && \
	rm -rf ../public && \
	cp -R build ../public

clean:
	rm -rf $(addprefix $(GO_BIN_PATH)/, $(GO_NF) $(CFG_GENERATOR))
	rm -rf $(addprefix $(GO_SRC_PATH)/, $(addsuffix /$(C_BUILD_PATH), $(C_NF)))
	rm -rf $(WEBCONSOLE)/$(GO_BIN_PATH)/$(WEBCONSOLE)

