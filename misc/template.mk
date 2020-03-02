PACKAGE := $(shell go list)
COMPONENT := $(notdir $(PACKAGE))
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
BUILD_DIR ?= build
COV_DIR ?= $(BUILD_DIR)/cov
DIST_DIR ?= $(BUILD_DIR)/dist
GOSUMDB := "off"
GOPROXY := "https://goproxy.cn/"

.PHONY: all dep build clean test coverage lint

all: build

dep: ## Get the dependencies
	@go get -v -d ./...
	## @go get -u github.com/golang/lint/golint

lint: ## Lint the files
	@golint ./...

test: ## Run unittests
	@go test -count=1 -short ./...

race: dep ## Run data race detector
	@go test -race -short ./...

msan: dep ## Run memory sanitizer
	@go test -msan -short ./...

coverage: ## Generate global code coverage report
	@mkdir -p $(COV_DIR);
	@go test -json -covermode=count -coverprofile $(COV_DIR)/$(COMPONENT).cov ./... > $(COV_DIR)/$(COMPONENT).out | true;
	@go tool cover -func=$(COV_DIR)/$(COMPONENT).cov ;
	@go tool cover -html=$(COV_DIR)/$(COMPONENT).cov -o $(COV_DIR)/$(COMPONENT).html ;

build: dep ## Build the binary file
	@mkdir -p $(DIST_DIR);
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${DIST_DIR} . && \
    ([ ! -e etc ] || cp -rf etc ${DIST_DIR}) &&  \
    ([ ! -e public ] || cp -rf public ${DIST_DIR})

clean:
	@rm -rf $(BUILD_DIR)