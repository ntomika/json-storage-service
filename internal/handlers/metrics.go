package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Объявление метрик
var (
	requestCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
	)

	requestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request latencies.",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// Регистрация метрик
func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
}

// Обработка запроса метрик Prometheus

func (s *Server) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	requestCount.Inc()

	// Запускаем таймер и измеряем длительность запроса
	timer := prometheus.NewTimer(requestDuration)
	defer timer.ObserveDuration()

	// Встроенный обработчик из пакета promhttp
	promhttp.Handler().ServeHTTP(w, r)
}
