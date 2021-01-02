PACKAGE_NAME=gh-wiki-asset-dl
PACKAGES_DIR=pkg

TARGET_OSs=darwin linux # windows
TARGET_ARCHs=amd64 386

GOPATH_BIN=$(GOPATH)/bin

$(GOPATH_BIN):
	mkdir -p $(GOPATH_BIN)

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test ./...

.PHONY: test-cache-clean
test-cache-clean:
	go clean -testcache

.PHONY: test-all
test-all: test-cache-clean test

.PHONY: clean
clean:
	rm -rf $(PACKAGES_DIR)

.PHONY: packages
packages: packages-main packages-windows

.PHONY: packages-main
packages-main:
	for os in ${TARGET_OSs}; do \
		for arch in ${TARGET_ARCHs}; do \
			GOOS=$$os \
			GOARCH=$$arch \
			CGO_ENABLED=0 \
			go build \
				-o $(PACKAGES_DIR)/$(PACKAGE_NAME)-$${os}-$${arch} \
				. \
		; done \
	; done

.PHONY: packages-windows
packages-windows:
	export os=windows; \
	for arch in ${TARGET_ARCHs}; do \
		GOOS=$$os \
		GOARCH=$$arch \
		CGO_ENABLED=0 \
		go build \
			-o $(PACKAGES_DIR)/$(PACKAGE_NAME)-$${os}-$${arch}.exe \
			. \
	; done \

VERSION=$(shell grep VERSION version.go | cut -d '"' -f 2)

.PHONY: version
version:
	@echo $(VERSION)

GHR=$(GOPATH_BIN)/ghr
$(GHR): $(GOPATH_BIN)
	go get -u github.com/tcnksm/ghr && \
	git checkout -- go.mod go.sum

RELEASE_TAG ?= v${VERSION}

.PHONY: release
release: release_setup push_release_tag upload-packages

.PHONY: upload-packages
upload-packages: $(GHR)
	$(GHR) $(RELEASE_TAG) $(PACKAGES_DIR)/

.PHONY: release_setup
release_setup: clean release_check build test-all packages

.PHONY: release_tag
release_tag:
	git tag $(RELEASE_TAG)

.PHONY: push_release_tag
push_release_tag: release_tag
	git push origin $(RELEASE_TAG)

.PHONY: release_check
release_check: release_check_with_git

.PHONY: release_check_with_git
release_check_with_git: exit_if_already_released exit_with_uncommited_changes exit_with_untracked_files

.PHONY: exit_if_already_released
exit_if_already_released:
	git fetch --tags && \
	git tag | grep -e '^$(RELEASE_TAG)$$' && exit 1 || echo "not released"

.PHONY: exit_with_uncommited_changes
exit_with_uncommited_changes:
	@git diff --exit-code > /dev/null

UNTRACKED_FILES := $(shell git ls-files . --exclude-standard --others)

.PHONY: exit_with_untracked_files
exit_with_untracked_files:
	if [ "$(UNTRACKED_FILES)" = "" ]; then \
	  echo "No untracked file" ; \
	else \
	  echo "There is untracked file" && exit 1 ; \
	fi
