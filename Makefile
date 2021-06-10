NAME 			      := mintscan-union
VERSION               := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT                := $(shell git log -1 --format='%H')
DESTDIR         	  ?= $(GOPATH)/bin/${NAME}
BUILD_FLAGS 		  := -ldflags "-w -s \
	-X github.com/cosmostation/cosmostation-cosmos/chain-exporter/exporter.Version=${VERSION} \
	-X github.com/cosmostation/cosmostation-cosmos/chain-exporter/exporter.Commit=${COMMIT}"

## Show all make target commands.
help:
	@make2help $(MAKEFILE_LIST)

## Print out application information.
version: 
	@echo "NAME: ${NAME}"
	@echo "VERSION: ${VERSION}"
	@echo "COMMIT: ${COMMIT}"	

## Build an executable file in $DESTDIR directory.
build: build_exporter build_mintscan

build_exporter: go.sum
	@echo "-> Building chain-exporter"
	@go build -mod=readonly $(BUILD_FLAGS) -o . ./cmd/chain-exporter

build_mintscan: go.sum
	@echo "-> Building mintscan"
	@go build -mod=readonly $(BUILD_FLAGS) -o . ./cmd/mintscan

## Install executable file in $GOBIN direcotry. 
install: install_exporter install_mintscan

install_exporter: go.sum
	@echo "-> Installing chain-exporter"
	@go install -mod=readonly ${BUILD_FLAG} ./cmd/chain-exporter

install_mintscan: go.sum
	@echo "-> Installing mintscan"
	@go install -mod=readonly ${BUILD_FLAG} ./cmd/mintscan


## Clean the executable file.
clean:
	@echo "-> Cleaning ${NAME} binary..."
	rm -f $(DESTDIR) 2> /dev/null

PHONY: build install make_service enable_service clean
