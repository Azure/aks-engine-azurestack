#!/bin/bash

sudo grep "* hard maxlogins 10" /etc/security/limits.conf || exit 1

sudo grep "disk_full_action = HALT" /etc/audit/auditd.conf || exit 1

sudo grep "TMOUT=600" /etc/profile.d/99-terminal_tmout.sh || exit 1

sudo grep "difok=8" /etc/security/pwquality.conf || exit 1
sudo grep "dictcheck=1" /etc/security/pwquality.conf || exit 1
sudo grep "minlen=15" /etc/security/pwquality.conf || exit 1
sudo grep "lcredit=-1" /etc/security/pwquality.conf || exit 1

sudo grep 'APT::Get::AllowUnauthenticated "false"' /etc/apt/apt.conf.d/01-vendor-ubuntu || exit 1
sudo grep 'Unattended-Upgrade::Remove-Unused-Dependencies "true"' /etc/apt/apt.conf.d/50unattended-upgrades || exit 1
sudo grep 'Unattended-Upgrade::Remove-Unused-Kernel-Packages "true"' /etc/apt/apt.conf.d/50unattended-upgrades || exit 1

set -x

CONFIGS=("ClientAliveInterval 120"
"ClientAliveCountMax 1"
"MACs hmac-sha2-512,hmac-sha2-256"
"KexAlgorithms curve25519-sha256@libssh.org"
"Ciphers aes256-ctr,aes192-ctr,aes128-ctr"
"HostKey /etc/ssh/ssh_host_rsa_key"
"HostKey /etc/ssh/ssh_host_dsa_key"
"HostKey /etc/ssh/ssh_host_ecdsa_key"
"HostKey /etc/ssh/ssh_host_ed25519_key"
"SyslogFacility AUTH"
"LogLevel INFO"
"LoginGraceTime 60"
"PermitRootLogin no"
"PermitUserEnvironment no"
"StrictModes yes"
"PubkeyAuthentication yes"
"IgnoreRhosts yes"
"HostbasedAuthentication no"
"X11Forwarding no"
"X11UseLocalhost yes"
"MaxAuthTries 4"
"AcceptEnv LANG LC_*"
"Subsystem sftp /usr/lib/openssh/sftp-server"
"UsePAM yes"
"UseDNS no"
"Banner /etc/issue-stig.net"
"GSSAPIAuthentication no")

for ((i = 0; i < ${#CONFIGS[@]}; i++))
do
    grep -i "${CONFIGS[$i]}" /etc/ssh/sshd_config || exit 1
done
