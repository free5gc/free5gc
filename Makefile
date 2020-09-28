GO_BIN_PATH = bin
GO_SRC_PATH = src
C_BUILD_PATH = build

NF = $(GO_NF) $(C_NF)
GO_NF = amf ausf nrf nssf pcf smf udm udr n3iwf
C_NF = upf

NF_GO_FILES = $(shell find $(GO_SRC_PATH)/$(%) -name "*.go" ! -name "*_test.go")

VERSION = $(shell git describe --tags)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_HASH = $(shell git submodule status | grep $(GO_SRC_PATH)/$(@F) | awk '{print $$(1)}' | cut -c1-8)
COMMIT_TIME = $(shell cd $(GO_SRC_PATH)/$(@F) && git log --pretty="%ai" -1 | awk '{time=$$(1)"T"$$(2)"Z"; print time}')
LDFLAGS = -X free5gc/src/$(@F)/version.VERSION=$(VERSION) \
          -X free5gc/src/$(@F)/version.BUILD_TIME=$(BUILD_TIME) \
          -X free5gc/src/$(@F)/version.COMMIT_HASH=$(COMMIT_HASH) \
          -X free5gc/src/$(@F)/version.COMMIT_TIME=$(COMMIT_TIME)


.PHONY: $(NF) clean

all: $(NF)

$(GO_NF): % : $(GO_BIN_PATH)/%

$(GO_BIN_PATH)/%: %.go $(NF_GO_FILES)
# $(@F): The file-within-directory part of the file name of the target.
	@echo "Start building $(@F)...."
	go build -ldflags "$(LDFLAGS)" -o $@ $<

vpath %.go $(addprefix $(GO_SRC_PATH)/, $(GO_NF))

$(C_NF): % :
	@echo "Start building $@...."
	cd $(GO_SRC_PATH)/$@ && \
	rm -rf $(C_BUILD_PATH) && \
	mkdir -p $(C_BUILD_PATH) && \
	cd ./$(C_BUILD_PATH) && \
	cmake .. && \
	make -j$(nproc)

clean:
	rm -rf $(addprefix $(GO_BIN_PATH)/, $(GO_NF))
	rm -rf $(addprefix $(GO_SRC_PATH)/, $(addsuffix /$(C_BUILD_PATH), $(C_NF)))

