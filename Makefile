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

install: go.sum
	@echo "--> Installing aurad"
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