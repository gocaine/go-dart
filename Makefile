VERSION=0.1.0-alpha
PROJECT_URL=https://github.com/gocaine/go-dart

SOURCES=$(shell git ls-files '*.go')

BUILD_IMAGE=go-dart.build:latest
RUN_IMAGE=docker run --rm -e GITHUB_TOKEN -v $(CURDIR)/dist:/go/src/github.com/gocaine/go-dart/dist -v $(CURDIR)/reports:/go/src/github.com/gocaine/go-dart/reports go-dart.build:latest

UI_BUILD_IMAGE=ggerbaud/node-bower-grunt:5
UI_RUN_IMAGE=docker run --rm -v $(PWD)/webapp:/data $(UI_BUILD_IMAGE)
UI_RUN_PRESTEP=

# Set the pi user
RPI_USER?=pi
# Set the rpi ip address to hostname rpi in /etc/hosts
RPI=rpi

PHONY: dev

dev: ## use local tools instead of docker containers
	@echo "Configuring dev build..."
	$(eval USE_LOCAL=local)
	$(eval RUN_IMAGE :=)
	$(eval UI_RUN_IMAGE :=)
	$(eval UI_RUN_PRESTEP := cd webapp &&)

all: binary

arm: ## build for ARM target
	$(eval GOARCH=GOOS=linux GOARCH=arm)

mock.ui: 
	if [ ! -e webapp/dist/index.html ]; then \
		mkdir -p webapp/dist; \
		echo "void starts here" > webapp/dist/index.html; \
	fi

binary-noui: mock.ui build.go ## package the core w/o ui

binary: build.ui build.go ## package the webui and the core

test: ## run all tests
	$(RUN_IMAGE) scripts/make.sh generate test-unit

test-coverage: ## run all tests w/ coverage
	$(RUN_IMAGE) scripts/make.sh generate test-coverage

test-coverage-report:
	$(RUN_IMAGE) scripts/make.sh generate test-coverage-report

build.go-image:
	@if [ "$(USE_LOCAL)" != "local" ]; then \
		docker build -t $(BUILD_IMAGE) -f Dockerfile.build . ;\
	fi

build.go: build.go-image
	$(RUN_IMAGE) $(GOARCH) scripts/make.sh generate binary

build.ui-image:
	@if [ "$(USE_LOCAL)" != "local" ]; then \
		echo "using remote image" ;\
	fi

build.ui: build.ui-image
	$(UI_RUN_PRESTEP) $(UI_RUN_IMAGE) npm install && \
	$(UI_RUN_IMAGE)  bower install && \
	$(UI_RUN_IMAGE)  grunt build

validate:
	$(RUN_IMAGE) scripts/make.sh validate-gofmt validate-govet validate-golint

format:
	scripts/make.sh format

release: ## create a release on github
	$(RUN_IMAGE) scripts/make.sh release


deploy: ## actually deploy on rpi
	scp -r shell/clean-i2c.sh boards dist/go-dart $(RPI_USER)@$(RPI):~/

help: ## this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
