#!/bin/bash

{{- if IsCustomCloudProfile}}
  {{- if not IsAzureStackCloud}}
ensureCustomCloudRootCertificates() {
    CUSTOM_CLOUD_ROOT_CERTIFICATES="{{GetCustomCloudRootCertificates}}"
    KUBE_CONTROLLER_MANAGER_FILE=/etc/kubernetes/manifests/kube-controller-manager.yaml

    if [ ! -z $CUSTOM_CLOUD_ROOT_CERTIFICATES ]; then
        # Replace placeholder for ssl binding
        if [ -f $KUBE_CONTROLLER_MANAGER_FILE ]; then
            sed -i "s|<volumessl>|- name: ssl\n      hostPath:\n        path: \\/etc\\/ssl\\/certs|g" $KUBE_CONTROLLER_MANAGER_FILE
            sed -i "s|<volumeMountssl>|- name: ssl\n          mountPath: \\/etc\\/ssl\\/certs\n          readOnly: true|g" $KUBE_CONTROLLER_MANAGER_FILE
        fi

        local i=1
        for cert in $(echo $CUSTOM_CLOUD_ROOT_CERTIFICATES | tr ',' '\n')
        do
            echo $cert | base64 -d > "/usr/local/share/ca-certificates/customCloudRootCertificate$i.crt"
            ((i++))
        done

        update-ca-certificates
    else
        if [ -f $KUBE_CONTROLLER_MANAGER_FILE ]; then
            # remove the placeholder for ssl binding
            sed -i "/<volumessl>/d" $KUBE_CONTROLLER_MANAGER_FILE
            sed -i "/<volumeMountssl>/d" $KUBE_CONTROLLER_MANAGER_FILE
        fi
    fi
}

ensureCustomCloudSourcesList() {
    CUSTOM_CLOUD_SOURCES_LIST="{{GetCustomCloudSourcesList}}"

    if [ ! -z $CUSTOM_CLOUD_SOURCES_LIST ]; then
        # Just in case, let's take a back up before we overwrite
        cp /etc/apt/sources.list /etc/apt/sources.list.backup
        echo $CUSTOM_CLOUD_SOURCES_LIST | base64 -d > /etc/apt/sources.list
    fi
}
  {{end}}

