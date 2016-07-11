# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w

TOPLEVEL_PKG=go-dart
GOARGS=GOARCH=arm GOOS=linux
DIST=dist
BINARY=go-dart
TARGET=$(DIST)/$(BINARY)

# Set the pi user
RPI_USER?=pi
# Set the rpi ip address to hostname rpi in /etc/hosts
RPI=rpi
all: build

local:
	$(eval GOARGS = )
	$(eval BIN_ARGS = "hardware" )

run-client: build
	$(TARGET) $(BIN_ARGS)

run-server: build
	$(TARGET) $(BIN_ARGS) server

build:
	mkdir -p $(DIST)
	$(GOARGS) $(GOBUILD) -o $(TARGET) $(TOPLEVEL_PKG)

deploy: build
	scp shell/clean-i2c.sh $(TARGET) $(RPI_USER)@$(RPI):~/
