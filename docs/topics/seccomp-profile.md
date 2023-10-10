# Seccomp Profile

Seccomp stands for secure computing mode and has been a feature of the Linux kernel since version 2.6.12. It can be used to sandbox the privileges of a process, restricting the calls it is able to make from userspace into the kernel. Kubernetes lets you automatically apply seccomp profiles loaded onto a node to your Pods and containers.

Starting from Kubernetes v1.27, the `--seccomp-default` flag will be automatically enabled for each node in the cluster. See [Kubernetes v1.27 Release Notes](https://kubernetes.io/blog/2023/04/11/kubernetes-v1-27-release/#seccompdefault-graduates-to-stable)

With it enabled, kubelet will use the `RuntimeDefault` seccomp profile by default, which is defined by the container runtime, instead of using the `Unconfined` (seccomp disabled) mode. The default profiles aim to provide a strong set of security defaults while preserving the functionality of the workload.

Use crictl to inspect the new seccomp profile:

```bash
CONTAINER_ID=$(sudo crictl ps -q --name=test-container)
sudo crictl inspect $CONTAINER_ID | jq .info.runtimeSpec.linux.seccomp
```
```json
{
  "defaultAction": "SCMP_ACT_ERRNO",
  "architectures": ["SCMP_ARCH_X86_64", "SCMP_ARCH_X86", "SCMP_ARCH_X32"],
  "syscalls": [
    {
      "names": ["_llseek", "_newselect", "accept", …, "write", "writev"],
      "action": "SCMP_ACT_ALLOW"
    },
    …
  ]
}

```

For more information, please see the following Kubernetes resources:

- [Seccomp Default in Kubernetes 1.27](https://kubernetes.io/blog/2021/08/25/seccomp-default/)
- [Seccomp Security Profiles in Kubernetes](https://kubernetes.io/docs/tutorials/security/seccomp/)

## Override the default seccomp profile
If applications are failing after the upgrade, it's possible that executed syscalls are now being blocked by the default profiles. If that's the case, then you can override the default by explicitly setting the pod or container to run as `Unconfined`. Alternatively, you can create a custom seccomp profile based on the default by adding the additional syscalls to the `"action": "SCMP_ACT_ALLOW"` section.
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
    - name: test-container-nginx
      image: nginx:1.21
      securityContext:
        seccompProfile:
          type: Unconfined
    - name: test-container-redis
      image: redis:6.2
```