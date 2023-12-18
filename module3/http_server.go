package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kit-Hung/cncamp/module3/util"
	"github.com/golang/glog"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

	// 为了测试先设置环境变量
	setEnv()
	// 启动服务
	Start(80)
}

func Start(port int) {
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

func healthz(w http.ResponseWriter, r *http.Request) {
	// 访问 localhost/healthz 时，返回 200
	funcName := "healthz"
	err := util.RequestHandler(&w, r, http.StatusOK)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeToResponse(funcName, &w, err.Error())
		return
	}
	writeToResponse(funcName, &w, "200")
}

func writeToResponse(funcName string, w *http.ResponseWriter, value string) {
	if writeString, err := io.WriteString(*w, value); err != nil {
		glog.Errorf("[%v] write to response error: %v", funcName, err)
	} else {
		glog.Info("[%v] write to response: %v", funcName, writeString)
	}
}

func setEnv() {
	// 设置环境变量
	if err := os.Setenv(versionKey, "kmq test version"); err != nil {
		glog.Error("set env error: %v", err)
	}
}
