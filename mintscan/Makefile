NAME 			    := mintscan
VERSION 			:= $(shell git describe --tags --abbrev=0)
COMMIT 				:= $(shell git log -1 --format='%H')
BINARY_DIR          ?= $(GOPATH)/bin/${NAME}
BUILD_FLAGS 		:= -ldflags "-X util.BuildVersion=${VERSION} -X main.buildHash=${COMMIT}"

# Commented for now because make2help is not working
# test:
# ifdef $(shell make2help 2> /dev/null)
# 	@make2help
# else
# 	@echo "-> Installing make2help executable file..."
# 	go get github.com/Songmu/make2help/cmd/make2help
# endif

## Show all make target commands
help:
	@make2help $(MAKEFILE_LIST)

## Print out application information
version: 
	@echo "NAME: ${NAME}"
	@echo "VERSION: ${VERSION}"
	@echo "COMMIT: ${COMMIT}"	

## Build executable file in $BINARY_DIR 
build: go.sum
	@echo "-> Building ${NAME} executable file..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(BINARY_DIR) .

## Install executable file in $GOPATH/bin 
install: go.sum
	@echo "-> Installing ${NAME} executable file..."
	@go install -mod=readonly $(BUILD_FLAGS) .

## Create systemd service file in /etc/systemd/system/{SERVICE_NAME}
make_service:
	@echo "-> Creating mintscan.service file..."
	@bash ./service.template
	@echo "-> Moving mintscan service to /etc/systemd/system/..."
	$(shell sudo mv /home/ubuntu/mintscan.service /etc/systemd/system/mintscan.service)
	@echo "-> Enabling mintscan.service..."
	$(shell sudo systemctl enable mintscan.service)

## Clean the executable file
clean:
	@echo "-> Cleaning ${NAME} executable file..."
	rm -f $(BINARY_DIR) 2> /dev/null

.PHONY: help version build install make_service clean