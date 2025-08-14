#!/usr/bin/env bash
#
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT license.

set -euo pipefail

# Function to install Azure CLI
install_azure_cli() {
    apt-get update
    apt-get install -y --no-install-recommends apt-transport-https ca-certificates curl gnupg lsb-release
    mkdir -p /etc/apt/keyrings
    curl -sLS https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor | tee /etc/apt/keyrings/microsoft.gpg > /dev/null
    chmod go+r /etc/apt/keyrings/microsoft.gpg
    echo "deb [arch=amd64 signed-by=/etc/apt/keyrings/microsoft.gpg] https://packages.microsoft.com/repos/azure-cli/ bullseye main" | tee /etc/apt/sources.list.d/azure-cli.list
    apt-get update
    apt-get install -y --allow-downgrades --no-install-recommends azure-cli="${AZCLI_VERSION}-1~bullseye"
    apt-mark hold azure-cli
    
    echo "Azure CLI ${AZCLI_VERSION} installed successfully"
}

# Check if Azure CLI is already installed
if ! command -v az >/dev/null 2>&1; then
    install_azure_cli
else
    AZ_PATH=$(command -v az)
    echo "Azure CLI already installed at ${AZ_PATH}, skipping installation"
fi