GO_BIN_PATH = bin

VERSION = $(shell git describe --tags)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_HASH = $(shell git submodule status | grep $(GO_SRC_PATH)/$(@F) | awk '{print $$(1)}' | cut -c1-8)
COMMIT_TIME = $(shell git log --pretty="@%at" -1 | xargs date -u +"%Y-%m-%dT%H:%M:%SZ" -d)
LDFLAGS = -X github.com/free5gc/util/version.VERSION=$(VERSION) \
          -X github.com/free5gc/util/version.BUILD_TIME=$(BUILD_TIME) \
          -X github.com/free5gc/util/version.COMMIT_HASH=$(COMMIT_HASH) \
          -X github.com/free5gc/util/version.COMMIT_TIME=$(COMMIT_TIME)

NF = nf
NF_GO_FILES = $(shell find -name "*.go" ! -name "*_test.go")

debug: GCFLAGS += -N -l

$(NF): $(GO_BIN_PATH)/$(NF)

$(GO_BIN_PATH)/$(NF): cmd/main.go $(NF_GO_FILES)
	@echo "Start building $(@F)...."
	CGO_ENABLED=0 go build -gcflags "$(GCFLAGS)" -ldflags "$(LDFLAGS)" -o $@ ./cmd/main.go

clean:
	rm -rf $(GO_BIN_PATH)/$(NF)