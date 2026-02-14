package middleware

import (
	"strconv"
	"time"

	"github.com/Kyrbanali/API-GateWay/internal/metrics"
	"github.com/gofiber/fiber/v2"
)

func Metrics() fiber.Handler {
	metrics.Register()

	return func(c *fiber.Ctx) error {
		metrics.HttpInFlight.Inc()
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Seconds()
		status := c.Response().StatusCode()
		method := string(c.Method())
		path := c.Route().Path

		metrics.HttpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, path, strconv.Itoa(status)).Observe(duration)
		metrics.HttpInFlight.Dec()

		return err
	}
}
