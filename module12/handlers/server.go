package handlers

import (
	"context"
	"fmt"
	"github.com/Kit-Hung/cncamp/module12/config"
	"github.com/Kit-Hung/cncamp/module12/log"
	"github.com/Kit-Hung/cncamp/module12/tracing"
	"github.com/Kit-Hung/cncamp/module12/tracing/metrics"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	tracer         trace.TracerProvider
	logger         log.Factory
	port           int
	ctx            context.Context
	srv            *http.Server
	metricsFactory metrics.Factory
}

func NewServer(ctx context.Context, port int, metricsFactory metrics.Factory, logger log.Factory) *Server {
	return &Server{
		ctx:            ctx,
		port:           port,
		tracer:         tracing.InitOTEL(fmt.Sprintf("%s-%s", config.Config.Service.Name, "server"), "otlp", metricsFactory, log.Logger),
		logger:         logger,
		metricsFactory: metricsFactory,
	}
}

func (s *Server) Run() error {
	s.logger.Bg().Info("Starting", zap.String("service", config.Config.Service.Name))
	// 启动 http 服务
	s.srv = &http.Server{
		Addr: ":" + strconv.Itoa(s.port),
		BaseContext: func(_ net.Listener) context.Context {
			return s.ctx
		},
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      s.newHTTPHandler(),
	}

	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.srv.Shutdown(s.ctx)
}
