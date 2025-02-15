package logging

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/jakubdrobny/speedcubingslovakia/backend/metrics"
)

func CustomLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	return logger
}

func GinLoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		latency := time.Since(startTime)

		labels := prometheus.Labels{
			"code":   strconv.Itoa(c.Writer.Status()),
			"method": c.Request.Method,
			"url":    c.Request.URL.Path,
		}
		metrics.RequestsTotal.With(labels).
			Inc()
		metrics.RequestDuration.With(labels).Observe(latency.Seconds())

		logger.LogAttrs(context.Background(), slog.LevelInfo, "HTTP request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.String("client_ip", c.ClientIP()),
			slog.Int("user_id", c.GetInt("uid")),
			slog.Duration("latency", latency),
			slog.String("user_agent", c.Request.UserAgent()),
		)
	}
}

func GinRecoveryMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.LogAttrs(context.Background(), slog.LevelError, "Panic recovered",
					slog.Any("error", err),
					slog.String("path", c.Request.URL.Path),
					slog.String("method", c.Request.Method),
					slog.String("client_ip", c.ClientIP()),
				)

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
