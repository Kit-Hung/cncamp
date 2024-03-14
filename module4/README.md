# 4.1 用 Kubeadm 安装 Kubernetes 集群
## 初始化
```shell
kubeadm init \
 --image-repository registry.aliyuncs.com/google_containers \
 --pod-network-cidr=192.168.0.0/16
```

## 拷贝 kubeconfig

```shell
$ mkdir -p $HOME/.kube
$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
```


## 节点加入集群
```shell
kubeadm join ${ip}:6443 --token ${token} \
	--discovery-token-ca-cert-hash sha256:${sha256}
```


## 安装网络插件
```shell
kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/tigera-operator.yaml
kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/custom-resources.yaml
```


# 4.2
## 通过配置文件生成 configmap
```shell
kubectl create configmap envoy-config --from-file=envoy.yaml
```


## Envoy 的启动配置从外部的配置文件 Mount 进 Pod, 启动一个 Envoy Deployment
```shell
kubectl apply -f envoy-deploy.yaml
```


## 进入 Pod 查看 Envoy 进程和配置
```shell
kubectl exec -it envoy-6b59fd4868-j4tt6 -- bash
cat /etc/envoy/envoy.yaml
```


## 更改配置的监听端口并测试访问入口的变化
```shell
kubectl edit configmap envoy-config
```

* 修改端口后一段时间配置会更新到 pod 里，但是访问还是旧端口
* 重启 pod 生效


## 通过非级联删除的方法逐个删除对象
```shell
kubectl delete deploy envoy
kubectl delete cm envoy-config
```