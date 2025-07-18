#!/bin/bash

copyPackerFiles() {
  SYSCTL_CONFIG_SRC=/home/packer/sysctl-d-60-CIS.conf
  SYSCTL_CONFIG_DEST=/etc/sysctl.d/60-CIS.conf
  RSYSLOG_CONFIG_SRC=/home/packer/rsyslog-d-60-CIS.conf
  RSYSLOG_CONFIG_DEST=/etc/rsyslog.d/60-CIS.conf
  ETC_ISSUE_CONFIG_SRC=/home/packer/etc-issue
  ETC_ISSUE_CONFIG_DEST=/etc/issue
  ETC_ISSUE_NET_CONFIG_SRC=/home/packer/etc-issue.net
  ETC_ISSUE_NET_CONFIG_DEST=/etc/issue.net
  ETC_ISSUE_STIG_NET_CONFIG_SRC=/home/packer/etc-issue-stig.net
  ETC_ISSUE_STIG_NET_CONFIG_DEST=/etc/issue-stig.net
  SSHD_CONFIG_SRC=/home/packer/sshd_config
  SSHD_CONFIG_DEST=/etc/ssh/sshd_config
  MODPROBE_CIS_SRC=/home/packer/modprobe-CIS.conf
  MODPROBE_CIS_DEST=/etc/modprobe.d/CIS.conf
  PWQUALITY_CONF_SRC=/home/packer/pwquality-CIS.conf
  PWQUALITY_CONF_DEST=/etc/security/pwquality.conf
  PAM_D_SU_SRC=/home/packer/pam-d-su
  PAM_D_SU_DEST=/etc/pam.d/su
  PROFILE_D_CIS_SH_SRC=/home/packer/profile-d-cis.sh
  PROFILE_D_CIS_SH_DEST=/etc/profile.d/CIS.sh
  LABEL_NODES_SRC=/home/packer/label-nodes.sh
  LABEL_NODES_DEST=/opt/azure/containers/label-nodes.sh
  LABEL_NODES_SERVICE_SRC=/home/packer/label-nodes.service
  LABEL_NODES_SERVICE_DEST=/etc/systemd/system/label-nodes.service
  CIS_SRC=/home/packer/cis.sh
  CIS_DEST=/opt/azure/containers/provision_cis.sh
  APT_PREFERENCES_SRC=/home/packer/apt-preferences
  APT_PREFERENCES_DEST=/etc/apt/preferences
  HEALTH_MONITOR_SRC=/home/packer/health-monitor.sh
  HEALTH_MONITOR_DEST=/usr/local/bin/health-monitor.sh
  DOCKER_CLEAR_MOUNT_PROPAGATION_FLAGS_SRC=/home/packer/docker_clear_mount_propagation_flags.conf
  DOCKER_CLEAR_MOUNT_PROPAGATION_FLAGS_DEST=/etc/systemd/system/docker.service.d/clear_mount_propagation_flags.conf
  GENERATE_PROXY_CERTS_SRC=/home/packer/generateproxycerts.sh
  GENERATE_PROXY_CERTS_DEST=/etc/kubernetes/generate-proxy-certs.sh
  COLLECT_LOGS_SRC=/home/packer/collect-logs.sh
  COLLECT_LOGS_DEST=/opt/azure/containers/collect-logs.sh
  NOTICE_SRC=/home/packer/NOTICE.txt
  NOTICE_DEST=/NOTICE.txt
  PROVISION_AZURESTACK_CNI_SRC=/home/packer/provision_azurestack_cni.sh
  PROVISION_AZURESTACK_CNI_DEST=/opt/azure/containers/provision_azurestack_cni.sh
  UBUNTU_STIG_SRC=/home/packer/provision_stig_ubuntu2204.sh
  UBUNTU_STIG_DEST=/opt/azure/containers/provision_stig_ubuntu2204.sh
  AUDITD_RULES_SRC=/home/packer/auditd-rules
  AUDITD_RULES_DEST=/etc/audit/rules.d/aks-engine.rules
  CREDENTIAL_PROVIDER_CONFIG_SRC=/home/packer/credential-provider-config.yaml
  CREDENTIAL_PROVIDER_CONFIG_DEST=/var/lib/kubelet/credential-provider-config.yaml
  if [[ ${UBUNTU_RELEASE} == "16.04" ]]; then
    SSHD_CONFIG_SRC=/home/packer/sshd_config_1604
  fi
  cpAndMode $SYSCTL_CONFIG_SRC $SYSCTL_CONFIG_DEST 644
  cpAndMode $RSYSLOG_CONFIG_SRC $RSYSLOG_CONFIG_DEST 644
  cpAndMode $ETC_ISSUE_CONFIG_SRC $ETC_ISSUE_CONFIG_DEST 644
  cpAndMode $ETC_ISSUE_NET_CONFIG_SRC $ETC_ISSUE_NET_CONFIG_DEST 644
  cpAndMode $ETC_ISSUE_STIG_NET_CONFIG_SRC $ETC_ISSUE_STIG_NET_CONFIG_DEST 644
  cpAndMode $SSHD_CONFIG_SRC $SSHD_CONFIG_DEST 644
  cpAndMode $MODPROBE_CIS_SRC $MODPROBE_CIS_DEST 644
  cpAndMode $PWQUALITY_CONF_SRC $PWQUALITY_CONF_DEST 600
  cpAndMode $PAM_D_SU_SRC $PAM_D_SU_DEST 644
  cpAndMode $PROFILE_D_CIS_SH_SRC $PROFILE_D_CIS_SH_DEST 755
  cpAndMode $LABEL_NODES_SRC $LABEL_NODES_DEST 744
  cpAndMode $LABEL_NODES_SERVICE_SRC $LABEL_NODES_SERVICE_DEST 644
  cpAndMode $CIS_SRC $CIS_DEST 744
  cpAndMode $APT_PREFERENCES_SRC $APT_PREFERENCES_DEST 644
  cpAndMode $HEALTH_MONITOR_SRC $HEALTH_MONITOR_DEST 544
  cpAndMode $DOCKER_CLEAR_MOUNT_PROPAGATION_FLAGS_SRC $DOCKER_CLEAR_MOUNT_PROPAGATION_FLAGS_DEST 644
  cpAndMode $GENERATE_PROXY_CERTS_SRC $GENERATE_PROXY_CERTS_DEST 744
  cpAndMode $COLLECT_LOGS_SRC $COLLECT_LOGS_DEST 744
  cpAndMode $NOTICE_SRC $NOTICE_DEST 444
  cpAndMode $PROVISION_AZURESTACK_CNI_SRC $PROVISION_AZURESTACK_CNI_DEST 744
  cpAndMode $UBUNTU_STIG_SRC $UBUNTU_STIG_DEST 744
  cpAndMode $AUDITD_RULES_SRC $AUDITD_RULES_DEST 744
  cpAndMode $CREDENTIAL_PROVIDER_CONFIG_SRC $CREDENTIAL_PROVIDER_CONFIG_DEST 644
}

cpAndMode() {
  src=$1; dest=$2; mode=$3
  DIR=$(dirname "$dest") && mkdir -p ${DIR} && cp $src $dest && chmod $mode $dest || exit 113
}
