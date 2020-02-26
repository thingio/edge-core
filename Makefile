PROJECT_NAME := "edge-core"
WORK_DIR ?= apiserver
IMAGE ?= xiao4er/thingio:${WORK_DIR}-latest
PKG_LIST := $(cd ${WORK_DIR} && go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
GOSUMDB := "off"
GOPROXY := "https://goproxy.cn/,https://goproxy.io/,https://mirrors.aliyun.com/goproxy,http://172.16.11.155:3000"
.PHONY: all dep build clean test coverage coverhtml lint
TOKEN := 39797b23-c6fe-4bac-97a9-de4d9f1c540b

all: build

lint: ## Lint the files
	cd ${WORK_DIR} && golint ${PKG_LIST}

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

login:
	@echo ${TOKEN} | docker login -u xiao4er --password-stdin

build-image: build login
	@cd ${WORK_DIR}; docker build -t ${IMAGE} .; docker push ${IMAGE}

