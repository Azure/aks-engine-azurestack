GO              ?= go
LOCALBIN		= $(shell pwd)/bin
SHELLCHECK_VERSION ?= v0.7.2

all: install

.PHONY: install
install: $(LOCALBIN)/go-bindata $(LOCALBIN)/gox $(LOCALBIN)/ginkgo $(LOCALBIN)/golangci-lint $(LOCALBIN)/pub $(LOCALBIN)/mockgen $(LOCALBIN)/shellcheck

$(LOCALBIN)/go-bindata:
    GOBIN=$(LOCALBIN) $(GO) install github.com/go-bindata/go-bindata/go-bindata@v3.1.2

$(LOCALBIN)/gox:
    GOBIN=$(LOCALBIN) $(GO) install github.com/mitchellh/gox@v1.0.1

$(LOCALBIN)/ginkgo:
    GOBIN=$(LOCALBIN) $(GO) install github.com/onsi/ginkgo/v2/ginkgo@v2.17.1

$(LOCALBIN)/golangci-lint:
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCALBIN) v1.51.2

$(LOCALBIN)/pub:
    GOBIN=$(LOCALBIN) $(GO) install github.com/devigned/pub@v0.2.6

$(LOCALBIN)/mockgen:
    GOBIN=$(LOCALBIN) $(GO) install github.com/golang/mock/mockgen@v1.2.0

$(LOCALBIN)/shellcheck:
    curl -sSL https://github.com/koalaman/shellcheck/releases/download/$(SHELLCHECK_VERSION)/shellcheck-$(SHELLCHECK_VERSION).linux.x86_64.tar.xz | tar -vxJ -C $(LOCALBIN) --strip=1

.PHONY: reload
reload: clean install

.PHONY: clean
clean:
	rm -rf $(LOCALBIN)
