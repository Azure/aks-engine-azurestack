#!/bin/bash

configureAzureStackInterfaces() {
  NI_FILE="/etc/kubernetes/network_interfaces.json"
  AZ_CNI_CFG_FILE="/etc/kubernetes/interfaces.json"
  ASC="/etc/kubernetes/azurestackcloud.json"
  SM_EP=$(jq -r '.serviceManagementEndpoint' ${ASC})
  AD_EP=$(jq -r '.activeDirectoryEndpoint' ${ASC})
  RM_EP=$(jq -r '.resourceManagerEndpoint' ${ASC})
  TOKEN_URL="${AD_EP}${TENANT_ID}/oauth2/token"

  if [[ ${IDENTITY_SYSTEM,,} == "adfs" ]]; then
    TOKEN_URL="${AD_EP}adfs/oauth2/token"
  fi

  set +x

  TOKEN=$(curl -s --retry 5 --retry-delay 10 --max-time 60 -f -X POST \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=client_credentials" \
    -d "client_id=$SERVICE_PRINCIPAL_CLIENT_ID" \
    --data-urlencode "client_secret=$SERVICE_PRINCIPAL_CLIENT_SECRET" \
    --data-urlencode "resource=$SM_EP" \
    ${TOKEN_URL} | jq '.access_token' | xargs)

  if [[ -z $TOKEN ]]; then
    echo "Error generating token for Azure Resource Manager"
    exit 120
  fi

  curl -s --retry 5 --retry-delay 10 --max-time 60 -f -X GET \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    "${RM_EP}subscriptions/$SUBSCRIPTION_ID/resourceGroups/$RESOURCE_GROUP/providers/Microsoft.Network/networkInterfaces?api-version=$NETWORK_API_VERSION" >${NI_FILE}

  if [[ ! -s ${NI_FILE} ]]; then
    echo "Error fetching network interface configuration for node"
    exit 121
  fi

  echo "Generating Azure CNI interface file"

  mapfile -t local_interfaces < <(cat /sys/class/net/*/address | tr -d : | sed 's/.*/\U&/g')

  SDN_IFS=$(jq ".value | map(select(.properties != null) | select(.properties.macAddress != null) | select(.properties.macAddress | inside(\"${local_interfaces[*]}\"))) | map(select((.properties.ipConfigurations | length) > 0))" ${NI_FILE})

  if [[ -z $SDN_IFS ]]; then
      echo "Error extracting the SDN interfaces from the network interfaces file"
      exit 123
  fi

  AZ_CNI_CFG=$(echo ${SDN_IFS} | jq "{Interfaces: [.[] | {MacAddress: .properties.macAddress, IsPrimary: .properties.primary, IPSubnets: [{Prefix: .properties.ipConfigurations[0].properties.subnet.id, IPAddresses: .properties.ipConfigurations | [.[] | {Address: .properties.privateIPAddress, IsPrimary: .properties.primary}]}]}]}")

  mapfile -t SUBNET_IDS < <(echo ${SDN_IFS} | jq '[.[].properties.ipConfigurations[0].properties.subnet.id] | unique | .[]' -r)

  for SUBNET_ID in "${SUBNET_IDS[@]}"; do
    SUBNET_PREFIX=$(curl -s --retry 5 --retry-delay 10 --max-time 60 -f -X GET \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      "${RM_EP}${SUBNET_ID:1}?api-version=$NETWORK_API_VERSION" |
      jq '.properties.addressPrefix' -r)

    if [[ -z $SUBNET_PREFIX ]]; then
      echo "Error fetching the subnet address prefix for a subnet ID"
      exit 122
    fi

    # shellcheck disable=SC2001
    AZ_CNI_CFG=$(echo ${AZ_CNI_CFG} | sed "s|$SUBNET_ID|$SUBNET_PREFIX|g")
  done

  echo ${AZ_CNI_CFG} >${AZ_CNI_CFG_FILE}

  chmod 0444 ${AZ_CNI_CFG_FILE}

  set -x
}
#EOF
