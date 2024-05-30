# Monitoring Kubernetes Clusters

Monitoring your Kubernetes cluster lets you see its health and performance. Statistics such as CPU, memory, and disk usage are available for both Linux and Windows nodes in your AKS Engine cluster.

Resource metrics are collected by the lightweight, in-memory [metrics-server][] component. Metrics-server discovers nodes and queries each one's kubelet for CPU and memory usage.

Metrics-server is part of every AKS Engine deployment. To verify that metrics-server is running in your cluster:

```shell
$ kubectl get pods --namespace=kube-system -l k8s-app=metrics-server
NAME                             READY   STATUS    RESTARTS   AGE
metrics-server-bb7db87bc-nm6vn   1/1     Running   2          140m
```

Tools like `kubectl` and the Kubernetes Dashboard use metrics-server, and it has an [API][metrics-server-api] to get metrics for your custom monitoring solution.

## Monitoring Options

1. [kubectl](#kubectl)
1. [Kubernetes Dashboard](#kubernetes-dashboard)
1. [Monitoring extension](#monitoring-extension)

## Kubectl

The `kubectl top` command is an easy way to see node or pod metrics in your terminal.

Use `kubectl top node` to see the resource consumption of nodes:

```shell
$ kubectl top node
NAME                                 CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%
k8s-agentpool1-26399701-vmss000000   67m          3%     729Mi           10%
k8s-agentpool1-26399701-vmss000001   80m          4%     787Mi           10%
k8s-master-26399701-0                201m         10%    1406Mi          19%
```

Use `kubectl top pod` to see the resource consumption of pods:

```shell
$ kubectl top pod --namespace=kube-system
NAME                                            CPU(cores)   MEMORY(bytes)
azure-cni-networkmonitor-7gfd4                  2m           15Mi
...
kube-proxy-mzlq5                                1m           18Mi
kube-scheduler-k8s-master-26399701-0            3m           16Mi
metrics-server-bb7db87bc-nm6vn                  1m           12Mi
```

## Kubernetes Dashboard

The [Kubernetes Dashboard][kubernetes-dashboard] is a web-based user interface that can visualize cluster metrics.

Describing all of the useful ways to use the dashboard project is out of scope of this documentation. See [here](https://github.com/kubernetes/dashboard) to learn more.

## Monitoring extension

A quick way to scaffold out cloud-native and open source monitoring components is to use the [aks-engine-azurestack monitoring extension](https://github.com/Azure/aks-engine/tree/master/extensions/prometheus-grafana-k8s). For details on how to use the monitoring extension, please refer to the [extension documentation](https://github.com/Azure/aks-engine/tree/master/extensions/prometheus-grafana-k8s). By embedding the extension in your apimodel, the extension will do much of the work to create a monitoring solution in your cluster, which includes the following:

- [cAdvisor](https://github.com/google/cadvisor) daemon set to publish container metrics
- [Prometheus](https://prometheus.io/) for metrics collection and storage
- [Grafana](https://grafana.com/) for dashboard and visualizations

The extension wires up these components together. Post-deployment of the Kubernetes cluster, you just have to retrieve Grafana admin password (Kubernetes secret) and target your browser to the Grafana endpoint. There is already a pre-loaded Kubernetes cluster monitoring dashboard, so out-of-the-box you will have meaningful monitoring points with the extensibility that Prometheus and Grafana offer you.

[creating-a-sample-user]: https://github.com/kubernetes/dashboard/blob/master/docs/user/access-control/creating-sample-user.md
[kubernetes-dashboard]: https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/
[metrics-server]: https://github.com/kubernetes-sigs/metrics-server
[metrics-server-api]: https://github.com/kubernetes/metrics/blob/master/pkg/apis/metrics/v1beta1/types.go
[web-ui-dashboard]: https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/
