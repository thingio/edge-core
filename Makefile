PROJECT_NAME := "edge-core"
WORK_DIR = apiserver
IMAGE ?= xiao4er/thingio:${WORK_DIR}-latest
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
GOSUMDB := "off"
GOPROXY := "https://goproxy.cn/,https://goproxy.io/,https://mirrors.aliyun.com/goproxy,http://172.16.11.155:3000"
.PHONY: all dep build clean test coverage coverhtml lint
TOKEN := 39797b23-c6fe-4bac-97a9-de4d9f1c540b

all: build

lint: ## Lint the files
	cd ${WORK_DIR} && golint ./...

test: ## Run unittests
	@cd ${WORK_DIR} && go test -short ./...

race: dep ## Run data race detector
	@cd ${WORK_DIR} && go test -race -short ./...

msan: dep ## Run memory sanitizer
	@cd ${WORK_DIR} && go test -msan -short ./...

coverage: ## Generate global code coverage report
	cp coverage.sh ${WORK_DIR} && cd ${WORK_DIR} && bash coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	cp coverage.sh ${WORK_DIR} && cd ${WORK_DIR} && bash coverage.sh html;

dep: ## Get the dependencies
	@cd ${WORK_DIR} && go get -v -d ./...

build: dep ## Build the binary file
	cd ${WORK_DIR} && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/${WORK_DIR} .

login:
	@echo ${TOKEN} | docker login -u xiao4er --password-stdin

build-image: build
	@cd ${WORK_DIR}; docker build -t ${IMAGE} .;

push-image: build-image login
	@docker push ${IMAGE};

