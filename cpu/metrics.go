package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"
)

const metricsPort = 8386

var metricsServer *http.Server

// StartMetricsServer 开启HttpServer把metrics信息提供给prometheus
func StartMetricsServer() {
	klog.Infof("Starting Metrics Server on [:%d]", metricsPort)
	metricsServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", metricsPort),
		Handler: promhttp.Handler(),
	}
	go func() {
		klog.Warning(metricsServer.ListenAndServe())
	}()
}

// StopMetricsServer 关闭
func StopMetricsServer() {
	if metricsServer != nil {
		klog.Info("Stopping Metrics Server")
		metricsServer.Close()
		metricsServer = nil
	}
}

// kubepods-burstable-pod503aa307_2ead_4099_be3a_6e824c92ab09.slice
type definedMetrics struct {
	// 请求数
	throttledTime *prometheus.GaugeVec
	burstTime     *prometheus.GaugeVec
}

var metrics = definedMetrics{
	throttledTime: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "throttled_time",
			Help: "The total number of requests throttled by the throttler.",
		},
		[]string{"pod", "container"},
	),
	burstTime: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "burst_time",
			Help: "The total number of requests throttled by the throttler.",
		},
		[]string{"pod", "container"},
	),
}

func init() {
	prometheus.MustRegister(metrics.throttledTime, metrics.burstTime)
}

type metricKey struct {
	pod       string
	container string
}

type metricValue struct {
	throttledTime uint64
	burstTime     uint64
}

func RecordMetrics(k metricKey, v metricValue) {
	metrics.throttledTime.With(prometheus.Labels{
		"pod":       k.pod,
		"container": k.container,
	}).Set(float64(v.throttledTime))
	metrics.burstTime.With(prometheus.Labels{
		"pod":       k.pod,
		"container": k.container,
	}).Set(float64(v.burstTime))
}
