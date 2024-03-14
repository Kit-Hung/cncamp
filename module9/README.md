# 9.1 测试对 CPU 的校验和准入行为
* 定义一个 Pod，并将该 Pod 中的 nodeName 属性直接写成集群中的节点名    -- 绕过调度器
* 将 Pod 的 CPU 的资源设置为超出计算节点的 CPU 的值
* 创建该 Pod 后 Pod 状态为 OutOfcpu                              -- 当前节点资源不足以运行此 Pod


# 9.2 用 Kubespray 安装集群
## 拉取镜像
```shell
docker pull quay.io/kubespray/kubespray:v2.24.1
```

## 下载代码
```shell
git clone https://github.com/kubernetes-sigs/kubespray.git
cd kubespray
git checkout v2.24.1
```

## 启动 kubespray 容器
```shell
docker run --net host --rm -it --mount type=bind,source="$(pwd)"/inventory/sample,dst=/inventory \
  --mount type=bind,source="${HOME}"/.ssh/id_rsa,dst=/root/.ssh/id_rsa \
  quay.io/kubespray/kubespray:v2.24.1 bash
```

### 免密登录
```shell
ssh-keygen -t rsa
ssh-copy-id -i ~/.ssh/id_rsa.pub root@172.28.58.199
```

### 修改并生成配置
```shell
cp -r inventory/sample inventory/mycluster
declare -a IPS=(172.28.58.199)

# 多 ip 情况
# declare -a IPS=(192.168.10.11 192.168.10.12)
CONFIG_FILE=inventory/mycluster/hosts.yml python3 contrib/inventory_builder/inventory.py ${IPS[@]}
```

### 修改镜像源
```shell
cat > inventory/mycluster/group_vars/k8s_cluster/vars.yml << EOF
gcr_image_repo: "registry.aliyuncs.com/google_containers"
kube_image_repo: "registry.aliyuncs.com/google_containers"
etcd_download_url: "https://mirror.ghproxy.com/https://github.com/coreos/etcd/releases/download/{{ etcd_version }}/etcd-{{ etcd_version }}-linux-{{ image_arch }}.tar.gz"
cni_download_url: "https://mirror.ghproxy.com/https://github.com/containernetworking/plugins/releases/download/{{ cni_version }}/cni-plugins-linux-{{ image_arch }}-{{ cni_version }}.tgz"
calicoctl_download_url: "https://mirror.ghproxy.com/https://github.com/projectcalico/calico/releases/download/{{ calico_ctl_version }}/calicoctl-linux-{{ image_arch }}"
calico_crds_download_url: "https://mirror.ghproxy.com/https://github.com/projectcalico/calico/archive/{{ calico_version }}.tar.gz"
crictl_download_url: "https://mirror.ghproxy.com/https://github.com/kubernetes-sigs/cri-tools/releases/download/{{ crictl_version }}/crictl-{{ crictl_version }}-{{ ansible_system | lower }}-{{ image_arch }}.tar.gz"
runc_download_url: "https://mirror.ghproxy.com/https://github.com/opencontainers/runc/releases/download/{{ runc_version }}/runc.{{ image_arch }}"
nerdctl_download_url: "https://mirror.ghproxy.com/https://github.com/containerd/nerdctl/releases/download/v{{ nerdctl_version }}/nerdctl-{{ nerdctl_version }}-{{ ansible_system | lower }}-{{ image_arch }}.tar.gz"
containerd_download_url: "https://mirror.ghproxy.com/https://github.com/containerd/containerd/releases/download/v{{ containerd_version }}/containerd-{{ containerd_version }}-linux-{{ image_arch }}.tar.gz"
nodelocaldns_image_repo: "registry.lank8s.cn/dns/k8s-dns-node-cache"
dnsautoscaler_image_repo: "registry.lank8s.cn/cpa/cluster-proportional-autoscaler"
EOF
```

### 如果需要指定用户（非必要）
```shell
vim ansible.cfg
add remote_user=cadmin to [default] section
```

### 部署
```shell
ansible-playbook -i inventory/mycluster/hosts.yml cluster.yml -b -vv \
  --private-key=~/.ssh/id_rsa
```


# 9.3 通过 Cluster API 搭建一个集群
## 安装 kind 并创建管理集群
### 获取 kind 二进制文件
```shell
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.22.0/kind-linux-amd64
```

## 通过 kind 安装管理集群，使用 docker 作为基础设施提供者
### 生成配置文件
```shell
cat > kind-cluster-with-extramounts.yaml <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  ipFamily: dual
nodes:
- role: control-plane
  extraMounts:
    - hostPath: /var/run/docker.sock
      containerPath: /var/run/docker.sock
EOF
```

### 创建管理集群
```shell
kind create cluster --config kind-cluster-with-extramounts.yaml
```



## 下载二进制
```shell
curl -L https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.6.2/clusterctl-linux-amd64 -o clusterctl
```

## 初始化 provider
```shell
# Enable the experimental Cluster topology feature.
export CLUSTER_TOPOLOGY=true

# Enable the experimental Machine Pool feature
export EXP_MACHINE_POOL=true

# Initialize the management cluster
clusterctl init --infrastructure docker
```

## 生成集群配置
```shell
# The list of service CIDR, default ["10.128.0.0/12"]
export SERVICE_CIDR=["10.96.0.0/12"]

# The list of pod CIDR, default ["192.168.0.0/16"]
export POD_CIDR=["192.168.0.0/16"]

# The service domain, default "cluster.local"
export SERVICE_DOMAIN="k8s.test"

clusterctl generate cluster capi-quickstart --flavor development \
  --kubernetes-version v1.29.2 \
  --control-plane-machine-count=1 \
  --worker-machine-count=1 \
  > capi-quickstart.yaml
```

## 创建工作集群
```shell
kubectl apply -f capi-quickstart.yaml
```

## 获取工作集群的 kubeconfig
```shell
clusterctl get kubeconfig capi-quickstart > capi-quickstart.kubeconfig
```

## 访问工作集群
```shell
docker ps | grep lb
# 01072213195c   kindest/haproxy:v20230510-486859a6   "haproxy -W -db -f /…"   29 minutes ago   Up 29 minutes   0/tcp, 0.0.0.0:32768->6443/tcp     capi-quickstart-lb

kubectl get nodes --kubeconfig capi-quickstart.kubeconfig --server https://127.0.0.1:32768
# NAME                                     STATUS     ROLES           AGE     VERSION
# capi-quickstart-md-0-blg4k-f82dl-rxl9x   NotReady   <none>          6m3s    v1.29.2
# capi-quickstart-nsxtv-pwdnw              NotReady   control-plane   8m12s   v1.29.2
# capi-quickstart-worker-o3fucl            NotReady   <none>          6m3s    v1.29.2

```

## 安装网络插件
```shell
kubectl --kubeconfig capi-quickstart.kubeconfig --server https://127.0.0.1:32768 apply -f https://github.com/projectcalico/calico/blob/v3.27.2/manifests/calico.yaml

kubectl --kubeconfig capi-quickstart.kubeconfig --server https://127.0.0.1:32768 get nodes
# NAME                                     STATUS   ROLES           AGE   VERSION
# capi-quickstart-md-0-blg4k-f82dl-rxl9x   Ready    <none>          11m   v1.29.2
# capi-quickstart-nsxtv-pwdnw              Ready    control-plane   13m   v1.29.2
# capi-quickstart-worker-o3fucl            Ready    <none>          11m   v1.29.2

```