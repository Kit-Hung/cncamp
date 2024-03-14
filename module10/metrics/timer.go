package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type ExecutionTimer struct {
	histogramVec *prometheus.HistogramVec
	start        time.Time
}

func (t *ExecutionTimer) ObserveTotal() {
	(*t.histogramVec).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}

func NewExecutionTimer(histogramVec *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histogramVec: histogramVec,
		start:        now,
	}
}

func NewTimer() *ExecutionTimer {
	return NewExecutionTimer(functionLatency)
}
