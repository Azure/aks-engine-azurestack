#!/bin/bash
source /home/packer/provision_installs.sh
source /home/packer/provision_source.sh
source /home/packer/packer_source.sh

VHD_LOGS_FILEPATH=/opt/azure/vhd-install.complete

echo "Starting build on " $(date) > ${VHD_LOGS_FILEPATH}

copyPackerFiles

echo ""
echo "Components downloaded in this VHD build (some of the below components might get deleted during cluster provisioning if they are not needed):" >> ${VHD_LOGS_FILEPATH}

AUDITD_ENABLED=true
MS_APT_REPO=packages.microsoft.com
installDeps
cat << EOF >> ${VHD_LOGS_FILEPATH}
apt packages:
  - apt-transport-https
  - auditd
  - ca-certificates
  - ceph-common
  - cgroup-lite
  - chrony
  - cifs-utils
  - conntrack
  - cracklib-runtime
  - dkms
  - dbus
  - ebtables
  - ethtool
  - fuse
  - gcc
  - git
  - glusterfs-client
  - htop
  - iftop
  - init-system-helpers
  - iotop
  - iproute2
  - ipset
  - iptables
  - jq
  - libpam-pkcs11
  - libpam-pwquality
  - libpwquality-tools
  - make
  - mount
  - net-tools
  - nfs-common
  - ntp
  - ntpstat
  - opensc-pkcs11
  - pigz
  - socat
  - sysstat
  - util-linux
  - vlock
  - xz-utils
  - zip
EOF

chmod a-x /etc/update-motd.d/??-{motd-news,release-upgrade}

overrideNetworkConfig

cat << EOF >> ${VHD_LOGS_FILEPATH}
Binaries:
EOF

apmz_version="v0.5.1"
ensureAPMZ "${apmz_version}"
echo "  - apmz $apmz_version" >> ${VHD_LOGS_FILEPATH}

installBpftrace
echo "  - bpftrace" >> ${VHD_LOGS_FILEPATH}

MOBY_VERSION="20.10.14"
CONTAINERD_VERSION="1.6.36"
installMoby
installRunc
systemctl_restart 100 5 30 docker || exit 1
echo "  - moby v${MOBY_VERSION}" >> ${VHD_LOGS_FILEPATH}
downloadGPUDrivers
echo "  - nvidia-container-runtime" >> ${VHD_LOGS_FILEPATH}

ETCD_VERSION="3.3.25"
ETCD_DOWNLOAD_URL="mcr.microsoft.com/oss/etcd-io/"
installEtcd "docker"
echo "  - etcd v${ETCD_VERSION}" >> ${VHD_LOGS_FILEPATH}

VNET_CNI_VERSIONS="
1.4.59
"
for VNET_CNI_VERSION in $VNET_CNI_VERSIONS; do
    VNET_CNI_PLUGINS_URL="https://packages.aks.azure.com/azure-cni/v${VNET_CNI_VERSION}/binaries/azure-vnet-cni-linux-amd64-v${VNET_CNI_VERSION}.tgz"
    downloadAzureCNI
    echo "  - Azure CNI version ${VNET_CNI_VERSION}" >> ${VHD_LOGS_FILEPATH}
done

CNI_PLUGIN_VERSIONS="
0.9.1
"
for CNI_PLUGIN_VERSION in $CNI_PLUGIN_VERSIONS; do
    CNI_PLUGINS_URL="https://packages.aks.azure.com/cni-plugins/v${CNI_PLUGIN_VERSION}/binaries/cni-plugins-linux-amd64-v${CNI_PLUGIN_VERSION}.tgz"
    downloadCNI
    echo "  - CNI plugin version ${CNI_PLUGIN_VERSION}" >> ${VHD_LOGS_FILEPATH}
done

ACR_CREDENTIAL_PROVIDER_VERSIONS="
1.30.8
1.29.9
"
for ACR_CREDENTIAL_PROVIDER_VERSION in $ACR_CREDENTIAL_PROVIDER_VERSIONS; do
    PROVIDER_ARTIFACT="https://github.com/kubernetes-sigs/cloud-provider-azure/releases/download/v${ACR_CREDENTIAL_PROVIDER_VERSION}/azure-acr-credential-provider-linux-amd64"
    downloadACR
    echo "  - ACR credential provider ${PROVIDER_ARTIFACT}" >> ${VHD_LOGS_FILEPATH}
