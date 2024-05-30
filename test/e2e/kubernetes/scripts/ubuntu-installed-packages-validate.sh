#!/bin/bash
set -x
OS=$(sort -r /etc/*-release | gawk 'match($0, /^ID=(.*))$/, a) { print toupper(a[2] a[3]); exit }')
UBUNTU_OS_NAME="UBUNTU"
FLATCAR_OS_NAME="FLATCAR"
DEBIAN_OS_NAME="DEBIAN"
if ! echo "${UBUNTU_OS_NAME} ${FLATCAR_OS_NAME} ${DEBIAN_OS_NAME}" | grep -q "${OS}"; then
    OS=$(sort -r /etc/*-release | gawk 'match($0, /^(ID_LIKE=(.*))$/, a) { print toupper(a[2] a[3]); exit }')
fi

ENSURE_NOT_INSTALLED="
postfix
"
for PACKAGE in ${ENSURE_NOT_INSTALLED}; do
    apt list --installed | grep -E "^${PACKAGE}" && exit 1
done

if [[ $OS == $UBUNTU_OS_NAME ]]; then
    ENSURE_INSTALLED_UBUNTU="
ceph-common
cgroup-lite
glusterfs-client
"
    for PACKAGE in ${ENSURE_INSTALLED_UBUNTU}; do
        apt list --installed | grep -E "^${PACKAGE}" || exit 1
    done
fi

ENSURE_INSTALLED="
apt-transport-https
ca-certificates
chrony
cifs-utils
conntrack
cracklib-runtime
dkms
dbus
ebtables
ethtool
fuse
gcc
git
htop
iftop
init-system-helpers
iotop
iproute2
ipset
iptables
jq
libpam-pkcs11
libpam-pwquality
libpwquality-tools
make
mount
net-tools
nfs-common
ntp
ntpstat
opensc-pkcs11
pigz
socat
sysstat
traceroute
util-linux
vlock
xz-utils
zip
"
for PACKAGE in ${ENSURE_INSTALLED}; do
    apt list --installed | grep -E "^${PACKAGE}" || exit 1
done
