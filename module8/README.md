# 1. 优雅启动

* postStart
* startupProbe
* livenessProbe
* readinessProbe
* readinessGates



# 2. 优雅终止
* preStop
* terminationGracePeriodSeconds



# 3. 资源需求和 Qos 保证
* Guaranteed
* Burstable
* BestEffort



# 4. 探活
* startupProbe
* livenessProbe
* readinessProbe
* readinessGates



# 5. 日常运维需求，日志等级
* 把日志输出到标准输出，由统一日志收集系统或插件收集
* 日志分级，通过不同的日志级别控制日志的输出量



# 6. 配置和代码分离
* configmap 存常规日志
* secret 存密钥相关
* download api 引用
* env 通过外部配置注入



# 7. 保证整个应用的高可用

## 单个 Pod 的高可用
* 设置合理的 resource.memory limits 防止 oom
* 设置合理的 emptydir.sizeLimit 防止被驱逐 


## 整体高可用
* 冗余部署，多副本
* 跨节点部署
* 跨机架
* 跨可用区部署
* 跨数据中心部署
* 跨提供商部署


## 实现方式
* topologyKeys
* 亲和性、反亲和性



# 8. 通过证书保证 httpServer 安全
* 外部访问入口统一收拢到 ingress
* 在 ingress 中配置证书