done

installImg
echo "  - img" >> ${VHD_LOGS_FILEPATH}

systemctl status docker --no-pager || exit 1
echo "Docker images pre-pulled:" >> ${VHD_LOGS_FILEPATH}

METRICS_SERVER_VERSIONS="
0.5.2
"
for METRICS_SERVER_VERSION in ${METRICS_SERVER_VERSIONS}; do
    CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes/metrics-server:v${METRICS_SERVER_VERSION}"
    loadContainerImage ${CONTAINER_IMAGE}
    echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

KUBE_ADDON_MANAGER_VERSIONS="
9.1.5
9.1.6
"
for KUBE_ADDON_MANAGER_VERSION in ${KUBE_ADDON_MANAGER_VERSIONS}; do
    CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes/kube-addon-manager:v${KUBE_ADDON_MANAGER_VERSION}"
    loadContainerImage ${CONTAINER_IMAGE}
    echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

MCR_PAUSE_VERSIONS="
3.8
"
for PAUSE_VERSION in ${MCR_PAUSE_VERSIONS}; do
    # Pull the arch independent MCR pause image which is built for Linux and Windows
    CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes/pause:${PAUSE_VERSION}"
    loadContainerImage "${CONTAINER_IMAGE}"
    echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

CLUSTER_AUTOSCALER_VERSIONS="
1.22.1
"
for CLUSTER_AUTOSCALER_VERSION in ${CLUSTER_AUTOSCALER_VERSIONS}; do
    CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes/autoscaler/cluster-autoscaler:v${CLUSTER_AUTOSCALER_VERSION}"
    loadContainerImage ${CONTAINER_IMAGE}
    echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

CORE_DNS_VERSIONS="
1.9.4
"
for CORE_DNS_VERSION in ${CORE_DNS_VERSIONS}; do
    CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes/coredns:v${CORE_DNS_VERSION}"
    loadContainerImage ${CONTAINER_IMAGE}
    echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

AZURE_CNIIMAGEBASE="mcr.microsoft.com/containernetworking"
AZURE_NPM_VERSIONS="
1.4.59
"
for AZURE_NPM_VERSION in ${AZURE_NPM_VERSIONS}; do
    CONTAINER_IMAGE="${AZURE_CNIIMAGEBASE}/azure-npm:v${AZURE_NPM_VERSION}"
    loadContainerImage ${CONTAINER_IMAGE}
    echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

NVIDIA_DEVICE_PLUGIN_VERSIONS="
1.0.0-beta6
"
for NVIDIA_DEVICE_PLUGIN_VERSION in ${NVIDIA_DEVICE_PLUGIN_VERSIONS}; do
    CONTAINER_IMAGE="mcr.microsoft.com/oss/nvidia/k8s-device-plugin:${NVIDIA_DEVICE_PLUGIN_VERSION}"
    loadContainerImage ${CONTAINER_IMAGE}
    echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

IP_MASQ_AGENT_VERSIONS="
2.8.0
"
for IP_MASQ_AGENT_VERSION in ${IP_MASQ_AGENT_VERSIONS}; do
    CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes/ip-masq-agent:v${IP_MASQ_AGENT_VERSION}"
    loadContainerImage ${CONTAINER_IMAGE}
    echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

KMS_PLUGIN_VERSIONS="0.0.10"
for KMS_PLUGIN_VERSION in ${KMS_PLUGIN_VERSIONS}; do
    CONTAINER_IMAGE="mcr.microsoft.com/k8s/kms/keyvault:v${KMS_PLUGIN_VERSION}"
    loadContainerImage ${CONTAINER_IMAGE}
    echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

BUSYBOX_IMAGE="mcr.microsoft.com/oss/busybox/busybox:1.33.1"
loadContainerImage ${BUSYBOX_IMAGE}
echo "  - ${BUSYBOX_IMAGE}" >> ${VHD_LOGS_FILEPATH}

