package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/Kit-Hung/cncamp/module12/config"
	"github.com/Kit-Hung/cncamp/module12/handlers"
	"github.com/Kit-Hung/cncamp/module12/log"
	"github.com/Kit-Hung/cncamp/module12/metrics"
	"github.com/Kit-Hung/cncamp/module12/tracing"
	tracingMetric "github.com/Kit-Hung/cncamp/module12/tracing/metrics"
	"github.com/Kit-Hung/cncamp/module12/tracing/metrics/prometheus"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"time"
)

const (
	versionKey   = "VERSION"
	closeTimeout = 3 * time.Second
)

func main() {
	// 解析命令行参数
	configFilePath := flag.String("config", "/etc/httpServer/config-service0.yaml", "the config file for http server")
	port := flag.Int("port", 80, "the config file for http server")
	flag.Parse()

	// 初始化配置
	config.InitGlobalConfig(*configFilePath)

	// 启动服务
	Start(*port)
}

func Start(port int) {
	// 优雅退出
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	// 为了测试先设置环境变量
	setEnv(ctx)

	// 注册监控指标收集
	metrics.Register(ctx)

	// 设置 open telemetry
	_, otelShutdown, err := tracing.SetupOTelSDK(ctx)
	metricsFactory := prometheus.New().Namespace(tracingMetric.NSOptions{Name: "http_server", Tags: nil})
	if err != nil {
		log.Logger.For(ctx).Fatal("setup otel sdk error: ", zap.Error(err))
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// 启动 http 服务
	srv := handlers.NewServer(ctx, port, metricsFactory, log.Logger)

	// 启动线程监听请求
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.Run()
	}()

	select {
	case err = <-srvErr:
		fmt.Printf("listen and serve error: %v\n", err)
		log.Logger.For(ctx).Fatal("listen and serve error: ", zap.Error(err))
	case <-ctx.Done():
		fmt.Printf("server is shutting down...\n")
		log.Logger.For(ctx).Info("server is shutting down...")

		// 留一点时间给应用收尾
		/*cancelCtx, cancel := context.WithTimeout(ctx, closeTimeout)
		<-cancelCtx.Done()
		defer cancel()*/
		time.Sleep(closeTimeout)
		handlers.ClearResources()
		log.Logger.For(ctx).Info("clear resource finish...")

		stop()
	}

	err = srv.Shutdown()
	log.Logger.For(ctx).Info("server shutdown completed")
	fmt.Println("server shutdown completed")
}

func setEnv(ctx context.Context) {
	// 设置环境变量
	if err := os.Setenv(versionKey, "kmq test version"); err != nil {
		log.Logger.For(ctx).Error("set env error: ", zap.Error(err))
	}
}
