package metrics

import (
	"context"
	"github.com/Kit-Hung/cncamp/module12/log"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	MetricNamespace = "http_server"
)

var (
	functionLatency = CreateExecutionTimeMetric(MetricNamespace, "time spent.")
)

func Register(ctx context.Context) {
	err := prometheus.Register(functionLatency)
	if err != nil {
		log.Logger.For(ctx).Error("prometheus register error: ", zap.Error(err))
	}
}

func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "execution_latency_seconds",
		Help:      help,
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
	}, []string{"step"})
}
