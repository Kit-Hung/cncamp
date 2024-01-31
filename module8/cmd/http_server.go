package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/Kit-Hung/cncamp/module8/config"
	"github.com/Kit-Hung/cncamp/module8/log"
	"github.com/Kit-Hung/cncamp/module8/util"
	"go.uber.org/zap"
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
	// 解析命令行参数
	configFilePath := flag.String("config", "/etc/httpServer/config.yaml", "the config file for http server")
	flag.Parse()

	// 初始化配置
	config.InitGlobalConfig(*configFilePath)

	// 为了测试先设置环境变量
	setEnv()
	// 启动服务
	Start(80)
}

func Start(port int) {
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/shutdown", shutdown)

	srv := &http.Server{
		Addr: ":" + strconv.Itoa(port),
	}

	// 启动线程监听请求
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("listen and serve error: %v\n", err)
			log.Logger.Panic("listen and serve error: ", zap.Error(err))
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	fmt.Printf("server is shutting down...\n")
	log.Logger.Info("server is shutting down...")

	// 留一点时间给应用收尾
	ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()
	defer clearResources()

	if err := srv.Shutdown(ctx); err != nil {
		log.Logger.Panic("server shutdown error: ", zap.Error(err))
	}
	<-ctx.Done()
	log.Logger.Info("server shutdown completed")
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
		log.Logger.Error("write to response error: ", zap.Any(funcName, err))
	} else {
		log.Logger.Info("write to response: ", zap.Any(funcName, writeString))
	}
}

func setEnv() {
	// 设置环境变量
	if err := os.Setenv(versionKey, "kmq test version"); err != nil {
		log.Logger.Error("set env error: ", zap.Error(err))
	}
}

func shutdown(w http.ResponseWriter, r *http.Request) {
	clearResources()

	funcName := "shutdown"
	err := util.RequestHandler(&w, r, http.StatusOK)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeToResponse(funcName, &w, err.Error())
		return
	}
	writeToResponse(funcName, &w, "ok")
}

func clearResources() {
	err := log.Logger.Sync()
	if err != nil {
		fmt.Printf("logger sync error: %v", err)
	}
}
