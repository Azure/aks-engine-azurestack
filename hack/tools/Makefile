GO              ?= go
LOCALBIN		= $(shell pwd)/bin
SHELLCHECK_VERSION ?= v0.8.0
AZCLI_VERSION   ?= 2.56.0
PYWINRM_VERSION ?= 0.4.3
KUBECTL_VERSION ?= v1.30.10

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
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCALBIN) v1.61.0

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
	apt-get update
	apt-get install -y --no-install-recommends apt-transport-https ca-certificates curl gnupg lsb-release
	mkdir -p /etc/apt/keyrings
	curl -sLS https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor | tee /etc/apt/keyrings/microsoft.gpg > /dev/null
	chmod go+r /etc/apt/keyrings/microsoft.gpg
	echo "deb [arch=amd64 signed-by=/etc/apt/keyrings/microsoft.gpg] https://packages.microsoft.com/repos/azure-cli/ bullseye main" | tee /etc/apt/sources.list.d/azure-cli.list
	apt-get update
	apt-get install -y --allow-downgrades --no-install-recommends azure-cli=$(AZCLI_VERSION)-1~bullseye
	apt-mark hold azure-cli

.PHONY: kubectl
kubectl:
	@echo "Entered kubectl"
	curl -LO https://dl.k8s.io/release/$(KUBECTL_VERSION)/bin/linux/amd64/kubectl
	@echo "Downloaded kubectl binary"
	install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
	@echo "Installed kubectl to /usr/local/bin/kubectl"
	install -o root -g root -m 0755 kubectl /usr/local/bin/k
	@echo "Installed kubectl to /usr/local/bin/k"
	kubectl version --client
	k version --client

azcopy:
	curl -sSL https://aka.ms/downloadazcopy-v10-linux | tar -vxz -C $(LOCALBIN) --strip=1
	cp $(LOCALBIN)/azcopy $(LOCALBIN)/azcopy-preview

.PHONY: reload
reload: clean install

.PHONY: clean
clean:
	rm -rf $(LOCALBIN)