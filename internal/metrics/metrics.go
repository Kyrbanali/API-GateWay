package metrics

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	once sync.Once

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

	CacheItems = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cache_items",
			Help: "Number of items in the cache.",
		},
	)
	CacheMemoryBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cache_memory_bytes",
			Help: "Approximate memory used by the cache in bytes.",
		},
	)
)

func Register() {
	once.Do(func() {
		mustRegister(collectors.NewGoCollector())
		mustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		mustRegister(HttpRequestsTotal, HttpInFlight, HttpRequestDuration, CacheItems, CacheMemoryBytes)
	})
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

func SetCacheStats(items int, bytes int64) {
	Register()
	CacheItems.Set(float64(items))
	CacheMemoryBytes.Set(float64(bytes))
}
