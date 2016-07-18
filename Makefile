VERSION=0.1.0-alpha
PROJECT_URL=https://github.com/gocaine/go-dart

SOURCES=$(shell git ls-files '*.go')
GO_DART_BUILD_IMAGE=go-dart.build:latest
GO_DART_RUN_IMAGE=docker run --rm go-dart.build:latest

GO_DART_UI_RUN_IMAGE=docker run --rm -v $(PWD)/webapp:/data ggerbaud/node-bower-grunt:5

# Set the pi user
RPI_USER?=pi
# Set the rpi ip address to hostname rpi in /etc/hosts
RPI=rpi

all: binary

binary: go-dart.ui.make go-dart.make

test:
	$(GO_DART_RUN_IMAGE) scripts/make.sh generate test-unit

go-dart.build-image:
	docker build -t $(GO_DART_BUILD_IMAGE) -f Dockerfile.build .

go-dart.make: go-dart.build-image
	$(GO_DART_RUN_IMAGE) scripts/make.sh generate binary

go-dart.ui.build-image:
	echo "using remote image"

go-dart.ui.make: go-dart.ui.build-image
	$(GO_DART_UI_RUN_IMAGE) npm install && \
	$(GO_DART_UI_RUN_IMAGE)  bower install && \
	$(GO_DART_UI_RUN_IMAGE)  grunt build

deploy:
	scp shell/clean-i2c.sh dist/go-dart $(RPI_USER)@$(RPI):~/
