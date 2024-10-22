GO              ?= go
LOCALBIN		= $(shell pwd)/bin
SHELLCHECK_VERSION ?= v0.8.0
AZCLI_VERSION   ?= 2.56.0
PYWINRM_VERSION ?= 0.4.3
KUBECTL_VERSION ?= v1.28.6
PACKER_VERSION  ?= 1.9.5

all: install

.PHONY: install
install: $(LOCALBIN)/go-bindata $(LOCALBIN)/gox $(LOCALBIN)/ginkgo $(LOCALBIN)/golangci-lint $(LOCALBIN)/pub $(LOCALBIN)/mockgen apt-get-install shellcheck azure-cli kubectl azcopy

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

apt-get-install:
	apt-get update && \
	apt-get upgrade -y --no-install-recommends && \
	apt-get install -y --no-install-recommends \
		bash \
		build-essential \
		ca-certificates \
		curl \
		git-core \
		jq \
		libc6 \
		libffi-dev \
		libssl-dev \
		libunwind8 \
		man \
		mercurial \
		net-tools \
		netcat \
		openssh-client \
		procps \
		python3 \
		python3-dev \
		python3-pip \
		python3-setuptools \
		rsync \
		ruby \
		unzip \
		util-linux \
		vim \
		wamerican \
		wget \
		zip \

shellcheck:
	curl -sSL https://github.com/koalaman/shellcheck/releases/download/$(SHELLCHECK_VERSION)/shellcheck-$(SHELLCHECK_VERSION).linux.x86_64.tar.xz | tar -vxJ -C $(LOCALBIN) --strip=1

azure-cli:
	pip3 install --disable-pip-version-check --no-cache-dir --upgrade pip --trusted-host pypi.org --trusted-host files.pythonhosted.org
	pip3 install --disable-pip-version-check --no-cache-dir --upgrade azure-cli==$(AZCLI_VERSION) PyJWT shyaml pywinrm==$(PYWINRM_VERSION) --trusted-host pypi.org --trusted-host files.pythonhosted.org

kubectl:
	curl -sSL https://storage.googleapis.com/kubernetes-release/release/$(KUBECTL_VERSION)/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl
	chmod +x /usr/local/bin/kubectl
	curl -o /usr/local/bin/k https://raw.githubusercontent.com/jakepearson/k/master/k
	chmod +x /usr/local/bin/k

azcopy:
	curl -sSL https://aka.ms/downloadazcopy-v10-linux | tar -vxz -C $(LOCALBIN) --strip=1
	mv $(LOCALBIN)/azcopy $(LOCALBIN)/azcopy-preview
	curl -sSL https://aka.ms/downloadazcopylinux64 | tar -vxz -C /tmp
	/tmp/install.sh

.PHONY: reload
reload: clean install

.PHONY: clean
clean:
	rm -rf $(LOCALBIN)
