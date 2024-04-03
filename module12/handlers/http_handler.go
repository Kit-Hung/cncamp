package handlers

import (
	"github.com/Kit-Hung/cncamp/module12/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

func (s *Server) newHTTPHandler() http.Handler {
	mux := tracing.NewServeMux(true, s.tracer, s.logger)

	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	// 注册 handler
	handleFunc("/healthz", s.healthz)
	handleFunc("/hello", s.hello)
	handleFunc("/shutdown", s.shutdown)
	mux.Handle("/metrics", promhttp.Handler())

	handler := otelhttp.NewHandler(mux, "/")
	return handler
}
