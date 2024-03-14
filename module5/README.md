## 5.1 
### 在本地构建一个单节点的基于 HTTPS 的 etcd 集群
```shell
# 安装 etcd 二进制
tar -zxvf etcd-v3.5.12-linux-amd64.tar.gz
cp etcd-v3.5.12-linux-amd64/etcd* /usr/local/bin/

# 生成证书
git clone https://github.com/etcd-io/etcd.git
cd etcd/hack/tls-setup/

# 只保留 127.0.0.1 和 localhost
vim config/req-csr.json

export infra0=127.0.0.1
export infra1=127.0.0.1
export infra2=127.0.0.1
make

# 把生成的证书移动到 /tmp目录
mkdir /tmp/etcd-certs
mv certs/ /tmp/etcd-certs/

# 单节点启动
etcd --listen-client-urls 'https://localhost:12379' \
 --advertise-client-urls 'https://localhost:12379' \
 --listen-peer-urls 'https://localhost:12380' \
 --initial-advertise-peer-urls 'https://localhost:12380' \
 --initial-cluster 'default=https://localhost:12380' \
 --client-cert-auth --trusted-ca-file=/tmp/etcd-certs/certs/ca.pem \
 --cert-file=/tmp/etcd-certs/certs/127.0.0.1.pem \
 --key-file=/tmp/etcd-certs/certs/127.0.0.1-key.pem \
 --peer-client-cert-auth --peer-trusted-ca-file=/tmp/etcd-certs/certs/ca.pem \
 --peer-cert-file=/tmp/etcd-certs/certs/127.0.0.1.pem \
 --peer-key-file=/tmp/etcd-certs/certs/127.0.0.1-key.pem
```

### 写一条数据
```shell
etcdctl --endpoints="https://127.0.0.1:12379" --cert="/tmp/etcd-certs/certs/127.0.0.1.pem" --key="/tmp/etcd-certs/certs/127.0.0.1-key.pem" --cacert="/tmp/etcd-certs/certs/ca.pem" put /a A
```

### 查看数据细节
```shell
etcdctl --endpoints="https://127.0.0.1:12379" --cert="/tmp/etcd-certs/certs/127.0.0.1.pem" --key="/tmp/etcd-certs/certs/127.0.0.1-key.pem" --cacert="/tmp/etcd-certs/certs/ca.pem" get /a -wjson
```

### 删除数据
```shell
etcdctl --endpoints="https://127.0.0.1:12379" --cert="/tmp/etcd-certs/certs/127.0.0.1.pem" --key="/tmp/etcd-certs/certs/127.0.0.1-key.pem" --cacert="/tmp/etcd-certs/certs/ca.pem" del /a
```


## 5.2 在 Kubernetes 集群中创建一个高可用的 etcd 集群
### 通过 helm 安装
```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm pull bitnami/etcd
tar -zxvf etcd-9.14.2.tgz

# 修改相关属性，如 replicaCount 修改为 2
vim etcd/values.yaml

# 安装到 k8s 集群
helm install my-etcd ./etcd

# 运行 client pod
kubectl run my-etcd-client --restart='Never' --image docker.io/bitnami/etcd:3.5.12-debian-12-r7 --env ROOT_PASSWORD=$(kubectl get secret --namespace default my-etcd -o jsonpath="{.data.etcd-root-password}" | base64 -d) --env ETCDCTL_ENDPOINTS="my-etcd.default.svc.cluster.local:2379" --namespace default --command -- sleep infinity

# 通过 client 进行数据操作
kubectl exec --namespace default -it my-etcd-client -- bash
etcdctl --user root:$ROOT_PASSWORD put /message Hello
etcdctl --user root:$ROOT_PASSWORD get /message

# 提供外部访问
 kubectl port-forward --namespace default svc/my-etcd 2379:2379 &
    echo "etcd URL: http://127.0.0.1:2379"
    
# 获取密码
export ETCD_ROOT_PASSWORD=$(kubectl get secret --namespace default my-etcd -o jsonpath="{.data.etcd-root-password}" | base64 -d)
```