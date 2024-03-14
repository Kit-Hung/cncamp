package util

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	versionKey = "VERSION"
)

func addReqHeaderToResp(w *http.ResponseWriter, r *http.Request) {
	// 将 request 中带的 header 写入 response header
	writer := *w
	if len(r.Header) > 0 {
		for key, value := range r.Header {
			respValue := strings.Join(value, ";")
			writer.Header().Set(key, respValue)
		}

	}
}

func readEnvAndSetToHeader(w *http.ResponseWriter, envKey string) {
	// 读取环境变量, 并设置到 response header
	envValue := os.Getenv(envKey)
	writer := *w
	writer.Header().Set(envKey, envValue)
}

func recordClientIpAndHttpCode(w *http.ResponseWriter, r *http.Request, httpCode int) error {
	// 记录客户端 ip 和 http 返回码
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return err
	}

	if net.ParseIP(host) != nil {
		// 客户端 ip 输出到标准输出
		fmt.Printf("client ip: %v\n", host)
	}

	// http 返回码输出到标准输出
	fmt.Printf("http response code: %v\n", httpCode)
	writer := *w
	writer.WriteHeader(httpCode)
	return nil
}

func RequestHandler(w *http.ResponseWriter, r *http.Request, httpCode int) error {
	/*
		接收客户端 request，并将 request 中带的 header 写入 response header
		读取当前系统的环境变量中的 VERSION 配置，并写入 response header
		Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	*/
	addReqHeaderToResp(w, r)
	readEnvAndSetToHeader(w, versionKey)
	return recordClientIpAndHttpCode(w, r, httpCode)
}
