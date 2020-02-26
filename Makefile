PROJECT_NAME := "edge-core"
WORK_DIR ?= apiserver
PKG_LIST := $(shell cd ${WORK_DIR} && go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
GOSUMDB := "off"
GOPROXY := "https://goproxy.cn/,https://goproxy.io/,https://mirrors.aliyun.com/goproxy,http://172.16.11.155:3000"
.PHONY: all dep build clean test coverage coverhtml lint

all: build

work: ## cd work directory
	@
lint: ## Lint the files
	@cd ${WORK_DIR} && golint ${PKG_LIST}

test: ## Run unittests
	@cd ${WORK_DIR} && go test -short ${PKG_LIST}

race: dep ## Run data race detector
	@cd ${WORK_DIR} && go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@cd ${WORK_DIR} && go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	cp coverage.sh ${WORK_DIR} && cd ${WORK_DIR} && bash coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	cp coverage.sh ${WORK_DIR} && cd ${WORK_DIR} && bash coverage.sh html;

dep: ## Get the dependencies
	@cd ${WORK_DIR} && go get -v -d ./...

build: dep ## Build the binary file
	@cd ${WORK_DIR} && go build -i -v .

clean: ## Remove previous build
	@rm -rf ${WORK_DIR}

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'