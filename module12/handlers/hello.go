package handlers

import (
	"context"
	"fmt"
	"github.com/Kit-Hung/cncamp/module12/config"
	"github.com/Kit-Hung/cncamp/module12/log"
	"github.com/Kit-Hung/cncamp/module12/metrics"
	"github.com/Kit-Hung/cncamp/module12/tracing"
	"github.com/Kit-Hung/cncamp/module12/util"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (s *Server) hello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	funcName := "hello handler"
	log.Logger.For(ctx).Info("entering", zap.String("handler", funcName))
	log.Logger.For(ctx).Info("headers", zap.Any("header", r.Header))

	// 增加执行时间统计
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()

	// 添加 0-2 秒的随机延时
	delay := util.RandInt(0, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))

	// 请求其他服务
	err := s.proxyService(ctx, r, "/hello")
	if err != nil {
		util.WriteToResponseAndHandleError(funcName, &w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取用户名
	user := r.URL.Query().Get("user")
	if user == "" {
		user = "stranger"
	}
	util.WriteToResponseAndHandleError(funcName, &w, r, fmt.Sprintf("hello [%s]", user), http.StatusOK)
	log.Logger.For(ctx).Info("Response： ", zap.String("handler", funcName), zap.Int("delay", delay))
}

func (s *Server) proxyService(ctx context.Context, r *http.Request, subPath string) error {
	log.Logger.For(ctx).Info("start proxy service: ", zap.Any("proxy", config.Config.Proxy), zap.String("subPath", subPath))
	if !config.Config.Proxy.Enabled {
		return nil
	}

	protocol := config.Config.Proxy.Protocol
	url := config.Config.Proxy.Url
	port := config.Config.Proxy.Port
	if protocol == "" || url == "" || port == "" {
		log.Logger.For(ctx).Error("request service failed, protocol or url or port is empty ", zap.String("protocol", protocol),
			zap.String("url", url), zap.String("port", port))
		return nil
	}

	fullUrl := fmt.Sprintf("%s://%s:%s%s", protocol, url, port, subPath)

	tracerName := fmt.Sprintf("%s-%s-%s", config.Config.Service.Name, "hello", "proxyService")
	tracer := tracing.InitOTEL(tracerName, "otlp", s.metricsFactory, log.Logger).Tracer(fmt.Sprintf("%s-%s", config.Config.Service.Name, "hello"))
	ctx, span := tracer.Start(ctx, "start proxy", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		semconv.PeerServiceKey.String("mysql"),
		attribute.
			Key("proxy.service").
			String(fullUrl),
	)
	span.SetAttributes(attribute.Bool("isTrue", true), attribute.String("stringAttr", "hi!"))
	defer span.End()

	resp, err := s.GetRequest(ctx, fullUrl, r)
	if err != nil {
		log.Logger.For(ctx).Error("request service error: ", zap.String("url", fullUrl), zap.Error(err))
		return err
	}

	log.Logger.For(ctx).Info("request service succeed", zap.Any("response", resp.Status))
	return nil
}

func (s *Server) GetRequest(ctx context.Context, url string, r *http.Request) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Logger.For(ctx).Error("new request error: ", zap.Error(err))
		return nil, err
	}

	/*headers := make(http.Header)
	if len(r.Header) > 0 {
		for key, value := range r.Header {
			headers[strings.ToLower(key)] = value
		}
	}*/

	client := http.Client{Transport: otelhttp.NewTransport(
		http.DefaultTransport,
		otelhttp.WithTracerProvider(s.tracer),
	)}
	//req.Header = headers
	resp, err := client.Do(req)
	if err != nil {
		log.Logger.For(ctx).Error("new request error: ", zap.Error(err))
		return nil, err
	}
	return resp, nil
}
