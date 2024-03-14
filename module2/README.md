# 2.1 接收客户端 request，并将 request 中带的 header 写入 response header
```go
for key, value := range r.Header {
    respValue := strings.Join(value, ";")
    w.Header().Set(key, respValue)
}
```


# 2.2 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
```go
version := os.Getenv(versionKey)
w.Header().Set(versionKey, version)
```


# 2.3 Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
```go
// 记录客户端 ip 和 http 返回码
host, _, err := net.SplitHostPort(r.RemoteAddr)
if err != nil {
    glog.Error("split host port error: %v", err)
    writeToResponse("rootHandler", &w, "split host port error")
    return
}

if net.ParseIP(host) != nil {
    // 客户端 ip 输出到标准输出
    fmt.Printf("client ip: %v\n", host)
}

// http 返回码输出到标准输出
fmt.Printf("http response code: %v\n", http.StatusOK)
w.WriteHeader(http.StatusOK)
```


# 2.4 当访问 localhost/healthz 时，应返回 200
```go
func healthz(w http.ResponseWriter, r *http.Request) {
	// 访问 localhost/healthz 时，返回 200
	w.WriteHeader(http.StatusOK)
	writeToResponse("healthz", &w, "200")
}
```