.DEFAULT_GOAL := build

APP_NAME := canu

OUTPUT_NAME := _canu

BINDIR := bin

BUILD_PATH = github.com/belitre/canu
PKG_PATH = $(BUILD_PATH)/pkg

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

GIT_SHORT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_TAG := $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "changes" || echo "clean")
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')

TMP_VERSION := canary

BINARY_VERSION := ""

ifndef VERSION
ifeq ($(GIT_DIRTY), clean)
ifdef GIT_TAG
	TMP_VERSION = $(GIT_TAG)
	BINARY_VERSION = $(GIT_TAG)
endif
endif
else
  BINARY_VERSION = $(VERSION)
endif

VERSION ?= $(TMP_VERSION)

LDFLAGS := -extldflags "-static"
LDFLAGS := -X $(PKG_PATH)/version.GitCommit=${GIT_SHORT_COMMIT}+${GIT_DIRTY} -X $(PKG_PATH)/version.BuildDate=${BUILD_DATE}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),"")
	LDFLAGS += -X $(PKG_PATH)/version.Version=$(VERSION)
endif

.PHONY: info
info:
	@echo "Git status:        $(GIT_DIRTY)"
	@echo "Version:           $(VERSION)"
	@echo "Git Tag:           $(GIT_TAG)"
	@echo "Git Commit:        $(GIT_SHORT_COMMIT)"
	@echo "binary:            $(BINARY_VERSION)"

.PHONY: test
test: build
	@echo "test target..."
	@go test ./pkg/... -v -count=1

.PHONY: clean-bin
clean-bin: 
	@rm -rf $(BINDIR)

.PHONY: tidy
tidy:
	@echo "tidy target..."
	@go mod tidy

.PHONY: go-cache
go-cache:
	@echo "go mod download..."
	@go mod download

.PHONY: fmt
fmt:
	@echo "fmt target..."
	@gofmt -l -w -s $(SRC)

build: clean-bin info tidy go-cache fmt 
	@echo "build target..."
	go build -o $(BINDIR)/$(OUTPUT_NAME) -ldflags "$(LDFLAGS)" $(BUILD_PATH)/cmd/$(APP_NAME)

# build cross platform

DIST_DIR := _dist

ARM64_ARCH = arm64
AMD64_ARCH = amd64

GOOSX = darwin
GOOSLINUX = linux

TARGET_DIRS = find * -type d -exec

.PHONY: darwin-arm64
darwin-arm64:
	GOARCH=$(ARM64_ARCH) GOOS=$(GOOSX) go build -o $(DIST_DIR)/$(GOOSX)-$(ARM64_ARCH)/$(OUTPUT_NAME) -ldflags "$(LDFLAGS)" $(BUILD_PATH)/cmd/$(APP_NAME)

.PHONY: darwin-amd64
darwin-amd64:
	GOARCH=$(AMD64_ARCH) GOOS=$(GOOSX) go build -o $(DIST_DIR)/$(GOOSX)-$(AMD64_ARCH)/$(OUTPUT_NAME) -ldflags "$(LDFLAGS)" $(BUILD_PATH)/cmd/$(APP_NAME)

.PHONY: linux-amd64
linux-amd64:
	GOARCH=$(AMD64_ARCH) GOOS=$(GOOSLINUX) go build -o $(DIST_DIR)/$(GOOSLINUX)-$(AMD64_ARCH)/$(OUTPUT_NAME) -ldflags "$(LDFLAGS)" $(BUILD_PATH)/cmd/$(APP_NAME)

.PHONY: clean-dist
clean-dist:
	@rm -rf $(DIST_DIR)

.PHONY: build-cross
build-cross: clean-dist darwin-arm64 darwin-amd64 linux-amd64

.PHONY: dist
dist: build-cross
	( \
		cd $(DIST_DIR) && \
		$(TARGET_DIRS) tar -zcf $(APP_NAME)-${VERSION}-{}.tar.gz -C {} . \; \
	)

# Semantic Release
.PHONY: semantic-release-dependencies
semantic-release-dependencies:
	@npm install --save-dev semantic-release
	@npm install @semantic-release/exec conventional-changelog-conventionalcommits -D

.PHONY: semantic-release 
semantic-release: semantic-release-dependencies
	@npm ci
	@npx semantic-release

.PHONY: semantic-release-ci
semantic-release-ci: semantic-release-dependencies
	@npx semantic-release

.PHONY: semantic-release-dry-run
semantic-release-dry-run: semantic-release-dependencies
	@npm ci
	@npx semantic-release -d
