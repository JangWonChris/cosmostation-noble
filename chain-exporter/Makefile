VERSION               := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT                := $(shell git log -1 --format='%H')
DATE 				  := $(shell date +%FT%T%z)
TOOLS_DESTDIR         ?= $(GOPATH)/bin/chain-exporter
BUILD_FLAGS 		  := -ldflags "-X github.com/cosmostation/cosmostation-cosmos/chain-exporter/exporter.Version=${VERSION} \
									-X github.com/cosmostation/cosmostation-cosmos/chain-exporter/exporter.Commit=${COMMIT}"

## Build an executable file in $BINARY_DIR directory.
build: go.sum
	@echo "-> Building chain-exporter binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(TOOLS_DESTDIR) .

## Install executable file in $GOBIN direcotry. 
install: go.sum
	@echo "-> Installing chain-exporter binary..."
	@go install -mod=readonly $(BUILD_FLAGS) .

## Create systemd service file in /etc/systemd/system/{SERVICE_NAME}.
make_service:
	@echo "-> Creating chain-exporter.service file..."
	@bash ./service.template

## Move systemd service file and enable it.
enable_service:
	@echo "-> Moving chain-exporter service to /etc/systemd/system/..."
	$(shell sudo mv /home/ubuntu/chain-exporter.service /etc/systemd/system/chain-exporter.service)
	@echo "-> Enabling chain-exporter.service..."
	$(shell sudo systemctl enable chain-exporter.service)

## Clean the executable file.
clean:
	@echo "-> Cleaning chain-exporter binary..."
	rm -f $(TOOLS_DESTDIR) 2> /dev/null

PHONY: build install make_service enable_service clean
