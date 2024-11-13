package logging

import (
	"net/http"
	"time"

	"bytes"
	"io"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error
type Middleware func(http.Handler) Handler

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return errors.WithStack(f(w, r))
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func HTTPRequestLoggingMiddleware(logger *zap.Logger) Middleware {
	return func(next http.Handler) Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			start := time.Now()

			// Генерируем уникальный идентификатор для каждого запроса
			reqID := uuid.New().String()
			r.Header.Set("X-Request-ID", reqID)

			// Логируем основные параметры запроса
			logger.Info("Incoming request",
				zap.String("request_id", reqID),
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()))

			// Логируем тело запроса, если оно есть
			if r.Body != nil {
				var buf bytes.Buffer
				tee := io.TeeReader(r.Body, &buf)
				body, err := io.ReadAll(tee)
				if err == nil {
					logger.Info("Request body",
						zap.String("request_id", reqID),
						zap.String("body", string(body)))
				}
				r.Body = io.NopCloser(&buf)
			}

			// Создаем обертку для ResponseWriter для логирования статуса ответа
			lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			err := HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				next.ServeHTTP(lw, r)
				return nil
			}).ServeHTTP(w, r)

			if err != nil {
				logger.Error("Handler error", zap.Error(err))
			}

			// Логируем статус ответа и время выполнения
			logger.Info("Request completed",
				zap.String("request_id", reqID),
				zap.Int("status", lw.statusCode),
				zap.Duration("duration", time.Since(start)))
			return nil
		})
	}
}
