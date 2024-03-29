# Using SGX with Kubernetes

> `DCv2-series` virtual machines are supported on Azure Stack Hub clouds.

<!-- TOC -->
- [Using SGX with Kubernetes](#using-sgx-with-kubernetes)
  - [Deploy a Kubernetes Cluster](#deploy-a-kubernetes-cluster)
  - [Deploying the SGX device plugin](#deploying-the-sgx-device-plugin)
    - [Device plugin installation](#device-plugin-installation)
  - [Scheduling Pods to TEE enabled Hardware](#scheduling-pods-to-tee-enabled-hardware)
<!-- /TOC -->

[Intel&reg; Secure Guard Extension](https://software.intel.com/en-us/sgx) (Intel&reg; SGX) is an architecture extension designed to increase the security of application code and data.
Developers may choose [Intel&reg; SGX SDK](https://software.intel.com/en-us/sgx-sdk) or [Open Enclave SDK](https://github.com/Microsoft/openenclave/) to create applications that leverage this technology.

Azure supports provisioning of SGX-enabled VMs under the umbrella of Azure Confidential Compute (ACC). You can create a Kubernetes cluster with one or multiple agent pool(s) running ACC VMs by specifying a [DCv2-series](https://docs.microsoft.com/en-us/azure/virtual-machines/dcv2-series) VM size with a supported distro from the table below.

## Deploy a Kubernetes Cluster
Refer to the [Quickstart Guide](../tutorials/quickstart.md) for details on how to provision a cluster using AKS-Engine. In order to use SGX enabled hardware we suggest updating the cluster model to include an additional agentpool with the supported operating system and virtual machine size. See below for further detail.


| OS           | distro              |
| ------------ | ------------------- |
| Ubuntu 18.04 | `ubuntu-18.04-gen2` |

The following example is a fragment of a cluster definition (apimodel) file declaring two ACC agent pools, one running `Ubuntu 18.04` image on `2 vCPU` nodes, and another running on `4 vCPU` nodes:

```
  "agentPoolProfiles": [
    {
      "name": "agentpool1",
      "count": 3,
      "distro": "ubuntu-18.04-gen2",
      "vmSize": "Standard_DC2s_v2"
    },
    {
      "name": "agentpool2",
      "count": 3,
      "distro": "ubuntu-18.04-gen2",
      "vmSize": "Standard_DC4s_v2"
    }
  ],
```

Note: ACC Gen2 images have Intel SGX DCAP driver [v1.3](https://01.org/intel-softwareguard-extensions/downloads/intel-sgx-dcap-linux-1.3-release) installed

## Deploying the SGX device plugin

You can install the SGX device plugin which surfaces the usage of Intel SGX’s Encrypted Page Cache (EPC) RAM as a schedulable resource for Kubernetes. This allows you to schedule pods that use the [Open Enclave SDK](https://github.com/openenclave/openenclave) onto hardware which supports Trusted Execution Environments.

### Device plugin installation

Using kubectl, deploy device plugin DaemonSet: 
1. For kubernetes versions before v1.17, use: <br>
    kubectl create -f [device-plugin-before-k8s-1-17.yaml](sgx/device-plugin-before-k8s-1-17.yaml)
2. For kubernetes v1.17 and onwards, use: <br> 
    kubectl create -f [device-plugin.yaml](sgx/device-plugin.yaml)

Confirm that the DaemonSet pods are running on each Intel SGX enabled node as follows:

```bash
$ kubectl get pods -n kube-system -l app=sgx-device-plugin
```

```bash
NAME                         READY   STATUS    RESTARTS   AGE
sgx-device-plugin-7d5l8   1/1     Running   0          12m
sgx-device-plugin-jzhk9   1/1     Running   0          12m
```

Confirm that the device plugin is advertising the available EPC RAM to the Kubernetes scheduler by running the following command:

```bash
$ kubectl get nodes <node-name> -o yaml
```

Under the status field you should see the total allocable resources with a name of `kubernetes.azure.com/sgx_epc_mem_in_MiB`
```bash
<snip>
status:
  allocatable:
    kubernetes.azure.com/sgx_epc_mem_in_MiB: "82"
<snip>
```

## Scheduling Pods to TEE enabled Hardware

The following pod specification demonstrates how you would schedule a pod to have access to a TEE by defining a limit on the specific EPC memory that is advertised to the Kubernetes scheduler by the device plugin

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: oe-deployment
spec:
  selector:
    matchLabels:
      app: oe-app
  replicas: 1
  template:
    metadata:
      labels:
        app: oe-app
    spec:
      tolerations:
      - key: kubernetes.azure.com/sgx_epc_mem_in_MiB
        operator: Exists
        effect: NoSchedule
      containers:
      - name: <image_name>
        image: <image_reference>
        command: <exec>
        resources:
          limits:
            kubernetes.azure.com/sgx_epc_mem_in_MiB: 10
```

You can use the following test workload to confirm that your cluster is correctly configured ([Dockerfile](https://github.com/microsoft/openenclave-aks/blob/master/k8s-sgxtest/Dockerfile) for sgx-test):

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: sgx-test
  labels:
    app: sgx-test
spec:
  template:
    metadata:
      labels:
        app: sgx-test
    spec:
      containers:
      - name: sgxtest
        image: oeciteam/sgx-test:1.0
        resources:
          limits:
            kubernetes.azure.com/sgx_epc_mem_in_MiB: 10
      restartPolicy: Never
  backoffLimit: 0
  ```

  You can confirm that the workload successfully created a Trusted Execution Environment by running the following commands:

```bash
$ kubectl get jobs -l app=sgx-test
```

```bash
$ kubectl get jobs -l app=sgx-test
NAME       COMPLETIONS   DURATION   AGE
sgx-test   1/1           1s         23s
```

```bash
$ kubectl get pods -l app=sgx-test
```

```bash
$ kubectl get pods -l app=sgx-test
NAME             READY   STATUS      RESTARTS   AGE
sgx-test-rchvg   0/1     Completed   0          25s
```

```bash
$ kubectl logs -l app=sgx-test
```

```bash
$ kubectl logs -l app=sgx-test
Hello world from the enclave
Enclave called into host to print: Hello World!
```
