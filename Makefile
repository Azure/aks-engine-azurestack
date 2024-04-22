TARGETS           = darwin/amd64 linux/amd64 windows/amd64
DIST_DIRS         = find * -type d -exec

.NOTPARALLEL:

.PHONY: bootstrap build test test_fmt validate-copyright-headers fmt lint ci

ifdef DEBUG
GOFLAGS   := -gcflags="-N -l" -mod=vendor
else
GOFLAGS   := -mod=vendor
endif

# go option
GO              ?= go
TAGS            :=
LDFLAGS         :=
BINDIR          := $(CURDIR)/bin
PROJECT         := aks-engine-azurestack
VERSION         ?= $(shell git rev-parse HEAD)
VERSION_SHORT   ?= $(shell git rev-parse --short HEAD)
GITTAG          := $(shell git describe --exact-match --tags $(shell git log -n1 --pretty='%h') 2> /dev/null)
GOBIN           ?= $(shell $(GO) env GOPATH)/bin
TOOLSBIN        := $(CURDIR)/hack/tools/bin
ACK_GINKGO_RC   := true
ifeq ($(GITTAG),)
GITTAG := $(VERSION_SHORT)
endif

DEV_ENV_IMAGE := mcr.microsoft.com/oss/azcu/go-dev:v1.36.2
DEV_ENV_WORK_DIR := /aks-engine-azurestack
DEV_ENV_OPTS := --rm -v $(GOPATH)/pkg/mod:/go/pkg/mod -v $(CURDIR):$(DEV_ENV_WORK_DIR) -w $(DEV_ENV_WORK_DIR) $(DEV_ENV_VARS)
DEV_ENV_CMD := docker run $(DEV_ENV_OPTS) $(DEV_ENV_IMAGE)
DEV_ENV_CMD_IT := docker run -it $(DEV_ENV_OPTS) $(DEV_ENV_IMAGE)
DEV_CMD_RUN := docker run $(DEV_ENV_OPTS)
ifdef DEBUG
LDFLAGS := -X main.version=$(VERSION)
else
LDFLAGS := -s -X main.version=$(VERSION)
endif
BINARY_DEST_DIR ?= bin

ifeq ($(OS),Windows_NT)
	EXTENSION = .exe
	SHELL     = cmd.exe
	CHECK     = where.exe
else
	EXTENSION =
	SHELL     = bash
	CHECK     = which
endif

# Active module mode, as we use go modules to manage dependencies
export GO111MODULE=on

# Add the tools bin to the front of the path
export PATH := $(TOOLSBIN):$(PATH)

all: build

.PHONY: dev
dev:
	$(DEV_ENV_CMD_IT) bash

.PHONY: validate-dependencies
validate-dependencies: bootstrap
	@./scripts/validate-dependencies.sh

.PHONY: validate-copyright-headers
validate-copyright-headers:
	@./scripts/validate-copyright-header.sh

.PHONY: validate-go
validate-go:
	@./scripts/validate-go.sh

.PHONY: validate-shell
validate-shell:
	@./scripts/validate-shell.sh

.PHONY: generate
generate: bootstrap
	@echo "$$(go-bindata --version)"
	go generate $(GOFLAGS) -v ./... > /dev/null 2>&1

.PHONY: build
build: generate go-build

.PHONY: go-build
go-build:
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -buildvcs=false -o $(BINDIR)/$(PROJECT)$(EXTENSION) $(REPO_PATH)

.PHONY: tidy
tidy:
	$(GO) mod tidy
	make -C ./test/e2e tidy

.PHONY: vendor
vendor: tidy
	$(GO) mod vendor

build-binary: generate
	go build $(GOFLAGS) -v -ldflags "$(LDFLAGS)" -buildvcs=false -o $(BINARY_DEST_DIR)/aks-engine-azurestack .

# usage: make clean build-cross dist VERSION=v0.4.0
.PHONY: build-cross
build-cross: build
build-cross: LDFLAGS += -extldflags "-static"
build-cross:
	CGO_ENABLED=0 gox -output="_dist/aks-engine-azurestack-$(GITTAG)-{{.OS}}-{{.Arch}}/{{.Dir}}" -osarch='$(TARGETS)' $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)'

.PHONY: dist
dist: build-cross
	( \
		cd _dist && \
		$(DIST_DIRS) cp ../LICENSE {} \; && \
		$(DIST_DIRS) cp ../README.md {} \; && \
		$(DIST_DIRS) tar -zcf {}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r {}.zip {} \; \
	)

.PHONY: checksum
checksum:
	for f in _dist/*.{gz,zip} ; do \
		shasum -a 256 "$${f}"  | awk '{print $$1}' > "$${f}.sha256" ; \
	done

.PHONY: clean
clean: tools-clean
	@rm -rf $(BINDIR) ./_dist ./pkg/helpers/unit_tests

GIT_BASEDIR    = $(shell git rev-parse --show-toplevel 2>/dev/null)
ifneq ($(GIT_BASEDIR),)
	LDFLAGS += -X github.com/Azure/aks-engine-azurestack/pkg/test.JUnitOutDir=$(GIT_BASEDIR)/test/junit
endif

ginkgoBuild: generate
	make -C ./test/e2e ginkgo-build

test: generate
	ginkgo -mod=vendor -skipPackage test/e2e -failFast -r -v -tags=fast -ldflags '$(LDFLAGS)' .

.PHONY: test-style
test-style: validate-go validate-shell validate-copyright-headers

.PHONY: ensure-generated
ensure-generated:
	@echo "==> Checking generated files <=="
	@scripts/ensure-generated.sh

.PHONY: test-e2e
test-e2e:
	@test/e2e.sh

HAS_GIT := $(shell $(CHECK) git)

.PHONY: bootstrap
bootstrap: tools-install
ifndef HAS_GIT
	$(error You must install Git)
endif

.PHONY: tools-reload
tools-reload:
	make -C hack/tools reload

.PHONY: tools-install
tools-install:
	make -C hack/tools/

.PHONY: tools-clean
tools-clean:
	make -C hack/tools/ clean

include versioning.mk
include test.mk
include packer.mk
