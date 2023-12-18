1. 构建本地镜像

`docker build -t kithung/cncamp/http-server:v1 .`

2. 编写 Dockerfile 将模块二作业编写的 httpserver 容器化
3. 将镜像推送至 docker 官方镜像仓库

`docker push kithung/cncamp/http-server:v1`

4. 通过 docker 命令本地启动 httpserver

`docker run -d kithung/cncamp/http-server:v1`

5. 通过 nsenter 进入容器查看 IP 配置

`docker inspect -f '{{.State.Pid}}' ec4da85f0931`

`nsenter -t 93774 -n ip addr`