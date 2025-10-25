package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// HTTPRequestsTotal - счётчик всех HTTP запросов
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)

	// HTTPRequestDuration - время обработки запросов
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// URLShortensTotal - счётчик созданных ссылок
	URLShortensTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "url_shortens_total",
			Help: "Total number of shortened URLs",
		},
	)

	// URLRedirectsTotal - счётчик редиректов
	URLRedirectsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "url_redirects_total",
			Help: "Total number of URL redirects",
		},
	)

	// ActiveURLsCount - количество активных ссылок
	ActiveURLsCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_urls_count",
			Help: "Current number of active URLs",
		},
	)
)

func Init() {
	prometheus.MustRegister(
		HTTPRequestsTotal,
		HTTPRequestDuration,
		URLShortensTotal,
		URLRedirectsTotal,
		ActiveURLsCount,
	)
}
