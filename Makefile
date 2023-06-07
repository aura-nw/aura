#!/usr/bin/make -f

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

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=aura \
	-X github.com/cosmos/cosmos-sdk/version.AppName=aurad \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TM_VERSION)

BUILD_FLAGS := -ldflags '$(ldflags)'

all: build install

install: check-go-version go.sum
	@echo "--> Installing aurad"
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/aurad

build: check-go-version go.sum
	@echo "--> Build aurad"
	@go build -mod=readonly $(BUILD_FLAGS) -o ./build/aurad ./cmd/aurad

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

test: check-go-version
	@go test -mod=readonly $(PACKAGES)

clean:
	@rm -rf build

# Add check to make sure we are using the proper Go version before proceeding with anything
check-go-version:
	@if ! go version | grep -q "go1.19"; then \
		echo "\033[0;31mERROR:\033[0m Go version 1.19 is required for compiling aurad. It looks like you are using" "$(shell go version) \nThere are potential consensus-breaking changes that can occur when running binaries compiled with different versions of Go. Please download Go version 1.19 and retry. Thank you!"; \
		exit 1; \
	fi