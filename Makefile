#!/usr/bin/make -f

LEDGER_ENABLED         ?= true

PACKAGES=$(shell go list ./... | grep -v '/simulation')

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# library versions
LIBWASM_VERSION = $(shell go list -m -f '{{ .Version }}' github.com/CosmWasm/wasmvm)

# docker 
DOCKER := $(shell which docker)

# don't override user values
ifeq (,$(VERSION))
	VERSION := $(shell git describe --tags)
	# if VERSION is empty, then populate it with branch's name and raw commit hash
	ifeq (,$(VERSION))
	VERSION := $(BRANCH)-$(COMMIT)
	endif
endif

BFT_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')

BUILD_TAGS += netgo
BUILD_TAGS := $(strip $(BUILD_TAGS))
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      BUILD_TAGS += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        BUILD_TAGS += ledger
      endif
    endif
  endif
endif

# Flags
WHITESPACE := $(subst ,, )
COMMA := ,
BUILD_TAGS_COMMA_SEP := $(subst $(WHITESPACE),$(COMMA),$(BUILD_TAGS))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=aura \
	-X github.com/cosmos/cosmos-sdk/version.AppName=aurad \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X github.com/cometbft/cometbft/version.TMCoreSemVer=$(BFT_VERSION) \
  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(BUILD_TAGS_COMMA_SEP)

ifeq ($(LINK_STATICALLY),true)
        ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(BUILD_TAGS)" -ldflags '$(ldflags)' -trimpath

# Go version
GO_MAJOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
MINIMUM_SUPPORTED_GO_MAJOR_VERSION = 1
MINIMUM_SUPPORTED_GO_MINOR_VERSION = 19
GO_VERSION_VALIDATION_ERR_MSG = Golang version is not supported, please update to at least $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION).$(MINIMUM_SUPPORTED_GO_MINOR_VERSION)

GORELEASER_VERSION = v1.20.0

all: build install

install: validate-go-version go.sum
	@echo "--> Installing aurad"
	@echo "go install -mod=readonly $(BUILD_FLAGS) ./cmd/aurad"
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/aurad

build: validate-go-version go.sum
	@echo "--> Build aurad"
	@go build -mod=readonly $(BUILD_FLAGS) -o ./build/aurad ./cmd/aurad

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

test:
	@go test -mod=readonly $(PACKAGES)

clean:
	@rm -rf build

validate-go-version: ## Validates the installed version of go against Mattermost's minimum requirement.
	@if [ $(GO_MAJOR_VERSION) -gt $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION) ]; then \
		exit 0 ;\
	elif [ $(GO_MAJOR_VERSION) -lt $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION) ]; then \
		echo '$(GO_VERSION_VALIDATION_ERR_MSG)';\
		exit 1; \
	elif [ $(GO_MINOR_VERSION) -lt $(MINIMUM_SUPPORTED_GO_MINOR_VERSION) ] ; then \
		echo '$(GO_VERSION_VALIDATION_ERR_MSG)';\
		exit 1; \
	fi

release:
	git tag $(VERSION)
	git push origin $(VERSION)
	$(DOCKER) run \
		--rm \
		-e LIBWASM_VERSION=$(LIBWASM_VERSION) \
		-e PRE_RELEASE=$(PRE_RELEASE) \
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-e VERSION="$(VERSION)" \
		-e COMMIT="$(COMMIT)" \
		-e BFT_VERSION="$(BFT_VERSION)" \
		-e PRE_RELEASE="$(PRE_RELEASE)" \
		-e BUILD_TAGS_COMMA_SEP="$(BUILD_TAGS_COMMA_SEP)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/github.com/aura-nw/aura \
		-w /go/src/github.com/aura-nw/aura \
		ghcr.io/goreleaser/goreleaser:$(GORELEASER_VERSION) \
		--clean --skip-validate