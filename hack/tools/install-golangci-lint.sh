#!/usr/bin/env bash
set -euo pipefail

if [[ "$#" -ne 2 || ! "$1" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ || -z "$2" ]]; then
    echo "Usage: $0 <vMAJOR.MINOR.PATCH> <installation directory>" >&2
    exit 1
fi

version="${1#v}"
destination="$2"
case "$(uname -s)" in
    Linux) os=linux ;;
    Darwin) os=darwin ;;
    *) echo "ERROR: Unsupported golangci-lint installer OS." >&2; exit 1 ;;
esac
case "$(uname -m)" in
    x86_64|amd64) arch=amd64 ;;
    aarch64|arm64) arch=arm64 ;;
    i386|i686) arch=386 ;;
    armv6l) arch=armv6 ;;
    armv7l) arch=armv7 ;;
    *) echo "ERROR: Unsupported golangci-lint installer architecture." >&2; exit 1 ;;
esac

name="golangci-lint-$version-$os-$arch"
archive="$name.tar.gz"
base_url="https://github.com/golangci/golangci-lint/releases/download/v$version"
temporary=$(mktemp -d)
trap 'rm -rf -- "$temporary"' EXIT

curl -fsSL --retry 3 "$base_url/$archive" -o "$temporary/$archive"
curl -fsSL --retry 3 "$base_url/golangci-lint-$version-checksums.txt" -o "$temporary/checksums.txt"
expected=$(awk -v file="$archive" '$2 == file { print $1; count++ } END { if (count != 1) exit 1 }' "$temporary/checksums.txt") || {
    echo "ERROR: Expected exactly one checksum for $archive." >&2
    exit 1
}
if [[ ! "$expected" =~ ^[[:xdigit:]]{64}$ ]]; then
    echo "ERROR: Invalid SHA-256 checksum for $archive." >&2
    exit 1
fi
if command -v sha256sum >/dev/null; then
    actual=$(sha256sum "$temporary/$archive")
else
    actual=$(shasum -a 256 "$temporary/$archive")
fi
if [[ "${actual%% *}" != "$expected" ]]; then
    echo "ERROR: SHA-256 verification failed for $archive." >&2
    exit 1
fi

tar -xzf "$temporary/$archive" -C "$temporary" "$name/golangci-lint"
mkdir -p "$destination"
install -m 0755 "$temporary/$name/golangci-lint" "$destination/golangci-lint"
