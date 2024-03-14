# 10.1 将 nginx 容器上传到 harbor
## Harbor 安装
### 下载 helm chart
```shell
helm repo add harbor https://helm.goharbor.io
helm fetch harbor/harbor --untar
```

### 修改配置
```shell
vim harbor/values.yaml

# 测试修改的配置
expose:
  type: nodePort
tls:
  commonName: 'core.harbor.domain'

persistence: false
```

### 安装 Harbor
```shell
helm install harbor ./harbor -n harbor
```

### 获取证书
```shell
# 页面
https://10.105.111.219/harbor/configs/setting

# 接口
https://10.105.111.219/api/v2.0/systeminfo/getcert
```

### 配置 docker
```shell
mkdir -p /etc/docker/certs.d/core.harbor.domain
cp ca.crt /etc/docker/certs.d/core.harbor.domain
systemctl restart docker

docker login -u admin -p Harbor12345 core.harbor.domain
```

### 配置 containerd
```shell
cp ca.crt /etc/containerd/certs.d/core.harbor.domain/
mkdir -p /etc/containerd/certs.d/core.harbor.domain

vim /etc/containerd/certs.d/core.harbor.domain/hosts.toml
server = "https://core.harbor.domain"

[host."http://core.harbor.domain"]
  capabilities = ["pull", "resolve", "push"]
  skip_verify = true
  ca = "ca.crt"

systemctl restart containerd
```

### 修改 host
```shell
# harbor svc 的 cluster ip
10.105.111.219 core.harbor.domain
```

### 查看 repositories 和 blobs
```shell
kubectl -n harbor exec -it harbor-registry-7886456f94-vkfv5 -- bash

ls -la /storage/docker/registry/v2/repositories/
ls -la /storage/docker/registry/v2/blobs/
```

## 查看数据库数据
```shell
kubectl exec -it harbor-database-0 -- bash

psql -U postgres -d postgres -h 127.0.0.1 -p 5432
\c registry
select * from harbor_user;
```


# 10.2
## 10.2.1 为 HTTPServer 添加 0-2 秒的随机延时
```go
delay := util.RandInt(0, 2000)
time.Sleep(time.Millisecond * time.Duration(delay))
```

## 10.2.2 为 HTTPServer 项目添加延时 Metric
### 注册指标
```go
prometheus.Register(functionLatency)
```

### 生成监控数据
```go
func (t *ExecutionTimer) ObserveTotal() {
	(*t.histogramVec).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}
```

### 提供接口访问
```go
http.Handle("/metrics", promhttp.Handler())
```

## 10.2.3 将 HTTPServer 部署至测试集群，并完成 Prometheus 配置
### 将 HTTPServer 部署至测试集群
#### 重新打包镜像并推送到仓库
```shell
make build-image
```

#### 修改 httpserver deployment 的镜像
```shell
- image: core.harbor.domain/http-server/http_server:v10
```

### 安装 loki-stack
#### 添加 helm repo
```shell
helm repo add grafana https://grafana.github.io/helm-charts
```

#### 更新
```shell
helm repo update
```

#### 安装
```shell
helm upgrade --install loki grafana/loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistence.enabled=false,prometheus.server.persistentVolume.enabled=false
```


## 10.2.4 从 Promethus 界面中查询延时指标数据
登录 prometheus 界面查询 http_server_execution_latency_seconds_bucket


## 10.2.5 创建一个 Grafana Dashboard 展现延时分配情况
### 获取 grafana 登录密码
```shell
kubectl get secret loki-grafana -o yaml
```

### 登录 grafana 界面并创建 dashboard
import resources/grafana-dashboard/httpserver-latency.json

