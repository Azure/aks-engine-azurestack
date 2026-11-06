#!/usr/bin/env bash
#
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT license.

set -euo pipefail

load_ca_profile() {
    local records
    if ! command -v python3 >/dev/null; then
        echo "ERROR: Reading the CA remediation profile requires python3." >&2
        return 1
    fi
    if ! records=$(python3 - "$1" <<'PY'
import json
import re
import sys

with open(sys.argv[1], encoding="utf-8-sig") as source:
    profile = json.load(source)
if type(profile.get("schemaVersion")) is not int or profile["schemaVersion"] != 1:
    raise ValueError("Unsupported CA remediation profile schemaVersion")
incident = profile.get("incident")
if not isinstance(incident, str) or not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9 ._-]{0,99}", incident):
    raise ValueError("Invalid CA remediation incident label")
linux = profile["linux"]
version = linux["ubuntuVersion"]
if not isinstance(version, str) or not re.fullmatch(r"[0-9]{2}\.[0-9]{2}", version):
    raise ValueError("Invalid Ubuntu version")
minimums = [linux["minimumCaCertificatesVersion"], linux["minimumLibnss3Version"]]
if any(not isinstance(value, str) or not re.fullmatch(r"[0-9][A-Za-z0-9.+:~_-]*", value) for value in minimums):
    raise ValueError("Invalid minimum package version")
roots = linux["withdrawnRootNames"]
if not isinstance(roots, list) or not roots or any(
    not isinstance(root, str) or not root.strip() or any(char in root for char in "\r\n,") for root in roots
):
    raise ValueError("withdrawnRootNames must contain nonempty, single-line common names without commas")
print("\n".join([incident, version] + minimums + roots))
PY
    ); then
        echo "ERROR: Cannot load CA remediation profile: $1" >&2
        return 1
    fi
    local -a fields
    mapfile -t fields <<< "$records"
    ca_incident="${fields[0]}"
    ca_ubuntu_version="${fields[1]}"
    ca_minimum="${fields[2]}"
    nss_minimum="${fields[3]}"
    ca_withdrawn_roots=("${fields[@]:4}")
    dpkg --validate-version "$ca_minimum"
    dpkg --validate-version "$nss_minimum"
}

verify_ca_os() {
    . /etc/os-release
    if [[ "${ID:-}" != ubuntu || "${VERSION_ID:-}" != "$ca_ubuntu_version" ]]; then
        echo "ERROR: CA remediation for $ca_incident requires Ubuntu $ca_ubuntu_version." >&2
        return 1
    fi
}

verify_ca_packages() {
    local records package status version minimum
    if ! records=$(dpkg-query -W -f='${binary:Package}\t${Status}\t${Version}\n' ca-certificates 'libnss3:*'); then
        echo "ERROR: Could not query installed CA remediation packages." >&2
        return 1
    fi
    local ca_count=0 nss_count=0
    while IFS=$'\t' read -r package status version; do
        case "${package%%:*}" in
            ca-certificates)
                minimum="$ca_minimum"
                ca_count=$((ca_count + 1))
                ;;
            libnss3)
                minimum="$nss_minimum"
                nss_count=$((nss_count + 1))
                ;;
            *)
                echo "ERROR: Unexpected package record: $package" >&2
                return 1
                ;;
        esac
        if [[ "$status" != "install ok installed" || -z "$version" ]]; then
            echo "ERROR: $package is not fully installed: $status" >&2
            return 1
        fi
        if ! dpkg --compare-versions "$version" ge "$minimum"; then
            echo "ERROR: $package $version is older than the $ca_incident fix $minimum." >&2
            return 1
        fi
        echo "Verified CA package remediation for $ca_incident: $package $version (minimum $minimum)"
    done <<< "$records"
    if [[ "$ca_count" -ne 1 || "$nss_count" -lt 1 ]]; then
        echo "ERROR: Missing ca-certificates or libnss3 package records." >&2
        return 1
    fi
}

verify_ca_bundle() {
    local bundle="${1:-/etc/ssl/certs/ca-certificates.crt}" subjects subject common_name root
    if [[ ! -s "$bundle" ]]; then
        echo "ERROR: System CA bundle is missing or empty: $bundle" >&2
        return 1
    fi
    if ! subjects=$(openssl crl2pkcs7 -nocrl -certfile "$bundle" | openssl pkcs7 -print_certs -noout); then
        echo "ERROR: Could not parse the system CA bundle." >&2
        return 1
    fi
    if [[ -z "$subjects" ]]; then
        echo "ERROR: System CA bundle contains no readable certificates." >&2
        return 1
    fi
    while IFS= read -r subject; do
        if [[ "$subject" =~ ^subject=.*CN[[:space:]]*=[[:space:]]*(.*) ]]; then
            common_name="${BASH_REMATCH[1]%%,*}"
            for root in "${ca_withdrawn_roots[@]}"; do
                if [[ "$common_name" == "$root" ]]; then
                    echo "ERROR: Root '$root' withdrawn for $ca_incident remains in the system CA bundle." >&2
                    return 1
                fi
            done
        fi
    done <<< "$subjects"
    echo "Verified regenerated system CA bundle: withdrawn root names are absent."
}

remediate_ca_trust() {
    local nss_records
    local -a nss_packages
    apt-get update -o APT::Update::Error-Mode=any
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates libnss3
    nss_records=$(dpkg-query -W -f='${binary:Package}\n' 'libnss3:*')
    mapfile -t nss_packages <<< "$nss_records"
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends "${nss_packages[@]}"
    update-ca-certificates --fresh
}

main() {
    local profile="$(dirname "${BASH_SOURCE[0]}")/ca-remediation.json" verify_only=false
    while [[ "$#" -gt 0 ]]; do
        case "$1" in
            --verify-only) verify_only=true; shift ;;
            --profile)
                if [[ "$#" -lt 2 || -z "$2" ]]; then
                    echo "ERROR: --profile requires a JSON file path." >&2
                    return 1
                fi
                profile="$2"; shift 2 ;;
            *)
                echo "Usage: $0 [--profile <JSON file>] [--verify-only]" >&2
                return 1 ;;
        esac
    done
    load_ca_profile "$profile"
    verify_ca_os
    if [[ "$verify_only" == false ]]; then
        if [[ "$EUID" -ne 0 ]]; then
            echo "ERROR: Run CA image remediation as root." >&2
            return 1
        fi
        remediate_ca_trust
    fi
    verify_ca_packages
    verify_ca_bundle
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    trap 'echo "ERROR: CA image remediation failed at line $LINENO." >&2' ERR
    main "$@"
fi
