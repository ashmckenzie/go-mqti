PROJECT_NAME ?= littlefly
VERSION=$(shell cat VERSION)

BUILD_IMAGE_NAME = $(PROJECT_NAME)-build

APP_DIR ?= $(PWD)
BIN_DIR=$(APP_DIR)/bin
RELEASE_DIR=$(APP_DIR)/release
VENDOR_DIR=$(APP_DIR)/vendor

BUILD_ARGS=-i -ldflags=""

DARWIN_RELEASE=${PROJECT_NAME}_darwin_$(VERSION)
LINUX_RELEASE=${PROJECT_NAME}_linux_$(VERSION)

GITHUB_RELEASE_ARGS=--user ${GITHUB_USERNAME} --repo ${GITHUB_REPO} --tag ${VERSION}

.PHONY: create_build_image deps update_deps make_dirs test run clean

build: make_dirs test
	go build ${BUILD_ARGS} -o $(BIN_DIR)/$(PROJECT_NAME) ./cmd

run: test
	go run cmd/main.go $(filter-out $@, $(MAKECMDGOALS))

release: make_dirs test
	# FIXME
	# CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo ${BUILD_ARGS} -o ${RELEASE_DIR}/${DARWIN_RELEASE} ./cmd
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ${BUILD_ARGS} -o ${RELEASE_DIR}/${LINUX_RELEASE} ./cmd

release_upload: release
	git tag $(VERSION)
	git push --tags

	github-release release ${GITHUB_RELEASE_ARGS}
	github-release upload ${GITHUB_RELEASE_ARGS} --name ${DARWIN_RELEASE} --file ${RELEASE_DIR}/${DARWIN_RELEASE}
	github-release upload ${GITHUB_RELEASE_ARGS} --name ${LINUX_RELEASE} --file ${RELEASE_DIR}/${LINUX_RELEASE}

test: deps
	go test -cover

deps: make_dirs
	go get github.com/kardianos/govendor
	go get github.com/aktau/github-release
	govendor fetch +missing
	go generate -v -x ./ ./cmd ./commands

update_deps:
	go get -u github.com/kardianos/govendor
	go get -u github.com/aktau/github-release
	govendor update +e

make_dirs:
	@mkdir -p ${BIN_DIR} ${RELEASE_DIR} ${VENDOR_DIR}

clean:
	rm -rf $(BIN_DIR) $(RELEASE_DIR)
