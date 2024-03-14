# 3.1 构建本地镜像
```shell
docker build -t kithung/cncamp/http-server:v1 .
```

``

# 3.2 编写 Dockerfile 将模块二作业编写的 httpserver 容器化
* 分阶段构建
* 一个编译、一个运行


# 3.3 将镜像推送至 docker 官方镜像仓库
```shell
docker push kithung/cncamp/http-server:v1
```


# 3.4 通过 docker 命令本地启动 httpserver
```shell
docker run -d kithung/cncamp/http-server:v1
```


# 3.5 通过 nsenter 进入容器查看 IP 配置
```shell
docker inspect -f '{{.State.Pid}}' ${container_id}

nsenter -t ${pid} -n ip addr
```
