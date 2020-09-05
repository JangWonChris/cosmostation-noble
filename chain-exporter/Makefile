NAME 			      := chain-exporter
VERSION               := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT                := $(shell git log -1 --format='%H')
DESTDIR         	  ?= $(GOPATH)/bin/${NAME}
BUILD_FLAGS 		  := -ldflags "-X github.com/cosmostation/cosmostation-cosmos/chain-exporter/exporter.Version=${VERSION} \
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
build: go.sum
	@echo "-> Building ${NAME} binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(DESTDIR) .

## Install executable file in $GOBIN direcotry. 
install: go.sum
	@echo "-> Installing ${NAME} binary..."
	@go install -mod=readonly $(BUILD_FLAGS) .

## Create systemd service file in /etc/systemd/system/{SERVICE_NAME}.
make_service:
	@echo "-> Creating ${NAME}.service file..."
	@bash ./service.template

## Move systemd service file and enable it.
enable_service:
	@echo "-> Moving ${NAME} service to /etc/systemd/system/..."
	$(shell sudo mv /home/ubuntu/${NAME}.service /etc/systemd/system/${NAME}.service)
	@echo "-> Enabling ${NAME}.service..."
	$(shell sudo systemctl enable ${NAME}.service)

## Clean the executable file.
clean:
	@echo "-> Cleaning ${NAME} binary..."
	rm -f $(DESTDIR) 2> /dev/null

PHONY: build install make_service enable_service clean