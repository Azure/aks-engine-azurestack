#!/bin/bash

{{- if IsCustomCloudProfile}}
  {{- if not IsAzureStackCloud}}
ensureCustomCloudRootCertificates() {
    CUSTOM_CLOUD_CERTS="{{GetCustomCloudRootCertificates}}"
    KCM_FILE=/etc/kubernetes/manifests/kube-controller-manager.yaml

    if [ ! -z $CUSTOM_CLOUD_CERTS ]; then
        # Replace placeholder for ssl binding
        if [ -f $KCM_FILE ]; then
            sed -i "s|<volumessl>|- name: ssl\n      hostPath:\n        path: \\/etc\\/ssl\\/certs|g" $KCM_FILE
            sed -i "s|<volumeMountssl>|- name: ssl\n          mountPath: \\/etc\\/ssl\\/certs\n          readOnly: true|g" $KCM_FILE
        fi

        local i=1
        for cert in $(echo $CUSTOM_CLOUD_CERTS | tr ',' '\n')
        do
            echo $cert | base64 -d > "/usr/local/share/ca-certificates/customCloudRootCertificate$i.crt"
            ((i++))
        done

        update-ca-certificates
    else
        if [ -f $KCM_FILE ]; then
            # remove the placeholder for ssl binding
            sed -i "/<volumessl>/d" $KCM_FILE
            sed -i "/<volumeMountssl>/d" $KCM_FILE
        fi
    fi
}

ensureCustomCloudSourcesList() {
    CUSTOM_CLOUD_LIST="{{GetCustomCloudSourcesList}}"

    if [ ! -z $CUSTOM_CLOUD_LIST ]; then
        # Just in case, let's take a back up before we overwrite
        cp /etc/apt/sources.list /etc/apt/sources.list.backup
        echo $CUSTOM_CLOUD_LIST | base64 -d > /etc/apt/sources.list
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
    SPN_DECODED=$(echo ${SERVICE_PRINCIPAL_CLIENT_SECRET} | base64 --decode)
    SPN_CERT=$(echo $SPN_DECODED | jq -r '.data | sub("^\"|\"$"; "")')
    SPN_PWD=$(echo $SPN_DECODED | jq -r '.password | sub("^\"|\"$"; "")')

    K8S_CLIENT_CERT="$(dirname ${azure_json_path})/k8s_auth_certificate.pfx"
    echo $SPN_CERT | base64 --decode >$K8S_CLIENT_CERT
    # shellcheck disable=SC2002,SC2005
    echo $(cat "${azure_json_path}" |
      jq --arg K8S_CLIENT_CERT ${K8S_CLIENT_CERT} '. + {aadClientCertPath:($K8S_CLIENT_CERT)}' |
      jq --arg SPN_PWD ${SPN_PWD} '. + {aadClientCertPassword:($SPN_PWD)}' |
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
  LOGIN_EP=$(jq -r '.activeDirectoryEndpoint | sub("^https://"; "") | sub("/$"; "")' /etc/kubernetes/azurestackcloud.json)
  timeout 10 nc -vz ${LOGIN_EP} 443 \
  && echo "login endpoint reachable: ${LOGIN_EP}" \
  || echo "error: login endpoint not reachable: ${LOGIN_EP}"
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