K8S_VERSIONS="
1.30.10
1.29.15
"
for KUBERNETES_VERSION in ${K8S_VERSIONS}; do
  for component in kube-apiserver kube-controller-manager kube-proxy kube-scheduler; do
    CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes/${component}:v${KUBERNETES_VERSION}"
    loadContainerImage ${CONTAINER_IMAGE}
    echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
  done
  KUBE_BINARY_URL="https://packages.aks.azure.com/kubernetes/v${KUBERNETES_VERSION}/binaries/kubernetes-node-linux-amd64.tar.gz"
  extractKubeBinaries
done

# Starting with 1.16 we pull cloud-controller-manager and cloud-node-manager
CLOUD_CONTROLLER_MANAGER_VERSIONS="
1.30.7
1.29.8
"
for CLOUD_CONTROLLER_MANAGER_VERSION in ${CLOUD_CONTROLLER_MANAGER_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes/azure-cloud-controller-manager:v${CLOUD_CONTROLLER_MANAGER_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

# Starting with 1.16 we pull cloud-controller-manager and cloud-node-manager
CLOUD_NODE_MANAGER_VERSIONS="
1.30.8
1.29.9
"
for CLOUD_NODE_MANAGER_VERSION in ${CLOUD_NODE_MANAGER_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v${CLOUD_NODE_MANAGER_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

AZUREDISK_CSI_VERSIONS="
1.29.1
1.31.5
"
for AZUREDISK_CSI_VERSION in ${AZUREDISK_CSI_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi:v${AZUREDISK_CSI_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

AZUREFILE_CSI_VERSIONS="
1.9.0
"
for AZUREFILE_CSI_VERSION in ${AZUREFILE_CSI_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes-csi/azurefile-csi:v${AZUREFILE_CSI_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

CSI_ATTACHER_VERSIONS="
4.8.0
4.3.0
"
for CSI_ATTACHER_VERSION in ${CSI_ATTACHER_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes-csi/csi-attacher:v${CSI_ATTACHER_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

CSI_NODE_DRIVER_REGISTRAR_VERSIONS="
2.13.0
2.8.0
"
for CSI_NODE_DRIVER_REGISTRAR_VERSION in ${CSI_NODE_DRIVER_REGISTRAR_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes-csi/csi-node-driver-registrar:v${CSI_NODE_DRIVER_REGISTRAR_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

CSI_PROVISIONER_VERSIONS="
5.2.0
3.5.0
"
for CSI_PROVISIONER_VERSION in ${CSI_PROVISIONER_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes-csi/csi-provisioner:v${CSI_PROVISIONER_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

LIVENESSPROBE_VERSIONS="
2.15.0
2.10.0
"
for LIVENESSPROBE_VERSION in ${LIVENESSPROBE_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes-csi/livenessprobe:v${LIVENESSPROBE_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

CSI_RESIZER_VERSIONS="
1.13.1
1.8.0
"
for CSI_RESIZER_VERSION in ${CSI_RESIZER_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes-csi/csi-resizer:v${CSI_RESIZER_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

CSI_SNAPSHOTTER_VERSIONS="
8.2.0
6.2.2
"
for CSI_SNAPSHOTTER_VERSION in ${CSI_SNAPSHOTTER_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes-csi/csi-snapshotter:v${CSI_SNAPSHOTTER_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

SNAPSHOT_CONTROLLER_VERSIONS="
8.2.0
6.2.2
"
for SNAPSHOT_CONTROLLER_VERSION in ${SNAPSHOT_CONTROLLER_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes-csi/snapshot-controller:v${SNAPSHOT_CONTROLLER_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

CSI_SECRETS_STORE_PROVIDER_AZURE_VERSIONS="
0.0.12
"
for CSI_SECRETS_STORE_PROVIDER_AZURE_VERSION in ${CSI_SECRETS_STORE_PROVIDER_AZURE_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/azure/secrets-store/provider-azure:${CSI_SECRETS_STORE_PROVIDER_AZURE_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

CSI_SECRETS_STORE_DRIVER_VERSIONS="
0.0.19
"
for CSI_SECRETS_STORE_DRIVER_VERSION in ${CSI_SECRETS_STORE_DRIVER_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes-csi/secrets-store/driver:v${CSI_SECRETS_STORE_DRIVER_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

AAD_POD_IDENTITY_MIC_VERSIONS="
1.6.1
"
for AAD_POD_IDENTITY_MIC_VERSION in ${AAD_POD_IDENTITY_MIC_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/k8s/aad-pod-identity/mic:${AAD_POD_IDENTITY_MIC_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

AAD_POD_IDENTITY_NMI_VERSIONS="
1.6.1
"
for AAD_POD_IDENTITY_NMI_VERSION in ${AAD_POD_IDENTITY_NMI_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/k8s/aad-pod-identity/nmi:${AAD_POD_IDENTITY_NMI_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

CLUSTER_PROPORTIONAL_AUTOSCALER_VERSIONS="
1.8.5
"
for CLUSTER_PROPORTIONAL_AUTOSCALER_VERSION in ${CLUSTER_PROPORTIONAL_AUTOSCALER_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/oss/kubernetes/autoscaler/cluster-proportional-autoscaler:${CLUSTER_PROPORTIONAL_AUTOSCALER_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

# for testing pipeline
NGINX_VERSIONS="
1.23
"

for NGINX_VERSION in ${NGINX_VERSIONS}; do
  CONTAINER_IMAGE="mcr.microsoft.com/mirror/docker/library/nginx:${NGINX_VERSION}"
  loadContainerImage ${CONTAINER_IMAGE}
  echo "  - ${CONTAINER_IMAGE}" >> ${VHD_LOGS_FILEPATH}
done

df -h

# warn at 75% space taken
[ -s $(df -P | grep '/dev/sda1' | awk '0+$5 >= 75 {print}') ] || echo "WARNING: 75% of /dev/sda1 is used" >> ${VHD_LOGS_FILEPATH}
# error at 95% space taken
[ -s $(df -P | grep '/dev/sda1' | awk '0+$5 >= 95 {print}') ] || exit 1

echo "Using kernel:" >> ${VHD_LOGS_FILEPATH}
tee -a ${VHD_LOGS_FILEPATH} < /proc/version
{ printf "Installed apt packages:\n"; apt list --installed | grep -v 'Listing...'; } >> ${VHD_LOGS_FILEPATH}

if [ "$#" -eq 4 ]; then
  feature_flags="$1"
  build_number="$2"
  build_id="$3"
  commit="$4"

  if [ -z "$feature_flags" ] || [ -z "$build_number" ] || [ -z "$build_id" ] || [ -z "$commit" ]; then
      echo "Arguments cannot be empty"
      exit 1
  fi

  {
    echo "Install completed successfully on " $(date)
    echo "VSTS Build NUMBER: $build_number"
    echo "VSTS Build ID: $build_id"
    echo "Commit: $commit"
    echo "Feature flags: $feature_flags"
  } >> ${VHD_LOGS_FILEPATH}
fi

VHD_CG_MANIFEST=/opt/azure/cgmanifest.json
apt list --installed \
| grep -v 'Listing...' \
| awk '{split($0, a, "/"); print a[1]","$2",Ubuntu,20.04"}' \
| jq --slurp --raw-input \
'{
  "$schema": "https://json.schemastore.org/component-detection-manifest.json",
  "Registrations": (split("\n") | map(split(",") | select(.[0] | length > 0)) | map(
  {
    "Component": {
      "Type": "linux",
      "Linux": {
        "Name": .[0],
        "Version": .[1],
        "Distribution": .[2],
        "Release": .[3]
      }
    }
  }
  ))
}' > ${VHD_CG_MANIFEST}

output_file_contents() {
    local filepath=$1

    if [ -f "$filepath" ]; then
        echo "===== Begin contents of $filepath ====="
        cat "$filepath"
        echo "===== End contents of $filepath ====="
    else
        echo "File not found: $filepath"
    fi
}
output_file_contents "$VHD_LOGS_FILEPATH"
output_file_contents "$VHD_CG_MANIFEST"