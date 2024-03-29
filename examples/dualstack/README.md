# AKS Engine - Dual-stack clusters

> The "dual-stack" feature is not supported on Azure Stack Hub clouds.

## Overview

AKS Engine enables you to create dual-stack (IPv4 *and* IPv6) Kubernetes clusters on Microsoft Azure.

- Dual-stack support is available for Kubernetes version 1.16.0 and later
- Dual-stack services support was added in 1.20.0 (recommended kubernetes version)

> Official docs are available here - https://kubernetes.io/docs/concepts/services-networking/dual-stack/

This example shows you how to configure a dual-stack cluster:

1. **kubernetes.json** - deploying and using [Kubernetes](kubernetes.json).

The default node CIDR mask size for IPv6 is /64 and the default clusters subnet is `fc00::/48`. In 1.17+ node CIDR mask size can be configured by:

```yaml
        "controllerManagerConfig" : {
                "--node-cidr-mask-size-ipv6": <value>
        }
```

> The difference between node CIDR mask size and cluster subnet mask size can't be > 16**

Things to try out after the cluster is deployed -

- Nodes are Kubernetes version 1.20.0 or later

```bash
$ kubectl get nodes
NAME                        STATUS   ROLES    AGE   VERSION
k8s-linuxpool1-20403072-0   Ready    agent    22m   v1.20.0
k8s-linuxpool1-20403072-1   Ready    agent    36m   v1.20.0
k8s-master-20403072-0       Ready    master   37m   v1.20.0
```

- Nodes have 2 internal IPs, one from each ip family

```bash
$ kubectl get nodes k8s-linuxpool1-20403072-0 -o go-template --template='{{range .status.addresses}}{{printf "%s: %s \n" .type .address}}{{end}}'
Hostname: k8s-pool1-12324934-0
InternalIP: 10.240.0.5
InternalIP: 2001:1234:5678:9abc::6
```

- Nodes have 2 PodCIDRs, one from each ip family

```bash
$ kubectl get nodes k8s-linuxpool1-20403072-0 -o go-template --template='{{range .spec.podCIDRs}}{{printf "%s\n" .}}{{end}}'
10.244.2.0/24
fc00::/24
```

- Pods have 2 PodIPs, one from each ip family

```bash
kubectl get pods nginx-pod -o go-template --template='{{range .status.podIPs}}{{printf "%s \n" .ip}}{{end}}'
10.244.2.6
fc00:200::7
```

- Able to reach other pods in cluster using IPv6

```bash
# inside the nginx-pod
# ifconfig eth0
eth0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 10.244.2.6  netmask 255.255.255.0  broadcast 0.0.0.0
        inet6 fc00:200::7  prefixlen 24  scopeid 0x0<global>
        inet6 fe80::8846:8cff:fe35:eaf0  prefixlen 64  scopeid 0x20<link>
        ether 8a:46:8c:35:ea:f0  txqueuelen 0  (Ethernet)
        RX packets 611  bytes 8685170 (8.2 MiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 415  bytes 35685 (34.8 KiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
# ping fd00:100::8
PING fc00:200::7(fc00:200::7) 56 data bytes
64 bytes from fc00:200::7: icmp_seq=1 ttl=62 time=0.798 ms
64 bytes from fc00:200::7: icmp_seq=2 ttl=62 time=0.762 ms
```

- Able to create services with IPv6 using `spec.ipFamilies=[IPv6]` in the service manifest:

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: nginx
  name: nginx-ipv6
spec:
  ipFamilies:
    - IPv6
  ipFamilyPolicy: SingleStack
  ports:
    - port: 80
      protocol: TCP
      targetPort: 80
  selector:
    app: nginx
  type: LoadBalancer
```

```bash
$ kubectl get svc
NAME          TYPE           CLUSTER-IP       EXTERNAL-IP          PORT(S)        AGE
kubernetes    ClusterIP      10.0.0.1         <none>               443/TCP        58m
nginx-ipv6    LoadBalancer   fd00::6283       2603:1030:805:3::3   80:31140/TCP   32s
```

## Limitations

- Dual-stack clusters are supported only with kubenet and azurecni.
  - Dual-stack cluster with azurecni are supported with transparent mode and bridge mode. Bridge mode will be deprecated in future. The recommendation is to use transparent mode.
- Dual-stack clusters are supported on Windows from version 2004 (kernel version 10.0.19041.610) and Kubernetes version 1.19
  - https://kubernetes.io/docs/setup/production-environment/windows/intro-windows-in-kubernetes/#ipv4-ipv6-dual-stack
- Dual-stack clusters are supported with
  - `ipvs` kube-proxy mode (Kubernetes version 1.16+)
  - `iptables` kube-proxy mode (Kubernetes version 1.18+).
- API model enables Azure Standard LB for dual-stack clusters. Azure Basic LBs have a limitation of only 1 IPv6 frontend configurations while Standard LB supports up to 600 IPv6 frontend configurations.
