# DISA Ubuntu 20.04 LTS STIG

AKS Engine is able to create cluster nodes that comply with the [DISA Ubuntu 20.04 LTS Security Technical Implementation Guide](https://public.cyber.mil/announcement/stig-update-disa-has-released-the-canonical-ubuntu-20-04-lts-stig/). A STIG describes how to minimize network-based attacks and prevent system access when the attacker is interfacing with the system, either physically at the machine or over a network. 

By default, AKS Engine-based clusters address most, but not all, the items required by the STIG. The remaining configuration items required by the STIG will be applied if the API Model sets the `enforceUbuntu2004DisaStig` feature flag:

```json
{
  "properties": {
    "masterProfile": {
      "auditDEnabled": true
    },
    "agentPoolProfiles": [
      {
          "osType": "Linux",
          "auditDEnabled": true
      }
    ],
    "featureFlags": {
      "enforceUbuntu2004DisaStig": true
    }
  }
}
```

The list of items required by the STIG can be found [here](https://www.stigviewer.com/stig/canonical_ubuntu_20.04_lts/).

The following script shows how to validate if a cluster node remain compliant:

```bash
# the zip file name can be different
STIG_URL="https://dl.dod.cyber.mil/wp-content/uploads/stigs/zip/U_CAN_Ubuntu_20-04_LTS_V1R7_STIG_Ansible.zip"
OUT_FILE="U_CAN_Ubuntu_20-04_LTS_STIG_Ansible.zip"
curl -fsSL $STIG_URL -o $OUT_FILE
# ensure zip is installed
unzip $OUT_FILE || exit 99
unzip ubuntu2004STIG-ansible.zip || exit 99

# adjust ansible files to naming insistencies
echo 'ubuntu2004STIG_stigrule_238213_ClientAliveInterval_Line: ClientAliveInterval 120' > vars.yml
echo 'ubuntu2004STIG_stigrule_238214__etc_issue_net_Dest: /etc/issue-stig.net' >> vars.yml

echo "  vars_files:" >> site.yml
echo "    - vars.yml" >> site.yml

sed -i 's|/01-vendor-Ubuntu$|/01-vendor-ubuntu|g' roles/ubuntu2004STIG/tasks/main.yml
sed -i 's|/stig.rules$|/aks-engine.rules|g' roles/ubuntu2004STIG/tasks/main.yml
sed -i 's|auid!=-1|auid!=4294967295|g' roles/ubuntu2004STIG/defaults/main.yml
sed -i 's|/etc/sysctl.conf|/etc/sysctl.d/11-aks-engine.conf|g' roles/ubuntu2004STIG/tasks/main.yml

# dry-run
export XML_PATH
XML_PATH=$(pwd)/results.xml
# ensure ansible is installed
ansible-playbook -v -b -i /dev/null --check site.yml
# check for findings
grep -B 1 fail $XML_PATH | grep -o -E 'V-[[:digit:]]{6}'
```
