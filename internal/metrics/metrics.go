package metrics

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	once sync.Once

	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	httpInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_inflight_requests",
			Help: "Current number of inflight HTTP requests.",
		},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	cacheItems = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cache_items",
			Help: "Number of items in the cache.",
		},
	)
	cacheMemoryBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cache_memory_bytes",
			Help: "Approximate memory used by the cache in bytes.",
		},
	)
)

func register() {
	once.Do(func() {
		mustRegister(prometheus.NewGoCollector())
		mustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
		mustRegister(httpRequestsTotal, httpInFlight, httpRequestDuration, cacheItems, cacheMemoryBytes)
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
	register()
	return adaptor.HTTPHandler(promhttp.Handler())
}

func Middleware() fiber.Handler {
	register()
	return func(c *fiber.Ctx) error {
		httpInFlight.Inc()
		start := time.Now()

		err := c.Next()

		dur := time.Since(start).Seconds()
		status := c.Response().StatusCode()
		method := string(c.Method())
		path := c.Route().Path

		httpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
		httpRequestDuration.WithLabelValues(method, path, strconv.Itoa(status)).Observe(dur)
		httpInFlight.Dec()

		return err
	}
}

func SetCacheStats(items int, bytes int64) {
	register()
	cacheItems.Set(float64(items))
	cacheMemoryBytes.Set(float64(bytes))
}
