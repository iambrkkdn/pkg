package logging

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func HTTPLoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Прокидываем запрос к следующему обработчику
			next.ServeHTTP(w, r)

			// Логгируем информацию о запросе после завершения
			duration := time.Since(start)
			logger.Info("HTTP Request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Duration("duration", duration),
			)
		})
	}
}
