# Patching Cluster Nodes

To favor a consistent behavior, AKS Engine's approach to infrastructure management is immutability.
By default, cluster nodes are not patched on bootstrap nor on a regular basis. 

To prevent node patching, AKS Engine creates file `/etc/apt/apt.conf.d/99periodic` with the following content:

```
APT::Periodic::Update-Package-Lists "0";
APT::Periodic::Download-Upgradeable-Packages "0";
APT::Periodic::AutocleanInterval "0";
APT::Periodic::Unattended-Upgrade "0";
```

## Enable Node Patching

AKS Engine allows changes to the default node patching configuration in situations where applying the latest security patches beats the benefits of immutable infrastructure.

### Enable Unattended Upgrades

Setting `linuxProfile.enableUnattendedUpgrades` to `"true"` configures each Linux node VM (including control plane node VMs) to run `/usr/bin/unattended-upgrade` in the background according to a daily schedule. If enabled, the default `unattended-upgrades` package configuration will be used as provided by the Ubuntu distro version running on the VM. More information [here](https://help.ubuntu.com/community/AutomaticSecurityUpdates).

### Run Unattended Upgrades On Bootstrap Only

An intermediate solution is to set `linuxProfile.runUnattendedUpgradesOnBootstrap` to `"true"`. This will invoke an unattended upgrade when each Linux node VM comes online for the first time.

Setting the `linuxProfile.runUnattendedUpgradesOnBootstrap` option allows cluster administrators to perform cluster upgrades using a [Blue Green approach](https://en.wikipedia.org/wiki/Blue-green_deployment).

## Orchestrate Node Reboots

Some OS patches are only effective after a node reboot. To ensure that the nodes are rebooted in a non-disruptive way, you can deploy the [kured](https://github.com/weaveworks/kured) daemonset.

The `kured` daemonset can be installed by following this [instructions](https://kured.dev/docs/installation/):

```bash
# Find the appropiate version for your cluster https://kured.dev/docs/installation/#kubernetes--os-compatibility
latest=$(curl -s https://api.github.com/repos/kubereboot/kured/releases | jq -r '.[0].tag_name')

# Check https://kured.dev/docs/configuration/ for configuration options
kubectl apply -f "https://github.com/kubereboot/kured/releases/download/$latest/kured-$latest-dockerhub.yaml"
```
