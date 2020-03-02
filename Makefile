PROJECT := "edge-core"
IMAGE ?= thingio/edge-core:latest
BUILD_DIR ?= ${PWD}/build
GOSUMDB := "off"
GOPROXY := "https://goproxy.cn/"
TOKEN := 789e9b8d-4fba-4f03-8cac-43a2f7150211

all: build

config:
	@cp misc/template.mk datahub/Makefile
	@cp misc/template.mk apiserver/Makefile
	@cp misc/template.mk deviceman/Makefile
	@cp misc/template.mk pipeman/Makefile
	@cp misc/template.mk bootman/Makefile
	@cp misc/template.mk funcman/Makefile

test: config
	mkdir -p build/cov
	make -C datahub     test
	make -C apiserver   test
	make -C deviceman   test
	make -C pipeman     test
	make -C bootman     test
	make -C funcman     test

build: config
	mkdir -p build/dist
	make -C datahub     build DIST_DIR=${BUILD_DIR}/dist
	make -C apiserver   build DIST_DIR=${BUILD_DIR}/dist
	make -C deviceman   build DIST_DIR=${BUILD_DIR}/dist
	make -C pipeman     build DIST_DIR=${BUILD_DIR}/dist
	make -C bootman     build DIST_DIR=${BUILD_DIR}/dist
	make -C funcman     build DIST_DIR=${BUILD_DIR}/dist

coverage: config
	mkdir -p build/cov
	make -C datahub     coverage COV_DIR=${BUILD_DIR}/cov
	make -C apiserver   coverage COV_DIR=${BUILD_DIR}/cov
	make -C deviceman   coverage COV_DIR=${BUILD_DIR}/cov
	make -C pipeman     coverage COV_DIR=${BUILD_DIR}/cov
	make -C bootman     coverage COV_DIR=${BUILD_DIR}/cov
	make -C funcman     coverage COV_DIR=${BUILD_DIR}/cov

clean:
	rm -rf build

docker-login:
	@echo ${TOKEN} | docker login -u thingio --password-stdin

docker-build: build
	@docker build -t ${IMAGE} .;

docker-image: docker-build docker-login
	@docker push ${IMAGE};

