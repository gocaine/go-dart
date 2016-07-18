VERSION=0.1.0-alpha
PROJECT_URL=https://github.com/gocaine/go-dart

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w

TOPLEVEL_PKG=go-dart
GCGLAGS=
GOARGS=
DIST=dist
BINARY=go-dart
TARGET=$(DIST)/$(BINARY)
LDFLAGS=-ldflags "-X go-dart/cmd.GitHash=`git rev-parse HEAD` -X go-dart/cmd.BuildDate=`date -u +"%Y-%m-%dT%H:%M:%SZ"` -X go-dart/cmd.Version=$(VERSION) -X go-dart/cmd.ProjectUrl=$(PROJECT_URL)"

SOURCES=$(shell git ls-files '*.go')

# Set the pi user
RPI_USER?=pi
# Set the rpi ip address to hostname rpi in /etc/hosts
RPI=rpi

all: clean format build test

bootstrap:
	glide install
	go get -u -v github.com/golang/lint/golint
	go get -u github.com/jteeuwen/go-bindata/...

clean:
	if [ -f ${TARGET} ] ; then rm ${TARGET} ; fi

verbose:
	$(eval GCGLAGS=-x -gcflags=-m)

arm:
	$(eval GOARGS=GOARCH=arm GOOS=linux)

test:
	$(GOTEST) -v `glide novendor`

format:
	gofmt -s -l -w $(SOURCES)

build-web:
	cd webapp && grunt build

generate: build-web
	$(GOCMD) generate

build: generate
	mkdir -p $(DIST)
	$(GOARGS) $(GOBUILD) $(LDFLAGS) $(GCGLAGS) -o $(TARGET) $(TOPLEVEL_PKG)

deploy: arm clean build
	scp shell/clean-i2c.sh $(TARGET) $(RPI_USER)@$(RPI):~/

.PHONY: arm clean
