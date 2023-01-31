#!/bin/bash

{{- if IsAzureStackCloud}}
configureAzureStackInterfaces() {
  NETWORK_INTERFACES_FILE="/etc/kubernetes/network_interfaces.json"
  AZURE_CNI_CONFIG_FILE="/etc/kubernetes/interfaces.json"
  AZURESTACK_ENVIRONMENT_JSON_PATH="/etc/kubernetes/azurestackcloud.json"
  SERVICE_MANAGEMENT_ENDPOINT=$(jq -r '.serviceManagementEndpoint' ${AZURESTACK_ENVIRONMENT_JSON_PATH})
  ACTIVE_DIRECTORY_ENDPOINT=$(jq -r '.activeDirectoryEndpoint' ${AZURESTACK_ENVIRONMENT_JSON_PATH})
  RESOURCE_MANAGER_ENDPOINT=$(jq -r '.resourceManagerEndpoint' ${AZURESTACK_ENVIRONMENT_JSON_PATH})
  TOKEN_URL="${ACTIVE_DIRECTORY_ENDPOINT}${TENANT_ID}/oauth2/token"

  if [[ ${IDENTITY_SYSTEM,,} == "adfs" ]]; then
    TOKEN_URL="${ACTIVE_DIRECTORY_ENDPOINT}adfs/oauth2/token"
  fi

  set +x

  TOKEN=$(curl -s --retry 5 --retry-delay 10 --max-time 60 -f -X POST \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=client_credentials" \
    -d "client_id=$SERVICE_PRINCIPAL_CLIENT_ID" \
    --data-urlencode "client_secret=$SERVICE_PRINCIPAL_CLIENT_SECRET" \
    --data-urlencode "resource=$SERVICE_MANAGEMENT_ENDPOINT" \
    ${TOKEN_URL} | jq '.access_token' | xargs)

  if [[ -z $TOKEN ]]; then
    echo "Error generating token for Azure Resource Manager"
    exit 120
  fi

  curl -s --retry 5 --retry-delay 10 --max-time 60 -f -X GET \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    "${RESOURCE_MANAGER_ENDPOINT}subscriptions/$SUBSCRIPTION_ID/resourceGroups/$RESOURCE_GROUP/providers/Microsoft.Network/networkInterfaces?api-version=$NETWORK_API_VERSION" >${NETWORK_INTERFACES_FILE}

  if [[ ! -s ${NETWORK_INTERFACES_FILE} ]]; then
    echo "Error fetching network interface configuration for node"
    exit 121
  fi

  echo "Generating Azure CNI interface file"

  mapfile -t local_interfaces < <(cat /sys/class/net/*/address | tr -d : | sed 's/.*/\U&/g')

  SDN_INTERFACES=$(jq ".value | map(select(.properties != null) | select(.properties.macAddress != null) | select(.properties.macAddress | inside(\"${local_interfaces[*]}\"))) | map(select((.properties.ipConfigurations | length) > 0))" ${NETWORK_INTERFACES_FILE})

  if [[ -z $SDN_INTERFACES ]]; then
      echo "Error extracting the SDN interfaces from the network interfaces file"
      exit 123
  fi

  AZURE_CNI_CONFIG=$(echo ${SDN_INTERFACES} | jq "{Interfaces: [.[] | {MacAddress: .properties.macAddress, IsPrimary: .properties.primary, IPSubnets: [{Prefix: .properties.ipConfigurations[0].properties.subnet.id, IPAddresses: .properties.ipConfigurations | [.[] | {Address: .properties.privateIPAddress, IsPrimary: .properties.primary}]}]}]}")

  mapfile -t SUBNET_IDS < <(echo ${SDN_INTERFACES} | jq '[.[].properties.ipConfigurations[0].properties.subnet.id] | unique | .[]' -r)

  for SUBNET_ID in "${SUBNET_IDS[@]}"; do
    SUBNET_PREFIX=$(curl -s --retry 5 --retry-delay 10 --max-time 60 -f -X GET \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      "${RESOURCE_MANAGER_ENDPOINT}${SUBNET_ID:1}?api-version=$NETWORK_API_VERSION" |
      jq '.properties.addressPrefix' -r)

    if [[ -z $SUBNET_PREFIX ]]; then
      echo "Error fetching the subnet address prefix for a subnet ID"
      exit 122
    fi

    # shellcheck disable=SC2001
    AZURE_CNI_CONFIG=$(echo ${AZURE_CNI_CONFIG} | sed "s|$SUBNET_ID|$SUBNET_PREFIX|g")
  done

  echo ${AZURE_CNI_CONFIG} >${AZURE_CNI_CONFIG_FILE}

  chmod 0444 ${AZURE_CNI_CONFIG_FILE}

  set -x
}
{{end}}
#EOF
