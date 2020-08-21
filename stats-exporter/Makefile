VERSION               := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT                := $(shell git log -1 --format='%H')
DATE 				  := $(shell date +%FT%T%z)
DESTDIR         	  ?= $(GOPATH)/bin/stats-exporter
BUILD_FLAGS 		  := -ldflags "-X github.com/cosmostation/cosmostation-cosmos/stats-exporter/exporter.Version=${VERSION} \
									-X github.com/cosmostation/cosmostation-cosmos/stats-exporter/exporter.Commit=${COMMIT}"

## Build an executable file in $DESTDIR directory.
build: go.sum
	@echo "-> Building stats-exporter binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(DESTDIR) .
	
## Install executable file in $GOBIN direcotry. 
install: go.sum
	@echo "-> Installing stats-exporter binary..."
	@go install -mod=readonly $(BUILD_FLAGS) .

## Create systemd service file in /etc/systemd/system/{SERVICE_NAME}.
make_service:
	@echo "-> Creating stats-exporter.service file..."
	@bash ./service.template

## Move systemd service file and enable it.
enable_service:
	@echo "-> Moving stats-exporter service to /etc/systemd/system/..."
	$(shell sudo mv /home/ubuntu/stats-exporter.service /etc/systemd/system/stats-exporter.service)
	@echo "-> Enabling stats-exporter.service..."
	$(shell sudo systemctl enable stats-exporter.service)

## Clean the executable file.
clean:
	@echo "cleaning stats-exporter binary..."
	rm -f $(TOOLS_DESTDIR) 2> /dev/null

PHONY: build install make_service enable_service clean
