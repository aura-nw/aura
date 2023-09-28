#!/usr/bin/make -f

LEDGER_ENABLED         ?= true

PACKAGES=$(shell go list ./... | grep -v '/simulation')

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# don't override user values
ifeq (,$(VERSION))
	VERSION := $(shell git describe --tags)
	# if VERSION is empty, then populate it with branch's name and raw commit hash
	ifeq (,$(VERSION))
	VERSION := $(BRANCH)-$(COMMIT)
	endif
endif

SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
TM_VERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::')

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
	-X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TM_VERSION) \
        -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(BUILD_TAGS_COMMA_SEP)

ifeq ($(LINK_STATICALLY),true)
        ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(BUILD_TAGS)" -ldflags '$(ldflags)' -trimpath

all: build install

install: go.sum
	@echo "--> Installing aurad"
	@echo "go install -mod=readonly $(BUILD_FLAGS) ./cmd/aurad"
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/aurad

build: go.sum
	@echo "--> Build aurad"
	@go build -mod=readonly $(BUILD_FLAGS) -o ./build/aurad ./cmd/aurad

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

test:
	@go test -mod=readonly $(PACKAGES)

clean:
	@rm -rf build
