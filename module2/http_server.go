package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

const (
	versionKey   = "VERSION"
	closeTimeout = 10 * time.Second
)

func main() {
	/*
		接收客户端 request，并将 request 中带的 header 写入 response header
		读取当前系统的环境变量中的 VERSION 配置，并写入 response header
		Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
		当访问 localhost/healthz 时，应返回 200
	*/
	Start(80)
}

func Start(port int) {
	http.HandleFunc("/access", access)
	http.HandleFunc("/healthz", healthz)

	srv := &http.Server{
		Addr: ":" + strconv.Itoa(port),
	}

	// 启动线程监听请求
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("listen and serve error: %v\n", err)
			glog.Fatalf("listen and serve error: %v", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	fmt.Printf("server is shutting down...\n")
	glog.Info("server is shutting down...")

	// 留一点时间给应用收尾
	ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("server shutdown error: %v\n", err)
		glog.Fatalf("server shutdown error: %v", err)
	}
	<-ctx.Done()
	glog.Info("server shutdown completed")
	fmt.Println("server shutdown completed")
}

func access(w http.ResponseWriter, r *http.Request) {
	// 写入 request header
	if len(r.Header) > 0 {
		for key, value := range r.Header {
			respValue := strings.Join(value, ";")
			w.Header().Set(key, respValue)
		}

	}

	// 设置环境变量
	if err := os.Setenv(versionKey, "kmq test version"); err != nil {
		glog.Error("set env error: %v", err)
		writeToResponse("rootHandler", &w, "set env error")
		return
	}
	// 读取环境变量, 并设置到 response header
	version := os.Getenv(versionKey)
	w.Header().Set(versionKey, version)

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
	writeToResponse("rootHandler", &w, "request success!")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	// 访问 localhost/healthz 时，返回 200
	writeToResponse("healthz", &w, "200")
}

func writeToResponse(funcName string, w *http.ResponseWriter, value string) {
	if writeString, err := io.WriteString(*w, value); err != nil {
		glog.Errorf("[%v] write to response error: %v", funcName, err)
	} else {
		glog.Info("[%v] write to response: %v", funcName, writeString)
	}
}
