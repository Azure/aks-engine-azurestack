# Cluster Certificate Rotation Feature Design

Design document for the certificate rotation feature in AKS Engine clusters.

## Problem

Until PR [#396](https://github.com/Azure/aks-engine-azurestack/pull/396), the default expiration of the certs generated by AKS Engine was 2 years. Clusters that were deployed with an AKS Engine version prior to that PR will soon reach expiration. We currently do not have an outlined process for rotating certificates in an existing cluster.

## Goals

Provide a process that is easy to follow for rotating CA, etcd, kubelet, kubeconfig and apiserver certificates in a cluster built with AKS Engine.

## Non-Goals

- Providing a certificate rotation tool for AKS clusters.
- Providing a tool for etcd backup.
- Rotating certificates in an existing cluster with no downtime.
- Rotating proxy certs.
- Providing a tool for recovery of unhealthy clusters.

## Alternative implementations considered

### Using upgrade to update the certificates

What: Implementing a new "certificate-rotation" command that runs upgrade under the hood. This command would perform a same k8s version upgrade on the cluster with newly generated certificates replacing the old ones in the apimodel.

Pros:

- No reboot required after the fact.

Cons:

- Requires a re-provisioning of all the VMs, causing an unnecessary time overhead.
- Current implementation does not support this: the newly built master nodes would not be able to join the etcd cluster, thus failing CSE validation.

### Bash scripts

What: Provide users with a bash script that they can run to rotate certificates on their cluster.

Pros:

- Easier to implement.
- More lightweight.
- Performs in-place cert update.

Cons:

- Requires more manual steps from the user. For example, they would have to set environemnent variables for their SSH key, KUBECONFIG, _output directory, etc.
- Not Windows friendly: would require running the script in a container, unless we provide a Powershell equivalent.

### New aks-engine-azurestack binary command

What: New `aks-engine-azurestack rotate-certs` command that uses a combination of Kubernetes client-go and SSH commands to access the cluster and rotate certificates.

Pros:

- Better UX.
- Faster than running a full upgrade.
- Easier to unit test.
- Works across platforms.
- Output directory can be one of the flag inputs (similar to `upgrade` and `scale`).
- Lets us re-use the certificate generation logic from deployment code.

Cons:

- Potentially more work than a bash script.

### Sources

https://kubernetes.io/docs/tasks/tls/certificate-rotation/

https://github.com/coreos/tls_rotate
