package metrics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	HttpInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_inflight_requests",
			Help: "Current number of inflight HTTP requests.",
		},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

type CacheMetrics struct {
	Items prometheus.Gauge
	Bytes prometheus.Gauge
}

func NewCache(cacheName string) CacheMetrics {
	labels := prometheus.Labels{"cache": cacheName}

	items := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "cache_items",
		Help:        "Number of items in the cache.",
		ConstLabels: labels,
	})
	bytes := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "cache_memory_bytes",
		Help:        "Approximate memory used by the cache in bytes.",
		ConstLabels: labels,
	})
	mustRegister(items, bytes)
	return CacheMetrics{Items: items, Bytes: bytes}
}

func Register() {
	mustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		HttpRequestsTotal, HttpInFlight, HttpRequestDuration,
	)
}

func mustRegister(cs ...prometheus.Collector) {
	for _, c := range cs {
		if err := prometheus.DefaultRegisterer.Register(c); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
				continue
			}
			panic(err)
		}
	}
}

func Handler() fiber.Handler {
	Register()
	return adaptor.HTTPHandler(promhttp.Handler())
}