configureK8sCustomCloud() {
  {{- if IsAzureStackCloud}}
  local azure_json_path="/etc/kubernetes/azure.json"
  export -f ensureAzureStackCertificates
  retrycmd 60 10 30 bash -c ensureAzureStackCertificates
  set +x
  # When AUTHENTICATION_METHOD is client_certificate, the certificate is stored into key valut,
  # And SERVICE_PRINCIPAL_CLIENT_SECRET will be the following json payload with based64 encode
  #{
  #    "data": "$pfxAsBase64EncodedString",
  #    "dataType" :"pfx",
  #    "password": "$password"
  #}
  if [[ ${AUTHENTICATION_METHOD,,} == "client_certificate" ]]; then
    SERVICE_PRINCIPAL_CLIENT_SECRET_DECODED=$(echo ${SERVICE_PRINCIPAL_CLIENT_SECRET} | base64 --decode)
    SERVICE_PRINCIPAL_CLIENT_SECRET_CERT=$(echo $SERVICE_PRINCIPAL_CLIENT_SECRET_DECODED | jq .data)
    SERVICE_PRINCIPAL_CLIENT_SECRET_PASSWORD=$(echo $SERVICE_PRINCIPAL_CLIENT_SECRET_DECODED | jq .password)

    # trim the starting and ending "
    SERVICE_PRINCIPAL_CLIENT_SECRET_CERT=${SERVICE_PRINCIPAL_CLIENT_SECRET_CERT#'"'}
    SERVICE_PRINCIPAL_CLIENT_SECRET_CERT=${SERVICE_PRINCIPAL_CLIENT_SECRET_CERT%'"'}

    SERVICE_PRINCIPAL_CLIENT_SECRET_PASSWORD=${SERVICE_PRINCIPAL_CLIENT_SECRET_PASSWORD#'"'}
    SERVICE_PRINCIPAL_CLIENT_SECRET_PASSWORD=${SERVICE_PRINCIPAL_CLIENT_SECRET_PASSWORD%'"'}

    KUBERNETES_FILE_DIR=$(dirname "${azure_json_path}")
    K8S_CLIENT_CERT_PATH="${KUBERNETES_FILE_DIR}/k8s_auth_certificate.pfx"
    echo $SERVICE_PRINCIPAL_CLIENT_SECRET_CERT | base64 --decode >$K8S_CLIENT_CERT_PATH
    # shellcheck disable=SC2002,SC2005
    echo $(cat "${azure_json_path}" |
      jq --arg K8S_CLIENT_CERT_PATH ${K8S_CLIENT_CERT_PATH} '. + {aadClientCertPath:($K8S_CLIENT_CERT_PATH)}' |
      jq --arg SERVICE_PRINCIPAL_CLIENT_SECRET_PASSWORD ${SERVICE_PRINCIPAL_CLIENT_SECRET_PASSWORD} '. + {aadClientCertPassword:($SERVICE_PRINCIPAL_CLIENT_SECRET_PASSWORD)}' |
      jq 'del(.aadClientSecret)') >${azure_json_path}
  fi

  if [[ ${IDENTITY_SYSTEM,,} == "adfs" ]]; then
    # update the tenent id for ADFS environment.
    # shellcheck disable=SC2002,SC2005
    echo $(cat "${azure_json_path}" | jq '.tenantId = "adfs"') >${azure_json_path}
  fi
  set -x

  {{/* Log whether the custom login endpoint is reachable to simplify troubleshooting. */}}
  {{/* CSE will finish successfully but kubelet will error out if not reachable. */}}
  LOGIN_ENDPOINT=$(jq -r .activeDirectoryEndpoint /etc/kubernetes/azurestackcloud.json)
  LOGIN_ENDPOINT=${LOGIN_ENDPOINT#'https://'}
  LOGIN_ENDPOINT=${LOGIN_ENDPOINT%'/'}
  timeout 10 nc -vz ${LOGIN_ENDPOINT} 443 \
  && echo "login endpoint reachable: ${LOGIN_ENDPOINT}" \
  || echo "error: login endpoint not reachable: ${LOGIN_ENDPOINT}"

  {{- if not IsAzureCNI}}
  # Decrease eth0 MTU to mitigate Azure Stack's NRP issue
  echo "iface eth0 inet dhcp" | sudo tee -a /etc/network/interfaces
  echo "    post-up /sbin/ifconfig eth0 mtu 1350" | sudo tee -a /etc/network/interfaces
  ifconfig eth0 mtu 1350
  {{end}}

  {{else}}
  ensureCustomCloudRootCertificates
  ensureCustomCloudSourcesList
  {{end}}
}
{{end}}

{{- if IsAzureStackCloud}}
ensureAzureStackCertificates() {
  ENV_JSON="/etc/kubernetes/azurestackcloud.json"
  ARM_EP=$(jq .resourceManagerEndpoint $ENV_JSON | tr -d '"')
  META_EP="$ARM_EP/metadata/endpoints?api-version=2015-01-01"
  curl $META_EP
  RET=$?
  KCM_FILE=/etc/kubernetes/manifests/kube-controller-manager.yaml
  if [ $RET != 0 ]; then
    if [ -f $KCM_FILE ]; then
      sed -i "s|<volumessl>|- name: ssl\n      hostPath:\n        path: \\/etc\\/ssl\\/certs|g" $KCM_FILE
      sed -i "s|<volumeMountssl>|- name: ssl\n          mountPath: \\/etc\\/ssl\\/certs\n          readOnly: true|g" $KCM_FILE
    fi
    cp /var/lib/waagent/Certificates.pem /usr/local/share/ca-certificates/azsCertificate.crt
    update-ca-certificates
  else
    if [ -f $KCM_FILE ]; then
      sed -i "/<volumessl>/d" $KCM_FILE
      sed -i "/<volumeMountssl>/d" $KCM_FILE
    fi
  fi
  curl $META_EP
  exit $?
}
{{end}}
#EOF
