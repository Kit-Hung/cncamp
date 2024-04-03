# bookinfo 部署

## 安装 istio

https://istio.io/latest/docs/setup/getting-started/

## 下载及安装
```shell
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.21.0
cp bin/istioctl /usr/local/bin/

# 此处选用 demo , 要用其他模式参考链接
# https://istio.io/latest/docs/setup/additional-setup/config-profiles/

istioctl install --set profile=demo -y
```


## 给命名空间打 label 自动注入
```shell
kubectl label ns default istio-injection=enabled
```


## 部署示例程序
```shell
kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml
```


## 通过 gateway 发布服务
```shell
kubectl apply -f samples/bookinfo/networking/bookinfo-gateway.yaml
```

## 访问服务
```shell
# 获取 gateway url
kubectl get svc istio-ingressgateway -n istio-system

# 访问页面
curl http://192.168.1.100:32000/productpage
```