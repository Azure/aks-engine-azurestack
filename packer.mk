build-packer:
	@packer build -var-file=vhd/packer/settings.json vhd/packer/vhd-image-builder.json

build-packer-windows:
	@packer build -var-file=vhd/packer/settings.json -var-file=vhd/packer/windows-${WINDOWS_SERVER_VERSION}-vars.json vhd/packer/windows-vhd-builder.json

init-packer:
	@./vhd/packer/init-variables.sh

az-login:
	az login --service-principal -u ${CLIENT_ID} -p ${CLIENT_SECRET} --tenant ${TENANT_ID}

run-packer: az-login
	@packer version && set -o pipefail && ($(MAKE) init-packer | tee packer-output) && ($(MAKE) build-packer | tee -a packer-output)

run-packer-windows: az-login
	@packer version && set -o pipefail && ($(MAKE) init-packer | tee packer-output) && ($(MAKE) build-packer-windows | tee -a packer-output)

az-copy: az-login
	azcopy-preview copy "${OS_DISK_SAS}" "${SA_CONTAINER_URL}?${SA_TOKEN}" --overwrite=false

new-az-copy: 
	azcopy-preview copy "${OS_DISK_SAS}" "${SA_CONTAINER_URL}" --overwrite=false

delete-sa: az-login
	az storage account delete -n ${PACKER_TEMP_SA} -g ${PACKER_TEMP_GROUP} --yes
generate-sas: az-login
	az storage container generate-sas --name ubuntu --permissions lr --connection-string "${CLASSIC_SA_CONNECTION_STRING}" --start ${START_DATE} --expiry ${EXPIRY_DATE} | tr -d '"' | tee -a vhd-sas && cat vhd-sas

windows-vhd-publishing-info: az-login
	@./vhd/packer/generate-windows-vhd-publishing-info.sh

sig-image-version: az-login
	az sig image-version create \
		--resource-group ${SIG_GROUP} \
		--gallery-name ${SIG_NAME} \
		--gallery-image-definition ${SIG_IMG_DEF} \
		--gallery-image-version ${VHD_VERSION} \
		--target-regions ${SIG_LOCATION} \
		--replica-count 1 \
		--os-vhd-uri ${SA_CONTAINER_URL}/${VHD_NAME} \
		--os-vhd-storage-account ${VHD_SA}
