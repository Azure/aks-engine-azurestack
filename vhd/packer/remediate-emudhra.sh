#!/usr/bin/env bash
#
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT license.

set -euo pipefail

verify_emudhra_os() {
    # These package thresholds are specific to Ubuntu 22.04.
    . /etc/os-release
    if [[ "${ID:-}" != ubuntu || "${VERSION_ID:-}" != 22.04 ]]; then
        echo "ERROR: eMudhra remediation requires Ubuntu 22.04." >&2
        return 1
    fi
}

verify_emudhra_packages() {
    local records package status version minimum
    if ! records=$(dpkg-query -W -f='${binary:Package}\t${Status}\t${Version}\n' ca-certificates 'libnss3:*'); then
        echo "ERROR: Could not query installed eMudhra remediation packages." >&2
        return 1
    fi
    local ca_count=0 nss_count=0
    while IFS=$'\t' read -r package status version; do
        case "${package%%:*}" in
            ca-certificates)
                minimum=20260816ubuntu0.22.04.1
                ca_count=$((ca_count + 1))
                ;;
            libnss3)
                minimum=2:3.98-0ubuntu0.22.04.5
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
            echo "ERROR: $package $version is older than the eMudhra fix $minimum." >&2
            return 1
        fi
        echo "Verified eMudhra package remediation: $package $version (minimum $minimum)"
    done <<< "$records"
    if [[ "$ca_count" -ne 1 || "$nss_count" -lt 1 ]]; then
        echo "ERROR: Missing ca-certificates or libnss3 package records." >&2
        return 1
    fi
}

verify_emudhra_bundle() {
    local bundle="${1:-/etc/ssl/certs/ca-certificates.crt}" subjects
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
    if grep -Eq '^subject=.*CN[[:space:]]*=[[:space:]]*(emSign ECC Root CA - G3|emSign Root TLS CA - G3|AAA Certificate Services)(,|$)' <<< "$subjects"; then
        echo "ERROR: A root withdrawn for eMudhra remediation remains in the system CA bundle." >&2
        return 1
    fi
    echo "Verified regenerated system CA bundle: withdrawn root names are absent."
}

remediate_emudhra() {
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
    if [[ "$#" -gt 1 || ( "$#" -eq 1 && "$1" != --verify-only ) ]]; then
        echo "Usage: $0 [--verify-only]" >&2
        return 1
    fi
    verify_emudhra_os
    if [[ "$#" -eq 0 ]]; then
        if [[ "$EUID" -ne 0 ]]; then
            echo "ERROR: Run eMudhra image remediation as root." >&2
            return 1
        fi
        remediate_emudhra
    fi
    verify_emudhra_packages
    verify_emudhra_bundle
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    trap 'echo "ERROR: eMudhra image remediation failed at line $LINENO." >&2' ERR
    main "$@"
fi
