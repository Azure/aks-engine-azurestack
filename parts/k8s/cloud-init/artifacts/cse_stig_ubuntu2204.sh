#!/bin/bash

setLoginDefs() {
  local f=/etc/login.defs
  sed -i '/^PASS_MAX_DAYS/d' ${f} || exit 115
  sed -i '$aPASS_MAX_DAYS 60' ${f} || exit 115
  sed -i '/^PASS_MIN_DAYS/d' ${f} || exit 115
  sed -i '$aPASS_MIN_DAYS 1' ${f} || exit 115
  sed -i '/^UMASK/d' ${f} || exit 115
  sed -i '$aUMASK 077' ${f} || exit 115
}
setPwqualityConf() {
  local f=/etc/security/pwquality.conf
  sed -i '/^difok/d' ${f} || exit 115
  sed -i '$adifok=8' ${f} || exit 115
  sed -i '/^dictcheck/d' ${f} || exit 115
  sed -i '$adictcheck=1' ${f} || exit 115
  sed -i '/^minlen/d' ${f} || exit 115
  sed -i '$aminlen=15' ${f} || exit 115
  sed -i '/^lcredit/d' ${f} || exit 115
  sed -i '$alcredit=-1' ${f} || exit 115
}
setTerminalTimeout() {
  local f=/etc/profile.d/99-terminal_tmout.sh
  {{/* STIG SV-238207r653796_rule */}}
  if [[ -f ${f} ]]; then
    sed -i '/^TMOUT/d' ${f} || exit 115
    sed -i '$aTMOUT=600' ${f} || exit 115
  else
    echo "TMOUT=600" > ${f}
    truncate -s -1 ${f}
  fi
}
setSSHDConfig() {
  local f=/etc/ssh/sshd_config
  {{/* STIG SV-238212r653811_rule */}}
  sed -i '/^ClientAliveCountMax/d' ${f} || exit 115
  sed -i '$aClientAliveCountMax 1' ${f} || exit 115
  {{/* STIG SV-238216r654316_rule */}}
  sed -i '/^MACs/d' ${f} || exit 115
  sed -i '$aMACs hmac-sha2-512,hmac-sha2-256' ${f} || exit 115
  {{/* STIG SV-238217r653826_rule */}}
  sed -i '/^Ciphers/d' ${f} || exit 115
  sed -i '$aCiphers aes256-ctr,aes192-ctr,aes128-ctr' ${f} || exit 115
  {{/* STIG SV-238220r653835_rule */}}
  sed -i '/^X11UseLocalhost/d' ${f} || exit 115
  sed -i '$aX11UseLocalhost yes' ${f} || exit 115
  {{/* STIG SV-238214r653817_rule */}}
  if [[ -f /etc/issue-stig.net ]]; then
    sed -i '/^Banner/d' ${f} || exit 115
    sed -i '$aBanner /etc/issue-stig.net' ${f} || exit 115
  fi
}
setAuditd() {
  local f=/etc/audit/auditd.conf
  {{/* STIG SV-238244r653907_rule */}}
  sed -i '/^disk_full_action/d' ${f} || exit 115
  sed -i '$adisk_full_action = HALT' ${f} || exit 115
}
setLimitsConf() {
  local f=/etc/security/limits.conf
  {{/* STIG SV-238323r654144_rule */}}
  sed -i '1s|^|* hard maxlogins 10\n|' ${f}
}
setAPTConfig() {
  local f=/etc/apt/apt.conf.d/01-vendor-ubuntu
  {{/* STIG SV-219155r610963_rule */}}
  sed -i '/^APT::Get::AllowUnauthenticated/d' ${f} || exit 115
  sed -i '$aAPT::Get::AllowUnauthenticated "false";' ${f} || exit 115
  local g=/etc/apt/apt.conf.d/50unattended-upgrades
  {{/* STIG SV-219156r610963_rule */}}
  sed -i '/^Unattended-Upgrade::Remove-Unused-Dependencies/d' ${g} || exit 115
  sed -i '$aUnattended-Upgrade::Remove-Unused-Dependencies "true";' ${g} || exit 115
  {{/* STIG SV-219156r610963_rule */}}
  sed -i '/^Unattended-Upgrade::Remove-Unused-Kernel-Packages/d' ${g} || exit 115
  sed -i '$aUnattended-Upgrade::Remove-Unused-Kernel-Packages "true";' ${g} || exit 115
}
setLoginDefs
setPwqualityConf
setTerminalTimeout
setSSHDConfig
setAuditd
setLimitsConf
setAPTConfig
#EOF
