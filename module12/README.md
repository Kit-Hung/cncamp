# 12 使用 istio 发布服务

## 安装 istio
https://istio.io/latest/docs/setup/getting-started/

### 下载
```shell
 curl -L https://istio.io/downloadIstio | sh -
```

### 安装
```shell
# 此处选用 demo , 要用其他模式参考链接
# https://istio.io/latest/docs/setup/additional-setup/config-profiles/

istioctl install --set profile=demo -y
```

### 给命名空间打 label 自动注入
```shell
kubectl label ns default istio-injection=enabled
```

## 12.1 把 httpserver 服务以 Istio Ingress Gateway 的形式发布出来

1. 创建 gateway
2. 创建 virtualservice


## 如何实现安全保证
1. 创建证书并导入到 secret
2. gateway 使用 tls 并指定证书

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: httpserver-gateway
spec:
  selector:
    istio: ingressgateway
  servers:
    - port:
        number: 443
        name: https-port
        protocol: HTTPS
      hosts:
        - "httpsserver.cncamp.io"
      tls:
        mode: SIMPLE
        credentialName: cncamp-io-tls
```


## 七层路由规则
```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: http-server
spec:
  gateways:
    - httpserver-gateway
  hosts:
    - httpsserver.cncamp.io
  http:
    - match:
        - uri:
            exact: "/service0/hello"
      rewrite:
        uri: "/hello"
      route:
        - destination:
            host: http-server-service0.default.svc.cluster.local
            port:
              number: 80
```

## 考虑 open tracing 的接入
1. 安装 Jaeger
2. 应用增加 header 转发和数据采集

```shell
istioctl upgrade --set meshConfig.defaultConfig.tracing.zipkin.address=jaeger-collector:9411
```